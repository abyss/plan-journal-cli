package output

import (
	"time"

	"github.com/fatih/color"
)

var colorsDisabled = false

// SetColorsDisabled controls whether color output is disabled
func SetColorsDisabled(disabled bool) {
	colorsDisabled = disabled
	if disabled {
		color.NoColor = true
	} else {
		color.NoColor = false
	}
}

// Define reusable color functions for consistent styling
var (
	// Headers and titles
	Header = color.New(color.Bold).SprintFunc() // Just bold, no color
	Bold   = color.New(color.Bold).SprintFunc()

	// Status colors (only the important ones)
	Success = color.New(color.FgGreen).SprintFunc()
	Error   = color.New(color.FgRed).SprintFunc()
	Warning = color.New(color.FgYellow).SprintFunc()
	Info    = func(s string) string { return s } // Default terminal color

	// Content colors - minimal approach
	DateBlue  = func(s string) string { return s }                // Default color for future dates
	DateGreen = color.New(color.FgGreen, color.Bold).SprintFunc() // Today - the most important highlight
	DateGray  = func(s string) string { return s }                // Default color for past dates
	FilePath  = func(s string) string { return s }                // Default color for file paths
	Number    = func(s string) string { return s }                // Default color for numbers
	Highlight = color.New(color.FgGreen).SprintFunc()             // Keep green for [today] label
)

// FormatDate colors a date string based on whether it's today, past, or future
func FormatDate(dateStr string, today time.Time) string {
	todayStr := today.Format("2006-01-02")

	if dateStr == todayStr {
		return DateGreen(dateStr)
	}

	// Parse the date to check if it's in the past
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return DateBlue(dateStr) // Default to blue if we can't parse
	}

	if date.Before(today) {
		return DateGray(dateStr)
	}

	return DateBlue(dateStr)
}

// FormatDayOfWeek returns the day of week in a consistent format
func FormatDayOfWeek(date time.Time) string {
	return Info("(" + date.Format("Mon") + ")")
}
