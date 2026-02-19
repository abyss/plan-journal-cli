package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// DefaultPreamble is the default preamble text for plan files (empty by default)
const DefaultPreamble = ""

// Config holds configuration loaded from file
type Config struct {
	Preamble   string
	Editor     string
	EditorType string
	Location   string
	NoColor    string
}

var loadedConfig *Config
var cachedConfigPath string
var cachedConfigFlag string

// GetConfigPath returns the config file path
// Priority: configFlag > PLAN_CONFIG env var > ~/plans/.config
func GetConfigPath(configFlag string) string {
	// Priority 1: Command-line flag
	if configFlag != "" {
		return expandPath(configFlag)
	}

	// Priority 2: Environment variable
	if envConfigPath := os.Getenv("PLAN_CONFIG"); envConfigPath != "" {
		return expandPath(envConfigPath)
	}

	// Priority 3: Default to ~/plans/.config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot determine home directory: %v", err)
	}
	return filepath.Join(homeDir, "plans", ".config")
}

// loadConfig loads configuration from config file
// Priority: configFlag > PLAN_CONFIG env var > ~/plans/.config
func loadConfig(configFlag string) *Config {
	// Get config path
	configPath := GetConfigPath(configFlag)

	// Return cached config if path and flag haven't changed
	if loadedConfig != nil && cachedConfigPath == configPath && cachedConfigFlag == configFlag {
		return loadedConfig
	}

	// Update cache
	cachedConfigPath = configPath
	cachedConfigFlag = configFlag
	loadedConfig = &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return loadedConfig
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "PLAN_PREAMBLE":
			loadedConfig.Preamble = value
		case "PLAN_EDITOR":
			loadedConfig.Editor = value
		case "PLAN_EDITOR_TYPE":
			loadedConfig.EditorType = value
		case "PLAN_LOCATION":
			loadedConfig.Location = value
		case "PLAN_NO_COLOR":
			loadedConfig.NoColor = value
		}
	}

	return loadedConfig
}

// GetPlansDirectory resolves the plans directory location
// Priority: flag > env > config file > default (~/plans/)
func GetPlansDirectory(configFlag, locationFlag string) string {
	// Priority 1: Command-line flag
	if locationFlag != "" {
		return expandPath(locationFlag)
	}

	// Priority 2: Environment variable
	if envPath := os.Getenv("PLAN_LOCATION"); envPath != "" {
		return expandPath(envPath)
	}

	// Priority 3: Config file
	cfg := loadConfig(configFlag)
	if cfg.Location != "" {
		return expandPath(cfg.Location)
	}

	// Priority 4: Default
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cannot determine home directory: %v", err)
	}
	return filepath.Join(homeDir, "plans")
}

// GetEditorCommand resolves the editor command template
// Priority: flag > env > config file > default (vim)
func GetEditorCommand(configFlag, editorFlag string) (string, error) {
	// Priority 1: Command-line flag
	if editorFlag != "" {
		return resolveEditorCommand(editorFlag)
	}

	// Priority 2: Environment variable
	if envEditor := os.Getenv("PLAN_EDITOR"); envEditor != "" {
		// Treat as custom template if it contains placeholders
		return envEditor, nil
	}

	// Priority 3: Config file
	cfg := loadConfig(configFlag)
	if cfg.Editor != "" {
		return resolveEditorCommand(cfg.Editor)
	}

	// Priority 4: Default to vim
	return BuiltInEditors["vim"].Command, nil
}

// GetEditorType resolves the editor type (terminal, gui, or auto)
// Priority: flag > env > config file > default (auto)
func GetEditorType(configFlag, editorTypeFlag string) string {
	// Priority 1: Command-line flag
	if editorTypeFlag != "" {
		return normalizeEditorType(editorTypeFlag)
	}

	// Priority 2: Environment variable
	if envEditorType := os.Getenv("PLAN_EDITOR_TYPE"); envEditorType != "" {
		return normalizeEditorType(envEditorType)
	}

	// Priority 3: Config file
	cfg := loadConfig(configFlag)
	if cfg.EditorType != "" {
		return normalizeEditorType(cfg.EditorType)
	}

	// Priority 4: Default (auto)
	return "auto"
}

// normalizeEditorType validates and normalizes the editor type
func normalizeEditorType(editorType string) string {
	normalized := strings.ToLower(strings.TrimSpace(editorType))
	switch normalized {
	case "terminal", "gui", "auto":
		return normalized
	default:
		// Invalid value, default to auto
		return "auto"
	}
}

// GetPreamble resolves the preamble text
// Priority: flag > env > config file > default (empty)
func GetPreamble(configFlag, preambleFlag string) string {
	// Priority 1: Command-line flag
	if preambleFlag != "" {
		return preambleFlag
	}

	// Priority 2: Environment variable
	if envPreamble := os.Getenv("PLAN_PREAMBLE"); envPreamble != "" {
		return envPreamble
	}

	// Priority 3: Config file
	cfg := loadConfig(configFlag)
	if cfg.Preamble != "" {
		return cfg.Preamble
	}

	// Priority 4: Default (empty)
	return DefaultPreamble
}

// resolveEditorCommand resolves an editor specification to a command template
// If it's a built-in editor name, return its template
// Otherwise, treat it as a custom template
func resolveEditorCommand(editor string) (string, error) {
	// Check if it's a built-in editor name
	if template, ok := BuiltInEditors[editor]; ok {
		return template.Command, nil
	}

	// Otherwise, treat as custom template
	return editor, nil
}

// GetNoColor resolves whether colors should be disabled
// Priority: flag > NO_COLOR env > PLAN_NO_COLOR env > config file > default (false)
func GetNoColor(configFlag, noColorFlag string) bool {
	// Priority 1: Command-line flag (if provided)
	if noColorFlag != "" {
		return isTruthy(noColorFlag)
	}

	// Priority 2: Standard NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return true // Any value means disable colors
	}

	// Priority 3: PLAN_NO_COLOR environment variable
	if envNoColor := os.Getenv("PLAN_NO_COLOR"); envNoColor != "" {
		return isTruthy(envNoColor)
	}

	// Priority 4: Config file
	cfg := loadConfig(configFlag)
	if cfg.NoColor != "" {
		return isTruthy(cfg.NoColor)
	}

	// Priority 5: Default (colors enabled)
	return false
}

// isTruthy checks if a string value should be considered true
// Accepts: "1", "true", "yes", "y" (case-insensitive)
func isTruthy(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("cannot determine home directory: %v", err)
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}
