package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/organogram"
)

// BaseProject contains shared fields for Project and Committee.
type BaseProject struct {
	Name           string    `json:"name"             gorm:"column:name;not null"`
	Description    string    `json:"description"      gorm:"column:description"`
	StartDate      time.Time `json:"start_date"       gorm:"column:start_date;not null"`
	EndDate        time.Time `json:"end_date"         gorm:"column:end_date;not null"`
	Deliverables   string    `json:"deliverables"     gorm:"column:deliverables"`
	ReviewPeriodID string    `json:"review_period_id" gorm:"column:review_period_id;not null"`
	DepartmentID   int       `json:"department_id"    gorm:"column:department_id;not null"`
	domain.BaseWorkFlow

	ReviewPeriod *PerformanceReviewPeriod `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
	Department   *organogram.Department   `json:"department"    gorm:"foreignKey:DepartmentID"`
}

// Project is an ad-hoc project within a review period.
type Project struct {
	ProjectID      string `json:"project_id"      gorm:"column:project_id;primaryKey"`
	ProjectManager string `json:"project_manager" gorm:"column:project_manager"`
	BaseProject

	ProjectMembers         []ProjectMember         `json:"project_members"          gorm:"foreignKey:ProjectID"`
	ProjectWorkProducts    []ProjectWorkProduct    `json:"project_work_products"    gorm:"foreignKey:ProjectID"`
	ProjectObjectives      []ProjectObjective      `json:"project_objectives"       gorm:"foreignKey:ProjectID"`
	ProjectAssignedWorkProducts []ProjectAssignedWorkProduct `json:"project_assigned_work_products" gorm:"foreignKey:ProjectID"`
}

func (Project) TableName() string { return "pms.projects" }

// Committee is an ad-hoc committee within a review period.
type Committee struct {
	CommitteeID string `json:"committee_id" gorm:"column:committee_id;primaryKey"`
	Chairperson string `json:"chairperson"  gorm:"column:chairperson"`
	BaseProject

	CommitteeMembers         []CommitteeMember         `json:"committee_members"          gorm:"foreignKey:CommitteeID"`
	CommitteeWorkProducts    []CommitteeWorkProduct    `json:"committee_work_products"    gorm:"foreignKey:CommitteeID"`
	CommitteeObjectives      []CommitteeObjective      `json:"committee_objectives"       gorm:"foreignKey:CommitteeID"`
	CommitteeAssignedWorkProducts []CommitteeAssignedWorkProduct `json:"committee_assigned_work_products" gorm:"foreignKey:CommitteeID"`
}

func (Committee) TableName() string { return "pms.committees" }

// ProjectMember assigns a staff member to a project.
type ProjectMember struct {
	ProjectMemberID    string `json:"project_member_id"    gorm:"column:project_member_id;primaryKey"`
	StaffID            string `json:"staff_id"             gorm:"column:staff_id;not null"`
	ProjectID          string `json:"project_id"           gorm:"column:project_id;not null"`
	PlannedObjectiveID string `json:"planned_objective_id" gorm:"column:planned_objective_id"`
	domain.BaseWorkFlow

	Project          *Project                                `json:"project"           gorm:"foreignKey:ProjectID"`
	PlannedObjective *ReviewPeriodIndividualPlannedObjective  `json:"planned_objective" gorm:"foreignKey:PlannedObjectiveID"`
}

func (ProjectMember) TableName() string { return "pms.project_members" }

// CommitteeMember assigns a staff member to a committee.
type CommitteeMember struct {
	CommitteeMemberID  string `json:"committee_member_id"  gorm:"column:committee_member_id;primaryKey"`
	StaffID            string `json:"staff_id"             gorm:"column:staff_id;not null"`
	CommitteeID        string `json:"committee_id"         gorm:"column:committee_id;not null"`
	PlannedObjectiveID string `json:"planned_objective_id" gorm:"column:planned_objective_id"`
	domain.BaseWorkFlow

	Committee        *Committee                              `json:"committee"         gorm:"foreignKey:CommitteeID"`
	PlannedObjective *ReviewPeriodIndividualPlannedObjective  `json:"planned_objective" gorm:"foreignKey:PlannedObjectiveID"`
}

func (CommitteeMember) TableName() string { return "pms.committee_members" }

// ProjectWorkProduct links a work product to a project.
type ProjectWorkProduct struct {
	ProjectWorkProductID       string `json:"project_work_product_id"        gorm:"column:project_work_product_id;primaryKey"`
	WorkProductID              string `json:"work_product_id"                gorm:"column:work_product_id;not null"`
	ProjectAssignedWorkProductID string `json:"project_assigned_work_product_id" gorm:"column:project_assigned_work_product_id"`
	ProjectID                  string `json:"project_id"                     gorm:"column:project_id;not null"`
	domain.BaseEntity

	WorkProduct              *WorkProduct              `json:"work_product"               gorm:"foreignKey:WorkProductID"`
	Project                  *Project                  `json:"project"                    gorm:"foreignKey:ProjectID"`
	ProjectAssignedWorkProduct *ProjectAssignedWorkProduct `json:"project_assigned_work_product" gorm:"foreignKey:ProjectAssignedWorkProductID"`
}

func (ProjectWorkProduct) TableName() string { return "pms.project_work_products" }

// CommitteeWorkProduct links a work product to a committee.
type CommitteeWorkProduct struct {
	CommitteeWorkProductID       string `json:"committee_work_product_id"        gorm:"column:committee_work_product_id;primaryKey"`
	WorkProductID                string `json:"work_product_id"                  gorm:"column:work_product_id;not null"`
	CommitteeAssignedWorkProductID string `json:"committee_assigned_work_product_id" gorm:"column:committee_assigned_work_product_id"`
	CommitteeID                  string `json:"committee_id"                     gorm:"column:committee_id;not null"`
	domain.BaseEntity

	WorkProduct                *WorkProduct                `json:"work_product"                 gorm:"foreignKey:WorkProductID"`
	Committee                  *Committee                  `json:"committee"                    gorm:"foreignKey:CommitteeID"`
	CommitteeAssignedWorkProduct *CommitteeAssignedWorkProduct `json:"committee_assigned_work_product" gorm:"foreignKey:CommitteeAssignedWorkProductID"`
}

func (CommitteeWorkProduct) TableName() string { return "pms.committee_work_products" }

// ProjectObjective links an enterprise objective to a project.
type ProjectObjective struct {
	ProjectObjectiveID string `json:"project_objective_id" gorm:"column:project_objective_id;primaryKey"`
	ObjectiveID        string `json:"objective_id"         gorm:"column:objective_id;not null"`
	ProjectID          string `json:"project_id"           gorm:"column:project_id;not null"`
	domain.BaseEntity

	Project   *Project              `json:"project"   gorm:"foreignKey:ProjectID"`
	Objective *EnterpriseObjective  `json:"objective" gorm:"foreignKey:ObjectiveID"`
}

func (ProjectObjective) TableName() string { return "pms.project_objectives" }

// CommitteeObjective links an enterprise objective to a committee.
type CommitteeObjective struct {
	CommitteeObjectiveID string `json:"committee_objective_id" gorm:"column:committee_objective_id;primaryKey"`
	ObjectiveID          string `json:"objective_id"           gorm:"column:objective_id;not null"`
	CommitteeID          string `json:"committee_id"           gorm:"column:committee_id;not null"`
	domain.BaseEntity

	Committee *Committee            `json:"committee" gorm:"foreignKey:CommitteeID"`
	Objective *EnterpriseObjective  `json:"objective" gorm:"foreignKey:ObjectiveID"`
}

func (CommitteeObjective) TableName() string { return "pms.committee_objectives" }

// ProjectAssignedWorkProduct defines a work product template for a project.
type ProjectAssignedWorkProduct struct {
	ProjectAssignedWorkProductID string    `json:"project_assigned_work_product_id" gorm:"column:project_assigned_work_product_id;primaryKey"`
	WorkProductDefinitionID      string    `json:"work_product_definition_id"       gorm:"column:work_product_definition_id"`
	Name                         string    `json:"name"                             gorm:"column:name;not null"`
	Description                  string    `json:"description"                      gorm:"column:description"`
	ProjectID                    string    `json:"project_id"                       gorm:"column:project_id;not null"`
	ReviewPeriodID               string    `json:"review_period_id"                 gorm:"column:review_period_id"`
	StartDate                    time.Time `json:"start_date"                       gorm:"column:start_date;not null"`
	EndDate                      time.Time `json:"end_date"                         gorm:"column:end_date;not null"`
	Deliverables                 string    `json:"deliverables"                     gorm:"column:deliverables"`
	domain.BaseWorkFlow

	Project      *Project                 `json:"project"       gorm:"foreignKey:ProjectID"`
	ReviewPeriod *PerformanceReviewPeriod  `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
}

