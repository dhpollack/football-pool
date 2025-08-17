package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/david/football-pool/internal/database"
)

func TestStart(t *testing.T) {
	database.Connect("file::memory:?cache=shared")
	go Start()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	req, err := http.NewRequest("GET", "http://localhost:8080/api/games?week=1&season=2023", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
