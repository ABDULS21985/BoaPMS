package enums

// Status represents the lifecycle state of a record.
type Status int

const (
	StatusDraft               Status = 1
	StatusPendingApproval     Status = 2
	StatusApprovedAndActive   Status = 3
	StatusReturned            Status = 4
	StatusRejected            Status = 5
	StatusAwaitingEvaluation  Status = 6
	StatusCompleted           Status = 7
	StatusPaused              Status = 8
	StatusSuspended           Status = StatusPaused
	StatusCancelled           Status = 9
	StatusBreached            Status = 10
	StatusDeactivated         Status = 11
	StatusAll                 Status = 12
	StatusClosed              Status = 13
	StatusPendingAcceptance   Status = 14
	StatusActive              Status = 15
	StatusPendingResolution   Status = 16
	StatusResolvedAwaitingFeedback    Status = 17
	StatusEscalated                   Status = 18
	StatusAwaitingRespondentComment   Status = 19
	StatusPendingHODReview            Status = 20
	StatusPendingBUHeadReview         Status = 21
	StatusPendingHRDReview            Status = 22
	StatusPendingHRDApproval          Status = 23
	StatusSuspensionPendingApproval   Status = 24
	StatusReEvaluate                  Status = 25
)

func (s Status) String() string {
	names := map[Status]string{
		StatusDraft: "Draft", StatusPendingApproval: "PendingApproval",
		StatusApprovedAndActive: "ApprovedAndActive", StatusReturned: "Returned",
		StatusRejected: "Rejected", StatusAwaitingEvaluation: "AwaitingEvaluation",
		StatusCompleted: "Completed", StatusPaused: "Paused",
		StatusCancelled: "Cancelled", StatusBreached: "Breached",
		StatusDeactivated: "Deactivated", StatusAll: "All",
		StatusClosed: "Closed", StatusPendingAcceptance: "PendingAcceptance",
		StatusActive: "Active", StatusPendingResolution: "PendingResolution",
		StatusResolvedAwaitingFeedback: "ResolvedAwaitingFeedback",
		StatusEscalated: "Escalated", StatusAwaitingRespondentComment: "AwaitingRespondentComment",
		StatusPendingHODReview: "PendingHODReview", StatusPendingBUHeadReview: "PendingBUHeadReview",
		StatusPendingHRDReview: "PendingHRDReview", StatusPendingHRDApproval: "PendingHRDApproval",
		StatusSuspensionPendingApproval: "SuspensionPendingApproval",
		StatusReEvaluate: "ReEvaluate",
	}
	if n, ok := names[s]; ok {
		return n
	}
	return "Unknown"
}

// ResolutionLevel represents grievance resolution levels.
type ResolutionLevel int

const (
	ResolutionLevelSBU        ResolutionLevel = 1
	ResolutionLevelDepartment ResolutionLevel = 2
	ResolutionLevelHRD        ResolutionLevel = 3
)

// ResolutionRemark represents the remark status on a resolution.
type ResolutionRemark int

const (
	ResolutionRemarkPending    ResolutionRemark = 1
	ResolutionRemarkAccepted   ResolutionRemark = 2
	ResolutionRemarkEscalated  ResolutionRemark = 3
	ResolutionRemarkReEvaluate ResolutionRemark = 4
	ResolutionRemarkClosed     ResolutionRemark = 5
)

// ReviewPeriodRange defines the range of a review period.
type ReviewPeriodRange int

const (
	ReviewPeriodRangeQuarterly ReviewPeriodRange = 1
	ReviewPeriodRangeBiAnnual  ReviewPeriodRange = 2
	ReviewPeriodRangeAnnual    ReviewPeriodRange = 3
)

// QuarterType identifies a fiscal quarter.
type QuarterType int

const (
	QuarterFirst  QuarterType = 1
	QuarterSecond QuarterType = 2
	QuarterThird  QuarterType = 3
	QuarterFourth QuarterType = 4
)

// PerformanceGrade categorizes performance levels.
type PerformanceGrade int

const (
	PerformanceGradeProbation    PerformanceGrade = 1
	PerformanceGradeDeveloping   PerformanceGrade = 2
	PerformanceGradeProgressive  PerformanceGrade = 3
	PerformanceGradeCompetent    PerformanceGrade = 4
	PerformanceGradeAccomplished PerformanceGrade = 5
	PerformanceGradeExemplary    PerformanceGrade = 6
)

