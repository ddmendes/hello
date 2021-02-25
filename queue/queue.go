package queue

import (
	"os"
	"time"

	"github.com/streadway/amqp"
)

const (
	rmqURIEnvKey         string = "HELLO_RMQ_URI"
	rmqGreetingsExchange string = "greetings-exchange"
	rmqHelloRoutingKey   string = "hello"
	rmqGreetingsQueue    string = "greetings-queue"
)

// Rmq holds necessary data for connection with rabbitmq
type Rmq struct {
	uri  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

// LoadURIOrDefault gets rabbitmq URI from HELLO_RMQ_URI env variable or
// uses the default at defaultURI parameter.
func (r *Rmq) LoadURIOrDefault(defaultURI string) {
	uri, ok := os.LookupEnv(rmqURIEnvKey)
	if ok {
		r.uri = uri
	} else {
		r.uri = defaultURI
	}
}

// Connect open connection with RabbitMQ
func (r *Rmq) Connect() {
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

// Setup exchange, queue and queue bind
func (r *Rmq) Setup() {
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

// Publish a body to greetings exchange and hello routing key
func (r *Rmq) Publish(body []byte) error {
	return r.ch.Publish(
		rmqGreetingsExchange, rmqHelloRoutingKey, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "application/json",
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
}
