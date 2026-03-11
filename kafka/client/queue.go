package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// Queue naming conventions
const (
	DLQSuffix   = "-dlq"
	RetrySuffix = "-retry"
)

// RetryConfig holds configuration for retry and DLQ behavior
type RetryConfig struct {
	MaxRetries     int           // Max retry attempts before sending to DLQ
	RetryDelay     time.Duration // Delay before reprocessing from retry queue
	Brokers        []string
	ConsumerGroupID string
}

// DefaultRetryConfig returns a sensible default configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:      3,
		RetryDelay:      5 * time.Second,
		Brokers:         []string{"localhost:29092"},
		ConsumerGroupID: "my-group",
	}
}

// QueueNames derives retry and DLQ topic names from the main topic
func QueueNames(mainTopic string) (retryTopic, dlqTopic string) {
	return mainTopic + RetrySuffix, mainTopic + DLQSuffix
}

// MessageProcessor processes a message. Return error to trigger retry/DLQ flow.
type MessageProcessor func(ctx context.Context, m kafka.Message) error

// Header keys for retry metadata
const (
	HeaderRetryAttempt = "x-retry-attempt"
	HeaderOriginalTopic = "x-original-topic"
	HeaderErrorMessage = "x-error-message"
)

// PublishToRetryQueue sends a failed message to the retry queue with attempt metadata
func PublishToRetryQueue(ctx context.Context, cfg RetryConfig, m kafka.Message, attempt int, errMsg string) error {
	retryTopic, _ := QueueNames(m.Topic)
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers[0]),
		Topic:    retryTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	headers := append(m.Headers,
		kafka.Header{Key: HeaderRetryAttempt, Value: []byte(strconv.Itoa(attempt))},
		kafka.Header{Key: HeaderOriginalTopic, Value: []byte(m.Topic)},
		kafka.Header{Key: HeaderErrorMessage, Value: []byte(errMsg)},
	)

	msg := kafka.Message{
		Key:     m.Key,
		Value:   m.Value,
		Headers: headers,
	}
	return w.WriteMessages(ctx, msg)
}

// PublishToDLQ sends a message to the Dead Letter Queue after max retries exhausted
func PublishToDLQ(ctx context.Context, cfg RetryConfig, m kafka.Message, errMsg string) error {
	_, dlqTopic := QueueNames(m.Topic)
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers[0]),
		Topic:    dlqTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	headers := append(m.Headers,
		kafka.Header{Key: HeaderOriginalTopic, Value: []byte(m.Topic)},
		kafka.Header{Key: HeaderErrorMessage, Value: []byte(errMsg)},
	)

	msg := kafka.Message{
		Key:     m.Key,
		Value:   m.Value,
		Headers: headers,
	}
	return w.WriteMessages(ctx, msg)
}

// GetRetryAttempt extracts retry attempt from message headers (0 if not present)
func GetRetryAttempt(m kafka.Message) int {
	for _, h := range m.Headers {
		if h.Key == HeaderRetryAttempt {
			if n, err := strconv.Atoi(string(h.Value)); err == nil {
				return n
			}
			return 0
		}
	}
	return 0
}

// ConsumeWithRetryAndDLQ consumes from main topic, processes messages, and routes
// failed messages to retry queue (with attempt limit) or DLQ when max retries exceeded.
func ConsumeWithRetryAndDLQ(mainTopic string, cfg RetryConfig, process MessageProcessor) {
	retryTopic, dlqTopic := QueueNames(mainTopic)
	log.Printf("Consuming from %s | retry: %s | dlq: %s", mainTopic, retryTopic, dlqTopic)

	readerConfig := func(topic string) kafka.ReaderConfig {
		return kafka.ReaderConfig{
			Brokers:  cfg.Brokers,
			Topic:    topic,
			GroupID:  cfg.ConsumerGroupID,
			MaxBytes: 10e6,
		}
	}

	// Consumer for main topic
	go consumeTopic(readerConfig(mainTopic), cfg, process, mainTopic)

	// Consumer for retry topic (delayed reprocessing)
	go consumeRetryTopic(readerConfig(retryTopic), cfg, process, mainTopic)

	select {} // block forever
}

func consumeTopic(rConfig kafka.ReaderConfig, cfg RetryConfig, process MessageProcessor, mainTopic string) {
	r := kafka.NewReader(rConfig)
	defer r.Close()

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Fetch error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if err := processAndCommit(context.Background(), r, m, cfg, process, 0); err != nil {
			log.Printf("Process error: %v", err)
		}
	}
}

func consumeRetryTopic(rConfig kafka.ReaderConfig, cfg RetryConfig, process MessageProcessor, mainTopic string) {
	r := kafka.NewReader(rConfig)
	defer r.Close()

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Printf("Retry fetch error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		attempt := GetRetryAttempt(m)
		time.Sleep(cfg.RetryDelay) // delay before reprocessing

		if err := processAndCommit(context.Background(), r, m, cfg, process, attempt); err != nil {
			log.Printf("Retry process error: %v", err)
		}
	}
}

func processAndCommit(ctx context.Context, r *kafka.Reader, m kafka.Message, cfg RetryConfig, process MessageProcessor, currentAttempt int) error {
	err := process(ctx, m)
	if err == nil {
		if commitErr := r.CommitMessages(ctx, m); commitErr != nil {
			return fmt.Errorf("commit: %w", commitErr)
		}
		fmt.Printf("Processed message at %s/%d/%d\n", m.Topic, m.Partition, m.Offset)
		return nil
	}

	nextAttempt := currentAttempt + 1
	errMsg := err.Error()

	if nextAttempt <= cfg.MaxRetries {
		// Send to retry queue
		if pubErr := PublishToRetryQueue(ctx, cfg, m, nextAttempt, errMsg); pubErr != nil {
			return fmt.Errorf("publish to retry: %w", pubErr)
		}
		if commitErr := r.CommitMessages(ctx, m); commitErr != nil {
			return fmt.Errorf("commit: %w", commitErr)
		}
		log.Printf("Message sent to retry queue (attempt %d/%d): %v", nextAttempt, cfg.MaxRetries, err)
	} else {
		// Send to DLQ
		if pubErr := PublishToDLQ(ctx, cfg, m, errMsg); pubErr != nil {
			return fmt.Errorf("publish to dlq: %w", pubErr)
		}
		if commitErr := r.CommitMessages(ctx, m); commitErr != nil {
			return fmt.Errorf("commit: %w", commitErr)
		}
		log.Printf("Message sent to DLQ after %d retries: %v", cfg.MaxRetries, err)
	}
	return nil
}
