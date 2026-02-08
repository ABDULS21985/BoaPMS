package repository

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/audit"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"gorm.io/gorm"
)

// PerformanceRepository provides data access for performance management entities.
// Covers: Strategy, StrategicTheme, PerformanceReviewPeriod, ObjectiveCategory,
// CategoryDefinition, EnterpriseObjective, DepartmentObjective, DivisionObjective,
// OfficeObjective, PeriodObjective, PeriodObjectiveEvaluation,
// PeriodObjectiveDepartmentEvaluation, ReviewPeriodIndividualPlannedObjective,
// WorkProduct, WorkProductTask, WorkProductEvaluation, EvaluationOption,
// OperationalObjectiveWorkProduct, PeriodScore, ReviewPeriodExtension,
// ReviewPeriod360Review, PmsConfiguration, Setting, SequenceNumber,
// WorkProductDefinition, CascadedWorkProduct.
type PerformanceRepository struct {
	db *gorm.DB
}

// NewPerformanceRepository creates a new performance management repository.
func NewPerformanceRepository(db *gorm.DB) *PerformanceRepository {
	return &PerformanceRepository{db: db}
}

func (r *PerformanceRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Where("soft_deleted = ?", false)
}

// ─── Strategy ────────────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetStrategyByID(ctx context.Context, id string) (*performance.Strategy, error) {
	var s performance.Strategy
	err := r.base(ctx).
		Preload("EnterpriseObjectives").
		Preload("StrategicThemes").
		First(&s, "strategy_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetStrategyByID: %w", err)
	}
	return &s, nil
}

func (r *PerformanceRepository) GetAllStrategies(ctx context.Context) ([]performance.Strategy, error) {
	var results []performance.Strategy
	err := r.base(ctx).Preload("StrategicThemes").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetAllStrategies: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetStrategiesByStatus(ctx context.Context, status enums.Status) ([]performance.Strategy, error) {
	q := r.base(ctx)
	if status != enums.StatusAll {
		q = q.Where("status = ?", status.String())
	}
	var results []performance.Strategy
	err := q.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetStrategiesByStatus: %w", err)
	}
	return results, nil
}

// ─── StrategicTheme ──────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetStrategicThemesByStrategy(ctx context.Context, strategyID string) ([]performance.StrategicTheme, error) {
	var results []performance.StrategicTheme
	err := r.base(ctx).
		Where("strategy_id = ?", strategyID).
		Preload("EnterpriseObjectives").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetStrategicThemesByStrategy: %w", err)
	}
	return results, nil
}

// ─── PerformanceReviewPeriod ─────────────────────────────────────────────────

func (r *PerformanceRepository) GetReviewPeriodByID(ctx context.Context, id string) (*performance.PerformanceReviewPeriod, error) {
	var rp performance.PerformanceReviewPeriod
	err := r.base(ctx).
		Preload("PeriodObjectives").
		Preload("Projects").
		Preload("Committees").
		Preload("ReviewPeriodExtensions").
		First(&rp, "period_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetReviewPeriodByID: %w", err)
	}
	return &rp, nil
}

func (r *PerformanceRepository) GetAllReviewPeriods(ctx context.Context) ([]performance.PerformanceReviewPeriod, error) {
	var results []performance.PerformanceReviewPeriod
	err := r.base(ctx).Order("start_date DESC").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetAllReviewPeriods: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetActiveReviewPeriod(ctx context.Context) (*performance.PerformanceReviewPeriod, error) {
	var rp performance.PerformanceReviewPeriod
	err := r.base(ctx).
		Where("is_active = ? AND is_approved = ?", true, true).
		First(&rp).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetActiveReviewPeriod: %w", err)
	}
	return &rp, nil
}

func (r *PerformanceRepository) GetReviewPeriodsByStatus(ctx context.Context, status enums.Status) ([]performance.PerformanceReviewPeriod, error) {
	q := r.base(ctx)
	if status != enums.StatusAll {
		q = q.Where("status = ?", status.String())
	}
	var results []performance.PerformanceReviewPeriod
	err := q.Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetReviewPeriodsByStatus: %w", err)
	}
	return results, nil
}

// ─── ObjectiveCategory ───────────────────────────────────────────────────────

func (r *PerformanceRepository) GetAllObjectiveCategories(ctx context.Context) ([]performance.ObjectiveCategory, error) {
	var results []performance.ObjectiveCategory
	err := r.base(ctx).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetAllObjectiveCategories: %w", err)
	}
	return results, nil
}

// ─── CategoryDefinition ─────────────────────────────────────────────────────

func (r *PerformanceRepository) GetCategoryDefinitions(ctx context.Context, reviewPeriodID string) ([]performance.CategoryDefinition, error) {
	var results []performance.CategoryDefinition
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Preload("ObjectiveCategory").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetCategoryDefinitions: %w", err)
	}
	return results, nil
}

