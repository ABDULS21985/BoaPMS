package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ===========================================================================
// Setting Request Models (RequestVms.cs)
// ===========================================================================

// AddSettingRequestModel is the create payload for a setting.
type AddSettingRequestModel struct {
	Name        string `json:"name"        validate:"required"`
	Value       string `json:"value"       validate:"required"`
	Type        string `json:"type"        validate:"required"`
	IsEncrypted bool   `json:"isEncrypted"`
}

// SettingRequestModel is the update payload for a setting.
type SettingRequestModel struct {
	SettingID   string `json:"settingId"   validate:"required"`
	Name        string `json:"name"        validate:"required"`
	Value       string `json:"value"       validate:"required"`
	Type        string `json:"type"        validate:"required"`
	IsEncrypted bool   `json:"isEncrypted"`
}

// AddPmsConfigurationRequestModel is the create payload for a PMS config.
type AddPmsConfigurationRequestModel struct {
	Name        string `json:"name"        validate:"required"`
	Value       string `json:"value"       validate:"required"`
	Type        string `json:"type"        validate:"required"`
	IsEncrypted bool   `json:"isEncrypted"`
}

// PmsConfigurationRequestModel is the update payload for a PMS config.
type PmsConfigurationRequestModel struct {
	PmsConfigurationID string `json:"pmsConfigurationId" validate:"required"`
	Name               string `json:"name"               validate:"required"`
	Value              string `json:"value"              validate:"required"`
	Type               string `json:"type"               validate:"required"`
	IsEncrypted        bool   `json:"isEncrypted"`
}

// ===========================================================================
// Review Period Request Models
// ===========================================================================

// ReviewPeriodRequestVm is the update payload for a review period.
type ReviewPeriodRequestVm struct {
	BaseWorkFlowVm
	PeriodID              string    `json:"periodId"              validate:"required"`
	Year                  int       `json:"year"`
	Range                 int       `json:"range"`
	RangeValue            int       `json:"rangeValue"            validate:"required"`
	Name                  string    `json:"name"                  validate:"required"`
	Description           string    `json:"description"`
	ShortName             string    `json:"shortName"`
	StartDate             time.Time `json:"startDate"`
	EndDate               time.Time `json:"endDate"`
	MaxPoints             float64   `json:"maxPoints"             validate:"required"`
	MinNoOfObjectives     int       `json:"minNoOfObjectives"     validate:"required"`
	MaxNoOfObjectives     int       `json:"maxNoOfObjectives"     validate:"required"`
	AllowObjectivePlanning bool    `json:"allowObjectivePlanning"`
	StrategyID            string    `json:"strategyId"            validate:"required"`
}

// ReviewPeriodListVm wraps a list of review period request VMs.
type ReviewPeriodListVm struct {
	BaseAPIResponse
	ReviewPeriods []ReviewPeriodRequestVm `json:"reviewPeriods"`
	TotalRecord   int                     `json:"totalRecord"`
}

// SearchReviewPeriodVm carries search filters for review periods.
type SearchReviewPeriodVm struct {
	BasePagedData
	CategoryID   *int   `json:"categoryId"`
	SearchString string `json:"searchString"`
	IsApproved   *bool  `json:"isApproved"`
	IsRejected   *bool  `json:"isRejected"`
	IsTechnical  *bool  `json:"isTechnical"`
}

// CreateNewReviewPeriodVm is the create payload for a review period.
type CreateNewReviewPeriodVm struct {
	Name                   string  `json:"name"                   validate:"required"`
	ShortName              string  `json:"shortName"`
	Description            string  `json:"description"`
	Year                   int     `json:"year"                   validate:"required"`
	Range                  int     `json:"range"                  validate:"required"`
	RangeValue             int     `json:"rangeValue"             validate:"required"`
	MaxPoints              float64 `json:"maxPoints"              validate:"required"`
	MinNoOfObjectives      int     `json:"minNoOfObjectives"      validate:"required"`
	MaxNoOfObjectives      int     `json:"maxNoOfObjectives"      validate:"required"`
	AllowObjectivePlanning bool    `json:"allowObjectivePlanning"`
	StrategyID             string  `json:"strategyId"             validate:"required"`
}

// ===========================================================================
// Period Objective Request Models
// ===========================================================================

