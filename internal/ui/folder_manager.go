package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

// FolderManager handles the folder management interface
type FolderManager struct {
	ui *UI
}

// NewFolderManager creates a new folder manager
func NewFolderManager(ui *UI) *FolderManager {
	return &FolderManager{ui: ui}
}

// Show displays the folder management interface
func (fm *FolderManager) Show() {
	fm.ui.state.CurrentFocus = FocusFolderManager
	fm.ui.state.ShowFolderManager = true
}

// Hide hides the folder management interface
func (fm *FolderManager) Hide() {
	fm.ui.state.CurrentFocus = FocusMain
	fm.ui.state.ShowFolderManager = false
}

// UpdateMainView updates the main folder manager view (background frame)
func (fm *FolderManager) UpdateMainView(v *gocui.View) {
	v.Clear()
	// This is now just the background frame - content is in sub-views
}

// UpdateWatchedFoldersView updates the "Currently Watching" sub-view (left panel)
func (fm *FolderManager) UpdateWatchedFoldersView(v *gocui.View) {
	v.Clear()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	// Get current watched folders directly from watcher
	roots := fm.getRealWatchedRoots()

	if len(roots) == 0 {
		_, _ = fmt.Fprintf(v, "\n%s   👁️  No folders being watched%s\n", red(""), red(""))
		_, _ = fmt.Fprintf(v, "%s   💡 Use Ctrl+F to add folders%s\n", yellow(""), yellow(""))
		return
	}

	// Display watched folders as simple list (use standard cursor highlighting)
	for _, root := range roots {
		// Get base name for display
		baseName := filepath.Base(root)
		if baseName == "." || baseName == "/" {
			baseName = root
		}

		// Get watched count for this specific root if available
		watchedCount := ""
		if counter, ok := fm.ui.watcher.(interface{ GetWatchedCountForRoot(string) int }); ok {
			rootWatched := counter.GetWatchedCountForRoot(root)
			if rootWatched > 0 {
				watchedCount = fmt.Sprintf(" %s(%d)%s", blue(""), rootWatched, blue(""))
			}
		}

		// Simple display format suitable for cursor highlighting
		_, _ = fmt.Fprintf(v, "  %s%s\n", magenta(baseName), green(watchedCount))
	}

	// Add some spacing and info
	_, _ = fmt.Fprintf(v, "\n")
	_, _ = fmt.Fprintf(v, "%s--- Stats ---%s\n", cyan(""), cyan(""))
	_, _ = fmt.Fprintf(v, " Roots: %s%d%s\n", yellow(""), len(roots), yellow(""))

	if counter, ok := fm.ui.watcher.(interface{ GetWatchedCount() int }); ok {
		totalWatched := counter.GetWatchedCount()
		_, _ = fmt.Fprintf(v, " Total: %s%d%s\n", yellow(""), totalWatched, yellow(""))
	}

	_, _ = fmt.Fprintf(v, "\n")
	_, _ = fmt.Fprintf(v, "%s--- Keys ---%s\n", cyan(""), cyan(""))
	_, _ = fmt.Fprintf(v, " %sUp/Down%s Nav %sR%s Del\n", blue(""), blue(""), blue(""), blue(""))

	// Set cursor position based on WatchedIdx
	if len(roots) > 0 {
		watchedIdx := fm.ui.state.FolderManager.WatchedIdx
		if watchedIdx >= len(roots) {
			watchedIdx = len(roots) - 1
			fm.ui.state.FolderManager.WatchedIdx = watchedIdx
		}
		if watchedIdx < 0 {
			watchedIdx = 0
			fm.ui.state.FolderManager.WatchedIdx = watchedIdx
		}
		v.SetCursor(0, watchedIdx)
	}
}

