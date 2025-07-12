package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/jroimartin/gocui"
)

// UI represents the terminal user interface
type UI struct {
	gui     *gocui.Gui
	state   *UIState
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
	return &UI{
		state: &UIState{
			Events:          make([]*FileEvent, 0),
			Filter:          Filter{ShowDirs: true, ShowFiles: true},
			SortOption:      SortByTime,
			MaxEvents:       1000,
			AggregateEvents: true,  // Enable aggregation by default
			ShowDetails:     false, // Details popup hidden by default
			SelectedEvent:   nil,   // No event selected by default
		},
		watcher:  watcher,
		rootPath: rootPath,
	}
}

// Run starts the TUI
func (ui *UI) Run() error {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return fmt.Errorf("failed to create GUI: %w", err)
	}
	defer gui.Close()

	ui.gui = gui
	gui.SetManagerFunc(ui.layout)

	// Set up keybindings
	if err := ui.keybindings(gui); err != nil {
		return fmt.Errorf("failed to set keybindings: %w", err)
	}

	// Start event watcher goroutine
	go ui.watchEvents()

	// Start the main loop
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("main loop error: %w", err)
	}

	return nil
}

// layout defines the layout of the TUI
func (ui *UI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Status view (top)
	if v, err := g.SetView(StatusView, 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Status "
		v.Frame = true
		ui.updateStatusView(v)
	}

	// Filter view (below status)
	if v, err := g.SetView(FilterView, 0, 3, maxX-1, 5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Filter "
		v.Frame = true
		ui.updateFilterView(v)
	}

	// Events view (main area)
	if v, err := g.SetView(EventsView, 0, 6, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Events "
		v.Frame = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		ui.updateEventsView(v)
	}

	// Help view (bottom)
	if v, err := g.SetView(HelpView, 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Help "
		v.Frame = true
		ui.updateHelpView(v)
	}

	// Details popup (overlay)
	if ui.state.ShowDetails && ui.state.SelectedEvent != nil {
		// Calculate popup size and position (centered)
		popupWidth := 60
		popupHeight := 12
		x0 := (maxX - popupWidth) / 2
		y0 := (maxY - popupHeight) / 2
		x1 := x0 + popupWidth
		y1 := y0 + popupHeight

		if v, err := g.SetView(DetailsView, x0, y0, x1, y1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = " Event Details "
			v.Frame = true
			v.BgColor = gocui.ColorBlack
			v.FgColor = gocui.ColorWhite
			ui.updateDetailsView(v)
		}
	} else {
		// Remove details view if not needed
		g.DeleteView(DetailsView)
	}

	// Set EventsView as the default active view (unless details are shown)
	if !ui.state.ShowDetails {
		if _, err := g.SetCurrentView(EventsView); err != nil {
			return err
		}
	} else {
		if _, err := g.SetCurrentView(DetailsView); err != nil {
			return err
		}
	}

	return nil
}

// updateStatusView updates the status view
func (ui *UI) updateStatusView(v *gocui.View) {
	v.Clear()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	_, _ = fmt.Fprintf(v, "Watching: %s | Events: %s | Sort: %s\n",
		cyan(ui.rootPath),
		yellow(len(ui.state.Events)),
		cyan(ui.getSortOptionName()))
}

// updateFilterView updates the filter view
func (ui *UI) updateFilterView(v *gocui.View) {
	v.Clear()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	dirsStatus := green("✓")
	if !ui.state.Filter.ShowDirs {
		dirsStatus = red("✗")
	}

	filesStatus := green("✓")
	if !ui.state.Filter.ShowFiles {
		filesStatus = red("✗")
	}

	aggregateStatus := green("✓")
	if !ui.state.AggregateEvents {
		aggregateStatus = red("✗")
	}

	_, _ = fmt.Fprintf(v, "Dirs: %s | Files: %s | Aggregate: %s | Path Filter: %s",
		dirsStatus, filesStatus, aggregateStatus, ui.state.Filter.PathFilter)
}

