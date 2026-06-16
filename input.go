package kalama

import "time"

// Key is a normalized key name, e.g. "left", "a", "space", "esc". It matches
// bubbletea's KeyMsg.String() values.
type Key string

// decayWindow is how long after a key's last press/repeat event the key is
// still considered held. OS key-repeat re-sends a held key faster than this, so
// a genuinely held key stays Held; once repeats stop, it decays to released.
const decayWindow = 150 * time.Millisecond

// Input is an immutable per-frame snapshot of key state.
type Input struct {
	held     map[Key]bool
	pressed  map[Key]bool
	released map[Key]bool
}

// Held reports whether k is down this frame.
func (in Input) Held(k Key) bool { return in.held[k] }

// Pressed reports whether k went down on this frame (was up last frame).
func (in Input) Pressed(k Key) bool { return in.pressed[k] }

// Released reports whether k went up on this frame (was down last frame).
func (in Input) Released(k Key) bool { return in.released[k] }

// inputTracker accumulates key events and produces Input snapshots. Times are
// durations since run start, supplied by the loop (or tests).
type inputTracker struct {
	lastSeen map[Key]time.Duration
	prevHeld map[Key]bool
}

func newInputTracker() *inputTracker {
	return &inputTracker{
		lastSeen: map[Key]time.Duration{},
		prevHeld: map[Key]bool{},
	}
}

// press records that key k produced an event (initial press or OS repeat) at
// elapsed time at.
func (t *inputTracker) press(k Key, at time.Duration) {
	t.lastSeen[k] = at
}

// snapshot computes the Input for the current frame at elapsed time now, then
// records the held state for next-frame edge detection.
func (t *inputTracker) snapshot(now time.Duration) Input {
	held := map[Key]bool{}
	pressed := map[Key]bool{}
	released := map[Key]bool{}
	for k, last := range t.lastSeen {
		isHeld := now-last <= decayWindow
		if isHeld {
			held[k] = true
			if !t.prevHeld[k] {
				pressed[k] = true
			}
		} else if t.prevHeld[k] {
			released[k] = true
		}
	}
	t.prevHeld = held
	return Input{held: held, pressed: pressed, released: released}
}
