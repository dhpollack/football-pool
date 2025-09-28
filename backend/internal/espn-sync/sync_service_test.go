package espnsync

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/david/football-pool/internal/api-espn"
	"github.com/david/football-pool/internal/config"
	"github.com/david/football-pool/internal/database"
)

// testConfig returns a test configuration for use in tests
func testConfig(t *testing.T) *config.Config {
	config := &config.Config{}
	config.ESPN.BaseURL = "https://test.api.espn.com"
	config.ESPN.CacheDir = t.TempDir()
	config.ESPN.SyncEnabled = false
	config.ESPN.SyncInterval = time.Hour
	config.ESPN.CacheExpiry = 24 * time.Hour
	config.ESPN.SeasonYear = 2025
	config.ESPN.Week1Date = time.Date(2025, 9, 4, 0, 0, 0, 0, time.UTC)
	return config
}

// MockHTTPClient is a mock implementation of the HTTP client for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// MockTimeProvider is a mock implementation of TimeProvider for testing
type MockTimeProvider struct {
	NowFunc func() time.Time
}

func (m MockTimeProvider) Now() time.Time {
	if m.NowFunc != nil {
		return m.NowFunc()
	}
	return time.Now()
}

func TestSyncService_StartDisabled(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	config := testConfig(t)

	service, err := NewSyncService(db, config)
	if err != nil {
		t.Fatalf("NewSyncService() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start should return immediately when disabled
	go service.Start(ctx, time.Hour)

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Service should not be running any sync operations
	// This test mainly verifies that the service doesn't crash when disabled
}

func TestSyncService_SyncNow(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Return a minimal valid response with all required fields
			return &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body: io.NopCloser(strings.NewReader(`{
					"season": {"type": 2, "year": 2025},
					"week": {"number": 1},
					"leagues": [{
						"abbreviation": "NFL",
						"name": "National Football League",
						"uid": "s:20~l:28"
					}],
					"events": [{
						"id": "event1",
						"name": "Test Game",
						"date": "2025-09-05T00:20Z",
						"competitions": [{
							"competitors": [{
								"homeAway": "home",
								"team": {"displayName": "Team A"},
								"score": "24"
							}, {
								"homeAway": "away",
								"team": {"displayName": "Team B"},
								"score": "17"
							}]
						}]
					}]
				}`)),
			}, nil
		},
	}

	config := testConfig(t)
	config.ESPN.BaseURL = "https://test.api.espn.com"
	config.ESPN.CacheDir = t.TempDir()
	config.ESPN.SyncEnabled = true
	config.ESPN.SeasonYear = 2025                                       // Match the season year in the sample file
	config.ESPN.Week1Date = time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC) // Set Week 1 date to match the mock response

	// Create the ESPN client with our mock HTTP client
	client, err := apiespn.NewClientWithResponses(config.ESPN.BaseURL, apiespn.WithHTTPClient(mockHTTPClient))
	if err != nil {
		t.Fatalf("Failed to create ESPN client: %v", err)
	}

	// Create a mock time provider that returns a specific time matching our mock response
	mockTimeProvider := MockTimeProvider{
		NowFunc: func() time.Time {
			return time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC)
		},
	}

	service, err := NewSyncServiceWithTimeProvider(db, config, mockTimeProvider)
	if err != nil {
		t.Fatalf("NewSyncServiceWithTimeProvider() error = %v", err)
	}

	// Replace the ESPN client with our mock
	service.espnClient = client

	ctx := context.Background()
	err = service.SyncNow(ctx)
	if err != nil {
		t.Fatalf("SyncNow() error = %v", err)
	}

	// Verify that games were stored in the database
	var games []database.Game
	if err := db.GetDB().Find(&games).Error; err != nil {
		t.Fatalf("Failed to query games: %v", err)
	}

	if len(games) == 0 {
		t.Error("Expected games to be stored in database")
	}

	// Verify the specific game was stored
	var game database.Game
	if err := db.GetDB().Where("season = ? AND week = ? AND favorite_team = ? AND underdog_team = ?",
		2025, 1, "Team A", "Team B").First(&game).Error; err != nil {
		t.Fatalf("Failed to find expected game: %v", err)
	}
}

