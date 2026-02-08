package performance

import "time"

// ===========================================================================
// BaseEntityVm – DTO mirror of the .NET BaseEntity (PMS schema).
// Used for VMs extending BaseEntity rather than BaseWorkFlow.
// ===========================================================================

// BaseEntityVm carries the entity-level audit fields for DTOs that
// extend BaseEntity on the .NET side.
type BaseEntityVm struct {
	ID           int        `json:"id"`
	RecordStatus string     `json:"recordStatus"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	CreatedBy    string     `json:"createdBy"`
	UpdatedBy    string     `json:"updatedBy"`
	IsActive     bool       `json:"isActive"`
}

// ===========================================================================
// CategoryDefinitionVm (CategoryDefinitionVm.cs)
// ===========================================================================

// CategoryDefinitionVm is the read/display DTO for an objective-category
// weight definition within a review period.
type CategoryDefinitionVm struct {
	BaseWorkFlowVm
	DefinitionID            string  `json:"definitionId"`
	ObjectiveCategoryID     string  `json:"objectiveCategoryId"     validate:"required"`
	Weight                  float64 `json:"weight"                  validate:"required"`
	MaxNoObjectives         int     `json:"maxNoObjectives"`
	MaxNoWorkProduct        int     `json:"maxNoWorkProduct"`
	MaxPoints               int     `json:"maxPoints"`
	IsCompulsory            bool    `json:"isCompulsory"`
	EnforceWorkProductLimit bool    `json:"enforceWorkProductLimit"`
	Description             string  `json:"description"             validate:"required"`
	GradeGroupID            int     `json:"gradeGroupId"`
	GroupName               string  `json:"groupName"`
	Category                string  `json:"category"`
}

// CreateCategoryDefinitionVm is the payload for creating a new category
// definition.
type CreateCategoryDefinitionVm struct {
	ObjectiveCategoryID     string  `json:"objectiveCategoryId"     validate:"required"`
	Weight                  float64 `json:"weight"                  validate:"required"`
	MaxNoObjectives         int     `json:"maxNoObjectives"         validate:"required"`
	MaxNoWorkProduct        int     `json:"maxNoWorkProduct"        validate:"required"`
	MaxPoints               int     `json:"maxPoints"               validate:"required"`
	IsCompulsory            bool    `json:"isCompulsory"            validate:"required"`
	EnforceWorkProductLimit bool    `json:"enforceWorkProductLimit" validate:"required"`
	Description             string  `json:"description"             validate:"required"`
	GradeGroupID            int     `json:"gradeGroupId"            validate:"required"`
}

// ===========================================================================
// PerformanceReviewPeriodVm (PerformanceReviewPeriodVm.cs)
// ===========================================================================

// PerformanceReviewPeriodVm is the read/write DTO for a review period.
type PerformanceReviewPeriodVm struct {
	BaseWorkFlowVm
	PeriodID                   string    `json:"periodId"`
	Year                       int       `json:"year"`
	Range                      int       `json:"range"                      validate:"required"`
	RangeValue                 int       `json:"rangeValue"                 validate:"required"`
	Name                       string    `json:"name"                       validate:"required"`
	ShortName                  string    `json:"shortName"                  validate:"required"`
	Description                string    `json:"description"                validate:"required"`
	StartDate                  time.Time `json:"startDate"`
	EndDate                    time.Time `json:"endDate"`
	MaxPoints                  float64   `json:"maxPoints"`
	MinNoOfObjectives          int       `json:"minNoOfObjectives"`
	MaxNoOfObjectives          int       `json:"maxNoOfObjectives"          validate:"required"`
	HasExtension               bool      `json:"hasExtension"`
	StrategyID                 string    `json:"strategyId"`
	AllowObjectivePlanning     bool      `json:"allowObjectivePlanning"`
	AllowWorkProductPlanning   bool      `json:"allowWorkProductPlanning"`
	AllowWorkProductEvaluation bool      `json:"allowWorkProductEvaluation"`
}

// ===========================================================================
// StrategyVm (StrategyVm.cs)
// ===========================================================================

// StrategyVm is the read/display DTO for a corporate strategy.
type StrategyVm struct {
	BaseWorkFlowVm
	StrategyID  string    `json:"strategyId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BankYearID  int       `json:"bankYearId"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	ImageFile   string    `json:"imageFile"`
}

