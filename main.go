// Command shaw is a monkeytype-style terminal typing trainer.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/justin06lee/shaw/internal/corpus"
	"github.com/justin06lee/shaw/internal/history"
	"github.com/justin06lee/shaw/internal/run"
	"github.com/justin06lee/shaw/internal/stats"
	"github.com/justin06lee/shaw/internal/tui"
)

func main() {
	timeFlag := flag.Int("time", 0, "timed mode: 15, 30, 60, or 120 seconds")
	wordsFlag := flag.Int("words", 0, "word mode: 10, 25, 50, or 100 words")
	zenFlag := flag.Bool("zen", false, "zen mode: type until Esc")
	histFlag := flag.Bool("history", false, "print progress chart and exit")
	flag.Parse()

	if *histFlag {
		printHistory()
		return
	}

	mode, target := run.ModeTime, 30
	switch {
	case *zenFlag:
		mode, target = run.ModeZen, 0
	case *wordsFlag > 0:
		mode, target = run.ModeWords, *wordsFlag
	case *timeFlag > 0:
		mode, target = run.ModeTime, *timeFlag
	}

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	files, err := corpus.Scan(dir)
	if err != nil {
		fail("cannot scan %s: %v", dir, err)
	}
	if len(files) == 0 {
		fail("no .txt files found in %s", dir)
	}

	stream := corpus.NewTextStream(files, rand.New(rand.NewSource(time.Now().UnixNano())))
	if _, ok := stream.Next(); !ok {
		fail("no usable text in %s (files empty or not UTF-8)", dir)
	}
	stream = corpus.NewTextStream(files, rand.New(rand.NewSource(time.Now().UnixNano())))

	m := tui.New(stream, mode, target, 80, 24)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fail("%v", err)
	}
}

// printHistory renders an ASCII chart of net WPM across past runs.
func printHistory() {
	recs, err := history.Load()
	if err != nil {
		fail("cannot read history: %v", err)
	}
	if len(recs) == 0 {
		fmt.Println("no runs recorded yet")
		return
	}
	samples := make([]float64, len(recs))
	for i, r := range recs {
		samples[i] = r.NetWPM
	}
	fmt.Printf("net wpm across %d runs:\n\n", len(recs))
	fmt.Println(stats.RenderChart(samples, 60, 12))
	fmt.Printf("\nlatest: %.0f wpm\n", samples[len(samples)-1])
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "shaw: "+format+"\n", args...)
	os.Exit(1)
}
