package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	DBConnectionString string
	ServerAddress      string
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig(file string) (*Config, error) {
	// Attempt to load environment variables from the .env file (optional)
	err := loadEnvFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found. Falling back to system environment variables.")
		} else {
			// Return any other errors (e.g., file permission issues)
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
	}

	// Initialize config from environment variables
	config, err := initConfigFromEnv()
	if err != nil {
		return nil, err
	}

	// Return the loaded and validated config
	return config, nil
}

// loadEnvFile loads environment variables from a .env file into the process environment
func loadEnvFile(file string) error {
	envFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines or comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Ignoring malformed line in .env file: %s\n", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set environment variable if not already set
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// initConfigFromEnv reads environment variables and initializes the Config struct
func initConfigFromEnv() (*Config, error) {
	config := &Config{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		ServerAddress:      os.Getenv("SERVER_ADDRESS"),
	}

	// Validate mandatory environment variables
	if config.DBConnectionString == "" {
		return nil, fmt.Errorf("DB_CONNECTION_STRING is required but not set")
	}

	// Set default values if optional env vars are not set
	if config.ServerAddress == "" {
		config.ServerAddress = "127.0.0.1:8080" // Default for local dev
	}

	return config, nil
}
