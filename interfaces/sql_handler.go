package interfaces

type SQLhandler interface {
	Execute(string, ...interface{}) (Result, error)
	Query(string, ...interface{}) (Rows, error)
	QueryRow(string, ...interface{}) Row
	Begin() (Tx, error)
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Rows interface {
	Scan(...interface{}) error
	Next() bool
	Close() error
}

type Row interface {
	Scan(...interface{}) error
}

type Tx interface {
	Execute(string, ...interface{}) (Result, error)
	Commit() error
	Rollback() error
}
