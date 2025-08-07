package ui

import (
	"github.com/jesseduffield/gocui"
)

// Keybindings manages all keyboard shortcuts and input handling
type Keybindings struct {
	ui *UI
}

// NewKeybindings creates a new Keybindings instance
func NewKeybindings(ui *UI) *Keybindings {
	return &Keybindings{
		ui: ui,
	}
}

// Setup configures all keybindings for the application
func (kb *Keybindings) Setup(g *gocui.Gui) error {
	// Global keybindings (always active)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, kb.quit); err != nil {
		return err
	}

	// Main interface keybindings (EventsView focused)
	if err := g.SetKeybinding(EventsView, gocui.KeyArrowUp, gocui.ModNone, kb.navigationMoveUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyArrowDown, gocui.ModNone, kb.navigationMoveDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyArrowLeft, gocui.ModNone, kb.navigationMoveLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyArrowRight, gocui.ModNone, kb.navigationMoveRight); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'k', gocui.ModNone, kb.navigationMoveUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'j', gocui.ModNone, kb.navigationMoveDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'h', gocui.ModNone, kb.navigationMoveLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'l', gocui.ModNone, kb.navigationMoveRight); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyPgup, gocui.ModNone, kb.navigationPageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyPgdn, gocui.ModNone, kb.navigationPageDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyHome, gocui.ModNone, kb.navigationMoveToTop); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyEnd, gocui.ModNone, kb.navigationMoveToBottom); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'g', gocui.ModNone, kb.navigationMoveToTop); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'G', gocui.ModNone, kb.navigationMoveToBottom); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyEnter, gocui.ModNone, kb.showEventDetails); err != nil {
		return err
	}

	if err := g.SetKeybinding(EventsView, 'q', gocui.ModNone, kb.quit); err != nil {
		return err
	}

	// Main interface actions (EventsView focused)
	if err := g.SetKeybinding(EventsView, 'f', gocui.ModNone, kb.toggleFiles); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'd', gocui.ModNone, kb.toggleDirs); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'a', gocui.ModNone, kb.toggleAggregate); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 's', gocui.ModNone, kb.cycleSort); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyCtrlE, gocui.ModNone, kb.exportEventsHandler); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyCtrlI, gocui.ModNone, kb.importEventsHandler); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyCtrlF, gocui.ModNone, kb.showFolderManager); err != nil {
		return err
	}

	// Global escape key for closing popups
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, kb.debugEscape); err != nil {
		return err
	}

	// Details popup keybindings
	if err := g.SetKeybinding(DetailsView, gocui.KeyEsc, gocui.ModNone, kb.hideEventDetails); err != nil {
		return err
	}
	if err := g.SetKeybinding(DetailsView, 'q', gocui.ModNone, kb.hideEventDetails); err != nil {
		return err
	}
	if err := g.SetKeybinding(DetailsView, gocui.KeyEnter, gocui.ModNone, kb.hideEventDetails); err != nil {
		return err
	}

	// File dialog keybindings
	if err := g.SetKeybinding(FileListView, gocui.KeyArrowUp, gocui.ModNone, kb.fileDialogUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(FileListView, gocui.KeyArrowDown, gocui.ModNone, kb.fileDialogDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(FileListView, 'k', gocui.ModNone, kb.fileDialogUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(FileListView, 'j', gocui.ModNone, kb.fileDialogDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(FileListView, gocui.KeyEnter, gocui.ModNone, kb.fileDialogEnter); err != nil {
		return err
	}
	if err := g.SetKeybinding(FileListView, gocui.KeyEsc, gocui.ModNone, kb.fileDialogCancel); err != nil {
		return err
	}
	if err := g.SetKeybinding(FileListView, 'q', gocui.ModNone, kb.fileDialogCancel); err != nil {
		return err
	}
	if err := g.SetKeybinding(FileListView, 'e', gocui.ModNone, kb.fileDialogEditFilename); err != nil {
		return err
	}

	// Folder manager keybindings
	if err := g.SetKeybinding(FolderListView, gocui.KeyArrowUp, gocui.ModNone, kb.folderManagerUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, gocui.KeyArrowDown, gocui.ModNone, kb.folderManagerDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, 'k', gocui.ModNone, kb.folderManagerUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, 'j', gocui.ModNone, kb.folderManagerDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, gocui.KeyEnter, gocui.ModNone, kb.folderManagerEnter); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, gocui.KeyEsc, gocui.ModNone, kb.folderManagerCancel); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, 'q', gocui.ModNone, kb.folderManagerCancel); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, 'a', gocui.ModNone, kb.folderManagerAdd); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, 'd', gocui.ModNone, kb.folderManagerRemove); err != nil {
		return err
	}

	// Keybindings for the "Currently Watching" list (watched_folders view)
	if err := g.SetKeybinding("watched_folders", gocui.KeyArrowUp, gocui.ModNone, kb.watchedFoldersUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", gocui.KeyArrowDown, gocui.ModNone, kb.watchedFoldersDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", 'k', gocui.ModNone, kb.watchedFoldersUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", 'j', gocui.ModNone, kb.watchedFoldersDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", 'r', gocui.ModNone, kb.watchedFoldersRemove); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", gocui.KeyEsc, gocui.ModNone, kb.folderManagerCancel); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", 'q', gocui.ModNone, kb.folderManagerCancel); err != nil {
		return err
	}

	// Panel switching keybindings for watched_folders panel
	if err := g.SetKeybinding("watched_folders", gocui.KeyTab, gocui.ModNone, kb.switchToNextPanel); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", gocui.KeyBacktab, gocui.ModNone, kb.switchToPreviousPanel); err != nil {
		return err
	}
	if err := g.SetKeybinding("watched_folders", gocui.KeyArrowRight, gocui.ModNone, kb.switchToRightPanel); err != nil {
		return err
	}

	// Panel switching keybindings for folder list panel
	if err := g.SetKeybinding(FolderListView, gocui.KeyTab, gocui.ModNone, kb.switchToNextPanel); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, gocui.KeyBacktab, gocui.ModNone, kb.switchToPreviousPanel); err != nil {
		return err
	}
	if err := g.SetKeybinding(FolderListView, gocui.KeyArrowLeft, gocui.ModNone, kb.switchToLeftPanel); err != nil {
		return err
	}

	return nil
}

