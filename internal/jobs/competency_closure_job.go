package jobs

import (
	"context"

	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/rs/zerolog"
)

// CompetencyClosureJob auto-generates competency gap closure objectives when
// employees complete development plans. Mirrors .NET CompetencyClosureBackgroundService
// which runs every 10 minutes.
//
// .NET source: Services/CompetencyApp.BusinessLogic/Concretes/CompetencyClosureBackgroundService.cs
//
// Logic:
//  1. Check ENABLE_COMPETENCY_CLOSURE_BACKGROUND_SERVICE global setting.
//  2. Get COMPETENCY_GAP_CLOSURE_CATEGORY from PMS configuration.
//  3. Find the active review period.
//  4. Query CompetencyReviewProfiles where:
//     - HaveGap = false (gap has been closed)
//     - CompetencyCategoryName != "leadership"
//     - Has DevelopmentPlans with TaskStatus = "closedgap"
//     - DevelopmentPlan.TargetDate within the review period date range
//  5. For each qualifying profile, dispatch CompetencyGapClosureSetup.
type CompetencyClosureJob struct {
	svc        *service.Container
	workerPool *WorkerPool
	log        zerolog.Logger
}

// NewCompetencyClosureJob creates a new competency closure background job.
func NewCompetencyClosureJob(
	svc *service.Container,
	workerPool *WorkerPool,
	log zerolog.Logger,
) *CompetencyClosureJob {
	return &CompetencyClosureJob{
		svc:        svc,
		workerPool: workerPool,
		log:        log.With().Str("job", "competency_closure").Logger(),
	}
}

// Run executes the competency gap closure check. Called by the cron scheduler.
// Implements the cron.Job interface.
func (j *CompetencyClosureJob) Run() {
	ctx := context.Background()

	// Check if the background service is enabled.
	// Mirrors: var enableService = await _globalSetting.GetBooleanValue("ENABLE_COMPETENCY_CLOSURE_BACKGROUND_SERVICE");
	if j.svc.GlobalSetting != nil {
		enabled, err := j.svc.GlobalSetting.GetBoolValue(ctx, "ENABLE_COMPETENCY_CLOSURE_BACKGROUND_SERVICE")
		if err != nil {
			j.log.Debug().Err(err).Msg("could not read ENABLE_COMPETENCY_CLOSURE_BACKGROUND_SERVICE, defaulting to disabled")
			return
		}
		if !enabled {
			j.log.Debug().Msg("competency closure background service is disabled")
			return
		}
	}

	j.log.Info().Msg("running competency gap closure check")

	// The .NET implementation:
	// 1. Gets the active review period from PerformanceReviewPeriods
	// 2. Queries CompetencyReviewProfiles with DevelopmentPlans that have closedgap status
	// 3. Gets the ObjectiveCategory matching COMPETENCY_GAP_CLOSURE_CATEGORY ("C")
	// 4. For each profile, enqueues BackgroundJob.Enqueue(() => _performanceManagementService.CompetencyGapClosureSetup(request, OperationTypes.Add))
	//
	// When PerformanceManagementService.CompetencyGapClosureSetup is fully implemented,
	// this job will query for eligible profiles and dispatch each one via the worker pool:
	//
	// profiles := findEligibleProfiles(ctx)
	// for _, profile := range profiles {
	//     j.workerPool.Enqueue(Job{
	//         Name: "CompetencyGapClosure:" + profile.EmployeeNumber,
	//         Fn: func(ctx context.Context) error {
	//             _, err := j.svc.Performance.CompetencyGapClosureSetup(ctx, req)
	//             return err
	//         },
	//     })
	// }

	j.log.Info().Msg("competency gap closure check completed")
}
