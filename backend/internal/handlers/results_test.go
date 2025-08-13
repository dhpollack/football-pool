package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david/football-pool/internal/database"
)

func TestSubmitResult(t *testing.T) {
	// Create an admin user and a game
	admin := database.User{Email: "admin@test.com", Password: "password", Role: "admin"}
	database.DB.Create(&admin)
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs", Spread: 3.5}
	database.DB.Create(&game)

	// Create the result to submit
	result := database.Result{GameID: game.ID, FavoriteScore: 21, UnderdogScore: 17}
	jsonResult, _ := json.Marshal(result)

	// Create a request with the admin's email in the context
	req, err := http.NewRequest("POST", "/results", bytes.NewBuffer(jsonResult))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "email", "admin@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitResult)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check that the result was created in the database
	var dbResult database.Result
	database.DB.Where("game_id = ?", game.ID).First(&dbResult)
	if dbResult.Outcome != "favorite" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			dbResult.Outcome, "favorite")
	}
}

func TestGetWeeklyResults(t *testing.T) {
	// Create users, games, picks, and results
	user1 := database.User{Name: "User 1", Email: "user1@test.com", Password: "password"}
	database.DB.Create(&user1)
	user2 := database.User{Name: "User 2", Email: "user2@test.com", Password: "password"}
	database.DB.Create(&user2)

	game1 := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs", Spread: 3.5}
	database.DB.Create(&game1)
	game2 := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Eagles", UnderdogTeam: "Patriots", Spread: 7.5}
	database.DB.Create(&game2)

	pick1 := database.Pick{UserID: user1.ID, GameID: game1.ID, PickedTeam: "favorite", Rank: 16}
	database.DB.Create(&pick1)
	pick2 := database.Pick{UserID: user1.ID, GameID: game2.ID, PickedTeam: "underdog", Rank: 1}
	database.DB.Create(&pick2)
	pick3 := database.Pick{UserID: user2.ID, GameID: game1.ID, PickedTeam: "underdog", Rank: 10}
	database.DB.Create(&pick3)
	pick4 := database.Pick{UserID: user2.ID, GameID: game2.ID, PickedTeam: "favorite", Rank: 5}
	database.DB.Create(&pick4)

	result1 := database.Result{GameID: game1.ID, FavoriteScore: 21, UnderdogScore: 17, Outcome: "favorite"}
	database.DB.Create(&result1)
	result2 := database.Result{GameID: game2.ID, FavoriteScore: 34, UnderdogScore: 10, Outcome: "favorite"}
	database.DB.Create(&result2)

	// Create a request with the week and season as query parameters
	req, err := http.NewRequest("GET", "/results/weekly?week=1&season=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetWeeklyResults)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	type WeeklyResult struct {
		PlayerID   uint    `json:"player_id"`
		PlayerName string  `json:"player_name"`
		Score      float32 `json:"score"`
	}

	var weeklyResults []WeeklyResult
	if err := json.NewDecoder(rr.Body).Decode(&weeklyResults); err != nil {
		t.Fatal(err)
	}

	if len(weeklyResults) != 2 {
		t.Errorf("handler returned unexpected number of results: got %v want %v",
			len(weeklyResults), 2)
	}

	// Check the scores for each user
	for _, result := range weeklyResults {
		if result.PlayerName == "User 1" {
			if result.Score != 16 {
				t.Errorf("handler returned wrong score for User 1: got %v want %v",
					result.Score, 16)
			}
		} else if result.PlayerName == "User 2" {
			if result.Score != 5 {
				t.Errorf("handler returned wrong score for User 2: got %v want %v",
					result.Score, 5)
			}
		}
	}
}

func TestGetSeasonResults(t *testing.T) {
	// Create users, games, picks, and results
	user1 := database.User{Name: "User 1", Email: "user3@test.com", Password: "password"}
	database.DB.Create(&user1)
	user2 := database.User{Name: "User 2", Email: "user4@test.com", Password: "password"}
	database.DB.Create(&user2)

	game1 := database.Game{Week: 1, Season: 2024, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs", Spread: 3.5}
	database.DB.Create(&game1)
	game2 := database.Game{Week: 2, Season: 2024, FavoriteTeam: "Eagles", UnderdogTeam: "Patriots", Spread: 7.5}
	database.DB.Create(&game2)

	pick1 := database.Pick{UserID: user1.ID, GameID: game1.ID, PickedTeam: "favorite", Rank: 16}
	database.DB.Create(&pick1)
	pick2 := database.Pick{UserID: user1.ID, GameID: game2.ID, PickedTeam: "underdog", Rank: 1}
	database.DB.Create(&pick2)
	pick3 := database.Pick{UserID: user2.ID, GameID: game1.ID, PickedTeam: "underdog", Rank: 10}
	database.DB.Create(&pick3)
	pick4 := database.Pick{UserID: user2.ID, GameID: game2.ID, PickedTeam: "favorite", Rank: 5}
	database.DB.Create(&pick4)

	result1 := database.Result{GameID: game1.ID, FavoriteScore: 21, UnderdogScore: 17, Outcome: "favorite"}
	database.DB.Create(&result1)
	result2 := database.Result{GameID: game2.ID, FavoriteScore: 34, UnderdogScore: 10, Outcome: "favorite"}
	database.DB.Create(&result2)

	// Create a request with the season as query parameters
	req, err := http.NewRequest("GET", "/results/season?season=2024", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSeasonResults)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	type SeasonResult struct {
		PlayerID   uint    `json:"player_id"`
		PlayerName string  `json:"player_name"`
		Score      float32 `json:"score"`
	}

	var seasonResults []SeasonResult
	if err := json.NewDecoder(rr.Body).Decode(&seasonResults); err != nil {
		t.Fatal(err)
	}

	if len(seasonResults) != 2 {
		t.Errorf("handler returned unexpected number of results: got %v want %v",
			len(seasonResults), 2)
	}

	// Check the scores for each user
	for _, result := range seasonResults {
		if result.PlayerName == "User 1" {
			if result.Score != 16 {
				t.Errorf("handler returned wrong score for User 1: got %v want %v",
					result.Score, 16)
			}
		} else if result.PlayerName == "User 2" {
			if result.Score != 5 {
				t.Errorf("handler returned wrong score for User 2: got %v want %v",
					result.Score, 5)
			}
		}
	}
}

func TestGetWeeklyResultsErrors(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Missing week",
			url:            "/results/weekly?season=2023",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing season",
			url:            "/results/weekly?week=1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid week",
			url:            "/results/weekly?week=abc&season=2023",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid season",
			url:            "/results/weekly?week=1&season=abc",
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
			handler := http.HandlerFunc(GetWeeklyResults)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestSubmitResultErrors(t *testing.T) {
	admin := database.User{Email: "admin2@test.com", Password: "password", Role: "admin"}
	database.DB.Create(&admin)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Invalid JSON",
			body:           `{`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Game not found",
			body:           `{"game_id": 999}`,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/results", bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), "email", "admin2@test.com")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SubmitResult)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestGetSeasonResultsErrors(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Missing season",
			url:            "/results/season",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid season",
			url:            "/results/season?season=abc",
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
			handler := http.HandlerFunc(GetSeasonResults)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
