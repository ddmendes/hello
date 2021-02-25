package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ddmendes/hello/queue"
	"github.com/gorilla/mux"
)

const (
	nameKey       string = "name"
	defaultRmqURI string = "amqp://guest:guest@localhost:5672/"
)

type message struct {
	Person string `json:"person"`
}

var r *queue.Rmq

func main() {
	r = initRmq()
	server := initMux()
	log.Fatal(server.ListenAndServe())
}

func initRmq() *queue.Rmq {
	r := &queue.Rmq{}
	r.LoadURIOrDefault(defaultRmqURI)
	r.Connect()
	r.Setup()
	return r
}

func initMux() *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", hello)
	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}

func hello(w http.ResponseWriter, rq *http.Request) {
	name := rq.FormValue(nameKey)
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing name parameter"))
		return
	}

	body, err := json.Marshal(message{name})
	if err != nil {
		panic(err)
	}

	err = r.Publish(body)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s!", name)))
}
