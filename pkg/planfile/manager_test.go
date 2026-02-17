package planfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEnsureMonthFile(t *testing.T) {
	tmpDir := t.TempDir()
	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	preamble := "Test preamble"

	err := EnsureMonthFile(date, tmpDir, preamble)
	if err != nil {
		t.Fatalf("EnsureMonthFile() error = %v", err)
	}

	// Calculate expected file path
	filePath := filepath.Join(tmpDir, "2026-02.plan")

	// Check file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("EnsureMonthFile() did not create file at %v", filePath)
	}

	// Check file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "# 2026-02") {
		t.Error("File does not contain month header")
	}
	if !strings.Contains(contentStr, preamble) {
		t.Error("File does not contain preamble")
	}
}

func TestEnsurePreamble(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "2026-02.plan")

	// Create file without preamble
	content := `# 2026-02

## 2026-02-13
* Entry 1
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add preamble
	preamble := "New preamble"
	if err := EnsurePreamble(testFile, preamble); err != nil {
		t.Fatalf("EnsurePreamble() error = %v", err)
	}

	// Check file now has preamble
	newContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(newContent), preamble) {
		t.Error("EnsurePreamble() did not add preamble")
	}
}

func TestEnsureDateHeader(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "2026-02.plan")

	// Create file with month header
	content := `# 2026-02

Test preamble
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	date := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)
	err := EnsureDateHeader(date, tmpDir)
	if err != nil {
		t.Fatalf("EnsureDateHeader() error = %v", err)
	}

	// Check file now has date header
	newContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(newContent), "## 2026-02-13") {
		t.Error("EnsureDateHeader() did not add date header")
	}
}

func TestFixPlanFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "2026-02.plan")

	// Create file with dates out of order
	content := `# 2026-02

## 2026-02-15
* Entry 1

## 2026-02-13
* Entry 2

## 2026-02-14
* Entry 3
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	preamble := "Test preamble"
	result, err := FixPlanFile("2026-02", tmpDir, preamble)
	if err != nil {
		t.Fatalf("FixPlanFile() error = %v", err)
	}

	// Should report fixes
	if !strings.Contains(result, "preamble") || !strings.Contains(result, "Reordered") {
		t.Errorf("FixPlanFile() result = %v, want to mention fixes", result)
	}

	// Check file is now ordered correctly
	newContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	contentStr := string(newContent)

	// Check order (2026-02-13 should come before 2026-02-14 and 2026-02-15)
	idx13 := strings.Index(contentStr, "## 2026-02-13")
	idx14 := strings.Index(contentStr, "## 2026-02-14")
	idx15 := strings.Index(contentStr, "## 2026-02-15")

	if idx13 == -1 || idx14 == -1 || idx15 == -1 {
		t.Error("FixPlanFile() lost date sections")
	}

	if idx13 > idx14 || idx14 > idx15 {
		t.Error("FixPlanFile() did not reorder dates correctly")
	}

	// Check preamble was added
	if !strings.Contains(contentStr, preamble) {
		t.Error("FixPlanFile() did not add preamble")
	}
}

func TestReadEntries(t *testing.T) {
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

	// Test reading specific date
	result, err := ReadEntries("2026-02-13", tmpDir)
	if err != nil {
		t.Fatalf("ReadEntries() error = %v", err)
	}

	if !strings.Contains(result, "2026-02-13") {
		t.Error("ReadEntries() does not contain date header")
	}
	if !strings.Contains(result, "Entry 1") {
		t.Error("ReadEntries() does not contain entries")
	}

	// Test reading entire month
	result, err = ReadEntries("2026-02", tmpDir)
	if err != nil {
		t.Fatalf("ReadEntries() error = %v", err)
	}

	if !strings.Contains(result, "# 2026-02") {
		t.Error("ReadEntries() for month does not contain month header")
	}
	if !strings.Contains(result, "Test preamble") {
		t.Error("ReadEntries() for month does not contain preamble")
	}

	// Test reading empty date section (header exists but no content)
	emptyFile := filepath.Join(tmpDir, "2026-03.plan")
	emptyContent := `# 2026-03

## 2026-03-01
`
	if err := os.WriteFile(emptyFile, []byte(emptyContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err = ReadEntries("2026-03-01", tmpDir)
	if err != nil {
		t.Fatalf("ReadEntries() error = %v", err)
	}

	// Should show header even with no content
	if !strings.Contains(result, "## 2026-03-01") {
		t.Error("ReadEntries() for empty date should still show header")
	}
	// Result should be just the header with no additional content
	expected := "## 2026-03-01"
	if strings.TrimSpace(result) != expected {
		t.Errorf("ReadEntries() for empty date = %q, want %q", result, expected)
	}

	// Test reading date that doesn't exist in file
	_, err = ReadEntries("2026-03-15", tmpDir)
	if err == nil {
		t.Error("ReadEntries() for non-existent date should return error")
	}
	if !strings.Contains(err.Error(), "no entries found") {
		t.Errorf("ReadEntries() error = %v, want error containing 'no entries found'", err)
	}
}
