package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ===========================================================================
// Frontend Custom VMs (FrontendCustomVms/ResponseFormVms.cs)
//
// These DTOs are purpose-built for frontend forms and composite responses.
// Types already defined in other DTO files are not repeated here.
// ===========================================================================

// ---------------------------------------------------------------------------
// Generic / Utility frontend VMs
// ---------------------------------------------------------------------------

// CustomGenericListVm is a generic list response with untyped list data.
type CustomGenericListVm struct {
	ListData    interface{} `json:"listData"`
	TotalRecord int         `json:"totalRecord"`
	IsSuccess   bool        `json:"isSuccess"`
	Message     string      `json:"message"`
	Errors      []string    `json:"errors"`
}

// CustomGenericVm is a generic single-item response with untyped data.
type CustomGenericVm struct {
	Data        interface{} `json:"data"`
	TotalRecord int         `json:"totalRecord"`
	IsSuccess   bool        `json:"isSuccess"`
	Message     string      `json:"message"`
	Errors      []string    `json:"errors"`
}

// CustomEnumTypeResponseVm is a select-list option DTO used for frontend
// dropdowns.
type CustomEnumTypeResponseVm struct {
	Text     string `json:"text"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
	Selected bool   `json:"selected"`
}

// CustomReviewPeriodRangeVm is a select-list option DTO for review period
// range dropdowns.
type CustomReviewPeriodRangeVm struct {
	Disabled bool   `json:"disabled"`
	Selected bool   `json:"selected"`
	Text     string `json:"text"`
	Value    string `json:"value"`
}

// ChartDataset carries data for a single chart series.
type ChartDataset struct {
	Label           string `json:"label"`
	BackgroundColor []string `json:"backgroundColor"`
	Data            []int    `json:"data"`
}

// ---------------------------------------------------------------------------
// Custom Review Period VMs
// ---------------------------------------------------------------------------

// CustomReviewPeriodResponseVm is a lightweight review period DTO for
// competency frontend forms.
type CustomReviewPeriodResponseVm struct {
	ReviewPeriodID int        `json:"reviewPeriodId"`
	BankYearID     int        `json:"bankYearId"`
	Name           string     `json:"name"`
	StartDate      time.Time  `json:"startDate"`
	EndDate        time.Time  `json:"endDate"`
	BankYearName   string     `json:"bankYearName"`
	ApprovedBy     string     `json:"approvedBy"`
	DateApproved   *time.Time `json:"dateApproved"`
	IsApproved     bool       `json:"isApproved"`
	IsActive       bool       `json:"isActive"`
}

// CustomBankYearResponseVm is a lightweight bank year DTO for dropdowns.
type CustomBankYearResponseVm struct {
	BankYearID int    `json:"bankYearId"`
	YearName   string `json:"yearName"`
	IsActive   bool   `json:"isActive"`
}

// ---------------------------------------------------------------------------
// Approve / Reject VMs
// ---------------------------------------------------------------------------

// ApproveRejectRequestVm is the batch approve/reject payload used by the
// frontend for bulk operations.
type ApproveRejectRequestVm struct {
	EntityType      string   `json:"entityType"`
	RecordIDs       []string `json:"recordIds"`
	RejectionReason string   `json:"rejectionReason"`
}

// ApproveRejectRequestSingleVm is the single-record approve/reject payload.
type ApproveRejectRequestSingleVm struct {
	EntityType      string  `json:"entityType"`
	RecordID        string  `json:"recordId"`
	RejectionReason *string `json:"rejectionReason"`
}

// ---------------------------------------------------------------------------
// Custom Operational Objective VMs
// ---------------------------------------------------------------------------

// CustomOperationalObjectiveVm is a lightweight operational objective DTO
// for the frontend objective planning view.
type CustomOperationalObjectiveVm struct {
	PlannedObjectiveID string `json:"plannedObjectiveId"`
	ReviewPeriodID     string `json:"reviewPeriodId"`
	ReviewPeriod       string `json:"reviewPeriod"`
	ObjectiveLevel     string `json:"objectiveLevel"`
	ObjectiveID        string `json:"objectiveId"`
	Objective          string `json:"objective"`
	CreatedBy          string `json:"createdBy"`
	Comment            string `json:"comment"`
	RecordStatus       int    `json:"recordStatus"`
	StaffID            string `json:"staffId"`
	IsActive           bool   `json:"isActive"`
}

// ObjectivePlanningListResponseVm wraps planned objectives for the
// objective planning screen.
type ObjectivePlanningListResponseVm struct {
	GenericListResponseVm
	PlannedObjectives []CustomOperationalObjectiveVm `json:"plannedObjectives"`
}

// OperationalObjectivesListResponseVm wraps operational objectives with
// detail data.
type OperationalObjectivesListResponseVm struct {
	GenericListResponseVm
	Objectives []OperationalObjectiveDataForm `json:"objectives"`
}

// OperationalObjectiveDataForm is the detailed operational objective DTO
// used in frontend forms.
type OperationalObjectiveDataForm struct {
	ReviewPeriodID      string                  `json:"reviewPeriodId"`
	ReviewPeriod        string                  `json:"reviewPeriod"`
	ObjectiveLevel      string                  `json:"objectiveLevel"`
	ObjectiveID         string                  `json:"objectiveId"`
	Objective           string                  `json:"objective"`
	RecordStatus        int                     `json:"recordStatus"`
	StaffID             string                  `json:"staffId"`
	Kpi                 string                  `json:"kpi"`
	Target              string                  `json:"target"`
	Description         string                  `json:"description"`
	EnterpriseObjective *EnterpriseObjectiveDataVm `json:"enterpriseObjective"`
}

// ---------------------------------------------------------------------------
// Custom Work Product VMs (frontend composites)
// ---------------------------------------------------------------------------

// CustomOperationalObjectiveWorkProduct links an operational objective to
// a work product.
type CustomOperationalObjectiveWorkProduct struct {
	OperationalObjectiveWorkProductID string           `json:"operationalObjectiveWorkProductId"`
	ObjectiveID                       string           `json:"objectiveId"`
	WorkProductEvaluationID           string           `json:"workProductEvaluationId"`
	WorkProduct                       *CustomWorkProductVm `json:"workProduct"`
}

// GetOperationalWorkProductListVm wraps operational objective work products.
type GetOperationalWorkProductListVm struct {
	OperationalObjectiveWorkProducts []CustomOperationalObjectiveWorkProduct `json:"operationalObjectiveWorkProducts"`
}

// GetEvaluationOptionListVm wraps evaluation options for a dropdown.
type GetEvaluationOptionListVm struct {
	EvaluationOptions []EvaluationOptionVm `json:"evaluationOptions"`
}

// GetWorkProductDetailVm wraps a single custom work product.
type GetWorkProductDetailVm struct {
	WorkProduct *CustomWorkProductVm `json:"workProduct"`
}

// GetWorkProductDetailForProjectVm wraps a single project work product.
type GetWorkProductDetailForProjectVm struct {
	WorkProduct *CustomProjectWorkProductVm `json:"workProduct"`
}

// GetWorkProductEvaluationVm wraps a single work product evaluation.
type GetWorkProductEvaluationVm struct {
	WorkProductEvaluation *WorkProductEvaluationVm `json:"workProductEvaluation"`
}

// ---------------------------------------------------------------------------
// Custom Statistics VMs
// ---------------------------------------------------------------------------

// CustomWorkProductStatisticsVm carries work product statistics for the
// dashboard.
type CustomWorkProductStatisticsVm struct {
	StaffID                          string `json:"staffId"`
	ReviewPeriodID                   string `json:"reviewPeriodId"`
	NoAllWorkProducts                int    `json:"noAllWorkProducts"`
	NoActiveWorkProducts             int    `json:"noActiveWorkProducts"`
	NoWorkProductsAwaitingEvaluation int    `json:"noWorkProductsAwaitingEvaluation"`
	NoWorkProductsClosed             int    `json:"noWorkProductsClosed"`
	NoWorkProductsPendingApproval    int    `json:"noWorkProductsPendingApproval"`
	TotalWorkProductTasks            int    `json:"totalWorkProductTasks"`
}

// CustomPerformanceStatisticsVm carries performance point statistics for the
// dashboard.
type CustomPerformanceStatisticsVm struct {
	StaffID           string  `json:"staffId"`
	ReviewPeriodID    string  `json:"reviewPeriodId"`
	MaxPoints         float64 `json:"maxPoints"`
	AccumulatedPoints float64 `json:"accumulatedPoints"`
	DeductedPoints    float64 `json:"deductedPoints"`
	ActualPoints      float64 `json:"actualPoints"`
}

// CustomRequestStatisticsVm carries feedback request statistics for the
// dashboard.
type CustomRequestStatisticsVm struct {
	StaffID                     string  `json:"staffId"`
	ReviewPeriodID              string  `json:"reviewPeriodId"`
	CompletedRequests           int     `json:"completedRequests"`
	PendingRequests             int     `json:"pendingRequests"`
	BreachedRequests            int     `json:"breachedRequests"`
	Pending360FeedbacksToTreat  int     `json:"pending360FeedbacksToTreat"`
	CompletedOverdueRequests    int     `json:"completedOverdueRequests"`
	PendingOverdueRequests      int     `json:"pendingOverdueRequests"`
	DeductedPoints              float64 `json:"deductedPoints"`
}

// ---------------------------------------------------------------------------
// Objective Work Product VMs (frontend composite)
// ---------------------------------------------------------------------------

// ObjectiveWorkProductVm links an objective to a work product definition for
// the frontend view.
type ObjectiveWorkProductVm struct {
	ReviewPeriodID          string                     `json:"reviewPeriodId"`
	ReviewPeriod            string                     `json:"reviewPeriod"`
	ObjectiveLevel          string                     `json:"objectiveLevel"`
	ObjectiveID             string                     `json:"objectiveId"`
	Objective               string                     `json:"objective"`
	WorkProductDefinitionID string                     `json:"workProductDefinitionId"`
	WorkProductName         string                     `json:"workProductName"`
	Deliverables            string                     `json:"deliverables"`
	StaffID                 string                     `json:"staffId"`
	Kpi                     string                     `json:"kpi"`
	Target                  string                     `json:"target"`
	EnterpriseObjective     *EnterpriseObjectiveDataVm `json:"enterpriseObjective"`
}

// GetObjectiveWorkProductsVm wraps objective-linked work products.
type GetObjectiveWorkProductsVm struct {
	ObjectiveWorkProducts []ObjectiveWorkProductVm `json:"objectiveWorkProducts"`
}

// ---------------------------------------------------------------------------
// Custom Project / Committee Member VMs
// ---------------------------------------------------------------------------

// CustomProjectMemberVm is the frontend DTO for a project member with
// display names.
type CustomProjectMemberVm struct {
	ProjectMemberID    string `json:"projectMemberId"`
	StaffID            string `json:"staffId"`
	ProjectID          string `json:"projectId"`
	StaffName          string `json:"staffName"`
	ObjectiveName      string `json:"objectiveName"`
	PlannedObjectiveID string `json:"plannedObjectiveId"`
}

// CustomCommitteeMemberVm is the frontend DTO for a committee member with
// display names.
type CustomCommitteeMemberVm struct {
	CommitteeMemberID  string `json:"committeeMemberId"`
	StaffID            string `json:"staffId"`
	CommitteeID        string `json:"committeeId"`
	StaffName          string `json:"staffName"`
	ObjectiveName      string `json:"objectiveName"`
	PlannedObjectiveID string `json:"plannedObjectiveId"`
}

// GetProjectMembersVm wraps project members for the frontend.
type GetProjectMembersVm struct {
	ProjectMembers []CustomProjectMemberVm `json:"projectMembers"`
}

// GetCommitteeMembersVm wraps committee members for the frontend.
type GetCommitteeMembersVm struct {
	CommitteeMembers []CustomCommitteeMemberVm `json:"committeeMembers"`
}

// ---------------------------------------------------------------------------
// Project / Committee Work Product frontend wrappers
// ---------------------------------------------------------------------------

// GetProjectWorkProductDefinitionsVm wraps project work product definitions.
type GetProjectWorkProductDefinitionsVm struct {
	WorkProductDefinitions []ProjectWorkProductVm `json:"workProductDefinitions"`
}

// GetProjectWorkProductsVm wraps project work products.
type GetProjectWorkProductsVm struct {
	ProjectWorkProducts []ProjectWorkProductVm `json:"projectWorkProducts"`
}

// GetProjectWorkProductsNewVm wraps project and committee work products
// in a single response.
type GetProjectWorkProductsNewVm struct {
	ProjectWorkProducts   []ProjectWorkProductVm `json:"projectWorkProducts"`
	CommitteeWorkProducts []ProjectWorkProductVm `json:"committeeWorkProducts"`
}

// GetCommitteeWorkProductsVm wraps committee work products.
type GetCommitteeWorkProductsVm struct {
	CommitteeWorkProducts []ProjectWorkProductVm `json:"committeeWorkProducts"`
}

// ---------------------------------------------------------------------------
// Project / Committee Objective Create VMs
// ---------------------------------------------------------------------------

// CreateProjectObjectiveVm is the request DTO for linking an objective to a
// project.
type CreateProjectObjectiveVm struct {
	ObjectiveID string `json:"objectiveId"`
	ProjectID   string `json:"projectId"`
}

// CreateCommitteeObjectiveVm is the request DTO for linking an objective to a
// committee.
type CreateCommitteeObjectiveVm struct {
	ObjectiveID string `json:"objectiveId"`
	CommitteeID string `json:"committeeId"`
}

// ---------------------------------------------------------------------------
// Department Evaluation VMs
// ---------------------------------------------------------------------------

// DepartmentEvaluationResponseVm is the DTO for a department objective
// evaluation result.
type DepartmentEvaluationResponseVm struct {
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

// GetDepartmentEvaluationsVm wraps department objective evaluations.
type GetDepartmentEvaluationsVm struct {
	DepartmentObjectiveEvaluations []DepartmentEvaluationResponseVm `json:"departmentObjectiveEvaluations"`
}

// GetObjectiveDepartmentsVm wraps departments for an enterprise objective.
type GetObjectiveDepartmentsVm struct {
	EnterpriseObjectiveID string        `json:"enterpriseObjectiveId"`
	Departments           []interface{} `json:"departments"`
}

// DepartmentOutcomeRequestVm is the request DTO for submitting a department
// outcome score.
type DepartmentOutcomeRequestVm struct {
	DepartmentID          int     `json:"departmentId"`
	OverallOutcomeScored  float64 `json:"overallOutcomeScored"`
	AllocatedOutcome      float64 `json:"allocatedOutcome"`
	OutcomeScore          float64 `json:"outcomeScore"`
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId"`
	ReviewPeriodID        string  `json:"reviewPeriodId"`
}

