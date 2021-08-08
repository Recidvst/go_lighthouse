package logs

import (
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	DebugLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	
	var cwd, _ = os.Getwd()

	// check folder exists, if not then create them otherwise cli command will error
	_, checkDirErr := os.Stat(cwd + "/logs"); os.IsNotExist(checkDirErr)
	if checkDirErr != nil {
		mkdirErr := os.Mkdir(cwd + "/logs", 0755)
		log.Printf("Error: failed to create log folder with error: : %s", mkdirErr)
	}

	// if the log file doesn't exist, create it or append to the file
	file, makeFileErr := os.OpenFile(cwd + "/logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if makeFileErr != nil {
		log.Printf("Error: failed to create log file with error: : %s", makeFileErr)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