// updateEventsView updates the events view
func (ui *UI) updateEventsView(v *gocui.View) {
	v.Clear()

	filteredEvents := ui.getFilteredEvents()
	lines := len(filteredEvents)

	if lines == 0 {
		_, _ = fmt.Fprintf(v, "No events to display")
		// Remettre le curseur en haut si la liste est vide
		_ = v.SetCursor(0, 0)
		return
	}

	for _, event := range filteredEvents {
		ui.renderEvent(v, event, false)
	}

	// Toujours borner le curseur dans la zone valide
	_, cy := v.Cursor()
	if cy < 0 {
		_ = v.SetCursor(0, 0)
	} else if cy >= lines {
		_ = v.SetCursor(0, lines-1)
	}
}

// updateHelpView updates the help view
func (ui *UI) updateHelpView(v *gocui.View) {
	v.Clear()
	_, _ = fmt.Fprintf(v, "q: Quit | f: Toggle files | d: Toggle dirs | a: Toggle aggregate | s: Sort | ↑↓←→/hjkl: Navigate | PgUp/PgDn/u/d: Page | Home/End/g/G: Top/Bottom | Enter: Details")
}

// updateDetailsView updates the details popup view
func (ui *UI) updateDetailsView(v *gocui.View) {
	v.Clear()
	if ui.state.SelectedEvent == nil {
		return
	}

	event := ui.state.SelectedEvent
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
	fileInfo, err := ui.getFileInfo(event.Path)
	fileSize := "N/A"
	fileMode := "N/A"
	if err == nil && fileInfo != nil {
		fileSize = fmt.Sprintf("%d bytes", fileInfo.Size())
		fileMode = fileInfo.Mode().String()
	}

	// Display detailed information
	_, _ = fmt.Fprintf(v, "%sEvent Details%s\n", cyan("="), cyan("="))
	_, _ = fmt.Fprintf(v, "\n")
	_, _ = fmt.Fprintf(v, "%s: %s\n", cyan("Operation"), operationStr)
	_, _ = fmt.Fprintf(v, "%s: %s\n", cyan("Path"), event.Path)
	_, _ = fmt.Fprintf(v, "%s: %s\n", cyan("Type"), yellow(func() string {
		if event.IsDir {
			return "Directory"
		}
		return "File"
	}()))
	_, _ = fmt.Fprintf(v, "%s: %s\n", cyan("Timestamp"), event.Timestamp.Format("2006-01-02 15:04:05.000"))
	_, _ = fmt.Fprintf(v, "%s: %d\n", cyan("Count"), event.Count)

	if err == nil && fileInfo != nil {
		_, _ = fmt.Fprintf(v, "%s: %s\n", cyan("Size"), fileSize)
		_, _ = fmt.Fprintf(v, "%s: %s\n", cyan("Permissions"), fileMode)
		_, _ = fmt.Fprintf(v, "%s: %s\n", cyan("Modified"), fileInfo.ModTime().Format("2006-01-02 15:04:05"))
	}

	_, _ = fmt.Fprintf(v, "\n")
	_, _ = fmt.Fprintf(v, "%sPress ESC or q to close%s", yellow(""), yellow(""))
}

// renderEvent renders a single event with colors
func (ui *UI) renderEvent(v *gocui.View, event *FileEvent, selected bool) {
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

	// Format path
	path := event.Path
	if event.IsDir {
		path = blue(path + "/")
	}

	// Format count
	countStr := ""
	if event.Count > 1 {
		countStr = fmt.Sprintf(" (%d)", event.Count)
	}

	line := fmt.Sprintf("%s %s %s%s", timestamp, operationStr, path, countStr)

	if selected {
		_, _ = fmt.Fprintf(v, "> %s\n", line)
	} else {
		_, _ = fmt.Fprintf(v, "  %s\n", line)
	}
}

// getFilteredEvents returns filtered and sorted events
func (ui *UI) getFilteredEvents() []*FileEvent {
	filtered := make([]*FileEvent, 0)

	for _, event := range ui.state.Events {
		// Apply path filter
		if ui.state.Filter.PathFilter != "" &&
			!strings.Contains(strings.ToLower(event.Path), strings.ToLower(ui.state.Filter.PathFilter)) {
			continue
		}

		// Apply operation filter
		if ui.state.Filter.OperationFilter != 0 && event.Operation != ui.state.Filter.OperationFilter {
			continue
		}

		// Apply type filter
		if event.IsDir && !ui.state.Filter.ShowDirs {
			continue
		}
		if !event.IsDir && !ui.state.Filter.ShowFiles {
			continue
		}

		filtered = append(filtered, event)
	}

	// Sort events
	ui.sortEvents(filtered)

	return filtered
}

