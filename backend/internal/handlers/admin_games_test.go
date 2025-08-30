
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david/football-pool/internal/database"
)

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

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		games := response["games"].([]interface{})
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

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		games := response["games"].([]interface{})
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
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Old Team", UnderdogTeam: "Old Underdog"}
	gormDB.Create(&game)

	handler := UpdateGame(gormDB)

	t.Run("successful update", func(t *testing.T) {
		updatePayload := []byte(`{"favorite_team": "New Team"}`)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/games/%d", game.ID), bytes.NewBuffer(updatePayload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedGame database.Game
		json.NewDecoder(rr.Body).Decode(&updatedGame)
		if updatedGame.FavoriteTeam != "New Team" {
			t.Errorf("expected favorite team to be 'New Team', got '%s'", updatedGame.FavoriteTeam)
		}
	})

	t.Run("game not found", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/admin/games/999", bytes.NewBuffer([]byte(`{}`)))
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
