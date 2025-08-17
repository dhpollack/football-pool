package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/david/football-pool/internal/database"
)

func GetWeeklyResults(w http.ResponseWriter, r *http.Request) {
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

	var gameIDs []uint
	for _, game := range games {
		gameIDs = append(gameIDs, game.ID)
	}

	var results []database.Result
	if result := database.DB.Where("game_id IN ?", gameIDs).Find(&results); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultsMap := make(map[uint]database.Result)
	for _, result := range results {
		resultsMap[result.GameID] = result
	}

	var picks []database.Pick
	if result := database.DB.Where("game_id IN ?", gameIDs).Find(&picks); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	playerScores := make(map[uint]float32)
	for _, pick := range picks {
		result, ok := resultsMap[pick.GameID]
		if !ok {
			continue
		}

		if (pick.PickedTeam == "favorite" && result.Outcome == "favorite") || (pick.PickedTeam == "underdog" && result.Outcome == "underdog") {
			playerScores[pick.UserID] += float32(pick.Rank)
		} else if result.Outcome == "push" {
			playerScores[pick.UserID] += float32(pick.Rank) / 2
		}
	}

	type WeeklyResult struct {
		PlayerID   uint    `json:"player_id"`
		PlayerName string  `json:"player_name"`
		Score      float32 `json:"score"`
	}

	var weeklyResults []WeeklyResult
	for userID, score := range playerScores {
		var user database.User
		if result := database.DB.First(&user, userID); result.Error != nil {
			continue
		}
		weeklyResults = append(weeklyResults, WeeklyResult{PlayerID: userID, PlayerName: user.Name, Score: score})
	}

	if err := json.NewEncoder(w).Encode(weeklyResults); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SubmitResult(w http.ResponseWriter, r *http.Request) {
	var result database.Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var game database.Game
	if dbResult := database.DB.First(&game, result.GameID); dbResult.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if float32(result.FavoriteScore)-game.Spread > float32(result.UnderdogScore) {
		result.Outcome = "favorite"
	} else if float32(result.UnderdogScore) > float32(result.FavoriteScore)-game.Spread {
		result.Outcome = "underdog"
	} else {
		result.Outcome = "push"
	}

	if dbResult := database.DB.Create(&result); dbResult.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetSeasonResults(w http.ResponseWriter, r *http.Request) {
	seasonStr := r.URL.Query().Get("season")

	if seasonStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	season, err := strconv.Atoi(seasonStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var games []database.Game
	if result := database.DB.Where("season = ?", season).Find(&games); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var gameIDs []uint
	for _, game := range games {
		gameIDs = append(gameIDs, game.ID)
	}

	var results []database.Result
	if result := database.DB.Where("game_id IN ?", gameIDs).Find(&results); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultsMap := make(map[uint]database.Result)
	for _, result := range results {
		resultsMap[result.GameID] = result
	}

	var picks []database.Pick
	if result := database.DB.Where("game_id IN ?", gameIDs).Find(&picks); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	playerScores := make(map[uint]float32)
	for _, pick := range picks {
		result, ok := resultsMap[pick.GameID]
		if !ok {
			continue
		}

		if (pick.PickedTeam == "favorite" && result.Outcome == "favorite") || (pick.PickedTeam == "underdog" && result.Outcome == "underdog") {
			playerScores[pick.UserID] += float32(pick.Rank)
		} else if result.Outcome == "push" {
			playerScores[pick.UserID] += float32(pick.Rank) / 2
		}
	}

	type SeasonResult struct {
		PlayerID   uint    `json:"player_id"`
		PlayerName string  `json:"player_name"`
		Score      float32 `json:"score"`
	}

	var seasonResults []SeasonResult
	for userID, score := range playerScores {
		var user database.User
		if result := database.DB.First(&user, userID); result.Error != nil {
			continue
		}
		seasonResults = append(seasonResults, SeasonResult{PlayerID: userID, PlayerName: user.Name, Score: score})
	}

	if err := json.NewEncoder(w).Encode(seasonResults); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
