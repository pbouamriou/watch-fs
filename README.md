# watch-fs

A Go file watcher utility that recursively monitors a directory and its subdirectories to detect file changes.

## Features

- Recursive directory monitoring
- Automatic detection of new subdirectories
- Real-time file event display
- Robust error handling

## Installation

```bash
go mod download
```

## Usage

```bash
go run main.go -path /path/to/directory
```

### Example

```bash
go run main.go -path ./my-project
```

## Options

- `-path` : The directory to watch (required)

## Monitored Events

The program monitors the following events:

- File and directory creation
- File modification
- File and directory deletion
- File and directory renaming

## Dependencies

- [fsnotify](https://github.com/fsnotify/fsnotify) - File watching library

## License

MIT
