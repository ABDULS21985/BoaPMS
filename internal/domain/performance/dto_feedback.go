package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ---------------------------------------------------------------------------
// Feedback Request Log DTOs  (source: ResponseViewModels.cs – FeedbackRequestLogVm)
// ---------------------------------------------------------------------------

// FeedbackRequestLogVm is the read/response DTO for a feedback request log.
type FeedbackRequestLogVm struct {
	BaseAuditVm
	FeedbackRequestLogID   string                    `json:"feedbackRequestLogId"`
	FeedbackRequestType    enums.FeedbackRequestType `json:"feedBackRequestType"`
	ReferenceID            string                    `json:"referenceId"`
	TimeInitiated          time.Time                 `json:"timeInitiated"`
	AssignedStaffID        string                    `json:"assignedStaffId"`
	AssignedStaffName      string                    `json:"assignedStaffName"`
	RequestOwnerStaffID    string                    `json:"requestOwnerStaffId"`
	RequestOwnerStaffName  string                    `json:"requestOwnerStaffName"`
	TimeCompleted          *time.Time                `json:"timeCompleted"`
	HasSLA                 bool                      `json:"hasSla"`
	IsBreached             bool                      `json:"isBreached"`
	ReviewPeriodID         string                    `json:"reviewPeriodId"`
	RequestOwnerComment    *string                   `json:"requestOwnerComment"`
	RequestOwnerAttachment *string                   `json:"requestOwnerAttachment"`
	AssignedStaffComment   *string                   `json:"assignedStaffComment"`
	AssignedStaffAttachment *string                  `json:"assignedStaffAttachment"`
}

// FeedbackRequestListResponseVm wraps a list of raw FeedbackRequestLog records.
type FeedbackRequestListResponseVm struct {
	GenericListResponseVm
	Requests []FeedbackRequestLog `json:"requests"`
}

// BreachedFeedbackRequestListResponseVm wraps breached feedback request logs.
type BreachedFeedbackRequestListResponseVm struct {
	GenericListResponseVm
	Requests []FeedbackRequestLogVm `json:"requests"`
}

// FeedbackRequestLogResponseVm wraps a single FeedbackRequestLog record.
type FeedbackRequestLogResponseVm struct {
	BaseAPIResponse
	Request *FeedbackRequestLog `json:"request"`
}

// FeedbackRequestModel is the request DTO for assigning a feedback request.
type FeedbackRequestModel struct {
	RequestID  string `json:"requestId"`
	AssigneeID string `json:"assigneeId"`
}

// TreatFeedbackRequestModel is the request DTO for treating a feedback request.
type TreatFeedbackRequestModel struct {
	RequestID     string              `json:"requestId" validate:"required"`
	OperationType enums.OperationType `json:"operationType" validate:"required"`
	Comment       string              `json:"comment" validate:"required"`
}

// AssignedRequestModel is the request DTO for an assigned feedback request.
type AssignedRequestModel struct {
	RequestID           string                    `json:"requestId" validate:"required"`
	ReviewPeriodID      string                    `json:"reviewPeriodId" validate:"required"`
	ReferenceID         string                    `json:"referenceId" validate:"required"`
	FeedbackRequestType enums.FeedbackRequestType `json:"feedBackRequestType" validate:"required"`
	Requestor           string                    `json:"requestor"`
}

// ---------------------------------------------------------------------------
// Feedback Request Dashboard  (source: ResponseViewModels.cs)
// ---------------------------------------------------------------------------

// FeedbackRequestDashboardResponseVm returns feedback request statistics.
type FeedbackRequestDashboardResponseVm struct {
	BaseAPIResponse
	StaffID                     string  `json:"staffId"`
	ReviewPeriodID              string  `json:"reviewPeriodId"`
	CompletedRequests           int     `json:"completedRequests"`
	CompletedOverdueRequests    int     `json:"completedOverdueRequests"`
	PendingRequests             int     `json:"pendingRequests"`
	PendingOverdueRequests      int     `json:"pendingOverdueRequests"`
	BreachedRequests            int     `json:"breachedRequests"`
	Pending360FeedbacksToTreat  int     `json:"pending360FeedbacksToTreat"`
	DeductedPoints              float64 `json:"deductedPoints"`
}

