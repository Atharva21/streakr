package cmd

import (
	"strings"

	"github.com/Atharva21/streakr/internal/service"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing habit and all its data",
	Long: `Delete an existing habit by its habit name can also specify , seperated queries
For example:
streakr delete morning walk
streakr delete gym,water,smoke,run
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
		}

		queries := strings.Split(args[0], ",")
		if len(args) > 1 {
			queries = append(queries, args[1:]...)
		}
		if len(queries) == 0 {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
		}
		for i, query := range queries {
			queries[i] = strings.TrimSpace(query)
			if query == "" {
				return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
			}
			if len(query) > 20 {
				return &se.StreakrError{TerminalMsg: "habit name cannot be > 20 chars"}
			}
		}

		return service.DeleteHabits(cmd.Context(), queries)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.InitDefaultHelpFlag()
	addCmd.Flags().Lookup("help").Shorthand = ""
}
