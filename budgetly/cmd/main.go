package main

import (
	"github.com/ZiadMansourM/budgetly/cmd/api"
	"github.com/ZiadMansourM/budgetly/pkg/middlewares"
	"github.com/ZiadMansourM/budgetly/pkg/settings"
)

func main() {
	// Load the application settings (includes logger and environment variables).
	settings, err := settings.Init(".env")
	if err != nil {
		panic(err)
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
