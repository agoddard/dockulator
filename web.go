package main

import (
	"code.google.com/p/go.net/websocket"
	"dockulator/models"
	"fmt"
	"encoding/json"
	"html/template"
	"net/http"
	"log"
	"os"
	"time"
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
	pruneEvery = 10 // minutes?
)

type Clients map[string]*websocket.Conn

func (c *Clients) Prune() {
	fmt.Println("Pruning dead clients")
	fmt.Println("Number of clients:", len(*c))
	for k, client := range *c {
		go func () {
			err := websocket.Message.Send(client, "ping")
			//error will be: use of closed network connection
			if err != nil {
				delete(*c, k)
				fmt.Printf("Deleted %v from list of clients.\n", k)
				fmt.Println(err)
			}
		}()
	}
}

func (c *Clients) SendAll(msg string, data interface{}) {
	for _, client := range *c {
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
			w.Header().Add("Content-type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			b, _ := json.Marshal(calc)
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
		if msg == "disconnecting" {
		}
		fmt.Println("Message Got: ", msg)
	}
}

func main() {
	// Might just be able to get rid of this entirely with ReadWriteDeadline or something?
	go func () {
		for {
			time.Sleep(pruneEvery * time.Second)
			if len(clients) > 0 {
				clients.Prune()
			}
		}
	}()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculations", calculationsHandler)
	http.Handle("/websock", websocket.Handler(websocketHandler))
	fmt.Println("Serving on port", port)
	http.ListenAndServe(":" + port, nil)
}
