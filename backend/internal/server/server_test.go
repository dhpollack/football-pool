package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david/football-pool/internal/config"
	"github.com/david/football-pool/internal/database"
)

func TestStart(t *testing.T) {
	t.Setenv("FOOTBALL_POOL_ENV", "test")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}
	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	server := NewServer(db, cfg)

	// Test that the router can be created without errors
	router := server.NewRouter()
	if router == nil {
		t.Error("Expected router to be created, got nil")
	}

	// Test that we can create a request to the games endpoint
	req, err := http.NewRequest("GET", "/api/games?week=1&season=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Use httptest recorder to test the handler without starting a real server
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestStartWithPortFromEnv(t *testing.T) {
	t.Setenv("FOOTBALL_POOL_PORT", "8081")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}
	db, err := database.New(cfg.Database.DSN)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	server := NewServer(db, cfg)

	// Test that the router can be created without errors
	router := server.NewRouter()
	if router == nil {
		t.Error("Expected router to be created, got nil")
	}

	// Test that we can create a request to the games endpoint
	req, err := http.NewRequest("GET", "/api/games?week=1&season=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Use httptest recorder to test the handler without starting a real server
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify that the configuration picked up the environment variable
	if cfg.Server.Port != "8081" {
		t.Errorf("Expected port 8081 from environment variable, got %s", cfg.Server.Port)
	}
}
