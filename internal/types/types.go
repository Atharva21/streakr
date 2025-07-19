package types

import (
	"time"

	"github.com/Atharva21/streakr/internal/store/generated"
)

type HabitInfo struct {
	Habit         generated.Habit
	CurrentStreak int64
	MaxStreak     int64
	LastLogged    time.Time
	LastMiss      time.Time
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
