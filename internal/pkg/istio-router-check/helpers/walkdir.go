package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

// WalkYAML walks the baseDirs and returns a list of all files found with yaml extension.
func WalkYAML(baseDirs ...string) ([]string, error) {
	return WalkFilter(func(path string, info os.FileInfo) bool {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
			return true
		}
		return false
	}, baseDirs...)
}

// WalkFilter walks the baseDirs and returns a list of all files found that return true in the filterFunc.
func WalkFilter(filterFunc func(path string, info os.FileInfo) bool, baseDirs ...string) ([]string, error) {
	var files []string
	for _, baseDir := range baseDirs {
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("could not read directory %q: %w", baseDir, err)
			}
			if info.IsDir() {
				return nil
			}
			if !filterFunc(path, info) {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("could not read directory %q: %w", baseDir, err)
		}
	}
	return files, nil
}
