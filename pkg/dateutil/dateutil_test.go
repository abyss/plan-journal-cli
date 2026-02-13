package dateutil

import (
	"testing"
	"time"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		wantErr   bool
		checkFunc func(time.Time) bool
	}{
		{
			name:    "parse today",
			target:  "today",
			wantErr: false,
			checkFunc: func(result time.Time) bool {
				now := time.Now()
				return result.Year() == now.Year() &&
					result.Month() == now.Month() &&
					result.Day() == now.Day()
			},
		},
		{
			name:    "parse tomorrow",
			target:  "tomorrow",
			wantErr: false,
			checkFunc: func(result time.Time) bool {
				tomorrow := time.Now().AddDate(0, 0, 1)
				return result.Year() == tomorrow.Year() &&
					result.Month() == tomorrow.Month() &&
					result.Day() == tomorrow.Day()
			},
		},
		{
			name:    "parse yesterday",
			target:  "yesterday",
			wantErr: false,
			checkFunc: func(result time.Time) bool {
				yesterday := time.Now().AddDate(0, 0, -1)
				return result.Year() == yesterday.Year() &&
					result.Month() == yesterday.Month() &&
					result.Day() == yesterday.Day()
			},
		},
		{
			name:    "parse YYYY-MM",
			target:  "2026-02",
			wantErr: false,
			checkFunc: func(result time.Time) bool {
				return result.Year() == 2026 && result.Month() == time.February
			},
		},
		{
			name:    "parse YYYY-MM-DD",
			target:  "2026-02-13",
			wantErr: false,
			checkFunc: func(result time.Time) bool {
				return result.Year() == 2026 &&
					result.Month() == time.February &&
					result.Day() == 13
			},
		},
		{
			name:    "invalid format",
			target:  "2026-99",
			wantErr: true,
		},
		{
			name:    "invalid date",
			target:  "2026-02-99",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTarget(tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !tt.checkFunc(result) {
				t.Errorf("ParseTarget() = %v, check failed", result)
			}
		})
	}
}

func TestFormatMonth(t *testing.T) {
	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	result := FormatMonth(date)
	expected := "2026-02"
	if result != expected {
		t.Errorf("FormatMonth() = %v, want %v", result, expected)
	}
}

func TestFormatDate(t *testing.T) {
	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	result := FormatDate(date)
	expected := "2026-02-13"
	if result != expected {
		t.Errorf("FormatDate() = %v, want %v", result, expected)
	}
}

func TestMonthFileName(t *testing.T) {
	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	result := MonthFileName(date)
	expected := "2026-02.plan"
	if result != expected {
		t.Errorf("MonthFileName() = %v, want %v", result, expected)
	}
}

func TestMonthHeader(t *testing.T) {
	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	result := MonthHeader(date)
	expected := "# 2026-02"
	if result != expected {
		t.Errorf("MonthHeader() = %v, want %v", result, expected)
	}
}

func TestDateHeader(t *testing.T) {
	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	result := DateHeader(date)
	expected := "## 2026-02-13"
	if result != expected {
		t.Errorf("DateHeader() = %v, want %v", result, expected)
	}
}

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		expected bool
	}{
		{"valid date", "2026-02-13", true},
		{"invalid month", "2026-13-01", false},
		{"invalid day", "2026-02-99", false},
		{"wrong format", "2026-2-13", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidDate(tt.dateStr)
			if result != tt.expected {
				t.Errorf("IsValidDate(%v) = %v, want %v", tt.dateStr, result, tt.expected)
			}
		})
	}
}

func TestIsValidMonth(t *testing.T) {
	tests := []struct {
		name     string
		monthStr string
		expected bool
	}{
		{"valid month", "2026-02", true},
		{"invalid month", "2026-13", false},
		{"wrong format", "2026-2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidMonth(tt.monthStr)
			if result != tt.expected {
				t.Errorf("IsValidMonth(%v) = %v, want %v", tt.monthStr, result, tt.expected)
			}
		})
	}
}

func TestCompareDates(t *testing.T) {
	tests := []struct {
		name     string
		date1    string
		date2    string
		expected int
	}{
		{"date1 before date2", "2026-02-13", "2026-02-14", -1},
		{"date1 after date2", "2026-02-14", "2026-02-13", 1},
		{"dates equal", "2026-02-13", "2026-02-13", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareDates(tt.date1, tt.date2)
			if result != tt.expected {
				t.Errorf("CompareDates(%v, %v) = %v, want %v", tt.date1, tt.date2, result, tt.expected)
			}
		})
	}
}
