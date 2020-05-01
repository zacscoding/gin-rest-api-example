package repository

import (
	"gin-rest-api-example/database"
	"log"
)

var (
	UserRepo UserRepository
)

func Init() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	//UserRepo = NewUserRepository(db)

}
