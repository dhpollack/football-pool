package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	Connect()
	if DB == nil {
		t.Fatal("failed to connect to database")
	}
}

func TestConnectTestDB(t *testing.T) {
	ConnectTestDB()
	if DB == nil {
		t.Fatal("failed to connect to test database")
	}
}

func TestAutoMigrate(t *testing.T) {
	ConnectTestDB()
	if err := DB.AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
}
