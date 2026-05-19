// Package tui renders the shaw terminal interface with bubbletea.
package tui

// Line is a half-open rune-index range [Start,End) into the target text.
type Line struct {
	Start, End int
}

// WrapLines breaks text into lines no wider than width, splitting on spaces.
// A word longer than width occupies its own (over-wide) line.
func WrapLines(text []rune, width int) []Line {
	if width < 1 {
		width = 1
	}
	var lines []Line
	lineStart := 0
	lastSpace := -1
	for i := 0; i < len(text); i++ {
		if text[i] == ' ' {
			lastSpace = i
		}
		if i-lineStart >= width {
			if lastSpace > lineStart {
				lines = append(lines, Line{lineStart, lastSpace})
				lineStart = lastSpace + 1
			} else {
				lines = append(lines, Line{lineStart, i})
				lineStart = i
			}
			lastSpace = -1
		}
	}
	if lineStart < len(text) || len(lines) == 0 {
		lines = append(lines, Line{lineStart, len(text)})
	}
	return lines
}

// LineOfCursor returns the index of the line containing cursor (an index into
// the text). A cursor at the very end maps to the last line.
func LineOfCursor(lines []Line, cursor int) int {
	for i, ln := range lines {
		if cursor >= ln.Start && cursor < ln.End {
			return i
		}
	}
	if len(lines) == 0 {
		return 0
	}
	return len(lines) - 1
}

// Viewport returns the start line index and count (<=3) for a 3-line window
// that keeps the cursor's line centered, clamped to the available lines.
func Viewport(lines []Line, cursorLine int) (start, count int) {
	const window = 3
	if len(lines) <= window {
		return 0, len(lines)
	}
	start = cursorLine - 1
	if start < 0 {
		start = 0
	}
	if start+window > len(lines) {
		start = len(lines) - window
	}
	return start, window
}
