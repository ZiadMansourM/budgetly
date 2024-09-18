package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// OpenDB creates and returns a new database connection pool
func OpenDB(dbType, dbConn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(dbType, dbConn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	return db, nil
}
