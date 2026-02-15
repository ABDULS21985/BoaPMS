package dto

// ===========================================================================
// Grievance Request VMs
// ===========================================================================

// CreateNewGrievanceVm is the request body for filing a new grievance.
type CreateNewGrievanceVm struct {
	GrievanceType              string `json:"grievance_type"`
	ReviewPeriodID             string `json:"review_period_id"`
	SubjectID                  string `json:"subject_id"`
	Subject                    string `json:"subject"`
	Description                string `json:"description"`
	ComplainantStaffID         string `json:"complainant_staff_id"`
	ComplainantEvidenceUpload  string `json:"complainant_evidence_upload,omitempty"`
	RespondentStaffID          string `json:"respondent_staff_id"`
}

// GrievanceVm extends the create model with resolution and response fields.
type GrievanceVm struct {
	GrievanceID                string `json:"grievance_id"`
	GrievanceType              string `json:"grievance_type"`
	ReviewPeriodID             string `json:"review_period_id"`
	SubjectID                  string `json:"subject_id"`
	Subject                    string `json:"subject"`
	Description                string `json:"description"`
	ComplainantStaffID         string `json:"complainant_staff_id"`
	ComplainantEvidenceUpload  string `json:"complainant_evidence_upload,omitempty"`
	RespondentStaffID          string `json:"respondent_staff_id"`
	RespondentComment          string `json:"respondent_comment,omitempty"`
	CurrentResolutionLevel     string `json:"current_resolution_level,omitempty"`
	CurrentMediatorStaffID     string `json:"current_mediator_staff_id,omitempty"`
	RespondentEvidenceUpload   string `json:"respondent_evidence_upload,omitempty"`
}

// GrievanceResolutionRequestVm is the request for submitting a grievance resolution.
type GrievanceResolutionRequestVm struct {
	ResolutionComment string `json:"resolution_comment"`
	Level             string `json:"level"`
	MediatorStaffID   string `json:"mediator_staff_id"`
	GrievanceID       string `json:"grievance_id"`
	EvidenceUpload    string `json:"evidence_upload,omitempty"`
}

// GrievanceResolutionVm extends the resolution request with feedback fields.
type GrievanceResolutionVm struct {
	GrievanceResolutionID string `json:"grievance_resolution_id"`
	ResolutionComment     string `json:"resolution_comment"`
	Level                 string `json:"level"`
	MediatorStaffID       string `json:"mediator_staff_id"`
	GrievanceID           string `json:"grievance_id"`
	EvidenceUpload        string `json:"evidence_upload,omitempty"`
	RespondentFeedback    string `json:"respondent_feedback,omitempty"`
	ComplainantFeedback   string `json:"complainant_feedback,omitempty"`
	ComplainantRemark     string `json:"complainant_remark,omitempty"`
	RespondentRemark      string `json:"respondent_remark,omitempty"`
}

// ===========================================================================
// Grievance Response VMs & Data Structs
// ===========================================================================

// GrievanceData holds grievance data for responses.
type GrievanceData struct {
	GrievanceID                string `json:"grievance_id"`
	GrievanceType              string `json:"grievance_type"`
	ReviewPeriodID             string `json:"review_period_id"`
	ReviewPeriodName           string `json:"review_period_name"`
	SubjectID                  string `json:"subject_id"`
	Subject                    string `json:"subject"`
	Description                string `json:"description"`
	ComplainantStaffID         string `json:"complainant_staff_id"`
	ComplainantName            string `json:"complainant_name"`
	ComplainantEvidenceUpload  string `json:"complainant_evidence_upload,omitempty"`
	RespondentStaffID          string `json:"respondent_staff_id"`
	RespondentName             string `json:"respondent_name"`
	RespondentComment          string `json:"respondent_comment,omitempty"`
	CurrentResolutionLevel     string `json:"current_resolution_level,omitempty"`
	CurrentMediatorStaffID     string `json:"current_mediator_staff_id,omitempty"`
	CurrentMediatorName        string `json:"current_mediator_name,omitempty"`
	RespondentEvidenceUpload   string `json:"respondent_evidence_upload,omitempty"`
	RecordStatus               string `json:"record_status"`
}

// GrievanceResponseVm wraps a single grievance in a standard response.
type GrievanceResponseVm struct {
	BaseAPIResponse
	Grievance GrievanceData `json:"grievance"`
}

// GrievanceListResponseVm wraps a list of grievances.
type GrievanceListResponseVm struct {
	GenericListResponseVm
	Grievances []GrievanceData `json:"grievances"`
}

