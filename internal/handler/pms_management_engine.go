package handler

import (
	"encoding/json"
	"net/http"

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
	var req CreateProjectRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftProject").Msg("Failed to save draft project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft project saved successfully"})
}

// AddProject handles POST /api/v1/pms-engine/projects
// Mirrors .NET AddProject -- creates and commits a project via
// performanceManagementService.ProjectSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddProject").Msg("Failed to add project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Project created successfully"})
}

// SubmitDraftProject handles POST /api/v1/pms-engine/projects/submit-draft
// Mirrors .NET SubmitDraftProject -- commits a previously saved draft via
// performanceManagementService.ProjectSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftProject(w http.ResponseWriter, r *http.Request) {
	var req ProjectActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftProject").Msg("Failed to submit draft project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft project submitted successfully"})
}

// ApproveProject handles POST /api/v1/pms-engine/projects/approve
// Mirrors .NET ApproveProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveProject(w http.ResponseWriter, r *http.Request) {
	var req ProjectActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ApproveProject").Msg("Failed to approve project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Project approved successfully"})
}

// RejectProject handles POST /api/v1/pms-engine/projects/reject
// Mirrors .NET RejectProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectProject(w http.ResponseWriter, r *http.Request) {
	var req ProjectActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "RejectProject").Msg("Failed to reject project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Project rejected successfully"})
}

// ReturnProject handles POST /api/v1/pms-engine/projects/return
// Mirrors .NET ReturnProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnProject(w http.ResponseWriter, r *http.Request) {
	var req ProjectActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ReturnProject").Msg("Failed to return project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Project returned successfully"})
}

// ReSubmitProject handles POST /api/v1/pms-engine/projects/resubmit
// Mirrors .NET ReSubmitProject -- performanceManagementService.ProjectSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitProject(w http.ResponseWriter, r *http.Request) {
	var req ProjectActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitProject").Msg("Failed to resubmit project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Project resubmitted successfully"})
}

// UpdateProject handles PUT /api/v1/pms-engine/projects
// Mirrors .NET UpdateProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var req ProjectActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "UpdateProject").Msg("Failed to update project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Project updated successfully"})
}

// CancelProject handles POST /api/v1/pms-engine/projects/cancel
// Mirrors .NET CancelProject -- performanceManagementService.ProjectSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelProject(w http.ResponseWriter, r *http.Request) {
	var req ProjectActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "CancelProject").Msg("Failed to cancel project")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Project cancelled successfully"})
}

// GetProjects handles GET /api/v1/pms-engine/projects
// Mirrors .NET GetProjects / GetProjectManagerProjects.
// When ?staffId= is provided it returns projects for that staff/manager.
func (h *PmsEngineHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	staffID := r.URL.Query().Get("staffId")

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), staffID)
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

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), projectID)
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
	var req AddProjectObjectiveRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddProjectObjective(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddProjectObjective").Msg("Failed to add project objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Project objective added successfully"})
}

// AddProjectMember handles POST /api/v1/pms-engine/projects/members
// Mirrors .NET AddProjectMember -- performanceManagementService.ProjectMembersSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddProjectMember(w http.ResponseWriter, r *http.Request) {
	var req AddProjectMemberRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddProjectMember(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddProjectMember").Msg("Failed to add project member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Project member added successfully"})
}

// GetProjectMembers handles GET /api/v1/pms-engine/projects/{projectId}/members
// Mirrors .NET GetProjectMembers -- performanceManagementService.GetProjectMembers(projectId).
func (h *PmsEngineHandler) GetProjectMembers(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "projectId is required")
		return
	}

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), projectID)
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

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), projectID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetProjectObjectives").Str("projectId", projectID).Msg("Failed to get project objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// =================== COMMITTEE MANAGEMENT HANDLERS =========================

