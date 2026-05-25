package shaw

import (
	"os"
	"path/filepath"
)

// DataDir returns and creates a per-game persistence directory. It uses
// $KALAMA_DATA_DIR/<game> when the env var is set, otherwise
// ~/.kalama/data/<game>. Games store their own scores/history here.
func DataDir(game string) (string, error) {
	base := os.Getenv("KALAMA_DATA_DIR")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".kalama", "data")
	}
	dir := filepath.Join(base, game)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}
