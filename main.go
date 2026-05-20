// Command shaw is a monkeytype-style terminal typing trainer.
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/justin06lee/shaw/internal/config"
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
	setDirFlag := flag.String("set-dir", "",
		"save the given path as shaw's default corpus directory, then exit")
	flag.Parse()

	if *setDirFlag != "" {
		setDefaultDir(*setDirFlag)
		return
	}

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

	dir := resolveCorpusDir()
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

// resolveCorpusDir picks the corpus folder using this priority:
// 1. positional arg (`shaw /some/dir`)
// 2. SHAW_DIR environment variable
// 3. DefaultDir from the saved config file
// Fails with a helpful message when none are set.
func resolveCorpusDir() string {
	if flag.NArg() > 0 {
		return flag.Arg(0)
	}
	if env := os.Getenv("SHAW_DIR"); env != "" {
		return env
	}
	cfg, err := config.Load()
	if err == nil && cfg.DefaultDir != "" {
		return cfg.DefaultDir
	}
	fail("no corpus directory configured.\n" +
		"  Pass one:        shaw /path/to/folder\n" +
		"  Set in env:      export SHAW_DIR=/path/to/folder\n" +
		"  Save as default: shaw --set-dir /path/to/folder")
	return "" // unreachable
}

// setDefaultDir validates the path then saves it to the config file.
func setDefaultDir(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		fail("cannot resolve %s: %v", path, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		fail("cannot read %s: %v", abs, err)
	}
	if !info.IsDir() {
		fail("%s is not a directory", abs)
	}
	files, err := corpus.Scan(abs)
	if err != nil {
		fail("cannot scan %s: %v", abs, err)
	}
	if len(files) == 0 {
		fail("no .txt files found in %s", abs)
	}
	if err := config.Save(config.Config{DefaultDir: abs}); err != nil {
		fail("cannot save config: %v", err)
	}
	fmt.Printf("default corpus directory set to %s\n", abs)
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
