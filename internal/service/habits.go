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

func GetHabitByName(appContext context.Context, name string) (generated.Habit, error) {
	habit, err := store.GetQueries().GetHabitByName(appContext, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return habit, &se.StreakrError{TerminalMsg: fmt.Sprintf("No habit with name %s", name)}
		}
		return habit, err
	}
	return habit, err
}

func AddHabit(appContext context.Context, name, description, habitType string) error {

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

	slog.Info("habit saved successfully", slog.Int64("habitId", habitId))

	return nil
}

func DeleteHabits(appContext context.Context, queries []string) error {
	habitsIDsToDelete := make([]int64, 0)
	for _, query := range queries {
		habit, err := GetHabitByName(appContext, query)
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
