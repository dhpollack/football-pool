package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/david/football-pool/internal/database"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)

	var user database.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var player database.Player
	if result := database.DB.Where("user_id = ?", user.ID).First(&player); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(player)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)

	var user database.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var player database.Player
	if result := database.DB.Where("user_id = ?", user.ID).First(&player); result.Error != nil {
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

	if result := database.DB.Save(&player); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(player)
}
