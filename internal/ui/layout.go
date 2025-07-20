package ui

import (
	"github.com/go-errors/errors"
	"github.com/pbouamriou/watch-fs/pkg/logger"

	"github.com/gdamore/tcell/v2"
	"github.com/jesseduffield/gocui"
)

// isUnknownViewError checks if the error is ErrUnknownView (handles wrapped errors)
func isUnknownViewError(err error) bool {
	return errors.Is(err, gocui.ErrUnknownView)
}

// Layout manages the layout and positioning of all UI views
type Layout struct {
	ui *UI
}

// NewLayout creates a new Layout instance
func NewLayout(ui *UI) *Layout {
	return &Layout{
		ui: ui,
	}
}

// Layout defines the layout of the TUI
func (l *Layout) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Layout main views
	if err := l.layoutMainViews(g, maxX, maxY); err != nil {
		return err
	}

	// Layout overlay views (popups and dialogs)
	if err := l.layoutOverlayViews(g, maxX, maxY); err != nil {
		return err
	}

	// Set focus based on current state
	return l.setFocus(g)
}

// layoutMainViews creates the main interface views
func (l *Layout) layoutMainViews(g *gocui.Gui, maxX, maxY int) error {
	// Status view (top)
	if v, err := g.SetView(StatusView, 0, 0, maxX-1, 2, 0); err != nil {
		if !isUnknownViewError(err) {
			return err
		}
		v.Title = " Status "
		v.Frame = true
		l.ui.views.UpdateStatusView(v)
	} else {
		l.ui.views.UpdateStatusView(v)
	}

	// Filter view (below status)
	if v, err := g.SetView(FilterView, 0, 3, maxX-1, 5, 0); err != nil {
		if !isUnknownViewError(err) {
			return err
		}
		v.Title = " Filter "
		v.Frame = true
		l.ui.views.UpdateFilterView(v)
	} else {
		l.ui.views.UpdateFilterView(v)
	}

	// Events view (main area)
	if v, err := g.SetView(EventsView, 0, 6, maxX-1, maxY-4, 0); err != nil {
		if !isUnknownViewError(err) {
			return err
		}
		v.Title = " Events "
		v.Frame = true
		v.Highlight = true
		v.SelBgColor = gocui.Attribute(tcell.ColorDarkGreen)
		v.SelFgColor = gocui.ColorBlack
		l.ui.views.UpdateEventsView(v)
	} else {
		l.ui.views.UpdateEventsView(v)
	}

	// Help view (bottom)
	if v, err := g.SetView(HelpView, 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if !isUnknownViewError(err) {
			return err
		}
		v.Title = " Help "
		v.Frame = true
		l.ui.views.UpdateHelpView(v)
	} else {
		l.ui.views.UpdateHelpView(v)
	}

	return nil
}

// layoutOverlayViews creates overlay views (popups and dialogs)
func (l *Layout) layoutOverlayViews(g *gocui.Gui, maxX, maxY int) error {
	// Layout details popup
	if err := l.layoutDetailsPopup(g, maxX, maxY); err != nil {
		return err
	}

	// Layout file dialog
	if err := l.layoutFileDialog(g, maxX, maxY); err != nil {
		return err
	}

	// Layout folder manager
	if err := l.layoutFolderManager(g, maxX, maxY); err != nil {
		return err
	}

	return nil
}

// layoutDetailsPopup creates the details popup overlay
func (l *Layout) layoutDetailsPopup(g *gocui.Gui, maxX, maxY int) error {
	if l.ui.state.ShowDetails && l.ui.state.SelectedEvent != nil {
		// Calculate popup size and position (centered)
		popupWidth := 60
		popupHeight := 12
		x0 := (maxX - popupWidth) / 2
		y0 := (maxY - popupHeight) / 2
		x1 := x0 + popupWidth
		y1 := y0 + popupHeight

		if v, err := g.SetView(DetailsView, x0, y0, x1, y1, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			v.Title = " Event Details "
			v.Frame = true
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			l.ui.views.UpdateDetailsView(v)
		}
	} else {
		// Remove details view if not needed
		if err := g.DeleteView(DetailsView); err != nil {
			logger.Error(err, "DeleteView error")
		}
	}
	return nil
}