// EnterpriseOutcomeRequestVm is the request DTO for submitting an enterprise
// outcome score.
type EnterpriseOutcomeRequestVm struct {
	TotalOutcomeScore     float64 `json:"totalOutcomeScore"`
	OutcomeScore          float64 `json:"outcomeScore"`
	EnterpriseObjectiveID string  `json:"enterpriseObjectiveId"`
	ReviewPeriodID        string  `json:"reviewPeriodId"`
}

// ===========================================================================
// Frontend Custom Request VMs (FrontendCustomVms/RequestFormVms.cs)
// ===========================================================================

// ---------------------------------------------------------------------------
// Custom Office Objective VM
// ---------------------------------------------------------------------------

// CustomOfficeObjectiveVm is the frontend request DTO for an office-level
// objective.
type CustomOfficeObjectiveVm struct {
	ObjectiveID        string  `json:"objectiveId"`
	StaffID            string  `json:"staffId"`
	ReviewPeriodID     string  `json:"reviewPeriodId"`
	ObjectiveLevel     *string `json:"objectiveLevel"`
	OfficeObjectiveID  string  `json:"officeObjectiveId"`
	Status             int     `json:"status"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Kpi                string  `json:"kpi"`
	OfficeID           int     `json:"officeId"`
	DivisionObjectiveID string `json:"divisionObjectiveId"`
	JobGradeGroupID    int     `json:"jobGradeGroupId"`
}

// ---------------------------------------------------------------------------
// Custom Review Period Request VMs
// ---------------------------------------------------------------------------

// CustomPerformanceReviewPeriodRequestVm is the frontend request DTO for
// creating/updating a performance review period.
type CustomPerformanceReviewPeriodRequestVm struct {
	PeriodID          string                   `json:"periodId"`
	Year              int                      `json:"year"`
	Range             enums.ReviewPeriodRange  `json:"range"`
	RangeValue        int                      `json:"rangeValue"`
	Name              string                   `json:"name"`
	StartDate         time.Time                `json:"startDate"`
	EndDate           time.Time                `json:"endDate"`
	MaxPoints         float64                  `json:"maxPoints"`
	MinNoOfObjectives int                      `json:"minNoOfObjectives"`
	MaxNoOfObjectives int                      `json:"maxNoOfObjectives"`
	RecordStatus      enums.Status             `json:"recordStatus"`
	StrategyID        string                   `json:"strategyId"`
}

// CustomReviewPeriodVm is a lightweight review period request DTO without
// dates.
type CustomReviewPeriodVm struct {
	PeriodID          string  `json:"periodId"`
	Year              int     `json:"year"`
	Range             int     `json:"range"`
	RangeValue        int     `json:"rangeValue"`
	Name              string  `json:"name"`
	MaxPoints         float64 `json:"maxPoints"`
	MinNoOfObjectives int     `json:"minNoOfObjectives"`
	MaxNoOfObjectives int     `json:"maxNoOfObjectives"`
	StrategyID        string  `json:"strategyId"`
}

// GetReviewPeriodVm is the full response DTO for a review period with all
// audit fields.
type GetReviewPeriodVm struct {
	PeriodID          string     `json:"periodId"`
	Year              int        `json:"year"`
	Range             int        `json:"range"`
	RangeValue        int        `json:"rangeValue"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	StartDate         time.Time  `json:"startDate"`
	EndDate           time.Time  `json:"endDate"`
	MaxPoints         float64    `json:"maxPoints"`
	MinNoOfObjectives int        `json:"minNoOfObjectives"`
	MaxNoOfObjectives int        `json:"maxNoOfObjectives"`
	RecordStatus      int        `json:"recordStatus"`
	StrategyID        string     `json:"strategyId"`
	DateApproved      *time.Time `json:"dateApproved"`
	IsApproved        bool       `json:"isApproved"`
	IsRejected        bool       `json:"isRejected"`
	DateRejected      *time.Time `json:"dateRejected"`
	ID                int        `json:"id"`
	CreatedAt         *time.Time `json:"createdAt"`
	SoftDeleted       bool       `json:"softDeleted"`
	CreatedBy         string     `json:"createdBy"`
	IsActive          bool       `json:"isActive"`
}

