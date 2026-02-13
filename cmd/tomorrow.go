package cmd

import (
	"github.com/spf13/cobra"
)

// NewTomorrowCmd creates the tomorrow command
func NewTomorrowCmd(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:   "tomorrow",
		Short: "Open tomorrow's plan file in editor",
		Long:  "Opens the plan file for tomorrow with cursor positioned at tomorrow's entry insertion point",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(*configFlag, *locationFlag, *editorFlag, *editorTypeFlag, *preambleFlag, "tomorrow")
		},
	}
}