// ─── Enterprise Objectives ───────────────────────────────────────────────────

func (r *PerformanceRepository) GetEnterpriseObjectiveByID(ctx context.Context, id string) (*performance.EnterpriseObjective, error) {
	var obj performance.EnterpriseObjective
	err := r.base(ctx).
		Preload("DepartmentObjectives").
		First(&obj, "enterprise_objective_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetEnterpriseObjectiveByID: %w", err)
	}
	return &obj, nil
}

func (r *PerformanceRepository) GetEnterpriseObjectivesByStrategy(ctx context.Context, strategyID string) ([]performance.EnterpriseObjective, error) {
	var results []performance.EnterpriseObjective
	err := r.base(ctx).
		Where("strategy_id = ?", strategyID).
		Preload("DepartmentObjectives").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetEnterpriseObjectivesByStrategy: %w", err)
	}
	return results, nil
}

// ─── Department Objectives ───────────────────────────────────────────────────

func (r *PerformanceRepository) GetDepartmentObjectivesByEnterprise(ctx context.Context, enterpriseObjID string) ([]performance.DepartmentObjective, error) {
	var results []performance.DepartmentObjective
	err := r.base(ctx).
		Where("enterprise_objective_id = ?", enterpriseObjID).
		Preload("DivisionObjectives").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetDepartmentObjectivesByEnterprise: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetDepartmentObjectivesByDept(ctx context.Context, deptID int) ([]performance.DepartmentObjective, error) {
	var results []performance.DepartmentObjective
	err := r.base(ctx).
		Where("department_id = ?", deptID).
		Preload("DivisionObjectives").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetDepartmentObjectivesByDept: %w", err)
	}
	return results, nil
}

// ─── Division Objectives ─────────────────────────────────────────────────────

func (r *PerformanceRepository) GetDivisionObjectivesByDept(ctx context.Context, deptObjID string) ([]performance.DivisionObjective, error) {
	var results []performance.DivisionObjective
	err := r.base(ctx).
		Where("department_objective_id = ?", deptObjID).
		Preload("OfficeObjectives").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetDivisionObjectivesByDept: %w", err)
	}
	return results, nil
}

// ─── Office Objectives ───────────────────────────────────────────────────────

func (r *PerformanceRepository) GetOfficeObjectivesByDivision(ctx context.Context, divObjID string) ([]performance.OfficeObjective, error) {
	var results []performance.OfficeObjective
	err := r.base(ctx).
		Where("division_objective_id = ?", divObjID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetOfficeObjectivesByDivision: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetOfficeObjectivesByOffice(ctx context.Context, officeID int) ([]performance.OfficeObjective, error) {
	var results []performance.OfficeObjective
	err := r.base(ctx).
		Where("office_id = ?", officeID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetOfficeObjectivesByOffice: %w", err)
	}
	return results, nil
}

// ─── PeriodObjective ─────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetPeriodObjectiveByID(ctx context.Context, id string) (*performance.PeriodObjective, error) {
	var po performance.PeriodObjective
	err := r.base(ctx).
		Preload("PeriodObjectiveEvaluations").
		Preload("PeriodObjectiveDepartmentEvaluations").
		First(&po, "period_objective_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetPeriodObjectiveByID: %w", err)
	}
	return &po, nil
}

func (r *PerformanceRepository) GetPeriodObjectivesByReviewPeriod(ctx context.Context, reviewPeriodID string) ([]performance.PeriodObjective, error) {
	var results []performance.PeriodObjective
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetPeriodObjectivesByReviewPeriod: %w", err)
	}
	return results, nil
}

// ─── Planned Objectives ─────────────────────────────────────────────────────

func (r *PerformanceRepository) GetPlannedObjectivesByStaff(ctx context.Context, staffID, reviewPeriodID string) ([]performance.ReviewPeriodIndividualPlannedObjective, error) {
	var results []performance.ReviewPeriodIndividualPlannedObjective
	err := r.base(ctx).
		Where("staff_id = ? AND review_period_id = ?", staffID, reviewPeriodID).
		Preload("OperationalObjectiveWorkProducts").
		Preload("OperationalObjectiveWorkProducts.WorkProduct").
		Preload("AssignedProjects").
		Preload("AssignedCommittees").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetPlannedObjectivesByStaff: %w", err)
	}
	return results, nil
}

// ─── WorkProduct ─────────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetWorkProductByID(ctx context.Context, id string) (*performance.WorkProduct, error) {
	var wp performance.WorkProduct
	err := r.base(ctx).
		Preload("WorkProductTasks").
		Preload("OperationalObjectiveWorkProducts").
		First(&wp, "work_product_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetWorkProductByID: %w", err)
	}
	return &wp, nil
}

