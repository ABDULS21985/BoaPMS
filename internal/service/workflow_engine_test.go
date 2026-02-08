package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/enterprise-pms/pms-api/internal/domain"
	"github.com/enterprise-pms/pms-api/internal/domain/enums"
	"github.com/rs/zerolog"
)

// nopLogger returns a no-op zerolog logger for tests.
func nopLogger() zerolog.Logger {
	return zerolog.Nop()
}

// ---------------------------------------------------------------------------
// BaseWorkflowEngine — valid transitions
// ---------------------------------------------------------------------------

func TestBaseWorkflow_DraftToPendingApproval(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.ValidateTransition(enums.StatusDraft, enums.StatusPendingApproval)
	if err != nil {
		t.Fatalf("expected Draft -> PendingApproval to be valid, got error: %v", err)
	}
}

func TestBaseWorkflow_PendingApprovalToApproved(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.ValidateTransition(enums.StatusPendingApproval, enums.StatusApprovedAndActive)
	if err != nil {
		t.Fatalf("expected PendingApproval -> ApprovedAndActive to be valid, got error: %v", err)
	}
}

func TestBaseWorkflow_PendingApprovalToRejected(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.ValidateTransition(enums.StatusPendingApproval, enums.StatusRejected)
	if err != nil {
		t.Fatalf("expected PendingApproval -> Rejected to be valid, got error: %v", err)
	}
}

func TestBaseWorkflow_PendingApprovalToReturned(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.ValidateTransition(enums.StatusPendingApproval, enums.StatusReturned)
	if err != nil {
		t.Fatalf("expected PendingApproval -> Returned to be valid, got error: %v", err)
	}
}

func TestBaseWorkflow_ReturnedToResubmit(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	// Returned -> PendingApproval (re-submit)
	err := engine.ValidateTransition(enums.StatusReturned, enums.StatusPendingApproval)
	if err != nil {
		t.Fatalf("expected Returned -> PendingApproval to be valid, got error: %v", err)
	}
}

func TestBaseWorkflow_InvalidTransition(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.ValidateTransition(enums.StatusDraft, enums.StatusApprovedAndActive)
	if err == nil {
		t.Fatal("expected Draft -> ApprovedAndActive to be invalid, got nil")
	}

	if !errors.Is(err, ErrInvalidWorkflowTransition) {
		t.Errorf("expected ErrInvalidWorkflowTransition, got: %v", err)
	}

	var wfErr *WorkflowTransitionError
	if !errors.As(err, &wfErr) {
		t.Fatal("expected error to be WorkflowTransitionError")
	}
	if wfErr.From != enums.StatusDraft || wfErr.To != enums.StatusApprovedAndActive {
		t.Errorf("expected From=Draft To=ApprovedAndActive, got From=%s To=%s",
			wfErr.From.String(), wfErr.To.String())
	}
}

func TestBaseWorkflow_ApprovedToDeactivated(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.ValidateTransition(enums.StatusApprovedAndActive, enums.StatusDeactivated)
	if err != nil {
		t.Fatalf("expected ApprovedAndActive -> Deactivated to be valid, got error: %v", err)
	}
}