// layoutFileDialog creates the file dialog overlay
func (l *Layout) layoutFileDialog(g *gocui.Gui, maxX, maxY int) error {
	if l.ui.state.ShowFileDialog {
		// Calculate popup size and position (centered, larger for file dialog)
		popupWidth := 80
		popupHeight := 20
		x0 := (maxX - popupWidth) / 2
		y0 := (maxY - popupHeight) / 2
		x1 := x0 + popupWidth
		y1 := y0 + popupHeight

		// Main dialog frame
		if v, err := g.SetView(FileDialogView, x0, y0, x1, y1, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			title := " Save File "
			if l.ui.state.FileDialog.Mode == ModeOpen {
				title = " Open File "
			}
			v.Title = title
			v.Frame = true
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			l.ui.updateFileDialogView(v)
		} else {
			l.ui.updateFileDialogView(v)
		}

		// Path display (top of dialog)
		pathHeight := 3
		if v, err := g.SetView(PathView, x0+1, y0+1, x1-1, y0+pathHeight, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			v.Frame = false
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorCyan
			l.ui.updatePathView(v)
		} else {
			l.ui.updatePathView(v)
		}

		// File list (main area) - smaller in save mode to make room for filename input
		listY0 := y0 + pathHeight + 1
		listY1 := y1 - 5 // Leave more space for filename input in save mode
		if l.ui.state.FileDialog.Mode == ModeOpen {
			listY1 = y1 - 3 // Full height for open mode
		}

		if v, err := g.SetView(FileListView, x0+1, listY0, x1-1, listY1, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			v.Frame = false
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			v.Highlight = true
			v.SelBgColor = gocui.Attribute(tcell.ColorDarkGreen)
			v.SelFgColor = gocui.ColorBlack
		}
		// Always update the file list view, whether it's new or existing
		if v, err := g.View(FileListView); err == nil {
			l.ui.updateFileListView(v)
		}

		// Filename input (only in save mode)
		if l.ui.state.FileDialog.Mode == ModeSave {
			filenameY0 := listY1 + 1
			filenameY1 := y1 - 2
			if v, err := g.SetView(FilenameView, x0+1, filenameY0, x1-1, filenameY1, 0); err != nil {
				if !isUnknownViewError(err) {
					return err
				}
				v.Frame = true
				v.BgColor = gocui.ColorBlack
				v.FgColor = gocui.ColorYellow
				v.Editable = true
				v.Editor = FilenameEditor{fd: l.ui.fileDialog}
			}
			// Always update the filename view
			if v, err := g.View(FilenameView); err == nil {
				g.Cursor = true
				l.ui.updateFilenameView(v)
			}
		}
	} else {
		// Remove file dialog views if not needed
		if err := g.DeleteView(FileDialogView); err != nil {
			logger.Error(err, "DeleteView error")
		}
		if err := g.DeleteView(FilenameView); err != nil {
			logger.Error(err, "DeleteView error")
		}
		if err := g.DeleteView(PathView); err != nil {
			logger.Error(err, "DeleteView error")
		}
		if err := g.DeleteView(FileListView); err != nil {
			logger.Error(err, "DeleteView error")
		}
	}
	return nil
}

// layoutFolderManager creates the folder manager overlay
func (l *Layout) layoutFolderManager(g *gocui.Gui, maxX, maxY int) error {
	if l.ui.state.ShowFolderManager {
		// Calculate popup size and position (centered, large for folder manager)
		popupWidth := 100
		popupHeight := 25
		x0 := (maxX - popupWidth) / 2
		y0 := (maxY - popupHeight) / 2
		x1 := x0 + popupWidth
		y1 := y0 + popupHeight

		// Main dialog frame
		if v, err := g.SetView(FolderManagerView, x0, y0, x1, y1, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			v.Title = " Folder Manager "
			v.Frame = true
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			l.ui.folderManager.UpdateMainView(v)
		} else {
			l.ui.folderManager.UpdateMainView(v)
		}

		// Split the dialog into two sections
		splitX := x0 + (x1-x0)/2

		// Left side: Currently watched folders
		leftY0 := y0 + 1
		leftY1 := y1 - 1
		if v, err := g.SetView("watched_folders", x0+1, leftY0, splitX-1, leftY1, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			v.Frame = true
			v.Title = " Currently Watching "
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			v.Highlight = true
			v.SelBgColor = gocui.Attribute(tcell.ColorDarkGreen)
			v.SelFgColor = gocui.ColorBlack
		}

		// Right side: Folder browser
		rightY0 := y0 + 1
		rightY1 := y1 - 1
		if v, err := g.SetView(FolderListView, splitX+1, rightY0, x1-1, rightY1, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			v.Frame = true
			v.Title = " Available Folders "
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			v.Highlight = true
			v.SelBgColor = gocui.Attribute(tcell.ColorDarkGreen)
			v.SelFgColor = gocui.ColorBlack
			l.ui.folderManager.UpdateFolderListView(v)
		} else {
			l.ui.folderManager.UpdateFolderListView(v)
		}

		// Path display (top of right side)
		pathHeight := 3
		if v, err := g.SetView("folder_path", splitX+1, rightY0, x1-1, rightY0+pathHeight, 0); err != nil {
			if !isUnknownViewError(err) {
				return err
			}
			v.Frame = false
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorCyan
			l.ui.folderManager.UpdatePathView(v)
		} else {
			l.ui.folderManager.UpdatePathView(v)
		}

	} else {
		// Remove folder manager views if not needed
		if err := g.DeleteView(FolderManagerView); err != nil {
			logger.Error(err, "DeleteView error")
		}
		if err := g.DeleteView(FolderListView); err != nil {
			logger.Error(err, "DeleteView error")
		}
		if err := g.DeleteView("watched_folders"); err != nil {
			logger.Error(err, "DeleteView error")
		}
		if err := g.DeleteView("folder_path"); err != nil {
			logger.Error(err, "DeleteView error")
		}
	}
	return nil
}

// setFocus sets the current focus based on the UI state
func (l *Layout) setFocus(g *gocui.Gui) error {
	// Set EventsView as the default active view (unless dialogs are shown)
	if !l.ui.state.ShowDetails && !l.ui.state.ShowFileDialog && !l.ui.state.ShowFolderManager {
		if _, err := g.SetCurrentView(EventsView); err != nil {
			return err
		}
	} else if l.ui.state.ShowDetails {
		if _, err := g.SetCurrentView(DetailsView); err != nil {
			return err
		}
	} else if l.ui.state.ShowFileDialog {
		if l.ui.state.FileDialog.Mode == ModeSave && l.ui.state.FileDialog.IsEditing {
			if _, err := g.SetCurrentView(FilenameView); err != nil {
				return err
			}
		} else {
			if _, err := g.SetCurrentView(FileListView); err != nil {
				return err
			}
		}
	} else if l.ui.state.ShowFolderManager {
		if _, err := g.SetCurrentView(FolderListView); err != nil {
			return err
		}
	}

	return nil
}
