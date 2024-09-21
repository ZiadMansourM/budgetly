package users

import (
	"time"

	"github.com/ZiadMansourM/budgetly/pkg/validate"
)

type User struct {
	ID             int       `db:"id"`
	Username       string    `db:"username"`
	Email          string    `db:"email"`
	PasswordHashed string    `db:"password_hashed"`
	CreatedAt      time.Time `db:"created_at"`
}

// UserRequest represents the input data for registering a new user.
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the UserRequest struct.
func (input *UserRequest) Validate() map[string]string {
	// Define validation rules for each field.
	validationFields := validate.ValidationFields{
		"Username": validate.Rules(
			validate.Required,
			validate.Min(3),
			validate.ErrorMessage("Username is required and must be at least 3 characters long"),
		),
		"Email": validate.Rules(
			validate.Required,
			validate.Email,
			validate.ErrorMessage("A valid email address is required"),
		),
		"Password": validate.Rules(
			validate.Required,
			validate.Min(6),
			validate.ErrorMessage("Password must be at least 6 characters long"),
		),
	}

	// Perform validation using the validate package.
	return validate.Validate(*input, validationFields)
}

// UserResponse represents the user data to return in responses (without sensitive info like password).
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ToResponse converts a User (from database) to a UserResponse (for API responses).
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}
}
