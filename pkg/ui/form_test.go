package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestShowLogo(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call ShowLogo
	ShowLogo()

	// Restore stdout and read output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify logo is displayed
	if output == "" {
		t.Error("ShowLogo() produced no output")
	}

	// Check for "PICK YOUR GO" text
	if !strings.Contains(output, "PICK YOUR GO") {
		t.Error("Logo missing 'PICK YOUR GO' text")
	}

	// Check for separator line (= characters)
	if !strings.Contains(output, "=") {
		t.Error("Logo missing separator line")
	}

	// Check for tagline
	if !strings.Contains(output, "Interactive Go Project Scaffolder") {
		t.Error("Logo missing tagline")
	}

	// Verify centering (output should be roughly centered within 80 char width)
	lines := strings.Split(output, "\n")
	var contentLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			contentLines = append(contentLines, trimmed)
		}
	}

	// Should have at least 3 content lines (logo, separator, tagline)
	if len(contentLines) < 3 {
		t.Errorf("Expected at least 3 content lines in logo output, got %d", len(contentLines))
	}

	t.Log("Logo output:")
	t.Log(output)
}

func TestShowLogoTerminalCompatibility(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call ShowLogo
	ShowLogo()

	// Restore stdout and read output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// List of problematic Unicode characters that should NOT be present
	problematicChars := []string{
		"╔", "═", "╗", "╚", "╝", "╠", "╣", "║", // Box drawing
		"│", "─", // Light box drawing
		"┌", "┐", "└", "┘", "├", "┤", "┬", "┴", "┼", // More box drawing
	}

	for _, char := range problematicChars {
		if strings.Contains(output, char) {
			t.Errorf("Logo contains problematic Unicode character: %s", char)
		}
	}

	// Verify only standard ASCII is used
	for i, r := range output {
		if r > 127 {
			t.Errorf("Logo contains non-ASCII character at position %d: %c (U+%04X)", i, r, r)
		}
	}
}

func TestShowLogoHasProperStructure(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call ShowLogo
	ShowLogo()

	// Restore stdout and read output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	lines := strings.Split(output, "\n")

	// Find the logo line
	var logoLineIndex int = -1
	var separatorLineIndex int = -1
	var taglineLineIndex int = -1

	for i, line := range lines {
		if strings.Contains(line, "PICK YOUR GO") {
			logoLineIndex = i
		}
		if strings.Contains(line, "======") {
			separatorLineIndex = i
		}
		if strings.Contains(line, "Interactive Go Project Scaffolder") {
			taglineLineIndex = i
		}
	}

	// Verify structure: logo -> separator -> tagline (in that order)
	if logoLineIndex == -1 {
		t.Error("Could not find logo line")
	}
	if separatorLineIndex == -1 {
		t.Error("Could not find separator line")
	}
	if taglineLineIndex == -1 {
		t.Error("Could not find tagline line")
	}

	// Check order
	if logoLineIndex != -1 && separatorLineIndex != -1 && separatorLineIndex <= logoLineIndex {
		t.Error("Separator should come after logo text")
	}
	if taglineLineIndex != -1 && separatorLineIndex != -1 && taglineLineIndex <= separatorLineIndex {
		t.Error("Tagline should come after separator")
	}
}
