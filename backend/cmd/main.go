// Package main is the entry point for the football pool backend server.
package main

import (
	"context"
	"fmt"
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
	dbConfig := cfg.Database.GetConfig()
	if dbConfig == nil {
		slog.Error("Invalid database configuration")
		os.Exit(1)
	}
	db, err := database.New(cfg.Database.Type, dbConfig.GetDSN())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Load and create configured users
	slog.Info("Loading user configuration")
	users, err := config.LoadUserConfig()
	if err != nil {
		slog.Error("Failed to load user configuration", "error", err)
		os.Exit(1)
	}

	// Create configured users if they don't exist
	if err := createConfiguredUsers(users, db); err != nil {
		slog.Error("Failed to create configured users", "error", err)
		os.Exit(1)
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

// initSyncService initializes the ESPN sync service with configuration.
// createConfiguredUsers creates users from configuration if they don't already exist.
func createConfiguredUsers(users []config.UserConfig, db *database.Database) error {
	for _, userCfg := range users {
		// Check if user already exists
		var existingUser database.User
		result := db.GetDB().Where("email = ?", userCfg.Email).First(&existingUser)

		// If no error, user already exists
		if result.Error == nil {
			slog.Info("User already exists", "email", userCfg.Email)
			continue
		}

		// If error is not "record not found", return the error
		if result.Error.Error() != "record not found" {
			return fmt.Errorf("error checking for user %s: %w", userCfg.Email, result.Error)
		}

		// User doesn't exist, create it
		slog.Info("Creating configured user", "email", userCfg.Email)

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCfg.Password), 8)
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", userCfg.Email, err)
		}

		// Create user
		user := database.User{
			Name:     userCfg.Name,
			Email:    userCfg.Email,
			Password: string(hashedPassword),
			Role:     userCfg.Role,
		}

		if createResult := db.GetDB().Create(&user); createResult.Error != nil {
			return fmt.Errorf("failed to create user %s: %w", userCfg.Email, createResult.Error)
		}

		// Create associated Player record
		player := database.Player{
			UserID:  user.ID,
			Name:    userCfg.Name,
			Address: "",
		}

		if createResult := db.GetDB().Create(&player); createResult.Error != nil {
			slog.Warn("Failed to create player record for user", "email", userCfg.Email, "error", createResult.Error)
			// Continue even if player creation fails - user can still function
		} else {
			slog.Info("Player record created successfully for user", "email", userCfg.Email)
		}

		slog.Info("User created successfully", "email", userCfg.Email)
	}

	return nil
}

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
