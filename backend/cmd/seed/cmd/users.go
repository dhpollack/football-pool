package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/david/football-pool/internal/database"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var seedUsersCmd = &cobra.Command{
	Use:   "users <json_file_path>",
	Short: "Seed the users table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		database.Connect(dsn)

		jsonFilePath := args[0]
		jsonFile, err := os.Open(jsonFilePath)
		if err != nil {
			fmt.Printf("Error opening JSON file: %v\n", err)
			return
		}
		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)

		var users []struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := json.Unmarshal(byteValue, &users); err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			return
		}

		for _, user := range users {
			var existingUser database.User
			result := database.DB.Where("email = ?", user.Email).First(&existingUser)

			if result.Error == nil {
				fmt.Printf("User with email %s already exists. Skipping.\n", user.Email)
				continue
			} else if result.Error != nil && result.Error.Error() != "record not found" {
				fmt.Printf("Error checking for existing user %s: %v\n", user.Email, result.Error)
				continue
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
			if err != nil {
				fmt.Printf("Error hashing password for user %s: %v\n", user.Email, err)
				continue
			}
			newUser := database.User{Name: user.Name, Email: user.Email, Password: string(hashedPassword), Role: user.Role}
			if result := database.DB.Create(&newUser); result.Error != nil {
				fmt.Printf("Error seeding user %s: %v\n", user.Email, result.Error)
			}
		}

		fmt.Println("Users seeded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(seedUsersCmd)
}
