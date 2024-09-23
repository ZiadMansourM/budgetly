package main

import (
	"fmt"

	"github.com/ZiadMansourM/budgetly/cmd/api"
	"github.com/ZiadMansourM/budgetly/pkg/middlewares"
	"github.com/ZiadMansourM/budgetly/pkg/settings"
)

func main() {
	// Initialize settings using the builder pattern
	settings, err := settings.NewSettingsBuilder().
		WithBaseDir().
		WithEnvironment().
		WithLogger().
		WithDBConnection().
		WithServerAddress().
		Build()

	if err != nil {
		panic(fmt.Sprintf("Error initializing settings: %v", err))
	}

	// Use the builder to assemble the server with plug-and-play apps
	// E.g. WithUserApp which encapsulate all its components (model, service, handler, routes).
	serverBuilder := api.NewServerBuilder(settings.Logger).
		WithDatabase("postgres", settings.DBConnectionString).
		WithUserApp().
		WithHealthCheck().
		Use(middlewares.LoggingMiddleware).
		BuildServer(settings.ServerAddress)

	// Start the server with graceful shutdown
	serverBuilder.StartServer()
}
