// Package espnsync provides functionality to sync season schedule and game scores from the ESPN API.
package espnsync

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dhpollack/football-pool/internal/api-espn"
	"github.com/dhpollack/football-pool/internal/config"
	"github.com/dhpollack/football-pool/internal/database"
	"github.com/dhpollack/football-pool/internal/odds-sync"
)

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

// SyncService orchestrates the fetching, transformation, and storage of ESPN data.
type SyncService struct {
	db           *database.Database
	espnClient   *apiespn.ClientWithResponses
	oddsService  *oddssync.OddsService
	cache        *Cache
	transformer  *Transformer
	syncEnabled  bool
	config       *config.Config
	timeProvider TimeProvider
}

// NewSyncService creates a new SyncService instance.
func NewSyncService(db *database.Database, config *config.Config) (*SyncService, error) {
	return NewSyncServiceWithTimeProvider(db, config, RealTimeProvider{})
}

// NewSyncServiceWithTimeProvider creates a new SyncService instance with a custom time provider.
// This is primarily for testing purposes.
func NewSyncServiceWithTimeProvider(db *database.Database, config *config.Config, timeProvider TimeProvider) (*SyncService, error) {
	// Create ESPN client
	client, err := apiespn.NewClientWithResponses(config.ESPN.BaseURL)
	if err != nil {
		return nil, err
	}

	// Create odds service
	oddsService, err := oddssync.NewOddsServiceWithTimeProvider(db, config, timeProvider)
	if err != nil {
		return nil, err
	}

	// Create cache
	cache := NewCache(config.ESPN.CacheDir, config.ESPN.CacheExpiry)

	// Create transformer
	transformer := NewTransformer(db)

	return &SyncService{
		db:           db,
		espnClient:   client,
		oddsService:  oddsService,
		cache:        cache,
		transformer:  transformer,
		syncEnabled:  config.ESPN.SyncEnabled,
		config:       config,
		timeProvider: timeProvider,
	}, nil
}

// Start begins the background synchronization process.
func (s *SyncService) Start(ctx context.Context, interval time.Duration) {
	if !s.syncEnabled {
		slog.Info("ESPN sync service is disabled")
		return
	}

	slog.Info("Starting ESPN sync service", "interval", interval.String())

	// Run initial sync
	s.syncData(ctx)

	// Backfill missing weeks
	go s.BackfillWeeks(ctx)

	// Check and update spreads for current and past week
	go s.CheckAndUpdateSpreads(ctx)

	// Start weekly spread updates at Monday 11pm EST
	go s.StartWeeklySpreadUpdates(ctx)

	// Start periodic sync
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping ESPN sync service")
			return
		case <-ticker.C:
			s.syncData(ctx)
		}
	}
}

// BackfillWeeks checks for missing games in the database for all weeks of the season and syncs them if necessary.
// This is a one-time operation at startup.
func (s *SyncService) BackfillWeeks(ctx context.Context) {
	slog.Info("Starting backfill of missing weeks")
	season := s.config.ESPN.SeasonYear
	for week := 1; week <= 18; week++ {
		hasGames, err := s.db.WeekHasGames(season, week)
		if err != nil {
			slog.Error("Failed to check if week has games", "season", season, "week", week, "error", err)
			continue
		}

		if !hasGames {
			slog.Info("Backfilling data for week", "season", season, "week", week)
			if err := s.SyncWeekData(ctx, season, week); err != nil {
				slog.Error("Failed to backfill week data", "season", season, "week", week, "error", err)
			}
		}
	}
	slog.Info("Finished backfill of missing weeks")
}

// syncData performs a single synchronization cycle.
func (s *SyncService) syncData(ctx context.Context) {
	slog.Info("Starting ESPN data sync")

	// Get current season and week
	currentSeason, currentWeek := s.getCurrentSeasonAndWeek()

	// Sync data for current week
	if err := s.SyncWeekData(ctx, currentSeason, currentWeek); err != nil {
		slog.Error("Failed to sync week data", "season", currentSeason, "week", currentWeek, "error", err)
		return
	}

	slog.Info("ESPN data sync completed successfully")
}

// SyncWeekData syncs data for a specific week and season.
func (s *SyncService) SyncWeekData(ctx context.Context, season, week int) error {
	slog.Info("Syncing week data", "season", season, "week", week)

	// Fetch events from ESPN API
	events, err := s.fetchEvents(ctx, season, week)
	if err != nil {
		return err
	}

	// Transform and store events
	if err := s.transformAndStoreEvents(events, season, week); err != nil {
		return err
	}

	return nil
}

