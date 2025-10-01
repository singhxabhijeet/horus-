package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct holds all configuration for the application
type Config struct {
	DatabaseURL string
	RabbitMQURL string
}

// LoadConfig reads configuration from environment variables
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	return Config{
		DatabaseURL: dbURL,
		RabbitMQURL: rabbitURL,
	}
}
