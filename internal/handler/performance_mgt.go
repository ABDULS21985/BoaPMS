package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// Request DTOs — inline structs mirroring the .NET ViewModels
// ---------------------------------------------------------------------------

// CreateStrategyRequest mirrors CreateNewStrategyVm.
type CreateStrategyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BankYearID  int    `json:"bankYearId"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	IsActive    *bool  `json:"isActive,omitempty"`
}

// UpdateStrategyRequest mirrors StrategyVm (update path).
type UpdateStrategyRequest struct {
	StrategyID  string `json:"strategyId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BankYearID  int    `json:"bankYearId"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

// CreateStrategicThemeRequest mirrors CreateStrategicThemeVm.
type CreateStrategicThemeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StrategyID  string `json:"strategyId"`
}

// UpdateStrategicThemeRequest mirrors StrategicThemeVm (update path).
type UpdateStrategicThemeRequest struct {
	StrategicThemeID string `json:"strategicThemeId"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	StrategyID       string `json:"strategyId"`
	StrategyName     string `json:"strategyName,omitempty"`
}

// CreateEnterpriseObjectiveRequest mirrors CreateEnterpriseObjectiveVm.
type CreateEnterpriseObjectiveRequest struct {
	Name                         string `json:"name"`
	Description                  string `json:"description,omitempty"`
	Kpi                          string `json:"kpi"`
	Target                       string `json:"target"`
	EnterpriseObjectivesCategoryID string `json:"enterpriseObjectivesCategoryId"`
	StrategyID                   string `json:"strategyId"`
}

// UpdateEnterpriseObjectiveRequest mirrors EnterpriseObjectiveVm (update path).
type UpdateEnterpriseObjectiveRequest struct {
	EnterpriseObjectiveID          string `json:"enterpriseObjectiveId"`
	Name                           string `json:"name"`
	Description                    string `json:"description,omitempty"`
	Kpi                            string `json:"kpi"`
	Target                         string `json:"target"`
	EnterpriseObjectivesCategoryID string `json:"enterpriseObjectivesCategoryId"`
	StrategyID                     string `json:"strategyId"`
}

// CreateDepartmentObjectiveRequest mirrors CreateDepartmentObjectiveVm.
type CreateDepartmentObjectiveRequest struct {
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	Kpi                    string `json:"kpi,omitempty"`
	Target                 string `json:"target,omitempty"`
	DepartmentID           int    `json:"departmentId"`
	EnterpriseObjectiveID  string `json:"enterpriseObjectiveId,omitempty"`
	SBUName                string `json:"sbuName,omitempty"`
	WorkProductName        string `json:"workProductName,omitempty"`
	WorkProductDescription string `json:"workProductDescription,omitempty"`
	WorkProductDeliverable string `json:"workProductDeliverable,omitempty"`
	JobGradeGroup          string `json:"jobGradeGroup,omitempty"`
}

// UpdateDepartmentObjectiveRequest mirrors DepartmentObjectiveVm (update path).
type UpdateDepartmentObjectiveRequest struct {
	DepartmentObjectiveID  string `json:"departmentObjectiveId"`
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	Kpi                    string `json:"kpi,omitempty"`
	Target                 string `json:"target,omitempty"`
	DepartmentID           int    `json:"departmentId"`
	EnterpriseObjectiveID  string `json:"enterpriseObjectiveId,omitempty"`
	SBUName                string `json:"sbuName,omitempty"`
	WorkProductName        string `json:"workProductName,omitempty"`
	WorkProductDescription string `json:"workProductDescription,omitempty"`
	WorkProductDeliverable string `json:"workProductDeliverable,omitempty"`
	JobGradeGroup          string `json:"jobGradeGroup,omitempty"`
}

// CreateDivisionObjectiveRequest mirrors CreateDivisionObjectiveVm.
type CreateDivisionObjectiveRequest struct {
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	Kpi                    string `json:"kpi,omitempty"`
	Target                 string `json:"target,omitempty"`
	DivisionID             int    `json:"divisionId"`
	DepartmentObjectiveID  string `json:"departmentObjectiveId"`
	DepartmentID           int    `json:"departmentId"`
	JobGradeGroup          string `json:"jobGradeGroup,omitempty"`
	SBUName                string `json:"sbuName,omitempty"`
	WorkProductName        string `json:"workProductName,omitempty"`
	WorkProductDescription string `json:"workProductDescription,omitempty"`
	WorkProductDeliverable string `json:"workProductDeliverable,omitempty"`
}

// UpdateDivisionObjectiveRequest mirrors DivisionObjectiveVm (update path).
type UpdateDivisionObjectiveRequest struct {
	DivisionObjectiveID    string `json:"divisionObjectiveId"`
	Name                   string `json:"name"`
	Description            string `json:"description,omitempty"`
	Kpi                    string `json:"kpi,omitempty"`
	Target                 string `json:"target,omitempty"`
	DivisionID             int    `json:"divisionId"`
	DepartmentObjectiveID  string `json:"departmentObjectiveId"`
	DepartmentID           int    `json:"departmentId"`
	JobGradeGroup          string `json:"jobGradeGroup,omitempty"`
	SBUName                string `json:"sbuName,omitempty"`
	WorkProductName        string `json:"workProductName,omitempty"`
	WorkProductDescription string `json:"workProductDescription,omitempty"`
	WorkProductDeliverable string `json:"workProductDeliverable,omitempty"`
}

// CreateOfficeObjectiveRequest mirrors CreateOfficeObjectiveVm.
type CreateOfficeObjectiveRequest struct {
	Name                          string `json:"name"`
	Description                   string `json:"description,omitempty"`
	Kpi                           string `json:"kpi,omitempty"`
	Target                        string `json:"target,omitempty"`
	OfficeID                      int    `json:"officeId"`
	DivisionObjectiveID           string `json:"divisionObjectiveId,omitempty"`
	JobGradeGroupID               int    `json:"jobGradeGroupId"`
	JobGradeGroupName             string `json:"jobGradeGroupName,omitempty"`
	SBUName                       string `json:"sbuName,omitempty"`
	WorkProductName               string `json:"workProductName,omitempty"`
	WorkProductDescription        string `json:"workProductDescription,omitempty"`
	WorkProductDeliverable        string `json:"workProductDeliverable,omitempty"`
	ParentDivisionID              int    `json:"parentDivisionId,omitempty"`
	ParentDivisionObjectiveName   string `json:"parentDivisionObjectiveName,omitempty"`
}

// UpdateOfficeObjectiveRequest mirrors OfficeObjectiveVm (update path).
type UpdateOfficeObjectiveRequest struct {
	OfficeObjectiveID             string `json:"officeObjectiveId"`
	Name                          string `json:"name"`
	Description                   string `json:"description,omitempty"`
	Kpi                           string `json:"kpi,omitempty"`
	Target                        string `json:"target,omitempty"`
	OfficeID                      int    `json:"officeId"`
	DivisionObjectiveID           string `json:"divisionObjectiveId,omitempty"`
	JobGradeGroupID               int    `json:"jobGradeGroupId"`
	JobGradeGroupName             string `json:"jobGradeGroupName,omitempty"`
	SBUName                       string `json:"sbuName,omitempty"`
	WorkProductName               string `json:"workProductName,omitempty"`
	WorkProductDescription        string `json:"workProductDescription,omitempty"`
	WorkProductDeliverable        string `json:"workProductDeliverable,omitempty"`
	ParentDivisionID              int    `json:"parentDivisionId,omitempty"`
	ParentDivisionObjectiveName   string `json:"parentDivisionObjectiveName,omitempty"`
}

// CreateObjectiveCategoryRequest mirrors CreateObjectiveCategoryVm.
type CreateObjectiveCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateObjectiveCategoryRequest mirrors ObjectiveCategoryVm (update path).
type UpdateObjectiveCategoryRequest struct {
	ObjectiveCategoryID string `json:"objectiveCategoryId"`
	Name                string `json:"name"`
	Description         string `json:"description,omitempty"`
}

// CreateCategoryDefinitionRequest mirrors CreateCategoryDefinitionVm.
type CreateCategoryDefinitionRequest struct {
	ObjectiveCategoryID    string  `json:"objectiveCategoryId"`
	Weight                 float64 `json:"weight"`
	MaxNoObjectives        int     `json:"maxNoObjectives"`
	MaxNoWorkProduct       int     `json:"maxNoWorkProduct"`
	MaxPoints              int     `json:"maxPoints"`
	IsCompulsory           bool    `json:"isCompulsory"`
	EnforceWorkProductLimit bool   `json:"enforceWorkProductLimit"`
	Description            string  `json:"description"`
	GradeGroupID           int     `json:"gradeGroupId"`
}

// UpdateCategoryDefinitionRequest mirrors CategoryDefinitionVm (update path).
type UpdateCategoryDefinitionRequest struct {
	DefinitionID           string  `json:"definitionId"`
	ObjectiveCategoryID    string  `json:"objectiveCategoryId"`
	Weight                 float64 `json:"weight"`
	MaxNoObjectives        int     `json:"maxNoObjectives"`
	MaxNoWorkProduct       int     `json:"maxNoWorkProduct"`
	MaxPoints              int     `json:"maxPoints"`
	IsCompulsory           bool    `json:"isCompulsory"`
	EnforceWorkProductLimit bool   `json:"enforceWorkProductLimit"`
	Description            string  `json:"description"`
	GradeGroupID           int     `json:"gradeGroupId"`
	GroupName              string  `json:"groupName,omitempty"`
	Category               string  `json:"category,omitempty"`
}

// CreatePmsCompetencyRequest mirrors CreateNewPmsCompetencyVm.
type CreatePmsCompetencyRequest struct {
	PmsCompetencyID  string `json:"pmsCompetencyId"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	ObjectCategoryID string `json:"objectCategoryId"`
}

