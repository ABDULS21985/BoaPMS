package dto

import "time"

// ===========================================================================
// Project Request VMs
// ===========================================================================

// CreateProjectRequestModel is the request body for creating a new project.
type CreateProjectRequestModel struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	ProjectManager string    `json:"project_manager"`
}

// ProjectRequestModel extends the create model with identifier and comment.
type ProjectRequestModel struct {
	ProjectID      string    `json:"project_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	ProjectManager string    `json:"project_manager"`
	Comment        string    `json:"comment,omitempty"`
}

// ProjectObjectiveRequestModel links an objective to a project.
type ProjectObjectiveRequestModel struct {
	ProjectID   string `json:"project_id"`
	ObjectiveID string `json:"objective_id"`
}

// CreateProjectMemberRequestModel is the request to add a member to a project.
type CreateProjectMemberRequestModel struct {
	ProjectID string `json:"project_id"`
	StaffID   string `json:"staff_id"`
}

// ProjectMemberRequestModel extends the create model with member details.
type ProjectMemberRequestModel struct {
	ProjectMemberID    string `json:"project_member_id"`
	ProjectID          string `json:"project_id"`
	StaffID            string `json:"staff_id"`
	PlannedObjectiveID string `json:"planned_objective_id"`
	Comment            string `json:"comment,omitempty"`
}

// ===========================================================================
// Committee Request VMs
// ===========================================================================

// CreateCommitteeRequestModel is the request body for creating a new committee.
type CreateCommitteeRequestModel struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	ChairPerson    string    `json:"chair_person"`
}

// CommitteeRequestModel extends the create model with identifier and comment.
type CommitteeRequestModel struct {
	CommitteeID    string    `json:"committee_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	ChairPerson    string    `json:"chair_person"`
	Comment        string    `json:"comment,omitempty"`
}

// CommitteeObjectiveRequestModel links an objective to a committee.
type CommitteeObjectiveRequestModel struct {
	CommitteeID string `json:"committee_id"`
	ObjectiveID string `json:"objective_id"`
}

// CreateCommitteeMemberRequestModel is the request to add a member to a committee.
type CreateCommitteeMemberRequestModel struct {
	CommitteeID string `json:"committee_id"`
	StaffID     string `json:"staff_id"`
}

// CommitteeMemberRequestModel extends the create model with member details.
type CommitteeMemberRequestModel struct {
	CommitteeMemberID  string `json:"committee_member_id"`
	CommitteeID        string `json:"committee_id"`
	StaffID            string `json:"staff_id"`
	PlannedObjectiveID string `json:"planned_objective_id"`
	Comment            string `json:"comment,omitempty"`
}

// ===========================================================================
// Work Product Request VMs
// ===========================================================================

// CreateWorkProductRequestModel is the request body for creating a work product.
type CreateWorkProductRequestModel struct {
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	MaxPoint                float64    `json:"max_point"`
	WorkProductType         string     `json:"work_product_type"`
	StaffID                 string     `json:"staff_id"`
	StartDate               time.Time  `json:"start_date"`
	EndDate                 time.Time  `json:"end_date"`
	Deliverables            string     `json:"deliverables"`
	PlannedObjectiveID      string     `json:"planned_objective_id"`
	WorkProductDefinitionID string     `json:"work_product_definition_id"`
	ReviewPeriodID          string     `json:"review_period_id"`
}

// WorkProductRequestModel extends the create model with identifier and approval fields.
type WorkProductRequestModel struct {
	WorkProductID           string     `json:"work_product_id"`
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	MaxPoint                float64    `json:"max_point"`
	WorkProductType         string     `json:"work_product_type"`
	StaffID                 string     `json:"staff_id"`
	StartDate               time.Time  `json:"start_date"`
	EndDate                 time.Time  `json:"end_date"`
	Deliverables            string     `json:"deliverables"`
	PlannedObjectiveID      string     `json:"planned_objective_id"`
	WorkProductDefinitionID string     `json:"work_product_definition_id"`
	ReviewPeriodID          string     `json:"review_period_id"`
	AcceptanceComment       string     `json:"acceptance_comment,omitempty"`
	ApproverComment         string     `json:"approver_comment,omitempty"`
	Comment                 string     `json:"comment,omitempty"`
}

