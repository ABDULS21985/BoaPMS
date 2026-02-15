package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ===========================================================================
// Generic Response Wrappers (ResponseViewModels.cs)
// ===========================================================================

// NOTE: GenericListResponseVm is defined in dto_feedback.go.

// ResponseViewModel is a generic API response carrying any data payload.
type ResponseViewModel struct {
	BaseAPIResponse
	Data interface{} `json:"data"`
}

// GenericListVm is a generic list response with untyped data.
type GenericListVm struct {
	BaseAPIResponse
	ListData    interface{} `json:"listData"`
	TotalRecord int         `json:"totalRecord"`
}

// PaginatedResult carries a page of items with a total count.
type PaginatedResult[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"totalCount"`
}

// ===========================================================================
// Approval VMs (ApprovalVm.cs)
// ===========================================================================

// ApprovalBase carries common fields for approval/rejection requests.
type ApprovalBase struct {
	Approval   enums.Approval `json:"approval"`
	EntityType string         `json:"entityType"`
	RecordIDs  []string       `json:"recordIds"`
}

// ApprovalRequestVm is the request payload for approving records.
type ApprovalRequestVm struct {
	ApprovalBase
}

// RejectionRequestVm is the request payload for rejecting records.
type RejectionRequestVm struct {
	ApprovalBase
	RejectionReason string `json:"rejectionReason"`
}

// ===========================================================================
// Review Period Response VMs
// ===========================================================================

// ReviewPeriodResponseVm wraps a single review period in an API response.
type ReviewPeriodResponseVm struct {
	BaseAPIResponse
	PerformanceReviewPeriod *PerformanceReviewPeriod `json:"performanceReviewPeriod"`
	HasExtension            bool                     `json:"hasExtension"`
	ExtensionEndDate        time.Time                `json:"extensionEndDate"`
}

// PerformanceReviewPeriodResponseVm wraps a review period VM.
type PerformanceReviewPeriodResponseVm struct {
	BaseAPIResponse
	PerformanceReviewPeriod *PerformanceReviewPeriodVm `json:"performanceReviewPeriod"`
}

// GetAllReviewPeriodResponseVm wraps a list of review period VMs.
type GetAllReviewPeriodResponseVm struct {
	GenericListResponseVm
	PerformanceReviewPeriods []PerformanceReviewPeriodVm `json:"performanceReviewPeriods"`
}

// GetStaffReviewPeriodResponseVm wraps review periods for a specific staff.
type GetStaffReviewPeriodResponseVm struct {
	GenericListResponseVm
	StaffID                  string                      `json:"staffId"`
	PerformanceReviewPeriods []PerformanceReviewPeriodVm `json:"performanceReviewPeriods"`
}

// ReviewPeriodExtensionResponseVm wraps a single review period extension.
type ReviewPeriodExtensionResponseVm struct {
	BaseAPIResponse
	ReviewPeriodExtension *ReviewPeriodExtension `json:"reviewPeriodExtension"`
}

// ReviewPeriodExtensionListResponseVm wraps a list of extension data.
type ReviewPeriodExtensionListResponseVm struct {
	GenericListResponseVm
	ReviewPeriodExtensions []ReviewPeriodExtensionData `json:"reviewPeriodExtensions"`
}

// ReviewPeriod360ReviewData is the flat DTO for a 360 review configuration.
type ReviewPeriod360ReviewData struct {
	ReviewPeriod360ReviewID string `json:"reviewPeriod360ReviewId"`
	ReviewPeriodID          string `json:"reviewPeriodId"`
	ReviewPeriod            string `json:"reviewPeriod"`
	TargetType              int    `json:"targetType"`
	TargetReference         string `json:"targetReference"`
	RecordStatus            string `json:"recordStatus"`
	IsActive                bool   `json:"isActive"`
}

// ReviewPeriod360ReviewListResponseVm wraps a list of 360 review records.
type ReviewPeriod360ReviewListResponseVm struct {
	GenericListResponseVm
	Reviews []ReviewPeriod360ReviewData `json:"reviews"`
}

// ReviewPeriodExtensionData is the flat DTO for a review period extension.
type ReviewPeriodExtensionData struct {
	ReviewPeriodExtensionID string       `json:"reviewPeriodExtensionId"`
	ReviewPeriodID          string       `json:"reviewPeriodId"`
	ReviewPeriod            string       `json:"reviewPeriod"`
	TargetType              int          `json:"targetType"`
	TargetTypeName          string       `json:"targetTypeName"`
	TargetReference         string       `json:"targetReference"`
	Description             string       `json:"description"`
	StartDate               time.Time    `json:"startDate"`
	EndDate                 time.Time    `json:"endDate"`
	RecordStatus            enums.Status `json:"recordStatus"`
}

// ===========================================================================
// Enterprise Objective Response VMs
// ===========================================================================

// ReviewPeriodObjectivesResponseVm wraps objectives for a review period.
type ReviewPeriodObjectivesResponseVm struct {
	GenericListResponseVm
	Objectives []EnterpriseObjectiveData `json:"objectives"`
}

// EnterpriseObjectiveResponseVm wraps a single enterprise objective.
type EnterpriseObjectiveResponseVm struct {
	BaseAPIResponse
	EnterpriseObjective *EnterpriseObjectiveData `json:"enterpriseObjective"`
}

// EnterpriseObjectiveData carries enterprise objective details used in
// response VMs.
type EnterpriseObjectiveData struct {
	BaseWorkFlowVm
	EnterpriseObjectiveID          string                   `json:"enterpriseObjectiveId"`
	PeriodObjectiveID              string                   `json:"periodObjectiveId"`
	Name                           string                   `json:"name"`
	Description                    string                   `json:"description"`
	Kpi                            string                   `json:"kpi"`
	Target                         string                   `json:"target"`
	EnterpriseObjectivesCategoryID string                   `json:"enterpriseObjectivesCategoryId"`
	StrategyID                     string                   `json:"strategyId"`
	HasEvaluation                  bool                     `json:"hasEvaluation"`
	OutcomeScore                   float64                  `json:"outcomeScore"`
	TotalOutcomeScore              float64                  `json:"totalOutcomeScore"`
	Evaluator                      string                   `json:"evaluator"`
	EvaluationDate                 time.Time                `json:"evaluationDate"`
	CategoryDefinitions            []CategoryDefinitionData `json:"categoryDefinitions"`
}

// ConsolidatedObjectiveListResponseVm wraps consolidated objectives.
type ConsolidatedObjectiveListResponseVm struct {
	GenericListResponseVm
	Objectives []ConsolidatedObjectiveVm `json:"objectives"`
}

// ===========================================================================
// Category Definition Response VMs
// ===========================================================================

// ReviewPeriodCategoryDefinitionResponseVm wraps category definitions.
type ReviewPeriodCategoryDefinitionResponseVm struct {
	GenericListResponseVm
	CategoryDefinitions []CategoryDefinitionData `json:"categoryDefinitions"`
}

// CategoryDefinitionData is the flat DTO for a category definition.
type CategoryDefinitionData struct {
	BaseWorkFlowVm
	DefinitionID            string  `json:"definitionId"`
	ObjectiveCategoryID     string  `json:"objectiveCategoryId"`
	ReviewPeriodID          string  `json:"reviewPeriodId"`
	Weight                  float64 `json:"weight"`
	MaxNoObjectives         int     `json:"maxNoObjectives"`
	MaxNoWorkProduct        int     `json:"maxNoWorkProduct"`
	MaxPoints               float64 `json:"maxPoints"`
	IsCompulsory            bool    `json:"isCompulsory"`
	EnforceWorkProductLimit bool    `json:"enforceWorkProductLimit"`
	Description             string  `json:"description"`
	GradeGroupID            int     `json:"gradeGroupId"`
	GradeGroupName          string  `json:"gradeGroupName"`
	CategoryName            string  `json:"categoryName"`
}

// CategoryDefinitionData2 is an alternative category definition DTO with
// required fields.
type CategoryDefinitionData2 struct {
	BaseWorkFlowVm
	DefinitionID            string  `json:"definitionId"`
	ObjectiveCategoryID     string  `json:"objectiveCategoryId"     validate:"required"`
	ReviewPeriodID          string  `json:"reviewPeriodId"          validate:"required"`
	Weight                  float64 `json:"weight"                  validate:"required"`
	MaxNoObjectives         int     `json:"maxNoObjectives"`
	MaxNoWorkProduct        int     `json:"maxNoWorkProduct"`
	MaxPoints               float64 `json:"maxPoints"`
	IsCompulsory            bool    `json:"isCompulsory"`
	EnforceWorkProductLimit bool    `json:"enforceWorkProductLimit"`
	Description             string  `json:"description"`
	GradeGroupID            int     `json:"gradeGroupId"`
}

// ===========================================================================
// Cascaded Objective Response VMs
// ===========================================================================

// CascadedObjectiveDataResponseVm wraps a single cascaded objective.
type CascadedObjectiveDataResponseVm struct {
	BaseAPIResponse
	Objective *CascadedObjectiveData `json:"objective"`
}

// CascadedObjectiveDataListResponseVm wraps a list of cascaded objectives.
type CascadedObjectiveDataListResponseVm struct {
	GenericListResponseVm
	Objectives []CascadedObjectiveData `json:"objectives"`
}

// CascadedObjectiveData is the flat DTO for a cascaded objective.
type CascadedObjectiveData struct {
	BaseWorkFlowVm
	ObjectiveID          string                   `json:"objectiveId"          validate:"required"`
	Name                 string                   `json:"name"                validate:"required"`
	Description          string                   `json:"description"`
	Kpi                  string                   `json:"kpi"                 validate:"required"`
	Target               string                   `json:"target"              validate:"required"`
	ObjectivesCategoryID string                   `json:"objectivesCategoryId" validate:"required"`
	StrategyID           string                   `json:"strategyId"          validate:"required"`
	CategoryDefinitions  []CategoryDefinitionData `json:"categoryDefinitions"`
}

// ===========================================================================
// Planned Objective Response VMs
// ===========================================================================

// PlannedOperationalObjectivesResponseVm wraps planned objectives.
type PlannedOperationalObjectivesResponseVm struct {
	GenericListResponseVm
	PlannedObjectives []PlannedObjectiveData `json:"plannedObjectives"`
}

// PlannedObjectiveResponseVm wraps a single planned objective.
type PlannedObjectiveResponseVm struct {
	BaseAPIResponse
	PlannedObjective *PlannedObjectiveData `json:"plannedObjective"`
}

// PlannedObjectiveData is the flat DTO for a planned objective.
type PlannedObjectiveData struct {
	PlannedObjectiveID  string                   `json:"plannedObjectiveId"`
	ReviewPeriodID      string                   `json:"reviewPeriodId"`
	ReviewPeriod        string                   `json:"reviewPeriod"`
	Year                int                      `json:"year"`
	ObjectiveLevel      string                   `json:"objectiveLevel"`
	ObjectiveID         string                   `json:"objectiveId"`
	Objective           string                   `json:"objective"`
	ObjectiveCategoryID string                   `json:"objectiveCategoryId"`
	RecordStatus        string                   `json:"recordStatus"`
	Status              enums.Status             `json:"status"`
	StaffID             string                   `json:"staffId"`
	Approver            string                   `json:"approver"`
	CreatedBy           string                   `json:"createdBy"`
	CreatedDate         *time.Time               `json:"createdDate"`
	IsApproved          bool                     `json:"isApproved"`
	IsRejected          bool                     `json:"isRejected"`
	IsActive            bool                     `json:"isActive"`
	Comment             string                   `json:"comment"`
	Kpi                 string                   `json:"kpi"`
	Target              string                   `json:"target"`
	Description         string                   `json:"description"`
	EnterpriseObjective *EnterpriseObjectiveData `json:"enterpriseObjective"`
}

// PlannedObjective is a lightweight planned objective DTO.
type PlannedObjective struct {
	PlannedObjectiveID string    `json:"plannedObjectiveId"`
	ReviewPeriodID     string    `json:"reviewPeriodId"`
	ReviewPeriod       string    `json:"reviewPeriod"`
	Year               int       `json:"year"`
	ObjectiveLevel     string    `json:"objectiveLevel"`
	ObjectiveID        string    `json:"objectiveId"`
	Objective          string    `json:"objective"`
	RecordStatus       string    `json:"recordStatus"`
	StaffID            string    `json:"staffId"`
	Approver           string    `json:"approver"`
	CreatedDate        time.Time `json:"createdDate"`
	IsApproved         bool      `json:"isApproved"`
	IsRejected         bool      `json:"isRejected"`
	IsActive           bool      `json:"isActive"`
}

// ===========================================================================
// Operational Objective Response VMs
// ===========================================================================

// OperationalObjectivesResponseVm wraps operational objectives.
type OperationalObjectivesResponseVm struct {
	GenericListResponseVm
	Objectives []OperationalObjectiveData `json:"objectives"`
}

// OperationalObjectiveData is the DTO for an operational objective.
type OperationalObjectiveData struct {
	ReviewPeriodID      string                   `json:"reviewPeriodId"`
	ReviewPeriod        string                   `json:"reviewPeriod"`
	ObjectiveLevel      string                   `json:"objectiveLevel"`
	ObjLevel            enums.ObjectiveLevel     `json:"objLevel"`
	ObjectiveID         string                   `json:"objectiveId"`
	Objective           string                   `json:"objective"`
	Kpi                 string                   `json:"kpi"`
	Target              string                   `json:"target"`
	Description         string                   `json:"description"`
	EnterpriseObjective *EnterpriseObjectiveData `json:"enterpriseObjective"`
	RecordStatus        string                   `json:"recordStatus"`
	StaffID             string                   `json:"staffId"`
	PlannedObjectiveID  string                   `json:"plannedObjectiveId"`
}

// OperationalObjective is a lightweight operational objective DTO.
type OperationalObjective struct {
	ReviewPeriodID string `json:"reviewPeriodId"`
	ReviewPeriod   string `json:"reviewPeriod"`
	ObjectiveLevel string `json:"objectiveLevel"`
	ObjectiveID    string `json:"objectiveId"`
	Objective      string `json:"objective"`
	RecordStatus   string `json:"recordStatus"`
	StaffID        string `json:"staffId"`
	IsActive       bool   `json:"isActive"`
}

// ===========================================================================
// Project Response VMs
// ===========================================================================

// ProjectResponseVm wraps a single project.
type ProjectResponseVm struct {
	BaseAPIResponse
	Project *ProjectData `json:"project"`
}

// ProjectListResponseVm wraps a list of project view models.
type ProjectListResponseVm struct {
	GenericListResponseVm
	Projects []ProjectViewModel `json:"projects"`
}

// ProjectAssignedListResponseVm wraps a list of project data.
type ProjectAssignedListResponseVm struct {
	GenericListResponseVm
	Projects []ProjectData `json:"projects"`
}

// ProjectViewModel is the read/write DTO for a project.
type ProjectViewModel struct {
	BaseWorkFlowVm
	ProjectID          string                 `json:"projectId"`
	ProjectManager     string                 `json:"projectManager"`
	ProjectManagerName string                 `json:"projectManagerName"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	StartDate          time.Time              `json:"startDate"`
	EndDate            time.Time              `json:"endDate"`
	Deliverables       string                 `json:"deliverables"`
	ReviewPeriodID     string                 `json:"reviewPeriodId"`
	Department         string                 `json:"department"`
	DepartmentID       int                    `json:"departmentId"`
	CreatedBy          string                 `json:"createdBy"`
	ProjectObjectives  []ProjectObjectiveData `json:"projectObjectives"`
}