// sortEvents sorts events based on current sort option
func (ui *UI) sortEvents(events []*FileEvent) {
	switch ui.state.SortOption {
	case SortByTime:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Timestamp.After(events[j].Timestamp)
		})
	case SortByPath:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Path < events[j].Path
		})
	case SortByOperation:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Operation < events[j].Operation
		})
	case SortByCount:
		sort.Slice(events, func(i, j int) bool {
			return events[i].Count > events[j].Count
		})
	}
}

// getSortOptionName returns the name of current sort option
func (ui *UI) getSortOptionName() string {
	switch ui.state.SortOption {
	case SortByTime:
		return "Time"
	case SortByPath:
		return "Path"
	case SortByOperation:
		return "Operation"
	case SortByCount:
		return "Count"
	default:
		return "Unknown"
	}
}

// addEvent adds a new event to the state
func (ui *UI) addEvent(path string, operation fsnotify.Op, isDir bool) {
	if ui.state.AggregateEvents {
		// Check if we already have an event for this path in the last second
		now := time.Now()
		for _, event := range ui.state.Events {
			if event.Path == path && event.Operation == operation &&
				now.Sub(event.Timestamp) < time.Second {
				event.Count++
				event.Timestamp = now
				return
			}
		}
	}

	// Add new event
	event := &FileEvent{
		Path:      path,
		Operation: operation,
		Timestamp: time.Now(),
		IsDir:     isDir,
		Count:     1,
	}

	ui.state.Events = append(ui.state.Events, event)

	// Limit the number of events
	if len(ui.state.Events) > ui.state.MaxEvents {
		ui.state.Events = ui.state.Events[1:]
	}

	// Update UI only if gui is initialized (not in tests)
	if ui.gui != nil {
		ui.gui.Update(func(g *gocui.Gui) error {
			if v, err := g.View(EventsView); err == nil {
				ui.updateEventsView(v)
			}
			if v, err := g.View(StatusView); err == nil {
				ui.updateStatusView(v)
			}
			if v, err := g.View(FilterView); err == nil {
				ui.updateFilterView(v)
			}
			return nil
		})
	}
}

// watchEvents watches for file system events
func (ui *UI) watchEvents() {
	for {
		select {
		case event, ok := <-ui.watcher.Events():
			if !ok {
				return
			}

			info, err := ui.getFileInfo(event.Name)
			isDir := false
			if err == nil {
				isDir = info.IsDir()
			}

			ui.addEvent(event.Name, event.Op, isDir)

		case err, ok := <-ui.watcher.Errors():
			if !ok {
				return
			}
			// Handle error silently for now
			_ = err
		}
	}
}

// debugKeyPress is a debug function to test if key events are being received
func (ui *UI) debugKeyPress(g *gocui.Gui, v *gocui.View) error {
	// This function will be called for any key press to help debug
	// We'll just return nil to continue normal processing
	return nil
}

