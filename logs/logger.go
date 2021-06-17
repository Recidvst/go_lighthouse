package logs

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	
	var cwd, _ = os.Getwd()

	// check folder exist, if not then create them otherwise cli command will error
	if _, checkDirErr := os.Stat(cwd + "/logs"); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir(cwd + "/logs", 0755)
		if mkdirErr != nil {
			return
		}
		if _, checkDirErr := os.Stat(cwd + "logs/files"); os.IsNotExist(checkDirErr) {
			mkdirErr := os.Mkdir(cwd + "/logs/files", 0755)
			if mkdirErr != nil {
				return
			}

			// If the file doesn't exist, create it or append to the file
			file, err := os.OpenFile(cwd + "/logs/files/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Fatal(err)
			}

			InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
			WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
			ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		}
	}
}
