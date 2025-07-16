package cmd

import (
	"fmt"

	"github.com/Atharva21/streakr/internal/service"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all habits being tracked",
	Long: `List all tracked habits
Example usage:

streakr list
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		habits, err := service.ListHabits(cmd.Context())
		if err != nil {
			return err
		}
		for _, habit := range habits {
			fmt.Println(habit.Name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.InitDefaultHelpFlag()
}
