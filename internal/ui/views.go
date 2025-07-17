package ui

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/jesseduffield/gocui"
)

// Views handles all view update operations
type Views struct {
	ui *UI
}

// NewViews creates a new Views instance
func NewViews(ui *UI) *Views {
	return &Views{ui: ui}
}

// UpdateStatusView updates the status view
func (v *Views) UpdateStatusView(view *gocui.View) {
	view.Clear()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	// Check for recent export files
	var exportInfo string
	sqliteFiles, _ := filepath.Glob("watch-fs-events_*.db")
	jsonFiles, _ := filepath.Glob("watch-fs-events_*.json")

	if len(sqliteFiles) > 0 {
		exportInfo = fmt.Sprintf(" | Export: %s", green("SQLite available"))
	} else if len(jsonFiles) > 0 {
		exportInfo = fmt.Sprintf(" | Export: %s", green("JSON available"))
	}

	_, _ = fmt.Fprintf(view, "Watching: %s | Events: %s | Sort: %s%s\n",
		cyan(v.ui.rootPath),
		yellow(len(v.ui.state.Events)),
		cyan(v.ui.getSortOptionName()),
		exportInfo)
}

// UpdateFilterView updates the filter view
func (v *Views) UpdateFilterView(view *gocui.View) {
	view.Clear()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	dirsStatus := green("✓")
	if !v.ui.state.Filter.ShowDirs {
		dirsStatus = red("✗")
	}

	filesStatus := green("✓")
	if !v.ui.state.Filter.ShowFiles {
		filesStatus = red("✗")
	}

	aggregateStatus := green("✓")
	if !v.ui.state.AggregateEvents {
		aggregateStatus = red("✗")
	}

	_, _ = fmt.Fprintf(view, "Dirs: %s | Files: %s | Aggregate: %s | Path Filter: %s",
		dirsStatus, filesStatus, aggregateStatus, v.ui.state.Filter.PathFilter)
}

// UpdateEventsView updates the events view
func (v *Views) UpdateEventsView(view *gocui.View) {
	view.Clear()

	filteredEvents := v.ui.getFilteredEvents()
	lines := len(filteredEvents)

	if lines == 0 {
		_, _ = fmt.Fprintf(view, "No events to display")
		// Remettre le curseur en haut si la liste est vide
		view.SetCursor(0, 0)
		return
	}

	for _, event := range filteredEvents {
		v.renderEvent(view, event)
	}

	// Toujours borner le curseur dans la zone valide
	_, cy := view.Cursor()
	if cy < 0 {
		view.SetCursor(0, 0)
	} else if cy >= lines {
		view.SetCursor(0, lines-1)
	}
}

// UpdateHelpView updates the help view
func (v *Views) UpdateHelpView(view *gocui.View) {
	view.Clear()
	var helpText string

	switch v.ui.state.CurrentFocus {
	case FocusMain:
		helpText = "q: Quit | f: Toggle files | d: Toggle dirs | a: Toggle aggregate | s: Sort | ↑↓←→/hjkl: Navigate | PgUp/PgDn: Page | Home/End/g/G: Top/Bottom | Enter: Details | Ctrl+E: Export | Ctrl+I: Import"

	case FocusDetails:
		helpText = "ESC/q: Close details | Enter: Close details"

	case FocusFileDialog:
		if v.ui.state.FileDialog.Mode == ModeSave {
			helpText = "↑↓/kj: Navigate | Enter: Select file | e: Edit filename | ESC/q: Cancel | Save mode"
		} else {
			helpText = "↑↓/kj: Navigate | Enter: Select file | ESC/q: Cancel | Open mode"
		}

	default:
		helpText = "q: Quit | Navigation: ↑↓←→/hjkl | Enter: Details | Ctrl+E: Export | Ctrl+I: Import"
	}

	_, _ = fmt.Fprint(view, helpText)
}

