import { get, post, put, del } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type {
  Strategy,
  StrategicTheme,
  ObjectiveCategory,
  EnterpriseObjective,
  DepartmentObjective,
  DivisionObjective,
  OfficeObjective,
  ConsolidatedObjective,
  EvaluationOption,
  WorkProductDefinition,
  FeedbackQuestionnaire,
  PmsConfiguration,
  Setting,
} from "@/types/performance";

// --- Strategies ---
export const getStrategies = () => get<BaseAPIResponse<Strategy[]>>("/performance/strategies");
export const createStrategy = (data: unknown) => post<ResponseVm>("/performance/strategies", data);
export const updateStrategy = (data: unknown) => put<ResponseVm>("/performance/strategies", data);

// --- Strategic Themes ---
export const getAllStrategicThemes = () =>
  get<BaseAPIResponse<StrategicTheme[]>>("/performance/strategic-themes");
export const getStrategicThemes = (strategyId?: string) =>
  get<BaseAPIResponse<StrategicTheme[]>>(`/performance/strategic-themes${strategyId ? `?strategyId=${strategyId}` : ""}`);
export const createStrategicTheme = (data: unknown) =>
  post<ResponseVm>("/performance/strategic-themes", data);
export const updateStrategicTheme = (data: unknown) =>
  put<ResponseVm>("/performance/strategic-themes", data);

// --- Objective Categories ---
export const getObjectiveCategories = () =>
  get<BaseAPIResponse<ObjectiveCategory[]>>("/performance/objective-categories");
export const createObjectiveCategory = (data: unknown) =>
  post<ResponseVm>("/performance/objective-categories", data);
export const updateObjectiveCategory = (data: unknown) =>
  put<ResponseVm>("/performance/objective-categories", data);

// --- Enterprise Objectives ---
export const getEnterpriseObjectives = () =>
  get<BaseAPIResponse<EnterpriseObjective[]>>("/performance/objectives/enterprise");
export const createEnterpriseObjective = (data: unknown) =>
  post<ResponseVm>("/performance/objectives/enterprise", data);
export const updateEnterpriseObjective = (data: unknown) =>
  put<ResponseVm>("/performance/objectives/enterprise", data);

// --- Department Objectives ---
export const getDepartmentObjectives = () =>
  get<BaseAPIResponse<DepartmentObjective[]>>("/performance/objectives/department");
export const createDepartmentObjective = (data: unknown) =>
  post<ResponseVm>("/performance/objectives/department", data);
export const updateDepartmentObjective = (data: unknown) =>
  put<ResponseVm>("/performance/objectives/department", data);

// --- Division Objectives ---
export const getDivisionObjectives = (divisionId?: number) =>
  get<BaseAPIResponse<DivisionObjective[]>>(`/performance/objectives/division${divisionId ? `?divisionId=${divisionId}` : ""}`);
export const createDivisionObjective = (data: unknown) =>
  post<ResponseVm>("/performance/objectives/division", data);
export const updateDivisionObjective = (data: unknown) =>
  put<ResponseVm>("/performance/objectives/division", data);

// --- Office Objectives ---
export const getOfficeObjectives = (officeId?: number) =>
  get<BaseAPIResponse<OfficeObjective[]>>(`/performance/objectives/office${officeId ? `?officeId=${officeId}` : ""}`);
export const createOfficeObjective = (data: unknown) =>
  post<ResponseVm>("/performance/objectives/office", data);
export const updateOfficeObjective = (data: unknown) =>
  put<ResponseVm>("/performance/objectives/office", data);

// --- Consolidated Objectives ---
export const getConsolidatedObjectives = () =>
  get<BaseAPIResponse<{ objectives: ConsolidatedObjective[]; totalRecords: number }>>("/performance/objectives/consolidated");
