package service

import (
	"errors"
	"testing"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ---------------------------------------------------------------------------
// WorkflowTransitionError
// ---------------------------------------------------------------------------

func TestWorkflowTransitionError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *WorkflowTransitionError
		contains []string
	}{
		{
			name: "basic transition message",
			err: &WorkflowTransitionError{
				From: enums.StatusDraft,
				To:   enums.StatusApprovedAndActive,
			},
			contains: []string{"Draft", "ApprovedAndActive"},
		},
		{
			name: "with entity ID",
			err: &WorkflowTransitionError{
				From:     enums.StatusDraft,
				To:       enums.StatusRejected,
				EntityID: "OBJ-001",
			},
			contains: []string{"Draft", "Rejected", "OBJ-001"},
		},
		{
			name: "with entity ID and reason",
			err: &WorkflowTransitionError{
				From:     enums.StatusPendingApproval,
				To:       enums.StatusApprovedAndActive,
				EntityID: "RP-042",
				Reason:   "approver not authorized",
			},
			contains: []string{"PendingApproval", "ApprovedAndActive", "RP-042", "approver not authorized"},
		},
		{
			name: "with reason but no entity ID",
			err: &WorkflowTransitionError{
				From:   enums.StatusReturned,
				To:     enums.StatusCompleted,
				Reason: "invalid path",
			},
			contains: []string{"Returned", "Completed", "invalid path"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.err.Error()
			for _, s := range tc.contains {
				if !containsStr(msg, s) {
					t.Errorf("expected error message to contain %q, got %q", s, msg)
				}
			}
		})
	}
}

func TestWorkflowTransitionError_Unwrap(t *testing.T) {
	err := &WorkflowTransitionError{
		From:     enums.StatusDraft,
		To:       enums.StatusApprovedAndActive,
		EntityID: "TEST-1",
		Reason:   "no matching rule",
	}

	if !errors.Is(err, ErrInvalidWorkflowTransition) {
		t.Error("expected WorkflowTransitionError to unwrap to ErrInvalidWorkflowTransition")
	}

	// Should NOT match other sentinel errors.
	if errors.Is(err, ErrWeightsNotBalanced) {
		t.Error("expected WorkflowTransitionError NOT to match ErrWeightsNotBalanced")
	}
}

// ---------------------------------------------------------------------------
// WeightValidationError
// ---------------------------------------------------------------------------

func TestWeightValidationError_Error(t *testing.T) {
	err := &WeightValidationError{
		ExpectedTotal: 100,
		ActualTotal:   90.5,
		CategoryID:    "CAT-A",
	}

	msg := err.Error()
	if !containsStr(msg, "90.50") {
		t.Errorf("expected actual total in message, got %q", msg)
	}
	if !containsStr(msg, "100.00") {
		t.Errorf("expected expected total in message, got %q", msg)
	}
	if !containsStr(msg, "CAT-A") {
		t.Errorf("expected category ID in message, got %q", msg)
	}
}

func TestWeightValidationError_Unwrap(t *testing.T) {
	err := &WeightValidationError{
		ExpectedTotal: 100,
		ActualTotal:   80,
		CategoryID:    "CAT-B",
	}

	if !errors.Is(err, ErrWeightsNotBalanced) {
		t.Error("expected WeightValidationError to unwrap to ErrWeightsNotBalanced")
	}

	if errors.Is(err, ErrInvalidWorkflowTransition) {
		t.Error("expected WeightValidationError NOT to match ErrInvalidWorkflowTransition")
	}
}

// ---------------------------------------------------------------------------
// ObjectiveLimitError
// ---------------------------------------------------------------------------

func TestObjectiveLimitError_Error(t *testing.T) {
	err := &ObjectiveLimitError{
		Max:      10,
		Current:  12,
		StaffID:  "STAFF-001",
		PeriodID: "RP-2024-Q1",
	}

	msg := err.Error()
	expectedFragments := []string{"12", "10", "STAFF-001", "RP-2024-Q1"}
	for _, frag := range expectedFragments {
		if !containsStr(msg, frag) {
			t.Errorf("expected error message to contain %q, got %q", frag, msg)
		}
	}
}

func TestObjectiveLimitError_Unwrap(t *testing.T) {
	err := &ObjectiveLimitError{
		Max:      5,
		Current:  6,
		StaffID:  "S1",
		PeriodID: "P1",
	}

	if !errors.Is(err, ErrMaxObjectivesExceeded) {
		t.Error("expected ObjectiveLimitError to unwrap to ErrMaxObjectivesExceeded")
	}

	if errors.Is(err, ErrWeightsNotBalanced) {
		t.Error("expected ObjectiveLimitError NOT to match ErrWeightsNotBalanced")
	}
}

// ---------------------------------------------------------------------------
// RangeValueError
// ---------------------------------------------------------------------------

func TestRangeValueError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *RangeValueError
		contains []string
	}{
		{
			name: "quarterly out of range",
			err: &RangeValueError{
				Range:    enums.ReviewPeriodRangeQuarterly,
				Value:    5,
				MaxValue: 4,
			},
			contains: []string{"5", "Quarterly", "4"},
		},
		{
			name: "bi-annual out of range",
			err: &RangeValueError{
				Range:    enums.ReviewPeriodRangeBiAnnual,
				Value:    3,
				MaxValue: 2,
			},
			contains: []string{"3", "BiAnnual", "2"},
		},
		{
			name: "annual out of range",
			err: &RangeValueError{
				Range:    enums.ReviewPeriodRangeAnnual,
				Value:    2,
				MaxValue: 1,
			},
			contains: []string{"2", "Annual", "1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.err.Error()
			for _, s := range tc.contains {
				if !containsStr(msg, s) {
					t.Errorf("expected error message to contain %q, got %q", s, msg)
				}
			}
		})
	}
}

