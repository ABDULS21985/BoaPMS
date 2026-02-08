package performance

import (
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ---------------------------------------------------------------------------
// Project DTOs  (source: PMSVms/ProjectVm.cs, PMSVms/SetupVm/ProjectVm.cs)
// ---------------------------------------------------------------------------

// ProjectMemberVm is the read/response DTO for a project member assignment.
type ProjectMemberVm struct {
	BaseAuditVm
	ProjectMemberID    string `json:"projectMemberId"`
	StaffID            string `json:"staffId"`
	ProjectID          string `json:"projectId"`
	PlannedObjectiveID string `json:"plannedObjectiveId"`
}

// WorkProductSummaryVm is a lightweight DTO for work products nested inside
// project and committee responses.
type WorkProductSummaryVm struct {
	BaseAuditVm
	WorkProductID   string               `json:"workProductId" validate:"required"`
	Name            string               `json:"name" validate:"required"`
	Description     string               `json:"description"`
	MaxPoint        float64              `json:"maxPoint" validate:"required"`
	WorkProductType enums.WorkProductType `json:"workProductType"`
	IsSelfCreated   bool                 `json:"isSelfCreated"`
	StaffID         string               `json:"staffId" validate:"required"`
	StartDate       time.Time            `json:"startDate" validate:"required"`
	EndDate         time.Time            `json:"endDate" validate:"required"`
	Deliverables    string               `json:"deliverables"`
	FinalScore      float64              `json:"finalScore"`
}

// PerformanceReviewPeriodSummaryVm is a lightweight DTO for the review period
// embedded in project and committee responses.
type PerformanceReviewPeriodSummaryVm struct {
	PeriodID  string    `json:"periodId"`
	Name      string    `json:"name"`
	ShortName string    `json:"shortName"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

// ProjectVm is the read/response DTO for a project record.
type ProjectVm struct {
	BaseAuditVm
	ProjectID      string                            `json:"projectId"`
	ProjectManager string                            `json:"projectManager"`
	ProjectMembers []ProjectMemberVm                 `json:"projectMembers"`
	Name           string                            `json:"name"`
	Description    string                            `json:"description"`
	StartDate      time.Time                         `json:"startDate"`
	EndDate        time.Time                         `json:"endDate"`
	Deliverables   string                            `json:"deliverables"`
	RecordStatus   enums.Status                      `json:"recordStatus"`
	ReviewPeriodID string                            `json:"reviewPeriodId"`
	ReviewPeriod   *PerformanceReviewPeriodSummaryVm `json:"reviewPeriod"`
	WorkProducts   []WorkProductSummaryVm            `json:"workProducts"`
}

// ---------------------------------------------------------------------------
// Committee DTOs  (source: PMSVms/SetupVm/CommitteeVm.cs)
// ---------------------------------------------------------------------------

// CommitteeMemberVm is the read/response DTO for a committee member assignment.
type CommitteeMemberVm struct {
	BaseAuditVm
	CommitteeMemberID string `json:"committeeMemberId" validate:"required"`
	StaffID           string `json:"staffId" validate:"required"`
	CommitteeID       string `json:"committeeId" validate:"required"`
	PlannedObjectiveID string `json:"plannedObjectiveId"`
}

// CommitteeVm is the read/response DTO for a committee record.
type CommitteeVm struct {
	BaseAuditVm
	CommitteeID      string                            `json:"committeeId"`
	Chairperson      string                            `json:"chairperson"`
	CommitteeMembers []CommitteeMemberVm               `json:"committeeMembers"`
	Name             string                            `json:"name"`
	Description      string                            `json:"description"`
	StartDate        time.Time                         `json:"startDate"`
	EndDate          time.Time                         `json:"endDate"`
	Deliverables     string                            `json:"deliverables"`
	RecordStatus     enums.Status                      `json:"recordStatus"`
	ReviewPeriodID   string                            `json:"reviewPeriodId"`
	ReviewPeriod     *PerformanceReviewPeriodSummaryVm `json:"reviewPeriod"`
	WorkProducts     []WorkProductSummaryVm            `json:"workProducts"`
}

// ---------------------------------------------------------------------------
// Project Milestone DTOs  (source: PMSVms/ProjectMilestoneVm.cs)
// ---------------------------------------------------------------------------

// ProjectMilestoneVm is the read/response DTO for a project milestone.
type ProjectMilestoneVm struct {
	BaseAuditVm
	ProjectMilestoneID string       `json:"projectMilestoneId"`
	Name               string       `json:"name"`
	Description        string       `json:"description"`
	ProjectID          string       `json:"projectId"`
	DueDate            time.Time    `json:"dueDate"`
	CompletionDate     *time.Time   `json:"completionDate"`
	RecordStatus       enums.Status `json:"recordStatus"`
}

// IsCompleted is a computed property that returns true when the milestone
// has a recorded completion date.
func (m ProjectMilestoneVm) IsCompleted() bool {
	return m.CompletionDate != nil
}

// IsOverdue is a computed property that returns true when the milestone
// has passed its due date without being completed.
func (m ProjectMilestoneVm) IsOverdue() bool {
	if m.CompletionDate != nil {
		return false
	}
	return time.Now().After(m.DueDate)
}
