# RabbitMQ Client

A simple command-line RabbitMQ client written in Go that supports publishing and subscribing to messages using different exchange types.

## Features

- Support for multiple exchange types:
  - Direct Exchange
  - Fanout Exchange
  - Topic Exchange
- Simple command-line interface
- Configurable connection settings
- Support for publishing and subscribing to messages

## Prerequisites

- Go 1.16 or higher
- RabbitMQ server running (default: localhost:5672)

## Installation

1. Clone the repository
2. Install dependencies:

```bash
go mod download
```

## Usage

The client supports two main actions: `publish` and `subscribe`.

### Publishing Messages

```bash
go run . -action publish -topic <topic_name> -message <message>
```

Example:

```bash
go run . -action publish -topic test -message "Hello, RabbitMQ!"
```

### Subscribing to Messages

```bash
go run . -action subscribe -topic <topic_name>
```

Example:

```bash
go run . -action subscribe -topic test
```

## Exchange Types

The client supports three types of exchanges:

1. **Direct Exchange** (Default)

   - Messages are routed based on exact routing key matches
   - Default configuration in the code

2. **Fanout Exchange**

   - Messages are broadcast to all queues bound to the exchange
   - Uncomment the fanout exchange configuration in both publisher.go and subscriber.go

3. **Topic Exchange**
   - Messages are routed based on pattern matching
   - Supports wildcards:
     - `*` (star) matches exactly one word
     - `#` (hash) matches zero or more words
   - Uncomment the topic exchange configuration in both publisher.go and subscriber.go

## Configuration

The default RabbitMQ connection settings are:

- Host: localhost
- Port: 5672
- Username: myuser
- Password: mypassword

To modify these settings, update the connection string in `main.go`.

## Examples

### Direct Exchange

```bash
# Terminal 1 - Subscribe
go run main.go -action subscribe -topic test

# Terminal 2 - Publish
go run main.go -action publish -topic test -message "Hello, Direct Exchange!"
```

### Fanout Exchange

1. Uncomment the fanout exchange configuration in both publisher.go and subscriber.go
2. Run multiple subscribers:

```bash
# Terminal 1
go run main.go -action subscribe -topic fanout_test

# Terminal 2
go run main.go -action subscribe -topic fanout_test

# Terminal 3 - Publish (all subscribers will receive the message)
go run main.go -action publish -topic fanout_test -message "Hello, Fanout Exchange!"
```

### Topic Exchange

1. Uncomment the topic exchange configuration in both publisher.go and subscriber.go
2. Run subscribers with different routing patterns:

```bash
# Terminal 1 - Subscribe to all messages
go run main.go -action subscribe -topic topic_test

# Terminal 2 - Publish
go run main.go -action publish -topic topic_test -message "Hello, Topic Exchange!"
```

## License

MIT License
