package competency

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/identity"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
)

// Competency represents a measurable skill or capability.
type Competency struct {
	CompetencyID   int    `json:"competency_id"   gorm:"column:competency_id;primaryKey;autoIncrement"`
	CompetencyCategoryID int `json:"competency_category_id" gorm:"column:competency_category_id;not null"`
	CompetencyName string `json:"competency_name" gorm:"column:competency_name;uniqueIndex;size:70;not null"`
	Description    string `json:"description"     gorm:"column:description"`
	domain.BaseWorkFlowData

	CompetencyCategory      *CompetencyCategory       `json:"competency_category"       gorm:"foreignKey:CompetencyCategoryID"`
	CompetencyReviews       []CompetencyReview        `json:"competency_reviews"        gorm:"foreignKey:CompetencyID"`
	CompetencyRatingDefinitions []CompetencyRatingDefinition `json:"competency_rating_definitions" gorm:"foreignKey:CompetencyID"`
}

func (Competency) TableName() string { return "CoreSchema.competencies" }

// CompetencyCategory groups competencies as Technical or Behavioral.
type CompetencyCategory struct {
	CompetencyCategoryID int    `json:"competency_category_id" gorm:"column:competency_category_id;primaryKey;autoIncrement"`
	CategoryName         string `json:"category_name"          gorm:"column:category_name;uniqueIndex;size:50;not null"`
	IsTechnical          bool   `json:"is_technical"           gorm:"column:is_technical;default:false"`
	domain.BaseAudit

	Competencies              []Competency              `json:"competencies"                gorm:"foreignKey:CompetencyCategoryID"`
	CompetencyCategoryGradings []CompetencyCategoryGrading `json:"competency_category_gradings" gorm:"foreignKey:CompetencyCategoryID"`
}

func (CompetencyCategory) TableName() string { return "CoreSchema.competency_categories" }

// CompetencyCategoryGrading assigns weights per review type.
type CompetencyCategoryGrading struct {
	CompetencyCategoryGradingID int     `json:"id"                      gorm:"column:competency_category_grading_id;primaryKey;autoIncrement"`
	CompetencyCategoryID        int     `json:"competency_category_id"  gorm:"column:competency_category_id;not null"`
	ReviewTypeID                int     `json:"review_type_id"          gorm:"column:review_type_id;not null"`
	WeightPercentage            float64 `json:"weight_percentage"       gorm:"column:weight_percentage"`
	domain.BaseAudit

	CompetencyCategory *CompetencyCategory `json:"competency_category" gorm:"foreignKey:CompetencyCategoryID"`
	ReviewType         *ReviewType         `json:"review_type"         gorm:"foreignKey:ReviewTypeID"`
}

func (CompetencyCategoryGrading) TableName() string { return "CoreSchema.competency_category_gradings" }

// CompetencyRatingDefinition defines what each rating level means for a competency.
type CompetencyRatingDefinition struct {
	CompetencyRatingDefinitionID int    `json:"id"            gorm:"column:competency_rating_definition_id;primaryKey;autoIncrement"`
	CompetencyID                 int    `json:"competency_id" gorm:"column:competency_id;not null"`
	RatingID                     int    `json:"rating_id"     gorm:"column:rating_id;not null"`
	Definition                   string `json:"definition"    gorm:"column:definition"`
	domain.BaseAudit

	Rating     *Rating     `json:"rating"     gorm:"foreignKey:RatingID"`
	Competency *Competency `json:"competency" gorm:"foreignKey:CompetencyID"`
}

func (CompetencyRatingDefinition) TableName() string { return "CoreSchema.competency_rating_definitions" }

