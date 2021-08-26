package database

import (
	"database/sql"
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
	currentTimestamp := currentTime.Format("2006-01-02 15:04:05")
	statementUpdateRecord, _ := database.Prepare("UPDATE sites SET name=?, url=?, description=?, date_edited=? WHERE id=?")
	mutex.Lock()
	statementUpdateRecord.Exec(newSitename, newUrl, newDescription, currentTimestamp, siteID)
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
	currentTimestamp := currentTime.Format("2006-01-02 15:04:05")
	statementUpdateRecord, _ := database.Prepare("UPDATE records SET records_data=?, date_edited=? WHERE id=?")
	mutex.Lock()
	statementUpdateRecord.Exec(newRecord, currentTimestamp, recordID)
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
