package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// FeedbackRequestLog tracks feedback/approval requests between staff.
type FeedbackRequestLog struct {
	FeedbackRequestLogID string                    `json:"feedback_request_log_id" gorm:"column:feedback_request_log_id;primaryKey"`
	FeedbackRequestType  enums.FeedbackRequestType `json:"feedback_request_type"   gorm:"column:feedback_request_type;not null"`
	ReferenceID          string                    `json:"reference_id"            gorm:"column:reference_id;not null;index"`
	TimeInitiated        time.Time                 `json:"time_initiated"          gorm:"column:time_initiated;not null"`
	AssignedStaffID      string                    `json:"assigned_staff_id"       gorm:"column:assigned_staff_id;not null;index"`
	AssignedStaffName    string                    `json:"assigned_staff_name"     gorm:"column:assigned_staff_name"`
	RequestOwnerStaffID  string                    `json:"request_owner_staff_id"  gorm:"column:request_owner_staff_id;not null"`
	RequestOwnerStaffName string                   `json:"request_owner_staff_name" gorm:"column:request_owner_staff_name"`
	TimeCompleted        *time.Time                `json:"time_completed"          gorm:"column:time_completed"`
	RequestOwnerComment  string                    `json:"request_owner_comment"   gorm:"column:request_owner_comment"`
	RequestOwnerAttachment string                  `json:"request_owner_attachment" gorm:"column:request_owner_attachment"`
	AssignedStaffComment string                    `json:"assigned_staff_comment"  gorm:"column:assigned_staff_comment"`
	AssignedStaffAttachment string                 `json:"assigned_staff_attachment" gorm:"column:assigned_staff_attachment"`
	HasSLA               bool                      `json:"has_sla"                 gorm:"column:has_sla;default:true"`
	ReviewPeriodID       string                    `json:"review_period_id"        gorm:"column:review_period_id"`
	domain.BaseEntity

	ReviewPeriod *PerformanceReviewPeriod `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
}

func (FeedbackRequestLog) TableName() string { return "pms.feedback_request_logs" }

// FeedbackQuestionaire defines a 360-feedback question for a PMS competency.
type FeedbackQuestionaire struct {
	FeedbackQuestionaireID string `json:"feedback_questionaire_id" gorm:"column:feedback_questionaire_id;primaryKey"`
	Question               string `json:"question"                 gorm:"column:question;not null"`
	Description            string `json:"description"              gorm:"column:description;not null"`
	PmsCompetencyID        string `json:"pms_competency_id"        gorm:"column:pms_competency_id;not null"`
	domain.BaseWorkFlow

	PmsCompetency *PmsCompetency              `json:"pms_competency" gorm:"foreignKey:PmsCompetencyID"`
	Options       []FeedbackQuestionaireOption `json:"options"        gorm:"foreignKey:QuestionID"`
}

func (FeedbackQuestionaire) TableName() string { return "pms.feedback_questionaires" }

// FeedbackQuestionaireOption is a selectable answer for a feedback question.
type FeedbackQuestionaireOption struct {
	FeedbackQuestionaireOptionID string  `json:"feedback_questionaire_option_id" gorm:"column:feedback_questionaire_option_id;primaryKey"`
	OptionStatement              string  `json:"option_statement"                gorm:"column:option_statement;not null"`
	Description                  string  `json:"description"                     gorm:"column:description"`
	Score                        float64 `json:"score"                           gorm:"column:score;type:decimal(18,2);not null"`
	QuestionID                   string  `json:"question_id"                     gorm:"column:question_id;not null"`
	domain.BaseEntity

	Question *FeedbackQuestionaire `json:"question" gorm:"foreignKey:QuestionID"`
}

func (FeedbackQuestionaireOption) TableName() string { return "pms.feedback_questionaire_options" }

// PmsCompetency is a competency dimension used in 360 feedback.
type PmsCompetency struct {
	PmsCompetencyID string `json:"pms_competency_id" gorm:"column:pms_competency_id;primaryKey"`
	Name            string `json:"name"              gorm:"column:name;not null"`
	Description     string `json:"description"       gorm:"column:description"`
	ObjectCategoryID string `json:"object_category_id" gorm:"column:object_category_id;not null"`
	domain.BaseWorkFlow

	ObjectiveCategory      *ObjectiveCategory        `json:"objective_category"       gorm:"foreignKey:ObjectCategoryID"`
	CompetencyReviewerRatings []CompetencyReviewerRating `json:"competency_reviewer_ratings" gorm:"foreignKey:PmsCompetencyID"`
	FeedbackQuestionaires  []FeedbackQuestionaire    `json:"feedback_questionaires"   gorm:"foreignKey:PmsCompetencyID"`
}

func (PmsCompetency) TableName() string { return "pms.pms_competencies" }

// CompetencyReviewFeedback aggregates 360 feedback for a staff member.
type CompetencyReviewFeedback struct {
	CompetencyReviewFeedbackID string  `json:"competency_review_feedback_id" gorm:"column:competency_review_feedback_id;primaryKey"`
	StaffID                    string  `json:"staff_id"                      gorm:"column:staff_id;not null;index"`
	MaxPoints                  float64 `json:"max_points"                    gorm:"column:max_points;type:decimal(18,2)"`
	FinalScore                 float64 `json:"final_score"                   gorm:"column:final_score;type:decimal(18,2)"`
	ReviewPeriodID             string  `json:"review_period_id"              gorm:"column:review_period_id;not null"`
	domain.BaseEntity

	ReviewPeriod        *PerformanceReviewPeriod `json:"review_period"        gorm:"foreignKey:ReviewPeriodID"`
	CompetencyReviewers []CompetencyReviewer     `json:"competency_reviewers" gorm:"foreignKey:CompetencyReviewFeedbackID"`
}

func (CompetencyReviewFeedback) TableName() string { return "pms.competency_review_feedbacks" }

// CompetencyReviewer is a single reviewer in a 360-feedback cycle.
type CompetencyReviewer struct {
	CompetencyReviewerID       string  `json:"competency_reviewer_id"        gorm:"column:competency_reviewer_id;primaryKey"`
	ReviewStaffID              string  `json:"review_staff_id"               gorm:"column:review_staff_id;not null"`
	FinalRating                float64 `json:"final_rating"                  gorm:"column:final_rating;type:decimal(18,2)"`
	CompetencyReviewFeedbackID string  `json:"competency_review_feedback_id" gorm:"column:competency_review_feedback_id;not null"`
	domain.BaseEntity

	CompetencyReviewFeedback *CompetencyReviewFeedback  `json:"competency_review_feedback" gorm:"foreignKey:CompetencyReviewFeedbackID"`
	CompetencyReviewerRatings []CompetencyReviewerRating `json:"competency_reviewer_ratings" gorm:"foreignKey:CompetencyReviewerID"`
}

func (CompetencyReviewer) TableName() string { return "pms.competency_reviewers" }

// CompetencyReviewerRating stores one rating from a reviewer for a competency.
type CompetencyReviewerRating struct {
	CompetencyReviewerRatingID   string  `json:"competency_reviewer_rating_id"    gorm:"column:competency_reviewer_rating_id;primaryKey"`
	PmsCompetencyID              string  `json:"pms_competency_id"                gorm:"column:pms_competency_id"`
	FeedbackQuestionaireOptionID string  `json:"feedback_questionaire_option_id"   gorm:"column:feedback_questionaire_option_id"`
	Rating                       float64 `json:"rating"                           gorm:"column:rating;type:decimal(18,2)"`
	CompetencyReviewerID         string  `json:"competency_reviewer_id"           gorm:"column:competency_reviewer_id;not null"`
	domain.BaseEntity

	FeedbackQuestionaireOption *FeedbackQuestionaireOption `json:"feedback_questionaire_option" gorm:"foreignKey:FeedbackQuestionaireOptionID"`
	PmsCompetency              *PmsCompetency              `json:"pms_competency"               gorm:"foreignKey:PmsCompetencyID"`
	CompetencyReviewer         *CompetencyReviewer         `json:"competency_reviewer"          gorm:"foreignKey:CompetencyReviewerID"`
}

func (CompetencyReviewerRating) TableName() string { return "pms.competency_reviewer_ratings" }

// CompetencyGapClosure tracks competency gap closure progress.
type CompetencyGapClosure struct {
	CompetencyGapClosureID string  `json:"competency_gap_closure_id" gorm:"column:competency_gap_closure_id;primaryKey"`
	StaffID                string  `json:"staff_id"                  gorm:"column:staff_id;not null;index"`
	MaxPoints              float64 `json:"max_points"                gorm:"column:max_points;type:decimal(18,2)"`
	FinalScore             float64 `json:"final_score"               gorm:"column:final_score;type:decimal(18,2)"`
	ReviewPeriodID         string  `json:"review_period_id"          gorm:"column:review_period_id;not null"`
	ObjectiveCategoryID    string  `json:"objective_category_id"     gorm:"column:objective_category_id;not null"`
	domain.BaseEntity

	ReviewPeriod      *PerformanceReviewPeriod `json:"review_period"      gorm:"foreignKey:ReviewPeriodID"`
	ObjectiveCategory *ObjectiveCategory       `json:"objective_category" gorm:"foreignKey:ObjectiveCategoryID"`
}

func (CompetencyGapClosure) TableName() string { return "pms.competency_gap_closures" }
