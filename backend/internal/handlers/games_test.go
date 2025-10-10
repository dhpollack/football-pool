package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dhpollack/football-pool/internal/api"
	"github.com/dhpollack/football-pool/internal/database"
)

func homeAndAway() (string, string) {
	return "Home", "Away"
}

func TestMain(m *testing.M) {
	// Database setup is now handled in individual tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateGame(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()
	handler := CreateGame(gormDB)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "create single game",
			body:           `[{"week": 1, "season": 2023, "home_team": "Team A", "away_team": "Team B", "favorite": "Home", "underdog": "Away", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}]`,
			expectedStatus: http.StatusCreated,
			expectedCount:  1,
		},
		{
			name:           "create multiple games",
			body:           `[{"week": 2, "season": 2023, "home_team": "Team C", "away_team": "Team D", "favorite": "Home", "underdog": "Away", "spread": 7.0, "start_time": "2023-09-18T20:15:00Z"}, {"week": 2, "season": 2023, "home_team": "Team E", "away_team": "Team F", "favorite": "Home", "underdog": "Away", "spread": 1.5, "start_time": "2023-09-18T21:25:00Z"}]`,
			expectedStatus: http.StatusCreated,
			expectedCount:  2,
		},
		{
			name:           "invalid json",
			body:           `[{"week": 1}]`,
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/admin/games/create", bytes.NewBuffer([]byte(tt.body)))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response []api.GameResponse
				json.NewDecoder(rr.Body).Decode(&response)
				if len(response) != tt.expectedCount {
					t.Errorf("expected %d games in response, got %d", tt.expectedCount, len(response))
				}

				var count int64
				gormDB.Model(&database.Game{}).Count(&count)
				if count != int64(tt.expectedCount) {
					t.Errorf("expected %d games in db, got %d", tt.expectedCount, count)
				}
			}
			// Clean up db after each test run
			gormDB.Exec("DELETE FROM games")
		})
	}
}

func TestGetGames(t *testing.T) {
	// Set up test database
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()
	home, away := homeAndAway()

	// Seed the database with some games
	game := database.Game{Week: 1, Season: 2023, HomeTeam: "Lions", AwayTeam: "Chiefs", Favorite: &home, Underdog: &away}
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
	var response api.GameListResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Games) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(response.Games), 1)
	}

	if response.Games[0].HomeTeam != "Lions" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Games[0].HomeTeam, "Lions")
	}
}

func TestGetGamesErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("sqlite", "file::memory:")
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

func TestAdminListGames(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()
	home, away := homeAndAway()

	// Seed data
	gormDB.Create(&database.Game{Week: 1, Season: 2023, HomeTeam: "Team A", AwayTeam: "Team B", Favorite: &home, Underdog: &away})
	gormDB.Create(&database.Game{Week: 2, Season: 2023, HomeTeam: "Team C", AwayTeam: "Team D", Favorite: &home, Underdog: &away})
	gormDB.Create(&database.Game{Week: 1, Season: 2024, HomeTeam: "Team E", AwayTeam: "Team F", Favorite: &home, Underdog: &away})

	handler := AdminListGames(gormDB)

	t.Run("list all games", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/games", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response api.GameListResponse
		json.NewDecoder(rr.Body).Decode(&response)
		games := response.Games
		if len(games) != 3 {
			t.Errorf("expected 3 games, got %d", len(games))
		}
	})

	t.Run("filter by week and season", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/games?week=1&season=2023", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response api.GameListResponse
		json.NewDecoder(rr.Body).Decode(&response)
		games := response.Games
		if len(games) != 1 {
			t.Errorf("expected 1 game, got %d", len(games))
		}
	})
}

func TestUpdateGame(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()
	game := database.Game{Week: 1, Season: 2023, HomeTeam: "Old Team", AwayTeam: "Old Underdog", Spread: 3.5, StartTime: time.Now()}
	gormDB.Create(&game)

	// Create router directly for testing
	mux := http.NewServeMux()
	mux.Handle("PUT /api/admin/games/{id}", UpdateGame(gormDB))
	mux.Handle("DELETE /api/admin/games/{id}", DeleteGame(gormDB))
	router := mux

	t.Run("successful update", func(t *testing.T) {
		updatePayload := []byte(`{"week": 1, "season": 2023, "home_team": "New Team", "away_team": "Old Underdog", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}`)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/games/%d", game.ID), bytes.NewBuffer(updatePayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedGame api.GameResponse
		json.NewDecoder(rr.Body).Decode(&updatedGame)
		if updatedGame.HomeTeam != "New Team" {
			t.Errorf("expected favorite team to be 'New Team', got '%s'", updatedGame.HomeTeam)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		updatePayload := []byte(`{"week": 1, "season": 2023, "home_team": "New Team", "away_team": "Old Underdog", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"`) // Invalid JSON
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/games/%d", game.ID), bytes.NewBuffer(updatePayload))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("game not found", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/admin/games/999", bytes.NewBuffer([]byte(`{"week": 1, "season": 2023, "home_team": "New Team", "away_team": "Old Underdog", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}`)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("invalid game ID format", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/admin/games/abc", bytes.NewBuffer([]byte(`{"week": 1, "season": 2023, "home_team": "New Team", "away_team": "Old Underdog", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}`)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("validation error - missing required fields", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/games/%d", game.ID), bytes.NewBuffer([]byte(`{"week": 1, "season": 2023, "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}`)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestDeleteGame(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create router for testing
	mux := http.NewServeMux()
	mux.Handle("DELETE /api/admin/games/{id}", DeleteGame(gormDB))
	router := mux

	t.Run("successful deletion", func(t *testing.T) {
		game := database.Game{Week: 1, Season: 2025, HomeTeam: "Deletable", AwayTeam: "Team"}
		gormDB.Create(&game)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/games/%d", game.ID), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	})

	t.Run("game not found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/admin/games/999", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("conflict with existing picks", func(t *testing.T) {
		game := database.Game{Week: 2, Season: 2025, HomeTeam: "Conflict", AwayTeam: "Team"}
		gormDB.Create(&game)
		user := database.User{Name: "Test User", Email: "test@example.com", Role: "user"}
		gormDB.Create(&user)
		pick := database.Pick{UserID: user.ID, GameID: game.ID, Picked: "Conflict"}
		gormDB.Create(&pick)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/games/%d", game.ID), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		}
	})

	t.Run("conflict with existing results", func(t *testing.T) {
		game := database.Game{Week: 3, Season: 2025, HomeTeam: "ConflictResult", AwayTeam: "Team"}
		gormDB.Create(&game)
		result := database.Result{GameID: game.ID, Outcome: "favorite"}
		gormDB.Create(&result)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/games/%d", game.ID), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		}
	})
}
