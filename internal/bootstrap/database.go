package bootstrap

import (
	"golang-boilerplate-example/database"
	"golang-boilerplate-example/module/note"
	"log"
	"os/user"

	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&note.Note{}, &user.User{})
	return db
}
