package ui

import (
	"time"

	"github.com/jesseduffield/gocui"
)

// Navigation handles all navigation and filtering operations
type Navigation struct {
	ui *UI
}

// NewNavigation creates a new Navigation instance
func NewNavigation(ui *UI) *Navigation {
	return &Navigation{ui: ui}
}

// moveUp moves the selection up
func (nav *Navigation) moveUp(_ *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	if cy > 0 {
		v.SetCursor(0, cy-1)
	}
	return nil
}

// moveDown moves the selection down
func (nav *Navigation) moveDown(_ *gocui.Gui, v *gocui.View) error {
	filteredEvents := nav.ui.getFilteredEvents()
	lines := len(filteredEvents)
	_, cy := v.Cursor()
	if cy < lines-1 {
		v.SetCursor(0, cy+1)
	}
	return nil
}

// moveLeft moves the selection left (previous page or previous item)
func (nav *Navigation) moveLeft(g *gocui.Gui, v *gocui.View) error {
	// For now, moveLeft behaves like moveUp for consistency
	return nav.moveUp(g, v)
}

// moveRight moves the selection right (next page or next item)
func (nav *Navigation) moveRight(g *gocui.Gui, v *gocui.View) error {
	// For now, moveRight behaves like moveDown for consistency
	return nav.moveDown(g, v)
}

// pageUp moves the selection up by a page (10 items)
func (nav *Navigation) pageUp(_ *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	pageSize := 10
	newY := cy - pageSize
	if newY < 0 {
		newY = 0
	}
	v.SetCursor(0, newY)
	return nil
}

// pageDown moves the selection down by a page (10 items)
func (nav *Navigation) pageDown(_ *gocui.Gui, v *gocui.View) error {
	filteredEvents := nav.ui.getFilteredEvents()
	lines := len(filteredEvents)
	_, cy := v.Cursor()
	pageSize := 10
	newY := cy + pageSize
	if newY >= lines {
		newY = lines - 1
	}
	v.SetCursor(0, newY)
	return nil
}

// moveToTop moves the selection to the top of the list
func (nav *Navigation) moveToTop(_ *gocui.Gui, v *gocui.View) error {
	v.SetCursor(0, 0)
	return nil
}

// moveToBottom moves the selection to the bottom of the list
func (nav *Navigation) moveToBottom(_ *gocui.Gui, v *gocui.View) error {
	filteredEvents := nav.ui.getFilteredEvents()
	lines := len(filteredEvents)
	if lines > 0 {
		v.SetCursor(0, lines-1)
	}
	return nil
}

// toggleFiles toggles file visibility
func (nav *Navigation) toggleFiles(g *gocui.Gui, _ *gocui.View) error {
	nav.ui.state.Filter.ShowFiles = !nav.ui.state.Filter.ShowFiles
	nav.ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		nav.ui.views.UpdateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		nav.ui.views.UpdateEventsView(v)
	}
	return nil
}

// toggleDirs toggles directory visibility
func (nav *Navigation) toggleDirs(g *gocui.Gui, _ *gocui.View) error {
	nav.ui.state.Filter.ShowDirs = !nav.ui.state.Filter.ShowDirs
	nav.ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		nav.ui.views.UpdateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		nav.ui.views.UpdateEventsView(v)
	}
	return nil
}

// toggleAggregate toggles event aggregation
func (nav *Navigation) toggleAggregate(g *gocui.Gui, _ *gocui.View) error {
	wasAggregated := nav.ui.state.AggregateEvents
	nav.ui.state.AggregateEvents = !nav.ui.state.AggregateEvents

	// If we're disabling aggregation, we need to "de-aggregate" existing events
	if !nav.ui.state.AggregateEvents {
		nav.deaggregateEvents()
	} else if !wasAggregated {
		// If we're re-enabling aggregation, we need to re-aggregate existing events
		nav.reaggregateEvents()
	}

	nav.ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		nav.ui.views.UpdateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		nav.ui.views.UpdateEventsView(v)
	}
	return nil
}

