package dto

import "time"

// ===========================================================================
// Feedback Request Log VMs
// ===========================================================================

// FeedbackRequestLogVm represents a feedback request log entry.
type FeedbackRequestLogVm struct {
	FeedbackRequestLogID string     `json:"feedback_request_log_id"`
	RequestType          string     `json:"request_type"`
	RequestDate          time.Time  `json:"request_date"`
	RequestedBy          string     `json:"requested_by"`
	RequestedByName      string     `json:"requested_by_name"`
	TargetStaffID        string     `json:"target_staff_id"`
	TargetStaffName      string     `json:"target_staff_name"`
	ReviewPeriodID       string     `json:"review_period_id"`
	ReviewPeriodName     string     `json:"review_period_name"`
	Status               string     `json:"status"`
	DueDate              *time.Time `json:"due_date,omitempty"`
	CompletedDate        *time.Time `json:"completed_date,omitempty"`
	Description          string     `json:"description"`
	RecordStatus         string     `json:"record_status"`
}

// FeedbackRequestLogData holds feedback request log data for responses.
type FeedbackRequestLogData struct {
	FeedbackRequestLogID string     `json:"feedback_request_log_id"`
	RequestType          string     `json:"request_type"`
	RequestDate          time.Time  `json:"request_date"`
	RequestedBy          string     `json:"requested_by"`
	RequestedByName      string     `json:"requested_by_name"`
	TargetStaffID        string     `json:"target_staff_id"`
	TargetStaffName      string     `json:"target_staff_name"`
	ReviewPeriodID       string     `json:"review_period_id"`
	ReviewPeriodName     string     `json:"review_period_name"`
	Status               string     `json:"status"`
	DueDate              *time.Time `json:"due_date,omitempty"`
	CompletedDate        *time.Time `json:"completed_date,omitempty"`
	Description          string     `json:"description"`
	RecordStatus         string     `json:"record_status"`
}

// FeedbackRequestListResponseVm wraps a list of feedback requests.
type FeedbackRequestListResponseVm struct {
	GenericListResponseVm
	FeedbackRequests []FeedbackRequestLogData `json:"feedback_requests"`
}

// BreachedFeedbackRequestListResponseVm wraps a list of breached feedback requests.
type BreachedFeedbackRequestListResponseVm struct {
	GenericListResponseVm
	FeedbackRequests []FeedbackRequestLogData `json:"feedback_requests"`
}

// FeedbackRequestLogResponseVm wraps a single feedback request log.
type FeedbackRequestLogResponseVm struct {
	BaseAPIResponse
	FeedbackRequest FeedbackRequestLogData `json:"feedback_request"`
}

// FeedbackRequestDashboardData holds dashboard summary data for feedback requests.
type FeedbackRequestDashboardData struct {
	TotalRequests    int `json:"total_requests"`
	PendingRequests  int `json:"pending_requests"`
	CompletedRequests int `json:"completed_requests"`
	BreachedRequests int `json:"breached_requests"`
}

// FeedbackRequestDashboardResponseVm wraps feedback request dashboard data.
type FeedbackRequestDashboardResponseVm struct {
	BaseAPIResponse
	Dashboard FeedbackRequestDashboardData `json:"dashboard"`
}

// ===========================================================================
// 360 Review Request VMs
// ===========================================================================

// Initiate360ReviewRequestModel is the request to initiate a 360-degree review.
type Initiate360ReviewRequestModel struct {
	StaffID        string `json:"staff_id"`
	ReviewPeriodID string `json:"review_period_id"`
}

// Complete360ReviewRequestModel is the request to complete a 360-degree review.
type Complete360ReviewRequestModel struct {
	CompetencyReviewFeedbackID string `json:"competency_review_feedback_id"`
}

// ===========================================================================
// Competency Review Feedback VMs
// ===========================================================================

