package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Atharva21/streakr/internal/store"
	"github.com/Atharva21/streakr/internal/store/generated"
	"github.com/Atharva21/streakr/internal/types"
	"github.com/Atharva21/streakr/internal/util"
)

func LogHabitsForToday(appContext context.Context, habitNames []string) (bool, error) {
	habitsToLogToday := make([]generated.Habit, 0)
	allQuittingHabits := true
	for _, habitName := range habitNames {
		habit, err := GetHabitByName(appContext, habitName)
		if err != nil {
			return allQuittingHabits, err
		}
		if habit.HabitType == store.HabitTypeImprove {
			allQuittingHabits = false
		}
		habitsToLogToday = append(habitsToLogToday, habit)
	}
	today := time.Now()
	yesterday := util.GetPrevDayOf(today)
	for _, habit := range habitsToLogToday {
		latestStreak, err := store.GetQueries().GetLatestStreakForHabit(appContext, habit.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return allQuittingHabits, err
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
				if err != nil {
					return allQuittingHabits, err
				}
				continue
			}
			// logic for first log of quitting habits here
			if util.IsSameDate(habit.CreatedAt, today) || util.IsSameDate(habit.CreatedAt, yesterday) {
				_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
					HabitID:     habit.ID,
					StreakStart: today,
					StreakEnd:   today,
				})
				if err != nil {
					return allQuittingHabits, err
				}
				continue
			}
			_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
				HabitID:     habit.ID,
				StreakStart: util.GetNextDayOf(habit.CreatedAt),
				StreakEnd:   today,
			})
			if err != nil {
				return allQuittingHabits, err
			}
			continue
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
						return allQuittingHabits, err
					}
				} else {
					_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
						HabitID:     habit.ID,
						StreakStart: today,
						StreakEnd:   today,
					})
					if err != nil {
						return allQuittingHabits, err
					}
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
					if err != nil {
						return allQuittingHabits, err
					}
				} else {
					_, err = store.GetQueries().AddStreak(appContext, generated.AddStreakParams{
						HabitID:     habit.ID,
						StreakStart: util.GetNextDayOf(latestStreak.StreakEnd),
						StreakEnd:   today,
					})
					if err != nil {
						return allQuittingHabits, err
					}
				}
			}

		}
	}
	return allQuittingHabits, nil
}

