package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dsn string) {
	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	database.AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{})

	DB = database
}
