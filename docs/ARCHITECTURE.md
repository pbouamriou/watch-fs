# Architecture

## Overview

`watch-fs` is a file system monitoring tool with a Terminal User Interface (TUI) built using the `gocui` library. The application monitors file system events in real-time and provides an interactive interface for viewing, filtering, and managing these events.

## Core Components

### 1. Watcher (`internal/watcher/`)

- **Purpose**: Monitors file system events using `fsnotify`
- **Key Features**:
  - Real-time file system event monitoring
  - Support for multiple directories
  - Event aggregation and deduplication
  - Error handling and recovery

### 2. UI (`internal/ui/`)

- **Purpose**: Provides the Terminal User Interface
- **Key Features**:
  - Event display and navigation
  - Filtering and sorting capabilities
  - Interactive dialogs and popups
  - Context-aware help system

### 3. Utils (`pkg/utils/`)

- **Purpose**: Common utility functions
- **Key Features**:
  - File system operations
  - Data formatting and conversion
  - Export/import functionality

## Focus System

The application uses a sophisticated focus-based input handling system built on top of gocui's native focus management.

### Focus Modes

The UI operates in different focus modes, each with its own set of available actions:

1. **FocusMain** - Main interface (default)

   - Navigation: Arrow keys, hjkl, Page Up/Down, Home/End
   - Filtering: f (files), d (directories), a (aggregate), s (sort)
   - Actions: Enter (details), Ctrl+E (export), Ctrl+I (import), q (quit)

2. **FocusDetails** - Event details popup

   - Actions: ESC/q/Enter (close details)

3. **FocusExport** - Export dialog

   - Actions: Enter (confirm), ESC (cancel)
   - Input: Filename editing

4. **FocusImport** - Import dialog

   - Actions: Enter (confirm), ESC (cancel)
   - Input: Filename editing

5. **FocusFileDialog** - File selection dialog
   - Navigation: Arrow keys, kj
   - Actions: Enter (select), ESC/q (cancel)

### Key Benefits

- **Context-Aware Help**: The help bar displays only relevant shortcuts for the current focus
- **Clean Input Handling**: No manual focus checks or input conflicts
- **Better UX**: Users always know what actions are available
- **Maintainable Code**: Clear separation of concerns and easier debugging

### Implementation

The focus system is implemented through:

1. **State Management**: `CurrentFocus` field in `UIState` tracks the current mode
2. **View-Specific Keybindings**: Each view has its own set of keybindings
3. **Automatic Focus Switching**: Focus changes automatically when opening/closing dialogs
4. **Context-Aware Help**: `updateHelpView()` displays help based on `CurrentFocus`

## Event Flow

```
File System Event → Watcher → UI State → Display Update
```

1. **Event Detection**: `fsnotify` detects file system changes
2. **Event Processing**: Events are processed and aggregated
3. **State Update**: UI state is updated with new events
4. **Display Refresh**: Views are updated to reflect changes

## Data Structures

### FileEvent

```go
type FileEvent struct {
    Path      string
    Operation fsnotify.Op
    Timestamp time.Time
    IsDir     bool
    Count     int
}
```

### UIState

```go
type UIState struct {
    Events          []*FileEvent
    Filter          Filter
    SortOption      SortOption
    CurrentFocus    FocusMode
    // ... other fields
}
```

## Export/Import System

The application supports exporting and importing events in multiple formats:

- **SQLite**: High-performance, structured storage
- **JSON**: Human-readable, portable format

### Export Process

1. User triggers export (Ctrl+E)
2. Export dialog opens with focus
3. User enters filename
4. Events are serialized to chosen format
5. File is saved to disk

### Import Process

1. User triggers import (Ctrl+I)
2. Import dialog opens with focus
3. User enters filename
4. Events are deserialized from file
5. Events are loaded into UI state

## Error Handling

- **Graceful Degradation**: UI remains responsive even if file operations fail
- **User Feedback**: Status messages inform users of operation results
- **Recovery**: Application can recover from most error conditions

## Performance Considerations

- **Event Aggregation**: Similar events are combined to reduce memory usage
- **Lazy Loading**: File information is loaded only when needed
- **Efficient Rendering**: Views are updated incrementally
- **Memory Management**: Old events are automatically pruned

## Testing

The application includes comprehensive testing:

- **Unit Tests**: Individual component testing
- **Integration Tests**: End-to-end functionality testing
- **Focus System Tests**: Automated testing of the focus system
- **Manual Testing**: Interactive testing scripts

## Future Enhancements

- **Plugin System**: Extensible event processing
- **Network Support**: Remote monitoring capabilities
- **Advanced Filtering**: Regex and complex filter expressions
- **Event Replay**: Historical event playback
- **Performance Metrics**: Detailed performance monitoring