// CompetencyReview records a single review of a staff member's competency.
type CompetencyReview struct {
	CompetencyReviewID int        `json:"competency_review_id" gorm:"column:competency_review_id;primaryKey;autoIncrement"`
	EmployeeNumber     string     `json:"employee_number"      gorm:"column:employee_number"`
	ReviewPeriodID     int        `json:"review_period_id"     gorm:"column:review_period_id;not null"`
	CompetencyID       int        `json:"competency_id"        gorm:"column:competency_id;not null"`
	ReviewTypeID       int        `json:"review_type_id"       gorm:"column:review_type_id;not null"`
	ExpectedRatingID   int        `json:"expected_rating_id"   gorm:"column:expected_rating_id;not null"`
	ReviewDate         *time.Time `json:"review_date"          gorm:"column:review_date"`
	ReviewerID         string     `json:"reviewer_id"          gorm:"column:reviewer_id"`
	ReviewerName       string     `json:"reviewer_name"        gorm:"column:reviewer_name"`
	ActualRatingID     int        `json:"actual_rating_id"     gorm:"column:actual_rating_id"`
	ActualRatingName   string     `json:"actual_rating_name"   gorm:"column:actual_rating_name"`
	ActualRatingValue  int        `json:"actual_rating_value"  gorm:"column:actual_rating_value"`
	IsTechnical        bool       `json:"is_technical"         gorm:"column:is_technical"`
	EmployeeName       string     `json:"employee_name"        gorm:"column:employee_name"`
	EmployeeInitial    string     `json:"employee_initial"     gorm:"column:employee_initial"`
	EmployeeGrade      string     `json:"employee_grade"       gorm:"column:employee_grade"`
	EmployeeDepartment string     `json:"employee_department"  gorm:"column:employee_department"`
	domain.BaseAudit

	ReviewType     *ReviewType    `json:"review_type"      gorm:"foreignKey:ReviewTypeID"`
	ReviewPeriod   *ReviewPeriod  `json:"review_period"    gorm:"foreignKey:ReviewPeriodID"`
	Competency     *Competency    `json:"competency"       gorm:"foreignKey:CompetencyID"`
	ExpectedRating *Rating        `json:"expected_rating"  gorm:"foreignKey:ExpectedRatingID"`
	DevelopmentPlans []DevelopmentPlan `json:"development_plans" gorm:"foreignKey:CompetencyReviewProfileID"`
}

func (CompetencyReview) TableName() string { return "CoreSchema.competency_reviews" }

// CompetencyReviewProfile is a summary of a staff member's competency review.
type CompetencyReviewProfile struct {
	CompetencyReviewProfileID int     `json:"competency_review_profile_id" gorm:"column:competency_review_profile_id;primaryKey;autoIncrement"`
	ReviewPeriodID            int     `json:"review_period_id"             gorm:"column:review_period_id"`
	ReviewPeriodName          string  `json:"review_period_name"           gorm:"column:review_period_name"`
	AverageRatingID           int     `json:"average_rating_id"            gorm:"column:average_rating_id"`
	AverageRatingName         string  `json:"average_rating_name"          gorm:"column:average_rating_name"`
	AverageRatingValue        int     `json:"average_rating_value"         gorm:"column:average_rating_value"`
	ExpectedRatingID          int     `json:"expected_rating_id"           gorm:"column:expected_rating_id"`
	ExpectedRatingName        string  `json:"expected_rating_name"         gorm:"column:expected_rating_name"`
	ExpectedRatingValue       int     `json:"expected_rating_value"        gorm:"column:expected_rating_value"`
	AverageScore              float64 `json:"average_score"                gorm:"column:average_score"`
	EmployeeNumber            string  `json:"employee_number"              gorm:"column:employee_number"`
	EmployeeName              string  `json:"employee_name"                gorm:"column:employee_name"`
	CompetencyID              int     `json:"competency_id"                gorm:"column:competency_id"`
	CompetencyName            string  `json:"competency_name"              gorm:"column:competency_name"`
	CompetencyCategoryName    string  `json:"competency_category_name"     gorm:"column:competency_category_name"`
	CompetencyGap             int     `json:"competency_gap"               gorm:"column:competency_gap"`
	HaveGap                   bool    `json:"have_gap"                     gorm:"column:have_gap"`
	OfficeID                  string  `json:"office_id"                    gorm:"column:office_id"`
	OfficeName                string  `json:"office_name"                  gorm:"column:office_name"`
	DivisionID                string  `json:"division_id"                  gorm:"column:division_id"`
	DivisionName              string  `json:"division_name"                gorm:"column:division_name"`
	DepartmentID              string  `json:"department_id"                gorm:"column:department_id"`
	DepartmentName            string  `json:"department_name"              gorm:"column:department_name"`
	JobRoleID                 string  `json:"job_role_id"                  gorm:"column:job_role_id"`
	JobRoleName               string  `json:"job_role_name"                gorm:"column:job_role_name"`
	GradeName                 string  `json:"grade_name"                   gorm:"column:grade_name"`
	domain.BaseAudit

	DevelopmentPlans []DevelopmentPlan `json:"development_plans" gorm:"foreignKey:CompetencyReviewProfileID"`
}

