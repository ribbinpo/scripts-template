# Kafka Client

A Go-based Kafka client implementation that provides both low-level and high-level APIs for interacting with Apache Kafka.

## Features

- Message publishing (both low-level and high-level APIs)
- Message subscription (both low-level and high-level APIs)
- Consumer group support
- Manual message commit capability
- Topic listing functionality
- Retry mechanism for failed message publishing
- Support for multiple topics

## Prerequisites

- Go 1.x
- Apache Kafka running on localhost:29092 (or update the broker address in the code)

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

## Usage

The client can be used via command-line interface with the following flags:

```bash
-action string    Action to perform (publish/subscribe)
-topic string     Kafka topic
-message string   Message to publish (required for publish action)
```

### Examples

1. Publish a message:
```bash
go run . -action publish -topic my-topic -message "Hello, Kafka!"
```

2. Subscribe to a topic:
```bash
go run . -action subscribe -topic my-topic
```

3. List all topics:
```bash
go run .
```

## API Overview

### Publisher APIs

1. Low-level API (`Publisher`):
   - Direct connection to Kafka broker
   - Simple message publishing
   - No retry mechanism

2. High-level API (`ProduceMessage`):
   - Automatic retry mechanism (3 attempts)
   - Support for multiple messages
   - Automatic topic creation
   - Better error handling

### Subscriber APIs

1. Low-level API (`Subscriber`):
   - Direct connection to Kafka broker
   - Batch message reading
   - Basic error handling

2. High-level APIs:
   - `ConsumeMessage`: Basic message consumption
   - `ConsumeGroupMessage`: Consumer group support
   - `ConsumeMessageManual`: Manual message commit capability

## Configuration

The client is configured to connect to Kafka at `localhost:29092` by default. To change this:

1. Update the broker address in the relevant files:
   - `publisher.go`
   - `subscriber.go`
   - `main.go`

## Error Handling

- The publisher implements a retry mechanism for failed message publishing
- The subscriber includes basic error handling for connection and reading issues
- All errors are logged with appropriate context

## Contributing

Feel free to submit issues and enhancement requests.

## License

This project is licensed under the MIT License.