// PeriodObjectiveRequestVm is the DTO for a period objective operation.
type PeriodObjectiveRequestVm struct {
	BaseWorkFlowVm
	PeriodObjectiveID string `json:"periodObjectiveId"`
	ObjectiveID       string `json:"objectiveId"`
	ReviewPeriodID    string `json:"reviewPeriodId"`
}

// AddPeriodObjectiveVm is the create payload for period objectives.
type AddPeriodObjectiveVm struct {
	ObjectiveIDs        []string `json:"objectiveIds"       validate:"required"`
	ReviewPeriodID      string   `json:"reviewPeriodId"     validate:"required"`
	ObjectiveID         *string  `json:"objectiveId"`
	PeriodObjectiveID   *string  `json:"periodObjectiveId"`
	ObjectiveCategoryID *string  `json:"objectiveCategoryId"`
}

// SaveDraftPeriodObjectiveVm is the draft-save payload for period objectives.
type SaveDraftPeriodObjectiveVm struct {
	ObjectiveIDs        []string `json:"objectiveIds"`
	ReviewPeriodID      string   `json:"reviewPeriodId"`
	ObjectiveID         *string  `json:"objectiveId"`
	ObjectiveCategoryID *string  `json:"objectiveCategoryId"`
}

// ===========================================================================
// Planned Objective Request Models
// ===========================================================================

// ReviewPeriodIndividualPlannedObjectiveRequestModel is the update payload.
type ReviewPeriodIndividualPlannedObjectiveRequestModel struct {
	BaseWorkFlowVm
	PlannedObjectiveID string               `json:"plannedObjectiveId"`
	ObjectiveID        string               `json:"objectiveId"        validate:"required"`
	StaffID            string               `json:"staffId"            validate:"required"`
	ObjectiveLevel     enums.ObjectiveLevel `json:"objectiveLevel"     validate:"required"`
	StaffJobRole       string               `json:"staffJobRole"`
	ReviewPeriodID     string               `json:"reviewPeriodId"     validate:"required"`
	Remark             string               `json:"remark"`
}

// AddReviewPeriodIndividualPlannedObjectiveRequestModel is the create payload.
type AddReviewPeriodIndividualPlannedObjectiveRequestModel struct {
	ObjectiveID    string `json:"objectiveId"    validate:"required"`
	StaffID        string `json:"staffId"        validate:"required"`
	ReviewPeriodID string `json:"reviewPeriodId" validate:"required"`
}

// ===========================================================================
// Category Definition Request Models
// ===========================================================================

// CategoryDefinitionRequestVm is the update payload for a category definition.
type CategoryDefinitionRequestVm struct {
	BaseWorkFlowVm
	DefinitionID            string  `json:"definitionId"`
	ObjectiveCategoryID     string  `json:"objectiveCategoryId"     validate:"required"`
	ReviewPeriodID          string  `json:"reviewPeriodId"          validate:"required"`
	Weight                  float64 `json:"weight"                  validate:"required"`
	MaxNoObjectives         int     `json:"maxNoObjectives"`
	MaxNoWorkProduct        int     `json:"maxNoWorkProduct"`
	MaxPoints               int     `json:"maxPoints"`
	IsCompulsory            bool    `json:"isCompulsory"`
	EnforceWorkProductLimit bool    `json:"enforceWorkProductLimit"`
	Description             string  `json:"description"`
	GradeGroupID            int     `json:"gradeGroupId"            validate:"required"`
}

// CreateCategoryDefinitionRequestVm is the create payload.
type CreateCategoryDefinitionRequestVm struct {
	DefinitionID            string  `json:"definitionId"`
	ObjectiveCategoryID     string  `json:"objectiveCategoryId"     validate:"required"`
	ReviewPeriodID          string  `json:"reviewPeriodId"          validate:"required"`
	Weight                  float64 `json:"weight"                  validate:"required"`
	MaxNoObjectives         int     `json:"maxNoObjectives"         validate:"required"`
	MaxNoWorkProduct        int     `json:"maxNoWorkProduct"        validate:"required"`
	IsCompulsory            bool    `json:"isCompulsory"`
	EnforceWorkProductLimit bool    `json:"enforceWorkProductLimit"`
	Description             string  `json:"description"`
	GradeGroupID            int     `json:"gradeGroupId"            validate:"required"`
}

// ===========================================================================
// Review Period Extension Request Models
// ===========================================================================

