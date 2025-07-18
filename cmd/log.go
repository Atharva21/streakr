package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Atharva21/streakr/internal/service"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log today's habit completion",
	Long: `Log today's habit completion to track your streak.
We can also log multiple habits at once seperated by ,
Examples:
 streakr log run
 streakr log read,run,gym,youtube
 
This updates your current streak.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
		}

		habitNames := strings.Split(args[0], ",")
		if len(args) > 1 {
			habitNames = append(habitNames, args[1:]...)
		}
		if len(habitNames) == 0 {
			return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
		}
		for i, habitName := range habitNames {
			habitNames[i] = strings.ToLower(strings.TrimSpace(habitName))
			if habitNames[i] == "" {
				return &se.StreakrError{TerminalMsg: "habit name cannot be empty"}
			}
			if len(habitNames[i]) > 20 {
				return &se.StreakrError{TerminalMsg: "habit name cannot be > 20 chars"}
			}
		}
		allQuittingHabits, err := service.LogHabitsForToday(cmd.Context(), habitNames)
		if err != nil {
			return err
		}
		loggedHabitCount, totalHabitCount, err := service.GetTodaysLoggedHabitCount(cmd.Context())
		if err != nil {
			slog.Error(err.Error())
			return nil
		}
		if allQuittingHabits {
			// this is when all the quitting habits are logged.
			fmt.Fprintln(
				os.Stdout,
				"✔️  logged",
			)
			return nil
		}
		fmt.Fprintf(
			os.Stdout,
			"✔️  logged %d/%d today\n",
			loggedHabitCount,
			totalHabitCount,
		)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.InitDefaultHelpFlag()
	logCmd.Flags().Lookup("help").Shorthand = ""
}
