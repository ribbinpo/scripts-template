# Kafka Docker Setup

This repository contains Docker Compose configurations for running Apache Kafka in two different modes:

1. KRaft Mode (3-node cluster) - Modern, ZooKeeper-less setup
2. ZooKeeper Mode - Traditional setup with ZooKeeper

## Prerequisites

- Docker
- Docker Compose

## KRaft Mode (3-node cluster)

This setup runs Kafka in KRaft mode (Kafka Raft) without ZooKeeper, using a 3-node cluster for high availability.

### Configuration

- 3 Kafka brokers (kafka1, kafka2, kafka3)
- Each broker acts as both broker and controller
- Replication factor: 3
- Transaction state log replication factor: 3
- Minimum ISR: 2

### Ports

- kafka1: 
  - Internal: 9092
  - External: 29092
- kafka2: 9093
- kafka3: 9094

### Running the Cluster

```bash
cd kafka
docker-compose up -d
```

### Accessing Kafka

- From within Docker network: Use `kafka1:9092`, `kafka2:9092`, or `kafka3:9092`
- From host machine: Use `localhost:29092`

## ZooKeeper Mode

This setup runs Kafka with ZooKeeper, suitable for development and testing.

### Configuration

- Single Kafka broker
- Single ZooKeeper instance
- Replication factor: 1
- Transaction state log replication factor: 1

### Ports

- Kafka: 9092
- ZooKeeper: 2181

### Running with ZooKeeper

```bash
cd kafka
docker-compose -f docker-compose.zookeeper.yaml up -d
```

### Accessing Kafka

- From within Docker network: Use `kafka:9092`
- From host machine: Use `localhost:9092`

## Network

Both configurations use a bridge network named `kafka_network` for internal communication.

## Volumes

Each service uses a local volume for data persistence:
- KRaft mode: `kafka1_data`, `kafka2_data`, `kafka3_data`
- ZooKeeper mode: `kafka_data`, `zookeeper_data`

## Security

Both configurations use PLAINTEXT protocol for simplicity. For production use, consider implementing SSL/TLS and SASL authentication.

## Notes

- KRaft mode is recommended for production use as it's the future of Kafka
- ZooKeeper mode is simpler and suitable for development/testing
- Both configurations allow anonymous access (ALLOW_PLAINTEXT_LISTENER: yes)
- Auto topic creation is enabled by default 