// BaseProjectData carries shared fields for project/committee data DTOs.
type BaseProjectData struct {
	BaseWorkFlowVm
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"reviewPeriodId"`
	DepartmentID   int       `json:"departmentId"`
}

// ProjectData is the detailed project DTO.
type ProjectData struct {
	BaseProjectData
	ProjectID                  string                         `json:"projectId"`
	ProjectManager             string                         `json:"projectManager"`
	ProjectManagerName         string                         `json:"projectManagerName"`
	DepartmentName             string                         `json:"departmentName"`
	ProjectAssignedWorkProducts []ProjectAssignedWorkProductData `json:"projectAssignedWorkProducts"`
	ProjectWorkProducts        []ProjectWorkProductDataResponse `json:"projectWorkProducts"`
	ProjectMembers             []ProjectMemberData            `json:"projectMembers"`
	ProjectObjectives          []ProjectObjectiveData         `json:"projectObjectives"`
}

// ProjectObjectiveListResponseVm wraps project objectives.
type ProjectObjectiveListResponseVm struct {
	GenericListResponseVm
	ProjectObjectives []ProjectObjectiveData `json:"projectObjectives"`
}

// ProjectObjectiveData is the DTO for a project objective.
type ProjectObjectiveData struct {
	BaseEntityVm
	ProjectObjectiveID string `json:"projectObjectiveId"`
	ObjectiveID        string `json:"objectiveId"`
	Objective          string `json:"objective"`
	Kpi                string `json:"kpi"`
	ProjectID          string `json:"projectId"`
	RecordStatusName   string `json:"recordStatusName"`
}

// ProjectMemberListResponseVm wraps project members.
type ProjectMemberListResponseVm struct {
	GenericListResponseVm
	ProjectMembers []ProjectMemberData `json:"projectMembers"`
}

// ProjectMemberData is the DTO for a project member.
type ProjectMemberData struct {
	BaseWorkFlowVm
	ProjectMemberID    string `json:"projectMemberId"`
	StaffID            string `json:"staffId"`
	StaffName          string `json:"staffName"`
	ProjectID          string `json:"projectId"`
	PlannedObjectiveID string `json:"plannedObjectiveId"`
	ObjectiveName      string `json:"objectiveName"`
}

// ProjectAssignedWorkProductListResponseVm wraps project assigned work products.
type ProjectAssignedWorkProductListResponseVm struct {
	GenericListResponseVm
	ProjectWorkProducts []ProjectAssignedWorkProductData `json:"projectWorkProducts"`
}

// ProjectAssignedWorkProductResponseVm wraps a single assigned work product.
type ProjectAssignedWorkProductResponseVm struct {
	BaseAPIResponse
	WorkProduct *ProjectAssignedWorkProductData `json:"workProduct"`
}

// ProjectAssignedWorkProductData is the DTO for a project assigned work product.
type ProjectAssignedWorkProductData struct {
	BaseEntityVm
	ProjectAssignedWorkProductID string    `json:"projectAssignedWorkProductId"`
	WorkProductDefinitionID      string    `json:"workProductDefinitionId"`
	WorkProductID                string    `json:"workProductId"`
	WorkProductType              int       `json:"workProductType"`
	Name                         string    `json:"name"`
	Description                  string    `json:"description"`
	ProjectID                    string    `json:"projectId"`
	ProjectManager               string    `json:"projectManager"`
	ReviewPeriodID               string    `json:"reviewPeriodId"`
	StartDate                    time.Time `json:"startDate"`
	EndDate                      time.Time `json:"endDate"`
	CompletionDate               time.Time `json:"completionDate"`
	Deliverables                 string    `json:"deliverables"`
	RejectionReason              string    `json:"rejectionReason"`
}

// ProjectWorkProductListResponseVm wraps project work products.
type ProjectWorkProductListResponseVm struct {
	GenericListResponseVm
	ProjectWorkProducts []ProjectWorkProductDataResponse `json:"projectWorkProducts"`
}

// ProjectWorkProductDataResponse is the DTO for a project work product.
type ProjectWorkProductDataResponse struct {
	BaseEntityVm
	ProjectWorkProductID string           `json:"projectWorkProductId"`
	WorkProductID        string           `json:"workProductId"`
	ProjectID            string           `json:"projectId"`
	AssignedEvaluator    string           `json:"assignedEvaluator"`
	WorkProduct          *WorkProductData `json:"workProduct"`
	Project              *ProjectData     `json:"project"`
}

// ===========================================================================
// Committee Response VMs
// ===========================================================================

// CommitteeResponseVm wraps a single committee.
type CommitteeResponseVm struct {
	BaseAPIResponse
	Committee *CommitteeData `json:"committee"`
}

// CommitteeListResponseVm wraps a list of committee view models.
type CommitteeListResponseVm struct {
	GenericListResponseVm
	Committees []CommitteeViewModel `json:"committees"`
}

// CommitteeAssignedListResponseVm wraps a list of committee data.
type CommitteeAssignedListResponseVm struct {
	GenericListResponseVm
	Committees []CommitteeData `json:"committees"`
}

// CommitteeViewModel is the read/write DTO for a committee.
type CommitteeViewModel struct {
	BaseWorkFlowVm
	CommitteeID        string                   `json:"committeeId"`
	Chairperson        string                   `json:"chairperson"`
	ChairpersonName    string                   `json:"chairpersonName"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	StartDate          time.Time                `json:"startDate"`
	EndDate            time.Time                `json:"endDate"`
	Deliverables       string                   `json:"deliverables"`
	ReviewPeriodID     string                   `json:"reviewPeriodId"`
	Department         string                   `json:"department"`
	DepartmentID       int                      `json:"departmentId"`
	CreatedBy          string                   `json:"createdBy"`
	CommitteeObjectives []CommitteeObjectiveData `json:"committeeObjectives"`
}

