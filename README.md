# Dockulator. A Docker powered calculator.

A slightly over the top calculator that uses Docker to run your calculation on a random OS + language combination.

# Architecture
1) Webserver that is responsible for
        1) Serving a UI
        2) Validating new calculations
        3) Adding calculations to the processing queue
2) Mongodb is responsible for
        1) Acting as a processing queue
        2) Acting as a datastore for completed calculations
        3) Acting as an error store for errored calculations
3) One or more pollers
        1) running calculations
        2) telling Montodb about a completed calculation

## Poller

Poller monitors Mongodb for new calculations and tells docker to run these calculations in a container, choosing a specific image and language to run the calculation with. Poller is also responsible for processing the output from the docker container, organizing the responses in Mongodb, and error handling

### Flags

* `--debug`: Prints extra information. (TODO: Probably rename to --verbose/-v)

### Setup

Setting up a poller environment can be tricky.

1. Have a working docker environment.
2. Clone this repo
3. Run `make docker`
4. Run `go run poller/poller.go`

## Compiling

Having a go compiler compiled to build linux binaries is required for building linux binaries. Under OS X, running `brew uninstall go && brew install go --cross-compile-common` worked well. Then simply run `make linux` for linux binaries and `make` for local binaries.

## Calculator APIs

Each calculator must adhere to the same API.

    Accepts input as arguments to program in the form of
    NUM OP NUM
    NUM = number
    OP = ["+", "-", "*", "/"]
    Examples: 1.4 + 4; -5.11 / 8; 83242 * -0.2
    Writes answer to STDOUT

## Technologies

### [Go](http://golang.org)

We use Go as the main language in this project. It is the webserver (which also
handles websocket connections) and it is the poller that acts as a message
organizer. We use [mgo](http://labix.org/mgo) as the MongoDB library.

### [MongoDB](http://www.mongodb.org/)

MongoDB acts as both a datastore and a queue. One collection is used as the
message queue. Calculations enter the "queue" and get picked up by the poller.
They then sort themselves into the correct resulting collection either 
completed or error.

### [Docker](http://docker.io)

Docker is responsible for helping us run our calculations across different
operating systems. We put our calculators in the docker image and run
calculations through Docker. This is a somewhat odd use of Docker as Docker
tends to be used as a running container that accepts input and gives output.
This actually starts an image and runs the calculation on the running
container. This will actually create a new container that needs to be cleaned
up since we don't care about the resulting state of the image.

## Calculation Model

    {
      Calculation string        // "1.21 / -32.33"
      OS          string        // "Ubuntu 12.04 LTS"
      Language    string        // "rb"
      Id          bson.ObjectId // # MongoID
      Answer      float64       // -0.0374265388184349
      Instance    string        // "dockulator-ubuntu"
      Time        time.Time     // # timestamp
      Error       string        // "" #(hopefully)
      Processing  bool          // true
    }
