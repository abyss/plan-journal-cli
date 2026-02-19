package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/abyss/plan-journal-cli/pkg/config"
	"github.com/abyss/plan-journal-cli/pkg/output"
	"github.com/abyss/plan-journal-cli/pkg/planfile"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command
func NewListCmd(configFlag, locationFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:     "list [filter]",
		Aliases: []string{"ls"},
		Short:   "List all available dates in plan files",
		Long: `List all dates that have plan entries, grouped by month.

The optional filter argument allows you to show only dates from a specific year or month:
  - YYYY: Show all dates from that year (e.g., 2026)
  - YYYY-MM: Show all dates from that month (e.g., 2026-02)

Examples:
  plan list              # Show all dates
  plan list 2026         # Show dates from 2026
  plan list 2026-02      # Show dates from February 2026`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filter := ""
			if len(args) > 0 {
				filter = args[0]
			}
			return runList(*configFlag, *locationFlag, filter)
		},
	}
}

func runList(configFlag, locationFlag, filter string) error {
	// Resolve configuration
	plansDir := config.GetPlansDirectory(configFlag, locationFlag)

	// Discover dates
	datesByMonth, err := planfile.DiscoverDates(plansDir, filter)
	if err != nil {
		return fmt.Errorf("failed to discover dates: %w", err)
	}

	// Check if any dates were found
	if len(datesByMonth) == 0 {
		if filter != "" {
			fmt.Println(output.Info(fmt.Sprintf("No dates found for %s", filter)))
		} else {
			fmt.Println(output.Info("No dates found"))
		}
		return nil
	}

	// Sort months chronologically
	months := make([]string, 0, len(datesByMonth))
	for month := range datesByMonth {
		months = append(months, month)
	}
	sort.Strings(months)

	// Display dates grouped by month
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	yesterday := today.AddDate(0, 0, -1)
	yesterdayStr := yesterday.Format("2006-01-02")
	tomorrow := today.AddDate(0, 0, 1)
	tomorrowStr := tomorrow.Format("2006-01-02")

	for i, month := range months {
		if i > 0 {
			fmt.Println() // Blank line between months
		}
		fmt.Printf("%s:\n", output.Header(month))
		for _, dateStr := range datesByMonth[month] {
			// Add today indicator
			indicator := "  "
			if dateStr == todayStr {
				indicator = output.DateGreen("> ")
			} else {
				indicator = "  "
			}

			// Add relative date label
			var label string
			if dateStr == todayStr {
				label = " " + output.Highlight("[today]")
			} else if dateStr == yesterdayStr {
				label = " [yesterday]"
			} else if dateStr == tomorrowStr {
				label = " [tomorrow]"
			}

			coloredDate := output.FormatDate(dateStr, today)
			fmt.Printf("%s%s%s\n", indicator, coloredDate, label)
		}
	}

	return nil
}
