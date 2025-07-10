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
			AggregateEvents: true, // Enable aggregation by default
		},
		watcher:  watcher,
		rootPath: rootPath,
	}
}

// Run starts the TUI
func (ui *UI) Run() error {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer gui.Close()

	ui.gui = gui
	gui.SetManagerFunc(ui.layout)

	// Set up keybindings
	if err := ui.keybindings(gui); err != nil {
		return err
	}

	// Start event watcher goroutine
	go ui.watchEvents()

	// Start the main loop
	return gui.MainLoop()
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

	return nil
}

// updateStatusView updates the status view
func (ui *UI) updateStatusView(v *gocui.View) {
	v.Clear()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Fprintf(v, "Watching: %s | Events: %s | Sort: %s\n",
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

	fmt.Fprintf(v, "Dirs: %s | Files: %s | Aggregate: %s | Path Filter: %s",
		dirsStatus, filesStatus, aggregateStatus, ui.state.Filter.PathFilter)
}

// updateEventsView updates the events view
func (ui *UI) updateEventsView(v *gocui.View) {
	v.Clear()

	filteredEvents := ui.getFilteredEvents()

	if len(filteredEvents) == 0 {
		fmt.Fprintf(v, "No events to display")
		return
	}

	// Apply scroll offset
	start := ui.state.ScrollOffset
	if start >= len(filteredEvents) {
		start = len(filteredEvents) - 1
	}
	if start < 0 {
		start = 0
	}

	end := start + 20 // Show 20 events at a time
	if end > len(filteredEvents) {
		end = len(filteredEvents)
	}

	for i := start; i < end; i++ {
		event := filteredEvents[i]
		ui.renderEvent(v, event, i == start) // Highlight first visible event
	}
}

// updateHelpView updates the help view
func (ui *UI) updateHelpView(v *gocui.View) {
	v.Clear()
	fmt.Fprintf(v, "q: Quit | f: Toggle files | d: Toggle dirs | a: Toggle aggregate | s: Sort | /: Filter | ↑↓: Navigate | Enter: Select")
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
	switch event.Operation {
	case fsnotify.Create:
		operationStr = green("CREATE")
	case fsnotify.Write:
		operationStr = yellow("WRITE")
	case fsnotify.Remove:
		operationStr = red("REMOVE")
	case fsnotify.Rename:
		operationStr = magenta("RENAME")
	case fsnotify.Chmod:
		operationStr = blue("CHMOD")
	default:
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
		fmt.Fprintf(v, "> %s\n", line)
	} else {
		fmt.Fprintf(v, "  %s\n", line)
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

	// Update UI
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

// keybindings sets up the key bindings
func (ui *UI) keybindings(g *gocui.Gui) error {
	// Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, ui.quit); err != nil {
		return err
	}

	// Navigation
	if err := g.SetKeybinding(EventsView, gocui.KeyArrowUp, gocui.ModNone, ui.moveUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(EventsView, gocui.KeyArrowDown, gocui.ModNone, ui.moveDown); err != nil {
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

	return nil
}

// quit quits the application
func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// moveUp moves the selection up
func (ui *UI) moveUp(g *gocui.Gui, v *gocui.View) error {
	if ui.state.ScrollOffset > 0 {
		ui.state.ScrollOffset--
		ui.updateEventsView(v)
	}
	return nil
}

// moveDown moves the selection down
func (ui *UI) moveDown(g *gocui.Gui, v *gocui.View) error {
	filteredEvents := ui.getFilteredEvents()
	if ui.state.ScrollOffset < len(filteredEvents)-1 {
		ui.state.ScrollOffset++
		ui.updateEventsView(v)
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
	ui.state.AggregateEvents = !ui.state.AggregateEvents
	ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		ui.updateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		ui.updateEventsView(v)
	}
	return nil
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
