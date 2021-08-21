package rest

import (
	"errors"
	"fmt"
	"sync"
	"time"

	CLI "go_svelte_lighthouse/cli"
	DATABASE "go_svelte_lighthouse/database"
	LOGS "go_svelte_lighthouse/logs"
)

// sites struct to contain slice of site structs
//type sites struct {
//	Sites []site `json:"sites"`
//}
//
//// site struct for the urls found in the site manifest
//type site struct {
//	Name string `json:"name"`
//	URL  string `json:"url"`
//}

// FetchStatus structs for the function return to handle errors and return the created report path
type FetchStatus struct {
	DidError bool
	Error    error
	Message  string
	Duration time.Duration
}

// ErrorStatus getters for FetchStatus struct
func (f FetchStatus) ErrorStatus() bool {
	return f.DidError
}

// GetError getter method
func (f FetchStatus) GetError() error {
	return f.Error
}

// GetMessage getter method
func (f FetchStatus) GetMessage() string {
	return f.Message
}

// GetDuration getter method
func (f FetchStatus) GetDuration() time.Duration {
	return f.Duration
}

// RefetchWebsite main function to trigger POST request to get site stats for specific named site
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

	ok, jsonResultString, err := CLI.CreateReport(url, false)
	fmt.Println(jsonResultString)

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

// RefetchWebsites main function to trigger POST request to get site stats for all available sites
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
	sitesSlice, err := DATABASE.ReturnSiteList()
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
					<-semaphoreTokens
				}()

				defer wg.Done()

				ok, jsonResultString, err := CLI.CreateReport(siteURL, false)
				fmt.Println(jsonResultString)

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
