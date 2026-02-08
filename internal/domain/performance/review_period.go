package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// Strategy represents a multi-year corporate strategy.
type Strategy struct {
	StrategyID       string     `json:"strategy_id"        gorm:"column:strategy_id;primaryKey"`
	Name             string     `json:"name"               gorm:"column:name;not null"`
	SmdReferenceCode string     `json:"smd_reference_code" gorm:"column:smd_reference_code"`
	Description      string     `json:"description"        gorm:"column:description"`
	BankYearID       int        `json:"bank_year_id"       gorm:"column:bank_year_id;not null"`
	StartDate        time.Time  `json:"start_date"         gorm:"column:start_date;not null"`
	EndDate          time.Time  `json:"end_date"           gorm:"column:end_date;not null"`
	FileImage        string     `json:"file_image"         gorm:"column:file_image"`
	domain.BaseWorkFlow

	EnterpriseObjectives []EnterpriseObjective `json:"enterprise_objectives" gorm:"foreignKey:StrategyID"`
	StrategicThemes      []StrategicTheme      `json:"strategic_themes"      gorm:"foreignKey:StrategyID"`
}

func (Strategy) TableName() string { return "pms.strategies" }

// StrategicTheme groups enterprise objectives under a strategy.
type StrategicTheme struct {
	StrategicThemeID string `json:"strategic_theme_id" gorm:"column:strategic_theme_id;primaryKey"`
	Name             string `json:"name"               gorm:"column:name;not null"`
	Description      string `json:"description"        gorm:"column:description"`
	StrategyID       string `json:"strategy_id"        gorm:"column:strategy_id;not null"`
	FileImage        string `json:"file_image"         gorm:"column:file_image"`
	domain.BaseWorkFlow

	Strategy             *Strategy             `json:"strategy"              gorm:"foreignKey:StrategyID"`
	EnterpriseObjectives []EnterpriseObjective `json:"enterprise_objectives" gorm:"foreignKey:StrategicThemeID"`
}

func (StrategicTheme) TableName() string { return "pms.strategic_themes" }

// PerformanceReviewPeriod represents a time-bound evaluation cycle.
type PerformanceReviewPeriod struct {
	PeriodID                    string                `json:"period_id"                      gorm:"column:period_id;primaryKey"`
	Year                        int                   `json:"year"                           gorm:"column:year;not null"`
	Range                       enums.ReviewPeriodRange `json:"range"                        gorm:"column:range"`
	RangeValue                  int                   `json:"range_value"                    gorm:"column:range_value"`
	Name                        string                `json:"name"                           gorm:"column:name;not null"`
	Description                 string                `json:"description"                    gorm:"column:description"`
	ShortName                   string                `json:"short_name"                     gorm:"column:short_name"`
	StartDate                   time.Time             `json:"start_date"                     gorm:"column:start_date;not null"`
	EndDate                     time.Time             `json:"end_date"                       gorm:"column:end_date;not null"`
	AllowObjectivePlanning      bool                  `json:"allow_objective_planning"        gorm:"column:allow_objective_planning;default:false"`
	AllowWorkProductPlanning    bool                  `json:"allow_work_product_planning"     gorm:"column:allow_work_product_planning;default:false"`
	AllowWorkProductEvaluation  bool                  `json:"allow_work_product_evaluation"   gorm:"column:allow_work_product_evaluation;default:false"`
	MaxPoints                   float64               `json:"max_points"                     gorm:"column:max_points;type:decimal(18,2);default:250"`
	MinNoOfObjectives           int                   `json:"min_no_of_objectives"            gorm:"column:min_no_of_objectives;default:1"`
	MaxNoOfObjectives           int                   `json:"max_no_of_objectives"            gorm:"column:max_no_of_objectives"`
	StrategyID                  string                `json:"strategy_id"                    gorm:"column:strategy_id"`
	domain.BaseWorkFlow

	Strategy               *Strategy                        `json:"strategy"                gorm:"foreignKey:StrategyID"`
	PeriodObjectives       []PeriodObjective                `json:"period_objectives"       gorm:"foreignKey:ReviewPeriodID"`
	Projects               []Project                        `json:"projects"                gorm:"foreignKey:ReviewPeriodID"`
	Committees             []Committee                      `json:"committees"              gorm:"foreignKey:ReviewPeriodID"`
	ReviewPeriodExtensions []ReviewPeriodExtension          `json:"review_period_extensions" gorm:"foreignKey:ReviewPeriodID"`
	FeedbackRequestLogs    []FeedbackRequestLog             `json:"feedback_request_logs"   gorm:"foreignKey:ReviewPeriodID"`
	PeriodScores           []PeriodScore                    `json:"period_scores"           gorm:"foreignKey:ReviewPeriodID"`
	CompetencyGapClosures  []CompetencyGapClosure           `json:"competency_gap_closures" gorm:"foreignKey:ReviewPeriodID"`
	ReviewPeriod360Reviews []ReviewPeriod360Review          `json:"review_period_360_reviews" gorm:"foreignKey:ReviewPeriodID"`
	CompetencyReviewFeedbacks []CompetencyReviewFeedback    `json:"competency_review_feedbacks" gorm:"foreignKey:ReviewPeriodID"`
}

