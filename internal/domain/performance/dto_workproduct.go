package performance

import "time"

// ---------------------------------------------------------------------------
// WorkProductVm mirrors the .NET WorkProductVm (extends BaseWorkFlow).
// ---------------------------------------------------------------------------

// WorkProductVm is the request/response DTO for a work product.
type WorkProductVm struct {
	BaseAuditVm
	WorkProductID     string     `json:"workProductId"`
	Name              string     `json:"name"              validate:"required"`
	Description       string     `json:"description"`
	MaxPoint          float64    `json:"maxPoint"          validate:"required"`
	WorkProductTypeID string     `json:"workProductTypeId"`
	IsSelfCreated     bool       `json:"isSelfCreated"`
	StaffID           string     `json:"staffId"           validate:"required"`
	AcceptanceComment string     `json:"acceptanceComment"`
	StartDate         time.Time  `json:"startDate"         validate:"required"`
	EndDate           time.Time  `json:"endDate"           validate:"required"`
	Deliverables      string     `json:"deliverables"`
	FinalScore        int        `json:"finalScore"`
	NoReturned        int        `json:"noReturned"`
	CompletionDate    time.Time  `json:"completionDate"`
	ApproverComment   string     `json:"approverComment"`
	RecordStatus      int        `json:"recordStatus"      validate:"required"`
	ProjectID         string     `json:"projectId"`
	CommitteeID       string     `json:"committeeId"`
	PlannedObjectiveID string   `json:"plannedObjectiveId"`

	// Workflow fields (from BaseWorkFlow).
	ApprovedBy      string     `json:"approvedBy"`
	DateApproved    time.Time  `json:"dateApproved"`
	IsApproved      bool       `json:"isApproved"`
	IsRejected      bool       `json:"isRejected"`
	RejectedBy      string     `json:"rejectedBy"`
	RejectionReason string     `json:"rejectionReason"`
	DateRejected    time.Time  `json:"dateRejected"`
}

// ---------------------------------------------------------------------------
// WorkProductEvaluationVm mirrors the .NET WorkProductEvaluationVm
// (extends BaseEntity).
// ---------------------------------------------------------------------------

// WorkProductEvaluationVm is the DTO for work product evaluation scores.
type WorkProductEvaluationVm struct {
	BaseAuditVm
	WorkProductEvaluationID      string     `json:"workProductEvaluationId"`
	WorkProductID                string     `json:"workProductId"`
	Timeliness                   float64    `json:"timeliness"`
	TimelinessEvaluationOptionID string     `json:"timelinessEvaluationOptionId"`
	Quality                      float64    `json:"quality"`
	QualityEvaluationOptionID    string     `json:"qualityEvaluationOptionId"`
	Output                       float64    `json:"output"`
	OutputEvaluationOptionID     string     `json:"outputEvaluationOptionId"`
	Outcome                      float64    `json:"outcome"`
	OutcomeEvaluationOptionID    string     `json:"outcomeEvaluationOptionId"`
	EvaluatorStaffID             string     `json:"evaluatorStaffId"`
	IsReEvaluated                bool       `json:"isReEvaluated"`
	CreatedAt                    *time.Time `json:"createdAt"`
}

// ---------------------------------------------------------------------------
// WorkProductTaskVm mirrors the .NET WorkProductTaskVm (extends BaseEntity).
// ---------------------------------------------------------------------------

// WorkProductTaskVm is the DTO for a sub-task of a work product.
type WorkProductTaskVm struct {
	BaseAuditVm
	WorkProductTaskID string `json:"workProductTaskId"`
	Name              string `json:"name"              validate:"required"`
	Description       string `json:"description"`
	WorkProductID     string `json:"workProductId"     validate:"required"`
	RecordStatus      int    `json:"recordStatus"`
}

// ---------------------------------------------------------------------------
// WorkProductDefinitionVm mirrors the .NET WorkProductDefinitionVm
// (extends BaseEntity).
// ---------------------------------------------------------------------------

// WorkProductDefinitionVm is the DTO for a work product definition.
type WorkProductDefinitionVm struct {
	BaseAuditVm
	WorkProductDefinitionID string `json:"workProductDefinitionId"`
	ReferenceNo             string `json:"referenceNo"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Deliverables            string `json:"deliverables"`
	ObjectiveID             string `json:"objectiveId"`
	ObjectiveLevel          string `json:"objectiveLevel"`
	SBUName                 string `json:"sbuName"`
	ObjectiveName           string `json:"objectiveName"`
	Grade                   string `json:"grade"`
}

// WorkProductDefinitionRequestVm is the create/update request for a work
// product definition.
type WorkProductDefinitionRequestVm struct {
	ReferenceNo    string `json:"referenceNo"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Deliverables   string `json:"deliverables"`
	ObjectiveID    string `json:"objectiveId"`
	ObjectiveLevel string `json:"objectiveLevel"`
	SBUName        string `json:"sbuName"`
	ObjectiveName  string `json:"objectiveName"`
}

