// Package corpus discovers .txt files and serves their words as an endless stream.
package corpus

import (
	"io/fs"
	"path/filepath"
	"strings"
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
