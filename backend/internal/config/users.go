// Package config provides configuration management for the football pool application.
package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// UserConfig represents a user to be created on application startup
type UserConfig struct {
	Name     string `mapstructure:"name"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
	Role     string `mapstructure:"role"`
}

// LoadUserConfig loads user configuration from a TOML file.
// The file location can be specified via FOOTBALL_POOL_USERS_CONFIG environment variable.
// If not specified, it will look for users.toml in the standard config paths.
func LoadUserConfig() ([]UserConfig, error) {
	// Check if a specific config file is specified via environment variable
	configFile := os.Getenv("FOOTBALL_POOL_USERS_CONFIG")

	// Create a new viper instance for user config
	userViper := viper.New()

	if configFile != "" {
		// Use the specified config file
		userViper.SetConfigFile(configFile)
	} else {
		// Use the default behavior - look for users.toml in standard paths
		userViper.SetConfigName("users")        // name of config file (without extension)
		userViper.SetConfigType("toml")         // or "yaml", "json", etc.
		userViper.AddConfigPath(".")            // look for config in the current directory
		userViper.AddConfigPath("./config")     // look for config in the config directory
		userViper.AddConfigPath("../")          // look for config in parent directory (for tests)
		userViper.AddConfigPath("../../")       // look for config in backend root directory (for tests)
		userViper.AddConfigPath("../../config") // look for config in backend config directory (for tests)
	}

	// Read config file (if exists)
	if err := userViper.ReadInConfig(); err != nil {
		// If file doesn't exist, return empty users (no error)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return []UserConfig{}, nil
		}
		return nil, fmt.Errorf("error reading user config file: %w", err)
	}

	var users []UserConfig
	if err := userViper.UnmarshalKey("users", &users); err != nil {
		return nil, fmt.Errorf("error unmarshaling user config: %w", err)
	}

	return users, nil
}