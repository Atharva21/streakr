# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Streakr is a CLI habit tracker built in Go that helps users track daily habits and maintain streaks. It supports two habit types:
- **Improve habits**: Track positive actions you want to do more (e.g., running, reading)
- **Quit habits**: Track things you want to avoid (e.g., smoking, junk food)

## Development Commands

### Building and Running
```bash
make build          # Builds binary for Linux to bin/linux/streakr
make install        # Builds and installs to $GOPATH/bin or $GOBIN
make run            # Builds and runs the binary
make clean-install  # Clean build artifacts and fresh install
```

### Code Generation
```bash
make generate       # Generates Go code from SQL using sqlc
```
Run this after modifying SQL queries in `internal/store/queries/` or schema in `internal/store/migrations/`.

### Testing and Formatting
```bash
make test           # Run all tests
make fmt            # Format code with go fmt
make tidy           # Tidy go.mod and go.sum
```

### Adding New Commands
```bash
cobra-cli add <command_name>  # Generates new command file in cmd/
```
Requires `github.com/spf13/cobra-cli@latest` to be installed.

## Architecture

### Core Application Flow

**Bootstrap Sequence** (triggered by `internal/streakr/streakr.go` init()):
1. Config initialization - Creates `~/.config/streakr/` directory structure
2. Logger setup - Configures file logging with lumberjack rotation
3. Store initialization - Opens SQLite database and runs migrations
4. Signal handling - Main context with graceful shutdown on SIGINT/SIGTERM

**Command Execution** ([cmd/root.go](cmd/root.go)):
- All commands built with cobra framework
- Custom error handling via `streakrerror.StreakrError` for user-friendly messages
- Context passed through for cancellation and graceful shutdown

### Database Design

**Schema** ([internal/store/migrations/000001_initial_schema.up.sql](internal/store/migrations/000001_initial_schema.up.sql)):
- `habits` table: Stores habit metadata (name, description, type, created_at)
- `streaks` table: Stores date ranges (streak_start, streak_end) for each habit

**Key Concept**: The `streaks` table represents different things based on habit type:
- **Improve habits**: Ranges are days the habit WAS performed. Gaps = missed days.
- **Quit habits**: Ranges are days the habit WAS NOT performed (clean days). Gaps = slip-ups.

**Code Generation**: Uses sqlc to generate type-safe Go code from SQL queries
- Queries: `internal/store/queries/*.sql`
- Generated: `internal/store/generated/*.go`
- Config: [sqlc.yaml](sqlc.yaml)

### Streak Calculation Logic

**Current Streak** ([internal/service/streaks.go](internal/service/streaks.go)):
- **Improve habits**: Latest streak counts if streak_end is today or yesterday
- **Quit habits**: Days since last streak_end (representing last slip-up), minus 1

**Max Streak**:
- **Improve habits**: `MAX(1 + streak_end - streak_start)` - inclusive of both dates
- **Quit habits**: `MAX(streak_end - streak_start)` - end date not inclusive

**Logging Today** ([internal/service/streaks.go:15-111](internal/service/streaks.go#L15-L111)):
- **Improve habits**: Extends previous streak if yesterday, otherwise creates new streak
- **Quit habits**: Complex logic accounting for "clean days" since last slip-up

### Package Structure

**internal/store/** - Database layer
- [store.go](internal/store/store.go): Bootstrap logic, database connection, migrations
- `queries/`: Raw SQL queries for sqlc
- `generated/`: Auto-generated type-safe query code (do not edit manually)
- `migrations/`: SQL schema migrations

**internal/service/** - Business logic layer
- [habits.go](internal/service/habits.go): CRUD operations for habits
- [streaks.go](internal/service/streaks.go): Streak calculation and logging logic

**internal/tui/** - Terminal UI components using Bubble Tea
- [listview.go](internal/tui/listview.go): Habit list display
- [calview.go](internal/tui/calview.go): Calendar heatmap view
- [overallstats.go](internal/tui/overallstats.go): Stats overview

**cmd/** - CLI command definitions (Cobra)
- One file per command (add, log, list, stats, delete)
- [root.go](cmd/root.go): Entry point with custom error handling

**internal/config/** - Application configuration
- Creates `~/.config/streakr/` for data and logs

**internal/util/** - Shared utilities
- [date.go](internal/util/date.go): Date comparison and manipulation helpers

### Important Patterns

**Singleton Bootstrap Pattern**: Most internal packages use `sync.Once` to ensure single initialization:
```go
var bootstrapOnce sync.Once
func BootstrapX() {
    bootstrapOnce.Do(func() { /* init code */ })
}
```

**Error Handling**: Use `streakrerror.StreakrError` for user-facing errors:
- `TerminalMsg`: User-friendly error message
- `ShowUsage`: Whether to display command usage
- Generic errors bubble up and get logged

**Database Access**: Never call `sql.Open` directly, always use `store.GetDB()` or `store.GetQueries()` after bootstrap.

**Date Handling**: All date comparisons use `util.IsSameDate()`, `util.CompareDate()`, `util.GetDayDiff()` to handle timezone-agnostic date logic.

## Data Storage

- Database: `~/.config/streakr/streakr.db` (SQLite)
- Logs: `~/.config/streakr/streakr.log` (Rotated by lumberjack)

## Dependencies

Key external libraries:
- `spf13/cobra`: CLI framework
- `charmbracelet/bubbletea`: TUI framework
- `mattn/go-sqlite3`: SQLite driver
- `golang-migrate/migrate`: Database migrations
- `sqlc`: SQL to Go code generation
