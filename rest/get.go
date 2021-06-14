package rest

import (
	"fmt"
)

func GetWebsiteStats(url string) {
	fmt.Println("get specific website")
	fmt.Println(url)
}

func GetAllWebsitesStats() {
	fmt.Println("get all websites")
}