// UpdatePmsCompetencyRequest mirrors PmsCompetencyRequestVm (update path).
type UpdatePmsCompetencyRequest struct {
	PmsCompetencyID  string `json:"pmsCompetencyId"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	ObjectCategoryID string `json:"objectCategoryId"`
	RecordStatus     int    `json:"recordStatus"`
}

// EvaluationOptionRequest mirrors EvaluationOptionVm.
type EvaluationOptionRequest struct {
	EvaluationOptionID string  `json:"evaluationOptionId,omitempty"`
	Name               string  `json:"name"`
	Description        string  `json:"description,omitempty"`
	RecordStatus       int     `json:"recordStatus"`
	Score              float64 `json:"score"`
	EvaluationType     int     `json:"evaluationType"`
}

// FeedbackQuestionnaireRequest mirrors FeedbackQuestionaireVm.
type FeedbackQuestionnaireRequest struct {
	FeedbackQuestionnaireID string                            `json:"feedbackQuestionnaireId,omitempty"`
	Question                string                            `json:"question"`
	Description             string                            `json:"description,omitempty"`
	PmsCompetencyID         string                            `json:"pmsCompetencyId,omitempty"`
	PmsCompetencyName       string                            `json:"pmsCompetencyName,omitempty"`
	RecordStatus            int                               `json:"recordStatus"`
	Category                string                            `json:"category,omitempty"`
	QuestionID              string                            `json:"questionId,omitempty"`
	OptionStatement         string                            `json:"optionStatement,omitempty"`
	Score                   float64                           `json:"score"`
	Options                 []FeedbackQuestionnaireOptionItem `json:"options,omitempty"`
}

// FeedbackQuestionnaireOptionItem mirrors FeedbackQuestionaireOptionVm.
type FeedbackQuestionnaireOptionItem struct {
	FeedbackQuestionnaireOptionID string  `json:"feedbackQuestionnaireOptionId,omitempty"`
	OptionStatement               string  `json:"optionStatement"`
	Description                   string  `json:"description,omitempty"`
	Score                         float64 `json:"score"`
	QuestionID                    string  `json:"questionId,omitempty"`
}

// WorkProductDefinitionRequest mirrors WorkProductDefinitionRequestVm.
type WorkProductDefinitionRequest struct {
	ReferenceNo    string `json:"referenceNo,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Deliverables   string `json:"deliverables,omitempty"`
	ObjectiveID    string `json:"objectiveId,omitempty"`
	ObjectiveLevel string `json:"objectiveLevel,omitempty"`
	SBUName        string `json:"sbuName,omitempty"`
	ObjectiveName  string `json:"objectiveName,omitempty"`
}

