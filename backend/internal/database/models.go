package database

import "gorm.io/gorm"

// User represents a user of the application
// swagger:model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Player represents a player in the football pool
// swagger:model
type Player struct {
	gorm.Model
	UserID  uint   `json:"user_id"`
	User    User   `json:"user"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// Game represents a single game in a week
// swagger:model
type Game struct {
	gorm.Model
	Week         int     `json:"week"`
	Season       int     `json:"season"`
	FavoriteTeam string  `json:"favorite_team"`
	UnderdogTeam string  `json:"underdog_team"`
	Spread       float32 `json:"spread"`
}

// Pick represents a user's pick for a game
// swagger:model
type Pick struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	User      User   `json:"user"`
	GameID    uint   `json:"game_id"`
	Game      Game   `json:"game"`
	Picked    string `json:"picked"`
	Rank      int    `json:"rank"`
	QuickPick bool   `json:"quick_pick"`
}

// Result represents the result of a game
// swagger:model
type Result struct {
	gorm.Model
	GameID        uint `json:"game_id"`
	Game          Game `json:"game"`
	FavoriteScore int  `json:"favorite_score"`
	UnderdogScore int  `json:"underdog_score"`
	Outcome       string `json:"outcome"`
}

// SurvivorPick represents a user's survivor pick for a week
// swagger:model
type SurvivorPick struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	User   User   `json:"user"`
	Week   int    `json:"week"`
	Team   string `json:"team"`
}