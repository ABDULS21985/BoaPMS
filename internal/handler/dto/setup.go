package dto

import "time"

// ---------------------------------------------------------------------------
// Strategy VMs
// ---------------------------------------------------------------------------

// CreateNewStrategyVm is the request body for creating a new strategy.
type CreateNewStrategyVm struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	SmdReferenceCode string    `json:"smd_reference_code"`
	BankYearID       string    `json:"bank_year_id"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}

// StrategyVm extends CreateNewStrategyVm with the persisted identifier.
type StrategyVm struct {
	StrategyID       string    `json:"strategy_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	SmdReferenceCode string    `json:"smd_reference_code"`
	BankYearID       string    `json:"bank_year_id"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}

// CreateStrategicThemeVm is the request body for creating a strategic theme.
type CreateStrategicThemeVm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StrategyID  string `json:"strategy_id"`
}

// StrategicThemeVm extends CreateStrategicThemeVm with its identifier.
type StrategicThemeVm struct {
	StrategicThemeID string `json:"strategic_theme_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	StrategyID       string `json:"strategy_id"`
}

// ---------------------------------------------------------------------------
// Enterprise Objective VMs
// ---------------------------------------------------------------------------

// CreateEnterpriseObjectiveVm is the request body for creating an enterprise objective.
type CreateEnterpriseObjectiveVm struct {
	Name                          string `json:"name"`
	Description                   string `json:"description"`
	Kpi                           string `json:"kpi"`
	Target                        string `json:"target"`
	SmdReferenceCode              string `json:"smd_reference_code"`
	EnterpriseObjectivesCategoryID string `json:"enterprise_objectives_category_id"`
	StrategyID                    string `json:"strategy_id"`
	StrategicThemeID              string `json:"strategic_theme_id"`
	Type                          string `json:"type"`
}

// EnterpriseObjectiveVm extends the create model with the persisted identifier.
type EnterpriseObjectiveVm struct {
	EnterpriseObjectiveID          string `json:"enterprise_objective_id"`
	Name                           string `json:"name"`
	Description                    string `json:"description"`
	Kpi                            string `json:"kpi"`
	Target                         string `json:"target"`
	SmdReferenceCode               string `json:"smd_reference_code"`
	EnterpriseObjectivesCategoryID string `json:"enterprise_objectives_category_id"`
	StrategyID                     string `json:"strategy_id"`
	StrategicThemeID               string `json:"strategic_theme_id"`
	Type                           string `json:"type"`
}

// ---------------------------------------------------------------------------
// Department Objective VMs
// ---------------------------------------------------------------------------

// CreateDepartmentObjectiveVm is the request body for creating a department objective.
type CreateDepartmentObjectiveVm struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Kpi                   string `json:"kpi"`
	Target                string `json:"target"`
	SmdReferenceCode      string `json:"smd_reference_code"`
	DepartmentID          string `json:"department_id"`
	EnterpriseObjectiveID string `json:"enterprise_objective_id"`
}

// DepartmentObjectiveVm extends the create model with the persisted identifier.
type DepartmentObjectiveVm struct {
	DepartmentObjectiveID string `json:"department_objective_id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Kpi                   string `json:"kpi"`
	Target                string `json:"target"`
	SmdReferenceCode      string `json:"smd_reference_code"`
	DepartmentID          string `json:"department_id"`
	EnterpriseObjectiveID string `json:"enterprise_objective_id"`
}

// ---------------------------------------------------------------------------
// Division Objective VMs
// ---------------------------------------------------------------------------

// CreateDivisionObjectiveVm is the request body for creating a division objective.
type CreateDivisionObjectiveVm struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Kpi                     string `json:"kpi"`
	Target                  string `json:"target"`
	SmdReferenceCode        string `json:"smd_reference_code"`
	DivisionID              string `json:"division_id"`
	DepartmentObjectiveID   string `json:"department_objective_id"`
}

// DivisionObjectiveVm extends the create model with the persisted identifier.
type DivisionObjectiveVm struct {
	DivisionObjectiveID   string `json:"division_objective_id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Kpi                   string `json:"kpi"`
	Target                string `json:"target"`
	SmdReferenceCode      string `json:"smd_reference_code"`
	DivisionID            string `json:"division_id"`
	DepartmentObjectiveID string `json:"department_objective_id"`
}