// CascadedObjectiveUploadRequest mirrors CascadedObjectiveUploadVm.
type CascadedObjectiveUploadRequest struct {
	StrategyID             string `json:"strategyId,omitempty"`
	StrategicThemeID       string `json:"strategicThemeId,omitempty"`
	EObjName               string `json:"eObjName,omitempty"`
	EObjDesc               string `json:"eObjDesc,omitempty"`
	EObjKPI                string `json:"eObjKPI,omitempty"`
	EObjTarget             string `json:"eObjTarget,omitempty"`
	EObjCategory           string `json:"eObjCategory,omitempty"`
	Dept                   string `json:"dept,omitempty"`
	DeptObjName            string `json:"deptObjName,omitempty"`
	DeptObjDesc            string `json:"deptObjDesc,omitempty"`
	DeptObjKPI             string `json:"deptObjKPI,omitempty"`
	DeptObjTarget          string `json:"deptObjTarget,omitempty"`
	Division               string `json:"division,omitempty"`
	DivObjName             string `json:"divObjName,omitempty"`
	DivObjDesc             string `json:"divObjDesc,omitempty"`
	DivObjKPI              string `json:"divObjKPI,omitempty"`
	DivObjTarget           string `json:"divObjTarget,omitempty"`
	Office                 string `json:"office,omitempty"`
	OffObjName             string `json:"offObjName,omitempty"`
	OffObjDesc             string `json:"offObjDesc,omitempty"`
	OffObjKPI              string `json:"offObjKPI,omitempty"`
	OffObjTarget           string `json:"offObjTarget,omitempty"`
	JobGradeGroup          string `json:"jobGradeGroup,omitempty"`
	SBULevel               string `json:"sbuLevel,omitempty"`
	WorkProductName        string `json:"workProductName,omitempty"`
	WorkProductDescription string `json:"workProductDescription,omitempty"`
	WorkProductDeliverable string `json:"workProductDeliverable,omitempty"`
}

// ConsolidatedObjectiveRequest mirrors ConsolidatedObjectiveVm for deactivate/reactivate.
type ConsolidatedObjectiveRequest struct {
	IsSelected bool `json:"isSelected"`
	// Embedded fields from ObjectiveBase would be here; kept generic for now.
	ObjectiveID    string `json:"objectiveId,omitempty"`
	ObjectiveLevel int    `json:"objectiveLevel,omitempty"`
	Name           string `json:"name,omitempty"`
}

// ApprovalRequest mirrors ApprovalRequestVm.
type ApprovalRequest struct {
	EntityType string   `json:"entityType"`
	RecordIDs  []string `json:"recordIds"`
}

// RejectionRequest mirrors RejectionRequestVm.
type RejectionRequest struct {
	EntityType      string   `json:"entityType"`
	RecordIDs       []string `json:"recordIds"`
	RejectionReason string   `json:"rejectionReason"`
}

// SearchObjectiveParams mirrors SearchObjectiveVm query parameters.
type SearchObjectiveParams struct {
	PageIndex       int    `json:"pageIndex"`
	PageSize        int    `json:"pageSize"`
	DepartmentID    *int   `json:"departmentId,omitempty"`
	DivisionID      *int   `json:"divisionId,omitempty"`
	OfficeID        *int   `json:"officeId,omitempty"`
	JobRoleID       *int   `json:"jobRoleId,omitempty"`
	SearchString    string `json:"searchString,omitempty"`
	IsApproved      bool   `json:"isApproved"`
	IsTechnical     bool   `json:"isTechnical"`
	TargetType      int    `json:"targetType"`
	TargetReference string `json:"targetReference,omitempty"`
	Status          string `json:"status,omitempty"`
}

// SelectListItem represents a key-value pair for enum select lists.
type SelectListItem struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

// ---------------------------------------------------------------------------
// Handler
// ---------------------------------------------------------------------------

// PerformanceMgtHandler handles all performance management HTTP endpoints.
// Mirrors the .NET PerformanceMgtController (main partial class).
type PerformanceMgtHandler struct {
	svc *service.Container
	log zerolog.Logger
}

// NewPerformanceMgtHandler creates a new performance management handler.
func NewPerformanceMgtHandler(svc *service.Container, log zerolog.Logger) *PerformanceMgtHandler {
	return &PerformanceMgtHandler{svc: svc, log: log}
}

// =========================================================================
// DATA LISTING — GET endpoints
// =========================================================================

