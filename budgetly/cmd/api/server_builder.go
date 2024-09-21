package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZiadMansourM/budgetly/internal/apps/users"
	"github.com/ZiadMansourM/budgetly/pkg/db"
	"github.com/ZiadMansourM/budgetly/utils"
	"github.com/jmoiron/sqlx"
)

type serverBuilder struct {
	dbType      string
	dbConn      string
	dbPool      *sqlx.DB
	router      *http.ServeMux
	httpServer  *http.Server
	middlewares []func(http.Handler) http.Handler
	logger      *slog.Logger
}

// NewServerBuilder initializes the serverBuilder
func NewServerBuilder(logger *slog.Logger) *serverBuilder {
	return &serverBuilder{
		router:      http.NewServeMux(),
		middlewares: []func(http.Handler) http.Handler{},
		logger:      logger,
	}
}

// WithDatabase sets up the database connection
func (b *serverBuilder) WithDatabase(dbType, dbConn string) *serverBuilder {
	pool, err := db.OpenDB(dbType, dbConn)
	if err != nil {
		slog.Error("Error opening database", "error", err)
	}
	b.dbType = dbType
	b.dbConn = dbConn
	b.dbPool = pool
	return b
}

// WithUserApp sets up the entire User application (model, service, handler, and routes)
func (b *serverBuilder) WithUserApp() *serverBuilder {
	userHandler := users.NewUserApp(b.dbPool, b.logger)
	userHandler.RegisterRoutes(b.router)
	return b
}

// WithHealthCheck adds a health check route
func (b *serverBuilder) WithHealthCheck() *serverBuilder {
	b.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJson(w, http.StatusOK, map[string]string{"message": "API is healthy"})
	})
	return b
}

// AddMiddleware adds middleware to the ordered middleware stack
func (b *serverBuilder) Use(mw func(http.Handler) http.Handler) *serverBuilder {
	b.middlewares = append(b.middlewares, mw)
	return b
}

// applyMiddlewares applies all registered middlewares in order
func (b *serverBuilder) applyMiddlewares(handler http.Handler) http.Handler {
	for i := len(b.middlewares) - 1; i >= 0; i-- {
		// Apply them in reverse order to maintain insertion order
		handler = b.middlewares[i](handler)
	}
	return handler
}

// BuildServer builds the HTTP server with graceful shutdown
func (b *serverBuilder) BuildServer(addr string) *serverBuilder {
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", b.router))

	b.httpServer = &http.Server{
		Addr:    addr,
		Handler: b.applyMiddlewares(v1),
	}

	return b
}

// StartServer starts the HTTP server with graceful shutdown
func (b *serverBuilder) StartServer() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("")
		slog.Info("Gracefully shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := b.httpServer.Shutdown(ctx); err != nil {
			slog.Error("Could not shutdown server", "error", err)
		}

		slog.Info("Server Exited Properly")
	}()

	slog.Info("Server Listening on", "address", b.httpServer.Addr)
	if err := b.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Could not listen on", "address", b.httpServer.Addr, "error", err)
	}
}
