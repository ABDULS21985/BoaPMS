package performance

import "github.com/enterprise-pms/pms-api/internal/domain/enums"

// ===========================================================================
// Enum Wrapper VMs (EnumVms/*.cs)
//
// These DTOs pair an enum value with a display-friendly description for use
// in dropdowns and lookup lists.  ObjectiveLevelVm, PerformanceGradeVm,
// WorkProductTypeVm, and GrievanceTypeVm are already defined in their
// respective DTO files.
// ===========================================================================

// StatusVm pairs a Status enum value with its display text.
type StatusVm struct {
	Status      enums.Status `json:"status"`
	Description string       `json:"description"`
}

// FeedbackRequestTypeVm pairs a FeedbackRequestType enum value with its
// display text.
type FeedbackRequestTypeVm struct {
	FeedbackRequestType enums.FeedbackRequestType `json:"feedBackRequestType"`
	Description         string                    `json:"description"`
}

// ReviewPeriodRangeVm pairs a ReviewPeriodRange enum value with its
// display text.
type ReviewPeriodRangeVm struct {
	ReviewPeriodRange enums.ReviewPeriodRange `json:"reviewPeriodRange"`
	Description       string                 `json:"description"`
}

// EvaluationTypeVm is used by the frontend as a select-list option for
// evaluation types.  Unlike the other enum VMs it uses Text/Value/Disabled/
// Selected instead of enum+Description.
type EvaluationTypeVm struct {
	Text     string `json:"text"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
	Selected bool   `json:"selected"`
}
