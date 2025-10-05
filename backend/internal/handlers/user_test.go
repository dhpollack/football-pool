package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dhpollack/football-pool/internal/auth"
	"github.com/dhpollack/football-pool/internal/database"
)

func TestGetProfile(t *testing.T) {
	// Set up test database
	db, err := database.New("sqlite", "file::memory:")
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
	db, err := database.New("sqlite", "file::memory:")
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
	db, err := database.New("sqlite", "file::memory:")
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
	db, err := database.New("sqlite", "file::memory:")
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
	db, err := database.New("sqlite", "file::memory:")
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

	// Create a request to the DeleteUser endpoint with path parameter
	pathParams := map[string]string{"id": fmt.Sprintf("%d", userToDelete.ID)}
	req := createRequestWithPathParams("DELETE", fmt.Sprintf("/admin/users/%d", userToDelete.ID), nil, pathParams)

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
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	tests := []struct {
		name           string
		userID         string
		setupUser      bool
		setupPlayer    bool
		expectedStatus int
	}{
		{
			name:           "Delete user without player",
			userID:         "1",
			setupUser:      true,
			setupPlayer:    false,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Delete non-existent user",
			userID:         "999",
			setupUser:      false,
			setupPlayer:    false,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid user ID format",
			userID:         "invalid",
			setupUser:      false,
			setupPlayer:    false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var userToDelete database.User
			if tt.setupUser {
				userToDelete = database.User{Email: fmt.Sprintf("test_%s@test.com", tt.name), Password: "password", Name: "Test User", Role: "player"}
				gormDB.Create(&userToDelete)

				if tt.setupPlayer {
					player := database.Player{UserID: userToDelete.ID, Name: "Test Player", Address: "123 Test St"}
					gormDB.Create(&player)
				}
			}

			// Create a request to the DeleteUser endpoint with path parameter
			var userID string
			if tt.setupUser {
				userID = fmt.Sprintf("%d", userToDelete.ID)
			} else {
				userID = tt.userID
			}
			pathParams := map[string]string{"id": userID}
			req := createRequestWithPathParams("DELETE", fmt.Sprintf("/admin/users/%s", userID), nil, pathParams)

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
				if result := gormDB.Where("id = ?", userToDelete.ID).First(&user); result.Error == nil {
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
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Seed data
	user1 := database.User{Name: "Admin User", Email: "admin@test.com", Role: "admin"}
	user2 := database.User{Name: "Player User", Email: "player@test.com", Role: "user"}
	user3 := database.User{Name: "Another Player", Email: "another@test.com", Role: "user"}
	gormDB.Create(&user1)
	gormDB.Create(&user2)
	gormDB.Create(&user3)

	// Create player records for some users
	player1 := database.Player{UserID: user1.ID, Name: "Admin Player", Address: "123 Admin St"}
	player2 := database.Player{UserID: user2.ID, Name: "Player One", Address: "456 Player St"}
	gormDB.Create(&player1)
	gormDB.Create(&player2)

	// Create games and picks for stats calculation
	game1 := database.Game{Week: 1, Season: 2024, FavoriteTeam: "Lions", UnderdogTeam: "Chiefs"}
	game2 := database.Game{Week: 2, Season: 2024, FavoriteTeam: "Packers", UnderdogTeam: "Bears"}
	gormDB.Create(&game1)
	gormDB.Create(&game2)

	// Create picks and results for stats
	pick1 := database.Pick{UserID: user1.ID, GameID: game1.ID, Picked: "favorite", Rank: 1}
	pick2 := database.Pick{UserID: user1.ID, GameID: game2.ID, Picked: "underdog", Rank: 2}
	pick3 := database.Pick{UserID: user2.ID, GameID: game1.ID, Picked: "favorite", Rank: 1}
	gormDB.Create(&pick1)
	gormDB.Create(&pick2)
	gormDB.Create(&pick3)

	// Create results (user1 has 1 win, user2 has 0 wins)
	result1 := database.Result{GameID: game1.ID, Outcome: "favorite"}
	result2 := database.Result{GameID: game2.ID, Outcome: "favorite"}
	gormDB.Create(&result1)
	gormDB.Create(&result2)

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
		if len(users) != 3 {
			t.Errorf("expected 3 users, got %d", len(users))
		}

		// Verify pagination data
		pagination := response["pagination"].(map[string]interface{})
		if pagination["page"].(float64) != 1 {
			t.Errorf("expected page 1, got %v", pagination["page"])
		}
		if pagination["limit"].(float64) != 20 {
			t.Errorf("expected limit 20, got %v", pagination["limit"])
		}
		if pagination["total"].(float64) != 3 {
			t.Errorf("expected total 3, got %v", pagination["total"])
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

	t.Run("search by email", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?search=admin", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 1 {
			t.Errorf("expected 1 user, got %d", len(users))
		}
	})

	t.Run("search by player name", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?search=Player", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 2 {
			t.Errorf("expected 2 users, got %d", len(users))
		}
	})

	t.Run("pagination - page 1 with limit 2", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?page=1&limit=2", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 2 {
			t.Errorf("expected 2 users on page 1, got %d", len(users))
		}

		pagination := response["pagination"].(map[string]interface{})
		if pagination["page"].(float64) != 1 {
			t.Errorf("expected page 1, got %v", pagination["page"])
		}
		if pagination["limit"].(float64) != 2 {
			t.Errorf("expected limit 2, got %v", pagination["limit"])
		}
		if pagination["pages"].(float64) != 2 {
			t.Errorf("expected 2 total pages, got %v", pagination["pages"])
		}
	})

	t.Run("pagination - page 2 with limit 2", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?page=2&limit=2", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 1 {
			t.Errorf("expected 1 user on page 2, got %d", len(users))
		}
	})

	t.Run("empty search results", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?search=nonexistent", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})

	t.Run("invalid page parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?page=invalid", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 3 {
			t.Errorf("expected 3 users with invalid page, got %d", len(users))
		}
	})

	t.Run("invalid limit parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?limit=invalid", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 3 {
			t.Errorf("expected 3 users with invalid limit, got %d", len(users))
		}
	})

	t.Run("limit too high", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users?limit=200", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		users := response["users"].([]interface{})
		if len(users) != 3 {
			t.Errorf("expected 3 users with high limit, got %d", len(users))
		}

		pagination := response["pagination"].(map[string]interface{})
		if pagination["limit"].(float64) != 20 {
			t.Errorf("expected limit to be capped at 20, got %v", pagination["limit"])
		}
	})

	t.Run("verify user stats calculation", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/users", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)

		// Check if users field exists and is a slice
		usersInterface, ok := response["users"]
		if !ok {
			t.Fatal("users field not found in response")
		}
		users, ok := usersInterface.([]interface{})
		if !ok {
			t.Fatal("users field is not a slice")
		}

		// Find admin user and verify stats
		var adminUser map[string]interface{}
		for _, u := range users {
			user, ok := u.(map[string]interface{})
			if !ok {
				t.Fatal("user item is not a map")
			}
			if user["email"] == "admin@test.com" {
				adminUser = user
				break
			}
		}

		if adminUser == nil {
			t.Fatal("admin user not found in response")
		}

		// Check pick_count (snake_case in JSON response)
		pickCountInterface, ok := adminUser["pick_count"]
		if !ok {
			t.Fatal("pick_count field not found in admin user")
		}
		pickCount, ok := pickCountInterface.(float64)
		if !ok {
			t.Fatalf("pick_count is not a float64, got %T", pickCountInterface)
		}
		if pickCount != 2 {
			t.Errorf("expected admin user to have 2 picks, got %v", pickCount)
		}

		// Check total_wins (snake_case in JSON response)
		totalWinsInterface, ok := adminUser["total_wins"]
		if !ok {
			t.Fatal("total_wins field not found in admin user")
		}
		totalWins, ok := totalWinsInterface.(float64)
		if !ok {
			t.Fatalf("total_wins is not a float64, got %T", totalWinsInterface)
		}
		if totalWins != 1 {
			t.Errorf("expected admin user to have 1 win, got %v", totalWins)
		}
	})
}

