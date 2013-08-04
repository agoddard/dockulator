package main

import (
	"net/http"
	"net/url"
	"dockulator/models"
	"math/rand"
	"time"
	"flag"
	"log"
)

const (
	maxJobs    = 5 // Run this many `docker` processes concurrently
	pollDelay  = 2 // in seconds
	serverUrl = "http://localhost:5000/poller"
)

var (
	throttle  = make(chan int, maxJobs)
	oses      = []string{"2b0268bd2e5b","035b5d4ff9f4"}
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

func notifyServer(id string) {
	if debug {
		log.Printf("Notifying server that %v is finished\n", id)
	}
	_, err := http.PostForm(serverUrl, url.Values{"calculationId": {id}});
	if err != nil {
		log.Printf("Error notifying the server: %v\n", err.Error())
	}
}

func startJob(calculation *models.Calculation) {
	err := calculation.Calculate()
	if debug {
		log.Printf("Calculation: %v\n", calculation.Answer)
	}
	if err != nil {
		log.Printf("Got an error while claculating: %v", err.Error())
	}
	err = calculation.GetOS()
	if debug {
		log.Printf("OS hint: %v\n", calculation.OS)
	}
	if err != nil {
		log.Printf("Got an error getting calculation's OS: %v", err.Error())
	}
	calculation.Save()
	notifyServer(calculation.Id.Hex())
}

