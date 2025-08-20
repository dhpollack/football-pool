package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
	"log/slog"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(auth.EmailKey).(string)

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

	if err := json.NewEncoder(w).Encode(player); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(auth.EmailKey).(string)

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

	if err := json.NewEncoder(w).Encode(player); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DebugGetUsers(w http.ResponseWriter, r *http.Request) {
	var users []database.User
	if result := database.DB.Find(&users); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DebugDeleteUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	slog.Debug("Attempting to delete user:", "email", email)
	if email == "" {
		slog.Debug("Error: Email is empty for delete request.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if result := database.DB.Where("email = ?", email).Delete(&database.User{}); result.Error != nil {
		slog.Debug("Error deleting user:", "email", email, "error", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Debug("User deleted successfully:", "email", email)
	w.WriteHeader(http.StatusOK)
}