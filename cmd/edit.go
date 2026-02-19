package cmd

import (
	"fmt"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/dateutil"
	"github.com/abyss/plan-journal-cli/pkg/editor"
	"github.com/abyss/plan-journal-cli/pkg/output"
	"github.com/abyss/plan-journal-cli/pkg/planfile"
	"github.com/spf13/cobra"
)

// NewEditCmd creates the edit command
func NewEditCmd(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:     "edit <target>",
		Aliases: []string{"open"},
		Short:   "Open a plan entry in editor",
		Long:    "Opens a plan file with cursor positioned at the specified date entry. Target can be 'yesterday', 'today', 'tomorrow', or a specific date (YYYY-MM-DD)",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(*configFlag, *locationFlag, *editorFlag, *editorTypeFlag, *preambleFlag, args[0])
		},
	}
}

func runEdit(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag, target string) error {
	// Resolve configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)
	editorCmd, err := config.GetEditorCommand(configFlag, editorFlag)
	if err != nil {
		return fmt.Errorf("failed to resolve editor: %w", err)
	}
	editorType := config.GetEditorType(configFlag, editorTypeFlag)
	preamble := config.GetPreamble(configFlag, preambleFlag)

	// Parse target date
	date, err := dateutil.ParseTarget(target)
	if err != nil {
		return fmt.Errorf("failed to parse target: %w", err)
	}

	// Ensure month file exists with preamble
	if err := planfile.EnsureMonthFile(date, plansDir, preamble); err != nil {
		return fmt.Errorf("failed to ensure month file: %w", err)
	}

	// Ensure date header exists
	if err := planfile.EnsureDateHeader(date, plansDir); err != nil {
		return fmt.Errorf("failed to ensure date header: %w", err)
	}

	// Find insertion point
	filePath, lineNum, err := planfile.FindInsertionPoint(date, plansDir)
	if err != nil {
		return fmt.Errorf("failed to find insertion point: %w", err)
	}

	// Launch editor
	if err := editor.LaunchEditor(editorCmd, filePath, lineNum, 0, editorType); err != nil {
		return fmt.Errorf("failed to launch editor: %w", err)
	}

	fmt.Printf("Opened %s at line %s\n",
		output.Bold(filePath),
		output.Bold(fmt.Sprintf("%d", lineNum)))
	return nil
}
