package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pbouamriou/watch-fs/internal/ui"
)

// TestHelper provides utilities for TUI testing
type TestHelper struct {
	ui       *ui.UI
	watcher  *MockWatcher
	tempDirs []string
}

// NewTestHelper creates a new test helper with a mock watcher
func NewTestHelper(t *testing.T) *TestHelper {
	mockWatcher := NewMockWatcher()
	uiInstance := ui.NewUI(mockWatcher, "/test/path")

	return &TestHelper{
		ui:      uiInstance,
		watcher: mockWatcher,
	}
}

// NewTestHelperWithTempDir creates a test helper with a real temporary directory
func NewTestHelperWithTempDir(t *testing.T, dirName string) *TestHelper {
	tempDir, err := os.MkdirTemp("", dirName)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mockWatcher := NewMockWatcher()
	uiInstance := ui.NewUI(mockWatcher, tempDir)

	helper := &TestHelper{
		ui:       uiInstance,
		watcher:  mockWatcher,
		tempDirs: []string{tempDir},
	}

	return helper
}

// Cleanup cleans up temporary directories created by the helper
func (th *TestHelper) Cleanup() {
	th.watcher.Close()
	for _, dir := range th.tempDirs {
		_ = os.RemoveAll(dir)
	}
}

// CreateTestDirectoryStructure creates a directory structure for testing
func (th *TestHelper) CreateTestDirectoryStructure(t *testing.T, structure map[string]interface{}) string {
	if len(th.tempDirs) == 0 {
		tempDir, err := os.MkdirTemp("", "watch-fs-test-structure")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		th.tempDirs = append(th.tempDirs, tempDir)
	}

	baseDir := th.tempDirs[0]
	th.createStructureRecursive(t, baseDir, structure)
	return baseDir
}

// createStructureRecursive recursively creates directory structure
func (th *TestHelper) createStructureRecursive(t *testing.T, basePath string, structure map[string]interface{}) {
	for name, value := range structure {
		fullPath := filepath.Join(basePath, name)

		switch v := value.(type) {
		case map[string]interface{}:
			// It's a directory
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", fullPath, err)
			}
			th.createStructureRecursive(t, fullPath, v)
		case string:
			// It's a file with content
			if err := os.WriteFile(fullPath, []byte(v), 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", fullPath, err)
			}
		case nil:
			// It's an empty file
			if err := os.WriteFile(fullPath, []byte(""), 0644); err != nil {
				t.Fatalf("Failed to create empty file %s: %v", fullPath, err)
			}
		}
	}
}

// AddTestEvents adds a series of test events to the UI
func (th *TestHelper) AddTestEvents(events []TestEvent) {
	for _, event := range events {
		th.ui.AddEvent(event.Path, event.Operation, event.IsDir)
		if event.Delay > 0 {
			time.Sleep(event.Delay)
		}
	}
}

// TestEvent represents a test event with optional delay
type TestEvent struct {
	Path      string
	Operation fsnotify.Op
	IsDir     bool
	Delay     time.Duration
}

// SimulateKeySequence simulates a sequence of key presses
func (th *TestHelper) SimulateKeySequence(keys []KeyPress) error {
	for _, key := range keys {
		switch key.Type {
		case KeyTypeChar:
			// Simulate character key
			// In a real implementation, this would trigger the appropriate handler
		case KeyTypeSpecial:
			// Simulate special keys like arrows, enter, etc.
			// In a real implementation, this would trigger the appropriate handler
		}

		if key.Delay > 0 {
			time.Sleep(key.Delay)
		}
	}
	return nil
}

// KeyPress represents a key press event
type KeyPress struct {
	Type  KeyType
	Char  rune
	Key   string // For special keys like "Enter", "Escape", "ArrowUp"
	Delay time.Duration
}

// KeyType represents the type of key press
type KeyType int

const (
	KeyTypeChar KeyType = iota
	KeyTypeSpecial
)

// AssertEventCount verifies the number of events in the UI
func (th *TestHelper) AssertEventCount(t *testing.T, expected int) {
	actual := len(th.ui.GetState().Events)
	if actual != expected {
		t.Errorf("Expected %d events, got %d", expected, actual)
	}
}