// ---------------------------------------------------------------------------
// Staff Pending Request  (source: FrontendCustomVms/ResponseFormVms.cs)
// ---------------------------------------------------------------------------

// StaffPendingRequestVm represents a pending request assigned to a staff member.
type StaffPendingRequestVm struct {
	FeedbackRequestLogID string    `json:"feedbackRequestLogId"`
	FeedbackRequestType  int       `json:"feedBackRequestType"`
	ReferenceID          string    `json:"referenceId"`
	TimeInitiated        time.Time `json:"timeInitiated"`
	AssignedStaffID      string    `json:"assignedStaffId"`
	RequestOwnerStaffID  string    `json:"requestOwnerStaffId"`
	RecordStatus         int       `json:"recordStatus"`
	HasSLA               bool      `json:"hasSla"`
	ID                   int       `json:"id"`
	IsActive             bool      `json:"isActive"`
}

// GetStaffPendingRequestVm wraps a list of pending requests for a staff member.
type GetStaffPendingRequestVm struct {
	GenericListResponseVm
	Requests []StaffPendingRequestVm `json:"requests"`
}

// ---------------------------------------------------------------------------
// Staff Search  (source: FrontendCustomVms – derived from common search patterns)
// ---------------------------------------------------------------------------

// StaffSearchVm is the search/filter DTO for staff-related queries.
type StaffSearchVm struct {
	BasePagedData
	StaffID      string `json:"staffId"`
	SearchString string `json:"searchString"`
	DepartmentID *int   `json:"departmentId"`
	DivisionID   *int   `json:"divisionId"`
	OfficeID     *int   `json:"officeId"`
}

// ---------------------------------------------------------------------------
// Feedback Log DTOs  (source: Feedbacks360Vms/FeedbackLogVm.cs)
// ---------------------------------------------------------------------------

// FeedbackLogVm is the read/response DTO for a feedback request log entry.
type FeedbackLogVm struct {
	BaseAuditVm
	FeedbackRequestLogID string                    `json:"feedbackRequestLogId"`
	FeedbackRequestType  enums.FeedbackRequestType `json:"feedBackRequestType"`
	ReferenceID          string                    `json:"referenceId"`
	TimeInitiated        time.Time                 `json:"timeInitiated"`
	AssignedStaffID      string                    `json:"assignedStaffId"`
	RequestOwnerStaffID  string                    `json:"requestOwnerStaffId"`
	TimeCompleted        *time.Time                `json:"timeCompleted"`
	RecordStatus         enums.Status              `json:"recordStatus"`
}

// ---------------------------------------------------------------------------
// Questionnaire DTOs  (source: Feedbacks360Vms/QuestionnaireVm.cs)
// ---------------------------------------------------------------------------

// QuestionnaireVm is the read/response DTO for a feedback questionnaire.
type QuestionnaireVm struct {
	BaseAuditVm
	FeedbackQuestionaireID string                       `json:"feedbackQuestionaireId"`
	Question               string                       `json:"question" validate:"required"`
	Description            string                       `json:"description" validate:"required"`
	RecordStatus           enums.Status                 `json:"recordStatus" validate:"required"`
	Options                []FeedbackQuestionaireOption `json:"options"`
}

// ---------------------------------------------------------------------------
// Questionnaire Option DTOs  (source: Feedbacks360Vms/QuestionnaireOptionVm.cs)
// ---------------------------------------------------------------------------

// QuestionnaireOptionVm is the read/response DTO for a feedback questionnaire option.
type QuestionnaireOptionVm struct {
	BaseAuditVm
	FeedbackQuestionaireOptionID string  `json:"feedbackQuestionaireOptionId"`
	OptionStatement              string  `json:"optionStatement" validate:"required"`
	Description                  string  `json:"description"`
	Score                        float64 `json:"score" validate:"required"`
	QuestionID                   string  `json:"questionId" validate:"required"`
}

// ---------------------------------------------------------------------------
// Feedback Questionaire Setup VMs  (source: SetupVm/FeedbackQuestionaireVm.cs)
// ---------------------------------------------------------------------------

