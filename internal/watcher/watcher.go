package watcher

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pbouamriou/watch-fs/pkg/logger"
)

// Watcher wraps fsnotify.Watcher with additional functionality
type Watcher struct {
	watcher *fsnotify.Watcher
	roots   []string        // Root directories being watched
	watched map[string]bool // Track all watched directories for removal
	mu      sync.RWMutex    // Protect concurrent access to roots and watched
}

// New creates a new file system watcher
func New(root string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher: watcher,
		roots:   []string{root},
		watched: make(map[string]bool),
	}, nil
}

// NewMultiRoot creates a new file system watcher with multiple root directories
func NewMultiRoot(roots []string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher: watcher,
		roots:   roots,
		watched: make(map[string]bool),
	}

	// Add all roots recursively
	if err := w.AddAllRootsRecursive(); err != nil {
		return nil, err
	}

	return w, nil
}

// Close closes the watcher
func (w *Watcher) Close() error {
	return w.watcher.Close()
}

// addRecursiveUnsafe adds a directory and all its subdirectories to the watcher
// This function is NOT thread-safe and assumes the caller holds the mutex
func (w *Watcher) addRecursiveUnsafe(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			err = w.watcher.Add(path)
			if err != nil {
				return err
			}
			// No mutex needed - caller must hold the lock
			w.watched[path] = true
		}
		return nil
	})
}

// AddRecursive adds a directory and all its subdirectories to the watcher
// This function is thread-safe
func (w *Watcher) AddRecursive(root string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.addRecursiveUnsafe(root)
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
	err := w.watcher.Add(path)
	if err == nil {
		w.mu.Lock()
		w.watched[path] = true
		w.mu.Unlock()
	}
	return err
}

// GetRoots returns all root directories being watched
func (w *Watcher) GetRoots() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// Return a copy to avoid race conditions
	roots := make([]string, len(w.roots))
	copy(roots, w.roots)
	return roots
}

// GetRoot returns the first root directory being watched (for backward compatibility)
func (w *Watcher) GetRoot() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.roots) > 0 {
		return w.roots[0]
	}
	return ""
}

// AddRoot adds a new root directory to watch
func (w *Watcher) AddRoot(root string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Check if root is already being watched
	for _, r := range w.roots {
		if r == root {
			return nil // Already watching this root
		}
	}

	// Add the root to our list
	w.roots = append(w.roots, root)

	// Add it recursively to the watcher (using unsafe version since we hold the lock)
	return w.addRecursiveUnsafe(root)
}

// RemoveRoot removes a root directory from watching
func (w *Watcher) RemoveRoot(root string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Find and remove the root from our list
	rootIndex := -1
	for i, r := range w.roots {
		if r == root {
			rootIndex = i
			break
		}
	}

	if rootIndex == -1 {
		return nil // Root not found, nothing to remove
	}

	// Remove from roots list
	w.roots = append(w.roots[:rootIndex], w.roots[rootIndex+1:]...)

	// Remove all subdirectories of this root from the watcher
	// We need to recreate the watcher to properly remove directories
	return w.recreateWatcherWithoutRoot(root)
}

// recreateWatcherWithoutRoot recreates the watcher without the specified root
func (w *Watcher) recreateWatcherWithoutRoot(rootToRemove string) error {
	// Store current watched directories
	oldWatched := make(map[string]bool)
	for path := range w.watched {
		oldWatched[path] = true
	}

	// Close old watcher
	if err := w.watcher.Close(); err != nil {
		logger.Error(err, "Failed to close watcher during recreation")
	}

	// Create new watcher
	newWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	w.watcher = newWatcher
	w.watched = make(map[string]bool)

	// Re-add all roots except the one to remove (using unsafe version since caller holds the lock)
	for _, root := range w.roots {
		if root != rootToRemove {
			if err := w.addRecursiveUnsafe(root); err != nil {
				return err
			}
		}
	}

	return nil
}

// IsWatching returns true if the given path is being watched
func (w *Watcher) IsWatching(path string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// Check if it's a root directory
	for _, root := range w.roots {
		if root == path {
			return true
		}
	}

	// Check if it's a watched subdirectory
	return w.watched[path]
}

// GetWatchedCount returns the number of directories being watched
func (w *Watcher) GetWatchedCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.watched)
}

// GetWatchedCountForRoot returns the number of directories being watched under a specific root
func (w *Watcher) GetWatchedCountForRoot(root string) int {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// Clean the root path to ensure consistent comparison
	cleanRoot := filepath.Clean(root)
	if !filepath.IsAbs(cleanRoot) {
		cleanRoot, _ = filepath.Abs(cleanRoot)
	}

	count := 0
	for watchedPath := range w.watched {
		cleanWatched := filepath.Clean(watchedPath)
		if !filepath.IsAbs(cleanWatched) {
			cleanWatched, _ = filepath.Abs(cleanWatched)
		}

		// Check if the watched path is under this root
		// Resolve symlinks for accurate comparison
		resolvedRoot, err := filepath.EvalSymlinks(cleanRoot)
		if err != nil {
			resolvedRoot = cleanRoot
		}

		resolvedWatched, err := filepath.EvalSymlinks(cleanWatched)
		if err != nil {
			resolvedWatched = cleanWatched
		}

		// Use filepath.Rel to check if watchedPath is under root
		if rel, err := filepath.Rel(resolvedRoot, resolvedWatched); err == nil && !filepath.IsAbs(rel) && rel != ".." {
			// Check if it doesn't start with "../" (which means it's not under the root)
			// Handle both Unix and Windows path separators
			parentDir := ".." + string(filepath.Separator)
			if runtime.GOOS == "windows" {
				// On Windows, also handle forward slashes that might be in the relative path
				rel = strings.ReplaceAll(rel, "/", string(filepath.Separator))
			}
			if len(rel) < len(parentDir) || !strings.HasPrefix(rel, parentDir) {
				count++
			}
		}
	}
	return count
}
