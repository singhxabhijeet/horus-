package internal

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Publisher handles sending messages to RabbitMQ
type Publisher struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

// NewPublisher creates a new publisher and connects to RabbitMQ
func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare a queue to make sure it exists.
	// This is idempotent - it will only be created if it doesn't exist.
	_, err = ch.QueueDeclare(
		"health_checks", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: conn, ch: ch}, nil
}

// Publish sends a health check result to the queue
func (p *Publisher) Publish(result HealthCheckResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.ch.PublishWithContext(ctx,
		"",              // exchange
		"health_checks", // routing key (queue name)
		false,           // mandatory
		false,           // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		return err
	}

	log.Printf(" [x] Sent message for site ID %d\n", result.SiteID)
	return nil
}

// Close closes the connection and channel
func (p *Publisher) Close() {
	p.ch.Close()
	p.conn.Close()
}
