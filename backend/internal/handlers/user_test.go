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

func TestGetProfile(t *testing.T) {
	// Create a user and a player
	user := database.User{Email: "profile_user@test.com", Password: "password"}
	database.DB.Create(&user)
	player := database.Player{UserID: user.ID, Name: "Test User", Address: "123 Test St"}
	database.DB.Create(&player)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("GET", "/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "email", "profile_user@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetProfile)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	var dbPlayer database.Player
	if err := json.NewDecoder(rr.Body).Decode(&dbPlayer); err != nil {
		t.Fatal(err)
	}

	if dbPlayer.Name != "Test User" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			dbPlayer.Name, "Test User")
	}
}

func TestUpdateProfile(t *testing.T) {
	// Create a user and a player
	user := database.User{Email: "profile_user2@test.com", Password: "password"}
	database.DB.Create(&user)
	player := database.Player{UserID: user.ID, Name: "Test User", Address: "123 Test St"}
	database.DB.Create(&player)

	// Create the updated profile to submit
	updates := struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}{
		Name:    "Updated User",
		Address: "456 Updated St",
	}
	jsonUpdates, _ := json.Marshal(updates)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("PUT", "/profile", bytes.NewBuffer(jsonUpdates))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), "email", "profile_user2@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateProfile)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the player was updated in the database
	var dbPlayer database.Player
	database.DB.Where("user_id = ?", user.ID).First(&dbPlayer)
	if dbPlayer.Name != "Updated User" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			dbPlayer.Name, "Updated User")
	}
}