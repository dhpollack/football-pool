package main

import (
	"log/slog"

	"github.com/david/football-pool/internal/database"
	"github.com/david/football-pool/internal/server"
)

func main() {
	slog.Info("Connecting to database")
	database.Connect()
	slog.Info("Starting server")
	server.Start()
}
