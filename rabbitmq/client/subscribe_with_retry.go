package main

import (
	"context"
	"fmt"
	"log"
	"math"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ribbinpo/scripts-template/rabbitmq/client/util"
)

func ProcessWithRetry(service string, d amqp.Delivery, ch *amqp.Channel) {
	maxRetries := 3

	retries := getRetryCount(d)

	err := doWork(d.Body)

	if err == nil {
		d.Ack(false)
		return
	}

	if retries >= maxRetries {
		totalAttempts := retries + 1
		log.Printf("Moving to DLQ after %d attempts", totalAttempts)
		dlxExchange := util.GetExchangeName(service, util.DLX)
		dlqName := util.GetQueueName(service, "main", util.DLQ)
		if pubErr := ch.PublishWithContext(
			context.Background(),
			dlxExchange,
			dlqName,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        d.Body,
				Headers:     d.Headers,
			},
		); pubErr != nil {
			log.Printf("Failed to publish to DLQ: %v; nacking for requeue", pubErr)
			_ = d.Nack(false, true)
			return
		}
	} else {
		delayDuration := int64(math.Pow(float64(retries+1), 2)) * 1000

		log.Printf("Retrying message after (%d/%d) attempts in %dms", retries+1, maxRetries+1, delayDuration)

		retryExchange := util.GetExchangeName(service, util.Retry)
		if pubErr := ch.PublishWithContext(
			context.Background(),
			retryExchange,
			d.RoutingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        d.Body,
				Headers:     d.Headers,
				Expiration:  fmt.Sprintf("%d", delayDuration),
			},
		); pubErr != nil {
			log.Printf("Failed to publish to retry exchange: %v; nacking for requeue", pubErr)
			_ = d.Nack(false, true)
			return
		}
	}
	d.Ack(false)
}

func doWork(body []byte) error {
	return nil
}

func getRetryCount(d amqp.Delivery) int {
	// x-death can be []interface{} (array of tables) or amqp.Table (single entry)
	var entries []interface{}
	switch v := d.Headers["x-death"].(type) {
	case []interface{}:
		entries = v
	case amqp.Table:
		entries = []interface{}{v}
	default:
		return 0
	}

	if len(entries) == 0 {
		return 0
	}

	// The first entry is usually the most recent "death" event
	deathInfo, ok := entries[0].(amqp.Table)
	if !ok {
		return 0
	}

	// count can be int64, int, uint64 depending on RabbitMQ/AMQP version
	switch v := deathInfo["count"].(type) {
	case int64:
		return int(v)
	case int:
		return v
	case int32:
		return int(v)
	case uint64:
		return int(v)
	case uint32:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}
