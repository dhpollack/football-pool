// Package config provides configuration management for the football pool application.
package config

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// DatabaseConfig is the interface for database configurations.
type DatabaseConfig interface {
	GetDSN() string
}

// SQLiteConfig holds SQLite-specific database configuration.
type SQLiteConfig struct {
	// File path for the SQLite database file
	File string `mapstructure:"file"`
}

// GetDSN returns the SQLite DSN string.
func (s SQLiteConfig) GetDSN() string {
	return s.File
}

// PostgresConfig holds PostgreSQL-specific database configuration.
type PostgresConfig struct {
	// Host is the database server hostname or IP address
	Host string `mapstructure:"host"`
	// Port is the database server port
	Port int `mapstructure:"port"`
	// User is the database username
	User string `mapstructure:"user"`
	// Password is the database password
	Password string `mapstructure:"password"`
	// DBName is the database name
	DBName string `mapstructure:"dbname"`
	// SSLMode is the SSL mode (disable, require, verify-ca, verify-full)
	SSLMode string `mapstructure:"sslmode"`
}

// GetDSN returns the PostgreSQL DSN string.
func (p PostgresConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

// DatabaseConfigWrapper is a wrapper that handles unmarshaling of database configurations.
type DatabaseConfigWrapper struct {
	Type   string         `mapstructure:"type"`
	Config DatabaseConfig `mapstructure:"config"`
}

// GetConfig returns the appropriate DatabaseConfig based on the type.
func (w DatabaseConfigWrapper) GetConfig() DatabaseConfig {
	return w.Config
}

// Config holds all configuration for the application.
type Config struct {
	// Server configuration
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	// Database configuration
	Database DatabaseConfigWrapper `mapstructure:"database"`

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
		SeasonYear   int           `mapstructure:"season_year"`
		Week1Date    time.Time     `mapstructure:"week1_date"`
	} `mapstructure:"espn"`

	// E2E testing configuration
	E2E struct {
		Test bool `mapstructure:"test"`
	} `mapstructure:"e2e"`

	// TheOddsAPI configuration
	TheOddsAPI struct {
		BaseURL string `mapstructure:"base_url"`
		APIKey  string `mapstructure:"api_key"`
		Region  string `mapstructure:"region"`
	} `mapstructure:"theoddsapi"`
}

// LoadConfig loads configuration from environment variables and config file.
func LoadConfig() (*Config, error) {
	// Determine environment
	env := os.Getenv("FOOTBALL_POOL_ENV")
	if env == "" {
		env = "prod"
	}

	viper.SetConfigName(env)            // name of config file (without extension)
	viper.SetConfigType("toml")         // or "yaml", "json", etc.
	viper.AddConfigPath(".")            // look for config in the current directory
	viper.AddConfigPath("./config")     // look for config in the config directory
	viper.AddConfigPath("../")          // look for config in parent directory (for tests)
	viper.AddConfigPath("../../")       // look for config in backend root directory (for tests)
	viper.AddConfigPath("../../config") // look for config in backend config directory (for tests)

	// Set default values
	setDefaults()

	// Bind environment variables
	bindEnvVars()

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		databaseDecodeHook(),
		mapstructure.StringToTimeDurationHookFunc(),
	))); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")

	// Database defaults (SQLite by default)
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.config.file", "football-pool.db")

	// Logging defaults
	viper.SetDefault("logging.level", "info")

	// ESPN sync defaults
	viper.SetDefault("espn.base_url", "https://site.api.espn.com/apis/site/v2/sports/football/nfl")
	viper.SetDefault("espn.cache_dir", "assets/cache")
	viper.SetDefault("espn.sync_enabled", false)
	viper.SetDefault("espn.sync_interval", "1h")
	viper.SetDefault("espn.cache_expiry", "24h")
	viper.SetDefault("espn.season_year", 2025)
	viper.SetDefault("espn.week1_date", time.Date(2025, 9, 4, 0, 0, 0, 0, time.UTC))

	// E2E testing defaults
	viper.SetDefault("e2e.test", false)

	// TheOddsAPI defaults
	viper.SetDefault("theoddsapi.base_url", "https://api.the-odds-api.com/v4")
	viper.SetDefault("theoddsapi.region", "us")
	viper.SetDefault("theoddsapi.api_key", "")
}

