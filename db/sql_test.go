package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	iface2 "github.com/nakamura244/databasesql/db/iface"
	"reflect"
	"sync"
	"testing"
)

type Mock struct {
	iface2.SQLAPI
	iface2.TxAPI
	iface2.ResultAPI
	iface2.RowsAPI
	iface2.RowAPI
	// 0 -> no error
	// 1 -> Begin error
	// 2 -> Exec(TX) error
	// 3 -> Commit error
	// 4 -> Rollback error
	// 5 -> Exec(SQL) error .. いらないかも
	// 6 -> Query error
	// 7 -> QueryRow error
	// 8 -> last insert id error
	// 9 -> rows affected error
	// 10 -> scan error
	// 11 -> next error
	// 12 -> close error
	errNo int
	err   error
}

// dummy
type driverResult struct {
	sync.Locker // the *driverConn
	resi        driver.Result
}

// dummy
func (dr driverResult) LastInsertId() (int64, error) {
	return 5, nil
}

// dummy
func (dr driverResult) RowsAffected() (int64, error) {
	return 1, nil
}


func (m Mock) Begin() (*sql.Tx, error) {
	if m.errNo == 1 {
		return &sql.Tx{}, errors.New("begin error")
	}
	return &sql.Tx{}, nil
}

func (m Mock) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.errNo == 2 {
		return driverResult{}, errors.New("exec error")
	}
	return driverResult{}, nil
}

func (m Mock) Commit() error {
	if m.errNo == 3 {
		return errors.New("commit error")
	}
	return nil
}

func (m Mock) Rollback() error {
	if m.errNo == 4 {
		return errors.New("rollback error")
	}
	return nil
}

func (m Mock) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.errNo == 6 {
		return &sql.Rows{}, errors.New("query error")
	}
	return &sql.Rows{}, nil
}

func (m Mock) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.errNo == 7 {
		// &sql.Row{err:} にerrorsをセットしたいができない
		return &sql.Row{}
	}
	return &sql.Row{}
}

func (m Mock) LastInsertId() (int64, error) {
	if m.errNo == 8 {
		return 0, errors.New("last insert id error")
	}
	return 2, nil
}

func (m Mock) RowsAffected() (int64, error) {
	if m.errNo == 9 {
		return 0, errors.New("rows affected error")
	}
	return 1, nil
}

func (m Mock) Scan(dest ...interface{}) error {
	if m.errNo == 10 {
		return errors.New("scan error")
	}
	return nil
}

func (m Mock) Next() bool {
	if m.errNo == 11 {
		return false
	}
	return true
}

func (m Mock) Close() error {
	if m.errNo == 12 {
		return errors.New("close error")
	}
	return nil
}

func TestNewConn(t *testing.T) {
	conn := NewConn()
	expected := "*sql.DB"
	if reflect.TypeOf(conn.Conn).String() != expected {
		t.Errorf("expected  %v, actual %v", expected, reflect.TypeOf(conn.Conn).String())
	}
}

func TestMysql_Begin(t *testing.T) {
	tests := []struct {
		i   iface2.SQLAPI
		r   string
		err error
	}{
		{
			i:Mock{},
			r: "*db.TX",
			err:nil,
		},
		{
			i:Mock{errNo:1},
			r: "",
			err:errors.New("begin error"),
		},
	}

	for i, test := range tests {
		m := Mysql{
			Conn:test.i,
		}
		tx, err := m.Begin()
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err, err)
			}
		} else {
			if reflect.TypeOf(tx).String() != test.r {
				t.Errorf("%d, expected  %v, actual %v", i, test.r, reflect.TypeOf(tx).Name())
			}
		}
	}
}

func TestTX_ExecuteWithTx(t *testing.T) {
	tests := []struct {
		i   iface2.TxAPI
		r   string
		err error
	}{
		{
			i:Mock{},
			r:"db.Result",
			err:nil,
		},
		{
			i:Mock{errNo:2},
			r:"",
			err:errors.New("exec error"),
		},
	}

	for i, test := range tests {
		m := TX{test.i}
		var p interface{}
		res, err := m.Execute("sql query", p)
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		} else {
			if reflect.TypeOf(res).String() != test.r {
				t.Errorf("%d, expected  %v, actual %v", i, test.r, reflect.TypeOf(res).String())
			}
		}
	}
}

func TestTX_Commit(t *testing.T) {
	tests := []struct {
		i   iface2.TxAPI
		err error
	}{
		{
			i:Mock{},
			err:nil,
		},
		{
			i:Mock{errNo:3},
			err:errors.New("commit error"),
		},
	}
	for i, test := range tests {
		m := TX{test.i}
		r := m.Commit()
		if test.err != nil {
			if r.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), r.Error())
			}
		} else {
			if r != nil {
				t.Errorf("%d, expected  %v, actual %v", i, test.err, r)
			}
		}
	}
}