// FeedbackQuestionaireVm is the setup DTO for a feedback questionnaire.
type FeedbackQuestionaireVm struct {
	BaseAuditVm
	FeedbackQuestionaireID string                      `json:"feedbackQuestionaireId"`
	Question               string                      `json:"question" validate:"required"`
	Description            string                      `json:"description"`
	PmsCompetencyID        string                      `json:"pmsCompetencyId"`
	PmsCompetencyName      string                      `json:"pmsCompetencyName"`
	RecordStatus           enums.Status                `json:"recordStatus"`
	Category               string                      `json:"category"`
	QuestionID             string                      `json:"questionId"`
	OptionStatement        string                      `json:"optionStatement"`
	Score                  float64                     `json:"score"`
	IsValidRecord          bool                        `json:"isValidRecord"`
	IsSuccess              *bool                       `json:"isSuccess"`
	Message                string                      `json:"message"`
	IsProcessed            *bool                       `json:"isProcessed"`
	IsSelected             bool                        `json:"isSelected"`
	Options                []FeedbackQuestionaireOptionVm `json:"options"`
}

// FeedbackQuestionaireOptionVm is the setup DTO for a questionnaire option.
type FeedbackQuestionaireOptionVm struct {
	BaseAuditVm
	FeedbackQuestionaireOptionID string  `json:"feedbackQuestionaireOptionId"`
	OptionStatement              string  `json:"optionStatement" validate:"required"`
	Description                  string  `json:"description"`
	Score                        float64 `json:"score" validate:"required"`
	QuestionID                   string  `json:"questionId"`
}

