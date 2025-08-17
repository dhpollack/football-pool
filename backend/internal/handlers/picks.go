package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
)

func GetPicks(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(auth.EmailKey).(string)

	var user database.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var picks []database.Pick
	if result := database.DB.Where("user_id = ?", user.ID).Find(&picks); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(picks); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SubmitPicks(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(auth.EmailKey).(string)

	var user database.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var picks []database.Pick
	if err := json.NewDecoder(r.Body).Decode(&picks); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := range picks {
		picks[i].UserID = user.ID
	}

	if result := database.DB.Create(&picks); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
