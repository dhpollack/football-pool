// Package handlers provides HTTP request handlers for week-related operations in the football pool application.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dhpollack/football-pool/internal/api"
	"github.com/dhpollack/football-pool/internal/database"
	"gorm.io/gorm"
)

// ListWeeks handles listing all weeks.
func ListWeeks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var weeks []database.Week
		if result := db.Find(&weeks); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to fetch weeks"})
			return
		}

		// Convert to API response
		weekResponses := make([]api.WeekResponse, len(weeks))
		for i, week := range weeks {
			weekResponses[i] = api.WeekToResponse(week)
		}

		_ = json.NewEncoder(w).Encode(api.WeekListResponse{Weeks: weekResponses})
	}
}

// CreateWeek handles creation of new week records.
func CreateWeek(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var weekRequest api.WeekRequest
		if err := json.NewDecoder(r.Body).Decode(&weekRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request body"})
			return
		}

		// Convert to database model
		week, err := api.WeekFromRequest(weekRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: err.Error()})
			return
		}

		if result := db.Create(&week); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to create week"})
			return
		}

		// Convert to API response
		weekResponse := api.WeekToResponse(week)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(weekResponse)
	}
}

// extractIDFromPath extracts the ID from URL path using PathValue.
func extractIDFromPath(r *http.Request) (uint, error) {
	// Use PathValue to get the "id" parameter from the path
	idStr := extractPathParam(r, "id")
	if idStr == "" {
		return 0, errors.New("no ID found in path")
	}

	// Parse the ID
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

// UpdateWeek handles updating existing week records.
func UpdateWeek(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract week ID from URL path
		id, err := extractIDFromPath(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid week ID"})
			return
		}

		var weekRequest api.WeekRequest
		if err := json.NewDecoder(r.Body).Decode(&weekRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request body"})
			return
		}

		// Find existing week
		var week database.Week
		if result := db.First(&week, id); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Week not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to fetch week"})
			}
			return
		}

		// Update week fields
		week.WeekNumber = weekRequest.WeekNumber
		week.Season = weekRequest.Season
		week.WeekStartTime = weekRequest.WeekStartTime
		week.WeekEndTime = weekRequest.WeekEndTime
		week.IsActive = weekRequest.IsActive

		if result := db.Save(&week); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to update week"})
			return
		}

		// Convert to API response
		weekResponse := api.WeekToResponse(week)

		_ = json.NewEncoder(w).Encode(weekResponse)
	}
}

// DeleteWeek handles deletion of week records.
func DeleteWeek(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract week ID from URL path
		id, err := extractIDFromPath(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid week ID"})
			return
		}

		// Find existing week
		var week database.Week
		if result := db.First(&week, id); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Week not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to fetch week"})
			}
			return
		}

		if result := db.Delete(&week); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to delete week"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ActivateWeek handles activating a week and deactivating all others.
func ActivateWeek(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract week ID from URL path
		id, err := extractIDFromPath(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid week ID"})
			return
		}

		// Find existing week
		var week database.Week
		if result := db.First(&week, id); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Week not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to fetch week"})
			}
			return
		}

		// Start a transaction to ensure atomicity
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Deactivate all weeks
		if result := tx.Model(&database.Week{}).Where("is_active = ?", true).Update("is_active", false); result.Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to deactivate other weeks"})
			return
		}

		// Activate the selected week
		week.IsActive = true
		if result := tx.Save(&week); result.Error != nil {
			tx.Rollback()
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to activate week"})
			return
		}

		if err := tx.Commit().Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to commit transaction"})
			return
		}

		// Convert to API response
		weekResponse := api.WeekToResponse(week)

		_ = json.NewEncoder(w).Encode(weekResponse)
	}
}