// CommitteeData is the detailed committee DTO.
type CommitteeData struct {
	BaseProjectData
	CommitteeID                    string                            `json:"committeeId"`
	Chairperson                    string                            `json:"chairperson"`
	ChairpersonName                string                            `json:"chairpersonName"`
	DepartmentName                 string                            `json:"departmentName"`
	CommitteeAssignedWorkProducts  []CommitteeAssignedWorkProductData `json:"committeeAssignedWorkProducts"`
	CommitteeWorkProducts          []CommitteeWorkProductDataResponse `json:"committeeWorkProducts"`
	CommitteeMembers               []CommitteeMemberData             `json:"committeeMembers"`
	CommitteeObjectives            []CommitteeObjectiveData          `json:"committeeObjectives"`
}

// CommitteeObjectiveListResponseVm wraps committee objectives.
type CommitteeObjectiveListResponseVm struct {
	GenericListResponseVm
	CommitteeObjectives []CommitteeObjectiveData `json:"committeeObjectives"`
}

// CommitteeObjectiveData is the DTO for a committee objective.
type CommitteeObjectiveData struct {
	BaseEntityVm
	CommitteeObjectiveID string `json:"committeeObjectiveId"`
	ObjectiveID          string `json:"objectiveId"`
	Objective            string `json:"objective"`
	Kpi                  string `json:"kpi"`
	CommitteeID          string `json:"committeeId"`
	RecordStatusName     string `json:"recordStatusName"`
}

