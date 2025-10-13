package api

import (
	"testing"
	"time"

	"github.com/dhpollack/football-pool/internal/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	HomeConst string = "Home"
	AwayConst string = "Away"
)

func TestUserToResponse(t *testing.T) {
	now := time.Now()
	user := database.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:  "Test User",
		Email: "test@test.com",
		Role:  "user",
		Player: database.Player{
			Model: gorm.Model{
				ID: 1,
			},
			UserID:  1,
			Name:    "Player Name",
			Address: "Player Address",
		},
	}

	response := UserToResponse(user)

	assert.Equal(t, user.ID, response.Id)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.Role, response.Role)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
	assert.Equal(t, user.UpdatedAt, response.UpdatedAt)
	assert.NotNil(t, response.Player)
	assert.Equal(t, user.Player.ID, response.Player.Id)
	assert.Equal(t, user.Player.UserID, response.Player.UserId)
	assert.Equal(t, user.Player.Name, response.Player.Name)
	assert.Equal(t, user.Player.Address, response.Player.Address)
}

func TestUserWithStatsFromUser(t *testing.T) {
	now := time.Now()
	user := database.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:  "Test User",
		Email: "test@test.com",
		Role:  "user",
	}

	response := UserWithStatsFromUser(user, 10, 5)

	assert.Equal(t, user.ID, response.Id)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.Role, response.Role)
	assert.Equal(t, 10, response.PickCount)
	assert.Equal(t, 5, response.TotalWins)
}

func TestUserFromRequest(t *testing.T) {
	password := "password"
	req := UserRequest{
		Name:     "Test User",
		Email:    "test@test.com",
		Role:     "user",
		Password: &password,
	}

	user, err := UserFromRequest(req)

	assert.NoError(t, err)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Role, user.Role)
	assert.Equal(t, *req.Password, user.Password)
}

func TestGameToResponse(t *testing.T) {
	now := time.Now()
	home := "Home"
	away := "Away"
	game := database.Game{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Week:      1,
		Season:    2023,
		HomeTeam:  "Lions",
		AwayTeam:  "Chiefs",
		Favorite:  &home,
		Underdog:  &away,
		Spread:    3.5,
		StartTime: now,
	}

	response := GameToResponse(game)
	expectedFav := ConvertStringPointerToTeamDesignationPointer(game.Favorite)
	expectedUnd := ConvertStringPointerToTeamDesignationPointer(game.Underdog)

	assert.Equal(t, game.ID, response.Id)
	assert.Equal(t, game.Week, response.Week)
	assert.Equal(t, game.Season, response.Season)
	assert.Equal(t, game.HomeTeam, response.HomeTeam)
	assert.Equal(t, game.AwayTeam, response.AwayTeam)
	assert.Equal(t, expectedFav, response.Favorite)
	assert.Equal(t, expectedUnd, response.Underdog)
	assert.Equal(t, game.Spread, response.Spread)
	assert.Equal(t, game.StartTime, response.StartTime)
}

func TestGameFromRequest(t *testing.T) {
	now := time.Now()
	favoriteHome := Home
	underdogAway := Away
	req := GameRequest{
		Week:      1,
		Season:    2023,
		HomeTeam:  "Lions",
		AwayTeam:  "Chiefs",
		Favorite:  &favoriteHome,
		Underdog:  &underdogAway,
		Spread:    3.5,
		StartTime: now,
	}

	game, err := GameFromRequest(req)
	expectedFav := ConvertTeamDesignationPointerToStringPointer(req.Favorite)
	expectedUnd := ConvertTeamDesignationPointerToStringPointer(req.Underdog)

	assert.NoError(t, err)
	assert.Equal(t, req.Week, game.Week)
	assert.Equal(t, req.Season, game.Season)
	assert.Equal(t, req.HomeTeam, game.HomeTeam)
	assert.Equal(t, req.AwayTeam, game.AwayTeam)
	assert.Equal(t, expectedFav, game.Favorite)
	assert.Equal(t, expectedUnd, game.Underdog)
	assert.Equal(t, req.Spread, game.Spread)
	assert.Equal(t, req.StartTime, game.StartTime)
}

func TestPickToResponse(t *testing.T) {
	now := time.Now()
	pick := database.Pick{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:    1,
		GameID:    1,
		Picked:    "favorite",
		Rank:      1,
		QuickPick: true,
	}

	response := PickToResponse(pick)

	assert.Equal(t, pick.ID, response.Id)
	assert.Equal(t, pick.UserID, response.UserId)
	assert.Equal(t, pick.GameID, response.GameId)
	assert.Equal(t, pick.Picked, response.Picked)
	assert.Equal(t, pick.Rank, response.Rank)
	assert.Equal(t, pick.QuickPick, response.QuickPick)
}

