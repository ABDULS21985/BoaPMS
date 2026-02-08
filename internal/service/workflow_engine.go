package service

import (
	"fmt"
	"time"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/rs/zerolog"
)

// ---------------------------------------------------------------------------
// Workflow state machine — generic engine for managing approval lifecycles.
//

// ---------------------------------------------------------------------------

// TransitionRule defines a single valid state transition and the operation
// that triggers it. The engine only permits transitions that appear in its
// rule table.
type TransitionRule struct {
	From      enums.Status
	To        enums.Status
	Operation enums.OperationType
}

// TransitionHook is a callback invoked before or after a transition executes.
// Returning a non-nil error from a before-hook aborts the transition.
type TransitionHook func(entityID string, from, to enums.Status, actorID string) error

// WorkflowEngine is a finite-state machine that validates and executes
// workflow transitions for PMS entities.
type WorkflowEngine struct {
	transitions []TransitionRule
	onBefore    TransitionHook
	onAfter     TransitionHook
	log         zerolog.Logger
}

// ---------------------------------------------------------------------------
// Engine constructors
// ---------------------------------------------------------------------------

// NewBaseWorkflowEngine returns an engine configured for single-level
// (line-manager) approval workflows used by most PMS entities.
func NewBaseWorkflowEngine(log zerolog.Logger) *WorkflowEngine {
	return &WorkflowEngine{
		log: log,
		transitions: []TransitionRule{
			// Submission
			{From: enums.StatusDraft, To: enums.StatusPendingApproval, Operation: enums.OperationCommitDraft},
			{From: enums.StatusDraft, To: enums.StatusPendingApproval, Operation: enums.OperationAdd},

			// Single-level approval
			{From: enums.StatusPendingApproval, To: enums.StatusApprovedAndActive, Operation: enums.OperationApprove},
			{From: enums.StatusPendingApproval, To: enums.StatusRejected, Operation: enums.OperationReject},
			{From: enums.StatusPendingApproval, To: enums.StatusReturned, Operation: enums.OperationReturn},

			// Re-submission after return
			{From: enums.StatusReturned, To: enums.StatusPendingApproval, Operation: enums.OperationReSubmit},

			// Deactivation and reactivation
			{From: enums.StatusApprovedAndActive, To: enums.StatusDeactivated, Operation: enums.OperationCancel},
			{From: enums.StatusApprovedAndActive, To: enums.StatusDeactivated, Operation: enums.OperationDelete},
			{From: enums.StatusDeactivated, To: enums.StatusApprovedAndActive, Operation: enums.OperationReactivate},

			// Closure
			{From: enums.StatusApprovedAndActive, To: enums.StatusClosed, Operation: enums.OperationClose},
			{From: enums.StatusActive, To: enums.StatusClosed, Operation: enums.OperationClose},

			// Completion
			{From: enums.StatusActive, To: enums.StatusCompleted, Operation: enums.OperationComplete},
		},
	}
}

// NewHrdWorkflowEngine returns an engine configured for two-level approval
// workflows where both a line manager and HRD must approve.
func NewHrdWorkflowEngine(log zerolog.Logger) *WorkflowEngine {
	return &WorkflowEngine{
		log: log,
		transitions: []TransitionRule{
			// Submission
			{From: enums.StatusDraft, To: enums.StatusPendingApproval, Operation: enums.OperationCommitDraft},
			{From: enums.StatusDraft, To: enums.StatusPendingApproval, Operation: enums.OperationAdd},

			// Line-manager approval escalates to HRD
			{From: enums.StatusPendingApproval, To: enums.StatusPendingHRDApproval, Operation: enums.OperationApprove},
			{From: enums.StatusPendingApproval, To: enums.StatusRejected, Operation: enums.OperationReject},
			{From: enums.StatusPendingApproval, To: enums.StatusReturned, Operation: enums.OperationReturn},

			// HRD-level approval
			{From: enums.StatusPendingHRDApproval, To: enums.StatusApprovedAndActive, Operation: enums.OperationApprove},
			{From: enums.StatusPendingHRDApproval, To: enums.StatusRejected, Operation: enums.OperationReject},

			// Re-submission after return
			{From: enums.StatusReturned, To: enums.StatusPendingApproval, Operation: enums.OperationReSubmit},

			// Deactivation and reactivation
			{From: enums.StatusApprovedAndActive, To: enums.StatusDeactivated, Operation: enums.OperationCancel},
			{From: enums.StatusApprovedAndActive, To: enums.StatusDeactivated, Operation: enums.OperationDelete},
			{From: enums.StatusDeactivated, To: enums.StatusApprovedAndActive, Operation: enums.OperationReactivate},

			// Closure
			{From: enums.StatusApprovedAndActive, To: enums.StatusClosed, Operation: enums.OperationClose},
			{From: enums.StatusActive, To: enums.StatusClosed, Operation: enums.OperationClose},

			// Completion
			{From: enums.StatusActive, To: enums.StatusCompleted, Operation: enums.OperationComplete},
		},
	}
}

