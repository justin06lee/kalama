package corpus

import (
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestScanFindsTxtRecursively(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), "alpha")
	writeFile(t, filepath.Join(dir, "sub", "b.txt"), "beta")
	writeFile(t, filepath.Join(dir, "sub", "c.md"), "ignored")

	got, err := Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	want := []string{
		filepath.Join(dir, "a.txt"),
		filepath.Join(dir, "sub", "b.txt"),
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestScanEmptyDir(t *testing.T) {
	got, err := Scan(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no files, got %v", got)
	}
}

func TestTextStreamNormalizesWhitespace(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), "one\ntwo\t three   four\n\nfive")
	files, _ := Scan(dir)

	s := NewTextStream(files, rand.New(rand.NewSource(1)))
	var got []string
	for i := 0; i < 5; i++ {
		w, ok := s.Next()
		if !ok {
			t.Fatalf("stream ended early at %d", i)
		}
		got = append(got, w)
	}
	want := []string{"one", "two", "three", "four", "five"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestTextStreamRollsOverToAnotherFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), "x")
	writeFile(t, filepath.Join(dir, "b.txt"), "y")
	files, _ := Scan(dir)

	s := NewTextStream(files, rand.New(rand.NewSource(1)))
	// Two single-word files ("x" and "y"): the stream is endless, so every
	// read must succeed. Across many reads it must also draw from BOTH files.
	seen := map[string]bool{}
	for i := 0; i < 30; i++ {
		w, ok := s.Next()
		if !ok {
			t.Fatalf("stream ended at read %d, expected endless", i)
		}
		seen[w] = true
	}
	if !seen["x"] || !seen["y"] {
		t.Fatalf("expected words from both files, saw %v", seen)
	}
}

func TestTextStreamSkipsNonUTF8(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.txt"), []byte{0xff, 0xfe, 0xfd}, 0o644); err != nil {
		t.Fatal(err)
	}
	files, _ := Scan(dir)
	s := NewTextStream(files, rand.New(rand.NewSource(1)))
	if _, ok := s.Next(); ok {
		t.Fatal("expected no words from a non-UTF8-only corpus")
	}
}
