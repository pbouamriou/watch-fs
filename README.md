# watch-fs

[![Go Report Card](https://goreportcard.com/badge/github.com/pbouamriou/watch-fs)](https://goreportcard.com/report/github.com/pbouamriou/watch-fs)
[![Go Version](https://img.shields.io/github/go-mod/go-version/pbouamriou/watch-fs)](https://go.dev/)
[![License](https://img.shields.io/github/license/pbouamriou/watch-fs)](https://github.com/pbouamriou/watch-fs/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/pbouamriou/watch-fs)](https://github.com/pbouamriou/watch-fs/releases)
[![Tests](https://github.com/pbouamriou/watch-fs/workflows/Tests/badge.svg)](https://github.com/pbouamriou/watch-fs/actions)
[![Codecov](https://codecov.io/gh/pbouamriou/watch-fs/branch/main/graph/badge.svg)](https://codecov.io/gh/pbouamriou/watch-fs)

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

## Project Structure

```
watch-fs/
├── cmd/watch-fs/         # Main application entry point
├── internal/ui/          # Terminal User Interface
├── internal/watcher/     # File system watcher
├── pkg/utils/           # Utility functions
├── docs/                # Documentation
└── [configuration files]
```

## Quick Start

### Installation

#### From Release (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/pbouamriou/watch-fs/releases).

#### From Source

```bash
git clone https://github.com/pbouamriou/watch-fs.git
cd watch-fs
make build
```

#### Using Go

```bash
go install github.com/pbouamriou/watch-fs/cmd/watch-fs@latest
```

### Usage

#### TUI Mode (Default)

```bash
watch-fs -path /path/to/directory
```

#### Console Mode (Simple Output)

```bash
watch-fs -path /path/to/directory -tui=false
```

#### Example

```bash
watch-fs -path ./my-project
```

## Options

- `-path` : The directory to watch (required)
- `-tui` : Use terminal user interface (default: true)

## TUI Controls

### Navigation

- **↑/↓/←/→** : Navigate through events (arrow keys)
- **h/j/k/l** : Alternative navigation (vim-style, for Mac French keyboards)
- **Page Up/Page Down** : Navigate by page (10 items at a time)
- **u/d** : Alternative page navigation (for Mac French keyboards)
- **Home/End** : Jump to top/bottom of the list
- **g/G** : Alternative top/bottom navigation (vim-style, for Mac French keyboards)
- **Enter** : Show event details popup
- **q** : Quit the application
- **Ctrl+C** : Quit the application

### Filtering

- **f** : Toggle file visibility
- **d** : Toggle directory visibility
- **a** : Toggle event aggregation

### Sorting

- **s** : Cycle through sort options (Time → Path → Operation → Count)

## Event Types and Colors

- **CREATE** (Green) : File or directory creation
- **WRITE** (Yellow) : File modification
- **REMOVE** (Red) : File or directory deletion
- **RENAME** (Magenta) : File or directory renaming
- **CHMOD** (Blue) : Permission changes

> **Note**: All fsnotify event types are properly supported, including combined operations. Events previously showing as "UNKNOWN" are now correctly identified.

## Event Details Popup

Press **Enter** on any event to view detailed information in a popup window. The popup shows:

- **Operation** : Type of file system operation with color coding
- **Path** : Full file or directory path
- **Type** : Whether it's a file or directory
- **Timestamp** : Exact time with milliseconds precision
- **Count** : Number of similar events (when aggregation is enabled)
- **Size** : File size in bytes (for files)
- **Permissions** : File permissions and mode
- **Modified** : Last modification time

Press **Enter**, **Escape**, or **q** to close the details popup. When the popup is open, **q** closes the popup instead of quitting the application.

## Sort Options

1. **Time** : Sort by event timestamp (newest first)
2. **Path** : Sort alphabetically by file path
3. **Operation** : Sort by operation type
4. **Count** : Sort by event frequency

## Event Aggregation

Similar events occurring within 1 second are automatically grouped with a counter, reducing noise and making it easier to track rapid changes. You can toggle this feature on/off using the **a** key.

**When enabled (default)**: Similar events are grouped together with a counter
**When disabled**: All individual events are shown separately

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for build scripts)

### Setup

```bash
git clone https://github.com/pbouamriou/watch-fs.git
cd watch-fs
make deps    # Install dependencies
make build   # Build the application
make test    # Run tests
```

### Testing

```bash
# Run all tests
go test ./test/...

# Run tests with coverage
go test -cover ./test/...

# Run linter
golangci-lint run
```

### Creating a Release

```bash
# Create a new release (requires clean working directory)
./scripts/release.sh 1.0.0
```

## Dependencies

- [fsnotify](https://github.com/fsnotify/fsnotify) - File watching library
- [gocui](https://github.com/jroimartin/gocui) - Terminal user interface library
- [fatih/color](https://github.com/fatih/color) - Color output library

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## Architecture

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for detailed architecture information.

## License

MIT License - see [LICENSE](LICENSE) file for details.
