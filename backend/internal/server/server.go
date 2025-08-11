package server

import (
	"net/http"

	"github.com/david/football-pool/internal/auth"
	"github.com/david/football-pool/internal/handlers"
)

func Start() {
	http.HandleFunc("/api/login", auth.Login)
	http.HandleFunc("/api/logout", auth.Logout)
	http.HandleFunc("/api/register", auth.Register)

	http.Handle("/api/users/me", auth.Middleware(http.HandlerFunc(handlers.GetProfile)))
	http.Handle("/api/users/me/update", auth.Middleware(http.HandlerFunc(handlers.UpdateProfile)))

	http.HandleFunc("/api/games", handlers.GetGames)

	http.Handle("/api/picks", auth.Middleware(http.HandlerFunc(handlers.GetPicks)))
	http.Handle("/api/picks/submit", auth.Middleware(http.HandlerFunc(handlers.SubmitPicks)))

	http.HandleFunc("/api/results/week", handlers.GetWeeklyResults)
	http.HandleFunc("/api/results/season", handlers.GetSeasonResults)

	http.Handle("/api/results", auth.Middleware(auth.AdminMiddleware(http.HandlerFunc(handlers.SubmitResult))))

	http.Handle("/api/survivor/picks", auth.Middleware(http.HandlerFunc(handlers.GetSurvivorPicks)))
	http.Handle("/api/survivor/picks/submit", auth.Middleware(http.HandlerFunc(handlers.SubmitSurvivorPick)))

	http.ListenAndServe(":8080", nil)
}
