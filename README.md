dockulator
==========

dockulator

####user hits web server
####user enters a calculation
####calculation is stored in mongo as follows:

####docker middleware polls redis for calculations with no answer
####docker middleware receives a calculation, spins up random docker OS/language combination
MVP - use shell-out and then switch to HTTP/API later
####docker receives request for OS/language, runs request, responds via HTTP PUT to ?
    PUT data-->Mongo unique ID, request ID, metadata, calculation response

## Basic Stack

* Webserver: [Go](http://golang.org) -- Handles UI, creation of calculation objects and display of completed calculations
* Datastore: [MongoDB](http://mongodb.org) -- Stores data
* Middleware: [Ruby](http://ruby-lang.org) -- Looks for empty calculations and passes them off to Docker
* Containerization: [Docker](http://docker.io) -- Does calculation

# Calculator APIs

Each calculator must adhere to the same API.

   Accepts input from STDIN
   Responds with input on STDOUT

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
    GET  /calculations     -- Returns all calculations
    POST /calculations    -- Creates a new calculation in MongoDB
    PUT  /calculations/:id -- Updates a document in MongoDB


## Calculation Model

    {
      "Calculation": "1 + 4",
      "_id":         23ab235feeda31098,
      "OS":          null, // until complete
      "Language":    null, // until complete
      "Answer":      null, // until complete
      "Time":        1234932849028342
    }
