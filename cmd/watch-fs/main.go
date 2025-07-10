package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	ui "github.com/pbouamriou/watch-fs/internal/ui"
	"github.com/pbouamriou/watch-fs/internal/watcher"
	"github.com/pbouamriou/watch-fs/pkg/utils"
)

func main() {
	var rootPath string
	var useTUI bool
	flag.StringVar(&rootPath, "path", "", "Directory to watch (required)")
	flag.BoolVar(&useTUI, "tui", true, "Use terminal user interface (default: true)")
	flag.Parse()

	if rootPath == "" {
		fmt.Println("Error: -path flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Validate directory
	if err := utils.ValidateDirectory(rootPath); err != nil {
		log.Fatalf("Invalid directory: %s\n", rootPath)
	}

	// Create watcher
	fileWatcher, err := watcher.New(rootPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fileWatcher.Close()

	// Add recursive watching
	if err := fileWatcher.AddRecursive(rootPath); err != nil {
		log.Fatal(err)
	}

	if useTUI {
		// Use TUI mode
		ui := ui.NewUI(fileWatcher, rootPath)
		if err := ui.Run(); err != nil {
			log.Fatal(err)
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
					log.Println("Watcher error:", err)
				}
			}
		}()

		<-done
	}
}
