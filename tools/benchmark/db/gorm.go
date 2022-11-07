package db

import (
	"benchmark/types"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	var err error
	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	db.AutoMigrate(&types.Job{})
}

func Get() *gorm.DB {
	if db == nil {
		Connect()
	}

	return db
}
