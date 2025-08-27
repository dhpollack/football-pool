package server

import (
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
		AllowedOrigins:   []string{"http://localhost:13000", "http://localhost:5173", "http://localhost:5174"},
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
	mux.Handle("/api/games/create", s.auth.Middleware(s.auth.AdminMiddleware(handlers.CreateGame(s.db.GetDB()))))

	mux.Handle("/api/picks", s.auth.Middleware(handlers.GetPicks(s.db.GetDB())))
	mux.Handle("/api/picks/submit", s.auth.Middleware(handlers.SubmitPicks(s.db.GetDB())))
	mux.Handle("/api/admin/picks/submit", s.auth.Middleware(s.auth.AdminMiddleware(handlers.AdminSubmitPicks(s.db.GetDB()))))

	mux.HandleFunc("/api/results/week", handlers.GetWeeklyResults(s.db.GetDB()))
	mux.HandleFunc("/api/results/season", handlers.GetSeasonResults(s.db.GetDB()))

	mux.Handle("/api/results", s.auth.Middleware(s.auth.AdminMiddleware(handlers.SubmitResult(s.db.GetDB()))))

	mux.Handle("/api/survivor/picks", s.auth.Middleware(handlers.GetSurvivorPicks(s.db.GetDB())))
	mux.Handle("/api/survivor/picks/submit", s.auth.Middleware(handlers.SubmitSurvivorPick(s.db.GetDB())))

	mux.Handle("/api/debug/users", s.auth.Middleware(s.auth.AdminMiddleware(handlers.DebugGetUsers(s.db.GetDB()))))
	mux.Handle("/api/debug/users/delete", s.auth.Middleware(s.auth.AdminMiddleware(handlers.DebugDeleteUser(s.db.GetDB()))))

	return c.Handler(mux)
}

func (s *Server) Start() {
	port := os.Getenv("FOOTBALL_POOL_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	if err := http.ListenAndServe(addr, s.NewRouter()); err != nil {
		panic(err)
	}
}
