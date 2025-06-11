# RabbitMQ Docker Setup

This directory contains a Docker Compose configuration for running RabbitMQ with management interface.

## Configuration

The setup uses RabbitMQ version 4.1.1 with the management plugin enabled, running on Alpine Linux for a smaller footprint.

### Default Credentials
- Username: `myuser`
- Password: `mypassword`

### Ports
- `5672`: AMQP protocol port (for message publishing/consuming)
- `15672`: Management UI port (web interface)

## Getting Started

1. Start the RabbitMQ container:
```bash
docker-compose up -d
```

2. Access the Management UI:
   - Open your browser and navigate to `http://localhost:15672`
   - Login using the default credentials above

3. Stop the container:
```bash
docker-compose down
```

## Data Persistence

The RabbitMQ data is persisted in the `./rabbitmq_data` directory, which is mounted as a volume in the container. This ensures that your queues, exchanges, and messages are preserved even if the container is stopped or removed.

## Environment Variables

The following environment variables are configured:
- `RABBITMQ_DEFAULT_USER`: Sets the default admin username
- `RABBITMQ_DEFAULT_PASS`: Sets the default admin password

## Security Notes

- The default credentials are for development purposes only
- For production use, please change the default credentials
- Consider using environment variables or Docker secrets for sensitive information

## Additional Resources

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [RabbitMQ Management Plugin](https://www.rabbitmq.com/management.html)
- [Docker Hub RabbitMQ Image](https://hub.docker.com/_/rabbitmq)
