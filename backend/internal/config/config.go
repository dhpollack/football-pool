package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	// Database configuration
	Database struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	// Logging configuration
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`

	// ESPN sync configuration
	ESPN struct {
		BaseURL      string        `mapstructure:"base_url"`
		CacheDir     string        `mapstructure:"cache_dir"`
		SyncEnabled  bool          `mapstructure:"sync_enabled"`
		SyncInterval time.Duration `mapstructure:"sync_interval"`
		CacheExpiry  time.Duration `mapstructure:"cache_expiry"`
	} `mapstructure:"espn"`

	// E2E testing configuration
	E2E struct {
		Test bool `mapstructure:"test"`
	} `mapstructure:"e2e"`
}

// LoadConfig loads configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	// Determine environment
	env := os.Getenv("FOOTBALL_POOL_ENV")
	if env == "" {
		env = "prod"
	}

	viper.SetConfigName(env)        // name of config file (without extension)
	viper.SetConfigType("toml")     // or "yaml", "json", etc.
	viper.AddConfigPath(".")        // look for config in the current directory
	viper.AddConfigPath("./config") // look for config in the config directory

	// Set default values
	setDefaults()

	// Bind environment variables
	bindEnvVars()

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Info("No config file found, using environment variables and defaults", "env", env)
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// NewTestConfig creates a configuration suitable for testing
func NewTestConfig() *Config {
	// Set test environment and load config
	os.Setenv("FOOTBALL_POOL_ENV", "test")
	cfg, err := LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load test configuration: %v", err))
	}
	return cfg
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")

	// Database defaults
	viper.SetDefault("database.dsn", "football-pool.db")

	// Logging defaults
	viper.SetDefault("logging.level", "info")

	// ESPN sync defaults
	viper.SetDefault("espn.base_url", "https://site.api.espn.com/apis/site/v2/sports/football/nfl")
	viper.SetDefault("espn.cache_dir", "assets/cache")
	viper.SetDefault("espn.sync_enabled", false)
	viper.SetDefault("espn.sync_interval", "1h")
	viper.SetDefault("espn.cache_expiry", "24h")

	// E2E testing defaults
	viper.SetDefault("e2e.test", false)
}

func bindEnvVars() {
	// Server environment variables
	viper.BindEnv("server.host", "FOOTBALL_POOL_HOST")
	viper.BindEnv("server.port", "FOOTBALL_POOL_PORT")

	// Database environment variables
	viper.BindEnv("database.dsn", "FOOTBALL_POOL_DSN")

	// Logging environment variables
	viper.BindEnv("logging.level", "FOOTBALL_POOL_LOG_LEVEL")

	// ESPN sync environment variables
	viper.BindEnv("espn.base_url", "ESPN_BASE_URL")
	viper.BindEnv("espn.cache_dir", "ESPN_CACHE_DIR")
	viper.BindEnv("espn.sync_enabled", "ESPN_SYNC_ENABLED")
	viper.BindEnv("espn.sync_interval", "ESPN_SYNC_INTERVAL")
	viper.BindEnv("espn.cache_expiry", "ESPN_CACHE_EXPIRY")

	// E2E testing environment variables
	viper.BindEnv("e2e.test", "E2E_TEST")
}
