# MQTT Client

A simple MQTT client implementation in Go that supports publishing and subscribing to MQTT topics.

## Prerequisites

- Go 1.16 or higher
- MQTT broker (e.g., Mosquitto, EMQ X, or any other MQTT broker)

## Installation

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```

## Usage

The client supports two main actions: publishing and subscribing to MQTT topics.

### Publishing Messages

To publish a message to a topic:

```bash
go run . -action="publish" -message="Hello" -topic="test"
```

Parameters:
- `-action`: Must be "publish"
- `-message`: The message to publish
- `-topic`: The MQTT topic to publish to

### Subscribing to Topics

To subscribe to a topic:

```bash
go run . -action="subscribe" -topic="test" -node=5
```

Parameters:
- `-action`: Must be "subscribe"
- `-topic`: The MQTT topic to subscribe to
- `-node`: (Optional) Node identifier

## Features

- Simple command-line interface
- Support for publishing messages
- Support for subscribing to topics
- Automatic reconnection handling
- QoS 0 message delivery

## Configuration

The MQTT client configuration is managed through the `config` package. You can modify the connection settings in the configuration file as needed.

## Error Handling

The client includes basic error handling for:
- Missing required parameters
- Connection failures
- Publishing errors
- Subscription errors