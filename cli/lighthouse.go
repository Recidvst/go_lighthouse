package cli

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"time"
)

func CreateReport(urlArg string, getDesktop bool) (outputPath string, err error) {
	// full url of the website to be checked
	urlToFetch := urlArg

	var urlAsFolder string

	// strip protocol, args etc. from the url
	parsedURL, err := url.Parse(urlToFetch)
	if err != nil {
		log.Fatal(err)
	}
	urlAsFolder = parsedURL.Host;

	// get the date, to use in outputPath
	currentDate := time.Now()
	dateAsFilename := currentDate.Format("02012006") // date as DDMMYYYY
	
	// get current working directory
	var cwd, _ = os.Getwd()

	// where the json output will be saved
	outputPath = cwd + "/reports/" + urlAsFolder + "/" + dateAsFilename + ".json"

	// check folders exist, if not then create them otherwise cli command will error
	if _, checkDirErr := os.Stat(cwd + "/reports"); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir(cwd + "/reports", 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return outputPath, err
		}
	}
	if _, checkDirErr := os.Stat(cwd + "/reports/" + urlAsFolder); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir(cwd + "/reports/"+ urlAsFolder, 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return outputPath, err
		}
	}

	// handle desktop vs mobile presets
	var presetFlag string = ""
	if getDesktop {
		presetFlag = "--preset=desktop"
	}

	// build slice of flags to pass to exec
	flags := []string{urlToFetch, "--output=json", "--chrome-flags='--headless'", "--output-path=" + outputPath, presetFlag}

	// run cli command
	cmd := exec.Command("lighthouse", flags...)
	if execErr := cmd.Run(); execErr != nil {
		err = execErr
		return outputPath, err
	}

	return outputPath, nil
}
