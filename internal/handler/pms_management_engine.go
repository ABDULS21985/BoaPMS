package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/middleware"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// PmsEngineHandler converts .NET PmsManagementEngineController (partial class
// of PerformanceMgtController) which handles project, committee, work product,
// evaluation, feedback, scoring, and individual-objective-planning endpoints.
// ---------------------------------------------------------------------------

// PmsEngineHandler handles PMS management engine HTTP endpoints.
type PmsEngineHandler struct {
	svc *service.Container
	log zerolog.Logger
}

// NewPmsEngineHandler creates a new PMS engine handler.
func NewPmsEngineHandler(svc *service.Container, log zerolog.Logger) *PmsEngineHandler {
	return &PmsEngineHandler{svc: svc, log: log}
}

// ============================= Request DTOs ================================

// --- Project ---

// CreateProjectRequest mirrors .NET CreateProjectRequestModel.
type CreateProjectRequest struct {
	ProjectName           string `json:"projectName"`
	Description           string `json:"description"`
	ProjectManagerID      string `json:"projectManagerId"`
	EnterpriseObjectiveID string `json:"enterpriseObjectiveId"`
	ReviewPeriodID        string `json:"reviewPeriodId"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	StaffID               string `json:"staffId"`
}

// ProjectActionRequest mirrors .NET ProjectRequestModel for lifecycle ops.
type ProjectActionRequest struct {
	ProjectID             string `json:"projectId"`
	ProjectName           string `json:"projectName,omitempty"`
	Description           string `json:"description,omitempty"`
	ProjectManagerID      string `json:"projectManagerId,omitempty"`
	EnterpriseObjectiveID string `json:"enterpriseObjectiveId,omitempty"`
	ReviewPeriodID        string `json:"reviewPeriodId,omitempty"`
	StartDate             string `json:"startDate,omitempty"`
	EndDate               string `json:"endDate,omitempty"`
	StaffID               string `json:"staffId,omitempty"`
	Comment               string `json:"comment,omitempty"`
}

// AddProjectObjectiveRequest mirrors .NET AddProjectObjectiveRequestModel.
type AddProjectObjectiveRequest struct {
	ProjectID   string  `json:"projectId"`
	ObjectiveID string  `json:"objectiveId"`
	Weight      float64 `json:"weight"`
	StaffID     string  `json:"staffId"`
}

// AddProjectMemberRequest mirrors .NET AddProjectMemberRequestModel.
type AddProjectMemberRequest struct {
	ProjectID string `json:"projectId"`
	StaffID   string `json:"staffId"`
	RoleID    string `json:"roleId,omitempty"`
}

// ProjectMemberActionRequest mirrors .NET ProjectMemberRequestModel.
type ProjectMemberActionRequest struct {
	ProjectMemberID string `json:"projectMemberId"`
	ProjectID       string `json:"projectId"`
	StaffID         string `json:"staffId"`
	RoleID          string `json:"roleId,omitempty"`
	Comment         string `json:"comment,omitempty"`
}

// --- Committee ---

// CreateCommitteeRequest mirrors .NET CreateCommitteeRequestModel.
type CreateCommitteeRequest struct {
	CommitteeName         string `json:"committeeName"`
	Description           string `json:"description"`
	ChairpersonID         string `json:"chairpersonId"`
	EnterpriseObjectiveID string `json:"enterpriseObjectiveId"`
	ReviewPeriodID        string `json:"reviewPeriodId"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	StaffID               string `json:"staffId"`
}

// CommitteeActionRequest mirrors .NET CommitteeRequestModel for lifecycle ops.
type CommitteeActionRequest struct {
	CommitteeID           string `json:"committeeId"`
	CommitteeName         string `json:"committeeName,omitempty"`
	Description           string `json:"description,omitempty"`
	ChairpersonID         string `json:"chairpersonId,omitempty"`
	EnterpriseObjectiveID string `json:"enterpriseObjectiveId,omitempty"`
	ReviewPeriodID        string `json:"reviewPeriodId,omitempty"`
	StartDate             string `json:"startDate,omitempty"`
	EndDate               string `json:"endDate,omitempty"`
	StaffID               string `json:"staffId,omitempty"`
	Comment               string `json:"comment,omitempty"`
}

// AddCommitteeMemberRequest mirrors .NET AddCommitteeMemberRequestModel.
type AddCommitteeMemberRequest struct {
	CommitteeID string `json:"committeeId"`
	StaffID     string `json:"staffId"`
	RoleID      string `json:"roleId,omitempty"`
}

// AddCommitteeObjectiveRequest mirrors .NET AddCommitteeObjectiveRequestModel.
type AddCommitteeObjectiveRequest struct {
	CommitteeID string  `json:"committeeId"`
	ObjectiveID string  `json:"objectiveId"`
	Weight      float64 `json:"weight"`
	StaffID     string  `json:"staffId"`
}

// --- Work Product ---

// CreateWorkProductRequest mirrors .NET CreateWorkProductRequestModel.
type CreateWorkProductRequest struct {
	WorkProductName         string  `json:"workProductName"`
	Description             string  `json:"description"`
	WorkProductDefinitionID string  `json:"workProductDefinitionId"`
	EnterpriseObjectiveID   string  `json:"enterpriseObjectiveId"`
	ReviewPeriodID          string  `json:"reviewPeriodId"`
	StaffID                 string  `json:"staffId"`
	AssignedToStaffID       string  `json:"assignedToStaffId"`
	Weight                  float64 `json:"weight"`
	StartDate               string  `json:"startDate"`
	EndDate                 string  `json:"endDate"`
	ReferenceID             string  `json:"referenceId,omitempty"`
	AdhocAssignmentType     string  `json:"adhocAssignmentType,omitempty"`
}

// WorkProductActionRequest mirrors .NET WorkProductRequestModel for lifecycle ops.
type WorkProductActionRequest struct {
	WorkProductID           string  `json:"workProductId"`
	WorkProductName         string  `json:"workProductName,omitempty"`
	Description             string  `json:"description,omitempty"`
	WorkProductDefinitionID string  `json:"workProductDefinitionId,omitempty"`
	EnterpriseObjectiveID   string  `json:"enterpriseObjectiveId,omitempty"`
	ReviewPeriodID          string  `json:"reviewPeriodId,omitempty"`
	StaffID                 string  `json:"staffId,omitempty"`
	AssignedToStaffID       string  `json:"assignedToStaffId,omitempty"`
	Weight                  float64 `json:"weight,omitempty"`
	StartDate               string  `json:"startDate,omitempty"`
	EndDate                 string  `json:"endDate,omitempty"`
	Comment                 string  `json:"comment,omitempty"`
	ReferenceID             string  `json:"referenceId,omitempty"`
	AdhocAssignmentType     string  `json:"adhocAssignmentType,omitempty"`
}

// AssignWorkProductRequest is for assigning a work product to a staff member.
type AssignWorkProductRequest struct {
	WorkProductID     string `json:"workProductId"`
	AssignedToStaffID string `json:"assignedToStaffId"`
	StaffID           string `json:"staffId"`
}

// EvaluateWorkProductRequest mirrors .NET AddWorkProductEvaluationRequestModel.
type EvaluateWorkProductRequest struct {
	WorkProductID  string  `json:"workProductId"`
	EvaluatorID    string  `json:"evaluatorId"`
	Score          float64 `json:"score"`
	Comment        string  `json:"comment"`
	ReviewPeriodID string  `json:"reviewPeriodId"`
	EvaluationDate string  `json:"evaluationDate,omitempty"`
}

// WorkProductEvaluationActionRequest mirrors .NET WorkProductEvaluationRequestModel.
type WorkProductEvaluationActionRequest struct {
	WorkProductEvaluationID string  `json:"workProductEvaluationId"`
	WorkProductID           string  `json:"workProductId"`
	EvaluatorID             string  `json:"evaluatorId"`
	Score                   float64 `json:"score"`
	Comment                 string  `json:"comment"`
	ReviewPeriodID          string  `json:"reviewPeriodId"`
}

// --- Evaluation (Period Objective Evaluation) ---

// CreateEvaluationRequest mirrors .NET AddPeriodObjectiveEvaluationRequestModel.
type CreateEvaluationRequest struct {
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId"`
	ReviewPeriodID        string  `json:"reviewPeriodId"`
	PercentageScore       float64 `json:"percentageScore"`
	StaffID               string  `json:"staffId"`
	Comment               string  `json:"comment,omitempty"`
}

// EvaluationActionRequest mirrors .NET PeriodObjectiveEvaluationRequestModel.
type EvaluationActionRequest struct {
	EvaluationID          string  `json:"evaluationId"`
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId"`
	ReviewPeriodID        string  `json:"reviewPeriodId"`
	PercentageScore       float64 `json:"percentageScore"`
	StaffID               string  `json:"staffId"`
	Comment               string  `json:"comment,omitempty"`
}

// --- Feedback ---

// FeedbackActionRequest mirrors .NET FeedbackRequestModel.
type FeedbackActionRequest struct {
	RequestID  string `json:"requestId"`
	AssigneeID string `json:"assigneeId,omitempty"`
}

// TreatFeedbackRequest mirrors .NET TreatFeedbackRequestModel.
type TreatFeedbackRequest struct {
	RequestID     string `json:"requestId"`
	OperationType string `json:"operationType"`
	Comment       string `json:"comment"`
}

// RequestFeedbackRequest is for requesting feedback.
type RequestFeedbackRequest struct {
	StaffID        string `json:"staffId"`
	ReviewPeriodID string `json:"reviewPeriodId"`
	RequestType    string `json:"requestType"`
	ReferenceID    string `json:"referenceId"`
	Comment        string `json:"comment,omitempty"`
}

// ProcessFeedbackRequest is for processing feedback.
type ProcessFeedbackRequest struct {
	RequestID     string `json:"requestId"`
	OperationType string `json:"operationType"`
	Comment       string `json:"comment"`
	StaffID       string `json:"staffId"`
}

// --- Individual Planned Objective ---

// CreateIndividualObjectiveRequest mirrors .NET individual objective create DTO.
type CreateIndividualObjectiveRequest struct {
	ObjectiveName         string  `json:"objectiveName"`
	Description           string  `json:"description"`
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId"`
	ReviewPeriodID        string  `json:"reviewPeriodId"`
	StaffID               string  `json:"staffId"`
	Weight                float64 `json:"weight"`
	TargetDate            string  `json:"targetDate,omitempty"`
}

// IndividualObjectiveActionRequest mirrors .NET individual objective action DTO.
type IndividualObjectiveActionRequest struct {
	IndividualObjectiveID string  `json:"individualObjectiveId"`
	ObjectiveName         string  `json:"objectiveName,omitempty"`
	Description           string  `json:"description,omitempty"`
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId,omitempty"`
	ReviewPeriodID        string  `json:"reviewPeriodId,omitempty"`
	StaffID               string  `json:"staffId,omitempty"`
	Weight                float64 `json:"weight,omitempty"`
	Comment               string  `json:"comment,omitempty"`
	TargetDate            string  `json:"targetDate,omitempty"`
}

// ============================= Helpers =====================================

// decodeJSON reads the request body into dst; on failure it writes an error
// response and returns false.
func (h *PmsEngineHandler) decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return false
	}
	return true
}

// requiredQuery extracts a required query parameter; returns empty string and
// writes a 400 response if missing.
func (h *PmsEngineHandler) requiredQuery(w http.ResponseWriter, r *http.Request, name string) string {
	val := r.URL.Query().Get(name)
	if val == "" {
		response.Error(w, http.StatusBadRequest, name+" is required")
	}
	return val
}

// =================== PROJECT MANAGEMENT HANDLERS ===========================