// keybindings sets up the key bindings
func (ui *UI) keybindings(g *gocui.Gui) error {
	// Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, ui.quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, ui.moveUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, ui.moveDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, ui.moveLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, ui.moveRight); err != nil {
		return err
	}

	// Navigation - Global alternative keys (fallback)
	if err := g.SetKeybinding("", 'k', gocui.ModNone, ui.moveUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, ui.moveDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone, ui.moveLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone, ui.moveRight); err != nil {
		return err
	}

	// Navigation - Page Up/Down (standard)
	if err := g.SetKeybinding(EventsView, gocui.KeyPgup, gocui.ModNone, ui.pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyPgdn, gocui.ModNone, ui.pageDown); err != nil {
		return err
	}

	// Navigation - Page Up/Down alternatives for Mac
	// u/d for page up/down (simple keys)
	if err := g.SetKeybinding(EventsView, 'u', gocui.ModNone, ui.pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'd', gocui.ModNone, ui.pageDown); err != nil {
		return err
	}

	// Navigation - Global page keys (fallback)
	if err := g.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, ui.pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, ui.pageDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'u', gocui.ModNone, ui.pageUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'd', gocui.ModNone, ui.pageDown); err != nil {
		return err
	}

	// Navigation - Home/End (standard)
	if err := g.SetKeybinding(EventsView, gocui.KeyHome, gocui.ModNone, ui.moveToTop); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyEnd, gocui.ModNone, ui.moveToBottom); err != nil {
		return err
	}

	// Navigation - Home/End alternatives for Mac
	// g/G for top/bottom (vim-style)
	if err := g.SetKeybinding(EventsView, 'g', gocui.ModNone, ui.moveToTop); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, 'G', gocui.ModNone, ui.moveToBottom); err != nil {
		return err
	}

	// Navigation - Global home/end keys (fallback)
	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, ui.moveToTop); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, ui.moveToBottom); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'g', gocui.ModNone, ui.moveToTop); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'G', gocui.ModNone, ui.moveToBottom); err != nil {
		return err
	}

	// Filtering
	if err := g.SetKeybinding("", 'f', gocui.ModNone, ui.toggleFiles); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'd', gocui.ModNone, ui.toggleDirs); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'a', gocui.ModNone, ui.toggleAggregate); err != nil {
		return err
	}

	// Sorting
	if err := g.SetKeybinding("", 's', gocui.ModNone, ui.cycleSort); err != nil {
		return err
	}

	// Event details - Escape to hide details (global)
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, ui.hideEventDetails); err != nil {
		return err
	}

	// Event details - Enter to hide details when popup is open (global)
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, ui.toggleEventDetails); err != nil {
		return err
	}

	// Event details - q to hide details (alternative, global)
	if err := g.SetKeybinding("", 'q', gocui.ModNone, ui.hideEventDetails); err != nil {
		return err
	}

	// Debug - test if any key is being received
	if err := g.SetKeybinding("", 't', gocui.ModNone, ui.debugKeyPress); err != nil {
		return err
	}

	return nil
}

// quit quits the application
func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
	// If details popup is shown, hide it instead of quitting
	if ui.state.ShowDetails {
		return ui.hideEventDetails(g, v)
	}

	return gocui.ErrQuit
}

// moveUp moves the selection up
func (ui *UI) moveUp(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if cy > 0 {
		if err := v.SetCursor(0, cy-1); err != nil {
			return err
		}
	}
	return nil
}

// moveDown moves the selection down
func (ui *UI) moveDown(g *gocui.Gui, v *gocui.View) error {
	filteredEvents := ui.getFilteredEvents()
	lines := len(filteredEvents)
	_, cy := v.Cursor()
	if cy < lines-1 {
		if err := v.SetCursor(0, cy+1); err != nil {
			return err
		}
	}
	return nil
}

// moveLeft moves the selection left (previous page or previous item)
func (ui *UI) moveLeft(g *gocui.Gui, v *gocui.View) error {
	// For now, moveLeft behaves like moveUp for consistency
	return ui.moveUp(g, v)
}

// moveRight moves the selection right (next page or next item)
func (ui *UI) moveRight(g *gocui.Gui, v *gocui.View) error {
	// For now, moveRight behaves like moveDown for consistency
	return ui.moveDown(g, v)
}

// pageUp moves the selection up by a page (10 items)
func (ui *UI) pageUp(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	pageSize := 10
	newY := cy - pageSize
	if newY < 0 {
		newY = 0
	}
	if err := v.SetCursor(0, newY); err != nil {
		return err
	}
	return nil
}

// pageDown moves the selection down by a page (10 items)
func (ui *UI) pageDown(g *gocui.Gui, v *gocui.View) error {
	filteredEvents := ui.getFilteredEvents()
	lines := len(filteredEvents)
	_, cy := v.Cursor()
	pageSize := 10
	newY := cy + pageSize
	if newY >= lines {
		newY = lines - 1
	}
	if err := v.SetCursor(0, newY); err != nil {
		return err
	}
	return nil
}

// moveToTop moves the selection to the top of the list
func (ui *UI) moveToTop(g *gocui.Gui, v *gocui.View) error {
	if err := v.SetCursor(0, 0); err != nil {
		return err
	}
	return nil
}

