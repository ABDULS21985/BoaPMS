package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserManagementService handles local user CRUD and password management.
// Mirrors the .NET UserManager<ApplicationUser> operations.
type UserManagementService struct {
	db  *gorm.DB
	log zerolog.Logger
}

// NewUserManagementService creates a user management service.
func NewUserManagementService(repos *repository.Container, log zerolog.Logger) *UserManagementService {
	return &UserManagementService{
		db:  repos.GormDB,
		log: log,
	}
}

// FindByUsername retrieves a user by username (case-insensitive).
func (s *UserManagementService) FindByUsername(ctx context.Context, username string) (*identity.ApplicationUser, error) {
	var user identity.ApplicationUser
	err := s.db.WithContext(ctx).
		Where("normalized_user_name = ?", strings.ToUpper(username)).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("finding user by username: %w", err)
	}
	return &user, nil
}

// FindByID retrieves a user by their primary key.
func (s *UserManagementService) FindByID(ctx context.Context, id string) (*identity.ApplicationUser, error) {
	var user identity.ApplicationUser
	err := s.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("finding user by id: %w", err)
	}
	return &user, nil
}

// CreateUser registers a new local user with a bcrypt-hashed password.
func (s *UserManagementService) CreateUser(ctx context.Context, user *identity.ApplicationUser, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	user.PasswordHash = string(hash)
	user.NormalizedUserName = strings.ToUpper(user.UserName)
	user.NormalizedEmail = strings.ToUpper(user.Email)
	user.SecurityStamp = generateSecurityStamp()
	user.IsActive = true

	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

// VerifyPassword checks a plaintext password against the stored hash.
func (s *UserManagementService) VerifyPassword(user *identity.ApplicationUser, password string) bool {
	if user.PasswordHash == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}

// UpdateUser persists changes to an existing user record.
func (s *UserManagementService) UpdateUser(ctx context.Context, user *identity.ApplicationUser) error {
	return s.db.WithContext(ctx).Save(user).Error
}

// DeactivateUser sets a user as inactive (soft deactivation).
func (s *UserManagementService) DeactivateUser(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).
		Model(&identity.ApplicationUser{}).
		Where("id = ?", userID).
		Update("is_active", false).Error
}

// IncrementAccessFailedCount increments the failed login counter and locks
// the account if max attempts are exceeded.
func (s *UserManagementService) IncrementAccessFailedCount(ctx context.Context, user *identity.ApplicationUser, maxAttempts int, lockoutMinutes int) error {
	user.AccessFailedCount++
	if maxAttempts > 0 && user.AccessFailedCount >= maxAttempts && user.LockoutEnabled {
		lockoutEnd := time.Now().UTC().Add(time.Duration(lockoutMinutes) * time.Minute)
		user.LockoutEnd = &lockoutEnd
	}
	return s.db.WithContext(ctx).Save(user).Error
}

// ResetAccessFailedCount resets the failed login counter after a successful login.
func (s *UserManagementService) ResetAccessFailedCount(ctx context.Context, user *identity.ApplicationUser) error {
	user.AccessFailedCount = 0
	user.LockoutEnd = nil
	return s.db.WithContext(ctx).Save(user).Error
}

// IsLockedOut checks if the user is currently locked out.
func (s *UserManagementService) IsLockedOut(user *identity.ApplicationUser) bool {
	if !user.LockoutEnabled || user.LockoutEnd == nil {
		return false
	}
	return user.LockoutEnd.After(time.Now().UTC())
}

// GetUserRoles retrieves the roles assigned to a user via the join table.
func (s *UserManagementService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var roleNames []string
	err := s.db.WithContext(ctx).
		Table(`"CoreSchema".asp_net_user_roles ur`).
		Joins(`JOIN "CoreSchema".asp_net_roles r ON r.id = ur.role_id`).
		Where("ur.user_id = ?", userID).
		Pluck("r.name", &roleNames).Error
	return roleNames, err
}

// AssignRole adds a role to a user (idempotent).
func (s *UserManagementService) AssignRole(ctx context.Context, userID, roleID string) error {
	// Use raw SQL for the join table since there's no GORM model for it
	return s.db.WithContext(ctx).Exec(
		`INSERT INTO "CoreSchema".asp_net_user_roles (user_id, role_id)
		 VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`, userID, roleID,
	).Error
}

// RemoveRole removes a role from a user.
func (s *UserManagementService) RemoveRole(ctx context.Context, userID, roleID string) error {
	return s.db.WithContext(ctx).Exec(
		`DELETE FROM "CoreSchema".asp_net_user_roles WHERE user_id = $1 AND role_id = $2`,
		userID, roleID,
	).Error
}

// FindRoleByName finds a role by its normalized name.
func (s *UserManagementService) FindRoleByName(ctx context.Context, name string) (*identity.ApplicationRole, error) {
	var role identity.ApplicationRole
	err := s.db.WithContext(ctx).
		Where("normalized_name = ?", strings.ToUpper(name)).
		First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetPermissionsByRoles retrieves all permission names assigned to any of the given roles.
// Queries the RolePermission junction table joined with Permission.
// Mirrors the .NET Permission/RolePermission fine-grained authorization model.
func (s *UserManagementService) GetPermissionsByRoles(ctx context.Context, roleNames []string) ([]string, error) {
	if len(roleNames) == 0 {
		return nil, nil
	}

	var permNames []string
	err := s.db.WithContext(ctx).
		Table(`"CoreSchema".role_permissions rp`).
		Joins(`JOIN "CoreSchema".asp_net_roles r ON r.id = rp.role_id`).
		Joins(`JOIN "CoreSchema".permissions p ON p.permission_id = rp.permission_id`).
		Where("r.name IN ? AND rp.soft_deleted = false AND p.soft_deleted = false", roleNames).
		Distinct().
		Pluck("p.name", &permNames).Error
	return permNames, err
}

// GetOrganizationalUnit retrieves the user's department/division/office via ERP lookup.
// Returns a string like "IT Department > Network Division > Support Office".
func (s *UserManagementService) GetOrganizationalUnit(ctx context.Context, erpSQL interface{}, staffID string) string {
	type orgUnit struct {
		DepartmentName string `db:"department_name"`
		DivisionName   string `db:"division_name"`
		OfficeName     string `db:"office_name"`
	}

	// If we have an ERP SQL connection, try to get the org info
	if db, ok := erpSQL.(interface {
		GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	}); ok {
		var ou orgUnit
		err := db.GetContext(ctx, &ou,
			`SELECT COALESCE(d.department_name, '') as department_name,
			        COALESCE(dv.division_name, '') as division_name,
			        COALESCE(o.office_name, '') as office_name
			 FROM employees e
			 LEFT JOIN departments d ON d.department_id = e.department_id
			 LEFT JOIN divisions dv ON dv.division_id = e.division_id
			 LEFT JOIN offices o ON o.office_id = e.office_id
			 WHERE e.staff_id = $1`, staffID)
		if err == nil {
			parts := make([]string, 0, 3)
			if ou.DepartmentName != "" {
				parts = append(parts, ou.DepartmentName)
			}
			if ou.DivisionName != "" {
				parts = append(parts, ou.DivisionName)
			}
			if ou.OfficeName != "" {
				parts = append(parts, ou.OfficeName)
			}
			if len(parts) > 0 {
				return strings.Join(parts, " > ")
			}
		}
	}
	return ""
}

func generateSecurityStamp() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
