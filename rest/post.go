package rest

import (
	"errors"
	"fmt"
	CLI "go_svelte_lighthouse/cli"
	DATABASE "go_svelte_lighthouse/database"
	LOGS "go_svelte_lighthouse/logs"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

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

// GetWebsiteStatistics main function to trigger POST request to get site stats for specific named site
func GetWebsiteStatistics(urlToFetch string) map[string]FetchStatus {

	statusMap := make(map[string]FetchStatus)

	reportStart := time.Now()

	if len(urlToFetch) < 1 {
		LOGS.DebugLogger.Println("Attempted to fetch a website without providing a URL")
		statusMap["nourl"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a website URL to fetch"),
			"Failure",
			time.Since(reportStart),
		}
	}

	// strip protocol, args etc. from the url to make nice sitename
	var sitename string
	parsedURL, err := url.Parse(urlToFetch)
	if err != nil {
		log.Fatal(err)
	}
	// just get the host part of the url
	sitename = parsedURL.Host
	// also remove www.
	sitename = strings.Replace(sitename, "www.", "", -1)

	// set a combined status variable (tracks CLI success plus DB write success)
	var combinedSuccessStatus bool

	// get site stats as a json string from the CLI tool
	ok, jsonResultString, err := CLI.CreateReport(urlToFetch, false)

	// add result to the database
	writeSuccess, writeError := DATABASE.InsertDatabaseRowRecord(sitename, urlToFetch, jsonResultString)
	if writeError != nil {
		LOGS.ErrorLogger.Fatalln(writeError)
	}

	// update the combined status variable for use in returned statusMap
	if ok && writeSuccess {
		combinedSuccessStatus = true
	}

	fmt.Println("combinedSuccessStatus")
	fmt.Println(combinedSuccessStatus)

	if err != nil {
		LOGS.ErrorLogger.Printf("Failure to fetch a report for %v", urlToFetch)
		statusMap[urlToFetch] = FetchStatus{
			!combinedSuccessStatus,
			err,
			"Failure to fetch a report for" + urlToFetch,
			time.Since(reportStart),
		}
	} else {
		LOGS.InfoLogger.Printf("Successfully fetched a report for %v", urlToFetch)
		statusMap[urlToFetch] = FetchStatus{
			!combinedSuccessStatus,
			nil,
			"Success",
			time.Since(reportStart),
		}
	}

	return statusMap
}

// GetAllWebsiteStatistics main function to trigger POST request to get site stats for all available sites
func GetAllWebsiteStatistics(cb func()) []map[string]FetchStatus {

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
