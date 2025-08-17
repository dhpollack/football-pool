package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
)

func GetSurvivorPicks(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(auth.EmailKey).(string)

	var user database.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var survivorPicks []database.SurvivorPick
	if result := database.DB.Where("user_id = ?", user.ID).Find(&survivorPicks); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(survivorPicks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SubmitSurvivorPick(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(auth.EmailKey).(string)

	var user database.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var pick database.SurvivorPick
	if err := json.NewDecoder(r.Body).Decode(&pick); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pick.UserID = user.ID

	if result := database.DB.Create(&pick); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
