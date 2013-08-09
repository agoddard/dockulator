# Distributed Dockulator. A laptop cloud.

A user boots a VM through a Vagrantfile we provide.

It asks them for an email.
        The email is sent to our auth server to register this user.
        The user is registered by email and is returned a key.
        This key is used to sign all of their requests. # As an environment variale or something.
The user now has poller running in the background.

## Questions for Anthony:

* Does the poller talk directly to Mongo? I think that's a bad idea.
* Does the poller just talk to a webserver?
        * Perhaps it polls dockulator.com/next?key=SOME_SECRET_KEY
        * then dockulator.com/next looks at the key and verifies it
                * dockulator.com/next pings mongo and gets a calculation back
                * dockulator.com/next responds to the poller with json
        * the poller runs the calculation and sends it to dockulator.com/add?key=SOME_SECRET_KEY
                * dockulator.com/add verifies the request
                * dockulator.com/add sends the result to either the error queue or the completed queue.
                * Note: we could even do something clever like any time a user returns an error status, 
                        we put it in the error queue and then add the same calculation back to the processing queue.

So we have

1) Webserver that is responsible for
        1) Serving a UI
        2) Validating new calculations
        3) Adding calculations to the processing queue
2) Mongodb is responsible for
        1) Acting as a processing queue
        2) Acting as a datastore for completed calculations
        3) Acting as an error store for errored calculations
3) An auth server that is responsible for
        1) registering laptops
        2) distributing keys
        3) validating requests
        4) getting items off the queue
        5) adding calculations to the completed/error collection
        6) collects stats on laptop pollers
3) A single laptop poller
        1) running calculations
        2) telling the auth server about a completed calculation

# Dockulator

A slightly over the top calculator that uses Docker to run your calculation on
a random OS + language combination.

## Poller

Poller is probably a misleading name at this point. It does indeed poll the
Mongo collection "queue" to find unprocessed calculations, but it also is
responsible for telling the calculations to calculate themselves as well as
organizing responses into the correct collections. For instance, any
calculation that has an Error will be sent to the "errors" collection.

### Flags

* `--debug`: Prints extra information. (TODO: Probably rename to --verbose/-v)

### Setup

Setting up a poller environment can be tricky.

1. Have a working docker environment.
2. Get this repo
3. Run `make docker`
4. Run `go run poller/poller.go` or upstart or something

## Compiling

Having a go compiler compiled to build linux binaries is required for building
linux binaries.  I found running `brew uninstall go && brew install go
--cross-compile-common` worked well.  Then simply run `make linux` for linux
binaries and `make` for local binaries.

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
