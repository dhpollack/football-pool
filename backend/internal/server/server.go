// Package server provides HTTP server configuration and routing for the football pool application.
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
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
		allowedOrigins = strings.Split(corsEnv, ",")
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", s.auth.Login)
	mux.HandleFunc("/api/logout", s.auth.Logout)
	mux.HandleFunc("/api/register", s.auth.Register)
	mux.HandleFunc("/api/health", handlers.HealthCheck(s.db.GetDB()))

	mux.Handle("/api/users/me", s.auth.Middleware(handlers.GetProfile(s.db.GetDB())))
	mux.Handle("/api/users/me/update", s.auth.Middleware(handlers.UpdateProfile(s.db.GetDB())))

	mux.HandleFunc("/api/games", handlers.GetGames(s.db.GetDB()))

	// Admin game management endpoints
	mux.Handle("/api/admin/games", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminListGames(s.db.GetDB()))))
	mux.Handle("/api/admin/games/create", s.auth.Middleware(s.auth.AdminMiddleware(handlers.CreateGame(s.db.GetDB()))))
	mux.Handle("/api/admin/games/", s.auth.Middleware(s.auth.AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			handlers.UpdateGame(s.db.GetDB())(w, r)
		case "DELETE":
			handlers.DeleteGame(s.db.GetDB())(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/picks", s.auth.Middleware(handlers.GetPicks(s.db.GetDB())))
	mux.Handle("/api/picks/submit", s.auth.Middleware(handlers.SubmitPicks(s.db.GetDB())))
	mux.Handle("/api/admin/picks/submit", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminSubmitPicks(s.db.GetDB()))))

	// Admin pick management endpoints
	mux.Handle("/api/admin/picks", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminListPicks(s.db.GetDB()))))
	mux.Handle("/api/admin/picks/week/", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminGetPicksByWeek(s.db.GetDB()))))
	mux.Handle("/api/admin/picks/user/", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminGetPicksByUser(s.db.GetDB()))))
	mux.Handle("/api/admin/picks/", s.auth.Middleware(s.auth.AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			handlers.AdminDeletePick(s.db.GetDB())(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))))

	mux.HandleFunc("/api/results/week", handlers.GetWeeklyResults(s.db.GetDB()))
	mux.HandleFunc("/api/results/season", handlers.GetSeasonResults(s.db.GetDB()))

	mux.Handle("/api/results", s.auth.Middleware(s.auth.AdminMiddleware(handlers.SubmitResult(s.db.GetDB()))))

	mux.Handle("/api/survivor/picks", s.auth.Middleware(handlers.GetSurvivorPicks(s.db.GetDB())))
	mux.Handle("/api/survivor/picks/submit", s.auth.Middleware(handlers.SubmitSurvivorPick(s.db.GetDB())))

	mux.Handle("/api/admin/users/delete", s.auth.Middleware(s.auth.AdminMiddleware(handlers.DeleteUser(s.db.GetDB()))))

	// Admin user management endpoints
	mux.Handle("/api/admin/users", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminListUsers(s.db.GetDB()))))
	mux.Handle("/api/admin/users/create", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminCreateUsers(s.db.GetDB()))))
	mux.Handle("/api/admin/users/", s.auth.Middleware(s.auth.AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.AdminGetUser(s.db.GetDB())(w, r)
		case "PUT":
			handlers.AdminUpdateUser(s.db.GetDB())(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))))

	// Week management endpoints
	mux.Handle("/api/admin/weeks", s.auth.Middleware(s.auth.AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.ListWeeks(s.db.GetDB())(w, r)
		case "POST":
			handlers.CreateWeek(s.db.GetDB())(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/admin/weeks/{id}", s.auth.Middleware(s.auth.AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			handlers.UpdateWeek(s.db.GetDB())(w, r)
		case "DELETE":
			handlers.DeleteWeek(s.db.GetDB())(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))))

	mux.Handle("/api/admin/weeks/{id}/activate", s.auth.Middleware(s.auth.AdminMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			handlers.ActivateWeek(s.db.GetDB())(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))))

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
