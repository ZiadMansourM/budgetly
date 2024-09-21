package config

import (
	"log/slog"
	"os"
	"strings"
)

// InitializeLogger sets up the global logger based on environment variables
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

	// Create a new logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level, // Set the log level
	})
	logger := slog.New(handler)

	// Optionally set as default logger
	slog.SetDefault(logger)

	return logger
}
