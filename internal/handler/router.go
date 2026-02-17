package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/auth"
	"github.com/enterprise-pms/pms-api/internal/middleware"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// jwtProtect wraps a handler with JWT authentication middleware.
func jwtProtect(mw *middleware.Stack, h http.HandlerFunc) http.Handler {
	return mw.JWTAuth(http.HandlerFunc(h))
}

// jwtRoleProtect wraps a handler with JWT authentication + role-based middleware.
func jwtRoleProtect(mw *middleware.Stack, h http.HandlerFunc, roles ...string) http.Handler {
	return mw.JWTAuth(mw.RequireRole(roles...)(http.HandlerFunc(h)))
}

// NewRouter sets up all HTTP routes and middleware chains.
func NewRouter(svc *service.Container, mw *middleware.Stack, cfg *config.Config, log zerolog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Health check — no auth required
	mux.HandleFunc("GET /health", healthCheck)
	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"status": "ready"})
	})

	// ----------------------------------------------------------------
	// Auth routes — public (no JWT required)
	// ----------------------------------------------------------------
	authHandler := NewAuthHandler(svc.Auth, log)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.RefreshToken)

	// Auth routes — JWT required
	mux.Handle("GET /api/v1/auth/validate", jwtProtect(mw, authHandler.ValidateToken))

	// ----------------------------------------------------------------
	// Performance Management routes — JWT required
	// ----------------------------------------------------------------
	perfHandler := NewPerformanceMgtHandler(svc, log)

	// -- Strategies --
	mux.Handle("GET /api/v1/performance/strategies", jwtProtect(mw, perfHandler.GetBankStrategies))
	mux.Handle("POST /api/v1/performance/strategies", jwtProtect(mw, perfHandler.CreateStrategy))
	mux.Handle("PUT /api/v1/performance/strategies", jwtProtect(mw, perfHandler.UpdateStrategy))

	// -- Strategic Themes --
	mux.Handle("GET /api/v1/performance/strategic-themes", jwtProtect(mw, perfHandler.GetBankStrategicThemes))
	mux.Handle("POST /api/v1/performance/strategic-themes", jwtProtect(mw, perfHandler.CreateStrategicTheme))
	mux.Handle("PUT /api/v1/performance/strategic-themes", jwtProtect(mw, perfHandler.UpdateStrategicTheme))

	// -- Enterprise Objectives --
	mux.Handle("GET /api/v1/performance/objectives/enterprise", jwtProtect(mw, perfHandler.GetEnterpriseObjectives))
	mux.Handle("POST /api/v1/performance/objectives/enterprise", jwtProtect(mw, perfHandler.CreateEnterpriseObjective))
	mux.Handle("PUT /api/v1/performance/objectives/enterprise", jwtProtect(mw, perfHandler.UpdateEnterpriseObjective))

	// -- Department Objectives --
	mux.Handle("GET /api/v1/performance/objectives/department", jwtProtect(mw, perfHandler.GetDepartmentObjectives))
	mux.Handle("POST /api/v1/performance/objectives/department", jwtProtect(mw, perfHandler.CreateDepartmentObjective))
	mux.Handle("PUT /api/v1/performance/objectives/department", jwtProtect(mw, perfHandler.UpdateDepartmentObjective))

	// -- Division Objectives --
	mux.Handle("GET /api/v1/performance/objectives/division", jwtProtect(mw, perfHandler.GetDivisionObjectives))
	mux.Handle("POST /api/v1/performance/objectives/division", jwtProtect(mw, perfHandler.CreateDivisionObjective))
	mux.Handle("PUT /api/v1/performance/objectives/division", jwtProtect(mw, perfHandler.UpdateDivisionObjective))

	// -- Office Objectives --
	mux.Handle("GET /api/v1/performance/objectives/office", jwtProtect(mw, perfHandler.GetOfficeObjectives))
	mux.Handle("POST /api/v1/performance/objectives/office", jwtProtect(mw, perfHandler.CreateOfficeObjective))
	mux.Handle("PUT /api/v1/performance/objectives/office", jwtProtect(mw, perfHandler.UpdateOfficeObjective))

	// -- Consolidated Objectives --
	mux.Handle("GET /api/v1/performance/objectives/consolidated", jwtProtect(mw, perfHandler.GetConsolidatedObjectives))
	mux.Handle("GET /api/v1/performance/objectives/consolidated/paginated", jwtProtect(mw, perfHandler.GetConsolidatedObjectivesPaginated))

	// -- Objective Categories --
	mux.Handle("GET /api/v1/performance/objective-categories", jwtProtect(mw, perfHandler.GetObjectiveCategories))
	mux.Handle("POST /api/v1/performance/objective-categories", jwtProtect(mw, perfHandler.CreateObjectiveCategory))
	mux.Handle("PUT /api/v1/performance/objective-categories", jwtProtect(mw, perfHandler.UpdateObjectiveCategory))

	// -- Category Definitions --
	mux.Handle("GET /api/v1/performance/category-definitions", jwtProtect(mw, perfHandler.GetCategoryDefinitions))
	mux.Handle("POST /api/v1/performance/category-definitions", jwtProtect(mw, perfHandler.CreateCategoryDefinition))
	mux.Handle("PUT /api/v1/performance/category-definitions", jwtProtect(mw, perfHandler.UpdateCategoryDefinition))

	// -- Evaluation Options --
	mux.Handle("GET /api/v1/performance/evaluation-options", jwtProtect(mw, perfHandler.GetEvaluationOptions))
	mux.Handle("POST /api/v1/performance/evaluation-options", jwtProtect(mw, perfHandler.SaveEvaluationOptions))

	// -- Feedback Questionnaires --
	mux.Handle("GET /api/v1/performance/feedback-questionnaires", jwtProtect(mw, perfHandler.GetFeedbackQuestionnaires))
	mux.Handle("POST /api/v1/performance/feedback-questionnaires", jwtProtect(mw, perfHandler.SaveFeedbackQuestionnaires))
	mux.Handle("POST /api/v1/performance/feedback-questionnaire-options", jwtProtect(mw, perfHandler.SaveFeedbackQuestionnaireOptions))

	// -- PMS Competencies --
	mux.Handle("GET /api/v1/performance/competencies", jwtProtect(mw, perfHandler.GetPmsCompetencies))
	mux.Handle("POST /api/v1/performance/competencies", jwtProtect(mw, perfHandler.CreatePmsCompetency))
	mux.Handle("PUT /api/v1/performance/competencies", jwtProtect(mw, perfHandler.UpdatePmsCompetency))

	// -- Work Product Definitions --
	mux.Handle("GET /api/v1/performance/work-product-definitions", jwtProtect(mw, perfHandler.GetObjectiveWorkProductDefinitions))
	mux.Handle("GET /api/v1/performance/work-product-definitions/all", jwtProtect(mw, perfHandler.GetAllWorkProductDefinitions))
	mux.Handle("GET /api/v1/performance/work-product-definitions/paginated", jwtProtect(mw, perfHandler.GetAllPaginatedWorkProductDefinitions))
	mux.Handle("POST /api/v1/performance/work-product-definitions", jwtProtect(mw, perfHandler.SaveWorkProductDefinition))

	// -- Objectives Upload / Activation --
	mux.Handle("POST /api/v1/performance/objectives/upload", jwtProtect(mw, perfHandler.UploadObjectives))
	mux.Handle("POST /api/v1/performance/objectives/deactivate", jwtProtect(mw, perfHandler.DeActivateObjectives))
	mux.Handle("POST /api/v1/performance/objectives/reactivate", jwtProtect(mw, perfHandler.ReActivateObjectives))

	// -- Approve / Reject --
	mux.Handle("POST /api/v1/performance/approve", jwtProtect(mw, perfHandler.ApproveRecords))
	mux.Handle("POST /api/v1/performance/reject", jwtProtect(mw, perfHandler.RejectRecords))

	// -- Performance Enum Select Lists --
	mux.Handle("GET /api/v1/performance/enums/objective-levels", jwtProtect(mw, perfHandler.GetObjectiveLevels))
	mux.Handle("GET /api/v1/performance/enums/extension-target-types", jwtProtect(mw, perfHandler.GetExtensionTargetTypes))
	mux.Handle("GET /api/v1/performance/enums/evaluation-types", jwtProtect(mw, perfHandler.GetEvaluationTypes))
	mux.Handle("GET /api/v1/performance/enums/work-product-types", jwtProtect(mw, perfHandler.GetWorkProductTypes))
	mux.Handle("GET /api/v1/performance/enums/grievance-types", jwtProtect(mw, perfHandler.GetGrievanceTypes))
	mux.Handle("GET /api/v1/performance/enums/feedback-request-types", jwtProtect(mw, perfHandler.GetFeedBackRequestTypes))
	mux.Handle("GET /api/v1/performance/enums/performance-grades", jwtProtect(mw, perfHandler.GetPerformanceGrades))
	mux.Handle("GET /api/v1/performance/enums/review-period-ranges", jwtProtect(mw, perfHandler.GetReviewPeriodRange))
	mux.Handle("GET /api/v1/performance/enums/statuses", jwtProtect(mw, perfHandler.GetStatuses))

	// ----------------------------------------------------------------
	// PMS Management Engine routes — JWT required
	// (PmsEngineHandler has its own RegisterRoutes method)
	// ----------------------------------------------------------------
	pmsEngineHandler := NewPmsEngineHandler(svc, log)
	pmsEngineHandler.RegisterRoutes(mux, mw)

	// ----------------------------------------------------------------
	// Review Period routes — JWT required
	// ----------------------------------------------------------------
	rpHandler := NewReviewPeriodHandler(svc, log)

	// -- Review Period Lifecycle --
	mux.Handle("POST /api/v1/review-periods/draft", jwtProtect(mw, rpHandler.SaveDraftReviewPeriod))
	mux.Handle("POST /api/v1/review-periods", jwtProtect(mw, rpHandler.AddReviewPeriod))
	mux.Handle("POST /api/v1/review-periods/submit-draft", jwtProtect(mw, rpHandler.SubmitDraftReviewPeriod))
	mux.Handle("POST /api/v1/review-periods/approve", jwtProtect(mw, rpHandler.ApproveReviewPeriod))
	mux.Handle("POST /api/v1/review-periods/reject", jwtProtect(mw, rpHandler.RejectReviewPeriod))
	mux.Handle("POST /api/v1/review-periods/return", jwtProtect(mw, rpHandler.ReturnReviewPeriod))
	mux.Handle("POST /api/v1/review-periods/resubmit", jwtProtect(mw, rpHandler.ReSubmitReviewPeriod))
	mux.Handle("PUT /api/v1/review-periods", jwtProtect(mw, rpHandler.UpdateReviewPeriod))
	mux.Handle("POST /api/v1/review-periods/cancel", jwtProtect(mw, rpHandler.CancelReviewPeriod))
	mux.Handle("POST /api/v1/review-periods/close", jwtProtect(mw, rpHandler.CloseReviewPeriod))

	// -- Review Period Toggles --
	mux.Handle("POST /api/v1/review-periods/enable-objective-planning", jwtProtect(mw, rpHandler.EnableObjectivePlanning))
	mux.Handle("POST /api/v1/review-periods/disable-objective-planning", jwtProtect(mw, rpHandler.DisableObjectivePlanning))
	mux.Handle("POST /api/v1/review-periods/enable-work-product-planning", jwtProtect(mw, rpHandler.EnableWorkProductPlanning))
	mux.Handle("POST /api/v1/review-periods/disable-work-product-planning", jwtProtect(mw, rpHandler.DisableWorkProductPlanning))
	mux.Handle("POST /api/v1/review-periods/enable-work-product-evaluation", jwtProtect(mw, rpHandler.EnableWorkProductEvaluation))
	mux.Handle("POST /api/v1/review-periods/disable-work-product-evaluation", jwtProtect(mw, rpHandler.DisableWorkProductEvaluation))

	// -- Review Period Queries --
	mux.Handle("GET /api/v1/review-periods/all", jwtProtect(mw, rpHandler.GetReviewPeriods))
	mux.Handle("GET /api/v1/review-periods/active", jwtProtect(mw, rpHandler.GetActiveReviewPeriod))
	mux.Handle("GET /api/v1/review-periods/staff-active", jwtProtect(mw, rpHandler.GetStaffActiveReviewPeriod))
	mux.Handle("GET /api/v1/review-periods/planned-objective", jwtProtect(mw, rpHandler.GetPlannedObjective))
	mux.Handle("GET /api/v1/review-periods/enterprise-objective", jwtProtect(mw, rpHandler.GetEnterpriseObjectiveByLevel))
	mux.Handle("GET /api/v1/review-periods/objectives-by-status", jwtProtect(mw, rpHandler.GetObjectivesByWorkproductStatus))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/category-definitions", jwtProtect(mw, rpHandler.GetReviewPeriodCategoryDefinitions))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/objectives-with-categories", jwtProtect(mw, rpHandler.GetReviewPeriodObjectivesWithCategoryDefinitions))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/planned-objectives", jwtProtect(mw, rpHandler.GetAllPlannedOperationalObjectives))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}", jwtProtect(mw, rpHandler.GetReviewPeriodDetails))

	// -- Review Period Objectives --
	mux.Handle("POST /api/v1/review-periods/objectives/draft", jwtProtect(mw, rpHandler.SaveDraftReviewPeriodObjective))
	mux.Handle("POST /api/v1/review-periods/objectives", jwtProtect(mw, rpHandler.AddReviewPeriodObjective))
	mux.Handle("POST /api/v1/review-periods/objectives/submit-draft", jwtProtect(mw, rpHandler.SubmitDraftReviewPeriodObjective))
	mux.Handle("POST /api/v1/review-periods/objectives/cancel", jwtProtect(mw, rpHandler.CancelReviewPeriodObjective))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/objectives", jwtProtect(mw, rpHandler.GetReviewPeriodObjectives))

	// -- Review Period Objective Category Definitions --
	mux.Handle("POST /api/v1/review-periods/category-definitions/draft", jwtProtect(mw, rpHandler.SaveDraftReviewPeriodObjectiveCategoryDefinition))
	mux.Handle("POST /api/v1/review-periods/category-definitions", jwtProtect(mw, rpHandler.AddReviewPeriodObjectiveCategoryDefinition))
	mux.Handle("POST /api/v1/review-periods/category-definitions/submit-draft", jwtProtect(mw, rpHandler.SubmitDraftReviewPeriodObjectiveCategoryDefinition))
	mux.Handle("POST /api/v1/review-periods/category-definitions/approve", jwtProtect(mw, rpHandler.ApproveReviewPeriodObjectiveCategoryDefinition))
	mux.Handle("POST /api/v1/review-periods/category-definitions/reject", jwtProtect(mw, rpHandler.RejectReviewPeriodObjectiveCategoryDefinition))

	// -- Review Period Extensions (full lifecycle) --
	mux.Handle("POST /api/v1/review-periods/extensions/draft", jwtProtect(mw, rpHandler.SaveDraftReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions/submit-draft", jwtProtect(mw, rpHandler.SubmitDraftReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions/approve", jwtProtect(mw, rpHandler.ApproveReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions/reject", jwtProtect(mw, rpHandler.RejectReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions/return", jwtProtect(mw, rpHandler.ReturnReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions/resubmit", jwtProtect(mw, rpHandler.ReSubmitReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions/cancel", jwtProtect(mw, rpHandler.CancelReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions/close", jwtProtect(mw, rpHandler.CloseReviewPeriodExtension))
	mux.Handle("PUT /api/v1/review-periods/extensions", jwtProtect(mw, rpHandler.UpdateReviewPeriodExtension))
	mux.Handle("POST /api/v1/review-periods/extensions", jwtProtect(mw, rpHandler.AddReviewPeriodExtension))
	mux.Handle("GET /api/v1/review-periods/extensions/all", jwtProtect(mw, rpHandler.GetAllReviewPeriodExtensions))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/extensions", jwtProtect(mw, rpHandler.GetReviewPeriodExtensions))

	// -- Review Period 360 Reviews --
	mux.Handle("POST /api/v1/review-periods/360-reviews", jwtProtect(mw, rpHandler.AddReviewPeriod360Review))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/360-reviews", jwtProtect(mw, rpHandler.GetReviewPeriod360Reviews))

	// -- Individual Planned Objectives --
	mux.Handle("POST /api/v1/review-periods/individual-objectives/draft", jwtProtect(mw, rpHandler.SaveDraftIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives", jwtProtect(mw, rpHandler.AddIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/submit-draft", jwtProtect(mw, rpHandler.SubmitDraftIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/approve", jwtProtect(mw, rpHandler.ApproveIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/reject", jwtProtect(mw, rpHandler.RejectIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/return", jwtProtect(mw, rpHandler.ReturnIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/cancel", jwtProtect(mw, rpHandler.CancelIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/accept", jwtProtect(mw, rpHandler.AcceptIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/reinstate", jwtProtect(mw, rpHandler.ReInstateIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/pause", jwtProtect(mw, rpHandler.PauseIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/suspend", jwtProtect(mw, rpHandler.SuspendIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/resume", jwtProtect(mw, rpHandler.ResumeIndividualPlannedObjective))
	mux.Handle("POST /api/v1/review-periods/individual-objectives/resubmit", jwtProtect(mw, rpHandler.ReSubmitIndividualPlannedObjective))
	mux.Handle("GET /api/v1/review-periods/individual-objectives", jwtProtect(mw, rpHandler.GetStaffIndividualPlannedObjectives))

	// -- Period Objective Evaluations --
	mux.Handle("POST /api/v1/review-periods/evaluations", jwtProtect(mw, rpHandler.CreatePeriodObjectiveEvaluation))
	mux.Handle("POST /api/v1/review-periods/evaluations/department", jwtProtect(mw, rpHandler.CreatePeriodObjectiveDepartmentEvaluation))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/evaluations", jwtProtect(mw, rpHandler.GetPeriodObjectiveEvaluations))
	mux.Handle("GET /api/v1/review-periods/{reviewPeriodId}/evaluations/department", jwtProtect(mw, rpHandler.GetPeriodObjectiveDepartmentEvaluations))

	// -- Period Scores --
	mux.Handle("GET /api/v1/review-periods/scores", jwtProtect(mw, rpHandler.GetStaffPeriodScore))

	// -- Archive Operations --
	mux.Handle("POST /api/v1/review-periods/archive-objectives", jwtProtect(mw, rpHandler.ArchiveCancelledObjectives))
	mux.Handle("POST /api/v1/review-periods/archive-workproducts", jwtProtect(mw, rpHandler.ArchiveCancelledWorkProducts))

	// ----------------------------------------------------------------
	// Competency Management routes — JWT required
	// ----------------------------------------------------------------
	compHandler := NewCompetencyMgtHandler(svc, log)

	// -- Competencies --
	mux.Handle("GET /api/v1/competency/competencies", jwtProtect(mw, compHandler.GetCompetencies))
	mux.Handle("POST /api/v1/competency/competencies", jwtProtect(mw, compHandler.SaveCompetency))
	mux.Handle("POST /api/v1/competency/competencies/approve", jwtProtect(mw, compHandler.ApproveCompetency))
	mux.Handle("POST /api/v1/competency/competencies/reject", jwtProtect(mw, compHandler.RejectCompetency))

	// -- Competency Categories --
	mux.Handle("GET /api/v1/competency/categories", jwtProtect(mw, compHandler.GetCompetencyCategories))
	mux.Handle("POST /api/v1/competency/categories", jwtProtect(mw, compHandler.SaveCompetencyCategory))

	// -- Competency Category Gradings --
	mux.Handle("GET /api/v1/competency/category-gradings", jwtProtect(mw, compHandler.GetCompetencyCategoryGradings))
	mux.Handle("POST /api/v1/competency/category-gradings", jwtProtect(mw, compHandler.SaveCompetencyCategoryGrading))

	// -- Competency Rating Definitions --
	mux.Handle("GET /api/v1/competency/rating-definitions", jwtProtect(mw, compHandler.GetCompetencyRatingDefinitions))
	mux.Handle("POST /api/v1/competency/rating-definitions", jwtProtect(mw, compHandler.SaveCompetencyRatingDefinition))

	// -- Competency Reviews --
	mux.Handle("GET /api/v1/competency/reviews", jwtProtect(mw, compHandler.GetCompetencyReviews))
	mux.Handle("GET /api/v1/competency/reviews/by-reviewer", jwtProtect(mw, compHandler.GetCompetencyReviewByReviewer))
	mux.Handle("GET /api/v1/competency/reviews/for-employee", jwtProtect(mw, compHandler.GetCompetencyReviewForEmployee))
	mux.Handle("GET /api/v1/competency/reviews/detail", jwtProtect(mw, compHandler.GetCompetencyReviewDetail))
	mux.Handle("POST /api/v1/competency/reviews", jwtProtect(mw, compHandler.SaveCompetencyReview))
	mux.Handle("GET /api/v1/competency/reviews/by-office", jwtProtect(mw, compHandler.GetOfficeCompetencyReviews))

	// -- Competency Review Profiles --
	mux.Handle("GET /api/v1/competency/review-profiles", jwtProtect(mw, compHandler.GetCompetencyReviewProfiles))
	mux.Handle("GET /api/v1/competency/review-profiles/group", jwtProtect(mw, compHandler.GetGroupCompetencyReviewProfiles))
	mux.Handle("GET /api/v1/competency/review-profiles/matrix", jwtProtect(mw, compHandler.GetCompetencyMatrixReviewProfiles))
	mux.Handle("GET /api/v1/competency/review-profiles/technical-matrix", jwtProtect(mw, compHandler.GetTechnicalCompetencyMatrixReviewProfiles))
	mux.Handle("POST /api/v1/competency/review-profiles", jwtProtect(mw, compHandler.SaveCompetencyReviewProfile))

	// -- Development Plans --
	mux.Handle("GET /api/v1/competency/development-plans", jwtProtect(mw, compHandler.GetDevelopmentPlans))
	mux.Handle("POST /api/v1/competency/development-plans", jwtProtect(mw, compHandler.SaveDevelopmentPlan))

	// -- Job Roles --
	mux.Handle("GET /api/v1/competency/job-roles", jwtProtect(mw, compHandler.GetJobRoles))
	mux.Handle("POST /api/v1/competency/job-roles", jwtProtect(mw, compHandler.SaveJobRole))
	mux.Handle("GET /api/v1/competency/office-job-roles", jwtProtect(mw, compHandler.GetOfficeJobRoles))
	mux.Handle("POST /api/v1/competency/office-job-roles", jwtProtect(mw, compHandler.SaveOfficeJobRole))
	mux.Handle("GET /api/v1/competency/job-role-competencies", jwtProtect(mw, compHandler.GetJobRoleCompetencies))
	mux.Handle("POST /api/v1/competency/job-role-competencies", jwtProtect(mw, compHandler.SaveJobRoleCompetency))

	// -- Behavioral Competencies --
	mux.Handle("GET /api/v1/competency/behavioral", jwtProtect(mw, compHandler.GetBehavioralCompetencies))
	mux.Handle("POST /api/v1/competency/behavioral", jwtProtect(mw, compHandler.SaveBehavioralCompetency))

	// -- Job Role Grades --
	mux.Handle("GET /api/v1/competency/job-role-grades", jwtProtect(mw, compHandler.GetJobRoleGrades))
	mux.Handle("POST /api/v1/competency/job-role-grades", jwtProtect(mw, compHandler.SaveJobRoleGrade))

	// -- Job Grades --
	mux.Handle("GET /api/v1/competency/job-grades", jwtProtect(mw, compHandler.GetJobGrades))
	mux.Handle("POST /api/v1/competency/job-grades", jwtProtect(mw, compHandler.SaveJobGrade))

	// -- Job Grade Groups --
	mux.Handle("GET /api/v1/competency/job-grade-groups", jwtProtect(mw, compHandler.GetJobGradeGroups))
	mux.Handle("POST /api/v1/competency/job-grade-groups", jwtProtect(mw, compHandler.SaveJobGradeGroup))
	mux.Handle("GET /api/v1/competency/assign-job-grade-groups", jwtProtect(mw, compHandler.GetAssignJobGradeGroups))
	mux.Handle("POST /api/v1/competency/assign-job-grade-groups", jwtProtect(mw, compHandler.SaveAssignJobGradeGroup))

	// -- Ratings --
	mux.Handle("GET /api/v1/competency/ratings", jwtProtect(mw, compHandler.GetRatings))
	mux.Handle("POST /api/v1/competency/ratings", jwtProtect(mw, compHandler.SaveRating))

	// -- Competency Review Periods --
	mux.Handle("GET /api/v1/competency/review-periods", jwtProtect(mw, compHandler.GetReviewPeriods))
	mux.Handle("POST /api/v1/competency/review-periods", jwtProtect(mw, compHandler.SaveReviewPeriod))
	mux.Handle("POST /api/v1/competency/review-periods/approve", jwtProtect(mw, compHandler.ApproveReviewPeriod))

	// -- Review Types --
	mux.Handle("GET /api/v1/competency/review-types", jwtProtect(mw, compHandler.GetReviewTypes))
	mux.Handle("POST /api/v1/competency/review-types", jwtProtect(mw, compHandler.SaveReviewType))

	// -- Bank Years --
	mux.Handle("GET /api/v1/competency/bank-years", jwtProtect(mw, compHandler.GetBankYears))
	mux.Handle("POST /api/v1/competency/bank-years", jwtProtect(mw, compHandler.SaveBankYear))

	// -- Training Types --
	mux.Handle("GET /api/v1/competency/training-types", jwtProtect(mw, compHandler.GetTrainingTypes))
	mux.Handle("POST /api/v1/competency/training-types", jwtProtect(mw, compHandler.SaveTrainingType))

	// -- Population & Calculation --
	mux.Handle("POST /api/v1/competency/populate/all-reviews", jwtProtect(mw, compHandler.PopulateAllReviews))
	mux.Handle("POST /api/v1/competency/populate/office-reviews", jwtProtect(mw, compHandler.PopulateOfficeReviews))
	mux.Handle("POST /api/v1/competency/populate/division-reviews", jwtProtect(mw, compHandler.PopulateDivisionReviews))
	mux.Handle("POST /api/v1/competency/populate/department-reviews", jwtProtect(mw, compHandler.PopulateDepartmentReviews))
	mux.Handle("POST /api/v1/competency/populate/employee-reviews", jwtProtect(mw, compHandler.PopulateReviewsByEmployeeId))
	mux.Handle("POST /api/v1/competency/calculate-reviews", jwtProtect(mw, compHandler.CalculateReviews))
	mux.Handle("POST /api/v1/competency/recalculate-review-profiles", jwtProtect(mw, compHandler.RecalculateReviewsProfiles))
	mux.Handle("POST /api/v1/competency/email-service", jwtProtect(mw, compHandler.EmailService))
	mux.Handle("POST /api/v1/competency/sync-job-role-soa", jwtProtect(mw, compHandler.SyncJobRoleUpdateSOA))

	// ----------------------------------------------------------------
	// Grievance Management routes — JWT required
	// ----------------------------------------------------------------
	grievanceHandler := NewGrievanceHandler(svc, log)

	mux.Handle("POST /api/v1/grievances", jwtProtect(mw, grievanceHandler.RaiseNewGrievance))
	mux.Handle("PUT /api/v1/grievances", jwtProtect(mw, grievanceHandler.UpdateGrievance))
	mux.Handle("POST /api/v1/grievances/resolution", jwtProtect(mw, grievanceHandler.CreateGrievanceResolution))
	mux.Handle("PUT /api/v1/grievances/resolution", jwtProtect(mw, grievanceHandler.UpdateGrievanceResolution))
	mux.Handle("GET /api/v1/grievances/staff", jwtProtect(mw, grievanceHandler.GetStaffGrievances))
	mux.Handle("GET /api/v1/grievances/report", jwtProtect(mw, grievanceHandler.GetGrievancesReport))

	// ----------------------------------------------------------------
	// PMS Setup routes — Admin only
	// ----------------------------------------------------------------
	setupHandler := NewPmsSetupHandler(svc, log)

	mux.Handle("POST /api/v1/setup/settings", jwtRoleProtect(mw, setupHandler.AddSetting, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("PUT /api/v1/setup/settings", jwtRoleProtect(mw, setupHandler.UpdateSetting, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("GET /api/v1/setup/settings/{settingId}", jwtRoleProtect(mw, setupHandler.GetSettingDetails, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("GET /api/v1/setup/settings", jwtRoleProtect(mw, setupHandler.ListAllSettings, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("POST /api/v1/setup/pms-configurations", jwtRoleProtect(mw, setupHandler.AddPmsConfiguration, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("PUT /api/v1/setup/pms-configurations", jwtRoleProtect(mw, setupHandler.UpdatePmsConfiguration, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("GET /api/v1/setup/pms-configurations/{configId}", jwtRoleProtect(mw, setupHandler.GetPmsConfigurationDetails, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("GET /api/v1/setup/pms-configurations", jwtRoleProtect(mw, setupHandler.ListAllPmsConfigurations, auth.RoleAdmin, auth.RoleSuperAdmin))

	// ----------------------------------------------------------------
	// Organogram routes — JWT required
	// ----------------------------------------------------------------
	orgHandler := NewOrganogramHandler(svc.Organogram, svc.ErpEmployee, log)

	// -- Directorates --
	mux.Handle("GET /api/v1/organogram/directorates", jwtProtect(mw, orgHandler.GetDirectorates))
	mux.Handle("POST /api/v1/organogram/directorates", jwtProtect(mw, orgHandler.SaveDirectorate))

	// -- Departments --
	mux.Handle("GET /api/v1/organogram/departments", jwtProtect(mw, orgHandler.GetDepartments))
	mux.Handle("POST /api/v1/organogram/departments", jwtProtect(mw, orgHandler.SaveDepartment))

	// -- Divisions --
	mux.Handle("GET /api/v1/organogram/divisions", jwtProtect(mw, orgHandler.GetDivisions))
	mux.Handle("POST /api/v1/organogram/divisions", jwtProtect(mw, orgHandler.SaveDivision))

	// -- Offices --
	mux.Handle("GET /api/v1/organogram/offices", jwtProtect(mw, orgHandler.GetOffices))
	mux.Handle("POST /api/v1/organogram/offices", jwtProtect(mw, orgHandler.SaveOffice))

	// -- ERP Organogram --
	mux.Handle("GET /api/v1/organogram/erp/departments", jwtProtect(mw, orgHandler.GetErpDepartments))
	mux.Handle("GET /api/v1/organogram/erp/divisions", jwtProtect(mw, orgHandler.GetErpDivisions))
	mux.Handle("GET /api/v1/organogram/erp/offices", jwtProtect(mw, orgHandler.GetErpOffices))

	// ----------------------------------------------------------------
	// Role Management routes — JWT required
	// ----------------------------------------------------------------
	roleMgtHandler := NewRoleMgtHandler(svc.RoleMgt, log)

	mux.Handle("GET /api/v1/rolemgmt/permissions", jwtProtect(mw, roleMgtHandler.GetPermissions))
	mux.Handle("GET /api/v1/rolemgmt/roles-with-permission", jwtProtect(mw, roleMgtHandler.GetAllRolesWithPermission))
	mux.Handle("POST /api/v1/rolemgmt/permissions", jwtProtect(mw, roleMgtHandler.AddPermissionToRole))
	mux.Handle("DELETE /api/v1/rolemgmt/permissions", jwtProtect(mw, roleMgtHandler.RemovePermissionFromRole))

	// ----------------------------------------------------------------
	// Staff Management routes — JWT required
	// ----------------------------------------------------------------
	staffsHandler := NewStaffsHandler(svc.StaffMgt, log)

	mux.Handle("POST /api/v1/staff", jwtProtect(mw, staffsHandler.AddStaff))
	mux.Handle("GET /api/v1/staff", jwtProtect(mw, staffsHandler.GetAllStaffs))
	mux.Handle("GET /api/v1/staff/roles", jwtProtect(mw, staffsHandler.GetAllRoles))
	mux.Handle("POST /api/v1/staff/roles", jwtProtect(mw, staffsHandler.AddRole))
	mux.Handle("DELETE /api/v1/staff/roles", jwtProtect(mw, staffsHandler.DeleteRole))
	mux.Handle("POST /api/v1/staff/roles/assign", jwtProtect(mw, staffsHandler.AddStaffToRole))
	mux.Handle("DELETE /api/v1/staff/roles/remove", jwtProtect(mw, staffsHandler.RemoveStaffFromRole))
	mux.Handle("GET /api/v1/staff/roles/by-staff", jwtProtect(mw, staffsHandler.GetStaffRoles))

	// ----------------------------------------------------------------
	// Employee Information routes — JWT required
	// ----------------------------------------------------------------
	empHandler := NewEmployeeInformationHandler(svc.ErpEmployee, log)

	mux.Handle("GET /api/v1/employees", jwtProtect(mw, empHandler.GetEmployeeDetail))
	mux.Handle("GET /api/v1/employees/head-subordinates", jwtProtect(mw, empHandler.GetHeadSubordinates))
	mux.Handle("GET /api/v1/employees/subordinates", jwtProtect(mw, empHandler.GetEmployeeSubordinates))
	mux.Handle("GET /api/v1/employees/peers", jwtProtect(mw, empHandler.GetEmployeePeers))
	mux.Handle("GET /api/v1/employees/by-department", jwtProtect(mw, empHandler.GetAllByDepartment))
	mux.Handle("GET /api/v1/employees/by-division", jwtProtect(mw, empHandler.GetAllByDivision))
	mux.Handle("GET /api/v1/employees/by-office", jwtProtect(mw, empHandler.GetAllByOffice))
	mux.Handle("GET /api/v1/employees/all", jwtProtect(mw, empHandler.GetAllEmployees))
	mux.Handle("GET /api/v1/employees/seed-organization", jwtRoleProtect(mw, empHandler.SeedOrganizationData, auth.RoleAdmin, auth.RoleSuperAdmin))
	mux.Handle("GET /api/v1/employees/staff-id-mask", jwtProtect(mw, empHandler.GetStaffIDMaskDetail))
	mux.Handle("POST /api/v1/employees/staff-job-role", jwtProtect(mw, empHandler.UpdateStaffJobRole))
	mux.Handle("GET /api/v1/employees/staff-job-role", jwtProtect(mw, empHandler.GetStaffJobRoleById))
	mux.Handle("POST /api/v1/employees/job-roles-by-office", jwtProtect(mw, empHandler.GetJobRolesByOffice))
	mux.Handle("GET /api/v1/employees/staff-job-role-requests", jwtProtect(mw, empHandler.GetStaffJobRoleRequests))
	mux.Handle("POST /api/v1/employees/approve-reject-staff-job-role", jwtProtect(mw, empHandler.ApproveRejectStaffJobRole))

	// ----------------------------------------------------------------
	// Enum Lists routes — JWT required (standalone functions)
	// ----------------------------------------------------------------
	mux.Handle("GET /api/v1/enums/objective-levels", jwtProtect(mw, GetObjectiveLevels))
	mux.Handle("GET /api/v1/enums/extension-target-types", jwtProtect(mw, GetExtensionTargetTypes))
	mux.Handle("GET /api/v1/enums/evaluation-types", jwtProtect(mw, GetEvaluationTypes))
	mux.Handle("GET /api/v1/enums/work-product-types", jwtProtect(mw, GetWorkProductTypes))
	mux.Handle("GET /api/v1/enums/grievance-types", jwtProtect(mw, GetGrievanceTypes))
	mux.Handle("GET /api/v1/enums/feedback-request-types", jwtProtect(mw, GetFeedbackRequestTypes))
	mux.Handle("GET /api/v1/enums/performance-grades", jwtProtect(mw, GetPerformanceGrades))
	mux.Handle("GET /api/v1/enums/review-period-ranges", jwtProtect(mw, GetReviewPeriodRanges))
	mux.Handle("GET /api/v1/enums/statuses", jwtProtect(mw, GetStatuses))

	// Apply global middleware chain (outermost first)
	var handler http.Handler = mux
	handler = mw.APIKeyAuth(handler)
	handler = mw.SecurityHeaders(handler)
	handler = mw.CORS(handler)
	handler = mw.Recover(handler)
	handler = mw.RequestLogger(handler)

	return handler
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "pms-api",
	})
}
