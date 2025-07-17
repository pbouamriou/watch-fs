package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/pbouamriou/watch-fs/pkg/logger"
)

const (
	fileNameLabel = "Filename: "
)

// FileDialog manages the file selection dialog
type FileDialog struct {
	ui *UI
}

// NewFileDialog creates a new FileDialog instance
func NewFileDialog(ui *UI) *FileDialog {
	return &FileDialog{ui: ui}
}

// Show opens the file dialog
func (fd *FileDialog) Show(mode FileDialogMode, filter string) {
	fd.ui.state.ShowFileDialog = true
	fd.ui.state.FileDialog.Mode = mode
	fd.ui.state.FileDialog.Filter = filter
	fd.ui.state.CurrentFocus = FocusFileDialog
	if err := fd.loadDirectory("."); err != nil {
		logger.Error(err, "loadDirectory error")
	}
	if fd.ui.gui != nil {
		fd.ui.gui.Update(func(g *gocui.Gui) error {
			_, err := g.SetCurrentView(FileListView)
			return err
		})
	}
}

// Hide closes the file dialog
func (fd *FileDialog) Hide() {
	fd.ui.state.ShowFileDialog = false
	fd.ui.state.FileDialog.Files = make([]*FileEntry, 0)
	fd.ui.state.FileDialog.SelectedIdx = 0
	fd.ui.state.CurrentFocus = FocusMain
	if fd.ui.gui != nil {
		fd.ui.gui.Update(func(g *gocui.Gui) error {
			_, err := g.SetCurrentView(EventsView)
			return err
		})
	}
}

// loadDirectory loads the contents of a directory for the file dialog
func (fd *FileDialog) loadDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var files []*FileEntry

	// Add parent directory entry (..)
	// Always add ".." except when we're at the root filesystem
	if path != "/" {
		// For "." we need to get the actual parent directory
		var parentPath string
		if path == "." {
			// Get current working directory and then its parent
			if cwd, err := os.Getwd(); err == nil {
				parentPath = filepath.Dir(cwd)
			} else {
				parentPath = ".."
			}
		} else {
			parentPath = filepath.Dir(path)
		}

		files = append(files, &FileEntry{
			Name:    "..",
			IsDir:   true,
			Path:    parentPath,
			ModTime: time.Now(),
		})
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Skip hidden files (except . and ..)
		if strings.HasPrefix(entry.Name(), ".") && entry.Name() != "." && entry.Name() != ".." {
			continue
		}

		// Apply filter for files
		if !entry.IsDir() {
			if fd.ui.state.FileDialog.Filter != "" && fd.ui.state.FileDialog.Filter != "*" {
				matched, _ := filepath.Match(fd.ui.state.FileDialog.Filter, entry.Name())
				if !matched {
					continue
				}
			}
		}

		fullPath := filepath.Join(path, entry.Name())
		files = append(files, &FileEntry{
			Name:    entry.Name(),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			Path:    fullPath,
		})
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return files[i].Name < files[j].Name
	})

	fd.ui.state.FileDialog.Files = files
	fd.ui.state.FileDialog.CurrentPath = path
	fd.ui.state.FileDialog.SelectedIdx = 0

	return nil
}

// Up moves selection up in file dialog
func (fd *FileDialog) Up(g *gocui.Gui, v *gocui.View) error {
	if fd.ui.state.FileDialog.SelectedIdx > 0 {
		fd.ui.state.FileDialog.SelectedIdx--
	}
	return fd.ui.layout.Layout(g)
}

// Down moves selection down in file dialog
func (fd *FileDialog) Down(g *gocui.Gui, v *gocui.View) error {
	if fd.ui.state.FileDialog.SelectedIdx < len(fd.ui.state.FileDialog.Files)-1 {
		fd.ui.state.FileDialog.SelectedIdx++
	}
	return fd.ui.layout.Layout(g)
}

