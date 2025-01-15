package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	webPort           = "80"
	maxRetries        = 5
	initialBackoff    = 1 * time.Second
	backoffMultiplier = 2
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var connection *amqp.Connection

	for attempt := 1; attempt <= maxRetries; attempt++ {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Printf("Attempt %d: RabbitMQ not ready. Retrying in %v...", attempt, initialBackoff)
		} else {
			log.Println("RabbitMQ connection established.")
			connection = c
			return connection, nil
		}

		// Exponential backoff
		backOff := time.Duration(math.Pow(backoffMultiplier, float64(attempt))) * initialBackoff
		time.Sleep(backOff)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d retries", maxRetries)
}