func (g PerformanceGrade) String() string {
	names := map[PerformanceGrade]string{
		PerformanceGradeProbation: "Probation", PerformanceGradeDeveloping: "Developing",
		PerformanceGradeProgressive: "Progressive", PerformanceGradeCompetent: "Competent",
		PerformanceGradeAccomplished: "Accomplished", PerformanceGradeExemplary: "Exemplary",
	}
	if n, ok := names[g]; ok {
		return n
	}
	return "Unknown"
}

// FeedbackRequestType categorizes feedback request types.
type FeedbackRequestType int

const (
	FeedbackRequestWorkProductEvaluation      FeedbackRequestType = 1
	FeedbackRequestObjectivePlanning          FeedbackRequestType = 2
	FeedbackRequestProjectPlanning            FeedbackRequestType = 3
	FeedbackRequestCommitteePlanning          FeedbackRequestType = 4
	FeedbackRequestWorkProductFeedback        FeedbackRequestType = 5
	FeedbackRequest360ReviewFeedback          FeedbackRequestType = 6
	FeedbackRequestWorkProductPlanning        FeedbackRequestType = 7
	FeedbackRequestCompetencyReview           FeedbackRequestType = 8
	FeedbackRequestReviewPeriod               FeedbackRequestType = 9
	FeedbackRequestReviewPeriodExtension      FeedbackRequestType = 10
	FeedbackRequestProjectMemberAssignment    FeedbackRequestType = 11
	FeedbackRequestCommitteeMemberAssignment  FeedbackRequestType = 12
	FeedbackRequestPeriodObjectiveOutcome     FeedbackRequestType = 13
	FeedbackRequestDeptObjectiveOutcome       FeedbackRequestType = 14
	FeedbackRequestReviewPeriod360Review      FeedbackRequestType = 15
	FeedbackRequestProjectWorkProductDef      FeedbackRequestType = 16
	FeedbackRequestCommitteeWorkProductDef    FeedbackRequestType = 17
)

// GrievanceType categorizes grievance subjects.
type GrievanceType int

const (
	GrievanceTypeNone                  GrievanceType = 0
	GrievanceTypeWorkProductEvaluation GrievanceType = 1
	GrievanceTypeWorkProductAssignment GrievanceType = 2
	GrievanceTypeWorkProductPlanning   GrievanceType = 3
	GrievanceTypeObjectivePlanning     GrievanceType = 4
)

// WorkProductType classifies work products.
type WorkProductType int

const (
	WorkProductTypeOperational WorkProductType = 1
	WorkProductTypeProject     WorkProductType = 2
	WorkProductTypeCommittee   WorkProductType = 3
)

// EvaluationType classifies evaluation dimensions.
type EvaluationType int

const (
	EvaluationTypeTimeliness EvaluationType = 1
	EvaluationTypeQuality    EvaluationType = 2
	EvaluationTypeOutput     EvaluationType = 3
)

// ObjectiveLevel classifies organisational hierarchy for objectives.
type ObjectiveLevel int

const (
	ObjectiveLevelDepartment ObjectiveLevel = 1
	ObjectiveLevelDivision   ObjectiveLevel = 2
	ObjectiveLevelOffice     ObjectiveLevel = 3
	ObjectiveLevelEnterprise ObjectiveLevel = 4
)

// ObjectiveType distinguishes enterprise vs operational objectives.
type ObjectiveType int

const (
	ObjectiveTypeEnterprise   ObjectiveType = 1
	ObjectiveTypeOperational  ObjectiveType = 2
)

// OperationType represents CRUD and workflow operations.
type OperationType int

const (
	OperationAdd OperationType = iota
	OperationUpdate
	OperationDelete
	OperationDraft
	OperationCommitDraft
	OperationApprove
	OperationReject
	OperationCancel
	OperationComplete
	OperationPause
	OperationClose
	OperationReSubmit
	OperationReturn
	OperationAccept
	OperationReEvaluate
	OperationEnableObjectivePlanning
	OperationDisableObjectivePlanning
	OperationAddWithoutApproval
	OperationReInstate
	OperationDrop
	OperationResume
	OperationReactivate
	OperationEnableWorkProductPlanning
	OperationDisableWorkProductPlanning
	OperationEnableWorkProductEvaluation
	OperationDisableWorkProductEvaluation
	OperationSuspend
)

