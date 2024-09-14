package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZiadMansourM/budgetly/internal/handlers"
	"github.com/ZiadMansourM/budgetly/internal/models"
	"github.com/ZiadMansourM/budgetly/internal/services"
	"github.com/ZiadMansourM/budgetly/pkg/db"
	"github.com/ZiadMansourM/budgetly/utils"
	"github.com/jmoiron/sqlx"
)

type ServerBuilder struct {
	dbType     string
	dbConn     string
	dbPool     *sqlx.DB
	router     *http.ServeMux
	httpServer *http.Server
}

// NewServerBuilder initializes the ServerBuilder
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{
		router: http.NewServeMux(),
	}
}

// WithDatabase sets up the database connection
func (b *ServerBuilder) WithDatabase(dbType, dbConn string) *ServerBuilder {
	pool, err := db.OpenDB(dbType, dbConn)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}
	b.dbType = dbType
	b.dbConn = dbConn
	b.dbPool = pool
	return b
}

// WithUserApp sets up the entire User application (model, service, handler, and routes)
func (b *ServerBuilder) WithUserApp() *ServerBuilder {
	// Initialize User Model
	userModel := &models.UserModel{DB: b.dbPool}

	// Initialize User Service
	userService := &services.UserService{UserRepo: userModel}

	// Initialize User Handler and Register Routes
	userHandler := &handlers.UserHandler{UserService: userService}
	userHandler.RegisterRoutes(b.router)

	return b
}

// WithHealthCheck adds a health check route
func (b *ServerBuilder) WithHealthCheck() *ServerBuilder {
	b.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJson(w, http.StatusOK, map[string]string{"message": "API is healthy"})
	})
	return b
}

// BuildServer builds the HTTP server with graceful shutdown
func (b *ServerBuilder) BuildServer(addr string) *ServerBuilder {
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", b.router))

	b.httpServer = &http.Server{
		Addr:    addr,
		Handler: v1,
	}

	return b
}

// StartServer starts the HTTP server with graceful shutdown
func (b *ServerBuilder) StartServer() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Gracefully shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := b.httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("Could not shutdown server: %v\n", err)
		}

		log.Println("Server Exited Properly")
	}()

	log.Printf("Server Listening on %s\n", b.httpServer.Addr)
	if err := b.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", b.httpServer.Addr, err)
	}
}
