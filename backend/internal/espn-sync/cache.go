// Package espnsync provides functionality to sync season schedule and game scores from the ESPN API.
package espnsync

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/david/football-pool/internal/api-espn"
)

// Cache manages caching of ESPN API responses.
type Cache struct {
	cacheDir string
	expiry   time.Duration
}

// NewCache creates a new Cache instance.
func NewCache(cacheDir string, expiry time.Duration) *Cache {
	return &Cache{
		cacheDir: cacheDir,
		expiry:   expiry,
	}
}

// cacheKey generates a cache key for a specific API call.
func (c *Cache) cacheKey(season, week int) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("espn-events-%d-%d", season, week)))
	return fmt.Sprintf("%x.json", hash[:8])
}

// cachePath returns the full path to a cache file.
func (c *Cache) cachePath(key string) string {
	return filepath.Join(c.cacheDir, key)
}

// Get retrieves cached data for a specific season and week.
func (c *Cache) Get(season, week int) ([]apiespn.Event, bool) {
	key := c.cacheKey(season, week)
	path := c.cachePath(key)

	// Check if cache file exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, false
	}

	// Check if cache is expired (skip if expiry is negative = never expire)
	if c.expiry >= 0 && time.Since(info.ModTime()) > c.expiry {
		os.Remove(path) // Remove expired cache
		return nil, false
	}

	// Read cache file
	file, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	// Parse cached data
	var events []apiespn.Event
	if err := json.NewDecoder(file).Decode(&events); err != nil {
		slog.Warn("Failed to parse cache file", "path", path, "error", err)
		return nil, false
	}

	slog.Debug("Cache hit", "season", season, "week", week)
	return events, true
}

// Set stores data in the cache for a specific season and week.
func (c *Cache) Set(season, week int, events []apiespn.Event) error {
	key := c.cacheKey(season, week)
	path := c.cachePath(key)

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
		return err
	}

	// Create cache file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to cache file
	if err := json.NewEncoder(file).Encode(events); err != nil {
		return err
	}

	slog.Debug("Cache set", "season", season, "week", week, "events", len(events))
	return nil
}

// ClearExpired removes all expired cache files.
func (c *Cache) ClearExpired() error {
	// Skip if expiry is negative (never expire)
	if c.expiry < 0 {
		return nil
	}

	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Cache directory doesn't exist yet
		}
		return err
	}

	now := time.Now()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(c.cacheDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > c.expiry {
			if err := os.Remove(path); err != nil {
				slog.Warn("Failed to remove expired cache file", "path", path, "error", err)
			} else {
				slog.Debug("Removed expired cache file", "path", path)
			}
		}
	}

	return nil
}

// ClearAll removes all cache files.
func (c *Cache) ClearAll() error {
	entries, err := os.ReadDir(c.cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Cache directory doesn't exist yet
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(c.cacheDir, entry.Name())
		if err := os.Remove(path); err != nil {
			slog.Warn("Failed to remove cache file", "path", path, "error", err)
		} else {
			slog.Debug("Removed cache file", "path", path)
		}
	}

	return nil
}