// WorkProductTaskRequestModel represents a task under a work product.
type WorkProductTaskRequestModel struct {
	WorkProductTaskID string    `json:"work_product_task_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	WorkProductID     string    `json:"work_product_id"`
}

// ProjectAssignedWorkProductRequestModel is the request for assigning a work product to a project.
type ProjectAssignedWorkProductRequestModel struct {
	WorkProductDefinitionID string    `json:"work_product_definition_id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	ProjectID               string    `json:"project_id"`
	ReviewPeriodID          string    `json:"review_period_id"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Deliverables            string    `json:"deliverables"`
}

// CommitteeAssignedWorkProductRequestModel is the request for assigning a work product to a committee.
type CommitteeAssignedWorkProductRequestModel struct {
	WorkProductDefinitionID string    `json:"work_product_definition_id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	CommitteeID             string    `json:"committee_id"`
	ReviewPeriodID          string    `json:"review_period_id"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Deliverables            string    `json:"deliverables"`
}

// ===========================================================================
// Work Product Evaluation VMs
// ===========================================================================

// WorkProductEvaluationRequestModel is the request for evaluating a work product.
type WorkProductEvaluationRequestModel struct {
	WorkProductID                string  `json:"work_product_id"`
	Timeliness                   float64 `json:"timeliness"`
	TimelinessEvaluationOptionID string  `json:"timeliness_evaluation_option_id"`
	Quality                      float64 `json:"quality"`
	QualityEvaluationOptionID    string  `json:"quality_evaluation_option_id"`
	Output                       float64 `json:"output"`
	OutputEvaluationOptionID     string  `json:"output_evaluation_option_id"`
	EvaluatorStaffID             string  `json:"evaluator_staff_id"`
}

// PeriodObjectiveEvaluationRequestModel is the request for evaluating a period objective.
type PeriodObjectiveEvaluationRequestModel struct {
	EnterpriseObjectiveID string  `json:"enterprise_objective_id"`
	ReviewPeriodID        string  `json:"review_period_id"`
	TotalOutcomeScore     float64 `json:"total_outcome_score"`
	OutcomeScore          float64 `json:"outcome_score"`
}

// PeriodObjectiveDepartmentEvaluationRequestModel extends PeriodObjectiveEvaluationRequestModel with department context.
type PeriodObjectiveDepartmentEvaluationRequestModel struct {
	EnterpriseObjectiveID string  `json:"enterprise_objective_id"`
	ReviewPeriodID        string  `json:"review_period_id"`
	TotalOutcomeScore     float64 `json:"total_outcome_score"`
	OutcomeScore          float64 `json:"outcome_score"`
	DepartmentID          string  `json:"department_id"`
	AllocatedOutcome      float64 `json:"allocated_outcome"`
	OverallOutcomeScored  float64 `json:"overall_outcome_scored"`
}

// ===========================================================================
// Project Response VMs & Data Structs
// ===========================================================================

// BaseProjectData holds the core project fields shared across responses.
type BaseProjectData struct {
	ProjectID      string    `json:"project_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	ProjectManager string    `json:"project_manager"`
	RecordStatus   string    `json:"record_status"`
}

// ProjectData extends BaseProjectData with computed or joined fields.
type ProjectData struct {
	BaseProjectData
	DepartmentName     string `json:"department_name"`
	ProjectManagerName string `json:"project_manager_name"`
	MemberCount        int    `json:"member_count"`
	ObjectiveCount     int    `json:"objective_count"`
}

// ProjectResponseVm wraps a single project in a standard response.
type ProjectResponseVm struct {
	BaseAPIResponse
	Project ProjectData `json:"project"`
}

// ProjectListResponseVm wraps a list of projects.
type ProjectListResponseVm struct {
	GenericListResponseVm
	Projects []ProjectData `json:"projects"`
}

// ---------------------------------------------------------------------------
// Committee Response VMs & Data Structs
// ---------------------------------------------------------------------------