func (CompetencyReviewProfile) TableName() string { return "CoreSchema.competency_review_profiles" }

// DevelopmentPlan records a training/development action to close a competency gap.
type DevelopmentPlan struct {
	DevelopmentPlanID         int        `json:"development_plan_id"          gorm:"column:development_plan_id;primaryKey;autoIncrement"`
	CompetencyReviewProfileID int        `json:"competency_review_profile_id" gorm:"column:competency_review_profile_id;not null"`
	TrainingTypeName          string     `json:"training_type_name"           gorm:"column:training_type_name"`
	Activity                  string     `json:"activity"                     gorm:"column:activity"`
	EmployeeNumber            string     `json:"employee_number"              gorm:"column:employee_number"`
	TargetDate                time.Time  `json:"target_date"                  gorm:"column:target_date"`
	CompletionDate            *time.Time `json:"completion_date"              gorm:"column:completion_date"`
	TaskStatus                string     `json:"task_status"                  gorm:"column:task_status"`
	LearningResource          string     `json:"learning_resource"            gorm:"column:learning_resource"`
	domain.BaseAudit

	CompetencyReviewProfile *CompetencyReviewProfile `json:"competency_review_profile" gorm:"foreignKey:CompetencyReviewProfileID"`
}

func (DevelopmentPlan) TableName() string { return "CoreSchema.development_plans" }

// JobRole defines a role within an office.
type JobRole struct {
	JobRoleID   int    `json:"job_role_id"  gorm:"column:job_role_id;primaryKey;autoIncrement"`
	JobRoleName string `json:"job_role_name" gorm:"column:job_role_name;not null"`
	Description string `json:"description"  gorm:"column:description"`
	domain.BaseAudit

	OfficeJobRoles      []OfficeJobRole      `json:"office_job_roles"      gorm:"foreignKey:JobRoleID"`
	JobRoleGrades       []JobRoleGrade       `json:"job_role_grades"       gorm:"foreignKey:JobRoleID"`
	JobRoleCompetencies []JobRoleCompetency  `json:"job_role_competencies" gorm:"foreignKey:JobRoleID"`
}

func (JobRole) TableName() string { return "CoreSchema.job_roles" }

// JobGrade represents a salary/grade band.
type JobGrade struct {
	JobGradeID int    `json:"job_grade_id" gorm:"column:job_grade_id;primaryKey;autoIncrement"`
	GradeCode  string `json:"grade_code"   gorm:"column:grade_code;uniqueIndex;size:10;not null"`
	GradeName  string `json:"grade_name"   gorm:"column:grade_name"`
	domain.BaseAudit

	AssignJobGradeGroups []AssignJobGradeGroup `json:"assign_job_grade_groups" gorm:"foreignKey:JobGradeID"`
}

func (JobGrade) TableName() string { return "CoreSchema.job_grades" }

// JobGradeGroup groups grades (Junior, Officer, Manager, Executive).
type JobGradeGroup struct {
	JobGradeGroupID int    `json:"job_grade_group_id" gorm:"column:job_grade_group_id;primaryKey;autoIncrement"`
	GroupName       string `json:"group_name"         gorm:"column:group_name;uniqueIndex;not null"`
	Order           int    `json:"order"              gorm:"column:order"`
	domain.BaseAudit

	AssignJobGradeGroups    []AssignJobGradeGroup    `json:"assign_job_grade_groups"    gorm:"foreignKey:JobGradeGroupID"`
	BehavioralCompetencies  []BehavioralCompetency   `json:"behavioral_competencies"    gorm:"foreignKey:JobGradeGroupID"`
}

