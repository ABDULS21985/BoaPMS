package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/competency"
	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"gorm.io/gorm"
)

// CompetencyRepository provides data access for competency domain entities.
// Covers: Competency, CompetencyCategory, CompetencyCategoryGrading,
// CompetencyRatingDefinition, CompetencyReview, CompetencyReviewProfile,
// DevelopmentPlan, Rating, ReviewType, TrainingType, BankYear,
// JobRole, JobGrade, JobGradeGroup, AssignJobGradeGroup, OfficeJobRole,
// JobRoleCompetency, BehavioralCompetency, JobRoleGrade, StaffJobRoles, ReviewPeriod.
type CompetencyRepository struct {
	db *gorm.DB
}

// NewCompetencyRepository creates a new competency repository.
func NewCompetencyRepository(db *gorm.DB) *CompetencyRepository {
	return &CompetencyRepository{db: db}
}

// ─── Competency ──────────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetCompetencyByID(ctx context.Context, id int) (*competency.Competency, error) {
	var c competency.Competency
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ?", false).
		Preload("CompetencyCategory").
		Preload("CompetencyRatingDefinitions").
		First(&c, "competency_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetCompetencyByID: %w", err)
	}
	return &c, nil
}

func (r *CompetencyRepository) GetAllCompetencies(ctx context.Context) ([]competency.Competency, error) {
	var results []competency.Competency
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ?", false).
		Preload("CompetencyCategory").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllCompetencies: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) GetCompetenciesByCategoryID(ctx context.Context, categoryID int) ([]competency.Competency, error) {
	var results []competency.Competency
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND competency_category_id = ?", false, categoryID).
		Preload("CompetencyRatingDefinitions").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetCompetenciesByCategoryID: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) SearchCompetencies(ctx context.Context, categoryID *int, search string, isApproved *bool, isTechnical *bool, offset, limit int) ([]competency.Competency, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&competency.Competency{}).
		Where("soft_deleted = ?", false).
		Preload("CompetencyCategory")

	if categoryID != nil {
		q = q.Where("competency_category_id = ?", *categoryID)
	}
	if search != "" {
		q = q.Where("competency_name ILIKE ?", "%"+search+"%")
	}
	if isApproved != nil {
		q = q.Where("is_approved = ?", *isApproved)
	}
	if isTechnical != nil {
		q = q.Joins("JOIN \"CoreSchema\".competency_categories cc ON cc.competency_category_id = competencies.competency_category_id").
			Where("cc.is_technical = ?", *isTechnical)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("competencyRepo.SearchCompetencies count: %w", err)
	}

	var results []competency.Competency
	if err := q.Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("competencyRepo.SearchCompetencies: %w", err)
	}
	return results, total, nil
}

// ─── CompetencyCategory ─────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllCategories(ctx context.Context) ([]competency.CompetencyCategory, error) {
	var results []competency.CompetencyCategory
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllCategories: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) GetCategoryByID(ctx context.Context, id int) (*competency.CompetencyCategory, error) {
	var c competency.CompetencyCategory
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).First(&c, "competency_category_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetCategoryByID: %w", err)
	}
	return &c, nil
}

// ─── CompetencyCategoryGrading ───────────────────────────────────────────────

func (r *CompetencyRepository) GetCategoryGradings(ctx context.Context) ([]competency.CompetencyCategoryGrading, error) {
	var results []competency.CompetencyCategoryGrading
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ?", false).
		Preload("CompetencyCategory").
		Preload("ReviewType").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetCategoryGradings: %w", err)
	}
	return results, nil
}

// ─── CompetencyRatingDefinition ──────────────────────────────────────────────

func (r *CompetencyRepository) GetRatingDefinitionsByCompetencyID(ctx context.Context, competencyID int) ([]competency.CompetencyRatingDefinition, error) {
	var results []competency.CompetencyRatingDefinition
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND competency_id = ?", false, competencyID).
		Preload("Rating").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetRatingDefinitionsByCompetencyID: %w", err)
	}
	return results, nil
}