// CommitteeData holds the committee fields for responses.
type CommitteeData struct {
	CommitteeID    string    `json:"committee_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	ChairPerson    string    `json:"chair_person"`
	ChairPersonName string   `json:"chair_person_name"`
	RecordStatus   string    `json:"record_status"`
	MemberCount    int       `json:"member_count"`
	ObjectiveCount int       `json:"objective_count"`
}

// CommitteeResponseVm wraps a single committee in a standard response.
type CommitteeResponseVm struct {
	BaseAPIResponse
	Committee CommitteeData `json:"committee"`
}

// CommitteeListResponseVm wraps a list of committees.
type CommitteeListResponseVm struct {
	GenericListResponseVm
	Committees []CommitteeData `json:"committees"`
}

// ---------------------------------------------------------------------------
// Project Member Response VMs & Data Structs
// ---------------------------------------------------------------------------

// ProjectMemberData holds project member data for responses.
type ProjectMemberData struct {
	ProjectMemberID    string `json:"project_member_id"`
	ProjectID          string `json:"project_id"`
	StaffID            string `json:"staff_id"`
	StaffName          string `json:"staff_name"`
	PlannedObjectiveID string `json:"planned_objective_id"`
	RecordStatus       string `json:"record_status"`
}

// ProjectMemberListResponseVm wraps a list of project members.
type ProjectMemberListResponseVm struct {
	GenericListResponseVm
	Members []ProjectMemberData `json:"members"`
}

// ---------------------------------------------------------------------------
// Committee Member Response VMs & Data Structs
// ---------------------------------------------------------------------------

// CommitteeMemberData holds committee member data for responses.
type CommitteeMemberData struct {
	CommitteeMemberID  string `json:"committee_member_id"`
	CommitteeID        string `json:"committee_id"`
	StaffID            string `json:"staff_id"`
	StaffName          string `json:"staff_name"`
	PlannedObjectiveID string `json:"planned_objective_id"`
	RecordStatus       string `json:"record_status"`
}

// CommitteeMemberListResponseVm wraps a list of committee members.
type CommitteeMemberListResponseVm struct {
	GenericListResponseVm
	Members []CommitteeMemberData `json:"members"`
}

// ---------------------------------------------------------------------------
// Project Objective Response VMs & Data Structs
// ---------------------------------------------------------------------------

// ProjectObjectiveData holds project objective data for responses.
type ProjectObjectiveData struct {
	ProjectObjectiveID string `json:"project_objective_id"`
	ProjectID          string `json:"project_id"`
	ObjectiveID        string `json:"objective_id"`
	ObjectiveName      string `json:"objective_name"`
	ObjectiveLevel     string `json:"objective_level"`
	RecordStatus       string `json:"record_status"`
}

// ProjectObjectiveListResponseVm wraps a list of project objectives.
type ProjectObjectiveListResponseVm struct {
	GenericListResponseVm
	Objectives []ProjectObjectiveData `json:"objectives"`
}

// ---------------------------------------------------------------------------
// Committee Objective Response VMs & Data Structs
// ---------------------------------------------------------------------------

// CommitteeObjectiveData holds committee objective data for responses.
type CommitteeObjectiveData struct {
	CommitteeObjectiveID string `json:"committee_objective_id"`
	CommitteeID          string `json:"committee_id"`
	ObjectiveID          string `json:"objective_id"`
	ObjectiveName        string `json:"objective_name"`
	ObjectiveLevel       string `json:"objective_level"`
	RecordStatus         string `json:"record_status"`
}

// CommitteeObjectiveListResponseVm wraps a list of committee objectives.
type CommitteeObjectiveListResponseVm struct {
	GenericListResponseVm
	Objectives []CommitteeObjectiveData `json:"objectives"`
}

// ---------------------------------------------------------------------------
// Work Product Response VMs & Data Structs
// ---------------------------------------------------------------------------

// WorkProductData holds work product data for responses.
type WorkProductData struct {
	WorkProductID           string     `json:"work_product_id"`
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	MaxPoint                float64    `json:"max_point"`
	WorkProductType         string     `json:"work_product_type"`
	StaffID                 string     `json:"staff_id"`
	StaffName               string     `json:"staff_name"`
	StartDate               time.Time  `json:"start_date"`
	EndDate                 time.Time  `json:"end_date"`
	Deliverables            string     `json:"deliverables"`
	PlannedObjectiveID      string     `json:"planned_objective_id"`
	WorkProductDefinitionID string     `json:"work_product_definition_id"`
	ReviewPeriodID          string     `json:"review_period_id"`
	RecordStatus            string     `json:"record_status"`
	AcceptanceComment       string     `json:"acceptance_comment,omitempty"`
	ApproverComment         string     `json:"approver_comment,omitempty"`
	AccumulatedPoints       float64    `json:"accumulated_points"`
	EvaluatedDate           *time.Time `json:"evaluated_date,omitempty"`
}

// WorkProductResponseVm wraps a single work product in a standard response.
type WorkProductResponseVm struct {
	BaseAPIResponse
	WorkProduct WorkProductData `json:"work_product"`
}

// ---------------------------------------------------------------------------
// Project Work Product Response VMs & Data Structs
// ---------------------------------------------------------------------------

// ProjectWorkProductData holds project work product data for responses.
type ProjectWorkProductData struct {
	WorkProductID string    `json:"work_product_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ProjectID     string    `json:"project_id"`
	ProjectName   string    `json:"project_name"`
	StaffID       string    `json:"staff_id"`
	StaffName     string    `json:"staff_name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	RecordStatus  string    `json:"record_status"`
	MaxPoint      float64   `json:"max_point"`
}

// ProjectWorkProductListResponseVm wraps a list of project work products.
type ProjectWorkProductListResponseVm struct {
	GenericListResponseVm
	WorkProducts []ProjectWorkProductData `json:"work_products"`
}

// ---------------------------------------------------------------------------
// Committee Work Product Response VMs & Data Structs
// ---------------------------------------------------------------------------

// CommitteeWorkProductData holds committee work product data for responses.
type CommitteeWorkProductData struct {
	WorkProductID string    `json:"work_product_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CommitteeID   string    `json:"committee_id"`
	CommitteeName string    `json:"committee_name"`
	StaffID       string    `json:"staff_id"`
	StaffName     string    `json:"staff_name"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	RecordStatus  string    `json:"record_status"`
	MaxPoint      float64   `json:"max_point"`
}

// CommitteeWorkProductListResponseVm wraps a list of committee work products.
type CommitteeWorkProductListResponseVm struct {
	GenericListResponseVm
	WorkProducts []CommitteeWorkProductData `json:"work_products"`
}

// ---------------------------------------------------------------------------
// Operational Objective Work Product Response VMs & Data Structs
// ---------------------------------------------------------------------------

// OperationalObjectiveWorkProductData holds work product data scoped to an operational objective.
type OperationalObjectiveWorkProductData struct {
	WorkProductID      string    `json:"work_product_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	ObjectiveID        string    `json:"objective_id"`
	ObjectiveName      string    `json:"objective_name"`
	PlannedObjectiveID string    `json:"planned_objective_id"`
	StaffID            string    `json:"staff_id"`
	StaffName          string    `json:"staff_name"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	RecordStatus       string    `json:"record_status"`
	MaxPoint           float64   `json:"max_point"`
}

// OperationalObjectiveWorkProductListResponseVm wraps a list of objective-scoped work products.
type OperationalObjectiveWorkProductListResponseVm struct {
	GenericListResponseVm
	WorkProducts []OperationalObjectiveWorkProductData `json:"work_products"`
}

// ---------------------------------------------------------------------------
// Staff Work Product & Objective Work Product Response VMs
// ---------------------------------------------------------------------------

// StaffWorkProductListResponseVm wraps a list of work products for a staff member.
type StaffWorkProductListResponseVm struct {
	GenericListResponseVm
	WorkProducts []WorkProductData `json:"work_products"`
}

// ObjectiveWorkProductData holds work product data associated with an objective.
type ObjectiveWorkProductData struct {
	WorkProductID      string    `json:"work_product_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	PlannedObjectiveID string    `json:"planned_objective_id"`
	ObjectiveID        string    `json:"objective_id"`
	ObjectiveName      string    `json:"objective_name"`
	StaffID            string    `json:"staff_id"`
	StaffName          string    `json:"staff_name"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	RecordStatus       string    `json:"record_status"`
	MaxPoint           float64   `json:"max_point"`
	AccumulatedPoints  float64   `json:"accumulated_points"`
}

// ObjectiveWorkProductListResponseVm wraps a list of work products for an objective.
type ObjectiveWorkProductListResponseVm struct {
	GenericListResponseVm
	WorkProducts []ObjectiveWorkProductData `json:"work_products"`
}

// ---------------------------------------------------------------------------
// Work Product Task Response VMs & Data Structs
// ---------------------------------------------------------------------------

// WorkProductTaskDetail represents the detailed view of a work product task.
type WorkProductTaskDetail struct {
	WorkProductTaskID string     `json:"work_product_task_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	WorkProductID     string     `json:"work_product_id"`
	RecordStatus      string     `json:"record_status"`
	CompletedDate     *time.Time `json:"completed_date,omitempty"`
}

