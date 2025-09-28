package espnsync

import (
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Config holds configuration for the SyncService.
type Config struct {
	ESPNBaseURL  string        `toml:"base_url"`
	CacheDir     string        `toml:"cache_dir"`
	SyncEnabled  bool          `toml:"sync_enabled"`
	SyncInterval time.Duration `toml:"sync_interval"`
	CacheExpiry  time.Duration `toml:"cache_expiry"`
	SeasonYear   int           `toml:"season_year"`
	Week1Date    time.Time     `toml:"week1_date"`
}

// LoadConfig loads configuration from a TOML file.
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config struct {
		ESPN_Sync Config `toml:"espn_sync"`
	}

	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config.ESPN_Sync, nil
}

// DefaultConfig returns a default configuration for testing.
func DefaultConfig() Config {
	return Config{
		ESPNBaseURL:  "https://site.api.espn.com",
		CacheDir:     "./cache",
		SyncEnabled:  false,
		SyncInterval: time.Hour,
		CacheExpiry:  24 * time.Hour,
		SeasonYear:   2024,
		Week1Date:    time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC),
	}
}
