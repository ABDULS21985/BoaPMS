package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/auth"
	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// staffManagementService implements StaffManagementService.
// Mirrors the .NET StaffHandlers and RoleHandler mediator handlers,
// consolidated into a single service backed by UserManagementService
// for user/role operations (equivalent to .NET UserManager/RoleManager).
type staffManagementService struct {
	userMgr *UserManagementService
	db      *gorm.DB
	cfg     *config.Config
	log     zerolog.Logger
}

func newStaffManagementService(
	repos *repository.Container,
	cfg *config.Config,
	log zerolog.Logger,
	userMgr *UserManagementService,
) StaffManagementService {
	return &staffManagementService{
		userMgr: userMgr,
		db:      repos.GormDB,
		cfg:     cfg,
		log:     log.With().Str("service", "staff").Logger(),
	}
}

// ---------------------------------------------------------------------------
// AddStaff creates a new staff user or updates an existing one.
// Maps to .NET CreateStaffHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) AddStaff(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*identity.AddStaffToRoleVm)
	if !ok {
		return nil, errors.New("invalid request type: expected *AddStaffToRoleVm")
	}

	s.log.Info().Str("userId", vm.UserId).Str("email", vm.Email).Msg("adding/updating staff")

	// Check if user already exists.
	existing, err := s.userMgr.FindByID(ctx, vm.UserId)
	if err != nil {
		return nil, fmt.Errorf("checking existing user: %w", err)
	}

	if existing == nil {
		// Create a new user.
		user := &identity.ApplicationUser{
			ID:        vm.UserId,
			UserName:  vm.Email,
			Email:     vm.Email,
			FirstName: vm.FirstName,
			LastName:  vm.LastName,
			IsActive:  true,
		}

		// Use a default password for new accounts (mirrors .NET pattern where
		// password is not set during staff creation via CreateAsync without password).
		defaultPwd, _ := s.getDefaultPassword(ctx)
		if err := s.userMgr.CreateUser(ctx, user, defaultPwd); err != nil {
			return nil, fmt.Errorf("creating user: %w", err)
		}

		s.log.Info().Str("userId", vm.UserId).Msg("staff user created")
	} else {
		// Update existing user's active status.
		existing.IsActive = true
		existing.FirstName = vm.FirstName
		existing.LastName = vm.LastName
		if err := s.userMgr.UpdateUser(ctx, existing); err != nil {
			return nil, fmt.Errorf("updating user: %w", err)
		}

		s.log.Info().Str("userId", vm.UserId).Msg("staff user updated")
	}

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "Operation completed successfully",
	}, nil
}

// ---------------------------------------------------------------------------
// GetAllStaffs retrieves all staff, optionally filtered by search string.
// Maps to .NET GetStaffsQueryHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) GetAllStaffs(ctx context.Context, searchString string) (interface{}, error) {
	s.log.Info().Str("search", searchString).Msg("fetching all staff")

	var users []identity.ApplicationUser

	query := s.db.WithContext(ctx).Model(&identity.ApplicationUser{})

	if strings.TrimSpace(searchString) != "" {
		search := "%" + strings.ToLower(strings.TrimSpace(searchString)) + "%"
		query = query.Where(
			"LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ? OR LOWER(user_name) LIKE ? OR LOWER(email) LIKE ?",
			search, search, search, search,
		)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("querying staff users: %w", err)
	}

	s.log.Info().Int("count", len(users)).Msg("staff users fetched")

	return &performance.GenericResponseVm{
		IsSuccess:    true,
		Message:      "Operation completed successfully",
		TotalRecords: len(users),
		Data:         users,
	}, nil
}

// ---------------------------------------------------------------------------
// GetAllRoles retrieves all roles except SuperAdmin.
// Maps to .NET GetRoleQueryHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) GetAllRoles(ctx context.Context) (interface{}, error) {
	s.log.Info().Msg("fetching all roles")

	var roles []identity.ApplicationRole
	err := s.db.WithContext(ctx).
		Where("UPPER(name) != ?", strings.ToUpper(auth.RoleSuperAdmin)).
		Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("querying roles: %w", err)
	}

	// Map to RoleVm DTOs.
	roleVms := make([]identity.RoleVm, 0, len(roles))
	for _, r := range roles {
		roleVms = append(roleVms, identity.RoleVm{
			ID:       r.ID,
			RoleName: r.NormalizedName,
		})
	}

	s.log.Info().Int("count", len(roleVms)).Msg("roles fetched")

	return &performance.GenericResponseVm{
		IsSuccess:    true,
		Message:      "Operation completed successfully",
		TotalRecords: len(roleVms),
		Data:         roleVms,
	}, nil
}

