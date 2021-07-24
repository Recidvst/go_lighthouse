package main

import (
	"database/sql"
	LOGS "go_svelte_lighthouse/logs"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func createDB(name string) (bool, error) {

	// init database driver
	var cwd, _ = os.Getwd()

	database, err := sql.Open("sqlite3", cwd + "/database/" + name)

	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
		return false, err
	}

	// prepare statements to create tables
	statementTableSites, err := database.Prepare("CREATE TABLE IF NOT EXISTS sites (id INTEGER PRIMARY KEY UNIQUE NOT NULL, name STRING NOT NULL);")
	statementTableRecords, err := database.Prepare("CREATE TABLE IF NOT EXISTS records (id INTEGER PRIMARY KEY UNIQUE NOT NULL, site_id INTEGER NOT NULL CONSTRAINT cx_results_site_id REFERENCES sites (id) ON DELETE NO ACTION ON UPDATE CASCADE, results_data STRING NOT NULL, date_fetched DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP), date_edited  DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP));")
	

	// add indexes to each table for faster querying
	statementIndexSites, err := database.Prepare("CREATE INDEX IF NOT EXISTS idx_sites_name ON sites (name ASC);")
	statementIndexRecords, err := database.Prepare("CREATE INDEX IF NOT EXISTS idx_records_site_id ON records (site_id ASC);")

	// execute statements
	statementTableSites.Exec();
	statementTableRecords.Exec();
	statementIndexSites.Exec();
	statementIndexRecords.Exec();

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func main() {
	createDB("sqlite__siteresults.db")
}
