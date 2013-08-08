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
	serverUrl = "http://dockulator.herokuapp.com/poller"
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
			log.Printf("Poll added a calculation to the queue")
		}

		result.Language = sample(languages)
		result.Instance = sample(oses)
		jobs <- &result
		time.Sleep(pollDelay * time.Second)
	}
}

func poll() models.Calculation {
	return models.GetNext()
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
		log.Printf("Calc command: %v\n", calculation.CalcCmd().Args)
		log.Printf("Calculation: %v\n", calculation.Answer)
	}
	if err != nil {
		log.Printf("Got an error while claculating: %v", err.Error())
	}
	err = calculation.GetOS()
	if debug {
		log.Printf("OS command: %v\n", calculation.OSCmd().Args)
		log.Printf("OS hint: %v\n", calculation.OS)
	}
	if err != nil {
		log.Printf("Got an error getting calculation's OS: %v", err.Error())
	}
	calculation.Save()
	notifyServer(calculation.Id.Hex())
}

