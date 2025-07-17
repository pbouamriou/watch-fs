package ui

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pbouamriou/watch-fs/pkg/logger"
)

// ExportImport manages export and import operations
type ExportImport struct {
	ui *UI
}

// NewExportImport creates a new ExportImport instance
func NewExportImport(ui *UI) *ExportImport {
	return &ExportImport{ui: ui}
}

// ExportEvents exports events to a file
func (ei *ExportImport) ExportEvents(filename string, format ExportFormat) error {
	switch format {
	case FormatSQLite:
		return ei.exportToSQLite(filename)
	case FormatJSON:
		return ei.exportToJSON(filename)
	default:
		return fmt.Errorf("unsupported export format")
	}
}

// ImportEvents imports events from a file
func (ei *ExportImport) ImportEvents(filename string, format ExportFormat) error {
	switch format {
	case FormatSQLite:
		return ei.importFromSQLite(filename)
	case FormatJSON:
		return ei.importFromJSON(filename)
	default:
		return fmt.Errorf("unsupported import format")
	}
}

// exportToSQLite exports events to SQLite database
func (ei *ExportImport) exportToSQLite(filename string) error {
	// Create database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err, "close error")
		}
	}()

	// Create events table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL,
		operation TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		is_dir BOOLEAN NOT NULL,
		count INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_events_path ON events(path);
	CREATE INDEX IF NOT EXISTS idx_events_operation ON events(operation);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	// Insert events
	insertSQL := `INSERT INTO events (path, operation, timestamp, is_dir, count) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			logger.Error(err, "stmt close error")
		}
	}()

	for _, event := range ei.ui.state.Events {
		_, err = stmt.Exec(event.Path, event.Operation.String(), event.Timestamp, event.IsDir, event.Count)
		if err != nil {
			return fmt.Errorf("failed to insert event: %w", err)
		}
	}

	return nil
}

// importFromSQLite imports events from SQLite database
func (ei *ExportImport) importFromSQLite(filename string) error {
	// Open database
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err, "close error")
		}
	}()

	// Query events
	rows, err := db.Query(`SELECT path, operation, timestamp, is_dir, count FROM events ORDER BY timestamp DESC`)
	if err != nil {
		return fmt.Errorf("failed to query events: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error(err, "rows close error")
		}
	}()

	var events []*FileEvent
	for rows.Next() {
		var path, operationStr string
		var timestamp time.Time
		var isDir bool
		var count int

		err := rows.Scan(&path, &operationStr, &timestamp, &isDir, &count)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse operation
		var operation fsnotify.Op
		switch operationStr {
		case "CREATE":
			operation = fsnotify.Create
		case "WRITE":
			operation = fsnotify.Write
		case "REMOVE":
			operation = fsnotify.Remove
		case "RENAME":
			operation = fsnotify.Rename
		case "CHMOD":
			operation = fsnotify.Chmod
		default:
			// Skip unknown operations
			continue
		}

		event := &FileEvent{
			Path:      path,
			Operation: operation,
			Timestamp: timestamp,
			IsDir:     isDir,
			Count:     count,
		}
		events = append(events, event)
	}

	// Replace current events
	ei.ui.state.Events = events

	return nil
}

// exportToJSON exports events to JSON file
func (ei *ExportImport) exportToJSON(filename string) error {
	// Create export data structure
	exportData := struct {
		Events []*FileEvent `json:"events"`
		Meta   struct {
			ExportTime time.Time `json:"export_time"`
			TotalCount int       `json:"total_count"`
		} `json:"meta"`
	}{
		Events: ei.ui.state.Events,
	}
	exportData.Meta.ExportTime = time.Now()
	exportData.Meta.TotalCount = len(ei.ui.state.Events)

	// Marshal to JSON
	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// importFromJSON imports events from JSON file
func (ei *ExportImport) importFromJSON(filename string) error {
	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal JSON
	var importData struct {
		Events []*FileEvent `json:"events"`
		Meta   struct {
			ExportTime time.Time `json:"export_time"`
			TotalCount int       `json:"total_count"`
		} `json:"meta"`
	}

	err = json.Unmarshal(data, &importData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Replace current events
	ei.ui.state.Events = importData.Events

	return nil
}

// GetRecentExportFiles returns the most recent export files
func (ei *ExportImport) GetRecentExportFiles() (string, string) {
	sqliteFiles, _ := filepath.Glob("watch-fs-events_*.db")
	jsonFiles, _ := filepath.Glob("watch-fs-events_*.json")

	var sqliteFile, jsonFile string
	if len(sqliteFiles) > 0 {
		sqliteFile = sqliteFiles[len(sqliteFiles)-1]
	}
	if len(jsonFiles) > 0 {
		jsonFile = jsonFiles[len(jsonFiles)-1]
	}

	return sqliteFile, jsonFile
}
