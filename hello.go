package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

const (
	nameKey              string = "name"
	rmqURIEnvKey         string = "HELLO_RMQ_URI"
	defaultRmqURI        string = "amqp://guest:guest@localhost:5672/"
	rmqGreetingsExchange string = "greetings-exchange"
	rmqHelloRoutingKey   string = "hello"
	rmqGreetingsQueue    string = "greetings-queue"
)

type rmq struct {
	uri  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func (r *rmq) loadURIOrDefault(defaultURI string) {
	uri, ok := os.LookupEnv(rmqURIEnvKey)
	if ok {
		r.uri = uri
	} else {
		r.uri = defaultURI
	}
}

func (r *rmq) connect() {
	conn, err := amqp.Dial(r.uri)
	if err != nil {
		panic(err)
	}
	r.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	r.ch = ch
}

func (r *rmq) setup() {
	err := r.ch.ExchangeDeclare(
		rmqGreetingsExchange, "fanout", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	_, err = r.ch.QueueDeclare(rmqGreetingsQueue, false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err = r.ch.QueueBind(rmqGreetingsQueue, rmqHelloRoutingKey, rmqGreetingsExchange, false, nil)
	if err != nil {
		panic(err)
	}
}

type message struct {
	Person string `json:"person"`
}

var r *rmq

func main() {
	r = initRmq()
	server := initMux()
	log.Fatal(server.ListenAndServe())
}

func initRmq() *rmq {
	r := &rmq{}
	r.loadURIOrDefault(defaultRmqURI)
	r.connect()
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

	err = r.ch.Publish(
		rmqGreetingsExchange, rmqHelloRoutingKey, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s!", name)))
}