// GrievanceResolutionData holds grievance resolution data for responses.
type GrievanceResolutionData struct {
	GrievanceResolutionID string `json:"grievance_resolution_id"`
	ResolutionComment     string `json:"resolution_comment"`
	Level                 string `json:"level"`
	MediatorStaffID       string `json:"mediator_staff_id"`
	MediatorName          string `json:"mediator_name"`
	GrievanceID           string `json:"grievance_id"`
	EvidenceUpload        string `json:"evidence_upload,omitempty"`
	RespondentFeedback    string `json:"respondent_feedback,omitempty"`
	ComplainantFeedback   string `json:"complainant_feedback,omitempty"`
	ComplainantRemark     string `json:"complainant_remark,omitempty"`
	RespondentRemark      string `json:"respondent_remark,omitempty"`
	RecordStatus          string `json:"record_status"`
}

// GrievanceResolutionResponseVm wraps a single grievance resolution.
type GrievanceResolutionResponseVm struct {
	BaseAPIResponse
	Resolution GrievanceResolutionData `json:"resolution"`
}

// GrievanceResolutionListResponseVm wraps a list of grievance resolutions.
type GrievanceResolutionListResponseVm struct {
	GenericListResponseVm
	Resolutions []GrievanceResolutionData `json:"resolutions"`
}

// ===========================================================================
// View Model DTOs (Vm suffix – handler-layer representations)
// ===========================================================================

// ---------------------------------------------------------------------------
// Grievance Create / Update View Models
// ---------------------------------------------------------------------------

// CreateGrievanceVm is an alias for CreateNewGrievanceVm providing a shorter name.
type CreateGrievanceVm = CreateNewGrievanceVm

// UpdateGrievanceVm is the request body for updating an existing grievance.
type UpdateGrievanceVm struct {
	GrievanceID               string `json:"grievance_id"`
	GrievanceType             string `json:"grievance_type"`
	ReviewPeriodID            string `json:"review_period_id"`
	SubjectID                 string `json:"subject_id"`
	Subject                   string `json:"subject"`
	Description               string `json:"description"`
	ComplainantStaffID        string `json:"complainant_staff_id"`
	ComplainantEvidenceUpload string `json:"complainant_evidence_upload,omitempty"`
	RespondentStaffID         string `json:"respondent_staff_id"`
	RespondentComment         string `json:"respondent_comment,omitempty"`
	RespondentEvidenceUpload  string `json:"respondent_evidence_upload,omitempty"`
}

// ---------------------------------------------------------------------------
// Grievance Resolution Create / Update View Models
// ---------------------------------------------------------------------------

// CreateGrievanceResolutionVm is the request body for creating a new grievance resolution.
type CreateGrievanceResolutionVm struct {
	ResolutionComment string `json:"resolution_comment"`
	Level             string `json:"level"`
	MediatorStaffID   string `json:"mediator_staff_id"`
	GrievanceID       string `json:"grievance_id"`
	EvidenceUpload    string `json:"evidence_upload,omitempty"`
}

// UpdateGrievanceResolutionVm is the request body for updating an existing grievance resolution.
type UpdateGrievanceResolutionVm struct {
	GrievanceResolutionID string `json:"grievance_resolution_id"`
	ResolutionComment     string `json:"resolution_comment"`
	Level                 string `json:"level"`
	MediatorStaffID       string `json:"mediator_staff_id"`
	GrievanceID           string `json:"grievance_id"`
	EvidenceUpload        string `json:"evidence_upload,omitempty"`
	RespondentFeedback    string `json:"respondent_feedback,omitempty"`
	ComplainantFeedback   string `json:"complainant_feedback,omitempty"`
	ComplainantRemark     string `json:"complainant_remark,omitempty"`
	RespondentRemark      string `json:"respondent_remark,omitempty"`
}

// ---------------------------------------------------------------------------
// Grievance Report View Model
// ---------------------------------------------------------------------------

// GrievanceReportVm represents a grievance report summary for dashboard and reporting.
type GrievanceReportVm struct {
	ReviewPeriodID         string `json:"review_period_id"`
	ReviewPeriodName       string `json:"review_period_name"`
	TotalGrievances        int    `json:"total_grievances"`
	ResolvedGrievances     int    `json:"resolved_grievances"`
	PendingGrievances      int    `json:"pending_grievances"`
	EscalatedGrievances    int    `json:"escalated_grievances"`
	AverageResolutionDays  int    `json:"average_resolution_days"`
}
