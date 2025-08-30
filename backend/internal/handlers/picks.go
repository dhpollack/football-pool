package handlers

import (
	"encoding/json"
	"net/http"

	//"strconv"
	//"strings"

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
