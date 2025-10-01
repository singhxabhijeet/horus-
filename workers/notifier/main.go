package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// HealthCheckResult is the structure of the message we expect to receive
type HealthCheckResult struct {
	SiteID         int   `json:"SiteID"`
	IsUp           bool  `json:"IsUp"`
	StatusCode     int   `json:"StatusCode"`
	ResponseTimeMs int64 `json:"ResponseTimeMs"`
}

func main() {
	// Get RabbitMQ URL from environment variable
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	// Connect to RabbitMQ (with retry logic)
	var conn *amqp091.Connection
	var err error
	for i := 0; i < 5; i++ {
		conn, err = amqp091.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in 5s... (%v)", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ after retries: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue to ensure it exists
	q, err := ch.QueueDeclare(
		"health_checks", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Start consuming messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack (message is considered delivered once received)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Notifier worker started. Waiting for messages...")

	// Create a channel to keep the application running
	forever := make(chan bool)

	// Process messages in a separate goroutine
	go func() {
		for d := range msgs {
			var result HealthCheckResult
			err := json.Unmarshal(d.Body, &result)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			// The core logic of our notifier!
			if !result.IsUp {
				log.Printf("🚨 ALERT: Site with ID %d is DOWN! Status Code: %d", result.SiteID, result.StatusCode)
			}
		}
	}()

	// Block forever
	<-forever
}
