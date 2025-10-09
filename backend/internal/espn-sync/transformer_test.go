package espnsync

import (
	"testing"
	"time"

	apiespn "github.com/dhpollack/football-pool/internal/api-espn"
	"github.com/dhpollack/football-pool/internal/database"
)

const favoriteRes = "favorite"

// espnDateTime creates an ESPNDateTime pointer from a time.Time.
func espnDateTime(t time.Time) *apiespn.ESPNDateTime {
	return &apiespn.ESPNDateTime{Time: t}
}

func TestTransformer_TransformEvent(t *testing.T) {
	db, err := database.New("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	transformer := NewTransformer(db)

	tests := []struct {
		name     string
		event    apiespn.Event
		season   int
		week     int
		wantGame *database.Game
		wantErr  bool
	}{
		{
			name: "valid event with two teams",
			event: apiespn.Event{
				Id:   &[]string{"event1"}[0],
				Name: &[]string{"Game 1"}[0],
				Date: espnDateTime(time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC)),
				Competitions: &[]apiespn.Competition{
					{
						Date: espnDateTime(time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC)),
						Competitors: &[]apiespn.Competitor{
							{
								HomeAway: &[]string{"home"}[0],
								Team: &apiespn.Team{
									DisplayName: &[]string{"Chiefs"}[0],
								},
							},
							{
								HomeAway: &[]string{"away"}[0],
								Team: &apiespn.Team{
									DisplayName: &[]string{"Lions"}[0],
								},
							},
						},
					},
				},
			},
			season: 2023,
			week:   1,
			wantGame: &database.Game{
				Week:      1,
				Season:    2023,
				HomeTeam:  "Chiefs",
				AwayTeam:  "Lions",
				Spread:    0.0,
				StartTime: time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "event with no competitions",
			event: apiespn.Event{
				Id:   &[]string{"event2"}[0],
				Name: &[]string{"Game 2"}[0],
			},
			season:   2023,
			week:     1,
			wantGame: nil,
			wantErr:  false,
		},
		{
			name: "event with empty competitions",
			event: apiespn.Event{
				Id:           &[]string{"event3"}[0],
				Name:         &[]string{"Game 3"}[0],
				Competitions: &[]apiespn.Competition{},
			},
			season:   2023,
			week:     1,
			wantGame: nil,
			wantErr:  false,
		},
		{
			name: "event with missing team names",
			event: apiespn.Event{
				Id:   &[]string{"event4"}[0],
				Name: &[]string{"Game 4"}[0],
				Competitions: &[]apiespn.Competition{
					{
						Competitors: &[]apiespn.Competitor{
							{
								HomeAway: &[]string{"home"}[0],
								Team:     &apiespn.Team{},
							},
							{
								HomeAway: &[]string{"away"}[0],
								Team:     &apiespn.Team{},
							},
						},
					},
				},
			},
			season:   2023,
			week:     1,
			wantGame: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, result, err := transformer.TransformEvent(tt.event, tt.season, tt.week)

			if (err != nil) != tt.wantErr {
				t.Errorf("TransformEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantGame == nil {
				if game != nil {
					t.Errorf("TransformEvent() game = %v, want nil", game)
				}
				return
			}

			if game == nil {
				t.Errorf("TransformEvent() game = nil, want %v", tt.wantGame)
				return
			}

			if game.Week != tt.wantGame.Week {
				t.Errorf("TransformEvent() week = %v, want %v", game.Week, tt.wantGame.Week)
			}
			if game.Season != tt.wantGame.Season {
				t.Errorf("TransformEvent() season = %v, want %v", game.Season, tt.wantGame.Season)
			}
			if game.HomeTeam != tt.wantGame.HomeTeam {
				t.Errorf("TransformEvent() homeTeam = %v, want %v", game.HomeTeam, tt.wantGame.HomeTeam)
			}
			if game.AwayTeam != tt.wantGame.AwayTeam {
				t.Errorf("TransformEvent() awayTeam = %v, want %v", game.AwayTeam, tt.wantGame.AwayTeam)
			}
			if game.Favorite != tt.wantGame.Favorite {
				t.Errorf("TransformEvent() favorite = %v, want %v", game.Favorite, tt.wantGame.Favorite)
			}
			if game.Underdog != tt.wantGame.Underdog {
				t.Errorf("TransformEvent() underdog = %v, want %v", game.Underdog, tt.wantGame.Underdog)
			}
			if game.Spread != tt.wantGame.Spread {
				t.Errorf("TransformEvent() spread = %v, want %v", game.Spread, tt.wantGame.Spread)
			}
			// For the "event with missing team names" test, don't compare exact start time
			// since it uses time.Now() as fallback
			if tt.name != "event with missing team names" && !game.StartTime.Equal(tt.wantGame.StartTime) {
				t.Errorf("TransformEvent() startTime = %v, want %v", game.StartTime, tt.wantGame.StartTime)
			}

			// Result should be nil since no scores are provided
			if result != nil {
				t.Errorf("TransformEvent() result = %v, want nil", result)
			}
		})
	}
}

