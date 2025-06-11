package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PublisherPayload struct {
	Channel *amqp.Channel
	Topic   string
	Message string
}

// Method 1: direct exchange
func Publisher(payload *PublisherPayload) {
	err := payload.Channel.ExchangeDeclare(
		payload.Topic, // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}

	// Publish the message
	err = payload.Channel.PublishWithContext(
		context.Background(),
		payload.Topic, // exchange
		payload.Topic, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(payload.Message),
			DeliveryMode: amqp.Persistent, // 0: transient, 1: persistent
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Published message to topic: %s | message: %s\n", payload.Topic, payload.Message)
}

// Method 2: topic exchange
func PublisherTopic(payload *PublisherPayload) {
	err := payload.Channel.ExchangeDeclare(
		payload.Topic, // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}

	// Publish the message
	err = payload.Channel.PublishWithContext(
		context.Background(),
		payload.Topic, // exchange
		payload.Topic, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(payload.Message),
			DeliveryMode: amqp.Persistent, // 0: transient, 1: persistent
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Published message to topic exchange: %s | message: %s\n", payload.Topic, payload.Message)
}

// Method 3: fanout exchange
func PublisherFanout(payload *PublisherPayload) {
	err := payload.Channel.ExchangeDeclare(
		payload.Topic, // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		panic(err)
	}

	// Publish the message
	err = payload.Channel.PublishWithContext(
		context.Background(),
		payload.Topic, // exchange
		"",            // routing key (empty for fanout)
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(payload.Message),
			DeliveryMode: amqp.Persistent, // 0: transient, 1: persistent
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Published message to fanout exchange: %s | message: %s\n", payload.Topic, payload.Message)
}
