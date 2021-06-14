package config

import (
	"encoding/json"
	LOGS "go_svelte_lighthouse/logs"
	"io/ioutil"
	"os"
)

// structs to parse the json data
type SiteEntry struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type SiteEntries struct {
	Sites []*SiteEntry `json:"sites"`
}

// handle global config vars
func GetAllRegisteredWebsites() map[string]interface{} {

	// get path
	var cwd, err = os.Getwd()
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}

	// open json file
	jsonFile, err := os.Open(cwd + "/config/sites.json")
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read the file
	jsonBytes, _ := ioutil.ReadAll(jsonFile)

	// initiate map to be returned
	var siteEntries map[string]interface{}

	// unmarshal json into map
	if err := json.Unmarshal(jsonBytes, &siteEntries); err != nil {
		LOGS.ErrorLogger.Fatalf("Failed to unmarshal json file, error: %v", err)
	}

	return siteEntries
}
