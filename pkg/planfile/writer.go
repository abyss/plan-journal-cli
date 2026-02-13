package planfile

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/abyss/plan-journal-cli/pkg/dateutil"
)

// GenerateFileContent generates the file content from a PlanFile structure
func GenerateFileContent(pf *PlanFile) (string, error) {
	// Build file content
	var lines []string

	// Month header
	lines = append(lines, pf.MonthHeader)
	lines = append(lines, "")

	// Preamble
	if pf.Preamble != "" {
		lines = append(lines, pf.Preamble)
		lines = append(lines, "", "")
	}

	// Sort dates chronologically
	sortedDates := make([]string, len(pf.DateOrder))
	copy(sortedDates, pf.DateOrder)
	sort.Slice(sortedDates, func(i, j int) bool {
		return dateutil.CompareDates(sortedDates[i], sortedDates[j]) < 0
	})

	// Write date sections
	for i, date := range sortedDates {
		lines = append(lines, "## "+date)

		// Add content lines (trim trailing empty lines first)
		if content, ok := pf.Dates[date]; ok {
			// Trim trailing empty lines from content
			trimmedContent := trimTrailingEmptyLines(content)
			lines = append(lines, trimmedContent...)
		}

		// Add two empty lines between date sections (but not after the last one)
		if i < len(sortedDates)-1 {
			lines = append(lines, "", "")
		}
	}

	// Return content with trailing newline
	content := strings.Join(lines, "\n") + "\n"
	return content, nil
}

// WritePlanFile writes a PlanFile structure to disk
func WritePlanFile(filePath string, pf *PlanFile) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Generate content
	content, err := GenerateFileContent(pf)
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, []byte(content), 0644)
}

// AppendDateSection appends a new date section to an existing file
func AppendDateSection(filePath, date string) error {
	// Read existing content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Add new date section
	newSection := "\n## " + date + "\n"

	// Append to file
	return os.WriteFile(filePath, append(content, []byte(newSection)...), 0644)
}

// EnsureDirectory ensures a directory exists
func EnsureDirectory(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// trimTrailingEmptyLines removes trailing empty lines from a slice of strings
func trimTrailingEmptyLines(lines []string) []string {
	// Find the last non-empty line
	lastNonEmpty := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] != "" {
			lastNonEmpty = i
			break
		}
	}

	if lastNonEmpty == -1 {
		// All lines are empty
		return []string{}
	}

	return lines[:lastNonEmpty+1]
}
