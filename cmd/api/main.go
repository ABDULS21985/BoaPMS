package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/enterprise-pms/pms-api/internal/handler"
	"github.com/enterprise-pms/pms-api/internal/jobs"
	"github.com/enterprise-pms/pms-api/internal/middleware"
	"github.com/enterprise-pms/pms-api/internal/repository"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logging)
	log.Info().Msg("Starting PMS API...")

	// Initialize repositories
	repos, err := repository.New(cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize repositories")
	}
	defer repos.Close()

	// Initialize services
	svc := service.New(repos, cfg, log)

	// Initialize and start background job scheduler
	// Replaces .NET Hangfire server + BackgroundService hosted services.
	scheduler := jobs.NewScheduler(svc, repos, cfg, log)
	scheduler.Start(context.Background())

	// Initialize middleware
	mw := middleware.New(cfg, log)

	// Initialize HTTP handlers and router
	router := handler.NewRouter(svc, mw, cfg, log)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Info().Int("port", cfg.Server.Port).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Shutting down server...")

	// Stop background jobs before shutting down HTTP server.
	scheduler.Stop()
	log.Info().Msg("Background jobs stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