// SearchFeedbackQuestionaireVm is the search/filter DTO for feedback questionnaires.
type SearchFeedbackQuestionaireVm struct {
	BasePagedData
	DepartmentID    *int    `json:"departmentId"`
	DivisionID      *int    `json:"divisionId"`
	OfficeID        *int    `json:"officeId"`
	JobRoleID       *int    `json:"jobRoleId"`
	SearchString    string  `json:"searchString"`
	Question        string  `json:"question"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	OptionStatement string  `json:"optionStatement"`
	Score           float64 `json:"score"`
}

// FeedbackQuestionaireListResponseVm wraps a list of feedback questionnaires.
type FeedbackQuestionaireListResponseVm struct {
	GenericListResponseVm
	FeedbackQuestionaires []FeedbackQuestionaireVm `json:"feedbackQuestionaires"`
}

// FeedbackQuestionaireResponseVm wraps a list of feedback questionnaires (alternate).
type FeedbackQuestionaireResponseVm struct {
	GenericListResponseVm
	FeedbackQuestionaires []FeedbackQuestionaireVm `json:"feedbackQuestionaires"`
}

// ---------------------------------------------------------------------------
// 360 Review – Competency Review Feedback DTOs  (source: ResponseViewModels.cs)
// ---------------------------------------------------------------------------

// CompetencyReviewFeedbackResponseVm wraps a single competency review feedback.
type CompetencyReviewFeedbackResponseVm struct {
	BaseAPIResponse
	CompetencyReviewFeedback *CompetencyReviewFeedbackData `json:"competencyReviewFeedback"`
}

// CompetencyReviewFeedbackDetailsResponseVm wraps detailed competency review feedback.
type CompetencyReviewFeedbackDetailsResponseVm struct {
	BaseAPIResponse
	CompetencyReviewFeedback *CompetencyReviewFeedbackDetails `json:"competencyReviewFeedback"`
}

// CompetencyReviewFeedbackDetails holds detailed competency feedback for a staff member.
type CompetencyReviewFeedbackDetails struct {
	CompetencyReviewFeedbackID string                                `json:"competencyReviewFeedbackId"`
	StaffID                    string                                `json:"staffId"`
	StaffName                  string                                `json:"staffName"`
	OfficeID                   int                                   `json:"officeID"`
	OfficeCode                 string                                `json:"officeCode"`
	OfficeName                 string                                `json:"officeName"`
	DivisionID                 int                                   `json:"divisionId"`
	DivisionCode               string                                `json:"divisionCode"`
	DivisionName               string                                `json:"divisionName"`
	DepartmentID               int                                   `json:"departmentId"`
	DepartmentCode             string                                `json:"departmentCode"`
	DepartmentName             string                                `json:"departmentName"`
	MaxPoints                  float64                               `json:"maxPoints"`
	FinalScore                 float64                               `json:"finalScore"`
	FinalScorePercentage       float64                               `json:"finalScorePercentage"`
	ReviewPeriodID             string                                `json:"reviewPeriodId"`
	RecordStatusName           string                                `json:"recordStatusName"`
	Ratings                    []CompetencyReviewerRatingSummaryData `json:"ratings"`
}

// CompetencyReviewerRatingSummaryData summarises ratings for a competency.
type CompetencyReviewerRatingSummaryData struct {
	PmsCompetencyID string  `json:"pmsCompetencyId"`
	PmsCompetency   string  `json:"pmsCompetency"`
	AverageRating   float64 `json:"averageRating"`
}

// CompetencyReviewFeedbackData is the entity-like DTO for competency review feedback.
type CompetencyReviewFeedbackData struct {
	BaseAuditVm
	CompetencyReviewFeedbackID string                  `json:"competencyReviewFeedbackId"`
	StaffID                    string                  `json:"staffId"`
	StaffName                  string                  `json:"staffName"`
	OfficeID                   int                     `json:"officeID"`
	OfficeCode                 string                  `json:"officeCode"`
	OfficeName                 string                  `json:"officeName"`
	DivisionID                 int                     `json:"divisionId"`
	DivisionCode               string                  `json:"divisionCode"`
	DivisionName               string                  `json:"divisionName"`
	DepartmentID               int                     `json:"departmentId"`
	DepartmentCode             string                  `json:"departmentCode"`
	DepartmentName             string                  `json:"departmentName"`
	MaxPoints                  float64                 `json:"maxPoints"`
	FinalScore                 float64                 `json:"finalScore"`
	FinalScorePercentage       float64                 `json:"finalScorePercentage"`
	ReviewPeriodID             string                  `json:"reviewPeriodId"`
	RecordStatusName           string                  `json:"recordStatusName"`
	CompetencyReviewers        []CompetencyReviewerData `json:"competencyReviewers"`
}

// CompetencyReviewFeedbackListResponseVm wraps a list of competency review feedbacks.
type CompetencyReviewFeedbackListResponseVm struct {
	GenericListResponseVm
	CompetencyReviewFeedbacks []CompetencyReviewFeedbackData `json:"competencyReviewFeedbacks"`
}

// ---------------------------------------------------------------------------
// 360 Review – Competency Reviewer DTOs  (source: ResponseViewModels.cs)
// ---------------------------------------------------------------------------

// CompetencyReviewerData is the entity-like DTO for a competency reviewer.
type CompetencyReviewerData struct {
	BaseAuditVm
	CompetencyReviewerID       string                          `json:"competencyReviewerId"`
	ReviewStaffID              string                          `json:"reviewStaffId"`
	FinalRating                float64                         `json:"finalRating"`
	CompetencyReviewFeedbackID string                          `json:"competencyReviewFeedbackId"`
	CompetencyReviewFeedback   *CompetencyReviewFeedbackData   `json:"competencyReviewFeedback"`
	RecordStatusName           string                          `json:"recordStatusName"`
	InitiatedDate              time.Time                       `json:"intitatedDate"`
	CompetencyReviewerRatings  []CompetencyReviewerRatingData  `json:"competencyReviewerRatings"`
}

// CompetencyReviewersListResponseVm wraps a list of competency reviewers.
type CompetencyReviewersListResponseVm struct {
	GenericListResponseVm
	CompetencyReviewers []CompetencyReviewerData `json:"competencyReviewers"`
}

// CompetencyReviewersResponseVm wraps a single competency reviewer.
type CompetencyReviewersResponseVm struct {
	BaseAPIResponse
	CompetencyReview *CompetencyReviewerData `json:"competencyReview"`
}

// ---------------------------------------------------------------------------
// 360 Review – Competency Reviewer Rating DTOs  (source: ResponseViewModels.cs)
// ---------------------------------------------------------------------------

// CompetencyReviewerRatingData is the entity-like DTO for a competency reviewer rating.
type CompetencyReviewerRatingData struct {
	BaseAuditVm
	CompetencyReviewerRatingID   string                         `json:"competencyReviewerRatingId"`
	PmsCompetencyID              string                         `json:"pmsCompetencyId"`
	FeedbackQuestionaireOptionID string                         `json:"feedbackQuestionaireOptionId"`
	Rating                       float64                        `json:"rating"`
	CompetencyReviewerID         string                         `json:"competencyReviewerId"`
	FeedbackQuestionaireOption   *FeedbackQuestionaireOptionData `json:"feedbackQuestionaireOption"`
	PmsCompetency                *PmsCompetencyData             `json:"pmsCompetency"`
}

// CompetencyRatingResponseVm wraps a single competency rating.
type CompetencyRatingResponseVm struct {
	BaseAPIResponse
	CompetencyReviewerRating *CompetencyReviewerRatingData `json:"competencyReviewerRating"`
}

// CompetencyRatingListResponseVm wraps a list of competency ratings.
type CompetencyRatingListResponseVm struct {
	GenericListResponseVm
	CompetencyReviewerRatings []CompetencyReviewerRatingData `json:"competencyReviewerRatings"`
}

// ---------------------------------------------------------------------------
// 360 Review – Questionnaire data DTOs  (source: ResponseViewModels.cs)
// ---------------------------------------------------------------------------

// PmsCompetencyData is the DTO for a PMS competency with its questionnaires.
type PmsCompetencyData struct {
	BaseAuditVm
	PmsCompetencyID  string                    `json:"pmsCompetencyId"`
	Name             string                    `json:"name"`
	Description      string                    `json:"description"`
	ObjectCategoryID string                    `json:"objectCategoryId"`
	FeedbackQuestionaires []FeedbackQuestionaireData `json:"feedbackQuestionaires"`
}

// FeedbackQuestionaireData is the DTO for a questionnaire under a competency.
type FeedbackQuestionaireData struct {
	BaseAuditVm
	FeedbackQuestionaireID string                          `json:"feedbackQuestionaireId"`
	Question               string                          `json:"question"`
	Description            string                          `json:"description"`
	PmsCompetencyID        string                          `json:"pmsCompetencyId"`
	Options                []FeedbackQuestionaireOptionData `json:"options"`
}

// FeedbackQuestionaireOptionData is the DTO for a questionnaire option.
type FeedbackQuestionaireOptionData struct {
	BaseAuditVm
	FeedbackQuestionaireOptionID string  `json:"feedbackQuestionaireOptionId"`
	OptionStatement              string  `json:"optionStatement"`
	Description                  string  `json:"description"`
	Score                        float64 `json:"score"`
	QuestionID                   string  `json:"questionId"`
}

// QuestionnaireListResponseVm wraps a list of PMS competencies with their questionnaires.
type QuestionnaireListResponseVm struct {
	GenericListResponseVm
	StaffID           string             `json:"staffId"`
	PmsCompetencyData []PmsCompetencyData `json:"pmsCompetencyData"`
}

// PmsCompetencyListResponseVm wraps a list of PMS competency setup records.
type PmsCompetencyListResponseVm struct {
	GenericListResponseVm
	PmsCompetencies []PmsCompetencyVm `json:"pmsCompetencies"`
}

// PmsCompetencyVm is the setup DTO for a PMS competency.
type PmsCompetencyVm struct {
	BaseAuditVm
	PmsCompetencyID  string `json:"pmsCompetencyId"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ObjectCategoryID string `json:"objectCategoryId"`
}

