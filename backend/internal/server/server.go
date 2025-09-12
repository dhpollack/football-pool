package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/database"
	"github.com/david/football-pool/internal/handlers"
	"github.com/rs/cors"
)

type Server struct {
	db   *database.Database
	auth *auth.Auth
}

func NewServer(db *database.Database) *Server {
	return &Server{
		db:   db,
		auth: auth.NewAuth(db),
	}
}

func (s *Server) NewRouter() http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:13000", "https://localhost:13001", "http://localhost:5173"},
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

	return c.Handler(mux)
}

func (s *Server) Start() {
	port := os.Getenv("FOOTBALL_POOL_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("FOOTBALL_POOL_HOST")
	if host == "" {
		host = "localhost"
	}

	addr := host + ":" + port

	// Log the server URL and port for debugging
	slog.Info(fmt.Sprintf("Starting server on http://%s", addr))
	slog.Debug(fmt.Sprintf("Health endpoint: http://%s/api/health", addr))

	if err := http.ListenAndServe(addr, s.NewRouter()); err != nil {
		panic(err)
	}
}
