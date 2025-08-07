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

// ExportFormat represents the format for import/export
type ExportFormat int

const (
	FormatSQLite ExportFormat = iota
	FormatJSON
)

// FileDialogMode represents the mode of the file dialog
type FileDialogMode int

const (
	ModeSave FileDialogMode = iota
	ModeOpen
)

// FileEntry represents a file or directory in the file dialog
type FileEntry struct {
	Name    string
	IsDir   bool
	Size    int64
	ModTime time.Time
	Path    string
}

// FileDialogState represents the state of the file dialog
type FileDialogState struct {
	CurrentPath string
	Files       []*FileEntry
	SelectedIdx int
	Mode        FileDialogMode
	Filter      string // File extension filter
	Filename    string // Custom filename for save mode
	Placeholder bool   // Whether the filename is a placeholder
	IsEditing   bool   // Whether we're editing the filename
}

// ViewNames for the TUI
const (
	EventsView        = "events"
	StatusView        = "status"
	FilterView        = "filter"
	HelpView          = "help"
	DetailsView       = "details"
	ExportView        = "export"
	ImportView        = "import"
	FileDialogView    = "filedialog"
	FileListView      = "filelist"
	FilenameView      = "filename"
	PathView          = "path"
	FolderManagerView = "foldermanager"
	FolderListView    = "folderlist"
)

// FocusMode represents the current focus mode of the UI
type FocusMode int

const (
	FocusMain FocusMode = iota
	FocusDetails
	FocusFileDialog
	FocusFolderManager
	FocusWatchedFolders // Focus on "Currently Watching" panel
	FocusFolderBrowser  // Focus on "Available Folders" panel
)

// FolderManagerState represents the state of the folder manager
type FolderManagerState struct {
	CurrentPath  string
	SelectedIdx  int       // Selected index in the "Available Folders" panel
	WatchedIdx   int       // Selected index in the "Currently Watching" panel
	ScrollOffset int       // Scroll offset for long directory lists
	ActivePanel  FocusMode // Which panel is currently focused (FocusWatchedFolders or FocusFolderBrowser)
}

// UIState represents the current state of the UI
type UIState struct {
	Events            []*FileEvent
	Filter            Filter
	SortOption        SortOption
	SelectedPath      string
	ScrollOffset      int
	MaxEvents         int
	AggregateEvents   bool               // Toggle for event aggregation
	ShowDetails       bool               // Toggle for details popup
	SelectedEvent     *FileEvent         // Currently selected event for details
	ExportFilename    string             // Current export filename
	ImportFilename    string             // Current import filename
	ShowFileDialog    bool               // Toggle for file dialog
	ShowFolderManager bool               // Toggle for folder manager
	FileDialog        FileDialogState    // File dialog state
	FolderManager     FolderManagerState // Folder manager state
	CurrentFocus      FocusMode          // Current focus mode
}