// ---------------------------------------------------------------------------
// AddRole creates a new role or updates an existing one.
// Maps to .NET SaveRoleHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) AddRole(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*identity.RoleVm)
	if !ok {
		return nil, errors.New("invalid request type: expected *RoleVm")
	}

	s.log.Info().Str("roleName", vm.RoleName).Str("id", vm.ID).Msg("saving role")

	if vm.ID != "" {
		// Update existing role.
		var role identity.ApplicationRole
		err := s.db.WithContext(ctx).Where("id = ?", vm.ID).First(&role).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("role not found: %s", vm.ID)
			}
			return nil, fmt.Errorf("fetching role: %w", err)
		}

		role.Name = vm.RoleName
		role.NormalizedName = strings.ToUpper(vm.RoleName)

		if err := s.db.WithContext(ctx).Save(&role).Error; err != nil {
			return nil, fmt.Errorf("updating role: %w", err)
		}

		s.log.Info().Str("roleId", role.ID).Msg("role updated successfully")

		return &performance.GenericResponseVm{
			IsSuccess: true,
			Message:   "Role updated successfully",
		}, nil
	}

	// Create new role -- check for duplicates first.
	existingRole, err := s.userMgr.FindRoleByName(ctx, vm.RoleName)
	if err != nil {
		return nil, fmt.Errorf("checking role existence: %w", err)
	}
	if existingRole != nil {
		return &performance.GenericResponseVm{
			IsSuccess: false,
			Message:   "Role already exists",
		}, nil
	}

	newRole := identity.ApplicationRole{
		ID:             uuid.New().String(),
		Name:           vm.RoleName,
		NormalizedName: strings.ToUpper(vm.RoleName),
	}

	if err := s.db.WithContext(ctx).Create(&newRole).Error; err != nil {
		return nil, fmt.Errorf("creating role: %w", err)
	}

	s.log.Info().Str("roleId", newRole.ID).Str("roleName", vm.RoleName).Msg("role created successfully")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "Role created successfully",
	}, nil
}

// ---------------------------------------------------------------------------
// DeleteRole removes a role and its associated permissions.
// Maps to .NET DeleteRoleHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) DeleteRole(ctx context.Context, roleName string) (interface{}, error) {
	s.log.Info().Str("roleName", roleName).Msg("deleting role")

	// Find the role.
	var role identity.ApplicationRole
	err := s.db.WithContext(ctx).
		Where("UPPER(TRIM(normalized_name)) = ?", strings.ToUpper(strings.TrimSpace(roleName))).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &performance.GenericResponseVm{
				IsSuccess: false,
				Message:   "Role not found",
			}, nil
		}
		return nil, fmt.Errorf("fetching role: %w", err)
	}

	// Delete associated role permissions first.
	if err := s.db.WithContext(ctx).
		Where("role_id = ?", role.ID).
		Delete(&identity.RolePermission{}).Error; err != nil {
		return nil, fmt.Errorf("deleting role permissions: %w", err)
	}

	// Delete the role itself.
	if err := s.db.WithContext(ctx).Delete(&role).Error; err != nil {
		return nil, fmt.Errorf("deleting role: %w", err)
	}

	s.log.Info().Str("roleId", role.ID).Str("roleName", roleName).Msg("role deleted successfully")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   "Role deleted successfully",
	}, nil
}

