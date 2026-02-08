package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// WorkProduct represents a deliverable assigned to a staff member.
type WorkProduct struct {
	WorkProductID          string              `json:"work_product_id"           gorm:"column:work_product_id;primaryKey"`
	Name                   string              `json:"name"                      gorm:"column:name;not null"`
	Description            string              `json:"description"               gorm:"column:description"`
	MaxPoint               float64             `json:"max_point"                 gorm:"column:max_point;type:decimal(18,2);not null"`
	WorkProductType        enums.WorkProductType `json:"work_product_type"       gorm:"column:work_product_type"`
	IsSelfCreated          bool                `json:"is_self_created"           gorm:"column:is_self_created;default:false"`
	StaffID                string              `json:"staff_id"                  gorm:"column:staff_id;not null;index"`
	AcceptanceComment      string              `json:"acceptance_comment"        gorm:"column:acceptance_comment"`
	StartDate              time.Time           `json:"start_date"                gorm:"column:start_date;not null"`
	EndDate                time.Time           `json:"end_date"                  gorm:"column:end_date;not null"`
	Deliverables           string              `json:"deliverables"              gorm:"column:deliverables"`
	FinalScore             float64             `json:"final_score"               gorm:"column:final_score;type:decimal(18,2)"`
	NoReturned             int                 `json:"no_returned"               gorm:"column:no_returned;default:0"`
	CompletionDate         *time.Time          `json:"completion_date"           gorm:"column:completion_date"`
	ApproverComment        string              `json:"approver_comment"          gorm:"column:approver_comment"`
	ReEvaluationReInitiated bool               `json:"re_evaluation_re_initiated" gorm:"column:re_evaluation_re_initiated;default:false"`
	Remark                 string              `json:"remark"                    gorm:"column:remark"`
	domain.BaseWorkFlow

	WorkProductTasks                []WorkProductTask                `json:"work_product_tasks"                 gorm:"foreignKey:WorkProductID"`
	OperationalObjectiveWorkProducts []OperationalObjectiveWorkProduct `json:"operational_objective_work_products" gorm:"foreignKey:WorkProductID"`
	ProjectWorkProducts             []ProjectWorkProduct             `json:"project_work_products"              gorm:"foreignKey:WorkProductID"`
	CommitteeWorkProducts           []CommitteeWorkProduct           `json:"committee_work_products"            gorm:"foreignKey:WorkProductID"`
}

func (WorkProduct) TableName() string { return "pms.work_products" }

// WorkProductTask is a sub-task of a WorkProduct.
type WorkProductTask struct {
	WorkProductTaskID string    `json:"work_product_task_id" gorm:"column:work_product_task_id;primaryKey"`
	Name              string    `json:"name"                 gorm:"column:name;not null"`
	Description       string    `json:"description"          gorm:"column:description"`
	StartDate         time.Time `json:"start_date"           gorm:"column:start_date;not null"`
	EndDate           time.Time `json:"end_date"             gorm:"column:end_date;not null"`
	CompletionDate    *time.Time `json:"completion_date"     gorm:"column:completion_date"`
	WorkProductID     string    `json:"work_product_id"      gorm:"column:work_product_id;not null"`
	domain.BaseEntity

	WorkProduct *WorkProduct `json:"work_product" gorm:"foreignKey:WorkProductID"`
}

func (WorkProductTask) TableName() string { return "pms.work_product_tasks" }

// WorkProductEvaluation stores timeliness/quality/output scores for a work product.
type WorkProductEvaluation struct {
	WorkProductEvaluationID      string  `json:"work_product_evaluation_id"       gorm:"column:work_product_evaluation_id;primaryKey"`
	WorkProductID                string  `json:"work_product_id"                  gorm:"column:work_product_id;not null"`
	Timeliness                   float64 `json:"timeliness"                       gorm:"column:timeliness;type:decimal(18,2);not null"`
	TimelinessEvaluationOptionID string  `json:"timeliness_evaluation_option_id"  gorm:"column:timeliness_evaluation_option_id"`
	Quality                      float64 `json:"quality"                          gorm:"column:quality;type:decimal(18,2);not null"`
	QualityEvaluationOptionID    string  `json:"quality_evaluation_option_id"     gorm:"column:quality_evaluation_option_id"`
	Output                       float64 `json:"output"                           gorm:"column:output;type:decimal(18,2);not null"`
	OutputEvaluationOptionID     string  `json:"output_evaluation_option_id"      gorm:"column:output_evaluation_option_id"`
	Outcome                      float64 `json:"outcome"                          gorm:"column:outcome;type:decimal(18,2)"`
	EvaluatorStaffID             string  `json:"evaluator_staff_id"               gorm:"column:evaluator_staff_id"`
	IsReEvaluated                bool    `json:"is_re_evaluated"                  gorm:"column:is_re_evaluated;default:false"`
	domain.BaseEntity

	WorkProduct              *WorkProduct      `json:"work_product"               gorm:"foreignKey:WorkProductID"`
	TimelinessEvaluationOption *EvaluationOption `json:"timeliness_evaluation_option" gorm:"foreignKey:TimelinessEvaluationOptionID"`
	QualityEvaluationOption  *EvaluationOption `json:"quality_evaluation_option"  gorm:"foreignKey:QualityEvaluationOptionID"`
	OutputEvaluationOption   *EvaluationOption `json:"output_evaluation_option"   gorm:"foreignKey:OutputEvaluationOptionID"`
}

func (WorkProductEvaluation) TableName() string { return "pms.work_product_evaluations" }

// EvaluationOption provides scoring options for evaluation dimensions.
type EvaluationOption struct {
	EvaluationOptionID string              `json:"evaluation_option_id" gorm:"column:evaluation_option_id;primaryKey"`
	Name               string              `json:"name"                 gorm:"column:name;not null"`
	Description        string              `json:"description"          gorm:"column:description"`
	RecordStatus       enums.Status        `json:"record_status"        gorm:"column:record_status;not null"`
	Score              float64             `json:"score"                gorm:"column:score;type:decimal(18,2);not null"`
	EvaluationType     enums.EvaluationType `json:"evaluation_type"     gorm:"column:evaluation_type;not null"`
	domain.BaseWorkFlow
}

func (EvaluationOption) TableName() string { return "pms.evaluation_options" }

// OperationalObjectiveWorkProduct links a work product to a planned objective.
type OperationalObjectiveWorkProduct struct {
	OperationalObjectiveWorkProductID string `json:"id"                    gorm:"column:operational_objective_work_product_id;primaryKey"`
	WorkProductID                     string `json:"work_product_id"       gorm:"column:work_product_id;not null"`
	WorkProductDefinitionID           string `json:"work_product_definition_id" gorm:"column:work_product_definition_id"`
	PlannedObjectiveID                string `json:"planned_objective_id"  gorm:"column:planned_objective_id;not null"`
	domain.BaseEntity

	WorkProduct      *WorkProduct                            `json:"work_product"      gorm:"foreignKey:WorkProductID"`
	PlannedObjective *ReviewPeriodIndividualPlannedObjective  `json:"planned_objective" gorm:"foreignKey:PlannedObjectiveID"`
}

func (OperationalObjectiveWorkProduct) TableName() string {
	return "pms.operational_objective_work_products"
}
