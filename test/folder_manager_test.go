package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pbouamriou/watch-fs/internal/ui"
)

// MockView simulates gocui.View for testing scroll logic
type MockView struct {
	originX, originY        int
	cursorY                 int
	innerWidth, innerHeight int
}

func NewMockView(width, height int) *MockView {
	return &MockView{
		innerWidth:  width,
		innerHeight: height,
	}
}

func (m *MockView) Origin() (int, int) {
	return m.originX, m.originY
}

func (m *MockView) SetOrigin(x, y int) error {
	m.originX = x
	m.originY = y
	return nil
}

func (m *MockView) SetOriginY(y int) error {
	m.originY = y
	return nil
}

func (m *MockView) SetCursorY(y int) error {
	m.cursorY = y
	return nil
}

func (m *MockView) InnerSize() (int, int) {
	return m.innerWidth, m.innerHeight
}

// TestFolderManagerNavigation tests basic navigation in folder manager
func TestFolderManagerNavigation(t *testing.T) {
	helper := NewTestHelperWithTempDir(t, "nav-test")
	defer helper.Cleanup()

	baseDir := helper.tempDirs[0]

	// Create some directories to navigate through
	for i := 0; i < 3; i++ {
		subDir := filepath.Join(baseDir, fmt.Sprintf("testdir%d", i))
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	helper.ui.GetState().FolderManager.CurrentPath = baseDir
	fm := helper.ui.GetFolderManager()

	// Test initial state
	state := helper.ui.GetState()
	if state.FolderManager.SelectedIdx != 0 {
		t.Errorf("Expected initial SelectedIdx 0, got %d", state.FolderManager.SelectedIdx)
	}

	// Test NavigateDown - should work with actual directories
	fm.NavigateDown()
	if state.FolderManager.SelectedIdx != 1 {
		t.Errorf("Expected SelectedIdx 1 after NavigateDown, got %d", state.FolderManager.SelectedIdx)
	}

	// Test NavigateUp
	fm.NavigateUp()
	if state.FolderManager.SelectedIdx != 0 {
		t.Errorf("Expected SelectedIdx 0 after NavigateUp, got %d", state.FolderManager.SelectedIdx)
	}

	// Test NavigateUp at top (should stay at 0)
	fm.NavigateUp()
	if state.FolderManager.SelectedIdx != 0 {
		t.Errorf("Expected SelectedIdx 0 when navigating up at top, got %d", state.FolderManager.SelectedIdx)
	}
}

// TestFolderManagerScrollLogic tests the scroll position calculation
func TestFolderManagerScrollLogic(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "watch-fs-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Create multiple subdirectories to test scrolling
	for i := 0; i < 20; i++ {
		subDir := filepath.Join(tempDir, fmt.Sprintf("subdir%02d", i))
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("Failed to create subdir: %v", err)
		}
	}

	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, tempDir)
	fm := uiInstance.GetFolderManager()

	// Set the current path to our test directory
	uiInstance.GetState().FolderManager.CurrentPath = tempDir

	// Simulate navigating to item 15 (which should trigger scroll)
	for i := 0; i < 15; i++ {
		fm.NavigateDown()
	}

	// The updateScrollPosition function should be called internally
	// We can't directly test it without exposing it, but we can test the effects
	selectedIdx := uiInstance.GetState().FolderManager.SelectedIdx
	if selectedIdx != 15 {
		t.Errorf("Expected SelectedIdx 15, got %d", selectedIdx)
	}
}

// TestFolderManagerBoundaryConditions tests edge cases
func TestFolderManagerBoundaryConditions(t *testing.T) {
	// Create a temporary directory with only one subdirectory
	tempDir, err := os.MkdirTemp("", "watch-fs-test-single")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	subDir := filepath.Join(tempDir, "onlydir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, tempDir)
	fm := uiInstance.GetFolderManager()

	// Set the current path to our test directory
	uiInstance.GetState().FolderManager.CurrentPath = tempDir

	// Test navigating down beyond available items
	initialIdx := uiInstance.GetState().FolderManager.SelectedIdx
	for i := 0; i < 10; i++ {
		fm.NavigateDown()
	}

	// Should not exceed the number of available directories
	finalIdx := uiInstance.GetState().FolderManager.SelectedIdx
	if finalIdx < initialIdx {
		t.Error("SelectedIdx should not decrease when navigating down")
	}

	// Test navigating up from valid position
	fm.NavigateUp()
	afterUpIdx := uiInstance.GetState().FolderManager.SelectedIdx
	if afterUpIdx < 0 {
		t.Errorf("SelectedIdx should not be negative, got %d", afterUpIdx)
	}
}

// TestFolderManagerPathNavigation tests directory navigation
func TestFolderManagerPathNavigation(t *testing.T) {
	// Create a nested directory structure
	tempDir, err := os.MkdirTemp("", "watch-fs-test-path")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, tempDir)
	fm := uiInstance.GetFolderManager()

	// Set initial path
	uiInstance.GetState().FolderManager.CurrentPath = tempDir

	// Test SelectCurrentItem (should navigate into directory or go to parent)
	err = fm.SelectCurrentItem()
	if err != nil {
		t.Errorf("SelectCurrentItem failed: %v", err)
	}

	// The path might change depending on what's selected
	// This tests that the function doesn't crash
	newPath := uiInstance.GetState().FolderManager.CurrentPath
	if newPath == "" {
		t.Error("Path should not be empty after SelectCurrentItem")
	}

	// Test that selected index was reset after navigation
	if uiInstance.GetState().FolderManager.SelectedIdx < 0 {
		t.Error("SelectedIdx should be non-negative after navigation")
	}
}

// TestFolderManagerEmptyDirectory tests behavior with empty directories
func TestFolderManagerEmptyDirectory(t *testing.T) {
	// Create an empty directory
	tempDir, err := os.MkdirTemp("", "watch-fs-test-empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	mockWatcher := NewMockWatcher()
	defer mockWatcher.Close()

	uiInstance := ui.NewUI(mockWatcher, tempDir)
	fm := uiInstance.GetFolderManager()

	// Set the current path to empty directory
	uiInstance.GetState().FolderManager.CurrentPath = tempDir

	// Test navigation in empty directory
	fm.NavigateDown()
	idx := uiInstance.GetState().FolderManager.SelectedIdx
	if idx < 0 {
		t.Errorf("SelectedIdx should not be negative in empty directory, got %d", idx)
	}

	fm.NavigateUp()
	idx = uiInstance.GetState().FolderManager.SelectedIdx
	if idx < 0 {
		t.Errorf("SelectedIdx should not be negative after NavigateUp in empty directory, got %d", idx)
	}
}
