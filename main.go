package main

import (
	"errors"
	"fmt"
	LOGS "go_svelte_lighthouse/logs"
	"log"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/gorilla/mux"

	CONFIG "go_svelte_lighthouse/config"
	"go_svelte_lighthouse/rest"
	REST "go_svelte_lighthouse/rest"

	// import cron although we don't fire it manually
	_ "go_svelte_lighthouse/cron"
)

// Website struct (holds the website)
type Website struct {
	ID       string    `json:"id"`
	Sitename string    `json:"sitename"`
	URL      string    `json:"url"`
	Results  []*Result `json:"website"`
}

// Results struct (array(slice) of results)
type Results struct {
	Website     *Website  `json:"website"`
	ResultItems []*Result `json:"resultitems"`
}

// Result struct (array(slice) of key value maps)
type Result struct {
	ID              string               `json:"id"`
	Datetime        string               `json:"datetime"`
	ResultContainer *Results             `json:"resultcontainer"`
	Contents        map[string]ResultMap `json:"contents"`
}

// ResultValue struct (map containing key value speed data)
type ResultMap struct {
	ResultParent *Result `json:"resultparent"`
	Key          string  `json:"key"`
	Value        string  `json:"value"`
}

var environmentType = CONFIG.GetEnvByKey("ENVIRONMENT")

// POST | refetch a specific website
func fetchSingleWebsite(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetch single website")

	// set headers
	w.Header().Set("Content-Type", "application/json")
	if environmentType != "production" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// get url param (in body)
	requestedUrl := r.FormValue("url")

	var status = false
	var statusErr error
	var duration int64

	// fetch website report
	if len(requestedUrl) > 0 {
		statusMap := REST.GetWebsiteStatistics(requestedUrl)

		duration = statusMap[requestedUrl].GetDuration().Milliseconds()

		if !statusMap[requestedUrl].ErrorStatus() {
			status = true
		} else {
			statusErr = statusMap[requestedUrl].GetError()
		}
	}

	// send a response depending on error or success
	if !status {
		json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "Success", "time_to_generate": strconv.FormatInt(duration, 10) + "ms"})
	}
}

// POST | refetch all tracked websites
func fetchAllWebsites(w http.ResponseWriter, r *http.Request) {

	// set headers
	w.Header().Set("Content-Type", "application/json")
	if environmentType != "production" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var status = false
	var statusErr error

	// fetch website report
	var statusMapsCollection = REST.GetAllWebsiteStatistics(func() {})

	if len(statusMapsCollection) < 1 {
		statusErr = errors.New("failed to refetch any websites")
	} else {
		status = true
	}

	// get meta map from collection
	var metaMap rest.FetchStatus
	var timeToGenerate int64
	for m := 0; m < len(statusMapsCollection); m++ {
		if val, ok := statusMapsCollection[m]["meta"]; ok {
			metaMap = val
			timeToGenerate = metaMap.GetDuration().Milliseconds()
		}
	}

	// send a response depending on error or success
	if !status {
		json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"status":                      "Success",
			"number_of_reports_generated": strconv.Itoa(len(statusMapsCollection) - 1), // minus 1 because we need to account for the extra 'timeStatusMap' meta entry
			"time_to_generate":            strconv.FormatInt(timeToGenerate, 10) + "ms",
		})
	}
}

// DELETE | delete all records for specific website
func deleteSingleWebsite(w http.ResponseWriter, r *http.Request) {
	// get site_id param (in body)
	var inwardID = r.FormValue("site_id")

	// string to int
	onwardID, err := strconv.ParseInt(inwardID, 10, 8)

	if err != nil {
		LOGS.ErrorLogger.Println(err)
	}

	if onwardID < 1 {
		LOGS.DebugLogger.Println("Request to deleteSingleWebsite missing valid site ID")
	}

	var status = false
	var statusErr error
	var duration int64

	// delete site
	statusMap := REST.DeleteWebsite(int(onwardID))

	duration = statusMap["deleted"].GetDuration().Milliseconds()

	if !statusMap["deleted"].ErrorStatus() {
		status = true
	} else {
		statusErr = statusMap["deleted"].GetError()
	}

	// send a response depending on error or success
	if !status {
		json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "Success", "time_to_delete": strconv.FormatInt(duration, 10) + "ms"})
	}
}

// DELETE | delete a single specified record
func deleteSingleRecord(w http.ResponseWriter, r *http.Request) {
	// get record_id param (in body)
	var inwardID = r.FormValue("record_id")

	// string to int
	onwardID, err := strconv.ParseInt(inwardID, 10, 8)

	if err != nil {
		LOGS.ErrorLogger.Println(err)
	}

	if onwardID < 1 {
		LOGS.DebugLogger.Println("Request to deleteSingleRecord missing valid record ID")
	}

	var status = false
	var statusErr error
	var duration int64

	// delete site
	statusMap := REST.DeleteRecord(int(onwardID))

	duration = statusMap["deleted"].GetDuration().Milliseconds()

	if !statusMap["deleted"].ErrorStatus() {
		status = true
	} else {
		statusErr = statusMap["deleted"].GetError()
	}

	// send a response depending on error or success
	if !status {
		json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "Success", "time_to_delete": strconv.FormatInt(duration, 10) + "ms"})
	}
}

