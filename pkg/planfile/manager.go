package planfile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/abyss/plan-journal-cli/pkg/dateutil"
)

// EnsureMonthFile ensures a month file exists with header and preamble
// Returns the file path
func EnsureMonthFile(date time.Time, plansDir, preamble string) (string, error) {
	// Ensure directory exists
	if err := EnsureDirectory(plansDir); err != nil {
		return "", fmt.Errorf("failed to create plans directory: %w", err)
	}

	// Build file path
	fileName := dateutil.MonthFileName(date)
	filePath := filepath.Join(plansDir, fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, ensure it has a preamble
		return filePath, EnsurePreamble(filePath, preamble)
	}

	// Create new file with month header and preamble
	monthHeader := dateutil.MonthHeader(date)
	content := monthHeader + "\n\n" + preamble + "\n"

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create month file: %w", err)
	}

	return filePath, nil
}

// EnsurePreamble ensures a file has the correct preamble
func EnsurePreamble(filePath, preamble string) error {
	pf, err := ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// If preamble matches, nothing to do
	if strings.TrimSpace(pf.Preamble) == strings.TrimSpace(preamble) {
		return nil
	}

	// Update preamble
	pf.Preamble = preamble
	return WritePlanFile(filePath, pf)
}

// EnsureDateHeader ensures a date header exists in the file
// Inserts it in chronological order if it doesn't exist
func EnsureDateHeader(date time.Time, plansDir string) (string, error) {
	// Ensure month file exists first
	filePath := filepath.Join(plansDir, dateutil.MonthFileName(date))
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("month file does not exist: %w", err)
	}

	// Parse file
	pf, err := ParseFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	// Check if date already exists
	dateStr := dateutil.FormatDate(date)
	if _, exists := pf.Dates[dateStr]; exists {
		return filePath, nil
	}

	// Add new date section
	pf.Dates[dateStr] = []string{}
	pf.DateOrder = append(pf.DateOrder, dateStr)

	// Write updated file (will be sorted chronologically)
	return filePath, WritePlanFile(filePath, pf)
}

// FindInsertionPoint returns the file path and line number for inserting new entries
func FindInsertionPoint(date time.Time, plansDir string) (string, int, error) {
	filePath := filepath.Join(plansDir, dateutil.MonthFileName(date))
	dateStr := dateutil.FormatDate(date)

	lineNum, err := FindInsertionLineForDate(filePath, dateStr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to find insertion point: %w", err)
	}

	return filePath, lineNum, nil
}

// ReadEntries reads and returns entries based on the target
// target can be "today", "YYYY-MM", or "YYYY-MM-DD"
func ReadEntries(target, plansDir string) (string, error) {
	// Parse target
	date, err := dateutil.ParseTarget(target)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(plansDir, dateutil.MonthFileName(date))

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("no plan file found for %s", target)
	}

	// For month target, return entire file
	if len(target) == 7 || (target == "today" && false) { // YYYY-MM format
		// Actually, let's check the target length more carefully
		parsedTarget := target
		if target == "today" {
			parsedTarget = dateutil.FormatDate(date)
		}

		if len(parsedTarget) == 7 { // YYYY-MM
			content, err := os.ReadFile(filePath)
			if err != nil {
				return "", err
			}
			return string(content), nil
		}
	}

	// For specific date, extract the date section
	dateStr := dateutil.FormatDate(date)
	content, err := ExtractDateContent(filePath, dateStr)
	if err != nil {
		return "", err
	}

	// If no content returned, date section doesn't exist
	if content == "" {
		return "", fmt.Errorf("no entries found for %s", dateStr)
	}

	return content, nil
}

// FixPlanFile repairs a plan file by reordering dates and updating preamble
func FixPlanFile(target, plansDir, preamble string) (string, error) {
	// Parse target to get file path
	date, err := dateutil.ParseTarget(target)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(plansDir, dateutil.MonthFileName(date))

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("no plan file found for %s", target)
	}

	// Parse file
	pf, err := ParseFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	// Track what was fixed
	fixes := []string{}

	// Check if preamble needs updating
	if strings.TrimSpace(pf.Preamble) != strings.TrimSpace(preamble) {
		pf.Preamble = preamble
		fixes = append(fixes, "Updated preamble")
	}

	// Check if dates need reordering
	sortedDates := make([]string, len(pf.DateOrder))
	copy(sortedDates, pf.DateOrder)
	sort.Slice(sortedDates, func(i, j int) bool {
		return dateutil.CompareDates(sortedDates[i], sortedDates[j]) < 0
	})

	needsReordering := false
	for i := range pf.DateOrder {
		if pf.DateOrder[i] != sortedDates[i] {
			needsReordering = true
			break
		}
	}

	if needsReordering {
		fixes = append(fixes, "Reordered date sections chronologically")
	}

	// Write fixed file
	if err := WritePlanFile(filePath, pf); err != nil {
		return "", fmt.Errorf("failed to write fixed file: %w", err)
	}

	if len(fixes) == 0 {
		return "No issues found", nil
	}

	return "Fixed: " + strings.Join(fixes, ", "), nil
}
