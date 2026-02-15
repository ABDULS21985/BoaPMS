package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/auth"
	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// authService orchestrates authentication using AD or local credentials,
// JWT token issuance, and dynamic role assignment.
// This mirrors the .NET AuthController's login flow.
type authService struct {
	db      *gorm.DB
	users   *UserManagementService
	jwt     *JWTService
	ad      ActiveDirectoryService
	gs      GlobalSettingService
	erpSQL  *repository.Container
	cfg     *config.Config
	log     zerolog.Logger
}

func newAuthService(repos *repository.Container, cfg *config.Config, log zerolog.Logger) AuthService {
	users := NewUserManagementService(repos, log)
	jwtSvc := NewJWTService(cfg.JWT, log)
	adSvc := newActiveDirectoryService(cfg.ActiveDirectory, log)
	gsSvc := newGlobalSettingService(repos, log)

	return &authService{
		db:     repos.GormDB,
		users:  users,
		jwt:    jwtSvc,
		ad:     adSvc,
		gs:     gsSvc,
		erpSQL: repos,
		cfg:    cfg,
		log:    log,
	}
}

// AuthenticateAD performs dual-mode authentication:
//  1. If AD is enabled (via GlobalSetting), authenticate against Active Directory first.
//  2. If AD auth fails or is disabled, fall back to local DB password verification.
//
// On success it performs dynamic role assignment and returns an AuthenticateResponse
// with JWT access and refresh tokens.
func (s *authService) AuthenticateAD(ctx context.Context, username, password string) (interface{}, error) {
	// Check if AD authentication is enabled
	adEnabled, _ := s.gs.GetBoolValue(ctx, auth.SettingEnableADAuth)

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("looking up user: %w", err)
	}

	// Check lockout
	if user != nil && s.users.IsLockedOut(user) {
		return nil, fmt.Errorf("account is locked out")
	}

	authenticated := false

	// 1) Try AD authentication if enabled
	if adEnabled {
		adOk, adErr := s.ad.Authenticate(username, password)
		if adErr != nil {
			s.log.Warn().Err(adErr).Str("user", username).Msg("AD authentication error, falling back to local")
		}
		if adOk {
			authenticated = true

			// If user doesn't exist locally, auto-provision from AD
			if user == nil {
				user, err = s.provisionADUser(ctx, username)
				if err != nil {
					return nil, fmt.Errorf("provisioning AD user: %w", err)
				}
			}
		}
	}

	// 2) Fall back to local DB authentication
	if !authenticated {
		if user == nil {
			return nil, fmt.Errorf("invalid username or password")
		}
		if !user.IsActive {
			return nil, fmt.Errorf("user account is deactivated")
		}
		if !s.users.VerifyPassword(user, password) {
			maxAttempts, _ := s.gs.GetIntValue(ctx, auth.SettingMaxFailedAttempts)
			lockoutMin, _ := s.gs.GetIntValue(ctx, auth.SettingLockoutDuration)
			if maxAttempts <= 0 {
				maxAttempts = 5
			}
			if lockoutMin <= 0 {
				lockoutMin = 15
			}
			_ = s.users.IncrementAccessFailedCount(ctx, user, maxAttempts, lockoutMin)
			return nil, fmt.Errorf("invalid username or password")
		}
		authenticated = true
	}

	if !authenticated || user == nil {
		return nil, fmt.Errorf("authentication failed")
	}

	// Reset failed attempts on success
	_ = s.users.ResetAccessFailedCount(ctx, user)

	// Dynamic role assignment (mirrors .NET AuthController role logic)
	roles, err := s.resolveRoles(ctx, user)
	if err != nil {
		s.log.Warn().Err(err).Str("user", username).Msg("Failed to resolve dynamic roles")
		roles = []string{auth.RoleStaff}
	}

	// Resolve permissions from role-permission junction table
	permissions, permErr := s.users.GetPermissionsByRoles(ctx, roles)
	if permErr != nil {
		s.log.Warn().Err(permErr).Str("user", username).Msg("Failed to resolve permissions")
	}

	// Resolve organizational unit from ERP data
	orgUnit := ""
	if s.erpSQL.ErpSQL != nil {
		orgUnit = s.users.GetOrganizationalUnit(ctx, s.erpSQL.ErpSQL, user.ID)
	}

	// Get token expiry from settings (with fallback to config)
	expiryMinutes, _ := s.gs.GetIntValue(ctx, auth.SettingTokenExpiryMinutes)

	// Generate tokens with all claims: UserId, Email, FullName, Roles, Permissions, OrganizationalUnit
	claims := TokenClaims{
		UserID:             user.ID,
		Email:              user.Email,
		Name:               user.FullName(),
		Roles:              roles,
		Permissions:        permissions,
		OrganizationalUnit: orgUnit,
	}
	accessToken, expiresAt, err := s.jwt.GenerateAccessToken(claims, expiryMinutes)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	// Persist the hashed refresh token in the database for validation and revocation.
	if err := s.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("storing refresh token: %w", err)
	}

	resp := &auth.AuthenticateResponse{
		UserID:             user.ID,
		Username:           user.UserName,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Email:              user.Email,
		Roles:              roles,
		Permissions:        permissions,
		OrganizationalUnit: orgUnit,
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		ExpiresAt:          expiresAt,
	}

	s.log.Info().Str("user", username).Strs("roles", roles).Msg("User authenticated")
	return resp, nil
}

