package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log/slog"
)

var DB *gorm.DB

func Connect(dsn string) {
	if DB != nil {
		slog.Debug("Database already connected. Skipping connection.")
		return
	}

	slog.Debug("Attempting to connect to database with DSN:", "dsn", dsn)
	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Debug("Failed to connect database:", "error", err)
		panic("failed to connect database")
	}
	slog.Debug("Database connected successfully.")

	// Migrate the schema
	slog.Debug("Attempting to migrate database schema...")
	err = database.AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{})
	if err != nil {
		slog.Debug("Failed to migrate database:", "error", err)
		panic("failed to migrate database")
	}
	slog.Debug("Database schema migrated successfully.")

	DB = database
}