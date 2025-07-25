package types

import (
	"time"

	"github.com/Atharva21/streakr/internal/store/generated"
)

type HabitInfo struct {
	Habit              generated.Habit
	CurrentStreak      int64
	MaxStreak          int64
	TotalPerformedDays int64
	TotalMissedDays    int64
}

type HabitStatsForRange struct {
	Habit                  generated.Habit
	Heatmap                []bool
	TotalStreakDaysInRange int
	TotalMissesInRange     int
	RangeStart             time.Time
	RangeEnd               time.Time
}

type OverallStats struct {
	HabitInfos []HabitInfo
}
