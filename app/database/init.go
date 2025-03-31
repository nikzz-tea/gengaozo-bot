package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Init() {
	var db, err = gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{}, &Map{})
}
