package main

import (
	"flag"
	"fmt"
	"os"

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
	var useTUI bool
	var showVersion bool
	flag.StringVar(&rootPath, "path", "", "Directory to watch (required)")
	flag.BoolVar(&useTUI, "tui", true, "Use terminal user interface (default: true)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("watch-fs version %s\n", version)
		os.Exit(0)
	}

	if rootPath == "" {
		fmt.Println("Error: -path flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate directory
	if err := utils.ValidateDirectory(rootPath); err != nil {
		logger.Error(err, "Invalid directory")
		os.Exit(1)
	}

	// Create watcher
	fileWatcher, err := watcher.New(rootPath)
	if err != nil {
		logger.Error(err, "Failed to create watcher")
		os.Exit(1)
	}

	// Add recursive watching
	if err := fileWatcher.AddRecursive(rootPath); err != nil {
		logger.Error(err, "Failed to add recursive watching")
		os.Exit(1)
	}
	defer func() { _ = fileWatcher.Close() }()

	if useTUI {
		// Use TUI mode
		ui := ui.NewUI(fileWatcher, rootPath)
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
