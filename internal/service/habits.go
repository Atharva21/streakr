package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Atharva21/streakr/internal/store"
	"github.com/Atharva21/streakr/internal/store/generated"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/mattn/go-sqlite3"
)

func getHabitByName(appContext context.Context, name string) (generated.Habit, error) {
	habit, err := store.GetQueries().GetHabitByName(appContext, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return habit, &se.StreakrError{TerminalMsg: fmt.Sprintf("No habit with name %s", name)}
		}
		return habit, err
	}
	return habit, err
}

func getHabitbyAlias(appContext context.Context, name string) (generated.Habit, error) {
	habit, err := store.GetQueries().GetHabitByAlias(appContext, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return habit, &se.StreakrError{TerminalMsg: fmt.Sprintf("No habit with alias %s", name)}
		}
		return habit, err
	}
	return habit, err
}

func getHabitByNameOrAlias(appContext context.Context, query string) (generated.Habit, error) {
	habit, err := getHabitByName(appContext, query)
	if err == nil {
		return habit, err
	}
	var streakrErr *se.StreakrError
	if !errors.As(err, &streakrErr) {
		return habit, err
	}
	habit, err = getHabitbyAlias(appContext, query)
	if err != nil && errors.As(err, &streakrErr) {
		err = &se.StreakrError{TerminalMsg: fmt.Sprintf("No habit with name or alias %s", query)}
	}
	return habit, err
}

func AddHabit(appContext context.Context, name, description, habitType string, aliases []string) error {

	// validate aliases are unique.
	for _, alias := range aliases {
		_, err := getHabitbyAlias(appContext, alias)
		if err == nil {
			return &se.StreakrError{TerminalMsg: fmt.Sprintf("Alias with name %s already exists", alias)}
		}
	}

	habitId, err := store.GetQueries().AddHabit(
		appContext,
		generated.AddHabitParams{
			Name: name,
			Description: sql.NullString{
				String: description,
				Valid:  description != "",
			},
			HabitType: habitType,
		},
	)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique || sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return &se.StreakrError{TerminalMsg: fmt.Sprintf("Cannot add habit with name %s as it already exists", name)}
			}
		}
		return err
	}

	for _, alias := range aliases {
		err = store.GetQueries().AddAliasForHabit(
			appContext,
			generated.AddAliasForHabitParams{
				Alias:   alias,
				HabitID: habitId,
			},
		)
		if err != nil {
			return err
		}
	}

	slog.Info("habit saved successfully", slog.Int64("habitId", habitId))

	return nil
}

func DeleteHabits(appContext context.Context, queries []string) error {
	habitsIDsToDelete := make([]int64, 0)
	for _, query := range queries {
		habit, err := getHabitByNameOrAlias(appContext, query)
		if err != nil {
			return err
		}
		habitsIDsToDelete = append(habitsIDsToDelete, habit.ID)
	}
	for _, habitID := range habitsIDsToDelete {
		err := store.GetQueries().DeleteHabit(appContext, habitID)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListHabits(appContext context.Context) ([]generated.Habit, error) {
	return store.GetQueries().ListHabits(appContext)
}
