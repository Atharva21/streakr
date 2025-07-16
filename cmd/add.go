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
Add cmd lets you specify aliases, description and the type of habit

Few examples:
streakr add morning 3km run --alias run
streakr add drink 3l water --alias water,h2o --description drink 3l of water daily to stay hydrated
streakr add smoking --alias smoke --type quit
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse habit name from args
		if len(args) == 0 {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty", ShowUsage: true}
		}

		// Clean and join habit name
		habitName := strings.TrimSpace(strings.Join(args, " "))
		if habitName == "" {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty", ShowUsage: true}
		}
		if len(habitName) > 50 {
			return &se.StreakrError{TerminalMsg: "habit name cannot exceed 50 characters", ShowUsage: true}
		}

		// Parse flags
		aliasStr, _ := cmd.Flags().GetString("alias")
		description, _ := cmd.Flags().GetString("description")
		habitType, _ := cmd.Flags().GetString("type")

		if len(description) > 200 {
			return &se.StreakrError{TerminalMsg: "description cannot exceed 200 characters", ShowUsage: true}
		}

		// Validate and default type
		if habitType == "" {
			habitType = "improve"
		}
		habitType = strings.ToLower(habitType)
		if habitType != "improve" && habitType != "quit" {
			return &se.StreakrError{TerminalMsg: "type must be either 'improve' or 'quit'", ShowUsage: true}
		}

		// Parse aliases
		var aliases []string
		if aliasStr != "" {
			aliasMap := make(map[string]bool)
			aliases = strings.Split(aliasStr, ",")
			if len(aliases) > 5 {
				return &se.StreakrError{TerminalMsg: "cannot add more than 5 aliases at once", ShowUsage: true}
			}
			for i, alias := range aliases {
				aliases[i] = strings.ToLower(strings.TrimSpace(alias))
				if len(aliases[i]) > 15 {
					return &se.StreakrError{TerminalMsg: "habit aliases cannot exceed 15 characters", ShowUsage: true}
				}
				if _, exists := aliasMap[aliases[i]]; exists {
					return &se.StreakrError{TerminalMsg: "habit aliases must be unique", ShowUsage: true}
				}
				aliasMap[aliases[i]] = true
			}
		}

		return service.AddHabit(cmd.Context(), habitName, description, habitType, aliases)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.InitDefaultHelpFlag()
	addCmd.Flags().Lookup("help").Shorthand = ""
	addCmd.PersistentFlags().StringP("alias", "a", "", "alias for the habit also supports comma seperated aliases")
	addCmd.PersistentFlags().StringP("description", "d", "", "description of the habit")
	addCmd.PersistentFlags().StringP("type", "t", "", "type of the habit (improve, quit) defaults to improve if unspecified")
}
