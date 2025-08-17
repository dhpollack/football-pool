package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	Connect("/tmp/database.db")
	if DB == nil {
		t.Fatal("failed to connect to database")
	}
}

func TestConnectTestDB(t *testing.T) {
	Connect("file::memory:?cache=shared")
	if DB == nil {
		t.Fatal("failed to connect to test database")
	}
}

func TestAutoMigrate(t *testing.T) {
	Connect("file::memory:?cache=shared")
	if err := DB.AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
}
