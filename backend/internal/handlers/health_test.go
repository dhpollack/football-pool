package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/david/football-pool/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	db, err := database.New(":memory:")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/api/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheck(db.GetDB()))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "", rr.Body.String())
}
