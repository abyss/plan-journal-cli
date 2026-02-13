package dateutil

import (
	"fmt"
	"time"
)

// ParseTarget parses a date target string into a time.Time
// Accepts: "yesterday", "today", "tomorrow", "YYYY-MM", or "YYYY-MM-DD"
func ParseTarget(target string) (time.Time, error) {
	if target == "today" {
		return time.Now(), nil
	}

	if target == "tomorrow" {
		return time.Now().AddDate(0, 0, 1), nil
	}

	if target == "yesterday" {
		return time.Now().AddDate(0, 0, -1), nil
	}

	// Try YYYY-MM-DD format
	t, err := time.Parse("2006-01-02", target)
	if err == nil {
		return t, nil
	}

	// Try YYYY-MM format
	t, err = time.Parse("2006-01", target)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s (expected 'yesterday', 'today', 'tomorrow', 'YYYY-MM', or 'YYYY-MM-DD')", target)
}

// FormatMonth returns the YYYY-MM format for a given time
func FormatMonth(t time.Time) string {
	return t.Format("2006-01")
}

// FormatDate returns the YYYY-MM-DD format for a given time
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// MonthFileName returns the filename for a given month (YYYY-MM.plan)
func MonthFileName(t time.Time) string {
	return FormatMonth(t) + ".plan"
}

// MonthHeader returns the markdown header for a month (# YYYY-MM)
func MonthHeader(t time.Time) string {
	return "# " + FormatMonth(t)
}

// DateHeader returns the markdown header for a date (## YYYY-MM-DD)
func DateHeader(t time.Time) string {
	return "## " + FormatDate(t)
}

// IsValidDate checks if a date string is valid
func IsValidDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// IsValidMonth checks if a month string is valid
func IsValidMonth(monthStr string) bool {
	_, err := time.Parse("2006-01", monthStr)
	return err == nil
}

// CompareDates compares two date strings chronologically
// Returns -1 if date1 < date2, 0 if equal, 1 if date1 > date2
func CompareDates(date1, date2 string) int {
	t1, err1 := time.Parse("2006-01-02", date1)
	t2, err2 := time.Parse("2006-01-02", date2)

	if err1 != nil || err2 != nil {
		return 0
	}

	if t1.Before(t2) {
		return -1
	} else if t1.After(t2) {
		return 1
	}
	return 0
}
