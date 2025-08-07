# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

watch-fs is a Go-based file system monitoring tool with a Terminal User Interface (TUI). It recursively monitors directories for file changes and provides real-time event display, filtering, sorting, and export capabilities.

## Common Development Commands

### Building and Testing

```bash
# Build the application
make build

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Install golangci-lint if needed
make lint-install

# Build for multiple platforms
make release-build

# Run in development mode
make dev

# Run built binary
make run
```

### Manual Testing

The project includes extensive manual testing scripts in the `test/` directory:

```bash
# Test multiple folder functionality
./test/test_multiple_folders.sh

# Test folder manager TUI
./test/test_folder_manager.sh

# Test import/export features
./test/test_import_export.sh

# Test keyboard navigation
./test/test_keyboard.sh
```

## Architecture

### Core Components

1. **Main Entry Point** (`cmd/watch-fs/main.go`)

   - Handles command-line arguments (`-path`, `-paths`, `-tui`)
   - Supports both single and multiple directory monitoring
   - Creates watcher instances and initializes UI

2. **Watcher** (`internal/watcher/watcher.go`)

   - Wraps `fsnotify.Watcher` with multi-root support
   - Thread-safe operations with `sync.RWMutex`
   - Manages root directories and watched subdirectories
   - Supports dynamic adding/removing of directories

3. **UI System** (`internal/ui/`)

   - **ui.go**: Main UI coordinator with component composition
   - **types.go**: Core data structures (FileEvent, UIState, FocusMode)
   - **events.go**: Event processing and aggregation
   - **navigation.go**: Keyboard navigation and filtering
   - **views.go**: View rendering and updates
   - **keybindings.go**: Key binding management
   - **layout.go**: TUI layout management
   - **folder_manager.go**: Dynamic folder management
   - **export_import.go**: SQLite/JSON export functionality
   - **file_dialog.go**: File browser dialog

