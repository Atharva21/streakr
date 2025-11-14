# Service Layer Tests

## Overview

Comprehensive unit tests for the service layer covering:
- **habits.go**: CRUD operations for habits
- **streaks.go**: Logging, streak calculation, and statistics

## Test Coverage

### Habits Service (habits_test.go)
- ✅ Add habit (improve and quit types)
- ✅ Add habit with/without description
- ✅ Duplicate habit detection
- ✅ Validation (name length, description length)
- ✅ Get habit by name (existing and non-existent)
- ✅ List habits (empty and multiple)
- ✅ Delete habits (single, multiple, cascade delete)
- ✅ Get today's logged habit count

### Streaks Service (streaks_test.go)
- ✅ Log improve habit (first time, consecutive days, missed days)
- ✅ Log quit habit (first time, subsequent logs, clean days)
- ✅ Duplicate log handling
- ✅ Multiple habits logging
- ✅ Get overall stats
- ✅ Get habit stats for date range
- ✅ Heatmap generation
- ✅ Edge cases (habits created mid-month, before/after date ranges)

## Running Tests

### Run All Service Tests
```bash
go test ./internal/service/...
```

### Run Specific Test
```bash
go test ./internal/service/... -run TestAddHabit
```

### Run with Verbose Output
```bash
go test ./internal/service/... -v
```

### Run with Coverage
```bash
go test ./internal/service/... -cover
```

## Known Limitations

Due to the singleton pattern used in `internal/store/store.go` (with `sync.Once` for `BootstrapStore`), all tests share the same database instance. This works well for individual test execution but can cause issues when running the full test suite due to temp directory cleanup.

### Workaround

Run tests individually or in smaller groups:
```bash
go test ./internal/service/... -run TestAddHabit
go test ./internal/service/... -run TestLogHabits
go test ./internal/service/... -run TestGetStats
```

### Future Improvement

To fully resolve this, consider:
1. Refactoring the store package to support dependency injection
2. Adding a test-specific bootstrap function that doesn't use sync.Once
3. Using interfaces for the database layer to enable mocking

## Test Structure

### Test Helpers (test_helpers.go)
- `SetupTestDB()`: Initializes test database
- `CreateTestHabit()`: Helper to create habits with custom timestamps
- `CreateTestStreak()`: Helper to create streak records

### Assertions
All tests use `testify/assert` and `testify/require` for clear, readable assertions.

### Test Organization
Tests are organized by function and include:
- **Positive cases**: Expected behavior with valid inputs
- **Negative cases**: Error handling with invalid inputs
- **Edge cases**: Boundary conditions and special scenarios

## Examples

### Testing Improve Habit
```go
func TestLogHabitsForToday_ImproveHabit_FirstLog(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Cleanup()

    ctx := context.Background()
    habit := testDB.CreateTestHabit(t, ctx, "running", "test", store.HabitTypeImprove, nil)

    allQuitting, err := LogHabitsForToday(ctx, []string{"running"})
    require.NoError(t, err)
    assert.False(t, allQuitting)
}
```

### Testing Quit Habit
```go
func TestLogHabitsForToday_QuitHabit_FirstLog(t *testing.T) {
    testDB := SetupTestDB(t)
    defer testDB.Cleanup()

    ctx := context.Background()
    createdAt := time.Now().AddDate(0, 0, -7)
    habit := testDB.CreateTestHabit(t, ctx, "smoking", "test", store.HabitTypeQuit, &createdAt)

    allQuitting, err := LogHabitsForToday(ctx, []string{"smoking"})
    require.NoError(t, err)
    assert.True(t, allQuitting)
}
```