func (o OperationType) String() string {
	names := map[OperationType]string{
		OperationAdd:                       "Add",
		OperationUpdate:                    "Update",
		OperationDelete:                    "Delete",
		OperationDraft:                     "Draft",
		OperationCommitDraft:               "CommitDraft",
		OperationApprove:                   "Approve",
		OperationReject:                    "Reject",
		OperationCancel:                    "Cancel",
		OperationComplete:                  "Complete",
		OperationPause:                     "Pause",
		OperationClose:                     "Close",
		OperationReSubmit:                  "ReSubmit",
		OperationReturn:                    "Return",
		OperationAccept:                    "Accept",
		OperationReEvaluate:                "ReEvaluate",
		OperationEnableObjectivePlanning:   "EnableObjectivePlanning",
		OperationDisableObjectivePlanning:  "DisableObjectivePlanning",
		OperationAddWithoutApproval:        "AddWithoutApproval",
		OperationReInstate:                 "ReInstate",
		OperationDrop:                      "Drop",
		OperationResume:                    "Resume",
		OperationReactivate:                "Reactivate",
		OperationEnableWorkProductPlanning: "EnableWorkProductPlanning",
		OperationDisableWorkProductPlanning: "DisableWorkProductPlanning",
		OperationEnableWorkProductEvaluation:  "EnableWorkProductEvaluation",
		OperationDisableWorkProductEvaluation: "DisableWorkProductEvaluation",
		OperationSuspend: "Suspend",
	}
	if n, ok := names[o]; ok {
		return n
	}
	return "Unknown"
}

// ReviewPeriodExtensionTargetType scopes who an extension applies to.
type ReviewPeriodExtensionTargetType int

const (
	ExtensionTargetBankwide   ReviewPeriodExtensionTargetType = 1
	ExtensionTargetDepartment ReviewPeriodExtensionTargetType = 2
	ExtensionTargetDivision   ReviewPeriodExtensionTargetType = 3
	ExtensionTargetOffice     ReviewPeriodExtensionTargetType = 4
	ExtensionTargetStaff      ReviewPeriodExtensionTargetType = 5
)

// ReviewPeriod360TargetType scopes who a 360 review applies to.
type ReviewPeriod360TargetType int

const (
	Review360TargetBankwide   ReviewPeriod360TargetType = 1
	Review360TargetDepartment ReviewPeriod360TargetType = 2
	Review360TargetDivision   ReviewPeriod360TargetType = 3
	Review360TargetOffice     ReviewPeriod360TargetType = 4
	Review360TargetStaff      ReviewPeriod360TargetType = 5
)

// JobGradeGroupType classifies grade groups.
type JobGradeGroupType int

const (
	JobGradeGroupJunior    JobGradeGroupType = 1
	JobGradeGroupOfficer   JobGradeGroupType = 2
	JobGradeGroupManager   JobGradeGroupType = 3
	JobGradeGroupExecutive JobGradeGroupType = 4
)

// AuditEventType classifies audit log entries.
type AuditEventType int

const (
	AuditEventAdded    AuditEventType = 1
	AuditEventDeleted  AuditEventType = 2
	AuditEventModified AuditEventType = 3
)

// SuspensionAction defines what happens on suspension.
type SuspensionAction int

const (
	SuspensionActionResumeLater          SuspensionAction = 1
	SuspensionActionSubmitForEvaluation  SuspensionAction = 2
)

// CompetencyLevel represents the proficiency level of a competency.
type CompetencyLevel int

const (
	CompetencyLevelBeginner     CompetencyLevel = 1
	CompetencyLevelIntermediate CompetencyLevel = 3
	CompetencyLevelExpert       CompetencyLevel = 5
)

// CompetencyCategoryType classifies competency categories.
type CompetencyCategoryType int

const (
	CompetencyCategoryTechnical      CompetencyCategoryType = 1
	CompetencyCategoryOrganisational CompetencyCategoryType = 2
	CompetencyCategoryProfessional   CompetencyCategoryType = 3
	CompetencyCategoryLeadership     CompetencyCategoryType = 4
)

// DurationType represents time-based duration units.
type DurationType int

const (
	DurationTypeDay   DurationType = 1
	DurationTypeWeek  DurationType = 2
	DurationTypeMonth DurationType = 3
)

// Priority represents task/item priority levels.
type Priority int

const (
	PriorityLow    Priority = 1
	PriorityMedium Priority = 2
	PriorityHigh   Priority = 3
)

// ProjectStatus represents the lifecycle state of a project.
type ProjectStatus int

const (
	ProjectStatusOngoing     ProjectStatus = 1
	ProjectStatusCompleted   ProjectStatus = 2
	ProjectStatusTerminated  ProjectStatus = 3
)

