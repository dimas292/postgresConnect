package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Import PostgreSQL driver for side effects
)

// ConnectPostgres connects to a PostgreSQL database using the provided parameters
func ConnectPostgres(host, port, user, pass, dbname string) (db *sql.DB, err error) {
	// Construct the data source name (DSN) string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

	// Open a connection to the database
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Return the database connection
	return db, nil
}