// ---------------------------------------------------------------------------
// 360 Review – Custom frontend DTOs  (source: FrontendCustomVms/ResponseFormVms.cs)
// ---------------------------------------------------------------------------

// CustomFeedbackQuestionaireOptionVm is a slim DTO for a questionnaire option.
type CustomFeedbackQuestionaireOptionVm struct {
	FeedbackQuestionaireOptionID string `json:"feedbackQuestionaireOptionId"`
	OptionStatement              string `json:"optionStatement"`
	Score                        int    `json:"score"`
	QuestionID                   string `json:"questionId"`
}

// CustomFeedbackQuestionaireVm is a slim DTO for a questionnaire.
type CustomFeedbackQuestionaireVm struct {
	FeedbackQuestionaireID string                               `json:"feedbackQuestionaireId"`
	Question               string                               `json:"question"`
	Description            string                               `json:"description"`
	RecordStatus           int                                  `json:"recordStatus"`
	PmsCompetencyID        string                               `json:"pmsCompetencyId"`
	Options                []CustomFeedbackQuestionaireOptionVm `json:"options"`
}

// CustomPmsCompetencyVm is a slim DTO for a PMS competency with questionnaires.
type CustomPmsCompetencyVm struct {
	PmsCompetencyID  string                         `json:"pmsCompetencyId"`
	Name             string                         `json:"name"`
	Description      string                         `json:"description"`
	ObjectCategoryID string                         `json:"objectCategoryId"`
	FeedbackQuestionaires []CustomFeedbackQuestionaireVm `json:"feedbackQuestionaires"`
}

