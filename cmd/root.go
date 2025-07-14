package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "streakr",
	Short: "streakr is a habit tracking CLI",
	Long: `streakr is a command-line tool for tracking habits and maintaining streaks.
It allows users to add, view, and manage their habits directly from the terminal.
It helps you improve good habits (like exercising, reading, etc.) and
also helps you quit bad habits (like doomscrolling, junkfood, etc.).
You can track your progress, set goals, and stay motivated with streakr.`,
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
