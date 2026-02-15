package dto

import "time"

// ---------------------------------------------------------------------------
// Review Period Request VMs
// ---------------------------------------------------------------------------

// CreateNewReviewPeriodVm is the request body for creating a new review period.
type CreateNewReviewPeriodVm struct {
	Year            int       `json:"year"`
	Range           string    `json:"range"`
	RangeValue      int       `json:"range_value"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ShortName       string    `json:"short_name"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	MaxPoints       float64   `json:"max_points"`
	MinNoObjectives int       `json:"min_no_objectives"`
	MaxNoObjectives int       `json:"max_no_objectives"`
	StrategyID      string    `json:"strategy_id"`
}

// ReviewPeriodRequestVm extends the create model with identifier and comment.
type ReviewPeriodRequestVm struct {
	PeriodID        string    `json:"period_id"`
	Year            int       `json:"year"`
	Range           string    `json:"range"`
	RangeValue      int       `json:"range_value"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ShortName       string    `json:"short_name"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	MaxPoints       float64   `json:"max_points"`
	MinNoObjectives int       `json:"min_no_objectives"`
	MaxNoObjectives int       `json:"max_no_objectives"`
	StrategyID      string    `json:"strategy_id"`
	Comment         string    `json:"comment,omitempty"`
}

// PerformanceReviewPeriodVm represents the full view of a performance review period.
type PerformanceReviewPeriodVm struct {
	PeriodID        string              `json:"period_id"`
	Year            int                 `json:"year"`
	Range           string              `json:"range"`
	RangeValue      int                 `json:"range_value"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	ShortName       string              `json:"short_name"`
	StartDate       time.Time           `json:"start_date"`
	EndDate         time.Time           `json:"end_date"`
	MaxPoints       float64             `json:"max_points"`
	MinNoObjectives int                 `json:"min_no_objectives"`
	MaxNoObjectives int                 `json:"max_no_objectives"`
	StrategyID      string              `json:"strategy_id"`
	RecordStatus    string              `json:"record_status"`
	PeriodObjectives []PeriodObjectiveVm `json:"period_objectives,omitempty"`
}

// ---------------------------------------------------------------------------
// Period Objective VMs
// ---------------------------------------------------------------------------

// AddPeriodObjectiveVm is the request to associate objectives with a review period.
type AddPeriodObjectiveVm struct {
	ReviewPeriodID string   `json:"review_period_id"`
	ObjectiveIDs   []string `json:"objective_ids"`
}

// PeriodObjectiveVm represents an objective linked to a review period.
type PeriodObjectiveVm struct {
	PeriodObjectiveID string `json:"period_objective_id"`
	ObjectiveID       string `json:"objective_id"`
	ReviewPeriodID    string `json:"review_period_id"`
}

// ---------------------------------------------------------------------------
// Planned Objective VMs
// ---------------------------------------------------------------------------

// AddReviewPeriodIndividualPlannedObjectiveRequestModel is the request to plan an objective for a staff member.
type AddReviewPeriodIndividualPlannedObjectiveRequestModel struct {
	ObjectiveID    string `json:"objective_id"`
	StaffID        string `json:"staff_id"`
	ObjectiveLevel string `json:"objective_level"`
	StaffJobRole   string `json:"staff_job_role"`
	ReviewPeriodID string `json:"review_period_id"`
	Remark         string `json:"remark,omitempty"`
}

// ReviewPeriodIndividualPlannedObjectiveRequestModel extends the add model with identifier and comment.
type ReviewPeriodIndividualPlannedObjectiveRequestModel struct {
	PlannedObjectiveID string `json:"planned_objective_id"`
	ObjectiveID        string `json:"objective_id"`
	StaffID            string `json:"staff_id"`
	ObjectiveLevel     string `json:"objective_level"`
	StaffJobRole       string `json:"staff_job_role"`
	ReviewPeriodID     string `json:"review_period_id"`
	Remark             string `json:"remark,omitempty"`
	Comment            string `json:"comment,omitempty"`
}

// ---------------------------------------------------------------------------
// Review Period Extension VMs
// ---------------------------------------------------------------------------

// ReviewPeriodExtensionRequestModel is the request for extending a review period.
type ReviewPeriodExtensionRequestModel struct {
	ReviewPeriodID  string    `json:"review_period_id"`
	TargetType      string    `json:"target_type"`
	TargetReference string    `json:"target_reference"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
}

