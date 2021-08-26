package rest

import (
	"errors"
	"fmt"
	DATABASE "go_svelte_lighthouse/database"
	LOGS "go_svelte_lighthouse/logs"
	"time"
)

// UpdateWebsite function to trigger PATCH request to update site details
func UpdateWebsite(siteID int, newSitename string, newUrl string, newDescription string) map[string]FetchStatus {

	statusMap := make(map[string]FetchStatus)

	reportStart := time.Now()

	if siteID < 1 {
		LOGS.DebugLogger.Println("Attempted to update a record without providing a valid record ID")
		statusMap["noid"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a valid record ID"),
			"Failure",
			time.Since(reportStart),
		}
	}

	if len(newSitename) < 1 {
		LOGS.DebugLogger.Println("Attempted to update a record without providing a valid sitename string")
		statusMap["empty_record"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a valid updated sitename string"),
			"Failure",
			time.Since(reportStart),
		}
	}

	if len(newUrl) < 1 {
		LOGS.DebugLogger.Println("Attempted to update a record without providing a valid url string")
		statusMap["empty_record"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a valid updated url string"),
			"Failure",
			time.Since(reportStart),
		}
	}

	// perform the PATCH operation
	updateSuccess, err := DATABASE.UpdateSiteRow(siteID, newSitename, newUrl, newDescription)
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}

	if err != nil {
		LOGS.ErrorLogger.Printf("Failure to update the site with ID of %v", siteID)
		statusMap["updated"] = FetchStatus{
			!updateSuccess,
			err,
			fmt.Sprintf("Failure to update the site with ID of %v", siteID),
			time.Since(reportStart),
		}
	} else {
		LOGS.InfoLogger.Printf("Successfully updated the site with ID of %v", siteID)
		statusMap["updated"] = FetchStatus{
			!updateSuccess,
			nil,
			"Success",
			time.Since(reportStart),
		}
	}

	return statusMap
}

// UpdateRecord function to trigger PATCH request to update singular record
func UpdateRecord(recordID int, newRecord string) map[string]FetchStatus {

	statusMap := make(map[string]FetchStatus)

	reportStart := time.Now()

	if recordID < 1 {
		LOGS.DebugLogger.Println("Attempted to update a record without providing a valid record ID")
		statusMap["noid"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a valid record ID"),
			"Failure",
			time.Since(reportStart),
		}
	}

	if len(newRecord) < 1 {
		LOGS.DebugLogger.Println("Attempted to update a record without providing a valid record string")
		statusMap["empty_record"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a valid record string"),
			"Failure",
			time.Since(reportStart),
		}
	}

	// perform the PATCH operation
	updateSuccess, err := DATABASE.UpdateRecordRow(recordID, newRecord)
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}

	if err != nil {
		LOGS.ErrorLogger.Printf("Failure to update the record with ID of %v", recordID)
		statusMap["updated"] = FetchStatus{
			!updateSuccess,
			err,
			fmt.Sprintf("Failure to update the record with ID of %v", recordID),
			time.Since(reportStart),
		}
	} else {
		LOGS.InfoLogger.Printf("Successfully updated the record with ID of %v", recordID)
		statusMap["updated"] = FetchStatus{
			!updateSuccess,
			nil,
			"Success",
			time.Since(reportStart),
		}
	}

	return statusMap
}
