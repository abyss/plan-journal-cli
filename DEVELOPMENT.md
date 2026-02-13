# Development Guide

This guide covers local development, testing, and contributing to Plan Journal CLI.

## Requirements

- Go 1.25 or later
- [Task](https://taskfile.dev/) - A task runner / build tool

Install Task:
```bash
# macOS
brew install go-task

# Or see https://taskfile.dev/installation/ for other platforms
```

## Building

```bash
# Build the binary to bin/plan
task build

# Install to $GOPATH/bin
task install

# Clean build artifacts
task clean
```

## Running Tests

```bash
# Run all tests
task test

# Or use go directly for more options
go test ./... -v           # verbose output
go test ./pkg/config       # specific package
go test ./... -cover       # with coverage
```

## Manual Testing

Use local directories for testing to keep your real plans separate from development data:

```bash
# Build the binary
task build

# Set up local test environment
export PLAN_CONFIG=./plans/.config
export PLAN_LOCATION=./plans/

# Create test config
cat > ./plans/.config << 'EOF'
# Test configuration
PLAN_PREAMBLE=Test preamble for development
EOF

# Test commands
./bin/plan today
./bin/plan read today
./bin/plan fix 2026-02
./bin/plan config
```

This keeps your development/testing data in `./plans/` while your real plans remain in `~/plans/`.

## Project Structure

```
plan-journal-cli/
├── main.go                      # CLI entry point
├── bin/                         # Build output (gitignored)
│   └── plan                     # Compiled binary
├── cmd/                         # Command implementations
│   ├── today.go
│   ├── read.go
│   ├── fix.go
│   └── config.go
└── pkg/
    ├── config/                  # Configuration resolution
    │   ├── config.go
    │   ├── config_test.go
    │   └── editors.go
    ├── dateutil/                # Date parsing utilities
    │   ├── dateutil.go
    │   └── dateutil_test.go
    ├── editor/                  # Editor launcher
    │   ├── launcher.go
    │   └── launcher_test.go
    └── planfile/                # Plan file management
        ├── manager.go
        ├── manager_test.go
        ├── parser.go
        ├── parser_test.go
        └── writer.go
```

## Test Coverage

Tests are organized by package:

### `pkg/dateutil`
- Date parsing (today, YYYY-MM, YYYY-MM-DD)
- Date formatting and validation
- Date comparison and ordering

### `pkg/config`
- Configuration priority resolution (flag > env > config > default)
- Config file loading and parsing
- Path expansion

### `pkg/editor`
- Command parsing with quotes and spaces
- Template placeholder substitution

### `pkg/planfile`
- File parsing and structure validation
- Month file creation and management
- Date section ordering
- Preamble management
- File repair operations

## Adding New Features

1. Write tests first (TDD approach recommended)
2. Implement the feature
3. Ensure all tests pass: `task test`
4. Update documentation if needed
5. Build and test manually with local plans directory: `task build`

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing or `task tidy` to tidy modules
- Keep functions focused and testable
- Use descriptive variable names
- Add comments for non-obvious logic

## Contributing

We welcome contributions! Here are some ways to help:

**Easy contributions:**
- Add support for new editors (see "Adding Support for New Editors" below)
- Improve documentation
- Fix typos or clarify README sections

**Code contributions:**
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes with tests
4. Ensure all tests pass: `task test`
5. Build and verify: `task build`
6. Commit your changes
7. Submit a pull request

Please include:
- Clear description of what you're adding/fixing
- Tests for new functionality
- Updates to documentation if applicable

## Common Development Tasks

### Adding a New Command

1. Create `cmd/yourcommand.go`
2. Implement `NewYourCmd()` function
3. Add command to `main.go` with `rootCmd.AddCommand()`
4. Write tests for the command logic

### Adding Configuration Options

1. Update `pkg/config/config.go` with new getter function
2. Add priority resolution (flag > env > config > default)
3. Update `cmd/config.go` to display the new setting
4. Write tests in `pkg/config/config_test.go`
5. Update README.md with usage examples

### Adding Support for New Editors

The tool includes predefined command templates for common editors. To add support for a new editor:

1. **Edit `pkg/config/editors.go`** and add your editor to the `BuiltInEditors` map:

```go
var BuiltInEditors = map[string]EditorTemplate{
    "vscode": {
        Name:    "Visual Studio Code",
        Command: "code --goto %file%:%line%:%column%",
    },
    "vim": {
        Name:    "Vim",
        Command: "vim +%line% %file%",
    },
    // Add your editor here
    "emacs": {
        Name:    "GNU Emacs",
        Command: "emacs +%line% %file%",
    },
}
```

2. **Use template placeholders** in your command:
   - `%file%` - The file path
   - `%line%` - The line number for cursor positioning
   - `%column%` - The column number for cursor positioning

3. **Test your editor** works correctly:

```bash
task build
./bin/plan --editor youreditor today
```

4. **Add tests** in `pkg/config/config_test.go`:

```go
{
    name:       "builtin youreditor name resolves",
    editorFlag: "youreditor",
    envEditor:  "",
    want:       "youreditor +%line% %file%",
},
```

5. **Update documentation** in README.md under "Predefined editor names"

6. **Submit a pull request** with:
   - A clear description of the editor being added
   - Confirmation that you've tested it works
   - Updates to both code and documentation

**Popular editors to consider contributing:**
- Emacs
- Sublime Text
- Nano
- Neovim (separate from Vim)
- IntelliJ/GoLand
- Helix
- Zed

### Modifying File Format

1. Update parser in `pkg/planfile/parser.go`
2. Update writer in `pkg/planfile/writer.go`
3. Update tests in `pkg/planfile/parser_test.go`
4. Ensure backward compatibility if possible
5. Update file format documentation in README.md

## Debugging

### Configuration Issues

Use `plan config` to see what configuration is being used and where it comes from:

```bash
./bin/plan config
```

### File Parsing Issues

Run tests with verbose output to see detailed parsing results:

```bash
go test ./pkg/planfile -v
```

### Editor Launch Issues

Test editor command parsing directly:

```bash
go test ./pkg/editor -v -run TestParseCommand
```

## Release Process

1. Ensure all tests pass: `task test`
2. Update version in code if applicable
3. Build binary: `task build`
4. Test binary with real usage: `./bin/plan`
5. Tag release
6. Build for multiple platforms if needed
