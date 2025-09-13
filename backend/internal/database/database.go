// Package database provides database connection and initialization for the football pool application.
package database

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database represents a database connection with GORM ORM.
type Database struct {
	db *gorm.DB
}

// New creates a new Database connection with the provided DSN.
func New(dsn string) (*Database, error) {
	slog.Debug("Attempting to connect to database with DSN:", "dsn", dsn)
	database, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Debug("Failed to connect database:", "error", err)
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	slog.Debug("Database connected successfully.")

	db := &Database{db: database}

	// Migrate the schema
	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// Migrate performs database schema migrations for all models.
func (d *Database) Migrate() error {
	slog.Debug("Attempting to migrate database schema...")
	err := d.db.AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{})
	if err != nil {
		slog.Debug("Failed to migrate database:", "error", err)
		return err
	}
	slog.Debug("Database schema migrated successfully.")
	return nil
}

// GetDB returns the underlying GORM database connection.
func (d *Database) GetDB() *gorm.DB {
	return d.db
}