// moveToBottom moves the selection to the bottom of the list
func (ui *UI) moveToBottom(g *gocui.Gui, v *gocui.View) error {
	filteredEvents := ui.getFilteredEvents()
	lines := len(filteredEvents)
	if lines > 0 {
		if err := v.SetCursor(0, lines-1); err != nil {
			return err
		}
	}
	return nil
}

// toggleFiles toggles file visibility
func (ui *UI) toggleFiles(g *gocui.Gui, v *gocui.View) error {
	ui.state.Filter.ShowFiles = !ui.state.Filter.ShowFiles
	ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		ui.updateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		ui.updateEventsView(v)
	}
	return nil
}

// toggleDirs toggles directory visibility
func (ui *UI) toggleDirs(g *gocui.Gui, v *gocui.View) error {
	ui.state.Filter.ShowDirs = !ui.state.Filter.ShowDirs
	ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		ui.updateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		ui.updateEventsView(v)
	}
	return nil
}

// toggleAggregate toggles event aggregation
func (ui *UI) toggleAggregate(g *gocui.Gui, v *gocui.View) error {
	wasAggregated := ui.state.AggregateEvents
	ui.state.AggregateEvents = !ui.state.AggregateEvents

	// If we're disabling aggregation, we need to "de-aggregate" existing events
	if !ui.state.AggregateEvents {
		ui.deaggregateEvents()
	} else if !wasAggregated {
		// If we're re-enabling aggregation, we need to re-aggregate existing events
		ui.reaggregateEvents()
	}

	ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		ui.updateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		ui.updateEventsView(v)
	}
	return nil
}

// deaggregateEvents splits aggregated events into individual events
func (ui *UI) deaggregateEvents() {
	var newEvents []*FileEvent

	for _, event := range ui.state.Events {
		if event.Count > 1 {
			// Create individual events for each count
			for i := 0; i < event.Count; i++ {
				newEvent := &FileEvent{
					Path:      event.Path,
					Operation: event.Operation,
					Timestamp: event.Timestamp,
					IsDir:     event.IsDir,
					Count:     1,
				}
				newEvents = append(newEvents, newEvent)
			}
		} else {
			// Keep single events as they are
			newEvents = append(newEvents, event)
		}
	}

	ui.state.Events = newEvents
}

// reaggregateEvents combines similar events that occurred within 1 second
func (ui *UI) reaggregateEvents() {
	if len(ui.state.Events) == 0 {
		return
	}

	var newEvents []*FileEvent
	eventMap := make(map[string]*FileEvent) // key: path+operation

	for _, event := range ui.state.Events {
		key := event.Path + "|" + event.Operation.String()

		if existingEvent, exists := eventMap[key]; exists {
			// Check if events are within 1 second of each other
			if event.Timestamp.Sub(existingEvent.Timestamp) < time.Second {
				existingEvent.Count++
				// Update timestamp to the most recent one
				if event.Timestamp.After(existingEvent.Timestamp) {
					existingEvent.Timestamp = event.Timestamp
				}
			} else {
				// Too much time has passed, create new event
				eventMap[key] = event
				newEvents = append(newEvents, event)
			}
		} else {
			eventMap[key] = event
			newEvents = append(newEvents, event)
		}
	}

	ui.state.Events = newEvents
}

// cycleSort cycles through sort options
func (ui *UI) cycleSort(g *gocui.Gui, v *gocui.View) error {
	ui.state.SortOption = (ui.state.SortOption + 1) % 4
	ui.state.ScrollOffset = 0

	if v, err := g.View(StatusView); err == nil {
		ui.updateStatusView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		ui.updateEventsView(v)
	}
	return nil
}

// getFileInfo gets file info for a path
func (ui *UI) getFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Public methods for testing

// GetState returns the current UI state
func (ui *UI) GetState() *UIState {
	return ui.state
}

// AddEvent adds an event (public version for testing)
func (ui *UI) AddEvent(path string, operation fsnotify.Op, isDir bool) {
	ui.addEvent(path, operation, isDir)
}

// ToggleAggregate toggles aggregation (public version for testing)
func (ui *UI) ToggleAggregate() {
	ui.state.AggregateEvents = !ui.state.AggregateEvents

	if !ui.state.AggregateEvents {
		ui.deaggregateEvents()
	} else {
		ui.reaggregateEvents()
	}
}

// ToggleFiles toggles file visibility (public version for testing)
func (ui *UI) ToggleFiles() {
	ui.state.Filter.ShowFiles = !ui.state.Filter.ShowFiles
}

