package database

import (
	"database/sql"
	LOGS "go_svelte_lighthouse/logs"
	"os"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func createDB(name string) (bool, error) {
	
	// mutex lock
	var mutex = &sync.Mutex{}

	var cwd, _ = os.Getwd()

	// init database driver
	database, err := sql.Open("sqlite3", cwd + "/database/" + name)
	defer database.Close()

	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
		return false, err
	}

	// prepare statements to create tables
	statementTableSites, err := database.Prepare("CREATE TABLE IF NOT EXISTS sites (id INTEGER PRIMARY KEY UNIQUE NOT NULL, name STRING NOT NULL);")
	mutex.Lock()
	statementTableSites.Exec();
	mutex.Unlock()

	statementTableRecords, err := database.Prepare("CREATE TABLE IF NOT EXISTS records (id INTEGER PRIMARY KEY UNIQUE NOT NULL, site_id INTEGER NOT NULL CONSTRAINT cx_results_site_id REFERENCES sites (id) ON DELETE NO ACTION ON UPDATE CASCADE, results_data STRING NOT NULL, date_fetched DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP), date_edited  DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP));")
	mutex.Lock()
	statementTableRecords.Exec();
	mutex.Unlock()

	// add indexes to each table for faster querying
	statementIndexSites, err := database.Prepare("CREATE INDEX IF NOT EXISTS idx_sites_name ON sites (name ASC);")
	mutex.Lock()
	statementIndexSites.Exec();
	mutex.Unlock()

	statementIndexRecords, err := database.Prepare("CREATE INDEX IF NOT EXISTS idx_records_site_id ON records (site_id ASC);")
	mutex.Lock()
	statementIndexRecords.Exec();
	mutex.Unlock()

	// add an example row to each table for testing purposes
	statementInsertTestRow, err := database.Prepare("INSERT INTO sites (name) VALUES (?);")
	mutex.Lock()
	statementInsertTestRow.Exec("example.com");
	mutex.Unlock()
	statementInsertTestRow2, err := database.Prepare("INSERT INTO records (site_id, results_data, date_fetched) VALUES (?, ?, ?);")
	mutex.Lock()
	statementInsertTestRow2.Exec(1, "{'name':'Sarah', 'age':25, 'car':'Ferrari}", "2021-07-20 09:15:03");
	mutex.Unlock()


	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func insertDatabaseRow(sitename string, record string) (bool, error) {

	// mutex lock
	var mutex = &sync.Mutex{}

	var cwd, _ = os.Getwd()

	// init database driver
	database, err := sql.Open("sqlite3", cwd + "/database/" + name)
	defer database.Close()

	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
		return false, err
	}

	// TODO - query for if site exists in table - if so, return site_id for use in record query. Otherwise, create fresh
	var siteid int

	// insert new site
	statementInsertSite, err := database.Prepare("INSERT INTO sites (name) VALUES (?);")
	mutex.Lock()
	statementInsertSite.Exec(sitename);
	mutex.Unlock()

	// insert data row
	statementInsertRecord, err := database.Prepare("INSERT INTO records (site_id, results_data, date_fetched) VALUES (?, ?, ?);")
	mutex.Lock()
	currentTime := time.Now()
	statementInsertRecord.Exec(siteid, record, currentTime.Format("2006-01-02 15:04:05"));
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// func updateRecordRow(record_id int, newrecord) (bool, error) {}

// func updateSiteRow(site_id int, newname string) (bool, error) {}

// func removeRecordRow(record_id int) (bool, error) {}

// func removeSiteRow(site_id int) (bool, error) {}

// init fn to create db tables
func init() {
	createDB("sqlite__siteresults.db")
}
