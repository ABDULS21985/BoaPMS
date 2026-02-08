package handler

import (
	"net/http"
	"strconv"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/pkg/response"
)

// SelectItem represents a single item in an enum select list.
type SelectItem struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

// newItem is a convenience helper to build a SelectItem from an int value and label.
func newItem(value int, text string) SelectItem {
	return SelectItem{Text: text, Value: strconv.Itoa(value)}
}

// ============================================================
// Enum List Handlers
// ============================================================

// GetObjectiveLevels handles GET /api/v1/enums/objective-levels
// Returns the objective level enum as a select list.
func GetObjectiveLevels(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.ObjectiveLevelDepartment), "Department"),
		newItem(int(enums.ObjectiveLevelDivision), "Division"),
		newItem(int(enums.ObjectiveLevelOffice), "Office"),
		newItem(int(enums.ObjectiveLevelEnterprise), "Enterprise"),
	}
	response.OK(w, items)
}

// GetExtensionTargetTypes handles GET /api/v1/enums/extension-target-types
// Returns the review period extension target type enum as a select list.
func GetExtensionTargetTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.ExtensionTargetBankwide), "Bankwide"),
		newItem(int(enums.ExtensionTargetDepartment), "Department"),
		newItem(int(enums.ExtensionTargetDivision), "Division"),
		newItem(int(enums.ExtensionTargetOffice), "Office"),
		newItem(int(enums.ExtensionTargetStaff), "Staff"),
	}
	response.OK(w, items)
}

// GetEvaluationTypes handles GET /api/v1/enums/evaluation-types
// Returns the evaluation type enum as a select list.
func GetEvaluationTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.EvaluationTypeTimeliness), "Timeliness"),
		newItem(int(enums.EvaluationTypeQuality), "Quality"),
		newItem(int(enums.EvaluationTypeOutput), "Output"),
	}
	response.OK(w, items)
}

// GetWorkProductTypes handles GET /api/v1/enums/work-product-types
// Returns the work product type enum as a select list.
func GetWorkProductTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.WorkProductTypeOperational), "Operational"),
		newItem(int(enums.WorkProductTypeProject), "Project"),
		newItem(int(enums.WorkProductTypeCommittee), "Committee"),
	}
	response.OK(w, items)
}

// GetGrievanceTypes handles GET /api/v1/enums/grievance-types
// Returns the grievance type enum as a select list.
func GetGrievanceTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.GrievanceTypeNone), "None"),
		newItem(int(enums.GrievanceTypeWorkProductEvaluation), "Work Product Evaluation"),
		newItem(int(enums.GrievanceTypeWorkProductAssignment), "Work Product Assignment"),
		newItem(int(enums.GrievanceTypeWorkProductPlanning), "Work Product Planning"),
		newItem(int(enums.GrievanceTypeObjectivePlanning), "Objective Planning"),
	}
	response.OK(w, items)
}

// GetFeedbackRequestTypes handles GET /api/v1/enums/feedback-request-types
// Returns the feedback request type enum as a select list.
func GetFeedbackRequestTypes(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.FeedbackRequestWorkProductEvaluation), "Work Product Evaluation"),
		newItem(int(enums.FeedbackRequestObjectivePlanning), "Objective Planning"),
		newItem(int(enums.FeedbackRequestProjectPlanning), "Project Planning"),
		newItem(int(enums.FeedbackRequestCommitteePlanning), "Committee Planning"),
		newItem(int(enums.FeedbackRequestWorkProductFeedback), "Work Product Feedback"),
		newItem(int(enums.FeedbackRequest360ReviewFeedback), "360 Review Feedback"),
		newItem(int(enums.FeedbackRequestWorkProductPlanning), "Work Product Planning"),
		newItem(int(enums.FeedbackRequestCompetencyReview), "Competency Review"),
		newItem(int(enums.FeedbackRequestReviewPeriod), "Review Period"),
		newItem(int(enums.FeedbackRequestReviewPeriodExtension), "Review Period Extension"),
		newItem(int(enums.FeedbackRequestProjectMemberAssignment), "Project Member Assignment"),
		newItem(int(enums.FeedbackRequestCommitteeMemberAssignment), "Committee Member Assignment"),
		newItem(int(enums.FeedbackRequestPeriodObjectiveOutcome), "Period Objective Outcome"),
		newItem(int(enums.FeedbackRequestDeptObjectiveOutcome), "Department Objective Outcome"),
		newItem(int(enums.FeedbackRequestReviewPeriod360Review), "Review Period 360 Review"),
		newItem(int(enums.FeedbackRequestProjectWorkProductDef), "Project Work Product Definition"),
		newItem(int(enums.FeedbackRequestCommitteeWorkProductDef), "Committee Work Product Definition"),
	}
	response.OK(w, items)
}

