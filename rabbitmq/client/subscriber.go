package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ribbinpo/scripts-template/rabbitmq/client/util"
)

// Main Queue -> DLX -> DLQ

type SubscriberPayload struct {
	Channel *amqp.Channel
	Topic   string
}

// Method 1: direct exchange
func Subscriber(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to topic: %s\n", payload.Topic)

	dlxExchange := util.GetExchangeName(payload.Topic, util.DLX)
	dlqName := util.GetQueueName(payload.Topic, "main", util.DLQ)
	mainExchange := util.GetExchangeName(payload.Topic, util.Events)
	mainQueueName := util.GetQueueName(payload.Topic, "main", util.NormalQueue)
	routingKey := mainQueueName

	// Declare the DLX exchange
	if err := payload.Channel.ExchangeDeclare(
		dlxExchange, // name
		"direct",    // type
		true,        // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	); err != nil {
		panic(err)
	}

	// Declare the DLQ queue
	if _, err := payload.Channel.QueueDeclare(
		dlqName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	); err != nil {
		panic(err)
	}

	// Bind the DLQ queue to the DLX exchange
	if err := payload.Channel.QueueBind(
		dlqName,     // queue name
		dlqName,     // routing key
		dlxExchange, // exchange
		false,       // no-wait
		nil,         // arguments
	); err != nil {
		panic(err)
	}

	// Declare the main exchange
	if err := payload.Channel.ExchangeDeclare(
		mainExchange, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		panic(err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange":    dlxExchange,
		"x-dead-letter-routing-key": dlqName,
	}

	// Declare the main queue
	q, err := payload.Channel.QueueDeclare(
		mainQueueName, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		args,          // arguments
	)
	if err != nil {
		panic(err)
	}

	// Bind the queue to the exchange
	if err := payload.Channel.QueueBind(
		q.Name,        // queue name
		routingKey,    // routing key
		mainExchange,  // exchange
		false,         // no-wait
		nil,           // arguments
	); err != nil {
		panic(err)
	}

	// Consume messages
	msgs, err := payload.Channel.Consume(
		q.Name,     // queue
		"worker-1", // consumer tag
		false,      // auto-ack = false → manual ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("consumed: %v", string(d.Body))
			d.Ack(false)
			// d.Nack(false, true) // requeue = true
			// d.Reject(true) // requeue = true
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Method 2: topic exchange
func SubscriberTopic(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to topic exchange: %s\n", payload.Topic)

	exchange := util.GetExchangeName(payload.Topic, util.Events)
	queueName := util.GetQueueName(payload.Topic, "main", util.NormalQueue)
	routingKey := queueName

	// Declare the exchange
	if err := payload.Channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		panic(err)
	}

	// Declare the queue
	q, err := payload.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}

	// Bind the queue to the exchange
	if err := payload.Channel.QueueBind(
		q.Name,      // queue name
		routingKey,  // routing key
		exchange,    // exchange
		false,       // no-wait
		nil,         // arguments
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
			// d.Nack(false, true) // requeue = true
			// d.Reject(true) // requeue = true
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// Method 3: fanout exchange
func SubscriberFanout(payload *SubscriberPayload) {
	fmt.Printf("Subscribing to fanout exchange: %s\n", payload.Topic)

	exchange := util.GetExchangeName(payload.Topic, util.Events)

	// Declare the exchange
	if err := payload.Channel.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
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
		q.Name,   // queue name
		"",       // routing key (empty for fanout)
		exchange, // exchange
		false,    // no-wait
		nil,      // arguments
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
			// d.Nack(false, true) // requeue = true
			// d.Reject(true) // requeue = true
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
