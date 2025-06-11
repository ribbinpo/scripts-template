package main

import (
	"flag"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
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

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://myuser:mypassword@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Set QoS to 1 to prefetch 1 message
	if err := ch.Qos(1, 0, false); err != nil {
		panic(err)
	}

	// Execute action based on flag
	switch *action {
	case "publish":
		payload := &PublisherPayload{
			Channel: ch,
			Topic:   *topic,
			Message: *message,
		}
		Publisher(payload)
	case "subscribe":
		payload := &SubscriberPayload{
			Channel: ch,
			Topic:   *topic,
		}
		Subscriber(payload)
	default:
		fmt.Printf("Error: Invalid action '%s'. Must be 'publish' or 'subscribe'\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}
