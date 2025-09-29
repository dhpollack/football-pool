// Package apiespn provides tests for the ESPN API client.
package apiespn

import (
	"encoding/json"
	"os"
	"testing"
)

func TestScoreboardUnmarshal(t *testing.T) {
	// Read the sample JSON file
	data, err := os.ReadFile("../../assets/test/scoreboard-sample.json")
	if err != nil {
		t.Fatalf("Failed to read sample file: %v", err)
	}

	// Try to unmarshal into Scoreboard struct
	// The ESPNDateTime type should handle the ESPN time format automatically
	var scoreboard Scoreboard
	err = json.Unmarshal(data, &scoreboard)
	if err != nil {
		t.Fatalf("Failed to unmarshal scoreboard: %v", err)
	}

	// Verify basic structure
	if scoreboard.Events == nil {
		t.Error("Scoreboard.Events should not be nil")
	}

	if scoreboard.Season == nil {
		t.Error("Scoreboard.Season should not be nil")
	}

	if scoreboard.Week == nil {
		t.Error("Scoreboard.Week should not be nil")
	}

	// Check if we have events
	if scoreboard.Events != nil && len(*scoreboard.Events) == 0 {
		t.Error("Expected at least one event in the scoreboard")
	}

	// Check season and week data
	if scoreboard.Season != nil && scoreboard.Season.Year == nil {
		t.Error("Season.Year should not be nil")
	}

	if scoreboard.Week != nil && scoreboard.Week.Number == nil {
		t.Error("Week.Number should not be nil")
	}

	t.Logf("Successfully unmarshalled scoreboard with %d events", len(*scoreboard.Events))
}

func TestEventUnmarshal(t *testing.T) {
	// Test with a minimal event structure
	minimalEventJSON := `{
		"id": "test-event",
		"name": "Test Game",
		"date": "2024-09-05T20:20:00Z",
		"competitions": [{
			"competitors": [{
				"homeAway": "home",
				"team": {"displayName": "Team A"},
				"score": "24"
			}, {
				"homeAway": "away",
				"team": {"displayName": "Team B"},
				"score": "17"
			}]
		}]
	}`

	var event Event
	err := json.Unmarshal([]byte(minimalEventJSON), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal minimal event: %v", err)
	}

	if event.Id == nil || *event.Id != "test-event" {
		t.Error("Event ID should be 'test-event'")
	}

	if event.Competitions == nil || len(*event.Competitions) == 0 {
		t.Error("Event should have competitions")
	}

	t.Log("Successfully unmarshalled minimal event")
}
