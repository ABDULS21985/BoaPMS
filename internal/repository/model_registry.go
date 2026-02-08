package repository

import (
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/audit"
	"github.com/enterprise-pms/pms-api/internal/domain/competency"
	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"gorm.io/gorm"
)

// AutoMigrateAll runs GORM AutoMigrate for every domain model registered
// against the primary PostgreSQL database. This mirrors the EF Core DbSet<T>
// declarations in CompetencyCoreDbContext plus the PMS schema models.
//
// NOTE: In production, prefer explicit SQL migrations (go-pms/migrations/).
// Use this for development/test bootstrapping only.
func AutoMigrateAll(db *gorm.DB) error {
	models := []interface{}{
		// ── Identity (CoreSchema) ───────────────────────────────────────
		&identity.ApplicationUser{},
		&identity.ApplicationRole{},
		&identity.Permission{},
		&identity.RolePermission{},
		&identity.BankYear{},

		// ── Organogram (CoreSchema) ─────────────────────────────────────
		&organogram.Directorate{},
		&organogram.Department{},
		&organogram.Division{},
		&organogram.Office{},

		// ── Competency (CoreSchema) ─────────────────────────────────────
		&competency.Competency{},
		&competency.CompetencyCategory{},
		&competency.CompetencyCategoryGrading{},
		&competency.CompetencyRatingDefinition{},
		&competency.CompetencyReview{},
		&competency.CompetencyReviewProfile{},
		&competency.DevelopmentPlan{},
		&competency.JobRole{},
		&competency.JobGrade{},
		&competency.JobGradeGroup{},
		&competency.AssignJobGradeGroup{},
		&competency.OfficeJobRole{},
		&competency.JobRoleCompetency{},
		&competency.BehavioralCompetency{},
		&competency.JobRoleGrade{},
		&competency.Rating{},
		&competency.ReviewType{},
		&competency.TrainingType{},
		&competency.StaffJobRoles{},
		&competency.ReviewPeriod{},

		// ── Performance (pms schema) ────────────────────────────────────
		&performance.Strategy{},
		&performance.StrategicTheme{},
		&performance.PerformanceReviewPeriod{},
		&performance.ObjectiveCategory{},
		&performance.CategoryDefinition{},
		&performance.EnterpriseObjective{},
		&performance.DepartmentObjective{},
		&performance.DivisionObjective{},
		&performance.OfficeObjective{},
		&performance.PeriodObjective{},
		&performance.PeriodObjectiveEvaluation{},
		&performance.PeriodObjectiveDepartmentEvaluation{},
		&performance.ReviewPeriodIndividualPlannedObjective{},
		&performance.WorkProduct{},
		&performance.WorkProductTask{},
		&performance.WorkProductEvaluation{},
		&performance.EvaluationOption{},
		&performance.OperationalObjectiveWorkProduct{},
		&performance.PeriodScore{},
		&performance.ReviewPeriodExtension{},
		&performance.ReviewPeriod360Review{},
		&performance.FeedbackRequestLog{},
		&performance.FeedbackQuestionaire{},
		&performance.FeedbackQuestionaireOption{},
		&performance.PmsCompetency{},
		&performance.CompetencyReviewFeedback{},
		&performance.CompetencyReviewer{},
		&performance.CompetencyReviewerRating{},
		&performance.CompetencyGapClosure{},
		&performance.Grievance{},
		&performance.GrievanceResolution{},
		&performance.Project{},
		&performance.Committee{},
		&performance.ProjectMember{},
		&performance.CommitteeMember{},
		&performance.ProjectWorkProduct{},
		&performance.CommitteeWorkProduct{},
		&performance.ProjectObjective{},
		&performance.CommitteeObjective{},
		&performance.ProjectAssignedWorkProduct{},
		&performance.CommitteeAssignedWorkProduct{},
		&performance.PmsConfiguration{},
		&performance.Setting{},
		&performance.WorkProductDefinition{},
		&performance.CascadedWorkProduct{},

		// ── Audit (pmsaudit schema) ─────────────────────────────────────
		&audit.AuditLog{},
		&audit.AuditableEntity{},
		&audit.AuditableAttribute{},
		&audit.SequenceNumber{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto-migrating models: %w", err)
	}
	return nil
}
