package main

// TODO: OS => human readable OS
// TODO: if the return value from the program is an error, send the calculation into an error state
// TODO: option for docker args, full command will be `docker run $os /opt/dockulator/calculattors/calc.$language '$calc'

import (
	"dockulator/db"
	"dockulator/models"
	"flag"
	"labix.org/v2/mgo/bson"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	dockerPath = "/usr/local/bin/docker"
	maxJobs    = 5 // Run this many `docker` processes concurrently
	pollDelay  = 2 // in seconds
)

var (
	throttle  = make(chan int, maxJobs)
	oses      = []string{"2b0268bd2e5b"}
	languages = []string{"rb"}
	debug     bool
)

func init() {
	for i := 0; i < maxJobs; i++ {
		throttle <- 1
	}

	flag.BoolVar(&debug, "debug", false, "set --debug to debug the docker output")
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	jobs := make(chan *models.Calculation)
	go ThrottledJobs(jobs)

	for {
		result := Poll()

		for i := 0; i < len(result); i++ {
			job := result[i]
			job.Language = PickString(languages)
			job.Instance = PickString(oses)
			jobs <- &job
		}
		time.Sleep(pollDelay * time.Second)
	}
}

func Poll() (result []models.Calculation) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)
	log.Printf("Polling Mongo")
	col.Find(bson.M{"instance": ""}).All(&result)
	log.Printf("Found %d calculations\n", len(result))
	return result
}

func PickString(slice []string) string {
	return slice[rand.Int()%len(slice)]
}

func ThrottledJobs(jobs chan *models.Calculation) {
	for job := range jobs {
		<-throttle
		log.Printf("Processing %s using %s on %s\n", job.Calculation, job.Language, job.Instance)
		go StartJob(job)
		throttle <- 1
	}
}

func StartJob(calculation *models.Calculation) {
	cmd := exec.Command(dockerPath, "run", calculation.Instance, "/opt/dockulator/calculators/calc."+calculation.Language, calculation.Calculation)
	if debug {
		log.Printf("args: %v", strings.Join(cmd.Args, " "))
		log.Println(cmd)
	}
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error from docker command: %s\n", err.Error())
	}
	floatVal := strings.TrimSpace(string(out))
	log.Printf("Value returned from docker: %s", string(floatVal))
	// update answer
	answer, err := strconv.ParseFloat(string(floatVal), 64)
	if err != nil {
		// send the calculation into the error state here.
		log.Printf("Could not convert answer to integer: %v", err.Error())
		return
	}
	calculation.Answer = answer
	calculation.Instance = calculation.Instance
	calculation.Save()
}
