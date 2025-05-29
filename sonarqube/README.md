# SonarQube Docker Setup

This scripts contains a Docker-based setup for SonarQube with PostgreSQL database. This setup allows you to run code quality analysis on your projects.

## Prerequisites

- Docker
- Docker Compose
- Make (optional, for using Makefile commands)

## Setup Instructions

1. Create a `.env` file in the sonarqube directory with the following `.env.example`:
```env
SONARQUBE_URL=
SONAR_TOKEN=
```

2. Start the SonarQube server and PostgreSQL database:
```bash
make up
```
or
```bash
docker-compose up -d
```

3. Wait for SonarQube to start (this may take a few minutes). You can access the SonarQube web interface at:
```
http://localhost:9000
```

4. On first login, use the default credentials:
- Username: admin
- Password: admin

## Usage

### Running Code Analysis

To analyze your code, use the following command:

```bash
make scan-cli input_dir=/path/to/your/code
```

Replace `/path/to/your/code` with the actual path to your project directory.

### Stopping the Services

To stop SonarQube and PostgreSQL:

```bash
make down
```
or
```bash
docker-compose down
```

## Configuration

The setup includes:
- SonarQube Community Edition
- PostgreSQL 17.2
- Persistent volumes for data storage
- Network configuration for service communication

### Ports
- SonarQube: 9000
- PostgreSQL: 5432

### Volumes
- `sonarqube-data`: Stores SonarQube data
- `sonarqube-db`: Stores PostgreSQL data