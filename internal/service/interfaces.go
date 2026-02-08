package service

import (
	"context"
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/enterprise-pms/pms-api/internal/domain/performance"
)

// --- Performance Management ---

// PerformanceManagementService handles all performance management operations.
// Methods are organised by the sub-service that implements them. The main
// performanceManagementService struct delegates each call to the appropriate
// composed sub-service.
type PerformanceManagementService interface {
	// =====================================================================
	// Strategies (delegated to strategyService -- separate agent)
	// =====================================================================
	GetStrategies(ctx context.Context) (interface{}, error)
	GetStrategicThemes(ctx context.Context) (interface{}, error)
	GetStrategicThemesById(ctx context.Context, strategyID string) (interface{}, error)
	CreateStrategy(ctx context.Context, req interface{}) (interface{}, error)
	UpdateStrategy(ctx context.Context, req interface{}) (interface{}, error)
	CreateStrategicTheme(ctx context.Context, req interface{}) (interface{}, error)
	UpdateStrategicTheme(ctx context.Context, req interface{}) (interface{}, error)

	// =====================================================================
	// Objectives (via objectiveService)
	// =====================================================================

	// Enterprise Objectives
	GetEnterpriseObjectives(ctx context.Context) (interface{}, error)
	CreateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error)
	UpdateEnterpriseObjective(ctx context.Context, req interface{}) (interface{}, error)

	// Department Objectives
	GetDepartmentObjectives(ctx context.Context) (interface{}, error)
	CreateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error)
	UpdateDepartmentObjective(ctx context.Context, req interface{}) (interface{}, error)

	// Division Objectives
	GetDivisionObjectives(ctx context.Context) (interface{}, error)
	GetDivisionObjectivesByDivisionId(ctx context.Context, divisionID int) (interface{}, error)
	CreateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error)
	UpdateDivisionObjective(ctx context.Context, req interface{}) (interface{}, error)

	// Office Objectives
	GetOfficeObjectives(ctx context.Context) (interface{}, error)
	GetOfficeObjectivesByOfficeId(ctx context.Context, officeID int) (interface{}, error)
	CreateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error)
	UpdateOfficeObjective(ctx context.Context, req interface{}) (interface{}, error)

	// Objective Categories & Definitions
	GetObjectiveCategories(ctx context.Context) (interface{}, error)
	CreateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error)
	UpdateObjectiveCategory(ctx context.Context, req interface{}) (interface{}, error)
	GetCategoryDefinitions(ctx context.Context, categoryID string) (interface{}, error)
	CreateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error)
	UpdateCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error)

	// Consolidated Objectives
	GetConsolidatedObjectives(ctx context.Context) (interface{}, error)
	GetConsolidatedObjectivesPaginated(ctx context.Context, params interface{}) (interface{}, error)
	ProcessObjectivesUpload(ctx context.Context, req interface{}) (interface{}, error)
	DeActivateOrReactivateObjectives(ctx context.Context, req interface{}, deactivate bool) (interface{}, error)

	// Evaluation Options & Questionnaires (via objectiveService)
	GetEvaluationOptions(ctx context.Context) (interface{}, error)
	SaveEvaluationOptions(ctx context.Context, req interface{}) (interface{}, error)
	GetFeedbackQuestionnaires(ctx context.Context) (interface{}, error)
	SaveFeedbackQuestionnaires(ctx context.Context, req interface{}) (interface{}, error)
	SaveFeedbackQuestionnaireOptions(ctx context.Context, req interface{}) (interface{}, error)

	// PMS Competencies (via objectiveService)
	GetPmsCompetencies(ctx context.Context) (interface{}, error)
	CreatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error)
	UpdatePmsCompetency(ctx context.Context, req interface{}) (interface{}, error)

	// Work Product Definitions (via objectiveService)
	GetObjectiveWorkProductDefinitions(ctx context.Context, objectiveID string, objectiveLevel int) (interface{}, error)
	GetAllWorkProductDefinitions(ctx context.Context) (interface{}, error)
	GetAllPaginatedWorkProductDefinitions(ctx context.Context, pageIndex, pageSize int, search string) (interface{}, error)
	SaveWorkProductDefinitions(ctx context.Context, req interface{}) (interface{}, error)

	// Approval (via objectiveService)
	ApproveRecords(ctx context.Context, req interface{}) (interface{}, error)
	RejectRecords(ctx context.Context, req interface{}) (interface{}, error)

	// =====================================================================
	// Projects (via projectService)
	// =====================================================================
	SetupProject(ctx context.Context, req interface{}) error
	GetProjectsByStaff(ctx context.Context, staffID string) (interface{}, error)
	AddProjectObjective(ctx context.Context, req interface{}) error
	AddProjectMember(ctx context.Context, req interface{}) error

	// Full project lifecycle
	ProjectSetup(ctx context.Context, req *performance.ProjectRequestModel) (performance.ResponseVm, error)
	GetProject(ctx context.Context, projectID string) (performance.ProjectResponseVm, error)
	GetProjects(ctx context.Context) (performance.ProjectListResponseVm, error)
	GetProjectsByManager(ctx context.Context, managerID string) (performance.ProjectListResponseVm, error)
	ProjectObjectiveSetup(ctx context.Context, req *performance.ProjectObjectiveRequestModel) (performance.ResponseVm, error)
	GetProjectObjectives(ctx context.Context, projectID string) (performance.ProjectObjectiveListResponseVm, error)
	ProjectMembersSetup(ctx context.Context, req *performance.ProjectMemberRequestModel) (performance.ResponseVm, error)
	GetProjectMembers(ctx context.Context, projectID string) (performance.ProjectMemberListResponseVm, error)
	GetProjectsAssigned(ctx context.Context, staffID string) (performance.ProjectAssignedListResponseVm, error)
	GetStaffProjects(ctx context.Context, staffID string) (performance.ProjectAssignedListResponseVm, error)
	GetProjectWorkProductStaffList(ctx context.Context, projectID string) ([]string, error)
	ChangeProjectLead(ctx context.Context, req *performance.ChangeAdhocLeadRequestModel) error
	ValidateStaffEligibilityForAdhoc(ctx context.Context, staffID, reviewPeriodID string) (performance.AdhocStaffResponseVm, error)

	// =====================================================================
	// Committees (via committeeService)
	// =====================================================================
	CommitteeSetup(ctx context.Context, req *performance.CommitteeRequestModel) (performance.ResponseVm, error)
	GetCommittee(ctx context.Context, committeeID string) (performance.CommitteeResponseVm, error)
	GetCommittees(ctx context.Context) (performance.CommitteeListResponseVm, error)
	GetCommitteesByChairperson(ctx context.Context, chairpersonID string) (performance.CommitteeListResponseVm, error)
	CommitteeObjectiveSetup(ctx context.Context, req *performance.CommitteeObjectiveRequestModel) (performance.ResponseVm, error)
	GetCommitteeObjectives(ctx context.Context, committeeID string) (performance.CommitteeObjectiveListResponseVm, error)
	CommitteeMembersSetup(ctx context.Context, req *performance.CommitteeMemberRequestModel) (performance.ResponseVm, error)
	GetCommitteeMembers(ctx context.Context, committeeID string) (performance.CommitteeMemberListResponseVm, error)
	GetCommitteesAssigned(ctx context.Context, staffID string) (performance.CommitteeAssignedListResponseVm, error)
	GetStaffCommittees(ctx context.Context, staffID string) (performance.CommitteeAssignedListResponseVm, error)
	GetCommitteeWorkProductStaffList(ctx context.Context, committeeID string) ([]string, error)
	ChangeCommitteeChairperson(ctx context.Context, req *performance.ChangeAdhocLeadRequestModel) error

	// =====================================================================
	// Work Products (via workProductService)
	// =====================================================================
	AddWorkProduct(ctx context.Context, req interface{}) error
	EvaluateWorkProduct(ctx context.Context, req interface{}) error

	// Full work product lifecycle
	WorkProductSetup(ctx context.Context, req *performance.WorkProductRequestModel) (performance.ResponseVm, error)
	ProjectAssignedWorkProductSetup(ctx context.Context, req *performance.ProjectAssignedWorkProductRequestModel) (performance.ResponseVm, error)
	CommitteeAssignedWorkProductSetup(ctx context.Context, req *performance.CommitteeAssignedWorkProductRequestModel) (performance.ResponseVm, error)
	GetWorkProduct(ctx context.Context, workProductID string) (performance.WorkProductResponseVm, error)
	GetProjectWorkProducts(ctx context.Context, projectID string) (performance.ProjectWorkProductListResponseVm, error)
	GetProjectAssignedWorkProducts(ctx context.Context, projectID string) (performance.ProjectAssignedWorkProductListResponseVm, error)
	GetCommitteeWorkProducts(ctx context.Context, committeeID string) (performance.CommitteeWorkProductListResponseVm, error)
	GetCommitteeAssignedWorkProducts(ctx context.Context, committeeID string) (performance.CommitteeAssignedWorkProductListResponseVm, error)
	GetOperationalWorkProducts(ctx context.Context, plannedObjectiveID string) (performance.OperationalObjectiveWorkProductListResponseVm, error)
	GetStaffWorkProducts(ctx context.Context, staffID, reviewPeriodID string) (performance.StaffWorkProductListResponseVm, error)
	GetAllStaffWorkProducts(ctx context.Context, staffID string) (performance.StaffWorkProductListResponseVm, error)
	GetObjectiveWorkProducts(ctx context.Context, objectiveID string) (performance.ObjectiveWorkProductListResponseVm, error)
	WorkProductTaskSetup(ctx context.Context, req *performance.WorkProductTaskRequestModel) (performance.ResponseVm, error)
	GetWorkProductTasks(ctx context.Context, workProductID string) (performance.WorkProductTaskListResponseVm, error)
	ReCalculateWorkProductPoints(ctx context.Context, staffID, reviewPeriodID string) (performance.RecalculateWorkProductResponseVm, error)
	WorkProductEvaluation(ctx context.Context, req *performance.WorkProductEvaluationRequestModel) (performance.EvaluationResponseVm, error)
	GetWorkProductEvaluation(ctx context.Context, workProductID string) (performance.WorkProductEvaluationResponseVm, error)
	InitiateWorkProductReEvaluation(ctx context.Context, workProductID string) (performance.ResponseVm, error)

	// =====================================================================
	// Feedback Requests (via feedbackRequestService)
	// =====================================================================
	RequestFeedback(ctx context.Context, req interface{}) error
	GetFeedbackRequests(ctx context.Context, staffID string) (interface{}, error)

	// Full feedback request lifecycle
	LogRequest(ctx context.Context, feedbackType enums.FeedbackRequestType, referenceID, assignedStaffID, requestOwnerStaffID, reviewPeriodID string, hasSLA bool) error
	LogAcceptanceRequest(ctx context.Context, feedbackType enums.FeedbackRequestType, referenceID, assignedStaffID, requestOwnerStaffID, reviewPeriodID string) error
	GetRequests(ctx context.Context, staffID string, feedbackType *enums.FeedbackRequestType, status *string) (performance.FeedbackRequestListResponseVm, error)
	GetRequestsByOwner(ctx context.Context, requestOwnerStaffID string, feedbackType *enums.FeedbackRequestType) (performance.FeedbackRequestListResponseVm, error)
	GetBreachedRequests(ctx context.Context, staffID, reviewPeriodID string) (performance.BreachedFeedbackRequestListResponseVm, error)
	GetPendingRequests(ctx context.Context, staffID string) (performance.GetStaffPendingRequestVm, error)
	GetFeedbackRequest(ctx context.Context, requestID string) (performance.FeedbackRequestLogResponseVm, error)
	GetRequestDetails(ctx context.Context, requestID string) (performance.FeedbackRequestLogResponseVm, error)
	UpdateRequest(ctx context.Context, requestID string, comment, attachment string) error
	ReassignRequest(ctx context.Context, requestID, newAssignedStaffID string) error
	ReassignSelfRequest(ctx context.Context, requestID, currentStaffID, newAssignedStaffID string) error
	CloseRequest(ctx context.Context, requestID string) error
	CloseReviewPeriodRequests(ctx context.Context, reviewPeriodID string) error
	ReactivateReviewPeriodRequest(ctx context.Context, reviewPeriodID string) error
	ReInitiateSameRequest(ctx context.Context, requestID string) error
	TreatAssignedRequest(ctx context.Context, req *performance.TreatFeedbackRequestModel) error
	HasLineManager(ctx context.Context, staffID string) (bool, error)
	HasVacationRule(ctx context.Context, staffID string) (bool, error)
	GetStaffLeaveDays(ctx context.Context, staffID string, startDate, endDate time.Time) (performance.LeaveResponseVm, error)
	GetPublicDays(ctx context.Context, startDate, endDate time.Time) (performance.PublicHolidaysResponseVm, error)
	AutoReassignAndLogRequest(ctx context.Context, requestID string) error

	// =====================================================================
	// Competency / 360 Review (via competencyReviewService)
	// =====================================================================
	GetCompetencyReview(ctx context.Context, staffID string) (interface{}, error)
	Submit360Feedback(ctx context.Context, req interface{}) error

	// Full competency review lifecycle
	CompetencyReviewFeedbackSetup(ctx context.Context, req *performance.CompetencyReviewFeedbackRequestModel) (performance.ResponseVm, error)
	CompetencyReviewerSetup(ctx context.Context, req *performance.CompetencyReviewerRequestModel) (performance.ResponseVm, error)
	CompetencyRatingSetup(ctx context.Context, req *performance.SavePmsCompetencyRequestVm) (performance.ResponseVm, error)
	GetCompetencyReviewFeedback(ctx context.Context, feedbackID string) (performance.CompetencyReviewFeedbackResponseVm, error)
	GetCompetencyReviewFeedbackDetails(ctx context.Context, feedbackID string) (performance.CompetencyReviewFeedbackDetailsResponseVm, error)
	GetAllCompetencyReviewFeedbacks(ctx context.Context, staffID string) (performance.CompetencyReviewFeedbackListResponseVm, error)
	GetCompetencyReviews(ctx context.Context, reviewerStaffID string) (performance.CompetencyReviewersListResponseVm, error)
	GetReviewerFeedbackDetails(ctx context.Context, reviewerID string) (performance.CompetencyReviewersResponseVm, error)
	GetQuestionnaire(ctx context.Context, staffID string) (performance.QuestionnaireListResponseVm, error)
	CompetencyGapClosureSetup(ctx context.Context, req *performance.CompetencyGapClosureRequestModel) (performance.ResponseVm, error)
	Initiate360Review(ctx context.Context, req *performance.Initiate360ReviewRequestModel) (performance.ResponseVm, error)
	Complete360Review(ctx context.Context, req *performance.Complete360ReviewRequestModel) (performance.ResponseVm, error)

	// =====================================================================
	// Period Objective Evaluations (via evaluationService)
	// =====================================================================
	ReviewPeriodObjectiveEvaluation(ctx context.Context, req *performance.PeriodObjectiveEvaluationRequestModel) (performance.PeriodObjectiveEvaluationResponseVm, error)
	GetReviewPeriodObjectiveEvaluation(ctx context.Context, reviewPeriodID, enterpriseObjectiveID string) (performance.PeriodObjectiveEvaluationResponseVm, error)
	GetReviewPeriodObjectiveEvaluations(ctx context.Context, reviewPeriodID string) (performance.PeriodObjectiveEvaluationListResponseVm, error)
	ReviewPeriodDepartmentObjectiveEvaluation(ctx context.Context, req *performance.PeriodObjectiveDepartmentEvaluationRequestModel) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error)
	GetReviewPeriodDepartmentObjectiveEvaluation(ctx context.Context, reviewPeriodID, enterpriseObjectiveID string, departmentID int) (performance.PeriodObjectiveDepartmentEvaluationResponseVm, error)
	GetReviewPeriodDepartmentObjectiveEvaluations(ctx context.Context, reviewPeriodID string) (performance.PeriodObjectiveDepartmentEvaluationListResponseVm, error)
	GetWorkProductsPerEnterpriseObjective(ctx context.Context, enterpriseObjectiveID, reviewPeriodID string) (performance.ObjectiveWorkProductListResponseVm, error)
	GetDepartmentsbyObjective(ctx context.Context, enterpriseObjectiveID string) (performance.DepartmentListResponseVm, error)

	// =====================================================================
	// Dashboard / Scoring / Reporting (via dashboardService & periodScoreService)
	// =====================================================================
	GetPerformanceScore(ctx context.Context, staffID string) (interface{}, error)
	GetDashboardStats(ctx context.Context, staffID string) (interface{}, error)
	GetStaffAnnualPerformanceScoreCardStatistics(ctx context.Context, staffID string, year int) (performance.StaffAnnualScoreCardResponseVm, error)
	GetSubordinatesStaffPerformanceScoreCardStatistics(ctx context.Context, managerStaffID, reviewPeriodID string) (performance.AllStaffScoreCardResponseVm, error)
	GetOrganogramPerformanceSummaryStatistics(ctx context.Context, referenceID, reviewPeriodID string, organogramLevel enums.OrganogramLevel) (performance.OrganogramPerformanceSummaryResponseVm, error)
	GetOrganogramPerformanceSummaryListStatistics(ctx context.Context, headOfUnitID, reviewPeriodID string, organogramLevel enums.OrganogramLevel) (performance.OrganogramPerformanceSummaryListResponseVm, error)

	// =====================================================================
	// Audit
	// =====================================================================
	LogAuditAction(ctx context.Context, action string, details interface{}) error
}