func TestBaseWorkflow_DeactivatedToReactivated(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	// Deactivated -> ApprovedAndActive (reactivation)
	err := engine.ValidateTransition(enums.StatusDeactivated, enums.StatusApprovedAndActive)
	if err != nil {
		t.Fatalf("expected Deactivated -> ApprovedAndActive to be valid, got error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// GetValidTransitions
// ---------------------------------------------------------------------------

func TestBaseWorkflow_GetValidTransitions(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	transitions := engine.GetValidTransitions(enums.StatusPendingApproval)

	// From PendingApproval, the base engine allows:
	// -> ApprovedAndActive (Approve), -> Rejected (Reject), -> Returned (Return)
	if len(transitions) < 3 {
		t.Errorf("expected at least 3 transitions from PendingApproval, got %d", len(transitions))
	}

	// Verify each expected target status is present.
	expectedTargets := map[enums.Status]bool{
		enums.StatusApprovedAndActive: false,
		enums.StatusRejected:          false,
		enums.StatusReturned:          false,
	}
	for _, tr := range transitions {
		if _, ok := expectedTargets[tr.To]; ok {
			expectedTargets[tr.To] = true
		}
	}
	for status, found := range expectedTargets {
		if !found {
			t.Errorf("expected transition to %s from PendingApproval, but not found", status.String())
		}
	}
}

// ---------------------------------------------------------------------------
// CanTransition
// ---------------------------------------------------------------------------

func TestBaseWorkflow_CanTransition(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	if !engine.CanTransition(enums.StatusDraft, enums.StatusPendingApproval) {
		t.Error("expected CanTransition(Draft, PendingApproval) to be true")
	}

	if engine.CanTransition(enums.StatusDraft, enums.StatusClosed) {
		t.Error("expected CanTransition(Draft, Closed) to be false")
	}
}

// ---------------------------------------------------------------------------
// ReviewPeriodWorkflowEngine
// ---------------------------------------------------------------------------

func TestReviewPeriodWorkflow_FullLifecycle(t *testing.T) {
	engine := NewReviewPeriodWorkflowEngine(nopLogger())

	// Draft -> PendingApproval
	if err := engine.ValidateTransition(enums.StatusDraft, enums.StatusPendingApproval); err != nil {
		t.Fatalf("Draft -> PendingApproval failed: %v", err)
	}

	// PendingApproval -> ApprovedAndActive
	if err := engine.ValidateTransition(enums.StatusPendingApproval, enums.StatusApprovedAndActive); err != nil {
		t.Fatalf("PendingApproval -> ApprovedAndActive failed: %v", err)
	}

	// ApprovedAndActive -> Closed
	if err := engine.ValidateTransition(enums.StatusApprovedAndActive, enums.StatusClosed); err != nil {
		t.Fatalf("ApprovedAndActive -> Closed failed: %v", err)
	}
}

func TestReviewPeriodWorkflow_RejectedToResubmit(t *testing.T) {
	engine := NewReviewPeriodWorkflowEngine(nopLogger())

	// Review periods allow re-submission from Rejected (unlike base workflow)
	err := engine.ValidateTransition(enums.StatusRejected, enums.StatusPendingApproval)
	if err != nil {
		t.Fatalf("expected Rejected -> PendingApproval to be valid for review periods, got error: %v", err)
	}
}

func TestReviewPeriodWorkflow_ApprovedToCancelled(t *testing.T) {
	engine := NewReviewPeriodWorkflowEngine(nopLogger())

	err := engine.ValidateTransition(enums.StatusApprovedAndActive, enums.StatusCancelled)
	if err != nil {
		t.Fatalf("expected ApprovedAndActive -> Cancelled to be valid for review periods, got error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Execute method
// ---------------------------------------------------------------------------

func TestBaseWorkflow_Execute_ValidTransition(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.Execute("ENT-1", enums.StatusDraft, enums.StatusPendingApproval, "actor-1")
	if err != nil {
		t.Fatalf("Execute for valid transition returned error: %v", err)
	}
}

func TestBaseWorkflow_Execute_InvalidTransition(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	err := engine.Execute("ENT-1", enums.StatusDraft, enums.StatusClosed, "actor-1")
	if err == nil {
		t.Fatal("Execute for invalid transition should return an error")
	}
	if !errors.Is(err, ErrInvalidWorkflowTransition) {
		t.Errorf("expected ErrInvalidWorkflowTransition, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Hooks
// ---------------------------------------------------------------------------

func TestWorkflow_BeforeHook_Called(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	called := false
	engine.SetBeforeHook(func(entityID string, from, to enums.Status, actorID string) error {
		called = true
		if entityID != "ENT-1" {
			t.Errorf("expected entityID ENT-1, got %s", entityID)
		}
		if actorID != "actor-A" {
			t.Errorf("expected actorID actor-A, got %s", actorID)
		}
		return nil
	})

	err := engine.Execute("ENT-1", enums.StatusDraft, enums.StatusPendingApproval, "actor-A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected before hook to be called")
	}
}

func TestWorkflow_AfterHook_Called(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	called := false
	engine.SetAfterHook(func(entityID string, from, to enums.Status, actorID string) error {
		called = true
		if from != enums.StatusDraft {
			t.Errorf("expected from=Draft, got %s", from.String())
		}
		if to != enums.StatusPendingApproval {
			t.Errorf("expected to=PendingApproval, got %s", to.String())
		}
		return nil
	})

	err := engine.Execute("ENT-2", enums.StatusDraft, enums.StatusPendingApproval, "actor-B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected after hook to be called")
	}
}

func TestWorkflow_BeforeHook_Error_StopsExecution(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	hookErr := fmt.Errorf("permission denied")
	afterCalled := false

	engine.SetBeforeHook(func(entityID string, from, to enums.Status, actorID string) error {
		return hookErr
	})
	engine.SetAfterHook(func(entityID string, from, to enums.Status, actorID string) error {
		afterCalled = true
		return nil
	})

	err := engine.Execute("ENT-3", enums.StatusDraft, enums.StatusPendingApproval, "actor-C")
	if err == nil {
		t.Fatal("expected error when before hook fails")
	}
	if !errors.Is(err, hookErr) {
		t.Errorf("expected wrapped hook error, got: %v", err)
	}
	if afterCalled {
		t.Error("after hook should NOT be called when before hook returns an error")
	}
}

func TestWorkflow_AfterHook_Error_PropagatesError(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	afterErr := fmt.Errorf("notification failed")
	engine.SetAfterHook(func(entityID string, from, to enums.Status, actorID string) error {
		return afterErr
	})

	err := engine.Execute("ENT-4", enums.StatusDraft, enums.StatusPendingApproval, "actor-D")
	if err == nil {
		t.Fatal("expected error when after hook fails")
	}
	if !errors.Is(err, afterErr) {
		t.Errorf("expected wrapped after-hook error, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Domain struct mutation helpers
// ---------------------------------------------------------------------------

func TestApplyApproval(t *testing.T) {
	wf := &domain.BaseWorkFlow{
		IsRejected:      true,
		RejectedBy:      "old-rejector",
		RejectionReason: "old reason",
	}

	ApplyApproval(wf, "approver-1")

	if !wf.IsApproved {
		t.Error("expected IsApproved=true")
	}
	if wf.ApprovedBy != "approver-1" {
		t.Errorf("expected ApprovedBy=approver-1, got %s", wf.ApprovedBy)
	}
	if wf.DateApproved == nil {
		t.Error("expected DateApproved to be set")
	}
	if wf.IsRejected {
		t.Error("expected IsRejected=false after approval")
	}
	if wf.RejectionReason != "" {
		t.Errorf("expected RejectionReason to be cleared, got %q", wf.RejectionReason)
	}
	if wf.RecordStatus != enums.StatusApprovedAndActive.String() {
		t.Errorf("expected RecordStatus=%s, got %s",
			enums.StatusApprovedAndActive.String(), wf.RecordStatus)
	}
}

func TestApplyRejection_WithReason(t *testing.T) {
	wf := &domain.BaseWorkFlow{
		IsApproved: true,
		ApprovedBy: "old-approver",
	}

	err := ApplyRejection(wf, "rejector-1", "inadequate documentation")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !wf.IsRejected {
		t.Error("expected IsRejected=true")
	}
	if wf.RejectedBy != "rejector-1" {
		t.Errorf("expected RejectedBy=rejector-1, got %s", wf.RejectedBy)
	}
	if wf.RejectionReason != "inadequate documentation" {
		t.Errorf("expected reason='inadequate documentation', got %q", wf.RejectionReason)
	}
	if wf.DateRejected == nil {
		t.Error("expected DateRejected to be set")
	}
	if wf.IsApproved {
		t.Error("expected IsApproved=false after rejection")
	}
	if wf.RecordStatus != enums.StatusRejected.String() {
		t.Errorf("expected RecordStatus=%s, got %s",
			enums.StatusRejected.String(), wf.RecordStatus)
	}
}

func TestApplyRejection_EmptyReason(t *testing.T) {
	wf := &domain.BaseWorkFlow{}

	err := ApplyRejection(wf, "rejector-2", "")
	if err == nil {
		t.Fatal("expected error for empty rejection reason")
	}
	if !errors.Is(err, ErrRejectionReasonRequired) {
		t.Errorf("expected ErrRejectionReasonRequired, got: %v", err)
	}

	// Verify the struct was NOT mutated.
	if wf.IsRejected {
		t.Error("struct should not be mutated when rejection fails")
	}
}

func TestApplyReturn(t *testing.T) {
	wf := &domain.BaseWorkFlow{
		IsApproved: true,
		IsRejected: true,
	}

	ApplyReturn(wf, "returner-1", "needs more detail")

	if wf.IsRejected {
		t.Error("expected IsRejected=false after return")
	}
	if wf.IsApproved {
		t.Error("expected IsApproved=false after return")
	}
	if wf.RejectionReason != "needs more detail" {
		t.Errorf("expected reason stored, got %q", wf.RejectionReason)
	}
	if wf.RecordStatus != enums.StatusReturned.String() {
		t.Errorf("expected RecordStatus=%s, got %s",
			enums.StatusReturned.String(), wf.RecordStatus)
	}
}

func TestResetWorkflow(t *testing.T) {
	wf := &domain.BaseWorkFlow{
		IsApproved:      true,
		ApprovedBy:      "someone",
		IsRejected:      true,
		RejectedBy:      "someone-else",
		RejectionReason: "reason",
	}

	ResetWorkflow(wf)

	if wf.IsApproved {
		t.Error("expected IsApproved=false")
	}
	if wf.ApprovedBy != "" {
		t.Errorf("expected ApprovedBy='', got %q", wf.ApprovedBy)
	}
	if wf.DateApproved != nil {
		t.Error("expected DateApproved=nil")
	}
	if wf.IsRejected {
		t.Error("expected IsRejected=false")
	}
	if wf.RejectedBy != "" {
		t.Errorf("expected RejectedBy='', got %q", wf.RejectedBy)
	}
	if wf.RejectionReason != "" {
		t.Errorf("expected RejectionReason='', got %q", wf.RejectionReason)
	}
	if wf.DateRejected != nil {
		t.Error("expected DateRejected=nil")
	}
}

// ---------------------------------------------------------------------------
// HrdWorkFlow helpers
// ---------------------------------------------------------------------------

func TestApplyHrdApproval(t *testing.T) {
	wf := &domain.HrdWorkFlow{
		HrdIsRejected:      true,
		HrdRejectionReason: "old reason",
	}

	ApplyHrdApproval(wf, "hrd-approver")

	if !wf.HrdIsApproved {
		t.Error("expected HrdIsApproved=true")
	}
	if wf.HrdApprovedBy != "hrd-approver" {
		t.Errorf("expected HrdApprovedBy=hrd-approver, got %s", wf.HrdApprovedBy)
	}
	if wf.HrdDateApproved == nil {
		t.Error("expected HrdDateApproved to be set")
	}
	if wf.HrdIsRejected {
		t.Error("expected HrdIsRejected=false")
	}
	if wf.HrdRejectionReason != "" {
		t.Errorf("expected HrdRejectionReason cleared, got %q", wf.HrdRejectionReason)
	}
	if wf.RecordStatus != enums.StatusApprovedAndActive.String() {
		t.Errorf("expected RecordStatus=%s, got %s",
			enums.StatusApprovedAndActive.String(), wf.RecordStatus)
	}
}

func TestApplyHrdRejection_WithReason(t *testing.T) {
	wf := &domain.HrdWorkFlow{
		HrdIsApproved: true,
	}

	err := ApplyHrdRejection(wf, "hrd-rejector", "policy violation")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !wf.HrdIsRejected {
		t.Error("expected HrdIsRejected=true")
	}
	if wf.HrdRejectedBy != "hrd-rejector" {
		t.Errorf("expected HrdRejectedBy=hrd-rejector, got %s", wf.HrdRejectedBy)
	}
	if wf.HrdRejectionReason != "policy violation" {
		t.Errorf("expected reason='policy violation', got %q", wf.HrdRejectionReason)
	}
	if wf.HrdIsApproved {
		t.Error("expected HrdIsApproved=false after rejection")
	}
}

func TestApplyHrdRejection_EmptyReason(t *testing.T) {
	wf := &domain.HrdWorkFlow{}

	err := ApplyHrdRejection(wf, "hrd-rejector", "")
	if err == nil {
		t.Fatal("expected error for empty HRD rejection reason")
	}
	if !errors.Is(err, ErrRejectionReasonRequired) {
		t.Errorf("expected ErrRejectionReasonRequired, got: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Table-driven transition tests for comprehensive coverage
// ---------------------------------------------------------------------------

func TestBaseWorkflow_AllValidTransitions(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	validTransitions := []struct {
		from enums.Status
		to   enums.Status
	}{
		{enums.StatusDraft, enums.StatusPendingApproval},
		{enums.StatusPendingApproval, enums.StatusApprovedAndActive},
		{enums.StatusPendingApproval, enums.StatusRejected},
		{enums.StatusPendingApproval, enums.StatusReturned},
		{enums.StatusReturned, enums.StatusPendingApproval},
		{enums.StatusApprovedAndActive, enums.StatusDeactivated},
		{enums.StatusDeactivated, enums.StatusApprovedAndActive},
		{enums.StatusApprovedAndActive, enums.StatusClosed},
		{enums.StatusActive, enums.StatusClosed},
		{enums.StatusActive, enums.StatusCompleted},
	}

	for _, tc := range validTransitions {
		name := fmt.Sprintf("%s->%s", tc.from.String(), tc.to.String())
		t.Run(name, func(t *testing.T) {
			if err := engine.ValidateTransition(tc.from, tc.to); err != nil {
				t.Errorf("expected valid transition %s -> %s, got error: %v",
					tc.from.String(), tc.to.String(), err)
			}
		})
	}
}

func TestBaseWorkflow_InvalidTransitions(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	invalidTransitions := []struct {
		from enums.Status
		to   enums.Status
	}{
		{enums.StatusDraft, enums.StatusApprovedAndActive},
		{enums.StatusDraft, enums.StatusRejected},
		{enums.StatusDraft, enums.StatusClosed},
		{enums.StatusRejected, enums.StatusApprovedAndActive},
		{enums.StatusRejected, enums.StatusPendingApproval}, // Not allowed in base engine
		{enums.StatusClosed, enums.StatusDraft},
		{enums.StatusClosed, enums.StatusPendingApproval},
		{enums.StatusCompleted, enums.StatusDraft},
	}

	for _, tc := range invalidTransitions {
		name := fmt.Sprintf("%s->%s", tc.from.String(), tc.to.String())
		t.Run(name, func(t *testing.T) {
			err := engine.ValidateTransition(tc.from, tc.to)
			if err == nil {
				t.Errorf("expected invalid transition %s -> %s to fail",
					tc.from.String(), tc.to.String())
			}
		})
	}
}

func TestBaseWorkflow_GetValidTransitions_Draft(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	transitions := engine.GetValidTransitions(enums.StatusDraft)
	// Draft should only lead to PendingApproval (via CommitDraft or Add)
	for _, tr := range transitions {
		if tr.To != enums.StatusPendingApproval {
			t.Errorf("unexpected transition from Draft to %s", tr.To.String())
		}
	}
}

func TestBaseWorkflow_GetValidTransitions_NoTransitions(t *testing.T) {
	engine := NewBaseWorkflowEngine(nopLogger())

	transitions := engine.GetValidTransitions(enums.StatusCompleted)
	// Completed has no outgoing transitions in the base engine
	if len(transitions) != 0 {
		t.Errorf("expected 0 transitions from Completed, got %d", len(transitions))
	}
}
