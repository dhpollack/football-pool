package server

import (
	"net/http"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/handlers"
	"github.com/rs/cors"
)

func Start() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:13000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", auth.Login)
	mux.HandleFunc("/api/logout", auth.Logout)
	mux.HandleFunc("/api/register", auth.Register)

	mux.Handle("/api/users/me", auth.Middleware(http.HandlerFunc(handlers.GetProfile)))
	mux.Handle("/api/users/me/update", auth.Middleware(http.HandlerFunc(handlers.UpdateProfile)))

	mux.HandleFunc("/api/games", handlers.GetGames)

	mux.Handle("/api/picks", auth.Middleware(http.HandlerFunc(handlers.SubmitPicks)))
	mux.Handle("/api/picks/submit", auth.Middleware(http.HandlerFunc(handlers.SubmitPicks)))

	mux.HandleFunc("/api/results/week", handlers.GetWeeklyResults)
	mux.HandleFunc("/api/results/season", handlers.GetSeasonResults)

	mux.Handle("/api/results", auth.Middleware(auth.AdminMiddleware(http.HandlerFunc(handlers.SubmitResult))))

	mux.Handle("/api/survivor/picks", auth.Middleware(http.HandlerFunc(handlers.GetSurvivorPicks)))
	mux.Handle("/api/survivor/picks/submit", auth.Middleware(http.HandlerFunc(handlers.SubmitSurvivorPick)))

	mux.Handle("/api/debug/users", auth.Middleware(auth.AdminMiddleware(http.HandlerFunc(handlers.DebugGetUsers))))
	mux.Handle("/api/debug/users/delete", auth.Middleware(auth.AdminMiddleware(http.HandlerFunc(handlers.DebugDeleteUser))))

	if err := http.ListenAndServe(":8080", c.Handler(mux)); err != nil {
		panic(err)
	}
}
