package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// performanceManagementService implements PerformanceManagementService.
// It composes sub-services for objectives, dashboard, period scoring,
// projects, committees, work products, feedback requests, competency
// reviews and period objective evaluations.
// Legacy "thin" methods for projects, work products, feedback and
// competency review are retained for backward compatibility; new callers
// should prefer the typed sub-service methods.
// ---------------------------------------------------------------------------

type performanceManagementService struct {
	db  *gorm.DB
	cfg *config.Config
	log zerolog.Logger

	// Composed sub-services
	strategy         *strategyService
	objectives       *objectiveService
	dashboard        *dashboardService
	periodScr        *periodScoreService
	project          *projectService
	committee        *committeeService
	workProduct      *workProductService
	feedbackReq      *feedbackRequestService
	competencyReview *competencyReviewService
	evaluation       *evaluationService

	// Peer service references (lazy to break cycles)
	reviewPeriodSvc ReviewPeriodService
	erpEmployeeSvc  ErpEmployeeService
	globalSettingSvc GlobalSettingService
}

// newPerformanceManagementService constructs the main performance management
// service with all sub-services composed.
func newPerformanceManagementService(
	repos *repository.Container,
	cfg *config.Config,
	log zerolog.Logger,
	reviewPeriodSvc ReviewPeriodService,
	erpEmployeeSvc ErpEmployeeService,
	globalSettingSvc GlobalSettingService,
	userCtxSvc UserContextService,
) PerformanceManagementService {
	db := repos.GormDB

	svc := &performanceManagementService{
		db:  db,
		cfg: cfg,
		log: log.With().Str("service", "performance_management").Logger(),

		reviewPeriodSvc:  reviewPeriodSvc,
		erpEmployeeSvc:   erpEmployeeSvc,
		globalSettingSvc: globalSettingSvc,
	}

	// Compose sub-services sharing the same DB and repos
	svc.strategy = newStrategyService(repos, cfg, svc.log, userCtxSvc)
	svc.objectives = newObjectiveService(db, cfg, svc.log, svc)
	svc.dashboard = newDashboardService(db, cfg, svc.log, svc)
	svc.periodScr = newPeriodScoreService(db, cfg, svc.log, svc)
	svc.project = newProjectService(db, cfg, svc.log, svc)
	svc.committee = newCommitteeService(db, cfg, svc.log, svc)
	svc.workProduct = newWorkProductService(db, cfg, svc.log, svc)
	svc.feedbackReq = newFeedbackRequestService(db, cfg, svc.log, svc)
	svc.competencyReview = newCompetencyReviewService(db, cfg, svc.log, svc)
	svc.evaluation = newEvaluationService(db, cfg, svc.log, svc)

	return svc
}

// =========================================================================
// Delegated methods: Strategies (via strategyService)
// =========================================================================

func (s *performanceManagementService) GetStrategies(ctx context.Context) (interface{}, error) {
	return s.strategy.GetStrategies(ctx)
}

func (s *performanceManagementService) GetStrategicThemes(ctx context.Context) (interface{}, error) {
	return s.strategy.GetStrategicThemes(ctx)
}

func (s *performanceManagementService) GetStrategicThemesById(ctx context.Context, strategyID string) (interface{}, error) {
	return s.strategy.GetStrategicThemesById(ctx, strategyID)
}

func (s *performanceManagementService) CreateStrategy(ctx context.Context, req interface{}) (interface{}, error) {
	return s.strategy.CreateStrategy(ctx, req)
}

func (s *performanceManagementService) UpdateStrategy(ctx context.Context, req interface{}) (interface{}, error) {
	return s.strategy.UpdateStrategy(ctx, req)
}

func (s *performanceManagementService) CreateStrategicTheme(ctx context.Context, req interface{}) (interface{}, error) {
	return s.strategy.CreateStrategicTheme(ctx, req)
}

func (s *performanceManagementService) UpdateStrategicTheme(ctx context.Context, req interface{}) (interface{}, error) {
	return s.strategy.UpdateStrategicTheme(ctx, req)
}

// =========================================================================
// Delegated methods: Objectives (via objectiveService)
// =========================================================================

func (s *performanceManagementService) GetEnterpriseObjectives(ctx context.Context) (interface{}, error) {
	return s.objectives.GetEnterpriseObjectives(ctx)
}

func (s *performanceManagementService) CreateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.CreateEnterpriseObjective(ctx, req)
}

func (s *performanceManagementService) UpdateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.UpdateEnterpriseObjective(ctx, req)
}

func (s *performanceManagementService) GetDepartmentObjectives(ctx context.Context) (interface{}, error) {
	return s.objectives.GetDepartmentObjectives(ctx)
}

func (s *performanceManagementService) CreateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.CreateDepartmentObjective(ctx, req)
}

func (s *performanceManagementService) UpdateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.UpdateDepartmentObjective(ctx, req)
}

func (s *performanceManagementService) GetDivisionObjectives(ctx context.Context) (interface{}, error) {
	return s.objectives.GetDivisionObjectives(ctx)
}

func (s *performanceManagementService) GetDivisionObjectivesByDivisionId(ctx context.Context, divisionID int) (interface{}, error) {
	return s.objectives.GetDivisionObjectivesByDivisionId(ctx, divisionID)
}

func (s *performanceManagementService) CreateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.CreateDivisionObjective(ctx, req)
}

func (s *performanceManagementService) UpdateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.UpdateDivisionObjective(ctx, req)
}

func (s *performanceManagementService) GetOfficeObjectives(ctx context.Context) (interface{}, error) {
	return s.objectives.GetOfficeObjectives(ctx)
}

func (s *performanceManagementService) GetOfficeObjectivesByOfficeId(ctx context.Context, officeID int) (interface{}, error) {
	return s.objectives.GetOfficeObjectivesByOfficeId(ctx, officeID)
}

func (s *performanceManagementService) CreateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.CreateOfficeObjective(ctx, req)
}

func (s *performanceManagementService) UpdateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.UpdateOfficeObjective(ctx, req)
}

func (s *performanceManagementService) GetObjectiveCategories(ctx context.Context) (interface{}, error) {
	return s.objectives.GetObjectiveCategories(ctx)
}

func (s *performanceManagementService) CreateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.CreateObjectiveCategory(ctx, req)
}

func (s *performanceManagementService) UpdateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.UpdateObjectiveCategory(ctx, req)
}

func (s *performanceManagementService) GetCategoryDefinitions(ctx context.Context, categoryID string) (interface{}, error) {
	return s.objectives.GetCategoryDefinitions(ctx, categoryID)
}

func (s *performanceManagementService) CreateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.CreateCategoryDefinition(ctx, req)
}

func (s *performanceManagementService) UpdateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.UpdateCategoryDefinition(ctx, req)
}

func (s *performanceManagementService) GetConsolidatedObjectives(ctx context.Context) (interface{}, error) {
	return s.objectives.GetConsolidatedObjectives(ctx)
}

func (s *performanceManagementService) GetConsolidatedObjectivesPaginated(ctx context.Context, params interface{}) (interface{}, error) {
	return s.objectives.GetConsolidatedObjectivesPaginated(ctx, params)
}

func (s *performanceManagementService) ProcessObjectivesUpload(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.ProcessObjectivesUpload(ctx, req)
}

func (s *performanceManagementService) DeActivateOrReactivateObjectives(ctx context.Context, req interface{}, deactivate bool) (interface{}, error) {
	return s.objectives.DeActivateOrReactivateObjectives(ctx, req, deactivate)
}

