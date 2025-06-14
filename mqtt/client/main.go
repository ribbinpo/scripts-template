package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ribbinpo/scripts-template/mqtt/client/config"
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

	client := config.NewMQTTConfig()

	token := client.Connect()
	defer client.Disconnect(250)

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Execute action based on flag
	switch *action {
	case "publish":
		payload := &PublisherPayload{
			Client:  client,
			Topic:   *topic,
			Message: *message,
		}
		Publisher(payload)
	case "subscribe":
		payload := &SubscriberPayload{
			Client: client,
			Topic:  *topic,
		}
		Subscriber(payload)
	default:
		fmt.Printf("Error: Invalid action '%s'. Must be 'publish' or 'subscribe'\n", *action)
		flag.Usage()
		os.Exit(1)
	}
}