// ---------------------------------------------------------------------------
// Office Objective VMs
// ---------------------------------------------------------------------------

// CreateOfficeObjectiveVm is the request body for creating an office objective.
type CreateOfficeObjectiveVm struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Kpi                   string `json:"kpi"`
	Target                string `json:"target"`
	SmdReferenceCode      string `json:"smd_reference_code"`
	OfficeID              string `json:"office_id"`
	DivisionObjectiveID   string `json:"division_objective_id"`
	JobGradeGroupID       string `json:"job_grade_group_id"`
}

// OfficeObjectiveVm extends the create model with the persisted identifier.
type OfficeObjectiveVm struct {
	OfficeObjectiveID   string `json:"office_objective_id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Kpi                 string `json:"kpi"`
	Target              string `json:"target"`
	SmdReferenceCode    string `json:"smd_reference_code"`
	OfficeID            string `json:"office_id"`
	DivisionObjectiveID string `json:"division_objective_id"`
	JobGradeGroupID     string `json:"job_grade_group_id"`
}

// ---------------------------------------------------------------------------
// Objective Category VMs
// ---------------------------------------------------------------------------

// CreateObjectiveCategoryVm is the request body for creating an objective category.
type CreateObjectiveCategoryVm struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ObjectiveCategoryVm extends the create model with the persisted identifier.
type ObjectiveCategoryVm struct {
	ObjectiveCategoryID string `json:"objective_category_id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
}

// ---------------------------------------------------------------------------
// Category Definition VMs
// ---------------------------------------------------------------------------

// CreateCategoryDefinitionVm is the request body for creating a category definition.
type CreateCategoryDefinitionVm struct {
	ObjectiveCategoryID    string  `json:"objective_category_id"`
	ReviewPeriodID         string  `json:"review_period_id"`
	Weight                 float64 `json:"weight"`
	MaxNoObjectives        int     `json:"max_no_objectives"`
	MaxNoWorkProduct       int     `json:"max_no_work_product"`
	MaxPoints              float64 `json:"max_points"`
	IsCompulsory           bool    `json:"is_compulsory"`
	EnforceWorkProductLimit bool   `json:"enforce_work_product_limit"`
	Description            string  `json:"description"`
	GradeGroupID           string  `json:"grade_group_id"`
}

// CategoryDefinitionVm extends the create model with the persisted identifier.
type CategoryDefinitionVm struct {
	DefinitionID           string  `json:"definition_id"`
	ObjectiveCategoryID    string  `json:"objective_category_id"`
	ReviewPeriodID         string  `json:"review_period_id"`
	Weight                 float64 `json:"weight"`
	MaxNoObjectives        int     `json:"max_no_objectives"`
	MaxNoWorkProduct       int     `json:"max_no_work_product"`
	MaxPoints              float64 `json:"max_points"`
	IsCompulsory           bool    `json:"is_compulsory"`
	EnforceWorkProductLimit bool   `json:"enforce_work_product_limit"`
	Description            string  `json:"description"`
	GradeGroupID           string  `json:"grade_group_id"`
}

// CreateCategoryDefinitionRequestVm is the request for creating a category definition within a review period context.
type CreateCategoryDefinitionRequestVm struct {
	ObjectiveCategoryID    string  `json:"objective_category_id"`
	ReviewPeriodID         string  `json:"review_period_id"`
	Weight                 float64 `json:"weight"`
	MaxNoObjectives        int     `json:"max_no_objectives"`
	MaxNoWorkProduct       int     `json:"max_no_work_product"`
	MaxPoints              float64 `json:"max_points"`
	IsCompulsory           bool    `json:"is_compulsory"`
	EnforceWorkProductLimit bool   `json:"enforce_work_product_limit"`
	Description            string  `json:"description"`
	GradeGroupID           string  `json:"grade_group_id"`
}

