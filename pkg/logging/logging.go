// Package logging is for logging the performance and prompting outputs of the pipelines.
package logging

import (
	"errors"
	"log"
	"os"

	"github.com/StoneG24/slape/pkg/vars"
)

// CreateLogFile is used to check and see if a logfile is already created.
// It then creates a logger for the log file and returns it.
func CreateLogFile() *os.File {
	// check if file is present
	// create if not, if present trucate for the current run.
	var logFile *os.File
	if _, err := os.Stat("./logs/" + vars.Logfilename); errors.Is(err, os.ErrNotExist) {
		// Open the log file for writing
		logFile, err = os.OpenFile("./logs/"+vars.Logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error creating the log file")
		}
	} else {
		// clear the file
		os.Truncate(vars.Logfilename, 0)
	}

	log.SetOutput(logFile)
	return logFile
}

// CloseLogging moves the current runs logs to a trunk log file.
func CloseLogging(file *os.File) {

	// read all of the contents of the log file
	content, err := os.ReadFile("./logs/" + vars.Logfilename)
	if err != nil {
		log.Println("Error trying to read ", vars.Logfilename)
	}

	// open the trunk file for appending
	trunk, err := os.OpenFile("./logs/"+vars.Trunkfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error creating the trunk file")
	}

	// move to the big log file
	n, err := trunk.Write(content)
	if err != nil || n != len(content) {
		log.Println("Error while trying to update old logs with current run")
	}

	trunk.Close()
}
