# Streaker is habit tracking cli

![streakrdemo.gif](./docs/images/streakrdemo.gif)

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)

 Track your habit streaks from the command line.
Build good habits. Break bad ones. Stay consistent.

For design please visit: [design.md](./docs/design.md)

## Features

- Track improvement habits (daily runs, reading, etc.)
- Monitor quit habits (smoking, junk food, etc.)
- View current and max streaks across all habits
- Simple CLI interface for quick daily logging

## Installation

<!-- ### Download Binary -->
<!-- Download the latest release from [GitHub Releases](https://github.com/yourusername/streakr/releases) -->

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

# View statistics
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

## Contributing

Contributions are welcome! Here's how to get started:

### Development Setup
```bash
git clone https://github.com/yourusername/streakr.git
cd streakr
make install
```
install cobra-cli via 
```bash
go install github.com/spf13/cobra-cli@latest
```
this lets you add a command or subcommand via 
```bash
cobra-cli add foo
```