// CompetencyReviewerRatingVm is the frontend DTO for a reviewer rating.
type CompetencyReviewerRatingVm struct {
	CompetencyReviewerRatingID   string                         `json:"competencyReviewerRatingId"`
	PmsCompetencyID              string                         `json:"pmsCompetencyId"`
	FeedbackQuestionaireOptionID string                         `json:"feedbackQuestionaireOptionId"`
	Rating                       float64                        `json:"rating"`
	CompetencyReviewerID         string                         `json:"competencyReviewerId"`
	FeedbackQuestionaireOption   *FeedbackQuestionaireOptionVm  `json:"feedbackQuestionaireOption"`
	PmsCompetency                *PmsCompetencyVm               `json:"pmsCompetency"`
}

// CustomCompetencyReviewFeedbackVm is a slim DTO for competency review feedback.
type CustomCompetencyReviewFeedbackVm struct {
	CompetencyReviewFeedbackID string                      `json:"competencyReviewFeedbackId"`
	StaffID                    string                      `json:"staffId"`
	MaxPoints                  float64                     `json:"maxPoints"`
	StaffName                  string                      `json:"staffName"`
	FinalScore                 float64                     `json:"finalScore"`
	ReviewPeriodID             string                      `json:"reviewPeriodId"`
	RecordStatus               int                         `json:"recordStatus"`
	CompetencyReviewers        []CustomCompetencyReviewerVm `json:"competencyReviewers"`
}

// CustomCompetencyReviewerVm is a slim DTO for a competency reviewer.
type CustomCompetencyReviewerVm struct {
	CompetencyReviewerID       string                          `json:"competencyReviewerId"`
	ReviewStaffID              string                          `json:"reviewStaffId"`
	FinalRating                float64                         `json:"finalRating"`
	CompetencyReviewFeedbackID string                          `json:"competencyReviewFeedbackId"`
	CompetencyReviewFeedback   *CustomCompetencyReviewFeedbackVm `json:"competencyReviewFeedback"`
	RecordStatus               enums.Status                    `json:"recordStatus"`
	InitiatedDate              time.Time                       `json:"intitatedDate"`
	CompetencyReviewerRatings  []CompetencyReviewerRatingVm    `json:"competencyReviewerRatings"`
}

// GetCompetenciesToReviewVm wraps reviewers for competency review.
type GetCompetenciesToReviewVm struct {
	CompetencyReviewers []CustomCompetencyReviewerVm `json:"competencyReviewers"`
}

// GetCompetenciesReviewDetailVm wraps a single reviewer detail.
type GetCompetenciesReviewDetailVm struct {
	CompetencyReview *CustomCompetencyReviewerVm `json:"competencyReview"`
}

// GetPMSQuestionnairesVm wraps PMS competency questionnaire data.
type GetPMSQuestionnairesVm struct {
	PmsCompetencyData []CustomPmsCompetencyVm `json:"pmsCompetencyData"`
}

