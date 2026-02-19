# Plan Journal CLI

A command-line tool for managing daily plan files organized by month. Files are stored as plain text with month headers and chronologically ordered date sections.

Inspired by [Matteo Landi's .plan files](https://matteolandi.net/plan-files.html) and John Carmack's early .plan files.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install abyss/tools/plan
```

### From Source

#### Requirements

- Go 1.25 or later
- [Task](https://taskfile.dev/) - A task runner / build tool

Install Task:
```bash
# macOS
brew install go-task

# Or see https://taskfile.dev/installation/ for other platforms
```

#### Building

```bash
# Build the binary to bin/plan
task build

# Install to $GOPATH/bin
task install
```

## Quick Start

```bash
# Open today's entry in your editor
plan today

# Open a specific date in your editor
plan edit 2026-02-13
plan edit yesterday

# Read today's entries
plan read today

# Read an entire month
plan read 2026-02

# Read a specific date
plan read 2026-02-13
```

## Commands

- **`plan edit <target>`** - Open a plan entry in your editor for the specified date
- **`plan today`** - Shortcut for `plan edit today`
- **`plan tomorrow`** - Shortcut for `plan edit tomorrow`
- **`plan read <target>`** - Display entries for a target (see below)
- **`plan list [filter]`** - List all dates with entries, optionally filtered by year (YYYY) or month (YYYY-MM)
- **`plan format <target>`** - Format file by reordering dates and updating preamble (target can be a date, file path, or filename)
- **`plan config`** - Show current configuration and sources

**Colors:** The CLI uses minimal color (green for today, red for errors). Disable with `NO_COLOR=1`, `PLAN_NO_COLOR=true`, or the `--no-color` flag. Test colors with `plan colors`.

### Special Dates

The `edit`, `read`, and `format` commands accept special date keywords:
- **`yesterday`** - Previous day
- **`today`** - Current day
- **`tomorrow`** - Next day

Examples:
```bash
plan edit yesterday
plan read today
plan format tomorrow
```

You can also use specific dates (`YYYY-MM-DD`) or entire months (`YYYY-MM`).

### File Paths

The `format` command also accepts file paths:
- **Absolute path**: `plan format /path/to/2026-01.plan`
- **Relative path**: `plan format ./plans/2026-01.plan`
- **Filename** (in plans directory): `plan format 2026-01.plan`

## Configuration

All configuration settings follow a consistent priority order:
1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Config file** at `~/plans/.config`
4. **Built-in defaults** (lowest priority)

Use `plan config` to see your current resolved configuration.

### Configuration Options

| Setting | Flag | Environment | Config File | Default |
|---------|------|-------------|-------------|---------|
| **Config File** | `--config` | `PLAN_CONFIG` | (none) | `~/plans/.config` |
| **Plans Directory** | `--location` | `PLAN_LOCATION` | `PLAN_LOCATION=` | `~/plans/` |
| **Editor** | `--editor` | `PLAN_EDITOR` | `PLAN_EDITOR=` | `vim` |
| **Editor Type** | `--editor-type` | `PLAN_EDITOR_TYPE` | `PLAN_EDITOR_TYPE=` | `auto` |
| **Preamble** | `--preamble` | `PLAN_PREAMBLE` | `PLAN_PREAMBLE=` | empty |
| **No Color** | `--no-color` | `NO_COLOR`, `PLAN_NO_COLOR` | `PLAN_NO_COLOR=` | `false` |

### Config File

Create `~/plans/.config` with your settings:

```bash
# Plans directory location
PLAN_LOCATION=~/plans

# Editor command (predefined: vim, vscode | or custom template with %file%, %line%, %column%)
PLAN_EDITOR=vim

# Editor type: terminal, gui, or auto
PLAN_EDITOR_TYPE=auto

# Preamble text for plan files
PLAN_PREAMBLE=Your custom preamble text here

# Disable color output (true/false)
PLAN_NO_COLOR=false
```

Override config file location with `--config` flag or `PLAN_CONFIG` environment variable.

## File Format

Files are named `YYYY-MM.plan` with month header (`# YYYY-MM`), optional preamble, and chronologically ordered date sections (`## YYYY-MM-DD`):

```markdown
# 2026-02

[Optional preamble text]

## 2026-02-13
Your entries for this day...

## 2026-02-14
Your entries for this day...
```

## Development

For information on building, testing, and contributing to this project, see [DEVELOPMENT.md](DEVELOPMENT.md).

## License
This project is licensed under the [ISC License](LICENSE).

Copyright (c) 2026 Abyss