// ToggleDirs toggles directory visibility (public version for testing)
func (ui *UI) ToggleDirs() {
	ui.state.Filter.ShowDirs = !ui.state.Filter.ShowDirs
}

// CycleSort cycles through sort options (public version for testing)
func (ui *UI) CycleSort() {
	ui.state.SortOption = (ui.state.SortOption + 1) % 4
}

// GetFilteredEvents returns filtered events (public version for testing)
func (ui *UI) GetFilteredEvents() []*FileEvent {
	return ui.getFilteredEvents()
}

// MoveUp moves the selection up (public version for testing)
func (ui *UI) MoveUp() {
	if ui.state.ScrollOffset > 0 {
		ui.state.ScrollOffset--
	}
}

// MoveDown moves the selection down (public version for testing)
func (ui *UI) MoveDown() {
	filteredEvents := ui.getFilteredEvents()
	if ui.state.ScrollOffset < len(filteredEvents)-1 {
		ui.state.ScrollOffset++
	}
}

// MoveLeft moves the selection left (public version for testing)
func (ui *UI) MoveLeft() {
	ui.MoveUp() // For now, same as MoveUp
}

// MoveRight moves the selection right (public version for testing)
func (ui *UI) MoveRight() {
	ui.MoveDown() // For now, same as MoveDown
}

// PageUp moves the selection up by a page (public version for testing)
func (ui *UI) PageUp() {
	pageSize := 10
	newOffset := ui.state.ScrollOffset - pageSize
	if newOffset < 0 {
		newOffset = 0
	}
	ui.state.ScrollOffset = newOffset
}

// PageDown moves the selection down by a page (public version for testing)
func (ui *UI) PageDown() {
	filteredEvents := ui.getFilteredEvents()
	pageSize := 10
	newOffset := ui.state.ScrollOffset + pageSize
	if newOffset >= len(filteredEvents) {
		newOffset = len(filteredEvents) - 1
	}
	if newOffset < 0 {
		newOffset = 0
	}
	ui.state.ScrollOffset = newOffset
}

// MoveToTop moves the selection to the top (public version for testing)
func (ui *UI) MoveToTop() {
	ui.state.ScrollOffset = 0
}

// MoveToBottom moves the selection to the bottom (public version for testing)
func (ui *UI) MoveToBottom() {
	filteredEvents := ui.getFilteredEvents()
	ui.state.ScrollOffset = len(filteredEvents) - 1
	if ui.state.ScrollOffset < 0 {
		ui.state.ScrollOffset = 0
	}
}

// showEventDetails shows the details popup for the currently selected event
func (ui *UI) showEventDetails(g *gocui.Gui, v *gocui.View) error {
	filteredEvents := ui.getFilteredEvents()
	if len(filteredEvents) == 0 {
		return nil
	}

	// Get current cursor position
	_, cy := v.Cursor()
	if cy < 0 || cy >= len(filteredEvents) {
		return nil
	}

	// Set the selected event and show details
	ui.state.SelectedEvent = filteredEvents[cy]
	ui.state.ShowDetails = true

	// Update the layout to show the popup
	return ui.layout(g)
}

// hideEventDetails hides the details popup
func (ui *UI) hideEventDetails(g *gocui.Gui, v *gocui.View) error {
	// Only hide details if popup is currently shown
	if ui.state.ShowDetails {
		ui.state.ShowDetails = false
		ui.state.SelectedEvent = nil

		// Update the layout to hide the popup
		return ui.layout(g)
	}

	// If popup is not shown, this might be a 'q' press in main view
	// Check if we're in the main view and should quit
	if v != nil && v.Name() == EventsView {
		return ui.quit(g, v)
	}

	return nil
}

// toggleEventDetails toggles the details popup (show if closed, hide if open)
func (ui *UI) toggleEventDetails(g *gocui.Gui, v *gocui.View) error {
	if ui.state.ShowDetails {
		// If popup is open, close it
		return ui.hideEventDetails(g, v)
	} else {
		// If popup is closed, open it (only if we're in EventsView)
		if v != nil && v.Name() == EventsView {
			return ui.showEventDetails(g, v)
		}
	}
	return nil
}
