package database

import (
	"database/sql"
	"fmt"
	LOGS "go_svelte_lighthouse/logs"
	"os"
	"sync"
	"time"
)

// UpdateSiteRow | update a specific site
func UpdateSiteRow(siteID int, newSitename string, newUrl string, newDescription string) (bool, error) {

	var mutex = &sync.Mutex{}

	var cwd, _ = os.Getwd()

	// init database driver
	database, err := sql.Open("sqlite3", cwd+"/database/"+DB_NAME)
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	defer database.Close()

	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
		return false, err
	}

	// update record by ID
	currentTime := time.Now()
	updateStatement := fmt.Sprintf("UPDATE sites SET name = %v, url = %v, description = %v, date_edited=%v WHERE id = %v", newSitename, newUrl, newDescription, currentTime.Format("2006-01-02 15:04:05"), siteID)
	statementUpdateRecord, _ := database.Prepare(updateStatement)
	mutex.Lock()
	statementUpdateRecord.Exec()
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// UpdateRecordRow | update a specific record
func UpdateRecordRow(recordID int, newRecord string) (bool, error) {

	// mutex lock
	var mutex = &sync.Mutex{}

	var cwd, _ = os.Getwd()

	// init database driver
	database, err := sql.Open("sqlite3", cwd+"/database/"+DB_NAME)
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	defer database.Close()

	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
		return false, err
	}

	// update record by ID
	currentTime := time.Now()
	updateStatement := fmt.Sprintf("UPDATE records SET records_data = %v, date_edited=%v WHERE id = %v", newRecord, currentTime.Format("2006-01-02 15:04:05"), recordID)
	statementUpdateRecord, _ := database.Prepare(updateStatement)
	mutex.Lock()
	statementUpdateRecord.Exec()
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