// PUT | update a single specified website
func updateSingleWebsite(w http.ResponseWriter, r *http.Request) {
	// get record_id param (in body)
	var inwardID = r.FormValue("record_id")
	var inwardOnwardSitename = r.FormValue("sitename")
	var inwardOnwardUrl = r.FormValue("url")
	var inwardOnwardDescription = r.FormValue("description")

	// string to int
	onwardID, err := strconv.ParseInt(inwardID, 10, 8)

	if err != nil {
		LOGS.ErrorLogger.Println(err)
	}

	var status = false
	var statusErr error
	var duration int64

	if onwardID < 1 {
		LOGS.DebugLogger.Println("Request to updateSingleWebsite missing valid record ID")
	}
	if len(inwardOnwardSitename) < 1 {
		LOGS.DebugLogger.Println("Request to updateSingleWebsite missing updated/confirmed site name")
	}
	if len(inwardOnwardUrl) < 1 {
		LOGS.DebugLogger.Println("Request to updateSingleWebsite missing updated/confirmed site url")
	}
	// description can be an empty string if required so no check here (tbc?)

	// update site
	statusMap := REST.UpdateWebsite(int(onwardID), inwardOnwardSitename, inwardOnwardUrl, inwardOnwardDescription)

	duration = statusMap["updated"].GetDuration().Milliseconds()

	if !statusMap["updated"].ErrorStatus() {
		status = true
	} else {
		statusErr = statusMap["updated"].GetError()
	}

	// send a response depending on error or success
	if !status {
		json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "Success", "time_to_update": strconv.FormatInt(duration, 10) + "ms"})
	}
}

// PATCH | update a single specified record
func updateSingleRecord(w http.ResponseWriter, r *http.Request) {
	// get record_id param (in body)
	var inwardID = r.FormValue("record_id")
	var inwardOnwardRecord = r.FormValue("record")

	// string to int
	onwardID, err := strconv.ParseInt(inwardID, 10, 8)

	if err != nil {
		LOGS.ErrorLogger.Println(err)
	}

	var status = false
	var statusErr error
	var duration int64

	if onwardID < 1 {
		LOGS.DebugLogger.Println("Request to updateSingleRecord missing valid record ID")
	}
	if len(inwardOnwardRecord) < 1 {
		LOGS.DebugLogger.Println("Request to updateSingleRecord missing new record string")
	}

	// update record
	statusMap := REST.UpdateRecord(int(onwardID), inwardOnwardRecord)

	duration = statusMap["updated"].GetDuration().Milliseconds()

	if !statusMap["updated"].ErrorStatus() {
		status = true
	} else {
		statusErr = statusMap["updated"].GetError()
	}

	// send a response depending on error or success
	if !status {
		json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "Success", "time_to_update": strconv.FormatInt(duration, 10) + "ms"})
	}
}

// GET | get details for specific website
func getSingleWebsite(w http.ResponseWriter, r *http.Request) {
	if environmentType != "production" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	// get request query
	url := r.URL.Query()
	// get requested site name
	var urlToTarget string
	urlToTarget = url.Get("url")

	// return if no website passed
	if len(urlToTarget) < 1 {
		json.NewEncoder(w).Encode(map[string]string{"status": "something went wrong"})
	}
	// TODO
	json.NewEncoder(w).Encode(map[string]string{"status": "TODO"})
}

// GET | get details for all tracked websites
func getAllWebsites(w http.ResponseWriter, r *http.Request) {
	if environmentType != "production" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	// TODO
}

// GET | view details for specific website in a browser, via html template
func viewSingleWebsite(w http.ResponseWriter, r *http.Request) {
	if environmentType != "production" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	// TODO
}

func main() {
	fmt.Println("API up")

	// Init router
	r := mux.NewRouter()

	// GET | root endpoint, test if API up
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if environmentType != "production" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET", "OPTIONS")

	// ENDPOINTS -------------------------------------------------------------- //

	// POST | refetch a specific website
	r.HandleFunc("/website", fetchSingleWebsite).Methods("POST", "OPTIONS")

	// POST | refetch all tracked websites
	r.HandleFunc("/websites", fetchAllWebsites).Methods("POST", "OPTIONS")

	// GET | get details for specific website
	r.HandleFunc("/website", getSingleWebsite).Methods("GET", "OPTIONS")

	// GET | get details for all tracked websites
	r.HandleFunc("/websites", getAllWebsites).Methods("GET", "OPTIONS")

	// GET | view details for specific website in a browser, via html template
	r.HandleFunc("/website", viewSingleWebsite).Methods("GET", "OPTIONS")

	// DELETE | delete site record and associated records
	r.HandleFunc("/website", deleteSingleWebsite).Methods("DELETE", "OPTIONS")

	// DELETE | delete specific record
	r.HandleFunc("/record", deleteSingleRecord).Methods("DELETE", "OPTIONS")

	// PUT | update specific site details
	r.HandleFunc("/website", updateSingleWebsite).Methods("PUT", "OPTIONS")

	// PUT | update specific record
	r.HandleFunc("/record", updateSingleRecord).Methods("PUT", "OPTIONS")

	// Start server
	log.Fatalln(http.ListenAndServe(":9999", r))
}
