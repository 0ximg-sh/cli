package cmd

import (
	"bytes"
	"strings"
	"testing"

	"0ximg.sh/cli/internal/models"
)

func TestParseLineRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		start   int
		end     int
		wantErr bool
	}{
		{name: "dash", input: "10-20", start: 10, end: 20},
		{name: "colon", input: "10:20", start: 10, end: 20},
		{name: "dots", input: "10..20", start: 10, end: 20},
		{name: "reverse", input: "20-10", wantErr: true},
		{name: "invalid", input: "10", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := parseLineRange(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseLineRange(%q) returned error: %v", tt.input, err)
			}
			if start != tt.start || end != tt.end {
				t.Fatalf("parseLineRange(%q) = (%d, %d), want (%d, %d)", tt.input, start, end, tt.start, tt.end)
			}
		})
	}
}

func TestApplyLineRange(t *testing.T) {
	code := "l1\nl2\nl3\nl4\nl5\n"

	sliced, highlights, offset, err := applyLineRange(code, "2-4", "2;3-4", false)
	if err != nil {
		t.Fatalf("applyLineRange returned error: %v", err)
	}
	if sliced != "l2\nl3\nl4" {
		t.Fatalf("unexpected sliced code: %q", sliced)
	}
	if highlights != "1;2-3" {
		t.Fatalf("unexpected highlights: %q", highlights)
	}
	if offset != 1 {
		t.Fatalf("unexpected offset: %d", offset)
	}
}

func TestApplyLineRangeKeepsExplicitOffset(t *testing.T) {
	code := "l1\nl2\nl3\n"

	_, _, offset, err := applyLineRange(code, "2-3", "", true)
	if err != nil {
		t.Fatalf("applyLineRange returned error: %v", err)
	}
	if offset != 0 {
		t.Fatalf("unexpected offset: %d", offset)
	}
}

func TestApplyLineRangeRejectsHighlightOutsideRange(t *testing.T) {
	code := "l1\nl2\nl3\nl4\n"

	_, _, _, err := applyLineRange(code, "2-3", "1", false)
	if err == nil {
		t.Fatal("expected error for highlight outside selected range")
	}
}

func TestListThemesFlag(t *testing.T) {
	oldListThemes := listThemes
	listThemes = true
	t.Cleanup(func() {
		listThemes = oldListThemes
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--list-themes"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Dracula\n") {
		t.Fatalf("expected Dracula in output, got %q", output)
	}
	if !strings.Contains(output, "Nord\n") {
		t.Fatalf("expected Nord in output, got %q", output)
	}
}

func TestDefaultPadValuesAreApplied(t *testing.T) {
	req := models.RenderRequest{}

	setIntValue(&req.PadHoriz, defaultPadHoriz)
	setIntValue(&req.PadVert, defaultPadVert)

	if req.PadHoriz == nil || *req.PadHoriz != defaultPadHoriz {
		t.Fatalf("unexpected horizontal pad: %#v", req.PadHoriz)
	}
	if req.PadVert == nil || *req.PadVert != defaultPadVert {
		t.Fatalf("unexpected vertical pad: %#v", req.PadVert)
	}
}

func TestLineOffsetIsIncludedWhenLinesSelected(t *testing.T) {
	req := models.RenderRequest{}
	lineRange = "4-11"
	lineOffset = 3
	t.Cleanup(func() {
		lineRange = ""
		lineOffset = 0
	})

	if strings.TrimSpace(lineRange) != "" {
		setIntValue(&req.LineOffset, lineOffset)
	} else {
		t.Fatal("expected lineRange to be set for test")
	}

	if req.LineOffset == nil || *req.LineOffset != 3 {
		t.Fatalf("unexpected line offset: %#v", req.LineOffset)
	}
}
