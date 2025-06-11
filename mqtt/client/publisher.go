package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PublisherPayload struct {
	Client  mqtt.Client
	Topic   string
	Message string
}

func Publisher(payload *PublisherPayload) {
	token := payload.Client.Publish(payload.Topic, 0, false, payload.Message)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("Published message to topic: %s | message: %s\n", payload.Topic, payload.Message)
}
