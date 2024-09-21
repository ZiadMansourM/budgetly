package users

import "github.com/ZiadMansourM/budgetly/pkg/validate"

// UserInput represents the input data for registering a new user.
type UserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the UserInput struct.
func (input *UserInput) Validate() map[string]string {
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
