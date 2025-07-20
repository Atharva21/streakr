package cmd

import (
	"github.com/Atharva21/streakr/internal/tui"
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
		return tui.RenderListView(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.InitDefaultHelpFlag()
	addCmd.Flags().Lookup("help").Shorthand = ""
}
