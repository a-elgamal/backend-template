package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/mock"

	"alielgamal.com/myservice/internal/db"
)

// DB mock for sql.DB that follows the interface alielgamal.com/myservice/internal/db#DB
type DB struct {
	mock.Mock
}

// NewDB Creates a new mocked DB
func NewDB(t *testing.T) *DB {
	m := DB{}
	m.Test(t)

	return &m
}

// PingContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *DB) PingContext(ctx context.Context) error {
	args := m.MethodCalled("PingContext", ctx)
	return args.Error(0)
}

// QueryRowContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *DB) QueryRowContext(ctx context.Context, query string, queryArgs ...any) *sql.Row {
	var params []interface{}
	params = append(params, ctx)
	params = append(params, query)
	params = append(params, queryArgs...)
	args := m.MethodCalled("QueryRowContext", params...)
	return args.Get(0).(*sql.Row)
}

// ExecContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *DB) ExecContext(ctx context.Context, query string, queryArgs ...any) (sql.Result, error) {
	var params []interface{}
	params = append(params, ctx)
	params = append(params, query)
	params = append(params, queryArgs...)
	args := m.MethodCalled("ExecContext", params...)
	return args.Get(0).(sql.Result), args.Error(1)
}

// QueryContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *DB) QueryContext(ctx context.Context, query string, queryArgs ...any) (*sql.Rows, error) {
	var params []interface{}
	params = append(params, ctx)
	params = append(params, query)
	params = append(params, queryArgs...)
	args := m.MethodCalled("QueryContext", params...)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

// BeginTx mirrors sql.DB#BeginTx
func (m *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (db.Tx, error) {
	args := m.MethodCalled("BeginTx", ctx, opts)
	return args.Get(0).(db.Tx), args.Error(1)
}

// Tx mock for sql.Tx that follows the interface alielgamal.com/myservice/internal/db#Tx
type Tx struct {
	mock.Mock
}

// NewTx Creates a new mocked DB
func NewTx(t *testing.T) *Tx {
	m := Tx{}
	m.Test(t)

	return &m
}

// PingContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *Tx) PingContext(ctx context.Context) error {
	args := m.MethodCalled("PingContext", ctx)
	return args.Error(0)
}

// QueryRowContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *Tx) QueryRowContext(ctx context.Context, query string, queryArgs ...any) *sql.Row {
	var params []interface{}
	params = append(params, ctx)
	params = append(params, query)
	params = append(params, queryArgs...)
	args := m.MethodCalled("QueryRowContext", params...)
	return args.Get(0).(*sql.Row)
}

// ExecContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *Tx) ExecContext(ctx context.Context, query string, queryArgs ...any) (sql.Result, error) {
	var params []interface{}
	params = append(params, ctx)
	params = append(params, query)
	params = append(params, queryArgs...)
	args := m.MethodCalled("ExecContext", params...)
	return args.Get(0).(sql.Result), args.Error(1)
}

// QueryContext follows the alielgamal.com/myservice/internal/db#DB interface
func (m *Tx) QueryContext(ctx context.Context, query string, queryArgs ...any) (*sql.Rows, error) {
	var params []interface{}
	params = append(params, ctx)
	params = append(params, query)
	params = append(params, queryArgs...)
	args := m.MethodCalled("QueryContext", params...)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

// Commit mirrors sql.Tx#Commit
func (m *Tx) Commit() error {
	args := m.MethodCalled("Commit")
	return args.Error(0)
}

// Rollback mirrors sql.Tx#Rollback
func (m *Tx) Rollback() error {
	args := m.MethodCalled("Rollback")
	return args.Error(0)
}
