package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Atharva21/streakr/internal/store"
	"github.com/Atharva21/streakr/internal/store/generated"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// TestDB holds test database state
type TestDB struct {
	DB      *sql.DB
	Queries *generated.Queries
	Cleanup func()
}

// SetupTestDB creates an in-memory SQLite database for each test
// This allows tests to run independently without interference
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Create in-memory database for this test
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
	require.NoError(t, err, "Failed to open test database")

	// Create schema
	schema := `
	CREATE TABLE habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE CHECK (length(name) <= 20),
		description TEXT CHECK (description IS NULL OR length(description) <= 200),
		habit_type TEXT CHECK (habit_type IN ('improve', 'quit')) NOT NULL DEFAULT 'improve',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
	);
	CREATE INDEX idx_habits_name ON habits(name);

	CREATE TABLE streaks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		habit_id INTEGER NOT NULL,
		streak_start DATE NOT NULL,
		streak_end DATE NOT NULL,
		FOREIGN KEY (habit_id) REFERENCES habits(id) ON DELETE CASCADE
	);
	CREATE INDEX idx_streaks_habit_id ON streaks(habit_id);
	`

	_, err = db.Exec(schema)
	require.NoError(t, err, "Failed to create schema")

	// Create queries for this test database
	queries := generated.New(db)

	// Inject test database and queries into the store package
	store.SetDBForTesting(db)
	store.SetQueriesForTesting(queries)

	cleanup := func() {
		db.Close()
	}

	return &TestDB{
		DB:      db,
		Queries: queries,
		Cleanup: cleanup,
	}
}

// CreateTestHabit creates a habit for testing with optional created_at time
func (tdb *TestDB) CreateTestHabit(t *testing.T, ctx context.Context, name, description, habitType string, createdAt *time.Time) generated.Habit {
	t.Helper()

	id, err := tdb.Queries.AddHabit(ctx, generated.AddHabitParams{
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
		HabitType: habitType,
	})
	require.NoError(t, err, "Failed to create test habit")

	// Update created_at to the specified time if provided
	if createdAt != nil {
		_, err = tdb.DB.ExecContext(ctx, "UPDATE habits SET created_at = ? WHERE id = ?", *createdAt, id)
		require.NoError(t, err, "Failed to update created_at")
	}

	habit, err := tdb.Queries.GetHabitByName(ctx, name)
	require.NoError(t, err, "Failed to get created habit")

	return habit
}

// CreateTestStreak creates a streak for testing
func (tdb *TestDB) CreateTestStreak(t *testing.T, ctx context.Context, habitID int64, start, end time.Time) {
	t.Helper()

	_, err := tdb.Queries.AddStreak(ctx, generated.AddStreakParams{
		HabitID:     habitID,
		StreakStart: start,
		StreakEnd:   end,
	})
	require.NoError(t, err, "Failed to create test streak")
}
