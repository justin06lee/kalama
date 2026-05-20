// Package config persists shaw's user-level settings as a JSON file.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds shaw's persistent user settings.
type Config struct {
	// DefaultDir is the corpus folder used when shaw is invoked with no
	// positional argument and no SHAW_DIR environment variable.
	DefaultDir string `json:"default_dir"`
}

// Path returns the config file location, honoring XDG_CONFIG_HOME.
func Path() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "shaw", "config.json"), nil
}

// Load reads the config. A missing or corrupt file yields a zero Config and
// no error so a broken settings file never blocks the program.
func Load() (Config, error) {
	p, err := Path()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, nil // tolerate corruption
	}
	return cfg, nil
}

// Save writes cfg to disk, creating directories and writing atomically via a
// temp file + rename so an interrupted write cannot corrupt the config.
func Save(cfg Config) error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), ".config-*.tmp")
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