func TestTransformer_TransformEventWithScores(t *testing.T) {
	db, err := database.New("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	transformer := NewTransformer(db)

	event := apiespn.Event{
		Id:   &[]string{"event1"}[0],
		Name: &[]string{"Game 1"}[0],
		Date: espnDateTime(time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC)),
		Competitions: &[]apiespn.Competition{
			{
				Date: espnDateTime(time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC)),
				Competitors: &[]apiespn.Competitor{
					{
						HomeAway: &[]string{"home"}[0],
						Team: &apiespn.Team{
							DisplayName: &[]string{"Chiefs"}[0],
						},
						Score: &[]string{"21"}[0],
					},
					{
						HomeAway: &[]string{"away"}[0],
						Team: &apiespn.Team{
							DisplayName: &[]string{"Lions"}[0],
						},
						Score: &[]string{"20"}[0],
					},
				},
			},
		},
	}

	game, result, err := transformer.TransformEvent(event, 2023, 1)
	if err != nil {
		t.Fatalf("TransformEvent() error = %v", err)
	}

	if game == nil {
		t.Fatal("TransformEvent() game = nil, want valid game")
	}

	if result == nil {
		t.Fatal("TransformEvent() result = nil, want valid result")
	}

	if result.FavoriteScore != 21 {
		t.Errorf("TransformEvent() favoriteScore = %v, want 21", result.FavoriteScore)
	}
	if result.UnderdogScore != 20 {
		t.Errorf("TransformEvent() underdogScore = %v, want 20", result.UnderdogScore)
	}
	if result.Outcome != favoriteRes {
		t.Errorf("TransformEvent() outcome = %v, want 'favorite'", result.Outcome)
	}
}

func TestTransformer_StoreGameAndResult(t *testing.T) {
	db, err := database.New("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	transformer := NewTransformer(db)

	game := &database.Game{
		Week:      1,
		Season:    2023,
		HomeTeam:  "Chiefs",
		AwayTeam:  "Lions",
		Spread:    0.0,
		StartTime: time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC),
	}

	// Test storing new game
	err = transformer.StoreGameAndResult(game, nil)
	if err != nil {
		t.Fatalf("StoreGameAndResult() error = %v", err)
	}

	// Verify game was stored
	var storedGame database.Game
	if err := db.GetDB().Where("season = ? AND week = ? AND favorite_team = ? AND underdog_team = ?",
		2023, 1, "Chiefs", "Lions").First(&storedGame).Error; err != nil {
		t.Fatalf("Failed to find stored game: %v", err)
	}

	if storedGame.ID == 0 {
		t.Error("Stored game should have ID")
	}

	// Test storing game with result
	result := &database.Result{
		Game:          *game,
		FavoriteScore: 21,
		UnderdogScore: 20,
		Outcome:       favoriteRes,
	}

	err = transformer.StoreGameAndResult(game, result)
	if err != nil {
		t.Fatalf("StoreGameAndResult() error = %v", err)
	}

	// Verify result was stored
	var storedResult database.Result
	if err := db.GetDB().Where("game_id = ?", storedGame.ID).First(&storedResult).Error; err != nil {
		t.Fatalf("Failed to find stored result: %v", err)
	}

	if storedResult.FavoriteScore != 21 {
		t.Errorf("Stored result favoriteScore = %v, want 21", storedResult.FavoriteScore)
	}
	if storedResult.UnderdogScore != 20 {
		t.Errorf("Stored result underdogScore = %v, want 20", storedResult.UnderdogScore)
	}
}

func TestTransformer_StoreGameAndResultUpdate(t *testing.T) {
	db, err := database.New("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	transformer := NewTransformer(db)

	game := &database.Game{
		Week:      1,
		Season:    2023,
		HomeTeam:  "Chiefs",
		AwayTeam:  "Lions",
		Spread:    0.0,
		StartTime: time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC),
	}

	// Store game initially
	err = transformer.StoreGameAndResult(game, nil)
	if err != nil {
		t.Fatalf("StoreGameAndResult() error = %v", err)
	}

	// Update game with result
	game.Spread = 3.5 // Update spread
	result := &database.Result{
		Game:          *game,
		FavoriteScore: 21,
		UnderdogScore: 20,
		Outcome:       favoriteRes,
	}

	err = transformer.StoreGameAndResult(game, result)
	if err != nil {
		t.Fatalf("StoreGameAndResult() error = %v", err)
	}

	// Verify game was updated
	var updatedGame database.Game
	if err := db.GetDB().Where("season = ? AND week = ? AND favorite_team = ? AND underdog_team = ?",
		2023, 1, "Chiefs", "Lions").First(&updatedGame).Error; err != nil {
		t.Fatalf("Failed to find updated game: %v", err)
	}

	if updatedGame.Spread != 3.5 {
		t.Errorf("Updated game spread = %v, want 3.5", updatedGame.Spread)
	}

	// Verify result was stored
	var storedResult database.Result
	if err := db.GetDB().Where("game_id = ?", updatedGame.ID).First(&storedResult).Error; err != nil {
		t.Fatalf("Failed to find stored result: %v", err)
	}

	if storedResult.FavoriteScore != 21 {
		t.Errorf("Stored result favoriteScore = %v, want 21", storedResult.FavoriteScore)
	}
}