func TestTX_Rollback(t *testing.T) {
	tests := []struct {
		i   iface2.TxAPI
		err error
	}{
		{
			i:Mock{},
			err:nil,
		},
		{
			i:Mock{errNo:4},
			err:errors.New("rollback error"),
		},
	}
	for i, test := range tests {
		m := TX{test.i}
		r := m.Rollback()
		if test.err != nil {
			if r.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), r.Error())
			}
		} else {
			if r != nil {
				t.Errorf("%d, expected  %v, actual %v", i, test.err, r)
			}
		}
	}
}

func TestMysql_Execute(t *testing.T) {
	tests := []struct {
		i   iface2.SQLAPI
		r   string
		err error
	}{
		{
			i:   Mock{errNo: 0},
			r:"db.Result",
			err: nil,
		},
		{
			i:   Mock{errNo: 2},
			r: "",
			err: errors.New("exec error"),
		},
	}
	for i, test := range tests {
		m := Mysql{test.i}
		var p interface{}
		res, err := m.Execute("sql query", p)
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		} else {
			if reflect.TypeOf(res).String() != test.r {
				t.Errorf("%d, expected  %v, actual %v", i, test.r, reflect.TypeOf(res).String())
			}
		}
	}
}

func TestMysql_Query(t *testing.T) {
	tests := []struct {
		i   iface2.SQLAPI
		r   string
		err error
	}{
		{
			i:   Mock{errNo: 0},
			r: "*db.Rows",
			err: nil,
		},
		{
			i:   Mock{errNo: 6},
			r: "*db.Rows",
			err: errors.New("query error"),
		},
	}
	for i, test := range tests {
		m := Mysql{test.i}
		var p interface{}
		res, err := m.Query("sql", p)
		if reflect.TypeOf(res).String() != test.r {
			t.Errorf("%d, expected  %v, actual %v", i, test.r, reflect.TypeOf(res).String())
		}
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		}
	}
}

func TestMysql_QueryRow(t *testing.T) {
	m := Mysql{Mock{}}
	var p interface{}
	res := m.QueryRow("sql", p)
	expected := "*db.Row"
	if reflect.TypeOf(res).String() != expected {
		t.Errorf("expected  %v, actual %v", expected, reflect.TypeOf(res).String())
	}
}

func TestResult_LastInsertId(t *testing.T) {
	tests := []struct {
		i   iface2.ResultAPI
		err error
		id  int64
	}{
		{
			i:   Mock{errNo: 0},
			err: nil,
			id:  2,
		},
		{
			i:   Mock{errNo: 8},
			err: errors.New("last insert id error"),
			id:  0,
		},
	}
	for i, test := range tests {
		m := Result{test.i}
		id, err := m.LastInsertId()
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		}
		if id != test.id {
			t.Errorf("%d, expected  %v, actual %v", i, test.id, id)
		}
	}
}

func TestResult_RowsAffected(t *testing.T) {
	tests := []struct {
		i   iface2.ResultAPI
		err error
		id  int64
	}{
		{
			i:   Mock{errNo: 0},
			err: nil,
			id:  1,
		},
		{
			i:   Mock{errNo: 9},
			err: errors.New("rows affected error"),
			id:  0,
		},
	}
	for i, test := range tests {
		m := Result{test.i}
		id, err := m.RowsAffected()
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		}
		if id != test.id {
			t.Errorf("%d, expected  %v, actual %v", i, test.id, id)
		}
	}
}

func TestRows_Scan(t *testing.T) {
	tests := []struct {
		i   iface2.RowsAPI
		err error
	}{
		{
			i:   Mock{errNo: 0},
			err: nil,
		},
		{
			i:   Mock{errNo: 10},
			err: errors.New("scan error"),
		},
	}
	for i, test := range tests {
		m := Rows{test.i}
		err := m.Scan()
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("%d, expected  %v, actual %v", i, nil, err)
			}
		}
	}
}

func TestRow_Scan(t *testing.T) {
	tests := []struct {
		i   iface2.RowAPI
		err error
	}{
		{
			i:   Mock{errNo: 0},
			err: nil,
		},
		{
			i:   Mock{errNo: 10},
			err: errors.New("scan error"),
		},
	}
	for i, test := range tests {
		m := Row{test.i}
		err := m.Scan()
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("%d, expected  %v, actual %v", i, nil, err)
			}
		}
	}
}

func TestRows_Next(t *testing.T) {
	tests := []struct {
		i iface2.RowsAPI
		r bool
	}{
		{
			i: Mock{errNo: 0},
			r: true,
		},
		{
			i: Mock{errNo: 11},
			r: false,
		},
	}
	for i, test := range tests {
		m := Rows{test.i}
		r := m.Next()
		if r != test.r {
			t.Errorf("%d, expected  %v, actual %v", i, test.r, r)
		}
	}

}

func TestRows_Close(t *testing.T) {
	tests := []struct {
		i   iface2.RowsAPI
		err error
	}{
		{
			i:   Mock{errNo: 0},
			err: nil,
		},
		{
			i:   Mock{errNo: 12},
			err: errors.New("close error"),
		},
	}
	for i, test := range tests {
		m := Rows{test.i}
		err := m.Close()
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err.Error(), err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("%d, expected  %v, actual %v", i, nil, err)
			}
		}
	}
}