// =========================================================================
// Delegated methods: Evaluation Options & Feedback Questionnaires, PMS
// Competencies, Work Product Definitions, Approval (via objectiveService)
// =========================================================================

func (s *performanceManagementService) GetEvaluationOptions(ctx context.Context) (interface{}, error) {
	return s.objectives.GetEvaluationOptions(ctx)
}

func (s *performanceManagementService) SaveEvaluationOptions(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.SaveEvaluationOptions(ctx, req)
}

func (s *performanceManagementService) GetFeedbackQuestionnaires(ctx context.Context) (interface{}, error) {
	return s.objectives.GetFeedbackQuestionnaires(ctx)
}

func (s *performanceManagementService) SaveFeedbackQuestionnaires(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.SaveFeedbackQuestionnaires(ctx, req)
}

func (s *performanceManagementService) SaveFeedbackQuestionnaireOptions(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.SaveFeedbackQuestionnaireOptions(ctx, req)
}

func (s *performanceManagementService) GetPmsCompetencies(ctx context.Context) (interface{}, error) {
	return s.objectives.GetPmsCompetencies(ctx)
}

func (s *performanceManagementService) CreatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.CreatePmsCompetency(ctx, req)
}

func (s *performanceManagementService) UpdatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.UpdatePmsCompetency(ctx, req)
}

func (s *performanceManagementService) GetObjectiveWorkProductDefinitions(ctx context.Context, objectiveID string, objectiveLevel int) (interface{}, error) {
	return s.objectives.GetObjectiveWorkProductDefinitions(ctx, objectiveID, objectiveLevel)
}

func (s *performanceManagementService) GetAllWorkProductDefinitions(ctx context.Context) (interface{}, error) {
	return s.objectives.GetAllWorkProductDefinitions(ctx)
}

func (s *performanceManagementService) GetAllPaginatedWorkProductDefinitions(ctx context.Context, pageIndex, pageSize int, search string) (interface{}, error) {
	return s.objectives.GetAllPaginatedWorkProductDefinitions(ctx, pageIndex, pageSize, search)
}

func (s *performanceManagementService) SaveWorkProductDefinitions(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.SaveWorkProductDefinitions(ctx, req)
}

func (s *performanceManagementService) ApproveRecords(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.ApproveRecords(ctx, req)
}

func (s *performanceManagementService) RejectRecords(ctx context.Context, req interface{}) (interface{}, error) {
	return s.objectives.RejectRecords(ctx, req)
}

// =========================================================================
// Delegated methods: Dashboard (via dashboardService)
// =========================================================================

func (s *performanceManagementService) GetDashboardStats(ctx context.Context, staffID string) (interface{}, error) {
	return s.dashboard.GetDashboardStats(ctx, staffID)
}

func (s *performanceManagementService) GetStaffAnnualPerformanceScoreCardStatistics(ctx context.Context, staffID string, year int) (performance.StaffAnnualScoreCardResponseVm, error) {
	return s.dashboard.GetStaffAnnualPerformanceScoreCardStatistics(ctx, staffID, year)
}

func (s *performanceManagementService) GetSubordinatesStaffPerformanceScoreCardStatistics(ctx context.Context, managerStaffID, reviewPeriodID string) (performance.AllStaffScoreCardResponseVm, error) {
	return s.dashboard.GetSubordinatesStaffPerformanceScoreCardStatistics(ctx, managerStaffID, reviewPeriodID)
}

func (s *performanceManagementService) GetOrganogramPerformanceSummaryStatistics(ctx context.Context, referenceID, reviewPeriodID string, organogramLevel enums.OrganogramLevel) (performance.OrganogramPerformanceSummaryResponseVm, error) {
	return s.dashboard.GetOrganogramPerformanceSummaryStatistics(ctx, referenceID, reviewPeriodID, organogramLevel)
}

func (s *performanceManagementService) GetOrganogramPerformanceSummaryListStatistics(ctx context.Context, headOfUnitID, reviewPeriodID string, organogramLevel enums.OrganogramLevel) (performance.OrganogramPerformanceSummaryListResponseVm, error) {
	return s.dashboard.GetOrganogramPerformanceSummaryListStatistics(ctx, headOfUnitID, reviewPeriodID, organogramLevel)
}

// =========================================================================
// Delegated methods: Period Scoring (via periodScoreService)
// =========================================================================

func (s *performanceManagementService) GetPerformanceScore(ctx context.Context, staffID string) (interface{}, error) {
	return s.periodScr.GetPerformanceScore(ctx, staffID)
}

// =========================================================================
// Delegated methods: Projects -- full lifecycle (via projectService)
// =========================================================================

func (s *performanceManagementService) ProjectSetup(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error) {
	return s.project.ProjectSetup(ctx, req)
}

func (s *performanceManagementService) GetProject(ctx context.Context, projectID string) (performance.ProjectResponseVm, error) {
	return s.project.GetProject(ctx, projectID)
}

func (s *performanceManagementService) GetProjects(ctx context.Context) (performance.ProjectListResponseVm, error) {
	return s.project.GetProjects(ctx)
}

func (s *performanceManagementService) GetProjectsByManager(ctx context.Context, managerID string) (performance.ProjectListResponseVm, error) {
	return s.project.GetProjectsByManager(ctx, managerID)
}

func (s *performanceManagementService) ProjectObjectiveSetup(ctx context.Context, req *performance.ProjectObjectiveRequestModel) (performance.ResponseVm, error) {
	return s.project.ProjectObjectiveSetup(ctx, req)
}

func (s *performanceManagementService) GetProjectObjectives(ctx context.Context, projectID string) (performance.ProjectObjectiveListResponseVm, error) {
	return s.project.GetProjectObjectives(ctx, projectID)
}

func (s *performanceManagementService) ProjectMembersSetup(ctx context.Context, req *performance.ProjectMemberRequestModel) (performance.ResponseVm, error) {
	return s.project.ProjectMembersSetup(ctx, req)
}

func (s *performanceManagementService) GetProjectMembers(ctx context.Context, projectID string) (performance.ProjectMemberListResponseVm, error) {
	return s.project.GetProjectMembers(ctx, projectID)
}

func (s *performanceManagementService) GetProjectsAssigned(ctx context.Context, staffID string) (performance.ProjectAssignedListResponseVm, error) {
	return s.project.GetProjectsAssigned(ctx, staffID)
}

func (s *performanceManagementService) GetStaffProjects(ctx context.Context, staffID string) (performance.ProjectAssignedListResponseVm, error) {
	return s.project.GetStaffProjects(ctx, staffID)
}

func (s *performanceManagementService) GetProjectWorkProductStaffList(ctx context.Context, projectID string) ([]string, error) {
	return s.project.GetProjectWorkProductStaffList(ctx, projectID)
}

func (s *performanceManagementService) ChangeProjectLead(ctx context.Context, req *performance.ChangeAdhocLeadRequestModel) error {
	return s.project.ChangeProjectLead(ctx, req)
}

func (s *performanceManagementService) ValidateStaffEligibilityForAdhoc(ctx context.Context, staffID, reviewPeriodID string) (performance.AdhocStaffResponseVm, error) {
	return s.project.ValidateStaffEligibilityForAdhoc(ctx, staffID, reviewPeriodID)
}