func (r *PerformanceRepository) GetWorkProductsByStaff(ctx context.Context, staffID string) ([]performance.WorkProduct, error) {
	var results []performance.WorkProduct
	err := r.base(ctx).
		Where("staff_id = ?", staffID).
		Preload("WorkProductTasks").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetWorkProductsByStaff: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetWorkProductsByStatus(ctx context.Context, staffID string, status enums.Status) ([]performance.WorkProduct, error) {
	q := r.base(ctx).Where("staff_id = ?", staffID)
	if status != enums.StatusAll {
		q = q.Where("status = ?", status.String())
	}
	var results []performance.WorkProduct
	err := q.Preload("WorkProductTasks").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetWorkProductsByStatus: %w", err)
	}
	return results, nil
}

// ─── WorkProductEvaluation ───────────────────────────────────────────────────

func (r *PerformanceRepository) GetEvaluationByWorkProductID(ctx context.Context, wpID string) (*performance.WorkProductEvaluation, error) {
	var eval performance.WorkProductEvaluation
	err := r.base(ctx).
		Where("work_product_id = ?", wpID).
		First(&eval).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetEvaluationByWorkProductID: %w", err)
	}
	return &eval, nil
}

// ─── EvaluationOption ────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetEvaluationOptions(ctx context.Context, evalType enums.EvaluationType) ([]performance.EvaluationOption, error) {
	var results []performance.EvaluationOption
	err := r.base(ctx).
		Where("evaluation_type = ?", evalType).
		Order("score ASC").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetEvaluationOptions: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetAllEvaluationOptions(ctx context.Context) ([]performance.EvaluationOption, error) {
	var results []performance.EvaluationOption
	err := r.base(ctx).Order("evaluation_type, score ASC").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetAllEvaluationOptions: %w", err)
	}
	return results, nil
}

// ─── PeriodScore ─────────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetPeriodScoresByStaff(ctx context.Context, staffID string) ([]performance.PeriodScore, error) {
	var results []performance.PeriodScore
	err := r.base(ctx).
		Where("staff_id = ?", staffID).
		Order("end_date DESC").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetPeriodScoresByStaff: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetPeriodScoresByReviewPeriod(ctx context.Context, reviewPeriodID string) ([]performance.PeriodScore, error) {
	var results []performance.PeriodScore
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetPeriodScoresByReviewPeriod: %w", err)
	}
	return results, nil
}

// ─── ReviewPeriodExtension ───────────────────────────────────────────────────

func (r *PerformanceRepository) GetExtensionsByReviewPeriod(ctx context.Context, reviewPeriodID string) ([]performance.ReviewPeriodExtension, error) {
	var results []performance.ReviewPeriodExtension
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetExtensionsByReviewPeriod: %w", err)
	}
	return results, nil
}

// ─── ReviewPeriod360Review ───────────────────────────────────────────────────

func (r *PerformanceRepository) Get360ReviewsByReviewPeriod(ctx context.Context, reviewPeriodID string) ([]performance.ReviewPeriod360Review, error) {
	var results []performance.ReviewPeriod360Review
	err := r.base(ctx).
		Where("review_period_id = ?", reviewPeriodID).
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.Get360ReviewsByReviewPeriod: %w", err)
	}
	return results, nil
}

// ─── Configuration ───────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetAllSettings(ctx context.Context) ([]performance.Setting, error) {
	var results []performance.Setting
	err := r.base(ctx).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetAllSettings: %w", err)
	}
	return results, nil
}

func (r *PerformanceRepository) GetSettingByName(ctx context.Context, name string) (*performance.Setting, error) {
	var s performance.Setting
	err := r.base(ctx).Where("name = ?", name).First(&s).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetSettingByName: %w", err)
	}
	return &s, nil
}

func (r *PerformanceRepository) GetAllPmsConfigurations(ctx context.Context) ([]performance.PmsConfiguration, error) {
	var results []performance.PmsConfiguration
	err := r.base(ctx).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("performanceRepo.GetAllPmsConfigurations: %w", err)
	}
	return results, nil
}

// ─── SequenceNumber ──────────────────────────────────────────────────────────

func (r *PerformanceRepository) GetNextSequence(ctx context.Context, seqType enums.SequenceNumberTypes) (int64, error) {
	var seq audit.SequenceNumber
	err := r.db.WithContext(ctx).
		Where("sequence_number_type = ?", int(seqType)).
		First(&seq).Error
	if err != nil {
		return 0, fmt.Errorf("performanceRepo.GetNextSequence: %w", err)
	}
	nextNum := seq.NextNumber
	// Atomically increment
	r.db.WithContext(ctx).
		Model(&audit.SequenceNumber{}).
		Where("sequence_number_type = ?", int(seqType)).
		Update("next_number", gorm.Expr("next_number + 1"))
	return nextNum, nil
}

// ─── Transaction helper ──────────────────────────────────────────────────────

func (r *PerformanceRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// DB exposes the underlying GORM instance for advanced queries.
func (r *PerformanceRepository) DB() *gorm.DB {
	return r.db
}
