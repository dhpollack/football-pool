package main

import (
	"github.com/david/football-pool/internal/database"
	"github.com/david/football-pool/internal/server"
)

func main() {
	database.Connect()
	server.Start()
}