// GetAllMyReviewed360CompetenciesVm wraps reviewed 360 competency feedbacks.
type GetAllMyReviewed360CompetenciesVm struct {
	CompetencyReviewFeedbacks []CustomCompetencyReviewFeedbackVm `json:"competencyReviewFeedbacks"`
}

// SavePmsCompetencyRequestVm is the request DTO for saving a competency rating.
type SavePmsCompetencyRequestVm struct {
	PmsCompetencyID              string `json:"pmsCompetencyId"`
	FeedbackQuestionaireOptionID string `json:"feedbackQuestionaireOptionId"`
	CompetencyReviewerID         string `json:"competencyReviewerId"`
	CompetencyReviewerRatingID   string `json:"competencyReviewerRatingId"`
}

// CompleteCompetencyReviewVm is the request DTO for completing a competency review.
type CompleteCompetencyReviewVm struct {
	ReviewStaffID              string `json:"reviewStaffId"`
	CompetencyReviewFeedbackID string `json:"competencyReviewFeedbackId"`
}

// CompleteFeedbackVm is the request DTO for completing a feedback cycle.
type CompleteFeedbackVm struct {
	ReviewStaffID              string `json:"reviewStaffId"`
	CompetencyReviewFeedbackID string `json:"competencyReviewFeedbackId"`
	CompetencyReviewerID       string `json:"competencyReviewerId"`
}

// ---------------------------------------------------------------------------
// 360 Review – Request VMs  (source: RequestVms.cs)
// ---------------------------------------------------------------------------

// Trigger360ReviewRequestModel is the request DTO for triggering a 360 review cycle.
type Trigger360ReviewRequestModel struct {
	Target         enums.ReviewPeriodExtensionTargetType `json:"target" validate:"required"`
	Reference      string                                `json:"reference"`
	ReviewPeriodID string                                `json:"reviewPeriodId" validate:"required"`
}

// Initiate360ReviewRequestModel is the request DTO for initiating 360 reviews.
type Initiate360ReviewRequestModel struct {
	StaffID        []string `json:"staffId" validate:"required"`
	ReviewPeriodID string   `json:"reviewPeriodId" validate:"required"`
}

// Complete360ReviewRequestModel is the request DTO for completing 360 reviews.
type Complete360ReviewRequestModel struct {
	ReviewPeriodID string `json:"reviewPeriodId" validate:"required"`
}

// ---------------------------------------------------------------------------
// Audit Log DTOs  (source: AuditLog.cs model + ResponseViewModels.cs)
// ---------------------------------------------------------------------------

// AuditLogVm is the read/response DTO for an audit log entry.
type AuditLogVm struct {
	BaseAuditVm
	UserName            string    `json:"userName"`
	AuditEventDateUTC   time.Time `json:"auditEventDateUTC"`
	AuditEventType      int       `json:"auditEventType"`
	TableName           string    `json:"tableName"`
	RecordID            string    `json:"recordId"`
	FieldName           string    `json:"fieldName"`
	OriginalValue       string    `json:"originalValue"`
	NewValue            string    `json:"newValue"`
}

// AuditLogListResponseVm wraps a paginated list of audit logs.
type AuditLogListResponseVm struct {
	GenericListResponseVm
	AuditLogs []AuditLogVm `json:"auditLogs"`
}

// AuditLogResponseVm wraps a single audit log entry.
type AuditLogResponseVm struct {
	BaseAPIResponse
	AuditLog *AuditLogVm `json:"auditLog"`
}

// AuditLogSearchVm is the search/filter DTO for querying audit logs.
type AuditLogSearchVm struct {
	BasePagedData
	SearchString string     `json:"searchString"`
	UserName     string     `json:"userName"`
	TableName    string     `json:"tableName"`
	StartDate    *time.Time `json:"startDate"`
	EndDate      *time.Time `json:"endDate"`
}

// ---------------------------------------------------------------------------
// Generic List Response  (source: ResponseViewModels.cs – GenericListResponseVm)
// ---------------------------------------------------------------------------

// GenericListResponseVm extends BaseAPIResponse with a total record count.
type GenericListResponseVm struct {
	BaseAPIResponse
	TotalRecords int `json:"totalRecords"`
}
