package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type SubscriberPayload struct {
	Channel *kafka.Conn
	Topic   string
}

// Low level API
func Subscriber(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to topic: %s\n", payload.Topic)
	defer payload.Channel.Close()
	// payload.Channel.Seek(0, kafka.SeekStart)
	for {
		buf := make([]byte, 1e3)
		batch := payload.Channel.ReadBatch(10e3, 1e6)
		if batch == nil {
			log.Println("No batch available")
			time.Sleep(1 * time.Second)
			continue
		}
		for {
			n, err := batch.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("read error: %v", err)
				break
			}
			fmt.Printf("Message: %s\n", string(buf[:n]))
		}
		batch.Close()
	}
}

// High level API - Consume Message
func ConsumeMessage(payload *SubscriberPayload) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     payload.Topic,
		Partition: 0,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

// High level API - Consume Group Message
func ConsumeGroupMessage(payload *SubscriberPayload) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    payload.Topic,
		GroupID:  "my-group",
		MaxBytes: 10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

// High level API - Consume Message Manual
func ConsumeMessageManual(payload *SubscriberPayload) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    payload.Topic,
		GroupID:  "my-group",
		MaxBytes: 10e6,
	})
	defer r.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := r.FetchMessage(ctx)
			if err != nil {
				log.Printf("Error fetching message: %v", err)
				continue
			}

			// Process message
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
				m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

			if err := r.CommitMessages(ctx, m); err != nil {
				log.Printf("Error committing message: %v", err)
				continue
			}
		}
	}
}

// ConsumeWithDLQ starts the retry + DLQ consumer for the given topic.
// Uses queue naming: {topic}-retry, {topic}-dlq
func ConsumeWithDLQ(payload *SubscriberPayload) {
	cfg := DefaultRetryConfig()
	process := func(ctx context.Context, m kafka.Message) error {
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		// Return error to simulate failure and trigger retry/DLQ flow
		// Replace with your actual processing logic
		return nil
	}
	ConsumeWithRetryAndDLQ(payload.Topic, cfg, process)
}

// ConsumeDLQ reads messages from the Dead Letter Queue for inspection/reprocessing
func ConsumeDLQ(payload *SubscriberPayload) {
	_, dlqTopic := QueueNames(payload.Topic)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    dlqTopic,
		GroupID:  "dlq-inspector",
		MaxBytes: 10e6,
	})
	defer r.Close()

	fmt.Printf("Reading from DLQ: %s\n", dlqTopic)
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("DLQ read error: %v", err)
			break
		}
		fmt.Printf("[DLQ] offset %d | key: %s | value: %s | headers: %v\n",
			m.Offset, string(m.Key), string(m.Value), m.Headers)
	}
}
