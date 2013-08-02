package main

import (
	"dockulator/models"
	"math/rand"
	"time"
)

const (
	maxJobs    = 5 // Run this many `docker` processes concurrently
	pollDelay  = 2 // in seconds
)

var (
	throttle  = make(chan int, maxJobs)
	oses      = []string{"2b0268bd2e5b"}
	languages = []string{"rb"}
)

func init() {
	// Fill up the throttle
	for i := 0; i < maxJobs; i++ {
		throttle <- 1
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	jobs := make(chan *models.Calculation)
	go processJobs(jobs)

	for {
		result := poll()

		for i := 0; i < len(result); i++ {
			calculation := result[i]
			calculation.Language = sample(languages)
			calculation.Instance = sample(oses)
			jobs <- &calculation
		}
		time.Sleep(pollDelay * time.Second)
	}
}

func poll() (result models.Calculations) {
	return models.UnsolvedCalculations()
}

func sample(slice []string) string {
	return slice[rand.Int()%len(slice)]
}

func processJobs(calculations chan *models.Calculation) {
	for calculation := range calculations {
		// wait for a value on throttle before starting
		<-throttle
		go startJob(calculation)
		throttle <- 1
	}
}

func startJob(calculation *models.Calculation) {
	calculation.Calculate()
	calculation.GetOS()
	calculation.Save()
}

