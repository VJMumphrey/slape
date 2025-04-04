// Package logging is for logging the performance and prompting outputs of the pipelines.
package logging

import (
	"os"
)

// CreateLogFile is used to check and see if a logfile is already created.
// If this is already done then func closes.
func CreateLogFile() error {
	err := os.WriteFile("logs.txt", []byte("--- SLaP-E Logging File ---\n"), 0644)
	if err != nil {
		return err
	}

    return nil
}

// WriteToFile writes to the log file using the io package in go.
// This is for efficiency and performance reasons.
func WriteToFile() {
    
}
