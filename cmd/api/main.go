package main

func main() {
	// Use the builder to assemble the server with plug-and-play apps
	// E.g. WithUserApp which encapsulate all its components (model, service, handler, routes).
	serverBuilder := NewServerBuilder().
		WithDatabase("postgres", "postgres://budgetly:p@ssw0rd@127.0.0.1:5432/budgetly?sslmode=disable").
		WithUserApp().
		WithHealthCheck().
		BuildServer(":8080")

	// Start the server with graceful shutdown
	serverBuilder.StartServer()
}