// --- PMS Setup ---

// PmsSetupService handles PMS system configuration.
type PmsSetupService interface {
	// Settings
	AddSetting(ctx context.Context, req interface{}) (interface{}, error)
	UpdateSetting(ctx context.Context, req interface{}) (interface{}, error)
	GetSettingDetails(ctx context.Context, settingID string) (interface{}, error)
	ListAllSettings(ctx context.Context) (interface{}, error)

	// PMS Configurations
	AddPmsConfiguration(ctx context.Context, req interface{}) (interface{}, error)
	UpdatePmsConfiguration(ctx context.Context, req interface{}) (interface{}, error)
	GetPmsConfigurationDetails(ctx context.Context, configID string) (interface{}, error)
	ListAllPmsConfigurations(ctx context.Context) (interface{}, error)
}

// --- Review Period ---

// ReviewPeriodService manages review periods and objectives planning.
type ReviewPeriodService interface {
	// Lifecycle
	SaveDraftReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	AddReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	SubmitDraftReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	ApproveReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	RejectReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	ReturnReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	ReSubmitReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	UpdateReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	CancelReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	CloseReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)

	// Toggles
	EnableObjectivePlanning(ctx context.Context, req interface{}) (interface{}, error)
	DisableObjectivePlanning(ctx context.Context, req interface{}) (interface{}, error)
	EnableWorkProductPlanning(ctx context.Context, req interface{}) (interface{}, error)
	DisableWorkProductPlanning(ctx context.Context, req interface{}) (interface{}, error)
	EnableWorkProductEvaluation(ctx context.Context, req interface{}) (interface{}, error)
	DisableWorkProductEvaluation(ctx context.Context, req interface{}) (interface{}, error)

	// Retrieval
	GetActiveReviewPeriod(ctx context.Context) (interface{}, error)
	GetStaffActiveReviewPeriod(ctx context.Context, staffID string) (interface{}, error)
	GetReviewPeriodDetails(ctx context.Context, reviewPeriodID string) (interface{}, error)

	// Period Objectives
	SaveDraftReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error)
	AddReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error)
	SubmitDraftReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error)
	CancelReviewPeriodObjective(ctx context.Context, req interface{}) (interface{}, error)
	GetReviewPeriodObjectives(ctx context.Context, reviewPeriodID string) (interface{}, error)

	// Category Definitions
	SaveDraftCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error)
	AddCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error)
	SubmitDraftCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error)
	ApproveCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error)
	RejectCategoryDefinition(ctx context.Context, req interface{}) (interface{}, error)

	// Extensions
	AddReviewPeriodExtension(ctx context.Context, req interface{}) (interface{}, error)
	GetReviewPeriodExtensions(ctx context.Context, reviewPeriodID string) (interface{}, error)

	// 360 Reviews
	AddReviewPeriod360Review(ctx context.Context, req interface{}) (interface{}, error)
	GetReviewPeriod360Reviews(ctx context.Context, reviewPeriodID string) (interface{}, error)

	// Individual Planned Objectives
	SaveDraftIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error)
	AddIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error)
	SubmitDraftIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error)
	ApproveIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error)
	RejectIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error)
	ReturnIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error)
	CancelIndividualPlannedObjective(ctx context.Context, req interface{}) (interface{}, error)
	GetStaffIndividualPlannedObjectives(ctx context.Context, staffID, reviewPeriodID string) (interface{}, error)

	// Period Objective Evaluations
	CreatePeriodObjectiveEvaluation(ctx context.Context, req interface{}) (interface{}, error)
	CreatePeriodObjectiveDepartmentEvaluation(ctx context.Context, req interface{}) (interface{}, error)
	GetPeriodObjectiveEvaluations(ctx context.Context, reviewPeriodID string) (interface{}, error)
	GetPeriodObjectiveDepartmentEvaluations(ctx context.Context, reviewPeriodID string) (interface{}, error)

	// Period Scores
	GetStaffPeriodScore(ctx context.Context, staffID, reviewPeriodID string) (interface{}, error)

	// Additional Retrieval (mirrors remaining .NET IReviewPeriodService methods)
	GetReviewPeriods(ctx context.Context) (interface{}, error)
	GetReviewPeriodCategoryDefinitions(ctx context.Context, reviewPeriodID string) (interface{}, error)
	GetPlannedObjective(ctx context.Context, plannedObjectiveID string) (interface{}, error)
	GetEnterpriseObjectiveByLevel(ctx context.Context, objectiveID string, objectiveLevel int) (interface{}, error)
	ArchiveCancelledObjectives(ctx context.Context, staffID string, reviewPeriodID string) (interface{}, error)
	ArchiveCancelledWorkProducts(ctx context.Context, staffID string, reviewPeriodID string) (interface{}, error)
}