// ReviewPeriodExtensionRequestModel is the update payload for an extension.
type ReviewPeriodExtensionRequestModel struct {
	BaseWorkFlowVm
	ReviewPeriodExtensionID string    `json:"reviewPeriodExtensionId"`
	ReviewPeriodID          string    `json:"reviewPeriodId" validate:"required"`
	TargetType              int       `json:"targetType"`
	TargetReference         string    `json:"targetReference"`
	Description             string    `json:"description"`
	StartDate               time.Time `json:"startDate"`
	EndDate                 time.Time `json:"endDate"`
}

// CreateReviewPeriodExtensionRequestModel is the create payload.
type CreateReviewPeriodExtensionRequestModel struct {
	ReviewPeriodID          string    `json:"reviewPeriodId"          validate:"required"`
	ReviewPeriodExtensionID string    `json:"reviewPeriodExtensionId"`
	TargetType              int       `json:"targetType"`
	TargetReference         string    `json:"targetReference"`
	Description             string    `json:"description"`
	StartDate               time.Time `json:"startDate"               validate:"required"`
	EndDate                 time.Time `json:"endDate"                 validate:"required"`
}

// ===========================================================================
// 360 Review Request Models
// ===========================================================================

// ReviewPeriod360ReviewRequestModel is the update payload.
type ReviewPeriod360ReviewRequestModel struct {
	BaseWorkFlowVm
	ReviewPeriod360ReviewID string `json:"reviewPeriod360ReviewId"`
	ReviewPeriodID          string `json:"reviewPeriodId" validate:"required"`
	TargetType              int    `json:"targetType"`
	TargetReference         string `json:"targetReference"`
}

// CreateReviewPeriod360ReviewRequestModel is the create payload.
type CreateReviewPeriod360ReviewRequestModel struct {
	ReviewPeriodID  string `json:"reviewPeriodId"  validate:"required"`
	TargetType      int    `json:"targetType"      validate:"required"`
	TargetReference string `json:"targetReference" validate:"required"`
}

// NOTE: Trigger360ReviewRequestModel, Initiate360ReviewRequestModel,
//       Complete360ReviewRequestModel are defined in dto_feedback.go.

// ResetPasswordRequestModel is the password reset request.
type ResetPasswordRequestModel struct {
	Username   string `json:"username"   validate:"required"`
	Password   string `json:"password"   validate:"required"`
	DeviceName string `json:"deviceName"`
	IPAddress  string `json:"ipAddress"`
}

// NOTE: FeedbackRequestModel and TreatFeedbackRequestModel are defined
//       in dto_feedback.go.

// ===========================================================================
// Project Request Models
// ===========================================================================

// ProjectRequestModel is the update payload for a project.
type ProjectRequestModel struct {
	BaseWorkFlowVm
	ProjectID      string    `json:"projectId"`
	ProjectManager string    `json:"projectManager" validate:"required"`
	Name           string    `json:"name"           validate:"required"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"startDate"      validate:"required"`
	EndDate        time.Time `json:"endDate"        validate:"required"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"reviewPeriodId" validate:"required"`
	DepartmentID   int       `json:"departmentId"   validate:"required"`
}

// CreateProjectRequestModel is the create payload for a project.
type CreateProjectRequestModel struct {
	ProjectManager string    `json:"projectManager" validate:"required"`
	Name           string    `json:"name"           validate:"required"`
	Description    string    `json:"description"`
	DepartmentID   int       `json:"departmentId"   validate:"required"`
	StartDate      time.Time `json:"startDate"      validate:"required"`
	EndDate        time.Time `json:"endDate"        validate:"required"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"reviewPeriodId" validate:"required"`
}

// ProjectMemberRequestModel is the update payload for a project member.
type ProjectMemberRequestModel struct {
	BaseWorkFlowVm
	ProjectMemberID    string `json:"projectMemberId"    validate:"required"`
	StaffID            string `json:"staffId"            validate:"required"`
	ProjectID          string `json:"projectId"          validate:"required"`
	PlannedObjectiveID string `json:"plannedObjectiveId" validate:"required"`
}

// AddProjectMemberRequestModel is the create payload for a project member.
type AddProjectMemberRequestModel struct {
	StaffID            string `json:"staffId"            validate:"required"`
	ProjectID          string `json:"projectId"          validate:"required"`
	PlannedObjectiveID string `json:"plannedObjectiveId" validate:"required"`
}

