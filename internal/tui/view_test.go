package tui

import (
	"testing"

	"github.com/justin06lee/shaw/internal/run"
)

func TestDisplayGlyph(t *testing.T) {
	cases := []struct {
		name  string
		ch    rune
		state run.CharState
		want  string
	}{
		{"untyped space shows placeholder", ' ', run.Untyped, spaceGlyph},
		{"incorrect space shows placeholder", ' ', run.Incorrect, spaceGlyph},
		{"correct space collapses to blank", ' ', run.Correct, " "},
		{"untyped letter is itself", 'a', run.Untyped, "a"},
		{"correct letter is itself", 'a', run.Correct, "a"},
		{"incorrect letter is itself", 'a', run.Incorrect, "a"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := displayGlyph(c.ch, c.state); got != c.want {
				t.Errorf("displayGlyph(%q, %v) = %q, want %q", c.ch, c.state, got, c.want)
			}
		})
	}
}
