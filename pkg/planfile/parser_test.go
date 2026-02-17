package planfile

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "2026-02.plan")

	content := `# 2026-02

Test preamble text

## 2026-02-13
* First entry
* Second entry

## 2026-02-14
* Third entry
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	pf, err := ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	// Check month header
	if pf.MonthHeader != "# 2026-02" {
		t.Errorf("MonthHeader = %v, want %v", pf.MonthHeader, "# 2026-02")
	}

	// Check preamble
	if pf.Preamble != "Test preamble text" {
		t.Errorf("Preamble = %v, want %v", pf.Preamble, "Test preamble text")
	}

	// Check dates exist
	if len(pf.Dates) != 2 {
		t.Errorf("len(Dates) = %v, want %v", len(pf.Dates), 2)
	}

	// Check date order
	expectedOrder := []string{"2026-02-13", "2026-02-14"}
	if !reflect.DeepEqual(pf.DateOrder, expectedOrder) {
		t.Errorf("DateOrder = %v, want %v", pf.DateOrder, expectedOrder)
	}

	// Check entries for first date
	entries := pf.Dates["2026-02-13"]
	if len(entries) < 2 {
		t.Errorf("len(entries for 2026-02-13) = %v, want at least 2", len(entries))
	}
}

func TestExtractDateContent(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "2026-02.plan")

	content := `# 2026-02

Test preamble

## 2026-02-13
* Entry 1
* Entry 2

## 2026-02-14
* Entry 3
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ExtractDateContent(testFile, "2026-02-13")
	if err != nil {
		t.Fatalf("ExtractDateContent() error = %v", err)
	}

	// Should contain the header
	if !strings.Contains(result, "## 2026-02-13") {
		t.Error("ExtractDateContent() should contain date header")
	}

	// Should contain the entries
	if !strings.Contains(result, "* Entry 1") {
		t.Error("ExtractDateContent() should contain Entry 1")
	}
	if !strings.Contains(result, "* Entry 2") {
		t.Error("ExtractDateContent() should contain Entry 2")
	}

	// Should not contain empty lines between content lines
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		// Skip header line, check content lines
		if i > 0 && line == "" {
			t.Error("ExtractDateContent() should not have empty lines in content")
		}
	}
}

func TestFindInsertionLineForDate(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "2026-02.plan")

	content := `# 2026-02

Test preamble

## 2026-02-13
* Entry 1
* Entry 2

## 2026-02-14
* Entry 3
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	line, err := FindInsertionLineForDate(testFile, "2026-02-13")
	if err != nil {
		t.Fatalf("FindInsertionLineForDate() error = %v", err)
	}

	// Should be after the last entry for that date
	if line <= 5 {
		t.Errorf("FindInsertionLineForDate() = %v, want > 5", line)
	}
}

func TestParseFileWithInvalidDate(t *testing.T) {
	// Create temporary test file with an invalid date
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "2026-02.plan")

	content := `# 2026-02

Test preamble

## 2026-02-13
* Valid entry

## 2026-02-99
* Invalid date entry

## 2026-02-15
* Another valid entry
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Capture stderr to check for warning
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Parse file
	pf, err := ParseFile(testFile)

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	var buf bytes.Buffer
	buf.ReadFrom(r)
	stderrOutput := buf.String()

	// Should not error
	if err != nil {
		t.Fatalf("ParseFile() should not error on invalid date, got error = %v", err)
	}

	// Should have warning in stderr
	if !strings.Contains(stderrOutput, "Warning: Invalid date '2026-02-99'") {
		t.Errorf("Expected warning about invalid date, got: %s", stderrOutput)
	}

	// Should still parse all three date sections (including invalid one)
	if len(pf.Dates) != 3 {
		t.Errorf("len(Dates) = %v, want %v", len(pf.Dates), 3)
	}

	// Should include the invalid date in the parsed data
	if _, exists := pf.Dates["2026-02-99"]; !exists {
		t.Error("Invalid date section should still be parsed")
	}

	// Check date order includes all dates
	expectedOrder := []string{"2026-02-13", "2026-02-99", "2026-02-15"}
	if !reflect.DeepEqual(pf.DateOrder, expectedOrder) {
		t.Errorf("DateOrder = %v, want %v", pf.DateOrder, expectedOrder)
	}
}
