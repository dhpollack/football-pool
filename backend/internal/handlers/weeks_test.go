package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dhpollack/football-pool/internal/api"
	"github.com/dhpollack/football-pool/internal/database"
)

func TestListWeeks(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create test weeks
	weeks := []database.Week{
		{
			WeekNumber:    1,
			Season:        2025,
			WeekStartTime: time.Date(2025, 9, 7, 0, 0, 0, 0, time.UTC),
			WeekEndTime:   time.Date(2025, 9, 10, 23, 59, 59, 0, time.UTC),
			IsActive:      false,
		},
		{
			WeekNumber:    2,
			Season:        2025,
			WeekStartTime: time.Date(2025, 9, 14, 0, 0, 0, 0, time.UTC),
			WeekEndTime:   time.Date(2025, 9, 17, 23, 59, 59, 0, time.UTC),
			IsActive:      true,
		},
	}

	for _, week := range weeks {
		if result := gormDB.Create(&week); result.Error != nil {
			t.Fatalf("Failed to create test week: %v", result.Error)
		}
	}

	req := httptest.NewRequest("GET", "/api/admin/weeks", nil)
	w := httptest.NewRecorder()

	ListWeeks(gormDB)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response api.WeekListResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response.Weeks) != 2 {
		t.Errorf("Expected 2 weeks, got %d", len(response.Weeks))
	}

	// Verify week data
	foundWeek1 := false
	foundWeek2 := false
	for _, week := range response.Weeks {
		if week.WeekNumber == 1 && week.Season == 2025 {
			foundWeek1 = true
		}
		if week.WeekNumber == 2 && week.Season == 2025 && week.IsActive {
			foundWeek2 = true
		}
	}

	if !foundWeek1 {
		t.Error("Week 1 not found in response")
	}
	if !foundWeek2 {
		t.Error("Week 2 not found in response")
	}
}

func TestCreateWeek(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	weekRequest := api.WeekRequest{
		WeekNumber:    3,
		Season:        2025,
		WeekStartTime: time.Date(2025, 9, 21, 0, 0, 0, 0, time.UTC),
		WeekEndTime:   time.Date(2025, 9, 24, 23, 59, 59, 0, time.UTC),
		IsActive:      false,
	}

	body, _ := json.Marshal(weekRequest)
	req := httptest.NewRequest("POST", "/api/admin/weeks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateWeek(gormDB)(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response api.WeekResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.WeekNumber != 3 {
		t.Errorf("Expected week number 3, got %d", response.WeekNumber)
	}
	if response.Season != 2025 {
		t.Errorf("Expected season 2025, got %d", response.Season)
	}
	if response.IsActive {
		t.Error("Expected week to be inactive")
	}

	// Verify week was created in database
	var week database.Week
	if result := gormDB.First(&week, response.Id); result.Error != nil {
		t.Fatalf("Failed to find created week: %v", result.Error)
	}
}

