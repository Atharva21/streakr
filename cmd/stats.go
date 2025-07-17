package cmd

import (
	"fmt"
	"strings"

	"github.com/Atharva21/streakr/internal/service"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Stats gives current streak, max streak",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
		}
		if len(args) > 1 {
			return &se.StreakrError{TerminalMsg: "habit name should not be more than 1 word"}
		}

		habitName := strings.TrimSpace(args[0])
		if habitName == "" {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
		}
		if len(habitName) > 20 {
			return &se.StreakrError{TerminalMsg: "habit name cannot exceed 20 characters"}
		}
		habitName = strings.ToLower(habitName)
		currentStreak, maxStreak, err := service.GetStatsForHabitName(cmd.Context(), habitName)
		if err != nil {
			return err
		}
		fmt.Printf("current streak: %d, max streak: %d\n", currentStreak, maxStreak)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.InitDefaultHelpFlag()
	statsCmd.Flags().Lookup("help").Shorthand = ""
}
