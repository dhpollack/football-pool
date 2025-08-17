package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/david/football-pool/internal/database"
)

func GetGames(w http.ResponseWriter, r *http.Request) {
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
	if result := database.DB.Where("week = ? AND season = ?", week, season).Find(&games); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(games); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
