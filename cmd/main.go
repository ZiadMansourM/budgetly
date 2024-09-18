package main

import "github.com/ZiadMansourM/budgetly/cmd/api"

func main() {
	// Use the builder to assemble the server with plug-and-play apps
	// E.g. WithUserApp which encapsulate all its components (model, service, handler, routes).
	serverBuilder := api.NewServerBuilder().
		WithDatabase("postgres", "postgres://budgetly:P@ssw0rd@127.0.0.1:5432/budgetly?sslmode=disable").
		WithUserApp().
		WithHealthCheck().
		BuildServer("0.0.0.0:8080")

	// Start the server with graceful shutdown
	serverBuilder.StartServer()
}