// DevelopmentTaskStatus represents the lifecycle of a development task.
type DevelopmentTaskStatus int

const (
	DevelopmentTaskStatusAssigned   DevelopmentTaskStatus = 1
	DevelopmentTaskStatusInitiated  DevelopmentTaskStatus = 2
	DevelopmentTaskStatusInProgress DevelopmentTaskStatus = 3
	DevelopmentTaskStatusCompleted  DevelopmentTaskStatus = 4
	DevelopmentTaskStatusClosedGap  DevelopmentTaskStatus = 5
)

// CompetencyGroup classifies competency groups.
type CompetencyGroup int

const (
	CompetencyGroupTechnical   CompetencyGroup = 0
	CompetencyGroupBehavioral  CompetencyGroup = 1
)

// ReviewTypeName identifies the relationship of a reviewer.
type ReviewTypeName int

const (
	ReviewTypeNameSupervisor   ReviewTypeName = 1
	ReviewTypeNameSuperior     ReviewTypeName = 2
	ReviewTypeNamePeers        ReviewTypeName = 3
	ReviewTypeNameSubordinates ReviewTypeName = 4
	ReviewTypeNameSelf         ReviewTypeName = 5
)

// Approval represents approval/rejection outcomes.
type Approval int

const (
	ApprovalApproved Approval = 1
	ApprovalRejected Approval = 2
)

// SequenceNumberTypes identifies the entity type for sequence number generation.
type SequenceNumberTypes int

const (
	SeqCategoryDefinitions                   SequenceNumberTypes = 1
	SeqCategoryMapping                       SequenceNumberTypes = 2
	SeqPerformanceGrade                      SequenceNumberTypes = 3
	SeqReviewPeriod                          SequenceNumberTypes = 4
	SeqObjective                             SequenceNumberTypes = 5
	SeqObjectivePeriodMapping                SequenceNumberTypes = 6
	SeqWorkProduct                           SequenceNumberTypes = 7
	SeqWorkProductAssignment                 SequenceNumberTypes = 8
	SeqWorkProductEvaluation                 SequenceNumberTypes = 9
	SeqFeedbackRequest                       SequenceNumberTypes = 10
	SeqFeedbackResponse                      SequenceNumberTypes = 11
	SeqEnterprisePriority                    SequenceNumberTypes = 12
	SeqReviewPeriodExtension                 SequenceNumberTypes = 13
	SeqJobGradeGroup                         SequenceNumberTypes = 14
	SeqGrievance                             SequenceNumberTypes = 15
	SeqGrievanceComment                      SequenceNumberTypes = 16
	SeqProject                               SequenceNumberTypes = 17
	SeqProjectMemberAssignment               SequenceNumberTypes = 18
	SeqProjectMilestone                      SequenceNumberTypes = 19
	SeqProjectObjective                      SequenceNumberTypes = 20
	SeqCommittee                             SequenceNumberTypes = 21
	SeqCommitteeMemberAssignment             SequenceNumberTypes = 22
	SeqCommitteeObjective                    SequenceNumberTypes = 23
	SeqWorkProductDefinition                 SequenceNumberTypes = 24
	SeqCompetencyCategory                    SequenceNumberTypes = 25
	SeqCompetency                            SequenceNumberTypes = 26
	SeqCompetencyRequirement                 SequenceNumberTypes = 27
	SeqCompetencyAssessment                  SequenceNumberTypes = 28
	SeqCompetencyDevelopmentPlan             SequenceNumberTypes = 29
	SeqCompetencyDevelopmentTask             SequenceNumberTypes = 30
	SeqCompetencyDevelopmentTaskUpdate       SequenceNumberTypes = 31
	SeqObjectiveOutcomePeriodMapping         SequenceNumberTypes = 32
	SeqDeptObjective                         SequenceNumberTypes = 33
	SeqDeptObjectiveOutcomePeriodMapping     SequenceNumberTypes = 34
	SeqReviewPeriod360                       SequenceNumberTypes = 35
	SeqReviewPeriod360Feedback               SequenceNumberTypes = 36
	SeqReviewPeriod360FeedbackResponse       SequenceNumberTypes = 37
	SeqConfigItem                            SequenceNumberTypes = 38
	SeqGlobalSetting                         SequenceNumberTypes = 39
	SeqLineManagerGrading                    SequenceNumberTypes = 40
	SeqLineManagerObjectiveCategory          SequenceNumberTypes = 41
	SeqReviewPeriodStaffPerformance          SequenceNumberTypes = 42
	SeqNormalization                         SequenceNumberTypes = 43
	SeqOrganogramPeriodObjectiveOutcomeGrade SequenceNumberTypes = 44
	SeqProjectWorkProductDefinition          SequenceNumberTypes = 45
	SeqCommitteeWorkProductDefinition        SequenceNumberTypes = 46
	SeqSuspension                            SequenceNumberTypes = 47
	SeqAuditLog                              SequenceNumberTypes = 48
	SeqStrategy                              SequenceNumberTypes = 49
	SeqStrategicTheme                        SequenceNumberTypes = 50
	SeqEnterpriseObjective                   SequenceNumberTypes = 51
	SeqDivisionObjective                     SequenceNumberTypes = 52
	SeqOfficeObjective                       SequenceNumberTypes = 53
	SeqObjectiveCategory                     SequenceNumberTypes = 54
	SeqPmsCompetency                         SequenceNumberTypes = 55
	SeqFeedbackQuestionaire                  SequenceNumberTypes = 56
	SeqFeedbackQuestionaireOption            SequenceNumberTypes = 57
	SeqEvaluationOption                      SequenceNumberTypes = 58
	SeqDepartmentObjective                   SequenceNumberTypes = 59
)

