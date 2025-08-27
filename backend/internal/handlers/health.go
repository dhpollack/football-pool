package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

func HealthCheck(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.Write([]byte("Healthy"))
	}
}
