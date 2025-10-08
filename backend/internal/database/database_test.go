package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	db, err := New("sqlite", "/tmp/database.db")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	if db.GetDB() == nil {
		t.Fatal("database connection is nil")
	}
}

func TestConnectTestDB(t *testing.T) {
	db, err := New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	if db.GetDB() == nil {
		t.Fatal("test database connection is nil")
	}
}

func TestAutoMigrate(t *testing.T) {
	db, err := New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.GetDB().AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
}

func TestWeekHasGames(t *testing.T) {
	db, err := New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	season := 2025
	week := 1

	// Initially, no games should exist for the week
	hasGames, err := db.WeekHasGames(season, week)
	if err != nil {
		t.Fatalf("failed to check if week has games: %v", err)
	}
	if hasGames {
		t.Fatal("expected week to have no games, but it did")
	}

	// Add a game for the week
	game := &Game{Season: season, Week: week, FavoriteTeam: "Team A", UnderdogTeam: "Team B"}
	if err := db.GetDB().Create(game).Error; err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	// Now, games should exist for the week
	hasGames, err = db.WeekHasGames(season, week)
	if err != nil {
		t.Fatalf("failed to check if week has games: %v", err)
	}
	if !hasGames {
		t.Fatal("expected week to have games, but it didn't")
	}
}