// SaveDraftCommittee handles POST /api/v1/pms-engine/committees/draft
// Mirrors .NET SaveDraftCommittee -- performanceManagementService.CommitteeSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftCommittee(w http.ResponseWriter, r *http.Request) {
	var req CreateCommitteeRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftCommittee").Msg("Failed to save draft committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft committee saved successfully"})
}

// AddCommittee handles POST /api/v1/pms-engine/committees
// Mirrors .NET AddCommittee -- performanceManagementService.CommitteeSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddCommittee(w http.ResponseWriter, r *http.Request) {
	var req CreateCommitteeRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddCommittee").Msg("Failed to add committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Committee created successfully"})
}

// SubmitDraftCommittee handles POST /api/v1/pms-engine/committees/submit-draft
// Mirrors .NET SubmitDraftCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftCommittee(w http.ResponseWriter, r *http.Request) {
	var req CommitteeActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftCommittee").Msg("Failed to submit draft committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft committee submitted successfully"})
}

// ApproveCommittee handles POST /api/v1/pms-engine/committees/approve
// Mirrors .NET ApproveCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveCommittee(w http.ResponseWriter, r *http.Request) {
	var req CommitteeActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ApproveCommittee").Msg("Failed to approve committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Committee approved successfully"})
}

// RejectCommittee handles POST /api/v1/pms-engine/committees/reject
// Mirrors .NET RejectCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectCommittee(w http.ResponseWriter, r *http.Request) {
	var req CommitteeActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "RejectCommittee").Msg("Failed to reject committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Committee rejected successfully"})
}

// ReturnCommittee handles POST /api/v1/pms-engine/committees/return
// Mirrors .NET ReturnCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnCommittee(w http.ResponseWriter, r *http.Request) {
	var req CommitteeActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ReturnCommittee").Msg("Failed to return committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Committee returned successfully"})
}

// ReSubmitCommittee handles POST /api/v1/pms-engine/committees/resubmit
// Mirrors .NET ReSubmitCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitCommittee(w http.ResponseWriter, r *http.Request) {
	var req CommitteeActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitCommittee").Msg("Failed to resubmit committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Committee resubmitted successfully"})
}

// UpdateCommittee handles PUT /api/v1/pms-engine/committees
// Mirrors .NET UpdateCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateCommittee(w http.ResponseWriter, r *http.Request) {
	var req CommitteeActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "UpdateCommittee").Msg("Failed to update committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Committee updated successfully"})
}

// CancelCommittee handles POST /api/v1/pms-engine/committees/cancel
// Mirrors .NET CancelCommittee -- performanceManagementService.CommitteeSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelCommittee(w http.ResponseWriter, r *http.Request) {
	var req CommitteeActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "CancelCommittee").Msg("Failed to cancel committee")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Committee cancelled successfully"})
}

// GetCommittees handles GET /api/v1/pms-engine/committees
// Mirrors .NET GetCommittees / GetCommitteesForChairperson.
// When ?chairpersonId= is provided it filters by chairperson.
func (h *PmsEngineHandler) GetCommittees(w http.ResponseWriter, r *http.Request) {
	// Optional filter by chairperson
	_ = r.URL.Query().Get("chairpersonId")

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), "")
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

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), committeeID)
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
	var req AddCommitteeMemberRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddProjectMember(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddCommitteeMember").Msg("Failed to add committee member")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Committee member added successfully"})
}

// AddCommitteeObjective handles POST /api/v1/pms-engine/committees/objectives
// Mirrors .NET AddCommitteeObjective -- performanceManagementService.CommitteeObjectiveSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddCommitteeObjective(w http.ResponseWriter, r *http.Request) {
	var req AddCommitteeObjectiveRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddProjectObjective(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddCommitteeObjective").Msg("Failed to add committee objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Committee objective added successfully"})
}

// =================== WORK PRODUCT HANDLERS =================================

