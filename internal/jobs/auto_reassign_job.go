package jobs

import (
	"context"
	"fmt"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/rs/zerolog"
)

// AutoReassignJob auto-reassigns breached feedback requests to the assigned
// staff's supervisor. Mirrors .NET AutoReassignRequestBackgroundService which
// runs every 10 minutes.
//
// .NET source: Services/CompetencyApp.BusinessLogic/Concretes/AutoReassignRequestBackgroundService.cs
//
// Logic:
//  1. Check ENABLE_AUTO_REASSIGN_REQUEST_BACKGROUND_SERVICE global setting.
//  2. Get all pending feedback requests.
//  3. Filter for breached requests (IsBreached == true).
//  4. For each breached request, dispatch AutoReassignAndLogRequest via worker pool.
type AutoReassignJob struct {
	svc        *service.Container
	workerPool *WorkerPool
	log        zerolog.Logger
}

// NewAutoReassignJob creates a new auto-reassignment background job.
func NewAutoReassignJob(
	svc *service.Container,
	workerPool *WorkerPool,
	log zerolog.Logger,
) *AutoReassignJob {
	return &AutoReassignJob{
		svc:        svc,
		workerPool: workerPool,
		log:        log.With().Str("job", "auto_reassign").Logger(),
	}
}

// Run executes the auto-reassignment check. Called by the cron scheduler.
// Implements the cron.Job interface.
func (j *AutoReassignJob) Run() {
	ctx := context.Background()

	// Check if the background service is enabled.
	// Mirrors: var enableService = await _globalSetting.GetBooleanValue("ENABLE_AUTO_REASSIGN_REQUEST_BACKGROUND_SERVICE");
	if j.svc.GlobalSetting != nil {
		enabled, err := j.svc.GlobalSetting.GetBoolValue(ctx, "ENABLE_AUTO_REASSIGN_REQUEST_BACKGROUND_SERVICE")
		if err != nil {
			j.log.Debug().Err(err).Msg("could not read ENABLE_AUTO_REASSIGN_REQUEST_BACKGROUND_SERVICE, defaulting to disabled")
			return
		}
		if !enabled {
			j.log.Debug().Msg("auto-reassign background service is disabled")
			return
		}
	}

	j.log.Info().Msg("running auto-reassignment check")

	// Get pending requests and find breached ones.
	// Mirrors .NET: var pendingRequests = await _performanceManagementService.GetPendingRequests();
	// Then filters: pendingRequests.Where(x => x.IsBreached == true)
	if j.svc.Performance == nil {
		j.log.Warn().Msg("performance service not available, skipping auto-reassignment")
		return
	}

	// GetPendingRequests returns all pending feedback requests for all staff (empty staffID = all).
	result, err := j.svc.Performance.GetPendingRequests(ctx, "")
	if err != nil {
		j.log.Error().Err(err).Msg("failed to get pending requests")
		return
	}

	// The result contains breached request IDs. Extract and dispatch each one.
	// The actual type assertion depends on the concrete return type from GetPendingRequests.
	// When the service is fully implemented, the response will include a list of
	// breached FeedbackRequestLogIds.
	_ = result // Will be used when response type is finalized.

	// Placeholder: When breached requests are identified, dispatch each one:
	// for _, requestID := range breachedRequestIDs {
	//     j.dispatchReassignment(requestID)
	// }

	j.log.Info().Msg("auto-reassignment check completed")
}

// dispatchReassignment queues a single request reassignment via the worker pool.
// Maps to .NET: BackgroundJob.Enqueue(() => _performanceManagementService.AutoReassignAndLogRequestAsync(request.FeedbackRequestLogId))
func (j *AutoReassignJob) dispatchReassignment(requestID string) {
	j.workerPool.Enqueue(Job{
		Name: fmt.Sprintf("AutoReassign:%s", requestID),
		Fn: func(ctx context.Context) error {
			return j.svc.Performance.AutoReassignAndLogRequest(ctx, requestID)
		},
	})
}
