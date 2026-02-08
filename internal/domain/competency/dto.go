package competency

import "time"

// ---------------------------------------------------------------------------
// Base DTO types (mirrors shared base types from domain package; defined
// locally to avoid circular-import issues when other packages reference
// these DTOs).
// ---------------------------------------------------------------------------

// BaseAuditVm carries common audit fields on every view-model.
type BaseAuditVm struct {
	CreatedBy   string     `json:"createdBy"`
	DateCreated time.Time  `json:"dateCreated"`
	IsActive    bool       `json:"isActive"`
	Status      string     `json:"status"`
	DateUpdated *time.Time `json:"dateUpdated"`
	UpdatedBy   string     `json:"updatedBy"`
}

// BaseAPIResponse is the standard envelope for list-style API responses.
type BaseAPIResponse struct {
	HasError bool   `json:"hasError"`
	Message  string `json:"message"`
}

// BasePagedData carries common paging parameters.
type BasePagedData struct {
	Skip     int `json:"skip"`
	PageSize int `json:"pageSize"`
}

// ---------------------------------------------------------------------------
// 1. CompetencyVm
// ---------------------------------------------------------------------------

// CompetencyVm is the read/display DTO for a single competency.
type CompetencyVm struct {
	BaseAuditVm
	CompetencyID         int        `json:"competencyId"`
	CompetencyCategoryID int        `json:"competencyCategoryId"`
	CompetencyName       string     `json:"competencyName"`
	Description          string     `json:"description"`
	CompetencyCategoryName string   `json:"competencyCategoryName"`
	ApprovedBy           string     `json:"approvedBy"`
	DateApproved         *time.Time `json:"dateApproved"`
	IsApproved           bool       `json:"isApproved"`
	IsRejected           bool       `json:"isRejected"`
	RejectedBy           string     `json:"rejectedBy"`
	RejectionReason      string     `json:"rejectionReason"`
	IsSelected           bool       `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// 2. ApproveCompetencyVm
// ---------------------------------------------------------------------------

// ApproveCompetencyVm is the payload for approving a competency.
type ApproveCompetencyVm struct {
	BaseAuditVm
	CompetencyID int        `json:"competencyId"`
	ApprovedBy   string     `json:"approvedBy"`
	DateApproved *time.Time `json:"dateApproved"`
	IsApproved   bool       `json:"isApproved"`
}

// ---------------------------------------------------------------------------
// 3. SaveCompetencyVm
// ---------------------------------------------------------------------------

// SaveCompetencyVm is the payload for creating or updating a competency.
type SaveCompetencyVm struct {
	BaseAuditVm
	CompetencyID                int                          `json:"competencyId"`
	CompetencyCategoryID        int                          `json:"competencyCategoryId"`
	CompetencyName              string                       `json:"competencyName"`
	Description                 string                       `json:"description"`
	CompetencyCategoryName      string                       `json:"competencyCategoryName"`
	CompetencyRatingDefinitions []CompetencyRatingDefinitionVm `json:"competencyRatingDefinitions"`
}

// ---------------------------------------------------------------------------
// 4. SaveUploadedCompetencyVm
// ---------------------------------------------------------------------------

// SaveUploadedCompetencyVm carries a single row from a bulk-upload spreadsheet.
type SaveUploadedCompetencyVm struct {
	BaseAuditVm
	ID                     int    `json:"id"`
	CompetencyID           int    `json:"competencyId"`
	CompetencyCategoryID   int    `json:"competencyCategoryId"`
	CompetencyName         string `json:"competencyName"`
	Description            string `json:"description"`
	CompetencyCategoryName string `json:"competencyCategoryName"`
	IsValidRecord          bool   `json:"isValidRecord"`
	IsSuccess              *bool  `json:"isSuccess"`
	Message                string `json:"message"`
	IsProcessed            *bool  `json:"isProcessed"`
	IsSelected             bool   `json:"isSelected"`
	EntryDefinition        string `json:"entryDefinition"`
	BasicDefinition        string `json:"basicDefinition"`
	IntermediateDefinition string `json:"intermediateDefinition"`
	ExpertDefinition       string `json:"expertDefinition"`
}

// ---------------------------------------------------------------------------
// 5. CompetencyListVm
// ---------------------------------------------------------------------------

// CompetencyListVm is the paged response for competency queries.
type CompetencyListVm struct {
	BaseAPIResponse
	Competencies []CompetencyVm `json:"competencies"`
	TotalRecord  int            `json:"totalRecord"`
}

// ---------------------------------------------------------------------------
// 6. SearchCompetencyVm
// ---------------------------------------------------------------------------

// SearchCompetencyVm carries filters for searching competencies.
type SearchCompetencyVm struct {
	BasePagedData
	CategoryID   *int   `json:"categoryId"`
	SearchString string `json:"searchString"`
	IsApproved   *bool  `json:"isApproved"`
	IsRejected   *bool  `json:"isRejected"`
	IsTechnical  *bool  `json:"isTechnical"`
}

// ---------------------------------------------------------------------------
// 7. SearchJobRoleCompetencyVm
// ---------------------------------------------------------------------------

// SearchJobRoleCompetencyVm carries filters for job-role competency queries.
type SearchJobRoleCompetencyVm struct {
	BasePagedData
	DepartmentID *int   `json:"departmentId"`
	DivisionID   *int   `json:"divisionId"`
	OfficeID     *int   `json:"officeId"`
	JobRoleID    *int   `json:"jobRoleId"`
	SearchString string `json:"searchString"`
}

// ---------------------------------------------------------------------------
// 8. RejectCompetencyVm
// ---------------------------------------------------------------------------

// RejectCompetencyVm is the payload for rejecting a competency.
type RejectCompetencyVm struct {
	BaseAuditVm
	CompetencyID    int        `json:"competencyId"`
	RejectedBy      string     `json:"rejectedBy"`
	RejectionReason string     `json:"rejectionReason" validate:"required"`
	DateRejected    *time.Time `json:"dateRejected"`
	IsRejected      bool       `json:"isRejected"`
}

// ---------------------------------------------------------------------------
// 9. CompetencyCategoryVm
// ---------------------------------------------------------------------------

// CompetencyCategoryVm is the read/display DTO for a competency category.
type CompetencyCategoryVm struct {
	BaseAuditVm
	CompetencyCategoryID int    `json:"competencyCategoryId"`
	CategoryName         string `json:"categoryName"`
	IsTechnical          bool   `json:"isTechnical"`
}

// ---------------------------------------------------------------------------
// 10. CompetencyCategoryGradingVm
// ---------------------------------------------------------------------------

// CompetencyCategoryGradingVm is the DTO for category-level grading weights.
type CompetencyCategoryGradingVm struct {
	BaseAuditVm
	CompetencyCategoryGradingID int     `json:"competencyCategoryGradingId"`
	CompetencyCategoryID        int     `json:"competencyCategoryId"        validate:"required"`
	ReviewTypeID                int     `json:"reviewTypeId"                validate:"required"`
	WeightPercentage            float64 `json:"weightPercentage"            validate:"required"`
	CompetencyCategoryName      string  `json:"competencyCategoryName"`
	ReviewTypeName              string  `json:"reviewTypeName"`
}

// ---------------------------------------------------------------------------
// 11. CompetencyRatingDefinitionVm
// ---------------------------------------------------------------------------

// CompetencyRatingDefinitionVm is the DTO for a competency's rating definition.
type CompetencyRatingDefinitionVm struct {
	BaseAuditVm
	CompetencyRatingDefinitionID int    `json:"competencyRatingDefinitionId"`
	CompetencyID                 int    `json:"competencyId"                 validate:"required"`
	RatingID                     int    `json:"ratingId"                     validate:"required"`
	Definition                   string `json:"definition"                   validate:"required"`
	RatingName                   string `json:"ratingName"`
	RatingValue                  int    `json:"ratingValue"`
	CompetencyName               string `json:"competencyName"`
}

// ---------------------------------------------------------------------------
// 12. CompetencyReviewVm
// ---------------------------------------------------------------------------

// CompetencyReviewVm is the full DTO for a competency review entry.
type CompetencyReviewVm struct {
	BaseAuditVm
	CompetencyReviewID          int                          `json:"competencyReviewId"`
	ReviewTypeID                int                          `json:"reviewTypeId"`
	ReviewPeriodID              int                          `json:"reviewPeriodId"`
	RatingID                    *int                         `json:"ratingId"`
	CompetencyID                int                          `json:"competencyId"`
	ExpectedRatingID            int                          `json:"expectedRatingId"`
	ReviewDate                  *time.Time                   `json:"reviewDate"`
	ReviewerID                  string                       `json:"reviewerId"`
	ReviewerName                string                       `json:"reviewerName"`
	EmployeeNumber              string                       `json:"employeeNumber"`
	ReviewTypeName              string                       `json:"reviewTypeName"`
	ReviewPeriodName            string                       `json:"reviewPeriodName"`
	CompetencyName              string                       `json:"competencyName"`
	CompetencyCategoryName      string                       `json:"competencyCategoryName"`
	CompetencyDefinition        string                       `json:"competencyDefinition"`
	IsTechnical                 bool                         `json:"isTechnical"`
	EmployeeName                string                       `json:"employeeName"`
	EmployeeInitial             string                       `json:"employeeInitial"`
	EmployeeGrade               string                       `json:"employeeGrade"`
	EmployeeDepartment          string                       `json:"employeeDepartment"`
	ActualRatingID              int                          `json:"actualRatingId"`
	ActualRatingName            string                       `json:"actualRatingName"`
	ActualRatingValue           int                          `json:"actualRatingValue"`
	ExpectedRatingName          string                       `json:"expectedRatingName"`
	ExpectedRatingValue         int                          `json:"expectedRatingValue"`
	CompetencyRatingDefinitions []CompetencyRatingDefinitionVm `json:"competencyRatingDefinitions"`
}

// ---------------------------------------------------------------------------
// 13. SearchForReviewDetailVm
// ---------------------------------------------------------------------------

// SearchForReviewDetailVm carries filters for finding review details.
type SearchForReviewDetailVm struct {
	ReviewTypeID   int    `json:"reviewTypeId"`
	ReviewPeriodID int    `json:"reviewPeriodId"`
	ReviewerID     string `json:"reviewerId"`
	EmployeeID     string `json:"employeeId"`
	IsTechnical    bool   `json:"isTechnical"`
}

// ---------------------------------------------------------------------------
// 14. CompetencyReviewDetailVm
// ---------------------------------------------------------------------------

// CompetencyReviewDetailVm wraps a list of competency reviews in an API response.
type CompetencyReviewDetailVm struct {
	BaseAPIResponse
	CompetencyReviews []CompetencyReviewVm `json:"competencyReviews"`
}

// ---------------------------------------------------------------------------
// 15. CalculateReviewProfileVm
// ---------------------------------------------------------------------------

// CalculateReviewProfileVm is the request payload for computing a review profile.
type CalculateReviewProfileVm struct {
	EmployeeNumber string `json:"employeeNumber"`
	ReviewPeriodID int    `json:"reviewPeriodId"`
	IsTechnical    bool   `json:"isTechnical"`
}

// ---------------------------------------------------------------------------
// 16. CompetencyReviewProfileVm
// ---------------------------------------------------------------------------

// CompetencyReviewProfileVm is the DTO for a competency review profile summary.
type CompetencyReviewProfileVm struct {
	BaseAuditVm
	CompetencyReviewProfileID int     `json:"competencyReviewProfileId"`
	ReviewPeriodID            int     `json:"reviewPeriodId"`
	ReviewPeriodName          string  `json:"reviewPeriodName"`
	AverageRatingID           int     `json:"averageRatingId"`
	AverageRatingName         string  `json:"averageRatingName"`
	AverageRatingValue        int     `json:"averageRatingValue"`
	ExpectedRatingID          int     `json:"expectedRatingId"`
	ExpectedRatingName        string  `json:"expectedRatingName"`
	ExpectedRatingValue       int     `json:"expectedRatingValue"`
	AverageScore              float64 `json:"averageScore"`
	EmployeeNumber            string  `json:"employeeNumber"`
	EmployeeFullName          string  `json:"employeeFullName"`
	CompetencyID              int     `json:"competencyId"`
	CompetencyName            string  `json:"competencyName"`
	CompetencyCategory        int     `json:"competencyCategory"`
	CompetencyCategoryName    string  `json:"competencyCategoryName"`
	NumberOfDevelopmentPlans   int    `json:"numberOfDevelopmentPlans"`
	ProgressCount             int     `json:"progressCount"`
	CompletedCount            int     `json:"completedCount"`
	OfficeID                  string  `json:"officeId"`
	OfficeName                string  `json:"officeName"`
	DivisionID                string  `json:"divisionId"`
	DivisionName              string  `json:"divisionName"`
	DepartmentID              string  `json:"departmentId"`
	DepartmentName            string  `json:"departmentName"`
	JobRoleID                 string  `json:"jobRoleId"`
	JobRoleName               string  `json:"jobRoleName"`
	GradeName                 string  `json:"gradeName"`
}

// CompetencyGap returns the gap between expected and average rating values.
func (p *CompetencyReviewProfileVm) CompetencyGap() int {
	return p.ExpectedRatingValue - p.AverageRatingValue
}

// HaveGap returns true when a competency gap exists (expected > average).
func (p *CompetencyReviewProfileVm) HaveGap() bool {
	return p.CompetencyGap() > 0
}

// ---------------------------------------------------------------------------
// 17. CompetencyProfileSummaryVm
// ---------------------------------------------------------------------------

// CompetencyProfileSummaryVm summarises competency profile totals per category.
type CompetencyProfileSummaryVm struct {
	CompetencyCategory string `json:"competencyCategory"`
	TotalExpected      int    `json:"totalExpected"`
	TotalActual        int    `json:"totalActual"`
	TotalCompetencies  int    `json:"totalCompetencies"`
	MatchCompetencies  int    `json:"matchCompetencies"`
	GapCompetencies    int    `json:"gapCompetencies"`
}

// Gap returns the difference between expected and actual totals.
func (s *CompetencyProfileSummaryVm) Gap() int {
	return s.TotalExpected - s.TotalActual
}

// ---------------------------------------------------------------------------
// 18. CompetencyGapClosureVm
// ---------------------------------------------------------------------------

// CompetencyGapClosureVm identifies a profile entry for gap closure processing.
type CompetencyGapClosureVm struct {
	EmployeeNumber            string `json:"employeeNumber"            validate:"required"`
	CompetencyReviewProfileID int    `json:"competencyReviewProfileId" validate:"required"`
}

// ---------------------------------------------------------------------------
// 19. DevelopmentPlanVm
// ---------------------------------------------------------------------------

// DevelopmentPlanVm is the DTO for a development plan entry.
type DevelopmentPlanVm struct {
	BaseAuditVm
	DevelopmentPlanID         int        `json:"developmentPlanId"`
	Activity                  string     `json:"activity"                  validate:"required,min=5,max=500"`
	EmployeeNumber            string     `json:"employeeNumber"`
	CompetencyReviewProfileID int        `json:"competencyReviewProfileId"`
	TargetDate                time.Time  `json:"targetDate"                validate:"required"`
	CompletionDate            *time.Time `json:"completionDate"`
	LearningResource          string     `json:"learningResource"          validate:"required,min=5,max=500"`
	CreatedBy                 string     `json:"createdBy"`
	CompetencyName            string     `json:"competencyName"`
	CompetencyCategoryName    string     `json:"competencyCategoryName"`
	TaskStatus                string     `json:"taskStatus"`
	ReviewPeriod              string     `json:"reviewPeriod"`
	CurrentGap                int        `json:"currentGap"`
	EmployeeName              string     `json:"employeeName"`
	TrainingTypeName          string     `json:"trainingTypeName"          validate:"required"`
}

// ---------------------------------------------------------------------------
// 20. JobRoleVm
// ---------------------------------------------------------------------------

// JobRoleVm is the read/display DTO for a job role.
type JobRoleVm struct {
	BaseAuditVm
	JobRoleID   int    `json:"jobRoleId"`
	JobRoleName string `json:"jobRoleName"`
	Description string `json:"description"`
}

// ---------------------------------------------------------------------------
// 21. JobGradeVm
// ---------------------------------------------------------------------------

// JobGradeVm is the read/display DTO for a job grade.
type JobGradeVm struct {
	BaseAuditVm
	JobGradeID int    `json:"jobGradeId"`
	GradeCode  string `json:"gradeCode"`
	GradeName  string `json:"gradeName"`
}

// ---------------------------------------------------------------------------
// 22. JobGradeGroupVm
// ---------------------------------------------------------------------------

// JobGradeGroupVm is the read/display DTO for a grade group.
type JobGradeGroupVm struct {
	BaseAuditVm
	JobGradeGroupID int    `json:"jobGradeGroupId"`
	GroupName       string `json:"groupName" validate:"required"`
	Order           int    `json:"order"     validate:"required"`
}

// ---------------------------------------------------------------------------
// 23. AssignJobGradeGroupVm
// ---------------------------------------------------------------------------

// AssignJobGradeGroupVm is the DTO for a grade-to-group assignment.
type AssignJobGradeGroupVm struct {
	BaseAuditVm
	AssignJobGradeGroupID int    `json:"assignJobGradeGroupId"`
	JobGradeGroupID       int    `json:"jobGradeGroupId"`
	JobGradeID            int    `json:"jobGradeId"`
	JobGradeName          string `json:"jobGradeName"`
	JobGradeGroupName     string `json:"jobGradeGroupName"`
}

// ---------------------------------------------------------------------------
// 24. JobRoleCompetencyVm
// ---------------------------------------------------------------------------

// JobRoleCompetencyVm is the read/display DTO for a job-role competency mapping.
type JobRoleCompetencyVm struct {
	BaseAuditVm
	JobRoleCompetencyID int    `json:"jobRoleCompetencyId"`
	DepartmentID        int    `json:"departmentId"`
	DivisionID          int    `json:"divisionId"`
	JobRoleID           int    `json:"jobRoleId"    validate:"required"`
	RatingID            int    `json:"ratingId"     validate:"required"`
	OfficeID            int    `json:"officeId"     validate:"required"`
	CompetencyID        int    `json:"competencyId" validate:"required"`
	CompetencyName      string `json:"competencyName"`
	JobRoleName         string `json:"jobRoleName"`
	GradeGroupName      string `json:"gradeGroupName"`
	RatingName          string `json:"ratingName"`
	OfficeName          string `json:"officeName"`
	DepartmentName      string `json:"departmentName"`
	DivisionName        string `json:"divisionName"`
}

// ---------------------------------------------------------------------------
// 25. PagedJobRoleCompetencyVm
// ---------------------------------------------------------------------------

// PagedJobRoleCompetencyVm is the paged response for job-role competency queries.
type PagedJobRoleCompetencyVm struct {
	BaseAPIResponse
	JobRoleCompetencies []JobRoleCompetencyVm `json:"jobRoleCompetencies"`
	TotalRecords        int                   `json:"totalRecords"`
}

// ---------------------------------------------------------------------------
// 26. SaveJobRoleCompetencyVm
// ---------------------------------------------------------------------------

// SaveJobRoleCompetencyVm is the payload for creating/updating job-role competencies.
type SaveJobRoleCompetencyVm struct {
	BaseAuditVm
	JobRoleCompetencyID      int                      `json:"jobRoleCompetencyId"`
	DepartmentID             *int                     `json:"departmentId"`
	DivisionID               *int                     `json:"divisionId"`
	JobRoleID                int                      `json:"jobRoleId"  validate:"required"`
	OfficeID                 *int                     `json:"officeId"   validate:"required"`
	CompetencyGroup          string                   `json:"competencyGroup"`
	JobRoleName              string                   `json:"jobRoleName"`
	OfficeName               string                   `json:"officeName"`
	JobRoleCompetencyRatings []JobRoleCompetencyRating `json:"jobRoleCompetencyRatings"`
}

// ---------------------------------------------------------------------------
// 27. JobRoleCompetencyRating
// ---------------------------------------------------------------------------

// JobRoleCompetencyRating represents a single competency+rating pair within
// a job-role competency save operation.
type JobRoleCompetencyRating struct {
	ID             int    `json:"id"`
	CompetencyID   int    `json:"competencyId"   validate:"required"`
	RatingID       int    `json:"ratingId"       validate:"required"`
	CompetencyName string `json:"competencyName"`
	RatingName     string `json:"ratingName"`
}

// ---------------------------------------------------------------------------
// 28. UploadJobRoleCompetencyVm
// ---------------------------------------------------------------------------

// UploadJobRoleCompetencyVm carries a single row from a bulk-upload for
// job-role competencies.
type UploadJobRoleCompetencyVm struct {
	BaseAuditVm
	ID                  int    `json:"id"`
	JobRoleCompetencyID int    `json:"jobRoleCompetencyId"`
	OfficeID            int    `json:"officeId"`
	JobRoleID           int    `json:"jobRoleId"`
	CompetencyID        int    `json:"competencyId"`
	RatingID            int    `json:"ratingId"`
	OfficeName          string `json:"officeName"`
	JobRoleName         string `json:"jobRoleName"`
	CompetencyName      string `json:"competencyName"`
	RatingName          string `json:"ratingName"`
	IsValidRecord       bool   `json:"isValidRecord"`
	IsSuccess           *bool  `json:"isSuccess"`
	Message             string `json:"message"`
	IsProcessed         *bool  `json:"isProcessed"`
	IsSelected          bool   `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// 29. JobRoleGradeVm
// ---------------------------------------------------------------------------

// JobRoleGradeVm is the DTO for a grade-to-job-role mapping.
type JobRoleGradeVm struct {
	BaseAuditVm
	JobRoleGradeID int    `json:"jobRoleGradeId"`
	JobRoleID      int    `json:"jobRoleId"`
	GradeID        string `json:"gradeId"`
	GradeName      string `json:"gradeName"`
	JobRoleName    string `json:"jobRoleName"`
}

// ---------------------------------------------------------------------------
// 30. OfficeJobRoleVm
// ---------------------------------------------------------------------------

// OfficeJobRoleVm is the DTO for an office-to-job-role mapping.
type OfficeJobRoleVm struct {
	BaseAuditVm
	OfficeJobRoleID int    `json:"officeJobRoleId"`
	OfficeID        int    `json:"officeId"`
	JobRoleID       int    `json:"jobRoleId"`
	OfficeName      string `json:"officeName"`
	JobRoleName     string `json:"jobRoleName"`
}

// ---------------------------------------------------------------------------
// 31. OfficeJobRoleListVm
// ---------------------------------------------------------------------------

// OfficeJobRoleListVm is the paged response for office-job-role queries.
type OfficeJobRoleListVm struct {
	BaseAPIResponse
	OfficeJobRoles []OfficeJobRoleVm `json:"officeJobRoles"`
	TotalRecord    int               `json:"totalRecord"`
}

// ---------------------------------------------------------------------------
// 32. SearchOfficeJobRoleVm
// ---------------------------------------------------------------------------

// SearchOfficeJobRoleVm carries filters for office-job-role searches.
type SearchOfficeJobRoleVm struct {
	Skip         int    `json:"skip"`
	PageSize     int    `json:"pageSize"`
	SearchString string `json:"searchString"`
	OfficeID     *int   `json:"officeId"`
}

// ---------------------------------------------------------------------------
// 33. BehavioralCompetencyVm
// ---------------------------------------------------------------------------

// BehavioralCompetencyVm is the read/display DTO for a behavioral competency.
type BehavioralCompetencyVm struct {
	BaseAuditVm
	BehavioralCompetencyID int    `json:"behavioralCompetencyId"`
	CompetencyID           int    `json:"competencyId"     validate:"required"`
	JobGradeGroupID        int    `json:"jobGradeGroupId"  validate:"required"`
	RatingID               int    `json:"ratingId"`
	CompetencyName         string `json:"competencyName"`
	RatingName             string `json:"ratingName"`
	JobGradeGroupName      string `json:"jobGradeGroupName"`
}

// ---------------------------------------------------------------------------
// 34. SaveBehavioralCompetencyVm
// ---------------------------------------------------------------------------

// SaveBehavioralCompetencyVm is the payload for creating/updating behavioral
// competency mappings.
type SaveBehavioralCompetencyVm struct {
	BaseAuditVm
	BehavioralCompetencyID   int                      `json:"behavioralCompetencyId"`
	CompetencyID             int                      `json:"competencyId"    validate:"required"`
	JobGradeGroupID          int                      `json:"jobGradeGroupId" validate:"required"`
	JobRoleCompetencyRatings []JobRoleCompetencyRating `json:"jobRoleCompetencyRatings"`
}

// ---------------------------------------------------------------------------
// 35. UploadBehavioralCompetencyVm
// ---------------------------------------------------------------------------

// UploadBehavioralCompetencyVm carries a single row from a bulk-upload for
// behavioral competencies.
type UploadBehavioralCompetencyVm struct {
	BaseAuditVm
	ID                     int    `json:"id"`
	BehavioralCompetencyID int    `json:"behavioralCompetencyId"`
	CompetencyID           int    `json:"competencyId"`
	JobGradeGroupID        int    `json:"jobGradeGroupId"`
	RatingID               int    `json:"ratingId"`
	CompetencyName         string `json:"competencyName"`
	RatingName             string `json:"ratingName"`
	JobGradeGroupName      string `json:"jobGradeGroupName"`
	IsValidRecord          bool   `json:"isValidRecord"`
	IsSuccess              *bool  `json:"isSuccess"`
	Message                string `json:"message"`
	IsProcessed            *bool  `json:"isProcessed"`
	IsSelected             bool   `json:"isSelected"`
}

// ---------------------------------------------------------------------------
// 36. RatingVm
// ---------------------------------------------------------------------------

// RatingVm is the read/display DTO for a rating level.
type RatingVm struct {
	BaseAuditVm
	RatingID int    `json:"ratingId"`
	Name     string `json:"name"  validate:"required"`
	Value    int    `json:"value" validate:"required"`
}

// ---------------------------------------------------------------------------
// 37. ReviewPeriodVm
// ---------------------------------------------------------------------------

// ReviewPeriodVm is the DTO for a competency review period.
type ReviewPeriodVm struct {
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

// ---------------------------------------------------------------------------
// 38. ReviewTypeVm
// ---------------------------------------------------------------------------

// ReviewTypeVm is the read/display DTO for a review type.
type ReviewTypeVm struct {
	BaseAuditVm
	ReviewTypeID   int    `json:"reviewTypeId"`
	ReviewTypeName string `json:"reviewTypeName"`
}

// ---------------------------------------------------------------------------
// 39. TrainingTypeVm
// ---------------------------------------------------------------------------

// TrainingTypeVm is the read/display DTO for a training type.
// Note: embeds BaseAuditVm (mirrors .NET BaseAudit embedding).
type TrainingTypeVm struct {
	BaseAuditVm
	TrainingTypeID   int    `json:"trainingTypeId"`
	TrainingTypeName string `json:"trainingTypeName"`
}

// ---------------------------------------------------------------------------
// 40. BankYearVm
// ---------------------------------------------------------------------------

// BankYearVm is the read/display DTO for a fiscal year.
type BankYearVm struct {
	BaseAuditVm
	BankYearID int    `json:"bankYearId"`
	YearName   string `json:"yearName"`
}
