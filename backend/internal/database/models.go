package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Role     string
}

type Player struct {
	gorm.Model
	UserID  uint
	User    User
	Name    string
	Address string
}

type Game struct {
	gorm.Model
	Week         int
	Season       int
	FavoriteTeam string
	UnderdogTeam string
	Spread       float32
}

type Pick struct {
	gorm.Model
	UserID     uint
	User       User
	GameID     uint
	Game       Game
	PickedTeam string
	Rank       int
	QuickPick  bool
}

type Result struct {
	gorm.Model
	GameID        uint
	Game          Game
	FavoriteScore int
	UnderdogScore int
	Outcome       string
}

type SurvivorPick struct {
	gorm.Model
	UserID uint
	User   User
	Week   int
	Team   string
}
