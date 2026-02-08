package service

import (
	"errors"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/enums"
)

// ---------------------------------------------------------------------------
// Domain-specific error types for workflow transitions, business rule
// violations, and scoring calculations.
//
// These errors allow callers to check failure reasons via errors.Is() and to
// extract structured context via errors.As().
// ---------------------------------------------------------------------------

// Sentinel errors — use errors.Is(err, ErrXxx) for type checking.
var (
	// Workflow errors
	ErrInvalidWorkflowTransition = errors.New("invalid workflow state transition")
	ErrRejectionReasonRequired   = errors.New("rejection reason is required")
	ErrUnauthorizedApprover      = errors.New("caller is not authorized to approve this record")
	ErrAlreadyApproved           = errors.New("record has already been approved")
	ErrAlreadyRejected           = errors.New("record has already been rejected")

	// Review period errors
	ErrDuplicateReviewPeriod = errors.New("a review period already exists for this range and year")
	ErrNoActiveReviewPeriod  = errors.New("no active review period found")
	ErrPeriodNotOpen         = errors.New("review period is not in the required state for this operation")
	ErrMultipleActivePeriods = errors.New("only one review period can be active at a time")

	// Objective errors
	ErrObjectiveNotFound     = errors.New("referenced objective not found")
	ErrMaxObjectivesExceeded = errors.New("maximum number of objectives exceeded")
	ErrWeightsNotBalanced    = errors.New("category weights must sum to 100%")
	ErrParentObjectiveNotFound = errors.New("parent objective not found for cascading")

	// Strategy errors
	ErrStrategyNotActive = errors.New("strategy must be in ApprovedAndActive status")

	// Range validation errors
	ErrInvalidRangeValue = errors.New("invalid range value for the specified period type")

	// Scoring errors
	ErrScoreOutOfRange     = errors.New("score value is outside the valid range")
	ErrNoScoreData         = errors.New("no score data available for calculation")
	ErrInvalidWeightConfig = errors.New("invalid weight configuration for scoring")
)

// ---------------------------------------------------------------------------
// Structured error types — use errors.As(err, &typedErr) for extraction.
// ---------------------------------------------------------------------------

// WorkflowTransitionError provides details about an invalid state transition.
type WorkflowTransitionError struct {
	From      enums.Status
	To        enums.Status
	EntityID  string
	Reason    string
}

func (e *WorkflowTransitionError) Error() string {
	msg := fmt.Sprintf("cannot transition from %s to %s", e.From.String(), e.To.String())
	if e.EntityID != "" {
		msg += fmt.Sprintf(" (entity: %s)", e.EntityID)
	}
	if e.Reason != "" {
		msg += fmt.Sprintf(": %s", e.Reason)
	}
	return msg
}

func (e *WorkflowTransitionError) Unwrap() error {
	return ErrInvalidWorkflowTransition
}

// WeightValidationError provides details about weight imbalance.
type WeightValidationError struct {
	ExpectedTotal float64
	ActualTotal   float64
	CategoryID    string
}

func (e *WeightValidationError) Error() string {
	return fmt.Sprintf("category weights sum to %.2f%%, expected %.2f%% (category: %s)",
		e.ActualTotal, e.ExpectedTotal, e.CategoryID)
}

func (e *WeightValidationError) Unwrap() error {
	return ErrWeightsNotBalanced
}

// ObjectiveLimitError provides details about objective count violations.
type ObjectiveLimitError struct {
	Max       int
	Current   int
	StaffID   string
	PeriodID  string
}

func (e *ObjectiveLimitError) Error() string {
	return fmt.Sprintf("maximum objectives exceeded: %d of %d allowed (staff: %s, period: %s)",
		e.Current, e.Max, e.StaffID, e.PeriodID)
}

func (e *ObjectiveLimitError) Unwrap() error {
	return ErrMaxObjectivesExceeded
}

// RangeValueError provides details about an invalid range value.
type RangeValueError struct {
	Range    enums.ReviewPeriodRange
	Value    int
	MaxValue int
}

func (e *RangeValueError) Error() string {
	return fmt.Sprintf("range value %d is invalid for %s (max: %d)",
		e.Value, rangeName(e.Range), e.MaxValue)
}

func (e *RangeValueError) Unwrap() error {
	return ErrInvalidRangeValue
}

// rangeName returns a human-readable name for a ReviewPeriodRange.
func rangeName(r enums.ReviewPeriodRange) string {
	switch r {
	case enums.ReviewPeriodRangeQuarterly:
		return "Quarterly"
	case enums.ReviewPeriodRangeBiAnnual:
		return "BiAnnual"
	case enums.ReviewPeriodRangeAnnual:
		return "Annual"
	default:
		return "Unknown"
	}
}
