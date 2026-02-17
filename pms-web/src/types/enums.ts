export enum Status {
  Draft = 1,
  PendingApproval = 2,
  ApprovedAndActive = 3,
  Rejected = 4,
  Returned = 5,
  Cancelled = 6,
  Closed = 7,
  Active = 8,
  Inactive = 9,
  Completed = 10,
  PendingAcceptance = 11,
  Accepted = 12,
  Paused = 13,
  Suspended = 14,
  SuspensionPendingApproval = 15,
  Resumed = 16,
  Resubmitted = 17,
  ReEvaluationInitiated = 18,
  PendingReEvaluation = 19,
  Archived = 20,
  InProgress = 21,
  Overdue = 22,
  PendingClosure = 23,
  ClosedAndInactive = 24,
  Deleted = 25,
}

export enum ReviewPeriodRange {
  Quarterly = 1,
  BiAnnual = 2,
  Annual = 3,
}

export enum QuarterType {
  First = 1,
  Second = 2,
  Third = 3,
  Fourth = 4,
}

export enum PerformanceGrade {
  Probation = 0,
  Developing = 1,
  Progressive = 2,
  Competent = 3,
  Accomplished = 4,
  Exemplary = 5,
}

export enum FeedbackRequestType {
  WorkProductEvaluation = 1,
  ObjectivePlanning = 2,
  WorkProductPlanning = 3,
  ProjectPlanning = 4,
  CommitteePlanning = 5,
  WorkProductCancellation = 6,
  WorkProductSuspension = 7,
  ObjectiveCancellation = 8,
  ProjectMemberRemoval = 9,
  CommitteeMemberRemoval = 10,
  ObjectiveSuspension = 11,
  ObjectiveReInstatement = 12,
  WorkProductReInstatement = 13,
  WorkProductResumption = 14,
  ObjectiveResumption = 15,
  WorkProductReEvaluation = 16,
  ReviewPeriodExtension = 17,
}

export enum GrievanceType {
  None = 0,
  WorkProductEvaluation = 1,
  WorkProductAssignment = 2,
  ObjectivePlanning = 3,
  General = 4,
}

export enum WorkProductType {
  Operational = 1,
  Project = 2,
  Committee = 3,
}

export enum EvaluationType {
  Timeliness = 1,
  Quality = 2,
  Output = 3,
}

export enum ObjectiveLevel {
  Department = 1,
  Division = 2,
  Office = 3,
  Enterprise = 4,
}

export enum ObjectiveType {
  Enterprise = 1,
  Operational = 2,
}

export enum OperationType {
  Add = 1,
  Update = 2,
  Delete = 3,
  Draft = 4,
  CommitDraft = 5,
  Approve = 6,
  Reject = 7,
  Return = 8,
  Cancel = 9,
  Close = 10,
  Accept = 11,
  ReInstate = 12,
  Pause = 13,
  Suspend = 14,
  Resume = 15,
  ReSubmit = 16,
  Submit = 17,
  Complete = 18,
  Activate = 19,
  Deactivate = 20,
  Enable = 21,
  Disable = 22,
  ReEvaluate = 23,
}

export enum ReviewPeriodExtensionTargetType {
  Bankwide = 1,
  Department = 2,
  Division = 3,
  Office = 4,
  Staff = 5,
}

export enum LineManagerPerformanceCategory {
  Objectives = 1,
  WorkProducts = 2,
  Projects = 3,
  Committees = 4,
  Feedbacks = 5,
  ScoreCards = 6,
  Grievances = 7,
  Evaluations = 8,
  Requests = 9,
}

export enum OrganogramLevel {
  Bankwide = 0,
  Department = 1,
  Division = 2,
  Office = 3,
  Directorate = 4,
}

export enum AuditEventType {
  Added = 1,
  Deleted = 2,
  Modified = 3,
}

export enum CompetencyGroup {
  Technical = 1,
  Behavioral = 2,
}

export enum ReviewTypeName {
  Supervisor = 1,
  Superior = 2,
  Peers = 3,
  Subordinates = 4,
  Self = 5,
}

export enum ProjectStatus {
  Ongoing = 1,
  Completed = 2,
  Terminated = 3,
}

export enum DevelopmentTaskStatus {
  Assigned = 1,
  Initiated = 2,
  InProgress = 3,
  Completed = 4,
  ClosedGap = 5,
}

// Status display helpers
export const statusLabels: Record<number, string> = {
  [Status.Draft]: "Draft",
  [Status.PendingApproval]: "Pending Approval",
  [Status.ApprovedAndActive]: "Approved",
  [Status.Rejected]: "Rejected",
  [Status.Returned]: "Returned",
  [Status.Cancelled]: "Cancelled",
  [Status.Closed]: "Closed",
  [Status.Active]: "Active",
  [Status.Inactive]: "Inactive",
  [Status.Completed]: "Completed",
  [Status.PendingAcceptance]: "Pending Acceptance",
  [Status.Accepted]: "Accepted",
  [Status.Paused]: "Paused",
  [Status.Suspended]: "Suspended",
  [Status.SuspensionPendingApproval]: "Suspension Pending",
  [Status.Resumed]: "Resumed",
  [Status.Resubmitted]: "Resubmitted",
  [Status.ReEvaluationInitiated]: "Re-Evaluation Initiated",
  [Status.PendingReEvaluation]: "Pending Re-Evaluation",
  [Status.Archived]: "Archived",
  [Status.InProgress]: "In Progress",
  [Status.Overdue]: "Overdue",
  [Status.PendingClosure]: "Pending Closure",
  [Status.ClosedAndInactive]: "Closed",
  [Status.Deleted]: "Deleted",
};

export type StatusVariant =
  | "default"
  | "secondary"
  | "destructive"
  | "outline";

export function getStatusVariant(status: Status): StatusVariant {
  switch (status) {
    case Status.Draft:
      return "secondary";
    case Status.PendingApproval:
    case Status.PendingAcceptance:
    case Status.SuspensionPendingApproval:
    case Status.PendingReEvaluation:
    case Status.PendingClosure:
    case Status.Resubmitted:
      return "outline";
    case Status.ApprovedAndActive:
    case Status.Active:
    case Status.Completed:
    case Status.Accepted:
    case Status.Resumed:
      return "default";
    case Status.Rejected:
    case Status.Cancelled:
    case Status.Deleted:
    case Status.Overdue:
      return "destructive";
    default:
      return "secondary";
  }
}
