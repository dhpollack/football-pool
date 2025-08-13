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

func TestGetPicks(t *testing.T) {
	// Create a user and a game
	user := database.User{Email: "test@test.com", Password: "password"}
	database.DB.Create(&user)
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs"}
	database.DB.Create(&game)

	// Create a pick for the user
	pick := database.Pick{UserID: user.ID, GameID: game.ID, PickedTeam: "Lions", Rank: 1}
	database.DB.Create(&pick)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("GET", "/picks", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "email", "test@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPicks)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var picks []database.Pick
	if err := json.NewDecoder(rr.Body).Decode(&picks); err != nil {
		t.Fatal(err)
	}

	if len(picks) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(picks), 1)
	}

	if picks[0].PickedTeam != "Lions" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			picks[0].PickedTeam, "Lions")
	}
}

func TestSubmitPicks(t *testing.T) {
	// Create a user and a game
	user := database.User{Email: "test2@test.com", Password: "password"}
	database.DB.Create(&user)
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Packers", UnderdogTeam: "Bears"}
	database.DB.Create(&game)

	// Create the picks to submit
	picks := []database.Pick{
		{GameID: game.ID, PickedTeam: "Packers", Rank: 1},
	}
	jsonPicks, _ := json.Marshal(picks)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("POST", "/picks", bytes.NewBuffer(jsonPicks))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "email", "test2@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SubmitPicks)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check that the picks were created in the database
	var dbPicks []database.Pick
	database.DB.Where("user_id = ?", user.ID).Find(&dbPicks)
	if len(dbPicks) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(dbPicks), 1)
	}

	if dbPicks[0].PickedTeam != "Packers" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			dbPicks[0].PickedTeam, "Packers")
	}
}

func TestGetPicksErrors(t *testing.T) {
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
			ctx := context.WithValue(req.Context(), "email", tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(GetPicks)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestSubmitPicksErrors(t *testing.T) {
	user := database.User{Email: "test3@test.com", Password: "password"}
	database.DB.Create(&user)

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
			ctx := context.WithValue(req.Context(), "email", tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(SubmitPicks)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
