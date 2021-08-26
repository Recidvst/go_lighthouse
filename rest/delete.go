package rest

import (
	"errors"
	"fmt"
	DATABASE "go_svelte_lighthouse/database"
	LOGS "go_svelte_lighthouse/logs"
	"time"
)

// DeleteWebsite function to trigger DELETE request to delete site and all associated records
func DeleteWebsite(siteID int) map[string]FetchStatus {

	statusMap := make(map[string]FetchStatus)

	reportStart := time.Now()

	if siteID < 1 {
		LOGS.DebugLogger.Println("Attempted to delete a website without providing a valid site ID")
		statusMap["noid"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a valid site ID"),
			"Failure",
			time.Since(reportStart),
		}
	}

	// perform the DELETE operation
	deleteSuccess, err := DATABASE.RemoveSiteRow(siteID)
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}

	if err != nil {
		LOGS.ErrorLogger.Printf("Failure to delete the site with ID of %v", siteID)
		statusMap["deleted"] = FetchStatus{
			!deleteSuccess,
			err,
			fmt.Sprintf("Failure to delete the site with ID of %v", siteID),
			time.Since(reportStart),
		}
	} else {
		LOGS.InfoLogger.Printf("Successfully deleted the site with ID of %v", siteID)
		statusMap["deleted"] = FetchStatus{
			!deleteSuccess,
			nil,
			"Success",
			time.Since(reportStart),
		}
	}

	return statusMap
}

// DeleteRecord function to trigger DELETE request to delete a singular record
func DeleteRecord(recordID int) map[string]FetchStatus {

	statusMap := make(map[string]FetchStatus)

	reportStart := time.Now()

	if recordID < 1 {
		LOGS.DebugLogger.Println("Attempted to delete a record without providing a valid record ID")
		statusMap["noid"] = FetchStatus{
			true, // true = did error
			errors.New("please provide a valid record ID"),
			"Failure",
			time.Since(reportStart),
		}
	}

	// perform the DELETE operation
	deleteSuccess, err := DATABASE.RemoveRecordRow(recordID)
	if err != nil {
		LOGS.ErrorLogger.Fatalln(err)
	}

	if err != nil {
		LOGS.ErrorLogger.Printf("Failure to delete the record with ID of %v", recordID)
		statusMap["deleted"] = FetchStatus{
			!deleteSuccess,
			err,
			fmt.Sprintf("Failure to delete the record with ID of %v", recordID),
			time.Since(reportStart),
		}
	} else {
		LOGS.InfoLogger.Printf("Successfully deleted the record with ID of %v", recordID)
		statusMap["deleted"] = FetchStatus{
			!deleteSuccess,
			nil,
			"Success",
			time.Since(reportStart),
		}
	}

	return statusMap
}
