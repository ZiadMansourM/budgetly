package settings

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/ZiadMansourM/budgetly/pkg/prettylog"
)

// Settings holds the application's configuration, including database connection and server address.
type Settings struct {
	BaseDir            string
	DBConnectionString string
	ServerAddress      string
	Logger             *slog.Logger
}

// Init loads configuration and initializes logging based on environment variables or .env file.
func Init() (*Settings, error) {
	baseDir, err := getBaseDir()
	if err != nil {
		return nil, fmt.Errorf("error getting base directory: %v", err)
	}

	// Load environment variables from the .env file (if it exists).
	err = loadEnvFile(filepath.Join(baseDir, ".env"))
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found. Falling back to system environment variables.")
		} else {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	// Initialize logger.
	logger := initializeLogger()

	// Load configuration from environment variables.
	settings, err := initSettingsFromEnv()
	if err != nil {
		logger.Error("Error initializing settings", "error", err)
		return nil, err
	}

	// Attach the logger to the settings.
	settings.Logger = logger

	settings.BaseDir = baseDir

	return settings, nil
}

// initializeLogger initializes the logger based on environment variables.
func initializeLogger() *slog.Logger {
	// Set the log level based on the environment variable.
	levelStr := os.Getenv("LOG_LEVEL")
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

	// Determine the log format based on environment variables.
	format := "pretty"
	if envFormat := os.Getenv("LOG_FRAMEWORK"); envFormat != "" {
		format = envFormat
	}

	var logger *slog.Logger
	if strings.ToLower(format) == "pretty" {
		logger = slog.New(prettylog.NewHandler(&slog.HandlerOptions{Level: level}))
	} else {
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
		logger = slog.New(handler)
	}

	// Set as the default logger.
	slog.SetDefault(logger)
	return logger
}

// initSettingsFromEnv reads environment variables and initializes the Settings struct.
func initSettingsFromEnv() (*Settings, error) {
	settings := &Settings{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		ServerAddress:      os.Getenv("SERVER_ADDRESS"),
	}

	// Validate mandatory environment variables.
	if settings.DBConnectionString == "" {
		return nil, fmt.Errorf("DB_CONNECTION_STRING is required but not set")
	}

	// Set default values if optional env vars are not set.
	if settings.ServerAddress == "" {
		settings.ServerAddress = "127.0.0.1:8080" // Default for local dev.
	}

	return settings, nil
}

// loadEnvFile loads environment variables from a .env file into the process environment.
func loadEnvFile(file string) error {
	envFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines or comments.
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value.
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Ignoring malformed line in .env file: %s\n", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set environment variable if not already set.
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// getBaseDir determines the base directory of the current file or project
func getBaseDir() (string, error) {
	baseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	return baseDir, nil
}
