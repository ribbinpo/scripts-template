package main

import (
	"fmt"

	"github.com/wagslane/go-rabbitmq"
)

type PublisherPayload struct {
	Client  *rabbitmq.Conn
	Topic   string
	Message string
}

func Publisher(payload *PublisherPayload) {
	// Method 1: direct exchange
	publisher, err := rabbitmq.NewPublisher(payload.Client,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(payload.Topic),
		rabbitmq.WithPublisherOptionsExchangeKind("direct"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)

	// Method 2: fanout exchange
	// publisher, err = rabbitmq.NewPublisher(payload.Client,
	// 	rabbitmq.WithPublisherOptionsLogging,
	// 	rabbitmq.WithPublisherOptionsExchangeName(payload.Topic),
	// 	rabbitmq.WithPublisherOptionsExchangeKind("fanout"),
	// )

	// Method 3: topic exchange
	// publisher, err = rabbitmq.NewPublisher(payload.Client,
	// 	rabbitmq.WithPublisherOptionsLogging,
	// 	rabbitmq.WithPublisherOptionsExchangeName(payload.Topic),
	// 	rabbitmq.WithPublisherOptionsExchangeDeclare,
	// 	rabbitmq.WithPublisherOptionsExchangeKind("topic"),
	// )

	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	if err := publisher.Publish([]byte(payload.Message), []string{payload.Topic}, rabbitmq.WithPublishOptionsContentType("application/json"), rabbitmq.WithPublishOptionsExchange("test")); err != nil {
		panic(err)
	}

	fmt.Printf("Published message to topic: %s | message: %s\n", payload.Topic, payload.Message)
}