// ─── CompetencyReview ────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetReviewsByEmployeeAndPeriod(ctx context.Context, employeeNumber string, reviewPeriodID int) ([]competency.CompetencyReview, error) {
	var results []competency.CompetencyReview
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND employee_number = ? AND review_period_id = ?", false, employeeNumber, reviewPeriodID).
		Preload("Competency").
		Preload("Competency.CompetencyCategory").
		Preload("ExpectedRating").
		Preload("ActualRating").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetReviewsByEmployeeAndPeriod: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) GetReviewsByReviewer(ctx context.Context, reviewerID string, reviewPeriodID int) ([]competency.CompetencyReview, error) {
	var results []competency.CompetencyReview
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND reviewer_id = ? AND review_period_id = ?", false, reviewerID, reviewPeriodID).
		Preload("Competency").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetReviewsByReviewer: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) GetPendingReviewerIDs(ctx context.Context, reviewPeriodID int) ([]string, error) {
	var ids []string
	err := r.db.WithContext(ctx).
		Model(&competency.CompetencyReview{}).
		Where("soft_deleted = ? AND actual_rating_value = 0 AND review_period_id = ?", false, reviewPeriodID).
		Distinct("reviewer_id").
		Pluck("reviewer_id", &ids).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetPendingReviewerIDs: %w", err)
	}
	return ids, nil
}

// ─── CompetencyReviewProfile ─────────────────────────────────────────────────

func (r *CompetencyRepository) GetReviewProfilesByEmployee(ctx context.Context, employeeNumber string, reviewPeriodID int) ([]competency.CompetencyReviewProfile, error) {
	var results []competency.CompetencyReviewProfile
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND employee_number = ? AND review_period_id = ?", false, employeeNumber, reviewPeriodID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetReviewProfilesByEmployee: %w", err)
	}
	return results, nil
}

// ─── DevelopmentPlan ─────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetDevelopmentPlansByProfile(ctx context.Context, profileID int) ([]competency.DevelopmentPlan, error) {
	var results []competency.DevelopmentPlan
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND competency_review_profile_id = ?", false, profileID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetDevelopmentPlansByProfile: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) GetDevelopmentPlansByEmployee(ctx context.Context, employeeNumber string) ([]competency.DevelopmentPlan, error) {
	var results []competency.DevelopmentPlan
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND employee_number = ?", false, employeeNumber).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetDevelopmentPlansByEmployee: %w", err)
	}
	return results, nil
}

// ─── Rating ──────────────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllRatings(ctx context.Context) ([]competency.Rating, error) {
	var results []competency.Rating
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Order("value ASC").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllRatings: %w", err)
	}
	return results, nil
}

// ─── ReviewType ──────────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllReviewTypes(ctx context.Context) ([]competency.ReviewType, error) {
	var results []competency.ReviewType
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllReviewTypes: %w", err)
	}
	return results, nil
}

// ─── TrainingType ────────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllTrainingTypes(ctx context.Context) ([]competency.TrainingType, error) {
	var results []competency.TrainingType
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllTrainingTypes: %w", err)
	}
	return results, nil
}

// ─── JobRole ─────────────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllJobRoles(ctx context.Context) ([]competency.JobRole, error) {
	var results []competency.JobRole
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllJobRoles: %w", err)
	}
	return results, nil
}

// ─── JobGrade ────────────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllJobGrades(ctx context.Context) ([]competency.JobGrade, error) {
	var results []competency.JobGrade
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllJobGrades: %w", err)
	}
	return results, nil
}

// ─── JobGradeGroup ───────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllJobGradeGroups(ctx context.Context) ([]competency.JobGradeGroup, error) {
	var results []competency.JobGradeGroup
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Order(`"order" ASC`).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllJobGradeGroups: %w", err)
	}
	return results, nil
}

// ─── AssignJobGradeGroup ─────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAssignmentsByGradeGroupID(ctx context.Context, gradeGroupID int) ([]competency.AssignJobGradeGroup, error) {
	var results []competency.AssignJobGradeGroup
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND job_grade_group_id = ?", false, gradeGroupID).
		Preload("JobGrade").
		Preload("JobGradeGroup").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAssignmentsByGradeGroupID: %w", err)
	}
	return results, nil
}

// ─── OfficeJobRole ───────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetOfficeJobRolesByOffice(ctx context.Context, officeID int) ([]competency.OfficeJobRole, error) {
	var results []competency.OfficeJobRole
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND office_id = ?", false, officeID).
		Preload("JobRole").
		Preload("Office").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetOfficeJobRolesByOffice: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) SearchOfficeJobRoles(ctx context.Context, officeID *int, search string, offset, limit int) ([]competency.OfficeJobRole, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&competency.OfficeJobRole{}).
		Where("soft_deleted = ?", false).
		Preload("JobRole").
		Preload("Office")

	if officeID != nil {
		q = q.Where("office_id = ?", *officeID)
	}
	if search != "" {
		q = q.Joins("JOIN \"CoreSchema\".job_roles jr ON jr.job_role_id = office_job_roles.job_role_id").
			Where("jr.job_role_name ILIKE ?", "%"+search+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("competencyRepo.SearchOfficeJobRoles count: %w", err)
	}
	var results []competency.OfficeJobRole
	if err := q.Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("competencyRepo.SearchOfficeJobRoles: %w", err)
	}
	return results, total, nil
}