// UpdateFolderListView updates the folder selection view
func (fm *FolderManager) UpdateFolderListView(v *gocui.View) {
	v.Clear()
	v.Title = "Available Folders"

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	currentPath := fm.ui.state.FolderManager.CurrentPath
	_, _ = fmt.Fprintf(v, "%sCurrent Path: %s%s\n", cyan(""), currentPath, cyan(""))
	_, _ = fmt.Fprintf(v, "\n")

	// List directories in current path
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		_, _ = fmt.Fprintf(v, "%sError reading directory: %v%s\n", red(""), err, red(""))
		return
	}

	// Add parent directory option
	parentDir := filepath.Dir(currentPath)
	if parentDir != currentPath {
		_, _ = fmt.Fprintf(v, "  %s\n", green(".."))
	}

	// List directories
	// (dirIndex calculation removed as it was unused)

	for _, entry := range entries {
		if entry.IsDir() && !fm.ui.ShouldIgnore(entry.Name()) {
			dirPath := filepath.Join(currentPath, entry.Name())
			isWatching := fm.ui.watcher.(interface{ IsWatching(string) bool }).IsWatching(dirPath)

			status := ""
			if isWatching {
				status = " [WATCHING]"
			}

			_, _ = fmt.Fprintf(v, "  %s%s\n", green(entry.Name()), status)
		}
	}

	// Update scroll position to ensure selected item is visible
	fm.updateScrollPosition(v)
}

// UpdatePathView updates the path display view
func (fm *FolderManager) UpdatePathView(v *gocui.View) {
	v.Clear()
	v.Title = "Path"

	cyan := color.New(color.FgCyan).SprintFunc()
	_, _ = fmt.Fprintf(v, "%s%s%s", cyan(""), fm.ui.state.FolderManager.CurrentPath, cyan(""))
}

// HandleInput handles keyboard input for the folder manager
func (fm *FolderManager) HandleInput(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case FolderManagerView:
		return fm.handleMainViewInput(g, v)
	case FolderListView:
		return fm.handleFolderListViewInput(g, v)
	}
	return nil
}

// handleMainViewInput handles input for the main folder manager view
func (fm *FolderManager) handleMainViewInput(_ *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case FolderManagerView:
		// Handle navigation and actions for watched folders
		return nil
	}
	return nil
}

// handleFolderListViewInput handles input for the folder list view
func (fm *FolderManager) handleFolderListViewInput(_ *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case FolderListView:
		// Navigation logic would go here if needed
		return nil
	}
	return nil
}

// NavigateUp moves selection up
func (fm *FolderManager) NavigateUp() {
	if fm.ui.state.FolderManager.SelectedIdx > 0 {
		fm.ui.state.FolderManager.SelectedIdx--
	}

	// Ensure the selected index stays within bounds
	totalItems := fm.getTotalDirectories()
	if fm.ui.state.FolderManager.SelectedIdx >= totalItems {
		fm.ui.state.FolderManager.SelectedIdx = totalItems - 1
	}
	if fm.ui.state.FolderManager.SelectedIdx < 0 {
		fm.ui.state.FolderManager.SelectedIdx = 0
	}
}

// NavigateDown moves selection down
func (fm *FolderManager) NavigateDown() {
	totalItems := fm.getTotalDirectories()
	if totalItems == 0 {
		return
	}

	if fm.ui.state.FolderManager.SelectedIdx < totalItems-1 {
		fm.ui.state.FolderManager.SelectedIdx++
	}

	// Ensure the selected index stays within bounds
	if fm.ui.state.FolderManager.SelectedIdx >= totalItems {
		fm.ui.state.FolderManager.SelectedIdx = totalItems - 1
	}
	if fm.ui.state.FolderManager.SelectedIdx < 0 {
		fm.ui.state.FolderManager.SelectedIdx = 0
	}
}

