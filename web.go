package main

import (
	"code.google.com/p/go.net/websocket"
	"dockulator/models"
	"dockulator/db"
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
	basePath   = "/calculations/"
	lenPath    = len(basePath)
	pruneEvery = 10 // minutes?
)

var (
	port    string
	funcMap = template.FuncMap{
		"lang": models.GetLanguage,
	}
	// Templates
	baseTmpl  = template.Must(template.ParseFiles("templates/base.html")).Funcs(funcMap)
	indexTmpl = template.Must(baseTmpl.ParseFiles("templates/index.html"))
	clients   = make(Clients)
	pollerSecret = os.Getenv("POLLER_SECRET")
)

type Clients map[string]*websocket.Conn

func (c *Clients) Prune() {
	fmt.Println("Pruning dead clients")
	fmt.Println("Number of clients:", len(*c))
	for k, client := range *c {
		err := websocket.Message.Send(client, "ping")
		//error will be: use of closed network connection
		if err != nil {
			delete(*c, k)
			fmt.Printf("Deleted %v from list of clients.\n", k)
			fmt.Println(err)
		}
	}
}

func (c *Clients) SendAll(msg string, data interface{}) {
	message := BuildMsg(msg, data)
	for _, client := range *c {
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
	fmt.Println("Connected client")
	ip := ws.Request().RemoteAddr
	_, ok := clients[ip]
	if ok {
		websocket.JSON.Send(ws, BuildMsg("error", "Your IP is already connected"))
		ws.Close()
		return
	}
	clients[ip] = ws
	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			break
		}
		fmt.Println("Message Got: ", msg)
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
	// Might just be able to get rid of this entirely with ReadWriteDeadline or something?
	go func() {
		for {
			time.Sleep(pruneEvery * time.Second)
			if len(clients) > 0 {
				clients.Prune()
			}
		}
	}()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/calculations", calculationsHandler)
	http.Handle("/websock", websocket.Handler(websocketHandler))
	http.HandleFunc("/poller", pollerHandler)
	http.HandleFunc("/", indexHandler)
	fmt.Println("Serving on port", port)
	http.ListenAndServe(":"+port, nil)
}