// ObjectiveWorkProductDefinitionRequestVm carries objective context when
// managing work product definitions.
type ObjectiveWorkProductDefinitionRequestVm struct {
	ObjectiveID    string `json:"objectiveId"`
	ObjectiveLevel int    `json:"objectiveLevel"`
	SBUName        string `json:"sbuName"`
	ObjectiveName  string `json:"objectiveName"`
}

// SearchWorkProductDefinitionVm is the search/filter DTO for work product
// definitions.  It embeds BasePagedData for pagination.
type SearchWorkProductDefinitionVm struct {
	BasePagedData
	SearchString   string `json:"searchString"`
	Name           string `json:"name"`
	ReferenceNo    string `json:"referenceNo"`
	Description    string `json:"description"`
	ObjectiveID    string `json:"objectiveId"`
	ObjectiveLevel string `json:"objectiveLevel"`
	DepartmentID   *int   `json:"departmentId"`
	Enterprise     string `json:"enterprise"`
	Department     string `json:"department"`
	DivisionID     *int   `json:"divisionId"`
	Division       string `json:"division"`
	OfficeID       *int   `json:"officeId"`
	Office         string `json:"office"`
	JobRoleID      *int   `json:"jobRoleId"`
}

// UploadWorkProductDefinitionRequestVm is the DTO for bulk-uploading work
// product definitions.
type UploadWorkProductDefinitionRequestVm struct {
	BaseAuditVm
	ReferenceNo    string `json:"referenceNo"`
	SBUName        string `json:"sbuName"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Deliverables   string `json:"deliverables"`
	ObjectiveID    string `json:"objectiveId"`
	ObjectiveName  string `json:"objectiveName"`
	ObjectiveLevel string `json:"objectiveLevel"`
	JobGradeGroup  string `json:"jobGradeGroup"`

	IsValidRecord bool   `json:"isValidRecord"`
	IsSuccess     *bool  `json:"isSuccess"`
	Message       string `json:"message"`
	IsProcessed   *bool  `json:"isProcessed"`
	IsSelected    bool   `json:"isSelected"`

	// Workflow fields (from BaseWorkFlow).
	ApprovedBy      string    `json:"approvedBy"`
	DateApproved    time.Time `json:"dateApproved"`
	IsApproved      bool      `json:"isApproved"`
	IsRejected      bool      `json:"isRejected"`
	RejectedBy      string    `json:"rejectedBy"`
	RejectionReason string    `json:"rejectionReason"`
	DateRejected    time.Time `json:"dateRejected"`
}

// DownloadWorkProductTemplateRequestVm carries parameters for downloading a
// work product definition template.
type DownloadWorkProductTemplateRequestVm struct {
	ObjectiveName  []string `json:"objectiveName"`
	ObjectiveLevel string   `json:"objectiveLevel"`
	SBUName        string   `json:"sbuName"`
}

// ---------------------------------------------------------------------------
// CascadedWorkProductVm mirrors the .NET CascadedWorkProduct model
// (extends BaseWorkFlow).
// ---------------------------------------------------------------------------

// CascadedWorkProductVm is the DTO for a cascaded work product.
type CascadedWorkProductVm struct {
	BaseAuditVm
	CascadedWorkProductID string `json:"cascadedWorkProductId" validate:"required"`
	SmdReferenceCode      string `json:"smdReferenceCode"`
	Name                  string `json:"name"                  validate:"required"`
	Description           string `json:"description"`
	ObjectiveID           string `json:"objectiveId"           validate:"required"`
	ObjectiveLevel        int    `json:"objectiveLevel"        validate:"required"`
	StaffJobRole          string `json:"staffJobRole"          validate:"required"`
	ReviewPeriodID        string `json:"reviewPeriodId"        validate:"required"`

	// Workflow fields (from BaseWorkFlow).
	ApprovedBy      string    `json:"approvedBy"`
	DateApproved    time.Time `json:"dateApproved"`
	IsApproved      bool      `json:"isApproved"`
	IsRejected      bool      `json:"isRejected"`
	RejectedBy      string    `json:"rejectedBy"`
	RejectionReason string    `json:"rejectionReason"`
	DateRejected    time.Time `json:"dateRejected"`
}

// ---------------------------------------------------------------------------
// CustomWorkProductVm mirrors the .NET CustomWorkProductVm.
// ---------------------------------------------------------------------------

// CustomWorkProductVm is the frontend-facing DTO for a work product with
// task and evaluation details.
type CustomWorkProductVm struct {
	WorkProductID           string                  `json:"workProductId"`
	PlannedObjectiveID      string                  `json:"plannedObjectiveId"`
	ObjectiveID             string                  `json:"objectivedId"`
	StaffID                 string                  `json:"staffId"`
	CreatedBy               string                  `json:"createdBy"`
	WorkProductDefinitionID string                  `json:"workProductDefinitionId"`
	Name                    string                  `json:"name"`
	Description             string                  `json:"description"`
	RejectionReason         string                  `json:"rejectionReason"`
	MaxPoint                float64                 `json:"maxPoint"`
	WorkProductTypeID       string                  `json:"workProductTypeId"`
	IsSelfCreated           bool                    `json:"isSelfCreated"`
	StartDate               time.Time               `json:"startDate"`
	EndDate                 time.Time               `json:"endDate"`
	Deliverables            string                  `json:"deliverables"`
	FinalScore              float64                 `json:"finalScore"`
	NoReturned              int                     `json:"noReturned"`
	CompletionDate          time.Time               `json:"completionDate"`
	RecordStatus            int                     `json:"recordStatus"`
	WorkProductStatus       string                  `json:"workProductStatus"`
	ObjectiveOwnID          string                  `json:"objectiveId"`
	WorkProductTasks        []CustomWorkProductTaskVm `json:"workProductTasks"`
	IsReEvaluated           bool                    `json:"isReEvaluated"`
	CreatedAt               *time.Time              `json:"createdAt"`
	ReEvaluationReInitiated bool                    `json:"reEvaluationReInitiated"`
	Kpi                     string                  `json:"kpi"`
	Target                  string                  `json:"target"`
	EnterpriseObjective     *EnterpriseObjectiveDataVm `json:"enterpriseObjective"`
}

// ---------------------------------------------------------------------------
// CustomWorkProductForProjectVm mirrors the .NET
// CustomWorkProductForProjectVm.
// ---------------------------------------------------------------------------

// CustomWorkProductForProjectVm is the project-specific work product DTO.
type CustomWorkProductForProjectVm struct {
	WorkProductID           string                  `json:"workProductId"`
	PlannedObjectiveID      string                  `json:"plannedObjectiveId"`
	ObjectiveID             string                  `json:"objectivedId"`
	StaffID                 string                  `json:"staffId"`
	CreatedBy               string                  `json:"createdBy"`
	WorkProductDefinitionID string                  `json:"workProductDefinitionId"`
	Name                    string                  `json:"name"`
	Description             string                  `json:"description"`
	RejectionReason         string                  `json:"rejectionReason"`
	FinalScore              float64                 `json:"finalScore"`
	MaxPoint                float64                 `json:"maxPoint"`
	WorkProductTypeID       string                  `json:"workProductTypeId"`
	IsSelfCreated           bool                    `json:"isSelfCreated"`
	StartDate               time.Time               `json:"startDate"`
	EndDate                 time.Time               `json:"endDate"`
	Deliverables            string                  `json:"deliverables"`
	NoReturned              int                     `json:"noReturned"`
	CompletionDate          time.Time               `json:"completionDate"`
	RecordStatus            int                     `json:"recordStatus"`
	WorkProductStatus       string                  `json:"workProductStatus"`
	ObjectiveOwnID          string                  `json:"objectiveId"`
	WorkProductTasks        []CustomWorkProductTaskVm `json:"workProductTasks"`
}

// ---------------------------------------------------------------------------
// CustomWorkProductTaskVm mirrors the .NET CutomWorkProductTaskVm
// (extends BaseEntity).
// ---------------------------------------------------------------------------

// CustomWorkProductTaskVm is the frontend-facing DTO for a work product task.
type CustomWorkProductTaskVm struct {
	BaseAuditVm
	WorkProductTaskID string    `json:"workProductTaskId"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	WorkProductID     string    `json:"workProductId"`
	EndDate           time.Time `json:"endDate"`
	StartDate         time.Time `json:"startDate"`
	CompletionDate    time.Time `json:"completionDate"`
	RecordStatus      int       `json:"recordStatus"`
}