// GetPerformanceReviewPeriodVm wraps a single PerformanceReviewPeriodVm.
type GetPerformanceReviewPeriodVm struct {
	PerformanceReviewPeriod *PerformanceReviewPeriodVm `json:"performanceReviewPeriod"`
}

// ---------------------------------------------------------------------------
// Staff Objective Planning Request VM
// ---------------------------------------------------------------------------

// StaffObjectivePlanningRequestVm is the request DTO for staff objective
// planning operations (approve/reject).
type StaffObjectivePlanningRequestVm struct {
	PlannedObjectiveID string               `json:"plannedObjectiveId"`
	ObjectiveID        string               `json:"objectiveId"`
	StaffID            string               `json:"staffId"`
	ObjectiveLevel     enums.ObjectiveLevel `json:"objectiveLevel"`
	ReviewPeriodID     string               `json:"reviewPeriodId"`
	RejectionReason    string               `json:"rejectionReason"`
}

// ---------------------------------------------------------------------------
// Custom Project / Committee Create VMs (frontend forms)
// ---------------------------------------------------------------------------

// ProjectObjectiveDataVm carries an objective reference within a project.
type ProjectObjectiveDataVm struct {
	ProjectObjectiveID string `json:"projectObjectiveId"`
	ObjectiveID        string `json:"objectiveId"`
	ProjectID          string `json:"projectId"`
}

