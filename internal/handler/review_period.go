package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// ReviewPeriodHandler handles review period HTTP endpoints.
// Mirrors the .NET PmsReviewPeriodController (partial class of PerformanceMgtController).
type ReviewPeriodHandler struct {
	svc *service.Container
	log zerolog.Logger
}

// NewReviewPeriodHandler creates a new review period handler.
func NewReviewPeriodHandler(svc *service.Container, log zerolog.Logger) *ReviewPeriodHandler {
	return &ReviewPeriodHandler{svc: svc, log: log}
}

// ---------------------------------------------------------------------------
// Request DTOs
// ---------------------------------------------------------------------------

// ReviewPeriodRequest is the payload for review period lifecycle operations
// (draft, add, submit, approve, reject, return, resubmit, update, cancel, close,
// and toggle operations such as enable/disable objective planning, work product
// planning, and work product evaluation).
// Mirrors .NET ReviewPeriodRequestVm / CreateNewReviewPeriodVm.
type ReviewPeriodRequest struct {
	ID                        string `json:"id,omitempty"`
	Name                      string `json:"name"`
	Description               string `json:"description,omitempty"`
	StartDate                 string `json:"startDate"`
	EndDate                   string `json:"endDate"`
	Year                      int    `json:"year,omitempty"`
	ReviewPeriodTypeID        string `json:"reviewPeriodTypeId,omitempty"`
	OrganizationalUnitID      string `json:"organizationalUnitId,omitempty"`
	ObjectivePlanningEnabled  *bool  `json:"objectivePlanningEnabled,omitempty"`
	WorkProductPlanningEnabled *bool  `json:"workProductPlanningEnabled,omitempty"`
	WorkProductEvaluationEnabled *bool `json:"workProductEvaluationEnabled,omitempty"`
	Remark                    string `json:"remark,omitempty"`
	CreatedBy                 string `json:"createdBy,omitempty"`
}

// AddPeriodObjectiveRequest is the payload for adding objectives to a review period.
// Mirrors .NET AddPeriodObjectiveVm.
type AddPeriodObjectiveRequest struct {
	ReviewPeriodID string   `json:"reviewPeriodId"`
	ObjectiveIds   []string `json:"objectiveIds"`
}

// CategoryDefinitionRequest is the payload for category definition operations.
// Mirrors .NET CategoryDefinitionRequestVm / CreateCategoryDefinitionRequestVm.
type CategoryDefinitionRequest struct {
	ID                        string  `json:"id,omitempty"`
	ReviewPeriodID            string  `json:"reviewPeriodId"`
	ObjectiveCategoryID       string  `json:"objectiveCategoryId,omitempty"`
	Name                      string  `json:"name"`
	Description               string  `json:"description,omitempty"`
	Weight                    float64 `json:"weight,omitempty"`
	Remark                    string  `json:"remark,omitempty"`
	CreatedBy                 string  `json:"createdBy,omitempty"`
}

// ReviewPeriodExtensionRequest is the payload for review period extension operations.
// Mirrors .NET ReviewPeriodExtensionRequestModel / CreateReviewPeriodExtensionRequestModel.
type ReviewPeriodExtensionRequest struct {
	ID             string `json:"id,omitempty"`
	ReviewPeriodID string `json:"reviewPeriodId"`
	ExtensionDate  string `json:"extensionDate"`
	Reason         string `json:"reason,omitempty"`
	Remark         string `json:"remark,omitempty"`
	CreatedBy      string `json:"createdBy,omitempty"`
}

// ReviewPeriod360ReviewRequest is the payload for 360 review operations.
type ReviewPeriod360ReviewRequest struct {
	ID             string `json:"id,omitempty"`
	ReviewPeriodID string `json:"reviewPeriodId"`
	ReviewerID     string `json:"reviewerId"`
	RevieweeID     string `json:"revieweeId"`
	Score          float64 `json:"score,omitempty"`
	Comment        string `json:"comment,omitempty"`
	CreatedBy      string `json:"createdBy,omitempty"`
}

