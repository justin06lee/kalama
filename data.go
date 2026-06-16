package kalama

import (
	"os"
	"path/filepath"
)

// DataDir returns and creates a per-game persistence directory. It uses
// $SHAW_DATA_DIR/<game> when the env var is set, otherwise
// ~/.shaw/data/<game>. Games store their own scores/history here.
func DataDir(game string) (string, error) {
	base := os.Getenv("SHAW_DATA_DIR")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".shaw", "data")
	}
	dir := filepath.Join(base, game)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}
