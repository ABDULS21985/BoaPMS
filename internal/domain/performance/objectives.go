package performance

import (
	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
)

// ObjectiveCategory groups enterprise objectives.
type ObjectiveCategory struct {
	ObjectiveCategoryID string `json:"objective_category_id" gorm:"column:objective_category_id;primaryKey"`
	Name                string `json:"name"                  gorm:"column:name;not null"`
	Description         string `json:"description"           gorm:"column:description"`
	domain.BaseWorkFlow

	PmsCompetencies      []PmsCompetency       `json:"pms_competencies"      gorm:"foreignKey:ObjectCategoryID"`
	EnterpriseObjectives []EnterpriseObjective  `json:"enterprise_objectives" gorm:"foreignKey:EnterpriseObjectivesCategoryID"`
	CategoryDefinitions  []CategoryDefinition   `json:"category_definitions"  gorm:"foreignKey:ObjectiveCategoryID"`
	CompetencyGapClosures []CompetencyGapClosure `json:"competency_gap_closures" gorm:"foreignKey:ObjectiveCategoryID"`
}

func (ObjectiveCategory) TableName() string { return "pms.objective_categories" }

// CategoryDefinition defines weight/limits for an objective category in a review period.
type CategoryDefinition struct {
	DefinitionID        string  `json:"definition_id"         gorm:"column:definition_id;primaryKey"`
	ObjectiveCategoryID string  `json:"objective_category_id" gorm:"column:objective_category_id;not null"`
	ReviewPeriodID      string  `json:"review_period_id"      gorm:"column:review_period_id;not null"`
	Weight              float64 `json:"weight"                gorm:"column:weight;type:decimal(18,2);not null"`
	MaxNoObjectives     int     `json:"max_no_objectives"     gorm:"column:max_no_objectives"`
	MaxNoWorkProduct    int     `json:"max_no_work_product"   gorm:"column:max_no_work_product"`
	MaxPoints           float64 `json:"max_points"            gorm:"column:max_points;type:decimal(18,2)"`
	IsCompulsory        bool    `json:"is_compulsory"         gorm:"column:is_compulsory;default:false"`
	EnforceWorkProductLimit bool `json:"enforce_work_product_limit" gorm:"column:enforce_work_product_limit;default:false"`
	Description         string  `json:"description"           gorm:"column:description"`
	GradeGroupID        int     `json:"grade_group_id"        gorm:"column:grade_group_id"`
	domain.BaseWorkFlow

	Category     *ObjectiveCategory       `json:"category"      gorm:"foreignKey:ObjectiveCategoryID"`
	ReviewPeriod *PerformanceReviewPeriod  `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
}

func (CategoryDefinition) TableName() string { return "pms.category_definitions" }

// EnterpriseObjective is a top-level strategic objective.
type EnterpriseObjective struct {
	EnterpriseObjectiveID        string            `json:"enterprise_objective_id"         gorm:"column:enterprise_objective_id;primaryKey"`
	Type                         enums.ObjectiveType `json:"type"                          gorm:"column:type;default:1"`
	EnterpriseObjectivesCategoryID string           `json:"enterprise_objectives_category_id" gorm:"column:enterprise_objectives_category_id;not null"`
	StrategicThemeID             string            `json:"strategic_theme_id"              gorm:"column:strategic_theme_id"`
	StrategyID                   string            `json:"strategy_id"                    gorm:"column:strategy_id;not null"`
	domain.ObjectiveBase

	Category             *ObjectiveCategory    `json:"category"              gorm:"foreignKey:EnterpriseObjectivesCategoryID"`
	Strategy             *Strategy             `json:"strategy"              gorm:"foreignKey:StrategyID"`
	StrategicTheme       *StrategicTheme       `json:"strategic_theme"       gorm:"foreignKey:StrategicThemeID"`
	DepartmentObjectives []DepartmentObjective `json:"department_objectives" gorm:"foreignKey:EnterpriseObjectiveID"`
}

func (EnterpriseObjective) TableName() string { return "pms.enterprise_objectives" }

// DepartmentObjective cascades an enterprise objective to a department.
type DepartmentObjective struct {
	DepartmentObjectiveID string `json:"department_objective_id" gorm:"column:department_objective_id;primaryKey"`
	DepartmentID          int    `json:"department_id"           gorm:"column:department_id;not null"`
	EnterpriseObjectiveID string `json:"enterprise_objective_id" gorm:"column:enterprise_objective_id;not null"`
	domain.ObjectiveBase

	Department          *organogram.Department `json:"department"           gorm:"foreignKey:DepartmentID"`
	EnterpriseObjective *EnterpriseObjective   `json:"enterprise_objective" gorm:"foreignKey:EnterpriseObjectiveID"`
	DivisionObjectives  []DivisionObjective    `json:"division_objectives"  gorm:"foreignKey:DepartmentObjectiveID"`
}

func (DepartmentObjective) TableName() string { return "pms.department_objectives" }

// DivisionObjective cascades a department objective to a division.
type DivisionObjective struct {
	DivisionObjectiveID   string `json:"division_objective_id"   gorm:"column:division_objective_id;primaryKey"`
	DivisionID            int    `json:"division_id"             gorm:"column:division_id;not null"`
	DepartmentObjectiveID string `json:"department_objective_id" gorm:"column:department_objective_id;not null"`
	domain.ObjectiveBase

	Division            *organogram.Division  `json:"division"             gorm:"foreignKey:DivisionID"`
	DepartmentObjective *DepartmentObjective  `json:"department_objective" gorm:"foreignKey:DepartmentObjectiveID"`
	OfficeObjectives    []OfficeObjective     `json:"office_objectives"    gorm:"foreignKey:DivisionObjectiveID"`
}

func (DivisionObjective) TableName() string { return "pms.division_objectives" }

// OfficeObjective cascades a division objective to an office (lowest level).
type OfficeObjective struct {
	OfficeObjectiveID   string `json:"office_objective_id"   gorm:"column:office_objective_id;primaryKey"`
	OfficeID            int    `json:"office_id"             gorm:"column:office_id;not null"`
	DivisionObjectiveID string `json:"division_objective_id" gorm:"column:division_objective_id;not null"`
	JobGradeGroupID     int    `json:"job_grade_group_id"    gorm:"column:job_grade_group_id;not null"`
	domain.ObjectiveBase

	Office            *organogram.Office   `json:"office"             gorm:"foreignKey:OfficeID"`
	DivisionObjective *DivisionObjective   `json:"division_objective" gorm:"foreignKey:DivisionObjectiveID"`
}

func (OfficeObjective) TableName() string { return "pms.office_objectives" }

// PeriodObjective links an enterprise objective to a review period.
type PeriodObjective struct {
	PeriodObjectiveID string `json:"period_objective_id" gorm:"column:period_objective_id;primaryKey"`
	ObjectiveID       string `json:"objective_id"        gorm:"column:objective_id;not null"`
	ReviewPeriodID    string `json:"review_period_id"    gorm:"column:review_period_id;not null"`
	domain.BaseWorkFlow

	Objective    *EnterpriseObjective     `json:"objective"     gorm:"foreignKey:ObjectiveID"`
	ReviewPeriod *PerformanceReviewPeriod  `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
	PeriodObjectiveEvaluations          []PeriodObjectiveEvaluation          `json:"evaluations"      gorm:"foreignKey:PeriodObjectiveID"`
	PeriodObjectiveDepartmentEvaluations []PeriodObjectiveDepartmentEvaluation `json:"dept_evaluations" gorm:"foreignKey:PeriodObjectiveID"`
}