// CommitteeObjectiveDataVm carries an objective reference within a committee.
type CommitteeObjectiveDataVm struct {
	CommitteeObjectiveID string `json:"committeeObjectiveId"`
	ObjectiveID          string `json:"objectiveId"`
	CommitteeID          string `json:"committeeId"`
}

// CommitteeMemberRefVm carries a committee member reference.
type CommitteeMemberRefVm struct {
	CommitteeMemberID  string `json:"committeeMemberId"`
	StaffID            string `json:"staffId"`
	CommitteeID        string `json:"committeeId"`
	PlannedObjectiveID string `json:"plannedObjectiveId"`
}

// CreateProjectVm is the frontend form DTO for creating/editing a project.
type CreateProjectVm struct {
	ProjectManager     string                    `json:"projectManager"`
	ProjectManagerName string                    `json:"projectManagerName"`
	Name               string                    `json:"name"`
	Description        string                    `json:"description"`
	DepartmentID       int                       `json:"departmentId"`
	StartDate          time.Time                 `json:"startDate"`
	EndDate            time.Time                 `json:"endDate"`
	ObjectiveID        string                    `json:"objectiveId"`
	Deliverables       string                    `json:"deliverables"`
	ReviewPeriodID     string                    `json:"reviewPeriodId"`
	RejectionReason    string                    `json:"rejectionReason"`
	ProjectID          string                    `json:"projectId"`
	RecordStatus       enums.Status              `json:"recordStatus"`
	ProjectObjectives  []ProjectObjectiveDataVm  `json:"projectObjectives"`
	CreatedBy          string                    `json:"createdBy"`
}

