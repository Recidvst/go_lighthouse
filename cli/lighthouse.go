package main // rename to cli after testing

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
    // full url of the website to be checked
    urlToFetch := "https://www.chris-snowden.me/"

    // remove the protocol and host, to use in outputPath
    urlSliced := strings.Split(urlToFetch, ".")
    urlAsFolder := strings.Join(urlSliced[1:], "_")

    // add a trailing slash to the url if it doesn't have one
    if (string(urlAsFolder[len(urlAsFolder)-1]) != string("/")) {
        urlAsFolder = urlAsFolder + "/"
    }

    // get the date, to use in outputPath
    currentDate := time.Now()
    dateAsFilename := currentDate.Format("02012006") // date as DDMMYYYY

    // where the json output will be saved
    outputPath := "../reports/" + urlAsFolder + dateAsFilename + ".json"
    fmt.Println(outputPath)

    // check folder exist, if not then create them otherwise cli command will error
    if _, checkDirErr := os.Stat("../reports/" + urlAsFolder); os.IsNotExist(checkDirErr) {
        mkdirErr := os.Mkdir("../reports/" + urlAsFolder, 0755)
        if (mkdirErr != nil) {
            log.Fatalln(mkdirErr)
        }   
    }

    // build slice of flags to pass to exec
    flags := []string{urlToFetch, "--output=json", "--output-path=" + outputPath}
    fmt.Println(flags)

    // run cli command
    cmd := exec.Command("lighthouse", flags...)
    if err := cmd.Run(); err != nil {
        log.Fatalln(err)
    }
}
