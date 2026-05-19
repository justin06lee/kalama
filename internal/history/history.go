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

// Append adds rec to the history file. A missing file is created. If the
// existing file is unparseable JSON, it is preserved as <path>.corrupt before
// a fresh file is written, so a damaged history is never silently destroyed.
func Append(rec Record) error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	recs, err := readForAppend(p)
	if err != nil {
		return err
	}
	recs = append(recs, rec)
	data, err := json.MarshalIndent(recs, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(p, data)
}

// readForAppend returns the records in the file at p. A missing file yields an
// empty slice. A file that does not parse as JSON is renamed to p+".corrupt"
// and an empty slice is returned, so the damaged data is preserved, not lost.
func readForAppend(p string) ([]Record, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var recs []Record
	if err := json.Unmarshal(data, &recs); err != nil {
		if rerr := os.Rename(p, p+".corrupt"); rerr != nil {
			return nil, rerr
		}
		return nil, nil
	}
	return recs, nil
}

// writeAtomic writes data to p via a temp file in p's directory followed by a
// rename, so an interrupted write cannot leave p corrupt.
func writeAtomic(p string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(p), ".history-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, p)
}
