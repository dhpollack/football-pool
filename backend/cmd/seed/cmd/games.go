package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/david/football-pool/internal/database"
	"github.com/spf13/cobra"
)

var seedGamesCmd = &cobra.Command{
	Use:   "games <json_file_path>",
	Short: "Seed the games table",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.New(dsn)
		if err != nil {
			fmt.Printf("Error connecting to database: %v\n", err)
			return
		}
		gormDB := db.GetDB()

		jsonFilePath := args[0]
		jsonFile, err := os.Open(jsonFilePath)
		if err != nil {
			fmt.Printf("Error opening JSON file: %v\n", err)
			return
		}
		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)

		var games []database.Game
		if err := json.Unmarshal(byteValue, &games); err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			return
		}

		for _, game := range games {
			var existingGame database.Game
			result := gormDB.Where("week = ? AND season = ? AND favorite_team = ? AND underdog_team = ?",
				game.Week, game.Season, game.FavoriteTeam, game.UnderdogTeam).First(&existingGame)

			if result.Error == nil {
				fmt.Printf("Game already exists: %s vs %s (Week %d, Season %d). Skipping.\n", game.FavoriteTeam, game.UnderdogTeam, game.Week, game.Season)
				continue
			} else if result.Error != nil && result.Error.Error() != "record not found" {
				fmt.Printf("Error checking for existing game: %v\n", result.Error)
				continue
			}

			if result := gormDB.Create(&game); result.Error != nil {
				fmt.Printf("Error seeding game: %v\n", result.Error)
			}
		}

		fmt.Println("Games seeded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(seedGamesCmd)
}
