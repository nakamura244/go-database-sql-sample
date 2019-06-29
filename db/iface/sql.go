package iface

import (
	"context"
	"database/sql"
)

// SQLAPI is interface
type SQLAPI interface {
	//Conn(ctx context.Context)  (*sql.Conn, error)
	//Driver() driver.Driver
	//Ping() error
	//PingContext(ctx context.Context) error
	Begin() (*sql.Tx, error)
	//BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	//Prepare(query string) (*sql.Stmt, error)
	//PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	//ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	//QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	//QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	//SetConnMaxLifetime(d time.Duration)
	//SetMaxIdleConns(n int)
	//SetMaxOpenConns(n int)
	//Stats() sql.DBStats
	//Close() error
}

// TxAPI is interface
type TxAPI interface {
	Commit() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Rollback() error
	Stmt(stmt *sql.Stmt) *sql.Stmt
	StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt
}

// ResultAPI is interface
type ResultAPI interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// RowsAPI is interface
type RowsAPI interface {
	Scan(dest ...interface{}) error
	Next() bool
	Close() error
}

// RowAPI is interface
type RowAPI interface {
	Scan(dest ...interface{}) error
}
