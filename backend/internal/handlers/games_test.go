package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/david/football-pool/internal/database"
)

func TestMain(m *testing.M) {
	// Database setup is now handled in individual tests
	code := m.Run()
	os.Exit(code)
}

func TestGetGames(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Seed the database with some games
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs"}
	gormDB.Create(&game)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/games?week=1&season=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler := GetGames(gormDB)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var games []database.Game
	if err := json.NewDecoder(rr.Body).Decode(&games); err != nil {
		t.Fatal(err)
	}

	if len(games) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(games), 1)
	}

	if games[0].FavoriteTeam != "Lions" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			games[0].FavoriteTeam, "Lions")
	}
}

func TestGetGamesErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Missing week",
			url:            "/games?season=2023",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing season",
			url:            "/games?week=1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid week",
			url:            "/games?week=abc&season=2023",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid season",
			url:            "/games?week=1&season=abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := GetGames(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
