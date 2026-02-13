package cmd

import (
	"fmt"
	"time"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/editor"
	"github.com/abyss/plan-journal-cli/pkg/planfile"
	"github.com/spf13/cobra"
)

// NewTomorrowCmd creates the tomorrow command
func NewTomorrowCmd(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:   "tomorrow",
		Short: "Open tomorrow's plan file in editor",
		Long:  "Opens the plan file for tomorrow with cursor positioned at tomorrow's entry insertion point",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTomorrow(*configFlag, *locationFlag, *editorFlag, *editorTypeFlag, *preambleFlag)
		},
	}
}

func runTomorrow(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag string) error {
	// Resolve configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)
	editorCmd, err := config.GetEditorCommand(configFlag, editorFlag)
	if err != nil {
		return fmt.Errorf("failed to resolve editor: %w", err)
	}
	editorType := config.GetEditorType(configFlag, editorTypeFlag)
	preamble := config.GetPreamble(configFlag, preambleFlag)

	// Get tomorrow's date
	now := time.Now().AddDate(0, 0, 1)

	// Ensure month file exists with preamble
	filePath, err := planfile.EnsureMonthFile(now, plansDir, preamble)
	if err != nil {
		return fmt.Errorf("failed to ensure month file: %w", err)
	}

	// Ensure date header exists
	filePath, err = planfile.EnsureDateHeader(now, plansDir)
	if err != nil {
		return fmt.Errorf("failed to ensure date header: %w", err)
	}

	// Find insertion point
	filePath, lineNum, err := planfile.FindInsertionPoint(now, plansDir)
	if err != nil {
		return fmt.Errorf("failed to find insertion point: %w", err)
	}

	// Launch editor
	if err := editor.LaunchEditor(editorCmd, filePath, lineNum, 0, editorType); err != nil {
		return fmt.Errorf("failed to launch editor: %w", err)
	}

	fmt.Printf("Opened %s at line %d\n", filePath, lineNum)
	return nil
}
