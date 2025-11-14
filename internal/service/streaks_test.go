package service

import (
	"context"
	"testing"
	"time"

	"github.com/Atharva21/streakr/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogHabitsForToday_ImproveHabit_FirstLog(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()
	today := time.Now()

	tests := []struct {
		name          string
		createdAt     time.Time
		expectStreaks int
	}{
		{
			name:          "first log for habit created today",
			createdAt:     today,
			expectStreaks: 1,
		},
		{
			name:          "first log for habit created yesterday",
			createdAt:     today.AddDate(0, 0, -1),
			expectStreaks: 1,
		},
		{
			name:          "first log for habit created last week",
			createdAt:     today.AddDate(0, 0, -7),
			expectStreaks: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB.Cleanup()
			testDB = SetupTestDB(t)

			habit := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, &tt.createdAt)

			allQuitting, err := LogHabitsForToday(ctx, []string{"running"})
			require.NoError(t, err)
			assert.False(t, allQuitting)

			// Verify streak was created
			var count int
			err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
			require.NoError(t, err)
			assert.Equal(t, tt.expectStreaks, count)
		})
	}
}

func TestLogHabitsForToday_ImproveHabit_ConsecutiveDays(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	habit := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)

	// Log yesterday
	testDB.CreateTestStreak(t, ctx, habit.ID, yesterday, yesterday)

	// Log today (should extend streak)
	allQuitting, err := LogHabitsForToday(ctx, []string{"running"})
	require.NoError(t, err)
	assert.False(t, allQuitting)

	// Verify streak was extended, not created new
	var count int
	err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify streak end was updated
	var streakEnd time.Time
	err = testDB.DB.QueryRowContext(ctx, "SELECT streak_end FROM streaks WHERE habit_id = ?", habit.ID).Scan(&streakEnd)
	require.NoError(t, err)
	assert.True(t, isSameDay(streakEnd, today))
}

func TestLogHabitsForToday_ImproveHabit_MissedDays(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()
	today := time.Now()
	threeDaysAgo := today.AddDate(0, 0, -3)

	habit := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)

	// Log 3 days ago
	testDB.CreateTestStreak(t, ctx, habit.ID, threeDaysAgo, threeDaysAgo)

	// Log today (should create new streak)
	allQuitting, err := LogHabitsForToday(ctx, []string{"running"})
	require.NoError(t, err)
	assert.False(t, allQuitting)

	// Verify new streak was created
	var count int
	err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestLogHabitsForToday_QuitHabit_FirstLog(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	tests := []struct {
		name      string
		createdAt time.Time
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "first log for quit habit created today",
			createdAt: today,
			wantStart: today,
			wantEnd:   today,
		},
		{
			name:      "first log for quit habit created yesterday",
			createdAt: yesterday,
			wantStart: today,
			wantEnd:   today,
		},
		{
			name:      "first log for quit habit created week ago",
			createdAt: today.AddDate(0, 0, -7),
			wantStart: today.AddDate(0, 0, -6), // day after creation
			wantEnd:   today,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB.Cleanup()
			testDB = SetupTestDB(t)

			habit := testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, &tt.createdAt)

			allQuitting, err := LogHabitsForToday(ctx, []string{"smoking"})
			require.NoError(t, err)
			assert.True(t, allQuitting)

			// Verify streak was created
			var streakStart, streakEnd time.Time
			err = testDB.DB.QueryRowContext(ctx,
				"SELECT streak_start, streak_end FROM streaks WHERE habit_id = ?",
				habit.ID).Scan(&streakStart, &streakEnd)
			require.NoError(t, err)

			assert.True(t, isSameDay(streakStart, tt.wantStart), "Expected start %v, got %v", tt.wantStart, streakStart)
			assert.True(t, isSameDay(streakEnd, tt.wantEnd), "Expected end %v, got %v", tt.wantEnd, streakEnd)
		})
	}
}

