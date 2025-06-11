package main

import (
	"fmt"
	"log"

	"github.com/wagslane/go-rabbitmq"
)

type SubscriberPayload struct {
	Client *rabbitmq.Conn
	Topic  string
}

func Subscriber(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to topic: %s\n", payload.Topic)
	// Method 1: direct exchange
	consumer, err := rabbitmq.NewConsumer(payload.Client,
		payload.Topic+"_queue",
		rabbitmq.WithConsumerOptionsRoutingKey(payload.Topic),
		rabbitmq.WithConsumerOptionsExchangeName(payload.Topic),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	// Method 2: fanout exchange
	// consumer, err = rabbitmq.NewConsumer(payload.Client,
	// 	payload.Topic+"_queue",
	// 	rabbitmq.WithConsumerOptionsExchangeName(payload.Topic),
	// 	rabbitmq.WithConsumerOptionsExchangeDeclare,
	// 	rabbitmq.WithConsumerOptionsExchangeKind("fanout"),
	// )

	// Method 3: topic exchange
	// - * means any word
	// - # means any word and any number of words
	// consumer, err = rabbitmq.NewConsumer(payload.Client,
	// 	payload.Topic+"_queue",
	// 	rabbitmq.WithConsumerOptionsRoutingKey(payload.Topic+"topic.exchange.*"),
	// 	rabbitmq.WithConsumerOptionsExchangeName(payload.Topic),
	// 	rabbitmq.WithConsumerOptionsExchangeDeclare,
	// 	rabbitmq.WithConsumerOptionsExchangeKind("topic"),
	// )

	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	if err := consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("consumed: %v", string(d.Body))
		// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		return rabbitmq.Ack
	}); err != nil {
		panic(err)
	}
}
