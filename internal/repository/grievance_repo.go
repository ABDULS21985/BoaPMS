package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"gorm.io/gorm"
)

// GrievanceRepository provides data access for grievance management entities.
type GrievanceRepository struct {
	db *gorm.DB
}

// NewGrievanceRepository creates a new grievance repository.
func NewGrievanceRepository(db *gorm.DB) *GrievanceRepository {
	return &GrievanceRepository{db: db}
}

func (r *GrievanceRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Where("soft_deleted = ?", false)
}

// ─── Grievance ───────────────────────────────────────────────────────────────

func (r *GrievanceRepository) GetGrievanceByID(ctx context.Context, id string) (*performance.Grievance, error) {
	var g performance.Grievance
	err := r.base(ctx).
		Preload("GrievanceResolutions").
		First(&g, "grievance_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("grievanceRepo.GetGrievanceByID: %w", err)
	}
	return &g, nil
}

func (r *GrievanceRepository) GetGrievancesByComplainant(ctx context.Context, complainantStaffID string) ([]performance.Grievance, error) {
	var results []performance.Grievance
	err := r.base(ctx).
		Where("complainant_staff_id = ?", complainantStaffID).
		Preload("GrievanceResolutions").
		Order("created_at DESC").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("grievanceRepo.GetGrievancesByComplainant: %w", err)
	}
	return results, nil
}

func (r *GrievanceRepository) GetGrievancesByRespondent(ctx context.Context, respondentStaffID string) ([]performance.Grievance, error) {
	var results []performance.Grievance
	err := r.base(ctx).
		Where("respondent_staff_id = ?", respondentStaffID).
		Preload("GrievanceResolutions").
		Order("created_at DESC").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("grievanceRepo.GetGrievancesByRespondent: %w", err)
	}
	return results, nil
}

func (r *GrievanceRepository) GetGrievancesByMediator(ctx context.Context, mediatorStaffID string) ([]performance.Grievance, error) {
	var results []performance.Grievance
	err := r.base(ctx).
		Where("current_mediator_staff_id = ?", mediatorStaffID).
		Preload("GrievanceResolutions").
		Order("created_at DESC").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("grievanceRepo.GetGrievancesByMediator: %w", err)
	}
	return results, nil
}

func (r *GrievanceRepository) GetGrievancesByStatus(ctx context.Context, status enums.Status) ([]performance.Grievance, error) {
	q := r.base(ctx)
	if status != enums.StatusAll {
		q = q.Where("status = ?", status.String())
	}
	var results []performance.Grievance
	err := q.Preload("GrievanceResolutions").Order("created_at DESC").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("grievanceRepo.GetGrievancesByStatus: %w", err)
	}
	return results, nil
}

func (r *GrievanceRepository) GetGrievancesByReviewPeriod(ctx context.Context, reviewPeriodID string) ([]performance.Grievance, error) {
	var results []performance.Grievance
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Preload("GrievanceResolutions").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("grievanceRepo.GetGrievancesByReviewPeriod: %w", err)
	}
	return results, nil
}

// ─── GrievanceResolution ─────────────────────────────────────────────────────

func (r *GrievanceRepository) GetResolutionsByGrievance(ctx context.Context, grievanceID string) ([]performance.GrievanceResolution, error) {
	var results []performance.GrievanceResolution
	err := r.base(ctx).
		Where("grievance_id = ?", grievanceID).
		Order("created_at ASC").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("grievanceRepo.GetResolutionsByGrievance: %w", err)
	}
	return results, nil
}

// ─── Transaction helper ──────────────────────────────────────────────────────

func (r *GrievanceRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *GrievanceRepository) DB() *gorm.DB {
	return r.db
}
