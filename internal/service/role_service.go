package service

import (
	"context"
	"fmt"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// roleManagementService handles role-permission CRUD operations.
// Mirrors the .NET RolePermissionMgt mediator handlers:
//   - GetRolePermissionQueryHandler      -> GetAllRolesWithPermission
//   - GetPermissionsByRoleQueryHandler   -> GetPermissions
//   - AddPermissionToRoleHandler         -> AddPermissionToRole
//   - RemovePermissionFromRoleCmdHandler -> RemovePermissionFromRole
type roleManagementService struct {
	permissionRepo     *repository.Repository[identity.Permission]
	rolePermissionRepo *repository.Repository[identity.RolePermission]
	db                 *gorm.DB
	log                zerolog.Logger
}

// newRoleManagementService creates a RoleManagementService with all required repositories.
func newRoleManagementService(repos *repository.Container, cfg *config.Config, log zerolog.Logger) RoleManagementService {
	return &roleManagementService{
		permissionRepo:     repository.NewRepository[identity.Permission](repos.GormDB),
		rolePermissionRepo: repository.NewRepository[identity.RolePermission](repos.GormDB),
		db:                 repos.GormDB,
		log:                log.With().Str("service", "roleManagement").Logger(),
	}
}

// ---------------------------------------------------------------------------
// GetPermissions
// ---------------------------------------------------------------------------

// GetPermissions retrieves permissions for a given role. If roleId is non-empty,
// it returns only the permissions assigned to that role (via RolePermission join).
// If roleId is empty, it returns all permissions in the system.
// Mirrors .NET GetPermissionsByRoleQueryQueryHandler.
func (s *roleManagementService) GetPermissions(ctx context.Context, roleId string) (interface{}, error) {
	if roleId != "" {
		// Fetch role with its permissions via eager-loading
		var role identity.ApplicationRole
		err := s.db.WithContext(ctx).
			Where("id = ?", roleId).
			Preload("RolePermissions", "soft_deleted = false").
			Preload("RolePermissions.Permission").
			First(&role).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				s.log.Warn().Str("roleId", roleId).Msg("role not found for permission lookup")
				return []identity.PermissionVm{}, nil
			}
			s.log.Error().Err(err).Str("roleId", roleId).Msg("failed to retrieve role permissions")
			return nil, fmt.Errorf("retrieving role permissions: %w", err)
		}

		vms := make([]identity.PermissionVm, 0, len(role.RolePermissions))
		for _, rp := range role.RolePermissions {
			if rp.Permission != nil {
				vms = append(vms, identity.PermissionVm{
					ID:          rp.Permission.PermissionID,
					Name:        rp.Permission.Name,
					Description: rp.Permission.Description,
				})
			}
		}

		s.log.Debug().Str("roleId", roleId).Int("count", len(vms)).Msg("role permissions retrieved")
		return vms, nil
	}

	// No roleId: return all permissions
	permissions, err := s.permissionRepo.GetAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve all permissions")
		return nil, fmt.Errorf("retrieving all permissions: %w", err)
	}

	vms := make([]identity.PermissionVm, 0, len(permissions))
	for _, p := range permissions {
		vms = append(vms, identity.PermissionVm{
			ID:          p.PermissionID,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	s.log.Debug().Int("count", len(vms)).Msg("all permissions retrieved")
	return vms, nil
}

// ---------------------------------------------------------------------------
// GetAllRolesWithPermission
// ---------------------------------------------------------------------------

// GetAllRolesWithPermission retrieves all system permissions alongside the
// permissions specifically assigned to the given role.
// Mirrors .NET GetRolePermissionQueryHandler which returns GetRolePermissionVm
// containing AllPermissions and RolesAndPermissions.
func (s *roleManagementService) GetAllRolesWithPermission(ctx context.Context, roleId string) (interface{}, error) {
	result := identity.GetRolePermissionVm{}

	// 1. Fetch all permissions
	allPerms, err := s.permissionRepo.GetAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to retrieve all permissions")
		return nil, fmt.Errorf("retrieving all permissions: %w", err)
	}

	allPermVms := make([]identity.PermissionVm, 0, len(allPerms))
	for _, p := range allPerms {
		allPermVms = append(allPermVms, identity.PermissionVm{
			ID:          p.PermissionID,
			Name:        p.Name,
			Description: p.Description,
		})
	}
	result.AllPermissions = allPermVms

	// 2. Fetch role with its assigned permissions (if roleId is provided)
	if roleId != "" {
		var role identity.ApplicationRole
		err := s.db.WithContext(ctx).
			Where("id = ?", roleId).
			Preload("RolePermissions", "soft_deleted = false").
			Preload("RolePermissions.Permission").
			First(&role).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			s.log.Error().Err(err).Str("roleId", roleId).Msg("failed to retrieve role with permissions")
			return nil, fmt.Errorf("retrieving role with permissions: %w", err)
		}

		if err == nil {
			rolePermVms := make([]identity.PermissionVm, 0, len(role.RolePermissions))
			for _, rp := range role.RolePermissions {
				if rp.Permission != nil {
					rolePermVms = append(rolePermVms, identity.PermissionVm{
						ID:          rp.Permission.PermissionID,
						Name:        rp.Permission.Name,
						Description: rp.Permission.Description,
					})
				}
			}

			result.RolesAndPermissions = &identity.RolePermissionVm{
				RoleId:      role.ID,
				RoleName:    role.Name,
				Permissions: rolePermVms,
			}
		}
	}

	s.log.Debug().
		Int("allPermissions", len(result.AllPermissions)).
		Str("roleId", roleId).
		Msg("roles with permissions retrieved")
	return &result, nil
}

