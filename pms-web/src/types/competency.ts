// --- Competency ---
export interface Competency {
  competencyId: number;
  competencyCategoryId: number;
  competencyName: string;
  description?: string;
  competencyCategoryName?: string;
  approvedBy?: string;
  dateApproved?: string;
  isApproved: boolean;
  isRejected: boolean;
  rejectedBy?: string;
  rejectionReason?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Competency Category ---
export interface CompetencyCategory {
  competencyCategoryId: number;
  categoryName: string;
  isTechnical: boolean;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Competency Category Grading ---
export interface CompetencyCategoryGrading {
  competencyCategoryGradingId: number;
  competencyCategoryId: number;
  reviewTypeId: number;
  weightPercentage: number;
  competencyCategoryName?: string;
  reviewTypeName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Competency Rating Definition ---
export interface CompetencyRatingDefinition {
  competencyRatingDefinitionId: number;
  competencyId: number;
  ratingId: number;
  definition: string;
  ratingName?: string;
  ratingValue?: number;
  competencyName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Competency Review ---
export interface CompetencyReview {
  competencyReviewId: number;
  reviewTypeId: number;
  reviewPeriodId: number;
  ratingId?: number;
  competencyId: number;
  expectedRatingId: number;
  reviewDate?: string;
  reviewerId: string;
  reviewerName?: string;
  employeeNumber: string;
  reviewTypeName?: string;
  reviewPeriodName?: string;
  competencyName?: string;
  competencyCategoryName?: string;
  competencyDefinition?: string;
  isTechnical: boolean;
  employeeName?: string;
  employeeInitial?: string;
  employeeGrade?: string;
  employeeDepartment?: string;
  actualRatingId: number;
  actualRatingName?: string;
  actualRatingValue: number;
  expectedRatingName?: string;
  expectedRatingValue: number;
  competencyRatingDefinitions?: CompetencyRatingDefinition[];
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Competency Review Profile ---
export interface CompetencyReviewProfile {
  competencyReviewProfileId: number;
  reviewPeriodId: number;
  reviewPeriodName?: string;
  averageRatingId?: number;
  averageRatingName?: string;
  averageRatingValue: number;
  expectedRatingId?: number;
  expectedRatingName?: string;
  expectedRatingValue: number;
  averageScore: number;
  employeeNumber: string;
  employeeFullName?: string;
  competencyId: number;
  competencyName?: string;
  competencyCategory?: number;
  competencyCategoryName?: string;
  numberOfDevelopmentPlans: number;
  progressCount: number;
  completedCount: number;
  officeId?: string;
  officeName?: string;
  divisionId?: string;
  divisionName?: string;
  departmentId?: string;
  departmentName?: string;
  jobRoleId?: string;
  jobRoleName?: string;
  gradeName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Development Plan ---
export interface DevelopmentPlan {
  developmentPlanId: number;
  competencyReviewProfileId: number;
  trainingTypeName?: string;
  activity: string;
  employeeNumber: string;
  targetDate: string;
  completionDate?: string;
  taskStatus?: string;
  learningResource?: string;
  competencyName?: string;
  competencyCategoryName?: string;
  reviewPeriod?: string;
  currentGap?: number;
  employeeName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Job Role ---
export interface JobRole {
  jobRoleId: number;
  jobRoleName: string;
  description?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Job Grade ---
export interface JobGrade {
  jobGradeId: number;
  gradeCode: string;
  gradeName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Job Grade Group ---
export interface JobGradeGroup {
  jobGradeGroupId: number;
  groupName: string;
  order: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Assign Job Grade Group ---
export interface AssignJobGradeGroup {
  assignJobGradeGroupId: number;
  jobGradeGroupId: number;
  jobGradeId: number;
  jobGradeName?: string;
  jobGradeGroupName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Job Role Competency ---
export interface JobRoleCompetency {
  jobRoleCompetencyId: number;
  officeId: number;
  jobRoleId: number;
  competencyId: number;
  ratingId: number;
  competencyName?: string;
  jobRoleName?: string;
  ratingName?: string;
  officeName?: string;
  departmentName?: string;
  divisionName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Job Role Grade ---
export interface JobRoleGrade {
  jobRoleGradeId: number;
  jobRoleId: number;
  gradeId?: string;
  gradeName?: string;
  jobRoleName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Behavioral Competency ---
export interface BehavioralCompetency {
  behavioralCompetencyId: number;
  competencyId: number;
  jobGradeGroupId: number;
  ratingId?: number;
  competencyName?: string;
  ratingName?: string;
  jobGradeGroupName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Office Job Role ---
export interface OfficeJobRole {
  officeJobRoleId: number;
  officeId: number;
  jobRoleId: number;
  officeName?: string;
  jobRoleName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Rating ---
export interface Rating {
  ratingId: number;
  name: string;
  value: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Review Period (Competency) ---
export interface CompetencyReviewPeriod {
  reviewPeriodId: number;
  bankYearId: number;
  name: string;
  startDate: string;
  endDate: string;
  bankYearName?: string;
  approvedBy?: string;
  dateApproved?: string;
  isApproved: boolean;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Review Type ---
export interface ReviewType {
  reviewTypeId: number;
  reviewTypeName: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Training Type ---
export interface TrainingType {
  trainingTypeId: number;
  trainingTypeName: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Bank Year ---
export interface BankYear {
  bankYearId: number;
  yearName: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Staff Job Role ---
export interface StaffJobRole {
  staffJobRoleId: number;
  employeeId: string;
  fullName?: string;
  departmentId?: number;
  divisionId?: number;
  officeId?: number;
  supervisorId?: string;
  jobRoleId: number;
  jobRoleName?: string;
  soaStatus: boolean;
  soaResponse?: string;
  isApproved: boolean;
  rejectionReason?: string;
  status?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

// --- Office Competency Review (Manager Overview) ---
export interface OfficeCompetencyReview {
  employeeNumber: string;
  employeeName?: string;
  gradeName?: string;
  jobRoleName?: string;
  departmentName?: string;
  isCompleted: boolean;
  officeId?: number;
}