func TestLogHabitsForToday_QuitHabit_SubsequentLogs(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	twoDaysAgo := today.AddDate(0, 0, -2)

	habit := testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, nil)

	t.Run("log after one clean day", func(t *testing.T) {
		// Log yesterday (first slip-up)
		testDB.CreateTestStreak(t, ctx, habit.ID, yesterday, yesterday)

		// Log today (second slip-up)
		allQuitting, err := LogHabitsForToday(ctx, []string{"smoking"})
		require.NoError(t, err)
		assert.True(t, allQuitting)

		// Should create new streak for today only
		var count int
		err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("log after multiple clean days", func(t *testing.T) {
		testDB.Cleanup()
		testDB = SetupTestDB(t)

		habit = testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, nil)

		// Log 2 days ago
		testDB.CreateTestStreak(t, ctx, habit.ID, twoDaysAgo, twoDaysAgo)

		// Log today (should create streak from day after last slip-up to today)
		allQuitting, err := LogHabitsForToday(ctx, []string{"smoking"})
		require.NoError(t, err)
		assert.True(t, allQuitting)

		// Should have created new streak
		var streakStart, streakEnd time.Time
		var count int
		err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count)

		// Get the latest streak
		err = testDB.DB.QueryRowContext(ctx,
			"SELECT streak_start, streak_end FROM streaks WHERE habit_id = ? ORDER BY streak_end DESC LIMIT 1",
			habit.ID).Scan(&streakStart, &streakEnd)
		require.NoError(t, err)

		assert.True(t, isSameDay(streakStart, yesterday)) // day after last slip-up
		assert.True(t, isSameDay(streakEnd, today))
	})
}

func TestLogHabitsForToday_DuplicateLog(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()
	today := time.Now()

	habit := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)

	// Log today
	testDB.CreateTestStreak(t, ctx, habit.ID, today, today)

	// Try to log again today (should be skipped)
	_, err := LogHabitsForToday(ctx, []string{"running"})
	require.NoError(t, err)

	// Verify only one streak exists
	var count int
	err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestLogHabitsForToday_MultipleHabits(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	habit1 := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)
	habit2 := testDB.CreateTestHabit(t, ctx, "reading", "test", store.HabitTypeImprove, nil)
	habit3 := testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, nil)

	allQuitting, err := LogHabitsForToday(ctx, []string{"running", "reading", "smoking"})
	require.NoError(t, err)
	assert.False(t, allQuitting) // Not all are quitting habits

	// Verify streaks created for all
	for _, habitID := range []int64{habit1.ID, habit2.ID, habit3.ID} {
		var count int
		err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habitID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	}
}

func TestLogHabitsForToday_NonExistentHabit(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	_, err := LogHabitsForToday(ctx, []string{"nonexistent"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "No habit with name")
}

func TestGetOverallStats(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	t.Run("no habits", func(t *testing.T) {
		stats, err := GetOverallStats(ctx)
		require.NoError(t, err)
		assert.Empty(t, stats.HabitInfos)
	})

	t.Run("multiple habits with streaks", func(t *testing.T) {
		today := time.Now()
		yesterday := today.AddDate(0, 0, -1)

		// Create improve habit with streak
		habit1 := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)
		testDB.CreateTestStreak(t, ctx, habit1.ID, yesterday, today)

		// Create quit habit with no logs
		testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, &yesterday)

		stats, err := GetOverallStats(ctx)
		require.NoError(t, err)
		assert.Len(t, stats.HabitInfos, 2)

		// Check improve habit stats
		var runningInfo, smokingInfo *struct {
			CurrentStreak      int64
			MaxStreak          int64
			TotalPerformedDays int64
		}

		for i := range stats.HabitInfos {
			if stats.HabitInfos[i].Habit.Name == "running" {
				runningInfo = &struct {
					CurrentStreak      int64
					MaxStreak          int64
					TotalPerformedDays int64
				}{
					CurrentStreak:      stats.HabitInfos[i].CurrentStreak,
					MaxStreak:          stats.HabitInfos[i].MaxStreak,
					TotalPerformedDays: stats.HabitInfos[i].TotalPerformedDays,
				}
			} else if stats.HabitInfos[i].Habit.Name == "smoking" {
				smokingInfo = &struct {
					CurrentStreak      int64
					MaxStreak          int64
					TotalPerformedDays int64
				}{
					CurrentStreak:      stats.HabitInfos[i].CurrentStreak,
					MaxStreak:          stats.HabitInfos[i].MaxStreak,
					TotalPerformedDays: stats.HabitInfos[i].TotalPerformedDays,
				}
			}
		}

		require.NotNil(t, runningInfo)
		assert.Equal(t, int64(2), runningInfo.CurrentStreak) // yesterday and today
		assert.Equal(t, int64(2), runningInfo.MaxStreak)
		assert.Equal(t, int64(2), runningInfo.TotalPerformedDays)

		require.NotNil(t, smokingInfo)
		assert.Equal(t, int64(0), smokingInfo.CurrentStreak) // today hasn't passed yet
	})
}

