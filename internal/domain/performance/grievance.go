package performance

import (
	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// Grievance records a staff complaint during a review period.
type Grievance struct {
	GrievanceID            string               `json:"grievance_id"             gorm:"column:grievance_id;primaryKey"`
	GrievanceType          enums.GrievanceType  `json:"grievance_type"           gorm:"column:grievance_type;not null"`
	ReviewPeriodID         string               `json:"review_period_id"         gorm:"column:review_period_id;not null"`
	SubjectID              string               `json:"subject_id"               gorm:"column:subject_id;not null"`
	Subject                string               `json:"subject"                  gorm:"column:subject"`
	Description            string               `json:"description"              gorm:"column:description;not null"`
	RespondentComment      string               `json:"respondent_comment"       gorm:"column:respondent_comment"`
	CurrentResolutionLevel enums.ResolutionLevel `json:"current_resolution_level" gorm:"column:current_resolution_level;default:1"`
	CurrentMediatorStaffID string               `json:"current_mediator_staff_id" gorm:"column:current_mediator_staff_id"`
	ComplainantStaffID     string               `json:"complainant_staff_id"     gorm:"column:complainant_staff_id;not null;index"`
	ComplainantEvidenceUpload string            `json:"complainant_evidence_upload" gorm:"column:complainant_evidence_upload"`
	RespondentStaffID      string               `json:"respondent_staff_id"      gorm:"column:respondent_staff_id;not null;index"`
	RespondentEvidenceUpload string             `json:"respondent_evidence_upload" gorm:"column:respondent_evidence_upload"`
	domain.BaseEntity

	GrievanceResolutions []GrievanceResolution `json:"grievance_resolutions" gorm:"foreignKey:GrievanceID"`
}

func (Grievance) TableName() string { return "pms.grievances" }

// GrievanceResolution records a mediation attempt for a grievance.
type GrievanceResolution struct {
	GrievanceResolutionID string                  `json:"grievance_resolution_id" gorm:"column:grievance_resolution_id;primaryKey"`
	ResolutionComment     string                  `json:"resolution_comment"      gorm:"column:resolution_comment;not null"`
	ResolutionLevel       string                  `json:"resolution_level"        gorm:"column:resolution_level"`
	Level                 enums.ResolutionLevel   `json:"level"                   gorm:"column:level;default:1"`
	MediatorStaffID       string                  `json:"mediator_staff_id"       gorm:"column:mediator_staff_id;not null"`
	EvidenceUpload        string                  `json:"evidence_upload"         gorm:"column:evidence_upload"`
	RespondentFeedback    string                  `json:"respondent_feedback"     gorm:"column:respondent_feedback"`
	ComplainantFeedback   string                  `json:"complainant_feedback"    gorm:"column:complainant_feedback"`
	ComplainantRemark     enums.ResolutionRemark  `json:"complainant_remark"      gorm:"column:complainant_remark;default:1"`
	RespondentRemark      enums.ResolutionRemark  `json:"respondent_remark"       gorm:"column:respondent_remark;default:1"`
	GrievanceID           string                  `json:"grievance_id"            gorm:"column:grievance_id;not null"`
	domain.BaseEntity

	Grievance *Grievance `json:"grievance" gorm:"foreignKey:GrievanceID"`
}

func (GrievanceResolution) TableName() string { return "pms.grievance_resolutions" }
