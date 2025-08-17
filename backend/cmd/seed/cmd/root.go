package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seed",
	Short: "A CLI to seed the football pool database",
	Long:  `A command-line interface application to seed the football pool database with initial data.`,
}

var dsn string // Declare dsn variable

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "football-pool.db", "Database DSN") // Add dsn flag
}