func getHabitInfoForHabit(appContext context.Context, habit generated.Habit) (*types.HabitInfo, error) {
	c, err := getCurrentStreakForHabit(appContext, habit)
	currentStreak := int64(c)
	if err != nil {
		return nil, err
	}
	pastMaxStreak, err := getPastMaxStreakForHabit(appContext, habit)
	if err != nil {
		return nil, err
	}
	if currentStreak >= pastMaxStreak {
		pastMaxStreak = currentStreak
	}
	daysSinceHabitCreation, err := store.GetQueries().GetDaysSinceHabitCreation(appContext, habit.ID)
	if err != nil {
		return nil, err
	}
	var totalStreakDays int64
	if habit.HabitType == store.HabitTypeImprove {
		totalStreakDays, err = store.GetQueries().GetTotalStreakDays(appContext, habit.ID)
		if err != nil {
			return nil, err
		}
	} else {
		totalStreakDays, err = store.GetQueries().GetTotalStreakDaysQuittingHabit(appContext, habit.ID)
		if err != nil {
			return nil, err
		}
		totalStreakDays += currentStreak
	}
	totalMissedDays := daysSinceHabitCreation - totalStreakDays
	return &types.HabitInfo{
		Habit:              habit,
		CurrentStreak:      currentStreak,
		MaxStreak:          pastMaxStreak,
		TotalPerformedDays: totalStreakDays,
		TotalMissedDays:    totalMissedDays,
	}, nil
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
		// for quitting habits, latest streak end represents the last slip-up
		// current streak = clean days since last slip-up
		daysDiff := util.GetDayDiff(latestStreak.StreakEnd, today)
		// Subtract 1 because today hasn't passed yet and the slip-up day doesn't count as clean
		currentStreak := daysDiff - 1
		if currentStreak < 0 {
			currentStreak = 0
		}
		return currentStreak, nil
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

func GetHabitStatsForRange(appContext context.Context, habitName string, startDate time.Time, endDate time.Time) (*types.HabitStatsForRange, error) {
	habit, err := GetHabitByName(appContext, habitName)
	if err != nil {
		return nil, err
	}
	streaksLst, err := store.GetQueries().GetStreaksInRange(appContext, generated.GetStreaksInRangeParams{
		StreakEnd:   startDate,
		StreakStart: endDate,
		HabitID:     habit.ID,
	})
	if err != nil {
		return nil, err
	}
	heatmap := make([]bool, util.GetDayDiff(startDate, endDate)+1)
	totalStreakDaysInRange := 0

	// For quit habits with no logs, all days from creation onwards are clean (true in heatmap)
	if habit.HabitType == store.HabitTypeQuit && len(streaksLst) == 0 {
		// No logs yet - all days from creation up to yesterday are clean
		// Today is not counted as completed yet since the day hasn't passed
		today := time.Now()
		yesterday := util.GetPrevDayOf(today)
		effectiveStart := startDate
		if util.CompareDate(habit.CreatedAt, startDate) == -1 {
			effectiveStart = habit.CreatedAt
		}
		effectiveEnd := endDate
		// Don't count today or future dates as completed
		if util.CompareDate(yesterday, endDate) == 1 {
			effectiveEnd = yesterday
		}
		for date := effectiveStart; util.CompareDate(date, effectiveEnd) >= 0; date = date.AddDate(0, 0, 1) {
			idx := util.GetDayDiff(startDate, date)
			heatmap[idx] = true
			totalStreakDaysInRange++
		}
	} else {
		// Process streaks from database
		today := time.Now()
		for _, streak := range streaksLst {
			for date := streak.StreakStart; util.CompareDate(date, streak.StreakEnd) >= 0; date = date.AddDate(0, 0, 1) {
				if util.CompareDate(date, startDate) == 1 || util.CompareDate(date, endDate) == -1 {
					// skip partial dates that may be out of range. (< startDate or > endDate)
					continue
				}
				// now date is >= startdate && date <= endDate
				idx := util.GetDayDiff(startDate, date)
				if util.IsSameDate(date, streak.StreakEnd) && habit.HabitType == store.HabitTypeQuit {
					// for quitting habits last day is not considered as a streak, but is just a marker
					continue
				}
				// Skip today for improve habits unless it's been logged
				// For improve habits, today is included if logged (streak_end >= today means today is logged)
				if habit.HabitType == store.HabitTypeImprove && util.IsSameDate(date, today) {
					// Only count today if this streak includes today (i.e., it was logged)
					if util.CompareDate(streak.StreakEnd, today) < 0 {
						// This streak ends before today, don't count today
						continue
					}
				}
				heatmap[idx] = true
				totalStreakDaysInRange++
			}
		}
	}

	// Calculate total days to consider (only from habit creation date onwards)
	today := time.Now()
	yesterday := util.GetPrevDayOf(today)

	effectiveStartDate := startDate
	// For quit habits, the creation day doesn't count - start from day after
	// UNLESS there's a slip-up logged on the creation day itself
	habitStartDate := habit.CreatedAt
	if habit.HabitType == store.HabitTypeQuit {
		// Check if any streak has a slip-up on the creation day
		hasCreationDaySlipup := false
		for _, streak := range streaksLst {
			if util.IsSameDate(streak.StreakEnd, habit.CreatedAt) {
				hasCreationDaySlipup = true
				break
			}
		}
		// Only skip creation day if there's no slip-up on that day
		if !hasCreationDaySlipup {
			habitStartDate = util.GetNextDayOf(habit.CreatedAt)
		}
	}

	if util.CompareDate(habitStartDate, startDate) == -1 {
		// habit tracking started after the start of the range
		effectiveStartDate = habitStartDate
	}

	effectiveEndDate := endDate
	// For quit habits: if there are no logs (no slip-ups), count up to yesterday
	// If there are logs (slip-ups), count up to today to include today's slip-up in missed count
	if habit.HabitType == store.HabitTypeQuit && len(streaksLst) == 0 {
		// No slip-ups logged yet, don't count today since day hasn't passed
		if util.CompareDate(yesterday, endDate) == 1 {
			effectiveEndDate = yesterday
		}
	} else {
		// Has logs or is an improve habit, count up to today
		if util.CompareDate(today, endDate) == 1 {
			effectiveEndDate = today
		}
	}

	// If effectiveStartDate is after effectiveEndDate, no days to count
	totalDaysInRange := 0
	if util.CompareDate(effectiveStartDate, effectiveEndDate) >= 0 {
		totalDaysInRange = util.GetDayDiff(effectiveStartDate, effectiveEndDate) + 1
	}

	hs := &types.HabitStatsForRange{
		Habit:                  habit,
		Heatmap:                heatmap,
		TotalStreakDaysInRange: totalStreakDaysInRange,
		TotalMissesInRange:     totalDaysInRange - totalStreakDaysInRange,
		RangeStart:             startDate,
		RangeEnd:               endDate,
	}
	return hs, nil
}

func GetOverallStats(appContext context.Context) (*types.OverallStats, error) {
	habits, err := ListHabits(appContext)
	if err != nil {
		return nil, err
	}
	habitInfos := make([]types.HabitInfo, 0)
	for _, habit := range habits {
		habitInfo, err := getHabitInfoForHabit(appContext, habit)
		if err != nil {
			return nil, err
		}
		habitInfos = append(habitInfos, *habitInfo)
	}
	return &types.OverallStats{
		HabitInfos: habitInfos,
	}, nil
}
