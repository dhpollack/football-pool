package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david/football-pool/internal/database"
	"github.com/david/football-pool/internal/server"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIntegration(t *testing.T) {
	// Set up the test database
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	database.DB = db

	// Migrate the database
	err = db.AutoMigrate(&database.User{}, &database.Game{}, &database.Pick{})
	assert.NoError(t, err)

	// Set up the router
	router := server.NewRouter()

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
	var games []database.Game
	err = json.NewDecoder(resp.Body).Decode(&games)
	assert.NoError(t, err)
	assert.Len(t, games, 2)
	assert.Equal(t, "Team A", games[0].FavoriteTeam)
	assert.Equal(t, "Team C", games[1].FavoriteTeam)
}

func createUser(t *testing.T, ts *httptest.Server, name, email, password, role string) {
	user := database.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", ts.URL+"/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func loginUser(t *testing.T, ts *httptest.Server, email, password string) string {
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", ts.URL+"/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenMap map[string]string
	json.NewDecoder(resp.Body).Decode(&tokenMap)
	return tokenMap["token"]
}

func createGames(t *testing.T, ts *httptest.Server, token string) {
	games := []database.Game{
		{Week: 1, Season: 2023, FavoriteTeam: "Team A", UnderdogTeam: "Team B", Spread: 3.5},
		{Week: 1, Season: 2023, FavoriteTeam: "Team C", UnderdogTeam: "Team D", Spread: 7.0},
	}
	body, _ := json.Marshal(games)
	req, _ := http.NewRequest("POST", ts.URL+"/api/games/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func submitPicks(t *testing.T, ts *httptest.Server, token string) {
	picks := []database.Pick{
		{UserID: 1, GameID: 1, Picked: "favorite", Rank: 1},
		{UserID: 1, GameID: 2, Picked: "underdog", Rank: 2},
	}
	body, _ := json.Marshal(picks)
	req, _ := http.NewRequest("POST", ts.URL+"/api/admin/picks/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
