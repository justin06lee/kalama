package run

import (
	"testing"
	"time"
)

func TestTypeMarksCharStates(t *testing.T) {
	r := New(ModeZen, 0)
	r.AppendWords([]string{"go"})
	r.Type('g')
	r.Type('x') // wrong
	states := r.States()
	if states[0] != Correct {
		t.Errorf("char 0: got %v, want Correct", states[0])
	}
	if states[1] != Incorrect {
		t.Errorf("char 1: got %v, want Incorrect", states[1])
	}
	if r.Cursor() != 2 {
		t.Errorf("cursor: got %d, want 2", r.Cursor())
	}
}

func TestBackspaceResetsChar(t *testing.T) {
	r := New(ModeZen, 0)
	r.AppendWords([]string{"go"})
	r.Type('g')
	r.Backspace()
	if r.Cursor() != 0 {
		t.Errorf("cursor: got %d, want 0", r.Cursor())
	}
	if r.States()[0] != Untyped {
		t.Errorf("char 0: got %v, want Untyped", r.States()[0])
	}
}

func TestBackspaceAtStartIsNoop(t *testing.T) {
	r := New(ModeZen, 0)
	r.AppendWords([]string{"go"})
	r.Backspace()
	if r.Cursor() != 0 {
		t.Errorf("cursor: got %d, want 0", r.Cursor())
	}
}

func TestTypePastEndIsNoop(t *testing.T) {
	r := New(ModeZen, 0)
	r.AppendWords([]string{"a"})
	r.Type('a')
	r.Type('b') // past end
	if r.Cursor() != 1 {
		t.Errorf("cursor: got %d, want 1", r.Cursor())
	}
}

func TestAppendWordsJoinsWithSpaces(t *testing.T) {
	r := New(ModeZen, 0)
	r.AppendWords([]string{"one", "two"})
	if string(r.Text()) != "one two" {
		t.Errorf("text: got %q, want %q", string(r.Text()), "one two")
	}
}

// fakeClock returns successive fixed times for the injectable Now function.
func fakeClock(times ...time.Time) func() time.Time {
	i := 0
	return func() time.Time {
		t := times[i]
		if i < len(times)-1 {
			i++
		}
		return t
	}
}

func TestGoalWordsReachedWhenTextFullyTyped(t *testing.T) {
	r := New(ModeWords, 1)
	r.AppendWords([]string{"hi"})
	if r.GoalReached() {
		t.Fatal("goal reached before typing")
	}
	r.Type('h')
	r.Type('i')
	if !r.GoalReached() {
		t.Fatal("goal not reached after typing all chars")
	}
}

func TestGoalTimeReachedWhenTargetSecondsElapse(t *testing.T) {
	base := time.Unix(0, 0)
	r := New(ModeTime, 30)
	r.Now = fakeClock(base, base.Add(30*time.Second))
	r.AppendWords([]string{"abc"})
	r.Type('a') // started at base
	if !r.GoalReached() {
		t.Fatal("goal not reached at +30s")
	}
}

func TestGoalTimeNotReachedEarly(t *testing.T) {
	base := time.Unix(0, 0)
	r := New(ModeTime, 30)
	r.Now = fakeClock(base, base.Add(5*time.Second))
	r.AppendWords([]string{"abc"})
	r.Type('a')
	if r.GoalReached() {
		t.Fatal("goal reached too early at +5s")
	}
}

func TestGoalZenNeverReached(t *testing.T) {
	r := New(ModeZen, 0)
	r.AppendWords([]string{"a"})
	r.Type('a')
	if r.GoalReached() {
		t.Fatal("zen mode should never auto-finish")
	}
}

func TestDurationIsLastKeystrokeTime(t *testing.T) {
	base := time.Unix(0, 0)
	r := New(ModeZen, 0)
	r.Now = fakeClock(base, base.Add(2*time.Second))
	r.AppendWords([]string{"ab"})
	r.Type('a') // sets started=base
	r.Type('b') // logged at +2s
	if r.Duration() != 2*time.Second {
		t.Fatalf("duration: got %v, want 2s", r.Duration())
	}
}
