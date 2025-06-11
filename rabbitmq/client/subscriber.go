package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SubscriberPayload struct {
	Channel *amqp.Channel
	Topic   string
}

// Method 1: direct exchange
func Subscriber(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to topic: %s\n", payload.Topic)

	// Declare the exchange
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

	// Declare the queue
	q, err := payload.Channel.QueueDeclare(
		payload.Topic+"_queue", // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		panic(err)
	}

	// Bind the queue to the exchange
	if err := payload.Channel.QueueBind(
		q.Name,        // queue name
		payload.Topic, // routing key
		payload.Topic, // exchange
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		panic(err)
	}

	// Consume messages
	msgs, err := payload.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack = false → manual ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("consumed: %v", string(d.Body))
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Method 2: topic exchange
func SubscriberTopic(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to topic exchange: %s\n", payload.Topic)

	// Declare the exchange
	if err := payload.Channel.ExchangeDeclare(
		payload.Topic, // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		panic(err)
	}

	// Declare the queue
	q, err := payload.Channel.QueueDeclare(
		payload.Topic+"_queue", // name
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		panic(err)
	}

	// Bind the queue to the exchange
	if err := payload.Channel.QueueBind(
		q.Name,        // queue name
		payload.Topic, // routing key
		payload.Topic, // exchange
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		panic(err)
	}

	// Consume messages
	msgs, err := payload.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack = false → manual ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("consumed from topic exchange: %v", string(d.Body))
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Method 3: fanout exchange
func SubscriberFanout(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to fanout exchange: %s\n", payload.Topic)

	// Declare the exchange
	if err := payload.Channel.ExchangeDeclare(
		payload.Topic, // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		panic(err)
	}

	// Declare the queue
	q, err := payload.Channel.QueueDeclare(
		"",    // name (empty for fanout to get a random queue name)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	// Bind the queue to the exchange
	if err := payload.Channel.QueueBind(
		q.Name,        // queue name
		"",            // routing key (empty for fanout)
		payload.Topic, // exchange
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		panic(err)
	}

	// Consume messages
	msgs, err := payload.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack = false → manual ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("consumed from fanout exchange: %v", string(d.Body))
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
