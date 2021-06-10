package cli

import (
	"os"
	"os/exec"
	"strings"
	"time"
)

func CreateReport(url string) (outputPath string, err error) {
	// full url of the website to be checked
	urlToFetch := url

	// remove the protocol and host, to use in outputPath
	urlSliced := strings.Split(urlToFetch, ".")
	urlAsFolder := strings.Join(urlSliced[1:], "_")

	// add a trailing slash to the url if it doesn't have one
	if string(urlAsFolder[len(urlAsFolder)-1]) != string("/") {
		urlAsFolder = urlAsFolder + "/"
	}

	// get the date, to use in outputPath
	currentDate := time.Now()
	dateAsFilename := currentDate.Format("02012006") // date as DDMMYYYY

	// where the json output will be saved
	outputPath = "reports/" + urlAsFolder + dateAsFilename + ".json"

	// check folder exist, if not then create them otherwise cli command will error
	if _, checkDirErr := os.Stat("reports"); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir("reports", 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return outputPath, err
		}
		if _, checkDirErr := os.Stat("reports/" + urlAsFolder); os.IsNotExist(checkDirErr) {
			mkdirErr := os.Mkdir("reports/"+urlAsFolder, 0755)
			if mkdirErr != nil {
				err = mkdirErr
				return outputPath, err
			}
		}
	}

	// build slice of flags to pass to exec
	flags := []string{urlToFetch, "--output=json", "--chrome-flags='--headless'", "--output-path=" + outputPath}

	// run cli command
	cmd := exec.Command("lighthouse", flags...)
	if execErr := cmd.Run(); execErr != nil {
		err = execErr
		return outputPath, err
	}

	return outputPath, nil
}
