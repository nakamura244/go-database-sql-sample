package interfaces

import (
	"errors"
)

type SQLRepository struct {
	SQLhandler
}

type User struct {
	ID    uint   // id
	Email string // email
}

func (repo *SQLRepository) FindUserByID(id uint) (*User, error) {
	const sqlstr = `SELECT ` +
		`id, email ` +
		`FROM users ` +
		`WHERE id = ? `
	row, err := repo.Query(sqlstr, id)
	if err != nil {
		return nil, err
	}
	if !row.Next() {
		return nil, errors.New("failed to row.Next()")
	}
	u := &User{}
	if err = row.Scan(&u.ID, &u.Email); err != nil {
		return nil, err
	}
	err = row.Close()
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (repo *SQLRepository) FindUsers() ([]*User, error) {
	const sqlstr = `SELECT ` +
		`id, email ` +
		`FROM users `
	q, err := repo.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	res := []*User{}
	for q.Next() {
		u := User{}
		err = q.Scan(&u.ID, &u.Email)
		if err != nil {
			return nil, err
		}
		res = append(res, &u)
	}
	err = q.Close()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *SQLRepository) InsertUser(u *User) (uint, error) {
	const sql = `INSERT INTO users ( ` +
		`email ` +
		`) VALUES (?) `
	res, err := repo.Execute(sql, u.Email)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (repo *SQLRepository) InsertUserWithTx(u *User) (uint, error) {
	tx, err := repo.Begin()
	if err != nil {
		return 0, err
	}

	const sql = `INSERT INTO users ( ` +
		`email ` +
		`) VALUES (?) `
	res, err := tx.Execute(sql, u.Email)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	rowsAffect, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if rowsAffect != 1 {
		tx.Rollback()
		return 0, errors.New("rowsAffect != 0 error")
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return uint(lastID), nil
}