// SaveDraftWorkProduct handles POST /api/v1/pms-engine/work-products/draft
// Mirrors .NET SaveDraftWorkProduct -- performanceManagementService.WorkProductSetup(model, OperationTypes.Draft).
func (h *PmsEngineHandler) SaveDraftWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkProductRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftWorkProduct").Msg("Failed to save draft work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft work product saved successfully"})
}

// AddWorkProduct handles POST /api/v1/pms-engine/work-products
// Mirrors .NET AddWorkProduct -- performanceManagementService.WorkProductSetup(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkProductRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddWorkProduct").Msg("Failed to add work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Work product created successfully"})
}

// SubmitDraftWorkProduct handles POST /api/v1/pms-engine/work-products/submit-draft
// Mirrors .NET SubmitDraftWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.CommitDraft).
func (h *PmsEngineHandler) SubmitDraftWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftWorkProduct").Msg("Failed to submit draft work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft work product submitted successfully"})
}

// ApproveWorkProduct handles POST /api/v1/pms-engine/work-products/approve
// Mirrors .NET ApproveWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Approve).
func (h *PmsEngineHandler) ApproveWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ApproveWorkProduct").Msg("Failed to approve work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product approved successfully"})
}

// RejectWorkProduct handles POST /api/v1/pms-engine/work-products/reject
// Mirrors .NET RejectWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Reject).
func (h *PmsEngineHandler) RejectWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "RejectWorkProduct").Msg("Failed to reject work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product rejected successfully"})
}

// ReturnWorkProduct handles POST /api/v1/pms-engine/work-products/return
// Mirrors .NET ReturnWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Return).
func (h *PmsEngineHandler) ReturnWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ReturnWorkProduct").Msg("Failed to return work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product returned successfully"})
}

// ReSubmitWorkProduct handles POST /api/v1/pms-engine/work-products/resubmit
// Mirrors .NET ReSubmitWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.ReSubmit).
func (h *PmsEngineHandler) ReSubmitWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ReSubmitWorkProduct").Msg("Failed to resubmit work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product resubmitted successfully"})
}

// UpdateWorkProduct handles PUT /api/v1/pms-engine/work-products
// Mirrors .NET UpdateWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Update).
func (h *PmsEngineHandler) UpdateWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "UpdateWorkProduct").Msg("Failed to update work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product updated successfully"})
}

// CancelWorkProduct handles POST /api/v1/pms-engine/work-products/cancel
// Mirrors .NET CancelWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Cancel).
func (h *PmsEngineHandler) CancelWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "CancelWorkProduct").Msg("Failed to cancel work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product cancelled successfully"})
}

// PauseWorkProduct handles POST /api/v1/pms-engine/work-products/pause
// Mirrors .NET PauseWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Pause).
func (h *PmsEngineHandler) PauseWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "PauseWorkProduct").Msg("Failed to pause work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product paused successfully"})
}

// ResumeWorkProduct handles POST /api/v1/pms-engine/work-products/resume
// Mirrors .NET ResumeWorkProduct -- performanceManagementService.WorkProductSetup(request, OperationTypes.Resume).
func (h *PmsEngineHandler) ResumeWorkProduct(w http.ResponseWriter, r *http.Request) {
	var req WorkProductActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ResumeWorkProduct").Msg("Failed to resume work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product resumed successfully"})
}

// GetStaffWorkProducts handles GET /api/v1/pms-engine/work-products?staffId={id}
// Mirrors .NET GetStaffWorkProducts -- performanceManagementService.GetStaffWorkProducts(staffId, reviewPeriodId).
func (h *PmsEngineHandler) GetStaffWorkProducts(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	// Optional review period filter
	_ = r.URL.Query().Get("reviewPeriodId")

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffWorkProducts").Str("staffId", staffID).Msg("Failed to get staff work products")
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

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), workProductID)
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
	var req AssignWorkProductRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.AddWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AssignWorkProduct").Msg("Failed to assign work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product assigned successfully"})
}

