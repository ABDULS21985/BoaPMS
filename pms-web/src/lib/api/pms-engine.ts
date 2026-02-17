import { get, post, put } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type {
  Project,
  Committee,
  WorkProduct,
  WorkProductTask,
  WorkProductEvaluation,
  IndividualPlannedObjective,
  ProjectMember,
  CommitteeMember,
  FeedbackRequestLog,
  BreachedFeedbackRequestLog,
  StaffPendingRequest,
  CompetencyReviewFeedback,
  CompetencyReviewFeedbackDetails,
  CompetencyReviewer,
  PmsCompetency,
} from "@/types/performance";

// --- Projects ---
export const getProjects = (params?: { staffId?: string; reviewPeriodId?: string }) => {
  const q = new URLSearchParams();
  if (params?.staffId) q.append("staffId", params.staffId);
  if (params?.reviewPeriodId) q.append("reviewPeriodId", params.reviewPeriodId);
  return get<BaseAPIResponse<Project[]>>(`/pms-engine/projects?${q.toString()}`);
};
export const getProjectDetails = (projectId: string) =>
  get<BaseAPIResponse<Project>>(`/pms-engine/projects/${projectId}`);
export const getStaffProjects = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<Project[]>>(`/pms-engine/projects/staff?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const getProjectsByManager = (managerId: string) =>
  get<BaseAPIResponse<Project[]>>(`/pms-engine/projects/by-manager?managerId=${managerId}`);
export const getProjectsAssigned = (staffId: string) =>
  get<BaseAPIResponse<Project[]>>(`/pms-engine/projects/assigned?staffId=${staffId}`);
export const createProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects", data);
export const updateProject = (data: unknown) => put<ResponseVm>("/pms-engine/projects", data);
export const submitDraftProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/submit-draft", data);
export const saveDraftProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/draft", data);
export const approveProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/approve", data);
export const rejectProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/reject", data);
export const returnProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/return", data);
export const cancelProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/cancel", data);
export const closeProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/close", data);
export const pauseProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/pause", data);
export const addProjectMember = (data: unknown) => post<ResponseVm>("/pms-engine/projects/members", data);
export const removeProjectMember = (data: unknown) => post<ResponseVm>("/pms-engine/projects/members/remove", data);
export const getProjectMembers = (projectId: string) =>
  get<BaseAPIResponse<ProjectMember[]>>(`/pms-engine/projects/${projectId}/members`);
export const getProjectObjectives = (projectId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/projects/${projectId}/objectives`);
export const getProjectWorkProducts = (projectId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/projects/${projectId}/work-products`);

// --- Committees ---
export const getCommittees = (params?: { staffId?: string; reviewPeriodId?: string }) => {
  const q = new URLSearchParams();
  if (params?.staffId) q.append("staffId", params.staffId);
  if (params?.reviewPeriodId) q.append("reviewPeriodId", params.reviewPeriodId);
  return get<BaseAPIResponse<Committee[]>>(`/pms-engine/committees?${q.toString()}`);
};
export const getCommitteeDetails = (committeeId: string) =>
  get<BaseAPIResponse<Committee>>(`/pms-engine/committees/${committeeId}`);
export const getStaffCommittees = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<Committee[]>>(`/pms-engine/committees/staff?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const getCommitteesByChairperson = (chairpersonId: string) =>
  get<BaseAPIResponse<Committee[]>>(`/pms-engine/committees/by-chairperson?chairpersonId=${chairpersonId}`);
export const getCommitteesAssigned = (staffId: string) =>
  get<BaseAPIResponse<Committee[]>>(`/pms-engine/committees/assigned?staffId=${staffId}`);
export const createCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees", data);
export const updateCommittee = (data: unknown) => put<ResponseVm>("/pms-engine/committees", data);
export const submitDraftCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/submit-draft", data);
export const saveDraftCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/draft", data);
export const approveCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/approve", data);
export const rejectCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/reject", data);
export const returnCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/return", data);
export const cancelCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/cancel", data);
export const closeCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/close", data);
export const pauseCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/pause", data);
export const addCommitteeMember = (data: unknown) => post<ResponseVm>("/pms-engine/committees/members", data);
export const removeCommitteeMember = (data: unknown) =>
  post<ResponseVm>("/pms-engine/committees/members/remove", data);
export const getCommitteeMembers = (committeeId: string) =>
  get<BaseAPIResponse<CommitteeMember[]>>(`/pms-engine/committees/${committeeId}/members`);
export const getCommitteeObjectives = (committeeId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/committees/${committeeId}/objectives`);
export const getCommitteeWorkProducts = (committeeId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/committees/${committeeId}/work-products`);

// --- Work Products ---
export const getStaffWorkProducts = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/work-products?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const getAllStaffWorkProducts = () =>
  get<BaseAPIResponse<WorkProduct[]>>("/pms-engine/work-products/all");
export const getWorkProductDetails = (workProductId: string) =>
  get<BaseAPIResponse<WorkProduct>>(`/pms-engine/work-products/${workProductId}`);
export const createWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products", data);
export const updateWorkProduct = (data: unknown) => put<ResponseVm>("/pms-engine/work-products", data);
export const submitDraftWorkProduct = (data: unknown) =>
  post<ResponseVm>("/pms-engine/work-products/submit-draft", data);
export const saveDraftWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/draft", data);
export const approveWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/approve", data);
export const rejectWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/reject", data);
export const returnWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/return", data);
export const reSubmitWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/resubmit", data);
export const cancelWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/cancel", data);
export const completeWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/complete", data);
export const pauseWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/pause", data);
export const resumeWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/resume", data);
export const suspendWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/suspend", data);
export const reInstateWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/reinstate", data);

// --- Work Product Tasks ---
export const getWorkProductTasks = (workProductId: string) =>
  get<BaseAPIResponse<WorkProductTask[]>>(`/pms-engine/work-products/${workProductId}/tasks`);
export const getWorkProductTaskDetail = (taskId: string) =>
  get<BaseAPIResponse<WorkProductTask>>(`/pms-engine/work-products/tasks/${taskId}`);
export const createWorkProductTask = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/tasks", data);
export const updateWorkProductTask = (data: unknown) => put<ResponseVm>("/pms-engine/work-products/tasks", data);
export const completeWorkProductTask = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/tasks/complete", data);
export const cancelWorkProductTask = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/tasks/cancel", data);

// --- Work Product Evaluation ---
export const evaluateWorkProduct = (data: unknown) =>
  post<ResponseVm>("/pms-engine/work-products/evaluation", data);
export const updateWorkProductEvaluation = (data: unknown) =>
  put<ResponseVm>("/pms-engine/work-products/evaluation", data);
export const getWorkProductEvaluation = (workProductId: string) =>
  get<BaseAPIResponse<WorkProductEvaluation>>(`/pms-engine/work-products/${workProductId}/evaluation`);
export const initiateReEvaluation = (workProductId: string) =>
  post<ResponseVm>(`/pms-engine/work-products/${workProductId}/re-evaluate`, {});

// --- Individual Objectives ---
export const getStaffIndividualObjectives = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<IndividualPlannedObjective[]>>(`/pms-engine/individual-objectives?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const getStaffObjectives = getStaffIndividualObjectives;
export const saveDraftIndividualObjective = (data: unknown) => post<ResponseVm>("/pms-engine/individual-objectives/draft", data);
export const addIndividualObjective = (data: unknown) => post<ResponseVm>("/pms-engine/individual-objectives", data);
export const submitDraftIndividualObjective = (data: unknown) => post<ResponseVm>("/pms-engine/individual-objectives/submit-draft", data);
export const approveIndividualObjective = (data: unknown) => post<ResponseVm>("/pms-engine/individual-objectives/approve", data);
export const rejectIndividualObjective = (data: unknown) => post<ResponseVm>("/pms-engine/individual-objectives/reject", data);
export const returnIndividualObjective = (data: unknown) => post<ResponseVm>("/pms-engine/individual-objectives/return", data);
export const cancelIndividualObjective = (data: unknown) => post<ResponseVm>("/pms-engine/individual-objectives/cancel", data);

// --- Evaluations ---
export const getStaffEvaluations = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/evaluations?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const saveDraftEvaluation = (data: unknown) => post<ResponseVm>("/pms-engine/evaluations/draft", data);
export const addEvaluation = (data: unknown) => post<ResponseVm>("/pms-engine/evaluations", data);
export const submitDraftEvaluation = (data: unknown) => post<ResponseVm>("/pms-engine/evaluations/submit-draft", data);
export const approveEvaluation = (data: unknown) => post<ResponseVm>("/pms-engine/evaluations/approve", data);
export const rejectEvaluation = (data: unknown) => post<ResponseVm>("/pms-engine/evaluations/reject", data);

// --- Feedback Requests ---
export const getStaffFeedbackRequests = (staffId: string) =>
  get<BaseAPIResponse<FeedbackRequestLog[]>>(`/pms-engine/feedback/requests?staffId=${staffId}`);
export const completeFeedbackRequest = (data: unknown) =>
  post<ResponseVm>("/pms-engine/feedback/process", data);
export const getStaffRequests = (staffId: string) =>
  get<BaseAPIResponse<FeedbackRequestLog[]>>(`/pms-engine/feedback/requests/staff?staffId=${staffId}`);
export const getBreachedRequests = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<BreachedFeedbackRequestLog[]>>(`/pms-engine/feedback/requests/breached?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const getPendingFeedbackActions = (staffId: string) =>
  get<BaseAPIResponse<StaffPendingRequest[]>>(`/pms-engine/feedback/pending?staffId=${staffId}`);
export const getAllRequests = (staffId: string) =>
  get<BaseAPIResponse<FeedbackRequestLog[]>>(`/pms-engine/feedback/requests/all?staffId=${staffId}`);
export const getRequestDetails = (requestId: string) =>
  get<BaseAPIResponse<FeedbackRequestLog>>(`/pms-engine/feedback/requests/${requestId}`);
export const getRequestsByStatus = (staffId: string, status: string) =>
  get<BaseAPIResponse<FeedbackRequestLog[]>>(`/pms-engine/feedback/requests/by-status?staffId=${staffId}&status=${status}`);
export const reassignRequest = (data: { requestId: string; newAssignedStaffId: string }) =>
  post<ResponseVm>("/pms-engine/feedback/requests/reassign", data);
export const reassignSelfRequest = (data: { requestId: string; currentStaffId: string; newAssignedStaffId: string }) =>
  post<ResponseVm>("/pms-engine/feedback/requests/reassign-self", data);
export const closeRequest = (data: { requestId: string }) =>
  post<ResponseVm>("/pms-engine/feedback/requests/close", data);
export const treatAssignedRequest = (data: { requestId: string; operationType: number; comment?: string }) =>
  post<ResponseVm>("/pms-engine/feedback/requests/treat", data);

// --- 360 Review ---
export const trigger360Review = (data: unknown) =>
  post<ResponseVm>("/pms-engine/360-review/trigger", data);
export const initiate360Review = (data: { staffId: string[]; reviewPeriodId: string }) =>
  post<ResponseVm>("/pms-engine/360-review/initiate", data);
export const complete360ReviewForStaff = (data: { reviewPeriodId: string }) =>
  post<ResponseVm>("/pms-engine/360-review/complete", data);
export const add360Rating = (data: unknown) =>
  post<ResponseVm>("/pms-engine/360-review/rating", data);
export const update360Rating = (data: unknown) =>
  put<ResponseVm>("/pms-engine/360-review/rating", data);
export const reviewerComplete360Review = (data: { reviewStaffId: string; competencyReviewFeedbackId: string }) =>
  post<ResponseVm>("/pms-engine/360-review/reviewer-complete", data);

// --- Competency Review ---
export const getCompetencyReviewDetail = (feedbackId: string) =>
  get<BaseAPIResponse<CompetencyReviewFeedback>>(`/pms-engine/competency-review/${feedbackId}`);
export const getCompetencyReviewFeedbackDetails = (feedbackId: string) =>
  get<BaseAPIResponse<CompetencyReviewFeedbackDetails>>(`/pms-engine/competency-review/${feedbackId}/details`);
export const getAllCompetencyReviewFeedbacks = (staffId: string) =>
  get<BaseAPIResponse<CompetencyReviewFeedback[]>>(`/pms-engine/competency-review/feedbacks?staffId=${staffId}`);
export const getCompetenciesToReview = (reviewerStaffId: string) =>
  get<BaseAPIResponse<CompetencyReviewer[]>>(`/pms-engine/competency-review/to-review?reviewerStaffId=${reviewerStaffId}`);
export const getMyReviewedCompetencies = (reviewerStaffId: string) =>
  get<BaseAPIResponse<CompetencyReviewer[]>>(`/pms-engine/competency-review/my-reviewed?reviewerStaffId=${reviewerStaffId}`);
export const getReviewerFeedbackDetails = (reviewerId: string) =>
  get<BaseAPIResponse<CompetencyReviewer>>(`/pms-engine/competency-review/reviewer/${reviewerId}`);
export const getQuestionnaire = (staffId: string) =>
  get<BaseAPIResponse<PmsCompetency[]>>(`/pms-engine/competency-review/questionnaire?staffId=${staffId}`);
export const competencyGapClosureSetup = (data: unknown) =>
  post<ResponseVm>("/pms-engine/competency-review/gap-closure", data);

// --- Line Manager & Staff ---
export const getLineManagerEmployees = (staffId: string, category: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/line-manager-employees?staffId=${staffId}&category=${category}`);
export const getMyStaff = () =>
  get<BaseAPIResponse<unknown[]>>("/pms-engine/my-staff");

// --- Scoring ---
export const getPerformanceScore = (staffId: string) =>
  get<BaseAPIResponse<unknown>>(`/pms-engine/scores?staffId=${staffId}`);
export const getPerformanceSummary = () =>
  get<BaseAPIResponse<unknown>>("/pms-engine/scores/summary");
export const getDashboardStats = (staffId: string) =>
  get<BaseAPIResponse<unknown>>(`/pms-engine/dashboard?staffId=${staffId}`);

// --- Staff Project/Committee Work Products ---
export const getStaffProjectWorkProducts = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/work-products/project/staff?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const getStaffCommitteeWorkProducts = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/work-products/committee/staff?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