// ReviewPeriod360ReviewRequestModel is the request for configuring 360 review on a period.
type ReviewPeriod360ReviewRequestModel struct {
	ReviewPeriodID  string `json:"review_period_id"`
	TargetType      string `json:"target_type"`
	TargetReference string `json:"target_reference"`
}

// ---------------------------------------------------------------------------
// Review Period Response VMs
// ---------------------------------------------------------------------------

// ReviewPeriodResponseVm wraps a single review period in a standard response.
type ReviewPeriodResponseVm struct {
	BaseAPIResponse
	ReviewPeriod PerformanceReviewPeriodVm `json:"review_period"`
}

// GetAllReviewPeriodResponseVm wraps a list of review periods.
type GetAllReviewPeriodResponseVm struct {
	GenericListResponseVm
	ReviewPeriods []PerformanceReviewPeriodVm `json:"review_periods"`
}

// PerformanceReviewPeriodResponseVm wraps a performance review period with objectives.
type PerformanceReviewPeriodResponseVm struct {
	BaseAPIResponse
	ReviewPeriod PerformanceReviewPeriodVm `json:"review_period"`
}

// ---------------------------------------------------------------------------
// Planned Objective Response VMs & Data Structs
// ---------------------------------------------------------------------------

// PlannedObjectiveData holds the detailed data for a planned objective.
type PlannedObjectiveData struct {
	PlannedObjectiveID string  `json:"planned_objective_id"`
	ObjectiveID        string  `json:"objective_id"`
	ObjectiveName      string  `json:"objective_name"`
	ObjectiveLevel     string  `json:"objective_level"`
	StaffID            string  `json:"staff_id"`
	StaffName          string  `json:"staff_name"`
	StaffJobRole       string  `json:"staff_job_role"`
	ReviewPeriodID     string  `json:"review_period_id"`
	Remark             string  `json:"remark"`
	RecordStatus       string  `json:"record_status"`
	Kpi                string  `json:"kpi"`
	Target             string  `json:"target"`
	CategoryID         string  `json:"category_id"`
	CategoryName       string  `json:"category_name"`
	MaxPoints          float64 `json:"max_points"`
	AccumulatedPoints  float64 `json:"accumulated_points"`
}

// PlannedOperationalObjectivesResponseVm wraps a list of planned operational objectives.
type PlannedOperationalObjectivesResponseVm struct {
	GenericListResponseVm
	PlannedObjectives []PlannedObjectiveData `json:"planned_objectives"`
}

// PlannedObjectiveResponseVm wraps a single planned objective.
type PlannedObjectiveResponseVm struct {
	BaseAPIResponse
	PlannedObjective PlannedObjectiveData `json:"planned_objective"`
}

// ---------------------------------------------------------------------------
// Review Period Objectives Response VMs & Data Structs
// ---------------------------------------------------------------------------

// EnterpriseObjectiveData holds enterprise-level objective data in a response.
type EnterpriseObjectiveData struct {
	EnterpriseObjectiveID string `json:"enterprise_objective_id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Kpi                   string `json:"kpi"`
	Target                string `json:"target"`
	CategoryID            string `json:"category_id"`
	CategoryName          string `json:"category_name"`
	StrategyID            string `json:"strategy_id"`
	StrategyName          string `json:"strategy_name"`
	RecordStatus          string `json:"record_status"`
}

// CascadedObjectiveData represents an objective that has been cascaded down the hierarchy.
type CascadedObjectiveData struct {
	ObjectiveID      string `json:"objective_id"`
	ObjectiveLevel   string `json:"objective_level"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Kpi              string `json:"kpi"`
	Target           string `json:"target"`
	ParentObjectiveID string `json:"parent_objective_id"`
	OrganisationID   string `json:"organisation_id"`
	OrganisationName string `json:"organisation_name"`
	RecordStatus     string `json:"record_status"`
}

