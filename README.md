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
- **Multiple Directory Monitoring**: Watch multiple directories simultaneously with the new `-paths` flag
- **Real-time Event Display**: Live updates with color-coded events
- **Recursive Directory Monitoring**: Automatically watches new subdirectories
- **Advanced Filtering**: Filter by file type, path, and operation
- **Multiple Sort Options**: Sort by time, path, operation, or event count
- **Event Aggregation**: Groups similar events with counters
- **Interactive Navigation**: Navigate through events with arrow keys
- **Event Details Popup**: View detailed information about any event
- **Import/Export Functionality**: Save and load events in SQLite or JSON format
- **Context-Aware Focus System**: Intelligent input handling with context-sensitive help
- **Robust Error Handling**: Graceful error handling and recovery
- **Backward Compatibility**: Original `-path` flag still supported

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

#### Single Directory (Original)

```bash
watch-fs -path /path/to/directory
```

#### Multiple Directories (New)

```bash
# Surveiller plusieurs dossiers
watch-fs -paths "/path/to/dir1,/path/to/dir2,/path/to/dir3"

# Avec des chemins relatifs
watch-fs -paths "./src,./tests,./docs"

# Avec des espaces (automatiquement supprimés)
watch-fs -paths "/path/to/dir1, /path/to/dir2"
```

#### Console Mode (Simple Output)

```bash
watch-fs -path /path/to/directory -tui=false
# ou
watch-fs -paths "/dir1,/dir2" -tui=false
```

#### Examples

```bash
# Surveillance d'un projet
watch-fs -path ./my-project

# Surveillance de plusieurs dossiers de développement
watch-fs -paths "./frontend/src,./backend/src,./shared"

# Surveillance de dossiers système
watch-fs -paths "/var/log,/tmp,/home/user/documents"
```

> **Note**: Le flag `-path` reste supporté pour la compatibilité. Le nouveau flag `-paths` permet de surveiller plusieurs dossiers simultanément.

## Options

- `-path` : The directory to watch (deprecated, use -paths instead)
- `-paths` : Comma-separated list of directories to watch
- `-tui` : Use terminal user interface (default: true)
- `-version` : Show version information

## TUI Controls

### Navigation

- **↑/↓/←/→** : Navigate through events (arrow keys)
- **h/j/k/l** : Alternative navigation (vim-style, for Mac French keyboards)
- **Page Up/Page Down** : Navigate by page (10 items at a time)
- **Home/End** : Jump to top/bottom of the list
- **g/G** : Alternative top/bottom navigation (vim-style, for Mac French keyboards)
- **Enter** : Show event details popup
- **Ctrl+E** : Export events to file (SQLite/JSON)
- **Ctrl+I** : Import events from file
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

## Import/Export Functionality

watch-fs supports importing and exporting file system events to external files for analysis, backup, or sharing.

### Export Formats

- **SQLite Database** (Recommended): Fast, indexed database format for large datasets
- **JSON Format**: Human-readable format for sharing and manual inspection

### Usage

- **Ctrl+E**: Open file dialog to save events (navigate and select location)
- **Ctrl+I**: Open file dialog to load events (browse and select file)
- **File Navigation**: Use arrow keys or hjkl to navigate directories
- **File Selection**: Enter to open directories or select files
- **Automatic Format Detection**: `.db` for SQLite, `.json` for JSON format
- **Status Bar**: Shows "Export: SQLite available" or "Export: JSON available" when files exist

### File Dialog Features

The file dialog provides a full-featured file browser with:

- **Navigation**: Arrow keys or hjkl to move through files and directories
- **Directory Browsing**: Enter to open directories, ".." to go back
- **File Filtering**: Automatically filters by file type (.db, .json)
- **Hidden Files**: Hidden files are automatically filtered out
- **File Information**: Shows file sizes and modification dates
- **Mode Awareness**: Different behavior for Save vs Open modes
- **Custom Filename (Save Mode)**: Press 'e' to edit a custom filename instead of overwriting existing files

### SQLite Analysis

Export files can be analyzed with any SQLite browser or command-line tools:

```bash
# Open exported database
sqlite3 /path/to/your/exported/events.db

# Example queries
SELECT COUNT(*) FROM events;
SELECT operation, COUNT(*) FROM events GROUP BY operation;
SELECT * FROM events WHERE path LIKE '%config%';
```

For detailed information, see [docs/IMPORT_EXPORT_FEATURE.md](docs/IMPORT_EXPORT_FEATURE.md).

## Sort Options

1. **Time** : Sort by event timestamp (newest first)
2. **Path** : Sort alphabetically by file path
3. **Operation** : Sort by operation type
4. **Count** : Sort by event frequency

## Event Aggregation

Similar events occurring within 1 second are automatically grouped with a counter, reducing noise and making it easier to track rapid changes. You can toggle this feature on/off using the **a** key.

**When enabled (default)**: Similar events are grouped together with a counter
**When disabled**: All individual events are shown separately

## Focus System

watch-fs uses an intelligent focus-based input handling system that provides context-aware controls and help. The interface automatically adapts to show only relevant shortcuts for the current context.

### Focus Modes

The application operates in different focus modes, each with its own set of available actions:

#### Main Interface (Default)

- **Navigation**: Arrow keys, hjkl, Page Up/Down, Home/End
- **Filtering**: f (files), d (directories), a (aggregate), s (sort)
- **Actions**: Enter (details), Ctrl+E (export), Ctrl+I (import), q (quit)

#### Event Details Popup

- **Actions**: ESC/q/Enter (close details)
- **Context**: View detailed event information

#### Export Dialog

- **Actions**: Enter (confirm), ESC (cancel)
- **Input**: Type filename with .db or .json extension

#### Import Dialog

- **Actions**: Enter (confirm), ESC (cancel)
- **Input**: Type filename to import

#### File Selection Dialog

- **Navigation**: Arrow keys, kj
- **Actions**: Enter (select), ESC/q (cancel)
- **Context**: Browse directories and select files

### Benefits

- **Context-Aware Help**: The help bar always shows relevant shortcuts for the current mode
- **No Input Conflicts**: Each mode has its own dedicated keybindings
- **Intuitive Navigation**: Focus automatically switches when opening/closing dialogs
- **Better UX**: Users always know what actions are available
- **Cleaner Code**: No manual focus checks or complex input handling

### Smart Quit Behavior

The **q** key intelligently adapts to the current context:

- **Main interface**: Quits the application
- **Details popup**: Closes the popup
- **Export/Import dialog**: Cancels the operation
- **File dialog**: Closes the file browser

This eliminates the need to remember different quit shortcuts for different contexts.

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

## Multiple Directory Monitoring

For detailed information about the multiple directory monitoring feature, see [docs/MULTIPLE_FOLDERS.md](docs/MULTIPLE_FOLDERS.md).

## License

MIT License - see [LICENSE](LICENSE) file for details.
