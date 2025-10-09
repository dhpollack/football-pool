// Package main provides integration tests for the football pool backend API.
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dhpollack/football-pool/internal/api"
	"github.com/dhpollack/football-pool/internal/config"
	"github.com/dhpollack/football-pool/internal/database"
	"github.com/dhpollack/football-pool/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	// Set up the test configuration and database
	t.Setenv("FOOTBALL_POOL_ENV", "test")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}
	dbConfig := cfg.Database.GetConfig()
	if dbConfig == nil {
		t.Fatalf("Invalid database configuration")
	}
	db, err := database.New(cfg.Database.Type, dbConfig.GetDSN())
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	// Set up the server with database
	srv := server.NewServer(db, cfg)
	router := srv.NewRouter()

	// Create a new test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Create a regular user
	createUser(t, ts, "testuser", "test@example.com", "password", "user")

	// Create an admin user
	createUser(t, ts, "adminuser", "admin@example.com", "password", "admin")

	// Log in as the admin user
	token := loginUser(t, ts, "admin@example.com", "password")

	// Create games
	createGames(t, ts, token)

	// Submit picks for the regular user
	submitPicks(t, ts, token)

	// Make a request to an endpoint to verify the data
	resp, err := http.Get(ts.URL + "/api/games?season=2023&week=1")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check the response body
	var gameListResponse api.GameListResponse
	err = json.NewDecoder(resp.Body).Decode(&gameListResponse)
	assert.NoError(t, err)
	assert.Len(t, gameListResponse.Games, 2)
	assert.Equal(t, "Team A", gameListResponse.Games[0].HomeTeam)
	assert.Equal(t, "Team C", gameListResponse.Games[1].HomeTeam)
}

func createUser(t *testing.T, ts *httptest.Server, name, email, password, role string) {
	user := api.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     &role,
	}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", ts.URL+"/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Print response status and body for debugging
	t.Logf("Create user status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Logf("Create user response: %s", string(bodyBytes))
		if err := resp.Body.Close(); err != nil {
			t.Logf("Error closing response body: %v", err)
		}
	}

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func loginUser(t *testing.T, ts *httptest.Server, email, password string) string {
	loginData := api.LoginRequest{
		Email:    email,
		Password: password,
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", ts.URL+"/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResponse api.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}
	return loginResponse.Token
}

func createGames(t *testing.T, ts *httptest.Server, token string) {
	favoriteHome := api.Home
	underdogAway := api.Away
	games := []api.GameRequest{
		{Week: 1, Season: 2023, HomeTeam: "Team A", AwayTeam: "Team B", Spread: 3.5, StartTime: time.Now(), Favorite: &favoriteHome, Underdog: &underdogAway},
		{Week: 1, Season: 2023, HomeTeam: "Team C", AwayTeam: "Team D", Spread: 7.0, StartTime: time.Now(), Favorite: &favoriteHome, Underdog: &underdogAway},
	}
	body, _ := json.Marshal(games)
	req, _ := http.NewRequest("POST", ts.URL+"/api/admin/games/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func submitPicks(t *testing.T, ts *httptest.Server, token string) {
	picks := []api.PickRequest{
		{GameId: 1, Picked: "favorite", Rank: 1, QuickPick: false},
		{GameId: 2, Picked: "underdog", Rank: 2, QuickPick: false},
	}
	body, _ := json.Marshal(picks)
	req, _ := http.NewRequest("POST", ts.URL+"/api/picks/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