// ProjectObjectiveRequestModel is the update payload for a project objective.
type ProjectObjectiveRequestModel struct {
	BaseEntityVm
	ProjectObjectiveID string `json:"projectObjectiveId" validate:"required"`
	ObjectiveID        string `json:"objectiveId"        validate:"required"`
	ProjectID          string `json:"projectId"          validate:"required"`
}

// AddProjectObjectiveRequestModel is the create payload.
type AddProjectObjectiveRequestModel struct {
	ObjectiveID string `json:"objectiveId" validate:"required"`
	ProjectID   string `json:"projectId"   validate:"required"`
}

// ChangeAdhocLeadRequestModel is the payload to change project/committee lead.
type ChangeAdhocLeadRequestModel struct {
	ReferenceID         string `json:"referenceId"         validate:"required"`
	StaffID             string `json:"staffId"             validate:"required"`
	AdhocAssignmentType int    `json:"adhocAssignmentType" validate:"required"`
}

// ===========================================================================
// Committee Request Models
// ===========================================================================

// CommitteeRequestModel is the update payload for a committee.
type CommitteeRequestModel struct {
	BaseWorkFlowVm
	CommitteeID    string    `json:"committeeId"`
	Chairperson    string    `json:"chairperson"    validate:"required"`
	Name           string    `json:"name"           validate:"required"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"startDate"      validate:"required"`
	EndDate        time.Time `json:"endDate"        validate:"required"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"reviewPeriodId" validate:"required"`
	DepartmentID   int       `json:"departmentId"   validate:"required"`
}

// CreateCommitteeRequestModel is the create payload for a committee.
type CreateCommitteeRequestModel struct {
	Chairperson    string    `json:"chairperson"    validate:"required"`
	Name           string    `json:"name"           validate:"required"`
	Description    string    `json:"description"`
	DepartmentID   int       `json:"departmentId"   validate:"required"`
	StartDate      time.Time `json:"startDate"      validate:"required"`
	EndDate        time.Time `json:"endDate"        validate:"required"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"reviewPeriodId" validate:"required"`
}

// CommitteeMemberRequestModel is the update payload for a committee member.
type CommitteeMemberRequestModel struct {
	BaseWorkFlowVm
	CommitteeMemberID  string `json:"committeeMemberId"  validate:"required"`
	StaffID            string `json:"staffId"            validate:"required"`
	CommitteeID        string `json:"committeeId"        validate:"required"`
	PlannedObjectiveID string `json:"plannedObjectiveId" validate:"required"`
}

// AddCommitteeMemberRequestModel is the create payload.
type AddCommitteeMemberRequestModel struct {
	StaffID            string `json:"staffId"            validate:"required"`
	CommitteeID        string `json:"committeeId"        validate:"required"`
	PlannedObjectiveID string `json:"plannedObjectiveId" validate:"required"`
}

// CommitteeObjectiveRequestModel is the update payload for a committee objective.
type CommitteeObjectiveRequestModel struct {
	BaseEntityVm
	CommitteeObjectiveID string `json:"committeeObjectiveId" validate:"required"`
	ObjectiveID          string `json:"objectiveId"          validate:"required"`
	CommitteeID          string `json:"committeeId"          validate:"required"`
}

// AddCommitteeObjectiveRequestModel is the create payload.
type AddCommitteeObjectiveRequestModel struct {
	ObjectiveID string `json:"objectiveId" validate:"required"`
	CommitteeID string `json:"committeeId" validate:"required"`
}

// ===========================================================================
// Work Product Request Models
// ===========================================================================

// WorkProductRequestModel is the update payload for a work product.
type WorkProductRequestModel struct {
	BaseWorkFlowVm
	WorkProductID                  string    `json:"workProductId"                  validate:"required"`
	WorkProductDefinitionID        *string   `json:"workProductDefinitionId"`
	ProjectAssignedWorkProductID   *string   `json:"projectAssignedWorkProductId"`
	CommitteeAssignedWorkProductID *string   `json:"committeeAssignedWorkProductId"`
	Name                           string    `json:"name"                           validate:"required"`
	Description                    string    `json:"description"`
	MaxPoint                       float64   `json:"maxPoint"`
	WorkProductType                int       `json:"workProductType"                validate:"required"`
	IsSelfCreated                  bool      `json:"isSelfCreated"`
	StaffID                        string    `json:"staffId"                        validate:"required"`
	AcceptanceComment              string    `json:"acceptanceComment"`
	StartDate                      time.Time `json:"startDate"                      validate:"required"`
	EndDate                        time.Time `json:"endDate"                        validate:"required"`
	Deliverables                   string    `json:"deliverables"`
	FinalScore                     int       `json:"finalScore"`
	NoReturned                     int       `json:"noReturned"`
	CompletionDate                 time.Time `json:"completionDate"`
	ApproverComment                string    `json:"approverComment"`
	ProjectID                      string    `json:"projectId"`
	CommitteeID                    string    `json:"committeeId"`
	PlannedObjectiveID             string    `json:"plannedObjectiveId"`
	Remark                         string    `json:"remark"`
	Attachment                     string    `json:"attachment"`
}