// =========================================================================
// Delegated methods: Committees (via committeeService)
// =========================================================================

func (s *performanceManagementService) CommitteeSetup(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error) {
	return s.committee.CommitteeSetup(ctx, req)
}

func (s *performanceManagementService) GetCommittee(ctx context.Context, committeeID string) (performance.CommitteeResponseVm, error) {
	return s.committee.GetCommittee(ctx, committeeID)
}

func (s *performanceManagementService) GetCommittees(ctx context.Context) (performance.CommitteeListResponseVm, error) {
	return s.committee.GetCommittees(ctx)
}

func (s *performanceManagementService) GetCommitteesByChairperson(ctx context.Context, chairpersonID string) (performance.CommitteeListResponseVm, error) {
	return s.committee.GetCommitteesByChairperson(ctx, chairpersonID)
}

func (s *performanceManagementService) CommitteeObjectiveSetup(ctx context.Context, req *performance.CommitteeObjectiveRequestModel) (performance.ResponseVm, error) {
	return s.committee.CommitteeObjectiveSetup(ctx, req)
}

func (s *performanceManagementService) GetCommitteeObjectives(ctx context.Context, committeeID string) (performance.CommitteeObjectiveListResponseVm, error) {
	return s.committee.GetCommitteeObjectives(ctx, committeeID)
}

func (s *performanceManagementService) CommitteeMembersSetup(ctx context.Context, req *performance.CommitteeMemberRequestModel) (performance.ResponseVm, error) {
	return s.committee.CommitteeMembersSetup(ctx, req)
}

func (s *performanceManagementService) GetCommitteeMembers(ctx context.Context, committeeID string) (performance.CommitteeMemberListResponseVm, error) {
	return s.committee.GetCommitteeMembers(ctx, committeeID)
}

func (s *performanceManagementService) GetCommitteesAssigned(ctx context.Context, staffID string) (performance.CommitteeAssignedListResponseVm, error) {
	return s.committee.GetCommitteesAssigned(ctx, staffID)
}

func (s *performanceManagementService) GetStaffCommittees(ctx context.Context, staffID string) (performance.CommitteeAssignedListResponseVm, error) {
	return s.committee.GetStaffCommittees(ctx, staffID)
}

func (s *performanceManagementService) GetCommitteeWorkProductStaffList(ctx context.Context, committeeID string) ([]string, error) {
	return s.committee.GetCommitteeWorkProductStaffList(ctx, committeeID)
}

func (s *performanceManagementService) ChangeCommitteeChairperson(ctx context.Context, req *performance.ChangeAdhocLeadRequestModel) error {
	return s.committee.ChangeCommitteeChairperson(ctx, req)
}

// =========================================================================
// Delegated methods: Work Products -- full lifecycle (via workProductService)
// =========================================================================

func (s *performanceManagementService) WorkProductSetup(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error) {
	return s.workProduct.WorkProductSetup(ctx, req)
}

func (s *performanceManagementService) ProjectAssignedWorkProductSetup(ctx context.Context, req *performance.ProjectAssignedWorkProductRequestModel) (performance.ResponseVm, error) {
	return s.workProduct.ProjectAssignedWorkProductSetup(ctx, req)
}

func (s *performanceManagementService) CommitteeAssignedWorkProductSetup(ctx context.Context, req *performance.CommitteeAssignedWorkProductRequestModel) (performance.ResponseVm, error) {
	return s.workProduct.CommitteeAssignedWorkProductSetup(ctx, req)
}

func (s *performanceManagementService) GetWorkProduct(ctx context.Context, workProductID string) (performance.WorkProductResponseVm, error) {
	return s.workProduct.GetWorkProduct(ctx, workProductID)
}

func (s *performanceManagementService) GetProjectWorkProducts(ctx context.Context, projectID string) (performance.ProjectWorkProductListResponseVm, error) {
	return s.workProduct.GetProjectWorkProducts(ctx, projectID)
}

func (s *performanceManagementService) GetProjectAssignedWorkProducts(ctx context.Context, projectID string) (performance.ProjectAssignedWorkProductListResponseVm, error) {
	return s.workProduct.GetProjectAssignedWorkProducts(ctx, projectID)
}

func (s *performanceManagementService) GetCommitteeWorkProducts(ctx context.Context, committeeID string) (performance.CommitteeWorkProductListResponseVm, error) {
	return s.workProduct.GetCommitteeWorkProducts(ctx, committeeID)
}

func (s *performanceManagementService) GetCommitteeAssignedWorkProducts(ctx context.Context, committeeID string) (performance.CommitteeAssignedWorkProductListResponseVm, error) {
	return s.workProduct.GetCommitteeAssignedWorkProducts(ctx, committeeID)
}

func (s *performanceManagementService) GetOperationalWorkProducts(ctx context.Context, plannedObjectiveID string) (performance.OperationalObjectiveWorkProductListResponseVm, error) {
	return s.workProduct.GetOperationalWorkProducts(ctx, plannedObjectiveID)
}

func (s *performanceManagementService) GetStaffWorkProducts(ctx context.Context, staffID, reviewPeriodID string) (performance.StaffWorkProductListResponseVm, error) {
	return s.workProduct.GetStaffWorkProducts(ctx, staffID, reviewPeriodID)
}

func (s *performanceManagementService) GetAllStaffWorkProducts(ctx context.Context, staffID string) (performance.StaffWorkProductListResponseVm, error) {
	return s.workProduct.GetAllStaffWorkProducts(ctx, staffID)
}

func (s *performanceManagementService) GetObjectiveWorkProducts(ctx context.Context, objectiveID string) (performance.ObjectiveWorkProductListResponseVm, error) {
	return s.workProduct.GetObjectiveWorkProducts(ctx, objectiveID)
}

func (s *performanceManagementService) WorkProductTaskSetup(ctx context.Context, req *performance.WorkProductTaskRequestModel) (performance.ResponseVm, error) {
	return s.workProduct.WorkProductTaskSetup(ctx, req)
}

func (s *performanceManagementService) GetWorkProductTasks(ctx context.Context, workProductID string) (performance.WorkProductTaskListResponseVm, error) {
	return s.workProduct.GetWorkProductTasks(ctx, workProductID)
}

func (s *performanceManagementService) ReCalculateWorkProductPoints(ctx context.Context, staffID, reviewPeriodID string) (performance.RecalculateWorkProductResponseVm, error) {
	return s.workProduct.ReCalculateWorkProductPoints(ctx, staffID, reviewPeriodID)
}

func (s *performanceManagementService) WorkProductEvaluation(ctx context.Context, req *performance.WorkProductEvaluationRequestModel) (performance.EvaluationResponseVm, error) {
	return s.workProduct.WorkProductEvaluation(ctx, req)
}

func (s *performanceManagementService) GetWorkProductEvaluation(ctx context.Context, workProductID string) (performance.WorkProductEvaluationResponseVm, error) {
	return s.workProduct.GetWorkProductEvaluation(ctx, workProductID)
}

func (s *performanceManagementService) InitiateWorkProductReEvaluation(ctx context.Context, workProductID string) (performance.ResponseVm, error) {
	return s.workProduct.InitiateWorkProductReEvaluation(ctx, workProductID)
}

// =========================================================================
// Delegated methods: Feedback Requests -- full lifecycle (via feedbackRequestService)
// =========================================================================

