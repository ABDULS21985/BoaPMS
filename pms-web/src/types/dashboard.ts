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
  departmentName?: string;
  divisionName?: string;
  officeName?: string;
  supervisorId?: string;
  supervisorName?: string;
  headOfOfficeId?: string;
  headOfOfficeName?: string;
  phone?: string;
  photoUrl?: string;
}
