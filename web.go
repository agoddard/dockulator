package main

import (
	"code.google.com/p/go.net/websocket"
	"dockulator/models"
	"fmt"
	"html/template"
	"net/http"
	"log"
	"os"
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

type Clients map[string]*websocket.Conn

func (c Clients) SendAll(msg string, data interface{}) {
	for _, client := range c {
		websocket.JSON.Send(client, BuildMsg(msg, data))
	}
}

func BuildMsg(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": msg,
		"data": data,
	}
}

var (
	port string
	// Templates
	indexTmpl  = template.Must(template.ParseFiles("templates/base.html", "templates/index.html"))
	listTmpl   = template.Must(template.ParseFiles("templates/base.html", "templates/calculations.html"))
	detailTmpl = template.Must(template.ParseFiles("templates/base.html", "templates/calculation_detail.html"))
	clients = make(Clients)
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
	ip := ws.Request().RemoteAddr
	_, ok := clients[ip]; if ok {
		websocket.JSON.Send(ws, BuildMsg("error", "Your IP is already connected"))
		ws.Close()
		return
	}
	clients[ip] = ws
	calcs := models.GetRecent(3)
	websocket.JSON.Send(ws, BuildMsg("initialData", calcs))
	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil{
			break
		}
		fmt.Println("Message Got: ", msg)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculations", calculationsHandler)
	http.Handle("/websock", websocket.Handler(websocketHandler))
	fmt.Println("Serving on port", port)
	http.ListenAndServe(":" + port, nil)
}
