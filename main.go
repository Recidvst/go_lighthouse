package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"

	CONFIG "go_svelte_lighthouse/config"
	LOGS "go_svelte_lighthouse/logs"
	REST "go_svelte_lighthouse/rest"
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

var EnvironmentType = CONFIG.GetEnvByKey("ENVIRONMENT")
var RegisteredWebsites = CONFIG.GetAllRegisteredWebsites()

func main() {
	fmt.Println("API up")
	fmt.Println(RegisteredWebsites)

	// Init router
	r := mux.NewRouter()

	// GET | root endpoint, test if API up
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if EnvironmentType != "production" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET", "OPTIONS")

	// POST | refetch a specific website
	r.HandleFunc("/website", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POST received")

		// set headers
		w.Header().Set("Content-Type", "application/json")
		if EnvironmentType != "production" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// get url param
		requestedUrl := r.FormValue("url")

		var status bool = false
		var statusErr error
		var statusPath string

		// fetch website report
		if len(requestedUrl) > 0 {
			statusMap := REST.RefetchWebsite(requestedUrl)
			if !statusMap[requestedUrl].ErrorStatus() {
				status = true
				statusPath = statusMap[requestedUrl].GetReportPath()
			} else {
				statusErr = statusMap[requestedUrl].GetError()
			}
		}

		// send a response depending on error or success
		if !status {
			json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
		} else {
			json.NewEncoder(w).Encode(map[string]string{"status": "Success", "outputPath": statusPath})
		}

	}).Methods("POST", "OPTIONS")

	// POST | refetch all tracked websites
	r.HandleFunc("/websites", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("POST received")

		// set headers
		w.Header().Set("Content-Type", "application/json")
		if EnvironmentType != "production" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		var status bool = false
		var statusErr error
		var statusPath string

		// fetch website report
		statusMap := REST.RefetchWebsites()
		if !statusMap[requestedUrl].ErrorStatus() {
			status = true
			statusPath = statusMap[requestedUrl].GetReportPath()
		} else {
			statusErr = statusMap[requestedUrl].GetError()
		}

		// send a response depending on error or success
		if !status {
			json.NewEncoder(w).Encode(map[string]string{"status": "Error", "error": statusErr.Error()})
		} else {
			json.NewEncoder(w).Encode(map[string]string{"status": "Success"})
		}

	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/website", func(w http.ResponseWriter, r *http.Request) {
		if EnvironmentType != "production" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		// get request query
		url := r.URL.Query()
		fmt.Printf("%+v\n", url)
		// get requested site name
		var urlToTarget string
		urlToTarget = url.Get("url")
		fmt.Println(urlToTarget)
		// return if no website passed
		if len(urlToTarget) < 1 {
			json.NewEncoder(w).Encode(map[string]string{"status": "something went wrong"})
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "TODO"})

		// fmt.Sprintf("trigger lighthouse for %d", website)
	}).Methods("GET", "OPTIONS")

	// Start server
	log.Fatalln(http.ListenAndServe(":9999", r))
}
