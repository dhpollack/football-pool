package database

import (
	"gorm.io/gorm"
	"time"
)

// User represents a user of the application
// swagger:model
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Role     string
	Player   Player
}

// Player represents a player in the football pool
// swagger:model
type Player struct {
	gorm.Model
	UserID  uint `gorm:"unique"`
	Name    string
	Address string
}

// Game represents a single game in a week
// swagger:model
type Game struct {
	gorm.Model
	Week         int `gorm:"index:idx_week_season"`
	Season       int `gorm:"index:idx_week_season"`
	FavoriteTeam string
	UnderdogTeam string
	Spread       float32
	StartTime    time.Time
}

// Pick represents a user's pick for a game
// swagger:model
type Pick struct {
	gorm.Model
	UserID    uint `gorm:"index:idx_user_game,unique"`
	User      User
	GameID    uint `gorm:"index:idx_user_game,unique"`
	Game      Game
	Picked    string
	Rank      int
	QuickPick bool
}

// Result represents the result of a game
// swagger:model
type Result struct {
	gorm.Model
	GameID        uint `gorm:"unique"`
	Game          Game
	FavoriteScore int
	UnderdogScore int
	Outcome       string
}

// SurvivorPick represents a user's survivor pick for a week
// swagger:model
type SurvivorPick struct {
	gorm.Model
	UserID uint `gorm:"index:idx_user_week,unique"`
	User   User
	Week   int `gorm:"index:idx_user_week,unique"`
	Team   string
}