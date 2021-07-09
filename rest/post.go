package rest

import (
	"errors"
	"sync"
	"time"

	CLI "go_svelte_lighthouse/cli"
	CONFIG "go_svelte_lighthouse/config"
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
	ReportPath string
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
func (f FetchStatus) GetReportPath() string {
	return f.ReportPath
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
		LOGS.WarningLogger.Println("Please provide a website URL to fetch")
		statusMap["nourl"] = FetchStatus{
			true,
			errors.New("Please provide a website URL to fetch"),
			"Failure",
			"",
			time.Since(reportStart),
		}
	}

	output, err := CLI.CreateReport(url, false)
	if err != nil {
		LOGS.WarningLogger.Printf("Failure to fetch a report for %v", url)
		statusMap[url] = FetchStatus{
			true,
			err,
			"Failure to fetch a report for" + url,
			"",
			time.Since(reportStart),
		}
	} else {
		statusMap[url] = FetchStatus{
			false,
			nil,
			"Success",
			output,
			time.Since(reportStart),
		}
	}

	return statusMap
}

func RefetchWebsites() []map[string]FetchStatus {

	var wg sync.WaitGroup

	// fn returns a slice of maps
	var statusMapSlice []map[string]FetchStatus

	// track time taken to fetch all sites
	startTime := time.Now()

	// grab all urls from manifest file (sites.json)
	allUrls := CONFIG.GetAllRegisteredWebsites()

	if len(allUrls) > 0 {

		// loops to get the map interfaces for the sites
		for _, site := range allUrls {
			sitesSlice := site.([]interface{})

			for _, siteMap := range sitesSlice {
				siteURLMap := siteMap.(map[string]interface{})
				siteURL := siteURLMap["url"].(string)

				// make a status map to be returned
				statusMap := make(map[string]FetchStatus)

				wg.Add(1)

				reportStart := time.Now()

				go func() {
					defer wg.Done()
					output, err := CLI.CreateReport(siteURL, false)

					if err != nil {
						LOGS.InfoLogger.Printf("Failure to fetch a report for %v", siteURL)
						statusMap[siteURL] = FetchStatus{
							true,
							err,
							"Failure to fetch a report for" + siteURL,
							"",
							time.Since(reportStart),
						}
					} else {
						statusMap[siteURL] = FetchStatus{
							false,
							nil,
							"Success",
							output,
							time.Since(reportStart),
						}
					}

					statusMapSlice = append(statusMapSlice, statusMap)
					
				}()

			}

		}

		wg.Wait()
		// send total time back to main
		finalTime := time.Since(startTime)
		timeStatusMap := make(map[string]FetchStatus)
		timeStatusMap["meta"] = FetchStatus{
			false,
			nil,
			"Success",
			"",
			finalTime,
		}
		statusMapSlice = append(statusMapSlice, timeStatusMap)

	}

	return statusMapSlice
}
