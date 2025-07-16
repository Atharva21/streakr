package cmd

import (
	"strings"

	"github.com/Atharva21/streakr/internal/service"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new habit to track",
	Long: `Add will add new habit and start tracking the streaks.
Add cmd lets you specify description and the type of habit

Few examples:
streakr add run --description "morning run 5kms"
streakr add drink --description "drink 3l of water daily to stay hydrated"
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
			habitType = "improve"
		}
		habitType = strings.ToLower(habitType)
		if habitType != "improve" && habitType != "quit" {
			return &se.StreakrError{TerminalMsg: "type must be either 'improve' or 'quit'"}
		}

		return service.AddHabit(cmd.Context(), habitName, description, habitType)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.InitDefaultHelpFlag()
	addCmd.Flags().Lookup("help").Shorthand = ""
	addCmd.PersistentFlags().StringP("description", "d", "", "description of the habit")
	addCmd.PersistentFlags().StringP("type", "t", "", "type of the habit (improve, quit) defaults to improve if unspecified")
}