func (ProjectAssignedWorkProduct) TableName() string { return "pms.project_assigned_work_products" }

// CommitteeAssignedWorkProduct defines a work product template for a committee.
type CommitteeAssignedWorkProduct struct {
	CommitteeAssignedWorkProductID string    `json:"committee_assigned_work_product_id" gorm:"column:committee_assigned_work_product_id;primaryKey"`
	WorkProductDefinitionID        string    `json:"work_product_definition_id"         gorm:"column:work_product_definition_id"`
	Name                           string    `json:"name"                               gorm:"column:name;not null"`
	Description                    string    `json:"description"                        gorm:"column:description"`
	CommitteeID                    string    `json:"committee_id"                       gorm:"column:committee_id;not null"`
	ReviewPeriodID                 string    `json:"review_period_id"                   gorm:"column:review_period_id"`
	StartDate                      time.Time `json:"start_date"                         gorm:"column:start_date;not null"`
	EndDate                        time.Time `json:"end_date"                           gorm:"column:end_date;not null"`
	Deliverables                   string    `json:"deliverables"                       gorm:"column:deliverables"`
	domain.BaseWorkFlow

	Committee    *Committee               `json:"committee"     gorm:"foreignKey:CommitteeID"`
	ReviewPeriod *PerformanceReviewPeriod  `json:"review_period" gorm:"foreignKey:ReviewPeriodID"`
}

func (CommitteeAssignedWorkProduct) TableName() string {
	return "pms.committee_assigned_work_products"
}
