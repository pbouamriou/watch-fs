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

func main() {
	// Initialise le logger
	if err := logger.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur d'initialisation du logger: %v\n", err)
		os.Exit(1)
	}

	var rootPath string
	var paths string // New flag for multiple paths
	var useTUI bool
	var showVersion bool
	flag.StringVar(&rootPath, "path", "", "Directory to watch (deprecated, use -paths instead)")
	flag.StringVar(&paths, "paths", "", "Comma-separated list of directories to watch")
	flag.BoolVar(&useTUI, "tui", true, "Use terminal user interface (default: true)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("watch-fs version %s\n", version)
		os.Exit(0)
	}

	// Parse paths
	var rootPaths []string
	if paths != "" {
		// Use the new -paths flag
		rootPaths = strings.Split(paths, ",")
		for i, path := range rootPaths {
			rootPaths[i] = strings.TrimSpace(path)
		}
	} else if rootPath != "" {
		// Use the deprecated -path flag for backward compatibility
		rootPaths = []string{rootPath}
	} else {
		fmt.Println("Error: either -path or -paths flag is required")
		fmt.Println("Usage:")
		fmt.Println("  watch-fs -path /single/directory")
		fmt.Println("  watch-fs -paths '/dir1,/dir2,/dir3'")
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
	defer func() { _ = fileWatcher.Close() }()

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