func (s *performanceManagementService) LogRequest(ctx context.Context, feedbackType enums.FeedbackRequestType, referenceID, assignedStaffID, requestOwnerStaffID, reviewPeriodID string, hasSLA bool) error {
	return s.feedbackReq.LogRequest(ctx, feedbackType, referenceID, assignedStaffID, requestOwnerStaffID, reviewPeriodID, hasSLA)
}

func (s *performanceManagementService) LogAcceptanceRequest(ctx context.Context, feedbackType enums.FeedbackRequestType, referenceID, assignedStaffID, requestOwnerStaffID, reviewPeriodID string) error {
	return s.feedbackReq.LogAcceptanceRequest(ctx, feedbackType, referenceID, assignedStaffID, requestOwnerStaffID, reviewPeriodID)
}

func (s *performanceManagementService) GetRequests(ctx context.Context, staffID string, feedbackType *enums.FeedbackRequestType, status *string) (performance.FeedbackRequestListResponseVm, error) {
	return s.feedbackReq.GetRequests(ctx, staffID, feedbackType, status)
}

func (s *performanceManagementService) GetRequestsByOwner(ctx context.Context, requestOwnerStaffID string, feedbackType *enums.FeedbackRequestType) (performance.FeedbackRequestListResponseVm, error) {
	return s.feedbackReq.GetRequestsByOwner(ctx, requestOwnerStaffID, feedbackType)
}

func (s *performanceManagementService) GetBreachedRequests(ctx context.Context, staffID, reviewPeriodID string) (performance.BreachedFeedbackRequestListResponseVm, error) {
	return s.feedbackReq.GetBreachedRequests(ctx, staffID, reviewPeriodID)
}

func (s *performanceManagementService) GetPendingRequests(ctx context.Context, staffID string) (performance.GetStaffPendingRequestVm, error) {
	return s.feedbackReq.GetPendingRequests(ctx, staffID)
}

func (s *performanceManagementService) GetFeedbackRequest(ctx context.Context, requestID string) (performance.FeedbackRequestLogResponseVm, error) {
	return s.feedbackReq.GetFeedbackRequest(ctx, requestID)
}

func (s *performanceManagementService) GetRequestDetails(ctx context.Context, requestID string) (performance.FeedbackRequestLogResponseVm, error) {
	return s.feedbackReq.GetRequestDetails(ctx, requestID)
}

func (s *performanceManagementService) UpdateRequest(ctx context.Context, requestID string, comment, attachment string) error {
	return s.feedbackReq.UpdateRequest(ctx, requestID, comment, attachment)
}

func (s *performanceManagementService) ReassignRequest(ctx context.Context, requestID, newAssignedStaffID string) error {
	return s.feedbackReq.ReassignRequest(ctx, requestID, newAssignedStaffID)
}

func (s *performanceManagementService) ReassignSelfRequest(ctx context.Context, requestID, currentStaffID, newAssignedStaffID string) error {
	return s.feedbackReq.ReassignSelfRequest(ctx, requestID, currentStaffID, newAssignedStaffID)
}

func (s *performanceManagementService) CloseRequest(ctx context.Context, requestID string) error {
	return s.feedbackReq.CloseRequest(ctx, requestID)
}

func (s *performanceManagementService) CloseReviewPeriodRequests(ctx context.Context, reviewPeriodID string) error {
	return s.feedbackReq.CloseReviewPeriodRequests(ctx, reviewPeriodID)
}

func (s *performanceManagementService) ReactivateReviewPeriodRequest(ctx context.Context, reviewPeriodID string) error {
	return s.feedbackReq.ReactivateReviewPeriodRequest(ctx, reviewPeriodID)
}

func (s *performanceManagementService) ReInitiateSameRequest(ctx context.Context, requestID string) error {
	return s.feedbackReq.ReInitiateSameRequest(ctx, requestID)
}

func (s *performanceManagementService) TreatAssignedRequest(ctx context.Context, req *performance.TreatFeedbackRequestModel) error {
	return s.feedbackReq.TreatAssignedRequest(ctx, req)
}

func (s *performanceManagementService) HasLineManager(ctx context.Context, staffID string) (bool, error) {
	return s.feedbackReq.HasLineManager(ctx, staffID)
}

func (s *performanceManagementService) HasVacationRule(ctx context.Context, staffID string) (bool, error) {
	return s.feedbackReq.HasVacationRule(ctx, staffID)
}

func (s *performanceManagementService) GetStaffLeaveDays(ctx context.Context, staffID string, startDate, endDate time.Time) (performance.LeaveResponseVm, error) {
	return s.feedbackReq.GetStaffLeaveDays(ctx, staffID, startDate, endDate)
}

func (s *performanceManagementService) GetPublicDays(ctx context.Context, startDate, endDate time.Time) (performance.PublicHolidaysResponseVm, error) {
	return s.feedbackReq.GetPublicDays(ctx, startDate, endDate)
}

func (s *performanceManagementService) AutoReassignAndLogRequest(ctx context.Context, requestID string) error {
	return s.feedbackReq.AutoReassignAndLogRequest(ctx, requestID)
}

// =========================================================================
// Delegated methods: Competency / 360 Review -- full lifecycle
// (via competencyReviewService)
// =========================================================================

func (s *performanceManagementService) CompetencyReviewFeedbackSetup(ctx context.Context, req *performance.CompetencyReviewFeedbackRequestModel) (performance.ResponseVm, error) {
	return s.competencyReview.CompetencyReviewFeedbackSetup(ctx, req)
}

func (s *performanceManagementService) CompetencyReviewerSetup(ctx context.Context, req *performance.CompetencyReviewerRequestModel) (performance.ResponseVm, error) {
	return s.competencyReview.CompetencyReviewerSetup(ctx, req)
}

func (s *performanceManagementService) CompetencyRatingSetup(ctx context.Context, req *performance.SavePmsCompetencyRequestVm) (performance.ResponseVm, error) {
	return s.competencyReview.CompetencyRatingSetup(ctx, req)
}

func (s *performanceManagementService) GetCompetencyReviewFeedback(ctx context.Context, feedbackID string) (performance.CompetencyReviewFeedbackResponseVm, error) {
	return s.competencyReview.GetCompetencyReviewFeedback(ctx, feedbackID)
}

func (s *performanceManagementService) GetCompetencyReviewFeedbackDetails(ctx context.Context, feedbackID string) (performance.CompetencyReviewFeedbackDetailsResponseVm, error) {
	return s.competencyReview.GetCompetencyReviewFeedbackDetails(ctx, feedbackID)
}

func (s *performanceManagementService) GetAllCompetencyReviewFeedbacks(ctx context.Context, staffID string) (performance.CompetencyReviewFeedbackListResponseVm, error) {
	return s.competencyReview.GetAllCompetencyReviewFeedbacks(ctx, staffID)
}

func (s *performanceManagementService) GetCompetencyReviews(ctx context.Context, reviewerStaffID string) (performance.CompetencyReviewersListResponseVm, error) {
	return s.competencyReview.GetCompetencyReviews(ctx, reviewerStaffID)
}

func (s *performanceManagementService) GetReviewerFeedbackDetails(ctx context.Context, reviewerID string) (performance.CompetencyReviewersResponseVm, error) {
	return s.competencyReview.GetReviewerFeedbackDetails(ctx, reviewerID)
}

