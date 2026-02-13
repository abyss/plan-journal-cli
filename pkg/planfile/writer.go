package planfile

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/abyss/plan-journal-cli/pkg/dateutil"
)

// WritePlanFile writes a PlanFile structure to disk
func WritePlanFile(filePath string, pf *PlanFile) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Build file content
	var lines []string

	// Month header
	lines = append(lines, pf.MonthHeader)
	lines = append(lines, "")

	// Preamble
	if pf.Preamble != "" {
		lines = append(lines, pf.Preamble)
		lines = append(lines, "")
	}

	// Sort dates chronologically
	sortedDates := make([]string, len(pf.DateOrder))
	copy(sortedDates, pf.DateOrder)
	sort.Slice(sortedDates, func(i, j int) bool {
		return dateutil.CompareDates(sortedDates[i], sortedDates[j]) < 0
	})

	// Write date sections
	for _, date := range sortedDates {
		lines = append(lines, "## "+date)

		// Add content lines
		if content, ok := pf.Dates[date]; ok {
			lines = append(lines, content...)
		}

		lines = append(lines, "")
	}

	// Write to file
	content := strings.Join(lines, "\n")
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
