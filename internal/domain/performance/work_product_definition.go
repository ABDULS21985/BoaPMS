package performance

import (
	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// WorkProductDefinition maps to pms.WorkProductDefinitions.
type WorkProductDefinition struct {
	domain.BaseEntity
	WorkProductDefinitionID string `json:"workProductDefinitionId" gorm:"column:work_product_definition_id;not null"`
	ReferenceNo             string `json:"referenceNo"             gorm:"column:reference_no"`
	Name                    string `json:"name"                    gorm:"column:name;not null"`
	Description             string `json:"description"             gorm:"column:description"`
	Deliverables            string `json:"deliverables"            gorm:"column:deliverables"`
	ObjectiveID             string `json:"objectiveId"             gorm:"column:objective_id"`
	ObjectiveLevel          string `json:"objectiveLevel"          gorm:"column:objective_level;default:'Office'"`
}

func (WorkProductDefinition) TableName() string { return "pms.WorkProductDefinitions" }

// CascadedWorkProduct maps to pms.CascadedWorkProducts.
type CascadedWorkProduct struct {
	domain.BaseWorkFlow
	CascadedWorkProductID string                   `json:"cascadedWorkProductId" gorm:"column:cascaded_work_product_id;primaryKey"`
	SmdReferenceCode      string                   `json:"smdReferenceCode"      gorm:"column:smd_reference_code"`
	Name                  string                   `json:"name"                  gorm:"column:name;not null"`
	Description           string                   `json:"description"           gorm:"column:description"`
	ObjectiveID           string                   `json:"objectiveId"           gorm:"column:objective_id;not null"`
	ObjectiveLevel        enums.ObjectiveLevel      `json:"objectiveLevel"        gorm:"column:objective_level;not null"`
	StaffJobRole          string                   `json:"staffJobRole"          gorm:"column:staff_job_role;not null"`
	ReviewPeriodID        string                   `json:"reviewPeriodId"        gorm:"column:review_period_id;not null"`
	ReviewPeriod          *PerformanceReviewPeriod  `json:"reviewPeriod,omitempty" gorm:"foreignKey:ReviewPeriodID"`
	WorkProducts          []WorkProduct            `json:"workProducts,omitempty" gorm:"foreignKey:CascadedWorkProductID"`
}

func (CascadedWorkProduct) TableName() string { return "pms.CascadedWorkProducts" }

// ObjectiveTemplateFields defines the objective upload template field names.
var ObjectiveTemplateFields = []string{
	"Strategy", "EObjName", "EObjDescription", "EObjKPI", "EObjTarget", "EObjCategory",
	"Dept", "DeptObjName", "DeptObjDescription", "DeptObjKPI", "DeptObjTarget",
	"Division", "DivObjName", "DivObjDescription", "DivObjKPI", "DivObjTarget",
	"Office", "OffObjName", "OffObjDescription", "OffObjKPI", "OffObjTarget",
}
