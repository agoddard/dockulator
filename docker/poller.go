package main

// TODO: OS => human readable OS
// TODO: if the return value from the program is an error, send the calculation into an error state
// TODO: option for docker args, full command will be `docker run $os /opt/dockulator/calculattors/calc.$language '$calc'

import (
	"fmt"
	calc "github.com/ChuckHa/calculations/calculations"
	"labix.org/v2/mgo/bson"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
	"flag"
	"strings"
)

const (
	dockerPath = "docker" // FIXME: should be full path to docker binary
	maxJobs    = 5        // Run this many `docker` processes concurrently
	pollDelay  = 2        // in seconds
)

var (
	throttle  = make(chan int, maxJobs)
	oses      = []string{"2b0268bd2e5b"}
	languages = []string{"sh"}
	c         = calc.GetCollection()
	debug bool
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

	var result []calc.Calculation
	jobs := make(chan calc.Calculation)
	go ThrottledJobs(jobs)

	for {
		fmt.Println("Polling Mongo")
		c.Find(bson.M{"instance": ""}).All(&result)
		fmt.Printf("Found %d calculations\n", len(result))

		for i := 0; i < len(result); i++ {
			job := result[i]
			job.Language = PickString(languages)
			job.OS = PickString(oses)
			jobs <- job
		}
		time.Sleep(pollDelay * time.Second)
	}
}

func PickString(slice []string) string {
	return slice[rand.Int()%len(slice)]
}

func ThrottledJobs(jobs chan calc.Calculation) {
	for job := range jobs {
		<-throttle
		fmt.Printf("Processing %s using %s on %s\n", job.Calculation, job.Language, job.OS)
		go StartJob(job)
		throttle <- 1
	}
}

func StartJob(calculation calc.Calculation) {
	cmd := exec.Command(dockerPath, "run", calculation.OS, "/opt/dockulator/calculators/calc."+calculation.Language, " \"", calculation.Calculation, "\"")
	if debug {
		fmt.Printf("Command: %v\n", strings.Join(cmd.Args, " "))
	}
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error from docker command: %s\n", err)
	}
	log.Println(string(out))
	// update answer
	answer, err := strconv.Atoi(string(out))
	if err != nil {
		// send the calculation into the error state here.
		log.Printf("Could not convert answer to integer")
	}
	calculation.Answer = answer
	calculation.Instance = calculation.OS
	calculation.Save(c)
}
