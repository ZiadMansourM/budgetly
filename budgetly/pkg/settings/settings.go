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
	EnvironmentMode    string
}

// SettingsBuilder is used to build the settings step by step.
type SettingsBuilder struct {
	settings *Settings
	err      error
}

// NewSettingsBuilder initializes the builder with default values.
func NewSettingsBuilder() *SettingsBuilder {
	return &SettingsBuilder{
		settings: &Settings{},
	}
}

// WithBaseDir sets the base directory.
func (b *SettingsBuilder) WithBaseDir() *SettingsBuilder {
	if b.err != nil {
		return b
	}

	baseDir, err := getBaseDir()
	if err != nil {
		b.err = fmt.Errorf("error getting base directory: %v", err)
		return b
	}
	b.settings.BaseDir = baseDir
	return b
}

// WithEnvironment loads the environment variables, prioritizing the correct .env file based on mode.
func (b *SettingsBuilder) WithEnvironment() *SettingsBuilder {
	if b.err != nil {
		return b
	}

	// Determine the mode (default is development).
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "development" // default to development if not set
	}
	b.settings.EnvironmentMode = mode

	// Load the appropriate .env file based on mode.
	var envPath string
	if mode == "production" {
		envPath = filepath.Join(b.settings.BaseDir, ".env")
	} else {
		envPath = ".env" // Load local .env for non-production
	}

	err := loadEnvFile(envPath)
	if err != nil && !os.IsNotExist(err) {
		b.err = fmt.Errorf("error loading .env file: %v", err)
		return b
	}

	return b
}

// WithLogger initializes the logger based on environment variables.
func (b *SettingsBuilder) WithLogger() *SettingsBuilder {
	if b.err != nil {
		return b
	}

	logger := initializeLogger()
	b.settings.Logger = logger
	return b
}

// WithDBConnection loads the DB connection string from environment variables.
func (b *SettingsBuilder) WithDBConnection() *SettingsBuilder {
	if b.err != nil {
		return b
	}

	dbConn := os.Getenv("DB_CONNECTION_STRING")
	if dbConn == "" {
		b.err = fmt.Errorf("DB_CONNECTION_STRING is required but not set")
		return b
	}
	b.settings.DBConnectionString = dbConn
	return b
}

// WithServerAddress loads the server address from environment variables.
func (b *SettingsBuilder) WithServerAddress() *SettingsBuilder {
	if b.err != nil {
		return b
	}

	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "127.0.0.1:8080" // Default for local dev.
	}
	b.settings.ServerAddress = serverAddr
	return b
}

// Build finalizes the settings creation and returns the Settings or an error.
func (b *SettingsBuilder) Build() (*Settings, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.settings, nil
}

// Initialize logger based on environment variables.
func initializeLogger() *slog.Logger {
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

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Ignoring malformed line in .env file: %s\n", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

// getBaseDir determines the base directory of the current file or project
func getBaseDir() (string, error) {
	baseDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return baseDir, nil
}