func TestCreateWeek_InvalidRequest(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Invalid JSON
	req := httptest.NewRequest("POST", "/api/admin/weeks", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateWeek(gormDB)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateWeek(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a test week
	week := database.Week{
		WeekNumber:    1,
		Season:        2025,
		WeekStartTime: time.Date(2025, 9, 7, 0, 0, 0, 0, time.UTC),
		WeekEndTime:   time.Date(2025, 9, 10, 23, 59, 59, 0, time.UTC),
		IsActive:      false,
	}
	if result := gormDB.Create(&week); result.Error != nil {
		t.Fatalf("Failed to create test week: %v", result.Error)
	}

	weekRequest := api.WeekRequest{
		WeekNumber:    1,
		Season:        2025,
		WeekStartTime: time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),     // Updated start time
		WeekEndTime:   time.Date(2025, 9, 11, 23, 59, 59, 0, time.UTC), // Updated end time
		IsActive:      true,                                            // Updated active status
	}

	body, _ := json.Marshal(weekRequest)
	pathParams := map[string]string{"id": "1"}
	req := createRequestWithPathParams("PUT", "/api/admin/weeks/1", bytes.NewBuffer(body), pathParams)

	w := httptest.NewRecorder()

	UpdateWeek(gormDB)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response api.WeekResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.IsActive {
		t.Error("Expected week to be active after update")
	}

	// Verify week was updated in database
	var updatedWeek database.Week
	if result := gormDB.First(&updatedWeek, week.ID); result.Error != nil {
		t.Fatalf("Failed to find updated week: %v", result.Error)
	}

	if !updatedWeek.IsActive {
		t.Error("Week should be active in database")
	}
}

func TestUpdateWeek_NotFound(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	weekRequest := api.WeekRequest{
		WeekNumber:    1,
		Season:        2025,
		WeekStartTime: time.Date(2025, 9, 7, 0, 0, 0, 0, time.UTC),
		WeekEndTime:   time.Date(2025, 9, 10, 23, 59, 59, 0, time.UTC),
		IsActive:      false,
	}

	body, _ := json.Marshal(weekRequest)
	pathParams := map[string]string{"id": "999"}
	req := createRequestWithPathParams("PUT", "/api/admin/weeks/999", bytes.NewBuffer(body), pathParams)
	w := httptest.NewRecorder()

	UpdateWeek(gormDB)(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteWeek(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create a test week
	week := database.Week{
		WeekNumber:    1,
		Season:        2025,
		WeekStartTime: time.Date(2025, 9, 7, 0, 0, 0, 0, time.UTC),
		WeekEndTime:   time.Date(2025, 9, 10, 23, 59, 59, 0, time.UTC),
		IsActive:      false,
	}
	if result := gormDB.Create(&week); result.Error != nil {
		t.Fatalf("Failed to create test week: %v", result.Error)
	}

	pathParams := map[string]string{"id": "1"}
	req := createRequestWithPathParams("DELETE", "/api/admin/weeks/1", nil, pathParams)
	w := httptest.NewRecorder()

	DeleteWeek(gormDB)(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify week was deleted from database
	var deletedWeek database.Week
	if result := gormDB.First(&deletedWeek, week.ID); result.Error == nil {
		t.Error("Week should have been deleted from database")
	}
}

func TestDeleteWeek_NotFound(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	pathParams := map[string]string{"id": "999"}
	req := createRequestWithPathParams("DELETE", "/api/admin/weeks/999", nil, pathParams)
	w := httptest.NewRecorder()

	DeleteWeek(gormDB)(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestActivateWeek(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	// Create test weeks
	weeks := []database.Week{
		{
			WeekNumber:    1,
			Season:        2025,
			WeekStartTime: time.Date(2025, 9, 7, 0, 0, 0, 0, time.UTC),
			WeekEndTime:   time.Date(2025, 9, 10, 23, 59, 59, 0, time.UTC),
			IsActive:      true, // Initially active
		},
		{
			WeekNumber:    2,
			Season:        2025,
			WeekStartTime: time.Date(2025, 9, 14, 0, 0, 0, 0, time.UTC),
			WeekEndTime:   time.Date(2025, 9, 17, 23, 59, 59, 0, time.UTC),
			IsActive:      false, // Initially inactive
		},
	}

	for i := range weeks {
		if result := gormDB.Create(&weeks[i]); result.Error != nil {
			t.Fatalf("Failed to create test week: %v", result.Error)
		}
	}

	pathParams := map[string]string{"id": "2"}
	req := createRequestWithPathParams("POST", "/api/admin/weeks/2/activate", nil, pathParams)
	w := httptest.NewRecorder()

	ActivateWeek(gormDB)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response api.WeekResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.IsActive {
		t.Error("Expected week 2 to be active")
	}

	// Verify week 2 is active and week 1 is inactive
	var week1 database.Week
	if result := gormDB.First(&week1, weeks[0].ID); result.Error != nil {
		t.Fatalf("Failed to find week 1: %v", result.Error)
	}
	if week1.IsActive {
		t.Error("Week 1 should be inactive after activating week 2")
	}

	var week2 database.Week
	if result := gormDB.First(&week2, weeks[1].ID); result.Error != nil {
		t.Fatalf("Failed to find week 2: %v", result.Error)
	}
	if !week2.IsActive {
		t.Error("Week 2 should be active")
	}
}

func TestActivateWeek_NotFound(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	pathParams := map[string]string{"id": "999"}
	req := createRequestWithPathParams("POST", "/api/admin/weeks/999/activate", nil, pathParams)
	w := httptest.NewRecorder()

	ActivateWeek(gormDB)(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestActivateWeek_InvalidID(t *testing.T) {
	db, err := database.New("sqlite", "file::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	gormDB := db.GetDB()

	req := httptest.NewRequest("POST", "/api/admin/weeks/invalid/activate", nil)
	w := httptest.NewRecorder()

	ActivateWeek(gormDB)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