func TestSyncService_SyncNowAPIError(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Return a proper error response without valid JSON structure
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("Internal Server Error")),
			}, nil
		},
	}

	config := testConfig(t)
	config.ESPN.BaseURL = "https://test.api.espn.com"
	config.ESPN.CacheDir = t.TempDir()
	config.ESPN.SyncEnabled = true
	config.ESPN.SeasonYear = 2024
	config.ESPN.Week1Date = time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC) // Set Week 1 date

	// Create the ESPN client with our mock HTTP client
	client, err := apiespn.NewClientWithResponses(config.ESPN.BaseURL, apiespn.WithHTTPClient(mockHTTPClient))
	if err != nil {
		t.Fatalf("Failed to create ESPN client: %v", err)
	}

	// Create a mock time provider that returns a date that would calculate to week 18
	mockTimeProvider := MockTimeProvider{
		NowFunc: func() time.Time {
			return time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC) // About 119 days after Week 1
		},
	}

	service, err := NewSyncServiceWithTimeProvider(db, config, mockTimeProvider)
	if err != nil {
		t.Fatalf("NewSyncServiceWithTimeProvider() error = %v", err)
	}

	service.espnClient = client

	ctx := context.Background()
	err = service.SyncNow(ctx)

	// The sync service should handle API errors gracefully and return nil
	if err != nil {
		t.Errorf("SyncNow() unexpected error: %v", err)
	}
}

func TestSyncService_GetCurrentSeasonAndWeek(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	config := testConfig(t)
	config.ESPN.BaseURL = "https://test.api.espn.com"
	config.ESPN.CacheDir = t.TempDir()
	config.ESPN.SyncEnabled = true
	config.ESPN.SeasonYear = 2023
	config.ESPN.Week1Date = time.Date(2023, 9, 7, 0, 0, 0, 0, time.UTC) // Week 1 Thursday

	tests := []struct {
		name         string
		now          time.Time
		expectedWeek int
	}{
		{
			name:         "Week 2",
			now:          time.Date(2023, 9, 15, 0, 0, 0, 0, time.UTC), // September 15
			expectedWeek: 2,
		},
		{
			name:         "Week 6",
			now:          time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC), // October 15
			expectedWeek: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock time provider for this specific test case
			mockTimeProvider := MockTimeProvider{
				NowFunc: func() time.Time {
					return tt.now
				},
			}

			service, err := NewSyncServiceWithTimeProvider(db, config, mockTimeProvider)
			if err != nil {
				t.Fatalf("NewSyncServiceWithTimeProvider() error = %v", err)
			}

			season, week := service.getCurrentSeasonAndWeek()

			if season != 2023 {
				t.Errorf("getCurrentSeasonAndWeek() season = %v, want 2023", season)
			}
			if week != tt.expectedWeek {
				t.Errorf("getCurrentSeasonAndWeek() week = %v, want %v", week, tt.expectedWeek)
			}
		})
	}
}

func TestSyncService_GetSyncStatus(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	config := testConfig(t)
	config.ESPN.BaseURL = "https://test.api.espn.com"
	config.ESPN.CacheDir = t.TempDir()
	config.ESPN.SyncEnabled = true

	service, err := NewSyncService(db, config)
	if err != nil {
		t.Fatalf("NewSyncService() error = %v", err)
	}

	status := service.GetSyncStatus()

	if enabled, ok := status["enabled"].(bool); !ok || !enabled {
		t.Error("GetSyncStatus() enabled should be true")
	}

	if _, ok := status["last_sync"].(time.Time); !ok {
		t.Error("GetSyncStatus() last_sync should be a time.Time")
	}
}

