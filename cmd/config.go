package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show current configuration",
		Long:  "Display the currently resolved configuration values, including sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfig(*configFlag, *locationFlag, *editorFlag, *editorTypeFlag, *preambleFlag)
		},
	}
}

func runConfig(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag string) error {
	// Resolve all configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)
	editorCmd, err := config.GetEditorCommand(configFlag, editorFlag)
	if err != nil {
		return fmt.Errorf("failed to resolve editor: %w", err)
	}
	editorType := config.GetEditorType(configFlag, editorTypeFlag)
	preamble := config.GetPreamble(configFlag, preambleFlag)

	// Display configuration
	fmt.Println("Current Configuration:")
	fmt.Println("=====================")
	fmt.Println()

	// Plans Directory
	fmt.Printf("Plans Directory: %s\n", plansDir)
	fmt.Printf("  Source: %s\n", getLocationSource(configFlag, locationFlag))
	fmt.Println()

	// Editor
	fmt.Printf("Editor: %s\n", editorCmd)
	fmt.Printf("  Source: %s\n", getEditorSource(configFlag, editorFlag))
	fmt.Println()

	// Editor Type
	fmt.Printf("Editor Type: %s\n", editorType)
	fmt.Printf("  Source: %s\n", getEditorTypeSource(configFlag, editorTypeFlag))
	fmt.Println()

	// Preamble
	if preamble == "" {
		fmt.Println("Preamble: (empty)")
	} else {
		fmt.Printf("Preamble: %s\n", preamble)
	}
	fmt.Printf("  Source: %s\n", getPreambleSource(configFlag, preambleFlag))
	fmt.Println()

	// Config file location
	configPath := config.GetConfigPath(configFlag)
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config File: %s (exists)\n", configPath)
	} else {
		fmt.Printf("Config File: %s (not found)\n", configPath)
	}
	fmt.Printf("  Source: %s\n", getConfigFileSource(configFlag))

	return nil
}

func getLocationSource(configFlag, locationFlag string) string {
	if locationFlag != "" {
		return "command-line flag"
	}
	if os.Getenv("PLAN_LOCATION") != "" {
		return "environment variable (PLAN_LOCATION)"
	}
	if hasConfigValue(configFlag, "PLAN_LOCATION") {
		return "config file"
	}
	return "default"
}

func getEditorSource(configFlag, editorFlag string) string {
	if editorFlag != "" {
		return "command-line flag"
	}
	if os.Getenv("PLAN_EDITOR") != "" {
		return "environment variable (PLAN_EDITOR)"
	}
	if hasConfigValue(configFlag, "PLAN_EDITOR") {
		return "config file"
	}
	return "default"
}

func getEditorTypeSource(configFlag, editorTypeFlag string) string {
	if editorTypeFlag != "" {
		return "command-line flag"
	}
	if os.Getenv("PLAN_EDITOR_TYPE") != "" {
		return "environment variable (PLAN_EDITOR_TYPE)"
	}
	if hasConfigValue(configFlag, "PLAN_EDITOR_TYPE") {
		return "config file"
	}
	return "default"
}

func getPreambleSource(configFlag, preambleFlag string) string {
	if preambleFlag != "" {
		return "command-line flag"
	}
	if os.Getenv("PLAN_PREAMBLE") != "" {
		return "environment variable (PLAN_PREAMBLE)"
	}
	if hasConfigValue(configFlag, "PLAN_PREAMBLE") {
		return "config file"
	}
	return "default"
}

func getConfigFileSource(configFlag string) string {
	if configFlag != "" {
		return "command-line flag"
	}
	if os.Getenv("PLAN_CONFIG") != "" {
		return "environment variable (PLAN_CONFIG)"
	}
	return "default"
}

func hasConfigValue(configFlag, key string) bool {
	configPath := config.GetConfigPath(configFlag)
	file, err := os.Open(configPath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Simple check if key exists in config file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, key+"=") {
			return true
		}
	}
	return false
}
