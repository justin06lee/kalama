package tui

import "testing"

func TestWrapLinesBreaksOnWordBoundary(t *testing.T) {
	// "aaa bbb ccc" width 7 => "aaa bbb" (7) then "ccc".
	text := []rune("aaa bbb ccc")
	lines := WrapLines(text, 7)
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %+v", len(lines), lines)
	}
	if string(text[lines[0].Start:lines[0].End]) != "aaa bbb" {
		t.Errorf("line 0: got %q", string(text[lines[0].Start:lines[0].End]))
	}
	if string(text[lines[1].Start:lines[1].End]) != "ccc" {
		t.Errorf("line 1: got %q", string(text[lines[1].Start:lines[1].End]))
	}
}

func TestWrapLinesSingleLineFits(t *testing.T) {
	lines := WrapLines([]rune("short"), 40)
	if len(lines) != 1 {
		t.Fatalf("got %d lines, want 1", len(lines))
	}
}

func TestLineOfCursor(t *testing.T) {
	lines := WrapLines([]rune("aaa bbb ccc"), 7) // line0 [0,7), line1 [8,11)
	if got := LineOfCursor(lines, 2); got != 0 {
		t.Errorf("cursor 2: got line %d, want 0", got)
	}
	if got := LineOfCursor(lines, 9); got != 1 {
		t.Errorf("cursor 9: got line %d, want 1", got)
	}
	if got := LineOfCursor(lines, 11); got != 1 {
		t.Errorf("cursor at end: got line %d, want 1", got)
	}
}

func TestViewportReturnsThreeLinesCentered(t *testing.T) {
	lines := []Line{{0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}}
	start, count := Viewport(lines, 2) // cursor on line 2 => window 1..3
	if start != 1 || count != 3 {
		t.Errorf("got start=%d count=%d, want 1,3", start, count)
	}
}

func TestViewportClampsAtTop(t *testing.T) {
	lines := []Line{{0, 1}, {1, 2}, {2, 3}}
	start, count := Viewport(lines, 0)
	if start != 0 || count != 3 {
		t.Errorf("got start=%d count=%d, want 0,3", start, count)
	}
}