// SaveDraftProject handles POST /api/v1/pms-engine/projects/draft
// Mirrors .NET SaveDraftProject -- creates a draft project via
// performanceManagementService.ProjectSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftProject").Msg("Failed to save draft project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddProject handles POST /api/v1/pms-engine/projects
// Mirrors .NET AddProject -- creates and commits a project via
// performanceManagementService.ProjectSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddProject").Msg("Failed to add project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// SubmitDraftProject handles POST /api/v1/pms-engine/projects/submit-draft
// Mirrors .NET SubmitDraftProject -- commits a previously saved draft via
// performanceManagementService.ProjectSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftProject").Msg("Failed to submit draft project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveProject handles POST /api/v1/pms-engine/projects/approve
// Mirrors .NET ApproveProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationApprove.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveProject").Msg("Failed to approve project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectProject handles POST /api/v1/pms-engine/projects/reject
// Mirrors .NET RejectProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReject.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectProject").Msg("Failed to reject project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReturnProject handles POST /api/v1/pms-engine/projects/return
// Mirrors .NET ReturnProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReturn.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnProject").Msg("Failed to return project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReSubmitProject handles POST /api/v1/pms-engine/projects/resubmit
// Mirrors .NET ReSubmitProject -- performanceManagementService.ProjectSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReSubmit.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitProject").Msg("Failed to resubmit project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// UpdateProject handles PUT /api/v1/pms-engine/projects
// Mirrors .NET UpdateProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationUpdate.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateProject").Msg("Failed to update project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelProject handles POST /api/v1/pms-engine/projects/cancel
// Mirrors .NET CancelProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCancel.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelProject").Msg("Failed to cancel project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjects handles GET /api/v1/pms-engine/projects
// Mirrors .NET GetProjects / GetProjectManagerProjects.
// When ?staffId= is provided it returns projects for that staff/manager.
func (h *PmsEngineHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staffId")

	if staffID != "" {
		result, err := h.svc.Performance.GetProjectsByManager(r.Context(), staffID)
		if err != nil {
			h.log.Error().Err(err).Str("action", "GetProjects").Str("staffId", staffID).Msg("Failed to get projects by manager")
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.OK(w, result)
		return
	}

	result, err := h.svc.Performance.GetProjects(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjects").Msg("Failed to get projects")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjectDetails handles GET /api/v1/pms-engine/projects/{projectId}
// Mirrors .NET GetProject -- retrieves project details by ID.
func (h *PmsEngineHandler) GetProjectDetails(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "projectId is required")
		return
	}

	result, err := h.svc.Performance.GetProject(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectDetails").Str("projectId", projectID).Msg("Failed to get project details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddProjectObjective handles POST /api/v1/pms-engine/projects/objectives
// Mirrors .NET AddProjectObjective -- performanceManagementService.ProjectObjectiveSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddProjectObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.ProjectObjectiveSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddProjectObjective").Msg("Failed to add project objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// AddProjectMember handles POST /api/v1/pms-engine/projects/members
// Mirrors .NET AddProjectMember -- performanceManagementService.ProjectMembersSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.ProjectMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddProjectMember").Msg("Failed to add project member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// GetProjectMembers handles GET /api/v1/pms-engine/projects/{projectId}/members
// Mirrors .NET GetProjectMembers -- performanceManagementService.GetProjectMembers(projectId).
func (h *PmsEngineHandler) GetProjectMembers(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "projectId is required")
		return
	}

	result, err := h.svc.Performance.GetProjectMembers(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectMembers").Str("projectId", projectID).Msg("Failed to get project members")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjectObjectives handles GET /api/v1/pms-engine/projects/{projectId}/objectives
// Mirrors .NET GetProjectObjectives -- performanceManagementService.GetProjectObjectives(projectId).
func (h *PmsEngineHandler) GetProjectObjectives(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "projectId is required")
		return
	}

	result, err := h.svc.Performance.GetProjectObjectives(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectObjectives").Str("projectId", projectID).Msg("Failed to get project objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CloseProject handles POST /api/v1/pms-engine/projects/close
// Mirrors .NET CloseProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Close).
func (h *PmsEngineHandler) CloseProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationClose.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CloseProject").Msg("Failed to close project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// PauseProject handles POST /api/v1/pms-engine/projects/pause
// Mirrors .NET PauseProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Pause).
func (h *PmsEngineHandler) PauseProject(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationPause.String()
	result, err := h.svc.Performance.ProjectSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "PauseProject").Msg("Failed to pause project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjectsByManager handles GET /api/v1/pms-engine/projects/by-manager?managerId={id}
// Mirrors .NET GetProjectManagerProjects -- dedicated endpoint for manager lookup.
func (h *PmsEngineHandler) GetProjectsByManager(w http.ResponseWriter, r *http.Request) {
	managerID := h.requiredQuery(w, r, "managerId")
	if managerID == "" {
		return
	}

	result, err := h.svc.Performance.GetProjectsByManager(r.Context(), managerID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectsByManager").Str("managerId", managerID).Msg("Failed to get projects by manager")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjectsAssigned handles GET /api/v1/pms-engine/projects/assigned?staffId={id}
// Mirrors .NET GetProjectsAssigned -- returns projects assigned to a staff member.
func (h *PmsEngineHandler) GetProjectsAssigned(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetProjectsAssigned(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectsAssigned").Str("staffId", staffID).Msg("Failed to get assigned projects")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffProjects handles GET /api/v1/pms-engine/projects/staff?staffId={id}
// Mirrors .NET GetStaffProjects -- returns all projects for a staff member.
func (h *PmsEngineHandler) GetStaffProjects(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetStaffProjects(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffProjects").Str("staffId", staffID).Msg("Failed to get staff projects")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjectWorkProductStaffList handles GET /api/v1/pms-engine/projects/{projectId}/work-product-staff
// Mirrors .NET GetProjectWorkProductStaffList.
func (h *PmsEngineHandler) GetProjectWorkProductStaffList(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "projectId is required")
		return
	}

	result, err := h.svc.Performance.GetProjectWorkProductStaffList(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectWorkProductStaffList").Str("projectId", projectID).Msg("Failed to get project work product staff list")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveDraftProjectMember handles POST /api/v1/pms-engine/projects/members/draft
// Mirrors .NET SaveDraftProjectMember -- ProjectMembersSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftProjectMember(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.ProjectMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftProjectMember").Msg("Failed to save draft project member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SubmitDraftProjectMember handles POST /api/v1/pms-engine/projects/members/submit-draft
// Mirrors .NET SubmitDraftProjectMember -- ProjectMembersSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftProjectMember(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.ProjectMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftProjectMember").Msg("Failed to submit draft project member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AcceptProjectMember handles POST /api/v1/pms-engine/projects/members/accept
// Mirrors .NET AcceptProjectMember -- ProjectMembersSetup(request, OperationTypes.Accept).
func (h *PmsEngineHandler) AcceptProjectMember(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAccept.String()
	result, err := h.svc.Performance.ProjectMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AcceptProjectMember").Msg("Failed to accept project member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveProjectMember handles POST /api/v1/pms-engine/projects/members/approve
// Mirrors .NET ApproveProjectMember -- ProjectMembersSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveProjectMember(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationApprove.String()
	result, err := h.svc.Performance.ProjectMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveProjectMember").Msg("Failed to approve project member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelProjectMember handles POST /api/v1/pms-engine/projects/members/cancel
// Mirrors .NET CancelProjectMember -- ProjectMembersSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelProjectMember(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCancel.String()
	result, err := h.svc.Performance.ProjectMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelProjectMember").Msg("Failed to cancel project member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelProjectObjective handles POST /api/v1/pms-engine/projects/objectives/cancel
// Mirrors .NET CancelProjectObjective -- ProjectObjectiveSetup with cancel status.
func (h *PmsEngineHandler) CancelProjectObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationCancel.String()
	result, err := h.svc.Performance.ProjectObjectiveSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelProjectObjective").Msg("Failed to cancel project objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ChangeAdhocAssignmentLead handles POST /api/v1/pms-engine/projects/change-lead
// Mirrors .NET ChangeAdhocAssignmentLead -- performanceManagementService.ChangeProjectLead(request).
func (h *PmsEngineHandler) ChangeAdhocAssignmentLead(w http.ResponseWriter, r *http.Request) {
	var req performance.ChangeAdhocLeadRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.ChangeProjectLead(r.Context(), &req); err != nil {
		h.log.Error().Err(err).Str("action", "ChangeAdhocAssignmentLead").Msg("Failed to change project lead")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, map[string]string{"message": "Project lead changed successfully"})
}

// ValidateStaffEligibilityForAdhoc handles GET /api/v1/pms-engine/projects/validate-eligibility?staffId={id}&reviewPeriodId={id}
// Mirrors .NET ValidateStaffEligibilityForAdhoc.
func (h *PmsEngineHandler) ValidateStaffEligibilityForAdhoc(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}

	result, err := h.svc.Performance.ValidateStaffEligibilityForAdhoc(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ValidateStaffEligibilityForAdhoc").Str("staffId", staffID).Msg("Failed to validate staff eligibility")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== COMMITTEE MANAGEMENT HANDLERS =========================

// SaveDraftCommittee handles POST /api/v1/pms-engine/committees/draft
// Mirrors .NET SaveDraftCommittee -- performanceManagementService.CommitteeSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftCommittee").Msg("Failed to save draft committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddCommittee handles POST /api/v1/pms-engine/committees
// Mirrors .NET AddCommittee -- performanceManagementService.CommitteeSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddCommittee").Msg("Failed to add committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// SubmitDraftCommittee handles POST /api/v1/pms-engine/committees/submit-draft
// Mirrors .NET SubmitDraftCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftCommittee").Msg("Failed to submit draft committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveCommittee handles POST /api/v1/pms-engine/committees/approve
// Mirrors .NET ApproveCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationApprove.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveCommittee").Msg("Failed to approve committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectCommittee handles POST /api/v1/pms-engine/committees/reject
// Mirrors .NET RejectCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReject.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectCommittee").Msg("Failed to reject committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReturnCommittee handles POST /api/v1/pms-engine/committees/return
// Mirrors .NET ReturnCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReturn.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnCommittee").Msg("Failed to return committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReSubmitCommittee handles POST /api/v1/pms-engine/committees/resubmit
// Mirrors .NET ReSubmitCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReSubmit.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitCommittee").Msg("Failed to resubmit committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// UpdateCommittee handles PUT /api/v1/pms-engine/committees
// Mirrors .NET UpdateCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationUpdate.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateCommittee").Msg("Failed to update committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelCommittee handles POST /api/v1/pms-engine/committees/cancel
// Mirrors .NET CancelCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCancel.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelCommittee").Msg("Failed to cancel committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommittees handles GET /api/v1/pms-engine/committees
// Mirrors .NET GetCommittees / GetCommitteesForChairperson.
// When ?chairpersonId= is provided it filters by chairperson.
func (h *PmsEngineHandler) GetCommittees(w http.ResponseWriter, r *http.Request) {
	chairpersonID := r.URL.Query().Get("chairpersonId")

	if chairpersonID != "" {
		result, err := h.svc.Performance.GetCommitteesByChairperson(r.Context(), chairpersonID)
		if err != nil {
			h.log.Error().Err(err).Str("action", "GetCommittees").Str("chairpersonId", chairpersonID).Msg("Failed to get committees by chairperson")
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.OK(w, result)
		return
	}

	result, err := h.svc.Performance.GetCommittees(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommittees").Msg("Failed to get committees")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteeDetails handles GET /api/v1/pms-engine/committees/{committeeId}
// Mirrors .NET GetCommittee -- performanceManagementService.GetCommittee(committeeId).
func (h *PmsEngineHandler) GetCommitteeDetails(w http.ResponseWriter, r *http.Request) {
	committeeID := r.PathValue("committeeId")
	if committeeID == "" {
		response.Error(w, http.StatusBadRequest, "committeeId is required")
		return
	}

	result, err := h.svc.Performance.GetCommittee(r.Context(), committeeID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteeDetails").Str("committeeId", committeeID).Msg("Failed to get committee details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddCommitteeMember handles POST /api/v1/pms-engine/committees/members
// Mirrors .NET AddCommitteeMember -- performanceManagementService.CommitteeMembersSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddCommitteeMember(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.CommitteeMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddCommitteeMember").Msg("Failed to add committee member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// AddCommitteeObjective handles POST /api/v1/pms-engine/committees/objectives
// Mirrors .NET AddCommitteeObjective -- performanceManagementService.CommitteeObjectiveSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddCommitteeObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.CommitteeObjectiveSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddCommitteeObjective").Msg("Failed to add committee objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// CloseCommittee handles POST /api/v1/pms-engine/committees/close
// Mirrors .NET CloseCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Close).
func (h *PmsEngineHandler) CloseCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationClose.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CloseCommittee").Msg("Failed to close committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// PauseCommittee handles POST /api/v1/pms-engine/committees/pause
// Mirrors .NET PauseCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Pause).
func (h *PmsEngineHandler) PauseCommittee(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationPause.String()
	result, err := h.svc.Performance.CommitteeSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "PauseCommittee").Msg("Failed to pause committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteesByChairperson handles GET /api/v1/pms-engine/committees/by-chairperson?chairpersonId={id}
// Mirrors .NET GetCommitteesForChairperson -- dedicated endpoint for chairperson lookup.
func (h *PmsEngineHandler) GetCommitteesByChairperson(w http.ResponseWriter, r *http.Request) {
	chairpersonID := h.requiredQuery(w, r, "chairpersonId")
	if chairpersonID == "" {
		return
	}

	result, err := h.svc.Performance.GetCommitteesByChairperson(r.Context(), chairpersonID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteesByChairperson").Str("chairpersonId", chairpersonID).Msg("Failed to get committees by chairperson")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteeMembers handles GET /api/v1/pms-engine/committees/{committeeId}/members
// Mirrors .NET GetCommitteeMembers -- performanceManagementService.GetCommitteeMembers(committeeId).
func (h *PmsEngineHandler) GetCommitteeMembers(w http.ResponseWriter, r *http.Request) {
	committeeID := r.PathValue("committeeId")
	if committeeID == "" {
		response.Error(w, http.StatusBadRequest, "committeeId is required")
		return
	}

	result, err := h.svc.Performance.GetCommitteeMembers(r.Context(), committeeID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteeMembers").Str("committeeId", committeeID).Msg("Failed to get committee members")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteesAssigned handles GET /api/v1/pms-engine/committees/assigned?staffId={id}
// Mirrors .NET GetCommitteesAssigned -- returns committees assigned to a staff member.
func (h *PmsEngineHandler) GetCommitteesAssigned(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetCommitteesAssigned(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteesAssigned").Str("staffId", staffID).Msg("Failed to get assigned committees")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffCommittees handles GET /api/v1/pms-engine/committees/staff?staffId={id}
// Mirrors .NET GetStaffCommittees -- returns all committees for a staff member.
func (h *PmsEngineHandler) GetStaffCommittees(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetStaffCommittees(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffCommittees").Str("staffId", staffID).Msg("Failed to get staff committees")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteeWorkProductStaffList handles GET /api/v1/pms-engine/committees/{committeeId}/work-product-staff
// Mirrors .NET GetCommitteeWorkProductStaffList.
func (h *PmsEngineHandler) GetCommitteeWorkProductStaffList(w http.ResponseWriter, r *http.Request) {
	committeeID := r.PathValue("committeeId")
	if committeeID == "" {
		response.Error(w, http.StatusBadRequest, "committeeId is required")
		return
	}

	result, err := h.svc.Performance.GetCommitteeWorkProductStaffList(r.Context(), committeeID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteeWorkProductStaffList").Str("committeeId", committeeID).Msg("Failed to get committee work product staff list")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteeObjectives handles GET /api/v1/pms-engine/committees/{committeeId}/objectives
// Mirrors .NET GetCommitteeObjectives -- performanceManagementService.GetCommitteeObjectives(committeeId).
func (h *PmsEngineHandler) GetCommitteeObjectives(w http.ResponseWriter, r *http.Request) {
	committeeID := r.PathValue("committeeId")
	if committeeID == "" {
		response.Error(w, http.StatusBadRequest, "committeeId is required")
		return
	}

	result, err := h.svc.Performance.GetCommitteeObjectives(r.Context(), committeeID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteeObjectives").Str("committeeId", committeeID).Msg("Failed to get committee objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveDraftCommitteeMember handles POST /api/v1/pms-engine/committees/members/draft
// Mirrors .NET SaveDraftCommitteeMember -- CommitteeMembersSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftCommitteeMember(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.CommitteeMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftCommitteeMember").Msg("Failed to save draft committee member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SubmitDraftCommitteeMember handles POST /api/v1/pms-engine/committees/members/submit-draft
// Mirrors .NET SubmitDraftCommitteeMember -- CommitteeMembersSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftCommitteeMember(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.CommitteeMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftCommitteeMember").Msg("Failed to submit draft committee member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelCommitteeMember handles POST /api/v1/pms-engine/committees/members/cancel
// Mirrors .NET CancelCommitteeMember -- CommitteeMembersSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelCommitteeMember(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeMemberRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCancel.String()
	result, err := h.svc.Performance.CommitteeMembersSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelCommitteeMember").Msg("Failed to cancel committee member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelCommitteeObjective handles POST /api/v1/pms-engine/committees/objectives/cancel
// Mirrors .NET CancelCommitteeObjective -- CommitteeObjectiveSetup with cancel status.
func (h *PmsEngineHandler) CancelCommitteeObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationCancel.String()
	result, err := h.svc.Performance.CommitteeObjectiveSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelCommitteeObjective").Msg("Failed to cancel committee objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ChangeCommitteeChairperson handles POST /api/v1/pms-engine/committees/change-chairperson
// Mirrors .NET ChangeCommitteeChairperson -- performanceManagementService.ChangeCommitteeChairperson(request).
func (h *PmsEngineHandler) ChangeCommitteeChairperson(w http.ResponseWriter, r *http.Request) {
	var req performance.ChangeAdhocLeadRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.ChangeCommitteeChairperson(r.Context(), &req); err != nil {
		h.log.Error().Err(err).Str("action", "ChangeCommitteeChairperson").Msg("Failed to change committee chairperson")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, map[string]string{"message": "Committee chairperson changed successfully"})
}

// =================== WORK PRODUCT HANDLERS =================================

// SaveDraftWorkProduct handles POST /api/v1/pms-engine/work-products/draft
// Mirrors .NET SaveDraftWorkProduct -- performanceManagementService.WorkProductSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftWorkProduct").Msg("Failed to save draft work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddWorkProduct handles POST /api/v1/pms-engine/work-products
// Mirrors .NET AddWorkProduct -- performanceManagementService.WorkProductSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddWorkProduct").Msg("Failed to add work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// SubmitDraftWorkProduct handles POST /api/v1/pms-engine/work-products/submit-draft
// Mirrors .NET SubmitDraftWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftWorkProduct").Msg("Failed to submit draft work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveWorkProduct handles POST /api/v1/pms-engine/work-products/approve
// Mirrors .NET ApproveWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationApprove.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveWorkProduct").Msg("Failed to approve work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectWorkProduct handles POST /api/v1/pms-engine/work-products/reject
// Mirrors .NET RejectWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReject.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectWorkProduct").Msg("Failed to reject work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReturnWorkProduct handles POST /api/v1/pms-engine/work-products/return
// Mirrors .NET ReturnWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReturn.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnWorkProduct").Msg("Failed to return work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReSubmitWorkProduct handles POST /api/v1/pms-engine/work-products/resubmit
// Mirrors .NET ReSubmitWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReSubmit.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitWorkProduct").Msg("Failed to resubmit work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// UpdateWorkProduct handles PUT /api/v1/pms-engine/work-products
// Mirrors .NET UpdateWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationUpdate.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateWorkProduct").Msg("Failed to update work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelWorkProduct handles POST /api/v1/pms-engine/work-products/cancel
// Mirrors .NET CancelWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCancel.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelWorkProduct").Msg("Failed to cancel work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// PauseWorkProduct handles POST /api/v1/pms-engine/work-products/pause
// Mirrors .NET PauseWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Pause).
func (h *PmsEngineHandler) PauseWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationPause.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "PauseWorkProduct").Msg("Failed to pause work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ResumeWorkProduct handles POST /api/v1/pms-engine/work-products/resume
// Mirrors .NET ResumeWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Resume).
func (h *PmsEngineHandler) ResumeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationResume.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ResumeWorkProduct").Msg("Failed to resume work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffWorkProducts handles GET /api/v1/pms-engine/work-products?staffId={id}
// Mirrors .NET GetStaffWorkProducts -- performanceManagementService.GetStaffWorkProducts(staffId, reviewPeriodId).
func (h *PmsEngineHandler) GetStaffWorkProducts(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	reviewPeriodID := r.URL.Query().Get("reviewPeriodId")

	if reviewPeriodID != "" {
		result, err := h.svc.Performance.GetStaffWorkProducts(r.Context(), staffID, reviewPeriodID)
		if err != nil {
			h.log.Error().Err(err).Str("action", "GetStaffWorkProducts").Str("staffId", staffID).Msg("Failed to get staff work products")
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.OK(w, result)
		return
	}

	result, err := h.svc.Performance.GetAllStaffWorkProducts(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffWorkProducts").Str("staffId", staffID).Msg("Failed to get all staff work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetWorkProductDetails handles GET /api/v1/pms-engine/work-products/{workProductId}
// Mirrors .NET GetWorkProduct -- performanceManagementService.GetWorkProduct(workProductId).
func (h *PmsEngineHandler) GetWorkProductDetails(w http.ResponseWriter, r *http.Request) {
	workProductID := r.PathValue("workProductId")
	if workProductID == "" {
		response.Error(w, http.StatusBadRequest, "workProductId is required")
		return
	}

	result, err := h.svc.Performance.GetWorkProduct(r.Context(), workProductID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetWorkProductDetails").Str("workProductId", workProductID).Msg("Failed to get work product details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AssignWorkProduct handles POST /api/v1/pms-engine/work-products/assign
// Mirrors .NET ProjectAssignedWorkProductSetup / CommitteeAssignedWorkProductSetup.
func (h *PmsEngineHandler) AssignWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AssignWorkProduct").Msg("Failed to assign work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAssignedWorkProducts handles GET /api/v1/pms-engine/work-products/assigned?staffId={id}
// Mirrors .NET GetProjectsAssigned / GetCommitteesAssigned.
func (h *PmsEngineHandler) GetAssignedWorkProducts(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetProjectsAssigned(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAssignedWorkProducts").Str("staffId", staffID).Msg("Failed to get assigned work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// EvaluateWorkProduct handles POST /api/v1/pms-engine/work-products/evaluate
// Mirrors .NET AddWorkProductEvaluation -- performanceManagementService.WorkProductEvaluation(model, OperationTypes.Add).
func (h *PmsEngineHandler) EvaluateWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.WorkProductEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "EvaluateWorkProduct").Msg("Failed to evaluate work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CompleteWorkProduct handles POST /api/v1/pms-engine/work-products/complete
// Mirrors .NET CompleteWorkProduct -- WorkProductSetup(request, OperationTypes.Complete).
func (h *PmsEngineHandler) CompleteWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationComplete.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CompleteWorkProduct").Msg("Failed to complete work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SuspendWorkProduct handles POST /api/v1/pms-engine/work-products/suspend
// Mirrors .NET SuspendWorkProduct -- WorkProductSetup(request, OperationTypes.Suspend).
func (h *PmsEngineHandler) SuspendWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationSuspend.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SuspendWorkProduct").Msg("Failed to suspend work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReInstateWorkProduct handles POST /api/v1/pms-engine/work-products/reinstate
// Mirrors .NET ReInstateWorkProduct -- WorkProductSetup(request, OperationTypes.ReInstate).
func (h *PmsEngineHandler) ReInstateWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReInstate.String()
	result, err := h.svc.Performance.WorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReInstateWorkProduct").Msg("Failed to reinstate work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== PROJECT ASSIGNED WORK PRODUCT HANDLERS ================

// SaveDraftProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/draft
// Mirrors .NET SaveDraftProjectWorkProduct -- ProjectAssignedWorkProductSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftProjectWorkProduct").Msg("Failed to save draft project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project
// Mirrors .NET AddProjectWorkProduct -- ProjectAssignedWorkProductSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddProjectWorkProduct").Msg("Failed to add project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// SubmitDraftProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/submit-draft
// Mirrors .NET SubmitDraftProjectWorkProduct -- ProjectAssignedWorkProductSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftProjectWorkProduct").Msg("Failed to submit draft project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/approve
// Mirrors .NET ApproveProjectWorkProduct -- ProjectAssignedWorkProductSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationApprove.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveProjectWorkProduct").Msg("Failed to approve project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/reject
// Mirrors .NET RejectProjectWorkProduct -- ProjectAssignedWorkProductSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReject.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectProjectWorkProduct").Msg("Failed to reject project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReturnProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/return
// Mirrors .NET ReturnProjectWorkProduct -- ProjectAssignedWorkProductSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReturn.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnProjectWorkProduct").Msg("Failed to return project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReSubmitProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/resubmit
// Mirrors .NET ReSubmitProjectWorkProduct -- ProjectAssignedWorkProductSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReSubmit.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitProjectWorkProduct").Msg("Failed to resubmit project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/cancel
// Mirrors .NET CancelProjectWorkProduct -- ProjectAssignedWorkProductSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCancel.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelProjectWorkProduct").Msg("Failed to cancel project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CloseProjectWorkProduct handles POST /api/v1/pms-engine/work-products/project/close
// Mirrors .NET CloseProjectWorkProduct -- ProjectAssignedWorkProductSetup(request, OperationTypes.Close).
func (h *PmsEngineHandler) CloseProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.ProjectAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationClose.String()
	result, err := h.svc.Performance.ProjectAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CloseProjectWorkProduct").Msg("Failed to close project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== COMMITTEE ASSIGNED WORK PRODUCT HANDLERS ==============

// SaveDraftCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/draft
// Mirrors .NET SaveDraftCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftCommitteeWorkProduct").Msg("Failed to save draft committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee
// Mirrors .NET AddCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddCommitteeWorkProduct").Msg("Failed to add committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// SubmitDraftCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/submit-draft
// Mirrors .NET SubmitDraftCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftCommitteeWorkProduct").Msg("Failed to submit draft committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/approve
// Mirrors .NET ApproveCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationApprove.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveCommitteeWorkProduct").Msg("Failed to approve committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/reject
// Mirrors .NET RejectCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReject.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectCommitteeWorkProduct").Msg("Failed to reject committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReturnCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/return
// Mirrors .NET ReturnCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReturn.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnCommitteeWorkProduct").Msg("Failed to return committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReSubmitCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/resubmit
// Mirrors .NET ReSubmitCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReSubmit.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitCommitteeWorkProduct").Msg("Failed to resubmit committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/cancel
// Mirrors .NET CancelCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCancel.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelCommitteeWorkProduct").Msg("Failed to cancel committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CloseCommitteeWorkProduct handles POST /api/v1/pms-engine/work-products/committee/close
// Mirrors .NET CloseCommitteeWorkProduct -- CommitteeAssignedWorkProductSetup(request, OperationTypes.Close).
func (h *PmsEngineHandler) CloseCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req performance.CommitteeAssignedWorkProductRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationClose.String()
	result, err := h.svc.Performance.CommitteeAssignedWorkProductSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CloseCommitteeWorkProduct").Msg("Failed to close committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== WORK PRODUCT RETRIEVAL HANDLERS =======================

// GetProjectAssignedWorkProductDetails handles GET /api/v1/pms-engine/work-products/project/{id}
// Mirrors .NET GetProjectAssignedWorkProductDetails -- retrieves work product details by ID.
func (h *PmsEngineHandler) GetProjectAssignedWorkProductDetails(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	result, err := h.svc.Performance.GetWorkProduct(r.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectAssignedWorkProductDetails").Str("id", id).Msg("Failed to get project assigned work product details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjectAssignedWorkProducts handles GET /api/v1/pms-engine/work-products/project?projectId={id}
// Mirrors .NET GetProjectAssignedWorkProducts -- retrieves assigned work products for a project.
func (h *PmsEngineHandler) GetProjectAssignedWorkProducts(w http.ResponseWriter, r *http.Request) {
	projectID := h.requiredQuery(w, r, "projectId")
	if projectID == "" {
		return
	}

	result, err := h.svc.Performance.GetProjectAssignedWorkProducts(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectAssignedWorkProducts").Str("projectId", projectID).Msg("Failed to get project assigned work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetProjectWorkProduct handles GET /api/v1/pms-engine/work-products/project/single?projectId={id}&workProductId={id}
// Mirrors .NET GetProjectWorkProduct -- retrieves a single work product by ID.
func (h *PmsEngineHandler) GetProjectWorkProduct(w http.ResponseWriter, r *http.Request) {
	workProductID := h.requiredQuery(w, r, "workProductId")
	if workProductID == "" {
		return
	}

	result, err := h.svc.Performance.GetWorkProduct(r.Context(), workProductID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectWorkProduct").Str("workProductId", workProductID).Msg("Failed to get project work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllProjectWorkProducts handles GET /api/v1/pms-engine/work-products/project/all?projectId={id}
// Mirrors .NET GetAllProjectWorkProducts -- retrieves all work products for a project.
func (h *PmsEngineHandler) GetAllProjectWorkProducts(w http.ResponseWriter, r *http.Request) {
	projectID := h.requiredQuery(w, r, "projectId")
	if projectID == "" {
		return
	}

	result, err := h.svc.Performance.GetProjectWorkProducts(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllProjectWorkProducts").Str("projectId", projectID).Msg("Failed to get all project work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffProjectWorkProducts handles GET /api/v1/pms-engine/work-products/project/staff?projectId={id}&staffId={id}
// Mirrors .NET GetStaffProjectWorkProducts -- retrieves work products for a project filtered by staff.
func (h *PmsEngineHandler) GetStaffProjectWorkProducts(w http.ResponseWriter, r *http.Request) {
	projectID := h.requiredQuery(w, r, "projectId")
	if projectID == "" {
		return
	}

	result, err := h.svc.Performance.GetProjectWorkProducts(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffProjectWorkProducts").Str("projectId", projectID).Msg("Failed to get staff project work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteeAssignedWorkProductDetails handles GET /api/v1/pms-engine/work-products/committee/{id}
// Mirrors .NET GetCommitteeAssignedWorkProductDetails -- retrieves work product details by ID.
func (h *PmsEngineHandler) GetCommitteeAssignedWorkProductDetails(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	result, err := h.svc.Performance.GetWorkProduct(r.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteeAssignedWorkProductDetails").Str("id", id).Msg("Failed to get committee assigned work product details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteeAssignedWorkProducts handles GET /api/v1/pms-engine/work-products/committee?committeeId={id}
// Mirrors .NET GetCommitteeAssignedWorkProducts -- retrieves assigned work products for a committee.
func (h *PmsEngineHandler) GetCommitteeAssignedWorkProducts(w http.ResponseWriter, r *http.Request) {
	committeeID := h.requiredQuery(w, r, "committeeId")
	if committeeID == "" {
		return
	}

	result, err := h.svc.Performance.GetCommitteeAssignedWorkProducts(r.Context(), committeeID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteeAssignedWorkProducts").Str("committeeId", committeeID).Msg("Failed to get committee assigned work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCommitteeWorkProduct handles GET /api/v1/pms-engine/work-products/committee/single?committeeId={id}&workProductId={id}
// Mirrors .NET GetCommitteeWorkProduct -- retrieves a single work product by ID.
func (h *PmsEngineHandler) GetCommitteeWorkProduct(w http.ResponseWriter, r *http.Request) {
	workProductID := h.requiredQuery(w, r, "workProductId")
	if workProductID == "" {
		return
	}

	result, err := h.svc.Performance.GetWorkProduct(r.Context(), workProductID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCommitteeWorkProduct").Str("workProductId", workProductID).Msg("Failed to get committee work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllCommitteeWorkProducts handles GET /api/v1/pms-engine/work-products/committee/all?committeeId={id}
// Mirrors .NET GetAllCommitteeWorkProducts -- retrieves all work products for a committee.
func (h *PmsEngineHandler) GetAllCommitteeWorkProducts(w http.ResponseWriter, r *http.Request) {
	committeeID := h.requiredQuery(w, r, "committeeId")
	if committeeID == "" {
		return
	}

	result, err := h.svc.Performance.GetCommitteeWorkProducts(r.Context(), committeeID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllCommitteeWorkProducts").Str("committeeId", committeeID).Msg("Failed to get all committee work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffCommitteeWorkProducts handles GET /api/v1/pms-engine/work-products/committee/staff?committeeId={id}&staffId={id}
// Mirrors .NET GetStaffCommitteeWorkProducts -- retrieves work products for a committee filtered by staff.
func (h *PmsEngineHandler) GetStaffCommitteeWorkProducts(w http.ResponseWriter, r *http.Request) {
	committeeID := h.requiredQuery(w, r, "committeeId")
	if committeeID == "" {
		return
	}

	result, err := h.svc.Performance.GetCommitteeWorkProducts(r.Context(), committeeID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffCommitteeWorkProducts").Str("committeeId", committeeID).Msg("Failed to get staff committee work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetOperationalWorkProducts handles GET /api/v1/pms-engine/work-products/operational?staffId={id}&reviewPeriodId={id}
// Mirrors .NET GetOperationalWorkProducts -- retrieves work products for a staff member in a review period.
func (h *PmsEngineHandler) GetOperationalWorkProducts(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}

	result, err := h.svc.Performance.GetStaffWorkProducts(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetOperationalWorkProducts").Str("staffId", staffID).Msg("Failed to get operational work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetObjectiveWorkProducts handles GET /api/v1/pms-engine/work-products/by-objective?staffId={id}&objectiveId={id}
// Mirrors .NET GetObjectiveWorkProducts -- retrieves work products linked to an objective.
func (h *PmsEngineHandler) GetObjectiveWorkProducts(w http.ResponseWriter, r *http.Request) {
	objectiveID := h.requiredQuery(w, r, "objectiveId")
	if objectiveID == "" {
		return
	}

	result, err := h.svc.Performance.GetObjectiveWorkProducts(r.Context(), objectiveID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetObjectiveWorkProducts").Str("objectiveId", objectiveID).Msg("Failed to get objective work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllStaffWorkProducts handles GET /api/v1/pms-engine/work-products/all?staffId={id}
// Mirrors .NET GetAllStaffWorkProducts -- retrieves all work products for a staff member.
func (h *PmsEngineHandler) GetAllStaffWorkProducts(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetAllStaffWorkProducts(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllStaffWorkProducts").Str("staffId", staffID).Msg("Failed to get all staff work products")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== WORK PRODUCT TASK HANDLERS ============================

// AddWorkProductTask handles POST /api/v1/pms-engine/work-products/tasks
// Mirrors .NET AddWorkProductTask -- WorkProductTaskSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddWorkProductTask(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductTaskRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationAdd.String()
	result, err := h.svc.Performance.WorkProductTaskSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddWorkProductTask").Msg("Failed to add work product task")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateWorkProductTask handles PUT /api/v1/pms-engine/work-products/tasks
// Mirrors .NET UpdateWorkProductTask -- WorkProductTaskSetup(model, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateWorkProductTask(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductTaskRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationUpdate.String()
	result, err := h.svc.Performance.WorkProductTaskSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateWorkProductTask").Msg("Failed to update work product task")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelWorkProductTask handles POST /api/v1/pms-engine/work-products/tasks/cancel
// Mirrors .NET CancelWorkProductTask -- WorkProductTaskSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelWorkProductTask(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductTaskRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationCancel.String()
	result, err := h.svc.Performance.WorkProductTaskSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelWorkProductTask").Msg("Failed to cancel work product task")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CompleteWorkProductTask handles POST /api/v1/pms-engine/work-products/tasks/complete
// Mirrors .NET CompleteWorkProductTask -- WorkProductTaskSetup(request, OperationTypes.Complete).
func (h *PmsEngineHandler) CompleteWorkProductTask(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductTaskRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationComplete.String()
	result, err := h.svc.Performance.WorkProductTaskSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CompleteWorkProductTask").Msg("Failed to complete work product task")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetWorkProductTaskDetail handles GET /api/v1/pms-engine/work-products/tasks/{taskId}
// Mirrors .NET GetWorkProductTaskDetail -- retrieves a single task by ID.
func (h *PmsEngineHandler) GetWorkProductTaskDetail(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskId")
	if taskID == "" {
		response.Error(w, http.StatusBadRequest, "taskId is required")
		return
	}

	result, err := h.svc.Performance.GetWorkProductTaskDetail(r.Context(), taskID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetWorkProductTaskDetail").Str("taskId", taskID).Msg("Failed to get work product task detail")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetWorkProductTasks handles GET /api/v1/pms-engine/work-products/{workProductId}/tasks
// Mirrors .NET GetWorkProductTasks -- retrieves all tasks for a work product.
func (h *PmsEngineHandler) GetWorkProductTasks(w http.ResponseWriter, r *http.Request) {
	workProductID := r.PathValue("workProductId")
	if workProductID == "" {
		response.Error(w, http.StatusBadRequest, "workProductId is required")
		return
	}

	result, err := h.svc.Performance.GetWorkProductTasks(r.Context(), workProductID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetWorkProductTasks").Str("workProductId", workProductID).Msg("Failed to get work product tasks")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== WORK PRODUCT EVALUATION HANDLERS ======================

// AddWorkProductEvaluation handles POST /api/v1/pms-engine/work-products/evaluation
// Mirrors .NET AddWorkProductEvaluation -- WorkProductEvaluation(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddWorkProductEvaluation(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationAdd.String()
	result, err := h.svc.Performance.WorkProductEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddWorkProductEvaluation").Msg("Failed to add work product evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateWorkProductEvaluation handles PUT /api/v1/pms-engine/work-products/evaluation
// Mirrors .NET UpdateWorkProductEvaluation -- WorkProductEvaluation(model, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateWorkProductEvaluation(w http.ResponseWriter, r *http.Request) {
	var req performance.WorkProductEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationUpdate.String()
	result, err := h.svc.Performance.WorkProductEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateWorkProductEvaluation").Msg("Failed to update work product evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetWorkProductEvaluation handles GET /api/v1/pms-engine/work-products/{workProductId}/evaluation
// Mirrors .NET GetWorkProductEvaluation -- retrieves evaluation for a work product.
func (h *PmsEngineHandler) GetWorkProductEvaluation(w http.ResponseWriter, r *http.Request) {
	workProductID := r.PathValue("workProductId")
	if workProductID == "" {
		response.Error(w, http.StatusBadRequest, "workProductId is required")
		return
	}

	result, err := h.svc.Performance.GetWorkProductEvaluation(r.Context(), workProductID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetWorkProductEvaluation").Str("workProductId", workProductID).Msg("Failed to get work product evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// InitiateWorkProductReEvaluation handles POST /api/v1/pms-engine/work-products/{workProductId}/re-evaluate
// Mirrors .NET InitiateWorkProductReEvaluation.
func (h *PmsEngineHandler) InitiateWorkProductReEvaluation(w http.ResponseWriter, r *http.Request) {
	workProductID := r.PathValue("workProductId")
	if workProductID == "" {
		response.Error(w, http.StatusBadRequest, "workProductId is required")
		return
	}

	result, err := h.svc.Performance.InitiateWorkProductReEvaluation(r.Context(), workProductID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "InitiateWorkProductReEvaluation").Str("workProductId", workProductID).Msg("Failed to initiate re-evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReCalculateWorkProductPoints handles POST /api/v1/pms-engine/work-products/recalculate?staffId={id}&reviewPeriodId={id}
// Mirrors .NET ReCalculateWorkProductPoints.
func (h *PmsEngineHandler) ReCalculateWorkProductPoints(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}

	result, err := h.svc.Performance.ReCalculateWorkProductPoints(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReCalculateWorkProductPoints").Str("staffId", staffID).Msg("Failed to recalculate work product points")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== PERIOD OBJECTIVE EVALUATION HANDLERS ==================

// SaveDraftEvaluation handles POST /api/v1/pms-engine/evaluations/draft
// Mirrors .NET AddObjectiveOutcomeScore with Draft operation.
func (h *PmsEngineHandler) SaveDraftEvaluation(w http.ResponseWriter, r *http.Request) {
	var req performance.PeriodObjectiveEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationDraft.String()
	result, err := h.svc.Performance.ReviewPeriodObjectiveEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftEvaluation").Msg("Failed to save draft evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddEvaluation handles POST /api/v1/pms-engine/evaluations
// Mirrors .NET AddObjectiveOutcomeScore -- performanceManagementService.ReviewPeriodObjectiveEvaluation(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddEvaluation(w http.ResponseWriter, r *http.Request) {
	var req performance.PeriodObjectiveEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationAdd.String()
	result, err := h.svc.Performance.ReviewPeriodObjectiveEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddEvaluation").Msg("Failed to add evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// SubmitDraftEvaluation handles POST /api/v1/pms-engine/evaluations/submit-draft
// Mirrors .NET commit draft evaluation flow.
func (h *PmsEngineHandler) SubmitDraftEvaluation(w http.ResponseWriter, r *http.Request) {
	var req performance.PeriodObjectiveEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationCommitDraft.String()
	result, err := h.svc.Performance.ReviewPeriodObjectiveEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftEvaluation").Msg("Failed to submit draft evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveEvaluation handles POST /api/v1/pms-engine/evaluations/approve
// Mirrors .NET ApproveObjectiveOutcomeScore.
func (h *PmsEngineHandler) ApproveEvaluation(w http.ResponseWriter, r *http.Request) {
	var req performance.PeriodObjectiveEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationApprove.String()
	result, err := h.svc.Performance.ReviewPeriodObjectiveEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveEvaluation").Msg("Failed to approve evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectEvaluation handles POST /api/v1/pms-engine/evaluations/reject
// Mirrors .NET ReturnObjectiveOutcomeScore (return/reject evaluation).
func (h *PmsEngineHandler) RejectEvaluation(w http.ResponseWriter, r *http.Request) {
	var req performance.PeriodObjectiveEvaluationRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.Status = enums.OperationReject.String()
	result, err := h.svc.Performance.ReviewPeriodObjectiveEvaluation(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectEvaluation").Msg("Failed to reject evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffEvaluations handles GET /api/v1/pms-engine/evaluations?reviewPeriodId={id}
// Mirrors .NET GetReviewPeriodObjectiveEvaluation.
func (h *PmsEngineHandler) GetStaffEvaluations(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}

	result, err := h.svc.Performance.GetReviewPeriodObjectiveEvaluations(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffEvaluations").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get staff evaluations")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== FEEDBACK HANDLERS =====================================

// RequestFeedback handles POST /api/v1/pms-engine/feedback/request
// Mirrors .NET feedback request initiation flow.
func (h *PmsEngineHandler) RequestFeedback(w http.ResponseWriter, r *http.Request) {
	var req performance.TreatFeedbackRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.TreatAssignedRequest(r.Context(), &req); err != nil {
		h.log.Error().Err(err).Str("action", "RequestFeedback").Msg("Failed to request feedback")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, map[string]string{"message": "Feedback requested successfully"})
}

// GetFeedbackRequests handles GET /api/v1/pms-engine/feedback/requests?staffId={id}
// Mirrors .NET GetStaffRequests -- performanceManagementService.GetRequests(staffId).
func (h *PmsEngineHandler) GetFeedbackRequests(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetRequests(r.Context(), staffID, nil, nil)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetFeedbackRequests").Str("staffId", staffID).Msg("Failed to get feedback requests")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ProcessFeedback handles POST /api/v1/pms-engine/feedback/process
// Mirrors .NET TreatAssignedRequest / CloseRequest / ReassignRequest.
func (h *PmsEngineHandler) ProcessFeedback(w http.ResponseWriter, r *http.Request) {
	var req performance.TreatFeedbackRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.TreatAssignedRequest(r.Context(), &req); err != nil {
		h.log.Error().Err(err).Str("action", "ProcessFeedback").Msg("Failed to process feedback")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, map[string]string{"message": "Feedback processed successfully"})
}

// GetPendingFeedbackActions handles GET /api/v1/pms-engine/feedback/pending?staffId={id}
// Mirrors .NET GetStaffRequestsByStatus with pending status filter.
func (h *PmsEngineHandler) GetPendingFeedbackActions(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetPendingRequests(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPendingFeedbackActions").Str("staffId", staffID).Msg("Failed to get pending feedback actions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== SCORING HANDLERS ======================================

// GetPerformanceScore handles GET /api/v1/pms-engine/scores?staffId={id}&reviewPeriodId={id}
// Mirrors .NET GetStaffPerformanceScoreCardStatistics.
func (h *PmsEngineHandler) GetPerformanceScore(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetPerformanceScore(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPerformanceScore").Str("staffId", staffID).Msg("Failed to get performance score")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetDashboardStats handles GET /api/v1/pms-engine/dashboard?staffId={id}
// Mirrors .NET GetRequestStatistics / GetStaffPerformanceStatistics /
// GetStaffWorkProductsStatistics -- aggregated dashboard data.
func (h *PmsEngineHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	result, err := h.svc.Performance.GetDashboardStats(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetDashboardStats").Str("staffId", staffID).Msg("Failed to get dashboard stats")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetPerformanceSummary handles GET /api/v1/pms-engine/scores/summary?reviewPeriodId={id}&referenceId={id}&organogramLevel={level}
// Mirrors .NET GetPeriodScores / GetOrganogramPerformanceSummaryStatistics.
func (h *PmsEngineHandler) GetPerformanceSummary(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}

	referenceID := r.URL.Query().Get("referenceId")
	organogramLevelStr := r.URL.Query().Get("organogramLevel")

	level, err := strconv.Atoi(organogramLevelStr)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPerformanceSummary").Str("organogramLevel", organogramLevelStr).Msg("Invalid organogram level")
		response.Error(w, http.StatusBadRequest, "organogramLevel must be a valid integer")
		return
	}

	result, err := h.svc.Performance.GetOrganogramPerformanceSummaryStatistics(r.Context(), referenceID, reviewPeriodID, enums.OrganogramLevel(level))
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPerformanceSummary").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get performance summary")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== INDIVIDUAL OBJECTIVE PLANNING HANDLERS ================

// SaveDraftIndividualPlannedObjective handles POST /api/v1/pms-engine/individual-objectives/draft
// Mirrors .NET individual objective draft creation.
func (h *PmsEngineHandler) SaveDraftIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.SaveDraftIndividualPlannedObjective(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftIndividualPlannedObjective").Msg("Failed to save draft individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// AddIndividualPlannedObjective handles POST /api/v1/pms-engine/individual-objectives
// Mirrors .NET individual objective creation with commit.
func (h *PmsEngineHandler) AddIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.AddReviewPeriodIndividualPlannedObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.AddIndividualPlannedObjective(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "AddIndividualPlannedObjective").Msg("Failed to add individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// SubmitDraftIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/submit-draft
// Mirrors .NET commit draft individual objective.
func (h *PmsEngineHandler) SubmitDraftIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.SubmitDraftIndividualPlannedObjective(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftIndividualObjective").Msg("Failed to submit draft individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ApproveIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/approve
// Mirrors .NET approve individual objective.
func (h *PmsEngineHandler) ApproveIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.ApproveIndividualPlannedObjective(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveIndividualObjective").Msg("Failed to approve individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/reject
// Mirrors .NET reject individual objective.
func (h *PmsEngineHandler) RejectIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.RejectIndividualPlannedObjective(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectIndividualObjective").Msg("Failed to reject individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReturnIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/return
// Mirrors .NET return individual objective.
func (h *PmsEngineHandler) ReturnIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.ReturnIndividualPlannedObjective(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReturnIndividualObjective").Msg("Failed to return individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CancelIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/cancel
// Mirrors .NET cancel individual objective.
func (h *PmsEngineHandler) CancelIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req performance.ReviewPeriodIndividualPlannedObjectiveRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.CancelIndividualPlannedObjective(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CancelIndividualObjective").Msg("Failed to cancel individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffIndividualObjectives handles GET /api/v1/pms-engine/individual-objectives?staffId={id}&reviewPeriodId={id}
// Mirrors .NET GetIndividualObjectives.
func (h *PmsEngineHandler) GetStaffIndividualObjectives(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}

	result, err := h.svc.ReviewPeriod.GetStaffIndividualPlannedObjectives(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffIndividualObjectives").Str("staffId", staffID).Msg("Failed to get staff individual objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== 360 REVIEW HANDLERS ===================================

// Trigger360Review handles POST /api/v1/pms-engine/360-review/trigger
// Mirrors .NET Trigger360Review — creates a 360 review configuration for a review period.
func (h *PmsEngineHandler) Trigger360Review(w http.ResponseWriter, r *http.Request) {
	var req performance.CreateReviewPeriod360ReviewRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.ReviewPeriod.AddReviewPeriod360Review(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "Trigger360Review").Msg("Failed to trigger 360 review")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// Initiate360Review handles POST /api/v1/pms-engine/360-review/initiate
// Mirrors .NET Initiate360Review — initiates 360 review for selected staff.
func (h *PmsEngineHandler) Initiate360Review(w http.ResponseWriter, r *http.Request) {
	var req performance.Initiate360ReviewRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.Initiate360Review(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "Initiate360Review").Msg("Failed to initiate 360 review")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// Complete360ReviewForStaff handles POST /api/v1/pms-engine/360-review/complete
// Mirrors .NET Complete360Review — completes 360 review for a review period.
func (h *PmsEngineHandler) Complete360ReviewForStaff(w http.ResponseWriter, r *http.Request) {
	var req performance.Complete360ReviewRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.Complete360Review(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "Complete360ReviewForStaff").Msg("Failed to complete 360 review")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== 360 RATING MANAGEMENT HANDLERS ========================

// Add360Rating handles POST /api/v1/pms-engine/360-review/rating
// Mirrors .NET Add360Rating — saves a competency rating for 360 review.
func (h *PmsEngineHandler) Add360Rating(w http.ResponseWriter, r *http.Request) {
	var req performance.SavePmsCompetencyRequestVm
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.CompetencyRatingSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "Add360Rating").Msg("Failed to add 360 rating")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// Update360Rating handles PUT /api/v1/pms-engine/360-review/rating
// Mirrors .NET Update360Rating — updates a competency rating for 360 review.
func (h *PmsEngineHandler) Update360Rating(w http.ResponseWriter, r *http.Request) {
	var req performance.SavePmsCompetencyRequestVm
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.CompetencyRatingSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "Update360Rating").Msg("Failed to update 360 rating")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReviewerComplete360Review handles POST /api/v1/pms-engine/360-review/reviewer-complete
// Mirrors .NET ReviewerComplete360Review — reviewer marks their 360 review as complete.
func (h *PmsEngineHandler) ReviewerComplete360Review(w http.ResponseWriter, r *http.Request) {
	var req performance.CompetencyReviewerRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	req.RecordStatus = enums.OperationComplete.String()
	result, err := h.svc.Performance.CompetencyReviewerSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReviewerComplete360Review").Msg("Failed to complete reviewer 360 review")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== COMPETENCY REVIEW RETRIEVAL HANDLERS ==================

// GetCompetencyReviewDetail handles GET /api/v1/pms-engine/competency-review/{feedbackId}
// Mirrors .NET GetCompetencyReviewDetail — retrieves a competency review feedback by ID.
func (h *PmsEngineHandler) GetCompetencyReviewDetail(w http.ResponseWriter, r *http.Request) {
	feedbackID := r.PathValue("feedbackId")
	if feedbackID == "" {
		response.Error(w, http.StatusBadRequest, "feedbackId is required")
		return
	}
	result, err := h.svc.Performance.GetCompetencyReviewFeedback(r.Context(), feedbackID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCompetencyReviewDetail").Str("feedbackId", feedbackID).Msg("Failed to get competency review detail")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCompetencyReviewFeedbackDetails handles GET /api/v1/pms-engine/competency-review/{feedbackId}/details
// Mirrors .NET GetCompetencyReviewFeedbackDetails — retrieves detailed feedback including reviewers.
func (h *PmsEngineHandler) GetCompetencyReviewFeedbackDetails(w http.ResponseWriter, r *http.Request) {
	feedbackID := r.PathValue("feedbackId")
	if feedbackID == "" {
		response.Error(w, http.StatusBadRequest, "feedbackId is required")
		return
	}
	result, err := h.svc.Performance.GetCompetencyReviewFeedbackDetails(r.Context(), feedbackID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCompetencyReviewFeedbackDetails").Str("feedbackId", feedbackID).Msg("Failed to get competency review feedback details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllCompetencyReviewFeedbacksByReviewPeriod handles GET /api/v1/pms-engine/competency-review/feedbacks?staffId={id}
// Mirrors .NET GetAllCompetencyReviewFeedbacksByReviewPeriod — retrieves all feedbacks for a staff member.
func (h *PmsEngineHandler) GetAllCompetencyReviewFeedbacksByReviewPeriod(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetAllCompetencyReviewFeedbacks(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllCompetencyReviewFeedbacksByReviewPeriod").Str("staffId", staffID).Msg("Failed to get competency review feedbacks")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllMyReviewedCompetencies handles GET /api/v1/pms-engine/competency-review/my-reviewed?reviewerStaffId={id}
// Mirrors .NET GetAllMyReviewedCompetencies — retrieves competencies reviewed by the current reviewer.
func (h *PmsEngineHandler) GetAllMyReviewedCompetencies(w http.ResponseWriter, r *http.Request) {
	reviewerStaffID := h.requiredQuery(w, r, "reviewerStaffId")
	if reviewerStaffID == "" {
		return
	}
	result, err := h.svc.Performance.GetCompetencyReviews(r.Context(), reviewerStaffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllMyReviewedCompetencies").Str("reviewerStaffId", reviewerStaffID).Msg("Failed to get reviewed competencies")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCompetenciesToReview handles GET /api/v1/pms-engine/competency-review/to-review?reviewerStaffId={id}
// Mirrors .NET GetCompetenciesToReview — retrieves competencies pending review by a reviewer.
func (h *PmsEngineHandler) GetCompetenciesToReview(w http.ResponseWriter, r *http.Request) {
	reviewerStaffID := h.requiredQuery(w, r, "reviewerStaffId")
	if reviewerStaffID == "" {
		return
	}
	result, err := h.svc.Performance.GetCompetencyReviews(r.Context(), reviewerStaffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCompetenciesToReview").Str("reviewerStaffId", reviewerStaffID).Msg("Failed to get competencies to review")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetReviewerFeedbackDetails handles GET /api/v1/pms-engine/competency-review/reviewer/{reviewerId}
// Mirrors .NET GetReviewerFeedbackDetails — retrieves detailed feedback for a specific reviewer.
func (h *PmsEngineHandler) GetReviewerFeedbackDetails(w http.ResponseWriter, r *http.Request) {
	reviewerID := r.PathValue("reviewerId")
	if reviewerID == "" {
		response.Error(w, http.StatusBadRequest, "reviewerId is required")
		return
	}
	result, err := h.svc.Performance.GetReviewerFeedbackDetails(r.Context(), reviewerID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetReviewerFeedbackDetails").Str("reviewerId", reviewerID).Msg("Failed to get reviewer feedback details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetQuestionnaire handles GET /api/v1/pms-engine/competency-review/questionnaire?staffId={id}
// Mirrors .NET GetQuestionnaire — retrieves the questionnaire for a staff member.
func (h *PmsEngineHandler) GetQuestionnaire(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetQuestionnaire(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetQuestionnaire").Str("staffId", staffID).Msg("Failed to get questionnaire")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== FEEDBACK REQUEST RETRIEVAL HANDLERS ===================

// GetStaffRequests handles GET /api/v1/pms-engine/feedback/requests/staff?staffId={id}
// Mirrors .NET GetStaffRequests — retrieves all feedback requests assigned to a staff member.
func (h *PmsEngineHandler) GetStaffRequests(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetRequests(r.Context(), staffID, nil, nil)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffRequests").Str("staffId", staffID).Msg("Failed to get staff requests")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetBreachedRequests handles GET /api/v1/pms-engine/feedback/requests/breached?staffId={id}&reviewPeriodId={id}
// Mirrors .NET GetBreachedRequests — retrieves breached (SLA-violated) feedback requests.
func (h *PmsEngineHandler) GetBreachedRequests(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}
	result, err := h.svc.Performance.GetBreachedRequests(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetBreachedRequests").Str("staffId", staffID).Msg("Failed to get breached requests")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffRequestsByStatus handles GET /api/v1/pms-engine/feedback/requests/staff/by-status?staffId={id}&status={status}
// Mirrors .NET GetStaffRequestsByStatus — retrieves feedback requests filtered by status.
func (h *PmsEngineHandler) GetStaffRequestsByStatus(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	status := h.requiredQuery(w, r, "status")
	if status == "" {
		return
	}
	result, err := h.svc.Performance.GetRequests(r.Context(), staffID, nil, &status)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffRequestsByStatus").Str("staffId", staffID).Str("status", status).Msg("Failed to get staff requests by status")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllRequests handles GET /api/v1/pms-engine/feedback/requests/all?staffId={id}
// Mirrors .NET GetAllRequests — retrieves all feedback requests owned by a staff member.
func (h *PmsEngineHandler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetRequestsByOwner(r.Context(), staffID, nil)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllRequests").Str("staffId", staffID).Msg("Failed to get all requests")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetRequestsByStatus handles GET /api/v1/pms-engine/feedback/requests/by-status?staffId={id}&status={status}
// Mirrors .NET GetRequestsByStatus — retrieves feedback requests filtered by status (owner perspective).
func (h *PmsEngineHandler) GetRequestsByStatus(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	status := h.requiredQuery(w, r, "status")
	if status == "" {
		return
	}
	result, err := h.svc.Performance.GetRequests(r.Context(), staffID, nil, &status)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetRequestsByStatus").Str("staffId", staffID).Str("status", status).Msg("Failed to get requests by status")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetRequestDetails handles GET /api/v1/pms-engine/feedback/requests/{requestId}
// Mirrors .NET GetRequestDetails — retrieves details for a specific feedback request.
func (h *PmsEngineHandler) GetRequestDetails(w http.ResponseWriter, r *http.Request) {
	requestID := r.PathValue("requestId")
	if requestID == "" {
		response.Error(w, http.StatusBadRequest, "requestId is required")
		return
	}
	result, err := h.svc.Performance.GetRequestDetails(r.Context(), requestID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetRequestDetails").Str("requestId", requestID).Msg("Failed to get request details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== FEEDBACK REQUEST ACTION HANDLERS ======================

// ReassignRequest handles POST /api/v1/pms-engine/feedback/requests/reassign
// Mirrors .NET ReassignRequest — reassigns a feedback request to a different staff member.
func (h *PmsEngineHandler) ReassignRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestID         string `json:"requestId"`
		NewAssignedStaffID string `json:"newAssignedStaffId"`
	}
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.ReassignRequest(r.Context(), req.RequestID, req.NewAssignedStaffID); err != nil {
		h.log.Error().Err(err).Str("action", "ReassignRequest").Str("requestId", req.RequestID).Msg("Failed to reassign request")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, map[string]string{"message": "Request reassigned successfully"})
}

// ReassignSelfRequest handles POST /api/v1/pms-engine/feedback/requests/reassign-self
// Mirrors .NET ReassignSelfRequest — staff member reassigns their own feedback request.
func (h *PmsEngineHandler) ReassignSelfRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestID         string `json:"requestId"`
		CurrentStaffID    string `json:"currentStaffId"`
		NewAssignedStaffID string `json:"newAssignedStaffId"`
	}
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.ReassignSelfRequest(r.Context(), req.RequestID, req.CurrentStaffID, req.NewAssignedStaffID); err != nil {
		h.log.Error().Err(err).Str("action", "ReassignSelfRequest").Str("requestId", req.RequestID).Msg("Failed to self-reassign request")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, map[string]string{"message": "Request reassigned successfully"})
}

// CloseRequest handles POST /api/v1/pms-engine/feedback/requests/close
// Mirrors .NET CloseRequest — closes a feedback request.
func (h *PmsEngineHandler) CloseRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RequestID string `json:"requestId"`
	}
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.CloseRequest(r.Context(), req.RequestID); err != nil {
		h.log.Error().Err(err).Str("action", "CloseRequest").Str("requestId", req.RequestID).Msg("Failed to close request")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, map[string]string{"message": "Request closed successfully"})
}

// TreatAssignedRequest handles POST /api/v1/pms-engine/feedback/requests/treat
// Mirrors .NET TreatAssignedRequest — treats (processes) an assigned feedback request.
func (h *PmsEngineHandler) TreatAssignedRequest(w http.ResponseWriter, r *http.Request) {
	var req performance.TreatFeedbackRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if err := h.svc.Performance.TreatAssignedRequest(r.Context(), &req); err != nil {
		h.log.Error().Err(err).Str("action", "TreatAssignedRequest").Str("requestId", req.RequestID).Msg("Failed to treat assigned request")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, map[string]string{"message": "Request treated successfully"})
}

// =================== COMPETENCY GAP CLOSURE HANDLERS =======================

// CompetencyGapClosureSetup handles POST /api/v1/pms-engine/competency-review/gap-closure
// Mirrors .NET CompetencyGapClosureSetup — creates or updates competency gap closure records.
func (h *PmsEngineHandler) CompetencyGapClosureSetup(w http.ResponseWriter, r *http.Request) {
	var req performance.CompetencyGapClosureRequestModel
	if !h.decodeJSON(w, r, &req) {
		return
	}
	result, err := h.svc.Performance.CompetencyGapClosureSetup(r.Context(), &req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CompetencyGapClosureSetup").Msg("Failed to setup competency gap closure")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== DASHBOARD & STATISTICS HANDLERS =======================

// GetRequestStatistics handles GET /api/v1/pms-engine/stats/requests?staffId=X
// Mirrors .NET GetRequestStatistics — feedback request SLA dashboard statistics.
func (h *PmsEngineHandler) GetRequestStatistics(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetRequestStatistics(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetRequestStatistics").Str("staffId", staffID).Msg("Failed to get request statistics")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffPerformanceStatistics handles GET /api/v1/pms-engine/stats/performance?staffId=X
// Mirrors .NET GetStaffPerformanceStatistics — points dashboard statistics.
func (h *PmsEngineHandler) GetStaffPerformanceStatistics(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetStaffPerformanceStatistics(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffPerformanceStatistics").Str("staffId", staffID).Msg("Failed to get performance statistics")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffWorkProductsStatistics handles GET /api/v1/pms-engine/stats/work-products?staffId=X
// Mirrors .NET GetStaffWorkProductsStatistics — work product count summary.
func (h *PmsEngineHandler) GetStaffWorkProductsStatistics(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetStaffWorkProductsStatistics(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffWorkProductsStatistics").Str("staffId", staffID).Msg("Failed to get work products statistics")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffWorkProductsDetailsStatistics handles GET /api/v1/pms-engine/stats/work-products-details?staffId=X
// Mirrors .NET GetStaffWorkProductsDetailsStatistics — detailed work product breakdown.
func (h *PmsEngineHandler) GetStaffWorkProductsDetailsStatistics(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetStaffWorkProductsDetailsStatistics(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffWorkProductsDetailsStatistics").Str("staffId", staffID).Msg("Failed to get work products details statistics")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== SCORECARD STATISTICS HANDLERS =========================

// GetStaffPerformanceScoreCardStatistics handles GET /api/v1/pms-engine/scorecard?staffId=X&reviewPeriodId=Y
// Mirrors .NET GetStaffPerformanceScoreCardStatistics — full score card.
func (h *PmsEngineHandler) GetStaffPerformanceScoreCardStatistics(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}
	result, err := h.svc.Performance.GetStaffPerformanceScoreCardStatistics(r.Context(), staffID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffPerformanceScoreCardStatistics").Str("staffId", staffID).Msg("Failed to get scorecard statistics")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffAnnualPerformanceScoreCardStatistics handles GET /api/v1/pms-engine/scorecard/annual?staffId=X&year=Y
// Mirrors .NET GetStaffAnnualPerformanceScoreCardStatistics — annual score card.
func (h *PmsEngineHandler) GetStaffAnnualPerformanceScoreCardStatistics(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	yearStr := h.requiredQuery(w, r, "year")
	if yearStr == "" {
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "year must be a valid integer")
		return
	}
	result, err := h.svc.Performance.GetStaffAnnualPerformanceScoreCardStatistics(r.Context(), staffID, year)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffAnnualPerformanceScoreCardStatistics").Str("staffId", staffID).Int("year", year).Msg("Failed to get annual scorecard")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetSubordinatesStaffPerformanceScoreCardStatistics handles GET /api/v1/pms-engine/scorecard/subordinates?managerId=X&reviewPeriodId=Y
// Mirrors .NET GetSurbodinatesStaffPerformanceScoreCardStatistics — subordinates score cards.
func (h *PmsEngineHandler) GetSubordinatesStaffPerformanceScoreCardStatistics(w http.ResponseWriter, r *http.Request) {
	managerID := h.requiredQuery(w, r, "managerId")
	if managerID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}
	result, err := h.svc.Performance.GetSubordinatesStaffPerformanceScoreCardStatistics(r.Context(), managerID, reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetSubordinatesScoreCard").Str("managerId", managerID).Msg("Failed to get subordinates scorecard")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== ORGANOGRAM PERFORMANCE HANDLERS =======================

// GetOrganogramPerformanceSummaryStatistics handles GET /api/v1/pms-engine/organogram-performance?referenceId=X&reviewPeriodId=Y&level=Z
// Mirrors .NET GetOrganogramPerformanceSummaryStatistics — single org unit performance summary.
func (h *PmsEngineHandler) GetOrganogramPerformanceSummaryStatistics(w http.ResponseWriter, r *http.Request) {
	referenceID := h.requiredQuery(w, r, "referenceId")
	if referenceID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}
	levelStr := r.URL.Query().Get("level")
	level := enums.OrganogramLevelOffice // default to Office per .NET
	if levelStr != "" {
		l, err := strconv.Atoi(levelStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "level must be a valid integer")
			return
		}
		level = enums.OrganogramLevel(l)
	}
	result, err := h.svc.Performance.GetOrganogramPerformanceSummaryStatistics(r.Context(), referenceID, reviewPeriodID, level)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetOrganogramPerformanceSummary").Str("referenceId", referenceID).Msg("Failed to get organogram performance summary")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetOrganogramPerformanceSummaryListStatistics handles GET /api/v1/pms-engine/organogram-performance/list?headOfUnitId=X&reviewPeriodId=Y&level=Z
// Mirrors .NET GetOrganogramPerformanceSummaryListStatistics — list of org unit performance summaries.
func (h *PmsEngineHandler) GetOrganogramPerformanceSummaryListStatistics(w http.ResponseWriter, r *http.Request) {
	headOfUnitID := h.requiredQuery(w, r, "headOfUnitId")
	if headOfUnitID == "" {
		return
	}
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}
	levelStr := r.URL.Query().Get("level")
	level := enums.OrganogramLevelDivision // default to Division per .NET
	if levelStr != "" {
		l, err := strconv.Atoi(levelStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "level must be a valid integer")
			return
		}
		level = enums.OrganogramLevel(l)
	}
	result, err := h.svc.Performance.GetOrganogramPerformanceSummaryListStatistics(r.Context(), headOfUnitID, reviewPeriodID, level)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetOrganogramPerformanceSummaryList").Str("headOfUnitId", headOfUnitID).Msg("Failed to get organogram performance summary list")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== PERIOD SCORE HANDLERS =================================

// GetPeriodScoreDetails handles GET /api/v1/pms-engine/period-scores?reviewPeriodId=X&staffId=Y
// Mirrors .NET GetPeriodScoreDetails — single staff + review period score.
func (h *PmsEngineHandler) GetPeriodScoreDetails(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetPeriodScoreDetails(r.Context(), reviewPeriodID, staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPeriodScoreDetails").Str("staffId", staffID).Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get period score details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetPeriodScores handles GET /api/v1/pms-engine/period-scores/all?reviewPeriodId=X
// Mirrors .NET GetPeriodScores — all scores for a review period.
func (h *PmsEngineHandler) GetPeriodScores(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}
	result, err := h.svc.Performance.GetPeriodScores(r.Context(), reviewPeriodID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPeriodScores").Str("reviewPeriodId", reviewPeriodID).Msg("Failed to get period scores")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetStaffReviewPeriods handles GET /api/v1/pms-engine/staff-review-periods?staffId=X
// Mirrors .NET GetStaffReviewPeriods — review periods a staff has participated in.
func (h *PmsEngineHandler) GetStaffReviewPeriods(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	result, err := h.svc.Performance.GetStaffReviewPeriods(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffReviewPeriods").Str("staffId", staffID).Msg("Failed to get staff review periods")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== AUDIT LOG HANDLERS ====================================

// GetAuditLogs handles GET /api/v1/pms-engine/audit-logs
// Mirrors .NET GetAuditLogs — retrieves recent audit log entries.
func (h *PmsEngineHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetAuditLogs(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAuditLogs").Msg("Failed to get audit logs")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAuditLogDetails handles GET /api/v1/pms-engine/audit-logs/{id}
// Mirrors .NET GetAuditLogDetails — retrieves a single audit log entry.
func (h *PmsEngineHandler) GetAuditLogDetails(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "id must be a valid integer")
		return
	}
	result, err := h.svc.Performance.GetAuditLogDetails(r.Context(), id)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAuditLogDetails").Int("id", id).Msg("Failed to get audit log details")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== LINE MANAGER & STAFF HANDLERS =========================

// GetLineManagerEmployees handles GET /api/v1/pms-engine/line-manager-employees?staffId=X&category=Y
// Mirrors .NET GetLineManagerEmployees — retrieves employees under a line manager.
func (h *PmsEngineHandler) GetLineManagerEmployees(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}
	categoryStr := h.requiredQuery(w, r, "category")
	if categoryStr == "" {
		return
	}
	categoryInt, err := strconv.Atoi(categoryStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "category must be a valid integer")
		return
	}
	category := enums.LineManagerPerformanceCategory(categoryInt)
	result, err := h.svc.Performance.GetLineManagerEmployees(r.Context(), staffID, category)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetLineManagerEmployees").Str("staffId", staffID).Msg("Failed to get line manager employees")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAdhocAssignmentEmployees handles GET /api/v1/pms-engine/adhoc-employees?leadStaffId=X&category=Y
// Mirrors .NET GetAdhocAssignmentEmployees — retrieves adhoc assignment employees.
func (h *PmsEngineHandler) GetAdhocAssignmentEmployees(w http.ResponseWriter, r *http.Request) {
	leadStaffID := h.requiredQuery(w, r, "leadStaffId")
	if leadStaffID == "" {
		return
	}
	categoryStr := h.requiredQuery(w, r, "category")
	if categoryStr == "" {
		return
	}
	categoryInt, err := strconv.Atoi(categoryStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "category must be a valid integer")
		return
	}
	category := enums.LineManagerPerformanceCategory(categoryInt)
	result, err := h.svc.Performance.GetAdhocAssignmentEmployees(r.Context(), leadStaffID, category)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAdhocAssignmentEmployees").Str("leadStaffId", leadStaffID).Msg("Failed to get adhoc assignment employees")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetMyStaff handles GET /api/v1/pms-engine/my-staff?managerId=X
// Mirrors .NET GetMyStaff — retrieves staff under a manager.
func (h *PmsEngineHandler) GetMyStaff(w http.ResponseWriter, r *http.Request) {
	managerID := h.requiredQuery(w, r, "managerId")
	if managerID == "" {
		return
	}
	result, err := h.svc.Performance.GetMyStaff(r.Context(), managerID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetMyStaff").Str("managerId", managerID).Msg("Failed to get my staff")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== PASSWORD MANAGEMENT HANDLERS ==========================

// ResetUserPassword handles POST /api/v1/pms-engine/reset-password
// Mirrors .NET ResetUserPassword — resets a user's password (AllowAnonymous).
func (h *PmsEngineHandler) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		DeviceName string `json:"deviceName"`
		IPAddress  string `json:"ipAddress"`
	}
	if !h.decodeJSON(w, r, &req) {
		return
	}
	if req.Username == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "username and password are required")
		return
	}
	result, err := h.svc.Performance.ResetUserPassword(r.Context(), req.Username, req.Password, req.IPAddress, req.DeviceName)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ResetUserPassword").Str("username", req.Username).Msg("Failed to reset user password")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =================== ROUTE REGISTRATION ====================================

// RegisterRoutes registers all PMS engine routes on the given mux.
// Uses Go 1.22+ method-aware patterns. All routes require JWT authentication.
func (h *PmsEngineHandler) RegisterRoutes(mux *http.ServeMux, mw *middleware.Stack) {
	base := "/api/v1/pms-engine"
	jwt := func(hf http.HandlerFunc) http.Handler { return mw.JWTAuth(http.HandlerFunc(hf)) }

	// --- Project Management ---
	mux.Handle("POST "+base+"/projects/draft", jwt(h.SaveDraftProject))
	mux.Handle("POST "+base+"/projects", jwt(h.AddProject))
	mux.Handle("POST "+base+"/projects/submit-draft", jwt(h.SubmitDraftProject))
	mux.Handle("POST "+base+"/projects/approve", jwt(h.ApproveProject))
	mux.Handle("POST "+base+"/projects/reject", jwt(h.RejectProject))
	mux.Handle("POST "+base+"/projects/return", jwt(h.ReturnProject))
	mux.Handle("POST "+base+"/projects/resubmit", jwt(h.ReSubmitProject))
	mux.Handle("PUT "+base+"/projects", jwt(h.UpdateProject))
	mux.Handle("POST "+base+"/projects/cancel", jwt(h.CancelProject))
	mux.Handle("GET "+base+"/projects", jwt(h.GetProjects))
	mux.Handle("GET "+base+"/projects/{projectId}", jwt(h.GetProjectDetails))
	mux.Handle("POST "+base+"/projects/objectives", jwt(h.AddProjectObjective))
	mux.Handle("POST "+base+"/projects/members", jwt(h.AddProjectMember))
	mux.Handle("GET "+base+"/projects/{projectId}/members", jwt(h.GetProjectMembers))
	mux.Handle("GET "+base+"/projects/{projectId}/objectives", jwt(h.GetProjectObjectives))
	mux.Handle("POST "+base+"/projects/close", jwt(h.CloseProject))
	mux.Handle("POST "+base+"/projects/pause", jwt(h.PauseProject))
	mux.Handle("GET "+base+"/projects/by-manager", jwt(h.GetProjectsByManager))
	mux.Handle("GET "+base+"/projects/assigned", jwt(h.GetProjectsAssigned))
	mux.Handle("GET "+base+"/projects/staff", jwt(h.GetStaffProjects))
	mux.Handle("GET "+base+"/projects/{projectId}/work-product-staff", jwt(h.GetProjectWorkProductStaffList))
	mux.Handle("POST "+base+"/projects/members/draft", jwt(h.SaveDraftProjectMember))
	mux.Handle("POST "+base+"/projects/members/submit-draft", jwt(h.SubmitDraftProjectMember))
	mux.Handle("POST "+base+"/projects/members/accept", jwt(h.AcceptProjectMember))
	mux.Handle("POST "+base+"/projects/members/approve", jwt(h.ApproveProjectMember))
	mux.Handle("POST "+base+"/projects/members/cancel", jwt(h.CancelProjectMember))
	mux.Handle("POST "+base+"/projects/objectives/cancel", jwt(h.CancelProjectObjective))
	mux.Handle("POST "+base+"/projects/change-lead", jwt(h.ChangeAdhocAssignmentLead))
	mux.Handle("GET "+base+"/projects/validate-eligibility", jwt(h.ValidateStaffEligibilityForAdhoc))

	// --- Committee Management ---
	mux.Handle("POST "+base+"/committees/draft", jwt(h.SaveDraftCommittee))
	mux.Handle("POST "+base+"/committees", jwt(h.AddCommittee))
	mux.Handle("POST "+base+"/committees/submit-draft", jwt(h.SubmitDraftCommittee))
	mux.Handle("POST "+base+"/committees/approve", jwt(h.ApproveCommittee))
	mux.Handle("POST "+base+"/committees/reject", jwt(h.RejectCommittee))
	mux.Handle("POST "+base+"/committees/return", jwt(h.ReturnCommittee))
	mux.Handle("POST "+base+"/committees/resubmit", jwt(h.ReSubmitCommittee))
	mux.Handle("PUT "+base+"/committees", jwt(h.UpdateCommittee))
	mux.Handle("POST "+base+"/committees/cancel", jwt(h.CancelCommittee))
	mux.Handle("GET "+base+"/committees", jwt(h.GetCommittees))
	mux.Handle("GET "+base+"/committees/{committeeId}", jwt(h.GetCommitteeDetails))
	mux.Handle("POST "+base+"/committees/members", jwt(h.AddCommitteeMember))
	mux.Handle("POST "+base+"/committees/objectives", jwt(h.AddCommitteeObjective))
	mux.Handle("POST "+base+"/committees/close", jwt(h.CloseCommittee))
	mux.Handle("POST "+base+"/committees/pause", jwt(h.PauseCommittee))
	mux.Handle("GET "+base+"/committees/by-chairperson", jwt(h.GetCommitteesByChairperson))
	mux.Handle("GET "+base+"/committees/{committeeId}/members", jwt(h.GetCommitteeMembers))
	mux.Handle("GET "+base+"/committees/assigned", jwt(h.GetCommitteesAssigned))
	mux.Handle("GET "+base+"/committees/staff", jwt(h.GetStaffCommittees))
	mux.Handle("GET "+base+"/committees/{committeeId}/work-product-staff", jwt(h.GetCommitteeWorkProductStaffList))
	mux.Handle("GET "+base+"/committees/{committeeId}/objectives", jwt(h.GetCommitteeObjectives))
	mux.Handle("POST "+base+"/committees/members/draft", jwt(h.SaveDraftCommitteeMember))
	mux.Handle("POST "+base+"/committees/members/submit-draft", jwt(h.SubmitDraftCommitteeMember))
	mux.Handle("POST "+base+"/committees/members/cancel", jwt(h.CancelCommitteeMember))
	mux.Handle("POST "+base+"/committees/objectives/cancel", jwt(h.CancelCommitteeObjective))
	mux.Handle("POST "+base+"/committees/change-chairperson", jwt(h.ChangeCommitteeChairperson))

	// --- Work Product Management ---
	mux.Handle("POST "+base+"/work-products/draft", jwt(h.SaveDraftWorkProduct))
	mux.Handle("POST "+base+"/work-products", jwt(h.AddWorkProduct))
	mux.Handle("POST "+base+"/work-products/submit-draft", jwt(h.SubmitDraftWorkProduct))
	mux.Handle("POST "+base+"/work-products/approve", jwt(h.ApproveWorkProduct))
	mux.Handle("POST "+base+"/work-products/reject", jwt(h.RejectWorkProduct))
	mux.Handle("POST "+base+"/work-products/return", jwt(h.ReturnWorkProduct))
	mux.Handle("POST "+base+"/work-products/resubmit", jwt(h.ReSubmitWorkProduct))
	mux.Handle("PUT "+base+"/work-products", jwt(h.UpdateWorkProduct))
	mux.Handle("POST "+base+"/work-products/cancel", jwt(h.CancelWorkProduct))
	mux.Handle("POST "+base+"/work-products/pause", jwt(h.PauseWorkProduct))
	mux.Handle("POST "+base+"/work-products/resume", jwt(h.ResumeWorkProduct))
	mux.Handle("GET "+base+"/work-products", jwt(h.GetStaffWorkProducts))
	mux.Handle("GET "+base+"/work-products/{workProductId}", jwt(h.GetWorkProductDetails))
	mux.Handle("POST "+base+"/work-products/assign", jwt(h.AssignWorkProduct))
	mux.Handle("GET "+base+"/work-products/assigned", jwt(h.GetAssignedWorkProducts))
	mux.Handle("POST "+base+"/work-products/evaluate", jwt(h.EvaluateWorkProduct))
	mux.Handle("POST "+base+"/work-products/complete", jwt(h.CompleteWorkProduct))
	mux.Handle("POST "+base+"/work-products/suspend", jwt(h.SuspendWorkProduct))
	mux.Handle("POST "+base+"/work-products/reinstate", jwt(h.ReInstateWorkProduct))

	// --- Project Assigned Work Products ---
	mux.Handle("POST "+base+"/work-products/project/draft", jwt(h.SaveDraftProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project", jwt(h.AddProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project/submit-draft", jwt(h.SubmitDraftProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project/approve", jwt(h.ApproveProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project/reject", jwt(h.RejectProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project/return", jwt(h.ReturnProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project/resubmit", jwt(h.ReSubmitProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project/cancel", jwt(h.CancelProjectWorkProduct))
	mux.Handle("POST "+base+"/work-products/project/close", jwt(h.CloseProjectWorkProduct))

	// --- Committee Assigned Work Products ---
	mux.Handle("POST "+base+"/work-products/committee/draft", jwt(h.SaveDraftCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee", jwt(h.AddCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee/submit-draft", jwt(h.SubmitDraftCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee/approve", jwt(h.ApproveCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee/reject", jwt(h.RejectCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee/return", jwt(h.ReturnCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee/resubmit", jwt(h.ReSubmitCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee/cancel", jwt(h.CancelCommitteeWorkProduct))
	mux.Handle("POST "+base+"/work-products/committee/close", jwt(h.CloseCommitteeWorkProduct))

	// --- Work Product Retrieval ---
	mux.Handle("GET "+base+"/work-products/project/{id}", jwt(h.GetProjectAssignedWorkProductDetails))
	mux.Handle("GET "+base+"/work-products/project", jwt(h.GetProjectAssignedWorkProducts))
	mux.Handle("GET "+base+"/work-products/project/single", jwt(h.GetProjectWorkProduct))
	mux.Handle("GET "+base+"/work-products/project/all", jwt(h.GetAllProjectWorkProducts))
	mux.Handle("GET "+base+"/work-products/project/staff", jwt(h.GetStaffProjectWorkProducts))
	mux.Handle("GET "+base+"/work-products/committee/{id}", jwt(h.GetCommitteeAssignedWorkProductDetails))
	mux.Handle("GET "+base+"/work-products/committee", jwt(h.GetCommitteeAssignedWorkProducts))
	mux.Handle("GET "+base+"/work-products/committee/single", jwt(h.GetCommitteeWorkProduct))
	mux.Handle("GET "+base+"/work-products/committee/all", jwt(h.GetAllCommitteeWorkProducts))
	mux.Handle("GET "+base+"/work-products/committee/staff", jwt(h.GetStaffCommitteeWorkProducts))
	mux.Handle("GET "+base+"/work-products/operational", jwt(h.GetOperationalWorkProducts))
	mux.Handle("GET "+base+"/work-products/by-objective", jwt(h.GetObjectiveWorkProducts))
	mux.Handle("GET "+base+"/work-products/all", jwt(h.GetAllStaffWorkProducts))

	// --- Work Product Tasks ---
	mux.Handle("POST "+base+"/work-products/tasks", jwt(h.AddWorkProductTask))
	mux.Handle("PUT "+base+"/work-products/tasks", jwt(h.UpdateWorkProductTask))
	mux.Handle("POST "+base+"/work-products/tasks/cancel", jwt(h.CancelWorkProductTask))
	mux.Handle("POST "+base+"/work-products/tasks/complete", jwt(h.CompleteWorkProductTask))
	mux.Handle("GET "+base+"/work-products/tasks/{taskId}", jwt(h.GetWorkProductTaskDetail))
	mux.Handle("GET "+base+"/work-products/{workProductId}/tasks", jwt(h.GetWorkProductTasks))

	// --- Work Product Evaluation ---
	mux.Handle("POST "+base+"/work-products/evaluation", jwt(h.AddWorkProductEvaluation))
	mux.Handle("PUT "+base+"/work-products/evaluation", jwt(h.UpdateWorkProductEvaluation))
	mux.Handle("GET "+base+"/work-products/{workProductId}/evaluation", jwt(h.GetWorkProductEvaluation))
	mux.Handle("POST "+base+"/work-products/{workProductId}/re-evaluate", jwt(h.InitiateWorkProductReEvaluation))
	mux.Handle("POST "+base+"/work-products/recalculate", jwt(h.ReCalculateWorkProductPoints))

	// --- Period Objective Evaluation ---
	mux.Handle("POST "+base+"/evaluations/draft", jwt(h.SaveDraftEvaluation))
	mux.Handle("POST "+base+"/evaluations", jwt(h.AddEvaluation))
	mux.Handle("POST "+base+"/evaluations/submit-draft", jwt(h.SubmitDraftEvaluation))
	mux.Handle("POST "+base+"/evaluations/approve", jwt(h.ApproveEvaluation))
	mux.Handle("POST "+base+"/evaluations/reject", jwt(h.RejectEvaluation))
	mux.Handle("GET "+base+"/evaluations", jwt(h.GetStaffEvaluations))

	// --- Feedback ---
	mux.Handle("POST "+base+"/feedback/request", jwt(h.RequestFeedback))
	mux.Handle("GET "+base+"/feedback/requests", jwt(h.GetFeedbackRequests))
	mux.Handle("POST "+base+"/feedback/process", jwt(h.ProcessFeedback))
	mux.Handle("GET "+base+"/feedback/pending", jwt(h.GetPendingFeedbackActions))

	// --- Scoring ---
	mux.Handle("GET "+base+"/scores", jwt(h.GetPerformanceScore))
	mux.Handle("GET "+base+"/dashboard", jwt(h.GetDashboardStats))
	mux.Handle("GET "+base+"/scores/summary", jwt(h.GetPerformanceSummary))

	// --- Period Objective Planning (Individual Objectives) ---
	mux.Handle("POST "+base+"/individual-objectives/draft", jwt(h.SaveDraftIndividualPlannedObjective))
	mux.Handle("POST "+base+"/individual-objectives", jwt(h.AddIndividualPlannedObjective))
	mux.Handle("POST "+base+"/individual-objectives/submit-draft", jwt(h.SubmitDraftIndividualObjective))
	mux.Handle("POST "+base+"/individual-objectives/approve", jwt(h.ApproveIndividualObjective))
	mux.Handle("POST "+base+"/individual-objectives/reject", jwt(h.RejectIndividualObjective))
	mux.Handle("POST "+base+"/individual-objectives/return", jwt(h.ReturnIndividualObjective))
	mux.Handle("POST "+base+"/individual-objectives/cancel", jwt(h.CancelIndividualObjective))
	mux.Handle("GET "+base+"/individual-objectives", jwt(h.GetStaffIndividualObjectives))

	// --- 360 Review ---
	mux.Handle("POST "+base+"/360-review/trigger", jwt(h.Trigger360Review))
	mux.Handle("POST "+base+"/360-review/initiate", jwt(h.Initiate360Review))
	mux.Handle("POST "+base+"/360-review/complete", jwt(h.Complete360ReviewForStaff))
	mux.Handle("POST "+base+"/360-review/rating", jwt(h.Add360Rating))
	mux.Handle("PUT "+base+"/360-review/rating", jwt(h.Update360Rating))
	mux.Handle("POST "+base+"/360-review/reviewer-complete", jwt(h.ReviewerComplete360Review))

	// --- Competency Review ---
	mux.Handle("GET "+base+"/competency-review/{feedbackId}/details", jwt(h.GetCompetencyReviewFeedbackDetails))
	mux.Handle("GET "+base+"/competency-review/{feedbackId}", jwt(h.GetCompetencyReviewDetail))
	mux.Handle("GET "+base+"/competency-review/feedbacks", jwt(h.GetAllCompetencyReviewFeedbacksByReviewPeriod))
	mux.Handle("GET "+base+"/competency-review/my-reviewed", jwt(h.GetAllMyReviewedCompetencies))
	mux.Handle("GET "+base+"/competency-review/to-review", jwt(h.GetCompetenciesToReview))
	mux.Handle("GET "+base+"/competency-review/reviewer/{reviewerId}", jwt(h.GetReviewerFeedbackDetails))
	mux.Handle("GET "+base+"/competency-review/questionnaire", jwt(h.GetQuestionnaire))
	mux.Handle("POST "+base+"/competency-review/gap-closure", jwt(h.CompetencyGapClosureSetup))

	// --- Feedback Requests (Extended) ---
	mux.Handle("GET "+base+"/feedback/requests/staff", jwt(h.GetStaffRequests))
	mux.Handle("GET "+base+"/feedback/requests/breached", jwt(h.GetBreachedRequests))
	mux.Handle("GET "+base+"/feedback/requests/staff/by-status", jwt(h.GetStaffRequestsByStatus))
	mux.Handle("GET "+base+"/feedback/requests/all", jwt(h.GetAllRequests))
	mux.Handle("GET "+base+"/feedback/requests/by-status", jwt(h.GetRequestsByStatus))
	mux.Handle("GET "+base+"/feedback/requests/{requestId}", jwt(h.GetRequestDetails))
	mux.Handle("POST "+base+"/feedback/requests/reassign", jwt(h.ReassignRequest))
	mux.Handle("POST "+base+"/feedback/requests/reassign-self", jwt(h.ReassignSelfRequest))
	mux.Handle("POST "+base+"/feedback/requests/close", jwt(h.CloseRequest))
	mux.Handle("POST "+base+"/feedback/requests/treat", jwt(h.TreatAssignedRequest))

	// --- Dashboard & Statistics ---
	mux.Handle("GET "+base+"/stats/requests", jwt(h.GetRequestStatistics))
	mux.Handle("GET "+base+"/stats/performance", jwt(h.GetStaffPerformanceStatistics))
	mux.Handle("GET "+base+"/stats/work-products", jwt(h.GetStaffWorkProductsStatistics))
	mux.Handle("GET "+base+"/stats/work-products-details", jwt(h.GetStaffWorkProductsDetailsStatistics))

	// --- ScoreCard Statistics ---
	mux.Handle("GET "+base+"/scorecard", jwt(h.GetStaffPerformanceScoreCardStatistics))
	mux.Handle("GET "+base+"/scorecard/annual", jwt(h.GetStaffAnnualPerformanceScoreCardStatistics))
	mux.Handle("GET "+base+"/scorecard/subordinates", jwt(h.GetSubordinatesStaffPerformanceScoreCardStatistics))

	// --- Organogram Performance ---
	mux.Handle("GET "+base+"/organogram-performance/list", jwt(h.GetOrganogramPerformanceSummaryListStatistics))
	mux.Handle("GET "+base+"/organogram-performance", jwt(h.GetOrganogramPerformanceSummaryStatistics))

	// --- Period Scores ---
	mux.Handle("GET "+base+"/period-scores/all", jwt(h.GetPeriodScores))
	mux.Handle("GET "+base+"/period-scores", jwt(h.GetPeriodScoreDetails))
	mux.Handle("GET "+base+"/staff-review-periods", jwt(h.GetStaffReviewPeriods))

	// --- Audit Logs ---
	mux.Handle("GET "+base+"/audit-logs/{id}", jwt(h.GetAuditLogDetails))
	mux.Handle("GET "+base+"/audit-logs", jwt(h.GetAuditLogs))

	// --- Line Manager & Staff ---
	mux.Handle("GET "+base+"/line-manager-employees", jwt(h.GetLineManagerEmployees))
	mux.Handle("GET "+base+"/adhoc-employees", jwt(h.GetAdhocAssignmentEmployees))
	mux.Handle("GET "+base+"/my-staff", jwt(h.GetMyStaff))

	// --- Password Management (AllowAnonymous) ---
	mux.Handle("POST "+base+"/reset-password", http.HandlerFunc(h.ResetUserPassword))
}
