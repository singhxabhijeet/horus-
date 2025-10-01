package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// HealthCheckResult is the structure of the message we expect to receive
type HealthCheckResult struct {
	SiteID         int    `json:"SiteID"`
	URL            string `json:"URL"` // Let's assume the URL is now in the message
	IsUp           bool   `json:"IsUp"`
	StatusCode     int    `json:"StatusCode"`
	ResponseTimeMs int64  `json:"ResponseTimeMs"`
}

// Discord message structures for a rich embed
type DiscordEmbedField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type DiscordEmbed struct {
	Title  string              `json:"title"`
	Color  int                 `json:"color"`
	Fields []DiscordEmbedField `json:"fields"`
}
type DiscordWebhookPayload struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

func main() {
	// Get environment variables
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		log.Fatal("DISCORD_WEBHOOK_URL is not set")
	}

	// ... (RabbitMQ connection logic is exactly the same as before)
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
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare("health_checks", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Notifier worker started. Waiting for messages...")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var result HealthCheckResult
			err := json.Unmarshal(d.Body, &result)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			// If the site is down, send a Discord notification
			if !result.IsUp {
				log.Printf("🚨 ALERT: Site with ID %d is DOWN! Status Code: %d. Sending notification...", result.SiteID, result.StatusCode)
				sendDiscordNotification(webhookURL, result)
			}
		}
	}()
	<-forever
}

func sendDiscordNotification(webhookURL string, result HealthCheckResult) {
	// Construct the rich embed message
	payload := DiscordWebhookPayload{
		Embeds: []DiscordEmbed{
			{
				Title: "🚨 Website Downtime Alert 🚨",
				Color: 15158332, // Red color
				Fields: []DiscordEmbedField{
					{Name: "Site ID", Value: strconv.Itoa(result.SiteID)},
					{Name: "Status", Value: "DOWN"},
					{Name: "HTTP Status Code", Value: strconv.Itoa(result.StatusCode)},
					{Name: "Checked At", Value: time.Now().Format(time.RFC1123)},
				},
			},
		},
	}

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling Discord payload: %v", err)
		return
	}

	// Send the HTTP POST request to the webhook URL
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("Error creating Discord request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending Discord notification: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("Discord returned non-2xx status: %d", resp.StatusCode)
	} else {
		log.Printf("Successfully sent Discord notification for site ID %d.", result.SiteID)
	}
}
