package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/esavich/otus_project/internal/config"
)

func SetupLogger(cfg *config.Config) {
	// Set the log level based on the configuration
	var level slog.Level
	switch cfg.App.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opts := slog.HandlerOptions{}
	opts.Level = level
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts)).
		With(slog.String("service", cfg.App.ServiceName))
	slog.SetDefault(logger)
	slog.Info("Logger initialized")
	slog.Info(fmt.Sprintf("Logger set to level: %s", opts.Level))
}
