package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/david/football-pool/internal/api"
	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "User not found"})
			return
		}

		var player database.Player
		if result := db.Where("user_id = ?", user.ID).First(&player); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Player profile not found"})
			return
		}

		// Convert to API response
		response := api.PlayerToResponse(player)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

func UpdateProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "User not found"})
			return
		}

		var player database.Player
		if result := db.Where("user_id = ?", user.ID).First(&player); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Player profile not found"})
			return
		}

		var updateRequest api.PlayerRequest
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request body"})
			return
		}

		player.Name = updateRequest.Name
		player.Address = updateRequest.Address

		if result := db.Save(&player); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to update profile"})
			return
		}

		// Convert to API response
		response := api.PlayerToResponse(player)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		email := r.URL.Query().Get("email")
		slog.Debug("Attempting to delete user:", "email", email)
		if email == "" {
			slog.Debug("Error: Email is empty for delete request.")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Email parameter is required"})
			return
		}

		// First, find the user to get their ID
		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			slog.Debug("User not found:", "email", email, "error", result.Error)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "User not found"})
			return
		}

		// Delete associated player record first (due to foreign key constraints)
		if result := db.Unscoped().Where("user_id = ?", user.ID).Delete(&database.Player{}); result.Error != nil {
			slog.Debug("Error deleting player:", "email", email, "error", result.Error)
			// Continue with user deletion even if player deletion fails
		}

		// Delete the user record
		if result := db.Unscoped().Where("email = ?", email).Delete(&database.User{}); result.Error != nil {
			slog.Debug("Error deleting user:", "email", email, "error", result.Error)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to delete user"})
			return
		}
		slog.Debug("User and associated player deleted successfully:", "email", email)
		w.WriteHeader(http.StatusNoContent)
	}
}

// AdminListUsers lists all users with pagination and search
func AdminListUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse query parameters
		query := db.Model(&database.User{}).Preload("Player")

		// Search by email or name
		if search := r.URL.Query().Get("search"); search != "" {
			query = query.Where("email LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
		}

		// Filter by role
		if role := r.URL.Query().Get("role"); role != "" {
			query = query.Where("role = ?", role)
		}

		// Pagination
		page := 1
		limit := 20 // Default limit for users

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		offset := (page - 1) * limit

		// Get total count
		var total int64
		if err := query.Count(&total).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			return
		}

		// Get users with pagination
		var users []database.User
		if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			return
		}

		// Calculate stats for each user and convert to API response
		usersWithStats := make([]api.UserWithStats, len(users))
		for i, user := range users {
			// Get pick count and wins for this user
			var pickCount, totalWins int64
			db.Model(&database.Pick{}).Where("user_id = ?", user.ID).Count(&pickCount)
			db.Table("picks").
				Joins("JOIN results ON picks.game_id = results.game_id").
				Where("picks.user_id = ? AND picks.picked = results.outcome", user.ID).
				Count(&totalWins)

			usersWithStats[i] = api.UserWithStatsFromUser(user, int(pickCount), int(totalWins))
		}

		// Create structured response
		response := api.UserListResponse{
			Users: usersWithStats,
			Pagination: api.PaginationResponse{
				Page:  page,
				Limit: limit,
				Total: total,
				Pages: (total + int64(limit) - 1) / int64(limit),
			},
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

// AdminGetUser gets details for a specific user
func AdminGetUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract user ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
		idStr := strings.Split(path, "/")[0]

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Get user with player details
		var user database.User
		if err := db.Preload("Player").First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "User not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			}
			return
		}

		// Get additional stats for the user
		var pickCount, totalWins int64
		db.Model(&database.Pick{}).Where("user_id = ?", id).Count(&pickCount)

		// Count wins (this is a simplified version - you may want to enhance this)
		db.Table("picks").
			Joins("JOIN results ON picks.game_id = results.game_id").
			Where("picks.user_id = ? AND picks.picked = results.outcome", id).
			Count(&totalWins)

		// Convert to API response
		response := api.UserWithStatsFromUser(user, int(pickCount), int(totalWins))

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

// AdminUpdateUser updates a specific user
func AdminUpdateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract user ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/users/")
		idStr := strings.Split(path, "/")[0]

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid user ID"})
			return
		}

		// Check if user exists
		var existingUser database.User
		if err := db.First(&existingUser, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "User not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			}
			return
		}

		// Parse update data with both user and player fields
		var updateData struct {
			Name          string `json:"name,omitempty"`
			Email         string `json:"email,omitempty"`
			Role          string `json:"role,omitempty"`
			PlayerName    string `json:"player_name,omitempty"`
			PlayerAddress string `json:"player_address,omitempty"`
		}

		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Update user fields
		if updateData.Name != "" {
			existingUser.Name = updateData.Name
		}
		if updateData.Email != "" {
			existingUser.Email = updateData.Email
		}
		if updateData.Role != "" {
			existingUser.Role = updateData.Role
		}

		// Save user updates
		if err := db.Save(&existingUser).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to update user"})
			return
		}

		// Update player info if provided
		if updateData.PlayerName != "" || updateData.PlayerAddress != "" {
			var player database.Player
			if err := db.Where("user_id = ?", id).First(&player).Error; err != nil {
				// Create player if doesn't exist
				player = database.Player{
					UserID:  uint(id),
					Name:    updateData.PlayerName,
					Address: updateData.PlayerAddress,
				}
				db.Create(&player)
			} else {
				// Update existing player
				if updateData.PlayerName != "" {
					player.Name = updateData.PlayerName
				}
				if updateData.PlayerAddress != "" {
					player.Address = updateData.PlayerAddress
				}
				db.Save(&player)
			}
		}

		// Return updated user as API response
		db.Preload("Player").First(&existingUser, id)
		response := api.UserToResponse(existingUser)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

// AdminCreateUsers creates multiple users
func AdminCreateUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var userRequests []api.UserRequest
		if err := json.NewDecoder(r.Body).Decode(&userRequests); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request body"})
			return
		}

		var userResponses []api.UserResponse

		tx := db.Begin()
		if tx.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to start transaction"})
			return
		}

		for _, req := range userRequests {
			user, err := api.UserFromRequest(req)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: err.Error()})
				return
			}

			if req.Password == nil || *req.Password == "" {
				tx.Rollback()
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Password is required for all users"})
				return
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), 8)
			if err != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to hash password"})
				return
			}

			user.Password = string(hashedPassword)

			if result := tx.Create(&user); result.Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to create user"})
				return
			}

			player := database.Player{
				UserID: user.ID,
				Name:   user.Name,
			}
			if result := tx.Create(&player); result.Error != nil {
				tx.Rollback()
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to create player"})
				return
			}
			user.Player = player
			userResponses = append(userResponses, api.UserToResponse(user))
		}

		if err := tx.Commit().Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to commit transaction"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(userResponses)
	}
}
