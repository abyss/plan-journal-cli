package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetPlansDirectory(t *testing.T) {
	// Save original env vars
	origLocation := os.Getenv("PLAN_LOCATION")
	origConfig := os.Getenv("PLAN_CONFIG")
	defer func() {
		os.Setenv("PLAN_LOCATION", origLocation)
		os.Setenv("PLAN_CONFIG", origConfig)
		loadedConfig = nil
		cachedConfigPath = ""
		cachedConfigFlag = ""
	}()

	tests := []struct {
		name         string
		locationFlag string
		envLocation  string
		wantContains string
	}{
		{
			name:         "flag takes priority",
			locationFlag: "/tmp/test-plans",
			envLocation:  "/tmp/env-plans",
			wantContains: "/tmp/test-plans",
		},
		{
			name:         "env when no flag",
			locationFlag: "",
			envLocation:  "/tmp/env-plans",
			wantContains: "/tmp/env-plans",
		},
		{
			name:         "default when no flag or env",
			locationFlag: "",
			envLocation:  "",
			wantContains: "plans",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cached config
			loadedConfig = nil
			cachedConfigPath = ""
			cachedConfigFlag = ""

			// Set env vars
			if tt.envLocation != "" {
				os.Setenv("PLAN_LOCATION", tt.envLocation)
			} else {
				os.Unsetenv("PLAN_LOCATION")
			}
			os.Unsetenv("PLAN_CONFIG") // Don't use config file

			result := GetPlansDirectory("", tt.locationFlag)
			if result != tt.wantContains && !contains(result, tt.wantContains) {
				t.Errorf("GetPlansDirectory() = %v, want to contain %v", result, tt.wantContains)
			}
		})
	}
}

func TestGetEditorCommand(t *testing.T) {
	// Save original env vars
	origEditor := os.Getenv("PLAN_EDITOR")
	origConfig := os.Getenv("PLAN_CONFIG")
	defer func() {
		os.Setenv("PLAN_EDITOR", origEditor)
		os.Setenv("PLAN_CONFIG", origConfig)
		loadedConfig = nil
		cachedConfigPath = ""
		cachedConfigFlag = ""
	}()

	tests := []struct {
		name       string
		editorFlag string
		envEditor  string
		want       string
	}{
		{
			name:       "flag takes priority",
			editorFlag: "vim",
			envEditor:  "emacs",
			want:       "vim +%line% %file%",
		},
		{
			name:       "env when no flag",
			editorFlag: "",
			envEditor:  "custom-editor",
			want:       "custom-editor",
		},
		{
			name:       "builtin vscode name resolves",
			editorFlag: "vscode",
			envEditor:  "",
			want:       "code --goto %file%:%line%:%column%",
		},
		{
			name:       "builtin vim name resolves",
			editorFlag: "vim",
			envEditor:  "",
			want:       "vim +%line% %file%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cached config
			loadedConfig = nil
			cachedConfigPath = ""
			cachedConfigFlag = ""

			// Set env vars
			if tt.envEditor != "" {
				os.Setenv("PLAN_EDITOR", tt.envEditor)
			} else {
				os.Unsetenv("PLAN_EDITOR")
			}
			os.Unsetenv("PLAN_CONFIG") // Don't use config file

			result, err := GetEditorCommand("", tt.editorFlag)
			if err != nil {
				t.Errorf("GetEditorCommand() error = %v", err)
				return
			}
			if result != tt.want {
				t.Errorf("GetEditorCommand() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestGetPreamble(t *testing.T) {
	// Save original env vars
	origPreamble := os.Getenv("PLAN_PREAMBLE")
	origConfig := os.Getenv("PLAN_CONFIG")
	defer func() {
		os.Setenv("PLAN_PREAMBLE", origPreamble)
		os.Setenv("PLAN_CONFIG", origConfig)
		loadedConfig = nil
		cachedConfigPath = ""
		cachedConfigFlag = ""
	}()

	tests := []struct {
		name         string
		preambleFlag string
		envPreamble  string
		want         string
	}{
		{
			name:         "flag takes priority",
			preambleFlag: "Flag preamble",
			envPreamble:  "Env preamble",
			want:         "Flag preamble",
		},
		{
			name:         "env when no flag",
			preambleFlag: "",
			envPreamble:  "Env preamble",
			want:         "Env preamble",
		},
		{
			name:         "default empty when no flag or env",
			preambleFlag: "",
			envPreamble:  "",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cached config
			loadedConfig = nil
			cachedConfigPath = ""
			cachedConfigFlag = ""

			// Set env vars
			if tt.envPreamble != "" {
				os.Setenv("PLAN_PREAMBLE", tt.envPreamble)
			} else {
				os.Unsetenv("PLAN_PREAMBLE")
			}
			// Point to non-existent config file to ensure we get defaults
			os.Setenv("PLAN_CONFIG", "/tmp/nonexistent-config-file-for-testing-12345")

			result := GetPreamble("", tt.preambleFlag)
			if result != tt.want {
				t.Errorf("GetPreamble() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".config")

	configContent := `# Test config
PLAN_PREAMBLE=Test preamble from file
PLAN_EDITOR=vim +%line% %file%
PLAN_LOCATION=/test/location
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Save and set env vars
	origConfig := os.Getenv("PLAN_CONFIG")
	defer func() {
		os.Setenv("PLAN_CONFIG", origConfig)
		loadedConfig = nil
		cachedConfigPath = ""
		cachedConfigFlag = ""
	}()

	os.Setenv("PLAN_CONFIG", configPath)
	os.Unsetenv("PLAN_PREAMBLE")
	os.Unsetenv("PLAN_EDITOR")
	os.Unsetenv("PLAN_LOCATION")

	// Reset cached config
	loadedConfig = nil
	cachedConfigPath = ""

	// Test preamble from config
	preamble := GetPreamble("", "")
	if preamble != "Test preamble from file" {
		t.Errorf("GetPreamble() from config = %v, want %v", preamble, "Test preamble from file")
	}

	// Reset for editor test
	loadedConfig = nil
	cachedConfigPath = ""

	// Test editor from config
	editor, err := GetEditorCommand("", "")
	if err != nil {
		t.Errorf("GetEditorCommand() error = %v", err)
	}
	if editor != "vim +%line% %file%" {
		t.Errorf("GetEditorCommand() from config = %v, want %v", editor, "vim +%line% %file%")
	}
}

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "expand tilde",
			path: "~/plans",
			want: filepath.Join(homeDir, "plans"),
		},
		{
			name: "no expansion needed",
			path: "/absolute/path",
			want: "/absolute/path",
		},
		{
			name: "relative path",
			path: "./relative",
			want: "./relative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.path)
			if result != tt.want {
				t.Errorf("expandPath() = %v, want %v", result, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || filepath.Base(s) == substr || filepath.Dir(s) == substr)
}