// ConCat represents concatenation direction.
type ConCat int

const (
	ConCatBefore ConCat = 0
	ConCatAfter  ConCat = 1
)

// LineManagerPerformanceCategory classifies line-manager performance areas.
type LineManagerPerformanceCategory int

const (
	LMPerfCategoryObjectivePlanning       LineManagerPerformanceCategory = 1
	LMPerfCategoryWorkProductPlanning     LineManagerPerformanceCategory = 2
	LMPerfCategoryWorkProductEvaluation   LineManagerPerformanceCategory = 3
	LMPerfCategoryCompetencyAssessment    LineManagerPerformanceCategory = 4
	LMPerfCategoryCompetencyDevelopment   LineManagerPerformanceCategory = 5
	LMPerfCategoryObjectiveOutcome        LineManagerPerformanceCategory = 6
	LMPerfCategoryDeptObjectiveOutcome    LineManagerPerformanceCategory = 7
	LMPerfCategory360Review               LineManagerPerformanceCategory = 8
	LMPerfCategoryPerformanceGrading      LineManagerPerformanceCategory = 9
)

// OrganogramLevel represents levels in the organisational hierarchy.
type OrganogramLevel int

const (
	OrganogramLevelBankwide    OrganogramLevel = 1
	OrganogramLevelDepartment  OrganogramLevel = 2
	OrganogramLevelDivision    OrganogramLevel = 3
	OrganogramLevelOffice      OrganogramLevel = 4
	OrganogramLevelDirectorate OrganogramLevel = 5
)

// AdhocAssignmentType classifies ad-hoc assignment types.
type AdhocAssignmentType int

const (
	AdhocAssignmentCommittee AdhocAssignmentType = 1
	AdhocAssignmentProject   AdhocAssignmentType = 2
)

// BiAnnualType identifies bi-annual review periods.
type BiAnnualType int

const (
	BiAnnualFirst  BiAnnualType = 1
	BiAnnualSecond BiAnnualType = 2
)

// ---------------------------------------------------------------------------
// Scalar constants
// ---------------------------------------------------------------------------

// ActiveStaffType is the employee type code for active staff.
const ActiveStaffType = 1120

// ERP organisation type codes.
const (
	ErpOrgBranch         = "BRN"
	ErpOrgDepartment     = "DEPT"
	ErpOrgDeputyGovernor = "DGOV"
	ErpOrgDivision       = "DIV"
	ErpOrgOffice         = "OFF"
	ErpOrgUnit           = "UNI"
)

// Organisation IDs for well-known business units.
const (
	OrgIDSMD                    = 427
	OrgIDFND                    = 347
	OrgIDITD                    = 377
	OrgIDHRD                    = 355
	OrgIDRSD                    = 737
	OrgIDMPD                    = 904
	OrgIDLagosBranch            = 145
	OrgIDAbujaBranch            = 106
	OrgIDAbujaLocation          = 147
	OrgIDBudgetOfficeFND        = 2333
	OrgIDRewardsAndBenefitHRD   = 364
)

// EmployeeLocationHQ is the location ID for headquarters.
const EmployeeLocationHQ = 147
