import { get, post, put, del } from "@/lib/api-client";
import type { BaseAPIResponse, ResponseVm } from "@/types/common";
import type {
  Competency,
  CompetencyCategory,
  CompetencyCategoryGrading,
  CompetencyReview,
  CompetencyReviewProfile,
  DevelopmentPlan,
  JobRole,
  JobGrade,
  JobGradeGroup,
  JobRoleCompetency,
  BehavioralCompetency,
  Rating,
  ReviewPeriod,
  ReviewType,
  TrainingType,
  StaffJobRole,
  OfficeJobRole,
} from "@/types/competency";

// --- Competencies ---
export const getCompetencies = () => get<BaseAPIResponse<Competency[]>>("/competency/competencies");
export const saveCompetency = (data: unknown) => post<ResponseVm>("/competency/competencies", data);
export const approveCompetency = (data: unknown) => post<ResponseVm>("/competency/competencies/approve", data);
export const rejectCompetency = (data: unknown) => post<ResponseVm>("/competency/competencies/reject", data);
export const getPendingCompetencies = () =>
  get<BaseAPIResponse<Competency[]>>("/competency/competencies/pending");

// --- Categories ---
export const getCompetencyCategories = () =>
  get<BaseAPIResponse<CompetencyCategory[]>>("/competency/categories");
export const saveCompetencyCategory = (data: unknown) =>
  post<ResponseVm>("/competency/categories", data);

// --- Category Grading ---
export const getCategoryGradings = (categoryId: number) =>
  get<BaseAPIResponse<CompetencyCategoryGrading[]>>(`/competency/categories/${categoryId}/gradings`);
export const saveCategoryGrading = (data: unknown) =>
  post<ResponseVm>("/competency/category-gradings", data);

// --- Ratings ---
export const getRatings = () => get<BaseAPIResponse<Rating[]>>("/competency/ratings");
export const saveRating = (data: unknown) => post<ResponseVm>("/competency/ratings", data);
export const getRatingDefinitions = (competencyId: number) =>
  get<BaseAPIResponse<unknown[]>>(`/competency/competencies/${competencyId}/rating-definitions`);
export const saveRatingDefinition = (data: unknown) =>
  post<ResponseVm>("/competency/rating-definitions", data);

// --- Review Types ---
export const getReviewTypes = () => get<BaseAPIResponse<ReviewType[]>>("/competency/review-types");
export const saveReviewType = (data: unknown) => post<ResponseVm>("/competency/review-types", data);

// --- Review Periods ---
export const getCompetencyReviewPeriods = () =>
  get<BaseAPIResponse<ReviewPeriod[]>>("/competency/review-periods");
export const saveCompetencyReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/competency/review-periods", data);
export const approveCompetencyReviewPeriod = (data: unknown) =>
  post<ResponseVm>("/competency/review-periods/approve", data);
export const getPendingReviewPeriods = () =>
  get<BaseAPIResponse<ReviewPeriod[]>>("/competency/review-periods/pending");

// --- Competency Reviews ---
export const getCompetencyReviews = (employeeNumber: string, reviewPeriodId: number) =>
  get<BaseAPIResponse<CompetencyReview[]>>(
    `/competency/reviews?employeeNumber=${employeeNumber}&reviewPeriodId=${reviewPeriodId}`
  );
export const saveCompetencyReview = (data: unknown) =>
  post<ResponseVm>("/competency/reviews", data);

// --- Competency Profiles ---
export const getCompetencyProfiles = (employeeNumber: string) =>
  get<BaseAPIResponse<CompetencyReviewProfile[]>>(
    `/competency/profiles?employeeNumber=${employeeNumber}`
  );
export const getCompetencyProfilesByReviewer = (reviewerNumber: string) =>
  get<BaseAPIResponse<CompetencyReviewProfile[]>>(
    `/competency/profiles/reviewer?reviewerNumber=${reviewerNumber}`
  );

// --- Development Plans ---
export const getDevelopmentPlans = (employeeNumber: string) =>
  get<BaseAPIResponse<DevelopmentPlan[]>>(
    `/competency/development-plans?employeeNumber=${employeeNumber}`
  );
export const saveDevelopmentPlan = (data: unknown) =>
  post<ResponseVm>("/competency/development-plans", data);

// --- Job Roles ---
export const getJobRoles = () => get<BaseAPIResponse<JobRole[]>>("/competency/job-roles");
export const saveJobRole = (data: unknown) => post<ResponseVm>("/competency/job-roles", data);
export const getOfficeJobRoles = (officeId: number) =>
  get<BaseAPIResponse<OfficeJobRole[]>>(`/competency/offices/${officeId}/job-roles`);
export const saveOfficeJobRole = (data: unknown) =>
  post<ResponseVm>("/competency/office-job-roles", data);

// --- Job Role Competencies ---
export const getJobRoleCompetencies = (jobRoleId: number, officeId: number) =>
  get<BaseAPIResponse<JobRoleCompetency[]>>(
    `/competency/job-role-competencies?jobRoleId=${jobRoleId}&officeId=${officeId}`
  );
export const saveJobRoleCompetency = (data: unknown) =>
  post<ResponseVm>("/competency/job-role-competencies", data);

// --- Job Grades ---
export const getJobGrades = () => get<BaseAPIResponse<JobGrade[]>>("/competency/job-grades");
export const saveJobGrade = (data: unknown) => post<ResponseVm>("/competency/job-grades", data);

// --- Job Grade Groups ---
export const getJobGradeGroups = () => get<BaseAPIResponse<JobGradeGroup[]>>("/competency/grade-groups");
export const saveJobGradeGroup = (data: unknown) => post<ResponseVm>("/competency/grade-groups", data);
export const assignJobGradeGroup = (data: unknown) =>
  post<ResponseVm>("/competency/assign-grade-groups", data);

// --- Behavioral Competencies ---
export const getBehavioralCompetencies = (gradeGroupId: number) =>
  get<BaseAPIResponse<BehavioralCompetency[]>>(
    `/competency/behavioral-competencies?gradeGroupId=${gradeGroupId}`
  );
export const saveBehavioralCompetency = (data: unknown) =>
  post<ResponseVm>("/competency/behavioral-competencies", data);

// --- Training Types ---
export const getTrainingTypes = () =>
  get<BaseAPIResponse<TrainingType[]>>("/competency/training-types");
export const saveTrainingType = (data: unknown) =>
  post<ResponseVm>("/competency/training-types", data);

// --- Staff Job Roles ---
export const getStaffJobRoles = (employeeId: string) =>
  get<BaseAPIResponse<StaffJobRole[]>>(`/competency/staff/${employeeId}/job-roles`);
export const saveStaffJobRole = (data: unknown) =>
  post<ResponseVm>("/competency/staff-job-roles", data);

// --- Bank Years ---
export const getBankYears = () => get<BaseAPIResponse<unknown[]>>("/competency/bank-years");
export const saveBankYear = (data: unknown) => post<ResponseVm>("/competency/bank-years", data);
