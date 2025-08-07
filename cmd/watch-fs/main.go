package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pbouamriou/watch-fs/internal/ui"
	"github.com/pbouamriou/watch-fs/internal/watcher"
	"github.com/pbouamriou/watch-fs/pkg/logger"
	"github.com/pbouamriou/watch-fs/pkg/utils"
)

// Version will be set by the linker during build
var version = "dev"

// Custom flag type for multiple --path flags
type pathsFlag []string

func (p *pathsFlag) String() string {
	return strings.Join(*p, ",")
}

func (p *pathsFlag) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func main() {
	// Initialise le logger
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur d'initialisation du logger: %v\n", err)
		os.Exit(1)
	}

	var paths string // Legacy flag for comma-separated paths
	var useTUI bool
	var showVersion bool
	var pathsVar pathsFlag
	flag.Var(&pathsVar, "path", "Directory to watch (can be used multiple times)")
	flag.StringVar(&paths, "paths", "", "Comma-separated list of directories to watch (legacy)")
	flag.BoolVar(&useTUI, "tui", true, "Use terminal user interface (default: true)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("watch-fs version %s\n", version)
		os.Exit(0)
	}

	// Parse paths - priority: multiple --path flags > legacy --paths > error
	var rootPaths []string
	if len(pathsVar) > 0 {
		// Use multiple --path flags (preferred method)
		rootPaths = []string(pathsVar)
		// Trim whitespace from each path
		for i, path := range rootPaths {
			rootPaths[i] = strings.TrimSpace(path)
		}
	} else if paths != "" {
		// Use the legacy --paths flag (comma-separated)
		rootPaths = strings.Split(paths, ",")
		for i, path := range rootPaths {
			rootPaths[i] = strings.TrimSpace(path)
		}
	} else {
		fmt.Println("Error: at least one --path flag is required")
		fmt.Println("Usage:")
		fmt.Println("  watch-fs --path /single/directory")
		fmt.Println("  watch-fs --path /dir1 --path /dir2 --path /dir3")
		fmt.Println("  watch-fs --paths '/dir1,/dir2,/dir3'  (legacy)")
		flag.Usage()
		os.Exit(1)
	}

	// Validate all directories
	for _, path := range rootPaths {
		if err := utils.ValidateDirectory(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid directory '%s': %v\n", path, err)
			os.Exit(1)
		}
	}

	// Create watcher
	var fileWatcher *watcher.Watcher
	var err error

	if len(rootPaths) == 1 {
		// Use the original constructor for single directory (backward compatibility)
		fileWatcher, err = watcher.New(rootPaths[0])
	} else {
		// Use the new multi-root constructor for multiple directories
		fileWatcher, err = watcher.NewMultiRoot(rootPaths)
	}

	if err != nil {
		logger.Error(err, "Failed to create watcher")
		os.Exit(1)
	}

	// Add recursive watching for all roots
	if err := fileWatcher.AddAllRootsRecursive(); err != nil {
		logger.Error(err, "Failed to add recursive watching")
		os.Exit(1)
	}
	defer func() {
		if err := fileWatcher.Close(); err != nil {
			logger.Error(err, "Failed to close watcher")
		}
	}()

	// Get the primary root path for UI (first one for backward compatibility)
	primaryRootPath := rootPaths[0]

	if useTUI {
		// Use TUI mode
		ui := ui.NewUI(fileWatcher, primaryRootPath)
		if err := ui.Run(); err != nil {
			logger.Error(err, "TUI exited with error")
			os.Exit(1)
		}
	} else {
		// Use simple console mode (original behavior)
		done := make(chan bool)

		go func() {
			for {
				select {
				case event, ok := <-fileWatcher.Events():
					if !ok {
						return
					}
					fmt.Println("Event:", event)

					if event.Op&fsnotify.Create == fsnotify.Create {
						info, err := os.Stat(event.Name)
						if err == nil && info.IsDir() {
							_ = fileWatcher.AddDirectory(event.Name)
						}
					}

				case err, ok := <-fileWatcher.Errors():
					if !ok {
						return
					}
					logger.Error(err, "Watcher error")
				}
			}
		}()

		<-done
	}
}
