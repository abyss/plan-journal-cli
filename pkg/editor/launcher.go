package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// LaunchEditor launches an editor with the given template, file, line, and column
// Template placeholders: %file%, %line%, %column%
// editorType can be "terminal", "gui", or "auto"
func LaunchEditor(template, filePath string, line, column int, editorType string) error {
	// Substitute placeholders
	cmd := strings.ReplaceAll(template, "%file%", filePath)
	cmd = strings.ReplaceAll(cmd, "%line%", fmt.Sprintf("%d", line))
	cmd = strings.ReplaceAll(cmd, "%column%", fmt.Sprintf("%d", column))

	// Parse command into parts
	parts := parseCommand(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// Execute command
	command := exec.Command(parts[0], parts[1:]...)

	// Determine if this is a terminal editor
	isTerminal := false
	switch editorType {
	case "terminal":
		isTerminal = true
	case "gui":
		isTerminal = false
	case "auto":
		// Auto-detect based on editor binary name
		isTerminal = isTerminalEditor(parts[0])
	}

	if isTerminal {
		// For terminal editors, attach to stdin/stdout/stderr and wait
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		return command.Run()
	}

	// For GUI editors, launch without waiting
	return command.Start()
}

// isTerminalEditor checks if the editor binary is a known terminal editor
func isTerminalEditor(editorBinary string) bool {
	// Do you use a terminal editor not listed here? Open an issue or PR!
	terminalEditors := []string{
		"vim", "vi", "nvim", "neovim",
		"nano",
		"micro", "helix", "hx",
		"joe", "jed", "mcedit",
	}

	for _, name := range terminalEditors {
		if strings.HasSuffix(editorBinary, name) || strings.Contains(editorBinary, "/"+name) {
			return true
		}
	}

	return false
}

// parseCommand parses a command string into parts, handling quoted arguments
func parseCommand(cmdStr string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, char := range cmdStr {
		switch {
		case char == '"' || char == '\'':
			if inQuotes {
				if char == quoteChar {
					// End of quoted section
					inQuotes = false
					quoteChar = 0
				} else {
					// Different quote char inside quotes
					current.WriteRune(char)
				}
			} else {
				// Start of quoted section
				inQuotes = true
				quoteChar = char
			}

		case char == ' ' && !inQuotes:
			// Space outside quotes - end of argument
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(char)
		}
	}

	// Add final part
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
