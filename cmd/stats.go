package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Atharva21/streakr/internal/service"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/Atharva21/streakr/internal/tui"
	"github.com/Atharva21/streakr/internal/util"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Stats gives current streak, max streak",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// TODO overall stats logic here.
			return nil
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

		yearStr, _ := cmd.Flags().GetString("year")
		monthStr, _ := cmd.Flags().GetString("month")

		// Validate and convert year
		currentYear := time.Now().Year()
		year := currentYear
		var err error
		if yearStr != "" {
			year, err = strconv.Atoi(yearStr)
			if err != nil {
				return &se.StreakrError{TerminalMsg: fmt.Sprintf("invalid year '%s': must be a valid integer", yearStr)}
			}
		}

		// Validate year range (reasonable bounds)
		if year < 1900 || year > currentYear {
			return &se.StreakrError{TerminalMsg: fmt.Sprintf("invalid year %d: must be between 1900 and %d", year, currentYear)}
		}

		// Validate and convert month
		if monthStr == "" {
			monthStr = time.Now().Month().String()
		}
		var month time.Month

		// Try parsing as number first
		if monthNum, err := strconv.Atoi(monthStr); err == nil {
			if monthNum < 1 || monthNum > 12 {
				return &se.StreakrError{TerminalMsg: fmt.Sprintf("invalid month %d: must be between 1-12", monthNum)}
			}
			month = time.Month(monthNum)
		} else {
			// Try parsing as month name
			switch strings.ToLower(monthStr) {
			case "january", "jan":
				month = time.January
			case "february", "feb":
				month = time.February
			case "march", "mar":
				month = time.March
			case "april", "apr":
				month = time.April
			case "may":
				month = time.May
			case "june", "jun":
				month = time.June
			case "july", "jul":
				month = time.July
			case "august", "aug":
				month = time.August
			case "september", "sep":
				month = time.September
			case "october", "oct":
				month = time.October
			case "november", "nov":
				month = time.November
			case "december", "dec":
				month = time.December
			default:
				return &se.StreakrError{TerminalMsg: fmt.Sprintf("invalid month '%s': must be 1-12 or month name", monthStr)}
			}
		}
		startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)
		habit, err := service.GetHabitByName(cmd.Context(), habitName)
		if err != nil {
			return err
		}
		startRange, endRange := habit.CreatedAt, time.Now()
		if err != nil {
			return err
		}
		// fmt.Printf("ms: %v, me: %v, rs: %v, re: %v\n", startOfMonth, endOfMonth, startRange, endRange)
		if startOfMonth.Year() < startRange.Year() {
			return &se.StreakrError{TerminalMsg: "Cannot get stats before habit creation date."}
		}
		if startOfMonth.Year() == startRange.Year() {
			if startOfMonth.Month() < startRange.Month() {
				return &se.StreakrError{TerminalMsg: "Cannot get stats before habit creation date."}
			}
		}
		if endOfMonth.Year() > endRange.Year() {
			return &se.StreakrError{TerminalMsg: "Cannot get stats after latest streak"}
		}
		if endOfMonth.Year() == endRange.Year() {
			if endOfMonth.Month() > endRange.Month() {
				return &se.StreakrError{TerminalMsg: "Cannot get stats after latest streak"}
			}
		}
		rangedStreaks, err := service.GetHabitStatsForRange(cmd.Context(), habitName, startOfMonth, endOfMonth)
		if err != nil {
			return err
		}
		sm := tui.StatsModel{
			Ctx:                 cmd.Context(),
			Habit:               rangedStreaks.Habit,
			HeatMap:             rangedStreaks.Heatmap,
			TotalStreaksInMonth: rangedStreaks.TotalStreakDaysInRange,
			TotalMissesInMonth:  rangedStreaks.TotalMissesInRange,
			Today:               time.Now(),
			FirstDayOfSetMonth:  startOfMonth,
			ExitError:           nil,
			HasPreviousNbr:      util.AtLeastOneMonthOlder(habit.CreatedAt, startOfMonth),
			HasNxtNbr:           util.AtLeastOneMonthOlder(startOfMonth, time.Now()),
		}
		err = tui.RenderStatsView(&sm)
		return err
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.InitDefaultHelpFlag()
	statsCmd.Flags().Lookup("help").Shorthand = ""
	statsCmd.PersistentFlags().StringP("month", "m", "", "Specify the month to view stats of")
	statsCmd.PersistentFlags().StringP("year", "y", "", "Specify the year to view stats of")
}
