# watch-fs

A Go file watcher utility that recursively monitors a directory and its subdirectories to detect file changes with a beautiful terminal user interface (TUI) similar to lazygit.

## Features

- **Beautiful TUI Interface**: Modern terminal user interface with multiple views
- **Real-time Event Display**: Live updates with color-coded events
- **Recursive Directory Monitoring**: Automatically watches new subdirectories
- **Advanced Filtering**: Filter by file type, path, and operation
- **Multiple Sort Options**: Sort by time, path, operation, or event count
- **Event Aggregation**: Groups similar events with counters
- **Interactive Navigation**: Navigate through events with arrow keys
- **Robust Error Handling**: Graceful error handling and recovery

## Screenshots

The interface includes:

- **Status Bar**: Shows current directory, event count, and sort option
- **Filter Bar**: Displays current filters and toggle options
- **Events View**: Main area showing file events with colors and timestamps
- **Help Bar**: Keyboard shortcuts and navigation help

## Installation

```bash
go mod download
```

## Usage

### TUI Mode (Default)

```bash
go run main.go -path /path/to/directory
```

### Console Mode (Simple Output)

```bash
go run main.go -path /path/to/directory -tui=false
```

### Example

```bash
go run main.go -path ./my-project
```

## Options

- `-path` : The directory to watch (required)
- `-tui` : Use terminal user interface (default: true)

## TUI Controls

### Navigation

- **↑/↓** : Navigate through events
- **q** : Quit the application
- **Ctrl+C** : Quit the application

### Filtering

- **f** : Toggle file visibility
- **d** : Toggle directory visibility

### Sorting

- **s** : Cycle through sort options (Time → Path → Operation → Count)

## Event Types and Colors

- **CREATE** (Green) : File or directory creation
- **WRITE** (Yellow) : File modification
- **REMOVE** (Red) : File or directory deletion
- **RENAME** (Magenta) : File or directory renaming
- **CHMOD** (Blue) : Permission changes

## Sort Options

1. **Time** : Sort by event timestamp (newest first)
2. **Path** : Sort alphabetically by file path
3. **Operation** : Sort by operation type
4. **Count** : Sort by event frequency

## Event Aggregation

Similar events occurring within 1 second are automatically grouped with a counter, reducing noise and making it easier to track rapid changes.

## Dependencies

- [fsnotify](https://github.com/fsnotify/fsnotify) - File watching library
- [gocui](https://github.com/jroimartin/gocui) - Terminal user interface library
- [fatih/color](https://github.com/fatih/color) - Color output library

## Building

```bash
make build
```

## Development

```bash
make deps    # Install dependencies
make test    # Run tests
make clean   # Clean build artifacts
```

## License

MIT
