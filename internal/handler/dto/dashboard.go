package dto

import "time"

// ===========================================================================
// Review Period Points Dashboard
// ===========================================================================

// ReviewPeriodPointsDashboardResponseVm holds point accumulation data for a staff member in a review period.
type ReviewPeriodPointsDashboardResponseVm struct {
	StaffID           string  `json:"staff_id"`
	ReviewPeriodID    string  `json:"review_period_id"`
	MaxPoints         float64 `json:"max_points"`
	AccumulatedPoints float64 `json:"accumulated_points"`
	DeductedPoints    float64 `json:"deducted_points"`
	ActualPoints      float64 `json:"actual_points"`
}

// ===========================================================================
// Review Period Work Product Dashboard
// ===========================================================================

// ReviewPeriodWorkProductDashboardResponseVm holds work product counters for a staff member.
type ReviewPeriodWorkProductDashboardResponseVm struct {
	StaffID                  string `json:"staff_id"`
	ReviewPeriodID           string `json:"review_period_id"`
	TotalWorkProducts        int    `json:"total_work_products"`
	PendingWorkProducts      int    `json:"pending_work_products"`
	ApprovedWorkProducts     int    `json:"approved_work_products"`
	RejectedWorkProducts     int    `json:"rejected_work_products"`
	EvaluatedWorkProducts    int    `json:"evaluated_work_products"`
	CompletedWorkProducts    int    `json:"completed_work_products"`
	OverdueWorkProducts      int    `json:"overdue_work_products"`
}

// WorkProductDashDetails holds detailed info for a single work product in the dashboard.
type WorkProductDashDetails struct {
	WorkProductID   string     `json:"work_product_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	WorkProductType string     `json:"work_product_type"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	MaxPoint        float64    `json:"max_point"`
	AccumulatedPoints float64  `json:"accumulated_points"`
	RecordStatus    string     `json:"record_status"`
	EvaluatedDate   *time.Time `json:"evaluated_date,omitempty"`
	ObjectiveName   string     `json:"objective_name"`
	CategoryName    string     `json:"category_name"`
}

// ReviewPeriodWorkProductDetailsDashboardResponseVm extends the counters with detailed work product data.
type ReviewPeriodWorkProductDetailsDashboardResponseVm struct {
	StaffID                  string                   `json:"staff_id"`
	ReviewPeriodID           string                   `json:"review_period_id"`
	TotalWorkProducts        int                      `json:"total_work_products"`
	PendingWorkProducts      int                      `json:"pending_work_products"`
	ApprovedWorkProducts     int                      `json:"approved_work_products"`
	RejectedWorkProducts     int                      `json:"rejected_work_products"`
	EvaluatedWorkProducts    int                      `json:"evaluated_work_products"`
	CompletedWorkProducts    int                      `json:"completed_work_products"`
	OverdueWorkProducts      int                      `json:"overdue_work_products"`
	WorkProductDetails       []WorkProductDashDetails  `json:"work_product_details"`
}

// ===========================================================================
// Staff Score Card
// ===========================================================================

// StaffScoreCardDetails holds the score breakdown for a staff member.
type StaffScoreCardDetails struct {
	CategoryID        string  `json:"category_id"`
	CategoryName      string  `json:"category_name"`
	Weight            float64 `json:"weight"`
	MaxPoints         float64 `json:"max_points"`
	AccumulatedPoints float64 `json:"accumulated_points"`
	WeightedScore     float64 `json:"weighted_score"`
	ObjectiveCount    int     `json:"objective_count"`
	WorkProductCount  int     `json:"work_product_count"`
}

// StaffScoreCardResponseVm wraps a single staff score card.
type StaffScoreCardResponseVm struct {
	BaseAPIResponse
	StaffID        string                  `json:"staff_id"`
	StaffName      string                  `json:"staff_name"`
	ReviewPeriodID string                  `json:"review_period_id"`
	TotalScore     float64                 `json:"total_score"`
	Details        []StaffScoreCardDetails `json:"details"`
}

// AllStaffScoreCardResponseVm wraps score cards for all staff members.
type AllStaffScoreCardResponseVm struct {
	GenericListResponseVm
	ScoreCards []StaffScoreCardResponseVm `json:"score_cards"`
}

// ===========================================================================
// Staff Annual Score Card
// ===========================================================================

