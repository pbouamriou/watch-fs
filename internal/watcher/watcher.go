package watcher

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watcher wraps fsnotify.Watcher with additional functionality
type Watcher struct {
	watcher *fsnotify.Watcher
	root    string
}

// New creates a new file system watcher
func New(root string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: watcher,
		root:    root,
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

// GetRoot returns the root directory being watched
func (w *Watcher) GetRoot() string {
	return w.root
}
