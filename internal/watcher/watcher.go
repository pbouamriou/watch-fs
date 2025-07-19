package watcher

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watcher wraps fsnotify.Watcher with additional functionality
type Watcher struct {
	watcher *fsnotify.Watcher
	roots   []string // Changed from single root to multiple roots
}

// New creates a new file system watcher
func New(root string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: watcher,
		roots:   []string{root}, // Initialize with single root for backward compatibility
	}, nil
}

// NewMultiRoot creates a new file system watcher with multiple root directories
func NewMultiRoot(roots []string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: watcher,
		roots:   roots,
	}, nil
}

// Close closes the watcher
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

// AddRecursive adds a directory and all its subdirectories to the watcher
func (w *Watcher) AddRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			err = w.watcher.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// AddAllRootsRecursive adds all root directories and their subdirectories to the watcher
func (w *Watcher) AddAllRootsRecursive() error {
	for _, root := range w.roots {
		if err := w.AddRecursive(root); err != nil {
			return err
		}
	}
	return nil
}

// Events returns the events channel
func (w *Watcher) Events() <-chan fsnotify.Event {
	return w.watcher.Events
}

// Errors returns the errors channel
func (w *Watcher) Errors() <-chan error {
	return w.watcher.Errors
}

// AddDirectory adds a new directory to the watcher (for newly created directories)
func (w *Watcher) AddDirectory(path string) error {
	return w.watcher.Add(path)
}

// GetRoots returns all root directories being watched
func (w *Watcher) GetRoots() []string {
	return w.roots
}

// GetRoot returns the first root directory being watched (for backward compatibility)
func (w *Watcher) GetRoot() string {
	if len(w.roots) > 0 {
		return w.roots[0]
	}
	return ""
}

// AddRoot adds a new root directory to watch
func (w *Watcher) AddRoot(root string) error {
	// Add the root to our list
	w.roots = append(w.roots, root)

	// Add it recursively to the watcher
	return w.AddRecursive(root)
}

// RemoveRoot removes a root directory from watching
func (w *Watcher) RemoveRoot(root string) error {
	// Find and remove the root from our list
	for i, r := range w.roots {
		if r == root {
			w.roots = append(w.roots[:i], w.roots[i+1:]...)
			break
		}
	}

	// Note: fsnotify doesn't support removing individual directories
	// We would need to recreate the watcher to remove a directory
	// For now, we'll just remove it from our list
	return nil
}
