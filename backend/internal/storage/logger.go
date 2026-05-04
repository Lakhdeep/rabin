package storage

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// LoggingDB wraps sql.DB to add query logging
type LoggingDB struct {
	*sql.DB
}

// QueryContext executes a query with logging
func (db *LoggingDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := db.DB.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	log.Printf("[DB Query] duration=%v query=%q args=%v error=%v", duration, query, args, err)

	return rows, err
}

// QueryRowContext executes a query that returns a single row with logging
func (db *LoggingDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := db.DB.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	log.Printf("[DB QueryRow] duration=%v query=%q args=%v", duration, query, args)

	return row
}

// ExecContext executes a query without returning rows with logging
func (db *LoggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := db.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	log.Printf("[DB Exec] duration=%v query=%q args=%v error=%v", duration, query, args, err)

	return result, err
}

// NewLoggingDB wraps a sql.DB with logging capabilities
func NewLoggingDB(db *sql.DB) *LoggingDB {
	return &LoggingDB{DB: db}
}
