package espnsync

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/david/football-pool/internal/api-espn"
)

func TestCache_GetSet(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir, time.Hour)

	// Test data
	events := []apiespn.Event{
		{
			Id:   &[]string{"event1"}[0],
			Name: &[]string{"Game 1"}[0],
		},
		{
			Id:   &[]string{"event2"}[0],
			Name: &[]string{"Game 2"}[0],
		},
	}

	// Test setting cache
	err := cache.Set(2023, 1, events)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Test getting cache
	cachedEvents, found := cache.Get(2023, 1)
	if !found {
		t.Fatal("Cache should be found")
	}

	if len(cachedEvents) != len(events) {
		t.Errorf("Expected %d events, got %d", len(events), len(cachedEvents))
	}

	if *cachedEvents[0].Id != "event1" {
		t.Errorf("Expected event ID 'event1', got '%s'", *cachedEvents[0].Id)
	}
}

func TestCache_Expired(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir, time.Millisecond) // Very short expiry

	events := []apiespn.Event{
		{
			Id:   &[]string{"event1"}[0],
			Name: &[]string{"Game 1"}[0],
		},
	}

	// Set cache
	err := cache.Set(2023, 1, events)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Wait for cache to expire
	time.Sleep(2 * time.Millisecond)

	// Cache should be expired and not found
	cachedEvents, found := cache.Get(2023, 1)
	if found {
		t.Error("Cache should be expired and not found")
	}
	if len(cachedEvents) != 0 {
		t.Error("Expired cache should return empty events")
	}
}

func TestCache_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir, time.Hour)

	// Try to get non-existent cache
	cachedEvents, found := cache.Get(2023, 1)
	if found {
		t.Error("Cache should not be found for non-existent key")
	}
	if len(cachedEvents) != 0 {
		t.Error("Non-existent cache should return empty events")
	}
}

func TestCache_ClearExpired(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir, time.Millisecond)

	events := []apiespn.Event{
		{
			Id:   &[]string{"event1"}[0],
			Name: &[]string{"Game 1"}[0],
		},
	}

	// Set cache that will expire
	err := cache.Set(2023, 1, events)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Wait for cache to expire
	time.Sleep(2 * time.Millisecond)

	// Clear expired caches
	err = cache.ClearExpired()
	if err != nil {
		t.Fatalf("Failed to clear expired cache: %v", err)
	}

	// Verify cache file is removed
	key := cache.cacheKey(2023, 1)
	path := cache.cachePath(key)
	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		t.Error("Expired cache file should be removed")
	}
}

func TestCache_ClearAll(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir, time.Hour)

	events := []apiespn.Event{
		{
			Id:   &[]string{"event1"}[0],
			Name: &[]string{"Game 1"}[0],
		},
	}

	// Set multiple caches
	err := cache.Set(2023, 1, events)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}
	err = cache.Set(2023, 2, events)
	if err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Clear all caches
	err = cache.ClearAll()
	if err != nil {
		t.Fatalf("Failed to clear all cache: %v", err)
	}

	// Verify cache files are removed
	key1 := cache.cacheKey(2023, 1)
	path1 := cache.cachePath(key1)
	_, err = os.Stat(path1)
	if !os.IsNotExist(err) {
		t.Error("Cache file should be removed")
	}

	key2 := cache.cacheKey(2023, 2)
	path2 := cache.cachePath(key2)
	_, err = os.Stat(path2)
	if !os.IsNotExist(err) {
		t.Error("Cache file should be removed")
	}
}

func TestCache_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	cache := NewCache(tempDir, time.Hour)

	// Create invalid JSON file
	key := cache.cacheKey(2023, 1)
	path := cache.cachePath(key)

	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}

	err = os.WriteFile(path, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid cache file: %v", err)
	}

	// Try to get cache with invalid JSON
	cachedEvents, found := cache.Get(2023, 1)
	if found {
		t.Error("Cache with invalid JSON should not be found")
	}
	if len(cachedEvents) != 0 {
		t.Error("Invalid cache should return empty events")
	}
}
