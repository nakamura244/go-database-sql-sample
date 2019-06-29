package repository

import (
	"github.com/nakamura244/databasesql/interfaces"
)

type DBRepository interface {
	FindUserByID(id uint) (*interfaces.User, error)
	FindUsers() ([]*interfaces.User, error)
	InsertUser(u *interfaces.User) (uint, error)
	InsertUserWithTx(u *interfaces.User) (uint, error)
}