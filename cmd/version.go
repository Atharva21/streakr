package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set via ldflags during build
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of streakr",
	Long:  `All software has versions. This is streakr's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("streakr v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
