package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
