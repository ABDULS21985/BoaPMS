package performance

import (
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ---------------------------------------------------------------------------
// Grievance DTOs  (source: GrievanceVms/GreivanceVm.cs)
// ---------------------------------------------------------------------------

// GrievanceVm is the read/response DTO for a grievance record.
type GrievanceVm struct {
	BaseAuditVm
	GrievanceID              string                    `json:"grievanceId" validate:"required"`
	GrievanceType            int                       `json:"grievanceType" validate:"required"`
	Description              string                    `json:"description" validate:"required"`
	RespondentComment        string                    `json:"respondentComment"`
	SubjectID                string                    `json:"subjectId" validate:"required"`
	Subject                  string                    `json:"subject"`
	ReviewPeriodID           string                    `json:"reviewPeriodId" validate:"required"`
	ComplainantStaffID       string                    `json:"complainantStaffId" validate:"required"`
	ComplainantStaff         string                    `json:"complainantStaff"`
	ComplainantEvidenceUpload string                   `json:"complainantEvidenceUpload"`
	CurrentResolutionLevel   enums.ResolutionLevel     `json:"currentResolutionLevel"`
	CurrentMediatorStaffID   string                    `json:"currentMediatorStaffId"`
	CurrentMediatorStaff     string                    `json:"currentMediatorStaff"`
	RespondentStaffID        string                    `json:"respondentStaffId" validate:"required"`
	RespondentStaff          string                    `json:"respondentStaff"`
	RespondentEvidenceUpload string                    `json:"respondentEvidenceUpload"`
	FinalResolution          *GrievanceResolutionVm    `json:"finalResolution"`
	IsResolved               bool                      `json:"isResolved"`
	GrievanceResolutions     []GrievanceResolutionVm   `json:"grievanceResolutions"`
}

// CreateNewGrievanceVm is the request DTO for creating a new grievance.
type CreateNewGrievanceVm struct {
	GrievanceType            int    `json:"grievanceType" validate:"required"`
	Description              string `json:"description" validate:"required"`
	ReviewPeriodID           string `json:"reviewPeriodId" validate:"required"`
	Subject                  string `json:"subject"`
	SubjectID                string `json:"subjectId" validate:"required"`
	ComplainantStaffID       string `json:"complainantStaffId" validate:"required"`
	ComplainantEvidenceUpload string `json:"complainantEvidenceUpload"`
}

// ---------------------------------------------------------------------------
// Grievance Resolution DTOs  (source: GrievanceVms/GrievanceResolutionVm.cs)
// ---------------------------------------------------------------------------

// GrievanceResolutionVm is the read/response DTO for a resolution attempt.
type GrievanceResolutionVm struct {
	BaseAuditVm
	GrievanceResolutionID string                `json:"grievanceResolutionId"`
	ResolutionComment     string                `json:"resolutionComment" validate:"required"`
	ResolutionLevel       string                `json:"resolutionLevel"`
	Level                 enums.ResolutionLevel `json:"level" validate:"required"`
	MediatorStaffID       string                `json:"mediatorStaffId" validate:"required"`
	MediatorStaff         string                `json:"mediatorStaff"`
	EvidenceUpload        string                `json:"evidenceUpload"`
	RespondentFeedback    string                `json:"respondentFeedback"`
	ComplainantFeedback   string                `json:"complainantFeedback"`
	ComplainantRemark     enums.ResolutionRemark `json:"complainantRemark"`
	RespondentRemark      enums.ResolutionRemark `json:"respondentRemark"`
	GrievanceID           string                `json:"grievanceId" validate:"required"`
	ActionStatus          *enums.Status         `json:"actionStatus"`
}

// IsAcceptedByComplainant is a computed property (NotMapped in C#).
func (v GrievanceResolutionVm) IsAcceptedByComplainant() bool {
	return v.ComplainantRemark == enums.ResolutionRemarkAccepted
}

// IsAcceptedByRespondent is a computed property (NotMapped in C#).
func (v GrievanceResolutionVm) IsAcceptedByRespondent() bool {
	return v.RespondentRemark == enums.ResolutionRemarkAccepted
}

// GrievanceResolutionRequestVm is the request DTO for creating/updating a resolution.
type GrievanceResolutionRequestVm struct {
	ResolutionComment string                `json:"resolutionComment" validate:"required"`
	Level             enums.ResolutionLevel `json:"level" validate:"required"`
	MediatorStaffID   string                `json:"mediatorStaffId" validate:"required"`
	EvidenceUpload    string                `json:"evidenceUpload"`
	GrievanceID       string                `json:"grievanceId" validate:"required"`
	ComplainantRemark enums.ResolutionRemark `json:"complainantRemark"`
	RespondentRemark  enums.ResolutionRemark `json:"respondentRemark"`
}

// ---------------------------------------------------------------------------
// Grievance Type enum DTO  (source: EnumVms/GrievanceTypeVm.cs)
// ---------------------------------------------------------------------------

// GrievanceTypeVm pairs a GrievanceType enum value with its display text.
type GrievanceTypeVm struct {
	GrievanceType enums.GrievanceType `json:"grievanceType"`
	Description   string              `json:"description"`
}