// Enter handles Enter key in file dialog
func (fd *FileDialog) Enter(g *gocui.Gui, v *gocui.View) error {
	if len(fd.ui.state.FileDialog.Files) == 0 {
		return nil
	}

	selected := fd.ui.state.FileDialog.Files[fd.ui.state.FileDialog.SelectedIdx]

	if selected.IsDir {
		// Navigate into directory
		if selected.Name == ".." {
			// Go to parent directory
			var parentPath string
			if fd.ui.state.FileDialog.CurrentPath == "." {
				// Get current working directory and then its parent
				if cwd, err := os.Getwd(); err == nil {
					parentPath = filepath.Dir(cwd)
				} else {
					parentPath = ".."
				}
			} else {
				parentPath = filepath.Dir(fd.ui.state.FileDialog.CurrentPath)
			}
			if err := fd.loadDirectory(parentPath); err != nil {
				logger.Error(err, "loadDirectory error")
			}
			// Force layout update after navigation
			return fd.ui.layout.Layout(g)
		} else {
			// Go into directory
			if err := fd.loadDirectory(selected.Path); err != nil {
				logger.Error(err, "loadDirectory error")
			}
			// Force layout update after navigation
			return fd.ui.layout.Layout(g)
		}
	} else {
		// File selected
		if fd.ui.state.FileDialog.Mode == ModeOpen {
			// Import mode - load the file
			format := FormatSQLite
			if strings.HasSuffix(selected.Path, ".json") {
				format = FormatJSON
			}

			err := fd.ui.ImportEvents(selected.Path, format)
			if err != nil {
				// Could show error, but for now just hide dialog
				fd.Hide()
			} else {
				fd.Hide()
				// Update UI
				if v, err := g.View(StatusView); err == nil {
					fd.ui.views.UpdateStatusView(v)
				}
				if v, err := g.View(EventsView); err == nil {
					fd.ui.views.UpdateEventsView(v)
				}
			}
		} else {
			// Save mode - offer to edit filename or use selected file
			if fd.ui.state.FileDialog.IsEditing {
				// If we're editing, use the selected file as base name
				fd.ui.state.FileDialog.Filename = selected.Name
				// Switch to filename editing mode
				if fd.ui.gui != nil {
					fd.ui.gui.Update(func(g *gocui.Gui) error {
						_, err := g.SetCurrentView("filename")
						return err
					})
				}
			} else {
				// Ask user if they want to edit the filename
				// For now, just use the selected file
				format := FormatSQLite
				if strings.HasSuffix(selected.Path, ".json") {
					format = FormatJSON
				}

				err := fd.ui.ExportEvents(selected.Path, format)
				if err != nil {
					// Could show error, but for now just hide dialog
					fd.Hide()
				} else {
					fd.Hide()
					// Update UI
					if v, err := g.View(StatusView); err == nil {
						fd.ui.views.UpdateStatusView(v)
					}
				}
			}
		}
	}

	return fd.ui.layout.Layout(g)
}

// EditFilename switches to filename editing mode
func (fd *FileDialog) EditFilename(g *gocui.Gui, v *gocui.View) error {
	if fd.ui.state.FileDialog.Mode == ModeSave {
		fd.ui.state.FileDialog.IsEditing = true
		if fd.ui.gui != nil {
			fd.ui.gui.Update(func(g *gocui.Gui) error {
				return fd.ui.layout.Layout(g)
			})
		}
	}
	return nil
}

// Cancel cancels the file dialog
func (fd *FileDialog) Cancel(g *gocui.Gui, v *gocui.View) error {
	fd.Hide()
	return fd.ui.layout.Layout(g)
}

// UpdateMainView updates the main file dialog view
func (fd *FileDialog) UpdateMainView(v *gocui.View) {
	v.Clear()
	// This view is mainly for the frame, content is in other views
}

// UpdatePathView updates the path display view
func (fd *FileDialog) UpdatePathView(v *gocui.View) {
	v.Clear()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	mode := "Save"
	if fd.ui.state.FileDialog.Mode == ModeOpen {
		mode = "Open"
	}

	_, _ = fmt.Fprintf(v, "%s: %s\n", yellow(mode), cyan(fd.ui.state.FileDialog.CurrentPath))
	_, _ = fmt.Fprintf(v, "Filter: %s\n", yellow(fd.ui.state.FileDialog.Filter))
}

// UpdateFileListView updates the file list view
func (fd *FileDialog) UpdateFileListView(v *gocui.View) {
	v.Clear()

	if len(fd.ui.state.FileDialog.Files) == 0 {
		_, _ = fmt.Fprintf(v, "No files found")
		return
	}

	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	for i, file := range fd.ui.state.FileDialog.Files {
		// Highlight selected item
		if i == fd.ui.state.FileDialog.SelectedIdx {
			_, _ = fmt.Fprintf(v, "> ")
		} else {
			_, _ = fmt.Fprintf(v, "  ")
		}

		// Format file/directory name
		if file.IsDir {
			if file.Name == ".." {
				_, _ = fmt.Fprintf(v, "%s\n", blue(file.Name))
			} else {
				_, _ = fmt.Fprintf(v, "%s/\n", blue(file.Name))
			}
		} else {
			// Show file size for files
			sizeStr := formatFileSize(file.Size)
			_, _ = fmt.Fprintf(v, "%s %s\n", file.Name, yellow(sizeStr))
		}
	}

	// Set cursor position to selected item and ensure it's visible
	selectedIdx := fd.ui.state.FileDialog.SelectedIdx
	if selectedIdx >= len(fd.ui.state.FileDialog.Files) {
		selectedIdx = len(fd.ui.state.FileDialog.Files) - 1
	}
	if selectedIdx < 0 {
		selectedIdx = 0
	}

	// Get view dimensions to calculate scroll
	_, viewHeight := v.InnerSize()
	if viewHeight <= 0 {
		viewHeight = 10 // Default height if we can't get it
	}

	// Calculate the visible range
	_, currentY := v.Origin()
	visibleStart := currentY
	visibleEnd := currentY + viewHeight - 1

	// If selected item is not visible, adjust scroll
	if selectedIdx < visibleStart {
		// Item is above visible area, scroll up
		v.SetOrigin(0, selectedIdx)
	} else if selectedIdx > visibleEnd {
		// Item is below visible area, scroll down
		newOrigin := max(selectedIdx-viewHeight+1, 0)
		v.SetOriginY(newOrigin)
	}

	// Set cursor to selected item (relative to origin)
	cursorY := max(selectedIdx, 0)
	v.SetCursorY(cursorY - visibleStart)
}

