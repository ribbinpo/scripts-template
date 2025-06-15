package main

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

func ListTopic() {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err)
	}

	m := map[string]struct{}{}

	for _, partition := range partitions {
		m[partition.Topic] = struct{}{}
	}

	for topic := range m {
		fmt.Println(topic)
	}
}
