package main

import (
	"github.com/ZiadMansourM/budgetly/cmd/api"
	"github.com/ZiadMansourM/budgetly/pkg/config"
	"github.com/ZiadMansourM/budgetly/pkg/middlewares"
)

func main() {
	// Initialize logger from environment variables
	logger := config.InitializeLogger()

	// Load configuration from environment variables or .env file
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Error("Error loading configuration", "error", err)
		return
	}

	// Use the builder to assemble the server with plug-and-play apps
	// E.g. WithUserApp which encapsulate all its components (model, service, handler, routes).
	serverBuilder := api.NewServerBuilder().
		WithDatabase("postgres", cfg.DBConnectionString).
		WithUserApp().
		WithHealthCheck().
		Use(middlewares.LoggingMiddleware).
		BuildServer(cfg.ServerAddress)

	// Start the server with graceful shutdown
	serverBuilder.StartServer()
}
