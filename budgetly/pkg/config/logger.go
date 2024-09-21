package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ZiadMansourM/budgetly/pkg/prettylog"
)

// InitializeLogger initializes the logger based on the environment variables
func InitializeLogger() *slog.Logger {
	// Set the log level based on the environment variable
	levelStr := os.Getenv("LOG_LEVEL")

	// Default to info level
	var level slog.Level

	switch strings.ToLower(levelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "info":
		fallthrough
	default:
		level = slog.LevelInfo
	}

	// Determine the log format based on environment variable
	format := "pretty"
	if envFormat := os.Getenv("LOG_FRAMEWORK"); envFormat != "" {
		format = envFormat
	}
	var logger *slog.Logger

	// Use `prettylog` if LOG_FRAMEWORK is set to "pretty", otherwise use slog
	if strings.ToLower(format) == "pretty" {
		logger = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level: level,
		}))
	} else {
		// Default to the text handler (slog's default)
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
		logger = slog.New(handler)
	}

	// Optionally set as default logger
	slog.SetDefault(logger)

	return logger
}