// IndividualPlannedObjectiveRequest is the payload for individual planned objective
// operations. Mirrors .NET ReviewPeriodIndividualPlannedObjectiveRequestModel /
// AddReviewPeriodIndividualPlannedObjectiveRequestModel.
type IndividualPlannedObjectiveRequest struct {
	ID                  string  `json:"id,omitempty"`
	ReviewPeriodID      string  `json:"reviewPeriodId"`
	StaffID             string  `json:"staffId"`
	ObjectiveID         string  `json:"objectiveId,omitempty"`
	Title               string  `json:"title,omitempty"`
	Description         string  `json:"description,omitempty"`
	Weight              float64 `json:"weight,omitempty"`
	TargetDate          string  `json:"targetDate,omitempty"`
	KeyPerformanceIndicator string `json:"keyPerformanceIndicator,omitempty"`
	Remark              string  `json:"remark,omitempty"`
	CreatedBy           string  `json:"createdBy,omitempty"`
}

// PeriodObjectiveEvaluationRequest is the payload for objective evaluation operations.
type PeriodObjectiveEvaluationRequest struct {
	ID             string  `json:"id,omitempty"`
	ReviewPeriodID string  `json:"reviewPeriodId"`
	ObjectiveID    string  `json:"objectiveId"`
	StaffID        string  `json:"staffId,omitempty"`
	DepartmentID   string  `json:"departmentId,omitempty"`
	Score          float64 `json:"score,omitempty"`
	Comment        string  `json:"comment,omitempty"`
	CreatedBy      string  `json:"createdBy,omitempty"`
}

// ===========================================================================
// REVIEW PERIOD LIFECYCLE
// ===========================================================================

