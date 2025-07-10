package test

import (
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pbouamriou/watch-fs/internal/ui"
)

// MockWatcher implements the watcher interface for testing
type MockWatcher struct {
	events chan fsnotify.Event
	errors chan error
}

func NewMockWatcher() *MockWatcher {
	return &MockWatcher{
		events: make(chan fsnotify.Event, 100),
		errors: make(chan error, 10),
	}
}

func (m *MockWatcher) Events() <-chan fsnotify.Event {
	return m.events
}

func (m *MockWatcher) Errors() <-chan error {
	return m.errors
}

func (m *MockWatcher) AddDirectory(path string) error {
	return nil
}

func (m *MockWatcher) Close() {
	close(m.events)
	close(m.errors)
}

// TestEventAggregation tests the event aggregation functionality
func TestEventAggregation(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	ui := ui.NewUI(mockWatcher, "/test/path")

	// Test 1: Aggregation enabled (default)
	if !ui.GetState().AggregateEvents {
		t.Error("Aggregation should be enabled by default")
	}

	// Test 2: Add multiple similar events
	ui.AddEvent("/test/file.txt", fsnotify.Write, false)
	ui.AddEvent("/test/file.txt", fsnotify.Write, false)
	ui.AddEvent("/test/file.txt", fsnotify.Write, false)

	events := ui.GetState().Events
	if len(events) != 1 {
		t.Errorf("Expected 1 aggregated event, got %d", len(events))
	}

	if events[0].Count != 3 {
		t.Errorf("Expected count 3, got %d", events[0].Count)
	}

	// Test 3: Toggle aggregation off
	ui.ToggleAggregate()
	if ui.GetState().AggregateEvents {
		t.Error("Aggregation should be disabled after toggle")
	}

	// Should have 3 individual events now
	events = ui.GetState().Events
	if len(events) != 3 {
		t.Errorf("Expected 3 individual events, got %d", len(events))
	}

	for _, event := range events {
		if event.Count != 1 {
			t.Errorf("Expected count 1 for individual event, got %d", event.Count)
		}
	}

	// Test 4: Toggle aggregation back on
	ui.ToggleAggregate()
	if !ui.GetState().AggregateEvents {
		t.Error("Aggregation should be enabled after toggle")
	}

	// Should have 1 aggregated event again
	events = ui.GetState().Events
	if len(events) != 1 {
		t.Errorf("Expected 1 aggregated event after re-enabling, got %d", len(events))
	}

	if events[0].Count != 3 {
		t.Errorf("Expected count 3 after re-aggregation, got %d", events[0].Count)
	}
}

// TestEventFiltering tests the event filtering functionality
func TestEventFiltering(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	ui := ui.NewUI(mockWatcher, "/test/path")

	// Add mixed events
	ui.AddEvent("/test/file.txt", fsnotify.Write, false)
	ui.AddEvent("/test/dir/", fsnotify.Create, true)
	ui.AddEvent("/test/another.txt", fsnotify.Remove, false)

	// Test file filtering
	ui.ToggleFiles()
	if ui.GetState().Filter.ShowFiles {
		t.Error("Files should be hidden after toggle")
	}

	filteredEvents := ui.GetFilteredEvents()
	for _, event := range filteredEvents {
		if !event.IsDir {
			t.Error("Should only show directory events when files are hidden")
		}
	}

	// Test directory filtering
	ui.ToggleDirs()
	if ui.GetState().Filter.ShowDirs {
		t.Error("Dirs should be hidden after toggle")
	}

	filteredEvents = ui.GetFilteredEvents()
	if len(filteredEvents) != 0 {
		t.Error("Should show no events when both files and dirs are hidden")
	}
}

// TestEventSorting tests the event sorting functionality
func TestEventSorting(t *testing.T) {
	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	ui := ui.NewUI(mockWatcher, "/test/path")

	// Add events with different timestamps
	ui.AddEvent("/test/a.txt", fsnotify.Write, false)
	time.Sleep(10 * time.Millisecond)
	ui.AddEvent("/test/b.txt", fsnotify.Write, false)
	time.Sleep(10 * time.Millisecond)
	ui.AddEvent("/test/c.txt", fsnotify.Write, false)

	// Test time sorting (newest first)
	events := ui.GetFilteredEvents()
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}

	// Check that events are sorted by time (newest first)
	if events[0].Path != "/test/c.txt" {
		t.Error("First event should be the newest (c.txt)")
	}

	// Test path sorting
	ui.CycleSort() // Switch to path sorting
	events = ui.GetFilteredEvents()
	if events[0].Path != "/test/a.txt" {
		t.Error("First event should be a.txt when sorted by path")
	}
}