// CreateWorkProductRequestModel is the create payload for a work product.
type CreateWorkProductRequestModel struct {
	BaseWorkFlowVm
	WorkProductDefinitionID        *string   `json:"workProductDefinitionId"        validate:"required"`
	ProjectAssignedWorkProductID   *string   `json:"projectAssignedWorkProductId"`
	CommitteeAssignedWorkProductID *string   `json:"committeeAssignedWorkProductId"`
	Name                           string    `json:"name"`
	Description                    string    `json:"description"`
	WorkProductType                int       `json:"workProductType"                validate:"required"`
	StaffID                        string    `json:"staffId"                        validate:"required"`
	StartDate                      time.Time `json:"startDate"                      validate:"required"`
	EndDate                        time.Time `json:"endDate"                        validate:"required"`
	Deliverables                   string    `json:"deliverables"`
	ProjectID                      string    `json:"projectId"`
	CommitteeID                    string    `json:"committeeId"`
	PlannedObjectiveID             string    `json:"plannedObjectiveId"`
	PlannedObjectiveName           string    `json:"plannedObjectiveName"`
	StrObjectiveID                 string    `json:"strObjectiveId"`
}

// ProjectAssignedWorkProductRequestModel is the update payload.
type ProjectAssignedWorkProductRequestModel struct {
	BaseWorkFlowVm
	ProjectAssignedWorkProductID string    `json:"projectAssignedWorkProductId" validate:"required"`
	WorkProductDefinitionID      string    `json:"workProductDefinitionId"      validate:"required"`
	Name                         string    `json:"name"                         validate:"required"`
	Description                  string    `json:"description"`
	ProjectID                    string    `json:"projectId"                    validate:"required"`
	StartDate                    time.Time `json:"startDate"                    validate:"required"`
	EndDate                      time.Time `json:"endDate"                      validate:"required"`
	Deliverables                 string    `json:"deliverables"`
}

// CreateProjectAssignedWorkProductRequestModel is the create payload.
type CreateProjectAssignedWorkProductRequestModel struct {
	Name         string    `json:"name"         validate:"required"`
	Description  string    `json:"description"`
	ProjectID    string    `json:"projectId"    validate:"required"`
	StartDate    time.Time `json:"startDate"    validate:"required"`
	EndDate      time.Time `json:"endDate"      validate:"required"`
	Deliverables string    `json:"deliverables"`
}

// CommitteeAssignedWorkProductRequestModel is the update payload.
type CommitteeAssignedWorkProductRequestModel struct {
	BaseWorkFlowVm
	CommitteeAssignedWorkProductID string    `json:"committeeAssignedWorkProductId" validate:"required"`
	WorkProductDefinitionID        string    `json:"workProductDefinitionId"        validate:"required"`
	Name                           string    `json:"name"                           validate:"required"`
	Description                    string    `json:"description"`
	CommitteeID                    string    `json:"committeeId"                    validate:"required"`
	StartDate                      time.Time `json:"startDate"                      validate:"required"`
	EndDate                        time.Time `json:"endDate"                        validate:"required"`
	Deliverables                   string    `json:"deliverables"`
}

// CreateCommitteeAssignedWorkProductRequestModel is the create payload.
type CreateCommitteeAssignedWorkProductRequestModel struct {
	Name         string    `json:"name"         validate:"required"`
	Description  string    `json:"description"`
	CommitteeID  string    `json:"committeeId"  validate:"required"`
	StartDate    time.Time `json:"startDate"    validate:"required"`
	EndDate      time.Time `json:"endDate"      validate:"required"`
	Deliverables string    `json:"deliverables"`
}

// ===========================================================================
// Work Product Task Request Models
// ===========================================================================

