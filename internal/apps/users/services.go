package users

import (
	"errors"
	"time"
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
	// Validate the input
	if username == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
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
