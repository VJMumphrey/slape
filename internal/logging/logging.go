// Package logging is for logging the performance and prompting outputs of the pipelines.
package logging

import (
	"log"
	"log/slog"
	"os"
)

const (
	logfilename = "logs.txt"
)

// CreateLogFile is used to check and see if a logfile is already created.
// It then creates a logger for the log file and returns it.
func CreateLogFile() {
	// Open the log file for writing
	logFile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error creating the log file")
	}
	defer logFile.Close()
}

func CreateLogger() *slog.Logger {

	// Open the log file for writing
	logFile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error setting up logging on startup")
	}
	defer logFile.Close()

	// Create a terminal handler (for stdout)
	//consoleHandler := slog.NewTextHandler(os.Stdout, nil)

	// Create a file handler (for the log file)
	fileHandler := slog.NewTextHandler(logFile, nil)

	// Create the logger with the combined handler
	logger := slog.New(fileHandler)

	return logger
}
