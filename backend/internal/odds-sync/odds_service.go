// Package oddssync provides functionality to sync game spreads from The Odds API.
package oddssync

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	theoddsapi "github.com/dhpollack/football-pool/internal/api-the-odds-api"
	"github.com/dhpollack/football-pool/internal/config"
	"github.com/dhpollack/football-pool/internal/database"
)

// OddsService orchestrates the fetching and updating of game spreads from The Odds API.
type OddsService struct {
	db           *database.Database
	oddsClient   *theoddsapi.ClientWithResponses
	config       *config.Config
	timeProvider TimeProvider
}

// TimeProvider defines an interface for getting the current time.
// This allows for dependency injection in tests.
type TimeProvider interface {
	Now() time.Time
}

// RealTimeProvider provides the actual current time.
type RealTimeProvider struct{}

// Now returns the current time.
func (r RealTimeProvider) Now() time.Time {
	return time.Now()
}

// NewOddsService creates a new OddsService instance.
func NewOddsService(db *database.Database, config *config.Config) (*OddsService, error) {
	return NewOddsServiceWithTimeProvider(db, config, RealTimeProvider{})
}

// NewOddsServiceWithTimeProvider creates a new OddsService instance with a custom time provider.
// This is primarily for testing purposes.
func NewOddsServiceWithTimeProvider(db *database.Database, config *config.Config, timeProvider TimeProvider) (*OddsService, error) {
	// Create The Odds API client
	client, err := theoddsapi.NewClientWithResponses(config.TheOddsAPI.BaseURL)
	if err != nil {
		return nil, err
	}

	return &OddsService{
		db:           db,
		oddsClient:   client,
		config:       config,
		timeProvider: timeProvider,
	}, nil
}

// GameSpread represents the spread information for a game.
type GameSpread struct {
	HomeTeam string
	AwayTeam string
	Spread   float32
	Favorite string // "Home" or "Away"
	Underdog string // "Home" or "Away"
}

// FetchSpreadsForWeek fetches spreads for games in a specific week and season.
func (s *OddsService) FetchSpreadsForWeek(ctx context.Context, season, week int) ([]GameSpread, error) {
	slog.Info("Fetching spreads for week", "season", season, "week", week)

	// Calculate date range for the week
	weekStart, weekEnd := s.getWeekDateRange(season, week)

	// Fetch odds from The Odds API
	params := &theoddsapi.GetOddsParams{
		ApiKey:  s.config.TheOddsAPI.APIKey,
		Regions: theoddsapi.GetOddsParamsRegions(s.config.TheOddsAPI.Region),
		Markets: func() *theoddsapi.GetOddsParamsMarkets {
			market := theoddsapi.GetOddsParamsMarketsSpreads
			return &market
		}(),
		CommenceTimeFrom: func() *string {
			from := weekStart.Format(time.RFC3339)
			return &from
		}(),
		CommenceTimeTo: func() *string {
			to := weekEnd.Format(time.RFC3339)
			return &to
		}(),
	}

	response, err := s.oddsClient.GetOddsWithResponse(ctx, "americanfootball_nfl", params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch odds: %w", err)
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("the Odds API returned status %d: %s", response.StatusCode(), response.Status())
	}

	if response.JSON200 == nil {
		slog.Error("The Odds API returned nil JSON200 response", "status", response.StatusCode(), "body", string(response.Body))
		return nil, fmt.Errorf("the Odds API returned empty response")
	}

	// Parse the response and extract spreads
	spreads := make([]GameSpread, 0, len(*response.JSON200))
	for _, event := range *response.JSON200 {
		if event.HomeTeam == nil || event.AwayTeam == nil {
			slog.Debug("Skipping event with missing team information")
			continue
		}

		spread, err := s.extractSpreadFromEvent(event)
		if err != nil {
			slog.Warn("Failed to extract spread from event", "home_team", *event.HomeTeam, "away_team", *event.AwayTeam, "error", err)
			continue
		}

		spreads = append(spreads, *spread)
	}

	slog.Info("Fetched spreads from The Odds API", "count", len(spreads))
	return spreads, nil
}

