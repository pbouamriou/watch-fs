package utils

import (
	"os"
	"path/filepath"
)

// ValidateDirectory checks if a path is a valid directory
func ValidateDirectory(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}
	if !info.IsDir() {
		return &os.PathError{
			Op:   "validate",
			Path: path,
			Err:  os.ErrInvalid,
		}
	}
	return nil
}

// GetRelativePath returns the relative path from root to the given path
func GetRelativePath(root, path string) (string, error) {
	return filepath.Rel(root, path)
}

// IsHidden checks if a file or directory is hidden
func IsHidden(path string) bool {
	base := filepath.Base(path)
	return len(base) > 0 && base[0] == '.'
}

// ShouldIgnore checks if a path should be ignored
func ShouldIgnore(path string) bool {
	// Ignore hidden files and directories
	if IsHidden(path) {
		return true
	}

	// Ignore common system directories
	base := filepath.Base(path)
	ignoreDirs := []string{
		"node_modules",
		".git",
		".svn",
		".hg",
		"__pycache__",
		".DS_Store",
		"Thumbs.db",
	}

	for _, dir := range ignoreDirs {
		if base == dir {
			return true
		}
	}

	return false
}
