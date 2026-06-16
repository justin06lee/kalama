package kalama

import (
	"testing"
	"time"
)

func ms(n int) time.Duration { return time.Duration(n) * time.Millisecond }

func TestPressMakesKeyHeldAndPressed(t *testing.T) {
	tr := newInputTracker()
	tr.press("a", ms(0))
	in := tr.snapshot(ms(10))
	if !in.Held("a") {
		t.Error("Held(a) = false, want true within decay window")
	}
	if !in.Pressed("a") {
		t.Error("Pressed(a) = false, want true on first appearance")
	}
	if in.Released("a") {
		t.Error("Released(a) = true, want false")
	}
}

func TestPressedOnlyOnFirstFrame(t *testing.T) {
	tr := newInputTracker()
	tr.press("a", ms(0))
	_ = tr.snapshot(ms(10)) // first frame: pressed
	tr.press("a", ms(20))   // key-repeat keeps it alive
	in := tr.snapshot(ms(30))
	if !in.Held("a") {
		t.Error("Held(a) = false, want true")
	}
	if in.Pressed("a") {
		t.Error("Pressed(a) = true, want false on second held frame")
	}
}

func TestReleasedFiresOnceAfterDecay(t *testing.T) {
	tr := newInputTracker()
	tr.press("a", ms(0))
	_ = tr.snapshot(ms(10))    // held
	in := tr.snapshot(ms(300)) // far past decay window -> released
	if in.Held("a") {
		t.Error("Held(a) = true, want false after decay")
	}
	if !in.Released("a") {
		t.Error("Released(a) = false, want true on the frame it decays")
	}
	in2 := tr.snapshot(ms(310))
	if in2.Released("a") {
		t.Error("Released(a) = true again, want false (fires once)")
	}
}
