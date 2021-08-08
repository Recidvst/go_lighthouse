package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"
)

func CreateReport(urlArg string, getDesktop bool) (success bool, err error) {
	
	fmt.Println("CreateReport")
	// full url of the website to be checked
	urlToFetch := urlArg

	var urlAsFolder string

	// strip protocol, args etc. from the url
	parsedURL, err := url.Parse(urlToFetch)
	if err != nil {
		log.Fatal(err)
	}
	urlAsFolder = parsedURL.Host

	// get the date, to use in outputPath
	// currentDate := time.Now()
	// dateAsFilename := currentDate.Format("02012006") // date as DDMMYYYY

	// get current working directory
	var cwd, _ = os.Getwd()

	// check folders exist, if not then create them otherwise cli command will error
	if _, checkDirErr := os.Stat(cwd + "/reports"); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir(cwd+"/reports", 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return false, err
		}
	}
	if _, checkDirErr := os.Stat(cwd + "/reports/" + urlAsFolder); os.IsNotExist(checkDirErr) {
		mkdirErr := os.Mkdir(cwd+"/reports/"+urlAsFolder, 0755)
		if mkdirErr != nil {
			err = mkdirErr
			return false, err
		}
	}

	// handle desktop vs mobile presets
	var presetFlag string = ""
	if getDesktop {
		presetFlag = "--preset=desktop"
	}

	// build slice of flags to pass to exec
	flags := []string{urlToFetch, "--output=json", "-quiet", "--chrome-flags='--headless'", "--output-path=stdout", presetFlag}
	fmt.Println(flags)

	// run cli command
	// c, b := exec.Command("lighthouse", flags...), new(strings.Builder)
	// c.Stdout = b
	// c.Run()

	// get stdout as string var
	// var stdoutString = b.String()
	// fmt.Println(stdoutString)

	// capturing stdout
	// stdoutReader, execErr := cmd.StdoutPipe()

	// if execErr != nil {
	// 	log.Fatal(execErr)
	// 	return false, execErr
	// }

	// scanner := bufio.NewScanner(stdoutReader)
	// go func() {
	// 	for scanner.Scan() {
	// 		fmt.Printf("\t > %s\n", scanner.Text())
	// 	}
	// }()

	// err = cmd.Start()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
	// 	return false, err
	// }

	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
	// 	return false, err
	// }

	// TODO: now inject into Database

	return true, nil
}