// NewReviewPeriodWorkflowEngine returns an engine configured for review period
// lifecycle transitions. Review periods have a simpler deactivation model but
// allow re-submission from a rejected state (unlike most entities).
func NewReviewPeriodWorkflowEngine(log zerolog.Logger) *WorkflowEngine {
	return &WorkflowEngine{
		log: log,
		transitions: []TransitionRule{
			// Submission
			{From: enums.StatusDraft, To: enums.StatusPendingApproval, Operation: enums.OperationCommitDraft},
			{From: enums.StatusDraft, To: enums.StatusPendingApproval, Operation: enums.OperationAdd},

			// Approval
			{From: enums.StatusPendingApproval, To: enums.StatusApprovedAndActive, Operation: enums.OperationApprove},
			{From: enums.StatusPendingApproval, To: enums.StatusRejected, Operation: enums.OperationReject},
			{From: enums.StatusPendingApproval, To: enums.StatusReturned, Operation: enums.OperationReturn},

			// Re-submission — review periods can be re-submitted from both Returned and Rejected
			{From: enums.StatusReturned, To: enums.StatusPendingApproval, Operation: enums.OperationReSubmit},
			{From: enums.StatusRejected, To: enums.StatusPendingApproval, Operation: enums.OperationReSubmit},

			// Closure
			{From: enums.StatusApprovedAndActive, To: enums.StatusClosed, Operation: enums.OperationClose},

			// Active → ApprovedAndActive is used when an already-active period is formally approved
			{From: enums.StatusActive, To: enums.StatusApprovedAndActive, Operation: enums.OperationApprove},

			// Cancellation
			{From: enums.StatusApprovedAndActive, To: enums.StatusCancelled, Operation: enums.OperationCancel},
			{From: enums.StatusActive, To: enums.StatusCancelled, Operation: enums.OperationCancel},
		},
	}
}

// ---------------------------------------------------------------------------
// Engine methods
// ---------------------------------------------------------------------------

// SetBeforeHook registers a callback that runs before each transition.
// Returning an error from the hook prevents the transition from completing.
func (w *WorkflowEngine) SetBeforeHook(hook TransitionHook) {
	w.onBefore = hook
}

// SetAfterHook registers a callback that runs after a successful transition.
func (w *WorkflowEngine) SetAfterHook(hook TransitionHook) {
	w.onAfter = hook
}

// ValidateTransition checks whether the engine allows transitioning from one
// status to another. Returns a WorkflowTransitionError if the transition is
// not in the rule table.
func (w *WorkflowEngine) ValidateTransition(from, to enums.Status) error {
	for _, r := range w.transitions {
		if r.From == from && r.To == to {
			return nil
		}
	}
	return &WorkflowTransitionError{
		From:   from,
		To:     to,
		Reason: "no matching transition rule",
	}
}

// CanTransition returns true if the engine has at least one rule that allows
// moving from the given source status to the target status.
func (w *WorkflowEngine) CanTransition(from, to enums.Status) bool {
	return w.ValidateTransition(from, to) == nil
}

// GetValidTransitions returns all transition rules that originate from the
// given status. Useful for building UI action menus.
func (w *WorkflowEngine) GetValidTransitions(from enums.Status) []TransitionRule {
	var result []TransitionRule
	for _, r := range w.transitions {
		if r.From == from {
			result = append(result, r)
		}
	}
	return result
}

