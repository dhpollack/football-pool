package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
	"gorm.io/gorm"
	"log/slog"
)

func GetProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var player database.Player
		if result := db.Where("user_id = ?", user.ID).First(&player); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(player); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func UpdateProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var player database.Player
		if result := db.Where("user_id = ?", user.ID).First(&player); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var updates struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		player.Name = updates.Name
		player.Address = updates.Address

		if result := db.Save(&player); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(player); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func DebugGetUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []database.User
		if result := db.Find(&users); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Debug("Found users:", "count", len(users))
		for _, user := range users {
			slog.Debug("User:", "email", user.Email)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func DeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		slog.Debug("Attempting to delete user:", "email", email)
		if email == "" {
			slog.Debug("Error: Email is empty for delete request.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// First, find the user to get their ID
		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			slog.Debug("User not found:", "email", email, "error", result.Error)
			w.WriteHeader(http.StatusNotFound)
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
			return
		}
		slog.Debug("User and associated player deleted successfully:", "email", email)
		w.WriteHeader(http.StatusOK)
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
			return
		}
		
		// Get users with pagination
		var users []database.User
		if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		// Create response with user details including player info
		type UserResponse struct {
			database.User
			PlayerName    string `json:"player_name,omitempty"`
			PlayerAddress string `json:"player_address,omitempty"`
		}
		
		var response []UserResponse
		for _, user := range users {
			userResp := UserResponse{User: user}
			
			// Get player info if exists
			var player database.Player
			if err := db.Where("user_id = ?", user.ID).First(&player).Error; err == nil {
				userResp.PlayerName = player.Name
				userResp.PlayerAddress = player.Address
			}
			
			response = append(response, userResp)
		}
		
		// Response with pagination metadata
		result := map[string]interface{}{
			"users": response,
			"pagination": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		}
		
		if err := json.NewEncoder(w).Encode(result); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
			return
		}
		
		// Get user with player details
		var user database.User
		if err := db.Preload("Player").First(&user, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
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
		
		// Response with user details and stats
		response := map[string]interface{}{
			"user":       user,
			"pick_count": pickCount,
			"total_wins": totalWins,
		}
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
			return
		}
		
		// Check if user exists
		var existingUser database.User
		if err := db.First(&existingUser, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			}
			return
		}
		
		// Parse update data
		var updateData struct {
			Name    string `json:"name,omitempty"`
			Email   string `json:"email,omitempty"`
			Role    string `json:"role,omitempty"`
			PlayerName    string `json:"player_name,omitempty"`
			PlayerAddress string `json:"player_address,omitempty"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
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
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update user"})
			return
		}
		
		// Update player info if provided
		if updateData.PlayerName != "" || updateData.PlayerAddress != "" {
			var player database.Player
			if err := db.Where("user_id = ?", id).First(&player).Error; err != nil {
				// Create player if doesn't exist
				player = database.Player{
					UserID: uint(id),
					Name:   updateData.PlayerName,
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
		
		// Return updated user
		db.Preload("Player").First(&existingUser, id)
		if err := json.NewEncoder(w).Encode(existingUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
