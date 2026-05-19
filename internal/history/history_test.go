package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPathUsesXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdgtest")
	got, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("/tmp/xdgtest", "shaw", "history.json")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAppendThenLoadRoundTrips(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	rec := Record{
		Time: time.Unix(1000, 0).UTC(), Mode: "time", Target: 30,
		NetWPM: 55.5, RawWPM: 60, Accuracy: 0.95, Consistency: 80,
	}
	if err := Append(rec); err != nil {
		t.Fatal(err)
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].NetWPM != 55.5 || got[0].Mode != "time" {
		t.Fatalf("round trip mismatch: %+v", got)
	}
}

func TestAppendAccumulates(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := Append(Record{Mode: "words", Target: 25}); err != nil {
		t.Fatal(err)
	}
	if err := Append(Record{Mode: "zen", Target: 0}); err != nil {
		t.Fatal(err)
	}
	got, _ := Load()
	if len(got) != 2 {
		t.Fatalf("got %d records, want 2", len(got))
	}
}

func TestAppendBacksUpCorruptFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	p := filepath.Join(dir, "shaw", "history.json")
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	invalid := []byte("{not json")
	if err := os.WriteFile(p, invalid, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Append(Record{Mode: "time", Target: 30}); err != nil {
		t.Fatal(err)
	}

	backup, err := os.ReadFile(p + ".corrupt")
	if err != nil {
		t.Fatalf("expected corrupt backup file: %v", err)
	}
	if string(backup) != string(invalid) {
		t.Fatalf("backup contents %q, want %q", backup, invalid)
	}

	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d records, want 1", len(got))
	}
}

func TestLoadMissingFileIsEmpty(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestLoadCorruptFileIsEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	p := filepath.Join(dir, "shaw", "history.json")
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	_ = os.WriteFile(p, []byte("{not json"), 0o644)
	got, err := Load()
	if err != nil {
		t.Fatalf("corrupt file should not error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty on corrupt file, got %v", got)
	}
}