// CategoryDefinitionRequestVm extends the create request model with a definition identifier.
type CategoryDefinitionRequestVm struct {
	DefinitionID           string  `json:"definition_id"`
	ObjectiveCategoryID    string  `json:"objective_category_id"`
	ReviewPeriodID         string  `json:"review_period_id"`
	Weight                 float64 `json:"weight"`
	MaxNoObjectives        int     `json:"max_no_objectives"`
	MaxNoWorkProduct       int     `json:"max_no_work_product"`
	MaxPoints              float64 `json:"max_points"`
	IsCompulsory           bool    `json:"is_compulsory"`
	EnforceWorkProductLimit bool   `json:"enforce_work_product_limit"`
	Description            string  `json:"description"`
	GradeGroupID           string  `json:"grade_group_id"`
}

// ---------------------------------------------------------------------------
// Evaluation Option VMs
// ---------------------------------------------------------------------------

// EvaluationOptionVm represents an evaluation option used in scoring.
type EvaluationOptionVm struct {
	EvaluationOptionID string  `json:"evaluation_option_id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	RecordStatus       string  `json:"record_status"`
	Score              float64 `json:"score"`
	EvaluationType     string  `json:"evaluation_type"`
}

// ---------------------------------------------------------------------------
// Feedback Questionnaire VMs
// ---------------------------------------------------------------------------

// FeedbackQuestionaireVm represents a feedback questionnaire item.
type FeedbackQuestionaireVm struct {
	FeedbackQuestionaireID string `json:"feedback_questionaire_id"`
	Question               string `json:"question"`
	Description            string `json:"description"`
	PmsCompetencyID        string `json:"pms_competency_id"`
}

// FeedbackQuestionaireOptionVm represents an answer option for a questionnaire.
type FeedbackQuestionaireOptionVm struct {
	FeedbackQuestionaireOptionID string  `json:"feedback_questionaire_option_id"`
	OptionStatement              string  `json:"option_statement"`
	Description                  string  `json:"description"`
	Score                        float64 `json:"score"`
	QuestionID                   string  `json:"question_id"`
}

// ---------------------------------------------------------------------------
// PMS Competency VMs
// ---------------------------------------------------------------------------

// CreateNewPmsCompetencyVm is the request body for creating a PMS competency.
type CreateNewPmsCompetencyVm struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	ObjectCategoryID  string `json:"object_category_id"`
}

// PmsCompetencyRequestVm extends the create model with the persisted identifier.
type PmsCompetencyRequestVm struct {
	PmsCompetencyID  string `json:"pms_competency_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ObjectCategoryID string `json:"object_category_id"`
}

// ---------------------------------------------------------------------------
// Work Product Definition VMs
// ---------------------------------------------------------------------------

// WorkProductDefinitionVm represents a work product definition.
type WorkProductDefinitionVm struct {
	WorkProductDefinitionID string `json:"work_product_definition_id"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	ObjectiveID             string `json:"objective_id"`
	ObjectiveLevel          string `json:"objective_level"`
}

// WorkProductDefinitionRequestVm is the request model for a work product definition.
type WorkProductDefinitionRequestVm struct {
	WorkProductDefinitionID string `json:"work_product_definition_id"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	ObjectiveID             string `json:"objective_id"`
	ObjectiveLevel          string `json:"objective_level"`
}

// ObjectiveWorkProductDefinitionRequest carries the objective reference for work product definitions.
type ObjectiveWorkProductDefinitionRequest struct {
	ObjectiveID    string `json:"objective_id"`
	ObjectiveLevel string `json:"objective_level"`
}

// ---------------------------------------------------------------------------
// Setting VMs
// ---------------------------------------------------------------------------

// AddSettingRequestModel is the request body for creating a new setting.
type AddSettingRequestModel struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	IsEncrypted bool   `json:"is_encrypted"`
}

// SettingRequestModel extends the add model with the persisted identifier.
type SettingRequestModel struct {
	SettingID   string `json:"setting_id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	IsEncrypted bool   `json:"is_encrypted"`
}

// SettingData holds the data payload for a single setting response.
type SettingData struct {
	SettingID   string `json:"setting_id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	IsEncrypted bool   `json:"is_encrypted"`
}

// SettingResponse extends BaseAPIResponse with a single setting record.
type SettingResponse struct {
	BaseAPIResponse
	Setting SettingData `json:"setting"`
}

// ListSettingResponse extends GenericListResponseVm with a list of settings.
type ListSettingResponse struct {
	GenericListResponseVm
	Settings []SettingData `json:"settings"`
}

