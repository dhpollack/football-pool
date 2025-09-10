package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
)

func TestGetProfile(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
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
	db, err := database.New("file::memory:")
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
	db, err := database.New("file::memory:")
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
	db, err := database.New("file::memory:")
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

func TestDeleteUser(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user to delete
	userToDelete := database.User{Email: "delete_me@test.com", Password: "password", Name: "Delete User", Role: "player"}
	gormDB.Create(&userToDelete)

	// Create associated player record
	playerToDelete := database.Player{UserID: userToDelete.ID, Name: "Delete Player", Address: "123 Delete St"}
	gormDB.Create(&playerToDelete)

	// Create a request to the DeleteUser endpoint
	req, err := http.NewRequest("DELETE", "/admin/users/delete?email=delete_me@test.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := DeleteUser(gormDB)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code (DeleteUser returns NoContent for successful deletion)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify the user is deleted from the database
	var user database.User
	if result := gormDB.Where("email = ?", "delete_me@test.com").First(&user); result.Error == nil {
		t.Errorf("user was not deleted from the database")
	}

	// Verify the associated player is also deleted from the database
	var player database.Player
	if result := gormDB.Where("user_id = ?", userToDelete.ID).First(&player); result.Error == nil {
		t.Errorf("associated player was not deleted from the database")
	}
}

func TestDeleteUserEdgeCases(t *testing.T) {
	// Set up test database
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	tests := []struct {
		name           string
		email          string
		setupUser      bool
		setupPlayer    bool
		expectedStatus int
	}{
		{
			name:           "Delete user without player",
			email:          "no_player@test.com",
			setupUser:      true,
			setupPlayer:    false,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete non-existent user",
			email:          "nonexistent@test.com",
			setupUser:      false,
			setupPlayer:    false,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Empty email parameter",
			email:          "",
			setupUser:      false,
			setupPlayer:    false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var userToDelete database.User
			if tt.setupUser {
				userToDelete = database.User{Email: tt.email, Password: "password", Name: "Test User", Role: "player"}
				gormDB.Create(&userToDelete)

				if tt.setupPlayer {
					player := database.Player{UserID: userToDelete.ID, Name: "Test Player", Address: "123 Test St"}
					gormDB.Create(&player)
				}
			}

			// Create a request to the DeleteUser endpoint
			url := "/admin/users/delete"
			if tt.email != "" {
				url += "?email=" + tt.email
			}
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder
			rr := httptest.NewRecorder()
			handler := DeleteUser(gormDB)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// If deletion was successful, verify cleanup
			if tt.expectedStatus == http.StatusNoContent && tt.setupUser {
				var user database.User
				if result := gormDB.Where("email = ?", tt.email).First(&user); result.Error == nil {
					t.Errorf("user was not deleted from the database")
				}

				if tt.setupPlayer {
					var player database.Player
					if result := gormDB.Where("user_id = ?", userToDelete.ID).First(&player); result.Error == nil {
						t.Errorf("associated player was not deleted from the database")
					}
				}
			}
		})
	}
}

func TestAdminListUsers(t *testing.T) {
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Seed data
	user1 := database.User{Name: "Admin User", Email: "admin@test.com", Role: "admin"}
	user2 := database.User{Name: "Player User", Email: "player@test.com", Role: "user"}
	gormDB.Create(&user1)
	gormDB.Create(&user2)

	handler := AdminListUsers(gormDB)

	t.Run("list all users", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 2 {
			t.Errorf("expected 2 users, got %d", len(users))
		}
	})

	t.Run("filter by role", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?role=admin", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 1 {
			t.Errorf("expected 1 user, got %d", len(users))
		}
	})
}

func TestAdminGetUser(t *testing.T) {
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	user := database.User{Name: "Test User", Email: "test@test.com", Role: "user"}
	gormDB.Create(&user)

	handler := AdminGetUser(gormDB)

	t.Run("get existing user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/admin/users/%d", user.ID), nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		if response["email"] != user.Email {
			t.Errorf("expected user email %s, got %s", user.Email, response["email"])
		}
	})

	t.Run("user not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users/999", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

func TestAdminUpdateUser(t *testing.T) {
	db, err := database.New("file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	user := database.User{Name: "Original Name", Email: "original@test.com", Role: "user"}
	gormDB.Create(&user)

	handler := AdminUpdateUser(gormDB)

	t.Run("update existing user", func(t *testing.T) {
		updatePayload := []byte(`{"name": "Updated Name", "role": "admin"}`)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/admin/users/%d", user.ID), bytes.NewBuffer(updatePayload))
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var updatedUser database.User
		gormDB.First(&updatedUser, user.ID)
		if updatedUser.Name != "Updated Name" {
			t.Errorf("expected user name to be updated, got %s", updatedUser.Name)
		}
		if updatedUser.Role != "admin" {
			t.Errorf("expected user role to be updated, got %s", updatedUser.Role)
		}
	})
}