// Execute validates the requested transition, invokes the before-hook (if set),
// logs the transition, and invokes the after-hook (if set).
func (w *WorkflowEngine) Execute(entityID string, from, to enums.Status, actorID string) error {
	if err := w.ValidateTransition(from, to); err != nil {
		w.log.Warn().
			Str("entity_id", entityID).
			Str("actor_id", actorID).
			Str("from", from.String()).
			Str("to", to.String()).
			Msg("workflow transition denied")
		return err
	}

	if w.onBefore != nil {
		if err := w.onBefore(entityID, from, to, actorID); err != nil {
			return fmt.Errorf("before-hook failed: %w", err)
		}
	}

	w.log.Info().
		Str("entity_id", entityID).
		Str("actor_id", actorID).
		Str("from", from.String()).
		Str("to", to.String()).
		Msg("workflow transition executed")

	if w.onAfter != nil {
		if err := w.onAfter(entityID, from, to, actorID); err != nil {
			return fmt.Errorf("after-hook failed: %w", err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Domain struct mutation helpers
// ---------------------------------------------------------------------------

// ApplyApproval marks a BaseWorkFlow entity as approved and clears any prior
// rejection state. The caller is responsible for persisting the change.
func ApplyApproval(wf *domain.BaseWorkFlow, approverID string) {
	now := time.Now().UTC()
	wf.IsApproved = true
	wf.ApprovedBy = approverID
	wf.DateApproved = &now
	wf.IsRejected = false
	wf.RejectionReason = ""
	wf.RecordStatus = enums.StatusApprovedAndActive.String()
}

// ApplyRejection marks a BaseWorkFlow entity as rejected. A non-empty reason
// is required by business rules; this prevents accidental rejections without
// explanation to the submitter.
func ApplyRejection(wf *domain.BaseWorkFlow, rejectedBy, reason string) error {
	if reason == "" {
		return ErrRejectionReasonRequired
	}
	now := time.Now().UTC()
	wf.IsRejected = true
	wf.RejectedBy = rejectedBy
	wf.RejectionReason = reason
	wf.DateRejected = &now
	wf.IsApproved = false
	wf.RecordStatus = enums.StatusRejected.String()
	return nil
}

// ApplyReturn sends a BaseWorkFlow entity back to the submitter for revision.
// Unlike rejection, a returned record can be edited and re-submitted without
// starting a new record.
func ApplyReturn(wf *domain.BaseWorkFlow, returnedBy, reason string) {
	wf.IsRejected = false
	wf.IsApproved = false
	wf.RejectionReason = reason
	wf.RecordStatus = enums.StatusReturned.String()
}

// ApplyHrdApproval marks an HrdWorkFlow entity as approved at the HRD level,
// completing the two-level approval chain.
func ApplyHrdApproval(wf *domain.HrdWorkFlow, approverID string) {
	now := time.Now().UTC()
	wf.HrdIsApproved = true
	wf.HrdApprovedBy = approverID
	wf.HrdDateApproved = &now
	wf.HrdIsRejected = false
	wf.HrdRejectionReason = ""
	wf.RecordStatus = enums.StatusApprovedAndActive.String()
}

// ApplyHrdRejection marks an HrdWorkFlow entity as rejected at the HRD level.
// A non-empty reason is required.
func ApplyHrdRejection(wf *domain.HrdWorkFlow, rejectedBy, reason string) error {
	if reason == "" {
		return ErrRejectionReasonRequired
	}
	now := time.Now().UTC()
	wf.HrdIsRejected = true
	wf.HrdRejectedBy = rejectedBy
	wf.HrdRejectionReason = reason
	wf.HrdDateRejected = &now
	wf.HrdIsApproved = false
	wf.RecordStatus = enums.StatusRejected.String()
	return nil
}

// ResetWorkflow clears all approval and rejection fields on a BaseWorkFlow
// entity, returning it to its initial state. Used when re-drafting or
// re-submitting a record.
func ResetWorkflow(wf *domain.BaseWorkFlow) {
	wf.IsApproved = false
	wf.ApprovedBy = ""
	wf.DateApproved = nil
	wf.IsRejected = false
	wf.RejectedBy = ""
	wf.RejectionReason = ""
	wf.DateRejected = nil
}