// CommitteeMemberListResponseVm wraps committee members.
type CommitteeMemberListResponseVm struct {
	GenericListResponseVm
	CommitteeMembers []CommitteeMemberData `json:"committeeMembers"`
}

// CommitteeMemberData is the DTO for a committee member.
type CommitteeMemberData struct {
	BaseWorkFlowVm
	CommitteeMemberID  string `json:"committeeMemberId"`
	StaffID            string `json:"staffId"`
	StaffName          string `json:"staffName"`
	CommitteeID        string `json:"committeeId"`
	PlannedObjectiveID string `json:"plannedObjectiveId"`
	ObjectiveName      string `json:"objectiveName"`
}

// CommitteeAssignedWorkProductListResponseVm wraps committee assigned work products.
type CommitteeAssignedWorkProductListResponseVm struct {
	GenericListResponseVm
	CommitteeWorkProducts []CommitteeAssignedWorkProductData `json:"committeeWorkProducts"`
}

// CommitteeAssignedWorkProductResponseVm wraps a single committee assigned work product.
type CommitteeAssignedWorkProductResponseVm struct {
	BaseAPIResponse
	WorkProduct *CommitteeAssignedWorkProductData `json:"workProduct"`
}

// CommitteeAssignedWorkProductData is the DTO for a committee assigned work product.
type CommitteeAssignedWorkProductData struct {
	BaseEntityVm
	CommitteeAssignedWorkProductID string    `json:"committeeAssignedWorkProductId"`
	WorkProductDefinitionID        string    `json:"workProductDefinitionId"`
	Name                           string    `json:"name"`
	Description                    string    `json:"description"`
	CommitteeID                    string    `json:"committeeId"`
	Chairperson                    string    `json:"chairperson"`
	ReviewPeriodID                 string    `json:"reviewPeriodId"`
	StartDate                      time.Time `json:"startDate"`
	EndDate                        time.Time `json:"endDate"`
	Deliverables                   string    `json:"deliverables"`
	RejectionReason                string    `json:"rejectionReason"`
}

// CommitteeWorkProductListResponseVm wraps committee work products.
type CommitteeWorkProductListResponseVm struct {
	GenericListResponseVm
	CommitteeWorkProducts []CommitteeWorkProductDataResponse `json:"committeeWorkProducts"`
}

// CommitteeWorkProductDataResponse is the DTO for a committee work product.
type CommitteeWorkProductDataResponse struct {
	BaseEntityVm
	CommitteeWorkProductID string           `json:"committeeWorkProductId"`
	WorkProductID          string           `json:"workProductId"`
	CommitteeID            string           `json:"committeeId"`
	AssignedEvaluator      string           `json:"assignedEvaluator"`
	WorkProduct            *WorkProductData `json:"workProduct"`
	Committee              *CommitteeData   `json:"committee"`
}

// AdhocStaffResponseVm wraps an employee with a planned objective.
type AdhocStaffResponseVm struct {
	BaseAPIResponse
	Employee           interface{} `json:"employee"`
	PlannedObjectiveID string      `json:"plannedObjectiveId"`
}

// ===========================================================================
// Work Product Response VMs
// ===========================================================================

// WorkProductResponseVm wraps a single work product.
type WorkProductResponseVm struct {
	BaseAPIResponse
	WorkProduct *WorkProductData `json:"workProduct"`
}

// StaffWorkProductListResponseVm wraps work products for a staff member.
type StaffWorkProductListResponseVm struct {
	GenericListResponseVm
	WorkProducts []WorkProductData `json:"workProducts"`
}

