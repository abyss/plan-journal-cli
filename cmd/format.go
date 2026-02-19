package cmd

import (
	"fmt"
	"strings"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/output"
	"github.com/abyss/plan-journal-cli/pkg/planfile"
	"github.com/spf13/cobra"
)

// NewFormatCmd creates the format command
func NewFormatCmd(configFlag, locationFlag, preambleFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:     "format <target>",
		Aliases: []string{"fmt", "fix"},
		Short:   "Format plan file",
		Long:    "Format plan files by reordering date sections chronologically and updating/adding preamble. Target can be a date (YYYY-MM, YYYY-MM-DD, today, etc.), a file path, or a filename in the plans directory",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFormat(*configFlag, *locationFlag, *preambleFlag, args[0])
		},
	}
}

func runFormat(configFlag, locationFlag, preambleFlag, target string) error {
	// Resolve configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)
	preamble := config.GetPreamble(configFlag, preambleFlag)

	// Format plan file
	result, err := planfile.FormatPlanFile(target, plansDir, preamble)
	if err != nil {
		return fmt.Errorf("failed to format plan file: %w", err)
	}

	// Display result with color
	if strings.HasPrefix(result, "No changes") {
		fmt.Println(output.Info(result))
	} else {
		// Format "Changes: ..." messages in green
		if changes, found := strings.CutPrefix(result, "Changes: "); found {
			fmt.Printf("%s %s\n", output.Bold("Changes:"), output.Success(changes))
		} else {
			fmt.Println(output.Success(result))
		}
	}
	return nil
}
