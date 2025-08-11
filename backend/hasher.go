package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run hasher.go <password>")
		return
	}

	password := os.Args[1]
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(hashedPassword))
}
