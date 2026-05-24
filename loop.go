package shaw

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Action is what a Game asks the loop to do after a frame.
type Action int

const (
	Continue Action = iota // keep running
	Quit                   // stop the loop and restore the terminal
)

// Game is the interface a shaw game implements. Update advances one frame given
// the time since the previous frame and the current input; Draw paints the
// frame into the canvas.
type Game interface {
	Update(dt time.Duration, in Input) Action
	Draw(c *Canvas)
}

// Options configures Run. Zero values mean: auto-size the canvas to the terminal
// (Width and Height both 0), and default to 30 FPS.
type Options struct {
	Width, Height int
	FPS           int
	Title         string
}

const defaultFPS = 30

// frameMsg is the per-frame tick.
type frameMsg time.Time

type model struct {
	game    Game
	canvas  *Canvas
	tracker *inputTracker
	fps     int
	auto    bool
	start   time.Time
	last    time.Duration
	quit    bool
}

func newModel(g Game, opts Options) *model {
	fps := opts.FPS
	if fps <= 0 {
		fps = defaultFPS
	}
	auto := opts.Width == 0 && opts.Height == 0
	w, h := opts.Width, opts.Height
	if auto {
		w, h = 80, 48
	}
	return &model{
		game:    g,
		canvas:  NewCanvas(w, h),
		tracker: newInputTracker(),
		fps:     fps,
		auto:    auto,
		start:   time.Now(),
	}
}

func (m *model) elapsed() time.Duration { return time.Since(m.start) }

func (m *model) tick() tea.Cmd {
	d := time.Second / time.Duration(m.fps)
	return tea.Tick(d, func(t time.Time) tea.Msg { return frameMsg(t) })
}

func (m *model) Init() tea.Cmd { return m.tick() }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.quit = true
			return m, tea.Quit
		}
		m.tracker.press(Key(msg.String()), m.elapsed())
		return m, nil
	case tea.WindowSizeMsg:
		if m.auto {
			m.canvas = NewCanvas(msg.Width, msg.Height*2)
		}
		return m, nil
	case frameMsg:
		now := m.elapsed()
		in := m.tracker.snapshot(now)
		dt := now - m.last
		m.last = now
		if m.game.Update(dt, in) == Quit {
			m.quit = true
			return m, tea.Quit
		}
		return m, m.tick()
	}
	return m, nil
}

func (m *model) View() string {
	m.game.Draw(m.canvas)
	return m.canvas.Render()
}

// Run starts the game loop in the alternate screen and blocks until the game
// returns Quit or the user presses Ctrl+C, then restores the terminal.
func Run(g Game, opts Options) error {
	_, err := tea.NewProgram(newModel(g, opts), tea.WithAltScreen()).Run()
	return err
}