// SaveDraftReviewPeriod handles POST /api/v1/review-periods/draft
// Mirrors .NET PerformanceMgtController.SaveDraftReviewPeriod.
func (h *ReviewPeriodHandler) SaveDraftReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.CreateNewReviewPeriodVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.SaveDraftReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftReviewPeriod").Msg("Failed to save draft review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// AddReviewPeriod handles POST /api/v1/review-periods
// Mirrors .NET PerformanceMgtController.AddReviewPeriod.
func (h *ReviewPeriodHandler) AddReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.CreateNewReviewPeriodVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.AddReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddReviewPeriod").Msg("Failed to add review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// SubmitDraftReviewPeriod handles POST /api/v1/review-periods/submit-draft
// Mirrors .NET PerformanceMgtController.SubmitDraftReviewPeriod.
func (h *ReviewPeriodHandler) SubmitDraftReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.SubmitDraftReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftReviewPeriod").Msg("Failed to submit draft review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ApproveReviewPeriod handles POST /api/v1/review-periods/approve
// Mirrors .NET PerformanceMgtController.ApproveReviewPeriod.
func (h *ReviewPeriodHandler) ApproveReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.ApproveReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveReviewPeriod").Msg("Failed to approve review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// RejectReviewPeriod handles POST /api/v1/review-periods/reject
// Mirrors .NET PerformanceMgtController.RejectReviewPeriod.
func (h *ReviewPeriodHandler) RejectReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.RejectReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectReviewPeriod").Msg("Failed to reject review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ReturnReviewPeriod handles POST /api/v1/review-periods/return
// Mirrors .NET PerformanceMgtController.ReturnReviewPeriod.
func (h *ReviewPeriodHandler) ReturnReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.ReturnReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnReviewPeriod").Msg("Failed to return review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ReSubmitReviewPeriod handles POST /api/v1/review-periods/resubmit
// Mirrors .NET PerformanceMgtController.ReSubmitReviewPeriod.
func (h *ReviewPeriodHandler) ReSubmitReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.ReSubmitReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitReviewPeriod").Msg("Failed to resubmit review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// UpdateReviewPeriod handles PUT /api/v1/review-periods
// Mirrors .NET PerformanceMgtController.UpdateReviewPeriod.
func (h *ReviewPeriodHandler) UpdateReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.UpdateReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateReviewPeriod").Msg("Failed to update review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// CancelReviewPeriod handles POST /api/v1/review-periods/cancel
// Mirrors .NET PerformanceMgtController.CancelReviewPeriod.
func (h *ReviewPeriodHandler) CancelReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.CancelReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelReviewPeriod").Msg("Failed to cancel review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// CloseReviewPeriod handles POST /api/v1/review-periods/close
// Mirrors .NET PerformanceMgtController.CloseReviewPeriod.
func (h *ReviewPeriodHandler) CloseReviewPeriod(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.CloseReviewPeriod(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CloseReviewPeriod").Msg("Failed to close review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// REVIEW PERIOD TOGGLES
// ===========================================================================

// EnableObjectivePlanning handles POST /api/v1/review-periods/enable-objective-planning
// Mirrors .NET PerformanceMgtController.EnableObjectivePlanning.
func (h *ReviewPeriodHandler) EnableObjectivePlanning(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.EnableObjectivePlanning(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "EnableObjectivePlanning").Msg("Failed to enable objective planning")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// DisableObjectivePlanning handles POST /api/v1/review-periods/disable-objective-planning
// Mirrors .NET PerformanceMgtController.DisableObjectivePlanning.
func (h *ReviewPeriodHandler) DisableObjectivePlanning(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.DisableObjectivePlanning(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "DisableObjectivePlanning").Msg("Failed to disable objective planning")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// EnableWorkProductPlanning handles POST /api/v1/review-periods/enable-work-product-planning
// Mirrors .NET PerformanceMgtController.EnableWorkProductPlanning.
func (h *ReviewPeriodHandler) EnableWorkProductPlanning(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.EnableWorkProductPlanning(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "EnableWorkProductPlanning").Msg("Failed to enable work product planning")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// DisableWorkProductPlanning handles POST /api/v1/review-periods/disable-work-product-planning
// Mirrors .NET PerformanceMgtController.DisableWorkProductPlanning.
func (h *ReviewPeriodHandler) DisableWorkProductPlanning(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.DisableWorkProductPlanning(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "DisableWorkProductPlanning").Msg("Failed to disable work product planning")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// EnableWorkProductEvaluation handles POST /api/v1/review-periods/enable-work-product-evaluation
// Mirrors .NET PerformanceMgtController.EnableWorkProductEvaluation.
func (h *ReviewPeriodHandler) EnableWorkProductEvaluation(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.EnableWorkProductEvaluation(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "EnableWorkProductEvaluation").Msg("Failed to enable work product evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// DisableWorkProductEvaluation handles POST /api/v1/review-periods/disable-work-product-evaluation
// Mirrors .NET PerformanceMgtController.DisableWorkProductEvaluation.
func (h *ReviewPeriodHandler) DisableWorkProductEvaluation(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.DisableWorkProductEvaluation(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "DisableWorkProductEvaluation").Msg("Failed to disable work product evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// REVIEW PERIOD RETRIEVAL
// ===========================================================================

// GetActiveReviewPeriod handles GET /api/v1/review-periods/active
// Mirrors .NET PerformanceMgtController.GetActiveReviewPeriod.
func (h *ReviewPeriodHandler) GetActiveReviewPeriod(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.ReviewPeriod.GetActiveReviewPeriod(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetActiveReviewPeriod").Msg("Failed to get active review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetStaffActiveReviewPeriod handles GET /api/v1/review-periods/active/staff?staffId={id}
// Mirrors .NET PerformanceMgtController.GetStaffActiveReviewPeriod.
func (h *ReviewPeriodHandler) GetStaffActiveReviewPeriod(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staffId")
	if staffID == "" {
		response.Error(w, http.StatusBadRequest, "staffId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetStaffActiveReviewPeriod(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffActiveReviewPeriod").Str("staffId", staffID).Msg("Failed to get staff active review period")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetReviewPeriodDetails handles GET /api/v1/review-periods/{reviewPeriodId}
// Mirrors .NET PerformanceMgtController.GetReviewPeriodDetails.
func (h *ReviewPeriodHandler) GetReviewPeriodDetails(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := r.PathValue("reviewPeriodId")
	if reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetReviewPeriodDetails(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetReviewPeriodDetails").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get review period details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// REVIEW PERIOD OBJECTIVES
// ===========================================================================

// SaveDraftReviewPeriodObjective handles POST /api/v1/review-periods/objectives/draft
// Mirrors .NET PerformanceMgtController.SaveDraftReviewPeriodObjective.
// Iterates over ObjectiveIds and creates a draft for each.
func (h *ReviewPeriodHandler) SaveDraftReviewPeriodObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.SaveDraftPeriodObjectiveVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if vm.ReviewPeriodID == "" || len(vm.ObjectiveIDs) == 0 {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId and at least one objectiveId are required")
		return
	}

	result, err := h.svc.ReviewPeriod.SaveDraftReviewPeriodObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftReviewPeriodObjective").Msg("Failed to save draft review period objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// AddReviewPeriodObjective handles POST /api/v1/review-periods/objectives
// Mirrors .NET PerformanceMgtController.AddReviewPeriodObjective.
// Iterates over ObjectiveIds and adds each to the period.
func (h *ReviewPeriodHandler) AddReviewPeriodObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.AddPeriodObjectiveVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if vm.ReviewPeriodID == "" || len(vm.ObjectiveIDs) == 0 {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId and at least one objectiveId are required")
		return
	}

	result, err := h.svc.ReviewPeriod.AddReviewPeriodObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddReviewPeriodObjective").Msg("Failed to add review period objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// SubmitDraftReviewPeriodObjective handles POST /api/v1/review-periods/objectives/submit-draft
// Mirrors .NET PerformanceMgtController.SubmitDraftReviewPeriodObjective.
func (h *ReviewPeriodHandler) SubmitDraftReviewPeriodObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.PeriodObjectiveRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if vm.ReviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.SubmitDraftReviewPeriodObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftReviewPeriodObjective").Msg("Failed to submit draft review period objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// CancelReviewPeriodObjective handles POST /api/v1/review-periods/objectives/cancel
// Mirrors .NET PerformanceMgtController.CancelReviewPeriodObjective.
func (h *ReviewPeriodHandler) CancelReviewPeriodObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.PeriodObjectiveRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if vm.ReviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.CancelReviewPeriodObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelReviewPeriodObjective").Msg("Failed to cancel review period objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetReviewPeriodObjectives handles GET /api/v1/review-periods/{reviewPeriodId}/objectives
// Mirrors .NET PerformanceMgtController.GetReviewPeriodObjectivesAsync.
func (h *ReviewPeriodHandler) GetReviewPeriodObjectives(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := r.PathValue("reviewPeriodId")
	if reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetReviewPeriodObjectives(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetReviewPeriodObjectives").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get review period objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// CATEGORY DEFINITIONS
// ===========================================================================

// SaveDraftReviewPeriodObjectiveCategoryDefinition handles POST /api/v1/review-periods/category-definitions/draft
// Mirrors .NET PerformanceMgtController.SaveDraftReviewPeriodObjectiveCategoryDefinition.
func (h *ReviewPeriodHandler) SaveDraftReviewPeriodObjectiveCategoryDefinition(w http.ResponseWriter, r *http.Request) {
	var vm performance.CategoryDefinitionRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.SaveDraftCategoryDefinition(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftReviewPeriodObjectiveCategoryDefinition").Msg("Failed to save draft category definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// AddReviewPeriodObjectiveCategoryDefinition handles POST /api/v1/review-periods/category-definitions
// Mirrors .NET PerformanceMgtController.AddObjectiveCategoryDefinition.
func (h *ReviewPeriodHandler) AddReviewPeriodObjectiveCategoryDefinition(w http.ResponseWriter, r *http.Request) {
	var vm performance.CategoryDefinitionRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.AddCategoryDefinition(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddReviewPeriodObjectiveCategoryDefinition").Msg("Failed to add category definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// SubmitDraftReviewPeriodObjectiveCategoryDefinition handles POST /api/v1/review-periods/category-definitions/submit-draft
// Mirrors .NET PerformanceMgtController.SubmitDraftObjectiveCategoryDefinition.
func (h *ReviewPeriodHandler) SubmitDraftReviewPeriodObjectiveCategoryDefinition(w http.ResponseWriter, r *http.Request) {
	var vm performance.CategoryDefinitionRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.SubmitDraftCategoryDefinition(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftReviewPeriodObjectiveCategoryDefinition").Msg("Failed to submit draft category definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ApproveReviewPeriodObjectiveCategoryDefinition handles POST /api/v1/review-periods/category-definitions/approve
// Mirrors .NET PerformanceMgtController.ApproveObjectiveCategoryDefinition.
func (h *ReviewPeriodHandler) ApproveReviewPeriodObjectiveCategoryDefinition(w http.ResponseWriter, r *http.Request) {
	var vm performance.CategoryDefinitionRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.ApproveCategoryDefinition(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveReviewPeriodObjectiveCategoryDefinition").Msg("Failed to approve category definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// RejectReviewPeriodObjectiveCategoryDefinition handles POST /api/v1/review-periods/category-definitions/reject
// Mirrors .NET PerformanceMgtController.CancelObjectiveCategoryDefinition (reject operation).
func (h *ReviewPeriodHandler) RejectReviewPeriodObjectiveCategoryDefinition(w http.ResponseWriter, r *http.Request) {
	var vm performance.CategoryDefinitionRequestVm
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.RejectCategoryDefinition(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectReviewPeriodObjectiveCategoryDefinition").Msg("Failed to reject category definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// REVIEW PERIOD EXTENSIONS
// ===========================================================================

// AddReviewPeriodExtension handles POST /api/v1/review-periods/extensions
// Mirrors .NET PerformanceMgtController.AddReviewPeriodExtension.
func (h *ReviewPeriodHandler) AddReviewPeriodExtension(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodExtensionRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.AddReviewPeriodExtension(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddReviewPeriodExtension").Msg("Failed to add review period extension")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetReviewPeriodExtensions handles GET /api/v1/review-periods/{reviewPeriodId}/extensions
// Mirrors .NET PerformanceMgtController.GetAllReviewPeriodExtensions.
func (h *ReviewPeriodHandler) GetReviewPeriodExtensions(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := r.PathValue("reviewPeriodId")
	if reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetReviewPeriodExtensions(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetReviewPeriodExtensions").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get review period extensions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// 360 REVIEWS
// ===========================================================================

// AddReviewPeriod360Review handles POST /api/v1/review-periods/360-reviews
func (h *ReviewPeriodHandler) AddReviewPeriod360Review(w http.ResponseWriter, r *http.Request) {
	var vm performance.CreateReviewPeriod360ReviewRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.AddReviewPeriod360Review(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddReviewPeriod360Review").Msg("Failed to add 360 review")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetReviewPeriod360Reviews handles GET /api/v1/review-periods/{reviewPeriodId}/360-reviews
func (h *ReviewPeriodHandler) GetReviewPeriod360Reviews(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := r.PathValue("reviewPeriodId")
	if reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetReviewPeriod360Reviews(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetReviewPeriod360Reviews").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get 360 reviews")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// INDIVIDUAL PLANNED OBJECTIVES
// ===========================================================================

// SaveDraftIndividualPlannedObjective handles POST /api/v1/review-periods/individual-planned-objectives/draft
// Mirrors .NET PerformanceMgtController.SaveDraftReviewPeriodOperationalObjective.
func (h *ReviewPeriodHandler) SaveDraftIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.SaveDraftIndividualPlannedObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftIndividualPlannedObjective").Msg("Failed to save draft individual planned objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// AddIndividualPlannedObjective handles POST /api/v1/review-periods/individual-planned-objectives
// Mirrors .NET PerformanceMgtController.AddReviewPeriodOperationalObjective.
func (h *ReviewPeriodHandler) AddIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.AddIndividualPlannedObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddIndividualPlannedObjective").Msg("Failed to add individual planned objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// SubmitDraftIndividualPlannedObjective handles POST /api/v1/review-periods/individual-planned-objectives/submit-draft
// Mirrors .NET PerformanceMgtController.SubmitDraftReviewPeriodOperationalObjective.
func (h *ReviewPeriodHandler) SubmitDraftIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.SubmitDraftIndividualPlannedObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftIndividualPlannedObjective").Msg("Failed to submit draft individual planned objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ApproveIndividualPlannedObjective handles POST /api/v1/review-periods/individual-planned-objectives/approve
// Mirrors .NET PerformanceMgtController.ApproveReviewPeriodOperationalObjective.
func (h *ReviewPeriodHandler) ApproveIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.ApproveIndividualPlannedObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveIndividualPlannedObjective").Msg("Failed to approve individual planned objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// RejectIndividualPlannedObjective handles POST /api/v1/review-periods/individual-planned-objectives/reject
// Mirrors .NET PerformanceMgtController.RejectReviewPeriodOperationalObjective.
func (h *ReviewPeriodHandler) RejectIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.RejectIndividualPlannedObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectIndividualPlannedObjective").Msg("Failed to reject individual planned objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ReturnIndividualPlannedObjective handles POST /api/v1/review-periods/individual-planned-objectives/return
// Mirrors .NET PerformanceMgtController.ReturnReviewPeriodOperationalObjective.
func (h *ReviewPeriodHandler) ReturnIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.ReturnIndividualPlannedObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnIndividualPlannedObjective").Msg("Failed to return individual planned objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// CancelIndividualPlannedObjective handles POST /api/v1/review-periods/individual-planned-objectives/cancel
// Mirrors .NET PerformanceMgtController.CancelReviewPeriodOperationalObjective.
func (h *ReviewPeriodHandler) CancelIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var vm performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.CancelIndividualPlannedObjective(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelIndividualPlannedObjective").Msg("Failed to cancel individual planned objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetStaffIndividualPlannedObjectives handles GET /api/v1/review-periods/individual-planned-objectives?staffId={id}&reviewPeriodId={id}
// Mirrors .NET PerformanceMgtController.GetReviewPeriodIndividualObjectivesAsync.
func (h *ReviewPeriodHandler) GetStaffIndividualPlannedObjectives(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staffId")
	reviewPeriodID := r.URL.Query().Get("reviewPeriodId")

	if staffID == "" || reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "staffId and reviewPeriodId are required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetStaffIndividualPlannedObjectives(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffIndividualPlannedObjectives").Str("staffId", staffID).Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get staff individual planned objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// PERIOD OBJECTIVE EVALUATIONS
// ===========================================================================

// CreatePeriodObjectiveEvaluation handles POST /api/v1/review-periods/objective-evaluations
func (h *ReviewPeriodHandler) CreatePeriodObjectiveEvaluation(w http.ResponseWriter, r *http.Request) {
	var vm performance.AddPeriodObjectiveEvaluationRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.CreatePeriodObjectiveEvaluation(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreatePeriodObjectiveEvaluation").Msg("Failed to create period objective evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// CreatePeriodObjectiveDepartmentEvaluation handles POST /api/v1/review-periods/objective-evaluations/department
func (h *ReviewPeriodHandler) CreatePeriodObjectiveDepartmentEvaluation(w http.ResponseWriter, r *http.Request) {
	var vm performance.AddPeriodObjectiveDepartmentEvaluationRequestModel
	if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.ReviewPeriod.CreatePeriodObjectiveDepartmentEvaluation(r.Context(), &vm)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreatePeriodObjectiveDepartmentEvaluation").Msg("Failed to create department objective evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetPeriodObjectiveEvaluations handles GET /api/v1/review-periods/{reviewPeriodId}/objective-evaluations
func (h *ReviewPeriodHandler) GetPeriodObjectiveEvaluations(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := r.PathValue("reviewPeriodId")
	if reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetPeriodObjectiveEvaluations(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPeriodObjectiveEvaluations").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get period objective evaluations")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// GetPeriodObjectiveDepartmentEvaluations handles GET /api/v1/review-periods/{reviewPeriodId}/objective-evaluations/department
func (h *ReviewPeriodHandler) GetPeriodObjectiveDepartmentEvaluations(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := r.PathValue("reviewPeriodId")
	if reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "reviewPeriodId is required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetPeriodObjectiveDepartmentEvaluations(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPeriodObjectiveDepartmentEvaluations").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get department objective evaluations")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// ===========================================================================
// PERIOD SCORES
// ===========================================================================

// GetStaffPeriodScore handles GET /api/v1/review-periods/scores?staffId={id}&reviewPeriodId={id}
func (h *ReviewPeriodHandler) GetStaffPeriodScore(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staffId")
	reviewPeriodID := r.URL.Query().Get("reviewPeriodId")

	if staffID == "" || reviewPeriodID == "" {
		response.Error(w, http.StatusBadRequest, "staffId and reviewPeriodId are required")
		return
	}

	result, err := h.svc.ReviewPeriod.GetStaffPeriodScore(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffPeriodScore").Str("staffId", staffID).Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get staff period score")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}
