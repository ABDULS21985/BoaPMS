package jobs

import (
	"context"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

// Scheduler manages all background jobs: cron-based recurring tasks and
// on-demand jobs via a worker pool. It replaces the .NET Hangfire server
// and BackgroundService hosted services.
//
// .NET equivalents:
//   - Hangfire RecurringJob → cron.AddJob with SkipIfStillRunning
//   - Hangfire BackgroundJob.Enqueue → WorkerPool.Enqueue
//   - BackgroundService (hosted services) → cron jobs at @every 10m
//   - Hangfire queues (pmsexecutions, etc.) → single WorkerPool
type Scheduler struct {
	cron       *cron.Cron
	workerPool *WorkerPool
	mailSender *MailSenderWorker
	svc        *service.Container
	repos      *repository.Container
	cfg        *config.Config
	log        zerolog.Logger
	cancel     context.CancelFunc
}

// NewScheduler creates a new job scheduler with all dependencies wired up.
func NewScheduler(
	svc *service.Container,
	repos *repository.Container,
	cfg *config.Config,
	log zerolog.Logger,
) *Scheduler {
	return &Scheduler{
		svc:   svc,
		repos: repos,
		cfg:   cfg,
		log:   log.With().Str("component", "scheduler").Logger(),
	}
}

// Start initializes and starts all background workers:
//  1. Worker pool for on-demand job dispatch.
//  2. Cron scheduler with 3 recurring jobs (@every 10m).
//  3. Mail sender worker (polls for Status='New' emails).
func (s *Scheduler) Start(ctx context.Context) {
	ctx, s.cancel = context.WithCancel(ctx)

	// --- Worker Pool ---
	poolSize := s.cfg.Jobs.WorkerPoolSize
	if poolSize <= 0 {
		poolSize = 5
	}
	queueSize := s.cfg.Jobs.WorkerQueueSize
	if queueSize <= 0 {
		queueSize = 100
	}
	s.workerPool = NewWorkerPool(poolSize, queueSize, s.log)
	s.workerPool.Start()

	// --- Cron Scheduler ---
	// Use SkipIfStillRunning to prevent overlapping executions when a job
	// takes longer than the 10-minute interval.
	cronLogger := newCronLogger(s.log)
	s.cron = cron.New(
		cron.WithLogger(cronLogger),
		cron.WithChain(cron.SkipIfStillRunning(cronLogger)),
	)

	schedule := s.cfg.Jobs.CronSchedule
	if schedule == "" {
		schedule = "@every 10m"
	}

	// Register recurring jobs — mirrors .NET BackgroundService registrations.
	reviewPeriodJob := NewReviewPeriodJob(s.svc, s.log)
	competencyClosureJob := NewCompetencyClosureJob(s.svc, s.workerPool, s.log)
	autoReassignJob := NewAutoReassignJob(s.svc, s.workerPool, s.log)

	if _, err := s.cron.AddJob(schedule, reviewPeriodJob); err != nil {
		s.log.Error().Err(err).Msg("failed to register review period job")
	}
	if _, err := s.cron.AddJob(schedule, competencyClosureJob); err != nil {
		s.log.Error().Err(err).Msg("failed to register competency closure job")
	}
	if _, err := s.cron.AddJob(schedule, autoReassignJob); err != nil {
		s.log.Error().Err(err).Msg("failed to register auto-reassign job")
	}

	s.cron.Start()
	s.log.Info().Str("schedule", schedule).Msg("cron scheduler started with 3 recurring jobs")

	// --- Mail Sender Worker ---
	if s.repos.Email != nil {
		interval := s.cfg.Jobs.MailSenderInterval
		s.mailSender = NewMailSenderWorker(s.repos.Email, s.cfg.Email, interval, s.log)
		go s.mailSender.Run(ctx)
		s.log.Info().Msg("mail sender worker started")
	} else {
		s.log.Warn().Msg("email repository not configured, mail sender worker not started")
	}

	s.log.Info().Msg("all background workers started")
}

// Stop gracefully shuts down all background workers.
// The shutdown sequence is:
//  1. Stop accepting new cron triggers.
//  2. Cancel the context (stops mail sender and in-flight jobs).
//  3. Wait for the worker pool to drain.
func (s *Scheduler) Stop() {
	s.log.Info().Msg("stopping scheduler")

	// Stop cron — waits for running jobs to finish.
	if s.cron != nil {
		cronCtx := s.cron.Stop()
		<-cronCtx.Done()
		s.log.Info().Msg("cron scheduler stopped")
	}

	// Cancel context to stop mail sender and worker pool.
	if s.cancel != nil {
		s.cancel()
	}

	// Drain worker pool.
	if s.workerPool != nil {
		s.workerPool.Shutdown()
	}

	s.log.Info().Msg("scheduler stopped")
}

// cronLogger adapts zerolog to the cron.Logger interface.
type cronLogAdapter struct {
	log zerolog.Logger
}

func newCronLogger(log zerolog.Logger) cron.Logger {
	return &cronLogAdapter{log: log.With().Str("component", "cron").Logger()}
}

func (l *cronLogAdapter) Info(msg string, keysAndValues ...interface{}) {
	l.log.Info().Fields(kvToMap(keysAndValues)).Msg(msg)
}

func (l *cronLogAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	l.log.Error().Err(err).Fields(kvToMap(keysAndValues)).Msg(msg)
}

func kvToMap(keysAndValues []interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(keysAndValues)/2)
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		if key, ok := keysAndValues[i].(string); ok {
			m[key] = keysAndValues[i+1]
		}
	}
	return m
}