// fetchEvents retrieves events from ESPN API for a specific season and week.
func (s *SyncService) fetchEvents(ctx context.Context, season, week int) ([]apiespn.Event, error) {
	slog.Info("Fetching events from ESPN API", "season", season, "week", week)

	// Check cache first
	if cachedEvents, found := s.cache.Get(season, week); found {
		slog.Debug("Using cached events", "season", season, "week", week, "events", len(cachedEvents))
		return cachedEvents, nil
	}

	// Fetch from ESPN Site API
	// Note: seasontype 2 represents regular season in NFL
	params := &apiespn.GetScoreboardParams{
		Week:       &week,
		Seasontype: func() *int { i := 2; return &i }(), // seasontype 2 = regular season
	}

	response, err := s.espnClient.GetScoreboardWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("ESPN API returned status %d: %s", response.StatusCode(), response.Status())
	}

	if response.JSON200 == nil {
		slog.Error("ESPN API returned nil JSON200 response", "status", response.StatusCode(), "body", string(response.Body))
		return nil, fmt.Errorf("ESPN API returned empty response")
	}

	// Extract events directly from the scoreboard response
	var events []apiespn.Event
	if response.JSON200.Events != nil {
		events = *response.JSON200.Events
		slog.Debug("Extracted events from scoreboard", "event_count", len(events))

		// Debug: log event IDs to check for duplicates
		for i, event := range events {
			if event.Id != nil {
				slog.Debug("Event details", "index", i, "id", *event.Id, "name", event.Name)
			}
		}
	}

	slog.Info("Fetched events from ESPN API", "season", season, "week", week, "events", len(events))

	// Cache the results
	if err := s.cache.Set(season, week, events); err != nil {
		slog.Warn("Failed to cache events", "error", err)
	}

	return events, nil
}

// transformAndStoreEvents transforms ESPN events to database models and stores them.
func (s *SyncService) transformAndStoreEvents(events []apiespn.Event, season, week int) error {
	slog.Info("transformAndStoreEvents called", "count", len(events), "season", season, "week", week)

	// Deduplicate events by ID to handle duplicate events in the API response
	seenEvents := make(map[string]bool)
	seenGames := make(map[string]bool) // Track games by team combination to detect duplicates
	processedCount := 0

	for i, event := range events {
		if event.Id == nil {
			slog.Debug("Skipping event with no ID", "index", i)
			continue
		}

		// Skip duplicate events
		if seenEvents[*event.Id] {
			slog.Debug("Skipping duplicate event", "event_id", *event.Id, "index", i)
			continue
		}
		seenEvents[*event.Id] = true

		slog.Debug("Processing unique event", "index", i, "event_id", *event.Id, "name", event.Name, "total_seen", len(seenEvents))

		// Transform event to game and result
		game, result, err := s.transformer.TransformEvent(event, season, week)
		if err != nil {
			slog.Warn("Failed to transform event", "event_id", *event.Id, "error", err)
			continue
		}

		if game == nil {
			slog.Debug("Skipping event with no valid game data", "event_id", *event.Id)
			continue
		}

		// Check for duplicate games by team combination
		gameKey := fmt.Sprintf("%d-%d-%s-%s", season, week, game.HomeTeam, game.AwayTeam)
		if seenGames[gameKey] {
			slog.Warn("Skipping duplicate game", "season", season, "week", week, "home", game.HomeTeam, "away", game.AwayTeam)
			continue
		}
		seenGames[gameKey] = true

		// Store game and result in database
		slog.Debug("About to store game", "season", season, "week", week, "home", game.HomeTeam, "away", game.AwayTeam, "has_result", result != nil)
		if err := s.transformer.StoreGameAndResult(game, result); err != nil {
			slog.Error("Failed to store game and result", "game", game, "error", err)
			continue
		}

		slog.Debug("Successfully stored game", "season", season, "week", week, "home", game.HomeTeam, "away", game.AwayTeam)
		processedCount++
	}

	slog.Info("Completed transforming and storing events", "total_events", len(events), "unique_events", len(seenEvents), "processed", processedCount)
	return nil
}

// getCurrentSeasonAndWeek returns the current NFL season and week.
// This uses the configured season year and calculates the week based on the
// configured Week 1 date, accounting for weeks ending on Monday.
func (s *SyncService) getCurrentSeasonAndWeek() (int, int) {
	currentTime := s.timeProvider.Now()

	// Use configured season year
	season := s.config.ESPN.SeasonYear

	// Calculate week based on days since Week 1 date
	// NFL Week 1 starts with the first game (typically Thursday)
	// Weeks run from the first game through the following Monday
	daysSinceWeek1 := int(currentTime.Sub(s.config.ESPN.Week1Date).Hours() / 24)

	// If we're before Week 1 first game, we're in Week 0 (preseason)
	if daysSinceWeek1 < 0 {
		return season, 0
	}

	// Simple calculation: week = (days since first game) / 7 + 1
	// This works because each week is exactly 7 days from the first game to the following Monday
	week := (daysSinceWeek1 / 7) + 1

	// Ensure week is within reasonable bounds
	if week < 1 {
		week = 1
	}
	if week > 18 { // NFL regular season is 18 weeks
		week = 18
	}

	return season, week
}

