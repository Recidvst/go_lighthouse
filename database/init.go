package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	LOGS "go_svelte_lighthouse/logs"
	"os"
	"sync"
)

const DB_NAME string = "sqlite__siteresults.db"

func createDB(name string) (bool, error) {

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

	// prepare statements to create tables
	statementTableSites, err := database.Prepare("CREATE TABLE IF NOT EXISTS sites (id INTEGER PRIMARY KEY UNIQUE NOT NULL, name STRING UNIQUE NOT NULL, url STRING NOT NULL, description STRING);")
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	mutex.Lock()
	statementTableSites.Exec()
	mutex.Unlock()

	statementTableRecords, err := database.Prepare("CREATE TABLE IF NOT EXISTS records (id INTEGER PRIMARY KEY UNIQUE NOT NULL, site_id INTEGER NOT NULL CONSTRAINT cx_results_site_id REFERENCES sites (id) ON DELETE CASCADE ON UPDATE CASCADE, records_data STRING NOT NULL, date_fetched DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP), date_edited  DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP));")
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	mutex.Lock()
	statementTableRecords.Exec()
	mutex.Unlock()

	// add indexes to each table for faster querying
	statementIndexSites, err := database.Prepare("CREATE INDEX IF NOT EXISTS idx_sites_name ON sites (name ASC);")
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	mutex.Lock()
	statementIndexSites.Exec()
	mutex.Unlock()

	statementIndexRecords, err := database.Prepare("CREATE INDEX IF NOT EXISTS idx_results_site_id ON records (site_id ASC);")
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	mutex.Lock()
	statementIndexRecords.Exec()
	mutex.Unlock()

	// add an example row to each table for testing purposes (if empty)
	// no need to add a site, because the InsertDatabaseRowRecord will do this if no matching site is found
	checkEmptyRecords, err := database.Query("SELECT count(*) FROM records")
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	var countRecords int
	// check if rows present or not
	for checkEmptyRecords.Next() {
		_ = checkEmptyRecords.Scan(&countRecords)
	}
	checkEmptyRecords.Close()

	// add rows
	if countRecords < 1 {
		InsertDatabaseRowRecord("chris-snowden.me", "https://www.chris-snowden.me/", "{'name':'Sarah', 'age':25, 'car':'Ferrari}")
	}

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// init fn to create db tables
func init() {
	createDB(DB_NAME)
}
