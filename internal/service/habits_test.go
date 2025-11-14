package service

import (
	"context"
	"testing"
	"time"

	"github.com/Atharva21/streakr/internal/store"
	se "github.com/Atharva21/streakr/internal/streakrerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddHabit(t *testing.T) {
	tests := []struct {
		name        string
		habitName   string
		description string
		habitType   string
		wantErr     bool
		errContains string
	}{
		{
			name:        "add improve habit successfully",
			habitName:   "running",
			description: "5k morning run",
			habitType:   store.HabitTypeImprove,
			wantErr:     false,
		},
		{
			name:        "add quit habit successfully",
			habitName:   "smoking",
			description: "quit smoking",
			habitType:   store.HabitTypeQuit,
			wantErr:     false,
		},
		{
			name:        "add habit without description",
			habitName:   "meditation",
			description: "",
			habitType:   store.HabitTypeImprove,
			wantErr:     false,
		},
		{
			name:        "add habit with empty name should fail",
			habitName:   "",
			description: "test",
			habitType:   store.HabitTypeImprove,
			wantErr:     true,
		},
		{
			name:        "add habit with name > 20 chars should fail",
			habitName:   "this_is_a_very_long_habit_name_that_exceeds_limit",
			description: "test",
			habitType:   store.HabitTypeImprove,
			wantErr:     true,
		},
		{
			name:        "add habit with description > 200 chars should fail",
			habitName:   "test",
			description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur",
			habitType:   store.HabitTypeImprove,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB := SetupTestDB(t)
			defer testDB.Cleanup()

			ctx := context.Background()
			err := AddHabit(ctx, tt.habitName, tt.description, tt.habitType)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				// Verify habit was actually created
				habit, err := GetHabitByName(ctx, tt.habitName)
				require.NoError(t, err)
				assert.Equal(t, tt.habitName, habit.Name)
				assert.Equal(t, tt.habitType, habit.HabitType)
				if tt.description != "" {
					assert.True(t, habit.Description.Valid)
					assert.Equal(t, tt.description, habit.Description.String)
				} else {
					assert.False(t, habit.Description.Valid)
				}
			}
		})
	}
}

func TestAddHabit_Duplicate(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	// Add first habit
	err := AddHabit(ctx, "running", "test", store.HabitTypeImprove)
	require.NoError(t, err)

	// Try to add duplicate
	err = AddHabit(ctx, "running", "another description", store.HabitTypeImprove)
	require.Error(t, err)

	var streakrErr *se.StreakrError
	assert.ErrorAs(t, err, &streakrErr)
	assert.Contains(t, streakrErr.TerminalMsg, "already exists")
}

func TestGetHabitByName(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	// Create test habit
	err := AddHabit(ctx, "running", "5k run", store.HabitTypeImprove)
	require.NoError(t, err)

	tests := []struct {
		name        string
		habitName   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "get existing habit",
			habitName: "running",
			wantErr:   false,
		},
		{
			name:        "get non-existent habit",
			habitName:   "swimming",
			wantErr:     true,
			errContains: "No habit with name",
		},
		{
			name:        "get habit with empty name",
			habitName:   "",
			wantErr:     true,
			errContains: "No habit with name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			habit, err := GetHabitByName(ctx, tt.habitName)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					var streakrErr *se.StreakrError
					if assert.ErrorAs(t, err, &streakrErr) {
						assert.Contains(t, streakrErr.TerminalMsg, tt.errContains)
					}
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.habitName, habit.Name)
			}
		})
	}
}

func TestListHabits(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	t.Run("list empty habits", func(t *testing.T) {
		habits, err := ListHabits(ctx)
		require.NoError(t, err)
		assert.Empty(t, habits)
	})

	t.Run("list multiple habits", func(t *testing.T) {
		// Create multiple habits
		err := AddHabit(ctx, "running", "5k run", store.HabitTypeImprove)
		require.NoError(t, err)

		err = AddHabit(ctx, "meditation", "10 min", store.HabitTypeImprove)
		require.NoError(t, err)

		err = AddHabit(ctx, "smoking", "quit", store.HabitTypeQuit)
		require.NoError(t, err)

		habits, err := ListHabits(ctx)
		require.NoError(t, err)
		assert.Len(t, habits, 3)

		// Verify habits are returned
		habitNames := make(map[string]bool)
		for _, h := range habits {
			habitNames[h.Name] = true
		}
		assert.True(t, habitNames["running"])
		assert.True(t, habitNames["meditation"])
		assert.True(t, habitNames["smoking"])
	})
}

func TestDeleteHabits(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	// Create test habits
	err := AddHabit(ctx, "running", "test", store.HabitTypeImprove)
	require.NoError(t, err)

	err = AddHabit(ctx, "meditation", "test", store.HabitTypeImprove)
	require.NoError(t, err)

	tests := []struct {
		name        string
		queries     []string
		wantErr     bool
		errContains string
		remaining   int
	}{
		{
			name:      "delete single habit",
			queries:   []string{"running"},
			wantErr:   false,
			remaining: 1,
		},
		{
			name:        "delete non-existent habit",
			queries:     []string{"swimming"},
			wantErr:     true,
			errContains: "No habit with name",
		},
		{
			name:      "delete multiple habits",
			queries:   []string{"running", "meditation"},
			wantErr:   false,
			remaining: 0,
		},
		{
			name:      "delete empty list",
			queries:   []string{},
			wantErr:   false,
			remaining: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset habits for each test
			testDB.Cleanup()
			testDB = SetupTestDB(t)

			err := AddHabit(ctx, "running", "test", store.HabitTypeImprove)
			require.NoError(t, err)
			err = AddHabit(ctx, "meditation", "test", store.HabitTypeImprove)
			require.NoError(t, err)

			err = DeleteHabits(ctx, tt.queries)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				// Verify remaining habits
				habits, err := ListHabits(ctx)
				require.NoError(t, err)
				assert.Len(t, habits, tt.remaining)
			}
		})
	}
}

func TestDeleteHabits_CascadeDelete(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	// Create habit
	habit := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)

	// Create streaks for habit
	testDB.CreateTestStreak(t, ctx, habit.ID, time.Date(2025, 11, 10, 0, 0, 0, 0, time.Local), time.Date(2025, 11, 12, 0, 0, 0, 0, time.Local))

	// Verify streak exists
	var count int
	err := testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Delete habit
	err = DeleteHabits(ctx, []string{"running"})
	require.NoError(t, err)

	// Verify streaks are also deleted (CASCADE)
	err = testDB.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM streaks WHERE habit_id = ?", habit.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestGetTodaysLoggedHabitCount(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()

	t.Run("no habits exist", func(t *testing.T) {
		completed, total, err := GetTodaysLoggedHabitCount(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), completed)
		assert.Equal(t, int64(0), total)
	})

	t.Run("habits exist but none logged", func(t *testing.T) {
		// Create improve habits
		testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)
		testDB.CreateTestHabit(t, ctx, "reading", "test", store.HabitTypeImprove, nil)

		completed, total, err := GetTodaysLoggedHabitCount(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), completed)
		assert.Equal(t, int64(2), total)
	})

	t.Run("quit habits don't count in totals", func(t *testing.T) {
		testDB.Cleanup()
		testDB = SetupTestDB(t)

		testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)
		testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, nil)

		completed, total, err := GetTodaysLoggedHabitCount(ctx)
		require.NoError(t, err)
		assert.Equal(t, int64(0), completed)
		assert.Equal(t, int64(1), total) // Only improve habits count
	})
}

// Helper function to create time.Time from date
func timeDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}
