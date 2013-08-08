dockulator
==========

A slightly over the top calculator that uses Docker to run your calculation on a random OS + language combination.

## Compiling

Having a go compiler compiled to build linux binaries is required for building linux binaries.
I found running `brew uninstall go && brew install go --cross-compile-common` worked well.
Then simply run `make linux` for linux binaries and `make` for local binaries.

## Calculator APIs

Each calculator must adhere to the same API.

    Accepts input as arguments to program in the form of
    NUM OP NUM
    NUM = integer
    OP = ["+", "-", "*", "/"]
    Examples: 1 + 4; 5 / 8; 83242 * 288338
    Writes answer to STDOUT

## Basic Stack

* Webserver: [Go](http://golang.org) -- Handles UI, creation of calculation objects and display of completed calculations
* Datastore: [MongoDB](http://mongodb.org) -- Stores data
* Poller: Go -- Looks for empty calculations and passes them off to Docker
* Containerization: [Docker](http://docker.io) -- Does calculation

## Calculation Model

    {
      "Calculation": "1 + 4",
      "_id":         23ab235feeda31098,
      "OS":          null, // until complete
      "Language":    null, // until complete
      "Answer":      null, // until complete
      "Time":        1234932849028342
    }

## Poller

The poller performs a query to Mongdb for all calculations that have not been processed. 
For each non processd calculation, the poller will launch an interactive docker instance concurrently (max of 5 processes at once).

### Flags

* `--debug`: Prints extra information. (TODO: Probably rename to --verbose/-v)
        * the command that will be executed before it is executed

### Interactive Docker session
In this model, the poller will request an instance from docker and wait for the response before updating MongoDB with the response. This could potentially be blocking if not done concurrently, we will also need to manage error states - if the response from docker is an error, the poller will still try to parse it (though if the response simply doesn't match what we expect, it's ok to throw the calculation into an error state).
####daemon
In this model, the poller adds the '-d' flag to the docker command, which will daemonize the docker instance and return an ID. The poller will record this ID and then update the MongoDB calculation with this instance ID. When the instance has completed the calculation, the instance itself will update the calculation in MongoDB. This could be done either by implementing an HTTP or MongoDB call in each calculator, or by having the calculator send its results to a local go script which takes care of the mongo connection (this script might just recieve a json string from the calculator and an ID)

## UI


## TODO

* ~~Client connects and gets 20 most recent calculations~~
* ~~Opens websocket server~~
* ~~Client submits calculation~~
* ~~Webserver adds calculation to mongo~~
* ~~Poller calculations calculation~~
* Poller notifies web server it is done
* webserver notifies all open websocket clients
