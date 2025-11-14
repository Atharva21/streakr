# Streaker is habit tracking cli

![streakrdemo.gif](./docs/images/streakrdemo.gif)

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Commands](#basic-commands)
  - [Understanding Habit Types](#understanding-habit-types)
  - [Calendar View Navigation](#calendar-view-navigation)
  - [Data Storage](#data-storage)
- [Examples](#examples)
- [Contributing](#contributing)
- [Requirements](#requirements)

 Track your habit streaks from the command line.
Build good habits. Break bad ones. Stay consistent.

For design please visit: [design.md](./docs/design.md)

## Features

- **Two Habit Types**:
  - **Improve habits**: Track positive actions you want to do more (running, reading, etc.)
  - **Quit habits**: Track things you want to avoid (smoking, junk food, etc.)
- **Streak Tracking**: View current and max streaks for each habit
- **Calendar View**: Interactive monthly calendar showing your habit history
- **Statistics**: Detailed stats including completed days, missed days, and success rates
- **Simple CLI**: Quick daily logging with minimal commands
- **Local Storage**: All data stored locally in SQLite database (`~/.config/streakr/`)

## Installation

### Download Binary
Download the latest release from [GitHub Releases](https://github.com/Atharva21/streakr/releases/download/v0.1.0/streakr)

```bash
curl -L https://github.com/Atharva21/streakr/releases/download/v0.1.0/streakr -o streakr && chmod +x streakr && sudo mv streakr /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/Atharva21/streakr.git
cd streakr
make build
sudo mv bin/<your-os-type>/streakr /usr/local/bin/
```

or if you have go installed with configured `$GOPATH` or `$GOBIN`
```bash
make install
```

## Usage

### Basic Commands
```bash
# Add a new habit
streakr add <habit_name> [flags]

# Log today's completion
streakr log <habit_name>

# View all habits
streakr list

# View overall stats
streakr stats

# View habit-wise statistics
streakr stats <habit_name>

# Delete a habit
streakr delete <habit_name>
```
#### Flags

- `--type`, `-t`: Set habit type (`improve` or `quit`) - defaults to improve
- `--description`, `-d`: Add description to habit
```bash
# Track improvement habits
streakr add running --description "5k morning run"
streakr add reading -d "Read 30 minutes daily"

# Track quit habits  
streakr add smoking --type quit
streakr add junkfood -t quit -d "No processed snacks"

# Daily logging
streakr log running
streakr log smoking  # Log a slip-up

# View progress
streakr stats           # All habits overview
streakr stats running   # Specific habit details
```

### Understanding Habit Types

**Improve Habits** (default):
- Log each day you complete the habit
- Streaks are built by logging consecutive days
- Example: Log "running" each day you go for a run

**Quit Habits**:
- Log each day you slip up (do the thing you're trying to quit)
- Streaks represent consecutive days WITHOUT the habit
- Example: Log "smoking" only on days you smoke; gaps represent clean days

### Calendar View Navigation

When viewing stats for a specific habit (`streakr stats <habit_name>`):
- Use `←` / `→` arrow keys or `h` / `l` to navigate between months
- Press `q` to quit
- Press `esc` to return to the list view (if navigated from list)

### Data Storage

All your data is stored locally on your machine:
- Database: `~/.config/streakr/streakr.db`
- Logs: `~/.config/streakr/streakr.log`

To backup your data, simply copy the `~/.config/streakr/` directory.

## Contributing

Contributions are welcome! Here's how to get started:

### Development Setup

**First time setup:**
```bash
git clone https://github.com/yourusername/streakr.git
cd streakr

# Install development tools (cobra-cli, sqlc, migrate)
make bootstrap

# Add Go binaries to your PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:$HOME/go/bin

# Build and install
make install
```

The `make bootstrap` command installs:
- `cobra-cli` - For adding new commands
- `sqlc` - For generating type-safe Go code from SQL
- `migrate` - For database migrations

**Common development commands:**
```bash
make help          # Show all available commands
make build         # Build the binary
make test          # Run tests
make clean         # Clean build artifacts
cobra-cli add foo  # Add a new command
```

**Project Structure:**
- `cmd/` - CLI commands (using Cobra)
- `internal/service/` - Business logic for habits and streaks
- `internal/store/` - Database layer (SQLite with sqlc)
- `internal/tui/` - Terminal UI components (using Bubble Tea)
- `internal/store/queries/` - SQL queries (used by sqlc to generate Go code)
- `internal/store/migrations/` - Database migrations

For more details on the architecture, see [CLAUDE.md](CLAUDE.md)

## Examples

### Tracking a Morning Run Habit
```bash
# Add the habit
streakr add running -d "5k morning run"

# Log it each day you run
streakr log running

# Check your progress
streakr stats running
```

### Quitting Coffee
```bash
# Add as a quit habit
streakr add coffee --type quit -d "No caffeine after 2pm"

# Only log when you slip up (have coffee after 2pm)
streakr log coffee

# View your clean streak
streakr stats coffee
```

### Multiple Habits
```bash
# Add several habits
streakr add meditation -d "10 min daily"
streakr add journal -d "Evening reflection"
streakr add snacking -t quit -d "No snacks after dinner"

# View all at once
streakr list

# Check overall progress
streakr stats
```

## Requirements

- Go 1.24 or higher (for building from source)
- Linux/Unix-based system (tested on Linux; WSL2 supported)

## License

MIT License - see LICENSE file for details