// ---------------------------------------------------------------------------
// EnterpriseObjectiveDataVm mirrors the .NET EnterpriseObjectiveDataVm.
// ---------------------------------------------------------------------------

// EnterpriseObjectiveDataVm carries enterprise objective details used in
// work product responses.
type EnterpriseObjectiveDataVm struct {
	EnterpriseObjectiveID          string    `json:"enterpriseObjectiveId"`
	PeriodObjectiveID              string    `json:"periodObjectiveId"`
	Name                           string    `json:"name"`
	Description                    string    `json:"description"`
	Kpi                            string    `json:"kpi"`
	Target                         string    `json:"target"`
	EnterpriseObjectivesCategoryID string    `json:"enterpriseObjectivesCategoryId"`
	StrategyID                     string    `json:"strategyId"`
	HasEvaluation                  bool      `json:"hasEvaluation"`
	OutcomeScore                   float64   `json:"outcomeScore"`
	TotalOutcomeScore              float64   `json:"totalOutcomeScore"`
	Evaluator                      string    `json:"evaluator"`
	EvaluationDate                 time.Time `json:"evaluationDate"`
}

// ---------------------------------------------------------------------------
// ProjectWorkProductVm mirrors the .NET ProjectWorkProductVm.
// ---------------------------------------------------------------------------

