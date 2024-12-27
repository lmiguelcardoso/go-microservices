package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	maxRetries        = 5
	initialBackoff    = 1 * time.Second
	backoffMultiplier = 2
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	// start listening for message

	// create consumer

	// watch the queue and consume events
}

func connect() (*amqp.Connection, error) {
	var connection *amqp.Connection

	for attempt := 1; attempt <= maxRetries; attempt++ {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
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
