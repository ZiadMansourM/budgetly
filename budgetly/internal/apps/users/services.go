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
func (s *userService) register(input UserRequest) (*UserResponse, error) {
	// Validate the user input.
	if validationErrors := input.Validate(); len(validationErrors) > 0 {
		// Return the validation errors as a ValidationError.
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
		return nil, err
	}

	user.ID = userID

	// Convert the user entity to UserResponse and return
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
