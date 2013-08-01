package main

import (
	"dockulator/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	basePath = "/calculations/"
	lenPath  = len(basePath)
)

var (
	// A valid calculation
	calcRe = regexp.MustCompile(`^\s*\d+ [\+\-\*\/] \d+\s*$`)

	// Templates
	indexTmpl  = template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	listTmpl   = template.Must(template.ParseFiles("templates/base.html", "templates/calculations.html"))
	detailTmpl = template.Must(template.ParseFiles("templates/base.html", "templates/calculation_detail.html"))
)

// Handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, nil)
}

func calculationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		calculation := r.FormValue("calculation")
		found := calcRe.FindString(calculation)
		if found != "" {
			calc := models.NewCalculation(calculation)
			calc.Insert()
		} else {
			http.Error(w, "Invalid calculation", 400)
		}
		// Definitely change this
		http.Redirect(w, r, "/", http.StatusFound)
	}

	/* don't do this yet
	var results []models.Calculation
	iter := c.Find(nil).Iter()
	err := iter.All(&results)
	if err != nil {
		log.Println("Error getting calculations from mongodb:", err)
	}
	listTmpl.Execute(w, results)
	*/
}
func calculationsIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[lenPath:]
	calc, err := models.Get(id)
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