func (PerformanceReviewPeriod) TableName() string { return "pms.performance_review_periods" }

// ReviewPeriodExtension extends a review period for a target scope.
type ReviewPeriodExtension struct {
	ReviewPeriodExtensionID string                          `json:"review_period_extension_id" gorm:"column:review_period_extension_id;primaryKey"`
	ReviewPeriodID          string                          `json:"review_period_id"           gorm:"column:review_period_id;not null"`
	TargetType              enums.ReviewPeriodExtensionTargetType `json:"target_type"          gorm:"column:target_type"`
	TargetReference         string                          `json:"target_reference"           gorm:"column:target_reference"`
	Description             string                          `json:"description"                gorm:"column:description"`
	StartDate               time.Time                       `json:"start_date"                 gorm:"column:start_date"`
	EndDate                 time.Time                       `json:"end_date"                   gorm:"column:end_date"`
	domain.BaseWorkFlow

	ReviewPeriod *PerformanceReviewPeriod `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
}

func (ReviewPeriodExtension) TableName() string { return "pms.review_period_extensions" }

// ReviewPeriod360Review enables 360 feedback for a given review period and target.
type ReviewPeriod360Review struct {
	ReviewPeriod360ReviewID string                        `json:"review_period_360_review_id" gorm:"column:review_period_360_review_id;primaryKey"`
	ReviewPeriodID          string                        `json:"review_period_id"            gorm:"column:review_period_id;not null"`
	TargetType              enums.ReviewPeriod360TargetType `json:"target_type"               gorm:"column:target_type"`
	TargetReference         string                        `json:"target_reference"            gorm:"column:target_reference"`
	domain.BaseWorkFlow

	ReviewPeriod *PerformanceReviewPeriod `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
}

func (ReviewPeriod360Review) TableName() string { return "pms.review_period_360_reviews" }

// PeriodScore stores a staff member's final score for a review period.
type PeriodScore struct {
	PeriodScoreID     string                  `json:"period_score_id"     gorm:"column:period_score_id;primaryKey"`
	ReviewPeriodID    string                  `json:"review_period_id"    gorm:"column:review_period_id;not null"`
	StaffID           string                  `json:"staff_id"            gorm:"column:staff_id;not null;index"`
	FinalScore        float64                 `json:"final_score"         gorm:"column:final_score;type:decimal(18,2)"`
	ScorePercentage   float64                 `json:"score_percentage"    gorm:"column:score_percentage;type:decimal(18,2)"`
	FinalGrade        enums.PerformanceGrade  `json:"final_grade"         gorm:"column:final_grade;not null"`
	EndDate           time.Time               `json:"end_date"            gorm:"column:end_date"`
	OfficeID          int                     `json:"office_id"           gorm:"column:office_id"`
	MinNoOfObjectives int                     `json:"min_no_of_objectives" gorm:"column:min_no_of_objectives"`
	MaxNoOfObjectives int                     `json:"max_no_of_objectives" gorm:"column:max_no_of_objectives"`
	StrategyID        string                  `json:"strategy_id"         gorm:"column:strategy_id"`
	StaffGrade        string                  `json:"staff_grade"         gorm:"column:staff_grade"`
	LocationID        string                  `json:"location_id"         gorm:"column:location_id"`
	HRDDeductedPoints float64                 `json:"hrd_deducted_points" gorm:"column:hrd_deducted_points;type:decimal(18,2);default:0"`
	IsUnderPerforming bool                    `json:"is_under_performing" gorm:"column:is_under_performing;default:false"`
	domain.BaseEntity

	ReviewPeriod *PerformanceReviewPeriod `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
	Strategy     *Strategy               `json:"strategy"      gorm:"foreignKey:StrategyID"`
}

func (PeriodScore) TableName() string { return "pms.period_scores" }