// StrategyListVm is the paged response for strategy queries.
type StrategyListVm struct {
	BaseAPIResponse
	BankStrategies []StrategyVm `json:"bankStrategies"`
	TotalRecord    int          `json:"totalRecord"`
}

// SearchStrategyVm carries filters for searching strategies.
type SearchStrategyVm struct {
	BasePagedData
	CategoryID   *int   `json:"categoryId"`
	SearchString string `json:"searchString"`
	IsApproved   *bool  `json:"isApproved"`
	IsRejected   *bool  `json:"isRejected"`
	IsTechnical  *bool  `json:"isTechnical"`
}

// CreateNewStrategyVm is the payload for creating a new strategy.
type CreateNewStrategyVm struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BankYearID  int       `json:"bankYearId"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	ImageFile   string    `json:"imageFile"`
	IsActive    *bool     `json:"isActive"`
}

// StrategicThemeVm is the read/display DTO for a strategic theme.
type StrategicThemeVm struct {
	BaseWorkFlowVm
	StrategicThemeID string `json:"strategicThemeId"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	StrategyID       string `json:"strategyId"`
	StrategyName     string `json:"strategyName"`
	ImageFile        string `json:"imageFile"`
}

// CreateStrategicThemeVm is the payload for creating a strategic theme.
type CreateStrategicThemeVm struct {
	BaseWorkFlowVm
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description" validate:"required"`
	StrategyID  string `json:"strategyId"  validate:"required"`
	ImageFile   string `json:"imageFile"`
}

// StrategicThemeListVm is the paged response for strategic theme queries.
type StrategicThemeListVm struct {
	BaseAPIResponse
	StrategicThemes []StrategicThemeVm `json:"strategicThemes"`
	TotalRecord     int                `json:"totalRecord"`
}

// ===========================================================================
// PmsConfigurationVm (PmsConfigurationVm.cs)
// ===========================================================================

// PmsConfigurationVm is the read/display DTO for a PMS configuration entry.
type PmsConfigurationVm struct {
	BaseEntityVm
	PmsConfigurationID string `json:"pmsConfigurationId"`
	Name               string `json:"name"`
	Value              string `json:"value"`
	Type               string `json:"type"`
	IsEncrypted        bool   `json:"isEncrypted"`
}

// PmsConfigurationResponseVm wraps a single PMS configuration in an API
// response.
type PmsConfigurationResponseVm struct {
	BaseAPIResponse
	Data       *PmsConfigurationVm `json:"data"`
	StatusCode string              `json:"statusCode"`
	ActionCall string              `json:"actionCall"`
}

// ListPmsConfigurationResponseVm wraps a list of PMS configurations in an
// API response.
type ListPmsConfigurationResponseVm struct {
	BaseAPIResponse
	Data          []PmsConfigurationVm `json:"data"`
	TotalSettings int                  `json:"totalSettings"`
	StatusCode    string               `json:"statusCode"`
	ActionCall    string               `json:"actionCall"`
}

// ===========================================================================
// EvaluationOptionVm (EvaluationOptionVm.cs)
// ===========================================================================

// EvaluationOptionVm is the read/display DTO for an evaluation option.
type EvaluationOptionVm struct {
	BaseWorkFlowVm
	EvaluationOptionID string  `json:"evaluationOptionId"`
	Name               string  `json:"name"               validate:"required"`
	Description        string  `json:"description"`
	RecordStatus       int     `json:"recordStatus"       validate:"required"`
	Score              float64 `json:"score"              validate:"required"`
	EvaluationType     int     `json:"evaluationType"     validate:"required"`
}

// UploadEvaluationOptionVm carries a single row from a bulk-upload for
// evaluation options.
type UploadEvaluationOptionVm struct {
	BaseWorkFlowVm
	EvaluationOptionID string  `json:"evaluationOptionId"`
	Name               string  `json:"name"               validate:"required"`
	Description        string  `json:"description"`
	RecordStatus       int     `json:"recordStatus"`
	Score              float64 `json:"score"`
	IsValidRecord      bool    `json:"isValidRecord"`
	IsSuccess          *bool   `json:"isSuccess"`
	Message            string  `json:"message"`
	IsProcessed        *bool   `json:"isProcessed"`
	IsSelected         bool    `json:"isSelected"`
	EvaluationType     string  `json:"evaluationType"`
}

