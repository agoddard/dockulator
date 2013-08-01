package main

import (
	"os"
	"github.com/ChuckHa/calculations/calculations"
	"log"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

const (
	basePath   = "/calculations/"
	lenPath	   = len(basePath)
	collection = "calculations"
)

// Database collection
var c = calculations.GetSession().DB("").C(collection)

// A valid calculation
var calcRe = regexp.MustCompile(`^\s*\d+ [\+\-\*\/] \d+\s*$`)

// Templates
var indexTmpl = template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
var listTmpl = template.Must(template.ParseFiles("templates/base.html", "templates/calculations.html"))
var detailTmpl = template.Must(template.ParseFiles("templates/base.html", "templates/calculation_detail.html"))

// Handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, nil)
}

func calculationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		calculation := r.FormValue("calculation")
		found := calcRe.FindString(calculation)
		if found != "" {
			// Save the calculation in MongoDB
			calc := calculations.NewCalculation(calculation)
			calc.Insert(c)
		} else {
			http.Error(w, "Invalid calculation", 400)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
	var results []calculations.Calculation
	iter := c.Find(nil).Iter()
	err := iter.All(&results)
	if err != nil {
		log.Println("Error getting calculations from mongodb:", err)
	}
	listTmpl.Execute(w, results)
}
func calculationsIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[lenPath:]
	calc, err := calculations.Get(id, c)
	if err != nil {
		log.Println("Error getting a single calculation from mongodb:", err)
	}
	detailTmpl.Execute(w, calc)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculations", calculationsHandler)
	http.HandleFunc("/calculations/", calculationsIdHandler)
	listeningPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	http.ListenAndServe(listeningPort, nil)
}
