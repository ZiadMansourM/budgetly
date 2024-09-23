package users

import (
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// NewUserApp creates a new user application with the provided database connection
func NewUserApp(db *sqlx.DB, logger *slog.Logger, router *http.ServeMux) {
	userModel := newUserModel(db, logger)
	userService := newUserService(userModel, logger)
	newUserHandler(userService, logger, router)
}