// SelectCurrentItem selects the currently highlighted item
func (fm *FolderManager) SelectCurrentItem() error {
	currentPath := fm.ui.state.FolderManager.CurrentPath
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	var dirs []string
	parentDir := filepath.Dir(currentPath)
	if parentDir != currentPath {
		dirs = append(dirs, "..")
	}

	for _, entry := range entries {
		if entry.IsDir() && !fm.ui.ShouldIgnore(entry.Name()) {
			dirs = append(dirs, entry.Name())
		}
	}

	selectedIdx := fm.ui.state.FolderManager.SelectedIdx
	if selectedIdx < 0 || selectedIdx >= len(dirs) {
		return nil
	}

	selectedDir := dirs[selectedIdx]

	if selectedDir == ".." {
		// Navigate to parent directory
		fm.ui.state.FolderManager.CurrentPath = parentDir
		fm.ui.state.FolderManager.SelectedIdx = 0
	} else {
		// Navigate to selected directory
		newPath := filepath.Join(currentPath, selectedDir)
		fm.ui.state.FolderManager.CurrentPath = newPath
		fm.ui.state.FolderManager.SelectedIdx = 0
	}

	return nil
}

// AddCurrentFolder adds the currently selected folder to watch
func (fm *FolderManager) AddCurrentFolder() error {
	currentPath := fm.ui.state.FolderManager.CurrentPath
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	var dirs []string
	parentDir := filepath.Dir(currentPath)
	if parentDir != currentPath {
		dirs = append(dirs, "..")
	}

	for _, entry := range entries {
		if entry.IsDir() && !fm.ui.ShouldIgnore(entry.Name()) {
			dirs = append(dirs, entry.Name())
		}
	}

	selectedIdx := fm.ui.state.FolderManager.SelectedIdx
	if selectedIdx < 0 || selectedIdx >= len(dirs) {
		return nil
	}

	selectedDir := dirs[selectedIdx]

	if selectedDir == ".." {
		return nil // Don't add parent directory
	}

	// Add the selected directory to watch
	folderPath := filepath.Join(currentPath, selectedDir)

	// Check if already watching
	if fm.ui.watcher.(interface{ IsWatching(string) bool }).IsWatching(folderPath) {
		return nil // Already watching
	}

	// Add to watcher
	if err := fm.ui.watcher.(interface{ AddRoot(string) error }).AddRoot(folderPath); err != nil {
		return err
	}

	// Update UI state with real watcher roots
	if multiRootWatcher, ok := fm.ui.watcher.(interface{ GetRoots() []string }); ok {
		fm.ui.rootPaths = multiRootWatcher.GetRoots()
	}

	// Close the folder manager popup after successful addition
	fm.Hide()

	return nil
}

// RemoveSelectedFolder removes the selected folder from watch
func (fm *FolderManager) RemoveSelectedFolder() error {
	roots := fm.ui.GetRootPaths()
	selectedIdx := fm.ui.state.FolderManager.SelectedIdx

	if selectedIdx < 0 || selectedIdx >= len(roots) {
		return nil
	}

	folderToRemove := roots[selectedIdx]

	// Remove from watcher
	if err := fm.ui.watcher.(interface{ RemoveRoot(string) error }).RemoveRoot(folderToRemove); err != nil {
		return err
	}

	// Update UI state with real watcher roots
	if multiRootWatcher, ok := fm.ui.watcher.(interface{ GetRoots() []string }); ok {
		fm.ui.rootPaths = multiRootWatcher.GetRoots()
	}

	// Adjust selection index
	if selectedIdx >= len(fm.ui.rootPaths) {
		fm.ui.state.FolderManager.SelectedIdx = len(fm.ui.rootPaths) - 1
	}
	if fm.ui.state.FolderManager.SelectedIdx < 0 {
		fm.ui.state.FolderManager.SelectedIdx = 0
	}

	return nil
}

// Up moves selection up in the folder list
func (fm *FolderManager) Up(g *gocui.Gui, v *gocui.View) error {
	fm.NavigateUp()
	fm.updateScrollPosition(v)
	return nil
}

// Down moves selection down in the folder list
func (fm *FolderManager) Down(g *gocui.Gui, v *gocui.View) error {
	fm.NavigateDown()
	fm.updateScrollPosition(v)
	return nil
}

// Enter selects the current item (navigate into directory)
func (fm *FolderManager) Enter(g *gocui.Gui, v *gocui.View) error {
	err := fm.SelectCurrentItem()
	if err == nil {
		// After navigating to a new directory, update scroll position
		fm.updateScrollPosition(v)
	}
	return err
}

