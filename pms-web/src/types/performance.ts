import type {
  EvaluationType,
  ObjectiveLevel,
  ObjectiveType,
  PerformanceGrade,
  ReviewPeriodRange,
  WorkProductType,
} from "./enums";

// --- Strategy ---
export interface Strategy {
  strategyId: string;
  name: string;
  smdReferenceCode?: string;
  description?: string;
  bankYearId: number;
  startDate: string;
  endDate: string;
  fileImage?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface StrategicTheme {
  strategicThemeId: string;
  name: string;
  description?: string;
  strategyId: string;
  strategyName?: string;
  fileImage?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Objective Categories ---
export interface ObjectiveCategory {
  objectiveCategoryId: string;
  name: string;
  description?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface CategoryDefinition {
  definitionId: string;
  objectiveCategoryId: string;
  reviewPeriodId: string;
  weight: number;
  maxNoObjectives: number;
  maxNoWorkProduct: number;
  maxPoints: number;
  isCompulsory: boolean;
  enforceWorkProductLimit: boolean;
  description?: string;
  gradeGroupId: number;
  categoryName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Objectives ---
export interface EnterpriseObjective {
  enterpriseObjectiveId: string;
  type: ObjectiveType;
  enterpriseObjectivesCategoryId: string;
  strategicThemeId?: string;
  strategyId: string;
  name: string;
  description?: string;
  kpi?: string;
  target?: string;
  categoryName?: string;
  strategyName?: string;
  themeName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface DepartmentObjective {
  departmentObjectiveId: string;
  departmentId: number;
  enterpriseObjectiveId: string;
  name: string;
  description?: string;
  kpi?: string;
  target?: string;
  departmentName?: string;
  enterpriseObjectiveName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface DivisionObjective {
  divisionObjectiveId: string;
  divisionId: number;
  departmentObjectiveId: string;
  name: string;
  description?: string;
  kpi?: string;
  target?: string;
  divisionName?: string;
  departmentObjectiveName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface OfficeObjective {
  officeObjectiveId: string;
  officeId: number;
  divisionObjectiveId: string;
  jobGradeGroupId: number;
  name: string;
  description?: string;
  kpi?: string;
  target?: string;
  officeName?: string;
  divisionObjectiveName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Consolidated Objective ---
export interface ConsolidatedObjective {
  objectiveId: string;
  name: string;
  description?: string;
  objectiveLevel: string;
  kpi?: string;
  target?: string;
  sbuName?: string;
  smdReferenceCode?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
  isSelected?: boolean;
}

// --- Review Period ---
export interface PerformanceReviewPeriod {
  periodId: string;
  year: number;
  range: ReviewPeriodRange;
  rangeValue: number;
  name: string;
  description?: string;
  shortName?: string;
  startDate: string;
  endDate: string;
  allowObjectivePlanning: boolean;
  allowWorkProductPlanning: boolean;
  allowWorkProductEvaluation: boolean;
  maxPoints: number;
  minNoOfObjectives: number;
  maxNoOfObjectives: number;
  strategyId?: string;
  strategyName?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface ReviewPeriodExtension {
  reviewPeriodExtensionId: string;
  reviewPeriodId: string;
  targetType: number;
  targetReference?: string;
  description?: string;
  startDate: string;
  endDate: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Work Product ---
export interface WorkProduct {
  workProductId: string;
  name: string;
  description?: string;
  maxPoint: number;
  workProductType: WorkProductType;
  isSelfCreated: boolean;
  staffId: string;
  acceptanceComment?: string;
  startDate: string;
  endDate: string;
  deliverables?: string;
  finalScore: number;
  noReturned: number;
  completionDate?: string;
  approverComment?: string;
  reEvaluationReInitiated: boolean;
  remark?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface WorkProductTask {
  workProductTaskId: string;
  name: string;
  description?: string;
  startDate: string;
  endDate: string;
  completionDate?: string;
  workProductId: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface WorkProductEvaluation {
  workProductEvaluationId: string;
  workProductId: string;
  timeliness: number;
  timelinessEvaluationOptionId?: string;
  quality: number;
  qualityEvaluationOptionId?: string;
  output: number;
  outputEvaluationOptionId?: string;
  outcome: number;
  evaluatorStaffId?: string;
  isReEvaluated: boolean;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface EvaluationOption {
  evaluationOptionId: string;
  name: string;
  description?: string;
  recordStatus: number;
  score: number;
  evaluationType: EvaluationType;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface WorkProductDefinition {
  workProductDefinitionId: string;
  referenceNo?: string;
  name: string;
  description?: string;
  deliverables?: string;
  objectiveId?: string;
  objectiveLevel: ObjectiveLevel | string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Project & Committee ---
export interface Project {
  projectId: string;
  name: string;
  description?: string;
  startDate: string;
  endDate: string;
  deliverables?: string;
  reviewPeriodId: string;
  departmentId: number;
  projectManager?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Committee {
  committeeId: string;
  name: string;
  description?: string;
  startDate: string;
  endDate: string;
  deliverables?: string;
  reviewPeriodId: string;
  departmentId: number;
  chairperson?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Feedback ---
export interface FeedbackRequestLog {
  feedbackRequestLogId: string;
  feedbackRequestType: number;
  referenceId: string;
  timeInitiated: string;
  assignedStaffId: string;
  assignedStaffName?: string;
  requestOwnerStaffId: string;
  requestOwnerStaffName?: string;
  timeCompleted?: string;
  requestOwnerComment?: string;
  assignedStaffComment?: string;
  hasSla: boolean;
  reviewPeriodId?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface FeedbackQuestionnaire {
  feedbackQuestionaireId: string;
  question: string;
  description: string;
  pmsCompetencyId: string;
  pmsCompetencyName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface FeedbackQuestionnaireOption {
  feedbackQuestionaireOptionId: string;
  optionStatement: string;
  description?: string;
  score: number;
  questionId: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- PMS Competency ---
export interface PmsCompetency {
  pmsCompetencyId: string;
  name: string;
  description?: string;
  objectCategoryId: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
  feedbackQuestionaires?: FeedbackQuestionnaireWithOptions[];
}

export interface FeedbackQuestionnaireWithOptions {
  feedbackQuestionaireId: string;
  question: string;
  description: string;
  pmsCompetencyId: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
  options: FeedbackQuestionnaireOption[];
}

// --- Competency Review (360 Feedback) ---
export interface CompetencyReviewFeedback {
  competencyReviewFeedbackId: string;
  staffId: string;
  staffName?: string;
  maxPoints: number;
  finalScore: number;
  reviewPeriodId: string;
  recordStatus?: number;
  recordStatusName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface CompetencyReviewFeedbackDetails {
  competencyReviewFeedbackId: string;
  staffId: string;
  staffName?: string;
  officeId: number;
  officeCode?: string;
  officeName?: string;
  divisionId: number;
  divisionCode?: string;
  divisionName?: string;
  departmentId: number;
  departmentCode?: string;
  departmentName?: string;
  maxPoints: number;
  finalScore: number;
  finalScorePercentage: number;
  reviewPeriodId: string;
  recordStatusName?: string;
  ratings?: CompetencyReviewerRatingSummary[];
}

export interface CompetencyReviewerRatingSummary {
  pmsCompetencyId: string;
  pmsCompetencyName: string;
  averageRating: number;
  totalReviewers: number;
}

export interface CompetencyReviewer {
  competencyReviewerId: string;
  reviewStaffId: string;
  finalRating: number;
  competencyReviewFeedbackId: string;
  competencyReviewFeedback?: CompetencyReviewFeedback;
  recordStatusName?: string;
  initiatedDate?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
  competencyReviewerRatings?: CompetencyReviewerRating[];
}

export interface CompetencyReviewerRating {
  competencyReviewerRatingId: string;
  pmsCompetencyId: string;
  feedbackQuestionaireOptionId: string;
  rating: number;
  competencyReviewerId: string;
  feedbackQuestionaireOption?: FeedbackQuestionnaireOption;
  pmsCompetency?: PmsCompetency;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Breached Feedback Request ---
export interface BreachedFeedbackRequestLog extends FeedbackRequestLog {
  isBreached: boolean;
}

// --- Pending Request ---
export interface StaffPendingRequest {
  feedbackRequestLogId: string;
  feedbackRequestType: number;
  referenceId: string;
  timeInitiated: string;
  assignedStaffId: string;
  requestOwnerStaffId: string;
  recordStatus: number;
  hasSla: boolean;
  id: number;
  isActive: boolean;
}

// --- Grievance ---
export interface Grievance {
  grievanceId: string;
  grievanceType: number;
  reviewPeriodId: string;
  subjectId: string;
  subject?: string;
  description: string;
  respondentComment?: string;
  currentResolutionLevel: number;
  complainantStaffId: string;
  respondentStaffId: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Period Score ---
export interface PeriodScore {
  periodScoreId: string;
  reviewPeriodId: string;
  staffId: string;
  finalScore: number;
  scorePercentage: number;
  finalGrade: PerformanceGrade;
  endDate: string;
  officeId: number;
  staffGrade?: string;
  hrdDeductedPoints: number;
  isUnderPerforming: boolean;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Individual Planned Objective ---
export interface IndividualPlannedObjective {
  individualPlannedObjectiveId: string;
  reviewPeriodId: string;
  staffId: string;
  objectiveId?: string;
  objectiveName?: string;
  objectiveLevel?: string;
  categoryName?: string;
  title?: string;
  description?: string;
  weight: number;
  targetDate?: string;
  keyPerformanceIndicator?: string;
  remark?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Project & Committee Members ---
export interface ProjectMember {
  projectMemberId: string;
  staffId: string;
  staffName?: string;
  projectId: string;
  plannedObjectiveId?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface CommitteeMember {
  committeeMemberId: string;
  staffId: string;
  staffName?: string;
  committeeId: string;
  plannedObjectiveId?: string;
  recordStatus?: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- PMS Configuration ---
export interface PmsConfiguration {
  pmsConfigurationId: string;
  name: string;
  value?: string;
  type?: string;
  isEncrypted: boolean;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Setting {
  settingId: string;
  name: string;
  value?: string;
  type?: string;
  isEncrypted: boolean;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}
