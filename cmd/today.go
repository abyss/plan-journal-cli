package cmd

import (
	"github.com/spf13/cobra"
)

// NewTodayCmd creates the today command
func NewTodayCmd(configFlag, locationFlag, editorFlag, editorTypeFlag, preambleFlag *string) *cobra.Command {
	return &cobra.Command{
		Use:   "today",
		Short: "Open today's plan file in editor",
		Long:  "Opens the current month's plan file with cursor positioned at today's entry insertion point",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(*configFlag, *locationFlag, *editorFlag, *editorTypeFlag, *preambleFlag, "today")
		},
	}
}
