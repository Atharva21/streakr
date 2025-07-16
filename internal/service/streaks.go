package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Atharva21/streakr/internal/store"
	"github.com/Atharva21/streakr/internal/store/generated"
	"github.com/Atharva21/streakr/internal/util"
)

func LogHabitForToday(appContext context.Context, habitNames []string) error {
	habitsToLogToday := make([]generated.Habit, 0)
	for _, habitName := range habitNames {
		habit, err := getHabitByName(appContext, habitName)
		if err != nil {
			return err
		}
		habitsToLogToday = append(habitsToLogToday, habit)
	}
	for _, habit := range habitsToLogToday {
		latestStreak, err := store.GetQueries().GetLatestStreakForHabit(appContext, habit.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			// first time log. for improvement habits, streak_start and end should be the same day.
			// for quitting habits, streak_start should be habit.Created_at + 1, and streak_end should be y'day
			// (with y'day being >= created_at + 1)
			today := time.Now()
			yesterday := today.AddDate(0, 0, -1)
			if habit.HabitType == store.HabitTypeImprove {
				if util.IsSameDate(yesterday, latestStreak.StreakEnd) {
					return store.GetQueries().UpdateStreakEnd(appContext, generated.UpdateStreakEndParams{
						ID:        latestStreak.ID,
						StreakEnd: today,
					})
				}
				_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
					HabitID:     habit.ID,
					StreakStart: today,
					StreakEnd:   today,
				})
				return err
			}
			// logic for first log quitting habits here
		} else {
			// logic for subsequent logs both improving and quitting
		}
	}
	return nil
}