func TestPickFromRequest(t *testing.T) {
	userID := uint(1)
	req := PickRequest{
		GameId:    1,
		Picked:    "favorite",
		Rank:      1,
		QuickPick: true,
		UserId:    &userID,
	}

	pick, err := PickFromRequest(req)

	assert.NoError(t, err)
	assert.Equal(t, req.GameId, pick.GameID)
	assert.Equal(t, req.Picked, pick.Picked)
	assert.Equal(t, req.Rank, pick.Rank)
	assert.Equal(t, req.QuickPick, pick.QuickPick)
	assert.Equal(t, *req.UserId, pick.UserID)
}

func TestPicksFromRequest(t *testing.T) {
	userID := uint(1)
	reqs := []PickRequest{
		{
			GameId:    1,
			Picked:    "favorite",
			Rank:      1,
			QuickPick: true,
			UserId:    &userID,
		},
	}

	picks, err := PicksFromRequest(reqs)

	assert.NoError(t, err)
	assert.Len(t, picks, 1)
	assert.Equal(t, reqs[0].GameId, picks[0].GameID)
}

func TestResultToResponse(t *testing.T) {
	now := time.Now()
	result := database.Result{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		GameID:        1,
		FavoriteScore: 21,
		UnderdogScore: 17,
		Outcome:       "favorite",
	}

	response := ResultToResponse(result)

	assert.Equal(t, result.ID, response.Id)
	assert.Equal(t, result.GameID, response.GameId)
	assert.Equal(t, result.FavoriteScore, response.FavoriteScore)
	assert.Equal(t, result.UnderdogScore, response.UnderdogScore)
	assert.Equal(t, result.Outcome, response.Outcome)
}

func TestResultFromRequest(t *testing.T) {
	req := ResultRequest{
		GameId:        1,
		FavoriteScore: 21,
		UnderdogScore: 17,
		Outcome:       "favorite",
	}

	result := ResultFromRequest(req)

	assert.Equal(t, req.GameId, result.GameID)
	assert.Equal(t, req.FavoriteScore, result.FavoriteScore)
	assert.Equal(t, req.UnderdogScore, result.UnderdogScore)
	assert.Equal(t, req.Outcome, result.Outcome)
}

func TestSurvivorPickToResponse(t *testing.T) {
	now := time.Now()
	pick := database.SurvivorPick{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID: 1,
		Week:   1,
		Team:   "Lions",
	}

	response := SurvivorPickToResponse(pick)

	assert.Equal(t, pick.ID, response.Id)
	assert.Equal(t, pick.UserID, response.UserId)
	assert.Equal(t, pick.Week, response.Week)
	assert.Equal(t, pick.Team, response.Team)
}

func TestSurvivorPickFromRequest(t *testing.T) {
	userID := uint(1)
	req := SurvivorPickRequest{
		Week:   1,
		Team:   "Lions",
		UserId: &userID,
	}

	pick := SurvivorPickFromRequest(req)

	assert.Equal(t, req.Week, pick.Week)
	assert.Equal(t, req.Team, pick.Team)
	assert.Equal(t, *req.UserId, pick.UserID)
}

func TestPlayerToResponse(t *testing.T) {
	player := database.Player{
		Model: gorm.Model{
			ID: 1,
		},
		UserID:  1,
		Name:    "Player Name",
		Address: "Player Address",
	}

	response := PlayerToResponse(player)

	assert.Equal(t, player.ID, response.Id)
	assert.Equal(t, player.UserID, response.UserId)
	assert.Equal(t, player.Name, response.Name)
	assert.Equal(t, player.Address, response.Address)
}

func TestPlayerFromRequest(t *testing.T) {
	req := PlayerRequest{
		Name:    "Player Name",
		Address: "Player Address",
	}

	player := PlayerFromRequest(req, 1)

	assert.Equal(t, uint(1), player.UserID)
	assert.Equal(t, req.Name, player.Name)
	assert.Equal(t, req.Address, player.Address)
}

func TestUserFromRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     UserRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid request with password",
			request: UserRequest{
				Name:     "Test User",
				Email:    "test@test.com",
				Role:     "user",
				Password: stringPtr("password"),
			},
			expectError: false,
		},
		{
			name: "Valid request without password",
			request: UserRequest{
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "user",
			},
			expectError: true,
			errorMsg:    "Password",
		},
		{
			name: "Missing required name",
			request: UserRequest{
				Email: "test@test.com",
				Role:  "user",
			},
			expectError: true,
			errorMsg:    "Name",
		},
		{
			name: "Invalid email format",
			request: UserRequest{
				Name:  "Test User",
				Email: "invalid-email",
				Role:  "user",
			},
			expectError: true,
			errorMsg:    "Email",
		},
		{
			name: "Missing role",
			request: UserRequest{
				Name:  "Test User",
				Email: "test@test.com",
			},
			expectError: true,
			errorMsg:    "Role",
		},
		{
			name: "Empty name",
			request: UserRequest{
				Name:  "",
				Email: "test@test.com",
				Role:  "user",
			},
			expectError: true,
			errorMsg:    "Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := UserFromRequest(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.request.Name, user.Name)
				assert.Equal(t, tt.request.Email, user.Email)
				assert.Equal(t, tt.request.Role, user.Role)
				if tt.request.Password != nil {
					assert.Equal(t, *tt.request.Password, user.Password)
				}
			}
		})
	}
}

func TestGameFromRequestValidation(t *testing.T) {
	favoriteHome := Home
	underdogAway := Away

	tests := []struct {
		name        string
		request     GameRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid request",
			request: GameRequest{
				Week:      1,
				Season:    2023,
				HomeTeam:  "Lions",
				AwayTeam:  "Chiefs",
				Favorite:  &favoriteHome,
				Underdog:  &underdogAway,
				Spread:    3.5,
				StartTime: time.Now(),
			},
			expectError: false,
		},
		{
			name: "Zero week",
			request: GameRequest{
				Week:      0,
				Season:    2023,
				HomeTeam:  "Lions",
				AwayTeam:  "Chiefs",
				Favorite:  &favoriteHome,
				Underdog:  &underdogAway,
				Spread:    3.5,
				StartTime: time.Now(),
			},
			expectError: true,
			errorMsg:    "Week",
		},
		{
			name: "Zero season",
			request: GameRequest{
				Week:      1,
				Season:    0,
				HomeTeam:  "Lions",
				AwayTeam:  "Chiefs",
				Favorite:  &favoriteHome,
				Underdog:  &underdogAway,
				Spread:    3.5,
				StartTime: time.Now(),
			},
			expectError: true,
			errorMsg:    "Season",
		},
		{
			name: "Missing home team",
			request: GameRequest{
				Week:      1,
				Season:    2023,
				AwayTeam:  "Chiefs",
				Favorite:  &favoriteHome,
				Underdog:  &underdogAway,
				Spread:    3.5,
				StartTime: time.Now(),
			},
			expectError: true,
			errorMsg:    "HomeTeam",
		},
		{
			name: "Missing away team",
			request: GameRequest{
				Week:      1,
				Season:    2023,
				HomeTeam:  "Lions",
				Favorite:  &favoriteHome,
				Underdog:  &underdogAway,
				Spread:    3.5,
				StartTime: time.Now(),
			},
			expectError: true,
			errorMsg:    "AwayTeam",
		},
		{
			name: "Missing start time",
			request: GameRequest{
				Week:     1,
				Season:   2023,
				HomeTeam: "Lions",
				AwayTeam: "Chiefs",
				Favorite: &favoriteHome,
				Underdog: &underdogAway,
			},
			expectError: true,
			errorMsg:    "StartTime",
		},
		{
			name: "Negative spread",
			request: GameRequest{
				Week:      1,
				Season:    2023,
				HomeTeam:  "Lions",
				AwayTeam:  "Chiefs",
				Favorite:  &favoriteHome,
				Underdog:  &underdogAway,
				Spread:    -3.5,
				StartTime: time.Now(),
			},
			expectError: true,
			errorMsg:    "Spread",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game, err := GameFromRequest(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.request.Week, game.Week)
				assert.Equal(t, tt.request.Season, game.Season)
				assert.Equal(t, tt.request.HomeTeam, game.HomeTeam)
				assert.Equal(t, tt.request.AwayTeam, game.AwayTeam)
				assert.Equal(t, tt.request.Spread, game.Spread)
				assert.Equal(t, tt.request.StartTime, game.StartTime)
			}
		})
	}
}

func TestPickFromRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     PickRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid request with user ID",
			request: PickRequest{
				GameId:    1,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
				UserId:    uintPtr(1),
			},
			expectError: false,
		},
		{
			name: "Valid request without user ID",
			request: PickRequest{
				GameId:    1,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
			},
			expectError: false,
		},
		{
			name: "Missing game ID",
			request: PickRequest{
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
			},
			expectError: true,
			errorMsg:    "GameID",
		},
		{
			name: "Missing picked value",
			request: PickRequest{
				GameId:    1,
				Rank:      1,
				QuickPick: true,
			},
			expectError: true,
			errorMsg:    "Picked",
		},
		{
			name: "Missing rank",
			request: PickRequest{
				GameId:    1,
				Picked:    "favorite",
				QuickPick: true,
			},
			expectError: true,
			errorMsg:    "Rank",
		},
		{
			name: "Zero game ID",
			request: PickRequest{
				GameId:    0,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
			},
			expectError: true,
			errorMsg:    "GameID",
		},
		{
			name: "Zero user ID",
			request: PickRequest{
				GameId:    1,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
				UserId:    uintPtr(0),
			},
			expectError: true,
			errorMsg:    "UserID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pick, err := PickFromRequest(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.request.GameId, pick.GameID)
				assert.Equal(t, tt.request.Picked, pick.Picked)
				assert.Equal(t, tt.request.Rank, pick.Rank)
				assert.Equal(t, tt.request.QuickPick, pick.QuickPick)
				if tt.request.UserId != nil {
					assert.Equal(t, *tt.request.UserId, pick.UserID)
				}
			}
		})
	}
}

func TestPicksFromRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		requests    []PickRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid multiple picks",
			requests: []PickRequest{
				{
					GameId:    1,
					Picked:    "favorite",
					Rank:      1,
					QuickPick: true,
				},
				{
					GameId:    2,
					Picked:    "underdog",
					Rank:      2,
					QuickPick: false,
				},
			},
			expectError: false,
		},
		{
			name:        "Empty picks array",
			requests:    []PickRequest{},
			expectError: false,
		},
		{
			name: "One invalid pick in array",
			requests: []PickRequest{
				{
					GameId:    1,
					Picked:    "favorite",
					Rank:      1,
					QuickPick: true,
				},
				{
					Picked:    "favorite",
					Rank:      2,
					QuickPick: true,
				},
			},
			expectError: true,
			errorMsg:    "GameID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			picks, err := PicksFromRequest(tt.requests)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, picks, len(tt.requests))
			}
		})
	}
}

func TestUserToResponseEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		user     database.User
		expect   UserResponse
		checkNil bool
	}{
		{
			name: "User without player",
			user: database.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "user",
			},
			expect: UserResponse{
				Id:    1,
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "user",
			},
			checkNil: true,
		},
		{
			name: "User with empty player",
			user: database.User{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "user",
				Player: database.Player{
					Model: gorm.Model{
						ID: 0,
					},
				},
			},
			expect: UserResponse{
				Id:    1,
				Name:  "Test User",
				Email: "test@test.com",
				Role:  "user",
			},
			checkNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := UserToResponse(tt.user)

			assert.Equal(t, tt.expect.Id, response.Id)
			assert.Equal(t, tt.expect.Name, response.Name)
			assert.Equal(t, tt.expect.Email, response.Email)
			assert.Equal(t, tt.expect.Role, response.Role)

			if tt.checkNil {
				assert.Nil(t, response.Player)
			}
		})
	}
}

func TestPickToResponseEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		pick    database.Pick
		expect  PickResponse
		userNil bool
		gameNil bool
	}{
		{
			name: "Pick without preloaded user and game",
			pick: database.Pick{
				Model: gorm.Model{
					ID: 1,
				},
				UserID:    1,
				GameID:    1,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
			},
			expect: PickResponse{
				Id:        1,
				UserId:    1,
				GameId:    1,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
			},
			userNil: true,
			gameNil: true,
		},
		{
			name: "Pick with empty user and game",
			pick: database.Pick{
				Model: gorm.Model{
					ID: 1,
				},
				UserID:    1,
				GameID:    1,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
				User: database.User{
					Model: gorm.Model{
						ID: 0,
					},
				},
				Game: database.Game{
					Model: gorm.Model{
						ID: 0,
					},
				},
			},
			expect: PickResponse{
				Id:        1,
				UserId:    1,
				GameId:    1,
				Picked:    "favorite",
				Rank:      1,
				QuickPick: true,
			},
			userNil: true,
			gameNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := PickToResponse(tt.pick)

			assert.Equal(t, tt.expect.Id, response.Id)
			assert.Equal(t, tt.expect.UserId, response.UserId)
			assert.Equal(t, tt.expect.GameId, response.GameId)
			assert.Equal(t, tt.expect.Picked, response.Picked)
			assert.Equal(t, tt.expect.Rank, response.Rank)
			assert.Equal(t, tt.expect.QuickPick, response.QuickPick)

			if tt.userNil {
				assert.Nil(t, response.User)
			}
			if tt.gameNil {
				assert.Nil(t, response.Game)
			}
		})
	}
}

// Helper functions.
func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}
