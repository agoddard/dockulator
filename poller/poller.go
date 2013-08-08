package main

import (
	"net/http"
	"net/url"
	"dockulator/models"
	"math/rand"
	"time"
	"flag"
	"log"
	"bytes"
)

const (
	maxJobs    = 5 // Run this many `docker` processes concurrently
	pollDelay  = 2 // in seconds
	serverUrl = "http://dockulator.herokuapp.com/poller"
)

var (
	throttle  = make(chan int, maxJobs)
	oses      = []string{"3b89fe9d1dfe"}
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

		if result.Processing {
			result.Language = sample(languages)
			result.Instance = sample(oses)
			if debug {
				log.Printf("Poll added a calculation to the queue:\n%v\n", result)
			}
			jobs <- &result
		}
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
	err := calculation.GetCalc()
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
	err = CleanUp()
	if err != nil {
		log.Printf("Got an error cleaning up docker containers: %v", err.Error())
	}
	notifyServer(calculation.Id.Hex())
}

// Cleanup function for docker. Might get better with -cidfile.
func CleanUp() error {
	// Get a list of running processes
	cmd := models.NewCommand("docker", "ps", "-a", "-q")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	lines := bytes.Split(out, []byte("\n"))
	instances := make([]string, len(lines) + 1)
	// Hack because sometimes go's type system is less flexible than it seems
	instances[0] = "rm"
	for i, line := range lines {
		instances[i + 1] = string(line)
	}
	// This removes all running processes
	cmd = models.NewCommand("docker", instances...)
	out, err = cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
