package dashboard

// BasicStatVm holds high-level counts for the dashboard overview.
type BasicStatVm struct {
	StaffCount   int `json:"staffCount"`
	ProjectCount int `json:"projectCount"`
	SkillCount   int `json:"skillCount"`
}

// ChartDataVm represents a single data point for actual-vs-expected charts.
type ChartDataVm struct {
	Label    string  `json:"label"`
	Actual   float64 `json:"actual"`
	Expected float64 `json:"expected"`
}

// PMSChartDataVm represents a single data point for single-value PMS charts.
type PMSChartDataVm struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// ManagerOverviewReportVm summarises review completion for a manager's team.
type ManagerOverviewReportVm struct {
	NoOfCompletedReviews  int                 `json:"noOfCompletedReviews"`
	NoStartedReviews      int                 `json:"noStartedReviews"`
	NoNotStartedReviews   int                 `json:"noNotStartedReviews"`
	BasicEmployeeDatas    []BasicEmployeeData `json:"basicEmployeeDatas"`
}

// BasicEmployeeData holds basic employee info with review completion status.
type BasicEmployeeData struct {
	FullName                  string `json:"fullName"`
	EmployeeNumber            string `json:"employeeNumber"`
	Grade                     string `json:"grade"`
	Position                  string `json:"position"`
	Department                string `json:"department"`
	Office                    string `json:"office"`
	NoOfCompletedReviews      int    `json:"noOfCompletedReviews"`
	NoOfNotCompletedReviews   int    `json:"noOfNotCompletedReviews"`
}

// IsCompleted returns true when all reviews have been completed.
func (b BasicEmployeeData) IsCompleted() bool {
	return b.NoOfNotCompletedReviews == 0
}

// ReviewStatusVm represents the status of reviews for dashboard visualisation.
type ReviewStatusVm struct {
	StatusName string `json:"statusName"`
	Fill       string `json:"fill"`
	Users      int    `json:"users"`
	Text       string `json:"text"`
}

// GroupedCompetencyReviewProfileVm groups competency statistics by category.
type GroupedCompetencyReviewProfileVm struct {
	CategoryCompetencyDetailStats []CategoryCompetencyDetailStat `json:"categoryCompetencyDetailStats"`
	CategoryCompetencyStats       []CategoryCompetencyStat       `json:"categoryCompetencyStats"`
}

// CategoryCompetencyStat holds expected-vs-actual data for a competency category.
type CategoryCompetencyStat struct {
	CategoryName string  `json:"categoryName"`
	Expected     float64 `json:"expected"`
	Actual       float64 `json:"actual"`
}

// CategoryCompetencyDetailStat holds detailed rating statistics for a competency category.
type CategoryCompetencyDetailStat struct {
	CategoryName           string               `json:"categoryName"`
	AverageRating          float64              `json:"averageRating"`
	HighestRating          float64              `json:"highestRating"`
	LowestRating           float64              `json:"lowestRating"`
	MostCommonRating       float64              `json:"mostCommonRating"`
	CompetencyRatingStat   []CompetencyRatingStat `json:"competencyRatingStat"`
	GroupCompetencyRatings []ChartDataVm        `json:"groupCompetencyRatings"`
}

// CompetencyRatingStat describes the distribution of a specific rating level.
type CompetencyRatingStat struct {
	RatingOrder     int     `json:"ratingOrder"`
	RatingValue     float64 `json:"ratingValue"`
	RatingName      string  `json:"ratingName"`
	NumberOfStaff   int     `json:"numberOfStaff"`
	StaffPercentage float64 `json:"staffPercentage"`
}

// SearchGroupedCompetencyReviewProfileVm carries filter criteria for
// grouped competency review profile queries.
type SearchGroupedCompetencyReviewProfileVm struct {
	ReviewPeriodID int  `json:"reviewPeriodId"`
	OfficeID       *int `json:"officeId"`
	DivisionID     *int `json:"divisionId"`
	DepartmentID   *int `json:"departmentId"`
	JobRoleID      *int `json:"jobRoleId"`
}

// CompetencyMatrixReviewOverviewVm is the top-level overview of a competency matrix report.
type CompetencyMatrixReviewOverviewVm struct {
	CompetencyMatrixReviewProfiles []CompetencyMatrixReviewProfileVm `json:"competencyMatrixReviewProfiles"`
	CompetencyNames                []string                          `json:"competencyNames"`
}

// CompetencyMatrixReviewProfileVm holds a single employee's competency matrix data.
type CompetencyMatrixReviewProfileVm struct {
	EmployeeID             string                    `json:"employeeId"`
	EmployeeName           string                    `json:"employeeName"`
	OfficeName             string                    `json:"officeName"`
	DivisionName           string                    `json:"divisionName"`
	DepartmentName         string                    `json:"departmentName"`
	Position               string                    `json:"position"`
	Grade                  string                    `json:"grade"`
	GapCount               int                       `json:"gapCount"`
	NoOfCompetent          int                       `json:"noOfCompetent"`
	NoOfCompetencies       int                       `json:"noOfCompetencies"`
	OverallAverage         float64                   `json:"overallAverage"`
	CompetencyName         string                    `json:"competencyName"`
	CompetencyMatrixDetails []CompetencyMatrixDetailVm `json:"competencyMatrixDetails"`
}

// CompetencyMatrixDetailVm holds the detail for a single competency within the matrix.
type CompetencyMatrixDetailVm struct {
	CompetencyName      string  `json:"competencyName"`
	AverageScore        float64 `json:"averageScore"`
	ExpectedRatingValue float64 `json:"expectedRatingValue"`
}