// --- Grievance Management ---

// GrievanceManagementService handles staff grievances.
type GrievanceManagementService interface {
	RaiseNewGrievance(ctx context.Context, req interface{}) (interface{}, error)
	UpdateGrievance(ctx context.Context, req interface{}) (interface{}, error)
	CreateGrievanceResolution(ctx context.Context, req interface{}) (interface{}, error)
	UpdateGrievanceResolution(ctx context.Context, req interface{}) (interface{}, error)
	GetStaffGrievances(ctx context.Context, staffID string) (interface{}, error)
	GetGrievancesReport(ctx context.Context) (interface{}, error)

	// LogRequestAsync creates or updates a FeedbackRequestLog for a given
	// reference, handles vacation rule delegation, and sends email notification.
	// Mirrors the .NET GrievanceManagementService.LogRequestAsync method.
	LogRequestAsync(ctx context.Context, referenceID, assignerStaffID, assigneeStaffID string, feedbackRequestType enums.FeedbackRequestType, hasSLA bool) error

	// HasVacationRule checks whether a staff member has an active vacation
	// delegation rule. Returns the delegate staff ID if configured.
	// Mirrors the .NET GrievanceManagementService.HasVacationRule method.
	HasVacationRule(ctx context.Context, staffID string, startDate time.Time) (*VacationRuleResult, error)
}

