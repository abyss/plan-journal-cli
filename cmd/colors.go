package cmd

import (
	"fmt"

	"github.com/abyss/plan-journal-cli/pkg/output"
	"github.com/spf13/cobra"
)

// NewColorsCmd creates a hidden command to test color output
func NewColorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "colors",
		Hidden: true,
		Short:  "Display color palette test",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(output.Header("Color Palette Test"))
			fmt.Println()

			fmt.Println("Headers and Titles:")
			fmt.Printf("  Header: %s\n", output.Header("This is a header"))
			fmt.Printf("  Bold: %s\n", output.Bold("This is bold text"))
			fmt.Println()

			fmt.Println("Status Colors:")
			fmt.Printf("  Success: %s\n", output.Success("Operation successful"))
			fmt.Printf("  Error: %s\n", output.Error("An error occurred"))
			fmt.Printf("  Warning: %s\n", output.Warning("Warning message"))
			fmt.Printf("  Info: %s\n", output.Info("Informational text"))
			fmt.Println()

			fmt.Println("Date Colors:")
			fmt.Printf("  DateBlue (future): %s\n", output.DateBlue("2026-03-15"))
			fmt.Printf("  DateGreen (today): %s\n", output.DateGreen("2026-02-19"))
			fmt.Printf("  DateGray (past): %s\n", output.DateGray("2026-01-10"))
			fmt.Println()

			fmt.Println("Content Colors:")
			fmt.Printf("  FilePath: %s\n", output.FilePath("/path/to/file.txt"))
			fmt.Printf("  Number: %s\n", output.Number("42"))
			fmt.Printf("  Highlight: %s\n", output.Highlight("[important]"))
			fmt.Println()

			fmt.Println("List Command Example:")
			fmt.Printf("  %s\n", output.Header("2026-02:"))
			fmt.Printf("    %s\n", output.DateGray("2026-02-11"))
			fmt.Printf("    %s %s\n", output.DateGray("2026-02-18"), "[yesterday]")
			fmt.Printf("  %s%s %s\n", output.DateGreen("> "), output.DateGreen("2026-02-19"), output.Highlight("[today]"))
			fmt.Printf("    %s %s\n", output.DateBlue("2026-02-20"), "[tomorrow]")
			fmt.Println()

			fmt.Println("Config Command Example:")
			fmt.Printf("  %s: %s\n", output.Bold("Plans Directory"), output.FilePath("~/plans"))
			fmt.Printf("    %s: %s\n", output.Info("Source"), "config file")
			fmt.Printf("  %s: %s %s\n", output.Bold("Config File"), output.FilePath("/path/to/config"), output.Success("(exists)"))

			return nil
		},
	}
}