export const getConsolidatedObjectivesPaginated = (params: {
  pageIndex: number;
  pageSize: number;
  searchString?: string;
  departmentId?: number;
  divisionId?: number;
  officeId?: number;
  status?: string;
}) => {
  const q = new URLSearchParams();
  q.append("pageIndex", String(params.pageIndex));
  q.append("pageSize", String(params.pageSize));
  if (params.searchString) q.append("searchString", params.searchString);
  if (params.departmentId) q.append("departmentId", String(params.departmentId));
  if (params.divisionId) q.append("divisionId", String(params.divisionId));
  if (params.officeId) q.append("officeId", String(params.officeId));
  if (params.status) q.append("status", params.status);
  return get<BaseAPIResponse<{ objectives: ConsolidatedObjective[]; totalRecords: number }>>(
    `/performance/objectives/consolidated/paginated?${q.toString()}`
  );
};

// --- Approval / Rejection ---
export const approveRecords = (data: { entityType: string; recordIds: string[] }) =>
  post<ResponseVm>("/performance/approve", data);
export const rejectRecords = (data: { entityType: string; recordIds: string[]; rejectionReason: string }) =>
  post<ResponseVm>("/performance/reject", data);

// --- Evaluation Options ---
export const getEvaluationOptions = () =>
  get<BaseAPIResponse<EvaluationOption[]>>("/performance/evaluation-options");
export const saveEvaluationOptions = (data: unknown[]) =>
  post<ResponseVm>("/performance/evaluation-options", data);

// --- Work Product Definitions ---
export const getAllWorkProductDefinitions = () =>
  get<BaseAPIResponse<WorkProductDefinition[]>>("/performance/work-product-definitions/all");
export const getWorkProductDefinitionsPaginated = (params: { pageIndex: number; pageSize: number; search?: string }) =>
  get<BaseAPIResponse<{ items: WorkProductDefinition[]; totalRecords: number }>>(
    `/performance/work-product-definitions/paginated?pageIndex=${params.pageIndex}&pageSize=${params.pageSize}${params.search ? `&search=${params.search}` : ""}`
  );
export const saveWorkProductDefinitions = (data: unknown[]) =>
  post<ResponseVm>("/performance/work-product-definitions", data);

// --- Feedback Questionnaires ---
export const getFeedbackQuestionnaires = () =>
  get<BaseAPIResponse<FeedbackQuestionnaire[]>>("/performance/feedback-questionnaires");
export const saveFeedbackQuestionnaires = (data: unknown[]) =>
  post<ResponseVm>("/performance/feedback-questionnaires", data);

// --- PMS Competencies ---
export interface PmsCompetency {
  pmsCompetencyId: string;
  name: string;
  description?: string;
  objectCategoryId: string;
  objectCategoryName?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}
export const getPmsCompetencies = () =>
  get<BaseAPIResponse<PmsCompetency[]>>("/performance/competencies");
export const createPmsCompetency = (data: unknown) =>
  post<ResponseVm>("/performance/competencies", data);
export const updatePmsCompetency = (data: unknown) =>
  put<ResponseVm>("/performance/competencies", data);

// --- PMS Configurations ---
export const getPmsConfigurations = () =>
  get<BaseAPIResponse<PmsConfiguration[]>>("/setup/pms-configurations");
export const createPmsConfiguration = (data: unknown) =>
  post<ResponseVm>("/setup/pms-configurations", data);
export const updatePmsConfiguration = (data: unknown) =>
  put<ResponseVm>("/setup/pms-configurations", data);

// --- Settings ---
export const getSettings = () => get<BaseAPIResponse<Setting[]>>("/setup/settings");
export const createSetting = (data: unknown) => post<ResponseVm>("/setup/settings", data);
export const updateSetting = (data: unknown) => put<ResponseVm>("/setup/settings", data);

// --- Enums ---
export const getEvaluationTypes = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/evaluation-types");
export const getFeedbackRequestTypes = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/feedback-request-types");
export const getGrievanceTypes = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/grievance-types");
export const getObjectiveLevels = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/objective-levels");
export const getWorkProductTypes = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/work-product-types");
export const getPerformanceGrades = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/performance-grades");
export const getReviewPeriodRanges = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/review-period-ranges");
export const getExtensionTargetTypes = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/extension-target-types");
export const getStatuses = () => get<BaseAPIResponse<unknown[]>>("/performance/enums/statuses");
