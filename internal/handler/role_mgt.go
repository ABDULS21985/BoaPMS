package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// RoleMgtHandler handles role management HTTP endpoints.
// Converts the .NET RoleMgtController (route base: api/rolemgmt).
// All endpoints require [Authorize].
type RoleMgtHandler struct {
	roleSvc service.RoleManagementService
	log     zerolog.Logger
}

// NewRoleMgtHandler creates a new role management handler.
func NewRoleMgtHandler(roleSvc service.RoleManagementService, log zerolog.Logger) *RoleMgtHandler {
	return &RoleMgtHandler{roleSvc: roleSvc, log: log}
}

// --- Request DTOs ---

// AddPermissionToRoleRequest mirrors .NET AddPermissionToRoleVm.
type AddPermissionToRoleRequest struct {
	RoleID       string `json:"roleId"`
	PermissionID int    `json:"permissionId"`
}

// --- Handlers ---

// GetPermissions handles GET /api/rolemgmt/getPermissions?roleId={roleId}
// Mirrors .NET RoleMgtController.GetPermissions — retrieves all permissions,
// optionally filtered by roleId.
func (h *RoleMgtHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	roleId := r.URL.Query().Get("roleId") // optional
	result, err := h.roleSvc.GetPermissions(r.Context(), roleId)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPermissions").Str("roleId", roleId).Msg("Failed to get permissions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllRolesWithPermission handles GET /api/rolemgmt/getAllRolesWithPermission?roleId={roleId}
// Mirrors .NET RoleMgtController.GetAllRolesWithPermission — retrieves all permissions
// in the application alongside the permissions assigned to the specified role.
func (h *RoleMgtHandler) GetAllRolesWithPermission(w http.ResponseWriter, r *http.Request) {
	roleId := r.URL.Query().Get("roleId")
	if roleId == "" {
		response.Error(w, http.StatusBadRequest, "roleId is required")
		return
	}
	result, err := h.roleSvc.GetAllRolesWithPermission(r.Context(), roleId)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllRolesWithPermission").Str("roleId", roleId).Msg("Failed to get all roles with permission")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddPermissionToRole handles POST /api/rolemgmt/addPermissionToRole
// Mirrors .NET RoleMgtController.AddPermissionToRole — assigns a permission to a role.
// Expects JSON body: {"roleId": "...", "permissionId": 123}
func (h *RoleMgtHandler) AddPermissionToRole(w http.ResponseWriter, r *http.Request) {
	var req AddPermissionToRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.RoleID == "" {
		response.Error(w, http.StatusBadRequest, "roleId is required")
		return
	}
	if req.PermissionID == 0 {
		response.Error(w, http.StatusBadRequest, "permissionId is required")
		return
	}
	result, err := h.roleSvc.AddPermissionToRole(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddPermissionToRole").Str("roleId", req.RoleID).Int("permissionId", req.PermissionID).Msg("Failed to add permission to role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RemovePermissionFromRole handles DELETE /api/rolemgmt/removePermissionFromRole?roleId={roleId}&permissionId={permissionId}
// Mirrors .NET RoleMgtController.DeletePermissionInRole — removes a permission from a role.
func (h *RoleMgtHandler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	roleId := r.URL.Query().Get("roleId")
	permIdStr := r.URL.Query().Get("permissionId")
	if roleId == "" || permIdStr == "" {
		response.Error(w, http.StatusBadRequest, "roleId and permissionId are required")
		return
	}
	permissionId, err := strconv.Atoi(permIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid permissionId")
		return
	}
	result, err := h.roleSvc.RemovePermissionFromRole(r.Context(), roleId, permissionId)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RemovePermissionFromRole").Str("roleId", roleId).Int("permissionId", permissionId).Msg("Failed to remove permission from role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}
