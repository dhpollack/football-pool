// Package server provides HTTP server configuration and routing for the football pool application.
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dhpollack/football-pool/internal/auth"
	"github.com/dhpollack/football-pool/internal/config"
	"github.com/dhpollack/football-pool/internal/database"
	"github.com/dhpollack/football-pool/internal/handlers"
	"github.com/rs/cors"
)

// Server represents the HTTP server with database and authentication components.
type Server struct {
	db   *database.Database
	auth *auth.Auth
	cfg  *config.Config
}

// NewServer creates a new Server instance with the provided database connection.
func NewServer(db *database.Database, cfg *config.Config) *Server {
	return &Server{
		db:   db,
		auth: auth.NewAuth(db),
		cfg:  cfg,
	}
}

// NewRouter creates and configures the HTTP router with all application routes.
func (s *Server) NewRouter() http.Handler {
	// Get CORS allowed origins from environment variable or use defaults
	allowedOrigins := []string{"http://localhost:13000", "https://localhost:13001", "http://localhost:5173"}
	if corsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); corsEnv != "" {
		// Split comma-separated origins without using strings package
		origins := []string{}
		start := 0
		for i, char := range corsEnv {
			if char == ',' {
				origins = append(origins, corsEnv[start:i])
				start = i + 1
			}
		}
		// Add the last segment
		if start < len(corsEnv) {
			origins = append(origins, corsEnv[start:])
		}
		allowedOrigins = origins
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/login", s.auth.Login)
	mux.HandleFunc("POST /api/logout", s.auth.Logout)
	mux.HandleFunc("POST /api/register", s.auth.Register)
	mux.HandleFunc("GET /api/health", handlers.HealthCheck(s.db.GetDB()))

	mux.Handle("GET /api/users/me", s.auth.Middleware(handlers.GetProfile(s.db.GetDB())))
	mux.Handle("PUT /api/users/me/update", s.auth.Middleware(handlers.UpdateProfile(s.db.GetDB())))

	mux.HandleFunc("GET /api/games", handlers.GetGames(s.db.GetDB()))

	// Admin game management endpoints
	mux.Handle("GET /api/admin/games", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminListGames(s.db.GetDB()))))
	mux.Handle("POST /api/admin/games/create", s.auth.Middleware(s.auth.AdminMiddleware(handlers.CreateGame(s.db.GetDB()))))
	mux.Handle("PUT /api/admin/games/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.UpdateGame(s.db.GetDB()))))
	mux.Handle("DELETE /api/admin/games/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.DeleteGame(s.db.GetDB()))))

	mux.Handle("GET /api/picks", s.auth.Middleware(handlers.GetPicks(s.db.GetDB())))
	mux.Handle("POST /api/picks/submit", s.auth.Middleware(handlers.SubmitPicks(s.db.GetDB())))
	mux.Handle("POST /api/admin/picks/submit", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminSubmitPicks(s.db.GetDB()))))

	// Admin pick management endpoints
	mux.Handle("GET /api/admin/picks", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminListPicks(s.db.GetDB()))))
	mux.Handle("GET /api/admin/picks/week/{week}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminGetPicksByWeek(s.db.GetDB()))))
	mux.Handle("GET /api/admin/picks/user/{userID}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminGetPicksByUser(s.db.GetDB()))))
	mux.Handle("DELETE /api/admin/picks/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminDeletePick(s.db.GetDB()))))

	mux.HandleFunc("GET /api/results/week", handlers.GetWeeklyResults(s.db.GetDB()))
	mux.HandleFunc("GET /api/results/season", handlers.GetSeasonResults(s.db.GetDB()))

	mux.Handle("POST /api/results", s.auth.Middleware(s.auth.AdminMiddleware(handlers.SubmitResult(s.db.GetDB()))))

	mux.Handle("GET /api/survivor/picks", s.auth.Middleware(handlers.GetSurvivorPicks(s.db.GetDB())))
	mux.Handle("POST /api/survivor/picks/submit", s.auth.Middleware(handlers.SubmitSurvivorPick(s.db.GetDB())))

	mux.Handle("DELETE /api/admin/users/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.DeleteUser(s.db.GetDB()))))
	mux.Handle("DELETE /api/admin/users/delete", s.auth.Middleware(s.auth.AdminMiddleware(handlers.DeleteUserByEmail(s.db.GetDB()))))

	// Admin user management endpoints
	mux.Handle("GET /api/admin/users", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminListUsers(s.db.GetDB()))))
	mux.Handle("POST /api/admin/users/create", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminCreateUsers(s.db.GetDB()))))
	mux.Handle("GET /api/admin/users/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminGetUser(s.db.GetDB()))))
	mux.Handle("PUT /api/admin/users/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminUpdateUser(s.db.GetDB()))))

	// Week management endpoints
	mux.Handle("GET /api/admin/weeks", s.auth.Middleware(s.auth.AdminMiddleware(handlers.ListWeeks(s.db.GetDB()))))
	mux.Handle("POST /api/admin/weeks", s.auth.Middleware(s.auth.AdminMiddleware(handlers.CreateWeek(s.db.GetDB()))))
	mux.Handle("PUT /api/admin/weeks/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.UpdateWeek(s.db.GetDB()))))
	mux.Handle("DELETE /api/admin/weeks/{id}", s.auth.Middleware(s.auth.AdminMiddleware(handlers.DeleteWeek(s.db.GetDB()))))
	mux.Handle("POST /api/admin/weeks/{id}/activate", s.auth.Middleware(s.auth.AdminMiddleware(handlers.ActivateWeek(s.db.GetDB()))))

	return c.Handler(mux)
}

// Start begins listening for HTTP requests and serves the application.
func (s *Server) Start() {
	addr := s.cfg.Server.Host + ":" + s.cfg.Server.Port

	// Log the server URL and port for debugging
	slog.Info(fmt.Sprintf("Starting server on http://%s", addr))
	slog.Debug(fmt.Sprintf("Health endpoint: http://%s/api/health", addr))

	server := &http.Server{
		Addr:         addr,
		Handler:      s.NewRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
