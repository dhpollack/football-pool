// Package main is the entry point for the football pool backend server.
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/david/football-pool/internal/database"
	"github.com/david/football-pool/internal/espn-sync"
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

	// Initialize ESPN sync service
	syncService, err := initSyncService(db)
	if err != nil {
		slog.Error("Failed to initialize ESPN sync service", "error", err)
		// Continue without sync service - it's not critical for server startup
	}

	// Start background sync process if service was initialized
	if syncService != nil {
		ctx := context.Background()
		syncInterval := getSyncInterval()
		go syncService.Start(ctx, syncInterval)
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

// initSyncService initializes the ESPN sync service with configuration from environment variables.
func initSyncService(db *database.Database) (*espnsync.SyncService, error) {
	espnBaseURL := os.Getenv("ESPN_BASE_URL")
	if espnBaseURL == "" {
		espnBaseURL = "https://site.api.espn.com/apis/site/v2/sports/football/nfl"
	}

	cacheDir := os.Getenv("ESPN_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = "assets/cache"
	}

	// Sync service is disabled by default and must be explicitly enabled
	// This ensures it doesn't run during E2E tests unless specifically configured
	syncEnabled := os.Getenv("ESPN_SYNC_ENABLED") == "true"

	// Disable sync service during E2E tests to avoid interference
	if os.Getenv("E2E_TEST") == "true" {
		syncEnabled = false
		slog.Info("ESPN sync service disabled for E2E tests")
	}

	config := espnsync.Config{
		ESPNBaseURL: espnBaseURL,
		CacheDir:    cacheDir,
		SyncEnabled: syncEnabled,
		CacheExpiry: getCacheExpiry(),
	}

	slog.Info("Initializing ESPN sync service",
		"enabled", syncEnabled,
		"base_url", espnBaseURL,
		"cache_dir", cacheDir,
	)

	return espnsync.NewSyncService(db, config)
}

// getSyncInterval returns the sync interval from environment variables.
func getSyncInterval() time.Duration {
	syncIntervalStr := os.Getenv("ESPN_SYNC_INTERVAL")
	if syncIntervalStr == "" {
		return time.Hour // Default to 1 hour
	}

	duration, err := time.ParseDuration(syncIntervalStr)
	if err != nil {
		slog.Warn("Invalid ESPN_SYNC_INTERVAL, using default", "value", syncIntervalStr, "error", err)
		return time.Hour
	}

	return duration
}

// getCacheExpiry returns the cache expiry duration from environment variables.
func getCacheExpiry() time.Duration {
	cacheExpiryStr := os.Getenv("ESPN_CACHE_EXPIRY")
	if cacheExpiryStr == "" {
		return 24 * time.Hour // Default to 24 hours
	}

	duration, err := time.ParseDuration(cacheExpiryStr)
	if err != nil {
		slog.Warn("Invalid ESPN_CACHE_EXPIRY, using default", "value", cacheExpiryStr, "error", err)
		return 24 * time.Hour
	}

	return duration
}
