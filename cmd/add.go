package cmd

import (
	"fmt"
	"strings"

	"github.com/Atharva21/streakr/internal/service"
	"github.com/Atharva21/streakr/internal/store"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new habit to track",
	Long: `Add will add new habit and start tracking the streaks.
Specify name of the habit in single word followed by streakr add

Few examples:
streakr add run --description "morning run 5kms"
streakr add read --description "read 5 pages of any book"
streakr add smoking --type quit
`,
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

		description, _ := cmd.Flags().GetString("description")
		habitType, _ := cmd.Flags().GetString("type")

		if len(description) > 200 {
			return &se.StreakrError{TerminalMsg: "description cannot exceed 200 characters"}
		}

		if habitType == "" {
			habitType = store.HabitTypeImprove
		}
		habitType = strings.ToLower(habitType)
		if habitType != store.HabitTypeImprove && habitType != store.HabitTypeQuit {
			return &se.StreakrError{TerminalMsg: fmt.Sprintf("type must be either '%s' or '%s'", store.HabitTypeImprove, store.HabitTypeQuit)}
		}

		return service.AddHabit(cmd.Context(), habitName, description, habitType)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.InitDefaultHelpFlag()
	addCmd.Flags().Lookup("help").Shorthand = ""
	addCmd.PersistentFlags().StringP("description", "d", "", "description of the habit")
	addCmd.PersistentFlags().StringP("type", "t", "", fmt.Sprintf("type of the habit (%s, %s) defaults to %s if unspecified", store.HabitTypeImprove, store.HabitTypeQuit, store.HabitTypeImprove))
}
