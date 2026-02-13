package planfile

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// Section represents a parsed section of a plan file
type Section struct {
	Header  string   // The header line (e.g., "# 2026-02" or "## 2026-02-13")
	Content []string // Lines of content following the header
}

// PlanFile represents a parsed plan file structure
type PlanFile struct {
	MonthHeader string              // The month header (e.g., "# 2026-02")
	Preamble    string              // The preamble text
	Dates       map[string][]string // Map of date (YYYY-MM-DD) to content lines
	DateOrder   []string            // Ordered list of dates for chronological sorting
	DateHeaders map[string]string   // Full header text for each date (e.g., "## 2026-02-13 - Title")
}

// ParseFile parses a plan file into sections
func ParseFile(filePath string) (*PlanFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pf := &PlanFile{
		Dates:       make(map[string][]string),
		DateOrder:   []string{},
		DateHeaders: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	var currentSection string
	var preambleLines []string
	inPreamble := false

	for scanner.Scan() {
		line := scanner.Text()

		// Month header (# YYYY-MM)
		if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "## ") {
			pf.MonthHeader = line
			inPreamble = true
			continue
		}

		// Date header (## YYYY-MM-DD [optional text])
		if strings.HasPrefix(line, "## ") {
			inPreamble = false
			headerText := strings.TrimPrefix(line, "## ")

			// Extract just the date (YYYY-MM-DD) from the beginning
			datePattern := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})`)
			matches := datePattern.FindStringSubmatch(headerText)

			if len(matches) > 0 {
				date := matches[1]
				currentSection = date
				pf.DateOrder = append(pf.DateOrder, date)
				pf.Dates[date] = []string{}
				pf.DateHeaders[date] = line // Store full header line as-is
			}
			continue
		}

		// Content lines
		if currentSection != "" {
			// Add to current date section
			pf.Dates[currentSection] = append(pf.Dates[currentSection], line)
		} else if inPreamble && line != "" {
			// Collect preamble lines (skip empty lines immediately after month header)
			preambleLines = append(preambleLines, line)
		}
	}

	// Join preamble lines
	if len(preambleLines) > 0 {
		pf.Preamble = strings.Join(preambleLines, "\n")
	}

	return pf, scanner.Err()
}

// FindDateSectionLine finds the line number where a date section starts
// Returns 0 if not found
func FindDateSectionLine(filePath, date string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	datePrefix := "## " + date

	for scanner.Scan() {
		lineNum++
		if strings.HasPrefix(scanner.Text(), datePrefix) {
			return lineNum, nil
		}
	}

	return 0, scanner.Err()
}

// FindInsertionLineForDate finds the line number where new entries should be added for a date
// Returns the line after the last entry for that date
func FindInsertionLineForDate(filePath, date string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	datePrefix := "## " + date
	inTargetSection := false
	lastContentLine := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Found our date header
		if strings.HasPrefix(line, datePrefix) {
			inTargetSection = true
			lastContentLine = lineNum
			continue
		}

		// Found another date header, we're done
		if inTargetSection && strings.HasPrefix(line, "## ") {
			break
		}

		// Track last non-empty line in our section
		if inTargetSection && line != "" {
			lastContentLine = lineNum
		}
	}

	// Return line after last content (insertion point)
	return lastContentLine + 1, scanner.Err()
}

// ExtractDateContent extracts the content lines for a specific date
func ExtractDateContent(filePath, date string) (string, error) {
	pf, err := ParseFile(filePath)
	if err != nil {
		return "", err
	}

	// Check if date exists
	content, ok := pf.Dates[date]
	if !ok {
		return "", nil
	}

	// Get the header (use stored header or construct default)
	header, ok := pf.DateHeaders[date]
	if !ok {
		header = "## " + date
	}

	if len(content) == 0 {
		return header, nil
	}

	// Filter out empty lines and join with header
	var lines []string
	for _, line := range content {
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return header, nil
	}

	return header + "\n" + strings.Join(lines, "\n"), nil
}
