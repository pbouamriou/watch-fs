package logger

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// Init initializes the logger with file output
func Init() error {
	// Determine log file path based on OS
	var logDir string
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		logDir = filepath.Join(localAppData, "watch-fs")
	case "linux", "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		logDir = filepath.Join(homeDir, ".cache", "watch-fs")
	default:
		// Fallback to current directory
		logDir = "."
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Create log file
	logFile := filepath.Join(logDir, "watch-fs.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Logger = zerolog.New(file).With().Timestamp().Logger()

	// Set global logger
	log.Logger = Logger

	return nil
}

// Error logs an error message
func Error(err error, msg string) {
	Logger.Error().Err(err).Msg(msg)
}

// Info logs an info message
func Info(msg string) {
	Logger.Info().Msg(msg)
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	Logger.Warn().Msg(msg)
}
