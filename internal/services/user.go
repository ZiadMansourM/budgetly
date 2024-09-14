package services

import (
	"errors"
	"time"

	"github.com/ZiadMansourM/budgetly/internal/models"
)

type UserService struct {
	UserRepo *models.UserModel
}

// RegisterUser handles user registration
func (s *UserService) Register(username, email, password string) (*models.User, error) {
	// Validate the input
	if username == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	// Create a new user object
	user := &models.User{
		Username:       username,
		Email:          email,
		PasswordHashed: hashPassword(password), // Simplified password hashing for this example
		CreatedAt:      time.Now(),
	}

	// Save the user to the database using the UserModel
	userID, err := s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}

	user.ID = userID
	return user, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetByID(id int) (*models.User, error) {
	// Retrieve the user from the database using the UserModel
	return s.UserRepo.GetByID(id)
}

// hashPassword is a simplified password hashing function (use a proper password hashing library in production)
func hashPassword(password string) string {
	// Placeholder hash function for demonstration
	return "hashed_" + password
}
