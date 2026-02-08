package jobs

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog"
)

// Job represents a unit of work to be executed by the worker pool.
// This replaces Hangfire's BackgroundJob.Enqueue pattern from .NET.
type Job struct {
	Name string                           // Human-readable job identifier for logging.
	Fn   func(ctx context.Context) error  // The work to execute.
}

// WorkerPool manages a fixed-size pool of worker goroutines that process
// jobs from a buffered channel. It replaces the Hangfire queue system
// (pmsexecutions, feedbackreviews, workproductsevaluations, etc.).
type WorkerPool struct {
	jobs    chan Job
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	log     zerolog.Logger
	workers int
}

// NewWorkerPool creates a worker pool with the specified number of workers
// and job queue buffer size.
func NewWorkerPool(workers, queueSize int, log zerolog.Logger) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		jobs:    make(chan Job, queueSize),
		ctx:     ctx,
		cancel:  cancel,
		log:     log.With().Str("component", "worker_pool").Logger(),
		workers: workers,
	}
}

// Start launches the worker goroutines.
func (wp *WorkerPool) Start() {
	wp.log.Info().Int("workers", wp.workers).Msg("starting worker pool")
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Enqueue submits a job for asynchronous processing. If the pool is shutting
// down or the queue is full, the job is dropped with a warning log.
func (wp *WorkerPool) Enqueue(job Job) {
	select {
	case wp.jobs <- job:
		wp.log.Debug().Str("job", job.Name).Msg("job enqueued")
	case <-wp.ctx.Done():
		wp.log.Warn().Str("job", job.Name).Msg("worker pool shutting down, job dropped")
	}
}

// Shutdown signals workers to stop, closes the job channel, and waits for
// all in-flight jobs to complete.
func (wp *WorkerPool) Shutdown() {
	wp.log.Info().Msg("shutting down worker pool")
	wp.cancel()
	close(wp.jobs)
	wp.wg.Wait()
	wp.log.Info().Msg("worker pool stopped")
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	wp.log.Debug().Int("worker_id", id).Msg("worker started")

	for job := range wp.jobs {
		wp.executeJob(id, job)
	}

	wp.log.Debug().Int("worker_id", id).Msg("worker stopped")
}

func (wp *WorkerPool) executeJob(workerID int, job Job) {
	defer func() {
		if r := recover(); r != nil {
			wp.log.Error().
				Int("worker_id", workerID).
				Str("job", job.Name).
				Interface("panic", r).
				Str("stack", string(debug.Stack())).
				Msg("job panicked")
		}
	}()

	wp.log.Info().Int("worker_id", workerID).Str("job", job.Name).Msg("executing job")

	if err := job.Fn(wp.ctx); err != nil {
		wp.log.Error().
			Err(err).
			Int("worker_id", workerID).
			Str("job", job.Name).
			Msg("job failed")
		return
	}

	wp.log.Info().Int("worker_id", workerID).Str("job", job.Name).Msg("job completed")
}
