package main

import (
	"log/slog"

	"github.com/david/football-pool/internal/database"
	"github.com/david/football-pool/internal/server"
)

func main() {
	slog.Info("Connecting to database")
	database.Connect("football-pool.db")
	slog.Info("Starting server")
	server.Start()
}
