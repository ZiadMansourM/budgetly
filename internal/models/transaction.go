package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// TransactionModel wraps the database connection pool using sqlx
type TransactionModel struct {
	DB *sqlx.DB
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	AccountID   int       `db:"account_id"`
	Amount      int       `db:"amount"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

// Create inserts a new transaction into the database and returns the inserted transaction's ID
func (m *TransactionModel) Create(t *Transaction) (int, error) {
	// Use NamedExec for more readable query with named parameters
	query := `INSERT INTO transactions (user_id, account_id, amount, description, created_at) 
	VALUES (:user_id, :account_id, :amount, :description, :created_at) 
	RETURNING id`

	// Ensure CreatedAt is set to the current time
	t.CreatedAt = time.Now()

	// Execute the query and return the ID
	rows, err := m.DB.NamedQuery(query, t)
	if err != nil {
		return 0, err
	}

	if rows.Next() {
		err = rows.Scan(&t.ID)
		if err != nil {
			return 0, err
		}
	}

	return t.ID, nil
}

// GetByUserID returns all transactions for a userId
func (m *TransactionModel) GetByUserID(userID int) ([]*Transaction, error) {
	// Use Select to map all rows into a slice of Transaction structs
	query := `SELECT id, user_id, account_id, amount, description, created_at 
	FROM transactions WHERE user_id = $1`

	var transactions []*Transaction
	err := m.DB.Select(&transactions, query, userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
