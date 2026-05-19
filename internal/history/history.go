// Package history persists typing-run results as a JSON file.
package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Record is one persisted run result.
type Record struct {
	Time        time.Time `json:"time"`
	Mode        string    `json:"mode"`
	Target      int       `json:"target"`
	NetWPM      float64   `json:"net_wpm"`
	RawWPM      float64   `json:"raw_wpm"`
	Accuracy    float64   `json:"accuracy"`
	Consistency float64   `json:"consistency"`
}

// Path returns the history file location, honoring XDG_CONFIG_HOME.
func Path() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "shaw", "history.json"), nil
}

// Load reads all records. A missing or corrupt file yields an empty slice and
// no error, so a damaged history never blocks a run.
func Load() ([]Record, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var recs []Record
	if err := json.Unmarshal(data, &recs); err != nil {
		return nil, nil // tolerate corruption
	}
	return recs, nil
}

// Append adds rec to the history file, creating it and its directory if needed.
func Append(rec Record) error {
	p, err := Path()
	if err != nil {
		return err
	}
	recs, err := Load()
	if err != nil {
		return err
	}
	recs = append(recs, rec)
	data, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}
