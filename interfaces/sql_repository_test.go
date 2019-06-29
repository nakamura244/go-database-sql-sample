package interfaces

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
	"time"
)

type Mock struct {
	SQLhandler
	// 0 -> no error
	// 1 -> Query error
	// 2 -> Scan error
	// 3 -> Close error
	// 4 -> Next error
	// 5 -> Execute error
	// 6 -> LastInsertId error
	// 7 -> RowsAffected error
	// 8 -> Begin error
	// 9 -> Commit error
	// 10 -> Rollback error
	errNo   int
	nextCnt int
	err     error
}

var Gi int


func (m Mock) Query(string, ...interface{}) (Rows, error) {
	if m.errNo == 1 {
		return Mock{}, errors.New("error query")
	}
	return &m, nil
}


func (m Mock) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	if m.errNo == 2 {
		return errors.New("error scan")
	}
	for _, des := range dest {
		switch s := des.(type) {
		case *string:
			*s = "test string"
		case *int:
			*s = 2
		case *uint:
			*s = 2
		case *sql.NullString:
			*s = sql.NullString{String: "test string", Valid: true}
		case *sql.NullInt64:
			*s = sql.NullInt64{Int64: 3, Valid: true}
		case *int8:
			*s = int8(4)
		case *bool:
			*s = true
		case *sql.NullBool:
			*s = sql.NullBool{Bool: true, Valid: true}
		case *mysql.NullTime:
			t, _ := time.Parse("2006-01-02", "2014-12-31")
			*s = mysql.NullTime{Time: t, Valid: true}
		}
	}
	return nil
}


func (m Mock) Close() error {
	if m.errNo == 3 {
		return errors.New("error close")
	}
	return nil
}

func (m Mock) Next() bool {
	if m.errNo == 4 {
		return false
	}
	if m.nextCnt == 1 {
		Gi++
	}
	if Gi > (m.nextCnt + 1) {
		return false
	}
	return true
}


func (m Mock) Execute(string, ...interface{}) (Result, error) {
	if m.errNo == 5 {
		return &Mock{}, errors.New("error execute")
	}
	return &m, nil
}

func (m Mock) LastInsertId() (int64, error) {
	if m.errNo == 6 {
		return int64(0), errors.New("error last insert id")
	}
	return int64(1), nil
}

func (m Mock) RowsAffected() (int64, error) {
	if m.errNo == 7 {
		return int64(0), errors.New("error row affected")
	}
	return int64(1), nil
}

func (m Mock) Begin() (Tx, error) {
	if m.errNo == 8 {
		return &Mock{}, errors.New("error begin")
	}
	return &m, nil
}

func (m Mock) Commit() error {
	if m.errNo == 9 {
		return errors.New("error commit")
	}
	return nil
}

func (m Mock) Rollback() error {
	if m.errNo == 10 {
		return errors.New("error rollback")
	}
	return nil
}

func TestSQLRepository_FindUserByID(t *testing.T) {
	tests := []struct{
		m   SQLhandler
		err error
	} {
		{
			m:Mock{},
			err:nil,
		},
		{
			m:   Mock{errNo: 1},
			err: errors.New("error query"),
		},
		{
			m:   Mock{errNo: 4},
			err: errors.New("failed to row.Next()"),
		},
		{
			m:   Mock{errNo: 2},
			err: errors.New("error scan"),
		},
		{
			m:   Mock{errNo: 3},
			err: errors.New("error close"),
		},
	}

	for i, test := range tests {
		m := SQLRepository{test.m}
		r, err := m.FindUserByID(1)
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err, err)
			}
		} else {
			expected := &User{
				ID:2,
				Email:"test string",
			}
			if reflect.DeepEqual(r, expected) == false {
				t.Errorf("expected %+v, got %+v", expected, r)
			}
		}
	}
}

func TestSQLRepository_FindUsers(t *testing.T) {
	tests := []struct{
		m   SQLhandler
		err error
	} {
		{
			m:   Mock{nextCnt: 1},
			err: nil,
		},
		{
			m:   Mock{errNo: 1},
			err: errors.New("error query"),
		},
		{
			m:   Mock{nextCnt: 2, errNo: 2},
			err: errors.New("error scan"),
		},
		{
			m:   Mock{nextCnt: 1, errNo: 3},
			err: errors.New("error close"),
		},
	}
	for i, test := range tests {
		Gi = 0
		m := SQLRepository{test.m}
		r, err := m.FindUsers()
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err, err)
			}
		} else {
			expected := &User{
				ID:2,
				Email:"test string",
			}
			if reflect.DeepEqual(r[0], expected) == false {
				t.Errorf("expected %+v, got %+v", expected, r[0])
			}
		}
	}
}

func TestSQLRepository_InsertUser(t *testing.T) {
	tests := []struct{
		m SQLhandler
		err error
	} {
		{
			m:   Mock{errNo: 0},
			err: nil,
		},
		{
			m:   Mock{errNo: 5},
			err: errors.New("error execute"),
		},
		{
			m:   Mock{errNo: 6},
			err: errors.New("error last insert id"),
		},
	}
	for i, test := range tests {
		Gi = 0
		m := SQLRepository{
			test.m,
		}
		r, err := m.InsertUser(&User{})
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err, err)
			}
		} else {
			expected := uint(1)
			if expected != r {
				t.Errorf("expected %+v, got %+v", expected, r)
			}
		}
	}
}

func TestSQLRepository_InsertUserWithTx(t *testing.T) {
	tests := []struct{
		m SQLhandler
		err error
	} {
		{
			m:   Mock{errNo: 0},
			err: nil,
		},
		{
			m:   Mock{errNo: 5},
			err: errors.New("error execute"),
		},
		{
			m:   Mock{errNo: 6},
			err: errors.New("error last insert id"),
		},
		{
			m:   Mock{errNo: 7},
			err: errors.New("error row affected"),
		},
		{
			m:   Mock{errNo: 9},
			err: errors.New("error commit"),
		},
	}
	for i, test := range tests {
		Gi = 0
		m := SQLRepository{
			test.m,
		}
		r, err := m.InsertUserWithTx(&User{})
		if test.err != nil {
			if err.Error() != test.err.Error() {
				t.Errorf("%d, expected  %v, actual %v", i, test.err, err)
			}
		} else {
			expected := uint(1)
			if expected != r {
				t.Errorf("expected %+v, got %+v", expected, r)
			}
		}
	}
}