// StaffLivingTheValueRatingsDetails holds living-the-value rating data.
type StaffLivingTheValueRatingsDetails struct {
	CompetencyName string  `json:"competency_name"`
	MaxScore       float64 `json:"max_score"`
	ActualScore    float64 `json:"actual_score"`
	WeightedScore  float64 `json:"weighted_score"`
}

// StaffAnnualScoreCardResponseVm wraps the annual score card for a staff member.
type StaffAnnualScoreCardResponseVm struct {
	BaseAPIResponse
	StaffID                string                              `json:"staff_id"`
	StaffName              string                              `json:"staff_name"`
	Year                   int                                 `json:"year"`
	TotalAnnualScore       float64                             `json:"total_annual_score"`
	CategoryDetails        []StaffScoreCardDetails             `json:"category_details"`
	LivingTheValueRatings  []StaffLivingTheValueRatingsDetails `json:"living_the_value_ratings"`
}

// ===========================================================================
// Organogram Performance Summary
// ===========================================================================

// OrganogramPerformanceSummaryDetails holds performance summary data for an organisation unit.
type OrganogramPerformanceSummaryDetails struct {
	OrganisationID     string  `json:"organisation_id"`
	OrganisationName   string  `json:"organisation_name"`
	OrganisationLevel  string  `json:"organisation_level"`
	TotalStaff         int     `json:"total_staff"`
	AverageScore       float64 `json:"average_score"`
	HighestScore       float64 `json:"highest_score"`
	LowestScore        float64 `json:"lowest_score"`
	CompletionRate     float64 `json:"completion_rate"`
}

// OrganogramPerformanceSummaryResponseVm wraps a single organogram performance summary.
type OrganogramPerformanceSummaryResponseVm struct {
	BaseAPIResponse
	Summary OrganogramPerformanceSummaryDetails `json:"summary"`
}

// OrganogramPerformanceSummaryListResponseVm wraps a list of organogram performance summaries.
type OrganogramPerformanceSummaryListResponseVm struct {
	GenericListResponseVm
	Summaries []OrganogramPerformanceSummaryDetails `json:"summaries"`
}

// ===========================================================================
// Period Score
// ===========================================================================

// PeriodScoreData holds score data for a review period.
type PeriodScoreData struct {
	PeriodScoreID  string  `json:"period_score_id"`
	StaffID        string  `json:"staff_id"`
	StaffName      string  `json:"staff_name"`
	ReviewPeriodID string  `json:"review_period_id"`
	TotalScore     float64 `json:"total_score"`
	MaxPoints      float64 `json:"max_points"`
	Percentage     float64 `json:"percentage"`
	Rank           int     `json:"rank"`
}

// PeriodScoreResponseVm wraps a single period score.
type PeriodScoreResponseVm struct {
	BaseAPIResponse
	Score PeriodScoreData `json:"score"`
}

// PeriodScoreListResponseVm wraps a list of period scores.
type PeriodScoreListResponseVm struct {
	GenericListResponseVm
	Scores []PeriodScoreData `json:"scores"`
}

// ===========================================================================
// Get Staff Review Period
// ===========================================================================

// GetStaffReviewPeriodResponseVm wraps a staff member's review period data.
type GetStaffReviewPeriodResponseVm struct {
	BaseAPIResponse
	StaffID        string                    `json:"staff_id"`
	StaffName      string                    `json:"staff_name"`
	ReviewPeriod   PerformanceReviewPeriodVm `json:"review_period"`
}

// ===========================================================================
// Period Objective Evaluation
// ===========================================================================

// PeriodObjectiveEvaluationData holds period objective evaluation data.
type PeriodObjectiveEvaluationData struct {
	PeriodObjectiveEvaluationID string  `json:"period_objective_evaluation_id"`
	EnterpriseObjectiveID       string  `json:"enterprise_objective_id"`
	EnterpriseObjectiveName     string  `json:"enterprise_objective_name"`
	ReviewPeriodID              string  `json:"review_period_id"`
	TotalOutcomeScore           float64 `json:"total_outcome_score"`
	OutcomeScore                float64 `json:"outcome_score"`
	RecordStatus                string  `json:"record_status"`
}

// PeriodObjectiveEvaluationResponseVm wraps a single period objective evaluation.
type PeriodObjectiveEvaluationResponseVm struct {
	BaseAPIResponse
	Evaluation PeriodObjectiveEvaluationData `json:"evaluation"`
}

