import { get, post } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type {
  Competency,
  CompetencyCategory,
  CompetencyCategoryGrading,
  CompetencyRatingDefinition,
  CompetencyReview,
  CompetencyReviewProfile,
  DevelopmentPlan,
  JobRole,
  JobGrade,
  JobGradeGroup,
  AssignJobGradeGroup,
  JobRoleCompetency,
  JobRoleGrade,
  BehavioralCompetency,
  Rating,
  CompetencyReviewPeriod,
  ReviewType,
  TrainingType,
  OfficeJobRole,
  BankYear,
  OfficeCompetencyReview,
} from "@/types/competency";

// --- Competencies ---
export const getCompetencies = () =>
  get<BaseAPIResponse<Competency[]>>("/competency/competencies");
export const saveCompetency = (data: unknown) =>
  post<ResponseVm>("/competency/competencies", data);
export const approveCompetency = (data: unknown) =>
  post<ResponseVm>("/competency/competencies/approve", data);
export const rejectCompetency = (data: unknown) =>
  post<ResponseVm>("/competency/competencies/reject", data);

// --- Categories ---
export const getCompetencyCategories = () =>
  get<BaseAPIResponse<CompetencyCategory[]>>("/competency/categories");
export const saveCompetencyCategory = (data: unknown) =>
  post<ResponseVm>("/competency/categories", data);

// --- Category Grading ---
export const getCategoryGradings = () =>
  get<BaseAPIResponse<CompetencyCategoryGrading[]>>("/competency/category-gradings");
export const saveCategoryGrading = (data: unknown) =>
  post<ResponseVm>("/competency/category-gradings", data);

// --- Ratings ---
export const getRatings = () =>
  get<BaseAPIResponse<Rating[]>>("/competency/ratings");
export const saveRating = (data: unknown) =>
  post<ResponseVm>("/competency/ratings", data);

// --- Rating Definitions ---
export const getRatingDefinitions = (competencyId?: number) =>
  get<BaseAPIResponse<CompetencyRatingDefinition[]>>(
    `/competency/rating-definitions${competencyId ? `?competencyId=${competencyId}` : ""}`
  );
export const saveRatingDefinition = (data: unknown) =>
  post<ResponseVm>("/competency/rating-definitions", data);

// --- Review Types ---
export const getReviewTypes = () =>
  get<BaseAPIResponse<ReviewType[]>>("/competency/review-types");
export const saveReviewType = (data: unknown) =>
  post<ResponseVm>("/competency/review-types", data);

// --- Review Periods ---
export const getCompetencyReviewPeriods = () =>
  get<BaseAPIResponse<CompetencyReviewPeriod[]>>("/competency/review-periods");
export const saveCompetencyReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/competency/review-periods", data);
export const approveCompetencyReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/competency/review-periods/approve", data);

// --- Competency Reviews ---
export const getCompetencyReviews = () =>
  get<BaseAPIResponse<CompetencyReview[]>>("/competency/reviews");
export const getCompetencyReviewByReviewer = (reviewerId: string, reviewPeriodId?: number) =>
  get<BaseAPIResponse<CompetencyReview[]>>(
    `/competency/reviews/by-reviewer?reviewerId=${reviewerId}${reviewPeriodId ? `&reviewPeriodId=${reviewPeriodId}` : ""}`
  );
export const getCompetencyReviewForEmployee = (employeeNumber: string, reviewPeriodId?: number) =>
  get<BaseAPIResponse<CompetencyReview[]>>(
    `/competency/reviews/for-employee?employeeNumber=${employeeNumber}${reviewPeriodId ? `&reviewPeriodId=${reviewPeriodId}` : ""}`
  );
export const getCompetencyReviewDetail = (data: unknown) =>
  get<BaseAPIResponse<CompetencyReview>>("/competency/reviews/detail");
export const saveCompetencyReview = (data: unknown) =>
  post<ResponseVm>("/competency/reviews", data);
export const getOfficeCompetencyReviews = (officeId: number, reviewPeriodId?: number) =>
  get<BaseAPIResponse<OfficeCompetencyReview[]>>(
    `/competency/reviews/by-office?officeId=${officeId}${reviewPeriodId ? `&reviewPeriodId=${reviewPeriodId}` : ""}`
  );

// --- Competency Review Profiles ---
export const getCompetencyReviewProfiles = (employeeNumber: string, reviewPeriodId?: number) =>
  get<BaseAPIResponse<CompetencyReviewProfile[]>>(
    `/competency/review-profiles?employeeNumber=${employeeNumber}${reviewPeriodId ? `&reviewPeriodId=${reviewPeriodId}` : ""}`
  );
