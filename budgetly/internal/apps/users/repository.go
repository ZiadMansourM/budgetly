// This is the model layer, responsible for interacting
// with the database and returning data to the service layer
package users

import (
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

// userModel wraps the database connection pool using sqlx
type userModel struct {
	DB     *sqlx.DB
	logger *slog.Logger
}

func newUserModel(db *sqlx.DB, logger *slog.Logger) *userModel {
	return &userModel{
		DB:     db,
		logger: logger,
	}
}

// Create inserts a new user into the database and returns the inserted user's ID
func (m *userModel) create(u *User) (int, error) {
	// Use NamedExec for more readable query with named parameters
	query := `INSERT INTO users (username, email, password_hashed, created_at) 
	VALUES (:username, :email, :password_hashed, :created_at) 
	RETURNING id`

	// Ensure the user has a valid CreatedAt value
	u.CreatedAt = time.Now()

	// Execute the query and return the ID
	rows, err := m.DB.NamedQuery(query, u)
	if err != nil {
		m.logger.Error("Error inserting user", "error", err)
		// Do not expose the error to the client add new error Internal server error only
		return 0, ErrInternalServer
	}

	if rows.Next() {
		err = rows.Scan(&u.ID)
		if err != nil {
			m.logger.Error("Error scanning user ID", "error", err)
			return 0, ErrInternalServer
		}
	}

	m.logger.Debug("User created successfully", "id", u.ID)
	return u.ID, nil
}

// GetByID returns a user by ID
func (m *userModel) getByID(id int) (*User, error) {
	// Use Get to map a single result to a struct
	query := `SELECT id, username, email, password_hashed, created_at 
	FROM users WHERE id = $1`

	u := &User{}
	err := m.DB.Get(u, query, id)
	if err != nil {
		m.logger.Error("Error getting user by ID", "error", err)
		return nil, ErrInternalServer
	}

	m.logger.Debug("User retrieved successfully", "id", u.ID)
	return u, nil
}
