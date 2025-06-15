package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Define command line flags
	action := flag.String("action", "", "Action to perform (publish/subscribe)")
	topic := flag.String("topic", "", "MQTT topic")
	message := flag.String("message", "", "Message to publish")

	// Parse command line flags
	flag.Parse()

	// Validate required flags
	if *action == "" {
		fmt.Println("Error: -action flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if *topic == "" {
		fmt.Println("Error: -topic flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if *action == "publish" && *message == "" {
		fmt.Println("Error: -message flag is required for publish action")
		flag.Usage()
		os.Exit(1)
	}

	// Connect to Kafka
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", *topic, 0)
	if err != nil {
		panic(err)
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	defer conn.Close()

	// Execute action based on flag
	switch *action {
	case "publish":
		payload := &PublisherPayload{
			Channel: conn,
			Message: *message,
			Topic:   *topic,
		}
		Publisher(payload)
	case "subscribe":
		payload := &SubscriberPayload{
			Channel: conn,
			Topic:   *topic,
		}
		Subscriber(payload)
	default:
		fmt.Printf("Error: Invalid action '%s'. Must be 'publish' or 'subscribe'\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}