export const getGroupCompetencyReviewProfiles = (params: {
  reviewPeriodId?: number; officeId?: number; divisionId?: number; departmentId?: number;
}) => {
  const q = new URLSearchParams();
  if (params.reviewPeriodId) q.set("reviewPeriodId", String(params.reviewPeriodId));
  if (params.officeId) q.set("officeId", String(params.officeId));
  if (params.divisionId) q.set("divisionId", String(params.divisionId));
  if (params.departmentId) q.set("departmentId", String(params.departmentId));
  return get<BaseAPIResponse<CompetencyReviewProfile[]>>(`/competency/review-profiles/group?${q}`);
};
export const getCompetencyMatrixReviewProfiles = (params: {
  reviewPeriodId?: number; officeId?: number; divisionId?: number; departmentId?: number;
}) => {
  const q = new URLSearchParams();
  if (params.reviewPeriodId) q.set("reviewPeriodId", String(params.reviewPeriodId));
  if (params.officeId) q.set("officeId", String(params.officeId));
  if (params.divisionId) q.set("divisionId", String(params.divisionId));
  if (params.departmentId) q.set("departmentId", String(params.departmentId));
  return get<BaseAPIResponse<CompetencyReviewProfile[]>>(`/competency/review-profiles/matrix?${q}`);
};
export const getTechnicalCompetencyMatrixReviewProfiles = (reviewPeriodId: number, jobRoleId: number) =>
  get<BaseAPIResponse<CompetencyReviewProfile[]>>(
    `/competency/review-profiles/technical-matrix?reviewPeriodId=${reviewPeriodId}&jobRoleId=${jobRoleId}`
  );
export const saveCompetencyReviewProfile = (data: unknown) =>
  post<ResponseVm>("/competency/review-profiles", data);

// --- Development Plans ---
export const getDevelopmentPlans = (competencyProfileReviewId?: number) =>
  get<BaseAPIResponse<DevelopmentPlan[]>>(
    `/competency/development-plans${competencyProfileReviewId ? `?competencyProfileReviewId=${competencyProfileReviewId}` : ""}`
  );
export const saveDevelopmentPlan = (data: unknown) =>
  post<ResponseVm>("/competency/development-plans", data);

// --- Job Roles ---
export const getJobRoles = () =>
  get<BaseAPIResponse<JobRole[]>>("/competency/job-roles");
export const saveJobRole = (data: unknown) =>
  post<ResponseVm>("/competency/job-roles", data);

// --- Office Job Roles ---
export const getOfficeJobRoles = () =>
  get<BaseAPIResponse<OfficeJobRole[]>>("/competency/office-job-roles");
export const saveOfficeJobRole = (data: unknown) =>
  post<ResponseVm>("/competency/office-job-roles", data);

// --- Job Role Competencies ---
export const getJobRoleCompetencies = () =>
  get<BaseAPIResponse<JobRoleCompetency[]>>("/competency/job-role-competencies");
export const saveJobRoleCompetency = (data: unknown) =>
  post<ResponseVm>("/competency/job-role-competencies", data);

// --- Behavioral Competencies ---
export const getBehavioralCompetencies = () =>
  get<BaseAPIResponse<BehavioralCompetency[]>>("/competency/behavioral");
export const saveBehavioralCompetency = (data: unknown) =>
  post<ResponseVm>("/competency/behavioral", data);

// --- Job Role Grades ---
export const getJobRoleGrades = () =>
  get<BaseAPIResponse<JobRoleGrade[]>>("/competency/job-role-grades");
export const saveJobRoleGrade = (data: unknown) =>
  post<ResponseVm>("/competency/job-role-grades", data);

// --- Job Grades ---
export const getJobGrades = () =>
  get<BaseAPIResponse<JobGrade[]>>("/competency/job-grades");
export const saveJobGrade = (data: unknown) =>
  post<ResponseVm>("/competency/job-grades", data);

// --- Job Grade Groups ---
export const getJobGradeGroups = () =>
  get<BaseAPIResponse<JobGradeGroup[]>>("/competency/job-grade-groups");
export const saveJobGradeGroup = (data: unknown) =>
  post<ResponseVm>("/competency/job-grade-groups", data);

// --- Assign Job Grade Groups ---
export const getAssignJobGradeGroups = () =>
  get<BaseAPIResponse<AssignJobGradeGroup[]>>("/competency/assign-job-grade-groups");
export const saveAssignJobGradeGroup = (data: unknown) =>
  post<ResponseVm>("/competency/assign-job-grade-groups", data);

// --- Training Types ---
export const getTrainingTypes = () =>
  get<BaseAPIResponse<TrainingType[]>>("/competency/training-types");
export const saveTrainingType = (data: unknown) =>
  post<ResponseVm>("/competency/training-types", data);

// --- Bank Years ---
export const getBankYears = () =>
  get<BaseAPIResponse<BankYear[]>>("/competency/bank-years");
export const saveBankYear = (data: unknown) =>
  post<ResponseVm>("/competency/bank-years", data);

// --- Population & Calculation ---
export const populateAllReviews = () =>
  post<ResponseVm>("/competency/populate/all-reviews", {});
export const populateOfficeReviews = (officeId: number) =>
  post<ResponseVm>("/competency/populate/office-reviews", { officeId });
export const populateDivisionReviews = (divisionId: number) =>
  post<ResponseVm>("/competency/populate/division-reviews", { divisionId });
export const populateDepartmentReviews = (departmentId: number) =>
  post<ResponseVm>("/competency/populate/department-reviews", { departmentId });
export const populateEmployeeReviews = (employeeNumber: string) =>
  post<ResponseVm>("/competency/populate/employee-reviews", { employeeNumber });
export const calculateReviews = (data: unknown) =>
  post<ResponseVm>("/competency/calculate-reviews", data);
export const recalculateReviewProfiles = (data: unknown) =>
  post<ResponseVm>("/competency/recalculate-review-profiles", data);
export const syncJobRoleSoa = (data: unknown) =>
  post<ResponseVm>("/competency/sync-job-role-soa", data);
