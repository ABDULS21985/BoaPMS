package jobs

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/domain/performance"
)

// ---------------------------------------------------------------------------
// Typed dispatch helpers — wrap service method calls into Job structs for
// the worker pool. These replace Hangfire's BackgroundJob.Enqueue() calls
// found across the .NET controllers and background services.
// ---------------------------------------------------------------------------

// DispatchCompetencyGapClosure queues a competency gap closure setup.
// .NET: BackgroundJob.Enqueue(() => _performanceManagementService.CompetencyGapClosureSetup(request, OperationTypes.Add))
func (s *Scheduler) DispatchCompetencyGapClosure(req *performance.CompetencyGapClosureRequestModel) {
	s.workerPool.Enqueue(Job{
		Name: fmt.Sprintf("CompetencyGapClosure:%s", req.StaffID),
		Fn: func(ctx context.Context) error {
			_, err := s.svc.Performance.CompetencyGapClosureSetup(ctx, req)
			return err
		},
	})
}

// DispatchAutoReassign queues a feedback request reassignment.
// .NET: BackgroundJob.Enqueue(() => _performanceManagementService.AutoReassignAndLogRequestAsync(requestId))
func (s *Scheduler) DispatchAutoReassign(requestID string) {
	s.workerPool.Enqueue(Job{
		Name: fmt.Sprintf("AutoReassign:%s", requestID),
		Fn: func(ctx context.Context) error {
			return s.svc.Performance.AutoReassignAndLogRequest(ctx, requestID)
		},
	})
}

// DispatchInitiate360Review queues a 360-degree review initiation.
// .NET: BackgroundJob.Enqueue(() => _performanceManagementService.Initiate360Review(request))
// Queue: pmsexecutions
func (s *Scheduler) DispatchInitiate360Review(req *performance.Initiate360ReviewRequestModel) {
	s.workerPool.Enqueue(Job{
		Name: fmt.Sprintf("Initiate360Review:%s", req.ReviewPeriodID),
		Fn: func(ctx context.Context) error {
			_, err := s.svc.Performance.Initiate360Review(ctx, req)
			return err
		},
	})
}

// DispatchCloseReviewPeriodRequests queues closing all requests for a review period.
// .NET: BackgroundJob.Enqueue(() => _performanceManagementService.CloseReviewPeriodRequests(reviewPeriodId))
// Queue: requestclosure
func (s *Scheduler) DispatchCloseReviewPeriodRequests(reviewPeriodID string) {
	s.workerPool.Enqueue(Job{
		Name: fmt.Sprintf("CloseReviewPeriodRequests:%s", reviewPeriodID),
		Fn: func(ctx context.Context) error {
			return s.svc.Performance.CloseReviewPeriodRequests(ctx, reviewPeriodID)
		},
	})
}

// DispatchWorkProductSetup queues a work product creation/update.
// .NET: BackgroundJob.Enqueue(() => _performanceManagementService.WorkProductSetup(request, operationType))
// Queue: workproductsexecution
func (s *Scheduler) DispatchWorkProductSetup(req *performance.WorkProductRequestModel) {
	s.workerPool.Enqueue(Job{
		Name: fmt.Sprintf("WorkProductSetup:%s", req.StaffID),
		Fn: func(ctx context.Context) error {
			_, err := s.svc.Performance.WorkProductSetup(ctx, req)
			return err
		},
	})
}

// DispatchWorkProductEvaluation queues a work product evaluation.
// .NET: BackgroundJob.Enqueue(() => _performanceManagementService.WorkProductEvaluation(request, operationType))
// Queue: workproductsevaluations
func (s *Scheduler) DispatchWorkProductEvaluation(req *performance.WorkProductEvaluationRequestModel) {
	s.workerPool.Enqueue(Job{
		Name: "WorkProductEvaluation:" + req.WorkProductID,
		Fn: func(ctx context.Context) error {
			_, err := s.svc.Performance.WorkProductEvaluation(ctx, req)
			return err
		},
	})
}

// Enqueue exposes the worker pool's Enqueue method for ad-hoc job dispatch.
// This allows handlers and services to submit custom jobs without needing
// typed dispatch helpers.
func (s *Scheduler) Enqueue(job Job) {
	s.workerPool.Enqueue(job)
}
