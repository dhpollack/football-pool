// Package main is the entry point for the football pool backend server.
package main

import (
	"context"
	"log/slog"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/dhpollack/football-pool/internal/config"
	"github.com/dhpollack/football-pool/internal/database"
	espnsync "github.com/dhpollack/football-pool/internal/espn-sync"
	"github.com/dhpollack/football-pool/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup logging
	var level slog.Level
	switch cfg.Logging.Level {
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

	slog.Info("Connecting to database")
	db, err := database.New(cfg.Database.Type, cfg.Database.DSN)
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

	// Initialize ESPN sync service
	syncService, err := initSyncService(db, cfg)
	if err != nil {
		slog.Error("Failed to initialize ESPN sync service", "error", err)
		// Continue without sync service - it's not critical for server startup
	}

	// Start background sync process if service was initialized
	if syncService != nil {
		ctx := context.Background()
		go syncService.Start(ctx, cfg.ESPN.SyncInterval)
	}

	slog.Info("Starting server")
	srv := server.NewServer(db, cfg)
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

// initSyncService initializes the ESPN sync service with configuration.
func initSyncService(db *database.Database, cfg *config.Config) (*espnsync.SyncService, error) {
	// Disable sync service during E2E tests to avoid interference
	syncEnabled := cfg.ESPN.SyncEnabled
	if cfg.E2E.Test {
		syncEnabled = false
		slog.Info("ESPN sync service disabled for E2E tests")
	}

	// Temporarily override sync enabled status for E2E tests
	tempConfig := *cfg
	tempConfig.ESPN.SyncEnabled = syncEnabled

	slog.Info("Initializing ESPN sync service",
		"enabled", syncEnabled,
		"base_url", cfg.ESPN.BaseURL,
		"cache_dir", cfg.ESPN.CacheDir,
	)

	return espnsync.NewSyncService(db, &tempConfig)
}
