package cmd

import (
	"fmt"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/planfile"
	"github.com/spf13/cobra"
)

// NewFixCmd creates the fix command
func NewFixCmd(configFlag, locationFlag, preambleFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:   "fix <target>",
		Short: "Fix plan file issues",
		Long:  "Repair plan files by reordering date sections chronologically and updating/adding preamble. Target can be a month (YYYY-MM) or specific date (YYYY-MM-DD)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFix(*configFlag, *locationFlag, *preambleFlag, args[0])
		},
	}
}

func runFix(configFlag, locationFlag, preambleFlag, target string) error {
	// Resolve configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)
	preamble := config.GetPreamble(configFlag, preambleFlag)

	// Fix plan file
	result, err := planfile.FixPlanFile(target, plansDir, preamble)
	if err != nil {
		return fmt.Errorf("failed to fix plan file: %w", err)
	}

	// Display result
	fmt.Println(result)
	return nil
}
