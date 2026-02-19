package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/output"
	"github.com/abyss/plan-journal-cli/pkg/planfile"
	"github.com/spf13/cobra"
)

// NewReadCmd creates the read command
func NewReadCmd(configFlag, locationFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:     "read <target>",
		Aliases: []string{"view"},
		Short:   "Read plan entries",
		Long:    "Display plan entries for 'yesterday', 'today', 'tomorrow', a specific month (YYYY-MM), or a specific date (YYYY-MM-DD)",
		Args:    cobra.ExactArgs(1),
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
			fmt.Println(output.Info(err.Error()))
			return nil
		}
		return fmt.Errorf("failed to read entries: %w", err)
	}

	// Display content with colorized date headers
	colorized := colorizePlanContent(content)
	fmt.Println(colorized)
	return nil
}

// colorizePlanContent adds color to date headers in plan content
func colorizePlanContent(content string) string {
	// Regex to match date headers like "## 2026-02-19"
	dateHeaderRegex := regexp.MustCompile(`(?m)^## (\d{4}-\d{2}-\d{2})$`)

	// Replace date headers with colorized versions
	return dateHeaderRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract just the date part
		dateStr := strings.TrimPrefix(match, "## ")
		return "## " + output.DateBlue(dateStr)
	})
}