// AssertScrollOffset verifies the current scroll offset
func (th *TestHelper) AssertScrollOffset(t *testing.T, expected int) {
	actual := th.ui.GetState().ScrollOffset
	if actual != expected {
		t.Errorf("Expected scroll offset %d, got %d", expected, actual)
	}
}

// AssertFocus verifies the current focus state
func (th *TestHelper) AssertFocus(t *testing.T, expected ui.FocusMode) {
	actual := th.ui.GetState().CurrentFocus
	if actual != expected {
		t.Errorf("Expected focus %v, got %v", expected, actual)
	}
}

// AssertAggregationState verifies the aggregation state
func (th *TestHelper) AssertAggregationState(t *testing.T, expected bool) {
	actual := th.ui.GetState().AggregateEvents
	if actual != expected {
		t.Errorf("Expected aggregation %v, got %v", expected, actual)
	}
}

// GetEventsByPath gets all events matching a specific path pattern
func (th *TestHelper) GetEventsByPath(pathPattern string) []*ui.FileEvent {
	var matching []*ui.FileEvent
	for _, event := range th.ui.GetState().Events {
		if matched, _ := filepath.Match(pathPattern, event.Path); matched {
			matching = append(matching, event)
		}
	}
	return matching
}

// WaitForEventCount waits for a specific number of events with timeout
func (th *TestHelper) WaitForEventCount(t *testing.T, expected int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if len(th.ui.GetState().Events) == expected {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}

	actual := len(th.ui.GetState().Events)
	t.Errorf("Timeout waiting for %d events, got %d", expected, actual)
	return false
}

// Example test using the helper functions
func TestHelperFunctionality(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// Test directory structure creation
	structure := map[string]interface{}{
		"dir1": map[string]interface{}{
			"subdir1": map[string]interface{}{
				"file1.txt": "content1",
				"file2.txt": "content2",
			},
			"file3.txt": "content3",
		},
		"dir2": map[string]interface{}{
			"empty_file.txt": nil,
		},
		"root_file.txt": "root content",
	}

	baseDir := helper.CreateTestDirectoryStructure(t, structure)

	// Verify structure was created
	if _, err := os.Stat(filepath.Join(baseDir, "dir1", "subdir1", "file1.txt")); os.IsNotExist(err) {
		t.Error("Test structure was not created properly")
	}

	// Test event addition
	testEvents := []TestEvent{
		{Path: "/test/file1.txt", Operation: fsnotify.Write, IsDir: false, Delay: 1 * time.Millisecond},
		{Path: "/test/file2.txt", Operation: fsnotify.Create, IsDir: false, Delay: 1 * time.Millisecond},
		{Path: "/test/dir1", Operation: fsnotify.Create, IsDir: true, Delay: 1 * time.Millisecond},
	}

	helper.AddTestEvents(testEvents)
	helper.AssertEventCount(t, 3)

	// Test navigation
	helper.ui.MoveDown()
	helper.AssertScrollOffset(t, 1)

	helper.ui.MoveToTop()
	helper.AssertScrollOffset(t, 0)

	// Test aggregation
	helper.AssertAggregationState(t, true) // Should be enabled by default
	helper.ui.ToggleAggregate()
	helper.AssertAggregationState(t, false)
}

// TestPerformanceWithManyEvents tests UI performance with a large number of events
func TestPerformanceWithManyEvents(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	start := time.Now()

	// Add many events
	for i := 0; i < 1000; i++ {
		helper.ui.AddEvent(fmt.Sprintf("/test/file%d.txt", i), fsnotify.Write, false)
	}

	elapsed := time.Since(start)
	t.Logf("Added 1000 events in %v", elapsed)

	if elapsed > 1*time.Second {
		t.Errorf("Adding 1000 events took too long: %v", elapsed)
	}

	// Test navigation performance
	start = time.Now()
	helper.ui.MoveToBottom()
	helper.ui.MoveToTop()
	elapsed = time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("Navigation took too long: %v", elapsed)
	}

	// Test filtering performance
	start = time.Now()
	filtered := helper.ui.GetFilteredEvents()
	elapsed = time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("Filtering took too long: %v", elapsed)
	}

	// Verify we got the expected number of events
	if len(filtered) != 1000 {
		t.Errorf("Expected 1000 filtered events, got %d", len(filtered))
	}
}