func (PeriodObjective) TableName() string { return "pms.period_objectives" }

// PeriodObjectiveEvaluation stores outcome evaluation for a period objective.
type PeriodObjectiveEvaluation struct {
	PeriodObjectiveEvaluationID string  `json:"period_objective_evaluation_id" gorm:"column:period_objective_evaluation_id;primaryKey"`
	TotalOutcomeScore           float64 `json:"total_outcome_score"            gorm:"column:total_outcome_score;type:decimal(18,2)"`
	OutcomeScore                float64 `json:"outcome_score"                  gorm:"column:outcome_score;type:decimal(18,2)"`
	PeriodObjectiveID           string  `json:"period_objective_id"            gorm:"column:period_objective_id;not null"`
	domain.BaseWorkFlow

	PeriodObjective *PeriodObjective `json:"period_objective" gorm:"foreignKey:PeriodObjectiveID"`
}

func (PeriodObjectiveEvaluation) TableName() string { return "pms.period_objective_evaluations" }

// PeriodObjectiveDepartmentEvaluation stores outcome evaluation per department.
type PeriodObjectiveDepartmentEvaluation struct {
	PeriodObjectiveDepartmentEvaluationID string  `json:"id"                     gorm:"column:period_objective_department_evaluation_id;primaryKey"`
	OverallOutcomeScored                  float64 `json:"overall_outcome_scored" gorm:"column:overall_outcome_scored;type:decimal(18,2)"`
	AllocatedOutcome                      float64 `json:"allocated_outcome"      gorm:"column:allocated_outcome;type:decimal(18,2)"`
	OutcomeScore                          float64 `json:"outcome_score"          gorm:"column:outcome_score;type:decimal(18,2)"`
	DepartmentID                          int     `json:"department_id"          gorm:"column:department_id;not null"`
	PeriodObjectiveID                     string  `json:"period_objective_id"    gorm:"column:period_objective_id;not null"`
	domain.BaseWorkFlow

	PeriodObjective *PeriodObjective       `json:"period_objective" gorm:"foreignKey:PeriodObjectiveID"`
	Department      *organogram.Department `json:"department"       gorm:"foreignKey:DepartmentID"`
}

