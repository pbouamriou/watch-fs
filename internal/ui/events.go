package ui

import (
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jesseduffield/gocui"
)

// Events manages the logic of watcher events and their processing
// (adding, filtering, sorting, aggregation)
type Events struct {
	ui *UI
}

func NewEvents(ui *UI) *Events {
	return &Events{ui: ui}
}

// addEvent adds a new event to the state
func (e *Events) addEvent(path string, operation fsnotify.Op, isDir bool) {
	if e.ui.state.AggregateEvents {
		// Check if a similar event exists in the last second
		now := time.Now()
		for _, event := range e.ui.state.Events {
			if event.Path == path && event.Operation == operation &&
				now.Sub(event.Timestamp) < time.Second {
				event.Count++
				event.Timestamp = now
				return
			}
		}
	}

	// Add a new event
	event := &FileEvent{
		Path:      path,
		Operation: operation,
		Timestamp: time.Now(),
		IsDir:     isDir,
		Count:     1,
	}
	e.ui.state.Events = append(e.ui.state.Events, event)

	// Limit the number of events
	if len(e.ui.state.Events) > e.ui.state.MaxEvents {
		e.ui.state.Events = e.ui.state.Events[1:]
	}

	// Update the UI if initialized
	if e.ui.gui != nil {
		e.ui.gui.Update(func(g *gocui.Gui) error {
			if v, err := g.View(EventsView); err == nil {
				e.ui.views.UpdateEventsView(v)
			}
			if v, err := g.View(StatusView); err == nil {
				e.ui.views.UpdateStatusView(v)
			}
			if v, err := g.View(FilterView); err == nil {
				e.ui.views.UpdateFilterView(v)
			}
			return nil
		})
	}
}

// watchEvents listens to watcher events
func (e *Events) watchEvents() {
	for {
		select {
		case event, ok := <-e.ui.watcher.Events():
			if !ok {
				return
			}
			info, err := e.ui.getFileInfo(event.Name)
			isDir := false
			if err == nil {
				isDir = info.IsDir()
			}
			e.addEvent(event.Name, event.Op, isDir)
		case err, ok := <-e.ui.watcher.Errors():
			if !ok {
				return
			}
			// Ignore the error for now
			_ = err
		}
	}
}

// getFilteredEvents returns the filtered and sorted events
func (e *Events) getFilteredEvents() []*FileEvent {
	filtered := make([]*FileEvent, 0)
	for _, event := range e.ui.state.Events {
		// Filter path
		if e.ui.state.Filter.PathFilter != "" &&
			!strings.Contains(strings.ToLower(event.Path), strings.ToLower(e.ui.state.Filter.PathFilter)) {
			continue
		}
		// Filter operation
		if e.ui.state.Filter.OperationFilter != 0 && event.Operation != e.ui.state.Filter.OperationFilter {
			continue
		}
		// Filter type
		if event.IsDir && !e.ui.state.Filter.ShowDirs {
			continue
		}
		if !event.IsDir && !e.ui.state.Filter.ShowFiles {
			continue
		}
		filtered = append(filtered, event)
	}
	e.sortEvents(filtered)
	return filtered
}

// sortEvents sorts events according to the current option
func (e *Events) sortEvents(events []*FileEvent) {
	switch e.ui.state.SortOption {
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

// getSortOptionName returns the name of the current sort option
func (e *Events) getSortOptionName() string {
	switch e.ui.state.SortOption {
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