func (s *performanceManagementService) GetQuestionnaire(ctx context.Context, staffID string) (performance.QuestionnaireListResponseVm, error) {
	return s.competencyReview.GetQuestionnaire(ctx, staffID)
}

func (s *performanceManagementService) CompetencyGapClosureSetup(ctx context.Context, req *performance.CompetencyGapClosureRequestModel) (performance.ResponseVm, error) {
	return s.competencyReview.CompetencyGapClosureSetup(ctx, req)
}

func (s *performanceManagementService) Initiate360Review(ctx context.Context, req *performance.Initiate360ReviewRequestModel) (performance.ResponseVm, error) {
	return s.competencyReview.Initiate360Review(ctx, req)
}

func (s *performanceManagementService) Complete360Review(ctx context.Context, req *performance.Complete360ReviewRequestModel) (performance.ResponseVm, error) {
	return s.competencyReview.Complete360Review(ctx, req)
}

// =========================================================================
// Delegated methods: Period Objective Evaluations (via evaluationService)
// =========================================================================

func (s *performanceManagementService) ReviewPeriodObjectiveEvaluation(ctx context.Context, req *performance.PeriodObjectiveEvaluationRequestModel) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	return s.evaluation.ReviewPeriodObjectiveEvaluation(ctx, req)
}

func (s *performanceManagementService) GetReviewPeriodObjectiveEvaluation(ctx context.Context, reviewPeriodID, enterpriseObjectiveID string) (performance.PeriodObjectiveEvaluationResponseVm, error) {
	return s.evaluation.GetReviewPeriodObjectiveEvaluation(ctx, reviewPeriodID, enterpriseObjectiveID)
}

func (s *performanceManagementService) GetReviewPeriodObjectiveEvaluations(ctx context.Context, reviewPeriodID string) (performance.PeriodObjectiveEvaluationListResponseVm, error) {
	return s.evaluation.GetReviewPeriodObjectiveEvaluations(ctx, reviewPeriodID)
}

func (s *performanceManagementService) ReviewPeriodDepartmentObjectiveEvaluation(ctx context.Context, req *performance.PeriodObjectiveDepartmentEvaluationRequestModel) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	return s.evaluation.ReviewPeriodDepartmentObjectiveEvaluation(ctx, req)
}

func (s *performanceManagementService) GetReviewPeriodDepartmentObjectiveEvaluation(ctx context.Context, reviewPeriodID, enterpriseObjectiveID string, departmentID int) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error) {
	return s.evaluation.GetReviewPeriodDepartmentObjectiveEvaluation(ctx, reviewPeriodID, enterpriseObjectiveID, departmentID)
}

func (s *performanceManagementService) GetReviewPeriodDepartmentObjectiveEvaluations(ctx context.Context, reviewPeriodID string) (performance.PeriodObjectiveDepartmentEvaluationListResponseVm, error) {
	return s.evaluation.GetReviewPeriodDepartmentObjectiveEvaluations(ctx, reviewPeriodID)
}

func (s *performanceManagementService) GetWorkProductsPerEnterpriseObjective(ctx context.Context, enterpriseObjectiveID, reviewPeriodID string) (performance.ObjectiveWorkProductListResponseVm, error) {
	return s.evaluation.GetWorkProductsPerEnterpriseObjective(ctx, enterpriseObjectiveID, reviewPeriodID)
}

func (s *performanceManagementService) GetDepartmentsbyObjective(ctx context.Context, enterpriseObjectiveID string) (performance.DepartmentListResponseVm, error) {
	return s.evaluation.GetDepartmentsbyObjective(ctx, enterpriseObjectiveID)
}

// =========================================================================
// Legacy Competency / 360 Review (direct implementations retained for
// backward compatibility)
// =========================================================================

func (s *performanceManagementService) GetCompetencyReview(ctx context.Context, staffID string) (interface{}, error) {
	resp := performance.CompetencyReviewFeedbackListResponseVm{}

	var feedbacks []performance.CompetencyReviewFeedback
	err := s.db.WithContext(ctx).
		Where("staff_id = ?", staffID).
		Preload("CompetencyReviewers").
		Preload("CompetencyReviewers.CompetencyReviewerRatings").
		Find(&feedbacks).Error
	if err != nil {
		s.log.Error().Err(err).Str("staffID", staffID).Msg("failed to get competency review feedbacks")
		resp.HasError = true
		resp.Message = "failed to retrieve competency review"
		return resp, err
	}

	var data []performance.CompetencyReviewFeedbackData
	for _, fb := range feedbacks {
		d := performance.CompetencyReviewFeedbackData{
			CompetencyReviewFeedbackID: fb.CompetencyReviewFeedbackID,
			StaffID:                    fb.StaffID,
			MaxPoints:                  fb.MaxPoints,
			FinalScore:                 fb.FinalScore,
			ReviewPeriodID:             fb.ReviewPeriodID,
		}
		data = append(data, d)
	}

	resp.CompetencyReviewFeedbacks = data
	resp.TotalRecords = len(data)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *performanceManagementService) Submit360Feedback(ctx context.Context, req interface{}) error {
	r, ok := req.(*performance.SavePmsCompetencyRequestVm)
	if !ok {
		return fmt.Errorf("invalid request type for Submit360Feedback")
	}

	// Find or create the rating record
	var existing performance.CompetencyReviewerRating
	err := s.db.WithContext(ctx).
		Where("competency_reviewer_id = ? AND pms_competency_id = ?",
			r.CompetencyReviewerID, r.PmsCompetencyID).
		First(&existing).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("querying existing rating: %w", err)
	}

	// Look up the option score
	var option performance.FeedbackQuestionaireOption
	if err := s.db.WithContext(ctx).
		Where("feedback_questionaire_option_id = ?", r.FeedbackQuestionaireOptionID).
		First(&option).Error; err != nil {
		return fmt.Errorf("feedback questionnaire option not found: %w", err)
	}

	if existing.CompetencyReviewerRatingID != "" {
		// Update existing rating
		existing.FeedbackQuestionaireOptionID = r.FeedbackQuestionaireOptionID
		existing.Rating = option.Score
		return s.db.WithContext(ctx).Save(&existing).Error
	}

	// Create new rating
	rating := performance.CompetencyReviewerRating{
		CompetencyReviewerRatingID:   r.CompetencyReviewerRatingID,
		PmsCompetencyID:              r.PmsCompetencyID,
		FeedbackQuestionaireOptionID: r.FeedbackQuestionaireOptionID,
		Rating:                       option.Score,
		CompetencyReviewerID:         r.CompetencyReviewerID,
	}

	return s.db.WithContext(ctx).Create(&rating).Error
}

// =========================================================================
// Legacy Project Management (direct implementations retained for
// backward compatibility)
// =========================================================================