// WorkProductData is the detailed work product DTO.
type WorkProductData struct {
	BaseWorkFlowVm
	WorkProductID           string                  `json:"workProductId"           validate:"required"`
	Name                    string                  `json:"name"                    validate:"required"`
	Description             string                  `json:"description"`
	MaxPoint                float64                 `json:"maxPoint"`
	WorkProductType         int                     `json:"workProductType"`
	WorkProductTypeName     string                  `json:"workProductTypeName"`
	IsSelfCreated           bool                    `json:"isSelfCreated"`
	StaffID                 string                  `json:"staffId"`
	AcceptanceComment       string                  `json:"acceptanceComment"`
	StartDate               time.Time               `json:"startDate"`
	EndDate                 time.Time               `json:"endDate"`
	Deliverables            string                  `json:"deliverables"`
	FinalScore              float64                 `json:"finalScore"`
	NoReturned              int                     `json:"noReturned"`
	CompletionDate          time.Time               `json:"completionDate"`
	ApproverComment         string                  `json:"approverComment"`
	ReviewPeriodID          string                  `json:"reviewPeriodId"`
	PlannedObjectiveID      string                  `json:"plannedObjectiveId"`
	ObjectiveID             string                  `json:"objectiveId"`
	ObjectiveName           string                  `json:"objectiveName"`
	WorkProductDefinitionID string                  `json:"workProductDefinitionId"`
	AssignedEvaluator       string                  `json:"assignedEvaluator"`
	WorkProductTasks        []WorkProductTaskDetail `json:"workProductTasks"`
}

// WorkProductTaskResponseVm wraps a single work product task.
type WorkProductTaskResponseVm struct {
	BaseAPIResponse
	WorkProductTask *WorkProductTaskDetail `json:"workProductTask"`
}

// WorkProductTaskListResponseVm wraps a list of work product tasks.
type WorkProductTaskListResponseVm struct {
	GenericListResponseVm
	WorkProductTasks []WorkProductTaskData `json:"workProductTasks"`
}

// WorkProductTaskDetail is the detailed work product task DTO.
type WorkProductTaskDetail struct {
	BaseEntityVm
	WorkProductTaskID string    `json:"workProductTaskId" validate:"required"`
	Name              string    `json:"name"              validate:"required"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"startDate"`
	EndDate           time.Time `json:"endDate"`
	CompletionDate    time.Time `json:"completionDate"`
	WorkProductID     string    `json:"workProductId"     validate:"required"`
}

