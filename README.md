dockulator
==========

dockulator

## Compiling

Having a go compiler compiled to build linux binaries is required for building linux binaries.
I found running `brew uninstall go && brew install go --cross-compile-common` worked well.
Then simply run `make linux` for linux binaries and `make` for local binaries.


####user hits web server
####user enters a calculation
####calculation is stored in mongo as follows:

####docker middleware polls redis for calculations with no answer
####docker middleware receives a calculation, spins up random docker OS/language combination
MVP - use shell-out and then switch to HTTP/API later
####docker receives request for OS/language, runs request, responds via HTTP PUT to ?
    PUT data-->Mongo unique ID, request ID, metadata, calculation response

## Calculator APIs

Each calculator must adhere to the same API.

    Accepts input as arguments to program (e.g. python calc.py "3 + 4")
    Writes answer to STDOUT

## Basic Stack

* Webserver: [Go](http://golang.org) -- Handles UI, creation of calculation objects and display of completed calculations
* Datastore: [MongoDB](http://mongodb.org) -- Stores data
* Poller: Go -- Looks for empty calculations and passes them off to Docker
* Containerization: [Docker](http://docker.io) -- Does calculation

## Use Case

1. Bianca goes to dockulator.com
1. She sees a form to do a calculation
2. She types a calculation into the form
3. She submits the form
  1. The form goes to the webserver
  2. The webserver parses the input
  3. If valid, the input is sent to mongodb
4. Bianca's browser is constantly polling (or getting push data from a websocket) for new calculations
5. Eventually she sees her calculation complete

## WebServer API

    GET  /                 -- Home page
    GET  /calculations     -- Calculation list view
    GET  /calculations/:id -- Calculation detail view
    POST /calculations     -- Creates a new calculation in MongoDB


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

### Interactive Docker session
In this model, the poller will request an instance from docker and wait for the response before updating MongoDB with the response. This could potentially be blocking if not done concurrently, we will also need to manage error states - if the response from docker is an error, the poller will still try to parse it (though if the response simply doesn't match what we expect, it's ok to throw the calculation into an error state).
####daemon
In this model, the poller adds the '-d' flag to the docker command, which will daemonize the docker instance and return an ID. The poller will record this ID and then update the MongoDB calculation with this instance ID. When the instance has completed the calculation, the instance itself will update the calculation in MongoDB. This could be done either by implementing an HTTP or MongoDB call in each calculator, or by having the calculator send its results to a local go script which takes care of the mongo connection (this script might just recieve a json string from the calculator and an ID)

## UI