// CompetencyReviewFeedbackRequestModel is the request to create a competency review feedback.
type CompetencyReviewFeedbackRequestModel struct {
	StaffID        string  `json:"staff_id"`
	MaxPoints      float64 `json:"max_points"`
	ReviewPeriodID string  `json:"review_period_id"`
}

// CompetencyReviewerRequestModel is the request to assign a reviewer.
type CompetencyReviewerRequestModel struct {
	ReviewStaffID              string `json:"review_staff_id"`
	CompetencyReviewFeedbackID string `json:"competency_review_feedback_id"`
}

// CompetencyReviewerRatingRequestModel is the request to submit a reviewer rating.
type CompetencyReviewerRatingRequestModel struct {
	PmsCompetencyID              string  `json:"pms_competency_id"`
	FeedbackQuestionaireOptionID string  `json:"feedback_questionaire_option_id"`
	Rating                       float64 `json:"rating"`
	CompetencyReviewerID         string  `json:"competency_reviewer_id"`
}

// ===========================================================================
// Competency Review Feedback Response VMs & Data Structs
// ===========================================================================

// CompetencyReviewFeedbackData holds competency review feedback data.
type CompetencyReviewFeedbackData struct {
	CompetencyReviewFeedbackID string     `json:"competency_review_feedback_id"`
	StaffID                    string     `json:"staff_id"`
	StaffName                  string     `json:"staff_name"`
	MaxPoints                  float64    `json:"max_points"`
	FinalScore                 float64    `json:"final_score"`
	ReviewPeriodID             string     `json:"review_period_id"`
	ReviewPeriodName           string     `json:"review_period_name"`
	RecordStatus               string     `json:"record_status"`
	CompletedDate              *time.Time `json:"completed_date,omitempty"`
}

// CompetencyReviewFeedbackResponseVm wraps a single competency review feedback.
type CompetencyReviewFeedbackResponseVm struct {
	BaseAPIResponse
	Feedback CompetencyReviewFeedbackData `json:"feedback"`
}

// CompetencyReviewFeedbackListResponseVm wraps a list of competency review feedbacks.
type CompetencyReviewFeedbackListResponseVm struct {
	GenericListResponseVm
	Feedbacks []CompetencyReviewFeedbackData `json:"feedbacks"`
}

// CompetencyReviewerData holds reviewer data for competency reviews.
type CompetencyReviewerData struct {
	CompetencyReviewerID       string     `json:"competency_reviewer_id"`
	ReviewStaffID              string     `json:"review_staff_id"`
	ReviewStaffName            string     `json:"review_staff_name"`
	CompetencyReviewFeedbackID string     `json:"competency_review_feedback_id"`
	RecordStatus               string     `json:"record_status"`
	CompletedDate              *time.Time `json:"completed_date,omitempty"`
}

// CompetencyReviewerListResponseVm wraps a list of competency reviewers.
type CompetencyReviewerListResponseVm struct {
	GenericListResponseVm
	Reviewers []CompetencyReviewerData `json:"reviewers"`
}

// CompetencyReviewerResponseVm wraps a single competency reviewer.
type CompetencyReviewerResponseVm struct {
	BaseAPIResponse
	Reviewer CompetencyReviewerData `json:"reviewer"`
}

// CompetencyReviewerRatingData holds rating data from a reviewer.
type CompetencyReviewerRatingData struct {
	CompetencyReviewerRatingID   string  `json:"competency_reviewer_rating_id"`
	PmsCompetencyID              string  `json:"pms_competency_id"`
	PmsCompetencyName            string  `json:"pms_competency_name"`
	FeedbackQuestionaireOptionID string  `json:"feedback_questionaire_option_id"`
	OptionStatement              string  `json:"option_statement"`
	Rating                       float64 `json:"rating"`
	CompetencyReviewerID         string  `json:"competency_reviewer_id"`
}

