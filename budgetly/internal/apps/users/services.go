package users

import (
	"log/slog"
	"time"

	"github.com/ZiadMansourM/budgetly/pkg/validate"
)

type userService struct {
	userRepo *userModel
	logger   *slog.Logger
}

func newUserService(userRepo *userModel, logger *slog.Logger) *userService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// RegisterUser handles user registration
func (s *userService) register(input UserRequest) (*UserResponse, error) {
	// Validate the user input.
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		// Return the validation errors as a ValidationError.
		s.logger.Warn("User registration validation failed", "errors", validationErrors)
		return nil, &validate.ValidationError{Errors: validationErrors}
	}

	// Create a new user object
	user := &User{
		Username:       input.Username,
		Email:          input.Email,
		PasswordHashed: hashPassword(input.Password), // Simplified password hashing for this example
		CreatedAt:      time.Now(),
	}

	// Save the user to the database using the UserModel
	userID, err := s.userRepo.create(user)
	if err != nil {
		s.logger.Error("Error creating user", "error", err)
		return nil, err
	}

	s.logger.Debug("User created successfully", "id", userID)
	user.ID = userID
	return user.ToResponse(), nil
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
