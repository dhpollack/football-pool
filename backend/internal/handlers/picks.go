package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
	"gorm.io/gorm"
)

func GetPicks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var picks []database.Pick
		if result := db.Where("user_id = ?", user.ID).Find(&picks); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(picks); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func SubmitPicks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var picks []database.Pick
		if err := json.NewDecoder(r.Body).Decode(&picks); err != nil || len(picks) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for i := range picks {
			picks[i].UserID = user.ID
		}

		if result := db.Create(&picks); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func AdminSubmitPicks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var picks []database.Pick
		if err := json.NewDecoder(r.Body).Decode(&picks); err != nil || len(picks) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if result := db.Create(&picks); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// AdminListPicks lists all picks with optional filtering
func AdminListPicks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse query parameters
		query := db.Model(&database.Pick{}).Preload("User").Preload("Game")
		
		// Filter by user ID
		if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
			if userID, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
				query = query.Where("user_id = ?", userID)
			}
		}
		
		// Filter by game ID
		if gameIDStr := r.URL.Query().Get("game_id"); gameIDStr != "" {
			if gameID, err := strconv.ParseUint(gameIDStr, 10, 32); err == nil {
				query = query.Where("game_id = ?", gameID)
			}
		}
		
		// Filter by week (join with games table)
		if weekStr := r.URL.Query().Get("week"); weekStr != "" {
			if week, err := strconv.Atoi(weekStr); err == nil {
				query = query.Joins("JOIN games ON picks.game_id = games.id").Where("games.week = ?", week)
			}
		}
		
		// Filter by season (join with games table)
		if seasonStr := r.URL.Query().Get("season"); seasonStr != "" {
			if season, err := strconv.Atoi(seasonStr); err == nil {
				query = query.Joins("JOIN games ON picks.game_id = games.id").Where("games.season = ?", season)
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
		
		// Get picks with pagination
		var picks []database.Pick
		if err := query.Offset(offset).Limit(limit).Order("picks.created_at DESC").Find(&picks).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		// Response with pagination metadata
		response := map[string]interface{}{
			"picks": picks,
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

// AdminGetPicksByWeek gets all picks for a specific week
func AdminGetPicksByWeek(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Extract week from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/picks/week/")
		weekStr := strings.Split(path, "/")[0]
		
		week, err := strconv.Atoi(weekStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid week"})
			return
		}
		
		// Optional season filter
		season := 2024 // Default season
		if seasonStr := r.URL.Query().Get("season"); seasonStr != "" {
			if s, err := strconv.Atoi(seasonStr); err == nil {
				season = s
			}
		}
		
		// Get picks for the week with user and game details
		var picks []database.Pick
		if err := db.Preload("User").Preload("Game").
			Joins("JOIN games ON picks.game_id = games.id").
			Where("games.week = ? AND games.season = ?", week, season).
			Order("picks.user_id, games.id").
			Find(&picks).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(picks); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// AdminGetPicksByUser gets all picks for a specific user
func AdminGetPicksByUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Extract user ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/picks/user/")
		userIDStr := strings.Split(path, "/")[0]
		
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
			return
		}
		
		// Optional season filter
		query := db.Preload("User").Preload("Game").Where("user_id = ?", userID)
		if seasonStr := r.URL.Query().Get("season"); seasonStr != "" {
			if season, err := strconv.Atoi(seasonStr); err == nil {
				query = query.Joins("JOIN games ON picks.game_id = games.id").Where("games.season = ?", season)
			}
		}
		
		var picks []database.Pick
		if err := query.Order("created_at DESC").Find(&picks).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		if err := json.NewEncoder(w).Encode(picks); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// AdminDeletePick deletes a specific pick
func AdminDeletePick(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Extract pick ID from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/admin/picks/")
		idStr := strings.Split(path, "/")[0]
		
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid pick ID"})
			return
		}
		
		// Check if pick exists
		var pick database.Pick
		if err := db.First(&pick, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "Pick not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			}
			return
		}
		
		// Delete the pick
		if err := db.Delete(&pick).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete pick"})
			return
		}
		
		w.WriteHeader(http.StatusNoContent)
	}
}