// ---------------------------------------------------------------------------
// AddPermissionToRole
// ---------------------------------------------------------------------------

// AddPermissionToRole assigns a permission to a role.
// Mirrors .NET AddPermissionToRoleHandler: checks for duplicates before inserting.
func (s *roleManagementService) AddPermissionToRole(ctx context.Context, req interface{}) (interface{}, error) {
	vm, ok := req.(*identity.AddPermissionToRoleVm)
	if !ok {
		return nil, fmt.Errorf("invalid request type: expected *identity.AddPermissionToRoleVm")
	}

	// Check for existing assignment (duplicate prevention)
	existing, err := s.rolePermissionRepo.FirstOrDefault(ctx,
		"role_id = ? AND permission_id = ?", vm.RoleId, vm.PermissionId)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to check existing role permission")
		return nil, fmt.Errorf("checking existing role permission: %w", err)
	}
	if existing != nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   "Permission already assigned to role",
		}, nil
	}

	// Create the role-permission mapping
	now := time.Now().UTC()
	entity := &identity.RolePermission{
		RoleID:       vm.RoleId,
		PermissionID: vm.PermissionId,
	}
	entity.DateCreated = now
	entity.SoftDeleted = false

	if err := s.rolePermissionRepo.Create(ctx, entity); err != nil {
		s.log.Error().Err(err).
			Str("roleId", vm.RoleId).
			Int("permissionId", vm.PermissionId).
			Msg("failed to assign permission to role")
		return nil, fmt.Errorf("assigning permission to role: %w", err)
	}

	s.log.Info().
		Str("roleId", vm.RoleId).
		Int("permissionId", vm.PermissionId).
		Msg("permission assigned to role")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   "Permission assigned successfully",
	}, nil
}

// ---------------------------------------------------------------------------
// RemovePermissionFromRole
// ---------------------------------------------------------------------------

// RemovePermissionFromRole removes a permission assignment from a role.
// Mirrors .NET RemovePermissionFromRoleCmdHandler: checks existence then deletes.
func (s *roleManagementService) RemovePermissionFromRole(ctx context.Context, roleId string, permissionId int) (interface{}, error) {
	// Verify the assignment exists
	existing, err := s.rolePermissionRepo.FirstOrDefault(ctx,
		"role_id = ? AND permission_id = ?", roleId, permissionId)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to look up role permission for removal")
		return nil, fmt.Errorf("looking up role permission: %w", err)
	}
	if existing == nil {
		return map[string]interface{}{
			"isSuccess": false,
			"message":   "Permission to remove was not previously assigned to role",
		}, nil
	}

	// Hard delete the role-permission mapping (mirrors .NET DeleteAsync behaviour)
	err = s.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleId, permissionId).
		Delete(&identity.RolePermission{}).Error
	if err != nil {
		s.log.Error().Err(err).
			Str("roleId", roleId).
			Int("permissionId", permissionId).
			Msg("failed to remove permission from role")
		return nil, fmt.Errorf("removing permission from role: %w", err)
	}

	s.log.Info().
		Str("roleId", roleId).
		Int("permissionId", permissionId).
		Msg("permission removed from role")
	return map[string]interface{}{
		"isSuccess": true,
		"message":   "Permission removed successfully",
	}, nil
}