func TestSyncService_SyncNowWithSampleData(t *testing.T) {
	// Enable debug logging
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))

	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Read the sample JSON file
	sampleData, err := os.ReadFile("../../assets/test/scoreboard-sample.json")
	if err != nil {
		t.Fatalf("Failed to read sample data file: %v", err)
	}

	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// Return the actual sample data as response
			return &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(string(sampleData))),
			}, nil
		},
	}

	config := testConfig(t)
	config.ESPN.BaseURL = "https://test.api.espn.com"
	config.ESPN.CacheDir = t.TempDir()
	config.ESPN.SyncEnabled = true
	config.ESPN.SeasonYear = 2025                                       // Match the season year in the sample file
	config.ESPN.Week1Date = time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC) // Set Week 1 date to match the sample data

	// Create the ESPN client with our mock HTTP client
	client, err := apiespn.NewClientWithResponses(config.ESPN.BaseURL, apiespn.WithHTTPClient(mockHTTPClient))
	if err != nil {
		t.Fatalf("Failed to create ESPN client: %v", err)
	}

	// Create a mock time provider that returns a specific time matching our sample data
	mockTimeProvider := MockTimeProvider{
		NowFunc: func() time.Time {
			return time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC)
		},
	}

	service, err := NewSyncServiceWithTimeProvider(db, config, mockTimeProvider)
	if err != nil {
		t.Fatalf("NewSyncServiceWithTimeProvider() error = %v", err)
	}

	// Replace the ESPN client with our mock
	service.espnClient = client

	ctx := context.Background()
	err = service.SyncNow(ctx)
	if err != nil {
		t.Fatalf("SyncNow() error = %v", err)
	}

	// Verify that games were stored in the database
	var games []database.Game
	if err := db.GetDB().Find(&games).Error; err != nil {
		t.Fatalf("Failed to query games: %v", err)
	}

	// The sample file contains multiple games, so we should have several stored
	if len(games) == 0 {
		t.Error("Expected games to be stored in database")
	}

	// Debug: print all stored games
	t.Logf("Stored %d games:", len(games))
	for i, game := range games {
		t.Logf("Game %d: Season=%d, Week=%d, Favorite=%s, Underdog=%s",
			i+1, game.Season, game.Week, game.FavoriteTeam, game.UnderdogTeam)
	}

	// Verify that at least one specific game from the sample was stored
	// Let's check for the first game in the sample: Dallas Cowboys at Philadelphia Eagles
	var game database.Game
	if err := db.GetDB().Where("season = ? AND week = ? AND favorite_team = ? AND underdog_team = ?",
		2025, 1, "Philadelphia Eagles", "Dallas Cowboys").First(&game).Error; err != nil {
		t.Fatalf("Failed to find expected game: %v", err)
	}

	// Verify the game has the expected properties
	if game.Season != 2025 {
		t.Errorf("Expected season 2025, got %v", game.Season)
	}
	if game.Week != 1 {
		t.Errorf("Expected week 1, got %v", game.Week)
	}
	if game.FavoriteTeam != "Philadelphia Eagles" {
		t.Errorf("Expected favorite team 'Philadelphia Eagles', got '%v'", game.FavoriteTeam)
	}
	if game.UnderdogTeam != "Dallas Cowboys" {
		t.Errorf("Expected underdog team 'Dallas Cowboys', got '%v'", game.UnderdogTeam)
	}

	// Verify that results were stored for games with scores
	var results []database.Result
	if err := db.GetDB().Find(&results).Error; err != nil {
		t.Fatalf("Failed to query results: %v", err)
	}

	// Debug: print all results
	t.Logf("Found %d results in database:", len(results))
	for i, result := range results {
		t.Logf("Result %d: GameID=%d, FavoriteScore=%d, UnderdogScore=%d, Outcome=%s",
			i+1, result.GameID, result.FavoriteScore, result.UnderdogScore, result.Outcome)
	}

	// The sample file contains games with scores, so we should have some results stored
	if len(results) == 0 {
		t.Error("Expected results to be stored in database for games with scores")
	}

	// Verify that the specific game we checked has a result
	var result database.Result
	if err := db.GetDB().Where("game_id = ?", game.ID).First(&result).Error; err != nil {
		t.Fatalf("Failed to find result for game: %v", err)
	}

	// Verify the result has the expected properties
	// Philadelphia Eagles (home/favorite) scored 24, Dallas Cowboys (away/underdog) scored 20
	if result.FavoriteScore != 24 {
		t.Errorf("Expected favorite score 24, got %v", result.FavoriteScore)
	}
	if result.UnderdogScore != 20 {
		t.Errorf("Expected underdog score 20, got %v", result.UnderdogScore)
	}
	if result.Outcome != "favorite" {
		t.Errorf("Expected outcome 'favorite', got '%v'", result.Outcome)
	}

	t.Logf("Successfully stored %d games and %d results from sample data", len(games), len(results))
}