// --- Staff Management ---

// StaffManagementService handles staff and role management operations.
// Mirrors the .NET StaffsController service dependencies.
type StaffManagementService interface {
	AddStaff(ctx context.Context, req interface{}) (interface{}, error)
	GetAllStaffs(ctx context.Context, searchString string) (interface{}, error)
	GetAllRoles(ctx context.Context) (interface{}, error)
	AddRole(ctx context.Context, req interface{}) (interface{}, error)
	DeleteRole(ctx context.Context, roleName string) (interface{}, error)
	AddStaffToRole(ctx context.Context, req interface{}) (interface{}, error)
	RemoveStaffFromRole(ctx context.Context, userId string, roleName string) (interface{}, error)
	GetStaffRoles(ctx context.Context, id string) (interface{}, error)
}

// --- Organogram ---

// OrganogramService handles organizational structure operations (directorates, departments, divisions, offices).
type OrganogramService interface {
	// Directorates
	GetDirectorates(ctx context.Context) (interface{}, error)
	SaveDirectorate(ctx context.Context, req interface{}) (interface{}, error)
	DeleteDirectorate(ctx context.Context, id int, isSoftDelete bool) (interface{}, error)

	// Departments
	GetDepartments(ctx context.Context, directorateId *int) (interface{}, error)
	SaveDepartment(ctx context.Context, req interface{}) (interface{}, error)
	DeleteDepartment(ctx context.Context, id int, isSoftDelete bool) (interface{}, error)

	// Divisions
	GetDivisions(ctx context.Context, departmentId *int) (interface{}, error)
	SaveDivision(ctx context.Context, req interface{}) (interface{}, error)
	DeleteDivision(ctx context.Context, id int, isSoftDelete bool) (interface{}, error)

	// Offices
	GetOffices(ctx context.Context, divisionId *int) (interface{}, error)
	GetOfficeByCode(ctx context.Context, officeCode string) (interface{}, error)
	SaveOffice(ctx context.Context, req interface{}) (interface{}, error)
	DeleteOffice(ctx context.Context, id int, isSoftDelete bool) (interface{}, error)
}

