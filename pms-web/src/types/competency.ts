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
  rejectionReason?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface CompetencyCategory {
  competencyCategoryId: number;
  categoryName: string;
  isTechnical: boolean;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface CompetencyCategoryGrading {
  competencyCategoryGradingId: number;
  competencyCategoryId: number;
  reviewTypeId: number;
  weightPercentage: number;
  categoryName?: string;
  reviewTypeName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface CompetencyReview {
  competencyReviewId: number;
  reviewTypeId: number;
  reviewPeriodId: number;
  competencyId: number;
  expectedRatingId: number;
  reviewDate?: string;
  reviewerId: string;
  reviewerName: string;
  employeeNumber: string;
  reviewTypeName: string;
  reviewPeriodName: string;
  competencyName: string;
  competencyCategoryName: string;
  isTechnical: boolean;
  employeeName: string;
  employeeGrade: string;
  employeeDepartment: string;
  actualRatingId: number;
  actualRatingName: string;
  actualRatingValue: number;
  expectedRatingName: string;
  expectedRatingValue: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface CompetencyReviewProfile {
  competencyReviewProfileId: number;
  reviewPeriodId: number;
  reviewPeriodName: string;
  averageRatingValue: number;
  expectedRatingValue: number;
  averageScore: number;
  employeeNumber: string;
  employeeName: string;
  competencyId: number;
  competencyName: string;
  competencyCategoryName: string;
  competencyGap: number;
  haveGap: boolean;
  officeName: string;
  divisionName: string;
  departmentName: string;
  jobRoleName: string;
  gradeName: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

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

export interface JobRole {
  jobRoleId: number;
  jobRoleName: string;
  description?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface JobGrade {
  jobGradeId: number;
  gradeCode: string;
  gradeName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface JobGradeGroup {
  jobGradeGroupId: number;
  groupName: string;
  order: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface JobRoleCompetency {
  jobRoleCompetencyId: number;
  officeId: number;
  jobRoleId: number;
  competencyId: number;
  ratingId: number;
  competencyName: string;
  jobRoleName: string;
  ratingName: string;
  officeName: string;
  departmentName: string;
  divisionName: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface BehavioralCompetency {
  behavioralCompetencyId: number;
  competencyId: number;
  jobGradeGroupId: number;
  ratingId: number;
  competencyName?: string;
  ratingName?: string;
  gradeGroupName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface Rating {
  ratingId: number;
  name: string;
  value: number;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface ReviewPeriod {
  reviewPeriodId: number;
  bankYearId: number;
  name: string;
  startDate: string;
  endDate: string;
  bankYearName?: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface ReviewType {
  reviewTypeId: number;
  reviewTypeName: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

export interface TrainingType {
  trainingTypeId: number;
  trainingTypeName: string;
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

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
  isActive: boolean;
  createdBy: string;
  dateCreated: string;
}

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