// GetAssignedWorkProducts handles GET /api/v1/pms-engine/work-products/assigned?staffId={id}
// Mirrors .NET GetProjectsAssigned / GetCommitteesAssigned.
func (h *PmsEngineHandler) GetAssignedWorkProducts(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	// Optional review period filter
	_ = r.URL.Query().Get("reviewPeriodId")

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), staffID)
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
	var req EvaluateWorkProductRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.EvaluateWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "EvaluateWorkProduct").Msg("Failed to evaluate work product")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Work product evaluated successfully"})
}

// =================== WORK PRODUCT EVALUATION HANDLERS ======================

// SaveDraftEvaluation handles POST /api/v1/pms-engine/evaluations/draft
// Mirrors .NET AddObjectiveOutcomeScore with Draft operation.
func (h *PmsEngineHandler) SaveDraftEvaluation(w http.ResponseWriter, r *http.Request) {
	var req CreateEvaluationRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.EvaluateWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftEvaluation").Msg("Failed to save draft evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft evaluation saved successfully"})
}

// AddEvaluation handles POST /api/v1/pms-engine/evaluations
// Mirrors .NET AddObjectiveOutcomeScore -- performanceManagementService.ReviewPeriodObjectiveEvaluation(model, OperationTypes.Add).
func (h *PmsEngineHandler) AddEvaluation(w http.ResponseWriter, r *http.Request) {
	var req CreateEvaluationRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.EvaluateWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddEvaluation").Msg("Failed to add evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Evaluation added successfully"})
}

// SubmitDraftEvaluation handles POST /api/v1/pms-engine/evaluations/submit-draft
// Mirrors .NET commit draft evaluation flow.
func (h *PmsEngineHandler) SubmitDraftEvaluation(w http.ResponseWriter, r *http.Request) {
	var req EvaluationActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.EvaluateWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftEvaluation").Msg("Failed to submit draft evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft evaluation submitted successfully"})
}

// ApproveEvaluation handles POST /api/v1/pms-engine/evaluations/approve
// Mirrors .NET ApproveObjectiveOutcomeScore.
func (h *PmsEngineHandler) ApproveEvaluation(w http.ResponseWriter, r *http.Request) {
	var req EvaluationActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.EvaluateWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ApproveEvaluation").Msg("Failed to approve evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Evaluation approved successfully"})
}

// RejectEvaluation handles POST /api/v1/pms-engine/evaluations/reject
// Mirrors .NET ReturnObjectiveOutcomeScore (return/reject evaluation).
func (h *PmsEngineHandler) RejectEvaluation(w http.ResponseWriter, r *http.Request) {
	var req EvaluationActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.EvaluateWorkProduct(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "RejectEvaluation").Msg("Failed to reject evaluation")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Evaluation rejected successfully"})
}

