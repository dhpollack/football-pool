package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	db, err := New("/tmp/database.db")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	if db.GetDB() == nil {
		t.Fatal("database connection is nil")
	}
}

func TestConnectTestDB(t *testing.T) {
	db, err := New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	if db.GetDB() == nil {
		t.Fatal("test database connection is nil")
	}
}

func TestAutoMigrate(t *testing.T) {
	db, err := New("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.GetDB().AutoMigrate(&User{}, &Player{}, &Game{}, &Pick{}, &Result{}, &SurvivorPick{}); err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
}
