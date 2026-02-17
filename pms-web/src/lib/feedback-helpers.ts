import { FeedbackRequestType } from "@/types/enums";

export const feedbackRequestTypeLabels: Record<number, string> = {
  [FeedbackRequestType.WorkProductEvaluation]: "Work Product Evaluation",
  [FeedbackRequestType.ObjectivePlanning]: "Objective Planning",
  [FeedbackRequestType.WorkProductPlanning]: "Work Product Planning",
  [FeedbackRequestType.ProjectPlanning]: "Project Planning",
  [FeedbackRequestType.CommitteePlanning]: "Committee Planning",
  [FeedbackRequestType.WorkProductCancellation]: "Work Product Cancellation",
  [FeedbackRequestType.WorkProductSuspension]: "Work Product Suspension",
  [FeedbackRequestType.ObjectiveCancellation]: "Objective Cancellation",
  [FeedbackRequestType.ProjectMemberRemoval]: "Project Member Removal",
  [FeedbackRequestType.CommitteeMemberRemoval]: "Committee Member Removal",
  [FeedbackRequestType.ObjectiveSuspension]: "Objective Suspension",
  [FeedbackRequestType.ObjectiveReInstatement]: "Objective Reinstatement",
  [FeedbackRequestType.WorkProductReInstatement]: "Work Product Reinstatement",
  [FeedbackRequestType.WorkProductResumption]: "Work Product Resumption",
  [FeedbackRequestType.ObjectiveResumption]: "Objective Resumption",
  [FeedbackRequestType.WorkProductReEvaluation]: "Work Product Re-Evaluation",
  [FeedbackRequestType.ReviewPeriodExtension]: "Review Period Extension",
};

export function getFeedbackRequestTypeLabel(type: number): string {
  return feedbackRequestTypeLabels[type] ?? `Type ${type}`;
}

export function getSlaIndicatorColor(
  hasSla: boolean,
  isBreached?: boolean,
  timeInitiated?: string,
  timeCompleted?: string
): "green" | "yellow" | "red" | "none" {
  if (!hasSla) return "none";
  if (isBreached) return "red";
  if (timeCompleted) return "green";

  if (!timeInitiated) return "none";
  const initiated = new Date(timeInitiated);
  const now = new Date();
  const hoursElapsed = (now.getTime() - initiated.getTime()) / (1000 * 60 * 60);

  // > 48h without completion → yellow warning
  if (hoursElapsed > 48) return "yellow";
  return "green";
}

export function formatOverdueDuration(timeInitiated: string): string {
  const initiated = new Date(timeInitiated);
  const now = new Date();
  const diffMs = now.getTime() - initiated.getTime();
  const days = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  const hours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));

  if (days > 0) return `${days}d ${hours}h overdue`;
  return `${hours}h overdue`;
}