// ---------------------------------------------------------------------------
// PMS Configuration VMs
// ---------------------------------------------------------------------------

// AddPmsConfigurationRequestModel is the request body for creating a PMS configuration.
type AddPmsConfigurationRequestModel struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	IsEncrypted bool   `json:"is_encrypted"`
}

// PmsConfigurationRequestModel extends the add model with the persisted identifier.
type PmsConfigurationRequestModel struct {
	PmsConfigurationID string `json:"pms_configuration_id"`
	Name               string `json:"name"`
	Value              string `json:"value"`
	Type               string `json:"type"`
	IsEncrypted        bool   `json:"is_encrypted"`
}

// PmsConfigurationData holds the data payload for a single PMS configuration.
type PmsConfigurationData struct {
	PmsConfigurationID string `json:"pms_configuration_id"`
	Name               string `json:"name"`
	Value              string `json:"value"`
	Type               string `json:"type"`
	IsEncrypted        bool   `json:"is_encrypted"`
}

// PmsConfigurationResponseVm extends BaseAPIResponse with a single configuration.
type PmsConfigurationResponseVm struct {
	BaseAPIResponse
	Configuration PmsConfigurationData `json:"configuration"`
}

// ListPmsConfigurationResponseVm extends GenericListResponseVm with a list of configurations.
type ListPmsConfigurationResponseVm struct {
	GenericListResponseVm
	Configurations []PmsConfigurationData `json:"configurations"`
}

// ---------------------------------------------------------------------------
// Consolidated Objective VMs
// ---------------------------------------------------------------------------

// ConsolidatedObjectiveVm represents an objective across all organisational levels.
type ConsolidatedObjectiveVm struct {
	ObjectiveID    string `json:"objective_id"`
	ObjectiveLevel string `json:"objective_level"`
	Name           string `json:"name"`
	Kpi            string `json:"kpi"`
	Target         string `json:"target"`
	Description    string `json:"description"`
	CategoryID     string `json:"category_id"`
	StrategyID     string `json:"strategy_id"`
	IsActive       bool   `json:"is_active"`
}

// SearchObjectiveVm extends BasePagedData with search filters for objectives.
type SearchObjectiveVm struct {
	BasePagedData
	Search         string `json:"search"`
	ObjectiveLevel string `json:"objective_level"`
	CategoryID     string `json:"category_id"`
}

// ---------------------------------------------------------------------------
// Cascaded Objective Upload VM
// ---------------------------------------------------------------------------

// CascadedObjectiveUploadVm carries the data for a bulk cascaded objective upload.
type CascadedObjectiveUploadVm struct {
	EnterpriseObjectiveName string `json:"enterprise_objective_name"`
	EnterpriseObjectiveKpi  string `json:"enterprise_objective_kpi"`
	EnterpriseObjectiveTarget string `json:"enterprise_objective_target"`
	EnterpriseObjectiveDescription string `json:"enterprise_objective_description"`
	DepartmentObjectiveName string `json:"department_objective_name"`
	DepartmentObjectiveKpi  string `json:"department_objective_kpi"`
	DepartmentObjectiveTarget string `json:"department_objective_target"`
	DepartmentObjectiveDescription string `json:"department_objective_description"`
	DivisionObjectiveName   string `json:"division_objective_name"`
	DivisionObjectiveKpi    string `json:"division_objective_kpi"`
	DivisionObjectiveTarget string `json:"division_objective_target"`
	DivisionObjectiveDescription string `json:"division_objective_description"`
	OfficeObjectiveName     string `json:"office_objective_name"`
	OfficeObjectiveKpi      string `json:"office_objective_kpi"`
	OfficeObjectiveTarget   string `json:"office_objective_target"`
	OfficeObjectiveDescription string `json:"office_objective_description"`
	DepartmentID            string `json:"department_id"`
	DivisionID              string `json:"division_id"`
	OfficeID                string `json:"office_id"`
	StrategyID              string `json:"strategy_id"`
	StrategicThemeID        string `json:"strategic_theme_id"`
	CategoryID              string `json:"category_id"`
	JobGradeGroupID         string `json:"job_grade_group_id"`
	SmdReferenceCode        string `json:"smd_reference_code"`
}
