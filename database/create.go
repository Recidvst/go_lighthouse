package database

import (
	"database/sql"
	LOGS "go_svelte_lighthouse/logs"
	"os"
	"sync"
	"time"
)

// InsertDatabaseRowSite insert a new site row into the DB
func InsertDatabaseRowSite(sitename string, description string, url string) (bool, error) {

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

	// add new site
	statementInsertSite, _ := database.Prepare("INSERT INTO sites (id, name, description, url) VALUES (?, ?, ?, ?);")
	mutex.Lock()
	statementInsertSite.Exec(nil, sitename, description, url)
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// InsertDatabaseRowRecord insert a new record row into the DB
func InsertDatabaseRowRecord(sitename string, url string, record string) (bool, error) {

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

	// init two empty variables to hold the two possible site IDs
	// empty value of int is 0
	var returnedSiteId int
	var returnedNewSiteId int

	// check that the site exists
	siteQueryString := "SELECT id FROM sites WHERE name LIKE " + "'" + sitename + "'"
	// make the query
	doesSiteExist, err := database.Query(siteQueryString)
	if err != nil {
		LOGS.DebugLogger.Println(err)
	} else {
		for doesSiteExist.Next() {
			_ = doesSiteExist.Scan(&returnedSiteId)
		}
		doesSiteExist.Close()
	}

	// if site row not found then add it
	if returnedSiteId < 1 {
		InsertDatabaseRowSite(sitename, "Default description for site: "+sitename, url)
		// then, get the id for that site (we need this to insert a record row)
		siteIdQueryString := "SELECT id FROM sites WHERE name LIKE " + "'" + sitename + "'"
		// make the query
		siteIdQuery, err := database.Query(siteIdQueryString)
		if err != nil {
			LOGS.DebugLogger.Println(err)
		} else {
			for siteIdQuery.Next() {
				_ = siteIdQuery.Scan(&returnedNewSiteId)
			}
			siteIdQuery.Close()
		}
	}

	// insert data row if/once the site exists
	statementInsertRecord, err := database.Prepare("INSERT INTO records (id, site_id, records_data, date_fetched) VALUES (?, ?, ?, ?);")
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	mutex.Lock()
	currentTime := time.Now()
	if returnedSiteId < 1 {
		statementInsertRecord.Exec(nil, returnedNewSiteId, record, currentTime.Format("2006-01-02 15:04:05"))
	} else {
		statementInsertRecord.Exec(nil, returnedSiteId, record, currentTime.Format("2006-01-02 15:04:05"))
	}
	mutex.Unlock()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
