package database

import (
	"database/sql"
	LOGS "go_svelte_lighthouse/logs"
	"os"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// return all available sites
func ReturnSiteList() ([]map[string]string, error) {

	var sites []map[string]string

	// mutex lock
	var mutex = &sync.Mutex{}

	var cwd, _ = os.Getwd()

	// init database driver
	database, err := sql.Open("sqlite3", cwd + "/database/" + DB_NAME)
	defer database.Close()

	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
		return sites, err
	}

	// prepare query
	mutex.Lock()
	rows, _ := database.Query("SELECT id, name, url FROM sites;")
	mutex.Unlock()
	var id int
	var name string
	var url string

	// pull query results into vars and add to return slice
	for rows.Next() {
		var siteDetails = make(map[string]string)

		rows.Scan(&id, &name, &url)
		
		siteDetails["id"] = strconv.Itoa(id)
		siteDetails["name"] = name
		siteDetails["url"] = url

		sites = append(sites, siteDetails)
	}

	return sites, err
}