// --- ERP Employee ---

// ErpEmployeeService retrieves employee and organogram data from the ERP system.
// Mirrors the .NET EmployeeInformationController service dependencies.
type ErpEmployeeService interface {
	// Organogram lookups
	GetAllDepartments(ctx context.Context) (interface{}, error)
	GetAllDivisions(ctx context.Context, departmentId *int) (interface{}, error)
	GetAllOffices(ctx context.Context, divisionId *int) (interface{}, error)

	// Employee lookups
	GetEmployeeDetail(ctx context.Context, employeeNumber string) (interface{}, error)
	GetHeadSubordinates(ctx context.Context, employeeNumber string) (interface{}, error)
	GetEmployeeSubordinates(ctx context.Context, employeeNumber string) (interface{}, error)
	GetEmployeePeers(ctx context.Context, employeeNumber string) (interface{}, error)
	GetAllByDepartmentId(ctx context.Context, departmentId int) (interface{}, error)
	GetAllByDivisionId(ctx context.Context, divisionId int) (interface{}, error)
	GetAllByOfficeId(ctx context.Context, officeId int) (interface{}, error)
	GetAllEmployees(ctx context.Context) (interface{}, error)
	SeedOrganizationData(ctx context.Context) (interface{}, error)
	GetStaffIDMaskDetail(ctx context.Context, employeeNumber string) (interface{}, error)

	// Job role management
	UpdateStaffJobRole(ctx context.Context, req interface{}) (interface{}, error)
	GetStaffJobRoleById(ctx context.Context, employeeNumber string) (interface{}, error)
	GetJobRolesByOffice(ctx context.Context, req interface{}) (interface{}, error)
	GetStaffJobRoleRequests(ctx context.Context, employeeNumber string) (interface{}, error)
	ApproveRejectStaffJobRole(ctx context.Context, req interface{}) (interface{}, error)
}