// Main interface action functions (only called when EventsView has focus)
func (kb *Keybindings) toggleFiles(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.toggleFiles(g, v)
}

func (kb *Keybindings) toggleDirs(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.toggleDirs(g, v)
}

func (kb *Keybindings) toggleAggregate(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.toggleAggregate(g, v)
}

func (kb *Keybindings) cycleSort(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.cycleSort(g, v)
}

// Navigation wrapper functions
func (kb *Keybindings) navigationMoveUp(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.moveUp(g, v)
}

func (kb *Keybindings) navigationMoveDown(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.moveDown(g, v)
}

func (kb *Keybindings) navigationMoveLeft(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.moveLeft(g, v)
}

func (kb *Keybindings) navigationMoveRight(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.moveRight(g, v)
}

func (kb *Keybindings) navigationPageUp(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.pageUp(g, v)
}

func (kb *Keybindings) navigationPageDown(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.pageDown(g, v)
}

func (kb *Keybindings) navigationMoveToTop(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.moveToTop(g, v)
}

func (kb *Keybindings) navigationMoveToBottom(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.navigation.moveToBottom(g, v)
}

// Event details functions
func (kb *Keybindings) showEventDetails(g *gocui.Gui, v *gocui.View) error {
	if len(kb.ui.state.Events) == 0 {
		return nil
	}

	// Get current cursor position
	_, cy := v.Cursor()
	filteredEvents := kb.ui.getFilteredEvents()
	if cy >= 0 && cy < len(filteredEvents) {
		kb.ui.state.SelectedEvent = filteredEvents[cy]
		kb.ui.state.ShowDetails = true
		kb.ui.state.CurrentFocus = FocusDetails
		g.Update(func(g *gocui.Gui) error {
			var err error
			_, err = g.SetCurrentView(DetailsView)
			return err
		})
	}
	return nil
}

