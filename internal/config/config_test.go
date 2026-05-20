package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathUsesXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdgtest")
	got, err := Path()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("/tmp/xdgtest", "shaw", "config.json")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSaveLoadRoundTrips(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := Save(Config{DefaultDir: "/some/dir"}); err != nil {
		t.Fatal(err)
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.DefaultDir != "/some/dir" {
		t.Errorf("DefaultDir: got %q, want %q", got.DefaultDir, "/some/dir")
	}
}

func TestLoadMissingFileIsEmpty(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.DefaultDir != "" {
		t.Errorf("expected zero Config, got %+v", got)
	}
}

func TestLoadCorruptFileIsEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	p := filepath.Join(dir, "shaw", "config.json")
	_ = os.MkdirAll(filepath.Dir(p), 0o755)
	_ = os.WriteFile(p, []byte("{not json"), 0o644)
	got, err := Load()
	if err != nil {
		t.Fatalf("corrupt file should not error: %v", err)
	}
	if got.DefaultDir != "" {
		t.Errorf("expected zero Config on corrupt file, got %+v", got)
	}
}
