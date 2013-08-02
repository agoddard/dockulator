package main

import (
	"code.google.com/p/go.net/websocket"
	"dockulator/models"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
)

func init () {
	port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

}

const (
	basePath = "/calculations/"
	lenPath  = len(basePath)
)

var (
	port string
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
		found := models.CleanCalculation(calculation)
		if found != "error" {
			calc := models.NewCalculation(found)
			calc.Insert()
			w.WriteHeader(http.StatusCreated)
			return
		}
		log.Printf("Got weird input: %v", calculation)
		http.Error(w, "Invalid calculation", 400)
		return
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func websocketHandler(ws *websocket.Conn) {
	calcs := string(models.GetRecent(3).Json())
	websocket.JSON.Send(ws, []byte(calcs))
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculations", calculationsHandler)
	http.Handle("/websock", websocket.Handler(websocketHandler))
	fmt.Println("Serving on port", port)
	http.ListenAndServe(":" + port, nil)
}
