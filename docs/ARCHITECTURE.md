# Architecture

## Project Structure

```
watch-fs/
├── cmd/
│   └── watch-fs/          # Main application entry point
│       └── main.go
├── internal/
│   ├── ui/               # Terminal User Interface
│   │   ├── ui.go         # TUI implementation
│   │   └── types.go      # UI data structures
│   └── watcher/          # File system watcher
│       └── watcher.go    # Watcher wrapper
├── pkg/
│   └── utils/            # Utility functions
│       └── utils.go      # Common utilities
├── docs/                 # Documentation
├── scripts/              # Build and deployment scripts
├── test/                 # Test files
└── [configuration files]
```

## Package Organization

### `cmd/watch-fs/`

Contains the main application entry point. This follows Go's standard layout where executables are placed in `cmd/` directories.

### `internal/`

Contains packages that are private to this application and should not be imported by other projects.

#### `internal/ui/`

- **ui.go**: Implements the terminal user interface using gocui
- **types.go**: Defines data structures for the UI (FileEvent, Filter, etc.)

#### `internal/watcher/`

- **watcher.go**: Wraps fsnotify.Watcher with additional functionality for recursive watching

### `pkg/`

Contains packages that could potentially be reused by other projects.

#### `pkg/utils/`

- **utils.go**: Common utility functions for file system operations

## Design Patterns

### Interface Segregation

The UI package uses interfaces to decouple from the concrete watcher implementation:

```go
type WatcherInterface interface {
    Events() <-chan fsnotify.Event
    Errors() <-chan error
    AddDirectory(path string) error
}
```

### Separation of Concerns

- **UI Layer**: Handles user interaction and display
- **Watcher Layer**: Manages file system monitoring
- **Utils Layer**: Provides common functionality

## Dependencies

### External Libraries

- `github.com/fsnotify/fsnotify`: File system events
- `github.com/jroimartin/gocui`: Terminal UI framework
- `github.com/fatih/color`: Colored output

### Internal Dependencies

- `internal/watcher` → `fsnotify`
- `internal/ui` → `internal/watcher`, `gocui`, `color`
- `cmd/watch-fs` → `internal/ui`, `internal/watcher`, `pkg/utils`