func TestGetHabitStatsForRange_ImproveHabit(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	createdAt := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	habit := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, &createdAt)

	// Create streaks for Nov 1-5 and Nov 10-12
	testDB.CreateTestStreak(t, ctx, habit.ID,
		time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local),
		time.Date(2025, 11, 5, 0, 0, 0, 0, time.Local))

	testDB.CreateTestStreak(t, ctx, habit.ID,
		time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local),
		time.Date(2025, 11, 12, 0, 0, 0, 0, time.Local))

	startDate := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2025, 11, 30, 0, 0, 0, 0, time.Local)

	stats, err := GetHabitStatsForRange(ctx, "running", startDate, endDate)
	require.NoError(t, err)

	assert.Equal(t, 8, stats.TotalStreakDaysInRange) // 5 days + 3 days
	assert.Len(t, stats.Heatmap, 30)                 // November has 30 days

	// Verify heatmap has correct days marked
	assert.True(t, stats.Heatmap[0])  // Nov 1
	assert.True(t, stats.Heatmap[4])  // Nov 5
	assert.False(t, stats.Heatmap[5]) // Nov 6 (missed)
	assert.True(t, stats.Heatmap[9])  // Nov 10
}

func TestGetHabitStatsForRange_QuitHabit_NoLogs(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	// Create quit habit on Nov 14
	createdAt := time.Date(2025, 11, 14, 0, 0, 0, 0, time.Local)
	testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, &createdAt)

	startDate := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2025, 11, 30, 0, 0, 0, 0, time.Local)

	stats, err := GetHabitStatsForRange(ctx, "smoking", startDate, endDate)
	require.NoError(t, err)

	// With no logs and created on Nov 14, all days from Nov 14 to yesterday should be clean
	// Total depends on current date, but missed should be 0
	assert.Equal(t, 0, stats.TotalMissesInRange)
	assert.True(t, stats.TotalStreakDaysInRange >= 0)
}

func TestGetHabitStatsForRange_QuitHabit_WithSlipups(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	createdAt := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	habit := testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, &createdAt)

	// Log slip-up on Nov 5 (clean days were Nov 2-4, slip-up on Nov 5)
	testDB.CreateTestStreak(t, ctx, habit.ID,
		time.Date(2025, 11, 2, 0, 0, 0, 0, time.Local),
		time.Date(2025, 11, 5, 0, 0, 0, 0, time.Local))

	// Log slip-up on Nov 10 (clean days were Nov 6-9, slip-up on Nov 10)
	testDB.CreateTestStreak(t, ctx, habit.ID,
		time.Date(2025, 11, 6, 0, 0, 0, 0, time.Local),
		time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local))

	startDate := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local)

	stats, err := GetHabitStatsForRange(ctx, "smoking", startDate, endDate)
	require.NoError(t, err)

	// Clean days: Nov 2,3,4 (3) + Nov 6,7,8,9 (4) = 7 days
	// Slip-ups: Nov 5, Nov 10 = 2 days
	// Nov 1 doesn't count (day of creation)
	assert.Equal(t, 7, stats.TotalStreakDaysInRange)
	assert.Equal(t, 2, stats.TotalMissesInRange)
}

func TestGetHabitStatsForRange_BeforeCreation(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	// Create habit on Nov 15
	createdAt := time.Date(2025, 11, 15, 0, 0, 0, 0, time.Local)
	testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, &createdAt)

	// Query for Nov 1-30 (habit doesn't exist for first 14 days)
	startDate := time.Date(2025, 11, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2025, 11, 30, 0, 0, 0, 0, time.Local)

	stats, err := GetHabitStatsForRange(ctx, "running", startDate, endDate)
	require.NoError(t, err)

	// Should only count days from Nov 15 onwards
	// Missed count should not include Nov 1-14
	assert.True(t, stats.TotalMissesInRange <= 16) // Max 16 days from Nov 15 to Nov 30 (excluding today if future)
}

// Helper function to check if two times are on the same day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