// WorkProductTaskData holds task data for list responses.
type WorkProductTaskData struct {
	WorkProductTaskID string     `json:"work_product_task_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	WorkProductID     string     `json:"work_product_id"`
	WorkProductName   string     `json:"work_product_name"`
	RecordStatus      string     `json:"record_status"`
	CompletedDate     *time.Time `json:"completed_date,omitempty"`
}

// WorkProductTaskListResponseVm wraps a list of work product tasks.
type WorkProductTaskListResponseVm struct {
	GenericListResponseVm
	Tasks []WorkProductTaskData `json:"tasks"`
}

// WorkProductTaskResponseVm wraps a single work product task.
type WorkProductTaskResponseVm struct {
	BaseAPIResponse
	Task WorkProductTaskDetail `json:"task"`
}

// ---------------------------------------------------------------------------
// Work Product Evaluation Response VMs & Data Structs
// ---------------------------------------------------------------------------

// EvaluationOptionData holds the data for a single evaluation option.
type EvaluationOptionData struct {
	EvaluationOptionID string  `json:"evaluation_option_id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Score              float64 `json:"score"`
	EvaluationType     string  `json:"evaluation_type"`
}

// WorkProductEvaluationData holds a work product evaluation result.
type WorkProductEvaluationData struct {
	WorkProductEvaluationID      string  `json:"work_product_evaluation_id"`
	WorkProductID                string  `json:"work_product_id"`
	Timeliness                   float64 `json:"timeliness"`
	TimelinessEvaluationOptionID string  `json:"timeliness_evaluation_option_id"`
	Quality                      float64 `json:"quality"`
	QualityEvaluationOptionID    string  `json:"quality_evaluation_option_id"`
	Output                       float64 `json:"output"`
	OutputEvaluationOptionID     string  `json:"output_evaluation_option_id"`
	EvaluatorStaffID             string  `json:"evaluator_staff_id"`
	EvaluatorName                string  `json:"evaluator_name"`
	TotalScore                   float64 `json:"total_score"`
	RecordStatus                 string  `json:"record_status"`
}