// GetPerformanceGrades handles GET /api/v1/enums/performance-grades
// Returns the performance grade enum as a select list.
func GetPerformanceGrades(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.PerformanceGradeProbation), enums.PerformanceGradeProbation.String()),
		newItem(int(enums.PerformanceGradeDeveloping), enums.PerformanceGradeDeveloping.String()),
		newItem(int(enums.PerformanceGradeProgressive), enums.PerformanceGradeProgressive.String()),
		newItem(int(enums.PerformanceGradeCompetent), enums.PerformanceGradeCompetent.String()),
		newItem(int(enums.PerformanceGradeAccomplished), enums.PerformanceGradeAccomplished.String()),
		newItem(int(enums.PerformanceGradeExemplary), enums.PerformanceGradeExemplary.String()),
	}
	response.OK(w, items)
}

// GetReviewPeriodRanges handles GET /api/v1/enums/review-period-ranges
// Returns the review period range enum as a select list.
func GetReviewPeriodRanges(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.ReviewPeriodRangeQuarterly), "Quarterly"),
		newItem(int(enums.ReviewPeriodRangeBiAnnual), "Bi-Annual"),
		newItem(int(enums.ReviewPeriodRangeAnnual), "Annual"),
	}
	response.OK(w, items)
}

// GetStatuses handles GET /api/v1/enums/statuses
// Returns the status enum as a select list.
func GetStatuses(w http.ResponseWriter, r *http.Request) {
	items := []SelectItem{
		newItem(int(enums.StatusDraft), enums.StatusDraft.String()),
		newItem(int(enums.StatusPendingApproval), enums.StatusPendingApproval.String()),
		newItem(int(enums.StatusApprovedAndActive), enums.StatusApprovedAndActive.String()),
		newItem(int(enums.StatusReturned), enums.StatusReturned.String()),
		newItem(int(enums.StatusRejected), enums.StatusRejected.String()),
		newItem(int(enums.StatusAwaitingEvaluation), enums.StatusAwaitingEvaluation.String()),
		newItem(int(enums.StatusCompleted), enums.StatusCompleted.String()),
		newItem(int(enums.StatusPaused), enums.StatusPaused.String()),
		newItem(int(enums.StatusCancelled), enums.StatusCancelled.String()),
		newItem(int(enums.StatusBreached), enums.StatusBreached.String()),
		newItem(int(enums.StatusDeactivated), enums.StatusDeactivated.String()),
		newItem(int(enums.StatusAll), enums.StatusAll.String()),
		newItem(int(enums.StatusClosed), enums.StatusClosed.String()),
		newItem(int(enums.StatusPendingAcceptance), enums.StatusPendingAcceptance.String()),
		newItem(int(enums.StatusActive), enums.StatusActive.String()),
		newItem(int(enums.StatusPendingResolution), enums.StatusPendingResolution.String()),
		newItem(int(enums.StatusResolvedAwaitingFeedback), enums.StatusResolvedAwaitingFeedback.String()),
		newItem(int(enums.StatusEscalated), enums.StatusEscalated.String()),
		newItem(int(enums.StatusAwaitingRespondentComment), enums.StatusAwaitingRespondentComment.String()),
		newItem(int(enums.StatusPendingHODReview), enums.StatusPendingHODReview.String()),
		newItem(int(enums.StatusPendingBUHeadReview), enums.StatusPendingBUHeadReview.String()),
		newItem(int(enums.StatusPendingHRDReview), enums.StatusPendingHRDReview.String()),
		newItem(int(enums.StatusPendingHRDApproval), enums.StatusPendingHRDApproval.String()),
		newItem(int(enums.StatusSuspensionPendingApproval), enums.StatusSuspensionPendingApproval.String()),
		newItem(int(enums.StatusReEvaluate), enums.StatusReEvaluate.String()),
	}
	response.OK(w, items)
}
