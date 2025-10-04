package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Test loading with test environment
	t.Setenv("FOOTBALL_POOL_ENV", "test")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify config values from test.toml
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "8080", cfg.Server.Port)
	dbConfig := cfg.Database.GetConfig()
	require.NotNil(t, dbConfig)
	assert.Equal(t, ":memory:", dbConfig.GetDSN())
	assert.Equal(t, "error", cfg.Logging.Level)
	assert.Equal(t, "https://site.api.espn.com/apis/site/v2/sports/football/nfl", cfg.ESPN.BaseURL)
	assert.Equal(t, "/tmp/test-cache", cfg.ESPN.CacheDir)
	assert.False(t, cfg.ESPN.SyncEnabled)
	assert.Equal(t, 1*time.Hour, cfg.ESPN.SyncInterval)
	assert.Equal(t, 24*time.Hour, cfg.ESPN.CacheExpiry)
	assert.False(t, cfg.E2E.Test)
}

func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// Set environment variables that should override config file
	t.Setenv("FOOTBALL_POOL_ENV", "test")
	t.Setenv("FOOTBALL_POOL_PORT", "9999")
	t.Setenv("FOOTBALL_POOL_DB_FILE", "test_override.db")
	t.Setenv("FOOTBALL_POOL_LOG_LEVEL", "debug")
	t.Setenv("ESPN_SYNC_ENABLED", "true")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify environment variables override config file
	assert.Equal(t, "9999", cfg.Server.Port)
	dbConfig := cfg.Database.GetConfig()
	require.NotNil(t, dbConfig)
	assert.Equal(t, "test_override.db", dbConfig.GetDSN())
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.True(t, cfg.ESPN.SyncEnabled)
}

func TestLoadConfigWithMissingConfigFile(t *testing.T) {
	// Test loading with non-existent environment - should fail
	t.Setenv("FOOTBALL_POOL_ENV", "nonexistent")

	_, err := LoadConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "error reading config file")
}

func TestConfigDefaults(t *testing.T) {
	// Test that defaults are applied correctly with empty config
	t.Setenv("FOOTBALL_POOL_ENV", "defaults")

	cfg, err := LoadConfig()
	require.NoError(t, err)

	// Verify default values (should match setDefaults())
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "8080", cfg.Server.Port)
	dbConfig := cfg.Database.GetConfig()
	require.NotNil(t, dbConfig)
	assert.Equal(t, "football-pool.db", dbConfig.GetDSN())
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "https://site.api.espn.com/apis/site/v2/sports/football/nfl", cfg.ESPN.BaseURL)
	assert.Equal(t, "assets/cache", cfg.ESPN.CacheDir)
	assert.False(t, cfg.ESPN.SyncEnabled)
	assert.Equal(t, 1*time.Hour, cfg.ESPN.SyncInterval)
	assert.Equal(t, 24*time.Hour, cfg.ESPN.CacheExpiry)
	assert.False(t, cfg.E2E.Test)
}

func TestLoadConfigProd(t *testing.T) {
	// Test loading with prod environment
	t.Setenv("FOOTBALL_POOL_ENV", "prod")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify prod config values
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "8080", cfg.Server.Port)
	dbConfig := cfg.Database.GetConfig()
	require.NotNil(t, dbConfig)
	assert.Equal(t, "football-pool.db", dbConfig.GetDSN())
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "https://site.api.espn.com/apis/site/v2/sports/football/nfl", cfg.ESPN.BaseURL)
	assert.Equal(t, "assets/cache", cfg.ESPN.CacheDir)
	assert.True(t, cfg.ESPN.SyncEnabled)
	assert.Equal(t, 1*time.Hour, cfg.ESPN.SyncInterval)
	assert.Equal(t, 1000*time.Hour, cfg.ESPN.CacheExpiry)
	assert.False(t, cfg.E2E.Test)
}

func TestEnvironmentVariableBinding(t *testing.T) {
	// Test that all expected environment variables are bound
	t.Setenv("FOOTBALL_POOL_ENV", "test")
	t.Setenv("FOOTBALL_POOL_HOST", "test-host")
	t.Setenv("FOOTBALL_POOL_PORT", "1234")
	t.Setenv("FOOTBALL_POOL_DB_FILE", "test-dsn")
	t.Setenv("FOOTBALL_POOL_LOG_LEVEL", "warn")
	t.Setenv("ESPN_BASE_URL", "http://test.example.com")
	t.Setenv("ESPN_CACHE_DIR", "/test/cache")
	t.Setenv("ESPN_SYNC_ENABLED", "true")
	t.Setenv("ESPN_SYNC_INTERVAL", "2h")
	t.Setenv("ESPN_CACHE_EXPIRY", "48h")
	t.Setenv("E2E_TEST", "true")

	cfg, err := LoadConfig()
	require.NoError(t, err)

	// Verify environment variables are properly bound
	assert.Equal(t, "test-host", cfg.Server.Host)
	assert.Equal(t, "1234", cfg.Server.Port)
	assert.Equal(t, "test-dsn", cfg.Database.GetConfig().GetDSN())
	assert.Equal(t, "warn", cfg.Logging.Level)
	assert.Equal(t, "http://test.example.com", cfg.ESPN.BaseURL)
	assert.Equal(t, "/test/cache", cfg.ESPN.CacheDir)
	assert.True(t, cfg.ESPN.SyncEnabled)
	assert.Equal(t, 2*time.Hour, cfg.ESPN.SyncInterval)
	assert.Equal(t, 48*time.Hour, cfg.ESPN.CacheExpiry)
	assert.True(t, cfg.E2E.Test)
}

func TestPostgreSQLConfigurationWithStringPort(t *testing.T) {
	// Test PostgreSQL configuration with string port values from environment variables
	t.Setenv("FOOTBALL_POOL_ENV", "test")
	t.Setenv("FOOTBALL_POOL_DB_TYPE", "postgres")
	t.Setenv("FOOTBALL_POOL_DB_HOST", "localhost")
	t.Setenv("FOOTBALL_POOL_DB_PORT", "5432") // String port value
	t.Setenv("FOOTBALL_POOL_DB_USER", "testuser")
	t.Setenv("FOOTBALL_POOL_DB_PASSWORD", "testpass")
	t.Setenv("FOOTBALL_POOL_DB_NAME", "testdb")
	t.Setenv("FOOTBALL_POOL_DB_SSLMODE", "disable")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify PostgreSQL configuration is properly loaded
	assert.Equal(t, "postgres", cfg.Database.Type)

	postgresConfig, ok := cfg.Database.GetConfig().(PostgresConfig)
	require.True(t, ok, "Database config should be PostgresConfig")

	// Verify all PostgreSQL fields are properly set
	assert.Equal(t, "localhost", postgresConfig.Host)
	assert.Equal(t, 5432, postgresConfig.Port) // Should be converted from string "5432" to int 5432
	assert.Equal(t, "testuser", postgresConfig.User)
	assert.Equal(t, "testpass", postgresConfig.Password)
	assert.Equal(t, "testdb", postgresConfig.DBName)
	assert.Equal(t, "disable", postgresConfig.SSLMode)

	// Verify DSN is correctly formatted
	expectedDSN := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expectedDSN, postgresConfig.GetDSN())
}