// EvaluationResponseVm wraps evaluation options in a standard response.
type EvaluationResponseVm struct {
	BaseAPIResponse
	EvaluationOptions []EvaluationOptionData `json:"evaluation_options"`
}

// WorkProductEvaluationResponseVm wraps a work product evaluation result.
type WorkProductEvaluationResponseVm struct {
	BaseAPIResponse
	Evaluation WorkProductEvaluationData `json:"evaluation"`
}

// ---------------------------------------------------------------------------
// Project Assigned Work Product Response VMs & Data Structs
// ---------------------------------------------------------------------------

// ProjectAssignedWorkProductData holds data for a work product assigned to a project.
type ProjectAssignedWorkProductData struct {
	AssignedWorkProductID   string    `json:"assigned_work_product_id"`
	WorkProductDefinitionID string    `json:"work_product_definition_id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	ProjectID               string    `json:"project_id"`
	ProjectName             string    `json:"project_name"`
	ReviewPeriodID          string    `json:"review_period_id"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Deliverables            string    `json:"deliverables"`
	RecordStatus            string    `json:"record_status"`
}

// ProjectAssignedWorkProductResponseVm wraps a single assigned work product.
type ProjectAssignedWorkProductResponseVm struct {
	BaseAPIResponse
	AssignedWorkProduct ProjectAssignedWorkProductData `json:"assigned_work_product"`
}