// SyncNow performs an immediate synchronization.
func (s *SyncService) SyncNow(ctx context.Context) error {
	s.syncData(ctx)
	return nil
}

// CheckAndUpdateSpreads checks games for the current and past week and updates spreads if needed.
func (s *SyncService) CheckAndUpdateSpreads(ctx context.Context) {
	slog.Info("Checking and updating spreads for current and past week")

	currentSeason, currentWeek := s.getCurrentSeasonAndWeek()

	// Check current week and previous week
	for week := currentWeek - 1; week <= currentWeek; week++ {
		if week < 1 {
			continue
		}

		// Check if all spreads are 0 for this week
		allSpreadsZero, err := s.allSpreadsZero(currentSeason, week)
		if err != nil {
			slog.Error("Failed to check spreads", "season", currentSeason, "week", week, "error", err)
			continue
		}

		if allSpreadsZero {
			slog.Info("Updating spreads for week", "season", currentSeason, "week", week)
			if err := s.oddsService.UpdateGameSpreads(ctx, currentSeason, week); err != nil {
				slog.Error("Failed to update spreads", "season", currentSeason, "week", week, "error", err)
			}
		} else {
			slog.Info("Spreads already exist for week", "season", currentSeason, "week", week)
		}
	}
}

// allSpreadsZero checks if all games in a week have spread = 0.
func (s *SyncService) allSpreadsZero(season, week int) (bool, error) {
	var games []database.Game
	if err := s.db.GetDB().Where("season = ? AND week = ?", season, week).Find(&games).Error; err != nil {
		return false, err
	}

	if len(games) == 0 {
		return false, nil // No games to check
	}

	for _, game := range games {
		if game.Spread != 0 {
			return false, nil // Found a game with non-zero spread
		}
	}

	return true, nil // All spreads are 0
}

// StartWeeklySpreadUpdates starts a background process that updates spreads every Monday at 11pm EST.
func (s *SyncService) StartWeeklySpreadUpdates(ctx context.Context) {
	slog.Info("Starting weekly spread updates (Monday 11pm EST)")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping weekly spread updates")
			return
		default:
			// Calculate next Monday 11pm EST
			nextUpdate := s.nextMonday11pmEST()
			durationUntilNextUpdate := time.Until(nextUpdate)

			slog.Info("Next spread update scheduled", "time", nextUpdate, "duration", durationUntilNextUpdate)

			// Wait until next Monday 11pm EST
			timer := time.NewTimer(durationUntilNextUpdate)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				// Update spreads for upcoming week
				currentSeason, currentWeek := s.getCurrentSeasonAndWeek()
				upcomingWeek := currentWeek + 1

				// Only update if upcoming week is within the season (1-18)
				if upcomingWeek >= 1 && upcomingWeek <= 18 {
					slog.Info("Updating spreads for upcoming week", "season", currentSeason, "week", upcomingWeek)
					if err := s.oddsService.UpdateGameSpreads(ctx, currentSeason, upcomingWeek); err != nil {
						slog.Error("Failed to update spreads for upcoming week", "season", currentSeason, "week", upcomingWeek, "error", err)
					}
				} else {
					slog.Info("Upcoming week is outside season range, skipping spread update", "week", upcomingWeek)
				}
			}
		}
	}
}

// nextMonday11pmEST calculates the next Monday at 11pm EST.
func (s *SyncService) nextMonday11pmEST() time.Time {
	now := s.timeProvider.Now()

	// Convert to EST (UTC-5)
	estLocation, err := time.LoadLocation("America/New_York")
	if err != nil {
		slog.Warn("Failed to load EST timezone, using UTC", "error", err)
		estLocation = time.UTC
	}

	nowEST := now.In(estLocation)

	// Find next Monday
	daysUntilMonday := (1 - int(nowEST.Weekday()) + 7) % 7
	if daysUntilMonday == 0 && nowEST.Hour() >= 23 {
		// If it's already Monday after 11pm, schedule for next Monday
		daysUntilMonday = 7
	}

	nextMonday := nowEST.AddDate(0, 0, daysUntilMonday)
	nextMonday11pm := time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 23, 0, 0, 0, estLocation)

	return nextMonday11pm
}

// GetSyncStatus returns the current status of the sync service.
func (s *SyncService) GetSyncStatus() map[string]interface{} {
	return map[string]interface{}{
		"enabled":   s.syncEnabled,
		"last_sync": time.Now(), // TODO: Track actual last sync time
	}
}
