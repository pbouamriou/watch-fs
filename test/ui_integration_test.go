package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pbouamriou/watch-fs/internal/ui"
)

// MockGui simulates gocui.Gui for integration testing
type MockGui struct {
	currentView string
	views       map[string]*MockGuiView
	keyBindings map[string][]KeyBinding
}

type MockGuiView struct {
	name    string
	title   string
	content strings.Builder
	cursor  struct{ x, y int }
	origin  struct{ x, y int }
	size    struct{ width, height int }
}

type KeyBinding struct {
	key      rune
	callback func() error
}

func NewMockGui() *MockGui {
	return &MockGui{
		views:       make(map[string]*MockGuiView),
		keyBindings: make(map[string][]KeyBinding),
	}
}

func (mg *MockGui) SetCurrentView(name string) (*MockGuiView, error) {
	mg.currentView = name
	if view, exists := mg.views[name]; exists {
		return view, nil
	}
	// Create view if it doesn't exist
	view := &MockGuiView{
		name: name,
		size: struct{ width, height int }{width: 80, height: 24},
	}
	mg.views[name] = view
	return view, nil
}

func (mg *MockGui) GetCurrentView() *MockGuiView {
	return mg.views[mg.currentView]
}

func (mg *MockGui) SimulateKey(key rune) error {
	bindings := mg.keyBindings[mg.currentView]
	for _, binding := range bindings {
		if binding.key == key {
			return binding.callback()
		}
	}
	return nil
}

func (mv *MockGuiView) Clear() {
	mv.content.Reset()
}

func (mv *MockGuiView) Write(p []byte) (n int, err error) {
	return mv.content.Write(p)
}

func (mv *MockGuiView) GetContent() string {
	return mv.content.String()
}

func (mv *MockGuiView) SetTitle(title string) {
	mv.title = title
}

func (mv *MockGuiView) SetCursor(x, y int) error {
	mv.cursor.x = x
	mv.cursor.y = y
	return nil
}

func (mv *MockGuiView) SetOrigin(x, y int) error {
	mv.origin.x = x
	mv.origin.y = y
	return nil
}

func (mv *MockGuiView) InnerSize() (int, int) {
	return mv.size.width, mv.size.height
}

// TestUIFocusTransitions tests focus changes between different views
func TestUIFocusTransitions(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, "/test/path")

	// Test initial focus
	state := uiInstance.GetState()
	if state.CurrentFocus != ui.FocusMain {
		t.Errorf("Expected initial focus to be FocusMain, got %v", state.CurrentFocus)
	}

	// Test showing folder manager changes focus
	uiInstance.ShowFolderManager()
	if state.CurrentFocus != ui.FocusFolderManager {
		t.Errorf("Expected focus to be FocusFolderManager after showing folder manager, got %v", state.CurrentFocus)
	}

	// Test hiding folder manager returns focus to main
	uiInstance.HideFolderManager()
	if state.CurrentFocus != ui.FocusMain {
		t.Errorf("Expected focus to return to FocusMain after hiding folder manager, got %v", state.CurrentFocus)
	}
}

// TestUIStateConsistency tests that UI state remains consistent during operations
func TestUIStateConsistency(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, "/test/path")

	// Add events and test state consistency
	uiInstance.AddEvent("/test/file1.txt", fsnotify.Write, false)
	uiInstance.AddEvent("/test/file2.txt", fsnotify.Create, false)
	uiInstance.AddEvent("/test/dir1/", fsnotify.Create, true)

	state := uiInstance.GetState()

	// Test that events were added
	if len(state.Events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(state.Events))
	}

	// Test aggregation state consistency
	initialAggregation := state.AggregateEvents
	uiInstance.ToggleAggregate()
	if state.AggregateEvents == initialAggregation {
		t.Error("Aggregation state should change after toggle")
	}

	// Test filter state consistency
	initialShowFiles := state.Filter.ShowFiles
	uiInstance.ToggleFiles()
	if state.Filter.ShowFiles == initialShowFiles {
		t.Error("ShowFiles state should change after toggle")
	}

	initialShowDirs := state.Filter.ShowDirs
	uiInstance.ToggleDirs()
	if state.Filter.ShowDirs == initialShowDirs {
		t.Error("ShowDirs state should change after toggle")
	}
}

