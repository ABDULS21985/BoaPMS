import { get, post, put } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type {
  PerformanceReviewPeriod,
  CategoryDefinition,
  ReviewPeriodExtension,
} from "@/types/performance";

// --- Review Period CRUD & lifecycle ---
export const getReviewPeriods = () =>
  get<BaseAPIResponse<PerformanceReviewPeriod[]>>("/review-periods/all");
export const getActiveReviewPeriod = () =>
  get<BaseAPIResponse<PerformanceReviewPeriod>>("/review-periods/active");
export const getStaffActiveReviewPeriod = () =>
  get<BaseAPIResponse<PerformanceReviewPeriod>>("/review-periods/staff-active");
export const getReviewPeriodDetails = (id: string) =>
  get<BaseAPIResponse<PerformanceReviewPeriod>>(`/review-periods/${id}`);

export const saveDraftReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods/draft", data);
export const addReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods", data);
export const updateReviewPeriod = (data: unknown) =>
  put<ResponseVm>("/review-periods", data);
export const submitDraftReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods/submit-draft", data);
export const approveReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods/approve", data);
export const rejectReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods/reject", data);
export const returnReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods/return", data);
export const closeReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods/close", data);
export const cancelReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/review-periods/cancel", data);

// --- Flags ---
export const enableObjectivePlanning = (data: unknown) =>
  post<ResponseVm>("/review-periods/enable-objective-planning", data);
export const disableObjectivePlanning = (data: unknown) =>
  post<ResponseVm>("/review-periods/disable-objective-planning", data);
export const enableWorkProductPlanning = (data: unknown) =>
  post<ResponseVm>("/review-periods/enable-work-product-planning", data);
export const disableWorkProductPlanning = (data: unknown) =>
  post<ResponseVm>("/review-periods/disable-work-product-planning", data);
export const enableWorkProductEvaluation = (data: unknown) =>
  post<ResponseVm>("/review-periods/enable-work-product-evaluation", data);
export const disableWorkProductEvaluation = (data: unknown) =>
  post<ResponseVm>("/review-periods/disable-work-product-evaluation", data);

// --- Category Definitions ---
export const getReviewPeriodCategoryDefinitions = (reviewPeriodId: string) =>
  get<BaseAPIResponse<CategoryDefinition[]>>(`/review-periods/${reviewPeriodId}/category-definitions`);
export const saveDraftCategoryDefinition = (data: unknown) =>
  post<ResponseVm>("/review-periods/category-definitions/draft", data);
export const addCategoryDefinition = (data: unknown) =>
  post<ResponseVm>("/review-periods/category-definitions", data);
export const submitDraftCategoryDefinition = (data: unknown) =>
  post<ResponseVm>("/review-periods/category-definitions/submit-draft", data);
export const approveCategoryDefinition = (data: unknown) =>
  post<ResponseVm>("/review-periods/category-definitions/approve", data);
export const rejectCategoryDefinition = (data: unknown) =>
  post<ResponseVm>("/review-periods/category-definitions/reject", data);

// --- Objectives (Period) ---
export const getReviewPeriodObjectives = (reviewPeriodId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/review-periods/${reviewPeriodId}/objectives`);
export const saveDraftReviewPeriodObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/objectives/draft", data);
export const addReviewPeriodObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/objectives", data);
export const submitDraftReviewPeriodObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/objectives/submit-draft", data);

// --- Extensions ---
export const getReviewPeriodExtensions = (reviewPeriodId: string) =>
  get<BaseAPIResponse<ReviewPeriodExtension[]>>(`/review-periods/${reviewPeriodId}/extensions`);
export const getAllReviewPeriodExtensions = () =>
  get<BaseAPIResponse<ReviewPeriodExtension[]>>("/review-periods/extensions/all");
export const saveDraftExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions/draft", data);
export const addExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions", data);
export const updateExtension = (data: unknown) =>
  put<ResponseVm>("/review-periods/extensions", data);
export const submitDraftExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions/submit-draft", data);
export const approveExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions/approve", data);
export const rejectExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions/reject", data);
export const returnExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions/return", data);
export const cancelExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions/cancel", data);
export const closeExtension = (data: unknown) =>
  post<ResponseVm>("/review-periods/extensions/close", data);

// --- 360 Reviews ---
export const addReviewPeriod360Review = (data: unknown) =>
  post<ResponseVm>("/review-periods/360-reviews", data);
export const getReviewPeriod360Reviews = (reviewPeriodId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/review-periods/${reviewPeriodId}/360-reviews`);

// --- Individual Planned Objectives ---
export const saveDraftIndividualPlannedObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/individual-objectives/draft", data);
export const addIndividualPlannedObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/individual-objectives", data);
export const submitDraftIndividualPlannedObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/individual-objectives/submit-draft", data);
export const approveIndividualPlannedObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/individual-objectives/approve", data);
export const rejectIndividualPlannedObjective = (data: unknown) =>
  post<ResponseVm>("/review-periods/individual-objectives/reject", data);

// --- Period Objective Evaluations ---
export const createPeriodObjectiveEvaluation = (data: unknown) =>
  post<ResponseVm>("/review-periods/evaluations", data);
export const createPeriodObjectiveDepartmentEvaluation = (data: unknown) =>
  post<ResponseVm>("/review-periods/evaluations/department", data);
export const getPeriodObjectiveEvaluations = (reviewPeriodId: string) =>
  get<BaseAPIResponse<unknown[]>>(`/review-periods/${reviewPeriodId}/evaluations`);

// --- Period Scores ---
export const getStaffPeriodScore = (staffId: string, reviewPeriodId: string) =>
  get<BaseAPIResponse<unknown>>(`/review-periods/scores?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);

// --- Archive ---
export const archiveCancelledObjectives = (data: unknown) =>
  post<ResponseVm>("/review-periods/archive-objectives", data);
export const archiveCancelledWorkProducts = (data: unknown) =>
  post<ResponseVm>("/review-periods/archive-workproducts", data);
