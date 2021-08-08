package cron

import (
	LOGS "go_svelte_lighthouse/logs"
	REST "go_svelte_lighthouse/rest"

	cron "github.com/robfig/cron/v3"
)

func createCron(regularity string) (error) {
	var err error

	// init cron
	c := cron.New()

	// add function which cron should run
	c.AddFunc(regularity, func() { 
		// not sure about this callback approach in go...
		REST.RefetchWebsites(func() {
			LOGS.InfoLogger.Println("Cron was triggered successfully")
		}) 
	})

	// kick off
	go func() {
		c.Start()
	}()

	return err

}

// init fn to create cron
func init() {
	createCron("@weekly")
	// equivalent to: 
	// createCron("0 0 0 * * 0")
}