func (kb *Keybindings) hideEventDetails(g *gocui.Gui, v *gocui.View) error {
	kb.ui.state.ShowDetails = false
	kb.ui.state.SelectedEvent = nil
	kb.ui.state.CurrentFocus = FocusMain
	_, err := kb.ui.gui.SetCurrentView(EventsView)
	return err
}

// Export/Import handlers
func (kb *Keybindings) exportEventsHandler(g *gocui.Gui, v *gocui.View) error {
	// Ouvre le dialogue d'export (mode Save, filtre .db)
	kb.ui.showFileDialog(ModeSave, "*.db")
	return nil
}

func (kb *Keybindings) importEventsHandler(g *gocui.Gui, v *gocui.View) error {
	// Ouvre le dialogue d'import (mode Open, filtre .db)
	kb.ui.showFileDialog(ModeOpen, "*.db")
	return nil
}

// Folder manager keybinding
func (kb *Keybindings) showFolderManager(g *gocui.Gui, v *gocui.View) error {
	kb.ui.ShowFolderManager()
	return nil
}

// File dialog functions
func (kb *Keybindings) fileDialogUp(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.fileDialog.Up(g, v)
}

func (kb *Keybindings) fileDialogDown(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.fileDialog.Down(g, v)
}

func (kb *Keybindings) fileDialogEnter(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.fileDialog.Enter(g, v)
}

func (kb *Keybindings) fileDialogEditFilename(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.fileDialog.EditFilename(g, v)
}

func (kb *Keybindings) fileDialogCancel(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.fileDialog.Cancel(g, v)
}

// Folder manager keybindings
func (kb *Keybindings) folderManagerUp(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.Up(g, v)
}

func (kb *Keybindings) folderManagerDown(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.Down(g, v)
}

func (kb *Keybindings) folderManagerEnter(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.Enter(g, v)
}

func (kb *Keybindings) folderManagerCancel(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.Cancel(g, v)
}

func (kb *Keybindings) folderManagerAdd(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.Add(g, v)
}

func (kb *Keybindings) folderManagerRemove(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.Remove(g, v)
}

// Keybindings for "Currently Watching" navigation
func (kb *Keybindings) watchedFoldersUp(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.NavigateWatchedUp(g, v)
}

func (kb *Keybindings) watchedFoldersDown(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.NavigateWatchedDown(g, v)
}

func (kb *Keybindings) watchedFoldersRemove(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.RemoveWatchedFolder(g, v)
}

// Panel switching functions for folder manager
func (kb *Keybindings) switchToNextPanel(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.SwitchToNextPanel(g, v)
}

func (kb *Keybindings) switchToPreviousPanel(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.SwitchToPreviousPanel(g, v)
}

func (kb *Keybindings) switchToRightPanel(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.SwitchToRightPanel(g, v)
}

func (kb *Keybindings) switchToLeftPanel(g *gocui.Gui, v *gocui.View) error {
	return kb.ui.folderManager.SwitchToLeftPanel(g, v)
}

// Utility functions
func (kb *Keybindings) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (kb *Keybindings) debugEscape(g *gocui.Gui, v *gocui.View) error {
	// Handle escape key based on current focus
	if kb.ui.state.ShowDetails {
		return kb.hideEventDetails(g, v)
	} else if kb.ui.state.ShowFileDialog {
		return kb.fileDialogCancel(g, v)
	} else if kb.ui.state.ShowFolderManager {
		return kb.folderManagerCancel(g, v)
	}
	return nil
}
