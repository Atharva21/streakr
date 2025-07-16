package cmd

import (
	"context"
	"errors"

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
	if err != nil {
		var streakrErr *streakrerror.StreakrError
		if errors.As(err, &streakrErr) && streakrErr.TerminalMsg != "" {
			if streakrErr.ShowUsage {
				rootCmd.Usage()
			}
			util.ErrorAndExit(streakrErr.TerminalMsg)
		} else {
			util.ErrorAndExitGeneric(err)
		}
	}
}
