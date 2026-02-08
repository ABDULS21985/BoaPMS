package service

// strategy_service_stubs.go contains stub implementations for
// PerformanceManagementService methods that are NOT part of the .NET
// SMDService.cs. These will be implemented in separate service files
// as the conversion progresses (e.g., project service, work product
// service, feedback service, scoring service, audit service).

import (
	"context"
	"fmt"
)

// ---------------------------------------------------------------------------
// Projects
// ---------------------------------------------------------------------------

func (s *strategyService) SetupProject(ctx context.Context, req interface{}) error {
	return fmt.Errorf("not implemented: SetupProject")
}

func (s *strategyService) GetProjectsByStaff(ctx context.Context, staffID string) (interface{}, error) {
	return nil, fmt.Errorf("not implemented: GetProjectsByStaff")
}

func (s *strategyService) AddProjectObjective(ctx context.Context, req interface{}) error {
	return fmt.Errorf("not implemented: AddProjectObjective")
}

func (s *strategyService) AddProjectMember(ctx context.Context, req interface{}) error {
	return fmt.Errorf("not implemented: AddProjectMember")
}

// ---------------------------------------------------------------------------
// Work Products
// ---------------------------------------------------------------------------

func (s *strategyService) AddWorkProduct(ctx context.Context, req interface{}) error {
	return fmt.Errorf("not implemented: AddWorkProduct")
}

func (s *strategyService) EvaluateWorkProduct(ctx context.Context, req interface{}) error {
	return fmt.Errorf("not implemented: EvaluateWorkProduct")
}

// ---------------------------------------------------------------------------
// Feedback
// ---------------------------------------------------------------------------

func (s *strategyService) RequestFeedback(ctx context.Context, req interface{}) error {
	return fmt.Errorf("not implemented: RequestFeedback")
}

func (s *strategyService) GetFeedbackRequests(ctx context.Context, staffID string) (interface{}, error) {
	return nil, fmt.Errorf("not implemented: GetFeedbackRequests")
}

// ---------------------------------------------------------------------------
// Competency
// ---------------------------------------------------------------------------

func (s *strategyService) GetCompetencyReview(ctx context.Context, staffID string) (interface{}, error) {
	return nil, fmt.Errorf("not implemented: GetCompetencyReview")
}

func (s *strategyService) Submit360Feedback(ctx context.Context, req interface{}) error {
	return fmt.Errorf("not implemented: Submit360Feedback")
}

// ---------------------------------------------------------------------------
// Scoring
// ---------------------------------------------------------------------------

func (s *strategyService) GetPerformanceScore(ctx context.Context, staffID string) (interface{}, error) {
	return nil, fmt.Errorf("not implemented: GetPerformanceScore")
}

func (s *strategyService) GetDashboardStats(ctx context.Context, staffID string) (interface{}, error) {
	return nil, fmt.Errorf("not implemented: GetDashboardStats")
}

// ---------------------------------------------------------------------------
// Audit
// ---------------------------------------------------------------------------

func (s *strategyService) LogAuditAction(ctx context.Context, action string, details interface{}) error {
	return fmt.Errorf("not implemented: LogAuditAction")
}