4. **Utilities** (`pkg/`)
   - **logger/**: Structured logging
   - **utils/**: File system utilities and validation

### Focus-Based Input System

The UI uses a sophisticated focus system (`FocusMode` enum) for context-aware input handling:

- **FocusMain**: Main interface navigation and filtering
- **FocusDetails**: Event details popup
- **FocusFileDialog**: File browser for import/export
- **FocusFolderManager**: Dynamic folder management
- **FocusWatchedFolders**: Focus on "Currently Watching" panel within folder manager
- **FocusFolderBrowser**: Focus on "Available Folders" panel within folder manager

Each focus mode has dedicated keybindings and help displays, preventing input conflicts.

### Multi-Root Architecture

The watcher supports multiple directory monitoring with flexible command-line options:

- **Multiple `--path` flags** (preferred): `--path /dir1 --path /dir2 --path /dir3`
- **Legacy `--paths` flag**: comma-separated directory list `--paths "/dir1,/dir2,/dir3"`
- **Backward compatibility**: original `-path` flag still supported
- Dynamic folder management via TUI (Ctrl+F)
- Thread-safe operations for concurrent directory modifications

## Key Features to Understand

### Event Aggregation

Events occurring within 1 second are automatically grouped with counters to reduce noise. Toggle with `a` key.

### Export/Import System

- **SQLite Format**: High-performance structured storage (`.db` files)
- **JSON Format**: Human-readable portable format (`.json` files)
- File dialog with navigation and custom filename support

### Context-Aware Help

The help bar dynamically updates based on current focus mode, showing only relevant shortcuts.

### TUI Navigation

- Arrow keys, vim-style (hjkl), Page Up/Down, Home/End
- Mac French keyboard support with alternative bindings
- Smart quit behavior (`q` key adapts to context)

### Advanced Folder Manager

The folder manager (Ctrl+F) features a dual-panel interface with sophisticated focus management:

- **Left Panel**: "Currently Watching" - Shows active watched folders with individual file counts
- **Right Panel**: "Available Folders" - Browse filesystem to add new folders
- **Panel Switching**: Tab/Shift+Tab to switch panels, Left/Right arrows for directional switching
- **Visual Focus Indicators**: Focused panels have cyan titles and cyan frame borders
- **Per-folder Counts**: Each watched folder shows its specific monitored file count (not global total)
- **Standard Cursor Behavior**: Both panels use consistent gocui cursor highlighting

## Development Guidelines

### Code Patterns

1. **Component Composition**: UI uses composition pattern with specialized components
2. **Interface-Based Design**: Watcher uses interfaces for testability
3. **Thread Safety**: Use `sync.RWMutex` for shared data structures
4. **Error Handling**: Graceful degradation with user feedback
5. **Focus Management**: Always update `CurrentFocus` when changing UI state

### Adding New Features

1. **New TUI Views**: Add view constants to `types.go`, implement in `views.go`
2. **New Focus Modes**: Add to `FocusMode` enum, implement keybindings and help text
3. **New Export Formats**: Add to `ExportFormat` enum, implement in `export_import.go`
4. **Watcher Extensions**: Extend interface in `ui.go`, implement in `watcher.go`
5. **Folder Manager Enhancements**: Dual-panel system requires updating both focus modes and visual indicators
6. **Command-line Options**: Use custom flag types for complex argument parsing (see `pathsFlag` implementation)

### File Structure Conventions

- `internal/`: Private application code
- `pkg/`: Reusable packages
- `cmd/`: Application entry points
- `test/`: Manual testing scripts and unit tests
- `docs/`: Architecture and feature documentation

### Dependencies

- `github.com/fsnotify/fsnotify`: File system event notifications
- `github.com/jesseduffield/gocui`: Terminal UI framework (fork with enhanced features including `FrameColor` support)
- `github.com/fatih/color`: Colored terminal output  
- `github.com/mattn/go-sqlite3`: SQLite database driver

## Testing Strategy

### Automated Testing

The project includes comprehensive automated tests in the `test/` directory:

#### Unit Tests
- **Event Processing**: `ui_test.go` - Tests event aggregation, filtering, and sorting
- **Navigation Logic**: Tests for scroll position, cursor management, and boundary conditions
- **Folder Manager**: `folder_manager_test.go` - Tests directory navigation and selection logic
- **Scroll Mechanics**: `scroll_test.go` - Tests scroll position calculations and edge cases
- **Integration**: `ui_integration_test.go` - Tests UI state consistency and focus transitions

#### Test Helpers and Utilities
- **MockWatcher**: Simulates file system watcher for controlled testing
- **TestHelper**: Provides utilities for creating test directory structures
- **Performance Tests**: Tests UI responsiveness with large numbers of events
- **Concurrency Tests**: Validates thread-safe operations

#### Key Testing Patterns
```go
// Use TestHelper for directory-based tests
helper := NewTestHelperWithTempDir(t, "test-name")
defer helper.Cleanup()

// Create controlled directory structures
structure := map[string]interface{}{
    "dir1": map[string]interface{}{
        "file1.txt": "content",
        "subdir": map[string]interface{}{},
    },
}
baseDir := helper.CreateTestDirectoryStructure(t, structure)

// Test navigation and state consistency
helper.AssertEventCount(t, expectedCount)
helper.AssertScrollOffset(t, expectedOffset)
helper.AssertFocus(t, expectedFocus)
```

#### Running Tests
```bash
# Run all automated tests
go test ./test/... -v

# Run specific test suites
go test ./test/ -v -run TestFolderManager
go test ./test/ -v -run TestScroll

# Run with coverage
make test-coverage
```

### Manual Testing
Interactive TUI testing scripts for user experience validation:
- **Focus System Tests**: Automated focus state validation
- **Keyboard Navigation**: `test_navigation_fix.sh`
- **Multiple Folders**: `test_multiple_folders.sh` 
- **Import/Export**: `test_import_export.sh`

### TUI Testing Best Practices
1. **Mock Dependencies**: Use MockWatcher and MockView for isolated testing
2. **Real Directories**: Use temporary directories for realistic folder navigation tests
3. **State Validation**: Always verify UI state consistency after operations
4. **Boundary Testing**: Test edge cases like empty directories and long lists
5. **Concurrency Safety**: Test thread-safe operations with goroutines

## Release Process

```bash
# Create tagged release (requires clean working directory)
./scripts/release.sh 1.0.0

# Or prepare release artifacts only
make release
```

The release process includes testing, linting, and multi-platform binary generation.

## Technical Notes

- **UI Framework**: The UI is built on a fork of gocui maintained by the lazygit developer: github.com/jesseduffield/gocui
- **Documentation Language**: All documentation must be written in English (include CLAUDE.md and CLAUDE.local.md files)

## Advanced Implementation Techniques

### Dual-Panel Focus Management

The folder manager implements a sophisticated dual-panel focus system:

```go
// Key patterns for panel switching
type FocusMode int
const (
    FocusWatchedFolders   // Left panel focus
    FocusFolderBrowser    // Right panel focus
)

// Visual feedback synchronization
if activePanel == FocusWatchedFolders {
    v.TitleColor = gocui.ColorCyan
    v.FrameColor = gocui.ColorCyan  // Frame color matches title
}
```

**Critical Implementation Details:**
- Always synchronize `TitleColor` and `FrameColor` for consistent visual feedback
- Separate selection indices: `SelectedIdx` (right panel) vs `WatchedIdx` (left panel)
- Update focus in both layout creation AND existing view updates (two code paths)

### Per-Folder File Counting

Avoiding the common mistake of showing global totals instead of per-folder counts:

```go
// WRONG: Shows same total for all folders
totalWatched := counter.GetWatchedCount() // Global total

// CORRECT: Shows individual folder counts  
rootWatched := counter.GetWatchedCountForRoot(specificRoot) // Per-folder
```

**Implementation Pattern:**
- Use `filepath.Rel()` to determine if a path belongs to a specific root
- Avoid deprecated `filepath.HasPrefix()` - use proper relative path checking
- Clean and normalize paths before comparison

### Custom Flag Types for Complex CLI Arguments

For handling multiple identical flags (`--path` used multiple times):

```go
type pathsFlag []string
func (p *pathsFlag) String() string { return strings.Join(*p, ",") }
func (p *pathsFlag) Set(value string) error { 
    *p = append(*p, value); return nil 
}

// Usage: var paths pathsFlag; flag.Var(&paths, "path", "...")
```

**Benefits:** More intuitive than comma-separated strings, shell-friendly, maintains backward compatibility.

### gocui Visual Enhancement Patterns

The `jesseduffield/gocui` fork provides enhanced visual control:

```go
// Standard properties
v.Highlight = true
v.SelBgColor = gocui.Attribute(tcell.ColorDarkGreen)
v.SelFgColor = gocui.ColorBlack

// Enhanced frame control (fork-specific)
v.FrameColor = gocui.ColorCyan  // Border color
v.TitleColor = gocui.ColorCyan  // Title color
```

**Key Insight:** The fork's `FrameColor` property enables sophisticated visual focus indicators that standard gocui lacks.

### Testing Strategy for TUI Components

Effective patterns for testing complex UI interactions:

```bash
# Create focused test scripts for specific features
./test_enhanced_focus_visual.sh    # Visual focus system
./test_folder_count_fix.sh         # Per-folder counting
./test_multiple_path_flags.sh      # CLI argument parsing
```

**Testing Philosophy:** Combine automated unit tests with targeted manual TUI testing scripts for user experience validation.

## Development Practices

### Code Quality and Maintenance

- **Always use golangci-lint to check code rules after code modification.**