// WorkProductTaskData is the work product task DTO with navigation.
type WorkProductTaskData struct {
	BaseEntityVm
	WorkProductTaskID string    `json:"workProductTaskId"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	CompletionDate    time.Time `json:"completionDate"`
	WorkProductID     string    `json:"workProductId"`
}

// WorkProductDashDetails is the work product dashboard detail DTO.
type WorkProductDashDetails struct {
	BaseWorkFlowVm
	WorkProductID            string                  `json:"workProductId"     validate:"required"`
	Name                     string                  `json:"name"              validate:"required"`
	Description              string                  `json:"description"`
	MaxPoint                 float64                 `json:"maxPoint"`
	WorkProductType          int                     `json:"workProductType"`
	WorkProductTypeName      string                  `json:"workProductTypeName"`
	IsSelfCreated            bool                    `json:"isSelfCreated"`
	StaffID                  string                  `json:"staffId"`
	AcceptanceComment        string                  `json:"acceptanceComment"`
	StartDate                time.Time               `json:"startDate"`
	EndDate                  time.Time               `json:"endDate"`
	Deliverables             string                  `json:"deliverables"`
	FinalScore               float64                 `json:"finalScore"`
	NoReturned               int                     `json:"noReturned"`
	CompletionDate           time.Time               `json:"completionDate"`
	ApproverComment          string                  `json:"approverComment"`
	ReviewPeriodID           string                  `json:"reviewPeriodId"`
	PlannedObjectiveID       string                  `json:"plannedObjectiveId"`
	ObjectiveID              string                  `json:"objectiveId"`
	ObjectiveName            string                  `json:"objectiveName"`
	WorkProductDefinitionID  string                  `json:"workProductDefinitionId"`
	AssignedEvaluator        string                  `json:"assignedEvaluator"`
	PercentageTaskCompletion float64                 `json:"percentageTaskCompletion"`
	TasksCompleted           int                     `json:"tasksCompleted"`
	TotalTasks               int                     `json:"totalTasks"`
	WorkProductTasks         []WorkProductTaskDetail `json:"workProductTasks"`
}

// ===========================================================================
// Work Product Evaluation Response VMs
// ===========================================================================

// WorkProductEvaluationResponseVm wraps a single work product evaluation.
type WorkProductEvaluationResponseVm struct {
	BaseAPIResponse
	WorkProductEvaluation *WorkProductEvaluationDataResponse `json:"workProductEvaluation"`
}

// WorkProductEvaluationDataResponse is the detailed evaluation DTO.
type WorkProductEvaluationDataResponse struct {
	BaseEntityVm
	WorkProductEvaluationID      string              `json:"workProductEvaluationId" validate:"required"`
	WorkProductID                string              `json:"workProductId"           validate:"required"`
	Timeliness                   float64             `json:"timeliness"              validate:"required"`
	TimelinessEvaluationOptionID string              `json:"timelinessEvaluationOptionId"`
	Quality                      float64             `json:"quality"                 validate:"required"`
	QualityEvaluationOptionID    string              `json:"qualityEvaluationOptionId"`
	Output                       float64             `json:"output"                  validate:"required"`
	OutputEvaluationOptionID     string              `json:"outputEvaluationOptionId"`
	Outcome                      float64             `json:"outcome"`
	EvaluatorStaffID             string              `json:"evaluatorStaffId"`
	IsReEvaluated                bool                `json:"isReEvaluated"`
	WorkProduct                  *WorkProductData    `json:"workProduct"`
	TimelinessEvaluationOption   *EvaluationOptionData `json:"timelinessEvaluationOption"`
	QualityEvaluationOption      *EvaluationOptionData `json:"qualityEvaluationOption"`
	OutputEvaluationOption       *EvaluationOptionData `json:"outputEvaluationOption"`
}

// EvaluationOptionData is the flat DTO for an evaluation option.
type EvaluationOptionData struct {
	EvaluationOptionID string       `json:"evaluationOptionId" validate:"required"`
	Name               string       `json:"name"               validate:"required"`
	Description        string       `json:"description"`
	RecordStatus       enums.Status `json:"recordStatus"       validate:"required"`
	Score              float64      `json:"score"              validate:"required"`
	EvaluationType     int          `json:"evaluationType"     validate:"required"`
}

// EvaluationOptionResponseVm wraps a list of evaluation options.
type EvaluationOptionResponseVm struct {
	GenericListResponseVm
	EvaluationOptions []EvaluationOptionVm `json:"evaluationOptions"`
}

// EvaluationResponseVm wraps an evaluation with a score.
type EvaluationResponseVm struct {
	BaseAPIResponse
	IsSuccess bool    `json:"isSuccess"`
	Message   string  `json:"message"`
	Score     float64 `json:"score"`
}

// RecalculateWorkProductResponseVm wraps a recalculation response.
type RecalculateWorkProductResponseVm struct {
	BaseAPIResponse
	StaffID        string `json:"staffId"`
	ReviewPeriodID string `json:"reviewPeriodId"`
}

// ===========================================================================
// Operational Objective Work Product Response VMs
// ===========================================================================

// OperationalObjectiveWorkProductListResponseVm wraps operational objective work products.
type OperationalObjectiveWorkProductListResponseVm struct {
	GenericListResponseVm
	OperationalObjectiveWorkProducts []OperationalObjectiveWorkProductData `json:"operationalObjectiveWorkProducts"`
}

// CustomObjectiveWorkProductListResponseVm wraps custom work products.
type CustomObjectiveWorkProductListResponseVm struct {
	GenericListResponseVm
	WorkProducts []CustomWorkProductVm `json:"workProducts"`
}

// OperationalObjectiveWorkProductData is the DTO for an operational objective work product.
type OperationalObjectiveWorkProductData struct {
	BaseEntityVm
	OperationalObjectiveWorkProductID string `json:"operationalObjectiveWorkProductId"`
	WorkProductID                     string `json:"workProductId"`
	WorkProductDefinitionID           string `json:"workProductDefinitionId"`
	PlannedObjectiveID                string `json:"plannedObjectiveId"`
	ObjectiveID                       string `json:"objectiveId"`
	WorkProductEvaluationID           string `json:"workProductEvaluationId"`
}

// ObjectiveWorkProductListResponseVm wraps objective work products.
type ObjectiveWorkProductListResponseVm struct {
	GenericListResponseVm
	ObjectiveWorkProducts []ObjectiveWorkProductData `json:"objectiveWorkProducts"`
}

// ObjectiveWorkProductData is the DTO for an objective work product.
type ObjectiveWorkProductData struct {
	ReviewPeriodID          string `json:"reviewPeriodId"`
	ReviewPeriod            string `json:"reviewPeriod"`
	ObjectiveLevel          string `json:"objectiveLevel"`
	ObjectiveID             string `json:"objectiveId"`
	Objective               string `json:"objective"`
	WorkProductDefinitionID string `json:"workProductDefinitionId"`
	WorkProductName         string `json:"workProductName"`
	Description             string `json:"description"`
	Deliverables            string `json:"deliverables"`
	StaffID                 string `json:"staffId"`
}

// WorkProductDefinitionResponseVm wraps work product definitions.
type WorkProductDefinitionResponseVm struct {
	GenericListResponseVm
	WorkProductDefinitions []WorkProductDefinitionVm `json:"workProductDefinitions"`
}

// PaginatedWorkProductDefinitionResponseVm wraps paginated work product definitions.
type PaginatedWorkProductDefinitionResponseVm struct {
	BaseAPIResponse
	WorkProductDefinitions *PaginatedResult[WorkProductDefinitionVm] `json:"workProductDefinitions"`
}

// ===========================================================================
// Dashboard Response VMs
// ===========================================================================

// ReviewPeriodWorkProductDashboardResponseVm is the work product dashboard summary.
type ReviewPeriodWorkProductDashboardResponseVm struct {
	BaseAPIResponse
	StaffID                           string `json:"staffId"`
	ReviewPeriodID                    string `json:"reviewPeriodId"`
	NoAllWorkProducts                 int    `json:"noAllWorkProducts"`
	TotalWorkProductTasks             int    `json:"totalWorkProductTasks"`
	NoActiveWorkProducts              int    `json:"noActiveWorkProducts"`
	NoWorkProductsAwaitingEvaluation  int    `json:"noWorkProductsAwaitingEvaluation"`
	NoWorkProductsClosed              int    `json:"noWorkProductsClosed"`
	NoWorkProductsPendingApproval     int    `json:"noWorkProductsPendingApproval"`
}

// ReviewPeriodWorkProductDetailsDashboardResponseVm is the detailed work product dashboard.
type ReviewPeriodWorkProductDetailsDashboardResponseVm struct {
	BaseAPIResponse
	StaffID                          string                   `json:"staffId"`
	ReviewPeriodID                   string                   `json:"reviewPeriodId"`
	WorkProducts                     []WorkProductDashDetails `json:"workProducts"`
	NoAllWorkProducts                int                      `json:"noAllWorkProducts"`
	TotalWorkProductTasks            int                      `json:"totalWorkProductTasks"`
	NoActiveWorkProducts             int                      `json:"noActiveWorkProducts"`
	NoWorkProductsAwaitingEvaluation int                      `json:"noWorkProductsAwaitingEvaluation"`
	NoWorkProductsClosed             int                      `json:"noWorkProductsClosed"`
	NoWorkProductsPendingApproval    int                      `json:"noWorkProductsPendingApproval"`
}

// ReviewPeriodPointsDashboardResponseVm is the points dashboard response.
type ReviewPeriodPointsDashboardResponseVm struct {
	BaseAPIResponse
	StaffID            string  `json:"staffId"`
	ReviewPeriodID     string  `json:"reviewPeriodId"`
	MaxPoints          float64 `json:"maxPoints"`
	AccumulatedPoints  float64 `json:"accumulatedPoints"`
	DeductedPoints     float64 `json:"deductedPoints"`
	ActualPoints       float64 `json:"actualPoints"`
}

// LeaveResponseVm wraps leave days count.
type LeaveResponseVm struct {
	BaseAPIResponse
	NoLeaveDays int `json:"noLeaveDays"`
}

// PublicHolidaysResponseVm wraps public holiday days count.
type PublicHolidaysResponseVm struct {
	BaseAPIResponse
	NoPublicDays int `json:"noPublicDays"`
}

// ===========================================================================
// Organogram Performance Summary VMs
// ===========================================================================

// OrganogramPerformanceSummaryResponseVm wraps organogram performance summary.
type OrganogramPerformanceSummaryResponseVm struct {
	BaseAPIResponse
	ReferenceID                          string               `json:"referenceId"`
	ManagerID                            string               `json:"managerId"`
	ReferenceName                        string               `json:"referenceName"`
	ReviewPeriodID                       string               `json:"reviewPeriodId"`
	ReviewPeriod                         string               `json:"reviewPeriod"`
	ReviewPeriodShortName                string               `json:"reviewPeriodShortName"`
	MaxPoint                             float64              `json:"maxPoint"`
	Year                                 int                  `json:"year"`
	ActualScore                          float64              `json:"actualScore"`
	PerformanceScore                     float64              `json:"performanceScore"`
	TotalWorkProducts                    int                  `json:"totalWorkProducts"`
	TotalStaff                           int                  `json:"totalStaff"`
	LivingTheValuesRating                float64              `json:"livingtheValuesRating"`
	TotalWorkProductsCompletedOnSchedule int                  `json:"totalWorkProductsCompletedOnSchedule"`
	TotalWorkProductsBehindSchedule      int                  `json:"totalWorkProductsBehindSchedule"`
	Total360Feedbacks                    int                  `json:"total360Feedbacks"`
	Completed360FeedbacksToTreat         int                  `json:"completed360FeedbacksToTreat"`
	Pending360FeedbacksToTreat           int                  `json:"pending360FeedbacksToTreat"`
	TotalCompetencyGaps                  int                  `json:"totalCompetencyGaps"`
	PercentageGapsClosure                float64              `json:"percentageGapsClosure"`
	PercentageWorkProductsClosed         float64              `json:"percentageWorkProductsClosed"`
	PercentageWorkProductsPending        float64              `json:"percentageWorkProductsPending"`
	OrganogramLevel                      enums.OrganogramLevel `json:"organogramLevel"`
	EarnedPerformanceGrade               string               `json:"earnedPerformanceGrade"`
}

// OrganogramPerformanceSummaryDetails is the detail DTO for org performance.
type OrganogramPerformanceSummaryDetails struct {
	ReferenceID                          string  `json:"referenceId"`
	ReferenceName                        string  `json:"referenceName"`
	ManagerID                            string  `json:"managerId"`
	ManagerName                          string  `json:"managerName"`
	ActualScore                          float64 `json:"actualScore"`
	PerformanceScore                     float64 `json:"performanceScore"`
	TotalWorkProducts                    int     `json:"totalWorkProducts"`
	TotalStaff                           int     `json:"totalStaff"`
	LivingTheValuesRating                float64 `json:"livingtheValuesRating"`
	TotalWorkProductsCompletedOnSchedule int     `json:"totalWorkProductsCompletedOnSchedule"`
	TotalWorkProductsBehindSchedule      int     `json:"totalWorkProductsBehindSchedule"`
	Total360Feedbacks                    int     `json:"total360Feedbacks"`
	Completed360FeedbacksToTreat         int     `json:"completed360FeedbacksToTreat"`
	Pending360FeedbacksToTreat           int     `json:"pending360FeedbacksToTreat"`
	TotalCompetencyGaps                  int     `json:"totalCompetencyGaps"`
	PercentageGapsClosure                float64 `json:"percentageGapsClosure"`
	PercentageWorkProductsClosed         float64 `json:"percentageWorkProductsClosed"`
	PercentageWorkProductsPending        float64 `json:"percentageWorkProductsPending"`
	EarnedPerformanceGrade               string  `json:"earnedPerformanceGrade"`
}

// OrganogramPerformanceSummaryListResponseVm wraps org performance list.
type OrganogramPerformanceSummaryListResponseVm struct {
	GenericListResponseVm
	HeadOfUnitID              string                                `json:"headOfUnitId"`
	ReviewPeriodID            string                                `json:"reviewPeriodId"`
	ReviewPeriod              string                                `json:"reviewPeriod"`
	ReviewPeriodShortName     string                                `json:"reviewPeriodShortName"`
	MaxPoint                  float64                               `json:"maxPoint"`
	Year                      int                                   `json:"year"`
	OrganogramLevel           enums.OrganogramLevel                 `json:"organogramLevel"`
	OrganogramPerformances    []OrganogramPerformanceSummaryDetails `json:"organogramPerformances"`
}

// ===========================================================================
// Staff Score Card VMs
// ===========================================================================

// StaffScoreCardResponseVm wraps a single staff score card.
type StaffScoreCardResponseVm struct {
	BaseAPIResponse
	ScoreCard *StaffScoreCardDetails `json:"scoreCard"`
}

// AllStaffScoreCardResponseVm wraps all staff score cards.
type AllStaffScoreCardResponseVm struct {
	BaseAPIResponse
	StaffScoreCards []StaffScoreCardDetails `json:"staffScoreCards"`
}

// StaffScoreCardDetails is the detailed staff score card DTO.
type StaffScoreCardDetails struct {
	StaffID                              string                              `json:"staffId"`
	StaffName                            string                              `json:"staffName"`
	ReviewPeriodID                       string                              `json:"reviewPeriodId"`
	ReviewPeriod                         string                              `json:"reviewPeriod"`
	ReviewPeriodShortName                string                              `json:"reviewPeriodShortName"`
	Year                                 int                                 `json:"year"`
	TotalWorkProducts                    int                                 `json:"totalWorkProducts"`
	PercentageWorkProductsCompletion     float64                             `json:"percentageWorkProductsCompletion"`
	TotalCompetencyGaps                  int                                 `json:"totalCompetencyGaps"`
	TotalCompetencyGapsClosed            int                                 `json:"totalCompetencyGapsClosed"`
	PercentageGapsClosure                float64                             `json:"percentageGapsClosure"`
	PercentageGapsClosureScore           float64                             `json:"percentageGapsClosureScore"`
	MaxPoints                            float64                             `json:"maxPoints"`
	AccumulatedPoints                    float64                             `json:"accumulatedPoints"`
	DeductedPoints                       float64                             `json:"deductedPoints"`
	ActualPoints                         float64                             `json:"actualPoints"`
	PercentageScore                      float64                             `json:"percentageScore"`
	PmsCompetencyCategory                map[string]float64                  `json:"pmsCompetencyCategory"`
	StaffPerformanceGrade                string                              `json:"staffPerformanceGrade"`
	TotalWorkProductsCompletedOnSchedule int                                 `json:"totalWorkProductsCompletedOnSchedule"`
	TotalWorkProductsBehindSchedule      int                                 `json:"totalWorkProductsBehindSchedule"`
	PmsCompetencies                      []StaffLivingTheValueRatingsDetails `json:"pmsCompetencies"`
}

// StaffLivingTheValueRatingsDetails is the DTO for staff LTV ratings.
type StaffLivingTheValueRatingsDetails struct {
	StaffID             string  `json:"staffId"`
	ReviewPeriodID      string  `json:"reviewPeriodId"`
	PmsCompetencyID     string  `json:"pmsCompetencyId"`
	ObjectiveCategoryID string  `json:"objectiveCategoryId"`
	PmsCompetency       string  `json:"pmsCompetency"`
	RatingScore         float64 `json:"ratingScore"`
}

// StaffAnnualScoreCardResponseVm wraps annual score cards.
type StaffAnnualScoreCardResponseVm struct {
	BaseAPIResponse
	StaffID    string                  `json:"staffId"`
	Year       int                     `json:"year"`
	ScoreCards []StaffScoreCardDetails `json:"scoreCards"`
}

// ===========================================================================
// Period Score Response VMs
// ===========================================================================

// PeriodScoreResponseVm wraps a single period score.
type PeriodScoreResponseVm struct {
	BaseAPIResponse
	PeriodScore *PeriodScoreData `json:"periodScore"`
}

// PeriodScoreListResponseVm wraps a list of period scores.
type PeriodScoreListResponseVm struct {
	GenericListResponseVm
	PeriodScores []PeriodScoreData `json:"periodScores"`
}

// PeriodScoreData is the detailed period score DTO.
type PeriodScoreData struct {
	BaseEntityVm
	PeriodScoreID      string  `json:"periodScoreId"`
	ReviewPeriodID     string  `json:"reviewPeriodId"`
	ReviewPeriod       string  `json:"reviewPeriod"`
	Year               int     `json:"year"`
	StaffID            string  `json:"staffId"`
	StaffFullName      string  `json:"staffFullName"`
	FinalScore         float64 `json:"finalScore"`
	MaxPoint           float64 `json:"maxPoint"`
	ScorePercentage    float64 `json:"scorePercentage"`
	FinalGrade         int     `json:"finalGrade"`
	FinalGradeName     string  `json:"finalGradeName"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
	OfficeID           int     `json:"officeID"`
	OfficeCode         string  `json:"officeCode"`
	OfficeName         string  `json:"officeName"`
	DivisionID         int     `json:"divisionId"`
	DivisionCode       string  `json:"divisionCode"`
	DivisionName       string  `json:"divisionName"`
	DepartmentID       int     `json:"departmentId"`
	DepartmentCode     string  `json:"departmentCode"`
	DepartmentName     string  `json:"departmentName"`
	MinNoOfObjectives  int     `json:"minNoOfObjectives"`
	MaxNoOfObjectives  int     `json:"maxNoOfObjectives"`
	StrategyID         string  `json:"strategyId"`
	StrategyName       string  `json:"strategyName"`
	StaffGrade         string  `json:"staffGrade"`
	LocationID         string  `json:"locationId"`
	HRDDeductedPoints  float64 `json:"hrddeductedPoints"`
	IsUnderPerforming  bool    `json:"isUnderPerforming"`
}

