package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"gopkg.in/yaml.v2"
)

// Define a struct to match the YAML structure
type OrderMessage struct {
	OrderID      string  `yaml:"order_id"`
	CustomerName string  `yaml:"customer_name"`
	Amount       float64 `yaml:"amount"`
	Status       string  `yaml:"status"`
}

func main() {
	// Fetch projectID and topicID from environment variables
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GCP_PROJECT_ID environment variable is not set")
	}
	topicID := os.Getenv("PUBSUB_TOPIC_ID")
	if topicID == "" {
		log.Fatal("PUBSUB_TOPIC_ID environment variable is not set")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Read the YAML file
	yamlFile, err := os.ReadFile("order_message.yaml")
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	// Parse the YAML file into the OrderMessage struct
	var order OrderMessage
	err = yaml.Unmarshal(yamlFile, &order)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML file: %v", err)
	}

	// Convert the struct to JSON format for the Pub/Sub message
	jsonData, err := yaml.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order data to JSON: %v", err)
	}

	// Create Pub/Sub client
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Printf("Failed to close Pub/Sub client: %v", err)
		}
	}()

	// Get the topic reference
	topic := client.Topic(topicID)

	// Create a Pub/Sub message
	msg := &pubsub.Message{
		Data: jsonData,
	}

	// Publish the message
	result := topic.Publish(ctx, msg)

	// Wait for the publish to complete
	messageID, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Message published successfully with ID: %s", messageID)
}
