package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// GrievanceHandler handles grievance management HTTP endpoints.
// Mirrors the .NET PmsGrievanceController with grievance CRUD and resolution endpoints.
type GrievanceHandler struct {
	svc *service.Container
	log zerolog.Logger
}

// NewGrievanceHandler creates a new grievance handler.
func NewGrievanceHandler(svc *service.Container, log zerolog.Logger) *GrievanceHandler {
	return &GrievanceHandler{svc: svc, log: log}
}

// --- Request DTOs ---

// CreateGrievanceRequest is the payload for raising a new grievance.
type CreateGrievanceRequest struct {
	ComplainantStaffID int    `json:"complainantStaffId"`
	RespondentStaffID  int    `json:"respondentStaffId"`
	GrievanceType      int    `json:"grievanceType"`
	ReviewPeriodID     int    `json:"reviewPeriodId"`
	Description        string `json:"description"`
}

// GrievanceRequest is the payload for updating an existing grievance.
type GrievanceRequest struct {
	ID                 int    `json:"id"`
	ComplainantStaffID int    `json:"complainantStaffId"`
	RespondentStaffID  int    `json:"respondentStaffId"`
	GrievanceType      int    `json:"grievanceType"`
	ReviewPeriodID     int    `json:"reviewPeriodId"`
	Description        string `json:"description"`
}

// CreateGrievanceResolutionRequest is the payload for creating a grievance resolution.
type CreateGrievanceResolutionRequest struct {
	GrievanceID     int    `json:"grievanceId"`
	ResolutionLevel int    `json:"resolutionLevel"`
	ResolvedByID    int    `json:"resolvedById"`
	Resolution      string `json:"resolution"`
	Remark          int    `json:"remark"`
}

// GrievanceResolutionRequest is the payload for updating a grievance resolution.
type GrievanceResolutionRequest struct {
	ID              int    `json:"id"`
	GrievanceID     int    `json:"grievanceId"`
	ResolutionLevel int    `json:"resolutionLevel"`
	ResolvedByID    int    `json:"resolvedById"`
	Resolution      string `json:"resolution"`
	Remark          int    `json:"remark"`
}

// --- Handlers ---

// RaiseNewGrievance handles POST /api/v1/grievances
// Mirrors .NET PerformanceMgtController.RaiseNewGrievance — creates a new grievance record.
func (h *GrievanceHandler) RaiseNewGrievance(w http.ResponseWriter, r *http.Request) {
	var req CreateGrievanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ComplainantStaffID == 0 || req.RespondentStaffID == 0 {
		response.Error(w, http.StatusBadRequest, "Complainant and respondent staff IDs are required")
		return
	}

	if req.Description == "" {
		response.Error(w, http.StatusBadRequest, "Description is required")
		return
	}

	result, err := h.svc.Grievance.RaiseNewGrievance(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RaiseNewGrievance").Msg("Failed to raise new grievance")
		response.Error(w, http.StatusInternalServerError, "Failed to raise grievance")
		return
	}

	response.Created(w, result)
}

// UpdateGrievance handles PUT /api/v1/grievances
// Mirrors .NET PerformanceMgtController.UpadateGrievance — updates an existing grievance.
func (h *GrievanceHandler) UpdateGrievance(w http.ResponseWriter, r *http.Request) {
	var req GrievanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ID == 0 {
		response.Error(w, http.StatusBadRequest, "Grievance ID is required")
		return
	}

	result, err := h.svc.Grievance.UpdateGrievance(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateGrievance").Msg("Failed to update grievance")
		response.Error(w, http.StatusInternalServerError, "Failed to update grievance")
		return
	}

	response.OK(w, result)
}

// CreateGrievanceResolution handles POST /api/v1/grievances/resolutions
// Mirrors .NET PerformanceMgtController.CreateGrievanceResolution — logs a new resolution.
func (h *GrievanceHandler) CreateGrievanceResolution(w http.ResponseWriter, r *http.Request) {
	var req CreateGrievanceResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.GrievanceID == 0 {
		response.Error(w, http.StatusBadRequest, "Grievance ID is required")
		return
	}

	if req.Resolution == "" {
		response.Error(w, http.StatusBadRequest, "Resolution text is required")
		return
	}

	result, err := h.svc.Grievance.CreateGrievanceResolution(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateGrievanceResolution").Msg("Failed to create grievance resolution")
		response.Error(w, http.StatusInternalServerError, "Failed to create grievance resolution")
		return
	}

	response.Created(w, result)
}

// UpdateGrievanceResolution handles PUT /api/v1/grievances/resolutions
// Mirrors .NET PerformanceMgtController.UpadateGrievanceResolution — updates an existing resolution.
func (h *GrievanceHandler) UpdateGrievanceResolution(w http.ResponseWriter, r *http.Request) {
	var req GrievanceResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ID == 0 {
		response.Error(w, http.StatusBadRequest, "Resolution ID is required")
		return
	}

	result, err := h.svc.Grievance.UpdateGrievanceResolution(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateGrievanceResolution").Msg("Failed to update grievance resolution")
		response.Error(w, http.StatusInternalServerError, "Failed to update grievance resolution")
		return
	}

	response.OK(w, result)
}

// GetStaffGrievances handles GET /api/v1/grievances?staffId={id}
// Mirrors .NET PerformanceMgtController.GetStaffGrievances — returns all grievances for a specific staff member.
func (h *GrievanceHandler) GetStaffGrievances(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staffId")
	if staffID == "" {
		response.Error(w, http.StatusBadRequest, "Staff ID query parameter is required")
		return
	}

	result, err := h.svc.Grievance.GetStaffGrievances(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffGrievances").Str("staffId", staffID).Msg("Failed to get staff grievances")
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve staff grievances")
		return
	}

	response.OK(w, result)
}

// GetGrievancesReport handles GET /api/v1/grievances/report
// Mirrors .NET PerformanceMgtController.GetGrievancesReport — returns an aggregated grievances report.
func (h *GrievanceHandler) GetGrievancesReport(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Grievance.GetGrievancesReport(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetGrievancesReport").Msg("Failed to get grievances report")
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve grievances report")
		return
	}

	response.OK(w, result)
}
