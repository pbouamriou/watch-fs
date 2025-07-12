package ui

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

// FileEvent represents a file system event with additional metadata
type FileEvent struct {
	Path      string
	Operation fsnotify.Op
	Timestamp time.Time
	IsDir     bool
	Count     int // Number of events for this path in recent time
}

// Filter represents filtering options for events
type Filter struct {
	PathFilter      string
	OperationFilter fsnotify.Op
	ShowDirs        bool
	ShowFiles       bool
}

// SortOption represents sorting options
type SortOption int

const (
	SortByTime SortOption = iota
	SortByPath
	SortByOperation
	SortByCount
)

// UIState represents the current state of the UI
type UIState struct {
	Events          []*FileEvent
	Filter          Filter
	SortOption      SortOption
	SelectedPath    string
	ScrollOffset    int
	MaxEvents       int
	AggregateEvents bool       // Toggle for event aggregation
	ShowDetails     bool       // Toggle for details popup
	SelectedEvent   *FileEvent // Currently selected event for details
}

// ViewNames for the TUI
const (
	EventsView  = "events"
	StatusView  = "status"
	FilterView  = "filter"
	HelpView    = "help"
	DetailsView = "details"
)
