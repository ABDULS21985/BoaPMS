package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/enterprise-pms/pms-api/internal/config"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// New creates a zerolog.Logger with console and file sinks
// replicating Serilog's file + console pattern with rolling retention.
func New(cfg config.LoggingConfig) zerolog.Logger {
	level := parseLevel(cfg.Level)

	var writers []io.Writer

	// Console sink
	if cfg.ConsoleEnabled {
		if cfg.Format == "text" {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			})
		} else {
			writers = append(writers, os.Stdout)
		}
	}

	// File sink with rotation (mirrors Serilog's RollingInterval.Day + retainedFileCountLimit)
	if cfg.FileEnabled && cfg.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0o755); err == nil {
			lj := &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.MaxSizeMB,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAgeDays,
				Compress:   cfg.Compress,
			}
			writers = append(writers, lj)
		}
	}

	var w io.Writer
	switch len(writers) {
	case 0:
		w = os.Stdout
	case 1:
		w = writers[0]
	default:
		w = zerolog.MultiLevelWriter(writers...)
	}

	return zerolog.New(w).
		Level(level).
		With().
		Timestamp().
		Caller().
		Str("service", "pms-api").
		Logger()
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
