package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	// "math/rand"
	// "strconv"
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

func main() {
	fmt.Println("API up")

	// Init router
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("API hit")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	r.HandleFunc("/website", func(w http.ResponseWriter, r *http.Request) {
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
	}).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":9999", r))
}