func (PeriodObjectiveDepartmentEvaluation) TableName() string {
	return "pms.period_objective_department_evaluations"
}

// ReviewPeriodIndividualPlannedObjective maps a staff member to an objective for a period.
type ReviewPeriodIndividualPlannedObjective struct {
	PlannedObjectiveID string              `json:"planned_objective_id" gorm:"column:planned_objective_id;primaryKey"`
	ObjectiveID        string              `json:"objective_id"         gorm:"column:objective_id;not null"`
	StaffID            string              `json:"staff_id"             gorm:"column:staff_id;not null;index"`
	ObjectiveLevel     enums.ObjectiveLevel `json:"objective_level"     gorm:"column:objective_level;not null"`
	StaffJobRole       string              `json:"staff_job_role"       gorm:"column:staff_job_role;not null"`
	ReviewPeriodID     string              `json:"review_period_id"     gorm:"column:review_period_id;not null"`
	NoReturned         int                 `json:"no_returned"          gorm:"column:no_returned;default:0"`
	Remark             string              `json:"remark"               gorm:"column:remark"`
	domain.BaseWorkFlow

	ReviewPeriod                  *PerformanceReviewPeriod         `json:"review_period"                   gorm:"foreignKey:ReviewPeriodID"`
	OperationalObjectiveWorkProducts []OperationalObjectiveWorkProduct `json:"operational_objective_work_products" gorm:"foreignKey:PlannedObjectiveID"`
	AssignedProjects              []ProjectMember                  `json:"assigned_projects"               gorm:"foreignKey:PlannedObjectiveID"`
	AssignedCommittees            []CommitteeMember                `json:"assigned_committees"             gorm:"foreignKey:PlannedObjectiveID"`
}

func (ReviewPeriodIndividualPlannedObjective) TableName() string {
	return "pms.review_period_individual_planned_objectives"
}