// ─── JobRoleCompetency ───────────────────────────────────────────────────────

func (r *CompetencyRepository) GetJobRoleCompetencies(ctx context.Context, officeID, jobRoleID int) ([]competency.JobRoleCompetency, error) {
	var results []competency.JobRoleCompetency
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND office_id = ? AND job_role_id = ?", false, officeID, jobRoleID).
		Preload("Competency").
		Preload("Rating").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetJobRoleCompetencies: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) SearchJobRoleCompetencies(ctx context.Context, officeID, jobRoleID *int, search string, offset, limit int) ([]competency.JobRoleCompetency, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&competency.JobRoleCompetency{}).
		Where("soft_deleted = ?", false).
		Preload("Competency").
		Preload("Rating")

	if officeID != nil {
		q = q.Where("office_id = ?", *officeID)
	}
	if jobRoleID != nil {
		q = q.Where("job_role_id = ?", *jobRoleID)
	}
	if search != "" {
		q = q.Joins("JOIN \"CoreSchema\".competencies c ON c.competency_id = job_role_competencies.competency_id").
			Where("c.competency_name ILIKE ?", "%"+search+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("competencyRepo.SearchJobRoleCompetencies count: %w", err)
	}
	var results []competency.JobRoleCompetency
	if err := q.Offset(offset).Limit(limit).Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("competencyRepo.SearchJobRoleCompetencies: %w", err)
	}
	return results, total, nil
}

// ─── BehavioralCompetency ────────────────────────────────────────────────────

func (r *CompetencyRepository) GetBehavioralCompetencies(ctx context.Context, gradeGroupID int) ([]competency.BehavioralCompetency, error) {
	var results []competency.BehavioralCompetency
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND job_grade_group_id = ?", false, gradeGroupID).
		Preload("Competency").
		Preload("Rating").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetBehavioralCompetencies: %w", err)
	}
	return results, nil
}

// ─── StaffJobRoles ───────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetStaffJobRole(ctx context.Context, employeeID string) (*competency.StaffJobRoles, error) {
	var result competency.StaffJobRoles
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND employee_id = ?", false, employeeID).
		First(&result).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetStaffJobRole: %w", err)
	}
	return &result, nil
}

func (r *CompetencyRepository) GetStaffJobRoleRequests(ctx context.Context, supervisorID string) ([]competency.StaffJobRoles, error) {
	var results []competency.StaffJobRoles
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND supervisor_id = ?", false, supervisorID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetStaffJobRoleRequests: %w", err)
	}
	return results, nil
}

func (r *CompetencyRepository) UpsertStaffJobRole(ctx context.Context, sjr *competency.StaffJobRoles) error {
	return r.db.WithContext(ctx).Save(sjr).Error
}

func (r *CompetencyRepository) UpdateStaffJobRoleFields(ctx context.Context, employeeID string, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&competency.StaffJobRoles{}).
		Where("employee_id = ?", employeeID).
		Updates(updates).Error
}

// ─── ReviewPeriod (Competency) ───────────────────────────────────────────────

func (r *CompetencyRepository) GetActiveReviewPeriod(ctx context.Context) (*competency.ReviewPeriod, error) {
	var rp competency.ReviewPeriod
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ? AND is_active = ? AND is_approved = ?", false, true, true).
		First(&rp).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetActiveReviewPeriod: %w", err)
	}
	return &rp, nil
}

func (r *CompetencyRepository) GetAllReviewPeriods(ctx context.Context) ([]competency.ReviewPeriod, error) {
	var results []competency.ReviewPeriod
	err := r.db.WithContext(ctx).
		Where("soft_deleted = ?", false).
		Preload("BankYear").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllReviewPeriods: %w", err)
	}
	return results, nil
}

// ─── BankYear ────────────────────────────────────────────────────────────────

func (r *CompetencyRepository) GetAllBankYears(ctx context.Context) ([]identity.BankYear, error) {
	var results []identity.BankYear
	err := r.db.WithContext(ctx).Where("soft_deleted = ?", false).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("competencyRepo.GetAllBankYears: %w", err)
	}
	return results, nil
}

// ─── Transaction helper ──────────────────────────────────────────────────────

func (r *CompetencyRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// DB exposes the underlying GORM instance for advanced queries.
func (r *CompetencyRepository) DB() *gorm.DB {
	return r.db
}
