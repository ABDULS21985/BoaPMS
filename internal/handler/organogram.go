package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// OrganogramHandler handles organizational structure HTTP endpoints.
// Converts .NET OrganogramController — route base: api/organogram.
type OrganogramHandler struct {
	organogramSvc service.OrganogramService
	erpSvc        service.ErpEmployeeService
	log           zerolog.Logger
}

// NewOrganogramHandler creates a new organogram handler.
func NewOrganogramHandler(organogramSvc service.OrganogramService, erpSvc service.ErpEmployeeService, log zerolog.Logger) *OrganogramHandler {
	return &OrganogramHandler{organogramSvc: organogramSvc, erpSvc: erpSvc, log: log}
}

// ---------------------------------------------------------------------------
// Directorates
// ---------------------------------------------------------------------------

// GetDirectorates handles GET /api/organogram/directorates
// Retrieves all directorates.
func (h *OrganogramHandler) GetDirectorates(w http.ResponseWriter, r *http.Request) {
	result, err := h.organogramSvc.GetDirectorates(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get directorates")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveDirectorate handles POST /api/organogram/directorates
// Creates or updates a directorate.
func (h *OrganogramHandler) SaveDirectorate(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.organogramSvc.SaveDirectorate(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save directorate")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// Departments
// ---------------------------------------------------------------------------

// GetDepartments handles GET /api/organogram/departments
// Retrieves all departments, optionally filtered by directorateId query parameter.
func (h *OrganogramHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	var directorateId *int
	if dirStr := r.URL.Query().Get("directorateId"); dirStr != "" {
		id, err := strconv.Atoi(dirStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid directorateId")
			return
		}
		directorateId = &id
	}
	result, err := h.organogramSvc.GetDepartments(r.Context(), directorateId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get departments")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveDepartment handles POST /api/organogram/departments
// Creates or updates a department.
func (h *OrganogramHandler) SaveDepartment(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.organogramSvc.SaveDepartment(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save department")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// Divisions
// ---------------------------------------------------------------------------

// GetDivisions handles GET /api/organogram/divisions
// Retrieves divisions, optionally filtered by departmentId query parameter.
func (h *OrganogramHandler) GetDivisions(w http.ResponseWriter, r *http.Request) {
	var departmentId *int
	if deptStr := r.URL.Query().Get("departmentId"); deptStr != "" {
		id, err := strconv.Atoi(deptStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid departmentId")
			return
		}
		departmentId = &id
	}
	result, err := h.organogramSvc.GetDivisions(r.Context(), departmentId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get divisions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveDivision handles POST /api/organogram/divisions
// Creates or updates a division.
func (h *OrganogramHandler) SaveDivision(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.organogramSvc.SaveDivision(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save division")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// Offices
// ---------------------------------------------------------------------------

// GetOffices handles GET /api/organogram/offices
// Retrieves offices, optionally filtered by divisionId query parameter.
func (h *OrganogramHandler) GetOffices(w http.ResponseWriter, r *http.Request) {
	var divisionId *int
	if divStr := r.URL.Query().Get("divisionId"); divStr != "" {
		id, err := strconv.Atoi(divStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid divisionId")
			return
		}
		divisionId = &id
	}
	result, err := h.organogramSvc.GetOffices(r.Context(), divisionId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get offices")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveOffice handles POST /api/organogram/offices
// Creates or updates an office.
func (h *OrganogramHandler) SaveOffice(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := h.organogramSvc.SaveOffice(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to save office")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ---------------------------------------------------------------------------
// ERP Organogram
// ---------------------------------------------------------------------------

// GetErpDepartments handles GET /api/organogram/erp/departments
// Retrieves all departments from the ERP system.
func (h *OrganogramHandler) GetErpDepartments(w http.ResponseWriter, r *http.Request) {
	result, err := h.erpSvc.GetAllDepartments(r.Context())
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get ERP departments")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetErpDivisions handles GET /api/organogram/erp/divisions
// Retrieves divisions from the ERP system, optionally filtered by departmentId query parameter.
func (h *OrganogramHandler) GetErpDivisions(w http.ResponseWriter, r *http.Request) {
	var departmentId *int
	if deptStr := r.URL.Query().Get("departmentId"); deptStr != "" {
		id, err := strconv.Atoi(deptStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid departmentId")
			return
		}
		departmentId = &id
	}
	result, err := h.erpSvc.GetAllDivisions(r.Context(), departmentId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get ERP divisions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetErpOffices handles GET /api/organogram/erp/offices
// Retrieves offices from the ERP system, optionally filtered by divisionId query parameter.
func (h *OrganogramHandler) GetErpOffices(w http.ResponseWriter, r *http.Request) {
	var divisionId *int
	if divStr := r.URL.Query().Get("divisionId"); divStr != "" {
		id, err := strconv.Atoi(divStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid divisionId")
			return
		}
		divisionId = &id
	}
	result, err := h.erpSvc.GetAllOffices(r.Context(), divisionId)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get ERP offices")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}