// deaggregateEvents splits aggregated events into individual events
func (nav *Navigation) deaggregateEvents() {
	var newEvents []*FileEvent

	for _, event := range nav.ui.state.Events {
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

	nav.ui.state.Events = newEvents
}

// reaggregateEvents combines similar events that occurred within 1 second
func (nav *Navigation) reaggregateEvents() {
	if len(nav.ui.state.Events) == 0 {
		return
	}

	var newEvents []*FileEvent
	eventMap := make(map[string]*FileEvent) // key: path+operation

	for _, event := range nav.ui.state.Events {
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

	nav.ui.state.Events = newEvents
}

// cycleSort cycles through sort options
func (nav *Navigation) cycleSort(g *gocui.Gui, _ *gocui.View) error {
	nav.ui.state.SortOption = (nav.ui.state.SortOption + 1) % 4
	nav.ui.state.ScrollOffset = 0

	if v, err := g.View(FilterView); err == nil {
		nav.ui.views.UpdateFilterView(v)
	}
	if v, err := g.View(EventsView); err == nil {
		nav.ui.views.UpdateEventsView(v)
	}
	return nil
}

// Public methods for testing and external access

// MoveUp moves the selection up (public version for testing)
func (nav *Navigation) MoveUp() {
	if nav.ui.state.ScrollOffset > 0 {
		nav.ui.state.ScrollOffset--
	}
}

// MoveDown moves the selection down (public version for testing)
func (nav *Navigation) MoveDown() {
	filteredEvents := nav.ui.getFilteredEvents()
	if nav.ui.state.ScrollOffset < len(filteredEvents)-1 {
		nav.ui.state.ScrollOffset++
	}
}

// MoveLeft moves the selection left (public version for testing)
func (nav *Navigation) MoveLeft() {
	nav.MoveUp() // For now, same as MoveUp
}

// MoveRight moves the selection right (public version for testing)
func (nav *Navigation) MoveRight() {
	nav.MoveDown() // For now, same as MoveDown
}

// PageUp moves the selection up by a page (public version for testing)
func (nav *Navigation) PageUp() {
	pageSize := 10
	newOffset := nav.ui.state.ScrollOffset - pageSize
	if newOffset < 0 {
		newOffset = 0
	}
	nav.ui.state.ScrollOffset = newOffset
}

// PageDown moves the selection down by a page (public version for testing)
func (nav *Navigation) PageDown() {
	filteredEvents := nav.ui.getFilteredEvents()
	pageSize := 10
	newOffset := nav.ui.state.ScrollOffset + pageSize
	if newOffset >= len(filteredEvents) {
		newOffset = len(filteredEvents) - 1
	}
	if newOffset < 0 {
		newOffset = 0
	}
	nav.ui.state.ScrollOffset = newOffset
}

// MoveToTop moves the selection to the top (public version for testing)
func (nav *Navigation) MoveToTop() {
	nav.ui.state.ScrollOffset = 0
}

// MoveToBottom moves the selection to the bottom (public version for testing)
func (nav *Navigation) MoveToBottom() {
	filteredEvents := nav.ui.getFilteredEvents()
	nav.ui.state.ScrollOffset = len(filteredEvents) - 1
	if nav.ui.state.ScrollOffset < 0 {
		nav.ui.state.ScrollOffset = 0
	}
}

// ToggleAggregate toggles event aggregation (public version)
func (nav *Navigation) ToggleAggregate() {
	wasAggregated := nav.ui.state.AggregateEvents
	nav.ui.state.AggregateEvents = !nav.ui.state.AggregateEvents

	if !nav.ui.state.AggregateEvents {
		nav.deaggregateEvents()
	} else if !wasAggregated {
		nav.reaggregateEvents()
	}
}

// ToggleFiles toggles file visibility (public version)
func (nav *Navigation) ToggleFiles() {
	nav.ui.state.Filter.ShowFiles = !nav.ui.state.Filter.ShowFiles
}

// ToggleDirs toggles directory visibility (public version)
func (nav *Navigation) ToggleDirs() {
	nav.ui.state.Filter.ShowDirs = !nav.ui.state.Filter.ShowDirs
}

// CycleSort cycles through sort options (public version)
func (nav *Navigation) CycleSort() {
	nav.ui.state.SortOption = (nav.ui.state.SortOption + 1) % 4
}