// PeriodObjectiveEvaluationListResponseVm wraps a list of period objective evaluations.
type PeriodObjectiveEvaluationListResponseVm struct {
	GenericListResponseVm
	Evaluations []PeriodObjectiveEvaluationData `json:"evaluations"`
}

// ===========================================================================
// Period Objective Department Evaluation
// ===========================================================================

// PeriodObjectiveDepartmentEvaluationData holds department-level evaluation data.
type PeriodObjectiveDepartmentEvaluationData struct {
	PeriodObjectiveDepartmentEvaluationID string  `json:"period_objective_department_evaluation_id"`
	EnterpriseObjectiveID                 string  `json:"enterprise_objective_id"`
	EnterpriseObjectiveName               string  `json:"enterprise_objective_name"`
	ReviewPeriodID                        string  `json:"review_period_id"`
	DepartmentID                          string  `json:"department_id"`
	DepartmentName                        string  `json:"department_name"`
	TotalOutcomeScore                     float64 `json:"total_outcome_score"`
	OutcomeScore                          float64 `json:"outcome_score"`
	AllocatedOutcome                      float64 `json:"allocated_outcome"`
	OverallOutcomeScored                  float64 `json:"overall_outcome_scored"`
	RecordStatus                          string  `json:"record_status"`
}

// PeriodObjectiveDepartmentEvaluationResponseVm wraps a single department evaluation.
type PeriodObjectiveDepartmentEvaluationResponseVm struct {
	BaseAPIResponse
	Evaluation PeriodObjectiveDepartmentEvaluationData `json:"evaluation"`
}

// PeriodObjectiveDepartmentEvaluationListResponseVm wraps a list of department evaluations.
type PeriodObjectiveDepartmentEvaluationListResponseVm struct {
	GenericListResponseVm
	Evaluations []PeriodObjectiveDepartmentEvaluationData `json:"evaluations"`
}

// ===========================================================================
// Audit Log
// ===========================================================================

