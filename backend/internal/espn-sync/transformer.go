// Package espnsync provides functionality to sync season schedule and game scores from the ESPN API.
package espnsync

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	apiespn "github.com/dhpollack/football-pool/internal/api-espn"
	"github.com/dhpollack/football-pool/internal/database"
)

// Transformer handles the transformation of ESPN API data to database models.
type Transformer struct {
	db *database.Database
}

// NewTransformer creates a new Transformer instance.
func NewTransformer(db *database.Database) *Transformer {
	return &Transformer{db: db}
}

// TransformEvent transforms an ESPN Event to database Game and Result models.
func (t *Transformer) TransformEvent(event apiespn.Event, season, week int) (*database.Game, *database.Result, error) {
	if event.Competitions == nil || len(*event.Competitions) == 0 {
		return nil, nil, nil
	}

	competition := (*event.Competitions)[0]

	// Extract teams from competition
	favoriteTeam, underdogTeam, err := t.extractTeams(competition)
	if err != nil {
		return nil, nil, err
	}

	// Extract game time
	startTime, err := t.extractStartTime(competition, event)
	if err != nil {
		return nil, nil, err
	}

	// Create Game model
	game := &database.Game{
		Week:         week,
		Season:       season,
		FavoriteTeam: favoriteTeam,
		UnderdogTeam: underdogTeam,
		Spread:       0.0, // ESPN API doesn't provide spread information
		StartTime:    startTime,
	}

	// Create Result model if scores are available
	var result *database.Result
	if competition.Competitors != nil && len(*competition.Competitors) == 2 {
		result = t.extractResult(competition)
	}

	return game, result, nil
}

// extractTeams extracts favorite and underdog teams from a competition.
// Note: ESPN API doesn't explicitly state which team is favorite/underdog.
// This implementation assumes the home team is the favorite.
func (t *Transformer) extractTeams(competition apiespn.Competition) (string, string, error) {
	if competition.Competitors == nil || len(*competition.Competitors) != 2 {
		return "", "", fmt.Errorf("invalid competition: expected 2 competitors, got %d", len(*competition.Competitors))
	}

	var homeTeam, awayTeam string
	competitors := *competition.Competitors

	for _, competitor := range competitors {
		if competitor.Team == nil || competitor.Team.DisplayName == nil {
			return "", "", fmt.Errorf("competitor missing team name")
		}

		if competitor.HomeAway != nil && *competitor.HomeAway == "home" {
			homeTeam = *competitor.Team.DisplayName
		} else {
			awayTeam = *competitor.Team.DisplayName
		}
	}

	// Validate that we have both team names
	if homeTeam == "" || awayTeam == "" {
		return "", "", fmt.Errorf("missing team names: home='%s', away='%s'", homeTeam, awayTeam)
	}

	// Use home team as favorite (this is an assumption)
	favoriteTeam := homeTeam
	underdogTeam := awayTeam

	return favoriteTeam, underdogTeam, nil
}

// extractStartTime extracts the game start time from competition and event data.
func (t *Transformer) extractStartTime(competition apiespn.Competition, event apiespn.Event) (time.Time, error) {
	// Try competition date first
	if competition.Date != nil {
		return competition.Date.Time, nil
	}

	// Fallback to event date
	if event.Date != nil {
		return event.Date.Time, nil
	}

	// Default to current time if no date available
	return time.Time{}, fmt.Errorf("unable to determine the start date from competition or event")
}

// extractResult extracts game result from competition data.
func (t *Transformer) extractResult(competition apiespn.Competition) *database.Result {
	if competition.Competitors == nil || len(*competition.Competitors) != 2 {
		return nil
	}

	competitors := *competition.Competitors
	var favoriteScore, underdogScore int
	var outcome string

	// Debug: log competitor details
	slog.Debug("Extracting result from competition", "competitors_count", len(competitors))

	// Extract scores and determine outcome
	for i, competitor := range competitors {
		if competitor.Score == nil {
			slog.Debug("Competitor has no score", "index", i, "homeAway", competitor.HomeAway)
			continue
		}

		score := 0
		if *competitor.Score != "" {
			score, _ = strconv.Atoi(*competitor.Score)
		}

		slog.Debug("Competitor score", "index", i, "homeAway", competitor.HomeAway, "score", score)

		// Determine if this competitor is the favorite or underdog
		// This is based on our earlier assumption that home team = favorite
		if competitor.HomeAway != nil && *competitor.HomeAway == "home" {
			favoriteScore = score
		} else {
			underdogScore = score
		}
	}

	slog.Debug("Extracted scores", "favoriteScore", favoriteScore, "underdogScore", underdogScore)

	// Determine outcome
	if favoriteScore > 0 || underdogScore > 0 {
		spread := float32(favoriteScore - underdogScore)
		switch {
		case spread > 0:
			outcome = "favorite"
		case spread < 0:
			outcome = "underdog"
		default:
			outcome = "push"
		}
	}

	if outcome == "" {
		slog.Debug("No outcome determined - returning nil result")
		return nil // No result available yet
	}

	slog.Debug("Result created", "favoriteScore", favoriteScore, "underdogScore", underdogScore, "outcome", outcome)
	return &database.Result{
		FavoriteScore: favoriteScore,
		UnderdogScore: underdogScore,
		Outcome:       outcome,
	}
}

// StoreGameAndResult stores a game and its result in the database.
func (t *Transformer) StoreGameAndResult(game *database.Game, result *database.Result) error {
	// Check if game already exists
	var existingGame database.Game
	err := t.db.GetDB().Where("season = ? AND week = ? AND favorite_team = ? AND underdog_team = ?",
		game.Season, game.Week, game.FavoriteTeam, game.UnderdogTeam).First(&existingGame).Error

	slog.Debug("StoreGameAndResult: checking for existing game", "season", game.Season, "week", game.Week, "favorite", game.FavoriteTeam, "underdog", game.UnderdogTeam, "error", err, "existing_game_id", existingGame.ID)

	if err != nil {
		// Game doesn't exist, create it
		slog.Debug("StoreGameAndResult: creating new game")
		if err := t.db.GetDB().Create(game).Error; err != nil {
			return err
		}
	} else {
		// Game exists, update it
		slog.Debug("StoreGameAndResult: updating existing game", "existingID", existingGame.ID)
		game.ID = existingGame.ID
		if err := t.db.GetDB().Save(game).Error; err != nil {
			return err
		}
	}

	// Store result if available
	if result == nil {
		return nil
	}

	result.GameID = game.ID

	// Check if result already exists
	var existingResult database.Result
	err = t.db.GetDB().Where("game_id = ?", game.ID).First(&existingResult).Error
	if err != nil {
		// Result doesn't exist, create it
		return t.db.GetDB().Create(result).Error
	}

	// Result exists, update it
	result.ID = existingResult.ID
	return t.db.GetDB().Save(result).Error
}
