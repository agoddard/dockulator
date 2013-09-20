package main

import (
	"code.google.com/p/go.net/websocket"
	"dockulator/db"
	"dockulator/models"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
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
	port    string
	funcMap = template.FuncMap{
		"lang": models.GetLanguage,
	}
	// Templates
	baseTmpl     = template.Must(template.ParseFiles("templates/base.html")).Funcs(funcMap)
	indexTmpl    = template.Must(baseTmpl.ParseFiles("templates/index.html"))
	clients      = make(Clients, 0)
	pollerSecret = os.Getenv("POLLER_SECRET")
)

type Clients []*websocket.Conn

func (c *Clients) AddClient(ws *websocket.Conn) {
	c = append(c, ws)
}
func (c *Clients) RemoveClient(ws *websocket.Conn) {
	for i, client := range c {
		if client == ws {
			c = append(c[:i], c[i+1:]...)
			return
		}
	}
}

func (c *Clients) SendAll(msg string, data interface{}) {
	message := BuildMsg(msg, data)
	for _, client := range c {
		websocket.JSON.Send(client, message)
	}
}

func BuildMsg(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": msg,
		"data": data,
	}
}

// Handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {
	calcs := models.GetRecent(20)
	indexTmpl.Execute(w, calcs)
}

func calculationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		calculation := r.FormValue("calculation")
		found := models.CleanCalculation(calculation)
		if found != "error" {
			calc := models.NewCalculation(found)
			// Add them to the queue
			calc.Insert(db.Queue)
			w.Header().Add("Content-type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			b, err := calc.AsJson()
			if err != nil {
				log.Printf("Error marshalling calc as JSON: %v", err.Error())
			}
			clients.SendAll("new", calc)
			w.Write(b)
			return
		}
		log.Printf("Got weird input: %v", calculation)
		http.Error(w, "Invalid calculation", 400)
		return
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

// set a read deadline here probably
func websocketHandler(ws *websocket.Conn) {
	defer ws.Close()
	clients.AddClient(ws)
	defer clients.RemoveClient(ws)
	for {
		var msg string
		// Just chill. We don't expect to be receiving any messages.
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			break
		}
	}
}

func pollerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	/*
		secret := r.PostFormValue("secret")
		if secret != pollerSecret {
			time.Sleep(time.Second)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
	*/
	id := r.PostFormValue("calculationId")
	calc, _ := models.Get(id)
	clients.SendAll("update", calc)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/calculations", calculationsHandler)
	http.Handle("/websock", websocket.Handler(websocketHandler))
	http.HandleFunc("/poller", pollerHandler)
	http.HandleFunc("/", indexHandler)
	fmt.Println("Serving on port", port)
	http.ListenAndServe(":"+port, nil)
}