// GenerateTokenPair creates an access + refresh token pair for the given user.
func (s *authService) GenerateTokenPair(ctx context.Context, userID string, roles []string) (string, string, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return "", "", fmt.Errorf("user not found")
	}

	permissions, _ := s.users.GetPermissionsByRoles(ctx, roles)

	claims := TokenClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Name:        user.FullName(),
		Roles:       roles,
		Permissions: permissions,
	}

	expiryMinutes, _ := s.gs.GetIntValue(ctx, auth.SettingTokenExpiryMinutes)
	accessToken, _, err := s.jwt.GenerateAccessToken(claims, expiryMinutes)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Persist the hashed refresh token in the database.
	if err := s.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return "", "", fmt.Errorf("storing refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT and returns the embedded claims.
func (s *authService) ValidateToken(ctx context.Context, token string) (interface{}, error) {
	return s.jwt.ValidateToken(token)
}

// RefreshAccessToken validates a refresh token against the database, issues a
// new JWT access token, and rotates the refresh token (revoke old, create new).
func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*auth.TokenResponse, error) {
	tokenHash := hashToken(refreshToken)

	// Look up the refresh token by its SHA-256 hash.
	var stored auth.RefreshToken
	err := s.db.WithContext(ctx).
		Where("token = ?", tokenHash).
		First(&stored).Error
	if err != nil {
		s.log.Warn().Str("token_hash", tokenHash[:8]).Msg("Refresh token not found in database")
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Validate the token is still usable.
	if stored.Revoked {
		s.log.Warn().Str("user_id", stored.UserID).Msg("Attempt to use revoked refresh token")
		return nil, fmt.Errorf("refresh token has been revoked")
	}
	if stored.IsExpired() {
		s.log.Warn().Str("user_id", stored.UserID).Msg("Attempt to use expired refresh token")
		return nil, fmt.Errorf("refresh token has expired")
	}

	// Load the user associated with the token.
	user, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil {
		return nil, fmt.Errorf("looking up user for refresh: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found for refresh token")
	}
	if !user.IsActive {
		return nil, fmt.Errorf("user account is deactivated")
	}

	// Resolve roles and permissions for the new access token.
	roles, err := s.resolveRoles(ctx, user)
	if err != nil {
		s.log.Warn().Err(err).Str("user", user.UserName).Msg("Failed to resolve roles during refresh")
		roles = []string{auth.RoleStaff}
	}

	permissions, _ := s.users.GetPermissionsByRoles(ctx, roles)

	orgUnit := ""
	if s.erpSQL.ErpSQL != nil {
		orgUnit = s.users.GetOrganizationalUnit(ctx, s.erpSQL.ErpSQL, user.ID)
	}

	// Generate a new access token.
	expiryMinutes, _ := s.gs.GetIntValue(ctx, auth.SettingTokenExpiryMinutes)
	claims := TokenClaims{
		UserID:             user.ID,
		Email:              user.Email,
		Name:               user.FullName(),
		Roles:              roles,
		Permissions:        permissions,
		OrganizationalUnit: orgUnit,
	}
	newAccessToken, expiresAt, err := s.jwt.GenerateAccessToken(claims, expiryMinutes)
	if err != nil {
		return nil, fmt.Errorf("generating new access token: %w", err)
	}

	// Rotate: revoke the old refresh token and issue a new one.
	if err := s.db.WithContext(ctx).Model(&stored).Update("revoked", true).Error; err != nil {
		s.log.Error().Err(err).Str("token_id", stored.ID).Msg("Failed to revoke old refresh token")
		return nil, fmt.Errorf("revoking old refresh token: %w", err)
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generating new refresh token: %w", err)
	}

	if err := s.storeRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("storing rotated refresh token: %w", err)
	}

	s.log.Info().Str("user", user.UserName).Msg("Access token refreshed with token rotation")

	return &auth.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// hashToken produces a hex-encoded SHA-256 digest of the plaintext token.
func hashToken(plaintext string) string {
	h := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(h[:])
}

// storeRefreshToken hashes the plaintext refresh token and persists it
// in the database. The expiry is derived from the JWT config.
func (s *authService) storeRefreshToken(ctx context.Context, userID, plaintext string) error {
	expiry := s.cfg.JWT.RefreshTokenExpiry
	if expiry <= 0 {
		expiry = 7 * 24 * time.Hour // default 7 days
	}

	rt := auth.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     hashToken(plaintext),
		ExpiresAt: time.Now().UTC().Add(expiry),
		CreatedAt: time.Now().UTC(),
		Revoked:   false,
	}

	if err := s.db.WithContext(ctx).Create(&rt).Error; err != nil {
		return fmt.Errorf("inserting refresh token: %w", err)
	}
	return nil
}