// WorkProductTaskRequestModel is the update payload for a work product task.
type WorkProductTaskRequestModel struct {
	BaseEntityVm
	WorkProductTaskID string    `json:"workProductTaskId" validate:"required"`
	Name              string    `json:"name"              validate:"required"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"startDate"         validate:"required"`
	EndDate           time.Time `json:"endDate"           validate:"required"`
	CompletionDate    time.Time `json:"completionDate"`
	WorkProductID     string    `json:"workProductId"     validate:"required"`
	Remark            string    `json:"remark"`
	Attachment        string    `json:"attachment"`
}

// AddWorkProductTaskRequestModel is the create payload.
type AddWorkProductTaskRequestModel struct {
	Name          string    `json:"name"          validate:"required"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"startDate"     validate:"required"`
	EndDate       time.Time `json:"endDate"       validate:"required"`
	WorkProductID string    `json:"workProductId" validate:"required"`
}

// ===========================================================================
// Work Product Evaluation Request Models
// ===========================================================================

// WorkProductEvaluationRequestModel is the update payload for an evaluation.
type WorkProductEvaluationRequestModel struct {
	BaseEntityVm
	WorkProductEvaluationID      string  `json:"workProductEvaluationId"      validate:"required"`
	WorkProductID                string  `json:"workProductId"                validate:"required"`
	Timeliness                   float64 `json:"timeliness"                   validate:"required"`
	TimelinessEvaluationOptionID string  `json:"timelinessEvaluationOptionId" validate:"required"`
	Quality                      float64 `json:"quality"                      validate:"required"`
	QualityEvaluationOptionID    string  `json:"qualityEvaluationOptionId"    validate:"required"`
	Output                       float64 `json:"output"                       validate:"required"`
	OutputEvaluationOptionID     string  `json:"outputEvaluationOptionId"     validate:"required"`
	TotalOutcome                 float64 `json:"totalOutcome"`
	Outcome                      float64 `json:"outcome"`
	EvaluatorStaffID             string  `json:"evaluatorStaffId"`
}

// AddWorkProductEvaluationRequestModel is the create payload.
type AddWorkProductEvaluationRequestModel struct {
	WorkProductID                string  `json:"workProductId"                validate:"required"`
	TimelinessEvaluationOptionID string  `json:"timelinessEvaluationOptionId" validate:"required"`
	QualityEvaluationOptionID    string  `json:"qualityEvaluationOptionId"    validate:"required"`
	OutputEvaluationOptionID     string  `json:"outputEvaluationOptionId"     validate:"required"`
	WorkProductEvaluationID      string  `json:"workProductEvaluationId"`
	TotalOutcome                 float64 `json:"totalOutcome"`
	Outcome                      float64 `json:"outcome"`
	EvaluatorStaffID             string  `json:"evaluatorStaffId"`
}

// ===========================================================================
// Period Objective Evaluation Request Models
// ===========================================================================

// PeriodObjectiveEvaluationRequestModel is the update payload.
type PeriodObjectiveEvaluationRequestModel struct {
	BaseWorkFlowVm
	PeriodObjectiveEvaluationID string  `json:"periodObjectiveEvaluationId" validate:"required"`
	TotalOutcomeScore           float64 `json:"totalOutcomeScore"           validate:"required"`
	OutcomeScore                float64 `json:"outcomeScore"                validate:"required"`
	EnterpriseObjectiveID       string  `json:"enterpriseObjectiveId"       validate:"required"`
	ReviewPeriodID              string  `json:"reviewPeriodId"              validate:"required"`
}

// AddPeriodObjectiveEvaluationRequestModel is the create payload.
type AddPeriodObjectiveEvaluationRequestModel struct {
	TotalOutcomeScore     float64 `json:"totalOutcomeScore"     validate:"required"`
	OutcomeScore          float64 `json:"outcomeScore"          validate:"required"`
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId" validate:"required"`
	ReviewPeriodID        string  `json:"reviewPeriodId"        validate:"required"`
}

