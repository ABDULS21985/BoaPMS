package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"gorm.io/gorm"
)

// FeedbackRepository provides data access for feedback, 360 reviews,
// and competency review feedback entities.
type FeedbackRepository struct {
	db *gorm.DB
}

// NewFeedbackRepository creates a new feedback repository.
func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

func (r *FeedbackRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Where("soft_deleted = ?", false)
}

// ─── FeedbackRequestLog ──────────────────────────────────────────────────────

func (r *FeedbackRepository) GetFeedbackRequestByID(ctx context.Context, id string) (*performance.FeedbackRequestLog, error) {
	var req performance.FeedbackRequestLog
	err := r.base(ctx).First(&req, "feedback_request_log_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetFeedbackRequestByID: %w", err)
	}
	return &req, nil
}

func (r *FeedbackRepository) GetFeedbackRequestsByStaff(ctx context.Context, staffID string, reviewPeriodID string) ([]performance.FeedbackRequestLog, error) {
	var results []performance.FeedbackRequestLog
	q := r.base(ctx).Where("assigned_staff_id = ?", staffID)
	if reviewPeriodID != "" {
		q = q.Where("review_period_id = ?", reviewPeriodID)
	}
	err := q.Order("time_initiated DESC").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetFeedbackRequestsByStaff: %w", err)
	}
	return results, nil
}

func (r *FeedbackRepository) GetFeedbackRequestsByOwner(ctx context.Context, ownerStaffID string, reviewPeriodID string) ([]performance.FeedbackRequestLog, error) {
	var results []performance.FeedbackRequestLog
	q := r.base(ctx).Where("request_owner_staff_id = ?", ownerStaffID)
	if reviewPeriodID != "" {
		q = q.Where("review_period_id = ?", reviewPeriodID)
	}
	err := q.Order("time_initiated DESC").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetFeedbackRequestsByOwner: %w", err)
	}
	return results, nil
}

func (r *FeedbackRepository) GetFeedbackRequestsByType(ctx context.Context, reqType enums.FeedbackRequestType, referenceID string) ([]performance.FeedbackRequestLog, error) {
	var results []performance.FeedbackRequestLog
	err := r.base(ctx).
		Where("feedback_request_type = ? AND reference_id = ?", reqType, referenceID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetFeedbackRequestsByType: %w", err)
	}
	return results, nil
}

func (r *FeedbackRepository) GetBreachedFeedbackRequests(ctx context.Context, reviewPeriodID string) ([]performance.FeedbackRequestLog, error) {
	var results []performance.FeedbackRequestLog
	err := r.base(ctx).
		Where("status = ? AND review_period_id = ?", enums.StatusBreached.String(), reviewPeriodID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetBreachedFeedbackRequests: %w", err)
	}
	return results, nil
}

// ─── FeedbackQuestionaire ────────────────────────────────────────────────────

func (r *FeedbackRepository) GetAllQuestionaires(ctx context.Context) ([]performance.FeedbackQuestionaire, error) {
	var results []performance.FeedbackQuestionaire
	err := r.base(ctx).
		Preload("Options").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetAllQuestionaires: %w", err)
	}
	return results, nil
}

func (r *FeedbackRepository) GetQuestionairesByCompetency(ctx context.Context, pmsCompetencyID string) ([]performance.FeedbackQuestionaire, error) {
	var results []performance.FeedbackQuestionaire
	err := r.base(ctx).
		Where("pms_competency_id = ?", pmsCompetencyID).
		Preload("Options").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetQuestionairesByCompetency: %w", err)
	}
	return results, nil
}

// ─── PmsCompetency ───────────────────────────────────────────────────────────

func (r *FeedbackRepository) GetAllPmsCompetencies(ctx context.Context) ([]performance.PmsCompetency, error) {
	var results []performance.PmsCompetency
	err := r.base(ctx).
		Preload("FeedbackQuestionaires").
		Preload("FeedbackQuestionaires.Options").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetAllPmsCompetencies: %w", err)
	}
	return results, nil
}

func (r *FeedbackRepository) GetPmsCompetencyByID(ctx context.Context, id string) (*performance.PmsCompetency, error) {
	var c performance.PmsCompetency
	err := r.base(ctx).
		Preload("FeedbackQuestionaires").
		Preload("FeedbackQuestionaires.Options").
		Preload("CompetencyReviewerRatings").
		First(&c, "pms_competency_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetPmsCompetencyByID: %w", err)
	}
	return &c, nil
}

// ─── CompetencyReviewFeedback ────────────────────────────────────────────────

func (r *FeedbackRepository) GetCompetencyReviewFeedbackByStaff(ctx context.Context, staffID, reviewPeriodID string) (*performance.CompetencyReviewFeedback, error) {
	var crf performance.CompetencyReviewFeedback
	err := r.base(ctx).
		Where("staff_id = ? AND review_period_id = ?", staffID, reviewPeriodID).
		Preload("CompetencyReviewers").
		Preload("CompetencyReviewers.CompetencyReviewerRatings").
		First(&crf).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetCompetencyReviewFeedbackByStaff: %w", err)
	}
	return &crf, nil
}

func (r *FeedbackRepository) GetAllCompetencyReviewFeedbacks(ctx context.Context, reviewPeriodID string) ([]performance.CompetencyReviewFeedback, error) {
	var results []performance.CompetencyReviewFeedback
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Preload("CompetencyReviewers").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetAllCompetencyReviewFeedbacks: %w", err)
	}
	return results, nil
}

// ─── CompetencyGapClosure ────────────────────────────────────────────────────

func (r *FeedbackRepository) GetCompetencyGapClosures(ctx context.Context, staffID, reviewPeriodID string) ([]performance.CompetencyGapClosure, error) {
	var results []performance.CompetencyGapClosure
	err := r.base(ctx).
		Where("staff_id = ? AND review_period_id = ?", staffID, reviewPeriodID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("feedbackRepo.GetCompetencyGapClosures: %w", err)
	}
	return results, nil
}

// ─── Transaction helper ──────────────────────────────────────────────────────

func (r *FeedbackRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *FeedbackRepository) DB() *gorm.DB {
	return r.db
}