func TestAdminGetUser(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	user := database.User{Name: "Test User", Email: "test@test.com", Role: "user"}
	gormDB.Create(&user)

	// Create router for testing
	mux := http.NewServeMux()
	mux.Handle("GET /api/admin/users/{id}", AdminGetUser(gormDB))
	router := mux

	t.Run("get existing user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/admin/users/%d", user.ID), nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

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
		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

func TestAdminUpdateUser(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	user := database.User{Name: "Original Name", Email: "original@test.com", Role: "user"}
	gormDB.Create(&user)

	handler := AdminUpdateUser(gormDB)

	t.Run("update existing user", func(t *testing.T) {
		updatePayload := []byte(`{"name": "Updated Name", "role": "admin"}`)
		pathParams := map[string]string{"id": fmt.Sprintf("%d", user.ID)}
		req := createRequestWithPathParams("PUT", fmt.Sprintf("/api/admin/users/%d", user.ID), bytes.NewBuffer(updatePayload), pathParams)
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

	t.Run("user not found", func(t *testing.T) {
		updatePayload := []byte(`{"name": "Updated Name"}`)
		pathParams := map[string]string{"id": "999"}
		req := createRequestWithPathParams("PUT", "/api/admin/users/999", bytes.NewBuffer(updatePayload), pathParams)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		updatePayload := []byte(`{"name": "Updated Name"`) // Invalid JSON
		pathParams := map[string]string{"id": fmt.Sprintf("%d", user.ID)}
		req := createRequestWithPathParams("PUT", fmt.Sprintf("/api/admin/users/%d", user.ID), bytes.NewBuffer(updatePayload), pathParams)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestAdminCreateUsers(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	handler := AdminCreateUsers(gormDB)

	password := "password"

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "success",
			body:           fmt.Sprintf(`[{"name": "New User 1", "email": "new1@test.com", "password": "%s", "role": "user"}, {"name": "New User 2", "email": "new2@test.com", "password": "%s", "role": "admin"}]`, password, password),
			expectedStatus: http.StatusCreated,
			expectedCount:  2,
		},
		{
			name:           "invalid json",
			body:           `[{"name": "New User 1"}]`,
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "missing password",
			body:           `[{"name": "New User 3", "email": "new3@test.com", "role": "user"}]`,
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "empty password",
			body:           `[{"name": "New User 4", "email": "new4@test.com", "password": "", "role": "user"}]`,
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "duplicate email",
			body:           `[{"name": "New User 5", "email": "duplicate@test.com", "password": "password", "role": "user"}, {"name": "New User 6", "email": "duplicate@test.com", "password": "password", "role": "user"}]`,
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/admin/users/create", bytes.NewBuffer([]byte(tt.body)))
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response []map[string]interface{}
				json.NewDecoder(rr.Body).Decode(&response)
				if len(response) != tt.expectedCount {
					t.Errorf("expected %d users, got %d", tt.expectedCount, len(response))
				}

				// Check if users were actually created
				var count int64
				gormDB.Model(&database.User{}).Count(&count)
				if count != int64(tt.expectedCount) {
					t.Errorf("expected %d users in db, got %d", tt.expectedCount, count)
				}
			}
		})
	}
}

func TestUpdatePlayerInfo(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a test user
	user := database.User{Email: "player_test@test.com", Password: "password"}
	gormDB.Create(&user)

	t.Run("create new player", func(t *testing.T) {
		updatePlayerInfo(gormDB, user.ID, "Test Player", "123 Test St")

		var player database.Player
		if err := gormDB.Where("user_id = ?", user.ID).First(&player).Error; err != nil {
			t.Fatalf("Failed to find player: %v", err)
		}

		if player.Name != "Test Player" {
			t.Errorf("expected player name to be 'Test Player', got %s", player.Name)
		}
		if player.Address != "123 Test St" {
			t.Errorf("expected player address to be '123 Test St', got %s", player.Address)
		}
	})

	t.Run("update existing player name only", func(t *testing.T) {
		updatePlayerInfo(gormDB, user.ID, "Updated Player", "")

		var player database.Player
		if err := gormDB.Where("user_id = ?", user.ID).First(&player).Error; err != nil {
			t.Fatalf("Failed to find player: %v", err)
		}

		if player.Name != "Updated Player" {
			t.Errorf("expected player name to be 'Updated Player', got %s", player.Name)
		}
		// Address should remain unchanged
		if player.Address != "123 Test St" {
			t.Errorf("expected player address to remain '123 Test St', got %s", player.Address)
		}
	})

	t.Run("update existing player address only", func(t *testing.T) {
		updatePlayerInfo(gormDB, user.ID, "", "456 New St")

		var player database.Player
		if err := gormDB.Where("user_id = ?", user.ID).First(&player).Error; err != nil {
			t.Fatalf("Failed to find player: %v", err)
		}

		// Name should remain unchanged
		if player.Name != "Updated Player" {
			t.Errorf("expected player name to remain 'Updated Player', got %s", player.Name)
		}
		if player.Address != "456 New St" {
			t.Errorf("expected player address to be '456 New St', got %s", player.Address)
		}
	})

	t.Run("update both name and address", func(t *testing.T) {
		updatePlayerInfo(gormDB, user.ID, "Final Player", "789 Final St")

		var player database.Player
		if err := gormDB.Where("user_id = ?", user.ID).First(&player).Error; err != nil {
			t.Fatalf("Failed to find player: %v", err)
		}

		if player.Name != "Final Player" {
			t.Errorf("expected player name to be 'Final Player', got %s", player.Name)
		}
		if player.Address != "789 Final St" {
			t.Errorf("expected player address to be '789 Final St', got %s", player.Address)
		}
	})
}

func TestDeleteUserByEmail(t *testing.T) {
	// Set up test database
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a user to delete
	userToDelete := database.User{Email: "delete_by_email@test.com", Password: "password", Name: "Delete Email User", Role: "player"}
	gormDB.Create(&userToDelete)

	// Create associated player record
	playerToDelete := database.Player{UserID: userToDelete.ID, Name: "Delete Email Player", Address: "123 Delete Email St"}
	gormDB.Create(&playerToDelete)

	// Create a request to the DeleteUserByEmail endpoint with query parameter
	req, err := http.NewRequest("DELETE", "/api/admin/users/delete?email=delete_by_email@test.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()
	handler := DeleteUserByEmail(gormDB)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code (DeleteUserByEmail returns NoContent for successful deletion)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Verify the user is deleted from the database
	var user database.User
	if result := gormDB.Where("email = ?", "delete_by_email@test.com").First(&user); result.Error == nil {
		t.Errorf("user was not deleted from the database")
	}

	// Verify the associated player is also deleted from the database
	var player database.Player
	if result := gormDB.Where("user_id = ?", userToDelete.ID).First(&player); result.Error == nil {
		t.Errorf("associated player was not deleted from the database")
	}
}

func TestDeleteUserByEmailEdgeCases(t *testing.T) {
	// Set up test database
	db, err := database.New("sqlite", "file::memory:")
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
		expectedError  string
	}{
		{
			name:           "Missing email parameter",
			email:          "",
			setupUser:      false,
			setupPlayer:    false,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email parameter is required",
		},
		{
			name:           "User not found by email",
			email:          "nonexistent@test.com",
			setupUser:      false,
			setupPlayer:    false,
			expectedStatus: http.StatusNotFound,
			expectedError:  "User not found",
		},
		{
			name:           "Delete user without player",
			email:          "no_player@test.com",
			setupUser:      true,
			setupPlayer:    false,
			expectedStatus: http.StatusNoContent,
			expectedError:  "",
		},
		{
			name:           "Delete user with player",
			email:          "with_player@test.com",
			setupUser:      true,
			setupPlayer:    true,
			expectedStatus: http.StatusNoContent,
			expectedError:  "",
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

			// Create a request to the DeleteUserByEmail endpoint with query parameter
			url := "/api/admin/users/delete"
			if tt.email != "" {
				url = fmt.Sprintf("%s?email=%s", url, tt.email)
			}
			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a ResponseRecorder
			rr := httptest.NewRecorder()
			handler := DeleteUserByEmail(gormDB)

			// Serve the request
			handler.ServeHTTP(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check error response for error cases
			if tt.expectedStatus >= 400 && tt.expectedStatus < 500 {
				var response map[string]interface{}
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if errorMsg, ok := response["error"].(string); ok && errorMsg != tt.expectedError {
					t.Errorf("handler returned unexpected error message: got %v want %v",
						errorMsg, tt.expectedError)
				}
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