// provisionADUser creates a local user record from AD data.
func (s *authService) provisionADUser(ctx context.Context, username string) (*identity.ApplicationUser, error) {
	adUserRaw, err := s.ad.GetUser(username)
	if err != nil {
		return nil, err
	}

	adUser, ok := adUserRaw.(*auth.ADUser)
	if !ok || adUser == nil {
		// AD lookup didn't return user details; create with minimal info
		adUser = &auth.ADUser{
			Username: username,
			Email:    username + "@" + s.cfg.ActiveDirectory.Domain,
		}
	}

	user := &identity.ApplicationUser{
		ID:       username, // Use sAMAccountName as the ID for AD users
		UserName: username,
		Email:    adUser.Email,
		FirstName: adUser.FirstName,
		LastName:  adUser.LastName,
		IsActive:  true,
	}

	if err := s.users.CreateUser(ctx, user, ""); err != nil {
		return nil, err
	}
	return user, nil
}

// resolveRoles determines the set of roles for a user.
// This mirrors the .NET AuthController's dynamic role assignment logic
// that queries the ERP organisational hierarchy.
func (s *authService) resolveRoles(ctx context.Context, user *identity.ApplicationUser) ([]string, error) {
	// Start with persisted roles from the database
	roles, err := s.users.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// If no roles assigned, default to Staff
	if len(roles) == 0 {
		roles = append(roles, auth.RoleStaff)
	}

	// Dynamic role resolution from ERP data (if ERP DB is available)
	if s.erpSQL.ErpSQL != nil {
		dynamicRoles := s.resolveOrgHierarchyRoles(ctx, user.ID)
		roles = mergeRoles(roles, dynamicRoles)
	}

	return roles, nil
}

// resolveOrgHierarchyRoles queries the ERP database to determine if the user
// holds a leadership position (HeadOfOffice, HeadOfDivision, HeadOfDepartment, Supervisor).
// Mirrors the .NET AuthController's ERP-based role assignment.
func (s *authService) resolveOrgHierarchyRoles(ctx context.Context, staffID string) []string {
	var dynamicRoles []string

	// Check if user is head of an office
	var officeCount int
	err := s.erpSQL.ErpSQL.GetContext(ctx, &officeCount,
		`SELECT COUNT(*) FROM offices WHERE head_staff_id = $1`, staffID)
	if err == nil && officeCount > 0 {
		dynamicRoles = append(dynamicRoles, auth.RoleHeadOfOffice)
	}

	// Check if user is head of a division
	var divCount int
	err = s.erpSQL.ErpSQL.GetContext(ctx, &divCount,
		`SELECT COUNT(*) FROM divisions WHERE head_staff_id = $1`, staffID)
	if err == nil && divCount > 0 {
		dynamicRoles = append(dynamicRoles, auth.RoleHeadOfDivision)
	}

	// Check if user is head of a department
	var deptCount int
	err = s.erpSQL.ErpSQL.GetContext(ctx, &deptCount,
		`SELECT COUNT(*) FROM departments WHERE head_staff_id = $1`, staffID)
	if err == nil && deptCount > 0 {
		dynamicRoles = append(dynamicRoles, auth.RoleHeadOfDepartment)
	}

	// Check if user is a supervisor
	var supCount int
	err = s.erpSQL.ErpSQL.GetContext(ctx, &supCount,
		`SELECT COUNT(*) FROM employees WHERE supervisor_staff_id = $1`, staffID)
	if err == nil && supCount > 0 {
		dynamicRoles = append(dynamicRoles, auth.RoleSupervisor)
	}

	return dynamicRoles
}

// mergeRoles combines two role slices, deduplicating entries.
func mergeRoles(existing, additional []string) []string {
	set := make(map[string]struct{}, len(existing)+len(additional))
	for _, r := range existing {
		set[r] = struct{}{}
	}
	for _, r := range additional {
		set[r] = struct{}{}
	}
	merged := make([]string, 0, len(set))
	for r := range set {
		merged = append(merged, r)
	}
	return merged
}
