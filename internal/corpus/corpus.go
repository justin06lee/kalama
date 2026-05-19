// Package corpus discovers .txt files and serves their words as an endless stream.
package corpus

import (
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Scan walks dir recursively and returns the paths of all .txt files.
func Scan(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.EqualFold(filepath.Ext(path), ".txt") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// TextStream yields an endless sequence of words drawn from a set of files.
// When a file is exhausted it picks another random file. Files that are not
// valid UTF-8 are skipped permanently.
type TextStream struct {
	files []string
	rng   *rand.Rand
	buf   []string // unread words from the current file
	dead  map[int]bool
}

// NewTextStream creates a stream over files using rng for file selection.
func NewTextStream(files []string, rng *rand.Rand) *TextStream {
	return &TextStream{files: files, rng: rng, dead: map[int]bool{}}
}

// Next returns the next word. It returns ("", false) only when no file in the
// corpus yields usable words (empty corpus or every file is non-UTF8/empty).
func (s *TextStream) Next() (string, bool) {
	for len(s.buf) == 0 {
		if !s.loadRandomFile() {
			return "", false
		}
	}
	w := s.buf[0]
	s.buf = s.buf[1:]
	return w, true
}

// loadRandomFile fills buf from a random non-dead file. Returns false when no
// file can produce words.
func (s *TextStream) loadRandomFile() bool {
	alive := make([]int, 0, len(s.files))
	for i := range s.files {
		if !s.dead[i] {
			alive = append(alive, i)
		}
	}
	if len(alive) == 0 {
		return false
	}
	idx := alive[s.rng.Intn(len(alive))]
	data, err := os.ReadFile(s.files[idx])
	if err != nil || !utf8.Valid(data) {
		s.dead[idx] = true
		return s.loadRandomFile()
	}
	words := strings.Fields(string(data))
	if len(words) == 0 {
		s.dead[idx] = true
		return s.loadRandomFile()
	}
	s.buf = words
	return true
}
