package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
)

func TestGetPicks(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
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
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user and a game
	user := database.User{Email: "test2@test.com", Password: "password"}
	gormDB.Create(&user)
	game := database.Game{Week: 1, Season: 2023, FavoriteTeam: "Packers", UnderdogTeam: "Bears"}
	gormDB.Create(&game)

	// Create the picks to submit
	picks := []database.Pick{
		{GameID: game.ID, Picked: "favorite", Rank: 1},
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
	db, err := database.New("file::memory:?cache=shared")
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
	db, err := database.New("file::memory:?cache=shared")
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