// ProjectAssignedWorkProductListResponseVm wraps a list of assigned work products.
type ProjectAssignedWorkProductListResponseVm struct {
	GenericListResponseVm
	AssignedWorkProducts []ProjectAssignedWorkProductData `json:"assigned_work_products"`
}

// ---------------------------------------------------------------------------
// Committee Assigned Work Product Response VMs & Data Structs
// ---------------------------------------------------------------------------

// CommitteeAssignedWorkProductData holds data for a work product assigned to a committee.
type CommitteeAssignedWorkProductData struct {
	AssignedWorkProductID   string    `json:"assigned_work_product_id"`
	WorkProductDefinitionID string    `json:"work_product_definition_id"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	CommitteeID             string    `json:"committee_id"`
	CommitteeName           string    `json:"committee_name"`
	ReviewPeriodID          string    `json:"review_period_id"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Deliverables            string    `json:"deliverables"`
	RecordStatus            string    `json:"record_status"`
}

// CommitteeAssignedWorkProductResponseVm wraps a single committee-assigned work product.
type CommitteeAssignedWorkProductResponseVm struct {
	BaseAPIResponse
	AssignedWorkProduct CommitteeAssignedWorkProductData `json:"assigned_work_product"`
}

// CommitteeAssignedWorkProductListResponseVm wraps a list of committee-assigned work products.
type CommitteeAssignedWorkProductListResponseVm struct {
	GenericListResponseVm
	AssignedWorkProducts []CommitteeAssignedWorkProductData `json:"assigned_work_products"`
}

// ---------------------------------------------------------------------------
// Recalculate, Assigned List & Adhoc Response VMs
// ---------------------------------------------------------------------------

// RecalculateWorkProductResponseVm wraps the result of a work product recalculation.
type RecalculateWorkProductResponseVm struct {
	BaseAPIResponse
	RecalculatedCount int `json:"recalculated_count"`
}

// ProjectAssignedListResponseVm wraps a list of project assignments.
type ProjectAssignedListResponseVm struct {
	GenericListResponseVm
	Projects []ProjectData `json:"projects"`
}

// CommitteeAssignedListResponseVm wraps a list of committee assignments.
type CommitteeAssignedListResponseVm struct {
	GenericListResponseVm
	Committees []CommitteeData `json:"committees"`
}

// AdhocStaffResponseVm wraps an adhoc staff assignment response.
type AdhocStaffResponseVm struct {
	BaseAPIResponse
	StaffID   string `json:"staff_id"`
	StaffName string `json:"staff_name"`
	RoleName  string `json:"role_name"`
}

// ===========================================================================
// View Model DTOs (Vm suffix – handler-layer representations)
// ===========================================================================

// ---------------------------------------------------------------------------
// Project View Models
// ---------------------------------------------------------------------------

// ProjectVm represents the full view of a project for API responses.
type ProjectVm struct {
	ProjectID      string    `json:"project_id"`
	ProjectManager string    `json:"project_manager"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	RecordStatus   string    `json:"record_status"`
	Members        []ProjectMemberVm    `json:"members,omitempty"`
	Objectives     []ProjectObjectiveVm `json:"objectives,omitempty"`
}

// CreateProjectVm is the request body for creating or updating a project via the handler layer.
type CreateProjectVm struct {
	ProjectID      string    `json:"project_id,omitempty"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	ProjectManager string    `json:"project_manager"`
}

// ProjectObjectiveVm represents an objective linked to a project.
type ProjectObjectiveVm struct {
	ProjectObjectiveID string `json:"project_objective_id"`
	ObjectiveID        string `json:"objective_id"`
	ObjectiveName      string `json:"objective_name"`
	Kpi                string `json:"kpi,omitempty"`
	ProjectID          string `json:"project_id"`
	RecordStatus       string `json:"record_status"`
}

