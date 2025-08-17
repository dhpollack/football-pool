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
	database.Connect("file::memory:?cache=shared")
	code := m.Run()
	os.Exit(code)
}

func TestGetGames(t *testing.T) {
	// Seed the database with some games
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs"}
	database.DB.Create(&game)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/games?week=1&season=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetGames)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
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
			handler := http.HandlerFunc(GetGames)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
