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

func LogHabitsForToday(appContext context.Context, habitNames []string) error {
	habitsToLogToday := make([]generated.Habit, 0)
	for _, habitName := range habitNames {
		habit, err := getHabitByName(appContext, habitName)
		if err != nil {
			return err
		}
		habitsToLogToday = append(habitsToLogToday, habit)
	}
	today := time.Now()
	yesterday := util.GetPrevDayOf(today)
	for _, habit := range habitsToLogToday {
		latestStreak, err := store.GetQueries().GetLatestStreakForHabit(appContext, habit.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			// first time log. for improvement habits, streak_start and end should be the same day.
			// for quitting habits, streak_start should be habit.Created_at + 1, and streak_end should be y'day
			// (with y'day being >= created_at + 1)
			if habit.HabitType == store.HabitTypeImprove {
				_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
					HabitID:     habit.ID,
					StreakStart: today,
					StreakEnd:   today,
				})
				return err
			}
			// logic for first log of quitting habits here
			if util.IsSameDate(habit.CreatedAt, today) || util.IsSameDate(habit.CreatedAt, yesterday) {
				_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
					HabitID:     habit.ID,
					StreakStart: today,
					StreakEnd:   today,
				})
				return err
			}
			_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
				HabitID:     habit.ID,
				StreakStart: util.GetNextDayOf(habit.CreatedAt),
				StreakEnd:   today,
			})
			return err
		} else {
			// logic for subsequent logs both improvement and quitting

			if util.IsSameDate(latestStreak.StreakEnd, today) {
				// handle duplicate logging. (skip for now)
				continue
			}

			if habit.HabitType == store.HabitTypeImprove {
				// for improvement habits, if latest streak is of y'day update it to today. else add new streak
				if util.IsSameDate(latestStreak.StreakEnd, yesterday) {
					err = store.GetQueries().UpdateStreakEnd(appContext, generated.UpdateStreakEndParams{
						ID:        latestStreak.ID,
						StreakEnd: today,
					})
					if err != nil {
						return err
					}
				} else {
					_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
						HabitID:     habit.ID,
						StreakStart: today,
						StreakEnd:   today,
					})
					return err
				}
			} else {
				// for quitting habits, if latest.end == y'day, we do a today->today log
				// else we log from latest.end+1->today
				if util.IsSameDate(latestStreak.StreakEnd, yesterday) {
					_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
						HabitID:     habit.ID,
						StreakStart: today,
						StreakEnd:   today,
					})
					return err
				} else {
					_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
						HabitID:     habit.ID,
						StreakStart: util.GetNextDayOf(latestStreak.StreakEnd),
						StreakEnd:   today,
					})
					return err
				}
			}

		}
	}
	return nil
}

func GetStatsForHabitName(appContext context.Context, habitName string) (int64, int64, error) {
	habit, err := getHabitByName(appContext, habitName)
	if err != nil {
		return 0, 0, err
	}
	currentStreak, err := getCurrentStreakForHabit(appContext, habit)
	if err != nil {
		return 0, 0, err
	}
	pastMaxStreak, err := getPastMaxStreakForHabit(appContext, habit)
	if err != nil {
		return 0, 0, err
	}
	if int64(currentStreak) >= pastMaxStreak {
		pastMaxStreak = int64(currentStreak)
	}
	return int64(currentStreak), pastMaxStreak, nil
}

func getCurrentStreakForHabit(appContext context.Context, habit generated.Habit) (int, error) {
	today := time.Now()
	yesterday := util.GetPrevDayOf(today)
	latestStreak, err := store.GetQueries().GetLatestStreakForHabit(appContext, habit.ID)
	if err == nil {
		if habit.HabitType == store.HabitTypeImprove {
			// for improvement habits latest streak is whatever is going on (if its y'day) else 0.
			if util.IsSameDate(latestStreak.StreakEnd, today) || util.IsSameDate(latestStreak.StreakEnd, yesterday) {
				daysDiff := util.GetDayDiff(latestStreak.StreakStart, latestStreak.StreakEnd)
				return daysDiff + 1, nil
			} else {
				// no latest logs
				return 0, nil
			}
		}
		// for quitting habits,
		daysDiff := util.GetDayDiff(latestStreak.StreakEnd, today)
		if daysDiff-1 >= 0 {
			daysDiff--
		}
		return daysDiff, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	// no past logs, first occurance.
	if habit.HabitType == store.HabitTypeImprove {
		return 0, nil
	}
	if util.IsSameDate(today, habit.CreatedAt) {
		return 0, nil
	}
	cleanDays := util.GetDayDiff(habit.CreatedAt, today)
	if cleanDays-1 >= 0 {
		cleanDays--
	}
	return cleanDays, nil
}

func getPastMaxStreakForHabit(appContext context.Context, habit generated.Habit) (int64, error) {
	if habit.HabitType == store.HabitTypeImprove {
		return store.GetQueries().GetMaxStreakForHabit(appContext, habit.ID)
	}
	return store.GetQueries().GetMaxStreakQuittingHabit(appContext, habit.ID)
}