// ProjectMemberVm represents a member assigned to a project.
type ProjectMemberVm struct {
	ProjectMemberID    string `json:"project_member_id"`
	StaffID            string `json:"staff_id"`
	StaffName          string `json:"staff_name"`
	ProjectID          string `json:"project_id"`
	PlannedObjectiveID string `json:"planned_objective_id"`
	RecordStatus       string `json:"record_status"`
}

// ---------------------------------------------------------------------------
// Committee View Models
// ---------------------------------------------------------------------------

// CommitteeVm represents the full view of a committee for API responses.
type CommitteeVm struct {
	CommitteeID     string    `json:"committee_id"`
	Chairperson     string    `json:"chairperson"`
	ChairpersonName string    `json:"chairperson_name"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	Deliverables    string    `json:"deliverables"`
	ReviewPeriodID  string    `json:"review_period_id"`
	DepartmentID    string    `json:"department_id"`
	DepartmentName  string    `json:"department_name"`
	RecordStatus    string    `json:"record_status"`
	Members         []CommitteeMemberVm    `json:"members,omitempty"`
	Objectives      []CommitteeObjectiveVm `json:"objectives,omitempty"`
}

// CreateCommitteeVm is the request body for creating or updating a committee via the handler layer.
type CreateCommitteeVm struct {
	CommitteeID    string    `json:"committee_id,omitempty"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Deliverables   string    `json:"deliverables"`
	ReviewPeriodID string    `json:"review_period_id"`
	DepartmentID   string    `json:"department_id"`
	Chairperson    string    `json:"chairperson"`
}

// CommitteeObjectiveVm represents an objective linked to a committee.
type CommitteeObjectiveVm struct {
	CommitteeObjectiveID string `json:"committee_objective_id"`
	ObjectiveID          string `json:"objective_id"`
	ObjectiveName        string `json:"objective_name"`
	Kpi                  string `json:"kpi,omitempty"`
	CommitteeID          string `json:"committee_id"`
	RecordStatus         string `json:"record_status"`
}

// CommitteeMemberVm represents a member assigned to a committee.
type CommitteeMemberVm struct {
	CommitteeMemberID  string `json:"committee_member_id"`
	StaffID            string `json:"staff_id"`
	StaffName          string `json:"staff_name"`
	CommitteeID        string `json:"committee_id"`
	PlannedObjectiveID string `json:"planned_objective_id"`
	RecordStatus       string `json:"record_status"`
}

// ---------------------------------------------------------------------------
// Work Product View Models
// ---------------------------------------------------------------------------

// WorkProductVm represents the full view of a work product for API responses.
type WorkProductVm struct {
	WorkProductID           string     `json:"work_product_id"`
	Name                    string     `json:"name"`
	Description             string     `json:"description"`
	MaxPoint                float64    `json:"max_point"`
	WorkProductType         string     `json:"work_product_type"`
	IsSelfCreated           bool       `json:"is_self_created"`
	StaffID                 string     `json:"staff_id"`
	StaffName               string     `json:"staff_name"`
	AcceptanceComment       string     `json:"acceptance_comment,omitempty"`
	StartDate               time.Time  `json:"start_date"`
	EndDate                 time.Time  `json:"end_date"`
	Deliverables            string     `json:"deliverables"`
	FinalScore              float64    `json:"final_score"`
	NoReturned              int        `json:"no_returned"`
	CompletionDate          *time.Time `json:"completion_date,omitempty"`
	ApproverComment         string     `json:"approver_comment,omitempty"`
	RecordStatus            string     `json:"record_status"`
	ReviewPeriodID          string     `json:"review_period_id"`
	PlannedObjectiveID      string     `json:"planned_objective_id"`
	WorkProductDefinitionID string     `json:"work_product_definition_id"`
	Tasks                   []WorkProductTaskVm       `json:"tasks,omitempty"`
	Evaluation              *WorkProductEvaluationVm  `json:"evaluation,omitempty"`
}