// ===========================================================================
// Audit Response VMs
// NOTE: AuditLogListResponseVm and AuditLogResponseVm are defined in dto_feedback.go.
// ===========================================================================

// ===========================================================================
// Feedback/360 Response VMs
// NOTE: All competency review feedback, reviewer, rating, questionnaire,
//       and PMS competency response/data types are defined in dto_feedback.go.
// ===========================================================================

// PmsCompetencyVmResponse is the PMS competency VM used in list responses.
type PmsCompetencyVmResponse struct {
	BaseWorkFlowVm
	PmsCompetencyID  string `json:"pmsCompetencyId"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ObjectCategoryID string `json:"objectCategoryId"`
}

// FeedbackQuestionaireVmResponse is the feedback questionnaire VM used in responses.
type FeedbackQuestionaireVmResponse struct {
	BaseWorkFlowVm
	FeedbackQuestionaireID string `json:"feedbackQuestionaireId"`
	Question               string `json:"question"`
	Description            string `json:"description"`
	PmsCompetencyID        string `json:"pmsCompetencyId"`
}

// ===========================================================================
// Period Objective Evaluation Response VMs
// ===========================================================================

// PeriodObjectiveEvaluationData is the evaluation data DTO.
type PeriodObjectiveEvaluationData struct {
	BaseWorkFlowVm
	PeriodObjectiveEvaluationID string  `json:"periodObjectiveEvaluationId"`
	TotalOutcomeScore           float64 `json:"totalOutcomeScore"`
	OutcomeScore                float64 `json:"outcomeScore"`
	PeriodObjectiveID           string  `json:"periodObjectiveId"`
	EnterpriseObjectiveID       string  `json:"enterpriseObjectiveId"`
	EnterpriseObjective         string  `json:"enterpriseObjective"`
	ReviewPeriodID              string  `json:"reviewPeriodId"`
	ReviewPeriod                string  `json:"reviewPeriod"`
}

// PeriodObjectiveEvaluationListResponseVm wraps evaluation list.
type PeriodObjectiveEvaluationListResponseVm struct {
	GenericListResponseVm
	ObjectiveEvaluations []PeriodObjectiveEvaluationData `json:"objectiveEvaluations"`
}

// PeriodObjectiveEvaluationResponseVm wraps a single evaluation.
type PeriodObjectiveEvaluationResponseVm struct {
	BaseAPIResponse
	ObjectiveEvaluation *PeriodObjectiveEvaluationData `json:"objectiveEvaluation"`
}

// PeriodObjectiveDepartmentEvaluationData is the department evaluation DTO.
type PeriodObjectiveDepartmentEvaluationData struct {
	BaseWorkFlowVm
	PeriodObjectiveDepartmentEvaluationID string  `json:"periodObjectiveDepartmentEvaluationId"`
	OverallOutcomeScored                  float64 `json:"overallOutcomeScored"`
	AllocatedOutcome                      float64 `json:"allocatedOutcome"`
	OutcomeScore                          float64 `json:"outcomeScore"`
	DepartmentID                          int     `json:"departmentId"`
	DepartmentName                        string  `json:"departmentName"`
	PeriodObjectiveID                     string  `json:"periodObjectiveId"`
	EnterpriseObjectiveID                 string  `json:"enterpriseObjectiveId"`
	DepartmentObjectiveID                 string  `json:"departmentObjectiveId"`
	DepartmentObjective                   string  `json:"departmentObjective"`
	ReviewPeriodID                        string  `json:"reviewPeriodId"`
	ReviewPeriod                          string  `json:"reviewPeriod"`
}

// PeriodObjectiveDepartmentEvaluationListResponseVm wraps department evaluations.
type PeriodObjectiveDepartmentEvaluationListResponseVm struct {
	GenericListResponseVm
	DepartmentObjectiveEvaluations []PeriodObjectiveDepartmentEvaluationData `json:"departmentObjectiveEvaluations"`
}

// PeriodObjectiveDepartmentEvaluationResponseVm wraps a single department evaluation.
type PeriodObjectiveDepartmentEvaluationResponseVm struct {
	BaseAPIResponse
	DepartmentObjectiveEvaluation *PeriodObjectiveDepartmentEvaluationData `json:"departmentObjectiveEvaluation"`
}

// DepartmentListResponseVm wraps departments for an enterprise objective.
type DepartmentListResponseVm struct {
	GenericListResponseVm
	EnterpriseObjectiveID string        `json:"enterpriseObjectiveId"`
	Departments           []interface{} `json:"departments"`
}
