# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Event Details Popup**: Press Enter on any event to view detailed information
  - Shows operation type, path, file/directory type, timestamp, count, size, permissions, and modification time
  - Popup is centered on screen with proper styling
  - Press Enter, Escape, or q to close the popup (Enter acts as toggle)
  - Added comprehensive test script to validate the functionality
  - Updated help text to include Enter: Details instruction

### Fixed

- **UNKNOWN Events**: Fixed issue where combined fsnotify operations were showing as "UNKNOWN"
  - Replaced direct comparison with `Has()` method to properly handle combined operations
  - All fsnotify event types (Create, Write, Remove, Rename, Chmod) are now correctly identified
  - Added comprehensive test script to validate the fix
  - See [docs/UNKNOWN_EVENTS_FIX.md](docs/UNKNOWN_EVENTS_FIX.md) for detailed explanation

### Added

- Test script `scripts/test_unknown_events.sh` to validate UNKNOWN events fix
- Test script `scripts/test_event_details.sh` to validate event details popup functionality
- Documentation for UNKNOWN events fix in `docs/UNKNOWN_EVENTS_FIX.md`

## [1.0.0] - 2024-01-XX

### Added

- Initial release of watch-fs
- Beautiful TUI interface with gocui
- Real-time file system event monitoring
- Recursive directory watching
- Advanced filtering and sorting options
- Event aggregation feature
- Interactive navigation with arrow keys and vim-style shortcuts
- Support for all fsnotify event types
- Comprehensive error handling
- Cross-platform compatibility (Linux, macOS, Windows)

### Features

- **Navigation**: Arrow keys, hjkl (vim-style), Page Up/Down, Home/End, g/G
- **Filtering**: Toggle files/directories, path filtering, operation filtering
- **Sorting**: Time, Path, Operation, Count
- **Event Types**: CREATE, WRITE, REMOVE, RENAME, CHMOD with color coding
- **Aggregation**: Group similar events with counters
- **Status Display**: Current directory, event count, sort option
- **Help System**: Built-in keyboard shortcuts reference
