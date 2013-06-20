package main

import (
	"fmt"
	calc "calculations/calculations"
	"labix.org/v2/mgo/bson"
	"os/exec"
	"log"
	"math/rand"
	"time"
	"strconv"
)

const (
	dockerPath = "docker"
	maxJobs = 5
)

var throttle = make(chan int, maxJobs)
var oses = []string{"jdoi34j", "djwoief", "fj93jg4", "jfio23jf"}
var languages = []string{"rb", "pl", "py", "sh"}
var c = calc.GetCollection()

func init () {
	for i := 0; i < maxJobs; i++ {
		throttle <- 1
	}
}

func main () {
	rand.Seed( time.Now().UTC().UnixNano())

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
		time.Sleep(2 * time.Second)
	}
}

func PickString(slice []string) string {
	return slice[rand.Int() % len(slice)]
}

func ThrottledJobs(jobs chan calc.Calculation) {
	for job := range jobs{
		<-throttle
		fmt.Printf("Processing %s using %s on %s\n", job.Calculation, job.Language, job.OS)
		go StartJob(job)
		throttle <- 1
	}
}

func StartJob(calculation calc.Calculation) {
	// Do all of this in a goroutine
	args := []string {
		"-d",
		"/opt/calculate." + calculation.Language,
		"\"" + calculation.Calculation + "\"",
		"--id=" + calculation.OS,
	}
	// delete me
	time.Sleep(1 * time.Second)
	cmd := exec.Command(dockerPath, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error from docker command: %s\n", err)
	}
	log.Println(string(out))
	// update answer
	answer, err := strconv.Atoi(string(out))
	if err != nil {
		log.Printf("Could not convert answer to integer")
	}
	calculation.Answer = answer
	calculation.Save(c)
}
