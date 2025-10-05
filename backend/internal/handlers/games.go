// Package handlers provides HTTP request handlers for game-related operations in the football pool application.
package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dhpollack/football-pool/internal/api"
	"github.com/dhpollack/football-pool/internal/database"
	"gorm.io/gorm"
)

// CreateGame handles creation of new game records.
func CreateGame(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var gameRequests []api.GameRequest
		if err := json.NewDecoder(r.Body).Decode(&gameRequests); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid request body"})
			return
		}

		// Convert to database models
		games := make([]database.Game, len(gameRequests))
		for i, req := range gameRequests {
			game, err := api.GameFromRequest(req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: err.Error()})
				return
			}
			games[i] = game
		}

		if result := db.Create(&games); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to create games"})
			return
		}

		// Convert to API response
		gameResponses := make([]api.GameResponse, len(games))
		for i, game := range games {
			gameResponses[i] = api.GameToResponse(game)
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(gameResponses)
	}
}

// GetGames handles retrieval of game records with optional week and season filtering.
func GetGames(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		weekStr := r.URL.Query().Get("week")
		seasonStr := r.URL.Query().Get("season")

		if weekStr == "" || seasonStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Week and season parameters are required"}); err != nil {
				slog.Error("Error encoding error response", "error", err)
			}
			return
		}

		week, err := strconv.Atoi(weekStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid week parameter"}); err != nil {
				slog.Error("Error encoding error response", "error", err)
			}
			return
		}

		season, err := strconv.Atoi(seasonStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid season parameter"})
			return
		}

		var games []database.Game
		if result := db.Where("week = ? AND season = ?", week, season).Find(&games); result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			return
		}

		// Convert to API response
		gameResponses := make([]api.GameResponse, len(games))
		for i, game := range games {
			gameResponses[i] = api.GameToResponse(game)
		}

		// Create proper GameListResponse with pagination
		response := api.GameListResponse{
			Games: gameResponses,
			Pagination: api.PaginationResponse{
				Limit: len(games),
				Page:  1,
				Pages: 1,
				Total: int64(len(games)),
			},
		}

		if err = json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

// AdminListGames lists all games with optional pagination and filtering.
func AdminListGames(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse query parameters
		query := db.Model(&database.Game{})

		// Filter by week if provided
		if weekStr := r.URL.Query().Get("week"); weekStr != "" {
			if week, err := strconv.Atoi(weekStr); err == nil {
				query = query.Where("week = ?", week)
			}
		}

		// Filter by season if provided
		if seasonStr := r.URL.Query().Get("season"); seasonStr != "" {
			if season, err := strconv.Atoi(seasonStr); err == nil {
				query = query.Where("season = ?", season)
			}
		}

		// Pagination
		page := 1
		limit := 50 // Default limit

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		offset := (page - 1) * limit

		// Get total count
		var total int64
		if err := query.Count(&total).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			return
		}

		// Get games with pagination
		var games []database.Game
		if err := query.Offset(offset).Limit(limit).Order("season DESC, week DESC, id DESC").Find(&games).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			return
		}

		// Convert to API response
		gameResponses := make([]api.GameResponse, len(games))
		for i, game := range games {
			gameResponses[i] = api.GameToResponse(game)
		}

		// Create structured response
		response := api.GameListResponse{
			Games: gameResponses,
			Pagination: api.PaginationResponse{
				Page:  page,
				Limit: limit,
				Total: total,
				Pages: (total + int64(limit) - 1) / int64(limit),
			},
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

// UpdateGame updates a specific game.
func UpdateGame(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract game ID from URL path using PathValue
		idStr := extractPathParam(r, "id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "No game ID found in path"})
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid game ID"})
			return
		}

		// Check if game exists
		var existingGame database.Game
		if err := db.First(&existingGame, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Game not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			}
			return
		}

		// Parse updated game data
		var updateRequest api.GameRequest
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid JSON"})
			return
		}

		// Convert to database model
		updateData, err := api.GameFromRequest(updateRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(api.ErrorResponse{Error: err.Error()})
			return
		}

		// Update the game
		if err := db.Model(&existingGame).Updates(updateData).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to update game"})
			return
		}

		// Return updated game as API response
		response := api.GameToResponse(existingGame)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to encode response"})
		}
	}
}

// DeleteGame deletes a specific game.
func DeleteGame(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extract game ID from URL path using PathValue
		idStr := extractPathParam(r, "id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "No game ID found in path"})
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Invalid game ID"})
			return
		}

		// Check if game exists
		var game database.Game
		if err := db.First(&game, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Game not found"})
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Database error"})
			}
			return
		}

		// Check if game has associated picks or results
		var pickCount, resultCount int64
		db.Model(&database.Pick{}).Where("game_id = ?", id).Count(&pickCount)
		db.Model(&database.Result{}).Where("game_id = ?", id).Count(&resultCount)

		if pickCount > 0 || resultCount > 0 {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(api.ErrorResponse{
				Error: "Cannot delete game with associated picks or results",
			})
			return
		}

		// Delete the game
		if err := db.Delete(&game).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ErrorResponse{Error: "Failed to delete game"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