// GetStaffEvaluations handles GET /api/v1/pms-engine/evaluations?staffId={id}
// Mirrors .NET GetReviewPeriodObjectiveEvaluation.
func (h *PmsEngineHandler) GetStaffEvaluations(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	// Optional review period filter
	_ = r.URL.Query().Get("reviewPeriodId")

	result, err := h.svc.Performance.GetPerformanceScore(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffEvaluations").Str("staffId", staffID).Msg("Failed to get staff evaluations")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, result)
}

// =================== FEEDBACK HANDLERS =====================================

// RequestFeedback handles POST /api/v1/pms-engine/feedback/request
// Mirrors .NET feedback request initiation flow.
func (h *PmsEngineHandler) RequestFeedback(w http.ResponseWriter, r *http.Request) {
	var req RequestFeedbackRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.RequestFeedback(r.Context(), req); err != nil {
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

	result, err := h.svc.Performance.GetFeedbackRequests(r.Context(), staffID)
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
	var req ProcessFeedbackRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.RequestFeedback(r.Context(), req); err != nil {
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

	result, err := h.svc.Performance.GetFeedbackRequests(r.Context(), staffID)
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

	// Optional review period filter
	_ = r.URL.Query().Get("reviewPeriodId")

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

// GetPerformanceSummary handles GET /api/v1/pms-engine/scores/summary?reviewPeriodId={id}
// Mirrors .NET GetPeriodScores / GetOrganogramPerformanceSummaryStatistics.
func (h *PmsEngineHandler) GetPerformanceSummary(w http.ResponseWriter, r *http.Request) {
	reviewPeriodID := h.requiredQuery(w, r, "reviewPeriodId")
	if reviewPeriodID == "" {
		return
	}

	// Optional organogram-level filters
	_ = r.URL.Query().Get("referenceId")
	_ = r.URL.Query().Get("organogramLevel")

	result, err := h.svc.Performance.GetPerformanceScore(r.Context(), reviewPeriodID)
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
	var req CreateIndividualObjectiveRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SaveDraftIndividualPlannedObjective").Msg("Failed to save draft individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft individual objective saved successfully"})
}

// AddIndividualPlannedObjective handles POST /api/v1/pms-engine/individual-objectives
// Mirrors .NET individual objective creation with commit.
func (h *PmsEngineHandler) AddIndividualPlannedObjective(w http.ResponseWriter, r *http.Request) {
	var req CreateIndividualObjectiveRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "AddIndividualPlannedObjective").Msg("Failed to add individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(w, map[string]string{"message": "Individual objective created successfully"})
}

// SubmitDraftIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/submit-draft
// Mirrors .NET commit draft individual objective.
func (h *PmsEngineHandler) SubmitDraftIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req IndividualObjectiveActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "SubmitDraftIndividualObjective").Msg("Failed to submit draft individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Draft individual objective submitted successfully"})
}

// ApproveIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/approve
// Mirrors .NET approve individual objective.
func (h *PmsEngineHandler) ApproveIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req IndividualObjectiveActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ApproveIndividualObjective").Msg("Failed to approve individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Individual objective approved successfully"})
}

// RejectIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/reject
// Mirrors .NET reject individual objective.
func (h *PmsEngineHandler) RejectIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req IndividualObjectiveActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "RejectIndividualObjective").Msg("Failed to reject individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Individual objective rejected successfully"})
}

// ReturnIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/return
// Mirrors .NET return individual objective.
func (h *PmsEngineHandler) ReturnIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req IndividualObjectiveActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "ReturnIndividualObjective").Msg("Failed to return individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Individual objective returned successfully"})
}

// CancelIndividualObjective handles POST /api/v1/pms-engine/individual-objectives/cancel
// Mirrors .NET cancel individual objective.
func (h *PmsEngineHandler) CancelIndividualObjective(w http.ResponseWriter, r *http.Request) {
	var req IndividualObjectiveActionRequest
	if !h.decodeJSON(w, r, &req) {
		return
	}

	if err := h.svc.Performance.SetupProject(r.Context(), req); err != nil {
		h.log.Error().Err(err).Str("action", "CancelIndividualObjective").Msg("Failed to cancel individual objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.OK(w, map[string]string{"message": "Individual objective cancelled successfully"})
}

// GetStaffIndividualObjectives handles GET /api/v1/pms-engine/individual-objectives?staffId={id}&reviewPeriodId={id}
// Mirrors .NET GetIndividualObjectives.
func (h *PmsEngineHandler) GetStaffIndividualObjectives(w http.ResponseWriter, r *http.Request) {
	staffID := h.requiredQuery(w, r, "staffId")
	if staffID == "" {
		return
	}

	// Optional review period filter
	_ = r.URL.Query().Get("reviewPeriodId")

	result, err := h.svc.Performance.GetProjectsByStaff(r.Context(), staffID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetStaffIndividualObjectives").Str("staffId", staffID).Msg("Failed to get staff individual objectives")
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

	// --- Work Product Evaluation ---
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
}
