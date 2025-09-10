package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func HealthCheck(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log health check requests for debugging
		if r.Header.Get("User-Agent") != "" {
			slog.Debug(
				fmt.Sprintf("Health check request from: %s %s\n", r.Method, r.URL.Path),
			)
		}

		sqlDB, err := db.DB()
		if err != nil {
			http.Error(w, "Failed to get database connection", http.StatusInternalServerError)
			return
		}

		err = sqlDB.Ping()
		if err != nil {
			http.Error(w, "Database connection is not healthy", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}