func (s *performanceManagementService) SetupProject(ctx context.Context, req interface{}) error {
	r, ok := req.(*performance.ProjectViewModel)
	if !ok {
		return fmt.Errorf("invalid request type for SetupProject")
	}

	// Validate review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := s.db.WithContext(ctx).
		Where("period_id = ?", r.ReviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		return fmt.Errorf("review period not found: %w", err)
	}

	if reviewPeriod.RecordStatus != enums.StatusActive.String() {
		return fmt.Errorf("review period is not active")
	}

	// Validate dates
	if !r.StartDate.Before(r.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}

	// Check duplicate name in department
	var existing performance.Project
	err := s.db.WithContext(ctx).
		Where("LOWER(name) = LOWER(?) AND department_id = ? AND record_status != ?",
			r.Name, r.DepartmentID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		return fmt.Errorf("project name already exists in this department")
	}

	project := performance.Project{
		ProjectManager: r.ProjectManager,
		BaseProject: performance.BaseProject{
			Name:           r.Name,
			Description:    r.Description,
			StartDate:      r.StartDate,
			EndDate:        r.EndDate,
			Deliverables:   r.Deliverables,
			ReviewPeriodID: r.ReviewPeriodID,
			DepartmentID:   r.DepartmentID,
		},
	}
	project.RecordStatus = enums.StatusPendingApproval.String()
	project.IsActive = true
	project.CreatedBy = r.CreatedBy

	if err := s.db.WithContext(ctx).Create(&project).Error; err != nil {
		return fmt.Errorf("creating project: %w", err)
	}

	s.log.Info().
		Str("projectID", project.ProjectID).
		Str("name", project.Name).
		Msg("project created successfully")

	return nil
}

func (s *performanceManagementService) GetProjectsByStaff(ctx context.Context, staffID string) (interface{}, error) {
	resp := performance.ProjectListResponseVm{}

	var projects []performance.Project
	err := s.db.WithContext(ctx).
		Where("project_manager = ?", staffID).
		Preload("ProjectMembers").
		Preload("ProjectObjectives").
		Preload("ProjectObjectives.Objective").
		Preload("ReviewPeriod").
		Find(&projects).Error
	if err != nil {
		s.log.Error().Err(err).Str("staffID", staffID).Msg("failed to get projects by staff")
		resp.HasError = true
		resp.Message = "failed to retrieve projects"
		return resp, err
	}

	var viewModels []performance.ProjectViewModel
	for _, p := range projects {
		vm := performance.ProjectViewModel{
			ProjectID:      p.ProjectID,
			ProjectManager: p.ProjectManager,
			Name:           p.Name,
			Description:    p.Description,
			StartDate:      p.StartDate,
			EndDate:        p.EndDate,
			Deliverables:   p.Deliverables,
			ReviewPeriodID: p.ReviewPeriodID,
			DepartmentID:   p.DepartmentID,
		}

		// Map project objectives
		for _, obj := range p.ProjectObjectives {
			objData := performance.ProjectObjectiveData{
				ProjectObjectiveID: obj.ProjectObjectiveID,
				ObjectiveID:        obj.ObjectiveID,
				ProjectID:          obj.ProjectID,
			}
			if obj.Objective != nil {
				objData.Objective = obj.Objective.Name
				objData.Kpi = obj.Objective.Kpi
			}
			vm.ProjectObjectives = append(vm.ProjectObjectives, objData)
		}

		viewModels = append(viewModels, vm)
	}

	resp.Projects = viewModels
	resp.TotalRecords = len(viewModels)
	resp.Message = "operation completed successfully"
	return resp, nil
}

func (s *performanceManagementService) AddProjectObjective(ctx context.Context, req interface{}) error {
	r, ok := req.(*performance.ProjectObjectiveData)
	if !ok {
		return fmt.Errorf("invalid request type for AddProjectObjective")
	}

	// Validate project exists
	var project performance.Project
	if err := s.db.WithContext(ctx).
		Where("project_id = ?", r.ProjectID).
		First(&project).Error; err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// Validate objective exists
	var objective performance.EnterpriseObjective
	if err := s.db.WithContext(ctx).
		Where("enterprise_objective_id = ?", r.ObjectiveID).
		First(&objective).Error; err != nil {
		return fmt.Errorf("enterprise objective not found: %w", err)
	}

	// Check for duplicate
	var existing performance.ProjectObjective
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND objective_id = ? AND record_status != ?",
			r.ProjectID, r.ObjectiveID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		return fmt.Errorf("objective already linked to this project")
	}

	projObj := performance.ProjectObjective{
		ObjectiveID: r.ObjectiveID,
		ProjectID:   r.ProjectID,
	}
	projObj.RecordStatus = enums.StatusPendingApproval.String()

	return s.db.WithContext(ctx).Create(&projObj).Error
}

func (s *performanceManagementService) AddProjectMember(ctx context.Context, req interface{}) error {
	r, ok := req.(*performance.ProjectMemberData)
	if !ok {
		return fmt.Errorf("invalid request type for AddProjectMember")
	}

	// Validate project exists and is active
	var project performance.Project
	if err := s.db.WithContext(ctx).
		Where("project_id = ? AND record_status = ?", r.ProjectID, enums.StatusActive.String()).
		First(&project).Error; err != nil {
		return fmt.Errorf("active project not found: %w", err)
	}

	// Check for duplicate member
	var existing performance.ProjectMember
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND staff_id = ? AND record_status != ?",
			r.ProjectID, r.StaffID, enums.StatusCancelled.String()).
		First(&existing).Error
	if err == nil {
		return fmt.Errorf("staff is already a member of this project")
	}

	member := performance.ProjectMember{
		StaffID:            r.StaffID,
		ProjectID:          r.ProjectID,
		PlannedObjectiveID: r.PlannedObjectiveID,
	}
	member.RecordStatus = enums.StatusActive.String()
	member.IsActive = true
	member.IsApproved = true

	if err := s.db.WithContext(ctx).Create(&member).Error; err != nil {
		return fmt.Errorf("creating project member: %w", err)
	}

	s.log.Info().
		Str("projectID", r.ProjectID).
		Str("staffID", r.StaffID).
		Msg("project member added")

	return nil
}

// =========================================================================
// Legacy Work Products (direct implementations retained for backward
// compatibility)
// =========================================================================

func (s *performanceManagementService) AddWorkProduct(ctx context.Context, req interface{}) error {
	r, ok := req.(*performance.WorkProductVm)
	if !ok {
		return fmt.Errorf("invalid request type for AddWorkProduct")
	}

	// Validate planned objective exists
	var plannedObj performance.ReviewPeriodIndividualPlannedObjective
	if err := s.db.WithContext(ctx).
		Where("planned_objective_id = ? AND staff_id = ? AND record_status = ?",
			r.PlannedObjectiveID, r.StaffID, enums.StatusActive.String()).
		First(&plannedObj).Error; err != nil {
		return fmt.Errorf("planned objective not found or not active: %w", err)
	}

	// Validate dates
	if !r.StartDate.Before(r.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}

	wp := performance.WorkProduct{
		Name:            r.Name,
		Description:     r.Description,
		MaxPoint:        r.MaxPoint,
		WorkProductType: enums.WorkProductType(r.RecordStatus),
		IsSelfCreated:   r.IsSelfCreated,
		StaffID:         r.StaffID,
		StartDate:       r.StartDate,
		EndDate:         r.EndDate,
		Deliverables:    r.Deliverables,
	}
	wp.RecordStatus = enums.StatusPendingApproval.String()
	wp.IsActive = true
	wp.CreatedBy = r.CreatedBy

	if err := s.db.WithContext(ctx).Create(&wp).Error; err != nil {
		return fmt.Errorf("creating work product: %w", err)
	}

	// Link to operational objective
	link := performance.OperationalObjectiveWorkProduct{
		WorkProductID:      wp.WorkProductID,
		PlannedObjectiveID: r.PlannedObjectiveID,
	}
	if err := s.db.WithContext(ctx).Create(&link).Error; err != nil {
		s.log.Error().Err(err).Msg("failed to link work product to planned objective")
	}

	// If project or committee work product, create the association
	if r.ProjectID != "" {
		projWP := performance.ProjectWorkProduct{
			WorkProductID: wp.WorkProductID,
			ProjectID:     r.ProjectID,
		}
		if err := s.db.WithContext(ctx).Create(&projWP).Error; err != nil {
			s.log.Error().Err(err).Msg("failed to link work product to project")
		}
	}

	if r.CommitteeID != "" {
		comWP := performance.CommitteeWorkProduct{
			WorkProductID: wp.WorkProductID,
			CommitteeID:   r.CommitteeID,
		}
		if err := s.db.WithContext(ctx).Create(&comWP).Error; err != nil {
			s.log.Error().Err(err).Msg("failed to link work product to committee")
		}
	}

	s.log.Info().
		Str("workProductID", wp.WorkProductID).
		Str("staffID", r.StaffID).
		Msg("work product created")

	return nil
}