// CreateCommitteeVm is the frontend form DTO for creating/editing a
// committee.
type CreateCommitteeVm struct {
	CommitteeID        string                     `json:"committeeId"`
	Chairperson        string                     `json:"chairperson"`
	ChairpersonName    string                     `json:"chairpersonName"`
	CommitteeMembers   []CommitteeMemberRefVm     `json:"committeeMembers"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	StartDate          time.Time                  `json:"startDate"`
	ObjectiveID        string                     `json:"objectiveId"`
	DepartmentID       int                        `json:"departmentId"`
	RejectionReason    string                     `json:"rejectionReason"`
	EndDate            time.Time                  `json:"endDate"`
	Deliverables       string                     `json:"deliverables"`
	RecordStatus       enums.Status               `json:"recordStatus"`
	ReviewPeriodID     string                     `json:"reviewPeriodId"`
	CreatedBy          string                     `json:"createdBy"`
	Remark             string                     `json:"remark"`
	CommitteeObjectives []CommitteeObjectiveDataVm `json:"committeeObjectives"`
}

// GetCommitteesVm wraps a list of committees with a success flag.
type GetCommitteesVm struct {
	Committees []CreateCommitteeVm `json:"committees"`
	IsSuccess  bool                `json:"isSuccess"`
}

// GetProjectsVm wraps a list of projects.
type GetProjectsVm struct {
	Projects []CreateProjectVm `json:"projects"`
}