// extractSpreadFromEvent extracts spread information from an odds API event.
func (s *OddsService) extractSpreadFromEvent(event struct {
	AwayTeam   *theoddsapi.AwayTeam `json:"away_team"`
	Bookmakers *[]struct {
		Key        *string    `json:"key,omitempty"`
		LastUpdate *time.Time `json:"last_update,omitempty"`
		Link       *string    `json:"link"`
		Markets    *[]struct {
			Key        *theoddsapi.GetOdds200BookmakersMarketsKey `json:"key,omitempty"`
			LastUpdate *time.Time                                 `json:"last_update,omitempty"`
			Link       *string                                    `json:"link"`
			Outcomes   *[]theoddsapi.Outcome                      `json:"outcomes,omitempty"`
			Sid        *string                                    `json:"sid"`
		} `json:"markets,omitempty"`
		Sid   *string `json:"sid"`
		Title *string `json:"title,omitempty"`
	} `json:"bookmakers,omitempty"`
	CommenceTime *theoddsapi.CommenceTime `json:"commence_time,omitempty"`
	HomeTeam     *theoddsapi.HomeTeam     `json:"home_team"`
	Id           *theoddsapi.MatchId      `json:"id,omitempty"` //nolint:revive,staticcheck // Field name matches generated API
	SportKey     *theoddsapi.SportKey     `json:"sport_key,omitempty"`
	SportTitle   *theoddsapi.SportTitle   `json:"sport_title,omitempty"`
},
) (*GameSpread, error) {
	if event.Bookmakers == nil || len(*event.Bookmakers) == 0 {
		return nil, fmt.Errorf("no bookmakers found for event")
	}

	// Use the first bookmaker that has spreads data
	for _, bookmaker := range *event.Bookmakers {
		if bookmaker.Markets == nil {
			continue
		}

		for _, market := range *bookmaker.Markets {
			if market.Key == nil || *market.Key != "spreads" {
				continue
			}

			if market.Outcomes == nil || len(*market.Outcomes) != 2 {
				continue
			}

			// Extract spread information from outcomes
			outcomes := *market.Outcomes
			var homeSpread, awaySpread float32
			var homeTeam, awayTeam string

			for _, outcome := range outcomes {
				if outcome.Point == nil {
					continue
				}

				// Determine which team this outcome belongs to
				if outcome.Name != nil {
					switch *outcome.Name {
					case *event.HomeTeam:
						homeSpread = *outcome.Point
						homeTeam = *event.HomeTeam
					case *event.AwayTeam:
						awaySpread = *outcome.Point
						awayTeam = *event.AwayTeam
					}
				}
			}

			// Determine favorite and underdog based on spreads
			// Positive spread indicates the favorite
			switch {
			case homeSpread > 0 && awaySpread < 0:
				// Home team is favorite
				return &GameSpread{
					HomeTeam: homeTeam,
					AwayTeam: awayTeam,
					Spread:   homeSpread,
					Favorite: "Home",
					Underdog: "Away",
				}, nil
			case awaySpread > 0 && homeSpread < 0:
				// Away team is favorite
				return &GameSpread{
					HomeTeam: homeTeam,
					AwayTeam: awayTeam,
					Spread:   awaySpread,
					Favorite: "Away",
					Underdog: "Home",
				}, nil
			case homeSpread == 0 && awaySpread == 0:
				// Even spread, default to home team as favorite
				return &GameSpread{
					HomeTeam: homeTeam,
					AwayTeam: awayTeam,
					Spread:   0,
					Favorite: "Home",
					Underdog: "Away",
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("no valid spread data found for event")
}

// UpdateGameSpreads updates the spreads for games in the database.
func (s *OddsService) UpdateGameSpreads(ctx context.Context, season, week int) error {
	slog.Info("Updating game spreads", "season", season, "week", week)

	// Fetch spreads from The Odds API
	spreads, err := s.FetchSpreadsForWeek(ctx, season, week)
	if err != nil {
		return fmt.Errorf("failed to fetch spreads: %w", err)
	}

	// Update games in database
	updatedCount := 0
	for _, spread := range spreads {
		// Find the game by teams and week/season
		var game database.Game
		err := s.db.GetDB().Where("season = ? AND week = ? AND favorite_team = ? AND underdog_team = ?",
			season, week, spread.HomeTeam, spread.AwayTeam).First(&game).Error
		if err != nil {
			slog.Warn("Game not found for spread update", "season", season, "week", week, "home", spread.HomeTeam, "away", spread.AwayTeam)
			continue
		}

		// Update game with spread information
		game.Spread = spread.Spread
		game.Favorite = &spread.Favorite
		game.Underdog = &spread.Underdog

		if err := s.db.GetDB().Save(&game).Error; err != nil {
			slog.Error("Failed to update game spread", "game_id", game.ID, "error", err)
			continue
		}

		updatedCount++
	}

	slog.Info("Updated game spreads", "updated_count", updatedCount, "total_spreads", len(spreads))
	return nil
}

// getWeekDateRange calculates the start and end dates for a given week and season.
func (s *OddsService) getWeekDateRange(_, week int) (time.Time, time.Time) {
	// This is a simplified implementation - you may need to adjust based on your week calculation logic
	// For now, we'll assume weeks start on Tuesday and end on Monday
	weekStart := s.config.ESPN.Week1Date.Add(time.Duration((week-1)*7) * 24 * time.Hour)
	weekEnd := weekStart.Add(7 * 24 * time.Hour)

	return weekStart, weekEnd
}
