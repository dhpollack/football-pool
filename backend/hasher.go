// Package main provides a utility for hashing passwords.
package main

import (
	"log/slog"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("Usage: go run hasher.go <password>")
		return
	}

	password := os.Args[1]
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		return
	}

	slog.Info(string(hashedPassword))
}
