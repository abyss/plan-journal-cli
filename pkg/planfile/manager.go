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
func EnsureMonthFile(date time.Time, plansDir, preamble string) error {
	// Ensure directory exists
	if err := EnsureDirectory(plansDir); err != nil {
		return fmt.Errorf("failed to create plans directory: %w", err)
	}

	// Build file path
	fileName := dateutil.MonthFileName(date)
	filePath := filepath.Join(plansDir, fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, ensure it has a preamble
		return EnsurePreamble(filePath, preamble)
	}

	// Create new file with month header and preamble
	monthHeader := dateutil.MonthHeader(date)
	content := monthHeader + "\n\n" + preamble + "\n"

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create month file: %w", err)
	}

	return nil
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
func EnsureDateHeader(date time.Time, plansDir string) error {
	// Ensure month file exists first
	filePath := filepath.Join(plansDir, dateutil.MonthFileName(date))
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("failed to access month file: %w", err)
	}

	// Parse file
	pf, err := ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	// Check if date already exists
	dateStr := dateutil.FormatDate(date)
	if _, exists := pf.Dates[dateStr]; exists {
		return nil
	}

	// Add new date section
	pf.Dates[dateStr] = []string{}
	pf.DateOrder = append(pf.DateOrder, dateStr)

	// Write updated file (will be sorted chronologically)
	return WritePlanFile(filePath, pf)
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
// target can be "yesterday", "today", "tomorrow", "YYYY-MM", or "YYYY-MM-DD"
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

	// If target is a month format (YYYY-MM), return entire file
	if dateutil.IsValidMonth(target) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file: %w", err)
		}
		return string(content), nil
	}

	// For specific date, extract the date section
	dateStr := dateutil.FormatDate(date)
	content, err := ExtractDateContent(filePath, dateStr)
	if err != nil {
		return "", fmt.Errorf("failed to extract date content: %w", err)
	}

	// If no content returned, date section doesn't exist
	if content == "" {
		return "", fmt.Errorf("no entries found for %s", dateStr)
	}

	return content, nil
}

// resolveTargetToFilePath resolves a target (date string or file path) to an absolute file path
// target can be:
// - A date string (YYYY-MM, YYYY-MM-DD, today, yesterday, tomorrow)
// - An absolute file path
// - A relative file path
// - A filename (looked up in plansDir)
func resolveTargetToFilePath(target, plansDir string) (string, error) {
	// First, try to parse as a date
	date, err := dateutil.ParseTarget(target)
	if err == nil {
		// Valid date, construct file path
		filePath := filepath.Join(plansDir, dateutil.MonthFileName(date))
		if _, statErr := os.Stat(filePath); statErr != nil {
			return "", fmt.Errorf("no plan file found for %s", target)
		}
		return filePath, nil
	}

	// Not a date, treat as a file path
	var filePath string

	// Check if it's an absolute path
	if filepath.IsAbs(target) {
		filePath = target
	} else {
		// Check if the relative path exists from current directory
		if _, statErr := os.Stat(target); statErr == nil {
			absPath, absErr := filepath.Abs(target)
			if absErr == nil {
				filePath = absPath
			} else {
				filePath = target
			}
		} else {
			// Try as a filename in the plans directory
			filePath = filepath.Join(plansDir, target)
		}
	}

	// Verify the file exists
	if _, statErr := os.Stat(filePath); statErr != nil {
		return "", fmt.Errorf("file not found: %s (tried as date, absolute path, relative path, and filename in plans directory)", target)
	}

	return filePath, nil
}

// FormatPlanFile formats a plan file by reordering dates and updating preamble
// target can be a date string (YYYY-MM, YYYY-MM-DD, today, etc.) or a file path
func FormatPlanFile(target, plansDir, preamble string) (string, error) {
	// Resolve target to file path
	filePath, err := resolveTargetToFilePath(target, plansDir)
	if err != nil {
		return "", err
	}

	// Read original file content
	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Parse file
	pf, err := ParseFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	// Track what was changed
	changes := []string{}

	// Check if preamble needs updating
	if strings.TrimSpace(pf.Preamble) != strings.TrimSpace(preamble) {
		pf.Preamble = preamble
		changes = append(changes, "Updated preamble")
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
		changes = append(changes, "Reordered date sections chronologically")
	}

	// Generate what the formatted file would look like
	formattedContent := GenerateFileContent(pf)

	// Check if formatting/spacing needs changes
	if string(originalContent) != formattedContent {
		if !needsReordering && strings.TrimSpace(pf.Preamble) == strings.TrimSpace(preamble) {
			changes = append(changes, "Formatted spacing")
		}
	}

	// Only write if there are changes
	if string(originalContent) != formattedContent {
		if err := os.WriteFile(filePath, []byte(formattedContent), 0644); err != nil {
			return "", fmt.Errorf("failed to write formatted file: %w", err)
		}
	}

	if len(changes) == 0 {
		return "No changes needed", nil
	}

	return "Changes: " + strings.Join(changes, ", "), nil
}