// Cancel closes the folder manager
func (fm *FolderManager) Cancel(g *gocui.Gui, v *gocui.View) error {
	fm.Hide()
	return nil
}

// Add adds the currently selected folder to watch
func (fm *FolderManager) Add(g *gocui.Gui, v *gocui.View) error {
	return fm.AddCurrentFolder()
}

// Remove removes the selected folder from watch
func (fm *FolderManager) Remove(g *gocui.Gui, v *gocui.View) error {
	return fm.RemoveSelectedFolder()
}

// updateScrollPosition adjusts the scroll position to keep the selected item visible
func (fm *FolderManager) updateScrollPosition(v *gocui.View) {
	totalItems := fm.getTotalDirectories()
	if totalItems == 0 {
		return
	}

	selectedIdx := fm.ui.state.FolderManager.SelectedIdx
	if selectedIdx >= totalItems {
		selectedIdx = totalItems - 1
		fm.ui.state.FolderManager.SelectedIdx = selectedIdx
	}
	if selectedIdx < 0 {
		selectedIdx = 0
		fm.ui.state.FolderManager.SelectedIdx = selectedIdx
	}

	_, viewHeight := v.InnerSize()
	if viewHeight <= 0 {
		viewHeight = 10
	}

	effectiveHeight := viewHeight - 2
	if effectiveHeight <= 0 {
		effectiveHeight = 1
	}

	_, currentY := v.Origin()
	visibleStart := currentY
	visibleEnd := currentY + effectiveHeight - 1

	if selectedIdx < visibleStart {
		v.SetOrigin(0, selectedIdx)
	} else if selectedIdx > visibleEnd {
		newOrigin := max(selectedIdx-effectiveHeight+1, 0)
		v.SetOriginY(newOrigin)
	}

	_, newOriginY := v.Origin()
	cursorY := selectedIdx - newOriginY + 2
	if cursorY >= 0 && cursorY < viewHeight {
		v.SetCursorY(cursorY)
	}
}

// getTotalDirectories returns the total number of directory entries (including ".." if applicable)
func (fm *FolderManager) getTotalDirectories() int {
	currentPath := fm.ui.state.FolderManager.CurrentPath
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return 0
	}

	var totalItems int
	parentDir := filepath.Dir(currentPath)
	if parentDir != currentPath {
		totalItems = 1 // Parent directory ".."
	}

	for _, entry := range entries {
		if entry.IsDir() && !fm.ui.ShouldIgnore(entry.Name()) {
			totalItems++
		}
	}

	return totalItems
}

// Helper function for max since Go doesn't have built-in max for int
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getRealWatchedRoots gets the actual watched roots from the watcher
func (fm *FolderManager) getRealWatchedRoots() []string {
	if multiRootWatcher, ok := fm.ui.watcher.(interface{ GetRoots() []string }); ok {
		roots := multiRootWatcher.GetRoots()
		return roots
	}
	return []string{fm.ui.rootPath}
}

// NavigateWatchedUp moves selection up in the "Currently Watching" list
func (fm *FolderManager) NavigateWatchedUp(g *gocui.Gui, v *gocui.View) error {
	roots := fm.getRealWatchedRoots()
	if len(roots) == 0 {
		return nil
	}

	// Ensure we're operating on the watched folders list, not the directory browser
	if fm.ui.state.FolderManager.WatchedIdx > 0 {
		fm.ui.state.FolderManager.WatchedIdx--
	}

	return nil
}

// NavigateWatchedDown moves selection down in the "Currently Watching" list
func (fm *FolderManager) NavigateWatchedDown(g *gocui.Gui, v *gocui.View) error {
	roots := fm.getRealWatchedRoots()
	if len(roots) == 0 {
		return nil
	}

	// Ensure we're operating on the watched folders list, not the directory browser
	if fm.ui.state.FolderManager.WatchedIdx < len(roots)-1 {
		fm.ui.state.FolderManager.WatchedIdx++
	}

	return nil
}

