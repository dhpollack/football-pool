package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/david/football-pool/internal/api"
	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
	"gorm.io/gorm"
)

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