// UpdateDetailsView updates the details popup view
func (v *Views) UpdateDetailsView(view *gocui.View) {
	view.Clear()
	if v.ui.state.SelectedEvent == nil {
		return
	}

	event := v.ui.state.SelectedEvent
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	// Format operation with color
	var operationStr string
	if event.Operation.Has(fsnotify.Create) {
		operationStr = green("CREATE")
	} else if event.Operation.Has(fsnotify.Write) {
		operationStr = yellow("WRITE")
	} else if event.Operation.Has(fsnotify.Remove) {
		operationStr = red("REMOVE")
	} else if event.Operation.Has(fsnotify.Rename) {
		operationStr = magenta("RENAME")
	} else if event.Operation.Has(fsnotify.Chmod) {
		operationStr = blue("CHMOD")
	} else {
		operationStr = "UNKNOWN"
	}

	// Get file info if possible
	fileInfo, err := v.ui.getFileInfo(event.Path)
	fileSize := "N/A"
	fileMode := "N/A"
	if err == nil && fileInfo != nil {
		fileSize = fmt.Sprintf("%d bytes", fileInfo.Size())
		fileMode = fileInfo.Mode().String()
	}

	// Display detailed information
	_, _ = fmt.Fprintf(view, "%sEvent Details%s\n", cyan("="), cyan("="))
	_, _ = fmt.Fprintf(view, "\n")
	_, _ = fmt.Fprintf(view, "%s: %s\n", cyan("Operation"), operationStr)
	_, _ = fmt.Fprintf(view, "%s: %s\n", cyan("Path"), event.Path)
	_, _ = fmt.Fprintf(view, "%s: %s\n", cyan("Type"), yellow(func() string {
		if event.IsDir {
			return "Directory"
		}
		return "File"
	}()))
	_, _ = fmt.Fprintf(view, "%s: %s\n", cyan("Timestamp"), event.Timestamp.Format("2006-01-02 15:04:05.000"))
	_, _ = fmt.Fprintf(view, "%s: %d\n", cyan("Count"), event.Count)

	if err == nil && fileInfo != nil {
		_, _ = fmt.Fprintf(view, "%s: %s\n", cyan("Size"), fileSize)
		_, _ = fmt.Fprintf(view, "%s: %s\n", cyan("Permissions"), fileMode)
		_, _ = fmt.Fprintf(view, "%s: %s\n", cyan("Modified"), fileInfo.ModTime().Format("2006-01-02 15:04:05"))
	}

	_, _ = fmt.Fprintf(view, "\n")
	_, _ = fmt.Fprintf(view, "%sPress ESC or q to close%s", yellow(""), yellow(""))
}

// renderEvent renders a single event with colors
func (v *Views) renderEvent(view *gocui.View, event *FileEvent) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	// Format timestamp
	timestamp := event.Timestamp.Format("15:04:05")

	// Format operation with color
	var operationStr string
	// Handle combined operations by checking each bit
	if event.Operation.Has(fsnotify.Create) {
		operationStr = green("CREATE")
	} else if event.Operation.Has(fsnotify.Write) {
		operationStr = yellow("WRITE")
	} else if event.Operation.Has(fsnotify.Remove) {
		operationStr = red("REMOVE")
	} else if event.Operation.Has(fsnotify.Rename) {
		operationStr = magenta("RENAME")
	} else if event.Operation.Has(fsnotify.Chmod) {
		operationStr = blue("CHMOD")
	} else {
		operationStr = "UNKNOWN"
	}

	// Format type indicator
	typeIndicator := "F"
	if event.IsDir {
		typeIndicator = "D"
	}

	// Format count if > 1
	countStr := ""
	if event.Count > 1 {
		countStr = fmt.Sprintf(" (%d)", event.Count)
	}

	// Format path
	pathStr := event.Path
	if len(pathStr) > 50 {
		pathStr = "..." + pathStr[len(pathStr)-47:]
	}

	// Render the event line
	line := fmt.Sprintf("[%s] %s %s %s%s", timestamp, operationStr, typeIndicator, pathStr, countStr)
	_, _ = fmt.Fprintln(view, line)
}