func bindEnvVars() {
	// Server environment variables
	viper.BindEnv("server.host", "FOOTBALL_POOL_HOST")
	viper.BindEnv("server.port", "FOOTBALL_POOL_PORT")

	// Database environment variables
	viper.BindEnv("database.type", "FOOTBALL_POOL_DB_TYPE")
	viper.BindEnv("database.config.file", "FOOTBALL_POOL_DB_FILE")
	viper.BindEnv("database.config.host", "FOOTBALL_POOL_DB_HOST")
	viper.BindEnv("database.config.port", "FOOTBALL_POOL_DB_PORT")
	viper.BindEnv("database.config.user", "FOOTBALL_POOL_DB_USER")
	viper.BindEnv("database.config.password", "FOOTBALL_POOL_DB_PASSWORD")
	viper.BindEnv("database.config.dbname", "FOOTBALL_POOL_DB_NAME")
	viper.BindEnv("database.config.sslmode", "FOOTBALL_POOL_DB_SSLMODE")

	// Logging environment variables
	viper.BindEnv("logging.level", "FOOTBALL_POOL_LOG_LEVEL")

	// ESPN sync environment variables
	viper.BindEnv("espn.base_url", "ESPN_BASE_URL")
	viper.BindEnv("espn.cache_dir", "ESPN_CACHE_DIR")
	viper.BindEnv("espn.sync_enabled", "ESPN_SYNC_ENABLED")
	viper.BindEnv("espn.sync_interval", "ESPN_SYNC_INTERVAL")
	viper.BindEnv("espn.cache_expiry", "ESPN_CACHE_EXPIRY")
	viper.BindEnv("espn.season_year", "ESPN_SEASON_YEAR")
	viper.BindEnv("espn.week1_date", "ESPN_WEEK1_DATE")

	// E2E testing environment variables
	viper.BindEnv("e2e.test", "E2E_TEST")

	// TheOddsAPI environment variables
	viper.BindEnv("theoddsapi.base_url", "THEODDSAPI_BASE_URL")
	viper.BindEnv("theoddsapi.api_key", "THEODDSAPI_API_KEY")
	viper.BindEnv("theoddsapi.region", "THEODDSAPI_REGION")
}

func databaseDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		_ reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if t != reflect.TypeOf(DatabaseConfigWrapper{}) {
			return data, nil
		}

		m, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("data for DatabaseConfigWrapper is not a map")
		}

		dbType, ok := m["type"].(string)
		if !ok {
			// Default to sqlite if not present
			dbType = "sqlite"
		}

		var dbConfig DatabaseConfig
		configData, ok := m["config"]
		if !ok {
			configData = make(map[string]interface{})
		}

		switch dbType {
		case "sqlite":
			var sqliteCfg SQLiteConfig
			if err := mapstructure.Decode(configData, &sqliteCfg); err != nil {
				return nil, err
			}
			dbConfig = sqliteCfg
		case "postgres":
			var postgresCfg PostgresConfig
			// Use a custom decoder with string-to-int conversion for postgres config
			decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				Result:           &postgresCfg,
				WeaklyTypedInput: true, // This allows string-to-int conversion
			})
			if err != nil {
				return nil, err
			}
			if err := decoder.Decode(configData); err != nil {
				return nil, err
			}
			dbConfig = postgresCfg
		default:
			return nil, fmt.Errorf("unknown database type: %s", dbType)
		}

		return DatabaseConfigWrapper{
			Type:   dbType,
			Config: dbConfig,
		}, nil
	}
}
