package database

import (
	"database/sql"
	"fmt"
	LOGS "go_svelte_lighthouse/logs"
	"os"
	"sync"
)

// RemoveSiteRow | delete a specific site and all records associated with it
func RemoveSiteRow(siteID int) (bool, error) {

	// mutex lock
	var mutex = &sync.Mutex{}

	var cwd, _ = os.Getwd()

	// init database driver
	database, err := sql.Open("sqlite3", cwd+"/database/"+DB_NAME)
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	defer database.Close()

	if err != nil || siteID < 1 {
		LOGS.ErrorLogger.Fatalln(err)
		return false, err
	}

	// delete site by ID
	deleteSiteStatement := fmt.Sprintf("DELETE FROM sites WHERE id = %v", siteID)
	statementDeleteSite, _ := database.Prepare(deleteSiteStatement)
	mutex.Lock()
	statementDeleteSite.Exec()
	mutex.Unlock()

	// delete all records previously associated with that site
	deleteRecordsStatement := fmt.Sprintf("DELETE FROM records WHERE site_id = %v", siteID)
	statementDeleteSiteRecords, _ := database.Prepare(deleteRecordsStatement)
	mutex.Lock()
	statementDeleteSiteRecords.Exec()
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}

}

// RemoveRecordRow | delete a specific record
func RemoveRecordRow(recordID int) (bool, error) {

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

	// delete record by ID
	deleteStatement := fmt.Sprintf("DELETE FROM records WHERE id = %v", recordID)
	statementDeleteRecord, _ := database.Prepare(deleteStatement)
	mutex.Lock()
	statementDeleteRecord.Exec()
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
