package users

import (
	"time"

	"github.com/ZiadMansourM/budgetly/pkg/validate"
)

type userService struct {
	userRepo *userModel
}

func newUserService(userRepo *userModel) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

// RegisterUser handles user registration
func (s *userService) register(username, email, password string) (*User, error) {
	// Create a struct to validate the input data.
	userInput := struct {
		Username string
		Email    string
		Password string
	}{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Define validation rules for each field, including custom error messages.
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
	if validationErrors := validate.Validate(userInput, validationFields); len(validationErrors) > 0 {
		// Return the validation errors as a ValidationError.
		return nil, &validate.ValidationError{Errors: validationErrors}
	}

	// Create a new user object
	user := &User{
		Username:       username,
		Email:          email,
		PasswordHashed: hashPassword(password), // Simplified password hashing for this example
		CreatedAt:      time.Now(),
	}

	// Save the user to the database using the UserModel
	userID, err := s.userRepo.create(user)
	if err != nil {
		return nil, err
	}

	user.ID = userID
	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) getByID(id int) (*User, error) {
	// Retrieve the user from the database using the UserModel
	return s.userRepo.getByID(id)
}

// hashPassword is a simplified password hashing function (use a proper password hashing library in production)
func hashPassword(password string) string {
	// Placeholder hash function for demonstration
	return "hashed_" + password
}
