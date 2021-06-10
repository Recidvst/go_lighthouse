package rest

import (
	"errors"
	"fmt"

	CLI "go_svelte_lighthouse/cli"
)

// struct for the function return to handle errors and return the created report path
type FetchStatus struct {
	DidError   bool
	Error      error
	Message    string
	ReportPath string
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

func RefetchWebsite(url string) map[string]FetchStatus {

	statusMap := make(map[string]FetchStatus)

	if len(url) < 1 {
		statusMap["nourl"] = FetchStatus{
			true,
			errors.New("Please provide a website URL to fetch"),
			"Failure",
			"",
		}
	}

	output, err := CLI.CreateReport(url)
	if err != nil {
		statusMap[url] = FetchStatus{
			true,
			err,
			"Failure",
			"",
		}
	} else {
		statusMap[url] = FetchStatus{
			false,
			nil,
			"Success",
			output,
		}
	}

	return statusMap
}

func RefetchWebsites() {
	fmt.Println("refetch all websites")
}
