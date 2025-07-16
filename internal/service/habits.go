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
)

func isHabitWithNameOrAliasPresent(appContext context.Context, query string) (generated.Habit, error) {
	habit, err := store.GetQueries().GetHabitByName(appContext, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			habit, err = store.GetQueries().GetHabitByAlias(appContext, query)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return habit, &se.StreakrError{TerminalMsg: fmt.Sprintf("No habit with name or alias %s", query)}
				}
				return habit, err
			}
		} else {
			return habit, err
		}
	}
	return habit, err
}

func AddHabit(appContext context.Context, name, description, habitType string, aliases []string) error {
	_, err := store.GetQueries().GetHabitByName(appContext, name)
	if err == nil {
		return &se.StreakrError{
			TerminalMsg: "Cannot add another habit with same name, to list all habits: streakr list",
		}
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// validate aliases are unique.
	for _, alias := range aliases {
		_, err = store.GetQueries().GetHabitByAlias(
			appContext,
			alias,
		)
		if err == nil {
			return &se.StreakrError{
				TerminalMsg: "Cannot add another alias with same name, to get all aliases: streakr alias list",
			}
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
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
	habitsToDelete := make([]generated.Habit, 0)
	for _, query := range queries {
		habit, err := isHabitWithNameOrAliasPresent(appContext, query)
		if err != nil {
			return err
		}
		habitsToDelete = append(habitsToDelete, habit)
	}
	for _, habit := range habitsToDelete {
		err := store.GetQueries().DeleteHabit(appContext, habit.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func ListHabits(appContext context.Context) ([]generated.Habit, error) {
	return store.GetQueries().ListHabits(appContext)
}