// --- Role Management ---

// RoleManagementService handles role and permission management.
// Mirrors the .NET RoleMgtController / RolePermissionMgt mediator handlers.
type RoleManagementService interface {
	// GetPermissions retrieves all permissions, optionally filtered by roleId.
	GetPermissions(ctx context.Context, roleId string) (interface{}, error)
	// GetAllRolesWithPermission retrieves all application permissions alongside the
	// permissions assigned to the specified role.
	GetAllRolesWithPermission(ctx context.Context, roleId string) (interface{}, error)
	// AddPermissionToRole assigns a permission to a role.
	AddPermissionToRole(ctx context.Context, req interface{}) (interface{}, error)
	// RemovePermissionFromRole removes a permission from a role.
	RemovePermissionFromRole(ctx context.Context, roleId string, permissionId int) (interface{}, error)
}

// --- Global Settings ---

// GlobalSettingService provides typed access to global configuration values.
type GlobalSettingService interface {
	GetBoolValue(ctx context.Context, key string) (bool, error)
	GetStringValue(ctx context.Context, key string) (string, error)
	GetIntValue(ctx context.Context, key string) (int, error)
	GetFloatValue(ctx context.Context, key string) (float64, error)
}

// --- Authentication ---

// AuthService handles user authentication and token management.
type AuthService interface {
	AuthenticateAD(ctx context.Context, username, password string) (interface{}, error)
	GenerateTokenPair(ctx context.Context, userID string, roles []string) (accessToken string, refreshToken string, err error)
	ValidateToken(ctx context.Context, token string) (claims interface{}, err error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// --- Email ---

// EmailService handles sending emails.
// Mirrors .NET MailServices (IMailServices) and MailSender (IMailSender).
type EmailService interface {
	// SendEmail saves a single email record for async delivery.
	SendEmail(ctx context.Context, to string, subject string, body string) error
	// SendBulkEmail saves one email record per recipient.
	SendBulkEmail(ctx context.Context, to []string, subject string, body string) error
	// SendEmailWithCC saves an email with CC recipients.
	SendEmailWithCC(ctx context.Context, to string, cc []string, subject string, body string) error
	// ProcessEmail dispatches an email notification by title (mirrors .NET MailServices.ProcessEmail).
	ProcessEmail(ctx context.Context, req interface{}) (interface{}, error)
}

// --- File Storage ---

// FileStorageService handles file upload and storage.
type FileStorageService interface {
	SaveFile(ctx context.Context, fileName string, data []byte) (path string, err error)
	DeleteFile(ctx context.Context, path string) error
	GetFile(ctx context.Context, path string) ([]byte, error)
}

// --- Notification ---

// NotificationService handles email notifications using templates.
// Mirrors the .NET NotificationService + NotificationTemplates pattern.
type NotificationService interface {
	// Send sends a plain-text/HTML notification email.
	Send(ctx context.Context, userID string, message string) error
	// SendNewRequestNotification sends a notification for a new feedback request.
	// Uses SLA template if hasSLA is true, otherwise generic template.
	SendNewRequestNotification(ctx context.Context, recipientEmail, recipientName, requestName string, assignedDate time.Time, hasSLA bool, slaHours int) error
	// SendAssignerNewRequestNotification notifies the assigner that their request was initiated.
	SendAssignerNewRequestNotification(ctx context.Context, recipientEmail, recipientName, requestName string, assignedDate time.Time) error
	// SendRequestTreatedNotification notifies when a request is treated.
	SendRequestTreatedNotification(ctx context.Context, recipientEmail, recipientName, requestName string, treatedDate time.Time) error
	// SendAssignerRequestTreatedNotification notifies the assigner when their request is treated.
	SendAssignerRequestTreatedNotification(ctx context.Context, recipientEmail, recipientName, requestName string, treatedDate time.Time) error
}

// --- Encryption ---

// EncryptionService provides AES encryption utilities.
type EncryptionService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// --- Active Directory ---

// ActiveDirectoryService handles LDAP/AD user lookups.
type ActiveDirectoryService interface {
	Authenticate(username, password string) (bool, error)
	GetUser(username string) (interface{}, error)
}

// --- User Context ---

// UserContextService provides current user information from the request.
type UserContextService interface {
	GetUserID(ctx context.Context) string
	GetEmail(ctx context.Context) string
	GetRoles(ctx context.Context) []string
	IsAuthenticated(ctx context.Context) bool
	IsInRole(ctx context.Context, role string) bool
}

// --- Competency Management ---

// CompetencyService handles all competency management operations.
// Mirrors the .NET CompetencyMgtController service dependencies.
type CompetencyService interface {
	// Competencies
	GetCompetencies(ctx context.Context, req interface{}) (interface{}, error)
	SaveCompetency(ctx context.Context, req interface{}) (interface{}, error)
	ApproveCompetency(ctx context.Context, req interface{}) (interface{}, error)
	RejectCompetency(ctx context.Context, req interface{}) (interface{}, error)

	// Competency Categories
	GetCompetencyCategories(ctx context.Context) (interface{}, error)
	SaveCompetencyCategory(ctx context.Context, req interface{}) (interface{}, error)

	// Competency Category Gradings
	GetCompetencyCategoryGradings(ctx context.Context) (interface{}, error)
	SaveCompetencyCategoryGrading(ctx context.Context, req interface{}) (interface{}, error)

	// Competency Rating Definitions
	GetCompetencyRatingDefinitions(ctx context.Context, competencyId *int) (interface{}, error)
	SaveCompetencyRatingDefinition(ctx context.Context, req interface{}) (interface{}, error)

	// Competency Reviews
	GetCompetencyReviews(ctx context.Context) (interface{}, error)
	GetCompetencyReviewByReviewer(ctx context.Context, reviewerId string, reviewPeriodId *int) (interface{}, error)
	GetCompetencyReviewForEmployee(ctx context.Context, employeeNumber string, reviewPeriodId *int) (interface{}, error)
	GetCompetencyReviewDetail(ctx context.Context, req interface{}) (interface{}, error)
	SaveCompetencyReview(ctx context.Context, req interface{}) (interface{}, error)

	// Competency Review Profiles
	GetCompetencyReviewProfiles(ctx context.Context, employeeNumber string, reviewPeriodId *int) (interface{}, error)
	GetOfficeCompetencyReviews(ctx context.Context, officeId int, reviewPeriodId *int) (interface{}, error)
	GetGroupCompetencyReviewProfiles(ctx context.Context, reviewPeriodId, officeId, divisionId, departmentId *int) (interface{}, error)
	GetCompetencyMatrixReviewProfiles(ctx context.Context, reviewPeriodId, officeId, divisionId, departmentId *int) (interface{}, error)
	GetTechnicalCompetencyMatrixReviewProfiles(ctx context.Context, reviewPeriodId *int, jobRoleId int) (interface{}, error)
	SaveCompetencyReviewProfile(ctx context.Context, req interface{}) (interface{}, error)

	// Competency Gaps
	GetCompetencyGaps(ctx context.Context, employeeNumber string) (interface{}, error)
	CloseCompetencyGap(ctx context.Context, req interface{}) (interface{}, error)

	// Development Plans
	GetDevelopmentPlans(ctx context.Context, competencyProfileReviewId *int) (interface{}, error)
	SaveDevelopmentPlan(ctx context.Context, req interface{}) (interface{}, error)

	// Job Roles
	GetJobRoles(ctx context.Context) (interface{}, error)
	GetJobRoleByName(ctx context.Context, jobRoleName string) (interface{}, error)
	SaveJobRole(ctx context.Context, req interface{}) (interface{}, error)

	// Office Job Roles
	GetOfficeJobRoles(ctx context.Context, req interface{}) (interface{}, error)
	SaveOfficeJobRole(ctx context.Context, req interface{}) (interface{}, error)

	// Job Role Competencies
	GetJobRoleCompetencies(ctx context.Context, req interface{}) (interface{}, error)
	SaveJobRoleCompetency(ctx context.Context, req interface{}) (interface{}, error)

	// Behavioral Competencies
	GetBehavioralCompetencies(ctx context.Context) (interface{}, error)
	GetBehavioralCompetenciesByGradeName(ctx context.Context, gradeName string) (interface{}, error)
	SaveBehavioralCompetency(ctx context.Context, req interface{}) (interface{}, error)

	// Job Role Grades
	GetJobRoleGrades(ctx context.Context) (interface{}, error)
	SaveJobRoleGrade(ctx context.Context, req interface{}) (interface{}, error)

	// Job Grades
	GetJobGrades(ctx context.Context) (interface{}, error)
	SaveJobGrade(ctx context.Context, req interface{}) (interface{}, error)

	// Job Grade Groups
	GetJobGradeGroups(ctx context.Context) (interface{}, error)
	SaveJobGradeGroup(ctx context.Context, req interface{}) (interface{}, error)

	// Assign Job Grade Groups
	GetAssignJobGradeGroups(ctx context.Context) (interface{}, error)
	GetAssignJobGradeGroupByGradeName(ctx context.Context, gradeName string) (interface{}, error)
	SaveAssignJobGradeGroup(ctx context.Context, req interface{}) (interface{}, error)

	// Ratings
	GetRatings(ctx context.Context) (interface{}, error)
	SaveRating(ctx context.Context, req interface{}) (interface{}, error)

	// Review Periods
	GetReviewPeriods(ctx context.Context) (interface{}, error)
	GetCurrentReviewPeriod(ctx context.Context) (interface{}, error)
	SaveReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)
	ApproveReviewPeriod(ctx context.Context, req interface{}) (interface{}, error)

	// Review Types
	GetReviewTypes(ctx context.Context) (interface{}, error)
	SaveReviewType(ctx context.Context, req interface{}) (interface{}, error)

	// Bank Years
	GetBankYears(ctx context.Context) (interface{}, error)
	SaveBankYear(ctx context.Context, req interface{}) (interface{}, error)

	// Training Types
	GetTrainingTypes(ctx context.Context, isActive *bool) (interface{}, error)
	SaveTrainingType(ctx context.Context, req interface{}) (interface{}, error)

	// Population / Calculation
	PopulateAllReviews(ctx context.Context) (interface{}, error)
	PopulateOfficeReviews(ctx context.Context, officeId int) (interface{}, error)
	PopulateDivisionReviews(ctx context.Context, divisionId int) (interface{}, error)
	PopulateDepartmentReviews(ctx context.Context, departmentId int) (interface{}, error)
	PopulateReviewsByEmployeeId(ctx context.Context, employeeNumber string) (interface{}, error)
	CalculateReviews(ctx context.Context, req interface{}) (interface{}, error)
	RecalculateReviewsProfiles(ctx context.Context, req interface{}) (interface{}, error)

	// Email / Sync
	EmailService(ctx context.Context, req interface{}) (interface{}, error)
	SyncJobRoleUpdateSOA(ctx context.Context, req interface{}) (interface{}, error)
}

// --- Password Generator ---

// PasswordGenerator generates secure passwords.
type PasswordGenerator interface {
	Generate(length int) string
}