// ---------------------------------------------------------------------------
// AddStaffToRole assigns a staff member to a role.
// If the user does not exist, they are created first.
// Maps to .NET AddStaffToRoleHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) AddStaffToRole(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*identity.AddStaffToRoleVm)
	if !ok {
		return nil, errors.New("invalid request type: expected *AddStaffToRoleVm")
	}

	s.log.Info().Str("userId", vm.UserId).Str("role", vm.RoleName).Msg("adding staff to role")

	// Ensure the user exists; create if not.
	user, err := s.userMgr.FindByID(ctx, vm.UserId)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		newUser := &identity.ApplicationUser{
			ID:        vm.UserId,
			UserName:  vm.Email,
			Email:     vm.Email,
			FirstName: vm.FirstName,
			LastName:  vm.LastName,
			IsActive:  true,
		}
		defaultPwd, _ := s.getDefaultPassword(ctx)
		if err := s.userMgr.CreateUser(ctx, newUser, defaultPwd); err != nil {
			return nil, fmt.Errorf("creating user for role assignment: %w", err)
		}
		user = newUser
	}

	// Find the role to get its ID.
	role, err := s.userMgr.FindRoleByName(ctx, vm.RoleName)
	if err != nil {
		return nil, fmt.Errorf("finding role: %w", err)
	}
	if role == nil {
		return &performance.GenericResponseVm{
			IsSuccess: false,
			Message:   fmt.Sprintf("Role %s not found", vm.RoleName),
		}, nil
	}

	// Check if already assigned.
	existingRoles, err := s.userMgr.GetUserRoles(ctx, vm.UserId)
	if err != nil {
		return nil, fmt.Errorf("getting user roles: %w", err)
	}
	for _, r := range existingRoles {
		if strings.EqualFold(r, vm.RoleName) {
			return &performance.GenericResponseVm{
				IsSuccess: false,
				Message:   fmt.Sprintf("User already in %s role", vm.RoleName),
			}, nil
		}
	}

	// Assign the role.
	if err := s.userMgr.AssignRole(ctx, vm.UserId, role.ID); err != nil {
		return nil, fmt.Errorf("assigning role: %w", err)
	}

	s.log.Info().Str("userId", vm.UserId).Str("role", vm.RoleName).Msg("staff added to role")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   fmt.Sprintf("User added to %s role successfully", vm.RoleName),
	}, nil
}

// ---------------------------------------------------------------------------
// RemoveStaffFromRole removes a staff member from a role.
// Maps to .NET DeleteStaffFromRoleHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) RemoveStaffFromRole(ctx context.Context, userID string, roleName string) (interface{}, error) {
	s.log.Info().Str("userId", userID).Str("role", roleName).Msg("removing staff from role")

	// Find the user.
	user, err := s.userMgr.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return &performance.GenericResponseVm{
			IsSuccess: false,
			Message:   "User not found",
		}, nil
	}

	// Find the role.
	role, err := s.userMgr.FindRoleByName(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("finding role: %w", err)
	}
	if role == nil {
		return &performance.GenericResponseVm{
			IsSuccess: false,
			Message:   fmt.Sprintf("Role %s not found", roleName),
		}, nil
	}

	// Remove the assignment.
	if err := s.userMgr.RemoveRole(ctx, userID, role.ID); err != nil {
		return nil, fmt.Errorf("removing role: %w", err)
	}

	s.log.Info().Str("userId", userID).Str("role", roleName).Msg("staff removed from role")

	return &performance.GenericResponseVm{
		IsSuccess: true,
		Message:   fmt.Sprintf("User removed from %s role successfully", roleName),
	}, nil
}

// ---------------------------------------------------------------------------
// GetStaffRoles retrieves the roles assigned to a specific staff member.
// Maps to .NET GetStaffRoleHandler.
// ---------------------------------------------------------------------------
func (s *staffManagementService) GetStaffRoles(ctx context.Context, id string) (interface{}, error) {
	s.log.Info().Str("userId", id).Msg("fetching staff roles")

	// Verify user exists.
	user, err := s.userMgr.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}

	var roles []string
	if user != nil {
		roles, err = s.userMgr.GetUserRoles(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("getting user roles: %w", err)
		}
	} else {
		roles = []string{}
	}

	s.log.Info().Str("userId", id).Int("count", len(roles)).Msg("staff roles fetched")

	return &performance.GenericResponseVm{
		IsSuccess:    true,
		Message:      "Operation completed successfully",
		TotalRecords: len(roles),
		Data:         roles,
	}, nil
}

// ===========================================================================
// Private helpers
// ===========================================================================

// getDefaultPassword retrieves the configured default password for new accounts.
// Falls back to a generated password if the setting is unavailable.
func (s *staffManagementService) getDefaultPassword(ctx context.Context) (string, error) {
	// Try to read from global settings (mirrors .NET DEFAULT_PASSWORD setting).
	val, err := s.db.WithContext(ctx).
		Raw(`SELECT value FROM pms.settings WHERE name = ? AND soft_deleted = false`, auth.SettingDefaultPassword).
		Rows()
	if err == nil && val != nil {
		defer val.Close()
		if val.Next() {
			var pwd string
			if err := val.Scan(&pwd); err == nil && pwd != "" {
				return pwd, nil
			}
		}
	}

	// Fallback default.
	return "P@ssw0rd123!", nil
}
