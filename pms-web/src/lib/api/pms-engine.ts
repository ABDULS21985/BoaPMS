import { get, post, put } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type {
  Project,
  Committee,
  WorkProduct,
  WorkProductTask,
} from "@/types/performance";

// --- Projects ---
export const getStaffProjects = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<Project[]>>(`/pms-engine/projects?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const createProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects", data);
export const updateProject = (data: unknown) => put<ResponseVm>("/pms-engine/projects", data);
export const submitDraftProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/submit-draft", data);
export const saveDraftProject = (data: unknown) => post<ResponseVm>("/pms-engine/projects/draft", data);
export const addProjectMember = (data: unknown) => post<ResponseVm>("/pms-engine/projects/members", data);
export const removeProjectMember = (data: unknown) => post<ResponseVm>("/pms-engine/projects/members/remove", data);
export const getProjectMembers = (projectId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/projects/${projectId}/members`);
export const getProjectWorkProducts = (projectId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/projects/${projectId}/work-products`);

// --- Committees ---
export const getStaffCommittees = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<Committee[]>>(`/pms-engine/committees?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const createCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees", data);
export const updateCommittee = (data: unknown) => put<ResponseVm>("/pms-engine/committees", data);
export const submitDraftCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/submit-draft", data);
export const saveDraftCommittee = (data: unknown) => post<ResponseVm>("/pms-engine/committees/draft", data);
export const addCommitteeMember = (data: unknown) => post<ResponseVm>("/pms-engine/committees/members", data);
export const removeCommitteeMember = (data: unknown) =>
  post<ResponseVm>("/pms-engine/committees/members/remove", data);
export const getCommitteeMembers = (committeeId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/committees/${committeeId}/members`);
export const getCommitteeWorkProducts = (committeeId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/committees/${committeeId}/work-products`);

// --- Work Products ---
export const getStaffWorkProducts = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<WorkProduct[]>>(`/pms-engine/work-products?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const createWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products", data);
export const submitDraftWorkProduct = (data: unknown) =>
  post<ResponseVm>("/pms-engine/work-products/submit-draft", data);
export const saveDraftWorkProduct = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/draft", data);

// --- Work Product Tasks ---
export const getWorkProductTasks = (workProductId: string) =>
  get<BaseAPIResponse<WorkProductTask[]>>(`/pms-engine/work-products/${workProductId}/tasks`);
export const createWorkProductTask = (data: unknown) => post<ResponseVm>("/pms-engine/work-products/tasks", data);

// --- Work Product Evaluation ---
export const evaluateWorkProduct = (data: unknown) =>
  post<ResponseVm>("/pms-engine/work-products/evaluate", data);
export const getWorkProductEvaluation = (workProductId: string) =>
  get<BaseAPIResponse<unknown>>(`/pms-engine/work-products/${workProductId}/evaluation`);

// --- Feedback ---
export const getStaffFeedbackRequests = (staffId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/feedback-requests?staffId=${staffId}`);
export const completeFeedbackRequest = (data: unknown) =>
  post<ResponseVm>("/pms-engine/feedback-requests/complete", data);

// --- Individual Objectives ---
export const getStaffObjectives = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/pms-engine/objectives?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);
export const saveDraftObjective = (data: unknown) => post<ResponseVm>("/pms-engine/objectives/draft", data);
export const submitDraftObjective = (data: unknown) => post<ResponseVm>("/pms-engine/objectives/submit-draft", data);

// --- Scoring ---
export const getPerformanceScore = (staffId: string) =>
  get<BaseAPIResponse<unknown>>(`/pms-engine/scores?staffId=${staffId}`);
export const getDashboardStats = (staffId: string) =>
  get<BaseAPIResponse<unknown>>(`/pms-engine/dashboard?staffId=${staffId}`);
