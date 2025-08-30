package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/david/football-pool/internal/database"
	"gorm.io/gorm"
)

func CreateGame(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var games []database.Game
		if err := json.NewDecoder(r.Body).Decode(&games); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if result := db.Create(&games); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func GetGames(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		weekStr := r.URL.Query().Get("week")
		seasonStr := r.URL.Query().Get("season")

		if weekStr == "" || seasonStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		week, err := strconv.Atoi(weekStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		season, err := strconv.Atoi(seasonStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var games []database.Game
		if result := db.Where("week = ? AND season = ?", week, season).Find(&games); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(games); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// AdminListGames lists all games with optional pagination and filtering
func AdminListGames(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse query parameters
		query := db.Model(&database.Game{})
		
		// Filter by week if provided
		if weekStr := r.URL.Query().Get("week"); weekStr != "" {
			if week, err := strconv.Atoi(weekStr); err == nil {
				query = query.Where("week = ?", week)
			}
		}
		
		// Filter by season if provided
		if seasonStr := r.URL.Query().Get("season"); seasonStr != "" {
			if season, err := strconv.Atoi(seasonStr); err == nil {
				query = query.Where("season = ?", season)
			}
		}
		
		// Pagination
		page := 1
		limit := 50 // Default limit
		
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
		
		// Get games with pagination
		var games []database.Game
		if err := query.Offset(offset).Limit(limit).Order("season DESC, week DESC, id DESC").Find(&games).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		// Response with pagination metadata
		response := map[string]interface{}{
			"games": games,
			"pagination": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		}
		
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// UpdateGame updates a specific game
func UpdateGame(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Extract game ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/games/")
		idStr := strings.Split(path, "/")[0]
		
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid game ID"})
			return
		}
		
		// Check if game exists
		var existingGame database.Game
		if err := db.First(&existingGame, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Game not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			}
			return
		}
		
		// Parse updated game data
		var updateData database.Game
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
			return
		}
		
		// Update the game
		if err := db.Model(&existingGame).Updates(updateData).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update game"})
			return
		}
		
		// Return updated game
		if err := json.NewEncoder(w).Encode(existingGame); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// DeleteGame deletes a specific game
func DeleteGame(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Extract game ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/games/")
		idStr := strings.Split(path, "/")[0]
		
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid game ID"})
			return
		}
		
		// Check if game exists
		var game database.Game
		if err := db.First(&game, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Game not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			}
			return
		}
		
		// Check if game has associated picks or results
		var pickCount, resultCount int64
		db.Model(&database.Pick{}).Where("game_id = ?", id).Count(&pickCount)
		db.Model(&database.Result{}).Where("game_id = ?", id).Count(&resultCount)
		
		if pickCount > 0 || resultCount > 0 {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Cannot delete game with associated picks or results",
			})
			return
		}
		
		// Delete the game
		if err := db.Delete(&game).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete game"})
			return
		}
		
		w.WriteHeader(http.StatusNoContent)
	}
}
