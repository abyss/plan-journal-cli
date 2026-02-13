package cmd

import (
	"fmt"
	"strings"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/planfile"
	"github.com/spf13/cobra"
)

// NewReadCmd creates the read command
func NewReadCmd(configFlag, locationFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:   "read <target>",
		Short: "Read plan entries",
		Long:  "Display plan entries for 'today', a specific month (YYYY-MM), or a specific date (YYYY-MM-DD)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRead(*configFlag, *locationFlag, args[0])
		},
	}
}

func runRead(configFlag, locationFlag, target string) error {
	// Resolve configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)

	// Read entries
	content, err := planfile.ReadEntries(target, plansDir)
	if err != nil {
		// If it's just "no entries found", print without usage
		if strings.Contains(err.Error(), "no entries found") {
			fmt.Println(err.Error())
			return nil
		}
		return fmt.Errorf("failed to read entries: %w", err)
	}

	// Display content
	fmt.Println(content)
	return nil
}
