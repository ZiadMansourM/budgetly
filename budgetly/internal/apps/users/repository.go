// This is the model layer, responsible for interacting
// with the database and returning data to the service layer
package users

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// userModel wraps the database connection pool using sqlx
type userModel struct {
	DB *sqlx.DB
}

func newUserModel(db *sqlx.DB) *userModel {
	return &userModel{
		DB: db,
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
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&u.ID)
		if err != nil {
			return 0, err
		}
	}

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
		return nil, err
	}

	return u, nil
}