// AuditLogData holds audit log entry data.
type AuditLogData struct {
	AuditLogID   string    `json:"audit_log_id"`
	Action       string    `json:"action"`
	EntityType   string    `json:"entity_type"`
	EntityID     string    `json:"entity_id"`
	PerformedBy  string    `json:"performed_by"`
	PerformedByName string `json:"performed_by_name"`
	Description  string    `json:"description"`
	OldValues    string    `json:"old_values,omitempty"`
	NewValues    string    `json:"new_values,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// AuditLogResponseVm wraps a single audit log entry.
type AuditLogResponseVm struct {
	BaseAPIResponse
	AuditLog AuditLogData `json:"audit_log"`
}

// AuditLogListResponseVm wraps a list of audit log entries.
type AuditLogListResponseVm struct {
	GenericListResponseVm
	AuditLogs []AuditLogData `json:"audit_logs"`
}

// ===========================================================================
// View Model DTOs (Vm suffix – handler-layer representations)
// ===========================================================================

// ---------------------------------------------------------------------------
// Score Card View Models
// ---------------------------------------------------------------------------

// ScoreCardVm represents a single category score within a staff score card.
type ScoreCardVm struct {
	CategoryID                     string  `json:"category_id"`
	CategoryName                   string  `json:"category_name"`
	Weight                         float64 `json:"weight"`
	MaxPoints                      float64 `json:"max_points"`
	AccumulatedPoints              float64 `json:"accumulated_points"`
	WeightedScore                  float64 `json:"weighted_score"`
	WorkProductCount               int     `json:"work_product_count"`
	ObjectiveCount                 int     `json:"objective_count"`
	PercentageWorkProductsComplete float64 `json:"percentage_work_products_complete"`
}

// StaffAnnualScoreCardVm represents a staff member's annual score card.
type StaffAnnualScoreCardVm struct {
	StaffID          string        `json:"staff_id"`
	StaffName        string        `json:"staff_name"`
	Year             int           `json:"year"`
	TotalAnnualScore float64       `json:"total_annual_score"`
	PerformanceGrade string        `json:"performance_grade"`
	CategoryDetails  []ScoreCardVm `json:"category_details"`
}

// AllStaffScoreCardVm wraps score cards for all staff members in a review period.
type AllStaffScoreCardVm struct {
	ReviewPeriodID string                    `json:"review_period_id"`
	ReviewPeriod   string                    `json:"review_period"`
	Year           int                       `json:"year"`
	ScoreCards     []StaffAnnualScoreCardVm  `json:"score_cards"`
	TotalStaff     int                       `json:"total_staff"`
}

// ---------------------------------------------------------------------------
// Organogram Performance Summary View Model
// ---------------------------------------------------------------------------

// OrganogramPerformanceSummaryVm represents a performance summary for an organisational unit.
type OrganogramPerformanceSummaryVm struct {
	ReferenceID                          string  `json:"reference_id"`
	ReferenceName                        string  `json:"reference_name"`
	ManagerID                            string  `json:"manager_id"`
	ManagerName                          string  `json:"manager_name"`
	ReviewPeriodID                       string  `json:"review_period_id"`
	ReviewPeriodName                     string  `json:"review_period_name"`
	OrganogramLevel                      string  `json:"organogram_level"`
	TotalStaff                           int     `json:"total_staff"`
	ActualScore                          float64 `json:"actual_score"`
	PerformanceScore                     float64 `json:"performance_score"`
	MaxPoints                            float64 `json:"max_points"`
	TotalWorkProducts                    int     `json:"total_work_products"`
	TotalWorkProductsCompletedOnSchedule int     `json:"total_work_products_completed_on_schedule"`
	TotalWorkProductsBehindSchedule      int     `json:"total_work_products_behind_schedule"`
	PercentageWorkProductsClosed         float64 `json:"percentage_work_products_closed"`
	Total360Feedbacks                    int     `json:"total_360_feedbacks"`
	Completed360Feedbacks                int     `json:"completed_360_feedbacks"`
	Pending360Feedbacks                  int     `json:"pending_360_feedbacks"`
	TotalCompetencyGaps                  int     `json:"total_competency_gaps"`
	PercentageGapsClosure                float64 `json:"percentage_gaps_closure"`
	EarnedPerformanceGrade               string  `json:"earned_performance_grade"`
}

// ---------------------------------------------------------------------------
// Statistics View Models
// ---------------------------------------------------------------------------

// PerformanceStatisticsVm carries performance point statistics for a staff member's dashboard.
type PerformanceStatisticsVm struct {
	StaffID           string  `json:"staff_id"`
	ReviewPeriodID    string  `json:"review_period_id"`
	MaxPoints         float64 `json:"max_points"`
	AccumulatedPoints float64 `json:"accumulated_points"`
	DeductedPoints    float64 `json:"deducted_points"`
	ActualPoints      float64 `json:"actual_points"`
	PercentageScore   float64 `json:"percentage_score"`
}

// RequestStatisticsVm carries feedback request statistics for a staff member's dashboard.
type RequestStatisticsVm struct {
	StaffID                    string  `json:"staff_id"`
	ReviewPeriodID             string  `json:"review_period_id"`
	CompletedRequests          int     `json:"completed_requests"`
	PendingRequests            int     `json:"pending_requests"`
	BreachedRequests           int     `json:"breached_requests"`
	Pending360FeedbacksToTreat int     `json:"pending_360_feedbacks_to_treat"`
	CompletedOverdueRequests   int     `json:"completed_overdue_requests"`
	PendingOverdueRequests     int     `json:"pending_overdue_requests"`
	DeductedPoints             float64 `json:"deducted_points"`
}

// ---------------------------------------------------------------------------
// Period Score View Model
// ---------------------------------------------------------------------------

// PeriodScoreVm represents a staff member's score for a specific review period.
type PeriodScoreVm struct {
	PeriodScoreID   string    `json:"period_score_id"`
	StaffID         string    `json:"staff_id"`
	StaffName       string    `json:"staff_name"`
	ReviewPeriodID  string    `json:"review_period_id"`
	ReviewPeriod    string    `json:"review_period"`
	Year            int       `json:"year"`
	FinalScore      float64   `json:"final_score"`
	MaxPoints       float64   `json:"max_points"`
	ScorePercentage float64   `json:"score_percentage"`
	FinalGrade      string    `json:"final_grade"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	DepartmentID    string    `json:"department_id"`
	DepartmentName  string    `json:"department_name"`
	DivisionID      string    `json:"division_id"`
	DivisionName    string    `json:"division_name"`
	OfficeID        string    `json:"office_id"`
	OfficeName      string    `json:"office_name"`
	StaffGrade      string    `json:"staff_grade"`
}