// PeriodObjectiveDepartmentEvaluationRequestModel is the update payload.
type PeriodObjectiveDepartmentEvaluationRequestModel struct {
	BaseWorkFlowVm
	PeriodObjectiveDepartmentEvaluationID string  `json:"periodObjectiveDepartmentEvaluationId" validate:"required"`
	DepartmentID                          int     `json:"departmentId"                          validate:"required"`
	OverallOutcomeScored                  float64 `json:"overallOutcomeScored"                  validate:"required"`
	AllocatedOutcome                      float64 `json:"allocatedOutcome"                      validate:"required"`
	OutcomeScore                          float64 `json:"outcomeScore"                          validate:"required"`
	ReviewPeriodID                        string  `json:"reviewPeriodId"                        validate:"required"`
	EnterpriseObjectiveID                 string  `json:"enterpriseObjectiveId"                 validate:"required"`
}

// AddPeriodObjectiveDepartmentEvaluationRequestModel is the create payload.
type AddPeriodObjectiveDepartmentEvaluationRequestModel struct {
	DepartmentID          int     `json:"departmentId"          validate:"required"`
	OverallOutcomeScored  float64 `json:"overallOutcomeScored"  validate:"required"`
	AllocatedOutcome      float64 `json:"allocatedOutcome"      validate:"required"`
	OutcomeScore          float64 `json:"outcomeScore"          validate:"required"`
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId" validate:"required"`
	ReviewPeriodID        string  `json:"reviewPeriodId"        validate:"required"`
}

// ===========================================================================
// Competency Review Feedback Request Models
// ===========================================================================

// CompetencyReviewFeedbackRequestModel is the update payload.
type CompetencyReviewFeedbackRequestModel struct {
	BaseWorkFlowVm
	CompetencyReviewFeedbackID string  `json:"competencyReviewFeedbackId" validate:"required"`
	StaffID                    string  `json:"staffId"                    validate:"required"`
	FinalScore                 float64 `json:"finalScore"`
	ReviewPeriodID             string  `json:"reviewPeriodId"             validate:"required"`
}

// AddCompetencyReviewFeedbackRequestModel is the create payload.
type AddCompetencyReviewFeedbackRequestModel struct {
	StaffID        string  `json:"staffId"        validate:"required"`
	FinalScore     float64 `json:"finalScore"`
	ReviewPeriodID string  `json:"reviewPeriodId" validate:"required"`
}

// CompetencyGapClosureRequestModel is the gap closure payload.
type CompetencyGapClosureRequestModel struct {
	BaseEntityVm
	CompetencyGapClosureID string  `json:"competencyGapClosureId"`
	StaffID                string  `json:"staffId"                validate:"required"`
	MaxPoints              float64 `json:"maxPoints"`
	FinalScore             float64 `json:"finalScore"`
	ReviewPeriodID         string  `json:"reviewPeriodId"         validate:"required"`
	ObjectiveCategoryID    string  `json:"objectiveCategoryId"    validate:"required"`
}

// CompetencyReviewerRequestModel is the update payload for a reviewer.
type CompetencyReviewerRequestModel struct {
	BaseEntityVm
	CompetencyReviewerID       string `json:"competencyReviewerId"       validate:"required"`
	ReviewStaffID              string `json:"reviewStaffId"              validate:"required"`
	CompetencyReviewFeedbackID string `json:"competencyReviewFeedbackId" validate:"required"`
}

// AddCompetencyReviewerRequestModel is the create payload.
type AddCompetencyReviewerRequestModel struct {
	ReviewStaffID              string `json:"reviewStaffId"              validate:"required"`
	CompetencyReviewFeedbackID string `json:"competencyReviewFeedbackId" validate:"required"`
}

// CompetencyReviewerRatingRequestModel is the update payload for a rating.
type CompetencyReviewerRatingRequestModel struct {
	BaseEntityVm
	CompetencyReviewerRatingID   string  `json:"competencyReviewerRatingId"   validate:"required"`
	PmsCompetencyID              string  `json:"pmsCompetencyId"              validate:"required"`
	FeedbackQuestionaireOptionID string  `json:"feedbackQuestionaireOptionId" validate:"required"`
	Rating                       float64 `json:"rating"`
	CompetencyReviewerID         string  `json:"competencyReviewerId"         validate:"required"`
}

// AddCompetencyReviewerRatingRequestModel is the create payload.
type AddCompetencyReviewerRatingRequestModel struct {
	PmsCompetencyID              string `json:"pmsCompetencyId"              validate:"required"`
	FeedbackQuestionaireOptionID string `json:"feedbackQuestionaireOptionId" validate:"required"`
	CompetencyReviewerID         string `json:"competencyReviewerId"         validate:"required"`
}
