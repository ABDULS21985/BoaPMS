// --- Dashboard Statistics DTOs (matching Go response VMs) ---

export interface FeedbackRequestDashboardStats {
  staffId: string;
  reviewPeriodId: string;
  completedRequests: number;
  completedOverdueRequests: number;
  pendingRequests: number;
  pendingOverdueRequests: number;
  breachedRequests: number;
  pending360FeedbacksToTreat: number;
  deductedPoints: number;
}

export interface PerformancePointsStats {
  staffId: string;
  reviewPeriodId: string;
  maxPoints: number;
  accumulatedPoints: number;
  deductedPoints: number;
  actualPoints: number;
}

export interface WorkProductDashboardStats {
  staffId: string;
  reviewPeriodId: string;
  noAllWorkProducts: number;
  totalWorkProductTasks: number;
  noActiveWorkProducts: number;
  noWorkProductsAwaitingEvaluation: number;
  noWorkProductsClosed: number;
  noWorkProductsPendingApproval: number;
}

export interface WorkProductDashDetails {
  workProductId: string;
  name: string;
  description?: string;
  maxPoint: number;
  workProductType: number;
  workProductTypeName?: string;
  isSelfCreated: boolean;
  staffId: string;
  startDate: string;
  endDate: string;
  deliverables?: string;
  finalScore: number;
  completionDate?: string;
  reviewPeriodId: string;
  objectiveName?: string;
  percentageTaskCompletion: number;
  tasksCompleted: number;
  totalTasks: number;
  recordStatus?: number;
}

export interface WorkProductDetailsDashboardStats extends WorkProductDashboardStats {
  workProducts: WorkProductDashDetails[];
}

export interface StaffScoreCardDetails {
  staffId: string;
  staffName: string;
  reviewPeriodId: string;
  reviewPeriod: string;
  reviewPeriodShortName?: string;
  year: number;
  totalWorkProducts: number;
  percentageWorkProductsCompletion: number;
  totalCompetencyGaps: number;
  totalCompetencyGapsClosed: number;
  percentageGapsClosure: number;
  percentageGapsClosureScore: number;
  maxPoints: number;
  accumulatedPoints: number;
  deductedPoints: number;
  actualPoints: number;
  percentageScore: number;
  pmsCompetencyCategory: Record<string, number>;
  staffPerformanceGrade: string;
  totalWorkProductsCompletedOnSchedule: number;
  totalWorkProductsBehindSchedule: number;
  pmsCompetencies: StaffLivingTheValueRating[];
}

export interface StaffLivingTheValueRating {
  staffId: string;
  reviewPeriodId: string;
  pmsCompetencyId: string;
  objectiveCategoryId: string;
  pmsCompetency: string;
  ratingScore: number;
}

export interface StaffScoreCardResponse {
  isSuccess: boolean;
  message: string;
  scoreCard?: StaffScoreCardDetails;
}

export interface StaffAnnualScoreCardResponse {
  isSuccess: boolean;
  message: string;
  staffId: string;
  year: number;
  scoreCards: StaffScoreCardDetails[];
}

export interface OrganogramPerformanceSummary {
  referenceId: string;
  managerId: string;
  referenceName: string;
  reviewPeriodId: string;
  reviewPeriod: string;
  reviewPeriodShortName?: string;
  maxPoint: number;
  year: number;
  actualScore: number;
  performanceScore: number;
  totalWorkProducts: number;
  totalStaff: number;
  totalWorkProductsCompletedOnSchedule: number;
  totalWorkProductsBehindSchedule: number;
  total360Feedbacks: number;
  completed360FeedbacksToTreat: number;
  pending360FeedbacksToTreat: number;
  totalCompetencyGaps: number;
  percentageGapsClosure: number;
  percentageWorkProductsClosed: number;
  percentageWorkProductsPending: number;
  organogramLevel: number;
  earnedPerformanceGrade: string;
}

// --- Pending Actions ---
export interface PendingAction {
  id: string;
  name: string;
  description: string;
  type: string;
  referenceId: string;
  assignedDate: string;
  dueDate?: string;
  status: number;
  feedbackRequestType?: number;
}

// --- Employee ERP details ---
export interface EmployeeErpDetails {
  employeeNumber: string;
  firstName: string;
  lastName: string;
  email?: string;
  jobTitle?: string;
  gradeName?: string;
  departmentId?: number;
  departmentName?: string;
  divisionId?: number;
  divisionName?: string;
  officeId?: number;
  officeName?: string;
  supervisorId?: string;
  supervisorName?: string;
  headOfOfficeId?: string;
  headOfOfficeName?: string;
  headOfDivisionId?: string;
  headOfDepartmentId?: string;
  phone?: string;
  photoUrl?: string;
}

// --- Subordinates ScoreCard ---
export interface SubordinateScoreCard {
  staffId: string;
  staffName: string;
  totalWorkProducts: number;
  percentageWorkProductsCompletion: number;
  actualPoints: number;
  maxPoints: number;
  percentageScore: number;
  staffPerformanceGrade: string;
  totalWorkProductsCompletedOnSchedule: number;
  totalWorkProductsBehindSchedule: number;
}

export interface SubordinatesScoreCardResponse {
  isSuccess: boolean;
  message: string;
  managerId: string;
  reviewPeriodId: string;
  scoreCards: SubordinateScoreCard[];
}

// --- Period Score Data (for PMS Score Report) ---
export interface PeriodScoreData {
  periodScoreId: string;
  staffId: string;
  staffFullName: string;
  reviewPeriodId: string;
  reviewPeriod: string;
  year: number;
  startDate: string;
  endDate: string;
  finalScore: number;
  maxPoint: number;
  scorePercentage: number;
  finalGradeName: string;
  strategyName?: string;
  hrdDeductedPoints: number;
  departmentId?: number;
  departmentName?: string;
  divisionId?: number;
  divisionName?: string;
  officeId?: number;
  officeName?: string;
  staffGrade?: string;
  isUnderPerforming: boolean;
  minNoOfObjectives?: number;
  maxNoOfObjectives?: number;
  locationId?: string;
}

// --- Competency Group Report ---
export interface CompetencyRatingStat {
  ratingOrder: number;
  ratingName: string;
  numberOfStaff: number;
  staffPercentage: number;
}

export interface CategoryCompetencyDetailStat {
  categoryName: string;
  competencyRatingStat: CompetencyRatingStat[];
  averageRating: number;
  highestRating: number;
  lowestRating: number;
  mostCommonRating: number;
  groupCompetencyRatings: ChartDataVm[];
}

export interface CategoryCompetencyStat {
  categoryName: string;
  actual: number;
  expected: number;
}

export interface GroupedCompetencyReviewProfile {
  categoryCompetencyStats: CategoryCompetencyStat[];
  categoryCompetencyDetailStats: CategoryCompetencyDetailStat[];
}

export interface ChartDataVm {
  label: string;
  actual: number;
  expected: number;
}

// --- Competency Matrix ---
export interface CompetencyMatrixDetail {
  competencyName: string;
  averageScore: number;
  expectedRatingValue: number;
}

export interface CompetencyMatrixReviewProfile {
  employeeId: string;
  employeeName: string;
  position: string;
  grade: string;
  officeName: string;
  divisionName: string;
  departmentName: string;
  noOfCompetencies: number;
  noOfCompetent: number;
  gapCount: number;
  overallAverage: number;
  competencyMatrixDetails: CompetencyMatrixDetail[];
}

export interface CompetencyMatrixReviewOverview {
  competencyNames: string[];
  competencyMatrixReviewProfiles: CompetencyMatrixReviewProfile[];
}
