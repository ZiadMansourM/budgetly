package users

import "github.com/jmoiron/sqlx"

// NewUserApp creates a new user application with the provided database connection
func NewUserApp(db *sqlx.DB) *userHandler {
	userModel := newUserModel(db)
	userService := newUserService(userModel)
	return &userHandler{
		userService: userService,
	}
}