// SearchEvaluationOptionVm carries filters for searching evaluation options.
type SearchEvaluationOptionVm struct {
	BasePagedData
	SearchString    string  `json:"searchString"`
	Question        string  `json:"question"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	OptionStatement string  `json:"optionStatement"`
	Score           float64 `json:"score"`
}

// ===========================================================================
// PmsCompetencyRequestVm (PmsCompetencyVm.cs)
// ===========================================================================

// PmsCompetencyRequestVm is the request DTO for a PMS competency.
type PmsCompetencyRequestVm struct {
	BaseWorkFlowVm
	PmsCompetencyID  string `json:"pmsCompetencyId"  validate:"required"`
	Name             string `json:"name"             validate:"required"`
	Description      string `json:"description"`
	ObjectCategoryID string `json:"objectCategoryId" validate:"required"`
	RecordStatus     int    `json:"recordStatus"     validate:"required"`
}

// CreateNewPmsCompetencyVm is the payload for creating a new PMS competency.
type CreateNewPmsCompetencyVm struct {
	PmsCompetencyID  string `json:"pmsCompetencyId"  validate:"required"`
	Name             string `json:"name"             validate:"required"`
	Description      string `json:"description"`
	ObjectCategoryID string `json:"objectCategoryId" validate:"required"`
}

// ===========================================================================
// SettingVm (SettingVm.cs)
// ===========================================================================

// SettingVm is the read/display DTO for a PMS setting.
type SettingVm struct {
	BaseEntityVm
	SettingID   string `json:"settingId"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	IsEncrypted bool   `json:"isEncrypted"`
}

// SettingResponseDetail is the detail DTO for a setting response.
type SettingResponseDetail struct {
	BaseEntityVm
	SettingID   string `json:"settingId"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	IsEncrypted bool   `json:"isEncrypted"`
}

// SettingResponse wraps a single setting in an API response.
type SettingResponse struct {
	BaseAPIResponse
	Data       *SettingResponseDetail `json:"data"`
	StatusCode string                 `json:"statusCode"`
	ActionCall string                 `json:"actionCall"`
}

// ListSettingResponse wraps a list of settings in an API response.
type ListSettingResponse struct {
	BaseAPIResponse
	Data          []SettingResponseDetail `json:"data"`
	TotalSettings int                     `json:"totalSettings"`
	StatusCode    string                  `json:"statusCode"`
	ActionCall    string                  `json:"actionCall"`
}

// CommonExceptionResponse is a standard error envelope.
type CommonExceptionResponse struct {
	BaseAPIResponse
	StatusCode string `json:"statusCode"`
	ActionCall string `json:"actionCall"`
}

// ===========================================================================
// PerformanceGradeVm (PerformanceGradeVm.cs)
// ===========================================================================

// PerformanceGradeVm is the read/display DTO for a performance grade enum.
type PerformanceGradeVm struct {
	PerformanceGrade int    `json:"performanceGrade"`
	Description      string `json:"description"`
}

// ===========================================================================
// CompetencyReviewPeriodVm (CompetencyMgt ReviewPeriodVm.cs)
// ===========================================================================

// CompetencyReviewPeriodVm is the DTO for a competency review period.
// Mirrors CompetencyApp.Models.CompetencyMgt.ReviewPeriodVm.
type CompetencyReviewPeriodVm struct {
	BaseAuditVm
	ReviewPeriodID int        `json:"reviewPeriodId"`
	BankYearID     int        `json:"bankYearId"     validate:"required"`
	Name           string     `json:"name"`
	StartDate      time.Time  `json:"startDate"`
	EndDate        time.Time  `json:"endDate"`
	BankYearName   string     `json:"bankYearName"`
	ApprovedBy     string     `json:"approvedBy"`
	DateApproved   *time.Time `json:"dateApproved"`
	IsApproved     bool       `json:"isApproved"`
}