// CreateWorkProductVm is the request body for creating a work product via the handler layer.
type CreateWorkProductVm struct {
	WorkProductID           string    `json:"work_product_id,omitempty"`
	Name                    string    `json:"name"`
	Description             string    `json:"description"`
	MaxPoint                float64   `json:"max_point"`
	WorkProductType         string    `json:"work_product_type"`
	StaffID                 string    `json:"staff_id"`
	StartDate               time.Time `json:"start_date"`
	EndDate                 time.Time `json:"end_date"`
	Deliverables            string    `json:"deliverables"`
	PlannedObjectiveID      string    `json:"planned_objective_id"`
	WorkProductDefinitionID string    `json:"work_product_definition_id"`
	ReviewPeriodID          string    `json:"review_period_id"`
}

// WorkProductTaskVm represents a sub-task of a work product.
type WorkProductTaskVm struct {
	WorkProductTaskID string     `json:"work_product_task_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	CompletionDate    *time.Time `json:"completion_date,omitempty"`
	WorkProductID     string     `json:"work_product_id"`
	RecordStatus      string     `json:"record_status"`
}

// WorkProductEvaluationVm represents an evaluation of a work product.
type WorkProductEvaluationVm struct {
	WorkProductEvaluationID      string  `json:"work_product_evaluation_id"`
	WorkProductID                string  `json:"work_product_id"`
	Timeliness                   float64 `json:"timeliness"`
	TimelinessEvaluationOptionID string  `json:"timeliness_evaluation_option_id"`
	Quality                      float64 `json:"quality"`
	QualityEvaluationOptionID    string  `json:"quality_evaluation_option_id"`
	Output                       float64 `json:"output"`
	OutputEvaluationOptionID     string  `json:"output_evaluation_option_id"`
	Outcome                      float64 `json:"outcome"`
	EvaluatorStaffID             string  `json:"evaluator_staff_id"`
	EvaluatorName                string  `json:"evaluator_name"`
	IsReEvaluated                bool    `json:"is_re_evaluated"`
	TotalScore                   float64 `json:"total_score"`
	RecordStatus                 string  `json:"record_status"`
}

// ---------------------------------------------------------------------------
// Assigned Work Product View Models
// ---------------------------------------------------------------------------

// ProjectAssignedWorkProductVm represents a work product template assigned to a project.
type ProjectAssignedWorkProductVm struct {
	ProjectAssignedWorkProductID string    `json:"project_assigned_work_product_id"`
	WorkProductDefinitionID      string    `json:"work_product_definition_id"`
	Name                         string    `json:"name"`
	Description                  string    `json:"description"`
	ProjectID                    string    `json:"project_id"`
	ProjectName                  string    `json:"project_name"`
	ReviewPeriodID               string    `json:"review_period_id"`
	StartDate                    time.Time `json:"start_date"`
	EndDate                      time.Time `json:"end_date"`
	Deliverables                 string    `json:"deliverables"`
	RecordStatus                 string    `json:"record_status"`
}

// CommitteeAssignedWorkProductVm represents a work product template assigned to a committee.
type CommitteeAssignedWorkProductVm struct {
	CommitteeAssignedWorkProductID string    `json:"committee_assigned_work_product_id"`
	WorkProductDefinitionID        string    `json:"work_product_definition_id"`
	Name                           string    `json:"name"`
	Description                    string    `json:"description"`
	CommitteeID                    string    `json:"committee_id"`
	CommitteeName                  string    `json:"committee_name"`
	ReviewPeriodID                 string    `json:"review_period_id"`
	StartDate                      time.Time `json:"start_date"`
	EndDate                        time.Time `json:"end_date"`
	Deliverables                   string    `json:"deliverables"`
	RecordStatus                   string    `json:"record_status"`
}

// ---------------------------------------------------------------------------
// Adhoc Lead Request
// ---------------------------------------------------------------------------

// ChangeAdhocLeadRequestVm is the request body for changing a project manager or committee chairperson.
type ChangeAdhocLeadRequestVm struct {
	ReferenceID         string `json:"reference_id"`
	StaffID             string `json:"staff_id"`
	AdhocAssignmentType string `json:"adhoc_assignment_type"`
}
