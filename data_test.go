package kalama

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDataDirUsesEnvAndCreates(t *testing.T) {
	base := t.TempDir()
	t.Setenv("SHAW_DATA_DIR", base)
	dir, err := DataDir("fighter")
	if err != nil {
		t.Fatalf("DataDir: %v", err)
	}
	want := filepath.Join(base, "fighter")
	if dir != want {
		t.Errorf("dir = %q, want %q", dir, want)
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Errorf("expected created directory at %q, stat err = %v", dir, err)
	}
}