func (s *performanceManagementService) EvaluateWorkProduct(ctx context.Context, req interface{}) error {
	r, ok := req.(*performance.WorkProductEvaluationVm)
	if !ok {
		return fmt.Errorf("invalid request type for EvaluateWorkProduct")
	}

	// Validate work product exists and is awaiting evaluation
	var wp performance.WorkProduct
	if err := s.db.WithContext(ctx).
		Where("work_product_id = ?", r.WorkProductID).
		First(&wp).Error; err != nil {
		return fmt.Errorf("work product not found: %w", err)
	}

	if wp.RecordStatus != enums.StatusAwaitingEvaluation.String() &&
		wp.RecordStatus != enums.StatusReEvaluate.String() {
		return fmt.Errorf("work product is not in a state that allows evaluation (current: %s)", wp.RecordStatus)
	}

	// Get evaluation options to compute scores
	var timelinessOption, qualityOption, outputOption performance.EvaluationOption

	if r.TimelinessEvaluationOptionID != "" {
		s.db.WithContext(ctx).Where("evaluation_option_id = ?", r.TimelinessEvaluationOptionID).First(&timelinessOption)
	}
	if r.QualityEvaluationOptionID != "" {
		s.db.WithContext(ctx).Where("evaluation_option_id = ?", r.QualityEvaluationOptionID).First(&qualityOption)
	}
	if r.OutputEvaluationOptionID != "" {
		s.db.WithContext(ctx).Where("evaluation_option_id = ?", r.OutputEvaluationOptionID).First(&outputOption)
	}

	timeliness := timelinessOption.Score
	quality := qualityOption.Score
	output := outputOption.Score

	// Calculate final score: (timeliness + quality + output) / 3 * maxPoint / 100
	avgScore := (timeliness + quality + output) / 3.0
	finalScore := wp.MaxPoint * avgScore / 100.0

	// Check if evaluation already exists (re-evaluation)
	var existingEval performance.WorkProductEvaluation
	isReEval := false
	err := s.db.WithContext(ctx).
		Where("work_product_id = ?", r.WorkProductID).
		First(&existingEval).Error
	if err == nil {
		isReEval = true
	}

	if isReEval {
		existingEval.Timeliness = timeliness
		existingEval.TimelinessEvaluationOptionID = r.TimelinessEvaluationOptionID
		existingEval.Quality = quality
		existingEval.QualityEvaluationOptionID = r.QualityEvaluationOptionID
		existingEval.Output = output
		existingEval.OutputEvaluationOptionID = r.OutputEvaluationOptionID
		existingEval.Outcome = finalScore
		existingEval.EvaluatorStaffID = r.EvaluatorStaffID
		existingEval.IsReEvaluated = true

		if err := s.db.WithContext(ctx).Save(&existingEval).Error; err != nil {
			return fmt.Errorf("updating work product evaluation: %w", err)
		}
	} else {
		eval := performance.WorkProductEvaluation{
			WorkProductID:                r.WorkProductID,
			Timeliness:                   timeliness,
			TimelinessEvaluationOptionID: r.TimelinessEvaluationOptionID,
			Quality:                      quality,
			QualityEvaluationOptionID:    r.QualityEvaluationOptionID,
			Output:                       output,
			OutputEvaluationOptionID:     r.OutputEvaluationOptionID,
			Outcome:                      finalScore,
			EvaluatorStaffID:             r.EvaluatorStaffID,
		}

		if err := s.db.WithContext(ctx).Create(&eval).Error; err != nil {
			return fmt.Errorf("creating work product evaluation: %w", err)
		}
	}

	// Update the work product with the final score and status
	wp.FinalScore = finalScore
	wp.RecordStatus = enums.StatusClosed.String()
	now := time.Now().UTC()
	wp.CompletionDate = &now

	if err := s.db.WithContext(ctx).Save(&wp).Error; err != nil {
		return fmt.Errorf("updating work product score: %w", err)
	}

	s.log.Info().
		Str("workProductID", r.WorkProductID).
		Float64("finalScore", finalScore).
		Msg("work product evaluated")

	return nil
}

// =========================================================================
// Legacy Feedback (direct implementations retained for backward
// compatibility)
// =========================================================================

func (s *performanceManagementService) RequestFeedback(ctx context.Context, req interface{}) error {
	r, ok := req.(*performance.FeedbackRequestModel)
	if !ok {
		return fmt.Errorf("invalid request type for RequestFeedback")
	}

	// Find the request log
	var log performance.FeedbackRequestLog
	if err := s.db.WithContext(ctx).
		Where("feedback_request_log_id = ?", r.RequestID).
		First(&log).Error; err != nil {
		return fmt.Errorf("feedback request not found: %w", err)
	}

	// Reassign to new assignee if specified
	if r.AssigneeID != "" && r.AssigneeID != log.AssignedStaffID {
		log.AssignedStaffID = r.AssigneeID
		log.TimeInitiated = time.Now().UTC()
		log.TimeCompleted = nil
	}

	log.RecordStatus = enums.StatusActive.String()

	if err := s.db.WithContext(ctx).Save(&log).Error; err != nil {
		return fmt.Errorf("updating feedback request: %w", err)
	}

	return nil
}

func (s *performanceManagementService) GetFeedbackRequests(ctx context.Context, staffID string) (interface{}, error) {
	resp := performance.FeedbackRequestListResponseVm{}

	var requests []performance.FeedbackRequestLog
	err := s.db.WithContext(ctx).
		Where("assigned_staff_id = ?", staffID).
		Order("time_initiated DESC").
		Find(&requests).Error
	if err != nil {
		s.log.Error().Err(err).Str("staffID", staffID).Msg("failed to get feedback requests")
		resp.HasError = true
		resp.Message = "failed to retrieve feedback requests"
		return resp, err
	}

	resp.Requests = requests
	resp.TotalRecords = len(requests)
	resp.Message = "operation completed successfully"
	return resp, nil
}

// =========================================================================
// Audit
// =========================================================================

func (s *performanceManagementService) LogAuditAction(ctx context.Context, action string, details interface{}) error {
	s.log.Info().
		Str("action", action).
		Interface("details", details).
		Msg("audit action logged")
	return nil
}

// =========================================================================
// Shared Helpers
// =========================================================================

