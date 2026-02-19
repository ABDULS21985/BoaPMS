import { get } from "@/lib/api-client";
import type { BaseAPIResponse } from "@/types/common";
import type {
  FeedbackRequestDashboardStats,
  PerformancePointsStats,
  WorkProductDashboardStats,
  WorkProductDetailsDashboardStats,
  StaffScoreCardResponse,
  StaffAnnualScoreCardResponse,
  OrganogramPerformanceSummary,
  PendingAction,
  EmployeeErpDetails,
  SubordinatesScoreCardResponse,
  PeriodScoreData,
} from "@/types/dashboard";
import type { PerformanceReviewPeriod } from "@/types/performance";
import type { StaffIdMask } from "@/types/staff";

// --- Dashboard Stats ---
export const getDashboardStats = (staffId: string) =>
  get<StaffScoreCardResponse>(`/pms-engine/dashboard?staffId=${staffId}`);

export const getRequestStatistics = (staffId: string) =>
  get<BaseAPIResponse<FeedbackRequestDashboardStats>>(`/pms-engine/stats/requests?staffId=${staffId}`);

export const getPerformanceStatistics = (staffId: string) =>
  get<BaseAPIResponse<PerformancePointsStats>>(`/pms-engine/stats/performance?staffId=${staffId}`);

export const getWorkProductStatistics = (staffId: string) =>
  get<BaseAPIResponse<WorkProductDashboardStats>>(`/pms-engine/stats/work-products?staffId=${staffId}`);

export const getWorkProductDetailsStatistics = (staffId: string) =>
  get<BaseAPIResponse<WorkProductDetailsDashboardStats>>(`/pms-engine/stats/work-products-details?staffId=${staffId}`);

// --- Score Cards ---
export const getStaffScoreCard = (staffId: string, reviewPeriodId: string) =>
  get<StaffScoreCardResponse>(`/pms-engine/scorecard?staffId=${staffId}&reviewPeriodId=${reviewPeriodId}`);

export const getStaffAnnualScoreCard = (staffId: string, year: number) =>
  get<StaffAnnualScoreCardResponse>(`/pms-engine/scorecard/annual?staffId=${staffId}&year=${year}`);

// --- Organogram Performance ---
export const getOrganogramPerformanceSummary = (
  referenceId: string,
  reviewPeriodId: string,
  level?: number
) =>
  get<BaseAPIResponse<OrganogramPerformanceSummary>>(
    `/pms-engine/organogram-performance?referenceId=${referenceId}&reviewPeriodId=${reviewPeriodId}${level !== undefined ? `&level=${level}` : ""}`
  );

// --- Pending Actions ---
export const getStaffPendingRequests = (staffId: string) =>
  get<BaseAPIResponse<PendingAction[]>>(`/pms-engine/feedback-requests?staffId=${staffId}`);

// --- Active Review Period ---
export const getActiveReviewPeriod = () =>
  get<BaseAPIResponse<PerformanceReviewPeriod>>("/review-periods/active");

// --- Employee Details (ERP) ---
export const getEmployeeDetail = (userId: string) =>
  get<BaseAPIResponse<EmployeeErpDetails>>(`/employees/${userId}`);

export const getStaffIdMask = (userId: string) =>
  get<BaseAPIResponse<StaffIdMask>>(`/staff/${userId}/id-mask`);

// --- Subordinates ScoreCard ---
export const getSubordinatesScoreCard = (managerId: string, reviewPeriodId: string) =>
  get<SubordinatesScoreCardResponse>(
    `/pms-engine/scorecard/subordinates?managerId=${managerId}&reviewPeriodId=${reviewPeriodId}`
  );

// --- Organogram Performance List ---
export const getOrganogramPerformanceSummaryList = (
  headOfUnitId: string,
  reviewPeriodId: string,
  level?: number
) =>
  get<BaseAPIResponse<OrganogramPerformanceSummary[]>>(
    `/pms-engine/organogram-performance/list?headOfUnitId=${headOfUnitId}&reviewPeriodId=${reviewPeriodId}${level !== undefined ? `&level=${level}` : ""}`
  );

// --- Period Scores (for PMS Score Report) ---
export const getPeriodScores = (reviewPeriodId: string) =>
  get<BaseAPIResponse<PeriodScoreData[]>>(
    `/pms-engine/period-scores/all?reviewPeriodId=${reviewPeriodId}`
  );

export const getPeriodScoreDetails = (reviewPeriodId: string, staffId: string) =>
  get<BaseAPIResponse<PeriodScoreData>>(
    `/pms-engine/period-scores?reviewPeriodId=${reviewPeriodId}&staffId=${staffId}`
  );

// --- Head Subordinates ---
export const getHeadSubordinates = (staffId: string) =>
  get<BaseAPIResponse<unknown[]>>(
    `/pms-engine/line-manager-employees?staffId=${staffId}&category=direct`
  );

// --- Staff Review Periods ---
export const getStaffReviewPeriods = (staffId: string) =>
  get<BaseAPIResponse<import("@/types/performance").PerformanceReviewPeriod[]>>(
    `/pms-engine/staff-review-periods?staffId=${staffId}`
  );
