package cli

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

// CreateReport function handling the CLI command and returning result
func CreateReport(urlArg string, getDesktop bool) (success bool, resultString string, err error) {

	// mutex lock
	var mutex = &sync.Mutex{}

	// full url of the website to be checked
	urlToFetch := urlArg

	var urlAsFolder string

	// strip protocol, args etc. from the url
	parsedURL, err := url.Parse(urlToFetch)
	if err != nil {
		log.Fatal(err)
	}
	urlAsFolder = parsedURL.Host

	// get current working directory
	var cwd, _ = os.Getwd()

	// check folders exist, if not then create them otherwise cli command will error
	if _, checkDirErr := os.Stat(cwd + "/reports"); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir(cwd+"/reports", 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return false, "", err
		}
	}
	if _, checkDirErr := os.Stat(cwd + "/reports/" + urlAsFolder); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir(cwd+"/reports/"+urlAsFolder, 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return false, "", err
		}
	}

	// handle desktop vs mobile presets
	var presetFlag = ""
	if getDesktop {
		presetFlag = "--preset=desktop"
	}

	// set temporary output path for the file we will read into memory as a string
	temporaryOutputFileName := urlAsFolder + "__" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10) + ".temp"
	temporaryOutputFilePath := "--output-path=" + temporaryOutputFileName

	// build slice of flags to pass to exec
	flags := []string{urlToFetch, "--output=json", "-quiet", "--chrome-flags='--headless'", temporaryOutputFilePath, presetFlag}

	// run cli command
	cmd := exec.Command("lighthouse", flags...)

	// catch error
	if execErr := cmd.Run(); execErr != nil {
		err = execErr
		return false, "", err
	}

	// now get the temporary file we created and read into a string
	temporaryOutputFileAsBytes, err := ioutil.ReadFile(temporaryOutputFileName)
	if err != nil {
		return false, "", err
	}

	// convert bytes content to a string
	temporaryOutputFileAsString := string(temporaryOutputFileAsBytes)

	// delete temporary file we created to clean up
	mutex.Lock()
	deleteErr := os.Remove(temporaryOutputFileName)
	mutex.Unlock()
	if deleteErr != nil {
		return false, "", deleteErr
	}

	return true, temporaryOutputFileAsString, nil
}
