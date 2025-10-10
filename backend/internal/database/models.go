// Package database provides data models and database operations for the football pool application.
package database

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user of the application
// swagger:model
type User struct {
	gorm.Model
	Name     string `validate:"required"`
	Email    string `gorm:"unique" validate:"required,email"`
	Password string `validate:"required"`
	Role     string `validate:"required"`
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
	Week      int       `gorm:"index:idx_game_week_season" validate:"required,ne=0"`
	Season    int       `gorm:"index:idx_game_week_season" validate:"required,ne=0"`
	HomeTeam  string    `gorm:"column:favorite_team" validate:"required"`
	AwayTeam  string    `gorm:"column:underdog_team" validate:"required"`
	Favorite  *string   `validate:"omitempty,oneof=Home Away"`
	Underdog  *string   `validate:"omitempty,oneof=Home Away"`
	Spread    float32   `validate:"gte=0"`
	StartTime time.Time `validate:"required"`
}

// Pick represents a user's pick for a game
// swagger:model
type Pick struct {
	gorm.Model
	UserID    uint   `gorm:"index:idx_user_game,unique" validate:"required"`
	User      User   `validate:"-"`
	GameID    uint   `gorm:"index:idx_user_game,unique" validate:"required"`
	Game      Game   `validate:"-"`
	Picked    string `validate:"required"`
	Rank      int    `validate:"required"`
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

// Week represents a week in the football season
// swagger:model
type Week struct {
	gorm.Model
	WeekNumber    int       `gorm:"index:idx_week_number_season,unique" validate:"required,ne=0"`
	Season        int       `gorm:"index:idx_week_number_season,unique" validate:"required,ne=0"`
	WeekStartTime time.Time `validate:"required"`
	WeekEndTime   time.Time `validate:"required"`
	IsActive      bool      `gorm:"default:false"`
}