// RemoveWatchedFolder removes the currently selected watched folder
func (fm *FolderManager) RemoveWatchedFolder(g *gocui.Gui, v *gocui.View) error {
	roots := fm.getRealWatchedRoots()
	if len(roots) == 0 {
		return nil
	}

	selectedIdx := fm.ui.state.FolderManager.WatchedIdx
	if selectedIdx < 0 || selectedIdx >= len(roots) {
		return nil
	}

	folderToRemove := roots[selectedIdx]

	// Remove from watcher
	if err := fm.ui.watcher.(interface{ RemoveRoot(string) error }).RemoveRoot(folderToRemove); err != nil {
		return err
	}

	// Update UI state with real watcher roots
	if multiRootWatcher, ok := fm.ui.watcher.(interface{ GetRoots() []string }); ok {
		fm.ui.rootPaths = multiRootWatcher.GetRoots()
	}

	// Adjust selection index to stay within bounds
	newRoots := fm.getRealWatchedRoots()
	if selectedIdx >= len(newRoots) {
		fm.ui.state.FolderManager.WatchedIdx = len(newRoots) - 1
	}
	if fm.ui.state.FolderManager.WatchedIdx < 0 {
		fm.ui.state.FolderManager.WatchedIdx = 0
	}

	return nil
}

// SwitchToNextPanel switches focus to the next panel (Tab key)
func (fm *FolderManager) SwitchToNextPanel(g *gocui.Gui, v *gocui.View) error {
	switch fm.ui.state.FolderManager.ActivePanel {
	case FocusWatchedFolders:
		fm.ui.state.FolderManager.ActivePanel = FocusFolderBrowser
		fm.ui.state.CurrentFocus = FocusFolderBrowser
		if _, err := g.SetCurrentView(FolderListView); err != nil {
			return err
		}
	case FocusFolderBrowser:
		fm.ui.state.FolderManager.ActivePanel = FocusWatchedFolders
		fm.ui.state.CurrentFocus = FocusWatchedFolders
		if _, err := g.SetCurrentView("watched_folders"); err != nil {
			return err
		}
	}
	return nil
}

// SwitchToPreviousPanel switches focus to the previous panel (Shift+Tab key)
func (fm *FolderManager) SwitchToPreviousPanel(g *gocui.Gui, v *gocui.View) error {
	switch fm.ui.state.FolderManager.ActivePanel {
	case FocusWatchedFolders:
		fm.ui.state.FolderManager.ActivePanel = FocusFolderBrowser
		fm.ui.state.CurrentFocus = FocusFolderBrowser
		if _, err := g.SetCurrentView(FolderListView); err != nil {
			return err
		}
	case FocusFolderBrowser:
		fm.ui.state.FolderManager.ActivePanel = FocusWatchedFolders
		fm.ui.state.CurrentFocus = FocusWatchedFolders
		if _, err := g.SetCurrentView("watched_folders"); err != nil {
			return err
		}
	}
	return nil
}

// SwitchToRightPanel switches to right panel (Arrow Right key)
func (fm *FolderManager) SwitchToRightPanel(g *gocui.Gui, v *gocui.View) error {
	if fm.ui.state.FolderManager.ActivePanel == FocusWatchedFolders {
		fm.ui.state.FolderManager.ActivePanel = FocusFolderBrowser
		fm.ui.state.CurrentFocus = FocusFolderBrowser
		if _, err := g.SetCurrentView(FolderListView); err != nil {
			return err
		}
	}
	return nil
}

// SwitchToLeftPanel switches to left panel (Arrow Left key)
func (fm *FolderManager) SwitchToLeftPanel(g *gocui.Gui, v *gocui.View) error {
	if fm.ui.state.FolderManager.ActivePanel == FocusFolderBrowser {
		fm.ui.state.FolderManager.ActivePanel = FocusWatchedFolders
		fm.ui.state.CurrentFocus = FocusWatchedFolders
		if _, err := g.SetCurrentView("watched_folders"); err != nil {
			return err
		}
	}
	return nil
}