// CompetencyReviewerRatingListResponseVm wraps a list of reviewer ratings.
type CompetencyReviewerRatingListResponseVm struct {
	GenericListResponseVm
	Ratings []CompetencyReviewerRatingData `json:"ratings"`
}

// CompetencyReviewerRatingResponseVm wraps a single reviewer rating.
type CompetencyReviewerRatingResponseVm struct {
	BaseAPIResponse
	Rating CompetencyReviewerRatingData `json:"rating"`
}

// ===========================================================================
// Questionnaire Response VMs & Data Structs
// ===========================================================================

// PmsCompetencyData holds PMS competency data for responses.
type PmsCompetencyData struct {
	PmsCompetencyID  string `json:"pms_competency_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ObjectCategoryID string `json:"object_category_id"`
	RecordStatus     string `json:"record_status"`
}

// FeedbackQuestionaireData holds questionnaire data for responses.
type FeedbackQuestionaireData struct {
	FeedbackQuestionaireID string                         `json:"feedback_questionaire_id"`
	Question               string                         `json:"question"`
	Description            string                         `json:"description"`
	PmsCompetencyID        string                         `json:"pms_competency_id"`
	PmsCompetencyName      string                         `json:"pms_competency_name"`
	Options                []FeedbackQuestionaireOptionData `json:"options,omitempty"`
}

// FeedbackQuestionaireOptionData holds a single option for a questionnaire.
type FeedbackQuestionaireOptionData struct {
	FeedbackQuestionaireOptionID string  `json:"feedback_questionaire_option_id"`
	OptionStatement              string  `json:"option_statement"`
	Description                  string  `json:"description"`
	Score                        float64 `json:"score"`
	QuestionID                   string  `json:"question_id"`
}

// QuestionnaireListResponseVm wraps a list of questionnaires grouped by competency.
type QuestionnaireListResponseVm struct {
	GenericListResponseVm
	Competencies []PmsCompetencyData      `json:"competencies"`
	Questionnaires []FeedbackQuestionaireData `json:"questionnaires"`
}

// ===========================================================================
// Competency Gap Closure VMs
// ===========================================================================

// CompetencyGapClosureRequestModel is the request for recording a competency gap closure.
type CompetencyGapClosureRequestModel struct {
	StaffID             string  `json:"staff_id"`
	MaxPoints           float64 `json:"max_points"`
	FinalScore          float64 `json:"final_score"`
	ReviewPeriodID      string  `json:"review_period_id"`
	ObjectiveCategoryID string  `json:"objective_category_id"`
}

// CompetencyGapClosureData holds competency gap closure data for responses.
type CompetencyGapClosureData struct {
	CompetencyGapClosureID string  `json:"competency_gap_closure_id"`
	StaffID                string  `json:"staff_id"`
	StaffName              string  `json:"staff_name"`
	MaxPoints              float64 `json:"max_points"`
	FinalScore             float64 `json:"final_score"`
	ReviewPeriodID         string  `json:"review_period_id"`
	ObjectiveCategoryID    string  `json:"objective_category_id"`
	RecordStatus           string  `json:"record_status"`
}

// CompetencyGapClosureResponseVm wraps a single competency gap closure.
type CompetencyGapClosureResponseVm struct {
	BaseAPIResponse
	GapClosure CompetencyGapClosureData `json:"gap_closure"`
}

// CompetencyGapClosureListResponseVm wraps a list of competency gap closures.
type CompetencyGapClosureListResponseVm struct {
	GenericListResponseVm
	GapClosures []CompetencyGapClosureData `json:"gap_closures"`
}

// ===========================================================================
// View Model DTOs (Vm suffix – handler-layer representations)
// ===========================================================================

// ---------------------------------------------------------------------------
// Treat Feedback Request
// ---------------------------------------------------------------------------

// TreatFeedbackRequestVm is the request body for treating (approving/rejecting) a feedback request.
type TreatFeedbackRequestVm struct {
	FeedbackRequestLogID string `json:"feedback_request_log_id"`
	Comment              string `json:"comment"`
	Attachment           string `json:"attachment,omitempty"`
	Action               string `json:"action"`
}

// ---------------------------------------------------------------------------
// Competency Review Feedback View Models
// ---------------------------------------------------------------------------

// CompetencyReviewFeedbackVm represents the full view of a competency review feedback.
type CompetencyReviewFeedbackVm struct {
	CompetencyReviewFeedbackID string                `json:"competency_review_feedback_id"`
	StaffID                    string                `json:"staff_id"`
	StaffName                  string                `json:"staff_name"`
	MaxPoints                  float64               `json:"max_points"`
	FinalScore                 float64               `json:"final_score"`
	FinalScorePercentage       float64               `json:"final_score_percentage"`
	ReviewPeriodID             string                `json:"review_period_id"`
	ReviewPeriodName           string                `json:"review_period_name"`
	RecordStatus               string                `json:"record_status"`
	CompletedDate              *time.Time            `json:"completed_date,omitempty"`
	Reviewers                  []CompetencyReviewerVm `json:"reviewers,omitempty"`
}

// CompetencyReviewerVm represents a single reviewer in a 360-degree feedback cycle.
type CompetencyReviewerVm struct {
	CompetencyReviewerID       string     `json:"competency_reviewer_id"`
	ReviewStaffID              string     `json:"review_staff_id"`
	ReviewStaffName            string     `json:"review_staff_name"`
	FinalRating                float64    `json:"final_rating"`
	CompetencyReviewFeedbackID string     `json:"competency_review_feedback_id"`
	RecordStatus               string     `json:"record_status"`
	CompletedDate              *time.Time `json:"completed_date,omitempty"`
}

// ---------------------------------------------------------------------------
// 360 Review Initiation / Completion View Models
// ---------------------------------------------------------------------------

// Initiate360ReviewVm is the request body for initiating a 360-degree review cycle.
type Initiate360ReviewVm struct {
	StaffIDs       []string `json:"staff_ids"`
	ReviewPeriodID string   `json:"review_period_id"`
}

// Complete360ReviewVm is the request body for completing a 360-degree review cycle.
type Complete360ReviewVm struct {
	CompetencyReviewFeedbackID string `json:"competency_review_feedback_id"`
	ReviewPeriodID             string `json:"review_period_id"`
}

// ---------------------------------------------------------------------------
// Questionnaire View Model
// ---------------------------------------------------------------------------

// QuestionnaireVm represents a feedback questionnaire with its options.
type QuestionnaireVm struct {
	FeedbackQuestionaireID string                        `json:"feedback_questionaire_id"`
	Question               string                        `json:"question"`
	Description            string                        `json:"description"`
	PmsCompetencyID        string                        `json:"pms_competency_id"`
	PmsCompetencyName      string                        `json:"pms_competency_name"`
	RecordStatus           string                        `json:"record_status"`
	Options                []FeedbackQuestionaireOptionData `json:"options,omitempty"`
}

// ---------------------------------------------------------------------------
// Competency Gap Closure View Model
// ---------------------------------------------------------------------------

// CompetencyGapClosureVm represents the full view of a competency gap closure record.
type CompetencyGapClosureVm struct {
	CompetencyGapClosureID string  `json:"competency_gap_closure_id"`
	StaffID                string  `json:"staff_id"`
	StaffName              string  `json:"staff_name"`
	MaxPoints              float64 `json:"max_points"`
	FinalScore             float64 `json:"final_score"`
	ReviewPeriodID         string  `json:"review_period_id"`
	ReviewPeriodName       string  `json:"review_period_name"`
	ObjectiveCategoryID    string  `json:"objective_category_id"`
	ObjectiveCategoryName  string  `json:"objective_category_name"`
	RecordStatus           string  `json:"record_status"`
}
