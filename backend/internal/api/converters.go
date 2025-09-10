package api

import (
	"github.com/david/football-pool/internal/database"
)

// UserToResponse converts a database User to a UserResponse
func UserToResponse(user database.User) UserResponse {
	response := UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Include player if it exists
	if user.Player.ID != 0 {
		response.Player = &PlayerResponse{
			Id:      user.Player.ID,
			UserId:  user.Player.UserID,
			Name:    user.Player.Name,
			Address: user.Player.Address,
		}
	}

	return response
}

// UserWithStatsFromUser converts a database User to a UserWithStats
func UserWithStatsFromUser(user database.User, pickCount, totalWins int) UserWithStats {
	stats := UserWithStats{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		PickCount: pickCount,
		TotalWins: totalWins,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Include player if it exists
	if user.Player.ID != 0 {
		stats.Player = &PlayerResponse{
			Id:      user.Player.ID,
			UserId:  user.Player.UserID,
			Name:    user.Player.Name,
			Address: user.Player.Address,
		}
	}

	return stats
}

// UserFromRequest converts a UserRequest to a database User
func UserFromRequest(req UserRequest) database.User {
	user := database.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	// Handle optional password
	if req.Password != nil {
		user.Password = *req.Password
	}

	return user
}

// GameToResponse converts a database Game to a GameResponse
func GameToResponse(game database.Game) GameResponse {
	return GameResponse{
		Id:           game.ID,
		Week:         game.Week,
		Season:       game.Season,
		FavoriteTeam: game.FavoriteTeam,
		UnderdogTeam: game.UnderdogTeam,
		Spread:       game.Spread,
		StartTime:    game.StartTime,
		CreatedAt:    game.CreatedAt,
		UpdatedAt:    game.UpdatedAt,
	}
}

// GameFromRequest converts a GameRequest to a database Game
func GameFromRequest(req GameRequest) database.Game {
	return database.Game{
		Week:         req.Week,
		Season:       req.Season,
		FavoriteTeam: req.FavoriteTeam,
		UnderdogTeam: req.UnderdogTeam,
		Spread:       req.Spread,
		StartTime:    req.StartTime,
	}
}

// PickToResponse converts a database Pick to a PickResponse
func PickToResponse(pick database.Pick) PickResponse {
	response := PickResponse{
		Id:        pick.ID,
		UserId:    pick.UserID,
		GameId:    pick.GameID,
		Picked:    pick.Picked,
		Rank:      pick.Rank,
		QuickPick: pick.QuickPick,
		CreatedAt: pick.CreatedAt,
		UpdatedAt: pick.UpdatedAt,
	}

	// Include user if preloaded
	if pick.User.ID != 0 {
		user := UserToResponse(pick.User)
		response.User = &user
	}

	// Include game if preloaded
	if pick.Game.ID != 0 {
		game := GameToResponse(pick.Game)
		response.Game = &game
	}

	return response
}

// PickFromRequest converts a PickRequest to a database Pick
func PickFromRequest(req PickRequest) database.Pick {
	pick := database.Pick{
		GameID:    req.GameId,
		Picked:    req.Picked,
		Rank:      req.Rank,
		QuickPick: req.QuickPick,
	}

	// Handle optional UserId
	if req.UserId != nil {
		pick.UserID = *req.UserId
	}

	return pick
}

// PicksFromRequest converts multiple PickRequest to database Picks
func PicksFromRequest(reqs []PickRequest) []database.Pick {
	picks := make([]database.Pick, len(reqs))
	for i, req := range reqs {
		picks[i] = PickFromRequest(req)
	}
	return picks
}

// ResultToResponse converts a database Result to a ResultResponse
func ResultToResponse(result database.Result) ResultResponse {
	response := ResultResponse{
		Id:            result.ID,
		GameId:        result.GameID,
		FavoriteScore: result.FavoriteScore,
		UnderdogScore: result.UnderdogScore,
		Outcome:       result.Outcome,
		CreatedAt:     result.CreatedAt,
		UpdatedAt:     result.UpdatedAt,
	}

	// Include game if preloaded
	if result.Game.ID != 0 {
		game := GameToResponse(result.Game)
		response.Game = &game
	}

	return response
}

// ResultFromRequest converts a ResultRequest to a database Result
func ResultFromRequest(req ResultRequest) database.Result {
	return database.Result{
		GameID:        req.GameId,
		FavoriteScore: req.FavoriteScore,
		UnderdogScore: req.UnderdogScore,
		Outcome:       req.Outcome,
	}
}

// SurvivorPickToResponse converts a database SurvivorPick to a SurvivorPickResponse
func SurvivorPickToResponse(pick database.SurvivorPick) SurvivorPickResponse {
	response := SurvivorPickResponse{
		Id:        pick.ID,
		UserId:    pick.UserID,
		Week:      pick.Week,
		Team:      pick.Team,
		CreatedAt: pick.CreatedAt,
		UpdatedAt: pick.UpdatedAt,
	}

	// Include user if preloaded
	if pick.User.ID != 0 {
		user := UserToResponse(pick.User)
		response.User = &user
	}

	return response
}

// SurvivorPickFromRequest converts a SurvivorPickRequest to a database SurvivorPick
func SurvivorPickFromRequest(req SurvivorPickRequest) database.SurvivorPick {
	pick := database.SurvivorPick{
		Week: req.Week,
		Team: req.Team,
	}

	// Handle optional UserId
	if req.UserId != nil {
		pick.UserID = *req.UserId
	}

	return pick
}

// PlayerToResponse converts a database Player to a PlayerResponse
func PlayerToResponse(player database.Player) PlayerResponse {
	return PlayerResponse{
		Id:      player.ID,
		UserId:  player.UserID,
		Name:    player.Name,
		Address: player.Address,
	}
}

// PlayerFromRequest converts a PlayerRequest to a database Player
func PlayerFromRequest(req PlayerRequest, userID uint) database.Player {
	return database.Player{
		UserID:  userID,
		Name:    req.Name,
		Address: req.Address,
	}
}