// ProjectWorkProductVm is the DTO for a work product linked to a project or
// committee.
type ProjectWorkProductVm struct {
	Name                            string                    `json:"name"`
	ProjectAssignedWorkProductID    string                    `json:"projectAssignedWorkProductId"`
	CommitteeAssignedWorkProductID  string                    `json:"committeeAssignedWorkProductId"`
	Description                     string                    `json:"description"`
	WorkProductDefinitionID         string                    `json:"workProductDefinitionId"`
	WorkProductID                   string                    `json:"workProductId"`
	WorkProductType                 int                       `json:"workProductType"`
	ProjectWorkProductID            string                    `json:"projectWorkProductId"`
	WorkProduct                     *CustomProjectWorkProductVm `json:"workProduct"`
	ProjectID                       string                    `json:"projectId"`
	ProjectManager                  string                    `json:"projectManager"`
	CommitteeID                     string                    `json:"committeeId"`
	Chairperson                     string                    `json:"chairperson"`
	StartDate                       time.Time                 `json:"startDate"`
	EndDate                         time.Time                 `json:"endDate"`
	Deliverables                    string                    `json:"deliverables"`
	RejectionReason                 string                    `json:"rejectionReason"`
	ReferenceNo                     string                    `json:"referenceNo"`
	ObjectiveID                     string                    `json:"objectiveId"`
	ObjectiveLevel                  string                    `json:"objectiveLevel"`
	AssignedEvaluator               string                    `json:"assignedEvaluator"`
	RecordStatus                    int                       `json:"recordStatus"`
}

// ---------------------------------------------------------------------------
// CustomProjectWorkProductVm mirrors the .NET CustomProjectWorkProduct.
// ---------------------------------------------------------------------------

// CustomProjectWorkProductVm is the project-specific work product summary
// DTO.
type CustomProjectWorkProductVm struct {
	WorkProductID    string                    `json:"workProductId"`
	Name             string                    `json:"name"`
	Description      string                    `json:"description"`
	WorkProductType  int                       `json:"workProductType"`
	FinalScore       float64                   `json:"finalScore"`
	MaxPoint         float64                   `json:"maxPoint"`
	StaffID          string                    `json:"staffId"`
	StartDate        time.Time                 `json:"startDate"`
	EndDate          time.Time                 `json:"endDate"`
	CompletionDate   time.Time                 `json:"completionDate"`
	RecordStatus     int                       `json:"recordStatus"`
	WorkProductTasks []CustomWorkProductTaskVm `json:"workProductTasks"`
}

// ---------------------------------------------------------------------------
// DashboardPendingActions mirrors the .NET DashboardPendingActions.
// ---------------------------------------------------------------------------

// DashboardPendingActionsVm represents a single pending action on the
// dashboard.
type DashboardPendingActionsVm struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// ---------------------------------------------------------------------------
// WorkProductTypeVm mirrors the .NET WorkProductTypeVm enum wrapper.
// ---------------------------------------------------------------------------

// WorkProductTypeVm is the DTO for exposing a work product type with a
// human-readable description.
type WorkProductTypeVm struct {
	WorkProductType int    `json:"workProductType"`
	Description     string `json:"description"`
}

// ---------------------------------------------------------------------------
// CustomAddWorkProductVm mirrors the .NET CustomAddWorkProductVm.
// ---------------------------------------------------------------------------

// CustomAddWorkProductVm is the create request DTO for adding a work
// product.
type CustomAddWorkProductVm struct {
	WorkProductID      string    `json:"workProductId"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	WorkProductTypeRaw int       `json:"workProductType"`
	StaffID            string    `json:"staffId"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
	PlannedObjectiveID string    `json:"plannedObjectiveId"`
	AcceptanceComment  string    `json:"acceptanceComment"`
	Remark             string    `json:"remark"`
	Attachment         string    `json:"attachment"`
}