// GetBankStrategies handles GET /api/v1/performance/strategies
// Mirrors .NET GetBankStrategies — returns all bank strategies.
func (h *PerformanceMgtHandler) GetBankStrategies(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetStrategies(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetBankStrategies").Msg("failed to get bank strategies")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetBankStrategicThemes handles GET /api/v1/performance/strategic-themes
// Mirrors .NET GetBankStrategicThemes — returns all strategic themes.
// Also handles GET /api/v1/performance/strategic-themes?strategyId={id}
// when strategyId query param is provided, mirrors .NET GetBankStrategicThemesById.
func (h *PerformanceMgtHandler) GetBankStrategicThemes(w http.ResponseWriter, r *http.Request) {
	strategyID := r.URL.Query().Get("strategyId")

	if strategyID != "" {
		result, err := h.svc.Performance.GetStrategicThemesById(r.Context(), strategyID)
		if err != nil {
			h.log.Error().Err(err).Str("action", "GetBankStrategicThemesById").Msg("failed to get strategic themes by strategy id")
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.OK(w, result)
		return
	}

	result, err := h.svc.Performance.GetStrategicThemes(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetBankStrategicThemes").Msg("failed to get strategic themes")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetEnterpriseObjectives handles GET /api/v1/performance/objectives/enterprise
// Mirrors .NET GetOrganisationalObjectives.
func (h *PerformanceMgtHandler) GetEnterpriseObjectives(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetEnterpriseObjectives(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetEnterpriseObjectives").Msg("failed to get enterprise objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetDepartmentObjectives handles GET /api/v1/performance/objectives/department
// Mirrors .NET GetDepartmentObjectives.
func (h *PerformanceMgtHandler) GetDepartmentObjectives(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetDepartmentObjectives(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetDepartmentObjectives").Msg("failed to get department objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetDivisionObjectives handles GET /api/v1/performance/objectives/division
// Mirrors .NET GetDivisionObjectives and GetDivisionObjectivesByDivisionId.
// When divisionId query param is provided, filters by division.
func (h *PerformanceMgtHandler) GetDivisionObjectives(w http.ResponseWriter, r *http.Request) {
	divisionIDStr := r.URL.Query().Get("divisionId")

	if divisionIDStr != "" {
		divisionID, err := strconv.Atoi(divisionIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid divisionId parameter")
			return
		}
		result, err := h.svc.Performance.GetDivisionObjectivesByDivisionId(r.Context(), divisionID)
		if err != nil {
			h.log.Error().Err(err).Str("action", "GetDivisionObjectivesByDivisionId").Msg("failed to get division objectives by division id")
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.OK(w, result)
		return
	}

	result, err := h.svc.Performance.GetDivisionObjectives(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetDivisionObjectives").Msg("failed to get division objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetOfficeObjectives handles GET /api/v1/performance/objectives/office
// Mirrors .NET GetOfficeObjectives and GetOfficeObjectivesByOfficeId.
// When officeId query param is provided, filters by office.
func (h *PerformanceMgtHandler) GetOfficeObjectives(w http.ResponseWriter, r *http.Request) {
	officeIDStr := r.URL.Query().Get("officeId")

	if officeIDStr != "" {
		officeID, err := strconv.Atoi(officeIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid officeId parameter")
			return
		}
		result, err := h.svc.Performance.GetOfficeObjectivesByOfficeId(r.Context(), officeID)
		if err != nil {
			h.log.Error().Err(err).Str("action", "GetOfficeObjectivesByOfficeId").Msg("failed to get office objectives by office id")
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.OK(w, result)
		return
	}

	result, err := h.svc.Performance.GetOfficeObjectives(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetOfficeObjectives").Msg("failed to get office objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetEvaluationOptions handles GET /api/v1/performance/evaluation-options
// Mirrors .NET GetEvaluationOptions.
func (h *PerformanceMgtHandler) GetEvaluationOptions(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetEvaluationOptions(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetEvaluationOptions").Msg("failed to get evaluation options")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetFeedbackQuestionnaires handles GET /api/v1/performance/feedback-questionnaires
// Mirrors .NET GetFeedbackQuestionaires.
func (h *PerformanceMgtHandler) GetFeedbackQuestionnaires(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetFeedbackQuestionnaires(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetFeedbackQuestionnaires").Msg("failed to get feedback questionnaires")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetPmsCompetencies handles GET /api/v1/performance/competencies
// Mirrors .NET GetPmsCompetencies.
func (h *PerformanceMgtHandler) GetPmsCompetencies(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetPmsCompetencies(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetPmsCompetencies").Msg("failed to get pms competencies")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetObjectiveCategories handles GET /api/v1/performance/objective-categories
// Mirrors .NET GetObjectiveCategories.
func (h *PerformanceMgtHandler) GetObjectiveCategories(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetObjectiveCategories(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetObjectiveCategories").Msg("failed to get objective categories")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetCategoryDefinitions handles GET /api/v1/performance/category-definitions?categoryId={id}
// Mirrors .NET GetCategoryDefinitions.
func (h *PerformanceMgtHandler) GetCategoryDefinitions(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("categoryId")
	if categoryID == "" {
		response.Error(w, http.StatusBadRequest, "categoryId query parameter is required")
		return
	}

	result, err := h.svc.Performance.GetCategoryDefinitions(r.Context(), categoryID)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetCategoryDefinitions").Msg("failed to get category definitions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetObjectiveWorkProductDefinitions handles GET /api/v1/performance/work-product-definitions?objectiveId={id}&objectiveLevel={level}
// Mirrors .NET GetObjectiveWorkProductDefinitions.
func (h *PerformanceMgtHandler) GetObjectiveWorkProductDefinitions(w http.ResponseWriter, r *http.Request) {
	objectiveID := r.URL.Query().Get("objectiveId")
	objectiveLevelStr := r.URL.Query().Get("objectiveLevel")

	if objectiveID == "" || objectiveLevelStr == "" {
		response.Error(w, http.StatusBadRequest, "objectiveId and objectiveLevel query parameters are required")
		return
	}

	objectiveLevel, err := strconv.Atoi(objectiveLevelStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid objectiveLevel parameter")
		return
	}

	result, err := h.svc.Performance.GetObjectiveWorkProductDefinitions(r.Context(), objectiveID, objectiveLevel)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetObjectiveWorkProductDefinitions").Msg("failed to get objective work product definitions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllWorkProductDefinitions handles GET /api/v1/performance/work-product-definitions/all
// Mirrors .NET GetAllWorkProductDefinitions.
func (h *PerformanceMgtHandler) GetAllWorkProductDefinitions(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetAllWorkProductDefinitions(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllWorkProductDefinitions").Msg("failed to get all work product definitions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetAllPaginatedWorkProductDefinitions handles GET /api/v1/performance/work-product-definitions?pageIndex={p}&pageSize={s}&search={q}
// Mirrors .NET GetAllPaginatedWorkProductDefinitions.
func (h *PerformanceMgtHandler) GetAllPaginatedWorkProductDefinitions(w http.ResponseWriter, r *http.Request) {
	pageIndexStr := r.URL.Query().Get("pageIndex")
	pageSizeStr := r.URL.Query().Get("pageSize")
	search := r.URL.Query().Get("search")

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		pageIndex = 0
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	result, err := h.svc.Performance.GetAllPaginatedWorkProductDefinitions(r.Context(), pageIndex, pageSize, search)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetAllPaginatedWorkProductDefinitions").Msg("failed to get paginated work product definitions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetConsolidatedObjectives handles GET /api/v1/performance/objectives/consolidated
// Mirrors .NET GetConsolidatedObjectives.
func (h *PerformanceMgtHandler) GetConsolidatedObjectives(w http.ResponseWriter, r *http.Request) {
	result, err := h.svc.Performance.GetConsolidatedObjectives(r.Context())
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetConsolidatedObjectives").Msg("failed to get consolidated objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// GetConsolidatedObjectivesPaginated handles GET /api/v1/performance/objectives/consolidated/paginated?...
// Mirrors .NET GetConsolidatedObjectivesPaginated with SearchObjectiveVm query params.
func (h *PerformanceMgtHandler) GetConsolidatedObjectivesPaginated(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	params := SearchObjectiveParams{
		SearchString:    q.Get("searchString"),
		TargetReference: q.Get("targetReference"),
		Status:          q.Get("status"),
	}

	if v := q.Get("pageIndex"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.PageIndex = n
		}
	}
	if v := q.Get("pageSize"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.PageSize = n
		}
	}
	if v := q.Get("departmentId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.DepartmentID = &n
		}
	}
	if v := q.Get("divisionId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.DivisionID = &n
		}
	}
	if v := q.Get("officeId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.OfficeID = &n
		}
	}
	if v := q.Get("jobRoleId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.JobRoleID = &n
		}
	}
	if v := q.Get("isApproved"); v == "true" {
		params.IsApproved = true
	}
	if v := q.Get("isTechnical"); v == "true" {
		params.IsTechnical = true
	}
	if v := q.Get("targetType"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.TargetType = n
		}
	}

	result, err := h.svc.Performance.GetConsolidatedObjectivesPaginated(r.Context(), params)
	if err != nil {
		h.log.Error().Err(err).Str("action", "GetConsolidatedObjectivesPaginated").Msg("failed to get paginated consolidated objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =========================================================================
// CRUD — POST / PUT endpoints
// =========================================================================

// CreateStrategy handles POST /api/v1/performance/strategies
// Mirrors .NET CreateNewStrategy.
func (h *PerformanceMgtHandler) CreateStrategy(w http.ResponseWriter, r *http.Request) {
	var req CreateStrategyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "Name is required")
		return
	}

	result, err := h.svc.Performance.CreateStrategy(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateStrategy").Msg("failed to create strategy")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateStrategy handles PUT /api/v1/performance/strategies
// Mirrors .NET UpdateStrategy.
func (h *PerformanceMgtHandler) UpdateStrategy(w http.ResponseWriter, r *http.Request) {
	var req UpdateStrategyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.StrategyID == "" {
		response.Error(w, http.StatusBadRequest, "strategyId is required")
		return
	}

	result, err := h.svc.Performance.UpdateStrategy(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateStrategy").Msg("failed to update strategy")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreateStrategicTheme handles POST /api/v1/performance/strategic-themes
// Mirrors .NET CreateNewStrategicTheme.
func (h *PerformanceMgtHandler) CreateStrategicTheme(w http.ResponseWriter, r *http.Request) {
	var req CreateStrategicThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Description == "" || req.StrategyID == "" {
		response.Error(w, http.StatusBadRequest, "name, description, and strategyId are required")
		return
	}

	result, err := h.svc.Performance.CreateStrategicTheme(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateStrategicTheme").Msg("failed to create strategic theme")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateStrategicTheme handles PUT /api/v1/performance/strategic-themes
// Mirrors .NET UpdateStrategicTheme.
func (h *PerformanceMgtHandler) UpdateStrategicTheme(w http.ResponseWriter, r *http.Request) {
	var req UpdateStrategicThemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.StrategicThemeID == "" {
		response.Error(w, http.StatusBadRequest, "strategicThemeId is required")
		return
	}

	result, err := h.svc.Performance.UpdateStrategicTheme(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateStrategicTheme").Msg("failed to update strategic theme")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// UploadObjectives handles POST /api/v1/performance/objectives/upload
// Mirrors .NET UploadObjectives — bulk upload of cascaded objectives.
func (h *PerformanceMgtHandler) UploadObjectives(w http.ResponseWriter, r *http.Request) {
	var req []CascadedObjectiveUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.Performance.ProcessObjectivesUpload(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UploadObjectives").Msg("failed to upload objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// DeActivateObjectives handles POST /api/v1/performance/objectives/deactivate
// Mirrors .NET DeActivateObjectives.
func (h *PerformanceMgtHandler) DeActivateObjectives(w http.ResponseWriter, r *http.Request) {
	var req []ConsolidatedObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.Performance.DeActivateOrReactivateObjectives(r.Context(), req, true)
	if err != nil {
		h.log.Error().Err(err).Str("action", "DeActivateObjectives").Msg("failed to deactivate objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// ReActivateObjectives handles POST /api/v1/performance/objectives/reactivate
// Mirrors .NET ReActivateObjectives.
func (h *PerformanceMgtHandler) ReActivateObjectives(w http.ResponseWriter, r *http.Request) {
	var req []ConsolidatedObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.Performance.DeActivateOrReactivateObjectives(r.Context(), req, false)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ReActivateObjectives").Msg("failed to reactivate objectives")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveEvaluationOptions handles POST /api/v1/performance/evaluation-options
// Mirrors .NET SaveEvaluationOptions — bulk save evaluation options.
func (h *PerformanceMgtHandler) SaveEvaluationOptions(w http.ResponseWriter, r *http.Request) {
	var req []EvaluationOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.Performance.SaveEvaluationOptions(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveEvaluationOptions").Msg("failed to save evaluation options")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveFeedbackQuestionnaires handles POST /api/v1/performance/feedback-questionnaires
// Mirrors .NET SaveFeedbackQuestionaires — bulk save feedback questionnaires.
func (h *PerformanceMgtHandler) SaveFeedbackQuestionnaires(w http.ResponseWriter, r *http.Request) {
	var req []FeedbackQuestionnaireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.Performance.SaveFeedbackQuestionnaires(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveFeedbackQuestionnaires").Msg("failed to save feedback questionnaires")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveWorkProductDefinition handles POST /api/v1/performance/work-product-definitions
// Mirrors .NET SaveWorkProductDefinition — bulk save work product definitions.
func (h *PerformanceMgtHandler) SaveWorkProductDefinition(w http.ResponseWriter, r *http.Request) {
	var req []WorkProductDefinitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.Performance.SaveWorkProductDefinitions(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveWorkProductDefinition").Msg("failed to save work product definitions")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// SaveFeedbackQuestionnaireOptions handles POST /api/v1/performance/feedback-questionnaire-options
// Mirrors .NET SaveFeedbackQuestionaireOptions — bulk save feedback questionnaire options.
func (h *PerformanceMgtHandler) SaveFeedbackQuestionnaireOptions(w http.ResponseWriter, r *http.Request) {
	var req []FeedbackQuestionnaireOptionItem
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.svc.Performance.SaveFeedbackQuestionnaireOptions(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "SaveFeedbackQuestionnaireOptions").Msg("failed to save feedback questionnaire options")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreateEnterpriseObjective handles POST /api/v1/performance/objectives/enterprise
// Mirrors .NET CreateEnterpriseObjective.
func (h *PerformanceMgtHandler) CreateEnterpriseObjective(w http.ResponseWriter, r *http.Request) {
	var req CreateEnterpriseObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Kpi == "" || req.Target == "" || req.EnterpriseObjectivesCategoryID == "" || req.StrategyID == "" {
		response.Error(w, http.StatusBadRequest, "name, kpi, target, enterpriseObjectivesCategoryId, and strategyId are required")
		return
	}

	result, err := h.svc.Performance.CreateEnterpriseObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateEnterpriseObjective").Msg("failed to create enterprise objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateEnterpriseObjective handles PUT /api/v1/performance/objectives/enterprise
// Mirrors .NET UpdateEnterpriseObjective.
func (h *PerformanceMgtHandler) UpdateEnterpriseObjective(w http.ResponseWriter, r *http.Request) {
	var req UpdateEnterpriseObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.EnterpriseObjectiveID == "" {
		response.Error(w, http.StatusBadRequest, "enterpriseObjectiveId is required")
		return
	}

	result, err := h.svc.Performance.UpdateEnterpriseObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateEnterpriseObjective").Msg("failed to update enterprise objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreatePmsCompetency handles POST /api/v1/performance/competencies
// Mirrors .NET CreatePmsCompetency.
func (h *PerformanceMgtHandler) CreatePmsCompetency(w http.ResponseWriter, r *http.Request) {
	var req CreatePmsCompetencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.ObjectCategoryID == "" {
		response.Error(w, http.StatusBadRequest, "name and objectCategoryId are required")
		return
	}

	result, err := h.svc.Performance.CreatePmsCompetency(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreatePmsCompetency").Msg("failed to create pms competency")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdatePmsCompetency handles PUT /api/v1/performance/competencies
// Mirrors .NET UpdatePmsCompetency.
func (h *PerformanceMgtHandler) UpdatePmsCompetency(w http.ResponseWriter, r *http.Request) {
	var req UpdatePmsCompetencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PmsCompetencyID == "" {
		response.Error(w, http.StatusBadRequest, "pmsCompetencyId is required")
		return
	}

	result, err := h.svc.Performance.UpdatePmsCompetency(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdatePmsCompetency").Msg("failed to update pms competency")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreateDepartmentObjective handles POST /api/v1/performance/objectives/department
// Mirrors .NET CreateDepartmentObjective.
func (h *PerformanceMgtHandler) CreateDepartmentObjective(w http.ResponseWriter, r *http.Request) {
	var req CreateDepartmentObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	result, err := h.svc.Performance.CreateDepartmentObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateDepartmentObjective").Msg("failed to create department objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateDepartmentObjective handles PUT /api/v1/performance/objectives/department
// Mirrors .NET UpdateDepartmentObjective.
func (h *PerformanceMgtHandler) UpdateDepartmentObjective(w http.ResponseWriter, r *http.Request) {
	var req UpdateDepartmentObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DepartmentObjectiveID == "" {
		response.Error(w, http.StatusBadRequest, "departmentObjectiveId is required")
		return
	}

	result, err := h.svc.Performance.UpdateDepartmentObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateDepartmentObjective").Msg("failed to update department objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreateDivisionObjective handles POST /api/v1/performance/objectives/division
// Mirrors .NET CreateDivisionObjective.
func (h *PerformanceMgtHandler) CreateDivisionObjective(w http.ResponseWriter, r *http.Request) {
	var req CreateDivisionObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DepartmentObjectiveID == "" {
		response.Error(w, http.StatusBadRequest, "departmentObjectiveId is required")
		return
	}

	result, err := h.svc.Performance.CreateDivisionObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateDivisionObjective").Msg("failed to create division objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateDivisionObjective handles PUT /api/v1/performance/objectives/division
// Mirrors .NET UpdateDivisionObjective.
func (h *PerformanceMgtHandler) UpdateDivisionObjective(w http.ResponseWriter, r *http.Request) {
	var req UpdateDivisionObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DivisionObjectiveID == "" {
		response.Error(w, http.StatusBadRequest, "divisionObjectiveId is required")
		return
	}

	result, err := h.svc.Performance.UpdateDivisionObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateDivisionObjective").Msg("failed to update division objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreateOfficeObjective handles POST /api/v1/performance/objectives/office
// Mirrors .NET CreateOfficeObjective.
func (h *PerformanceMgtHandler) CreateOfficeObjective(w http.ResponseWriter, r *http.Request) {
	var req CreateOfficeObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	result, err := h.svc.Performance.CreateOfficeObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateOfficeObjective").Msg("failed to create office objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateOfficeObjective handles PUT /api/v1/performance/objectives/office
// Mirrors .NET UpdateOfficeObjective.
func (h *PerformanceMgtHandler) UpdateOfficeObjective(w http.ResponseWriter, r *http.Request) {
	var req UpdateOfficeObjectiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.OfficeObjectiveID == "" {
		response.Error(w, http.StatusBadRequest, "officeObjectiveId is required")
		return
	}

	result, err := h.svc.Performance.UpdateOfficeObjective(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateOfficeObjective").Msg("failed to update office objective")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreateObjectiveCategory handles POST /api/v1/performance/objective-categories
// Mirrors .NET CreateObjectiveCategory.
func (h *PerformanceMgtHandler) CreateObjectiveCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateObjectiveCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	result, err := h.svc.Performance.CreateObjectiveCategory(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateObjectiveCategory").Msg("failed to create objective category")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateObjectiveCategory handles PUT /api/v1/performance/objective-categories
// Mirrors .NET UpdateObjectiveCategory.
func (h *PerformanceMgtHandler) UpdateObjectiveCategory(w http.ResponseWriter, r *http.Request) {
	var req UpdateObjectiveCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ObjectiveCategoryID == "" {
		response.Error(w, http.StatusBadRequest, "objectiveCategoryId is required")
		return
	}

	result, err := h.svc.Performance.UpdateObjectiveCategory(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateObjectiveCategory").Msg("failed to update objective category")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// CreateCategoryDefinition handles POST /api/v1/performance/category-definitions
// Mirrors .NET CreateCategoryDefinition.
func (h *PerformanceMgtHandler) CreateCategoryDefinition(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryDefinitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ObjectiveCategoryID == "" || req.Description == "" {
		response.Error(w, http.StatusBadRequest, "objectiveCategoryId and description are required")
		return
	}

	result, err := h.svc.Performance.CreateCategoryDefinition(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "CreateCategoryDefinition").Msg("failed to create category definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(w, result)
}

// UpdateCategoryDefinition handles PUT /api/v1/performance/category-definitions
// Mirrors .NET UpdateCategoryDefinition.
func (h *PerformanceMgtHandler) UpdateCategoryDefinition(w http.ResponseWriter, r *http.Request) {
	var req UpdateCategoryDefinitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DefinitionID == "" {
		response.Error(w, http.StatusBadRequest, "definitionId is required")
		return
	}

	result, err := h.svc.Performance.UpdateCategoryDefinition(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "UpdateCategoryDefinition").Msg("failed to update category definition")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =========================================================================
// APPROVAL — POST endpoints
// =========================================================================

// ApproveRecords handles POST /api/v1/performance/approve
// Mirrors .NET ApproveRecords — approve one or more workflow records.
func (h *PerformanceMgtHandler) ApproveRecords(w http.ResponseWriter, r *http.Request) {
	var req ApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.EntityType == "" || len(req.RecordIDs) == 0 {
		response.Error(w, http.StatusBadRequest, "entityType and recordIds are required")
		return
	}

	result, err := h.svc.Performance.ApproveRecords(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "ApproveRecords").Msg("failed to approve records")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// RejectRecords handles POST /api/v1/performance/reject
// Mirrors .NET RejectRecords — reject one or more workflow records.
func (h *PerformanceMgtHandler) RejectRecords(w http.ResponseWriter, r *http.Request) {
	var req RejectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.EntityType == "" || len(req.RecordIDs) == 0 {
		response.Error(w, http.StatusBadRequest, "entityType and recordIds are required")
		return
	}

	result, err := h.svc.Performance.RejectRecords(r.Context(), req)
	if err != nil {
		h.log.Error().Err(err).Str("action", "RejectRecords").Msg("failed to reject records")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(w, result)
}

// =========================================================================
// SELECT LISTS — GET endpoints returning enum select lists
// =========================================================================

// GetObjectiveLevels handles GET /api/v1/performance/enums/objective-levels
// Mirrors .NET GetObjectiveLevels — returns ObjectiveLevel enum as select list.
func (h *PerformanceMgtHandler) GetObjectiveLevels(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Department", Value: strconv.Itoa(int(enums.ObjectiveLevelDepartment))},
		{Text: "Division", Value: strconv.Itoa(int(enums.ObjectiveLevelDivision))},
		{Text: "Office", Value: strconv.Itoa(int(enums.ObjectiveLevelOffice))},
		{Text: "Enterprise", Value: strconv.Itoa(int(enums.ObjectiveLevelEnterprise))},
	}
	response.OK(w, items)
}

// GetExtensionTargetTypes handles GET /api/v1/performance/enums/extension-target-types
// Mirrors .NET GetExtensionTargetTypes — returns ReviewPeriodExtensionTargetType enum as select list.
func (h *PerformanceMgtHandler) GetExtensionTargetTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Bankwide", Value: strconv.Itoa(int(enums.ExtensionTargetBankwide))},
		{Text: "Department", Value: strconv.Itoa(int(enums.ExtensionTargetDepartment))},
		{Text: "Division", Value: strconv.Itoa(int(enums.ExtensionTargetDivision))},
		{Text: "Office", Value: strconv.Itoa(int(enums.ExtensionTargetOffice))},
		{Text: "Staff", Value: strconv.Itoa(int(enums.ExtensionTargetStaff))},
	}
	response.OK(w, items)
}

// GetEvaluationTypes handles GET /api/v1/performance/enums/evaluation-types
// Mirrors .NET GetEvaluationTypes — returns EvaluationType enum as select list.
func (h *PerformanceMgtHandler) GetEvaluationTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Timeliness", Value: strconv.Itoa(int(enums.EvaluationTypeTimeliness))},
		{Text: "Quality", Value: strconv.Itoa(int(enums.EvaluationTypeQuality))},
		{Text: "Output", Value: strconv.Itoa(int(enums.EvaluationTypeOutput))},
	}
	response.OK(w, items)
}

// GetWorkProductTypes handles GET /api/v1/performance/enums/work-product-types
// Mirrors .NET GetWorkProductTypes — returns WorkProductType enum as select list.
func (h *PerformanceMgtHandler) GetWorkProductTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Operational", Value: strconv.Itoa(int(enums.WorkProductTypeOperational))},
		{Text: "Project", Value: strconv.Itoa(int(enums.WorkProductTypeProject))},
		{Text: "Committee", Value: strconv.Itoa(int(enums.WorkProductTypeCommittee))},
	}
	response.OK(w, items)
}

// GetGrievanceTypes handles GET /api/v1/performance/enums/grievance-types
// Mirrors .NET GetGrievanceTypes — returns GrievanceType enum as select list.
func (h *PerformanceMgtHandler) GetGrievanceTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "None", Value: strconv.Itoa(int(enums.GrievanceTypeNone))},
		{Text: "Work Product Evaluation", Value: strconv.Itoa(int(enums.GrievanceTypeWorkProductEvaluation))},
		{Text: "Work Product Assignment", Value: strconv.Itoa(int(enums.GrievanceTypeWorkProductAssignment))},
		{Text: "Work Product Planning", Value: strconv.Itoa(int(enums.GrievanceTypeWorkProductPlanning))},
		{Text: "Objective Planning", Value: strconv.Itoa(int(enums.GrievanceTypeObjectivePlanning))},
	}
	response.OK(w, items)
}

// GetFeedBackRequestTypes handles GET /api/v1/performance/enums/feedback-request-types
// Mirrors .NET GetFeedBackRequestTypes — returns FeedbackRequestType enum as select list.
func (h *PerformanceMgtHandler) GetFeedBackRequestTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Work Product Evaluation", Value: strconv.Itoa(int(enums.FeedbackRequestWorkProductEvaluation))},
		{Text: "Objective Planning", Value: strconv.Itoa(int(enums.FeedbackRequestObjectivePlanning))},
		{Text: "Project Planning", Value: strconv.Itoa(int(enums.FeedbackRequestProjectPlanning))},
		{Text: "Committee Planning", Value: strconv.Itoa(int(enums.FeedbackRequestCommitteePlanning))},
		{Text: "Work Product Feedback", Value: strconv.Itoa(int(enums.FeedbackRequestWorkProductFeedback))},
		{Text: "360 Review Feedback", Value: strconv.Itoa(int(enums.FeedbackRequest360ReviewFeedback))},
		{Text: "Work Product Planning", Value: strconv.Itoa(int(enums.FeedbackRequestWorkProductPlanning))},
		{Text: "Competency Review", Value: strconv.Itoa(int(enums.FeedbackRequestCompetencyReview))},
		{Text: "Review Period", Value: strconv.Itoa(int(enums.FeedbackRequestReviewPeriod))},
		{Text: "Review Period Extension", Value: strconv.Itoa(int(enums.FeedbackRequestReviewPeriodExtension))},
		{Text: "Project Member Assignment", Value: strconv.Itoa(int(enums.FeedbackRequestProjectMemberAssignment))},
		{Text: "Committee Member Assignment", Value: strconv.Itoa(int(enums.FeedbackRequestCommitteeMemberAssignment))},
		{Text: "Period Objective Outcome", Value: strconv.Itoa(int(enums.FeedbackRequestPeriodObjectiveOutcome))},
		{Text: "Dept Objective Outcome", Value: strconv.Itoa(int(enums.FeedbackRequestDeptObjectiveOutcome))},
		{Text: "Review Period 360 Review", Value: strconv.Itoa(int(enums.FeedbackRequestReviewPeriod360Review))},
		{Text: "Project Work Product Def", Value: strconv.Itoa(int(enums.FeedbackRequestProjectWorkProductDef))},
		{Text: "Committee Work Product Def", Value: strconv.Itoa(int(enums.FeedbackRequestCommitteeWorkProductDef))},
	}
	response.OK(w, items)
}

// GetPerformanceGrades handles GET /api/v1/performance/enums/performance-grades
// Mirrors .NET GetPerformanceGrades — returns PerformanceGrade enum as select list.
func (h *PerformanceMgtHandler) GetPerformanceGrades(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Probation", Value: strconv.Itoa(int(enums.PerformanceGradeProbation))},
		{Text: "Developing", Value: strconv.Itoa(int(enums.PerformanceGradeDeveloping))},
		{Text: "Progressive", Value: strconv.Itoa(int(enums.PerformanceGradeProgressive))},
		{Text: "Competent", Value: strconv.Itoa(int(enums.PerformanceGradeCompetent))},
		{Text: "Accomplished", Value: strconv.Itoa(int(enums.PerformanceGradeAccomplished))},
		{Text: "Exemplary", Value: strconv.Itoa(int(enums.PerformanceGradeExemplary))},
	}
	response.OK(w, items)
}

// GetReviewPeriodRange handles GET /api/v1/performance/enums/review-period-range
// Mirrors .NET GetReviewPeriodRange — returns ReviewPeriodRange enum as select list.
func (h *PerformanceMgtHandler) GetReviewPeriodRange(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Quarterly", Value: strconv.Itoa(int(enums.ReviewPeriodRangeQuarterly))},
		{Text: "Bi Annual", Value: strconv.Itoa(int(enums.ReviewPeriodRangeBiAnnual))},
		{Text: "Annual", Value: strconv.Itoa(int(enums.ReviewPeriodRangeAnnual))},
	}
	response.OK(w, items)
}

// GetStatuses handles GET /api/v1/performance/enums/statuses
// Mirrors .NET GetList — returns Status enum as select list.
func (h *PerformanceMgtHandler) GetStatuses(w http.ResponseWriter, r *http.Request) {
	items := []SelectListItem{
		{Text: "Draft", Value: strconv.Itoa(int(enums.StatusDraft))},
		{Text: "Pending Approval", Value: strconv.Itoa(int(enums.StatusPendingApproval))},
		{Text: "Approved And Active", Value: strconv.Itoa(int(enums.StatusApprovedAndActive))},
		{Text: "Returned", Value: strconv.Itoa(int(enums.StatusReturned))},
		{Text: "Rejected", Value: strconv.Itoa(int(enums.StatusRejected))},
		{Text: "Awaiting Evaluation", Value: strconv.Itoa(int(enums.StatusAwaitingEvaluation))},
		{Text: "Completed", Value: strconv.Itoa(int(enums.StatusCompleted))},
		{Text: "Paused", Value: strconv.Itoa(int(enums.StatusPaused))},
		{Text: "Cancelled", Value: strconv.Itoa(int(enums.StatusCancelled))},
		{Text: "Breached", Value: strconv.Itoa(int(enums.StatusBreached))},
		{Text: "Deactivated", Value: strconv.Itoa(int(enums.StatusDeactivated))},
		{Text: "All", Value: strconv.Itoa(int(enums.StatusAll))},
		{Text: "Closed", Value: strconv.Itoa(int(enums.StatusClosed))},
		{Text: "Pending Acceptance", Value: strconv.Itoa(int(enums.StatusPendingAcceptance))},
		{Text: "Active", Value: strconv.Itoa(int(enums.StatusActive))},
		{Text: "Pending Resolution", Value: strconv.Itoa(int(enums.StatusPendingResolution))},
		{Text: "Resolved Awaiting Feedback", Value: strconv.Itoa(int(enums.StatusResolvedAwaitingFeedback))},
		{Text: "Escalated", Value: strconv.Itoa(int(enums.StatusEscalated))},
		{Text: "Awaiting Respondent Comment", Value: strconv.Itoa(int(enums.StatusAwaitingRespondentComment))},
		{Text: "Pending HOD Review", Value: strconv.Itoa(int(enums.StatusPendingHODReview))},
		{Text: "Pending BU Head Review", Value: strconv.Itoa(int(enums.StatusPendingBUHeadReview))},
		{Text: "Pending HRD Review", Value: strconv.Itoa(int(enums.StatusPendingHRDReview))},
		{Text: "Pending HRD Approval", Value: strconv.Itoa(int(enums.StatusPendingHRDApproval))},
		{Text: "Suspension Pending Approval", Value: strconv.Itoa(int(enums.StatusSuspensionPendingApproval))},
		{Text: "Re Evaluate", Value: strconv.Itoa(int(enums.StatusReEvaluate))},
	}
	response.OK(w, items)
}