// UpdateFilenameView updates the filename input view
func (fd *FileDialog) UpdateFilenameView(v *gocui.View) {
	v.Clear()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Show current filename or default
	filename := fd.ui.state.FileDialog.Filename
	if filename == "" {
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename = fmt.Sprintf("watch-fs-events_%s.db", timestamp)
		fd.ui.state.FileDialog.Filename = filename
		fd.ui.state.FileDialog.Placeholder = true
	}

	// In editing mode, show only the filename without label
	if fd.ui.state.FileDialog.IsEditing {
		_, _ = fmt.Fprintf(v, "%s", filename)
		if fd.ui.state.FileDialog.Placeholder {
			v.SetCursor(0, 0)
		} else {
			// Position cursor at the end of the filename
			v.SetCursor(len(filename), 0)
		}
	} else {
		// In display mode, show with label
		_, _ = fmt.Fprintf(v, "%s%s%s%s", cyan(""), fileNameLabel, yellow(filename), cyan(""))
	}
}

func (fd *FileDialog) getFilenameFromView(v *gocui.View) string {
	if fd.ui.state.FileDialog.IsEditing {
		// In editing mode, the buffer contains only the filename
		return strings.TrimSpace(v.Buffer())
	} else {
		// In display mode, extract filename from label
		return strings.TrimSpace(strings.Replace(v.Buffer(), fileNameLabel, "", 1))
	}
}

// FilenameEditor handles editing in the filename input
func (fd *FileDialog) FilenameEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		gocui.DefaultEditor.Edit(v, key, ch, mod)
		fd.ui.state.FileDialog.Placeholder = false
	case key == gocui.KeySpace:
		gocui.DefaultEditor.Edit(v, key, ' ', mod)
		fd.ui.state.FileDialog.Placeholder = false
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		gocui.DefaultEditor.Edit(v, key, 0, mod)
		fd.ui.state.FileDialog.Placeholder = false
	case key == gocui.KeyDelete:
		gocui.DefaultEditor.Edit(v, key, 0, mod)
		fd.ui.state.FileDialog.Placeholder = false
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
		fd.ui.state.FileDialog.Placeholder = false
	case key == gocui.KeyEnter:
		// Get the filename from the view
		filename := fd.getFilenameFromView(v)
		if filename == "" {
			// Use default if empty
			timestamp := time.Now().Format("2006-01-02_15-04-05")
			filename = fmt.Sprintf("watch-fs-events_%s.db", timestamp)
		}

		// Build full path
		fullPath := filepath.Join(fd.ui.state.FileDialog.CurrentPath, filename)

		// Determine format based on extension
		var format ExportFormat
		if strings.HasSuffix(filename, ".json") {
			format = FormatJSON
		} else {
			format = FormatSQLite
			// Ensure .db extension
			if !strings.HasSuffix(filename, ".db") {
				fullPath += ".db"
			}
		}

		// Perform export
		err := fd.ui.ExportEvents(fullPath, format)
		if err != nil {
			// Could show error in status, but for now just hide dialog
			fd.Hide()
		} else {
			fd.Hide()
			// Update UI
			if fd.ui.gui != nil {
				fd.ui.gui.Update(func(g *gocui.Gui) error {
					if v, err := g.View(StatusView); err == nil {
						fd.ui.views.UpdateStatusView(v)
					}
					return fd.ui.layout.Layout(g)
				})
			}
		}
	case key == gocui.KeyEsc:
		// Cancel editing, return to file list
		fd.ui.state.FileDialog.IsEditing = false
		if fd.ui.gui != nil {
			fd.ui.gui.Update(func(g *gocui.Gui) error {
				_, err := g.SetCurrentView(FileListView)
				return err
			})
		}
	}
	fd.ui.state.FileDialog.Filename = fd.getFilenameFromView(v)
}

type FilenameEditor struct {
	fd *FileDialog
}

func (e FilenameEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) bool {
	e.fd.FilenameEditor(v, key, ch, mod)
	return false
}