// ReviewPeriodObjectivesResponseVm wraps a list of objectives for a review period.
type ReviewPeriodObjectivesResponseVm struct {
	GenericListResponseVm
	EnterpriseObjectives []EnterpriseObjectiveData `json:"enterprise_objectives"`
	CascadedObjectives   []CascadedObjectiveData   `json:"cascaded_objectives,omitempty"`
}

// ---------------------------------------------------------------------------
// Category Definition Response VMs & Data Structs
// ---------------------------------------------------------------------------

// CategoryDefinitionData holds a category definition in a response.
type CategoryDefinitionData struct {
	DefinitionID           string  `json:"definition_id"`
	ObjectiveCategoryID    string  `json:"objective_category_id"`
	ObjectiveCategoryName  string  `json:"objective_category_name"`
	ReviewPeriodID         string  `json:"review_period_id"`
	Weight                 float64 `json:"weight"`
	MaxNoObjectives        int     `json:"max_no_objectives"`
	MaxNoWorkProduct       int     `json:"max_no_work_product"`
	MaxPoints              float64 `json:"max_points"`
	IsCompulsory           bool    `json:"is_compulsory"`
	EnforceWorkProductLimit bool   `json:"enforce_work_product_limit"`
	Description            string  `json:"description"`
	GradeGroupID           string  `json:"grade_group_id"`
	RecordStatus           string  `json:"record_status"`
}

// ReviewPeriodCategoryDefinitionResponseVm wraps category definitions for a review period.
type ReviewPeriodCategoryDefinitionResponseVm struct {
	GenericListResponseVm
	CategoryDefinitions []CategoryDefinitionData `json:"category_definitions"`
}

// ---------------------------------------------------------------------------
// Operational Objectives Response VMs & Data Structs
// ---------------------------------------------------------------------------

// OperationalObjectiveData represents an operational objective in a response.
type OperationalObjectiveData struct {
	ObjectiveID    string `json:"objective_id"`
	ObjectiveLevel string `json:"objective_level"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Kpi            string `json:"kpi"`
	Target         string `json:"target"`
	CategoryID     string `json:"category_id"`
	CategoryName   string `json:"category_name"`
	RecordStatus   string `json:"record_status"`
}

// OperationalObjectivesResponseVm wraps a list of operational objectives.
type OperationalObjectivesResponseVm struct {
	GenericListResponseVm
	Objectives []OperationalObjectiveData `json:"objectives"`
}

// ---------------------------------------------------------------------------
// Review Period Extension Response VMs
// ---------------------------------------------------------------------------

// ReviewPeriodExtensionData holds extension data for a review period.
type ReviewPeriodExtensionData struct {
	ExtensionID     string    `json:"extension_id"`
	ReviewPeriodID  string    `json:"review_period_id"`
	TargetType      string    `json:"target_type"`
	TargetReference string    `json:"target_reference"`
	Description     string    `json:"description"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	RecordStatus    string    `json:"record_status"`
}

// ReviewPeriodExtensionResponseVm wraps a single extension record.
type ReviewPeriodExtensionResponseVm struct {
	BaseAPIResponse
	Extension ReviewPeriodExtensionData `json:"extension"`
}

// ReviewPeriodExtensionListResponseVm wraps a list of extension records.
type ReviewPeriodExtensionListResponseVm struct {
	GenericListResponseVm
	Extensions []ReviewPeriodExtensionData `json:"extensions"`
}
