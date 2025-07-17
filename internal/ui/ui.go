package ui

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/jesseduffield/gocui"
	_ "github.com/mattn/go-sqlite3"
)

// UI represents the terminal user interface
type UI struct {
	gui          *gocui.Gui
	state        *UIState
	fileDialog   *FileDialog
	exportImport *ExportImport
	navigation   *Navigation
	views        *Views
	keybindings  *Keybindings
	layout       *Layout
	events       *Events

	watcher interface {
		Events() <-chan fsnotify.Event
		Errors() <-chan error
		AddDirectory(path string) error
	}
	rootPath string
}

// NewUI creates a new UI instance
func NewUI(watcher interface {
	Events() <-chan fsnotify.Event
	Errors() <-chan error
	AddDirectory(path string) error
}, rootPath string) *UI {
	ui := &UI{
		state: &UIState{
			Events:          make([]*FileEvent, 0),
			Filter:          Filter{ShowDirs: true, ShowFiles: true},
			SortOption:      SortByTime,
			MaxEvents:       1000,
			AggregateEvents: true,      // Enable aggregation by default
			ShowDetails:     false,     // Details popup hidden by default
			SelectedEvent:   nil,       // No event selected by default
			ExportFilename:  "",        // No export filename by default
			ImportFilename:  "",        // No import filename by default
			ShowFileDialog:  false,     // File dialog hidden by default
			CurrentFocus:    FocusMain, // Start with main focus
			FileDialog: FileDialogState{
				CurrentPath: ".",
				Files:       make([]*FileEntry, 0),
				SelectedIdx: 0,
				Mode:        ModeSave,
				Filter:      "*.db",
			},
		},
		watcher:  watcher,
		rootPath: rootPath,
	}
	ui.fileDialog = NewFileDialog(ui)
	ui.exportImport = NewExportImport(ui)
	ui.navigation = NewNavigation(ui)
	ui.views = NewViews(ui)
	ui.keybindings = NewKeybindings(ui)
	ui.layout = NewLayout(ui)
	ui.events = NewEvents(ui)

	return ui
}

// Run starts the TUI
func (ui *UI) Run() error {
	gui, err := gocui.NewGui(gocui.NewGuiOpts{OutputMode: gocui.OutputTrue})
	if err != nil {
		return fmt.Errorf("failed to create GUI: %w", err)
	}
	defer gui.Close()

	ui.gui = gui
	gui.SetManagerFunc(ui.layout.Layout)

	// Set up keybindings
	if err := ui.keybindings.Setup(gui); err != nil {
		return fmt.Errorf("failed to set keybindings: %w", err)
	}

	// Start event watcher goroutine
	go ui.events.watchEvents()

	// Start the main loop
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("main loop error: %w", err)
	}

	return nil
}

// getFilteredEvents returns filtered and sorted events
func (ui *UI) getFilteredEvents() []*FileEvent {
	return ui.events.getFilteredEvents()
}

// getSortOptionName returns the name of current sort option
func (ui *UI) getSortOptionName() string {
	return ui.events.getSortOptionName()
}

// addEvent adds a new event to the state
func (ui *UI) addEvent(path string, operation fsnotify.Op, isDir bool) {
	ui.events.addEvent(path, operation, isDir)
}

// moveUp moves the selection up

// getFileInfo gets file info for a path
func (ui *UI) getFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Public methods for testing

// GetState returns the current UI state
func (ui *UI) GetState() *UIState {
	return ui.state
}

// GetRootPath returns the root path being watched
func (ui *UI) GetRootPath() string {
	return ui.rootPath
}

// AddEvent adds an event (public version for testing)
func (ui *UI) AddEvent(path string, operation fsnotify.Op, isDir bool) {
	ui.addEvent(path, operation, isDir)
}

// ToggleAggregate toggles aggregation (public version for testing)
func (ui *UI) ToggleAggregate() {
	ui.navigation.ToggleAggregate()
}

// ToggleFiles toggles file visibility (public version for testing)
func (ui *UI) ToggleFiles() {
	ui.navigation.ToggleFiles()
}

// ToggleDirs toggles directory visibility (public version for testing)
func (ui *UI) ToggleDirs() {
	ui.navigation.ToggleDirs()
}

// CycleSort cycles through sort options (public version for testing)
func (ui *UI) CycleSort() {
	ui.navigation.CycleSort()
}

// GetFilteredEvents returns filtered events (public version for testing)
func (ui *UI) GetFilteredEvents() []*FileEvent {
	return ui.getFilteredEvents()
}

// MoveUp moves the selection up (public version for testing)
func (ui *UI) MoveUp() {
	ui.navigation.MoveUp()
}

// MoveDown moves the selection down (public version for testing)
func (ui *UI) MoveDown() {
	ui.navigation.MoveDown()
}

// MoveLeft moves the selection left (public version for testing)
func (ui *UI) MoveLeft() {
	ui.navigation.MoveLeft()
}

// MoveRight moves the selection right (public version for testing)
func (ui *UI) MoveRight() {
	ui.navigation.MoveRight()
}

// PageUp moves the selection up by a page (public version for testing)
func (ui *UI) PageUp() {
	ui.navigation.PageUp()
}

// PageDown moves the selection down by a page (public version for testing)
func (ui *UI) PageDown() {
	ui.navigation.PageDown()
}

// MoveToTop moves the selection to the top (public version for testing)
func (ui *UI) MoveToTop() {
	ui.navigation.MoveToTop()
}

// MoveToBottom moves the selection to the bottom (public version for testing)
func (ui *UI) MoveToBottom() {
	ui.navigation.MoveToBottom()
}

// ExportEvents exports events to a file
func (ui *UI) ExportEvents(filename string, format ExportFormat) error {
	return ui.exportImport.ExportEvents(filename, format)
}

// ImportEvents imports events from a file
func (ui *UI) ImportEvents(filename string, format ExportFormat) error {
	return ui.exportImport.ImportEvents(filename, format)
}

// showFileDialog affiche le dialogue de fichiers et donne le focus
func (ui *UI) showFileDialog(mode FileDialogMode, filter string) {
	ui.fileDialog.Show(mode, filter)
}

// updateFileDialogView updates the main file dialog view
func (ui *UI) updateFileDialogView(v *gocui.View) {
	ui.fileDialog.UpdateMainView(v)
}

// updatePathView updates the path display view
func (ui *UI) updatePathView(v *gocui.View) {
	ui.fileDialog.UpdatePathView(v)
}

// updateFileListView updates the file list view
func (ui *UI) updateFileListView(v *gocui.View) {
	ui.fileDialog.UpdateFileListView(v)
}

// updateFilenameView updates the filename input view
func (ui *UI) updateFilenameView(v *gocui.View) {
	ui.fileDialog.UpdateFilenameView(v)
}

// formatFileSize formats file size in human readable format
func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.1f GB", float64(size)/(1024*1024*1024))
	}
}