// TestEventProcessingOrder tests that events are processed in the correct order
func TestEventProcessingOrder(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, "/test/path")

	// Add events with specific timing
	uiInstance.AddEvent("/test/file1.txt", fsnotify.Write, false)
	time.Sleep(1 * time.Millisecond)
	uiInstance.AddEvent("/test/file2.txt", fsnotify.Write, false)
	time.Sleep(1 * time.Millisecond)
	uiInstance.AddEvent("/test/file3.txt", fsnotify.Write, false)

	// Get filtered events (should be sorted by time, newest first)
	events := uiInstance.GetFilteredEvents()

	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}

	// Check that events are in reverse chronological order (newest first)
	if events[0].Path != "/test/file3.txt" {
		t.Errorf("Expected first event to be file3.txt, got %s", events[0].Path)
	}
	if events[2].Path != "/test/file1.txt" {
		t.Errorf("Expected last event to be file1.txt, got %s", events[2].Path)
	}
}

// TestNavigationBoundaries tests navigation at list boundaries
func TestNavigationBoundaries(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, "/test/path")

	// Add exactly 5 events for testing
	for i := 0; i < 5; i++ {
		uiInstance.AddEvent(fmt.Sprintf("/test/file%d.txt", i), fsnotify.Write, false)
	}

	state := uiInstance.GetState()

	// Test moving to bottom
	uiInstance.MoveToBottom()
	if state.ScrollOffset != 4 { // 5 events, 0-indexed
		t.Errorf("Expected scroll offset 4 at bottom, got %d", state.ScrollOffset)
	}

	// Test that we can't go beyond bottom
	uiInstance.MoveDown()
	if state.ScrollOffset != 4 {
		t.Errorf("Expected scroll offset to remain 4 when at bottom, got %d", state.ScrollOffset)
	}

	// Test moving to top
	uiInstance.MoveToTop()
	if state.ScrollOffset != 0 {
		t.Errorf("Expected scroll offset 0 at top, got %d", state.ScrollOffset)
	}

	// Test that we can't go beyond top
	uiInstance.MoveUp()
	if state.ScrollOffset != 0 {
		t.Errorf("Expected scroll offset to remain 0 when at top, got %d", state.ScrollOffset)
	}
}

// TestConcurrentEventHandling tests that events can be safely added concurrently
func TestConcurrentEventHandling(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, "/test/path")

	// Add events concurrently
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			uiInstance.AddEvent(fmt.Sprintf("/test/goroutine1/file%d.txt", i), fsnotify.Write, false)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			uiInstance.AddEvent(fmt.Sprintf("/test/goroutine2/file%d.txt", i), fsnotify.Create, false)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	state := uiInstance.GetState()

	// We should have some events (exact number depends on aggregation)
	if len(state.Events) == 0 {
		t.Error("Expected some events after concurrent addition")
	}

	// All events should be valid
	for _, event := range state.Events {
		if event.Path == "" {
			t.Error("Found event with empty path")
		}
		if event.Count < 1 {
			t.Errorf("Found event with invalid count: %d", event.Count)
		}
	}
}

// TestErrorRecovery tests that the UI can recover from various error conditions
func TestErrorRecovery(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, "/test/path")

	// Test with invalid paths - these shouldn't crash the application
	uiInstance.AddEvent("", fsnotify.Write, false)
	uiInstance.AddEvent("/nonexistent/path/file.txt", fsnotify.Write, false)

	state := uiInstance.GetState()

	// The UI should still be functional
	uiInstance.MoveDown()
	uiInstance.MoveUp()
	uiInstance.ToggleAggregate()
	uiInstance.CycleSort()

	// State should remain consistent
	if state.ScrollOffset < 0 {
		t.Errorf("ScrollOffset should not be negative: %d", state.ScrollOffset)
	}
}
