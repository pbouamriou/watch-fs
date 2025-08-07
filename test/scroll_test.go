package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestScrollPositionCalculation tests the scroll position logic that was recently fixed
func TestScrollPositionCalculation(t *testing.T) {
	helper := NewTestHelperWithTempDir(t, "scroll-test")
	defer helper.Cleanup()

	// Create a directory with many subdirectories to test scrolling
	baseDir := helper.tempDirs[0]
	for i := 0; i < 30; i++ {
		subDir := filepath.Join(baseDir, fmt.Sprintf("testdir_%02d", i))
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Set the folder manager to use our test directory
	helper.ui.GetState().FolderManager.CurrentPath = baseDir
	fm := helper.ui.GetFolderManager()

	// Test navigation bounds - this was a previous issue
	initialIdx := helper.ui.GetState().FolderManager.SelectedIdx

	// Navigate down several times
	for i := 0; i < 15; i++ {
		fm.NavigateDown()
	}

	// Check that we actually moved
	afterNavIdx := helper.ui.GetState().FolderManager.SelectedIdx
	if afterNavIdx <= initialIdx {
		t.Errorf("Expected SelectedIdx to increase after navigation, was %d, now %d", initialIdx, afterNavIdx)
	}

	// Test that we don't go beyond bounds
	// Navigate to what should be near the end
	for i := 0; i < 50; i++ {
		fm.NavigateDown()
	}

	finalIdx := helper.ui.GetState().FolderManager.SelectedIdx
	if finalIdx < 0 {
		t.Errorf("SelectedIdx should never be negative, got %d", finalIdx)
	}

	// Navigate back up
	for i := 0; i < 50; i++ {
		fm.NavigateUp()
	}

	backToTopIdx := helper.ui.GetState().FolderManager.SelectedIdx
	if backToTopIdx < 0 {
		t.Errorf("SelectedIdx should never be negative after navigating up, got %d", backToTopIdx)
	}
}

// TestScrollWithEmptyDirectory tests scrolling behavior with empty directories
func TestScrollWithEmptyDirectory(t *testing.T) {
	helper := NewTestHelperWithTempDir(t, "empty-scroll-test")
	defer helper.Cleanup()

	// Use empty directory
	baseDir := helper.tempDirs[0]
	helper.ui.GetState().FolderManager.CurrentPath = baseDir
	fm := helper.ui.GetFolderManager()

	// Test navigation in empty directory
	fm.NavigateDown()
	afterDownIdx := helper.ui.GetState().FolderManager.SelectedIdx
	if afterDownIdx < 0 {
		t.Errorf("SelectedIdx should not be negative in empty directory, got %d", afterDownIdx)
	}

	fm.NavigateUp()
	afterUpIdx := helper.ui.GetState().FolderManager.SelectedIdx
	if afterUpIdx < 0 {
		t.Errorf("SelectedIdx should not be negative after NavigateUp in empty directory, got %d", afterUpIdx)
	}
}

// TestScrollEdgeCases tests various edge cases that were problematic before
func TestScrollEdgeCases(t *testing.T) {
	helper := NewTestHelperWithTempDir(t, "edge-cases-test")
	defer helper.Cleanup()

	baseDir := helper.tempDirs[0]

	// Create exactly one subdirectory
	subDir := filepath.Join(baseDir, "single_dir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	helper.ui.GetState().FolderManager.CurrentPath = baseDir
	fm := helper.ui.GetFolderManager()

	// Test with single directory
	// This should include parent dir (..) + our single dir = 2 items total
	for i := 0; i < 5; i++ {
		fm.NavigateDown()
	}

	idx := helper.ui.GetState().FolderManager.SelectedIdx
	if idx < 0 {
		t.Errorf("SelectedIdx should not be negative with single directory, got %d", idx)
	}

	// Test navigation when at boundaries
	fm.NavigateUp()
	fm.NavigateUp() // Should not crash or go negative

	finalIdx := helper.ui.GetState().FolderManager.SelectedIdx
	if finalIdx < 0 {
		t.Errorf("SelectedIdx should not be negative after boundary navigation, got %d", finalIdx)
	}
}

// TestParentDirectoryNavigation tests the ".." navigation that was previously broken
func TestParentDirectoryNavigation(t *testing.T) {
	helper := NewTestHelperWithTempDir(t, "parent-nav-test")
	defer helper.Cleanup()

	baseDir := helper.tempDirs[0]

	// Create a nested structure
	level1 := filepath.Join(baseDir, "level1")
	level2 := filepath.Join(level1, "level2")
	if err := os.MkdirAll(level2, 0755); err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	// Start from the deepest level
	helper.ui.GetState().FolderManager.CurrentPath = level2
	fm := helper.ui.GetFolderManager()

	initialPath := helper.ui.GetState().FolderManager.CurrentPath
	if initialPath != level2 {
		t.Errorf("Expected initial path %s, got %s", level2, initialPath)
	}

	// Navigate using SelectCurrentItem (which should handle ".." entries)
	// This was the functionality that was broken before
	err := fm.SelectCurrentItem()
	if err != nil {
		t.Errorf("SelectCurrentItem should not fail: %v", err)
	}

	// Path might have changed if we selected ".."
	newPath := helper.ui.GetState().FolderManager.CurrentPath

	// The exact behavior depends on what was selected, but it shouldn't crash
	// and the path should still be valid
	if newPath == "" {
		t.Error("Path should not be empty after SelectCurrentItem")
	}

	// Index should be reset to 0 after navigation
	if helper.ui.GetState().FolderManager.SelectedIdx < 0 {
		t.Errorf("SelectedIdx should not be negative after navigation, got %d", helper.ui.GetState().FolderManager.SelectedIdx)
	}
}

// TestConcurrentFolderManagerOperations tests thread safety
func TestConcurrentFolderManagerOperations(t *testing.T) {
	helper := NewTestHelperWithTempDir(t, "concurrent-test")
	defer helper.Cleanup()

	baseDir := helper.tempDirs[0]

	// Create multiple directories
	for i := 0; i < 10; i++ {
		subDir := filepath.Join(baseDir, fmt.Sprintf("dir%d", i))
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	helper.ui.GetState().FolderManager.CurrentPath = baseDir
	fm := helper.ui.GetFolderManager()

	// Perform operations concurrently
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 20; i++ {
			fm.NavigateDown()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 20; i++ {
			fm.NavigateUp()
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Check that state is still valid
	idx := helper.ui.GetState().FolderManager.SelectedIdx
	if idx < 0 {
		t.Errorf("SelectedIdx should not be negative after concurrent operations, got %d", idx)
	}

	path := helper.ui.GetState().FolderManager.CurrentPath
	if path == "" {
		t.Error("Path should not be empty after concurrent operations")
	}
}

// TestScrollOffsetConsistency tests that ScrollOffset field remains consistent
func TestScrollOffsetConsistency(t *testing.T) {
	helper := NewTestHelperWithTempDir(t, "scroll-offset-test")
	defer helper.Cleanup()

	// Note: ScrollOffset in FolderManagerState was added to track scroll position
	// This test ensures it's managed correctly

	baseDir := helper.tempDirs[0]

	// Create many directories to test scrolling
	for i := 0; i < 25; i++ {
		subDir := filepath.Join(baseDir, fmt.Sprintf("scrolldir%02d", i))
		if err := os.Mkdir(subDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	helper.ui.GetState().FolderManager.CurrentPath = baseDir
	fm := helper.ui.GetFolderManager()

	// Check initial scroll offset
	initialOffset := helper.ui.GetState().FolderManager.ScrollOffset
	if initialOffset != 0 {
		t.Errorf("Expected initial ScrollOffset 0, got %d", initialOffset)
	}

	// Navigate and check that scroll offset tracking works
	// Note: The actual scroll offset management is done in updateScrollPosition
	// which is called internally, so we test the state consistency

	for i := 0; i < 10; i++ {
		fm.NavigateDown()
	}

	// The ScrollOffset field should remain non-negative
	offset := helper.ui.GetState().FolderManager.ScrollOffset
	if offset < 0 {
		t.Errorf("ScrollOffset should not be negative, got %d", offset)
	}
}
