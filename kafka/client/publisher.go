package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type PublisherPayload struct {
	Channel *kafka.Conn
	Message string
	Topic   string
}

// Low level API
func Publisher(payload *PublisherPayload) error {
	msg := kafka.Message{Value: []byte(payload.Message)}
	if _, err := payload.Channel.WriteMessages(msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	fmt.Printf("Published message to topic: %s | message: %s\n", payload.Topic, payload.Message)
	return nil
}

// High level API
func ProduceMessage(payload *PublisherPayload) error {
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:                  payload.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	defer w.Close()

	messages := []kafka.Message{
		{Value: []byte(payload.Message)},
		{Value: []byte("Hello World"), Topic: "topic-B"}, // specify topic in the message
	}

	// if err := w.WriteMessages(context.Background(), messages...); err != nil {
	// 	panic(err)
	// }

	// Retry if the topic is not available
	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = w.WriteMessages(ctx, messages...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			return fmt.Errorf("failed to publish message after %d retries: %w", i+1, err)
		}

		fmt.Printf("Published message to topic: %s | message: %s\n", payload.Topic, payload.Message)
		return nil
	}

	return fmt.Errorf("failed to publish message after %d retries: %w", retries, err)
}
