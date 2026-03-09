package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ribbinpo/scripts-template/rabbitmq/client/util"
)

type PublisherPayload struct {
	Channel *amqp.Channel
	Topic   string
	Message string
}

// Method 1: direct exchange
func Publisher(payload *PublisherPayload) {
	exchange := util.GetExchangeName(payload.Topic, util.Events)
	err := payload.Channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}

	// Publish the message
	routingKey := util.GetQueueName(payload.Topic, "main", util.NormalQueue)
	err = payload.Channel.PublishWithContext(
		context.Background(),
		exchange, // exchange
		routingKey, // routing key
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
	exchange := util.GetExchangeName(payload.Topic, util.Events)
	err := payload.Channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		panic(err)
	}

	// Publish the message
	routingKey := util.GetQueueName(payload.Topic, "main", util.NormalQueue)
	err = payload.Channel.PublishWithContext(
		context.Background(),
		exchange,    // exchange
		routingKey,  // routing key
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
	exchange := util.GetExchangeName(payload.Topic, util.Events)
	err := payload.Channel.ExchangeDeclare(
		exchange, // name
		"fanout", // type
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
		exchange, // exchange
		"",       // routing key (empty for fanout)
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
