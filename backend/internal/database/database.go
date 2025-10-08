// Package database provides database connection and initialization for the football pool application.
package database

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database represents a database connection with GORM ORM.
type Database struct {
	db *gorm.DB
}

// New creates a new Database connection with the provided database type and DSN.
func New(dbType, dsn string) (*Database, error) {
	slog.Debug("Attempting to connect to database", "type", dbType, "dsn", dsn)

	var database *gorm.DB
	var err error

	switch dbType {
	case "sqlite":
		database, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgres":
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

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
	err := d.db.AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{}, &Week{})
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

// WeekHasGames checks if a given week has any games in the database.
func (d *Database) WeekHasGames(season int, week int) (bool, error) {
	var count int64
	err := d.db.Model(&Game{}).Where("season = ? AND week = ?", season, week).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
