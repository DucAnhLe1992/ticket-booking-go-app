package store

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // optional: uncomment when integrating Postgres
)

// NewPostgres connects to Postgres and returns *sql.DB.
// Example DSN: postgres://user:pw@localhost:5432/dbname?sslmode=disable
func NewPostgres(dsn string) (*sql.DB, error) {
	// NOTE: the pq/pgx driver must be added in go.mod if you call this function
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// quick ping
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}
