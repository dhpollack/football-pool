package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dhpollack/football-pool/internal/auth"
	"github.com/dhpollack/football-pool/internal/database"
)

func TestGetSurvivorPicks(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user
	user := database.User{Email: "survivor_user@test.com", Password: "password"}
	gormDB.Create(&user)

	// Create a survivor pick for the user
	pick := database.SurvivorPick{UserID: user.ID, Week: 1, Team: "Lions"}
	gormDB.Create(&pick)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("GET", "/survivor", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), auth.EmailKey, "survivor_user@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := GetSurvivorPicks(gormDB)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var picks []database.SurvivorPick
	if err := json.NewDecoder(rr.Body).Decode(&picks); err != nil {
		t.Fatal(err)
	}

	if len(picks) != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			len(picks), 1)
	}

	if picks[0].Team != "Lions" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			picks[0].Team, "Lions")
	}
}

func TestSubmitSurvivorPick(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user
	user := database.User{Email: "survivor_user2@test.com", Password: "password"}
	gormDB.Create(&user)

	// Create the pick to submit
	pick := database.SurvivorPick{Week: 1, Team: "Packers"}
	jsonPick, _ := json.Marshal(pick)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("POST", "/survivor", bytes.NewBuffer(jsonPick))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), auth.EmailKey, "survivor_user2@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := SubmitSurvivorPick(gormDB)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check that the pick was created in the database
	var dbPick database.SurvivorPick
	gormDB.Where("user_id = ?", user.ID).First(&dbPick)
	if dbPick.Team != "Packers" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			dbPick.Team, "Packers")
	}
}

func TestGetSurvivorPicksErrors(t *testing.T) {
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
			req, err := http.NewRequest("GET", "/survivor", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), auth.EmailKey, tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := GetSurvivorPicks(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestSubmitSurvivorPickErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	user := database.User{Email: "survivor_user3@test.com", Password: "password"}
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
			body:           `{}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid JSON",
			email:          "survivor_user3@test.com",
			body:           `{`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/survivor", bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), auth.EmailKey, tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := SubmitSurvivorPick(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}