func (JobGradeGroup) TableName() string { return "CoreSchema.job_grade_groups" }

// AssignJobGradeGroup maps a grade to a grade group.
type AssignJobGradeGroup struct {
	AssignJobGradeGroupID int `json:"assign_job_grade_group_id" gorm:"column:assign_job_grade_group_id;primaryKey;autoIncrement"`
	JobGradeGroupID       int `json:"job_grade_group_id"        gorm:"column:job_grade_group_id;not null"`
	JobGradeID            int `json:"job_grade_id"              gorm:"column:job_grade_id;not null"`
	domain.BaseAudit

	JobGrade      *JobGrade      `json:"job_grade"       gorm:"foreignKey:JobGradeID"`
	JobGradeGroup *JobGradeGroup `json:"job_grade_group" gorm:"foreignKey:JobGradeGroupID"`
}

func (AssignJobGradeGroup) TableName() string { return "CoreSchema.assign_job_grade_groups" }

// JobRoleCompetency maps a required competency + expected rating to a job role in an office.
type JobRoleCompetency struct {
	JobRoleCompetencyID int `json:"job_role_competency_id" gorm:"column:job_role_competency_id;primaryKey;autoIncrement"`
	OfficeID            int `json:"office_id"              gorm:"column:office_id;uniqueIndex:idx_jrc_unique"`
	JobRoleID           int `json:"job_role_id"            gorm:"column:job_role_id;uniqueIndex:idx_jrc_unique"`
	CompetencyID        int `json:"competency_id"          gorm:"column:competency_id;uniqueIndex:idx_jrc_unique"`
	RatingID            int `json:"rating_id"              gorm:"column:rating_id"`
	domain.BaseAudit

	Competency *Competency       `json:"competency" gorm:"foreignKey:CompetencyID"`
	JobRole    *JobRole           `json:"job_role"   gorm:"foreignKey:JobRoleID"`
	Rating     *Rating            `json:"rating"     gorm:"foreignKey:RatingID"`
	Office     *organogram.Office `json:"office"     gorm:"foreignKey:OfficeID"`
}

func (JobRoleCompetency) TableName() string { return "CoreSchema.job_role_competencies" }

// JobRoleGrade maps a grade to a job role.
type JobRoleGrade struct {
	JobRoleGradeID int    `json:"job_role_grade_id" gorm:"column:job_role_grade_id;primaryKey;autoIncrement"`
	JobRoleID      int    `json:"job_role_id"       gorm:"column:job_role_id;not null"`
	GradeID        string `json:"grade_id"          gorm:"column:grade_id"`
	GradeName      string `json:"grade_name"        gorm:"column:grade_name"`
	domain.BaseAudit

	JobRole *JobRole `json:"job_role" gorm:"foreignKey:JobRoleID"`
}

func (JobRoleGrade) TableName() string { return "CoreSchema.job_role_grades" }

// BehavioralCompetency maps a required behavioral competency to a grade group.
type BehavioralCompetency struct {
	BehavioralCompetencyID int `json:"behavioral_competency_id" gorm:"column:behavioral_competency_id;primaryKey;autoIncrement"`
	CompetencyID           int `json:"competency_id"            gorm:"column:competency_id;uniqueIndex:idx_bc_unique"`
	JobGradeGroupID        int `json:"job_grade_group_id"       gorm:"column:job_grade_group_id;uniqueIndex:idx_bc_unique"`
	RatingID               int `json:"rating_id"                gorm:"column:rating_id"`
	domain.BaseAudit

	Competency    *Competency    `json:"competency"      gorm:"foreignKey:CompetencyID"`
	Rating        *Rating        `json:"rating"          gorm:"foreignKey:RatingID"`
	JobGradeGroup *JobGradeGroup `json:"job_grade_group" gorm:"foreignKey:JobGradeGroupID"`
}

func (BehavioralCompetency) TableName() string { return "CoreSchema.behavioral_competencies" }

