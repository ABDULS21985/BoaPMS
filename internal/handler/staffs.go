package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// StaffsHandler handles staff management HTTP endpoints.
// Converts StaffsController.cs. Route base: api/staff. All endpoints require [Authorize].
type StaffsHandler struct {
	staffSvc service.StaffManagementService
	log      zerolog.Logger
}

// NewStaffsHandler creates a new StaffsHandler.
func NewStaffsHandler(staffSvc service.StaffManagementService, log zerolog.Logger) *StaffsHandler {
	return &StaffsHandler{staffSvc: staffSvc, log: log}
}

// AddStaff handles POST /api/staff
// Creates a new staff member.
func (h *StaffsHandler) AddStaff(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.staffSvc.AddStaff(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to add staff")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllStaffs handles GET /api/staff?searchString={query}
// Retrieves all staff members, optionally filtered by search string.
func (h *StaffsHandler) GetAllStaffs(w http.ResponseWriter, r *http.Request) {
	searchString := r.URL.Query().Get("searchString")
	result, err := h.staffSvc.GetAllStaffs(r.Context(), searchString)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get all staffs")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllRoles handles GET /api/staff/roles
// Retrieves all available roles.
func (h *StaffsHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	result, err := h.staffSvc.GetAllRoles(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get all roles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddRole handles POST /api/staff/role
// Creates a new role.
func (h *StaffsHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.staffSvc.AddRole(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to add role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// DeleteRole handles DELETE /api/staff/role?roleName={name}
// Deletes a role by name.
func (h *StaffsHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleName := r.URL.Query().Get("roleName")
	if roleName == "" {
		response.Error(w, http.StatusBadRequest, "roleName is required")
		return
	}
	result, err := h.staffSvc.DeleteRole(r.Context(), roleName)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to delete role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddStaffToRole handles POST /api/staff/role/assign
// Assigns a staff member to a role.
func (h *StaffsHandler) AddStaffToRole(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.staffSvc.AddStaffToRole(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to add staff to role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RemoveStaffFromRole handles DELETE /api/staff/role/remove?userId={id}&roleName={name}
// Removes a staff member from a role.
func (h *StaffsHandler) RemoveStaffFromRole(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	roleName := r.URL.Query().Get("roleName")
	if userId == "" || roleName == "" {
		response.Error(w, http.StatusBadRequest, "userId and roleName are required")
		return
	}
	result, err := h.staffSvc.RemoveStaffFromRole(r.Context(), userId, roleName)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to remove staff from role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffRoles handles GET /api/staff/roles?id={staffId}
// Retrieves all roles assigned to a specific staff member.
func (h *StaffsHandler) GetStaffRoles(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := h.staffSvc.GetStaffRoles(r.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get staff roles")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}
