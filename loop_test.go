package shaw

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// fakeGame records Update calls and returns a configurable Action.
type fakeGame struct {
	updates int
	lastIn  Input
	action  Action
	drawn   bool
}

func (g *fakeGame) Update(dt time.Duration, in Input) Action {
	g.updates++
	g.lastIn = in
	return g.action
}
func (g *fakeGame) Draw(c *Canvas) { g.drawn = true }

func TestFrameAdvancesGameAndContinues(t *testing.T) {
	g := &fakeGame{action: Continue}
	m := newModel(g, Options{Width: 4, Height: 4})
	next, cmd := m.Update(frameMsg(time.Now()))
	m = next.(*model)
	if g.updates != 1 {
		t.Errorf("Update calls = %d, want 1", g.updates)
	}
	if m.quit {
		t.Error("quit = true, want false on Continue")
	}
	if cmd == nil {
		t.Error("cmd = nil, want a follow-up tick command on Continue")
	}
}

func TestFrameQuitStopsLoop(t *testing.T) {
	g := &fakeGame{action: Quit}
	m := newModel(g, Options{Width: 4, Height: 4})
	next, _ := m.Update(frameMsg(time.Now()))
	if !next.(*model).quit {
		t.Error("quit = false, want true when game returns Quit")
	}
}

func TestCtrlCQuits(t *testing.T) {
	g := &fakeGame{action: Continue}
	m := newModel(g, Options{Width: 4, Height: 4})
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if !next.(*model).quit {
		t.Error("quit = false, want true on Ctrl+C")
	}
}

func TestKeyEventBecomesHeldNextFrame(t *testing.T) {
	g := &fakeGame{action: Continue}
	m := newModel(g, Options{Width: 4, Height: 4})
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m.Update(frameMsg(time.Now()))
	if !g.lastIn.Held("a") {
		t.Error(`Held("a") = false, want true after key event then frame`)
	}
}

func TestViewRendersCanvas(t *testing.T) {
	g := &fakeGame{action: Continue}
	m := newModel(g, Options{Width: 1, Height: 2})
	out := m.View()
	if !g.drawn {
		t.Error("game.Draw was not called by View")
	}
	if out == "" {
		t.Error("View returned empty string")
	}
}
