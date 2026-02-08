package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ---------------------------------------------------------------------------
// Base DTO for ViewModels that mirror the .NET BaseWorkFlow model class.
// ---------------------------------------------------------------------------

// BaseWorkFlowVm carries the workflow-relevant fields exposed in VMs that
// extend BaseWorkFlow on the .NET side.
type BaseWorkFlowVm struct {
	ID              int        `json:"id"`
	RecordStatus    string     `json:"recordStatus"`
	CreatedAt       *time.Time `json:"createdAt"`
	Status          string     `json:"status"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	CreatedBy       string     `json:"createdBy"`
	UpdatedBy       string     `json:"updatedBy"`
	IsActive        bool       `json:"isActive"`
	ApprovedBy      string     `json:"approvedBy"`
	DateApproved    *time.Time `json:"dateApproved"`
	IsApproved      bool       `json:"isApproved"`
	IsRejected      bool       `json:"isRejected"`
	RejectedBy      string     `json:"rejectedBy"`
	RejectionReason string     `json:"rejectionReason"`
	DateRejected    *time.Time `json:"dateRejected"`
}

// ObjectiveBaseVm carries the shared fields from .NET ObjectiveBase.
type ObjectiveBaseVm struct {
	BaseWorkFlowVm
	ObjectiveID      string                `json:"objectiveId"`
	SBUName          string                `json:"sbuName"`
	Name             string                `json:"name" validate:"required"`
	SmdReferenceCode string                `json:"smdReferenceCode"`
	Description      string                `json:"description"`
	ObjectiveLevel   string                `json:"objectiveLevel"`
	Kpi              string                `json:"kpi" validate:"required"`
	Target           string                `json:"target"`
	WorkProducts     []WorkProductDefinition `json:"workProducts"`
}

// ---------------------------------------------------------------------------
// ObjectiveVm – mirrors .NET ObjectiveVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// ObjectiveVm represents a cascaded objective view across department,
// division, and office levels.
type ObjectiveVm struct {
	BaseWorkFlowVm
	DepartmentObjectiveID          string `json:"departmentObjectiveId" validate:"required"`
	DepartmentID                   int    `json:"departmentId"`
	DepartmentName                 string `json:"departmentName"`
	DepartmentObjectiveName        string `json:"departmentObjectiveName"`
	DepartmentObjectiveDescription string `json:"departmentObjectiveDescription"`
	DepartmentObjectiveKPI         string `json:"departmentObjectiveKPI"`
	DivisionObjectiveID            string `json:"divisionObjectiveId" validate:"required"`
	DivisionID                     int    `json:"divisionId"`
	DivisionName                   string `json:"divisionName"`
	DivisionObjectiveName          string `json:"divisionObjectiveName"`
	DivisionObjectiveDescription   string `json:"divisionObjectiveDescription"`
	DivisionObjectiveKPI           string `json:"divisionObjectiveKPI"`
	OfficeObjectiveID              string `json:"officeObjectiveId" validate:"required"`
	OfficeID                       int    `json:"officeId"`
	OfficeName                     string `json:"officeName"`
	OfficeObjectiveName            string `json:"officeObjectiveName"`
	OfficeObjectiveDescription     string `json:"officeObjectiveDescription"`
	OfficeObjectiveKPI             string `json:"officeObjectiveKPI"`
	JobGradeGroupID                int    `json:"jobGradeGroupId"`
	IsSelected                     bool   `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// SearchObjectiveVm – mirrors .NET SearchObjectiveVm : BasePagedData
// ---------------------------------------------------------------------------

// SearchObjectiveVm carries search/filter criteria with paging.
type SearchObjectiveVm struct {
	BasePagedData
	DepartmentID    *int               `json:"departmentId"`
	DivisionID      *int               `json:"divisionId"`
	OfficeID        *int               `json:"officeId"`
	JobRoleID       *int               `json:"jobRoleId"`
	SearchString    string             `json:"searchString"`
	IsApproved      bool               `json:"isApproved"`
	IsTechnical     bool               `json:"isTechnical"`
	TargetType      enums.ObjectiveLevel `json:"targetType"`
	TargetTypeRaw   string             `json:"_targetType"`
	TargetReference string             `json:"targetReference"`
	Status          string             `json:"status"`
}

// ---------------------------------------------------------------------------
// PagedObjectiveVm – mirrors .NET PagedObjectiveVm : BaseAPIResponse
// ---------------------------------------------------------------------------

// PagedObjectiveVm wraps a page of ObjectiveVm results.
type PagedObjectiveVm struct {
	BaseAPIResponse
	Objectives   []ObjectiveVm `json:"objectivess"`
	TotalRecords int           `json:"totalRecords"`
}

// ---------------------------------------------------------------------------
// SaveObjectiveVm – mirrors .NET SaveObjectiveVm : BaseAuditVm
// ---------------------------------------------------------------------------

// SaveObjectiveVm is used when saving/creating an objective assignment.
type SaveObjectiveVm struct {
	BaseAuditVm
	JobRoleCompetencyID      int                        `json:"jobRoleCompetencyId"`
	DepartmentID             *int                       `json:"departmentId"`
	DivisionID               *int                       `json:"divisionId"`
	JobRoleID                int                        `json:"jobRoleId" validate:"required"`
	OfficeID                 *int                       `json:"officeId" validate:"required"`
	CompetencyGroup          string                     `json:"competencyGroup"`
	JobRoleName              string                     `json:"jobRoleName"`
	OfficeName               string                     `json:"officeName"`
	JobRoleCompetencyRatings []JobRoleCompetencyRatingVm `json:"jobRoleCompetencyRatings"`
}

// ---------------------------------------------------------------------------
// JobRoleCompetencyRatingVm – mirrors .NET JobRoleCompetencyRating
// (defined in CompetencyMgtVm but referenced by SaveObjectiveVm)
// ---------------------------------------------------------------------------

// JobRoleCompetencyRatingVm represents a single competency-rating pair.
type JobRoleCompetencyRatingVm struct {
	ID             int    `json:"id"`
	CompetencyID   int    `json:"competencyId" validate:"required"`
	RatingID       int    `json:"ratingId" validate:"required"`
	CompetencyName string `json:"competencyName"`
	RatingName     string `json:"ratingName"`
}

// ---------------------------------------------------------------------------
// ObjectiveRatingVm – mirrors .NET ObjectiveRating
// ---------------------------------------------------------------------------

// ObjectiveRatingVm represents a required rating for an objective competency.
type ObjectiveRatingVm struct {
	ID             int    `json:"id"`
	CompetencyID   int    `json:"competencyId" validate:"required"`
	RatingID       int    `json:"ratingId" validate:"required"`
	CompetencyName string `json:"competencyName"`
	RatingName     string `json:"ratingName"`
}

// ---------------------------------------------------------------------------
// UploadObjectiveVm – mirrors .NET UploadObjectiveVm : BaseAuditVm
// ---------------------------------------------------------------------------

// UploadObjectiveVm carries a single row from an objective bulk-upload file.
type UploadObjectiveVm struct {
	BaseAuditVm
	ID                  int    `json:"id"`
	JobRoleCompetencyID int    `json:"jobRoleCompetencyId"`
	DepartmentID        int    `json:"departmentId"`
	DivisionID          int    `json:"divisionId"`
	JobRoleID           int    `json:"jobRoleId"`
	RatingID            int    `json:"ratingId" validate:"required"`
	OfficeID            int    `json:"officeId" validate:"required"`
	CompetencyID        int    `json:"competencyId" validate:"required"`
	CompetencyName      string `json:"competencyName"`
	JobRoleName         string `json:"jobRoleName"`
	RatingName          string `json:"ratingName"`
	OfficeName          string `json:"officeName"`
	DivisionName        string `json:"divisionName"`
	DepartmentName      string `json:"departmentName"`
	IsValidRecord       bool   `json:"isValidRecord"`
	IsSuccess           *bool  `json:"isSuccess"`
	Message             string `json:"message"`
	IsProcessed         *bool  `json:"isProcessed"`
	IsSelected          bool   `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// UploadCascadedObjectiveVm – mirrors .NET UploadCascadedObjectiveVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// UploadCascadedObjectiveVm carries a row from a cascaded-objective upload.
type UploadCascadedObjectiveVm struct {
	BaseWorkFlowVm
	StrategyID              string `json:"strategyId"`
	StrategyName            string `json:"strategyName"`
	EnterpriseObjectiveID   string `json:"enterpriseObjectiveId"`
	EObjName                string `json:"eObjName"`
	EObjDesc                string `json:"eObjDesc"`
	EObjKPI                 string `json:"eObjKPI"`
	EObjTarget              string `json:"eObjTarget"`
	EObjCategory            string `json:"eObjCategory"`
	DepartmentID            int    `json:"departmentId"`
	Dept                    string `json:"dept"`
	DepartmentObjectiveID   string `json:"departmentObjectiveId"`
	DeptObjName             string `json:"deptObjName"`
	DeptObjDesc             string `json:"deptObjDesc"`
	DeptObjKPI              string `json:"deptObjKPI"`
	DeptObjTarget           string `json:"deptObjTarget"`
	DivisionID              int    `json:"divisionId"`
	Division                string `json:"division"`
	DivisionObjectiveID     string `json:"divisionObjectiveId"`
	DivObjName              string `json:"divObjName"`
	DivObjDesc              string `json:"divObjDesc"`
	DivObjKPI               string `json:"divObjKPI"`
	DivObjTarget            string `json:"divObjTarget"`
	OfficeID                int    `json:"officeId"`
	Office                  string `json:"office"`
	OfficeObjectiveID       string `json:"officeObjectiveId"`
	OffObjName              string `json:"offObjName"`
	OffObjDesc              string `json:"offObjDesc"`
	OffObjKPI               string `json:"offObjKPI"`
	OffObjTarget            string `json:"offObjTarget"`
	StrategicThemeID        string `json:"strategicThemeId"`
	StrategicThemeName      string `json:"strategicThemeName"`
	SBUCode                 string `json:"sbuCode"`
	SBUName                 string `json:"sbuName"`
	SBULevel                string `json:"sbuLevel"`
	ParentObjectiveName     string `json:"parentObjectiveName"`
	ObjectiveName           string `json:"objectiveName"`
	ObjectiveDesc           string `json:"objectiveDesc"`
	ObjectiveKPI            string `json:"objectiveKPI"`
	ObjectiveTarget         string `json:"objectiveTarget"`
	WorkProductName         string `json:"workProductName"`
	WorkProductDescription  string `json:"workProductDescription"`
	WorkProductDeliverable  string `json:"workProductDeliverable"`
	JobGradeGroup           string `json:"jobGradeGroup"`
	JobGradeGroupID         int    `json:"jobGradeGroupId"`
	IsValidRecord           bool   `json:"isValidRecord"`
	IsSuccess               *bool  `json:"isSuccess"`
	Message                 string `json:"message"`
	IsProcessed             *bool  `json:"isProcessed"`
	IsSelected              bool   `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// UploadPMSObjectiveVm – mirrors .NET UploadPMSObjectiveVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// UploadPMSObjectiveVm carries a row from a PMS objective upload.
type UploadPMSObjectiveVm struct {
	BaseWorkFlowVm
	StrategicThemeID       string `json:"strategicThemeId"`
	StrategicThemeName     string `json:"strategicThemeName"`
	SBUCode                string `json:"sbuCode"`
	SBUName                string `json:"sbuName"`
	SBULevel               string `json:"sbuLevel"`
	ParentObjectiveName    string `json:"parentObjectiveName"`
	ObjectiveName          string `json:"objectiveName"`
	ObjectiveDesc          string `json:"objectiveDesc"`
	ObjectiveKPI           string `json:"objectiveKPI"`
	ObjectiveTarget        string `json:"objectiveTarget"`
	WorkProductName        string `json:"workProductName"`
	WorkProductDescription string `json:"workProductDescription"`
	WorkProductDeliverable string `json:"workProductDeliverable"`
	JobGradeGroup          string `json:"jobGradeGroup"`
	IsValidRecord          bool   `json:"isValidRecord"`
	IsSuccess              *bool  `json:"isSuccess"`
	Message                string `json:"message"`
	IsProcessed            *bool  `json:"isProcessed"`
	IsSelected             bool   `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// DownloadObjectiveTemplateVm – mirrors .NET DownloadObjectiveTemplate
// ---------------------------------------------------------------------------

// DownloadObjectiveTemplateVm is the shape of each row in the objective
// download template.
type DownloadObjectiveTemplateVm struct {
	StrategyID         string `json:"strategyId"`
	StrategyName       string `json:"strategyName"`
	StrategicThemeName string `json:"strategicThemeName"`
	EObjName           string `json:"eObjName"`
	EObjDesc           string `json:"eObjDesc"`
	EObjKPI            string `json:"eObjKPI"`
	EObjTarget         string `json:"eObjTarget"`
	EObjCategory       string `json:"eObjCategory"`
	Dept               string `json:"dept"`
	DeptObjName        string `json:"deptObjName"`
	DeptObjDesc        string `json:"deptObjDesc"`
	DeptObjKPI         string `json:"deptObjKPI"`
	DeptObjTarget      string `json:"deptObjTarget"`
	Division           string `json:"division"`
	DivObjName         string `json:"divObjName"`
	DivObjDesc         string `json:"divObjDesc"`
	DivObjKPI          string `json:"divObjKPI"`
	DivObjTarget       string `json:"divObjTarget"`
	Office             string `json:"office"`
	OffObjName         string `json:"offObjName"`
	OffObjDesc         string `json:"offObjDesc"`
	OffObjKPI          string `json:"offObjKPI"`
	OffObjTarget       string `json:"offObjTarget"`
	JobGradeGroup      string `json:"jobGradeGroup"`
	ObjectiveLevel     string `json:"objectiveLevel"`
	SBUName            string `json:"sbuName"`
}

// ---------------------------------------------------------------------------
// EnterpriseObjectiveVm – mirrors .NET EnterpriseObjectiveVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// EnterpriseObjectiveVmDTO is the DTO for an enterprise-level objective.
type EnterpriseObjectiveVmDTO struct {
	BaseWorkFlowVm
	EnterpriseObjectiveID          string `json:"enterpriseObjectiveId" validate:"required"`
	Name                           string `json:"name" validate:"required"`
	Description                    string `json:"description"`
	Kpi                            string `json:"kpi" validate:"required"`
	Target                         string `json:"target" validate:"required"`
	EnterpriseObjectivesCategoryID string `json:"enterpriseObjectivesCategoryId" validate:"required"`
	StrategyID                     string `json:"strategyId" validate:"required"`
}

// CreateEnterpriseObjectiveVm is the DTO for creating an enterprise objective.
type CreateEnterpriseObjectiveVm struct {
	Name                           string `json:"name" validate:"required"`
	Description                    string `json:"description"`
	Kpi                            string `json:"kpi" validate:"required"`
	Target                         string `json:"target" validate:"required"`
	EnterpriseObjectivesCategoryID string `json:"enterpriseObjectivesCategoryId" validate:"required"`
	StrategyID                     string `json:"strategyId" validate:"required"`
}

// ConsolidatedObjectiveVm extends ObjectiveBaseVm with a selection flag.
type ConsolidatedObjectiveVm struct {
	ObjectiveBaseVm
	IsSelected bool `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// DepartmentObjectiveVmDTO – mirrors .NET DepartmentObjectiveVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// DepartmentObjectiveVmDTO is the DTO for a department-level objective.
type DepartmentObjectiveVmDTO struct {
	BaseWorkFlowVm
	DepartmentObjectiveID  string  `json:"departmentObjectiveId"`
	Name                   string  `json:"name" validate:"required"`
	Description            string  `json:"description"`
	Kpi                    string  `json:"kpi"`
	Target                 string  `json:"target"`
	DepartmentID           int     `json:"departmentId"`
	EnterpriseObjectiveID  string  `json:"enterpriseObjectiveId"`
	SBUName                *string `json:"sbuName"`
	WorkProductName        *string `json:"workProductName"`
	WorkProductDescription *string `json:"workProductDescription"`
	WorkProductDeliverable *string `json:"workProductDeliverable"`
	JobGradeGroup          *string `json:"jobGradeGroup"`
}

// CreateDepartmentObjectiveVm is the DTO for creating a department objective.
type CreateDepartmentObjectiveVm struct {
	Name                   string  `json:"name" validate:"required"`
	Description            string  `json:"description"`
	Kpi                    string  `json:"kpi"`
	Target                 string  `json:"target"`
	DepartmentID           int     `json:"departmentId"`
	EnterpriseObjectiveID  string  `json:"enterpriseObjectiveId"`
	SBUName                *string `json:"sbuName"`
	WorkProductName        *string `json:"workProductName"`
	WorkProductDescription *string `json:"workProductDescription"`
	WorkProductDeliverable *string `json:"workProductDeliverable"`
	JobGradeGroup          *string `json:"jobGradeGroup"`
}

// ---------------------------------------------------------------------------
// DivisionObjectiveVmDTO – mirrors .NET DivisionObjectiveVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// DivisionObjectiveVmDTO is the DTO for a division-level objective.
type DivisionObjectiveVmDTO struct {
	BaseWorkFlowVm
	DivisionObjectiveID    string  `json:"divisionObjectiveId"`
	Name                   string  `json:"name" validate:"required"`
	Description            string  `json:"description"`
	Kpi                    string  `json:"kpi"`
	Target                 string  `json:"target"`
	DivisionID             int     `json:"divisionId" validate:"required"`
	DepartmentObjectiveID  string  `json:"departmentObjectiveId" validate:"required"`
	DepartmentID           int     `json:"departmentId"`
	JobGradeGroup          string  `json:"jobGradeGroup"`
	SBUName                *string `json:"sbuName"`
	WorkProductName        *string `json:"workProductName"`
	WorkProductDescription *string `json:"workProductDescription"`
	WorkProductDeliverable *string `json:"workProductDeliverable"`
}

// CreateDivisionObjectiveVm is the DTO for creating a division objective.
type CreateDivisionObjectiveVm struct {
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	Kpi                    string  `json:"kpi"`
	DivisionID             int     `json:"divisionId" validate:"required"`
	Target                 string  `json:"target"`
	DepartmentObjectiveID  string  `json:"departmentObjectiveId" validate:"required"`
	DepartmentID           int     `json:"departmentId"`
	JobGradeGroup          string  `json:"jobGradeGroup"`
	SBUName                *string `json:"sbuName"`
	WorkProductName        *string `json:"workProductName"`
	WorkProductDescription *string `json:"workProductDescription"`
	WorkProductDeliverable *string `json:"workProductDeliverable"`
}

// ---------------------------------------------------------------------------
// OfficeObjectiveVmDTO – mirrors .NET OfficeObjectiveVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// OfficeObjectiveVmDTO is the DTO for an office-level objective.
type OfficeObjectiveVmDTO struct {
	BaseWorkFlowVm
	OfficeObjectiveID              string  `json:"officeObjectiveId"`
	Name                           string  `json:"name" validate:"required"`
	Description                    string  `json:"description"`
	Kpi                            string  `json:"kpi"`
	Target                         string  `json:"target"`
	OfficeID                       int     `json:"officeId"`
	DivisionObjectiveID            string  `json:"divisionObjectiveId"`
	JobGradeGroupID                int     `json:"jobGradeGroupId"`
	JobGradeGroupName              string  `json:"jobGradeGroupName"`
	SBUName                        *string `json:"sbuName"`
	WorkProductName                *string `json:"workProductName"`
	WorkProductDescription         *string `json:"workProductDescription"`
	WorkProductDeliverable         *string `json:"workProductDeliverable"`
	ParentDivisionID               int     `json:"parentDivisionId"`
	ParentDivisionObjectiveName    string  `json:"parentDivisionObjectiveName"`
}

// CreateOfficeObjectiveVm is the DTO for creating an office objective.
type CreateOfficeObjectiveVm struct {
	OfficeObjectiveID           string  `json:"officeObjectiveId"`
	Name                        string  `json:"name" validate:"required"`
	Description                 string  `json:"description"`
	Kpi                         string  `json:"kpi"`
	Target                      string  `json:"target"`
	OfficeID                    int     `json:"officeId"`
	DivisionObjectiveID         string  `json:"divisionObjectiveId"`
	JobGradeGroupID             int     `json:"jobGradeGroupId"`
	JobGradeGroupName           string  `json:"jobGradeGroupName"`
	SBUName                     *string `json:"sbuName"`
	WorkProductName             *string `json:"workProductName"`
	WorkProductDescription      *string `json:"workProductDescription"`
	WorkProductDeliverable      *string `json:"workProductDeliverable"`
	ParentDivisionID            int     `json:"parentDivisionId"`
	ParentDivisionObjectiveName string  `json:"parentDivisionObjectiveName"`
}

// ---------------------------------------------------------------------------
// ObjectiveCategoryVmDTO – mirrors .NET ObjectiveCategoryVm : BaseWorkFlow
// ---------------------------------------------------------------------------

// ObjectiveCategoryVmDTO is the DTO for an objective category.
type ObjectiveCategoryVmDTO struct {
	BaseWorkFlowVm
	ObjectiveCategoryID string `json:"objectiveCategoryId"`
	Name                string `json:"name" validate:"required"`
	Description         string `json:"description"`
}

// CreateObjectiveCategoryVm is the DTO for creating an objective category.
type CreateObjectiveCategoryVm struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// ---------------------------------------------------------------------------
// CascadedObjectiveUploadVm – mirrors .NET CascadedObjectiveUploadVm
// ---------------------------------------------------------------------------

// CascadedObjectiveUploadVm carries a parsed row from the cascaded objective
// upload spreadsheet.
type CascadedObjectiveUploadVm struct {
	StrategyID             string  `json:"strategyId"`
	StrategicThemeID       string  `json:"strategicThemeId"`
	EObjName               string  `json:"eObjName"`
	EObjDesc               string  `json:"eObjDesc"`
	EObjKPI                string  `json:"eObjKPI"`
	EObjTarget             string  `json:"eObjTarget"`
	EObjCategory           string  `json:"eObjCategory"`
	Dept                   string  `json:"dept"`
	DeptObjName            string  `json:"deptObjName"`
	DeptObjDesc            string  `json:"deptObjDesc"`
	DeptObjKPI             string  `json:"deptObjKPI"`
	DeptObjTarget          string  `json:"deptObjTarget"`
	Division               string  `json:"division"`
	DivObjName             string  `json:"divObjName"`
	DivObjDesc             string  `json:"divObjDesc"`
	DivObjKPI              string  `json:"divObjKPI"`
	DivObjTarget           string  `json:"divObjTarget"`
	Office                 string  `json:"office"`
	OffObjName             string  `json:"offObjName"`
	OffObjDesc             string  `json:"offObjDesc"`
	OffObjKPI              string  `json:"offObjKPI"`
	OffObjTarget           string  `json:"offObjTarget"`
	JobGradeGroup          string  `json:"jobGradeGroup"`
	SBULevel               *string `json:"sbuLevel"`
	WorkProductName        *string `json:"workProductName"`
	WorkProductDescription *string `json:"workProductDescription"`
	WorkProductDeliverable *string `json:"workProductDeliverable"`
}

// ---------------------------------------------------------------------------
// ObjectiveLevelVm – mirrors .NET ObjectiveLevelVm
// ---------------------------------------------------------------------------

// ObjectiveLevelVm pairs an ObjectiveLevel enum value with its display name.
type ObjectiveLevelVm struct {
	ObjectiveLevel enums.ObjectiveLevel `json:"objectiveLevel"`
	Description    string               `json:"description"`
}

// ---------------------------------------------------------------------------
// ObjectivesUploadRequestModel – carries the parsed rows from a cascaded
// objective upload spreadsheet for bulk processing.
// ---------------------------------------------------------------------------

// ObjectivesUploadRequestModel wraps a list of UploadCascadedObjectiveVm rows
// together with the identity of the user performing the upload.
type ObjectivesUploadRequestModel struct {
	Objectives []UploadCascadedObjectiveVm `json:"objectives" validate:"required,min=1"`
	CreatedBy  string                      `json:"createdBy"  validate:"required"`
}
