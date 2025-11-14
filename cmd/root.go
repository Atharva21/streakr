package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Atharva21/streakr/internal/shutdown"
	"github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/Atharva21/streakr/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "streakr",
	Short:         "streakr is a habit tracking CLI",
	Long:          `streakr is a command-line tool for tracking habits and maintaining streaks...`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err == nil {
		return
	}
	var streakrErr *streakrerror.StreakrError
	if errors.As(err, &streakrErr) {
		if streakrErr.ShowUsage {
			rootCmd.Usage()
		}
		if streakrErr.TerminalMsg != "" {
			util.ErrorAndExit(streakrErr.TerminalMsg)
		} else {
			util.ErrorAndExitGeneric(streakrErr)
		}
	} else {
		// This is a Cobra error - manually print it
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		rootCmd.Usage()
		shutdown.GracefulShutdown(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.InitDefaultHelpFlag()
	addCmd.Flags().Lookup("help").Shorthand = ""
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("streakr v{{.Version}}\n")
}