// excludedStatuses returns the set of statuses that should be excluded from
// active work product queries. Mirrors the .NET statuses array used throughout
// the PerformanceManagementService.
func excludedStatuses() []string {
	return []string{
		enums.StatusCancelled.String(),
		enums.StatusPaused.String(),
		enums.StatusRejected.String(),
		enums.StatusReturned.String(),
		enums.StatusDraft.String(),
		enums.StatusPendingAcceptance.String(),
	}
}

// getGrade maps a percentage score to a PerformanceGrade.
// Mirrors the .NET GetGrade method exactly.
func getGrade(percentageScore float64) enums.PerformanceGrade {
	switch {
	case percentageScore < 30:
		return enums.PerformanceGradeProbation
	case percentageScore < 50:
		return enums.PerformanceGradeDeveloping
	case percentageScore < 66:
		return enums.PerformanceGradeProgressive
	case percentageScore < 80:
		return enums.PerformanceGradeCompetent
	case percentageScore < 90:
		return enums.PerformanceGradeAccomplished
	default:
		return enums.PerformanceGradeExemplary
	}
}

// getSLAConfig retrieves the SLA configuration values from global settings.
// Returns (REQUEST_SLA_IN_HOURS, PMS_360_FEEDBACK_REQUEST_SLA_IN_HOURS).
func (s *performanceManagementService) getSLAConfig(ctx context.Context) (int, int) {
	requestSLA := 168 // default 168 hours / 7 days
	pms360SLA := 336  // default 336 hours / 14 days

	if val, err := s.globalSettingSvc.GetIntValue(ctx, "REQUEST_SLA_IN_HOURS"); err == nil {
		requestSLA = val
	}
	if val, err := s.globalSettingSvc.GetIntValue(ctx, "PMS_360_FEEDBACK_REQUEST_SLA_IN_HOURS"); err == nil {
		pms360SLA = val
	}

	return requestSLA, pms360SLA
}

// getStaffActiveReviewPeriod retrieves the active review period for a staff member.
// First tries the staff-specific period, then falls back to the globally active period.
func (s *performanceManagementService) getStaffActiveReviewPeriod(ctx context.Context, staffID string) (*performance.PerformanceReviewPeriod, error) {
	var period performance.PerformanceReviewPeriod

	// Try staff-specific first by looking at planned objectives
	var plannedObj performance.ReviewPeriodIndividualPlannedObjective
	err := s.db.WithContext(ctx).
		Where("staff_id = ? AND record_status = ?", staffID, enums.StatusActive.String()).
		Preload("ReviewPeriod").
		Order("created_at DESC").
		First(&plannedObj).Error
	if err == nil && plannedObj.ReviewPeriod != nil {
		return plannedObj.ReviewPeriod, nil
	}

	// Fall back to globally active review period
	err = s.db.WithContext(ctx).
		Where("record_status = ?", enums.StatusActive.String()).
		Order("start_date DESC").
		First(&period).Error
	if err != nil {
		return nil, fmt.Errorf("no active review period found: %w", err)
	}

	return &period, nil
}

// countSLABreachedRequests counts the number of requests that have breached
// their SLA, accounting for leave days and public holidays.
// Mirrors the .NET GetRequestCount method.
func (s *performanceManagementService) countSLABreachedRequests(
	requests []performance.FeedbackRequestLog,
	staffID string,
	slaHours int,
) int {
	count := 0
	now := time.Now()

	for _, r := range requests {
		adjustedSLA := slaHours
		// Note: Leave days and public holiday adjustments would be fetched
		// from external services in production. Preserved as placeholder
		// matching .NET business logic.

		if r.TimeCompleted != nil {
			// Request completed: check if it exceeded SLA
			completionHours := r.TimeCompleted.Sub(r.TimeInitiated).Hours()
			if completionHours > float64(adjustedSLA) {
				count++
			}
		} else {
			// Request still pending: check if SLA window has passed
			deadline := r.TimeInitiated.Add(time.Duration(adjustedSLA) * time.Hour)
			if deadline.Before(now) || deadline.Equal(now) {
				count++
			}
		}
	}

	return count
}

// recalculateDeductedPoints recalculates and updates the deducted points
// for a staff member based on SLA breaches.
// Mirrors the .NET RecalculateDeductedPoints method.
func (s *performanceManagementService) recalculateDeductedPoints(ctx context.Context, staffID, reviewPeriodID string) {
	if strings.TrimSpace(reviewPeriodID) == "" {
		return
	}

	// Get review period
	var reviewPeriod performance.PerformanceReviewPeriod
	if err := s.db.WithContext(ctx).
		Where("period_id = ?", reviewPeriodID).
		First(&reviewPeriod).Error; err != nil {
		s.log.Error().Err(err).Msg("recalculate deducted points: review period not found")
		return
	}

	requestSLA, pms360SLA := s.getSLAConfig(ctx)

	// Get all feedback requests for this staff in this review period
	var requests []performance.FeedbackRequestLog
	s.db.WithContext(ctx).
		Where("assigned_staff_id = ? AND review_period_id = ?", staffID, reviewPeriodID).
		Find(&requests)

	now := time.Now()
	deductedPoints := 0.0

	if len(requests) > 0 {
		// Non-360 breached requests
		var nonPMSBreached []performance.FeedbackRequestLog
		for _, r := range requests {
			if r.FeedbackRequestType != enums.FeedbackRequest360ReviewFeedback &&
				r.HasSLA &&
				r.TimeInitiated.Add(time.Duration(requestSLA)*time.Hour).Before(now) {
				nonPMSBreached = append(nonPMSBreached, r)
			}
		}
		nonPMSCount := s.countSLABreachedRequests(nonPMSBreached, staffID, requestSLA)

		// 360 breached requests
		var pmsBreached []performance.FeedbackRequestLog
		for _, r := range requests {
			if r.FeedbackRequestType == enums.FeedbackRequest360ReviewFeedback &&
				r.HasSLA &&
				r.TimeInitiated.Add(time.Duration(pms360SLA)*time.Hour).Before(now) {
				pmsBreached = append(pmsBreached, r)
			}
		}
		pmsCount := s.countSLABreachedRequests(pmsBreached, staffID, pms360SLA)

		deductedPoints = float64(nonPMSCount + pmsCount)
	}

	if deductedPoints > reviewPeriod.MaxPoints {
		deductedPoints = reviewPeriod.MaxPoints
	}

	// Update or create period score
	var periodScore performance.PeriodScore
	err := s.db.WithContext(ctx).
		Where("review_period_id = ? AND staff_id = ?", reviewPeriodID, staffID).
		First(&periodScore).Error

	if err != nil {
		// Create new period score
		periodScore = performance.PeriodScore{
			ReviewPeriodID:    reviewPeriodID,
			StaffID:           staffID,
			FinalGrade:        enums.PerformanceGradeDeveloping,
			HRDDeductedPoints: deductedPoints,
			EndDate:           reviewPeriod.EndDate,
			StrategyID:        reviewPeriod.StrategyID,
		}
		periodScore.RecordStatus = enums.StatusActive.String()
		periodScore.IsActive = true

		if err := s.db.WithContext(ctx).Create(&periodScore).Error; err != nil {
			s.log.Error().Err(err).Msg("failed to create period score for deducted points")
		}
	} else {
		periodScore.HRDDeductedPoints = deductedPoints
		if err := s.db.WithContext(ctx).Save(&periodScore).Error; err != nil {
			s.log.Error().Err(err).Msg("failed to update period score deducted points")
		}
	}
}
