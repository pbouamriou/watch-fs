# Multiple Path Support - Usage Guide

## New Feature: Multiple `--path` Flags

You can now specify multiple directories to watch by using the `--path` flag multiple times.

## Usage Examples

### Single Directory (unchanged)
```bash
./watch-fs --path /Users/username/Documents
```

### Multiple Directories (NEW!)
```bash
./watch-fs --path /Users/username/Documents --path /Users/username/Downloads
```

```bash
./watch-fs --path /home/user/projects --path /home/user/photos --path /home/user/music
```

### Legacy Support (still works)
```bash
./watch-fs --paths "/dir1,/dir2,/dir3"
```

## Command Line Options

| Flag | Description | Example |
|------|-------------|---------|
| `--path` | Directory to watch (can be used multiple times) | `--path /dir1 --path /dir2` |
| `--paths` | Comma-separated directories (legacy) | `--paths "/dir1,/dir2"` |
| `--tui` | Use terminal interface (default: true) | `--tui=false` |
| `--version` | Show version information | `--version` |

## Priority Order

1. **Multiple `--path` flags** (preferred method)
2. **Legacy `--paths` flag** (comma-separated)
3. **Error** if none specified

## Folder Manager Integration

When using multiple paths:
- All directories appear in the "Currently Watching" panel
- Each directory shows its individual file count
- You can add/remove directories using the folder manager (Ctrl+F)
- Focus system works across all watched directories

## Benefits

- **More intuitive**: Natural command-line flag usage
- **Shell-friendly**: Easy to use with tab completion and scripting  
- **Flexible**: Mix different path types without escaping commas
- **Backward compatible**: Existing scripts continue to work

## Migration from Legacy

**Old way:**
```bash
./watch-fs --paths "/Users/john/Documents,/Users/john/Downloads"
```

**New way:**  
```bash
./watch-fs --path /Users/john/Documents --path /Users/john/Downloads
```

Both approaches work, but the new `--path` multiple flags approach is recommended for new usage.