package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/david/football-pool/internal/api"
	"github.com/david/football-pool/internal/database"
)

func TestMain(m *testing.M) {
	// Database setup is now handled in individual tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateGame(t *testing.T) {
	db, err := database.New("file::memory:")
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
			body:           `[{"week": 1, "season": 2023, "favorite_team": "Team A", "underdog_team": "Team B", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}]`,
			expectedStatus: http.StatusCreated,
			expectedCount:  1,
		},
		{
			name:           "create multiple games",
			body:           `[{"week": 2, "season": 2023, "favorite_team": "Team C", "underdog_team": "Team D", "spread": 7.0, "start_time": "2023-09-18T20:15:00Z"}, {"week": 2, "season": 2023, "favorite_team": "Team E", "underdog_team": "Team F", "spread": 1.5, "start_time": "2023-09-18T21:25:00Z"}]`,
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
	db, err := database.New("file::memory:")
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
	var response api.GameListResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response.Games) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(response.Games), 1)
	}

	if response.Games[0].FavoriteTeam != "Lions" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Games[0].FavoriteTeam, "Lions")
	}
}

func TestGetGamesErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
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
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Seed data
	gormDB.Create(&database.Game{Week: 1, Season: 2023, FavoriteTeam: "Team A", UnderdogTeam: "Team B"})
	gormDB.Create(&database.Game{Week: 2, Season: 2023, FavoriteTeam: "Team C", UnderdogTeam: "Team D"})
	gormDB.Create(&database.Game{Week: 1, Season: 2024, FavoriteTeam: "Team E", UnderdogTeam: "Team F"})

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
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Old Team", UnderdogTeam: "Old Underdog", Spread: 3.5, StartTime: time.Now()}
	gormDB.Create(&game)

	handler := UpdateGame(gormDB)

	t.Run("successful update", func(t *testing.T) {
		updatePayload := []byte(`{"week": 1, "season": 2023, "favorite_team": "New Team", "underdog_team": "Old Underdog", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}`)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/games/%d", game.ID), bytes.NewBuffer(updatePayload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedGame api.GameResponse
		json.NewDecoder(rr.Body).Decode(&updatedGame)
		if updatedGame.FavoriteTeam != "New Team" {
			t.Errorf("expected favorite team to be 'New Team', got '%s'", updatedGame.FavoriteTeam)
		}
	})

	t.Run("game not found", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/admin/games/999", bytes.NewBuffer([]byte(`{"week": 1, "season": 2023, "favorite_team": "New Team", "underdog_team": "Old Underdog", "spread": 3.5, "start_time": "2023-09-11T20:15:00Z"}`)))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

func TestDeleteGame(t *testing.T) {
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()
	handler := DeleteGame(gormDB)

	t.Run("successful deletion", func(t *testing.T) {
		game := database.Game{Week: 1, Season: 2025, FavoriteTeam: "Deletable", UnderdogTeam: "Team"}
		gormDB.Create(&game)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/games/%d", game.ID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}
	})

	t.Run("game not found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/admin/games/999", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("conflict with existing picks", func(t *testing.T) {
		game := database.Game{Week: 2, Season: 2025, FavoriteTeam: "Conflict", UnderdogTeam: "Team"}
		gormDB.Create(&game)
		user := database.User{Name: "Test User", Email: "test@example.com", Role: "user"}
		gormDB.Create(&user)
		pick := database.Pick{UserID: user.ID, GameID: game.ID, Picked: "Conflict"}
		gormDB.Create(&pick)

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/games/%d", game.ID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		}
	})
}

func TestCreateGameFromSeed(t *testing.T) {
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()
	handler := CreateGame(gormDB)

	// Read and process the seed file
	file, err := os.Open("../../seed/games.jsonl")
	if err != nil {
		t.Fatalf("Failed to open seed file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Failed to read seed file: %v", err)
	}
	seedBody := "[" + strings.Join(lines, ",") + "]"

	// Let's try to unmarshal it here to see if it's valid
	var gameRequests []api.GameRequest
	if err := json.Unmarshal([]byte(seedBody), &gameRequests); err != nil {
		t.Fatalf("Failed to unmarshal seed body: %v", err)
	}

	expectedCount := len(lines)

	req, _ := http.NewRequest("POST", "/api/admin/games/create", bytes.NewBuffer([]byte(seedBody)))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("handler returned wrong status code: got %v want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
	}

	var response []api.GameResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != expectedCount {
		t.Errorf("expected %d games in response, got %d", expectedCount, len(response))
	}

	var count int64
	gormDB.Model(&database.Game{}).Count(&count)
	if count != int64(expectedCount) {
		t.Errorf("expected %d games in db, got %d", expectedCount, count)
	}
}
