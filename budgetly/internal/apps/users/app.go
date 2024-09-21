package users

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

// NewUserApp creates a new user application with the provided database connection
func NewUserApp(db *sqlx.DB, logger *slog.Logger) *userHandler {
	userModel := newUserModel(db, logger)
	userService := newUserService(userModel, logger)
	return &userHandler{
		userService: userService,
		logger:      logger,
	}
}
