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

func TestGetProfile(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user and a player
	user := database.User{Email: "profile_user@test.com", Password: "password"}
	gormDB.Create(&user)
	player := database.Player{UserID: user.ID, Name: "Test User", Address: "123 Test St"}
	gormDB.Create(&player)

	// Create a request with the user's email in the context
	req, err := http.NewRequest("GET", "/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.WithValue(req.Context(), auth.EmailKey, "profile_user@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := GetProfile(gormDB)

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
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user and a player
	user := database.User{Email: "profile_user2@test.com", Password: "password"}
	gormDB.Create(&user)
	player := database.Player{UserID: user.ID, Name: "Test User", Address: "123 Test St"}
	gormDB.Create(&player)

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
	ctx := context.WithValue(req.Context(), auth.EmailKey, "profile_user2@test.com")
	req = req.WithContext(ctx)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := UpdateProfile(gormDB)

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
	gormDB.Where("user_id = ?", user.ID).First(&dbPlayer)
	if dbPlayer.Name != "Updated User" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			dbPlayer.Name, "Updated User")
	}
}

func TestGetProfileErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	user := database.User{Email: "profile_user3@test.com", Password: "password"}
	gormDB.Create(&user)

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
		{
			name:           "Player not found",
			email:          "profile_user3@test.com",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/profile", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), auth.EmailKey, tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := GetProfile(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestUpdateProfileErrors(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	userWithPlayer := database.User{Email: "profile_user4@test.com", Password: "password"}
	gormDB.Create(&userWithPlayer)
	player := database.Player{UserID: userWithPlayer.ID, Name: "Test User", Address: "123 Test St"}
	gormDB.Create(&player)

	userWithoutPlayer := database.User{Email: "profile_user5@test.com", Password: "password"}
	gormDB.Create(&userWithoutPlayer)

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
			name:           "Player not found",
			email:          "profile_user5@test.com",
			body:           `{}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid JSON",
			email:          "profile_user4@test.com",
			body:           `{`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/profile", bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(req.Context(), auth.EmailKey, tt.email)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler := UpdateProfile(gormDB)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}

func TestDebugGetUsers(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Clear existing users to ensure a clean state for this test
	gormDB.Exec("DELETE FROM users")

	// Create some dummy users
	user1 := database.User{Email: "user1@test.com", Password: "pass1", Name: "User One", Role: "player"}
	user2 := database.User{Email: "user2@test.com", Password: "pass2", Name: "User Two", Role: "player"}
	adminUser := database.User{Email: "admin@test.com", Password: "adminpass", Name: "Admin User", Role: "admin"}
	gormDB.Create(&user1)
	gormDB.Create(&user2)
	gormDB.Create(&adminUser)

	// Create a request to the DebugGetUsers endpoint
	req, err := http.NewRequest("GET", "/debug/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := DebugGetUsers(gormDB)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var users []database.User
	if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}

	// Basic check for user data
	foundUser1 := false
	foundUser2 := false
	foundAdmin := false
	for _, u := range users {
		if u.Email == user1.Email {
			foundUser1 = true
		}
		if u.Email == user2.Email {
			foundUser2 = true
		}
		if u.Email == adminUser.Email {
			foundAdmin = true
		}
	}

	if !foundUser1 || !foundUser2 || !foundAdmin {
		t.Errorf("expected all created users to be in the response")
	}
}

func TestDebugDeleteUser(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user to delete
	userToDelete := database.User{Email: "delete_me@test.com", Password: "password", Name: "Delete User", Role: "player"}
	gormDB.Create(&userToDelete)

	// Create a request to the DebugDeleteUser endpoint
	req, err := http.NewRequest("DELETE", "/debug/users/delete?email=delete_me@test.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := DebugDeleteUser(gormDB)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify the user is deleted from the database
	var user database.User
	if result := gormDB.Where("email = ?", "delete_me@test.com").First(&user); result.Error == nil {
		t.Errorf("user was not deleted from the database")
	}
}
