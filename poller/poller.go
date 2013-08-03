package main

import (
	"dockulator/models"
	"math/rand"
	"time"
	"flag"
	"log"
)

const (
	maxJobs    = 5 // Run this many `docker` processes concurrently
	pollDelay  = 2 // in seconds
)

var (
	throttle  = make(chan int, maxJobs)
	oses      = []string{"2b0268bd2e5b"}
	languages = []string{"rb"}
	debug bool
)

func init() {
	// Fill up the throttle
	for i := 0; i < maxJobs; i++ {
		throttle <- 1
	}
	flag.BoolVar(&debug, "debug", false, "set --debug to debug the docker output")
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	jobs := make(chan *models.Calculation)
	go processJobs(jobs)

	for {
		result := poll()
		if debug {
			log.Printf("Poll returned %v results\n", len(result))
		}

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
	if debug {
		log.Printf("Calculation: %v\n", calculation.Answer)
	}
	calculation.GetOS()
	if debug {
		log.Printf("OS hint: %v\n", calculation.OS)
	}
	calculation.Save()
}

