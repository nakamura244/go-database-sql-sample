package db

import (
	"database/sql"
	"log"

	"github.com/nakamura244/databasesql/db/iface"
	"github.com/nakamura244/databasesql/interfaces"

	// Mysql Driver
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Conn iface.SQLAPI
}

// Result is struct
type Result struct {
	Result iface.ResultAPI
}

// Rows is struct
type Rows struct {
	Rows iface.RowsAPI
}

// Row is struct for sql.Row
type Row struct {
	Row iface.RowAPI
}

// TX is struct for tx
type TX struct {
	Tx iface.TxAPI
}

func NewConn() *Mysql {
	conn, err := sql.Open("mysql", "vagrant:vagrant@tcp(192.168.33.10:3306)/geenie2")
	if err != nil {
		log.Fatal(err)
	}
	return &Mysql{Conn:conn}
}

// Begin is transaction begin
func (m *Mysql) Begin() (interfaces.Tx, error) {
	tx, err := m.Conn.Begin()
	if err != nil {
		return nil, err
	}
	sqlTx := new(TX)
	sqlTx.Tx = tx
	return sqlTx, nil
}

// Execute is execute sql
func (tx TX) Execute(statement string, args ...interface{}) (interfaces.Result, error) {
	res := Result{}
	result, err := tx.Tx.Exec(statement, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, nil
}

// Commit is transaction commit
func (tx TX) Commit() error {
	return tx.Tx.Commit()
}

// Rollback is transaction rollback
func (tx TX) Rollback() error {
	return tx.Tx.Rollback()
}

// Execute is exe to db
func (m *Mysql) Execute(statement string, args ...interface{}) (interfaces.Result, error) {
	res := Result{}
	result, err := m.Conn.Exec(statement, args...)
	if err != nil {
		return res, err
	}
	res.Result = result
	return res, nil
}

// Query is query to db
func (m *Mysql) Query(statement string, args ...interface{}) (interfaces.Rows, error) {
	rows, err := m.Conn.Query(statement, args...)
	if err != nil {
		return new(Rows), err
	}
	row := new(Rows)
	row.Rows = rows
	return row, nil
}

func (m *Mysql) QueryRow(statement string, args ...interface{}) interfaces.Row {
	r := m.Conn.QueryRow(statement, args...)
	row := new(Row)
	row.Row = r
	return row
}

// LastInsertId is get last insert id
func (r Result) LastInsertId() (int64, error) {
	return r.Result.LastInsertId()
}

// RowsAffected is get affect rows
func (r Result) RowsAffected() (int64, error) {
	return r.Result.RowsAffected()
}

// Scan is mapping for SQLRows
func (r Rows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

// Scan is mapping for SQLRow
func (r Row) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}

// Next is next row
func (r Rows) Next() bool {
	return r.Rows.Next()
}

// Close is close rows
func (r Rows) Close() error {
	return r.Rows.Close()
}
