package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/david/football-pool/internal/api"
	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
)

func TestGetPicks(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user and a game
	user := database.User{Email: "test@test.com", Password: "password"}
	gormDB.Create(&user)
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs"}
	gormDB.Create(&game)

	// Create a pick for the user
	pick := database.Pick{UserID: user.ID, GameID: game.ID, Picked: "favorite", Rank: 1}
	gormDB.Create(&pick)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("GET", "/picks", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), auth.EmailKey, "test@test.com")
	req = req.WithContext(ctx)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := GetPicks(gormDB)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var picks []database.Pick
	if err := json.NewDecoder(rr.Body).Decode(&picks); err != nil {
		t.Fatal(err)
	}

	if len(picks) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(picks), 1)
	}

	if picks[0].Picked != "favorite" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			picks[0].Picked, "favorite")
	}
}

func TestSubmitPicks(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user and a game
	user := database.User{Name: "testuser", Email: "test2@test.com", Password: "password", Role: "user"}
	gormDB.Create(&user)
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Packers", UnderdogTeam: "Bears", Spread: 3.5, StartTime: time.Now()}
	gormDB.Create(&game)

	// Create the picks to submit
	picks := []api.PickRequest{
		{GameId: game.ID, Picked: "favorite", Rank: 1},
	}
	jsonPicks, _ := json.Marshal(picks)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("POST", "/picks", bytes.NewBuffer(jsonPicks))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), auth.EmailKey, "test2@test.com")
	req = req.WithContext(ctx)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := SubmitPicks(gormDB)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		// Print the response body to see the error
		var errorResponse api.ErrorResponse
		if err := json.NewDecoder(rr.Body).Decode(&errorResponse); err == nil {
			t.Logf("Error response: %+v", errorResponse)
		}
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check that the picks were created in the database
	var dbPicks []database.Pick
	gormDB.Where("user_id = ?", user.ID).Find(&dbPicks)
	if len(dbPicks) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(dbPicks), 1)
	}

	if dbPicks[0].Picked != "favorite" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			dbPicks[0].Picked, "favorite")
	}
}

func TestGetPicksErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	tests := []struct {
		name           string
		email          string
		expectedStatus int
	}{
		{
			name:           "User not found",
			email:          "notfound@test.com",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/picks", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), auth.EmailKey, tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := GetPicks(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestSubmitPicksErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	user := database.User{Email: "test3@test.com", Password: "password"}
	gormDB.Create(&user)

	tests := []struct {
		name           string
		email          string
		body           string
		expectedStatus int
	}{
		{
			name:           "User not found",
			email:          "notfound@test.com",
			body:           `[]`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid JSON",
			email:          "test3@test.com",
			body:           `[`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/picks", bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), auth.EmailKey, tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := SubmitPicks(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestAdminListPicks(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create test data
	user1 := database.User{Email: "user1@test.com", Password: "password"}
	user2 := database.User{Email: "user2@test.com", Password: "password"}
	gormDB.Create(&user1)
	gormDB.Create(&user2)

	game1 := database.Game{Week: 1, Season: 2024, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs"}
	game2 := database.Game{Week: 2, Season: 2024, FavoriteTeam: "Packers", UnderdogTeam: "Bears"}
	gormDB.Create(&game1)
	gormDB.Create(&game2)

	pick1 := database.Pick{UserID: user1.ID, GameID: game1.ID, Picked: "favorite", Rank: 1}
	pick2 := database.Pick{UserID: user2.ID, GameID: game1.ID, Picked: "underdog", Rank: 1}
	pick3 := database.Pick{UserID: user1.ID, GameID: game2.ID, Picked: "favorite", Rank: 2}
	gormDB.Create(&pick1)
	gormDB.Create(&pick2)
	gormDB.Create(&pick3)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "List all picks",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "Filter by user_id",
			queryParams:    fmt.Sprintf("?user_id=%d", user1.ID),
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Filter by game_id",
			queryParams:    fmt.Sprintf("?game_id=%d", game1.ID),
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Filter by week",
			queryParams:    "?week=1",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "Pagination",
			queryParams:    "?page=1&limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/admin/picks"+tt.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := AdminListPicks(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatal(err)
				}

				picks := response["picks"].([]interface{})
				if len(picks) != tt.expectedCount {
					t.Errorf("handler returned unexpected pick count: got %v want %v",
						len(picks), tt.expectedCount)
				}
			}
		})
	}
}

func TestAdminDeletePick(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create test data
	user := database.User{Email: "user@test.com", Password: "password"}
	gormDB.Create(&user)

	game := database.Game{Week: 1, Season: 2024, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs"}
	gormDB.Create(&game)

	pick := database.Pick{UserID: user.ID, GameID: game.ID, Picked: "favorite", Rank: 1}
	gormDB.Create(&pick)

	// Test deleting the pick
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/admin/picks/%d", pick.ID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := AdminDeletePick(gormDB)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify the pick was deleted
	var deletedPick database.Pick
	result := gormDB.First(&deletedPick, pick.ID)
	if result.Error == nil {
		t.Error("pick was not deleted from database")
	}
}
