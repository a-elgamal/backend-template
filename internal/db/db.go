// Package db contains utility methods for accessing SQL database
package db

import (
	"context"
	"database/sql"
)

type querable interface {
	// QueryRowContext mirrors sql.DB#QueryRowContext
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	// ExecuteContext mirrors sql.DB#QueryContext
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// QueryContext mirrors sql.DB#QueryContext
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// DB interface to enable replacing the implementation of sql.DB instances. Helpful for testing
type DB interface {
	querable

	//PingContext mirrors sql.DB#PingContext
	PingContext(ctx context.Context) error

	// BeginTx mirrors sql.DB#BeginTx
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

// Tx interface to enable replacing implementation of sql.Tx instances. Helpful for testing
type Tx interface {
	querable

	// Commit mirrors sql.Tx#Commit
	Commit() error

	// Rollback mirrors sql.Tx#Rollback
	Rollback() error
}

// SQLDB Wraps sql.DB and implements internal DB interface
type SQLDB struct {
	DB *sql.DB
}

// PingContext mirrors sql.DB#PingContext
func (s *SQLDB) PingContext(ctx context.Context) error {
	return s.DB.PingContext(ctx)
}

// BeginTx mirrors sql.DB#BeginTx
func (s *SQLDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	return s.DB.BeginTx(ctx, opts)
}

// QueryRowContext mirrors sql.DB#QueryRowContext
func (s *SQLDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return s.DB.QueryRowContext(ctx, query, args...)
}

// ExecContext mirrors sql.DB#QueryContext
func (s *SQLDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.DB.ExecContext(ctx, query, args...)
}

// QueryContext mirrors sql.DB#QueryContext
func (s *SQLDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return s.DB.QueryContext(ctx, query, args...)
}
