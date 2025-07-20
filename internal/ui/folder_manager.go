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

// UpdateMainView updates the main folder manager view
func (fm *FolderManager) UpdateMainView(v *gocui.View) {
	v.Clear()
	v.Title = "Folder Manager"

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get current watched folders
	roots := fm.ui.GetRootPaths()

	_, _ = fmt.Fprintf(v, "%sCurrently Watching:%s\n", cyan("="), cyan("="))
	_, _ = fmt.Fprintf(v, "\n")

	if len(roots) == 0 {
		_, _ = fmt.Fprintf(v, "%sNo folders being watched%s\n", red(""), red(""))
	} else {
		for i, root := range roots {
			prefix := "  "
			if i == fm.ui.state.FolderManager.SelectedIdx {
				prefix = "> "
			}
			_, _ = fmt.Fprintf(v, "%s%s%s\n", prefix, green(root), "")
		}
	}

	_, _ = fmt.Fprintf(v, "\n")
	_, _ = fmt.Fprintf(v, "%sTotal: %d folder(s)%s\n", yellow(""), len(roots), yellow(""))
	_, _ = fmt.Fprintf(v, "%sWatched directories: %d%s\n", yellow(""), fm.ui.watcher.(interface{ GetWatchedCount() int }).GetWatchedCount(), yellow(""))
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
		prefix := "  "
		if fm.ui.state.FolderManager.SelectedIdx == 0 {
			prefix = "> "
		}
		_, _ = fmt.Fprintf(v, "%s%s%s\n", prefix, green(".."), "")
	}

	// List directories
	dirIndex := 1
	if parentDir != currentPath {
		dirIndex = 1
	} else {
		dirIndex = 0
	}

	for _, entry := range entries {
		if entry.IsDir() && !fm.ui.ShouldIgnore(entry.Name()) {
			prefix := "  "
			if fm.ui.state.FolderManager.SelectedIdx == dirIndex {
				prefix = "> "
			}

			dirPath := filepath.Join(currentPath, entry.Name())
			isWatching := fm.ui.watcher.(interface{ IsWatching(string) bool }).IsWatching(dirPath)

			status := ""
			if isWatching {
				status = " [WATCHING]"
			}

			_, _ = fmt.Fprintf(v, "%s%s%s%s\n", prefix, green(entry.Name()), status, "")
			dirIndex++
		}
	}
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
func (fm *FolderManager) handleMainViewInput(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case FolderManagerView:
		// Handle navigation and actions for watched folders
		return nil
	}
	return nil
}

// handleFolderListViewInput handles input for the folder list view
func (fm *FolderManager) handleFolderListViewInput(g *gocui.Gui, v *gocui.View) error {
	currentPath := fm.ui.state.FolderManager.CurrentPath
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	// Get directories
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
	if selectedIdx < 0 {
		selectedIdx = 0
	}
	if selectedIdx >= len(dirs) {
		selectedIdx = len(dirs) - 1
	}

	switch v.Name() {
	case FolderListView:
		// Navigation
		return nil
	}
	return nil
}

// NavigateUp moves selection up
func (fm *FolderManager) NavigateUp() {
	if fm.ui.state.FolderManager.SelectedIdx > 0 {
		fm.ui.state.FolderManager.SelectedIdx--
	}
}

// NavigateDown moves selection down
func (fm *FolderManager) NavigateDown() {
	// Get current directories to know the limit
	currentPath := fm.ui.state.FolderManager.CurrentPath
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return
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

	if fm.ui.state.FolderManager.SelectedIdx < len(dirs)-1 {
		fm.ui.state.FolderManager.SelectedIdx++
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

	// Update UI state
	fm.ui.rootPaths = fm.ui.GetRootPaths()

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

	// Update UI state
	fm.ui.rootPaths = fm.ui.GetRootPaths()

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
	return nil
}

// Down moves selection down in the folder list
func (fm *FolderManager) Down(g *gocui.Gui, v *gocui.View) error {
	fm.NavigateDown()
	return nil
}

// Enter selects the current item (navigate into directory)
func (fm *FolderManager) Enter(g *gocui.Gui, v *gocui.View) error {
	return fm.SelectCurrentItem()
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
