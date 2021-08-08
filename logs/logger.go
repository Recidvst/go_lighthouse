package logs

import (
	"fmt"
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

	// check folders exist, if not then create them otherwise cli command will error
	_, checkDirErr := os.Stat(cwd + "/logs/logfiles"); os.IsNotExist(checkDirErr)
	fmt.Println("======================")
	mkdirErr := os.Mkdir(cwd + "/logs/logfiles", 0755)
	if mkdirErr != nil {
		log.Printf("Error: failed to create log folder with error: : %s", mkdirErr)
	}

	// if the log file doesn't exist, create it or append to the file
	file, makeFileErr := os.OpenFile(cwd + "/logs/logfiles/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if makeFileErr != nil {
		log.Printf("Error: failed to create log file with error: : %s", makeFileErr)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
