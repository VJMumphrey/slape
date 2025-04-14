// Package logging is for logging the performance and prompting outputs of the pipelines.
package logging

import (
	"log"
	"os"

	"github.com/StoneG24/slape/pkg/vars"
)

// CreateLogFile is used to check and see if a logfile is already created.
// It then creates a logger for the log file and returns it.
func CreateLogFile() *os.File {

	// Open the log file for writing
	logFile, err := os.OpenFile(vars.Logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error creating the log file")
	}
	log.SetOutput(logFile)

    return logFile
}
