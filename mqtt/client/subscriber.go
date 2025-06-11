package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SubscriberPayload struct {
	Client mqtt.Client
	Topic  string
}

func Subscriber(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to topic: %s\n", payload.Topic)
	token := payload.Client.Subscribe(payload.Topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message on topic: %s | message: %s\n", msg.Topic(), string(msg.Payload()))
	})

	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	select {}
}
