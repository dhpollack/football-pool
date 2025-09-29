// Package handlers provides HTTP request handlers for survivor pool operations in the football pool application.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dhpollack/football-pool/internal/api"
	"github.com/dhpollack/football-pool/internal/auth"
	"github.com/dhpollack/football-pool/internal/database"
	"gorm.io/gorm"
)

// GetSurvivorPicks handles retrieval of survivor pool picks for the current user.
func GetSurvivorPicks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var survivorPicks []database.SurvivorPick
		if result := db.Where("user_id = ?", user.ID).Find(&survivorPicks); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := make([]api.SurvivorPickResponse, len(survivorPicks))
		for i, pick := range survivorPicks {
			response[i] = api.SurvivorPickToResponse(pick)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// SubmitSurvivorPick handles submission of survivor pool picks.
func SubmitSurvivorPick(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(auth.EmailKey).(string)

		var user database.User
		if result := db.Where("email = ?", email).First(&user); result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var pick database.SurvivorPick
		if err := json.NewDecoder(r.Body).Decode(&pick); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pick.UserID = user.ID

		if result := db.Create(&pick); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := api.SurvivorPickToResponse(pick)

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}
}
