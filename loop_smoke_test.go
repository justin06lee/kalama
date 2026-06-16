package kalama

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// box is a tiny live game used by the end-to-end smoke test: it moves in
// response to arrow keys and quits on esc, drawing a small bright block at its
// position each frame.
type box struct {
	x, y int
}

func (b *box) Update(dt time.Duration, in Input) Action {
	if in.Pressed("esc") {
		return Quit
	}
	if in.Held("right") {
		b.x += 2
	}
	if in.Held("left") {
		b.x -= 2
	}
	if in.Held("up") {
		b.y -= 2
	}
	if in.Held("down") {
		b.y += 2
	}
	return Continue
}

// boxColor is a bright magenta; its truecolor foreground SGR is a distinctive
// marker we can look for in the rendered ANSI output.
var boxColor = Color{R: 255, G: 0, B: 255, A: 255}

// boxColorFG is the exact foreground SGR sequence Canvas.Render emits for a
// pixel of boxColor.
const boxColorFG = "\x1b[38;2;255;0;255m"

func (b *box) Draw(c *Canvas) {
	c.Clear(Color{A: 255}) // opaque black
	for dy := 0; dy < 4; dy++ {
		for dx := 0; dx < 4; dx++ {
			c.Set(b.x+dx, b.y+dy, boxColor)
		}
	}
}

// TestSmokeRunLoopEndToEnd drives the real *model (the same type Run runs)
// through teatest: it feeds arrow-key presses so the loop ticks frames and the
// box moves, asserts the rendered ANSI contains the box's bright color (proving
// Update+Draw+Render ran live), then sends esc and asserts the program exits
// through our Update/Quit path.
func TestSmokeRunLoopEndToEnd(t *testing.T) {
	m := newModel(&box{x: 0, y: 0}, Options{Width: 40, Height: 20, FPS: 60})
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(40, 20))

	// Feed right-arrow presses spread over time. The events land between
	// frames; at 60 FPS frames are ~16ms apart and the held-decay window is
	// 150ms, so "right" stays held across several frames and the box moves.
	for i := 0; i < 8; i++ {
		tm.Send(tea.KeyMsg{Type: tea.KeyRight})
		time.Sleep(20 * time.Millisecond)
	}

	// Assert the program actually rendered frames containing the box. We poll
	// the cumulative output until it carries the bright-color foreground SGR.
	teatest.WaitFor(t, tm.Output(), func(b []byte) bool {
		return len(b) > 0 && bytes.Contains(b, []byte(boxColorFG))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(20*time.Millisecond))

	// Now quit through our own logic: esc -> Pressed("esc") -> Quit.
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})

	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	// Sanity-check the final model ran our code: it quit and the box moved
	// right from its start (x grew) in response to the held key.
	fm := tm.FinalModel(t).(*model)
	if !fm.quit {
		t.Error("final model quit = false, want true after esc")
	}
	if bx := fm.game.(*box); bx.x <= 0 {
		t.Errorf("box.x = %d, want > 0 (box should have moved right while held)", bx.x)
	}
}