// OfficeJobRole maps a job role to an office.
type OfficeJobRole struct {
	OfficeJobRoleID int `json:"office_job_role_id" gorm:"column:office_job_role_id;primaryKey;autoIncrement"`
	OfficeID        int `json:"office_id"          gorm:"column:office_id;not null"`
	JobRoleID       int `json:"job_role_id"        gorm:"column:job_role_id;not null"`
	domain.BaseAudit

	Office  *organogram.Office `json:"office"   gorm:"foreignKey:OfficeID"`
	JobRole *JobRole           `json:"job_role" gorm:"foreignKey:JobRoleID"`
}

func (OfficeJobRole) TableName() string { return "CoreSchema.office_job_roles" }

// Rating defines a numeric rating level.
type Rating struct {
	RatingID int    `json:"rating_id" gorm:"column:rating_id;primaryKey;autoIncrement"`
	Name     string `json:"name"      gorm:"column:name;uniqueIndex;size:50;not null"`
	Value    int    `json:"value"     gorm:"column:value"`
	domain.BaseAudit
}

func (Rating) TableName() string { return "CoreSchema.ratings" }

// ReviewPeriod (Competency module) represents a competency review cycle.
type ReviewPeriod struct {
	ReviewPeriodID int       `json:"review_period_id" gorm:"column:review_period_id;primaryKey;autoIncrement"`
	BankYearID     int       `json:"bank_year_id"     gorm:"column:bank_year_id;not null"`
	Name           string    `json:"name"             gorm:"column:name;not null"`
	StartDate      time.Time `json:"start_date"       gorm:"column:start_date"`
	EndDate        time.Time `json:"end_date"         gorm:"column:end_date"`
	domain.BaseWorkFlowData

	BankYear          *identity.BankYear  `json:"bank_year"          gorm:"foreignKey:BankYearID"`
	CompetencyReviews []CompetencyReview  `json:"competency_reviews" gorm:"foreignKey:ReviewPeriodID"`
}

func (ReviewPeriod) TableName() string { return "CoreSchema.review_periods" }

// ReviewType classifies reviews (e.g. Self, Peer, Manager).
type ReviewType struct {
	ReviewTypeID   int    `json:"review_type_id"   gorm:"column:review_type_id;primaryKey;autoIncrement"`
	ReviewTypeName string `json:"review_type_name" gorm:"column:review_type_name;uniqueIndex;size:50;not null"`
	domain.BaseAudit
}

func (ReviewType) TableName() string { return "CoreSchema.review_types" }

// TrainingType classifies development activities.
type TrainingType struct {
	TrainingTypeID   int    `json:"training_type_id"   gorm:"column:training_type_id;primaryKey;autoIncrement"`
	TrainingTypeName string `json:"training_type_name" gorm:"column:training_type_name;not null"`
	domain.BaseAudit
}

func (TrainingType) TableName() string { return "CoreSchema.training_types" }

// StaffJobRoles records a staff member's assigned job role with HRD approval workflow.
type StaffJobRoles struct {
	StaffJobRoleID int    `json:"staff_job_role_id" gorm:"column:staff_job_role_id;primaryKey;autoIncrement"`
	EmployeeID     string `json:"employee_id"       gorm:"column:employee_id;not null"`
	FullName       string `json:"full_name"         gorm:"column:full_name"`
	DepartmentID   int    `json:"department_id"     gorm:"column:department_id"`
	DivisionID     int    `json:"division_id"       gorm:"column:division_id"`
	OfficeID       int    `json:"office_id"         gorm:"column:office_id"`
	SupervisorID   string `json:"supervisor_id"     gorm:"column:supervisor_id"`
	JobRoleID      int    `json:"job_role_id"       gorm:"column:job_role_id"`
	JobRoleName    string `json:"job_role_name"     gorm:"column:job_role_name"`
	SoaStatus      bool   `json:"soa_status"        gorm:"column:soa_status;default:false"`
	SoaResponse    string `json:"soa_response"      gorm:"column:soa_response"`
	domain.HrdWorkFlowData
}

func (StaffJobRoles) TableName() string { return "CoreSchema.staff_job_roles" }
