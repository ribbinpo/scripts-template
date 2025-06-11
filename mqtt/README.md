# MQTT Broker Setup with Mosquitto

This repository contains a Docker-based MQTT broker setup using Eclipse Mosquitto. The setup includes authentication, ACL (Access Control List), and supports both WebSocket and MQTT protocols.

## Features

- MQTT over WebSockets (port 9001)
- Authentication enabled with username/password
- Access Control List (ACL) for topic-based permissions
- Persistent storage for messages
- Logging configuration
- TLS support (configured but disabled by default)

## Prerequisites

- Docker
- Docker Compose

## Directory Structure

```
mqtt/
├── docker-compose.yaml
├── mosquitto/
│   ├── mosquitto.conf
│   ├── acl_file
│   ├── data/
│   ├── log/
│   └── init/
│       └── password_file
```

## Default Configuration

- **Username**: myuser
- **Password**: mypassword
- **WebSocket Port**: 9001
- **MQTT Port**: 1883 (enabled by default)
- **TLS Port**: 8883 (disabled by default)

## Getting Started

1. Clone this repository
2. Navigate to the project directory
3. Start the MQTT broker:
   ```bash
   docker-compose up -d
   ```

## Access Control

The ACL file (`mosquitto/acl_file`) defines the following permissions:

- User `myuser` has read/write access to:
  - `home/myuser/#`
  - `test`
- All users have read access to `$SYS/#` topics (system monitoring)

- Change ACL permission file

```sh
chmod 0700 ./acl_file
```

## Connecting to the Broker

### WebSocket Connection

Connect to the broker using WebSocket on port 9001:

```javascript
// Example using MQTT.js
const client = mqtt.connect("ws://localhost:9001", {
  username: "myuser",
  password: "mypassword",
});
```

## Security Notes

- Authentication is enabled by default
- Anonymous access is disabled
- TLS is configured but disabled by default
- Passwords are stored in a password file

## Logging

Logs are stored in the `mosquitto/log` directory. The configuration includes:

- Error logs
- Warning logs
- Notice logs
- Information logs

## Persistence

Message persistence is enabled and stored in the `mosquitto/data` directory.

## Customization

To modify the configuration:

1. Edit `mosquitto/mosquitto.conf` for broker settings
2. Edit `mosquitto/acl_file` for access control rules
3. Update the password using the `init-password` service in docker-compose.yaml

## Stopping the Broker

```bash
docker-compose down
```

## License

This project is licensed under the MIT License.
