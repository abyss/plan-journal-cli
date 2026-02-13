package main

import (
	"fmt"
	"os"

	"github.com/abyss/plan-journal-cli/cmd"
	"github.com/spf13/cobra"
)

var (
	configFlag     string
	locationFlag   string
	editorFlag     string
	editorTypeFlag string
	preambleFlag   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "plan",
		Short: "Plan Journal CLI - Manage daily plan files",
		Long: `Plan Journal CLI helps you manage daily plan files organized by month.
Files are structured with month headers and chronologically ordered date sections.`,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "Override config file location (default: ~/plans/.config)")
	rootCmd.PersistentFlags().StringVar(&locationFlag, "location", "", "Override plans directory (default: ~/plans/)")
	rootCmd.PersistentFlags().StringVar(&editorFlag, "editor", "", "Override default editor (default: vim)")
	rootCmd.PersistentFlags().StringVar(&editorTypeFlag, "editor-type", "", "Override editor type: terminal, gui, or auto (default: auto)")
	rootCmd.PersistentFlags().StringVar(&preambleFlag, "preamble", "", "Override preamble text (default: empty)")

	// Add commands
	rootCmd.AddCommand(cmd.NewTodayCmd(&configFlag, &locationFlag, &editorFlag, &editorTypeFlag, &preambleFlag))
	rootCmd.AddCommand(cmd.NewTomorrowCmd(&configFlag, &locationFlag, &editorFlag, &editorTypeFlag, &preambleFlag))
	rootCmd.AddCommand(cmd.NewEditCmd(&configFlag, &locationFlag, &editorFlag, &editorTypeFlag, &preambleFlag))
	rootCmd.AddCommand(cmd.NewReadCmd(&configFlag, &locationFlag))
	rootCmd.AddCommand(cmd.NewFixCmd(&configFlag, &locationFlag, &preambleFlag))
	rootCmd.AddCommand(cmd.NewConfigCmd(&configFlag, &locationFlag, &editorFlag, &editorTypeFlag, &preambleFlag))

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
