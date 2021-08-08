package rest

import (
	"errors"
	"sync"
	"time"

	CLI "go_svelte_lighthouse/cli"
	// CONFIG "go_svelte_lighthouse/config"
	DATABASE "go_svelte_lighthouse/database"
	LOGS "go_svelte_lighthouse/logs"
)

// structs for the urls found in the site manifest
type Sites struct {
	Sites []Site `json:sites`
}
type Site struct {
	Name string `json:name`
	URL  string `json:url`
}

// structs for the function return to handle errors and return the created report path
type FetchStatus struct {
	DidError   bool
	Error      error
	Message    string
	Duration   time.Duration
}
type FetchStatusSlice struct {
	Statuses []FetchStatus
	DidError bool
	Error    error
	FullDuration	float64
}

// getters for FetchStatus struct
func (f FetchStatus) ErrorStatus() bool {
	return f.DidError
}
func (f FetchStatus) GetError() error {
	return f.Error
}
func (f FetchStatus) GetMessage() string {
	return f.Message
}
func (f FetchStatus) GetDuration() time.Duration {
	return f.Duration
}

// getters for FetchStatusSlice struct
func (f FetchStatusSlice) ErrorStatus() bool {
	return f.DidError
}
func (f FetchStatusSlice) GetError() error {
	return f.Error
}

func RefetchWebsite(url string) map[string]FetchStatus {

	statusMap := make(map[string]FetchStatus)

	reportStart := time.Now()

	if len(url) < 1 {
		LOGS.DebugLogger.Println("Attempted to fetch a website without providing a URL")
		statusMap["nourl"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a website URL to fetch"),
			"Failure",
			time.Since(reportStart),
		}
	}

	ok, err := CLI.CreateReport(url, false)

	if err != nil {
		LOGS.ErrorLogger.Printf("Failure to fetch a report for %v", url)
		statusMap[url] = FetchStatus{
			!ok,
			err,
			"Failure to fetch a report for" + url,
			time.Since(reportStart),
		}
	} else {
		LOGS.InfoLogger.Printf("Successfully fetched a report for %v", url)
		statusMap[url] = FetchStatus{
			!ok,
			nil,
			"Success",
			time.Since(reportStart),
		}
	}

	return statusMap
}

func RefetchWebsites(cb func()) []map[string]FetchStatus {

	// waitgroup to handle goroutines concurrent dispatch
	var wg sync.WaitGroup

	// counting semaphore to limit goroutines to 20
	var semaphoreTokens = make(chan struct{}, 20)

	// fn returns a slice of maps
	var statusMapSlice []map[string]FetchStatus

	// track time taken to fetch all sites
	startTime := time.Now()

	// grab all available sites and their urls from manifest file (sites.json)
	var allUrls []map[string]string
	sitesSlice, err := DATABASE.ReturnSiteList();
	if err != nil {
		LOGS.ErrorLogger.Println("Failure to fetch a list of available sites")
	}

	allUrls = append(allUrls, sitesSlice...)

	// allUrls = CONFIG.GetAllRegisteredWebsites()

	if len(allUrls) > 0 {

		// loops to get the map interfaces for the sites
		for _, siteMap := range allUrls {
			siteURL := siteMap["url"]

			// make a status map to be returned
			statusMap := make(map[string]FetchStatus)

			// get semaphore token
			semaphoreTokens <- struct{}{}

			// add to waitgroup
			wg.Add(1)

			// start timer
			reportStart := time.Now()

			go func() {

				// release semaphore token
				defer func() { 
					<- semaphoreTokens
				}()

				defer wg.Done()

				ok, err := CLI.CreateReport(siteURL, false)

				if err != nil {
					LOGS.ErrorLogger.Printf("Failure to fetch a report for %v", siteURL)
					statusMap[siteURL] = FetchStatus{
						!ok,
						err,
						"Failure to fetch a report for" + siteURL,
						time.Since(reportStart),
					}
				} else {
					LOGS.InfoLogger.Printf("Successfully fetched a report for %v", siteURL)
					statusMap[siteURL] = FetchStatus{
						!ok,
						nil,
						"Success",
						time.Since(reportStart),
					}
				}

				statusMapSlice = append(statusMapSlice, statusMap)
				
			}()

		}

		wg.Wait()
		// send total time back to main
		finalTime := time.Since(startTime)
		timeStatusMap := make(map[string]FetchStatus)
		timeStatusMap["meta"] = FetchStatus{
			false,
			nil,
			"Success",
			finalTime,
		}
		statusMapSlice = append(statusMapSlice, timeStatusMap)

	}
	
	// trigger optional callback
	cb()
	return statusMapSlice
}