func TestRangeValueError_Unwrap(t *testing.T) {
	err := &RangeValueError{
		Range:    enums.ReviewPeriodRangeQuarterly,
		Value:    5,
		MaxValue: 4,
	}

	if !errors.Is(err, ErrInvalidRangeValue) {
		t.Error("expected RangeValueError to unwrap to ErrInvalidRangeValue")
	}

	if errors.Is(err, ErrMaxObjectivesExceeded) {
		t.Error("expected RangeValueError NOT to match ErrMaxObjectivesExceeded")
	}
}

// ---------------------------------------------------------------------------
// Sentinel errors
// ---------------------------------------------------------------------------

func TestSentinelErrors(t *testing.T) {
	sentinels := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrInvalidWorkflowTransition", ErrInvalidWorkflowTransition, "invalid workflow state transition"},
		{"ErrRejectionReasonRequired", ErrRejectionReasonRequired, "rejection reason is required"},
		{"ErrUnauthorizedApprover", ErrUnauthorizedApprover, "caller is not authorized to approve this record"},
		{"ErrAlreadyApproved", ErrAlreadyApproved, "record has already been approved"},
		{"ErrAlreadyRejected", ErrAlreadyRejected, "record has already been rejected"},
		{"ErrDuplicateReviewPeriod", ErrDuplicateReviewPeriod, "a review period already exists for this range and year"},
		{"ErrNoActiveReviewPeriod", ErrNoActiveReviewPeriod, "no active review period found"},
		{"ErrPeriodNotOpen", ErrPeriodNotOpen, "review period is not in the required state for this operation"},
		{"ErrMultipleActivePeriods", ErrMultipleActivePeriods, "only one review period can be active at a time"},
		{"ErrObjectiveNotFound", ErrObjectiveNotFound, "referenced objective not found"},
		{"ErrMaxObjectivesExceeded", ErrMaxObjectivesExceeded, "maximum number of objectives exceeded"},
		{"ErrWeightsNotBalanced", ErrWeightsNotBalanced, "category weights must sum to 100%"},
		{"ErrParentObjectiveNotFound", ErrParentObjectiveNotFound, "parent objective not found for cascading"},
		{"ErrStrategyNotActive", ErrStrategyNotActive, "strategy must be in ApprovedAndActive status"},
		{"ErrInvalidRangeValue", ErrInvalidRangeValue, "invalid range value for the specified period type"},
		{"ErrScoreOutOfRange", ErrScoreOutOfRange, "score value is outside the valid range"},
		{"ErrNoScoreData", ErrNoScoreData, "no score data available for calculation"},
		{"ErrInvalidWeightConfig", ErrInvalidWeightConfig, "invalid weight configuration for scoring"},
	}

	for _, tc := range sentinels {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatalf("sentinel error %s is nil", tc.name)
			}
			if tc.err.Error() != tc.msg {
				t.Errorf("expected message %q, got %q", tc.msg, tc.err.Error())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// errors.As extraction
// ---------------------------------------------------------------------------

func TestErrorsAs_WorkflowTransitionError(t *testing.T) {
	var target *WorkflowTransitionError
	err := &WorkflowTransitionError{
		From:     enums.StatusDraft,
		To:       enums.StatusApprovedAndActive,
		EntityID: "E1",
		Reason:   "test",
	}

	if !errors.As(err, &target) {
		t.Fatal("errors.As should match WorkflowTransitionError")
	}
	if target.EntityID != "E1" {
		t.Errorf("expected EntityID E1, got %s", target.EntityID)
	}
}

func TestErrorsAs_WeightValidationError(t *testing.T) {
	var target *WeightValidationError
	err := &WeightValidationError{
		ExpectedTotal: 100,
		ActualTotal:   85,
		CategoryID:    "CAT-X",
	}

	if !errors.As(err, &target) {
		t.Fatal("errors.As should match WeightValidationError")
	}
	if target.CategoryID != "CAT-X" {
		t.Errorf("expected CategoryID CAT-X, got %s", target.CategoryID)
	}
}

func TestErrorsAs_ObjectiveLimitError(t *testing.T) {
	var target *ObjectiveLimitError
	err := &ObjectiveLimitError{Max: 5, Current: 6, StaffID: "S1", PeriodID: "P1"}

	if !errors.As(err, &target) {
		t.Fatal("errors.As should match ObjectiveLimitError")
	}
	if target.Max != 5 || target.Current != 6 {
		t.Errorf("expected Max=5 Current=6, got Max=%d Current=%d", target.Max, target.Current)
	}
}

func TestErrorsAs_RangeValueError(t *testing.T) {
	var target *RangeValueError
	err := &RangeValueError{
		Range:    enums.ReviewPeriodRangeQuarterly,
		Value:    5,
		MaxValue: 4,
	}

	if !errors.As(err, &target) {
		t.Fatal("errors.As should match RangeValueError")
	}
	if target.Value != 5 {
		t.Errorf("expected Value=5, got %d", target.Value)
	}
}

// ---------------------------------------------------------------------------
// rangeName helper (exercised through RangeValueError)
// ---------------------------------------------------------------------------

func TestRangeName_UnknownRange(t *testing.T) {
	err := &RangeValueError{
		Range:    enums.ReviewPeriodRange(99),
		Value:    1,
		MaxValue: 1,
	}

	msg := err.Error()
	if !containsStr(msg, "Unknown") {
		t.Errorf("expected 'Unknown' in error message for invalid range, got %q", msg)
	}
}

// ---------------------------------------------------------------------------
// helper
// ---------------------------------------------------------------------------

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchStr(s, substr)
}

func searchStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
