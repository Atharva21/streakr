package service

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/Atharva21/streakr/internal/store"
	"github.com/Atharva21/streakr/internal/store/generated"
	se "github.com/Atharva21/streakr/internal/streakrerror"
)

func AddHabit(appContext context.Context, name, description, habitType string, aliases []string) error {
	_, err := store.GetQueries().GetHabitByName(appContext, name)
	if err == nil {
		return &se.StreakrError{
			TerminalMsg: "Cannot add another habit with same name, to list all habits: streakr list",
			ShowUsage:   true,
		}
	}
	if err != sql.ErrNoRows {
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
				ShowUsage:   true,
			}
		}
		if err != sql.ErrNoRows {
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
