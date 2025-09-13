// Package main is the entry point for the football pool backend server.
package main

import (
	"log/slog"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/david/football-pool/internal/database"
	"github.com/david/football-pool/internal/server"
)

func main() {
	var level slog.Level
	switch os.Getenv("FOOTBALL_POOL_LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default to Info level
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))

	dsn := "football-pool.db"
	if dsnEnv, ok := os.LookupEnv("FOOTBALL_POOL_DSN"); ok {
		dsn = dsnEnv
	}
	slog.Info("Connecting to database")
	db, err := database.New(dsn)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Check if admin user exists, if not, create it
	var adminUser database.User
	result := db.GetDB().Where("email = ?", "admin@test.com").First(&adminUser)
	if result.Error != nil {
		if err := handleAdminUserNotFound(result.Error, db); err != nil {
			slog.Error("Error handling admin user creation", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Info("Admin user already exists.")
	}

	slog.Info("Starting server")
	srv := server.NewServer(db)
	srv.Start()
}

func handleAdminUserNotFound(err error, db *database.Database) error {
	if err.Error() != "record not found" {
		slog.Error("Error checking for admin user", "error", err)
		return err
	}

	slog.Info("Admin user not found, creating...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("adminpassword"), 8)
	if err != nil {
		slog.Error("Failed to hash admin password", "error", err)
		return err
	}

	adminUser := database.User{
		Name:     "Admin User",
		Email:    "admin@test.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if createResult := db.GetDB().Create(&adminUser); createResult.Error != nil {
		slog.Error("Failed to create admin user", "error", createResult.Error)
		return createResult.Error
	}

	// Create associated Player record for admin user
	adminPlayer := database.Player{
		UserID:  adminUser.ID,
		Name:    "Admin User",
		Address: "",
	}

	if createResult := db.GetDB().Create(&adminPlayer); createResult.Error != nil {
		slog.Error("Failed to create admin player record", "error", createResult.Error)
		// Continue even if player creation fails - admin can still function
	} else {
		slog.Info("Admin player record created successfully.")
	}

	slog.Info("Admin user created successfully.")
	return nil
}
