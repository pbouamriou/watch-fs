package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func watchRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("Watching:", path)
		}
		return nil
	})
}

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

	info, err := os.Stat(rootPath)
	if os.IsNotExist(err) || !info.IsDir() {
		log.Fatalf("Invalid directory: %s\n", rootPath)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watchRecursive(watcher, rootPath); err != nil {
		log.Fatal(err)
	}

	if useTUI {
		// Use TUI mode
		ui := NewUI(watcher, rootPath)
		if err := ui.Run(); err != nil {
			log.Fatal(err)
		}
	} else {
		// Use simple console mode (original behavior)
		done := make(chan bool)

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					fmt.Println("Event:", event)

					if event.Op&fsnotify.Create == fsnotify.Create {
						info, err := os.Stat(event.Name)
						if err == nil && info.IsDir() {
							_ = watchRecursive(watcher, event.Name)
						}
					}

				case err, ok := <-watcher.Errors:
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
