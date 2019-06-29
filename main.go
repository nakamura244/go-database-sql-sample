package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/nakamura244/databasesql/db"
	"github.com/nakamura244/databasesql/interfaces"
	"github.com/nakamura244/databasesql/repository"
)


func main() {
	conn := db.NewConn()

	var u *interfaces.User
	var sqlRepo repository.DBRepository
	sqlRepo = &interfaces.SQLRepository{
		SQLhandler:conn,
	}

	// insert
	u = &interfaces.User{Email:"email@example.com"}
	id, err := sqlRepo.InsertUser(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)

	// insert with tx
	u = &interfaces.User{Email:"email@example.com"}
	id, err = sqlRepo.InsertUserWithTx(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)


	// select
	user, err := sqlRepo.FindUserByID(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)


}
