package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/output"
	"github.com/spf13/cobra"
)

// NewConfigCmd creates the config command
func NewConfigCmd(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag, noColorFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show current configuration",
		Long:  "Display the currently resolved configuration values, including sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfig(*configFlag, *locationFlag, *editorFlag, *editorTypeFlag, *preambleFlag, *noColorFlag)
		},
	}
}

func runConfig(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag, noColorFlag string) error {
	// Resolve all configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)
	editorCmd, err := config.GetEditorCommand(configFlag, editorFlag)
	if err != nil {
		return fmt.Errorf("failed to resolve editor: %w", err)
	}
	editorType := config.GetEditorType(configFlag, editorTypeFlag)
	preamble := config.GetPreamble(configFlag, preambleFlag)
	noColor := config.GetNoColor(configFlag, noColorFlag)

	// Display configuration
	fmt.Println(output.Header("Current Configuration:"))
	fmt.Println(output.Header("====================="))
	fmt.Println()

	// Plans Directory
	fmt.Printf("%s: %s\n", output.Bold("Plans Directory"), output.FilePath(plansDir))
	fmt.Printf("  %s: %s\n", output.Info("Source"), getLocationSource(configFlag, locationFlag))
	fmt.Println()

	// Editor
	fmt.Printf("%s: %s\n", output.Bold("Editor"), editorCmd)
	fmt.Printf("  %s: %s\n", output.Info("Source"), getEditorSource(configFlag, editorFlag))
	fmt.Println()

	// Editor Type
	fmt.Printf("%s: %s\n", output.Bold("Editor Type"), editorType)
	fmt.Printf("  %s: %s\n", output.Info("Source"), getEditorTypeSource(configFlag, editorTypeFlag))
	fmt.Println()

	// Preamble
	if preamble == "" {
		fmt.Printf("%s: %s\n", output.Bold("Preamble"), output.Info("(empty)"))
	} else {
		fmt.Printf("%s: %s\n", output.Bold("Preamble"), preamble)
	}
	fmt.Printf("  %s: %s\n", output.Info("Source"), getPreambleSource(configFlag, preambleFlag))
	fmt.Println()

	// No Color
	noColorDisplay := "false"
	if noColor {
		noColorDisplay = "true"
	}
	fmt.Printf("%s: %s\n", output.Bold("No Color"), noColorDisplay)
	fmt.Printf("  %s: %s\n", output.Info("Source"), getNoColorSource(configFlag, noColorFlag))
	fmt.Println()

	// Config file location
	configPath := config.GetConfigPath(configFlag)
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("%s: %s %s\n", output.Bold("Config File"), output.FilePath(configPath), output.Success("(exists)"))
	} else {
		fmt.Printf("%s: %s %s\n", output.Bold("Config File"), output.FilePath(configPath), output.Warning("(not found)"))
	}
	fmt.Printf("  %s: %s\n", output.Info("Source"), getConfigFileSource(configFlag))

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

func getNoColorSource(configFlag, noColorFlag string) string {
	if noColorFlag != "" {
		return "command-line flag"
	}
	if os.Getenv("NO_COLOR") != "" {
		return "environment variable (NO_COLOR)"
	}
	if os.Getenv("PLAN_NO_COLOR") != "" {
		return "environment variable (PLAN_NO_COLOR)"
	}
	if hasConfigValue(configFlag, "PLAN_NO_COLOR") {
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
