package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// EmployeeInformationHandler handles employee information HTTP endpoints.
// Converts EmployeeInformationController.cs. Route base: api/employees.
// All endpoints require [Authorize].
type EmployeeInformationHandler struct {
	erpSvc service.ErpEmployeeService
	log    zerolog.Logger
}

// NewEmployeeInformationHandler creates a new EmployeeInformationHandler.
func NewEmployeeInformationHandler(erpSvc service.ErpEmployeeService, log zerolog.Logger) *EmployeeInformationHandler {
	return &EmployeeInformationHandler{erpSvc: erpSvc, log: log}
}

// GetEmployeeDetail handles GET /api/employees?employeeNumber=...
// Retrieves employee details by employee number.
func (h *EmployeeInformationHandler) GetEmployeeDetail(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.erpSvc.GetEmployeeDetail(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("Failed to get employee detail")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetHeadSubordinates handles GET /api/employees/head-subordinates?employeeNumber=...
// Retrieves subordinates for heads by employee number.
func (h *EmployeeInformationHandler) GetHeadSubordinates(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.erpSvc.GetHeadSubordinates(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("Failed to get head subordinates")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetEmployeeSubordinates handles GET /api/employees/subordinates?employeeNumber=...
// Retrieves subordinates for a given employee.
func (h *EmployeeInformationHandler) GetEmployeeSubordinates(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.erpSvc.GetEmployeeSubordinates(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("Failed to get employee subordinates")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetEmployeePeers handles GET /api/employees/peers?employeeNumber=...
// Retrieves peers for a given employee based on grade and office.
func (h *EmployeeInformationHandler) GetEmployeePeers(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.erpSvc.GetEmployeePeers(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("Failed to get employee peers")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllByDepartment handles GET /api/employees/by-department?departmentId=...
// Retrieves all employees in a department.
func (h *EmployeeInformationHandler) GetAllByDepartment(w http.ResponseWriter, r *http.Request) {
	deptIdStr := r.URL.Query().Get("departmentId")
	if deptIdStr == "" {
		response.Error(w, http.StatusBadRequest, "departmentId is required")
		return
	}
	departmentId, err := strconv.Atoi(deptIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid departmentId")
		return
	}
	result, err := h.erpSvc.GetAllByDepartmentId(r.Context(), departmentId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get employees by department")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllByDivision handles GET /api/employees/by-division?divisionId=...
// Retrieves all employees in a division.
func (h *EmployeeInformationHandler) GetAllByDivision(w http.ResponseWriter, r *http.Request) {
	divIdStr := r.URL.Query().Get("divisionId")
	if divIdStr == "" {
		response.Error(w, http.StatusBadRequest, "divisionId is required")
		return
	}
	divisionId, err := strconv.Atoi(divIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid divisionId")
		return
	}
	result, err := h.erpSvc.GetAllByDivisionId(r.Context(), divisionId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get employees by division")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllByOffice handles GET /api/employees/by-office?officeId=...
// Retrieves all employees in an office.
func (h *EmployeeInformationHandler) GetAllByOffice(w http.ResponseWriter, r *http.Request) {
	offIdStr := r.URL.Query().Get("officeId")
	if offIdStr == "" {
		response.Error(w, http.StatusBadRequest, "officeId is required")
		return
	}
	officeId, err := strconv.Atoi(offIdStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid officeId")
		return
	}
	result, err := h.erpSvc.GetAllByOfficeId(r.Context(), officeId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get employees by office")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllEmployees handles GET /api/employees/all
// Retrieves all employees.
func (h *EmployeeInformationHandler) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	result, err := h.erpSvc.GetAllEmployees(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get all employees")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SeedOrganizationData handles GET /api/employees/seed-organization
// Seeds the database with organization data from ERP.
func (h *EmployeeInformationHandler) SeedOrganizationData(w http.ResponseWriter, r *http.Request) {
	result, err := h.erpSvc.SeedOrganizationData(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to seed organization data")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffIDMaskDetail handles GET /api/employees/staff-id-mask?employeeNumber=...
// Retrieves staff ID mask detail for a given employee.
func (h *EmployeeInformationHandler) GetStaffIDMaskDetail(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.erpSvc.GetStaffIDMaskDetail(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("Failed to get staff ID mask detail")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// UpdateStaffJobRole handles POST /api/employees/staff-job-role
// Creates or updates a staff job role.
func (h *EmployeeInformationHandler) UpdateStaffJobRole(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.erpSvc.UpdateStaffJobRole(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to update staff job role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffJobRoleById handles GET /api/employees/staff-job-role?employeeNumber=...
// Retrieves the job role for a specific employee.
func (h *EmployeeInformationHandler) GetStaffJobRoleById(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.erpSvc.GetStaffJobRoleById(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("Failed to get staff job role by id")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetJobRolesByOffice handles POST /api/employees/job-roles-by-office
// Retrieves job roles filtered by office.
func (h *EmployeeInformationHandler) GetJobRolesByOffice(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.erpSvc.GetJobRolesByOffice(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get job roles by office")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffJobRoleRequests handles GET /api/employees/staff-job-role-requests?employeeNumber=...
// Retrieves staff job role update requests for a given employee.
func (h *EmployeeInformationHandler) GetStaffJobRoleRequests(w http.ResponseWriter, r *http.Request) {
	employeeNumber := r.URL.Query().Get("employeeNumber")
	if employeeNumber == "" {
		response.Error(w, http.StatusBadRequest, "employeeNumber is required")
		return
	}
	result, err := h.erpSvc.GetStaffJobRoleRequests(r.Context(), employeeNumber)
	if err != nil {
		h.log.Error().Err(err).Str("employeeNumber", employeeNumber).Msg("Failed to get staff job role requests")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveRejectStaffJobRole handles POST /api/employees/approve-reject-staff-job-role
// Approves or rejects a staff job role update.
func (h *EmployeeInformationHandler) ApproveRejectStaffJobRole(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.erpSvc.ApproveRejectStaffJobRole(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to approve/reject staff job role")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}
