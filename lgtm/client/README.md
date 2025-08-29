# LGTM Client - Express Log Generator

A simple Express.js application that generates logs for testing the LGTM (Loki, Grafana, Tempo, Mimir) observability stack.

## Features

- Winston-based structured logging
- **OpenTelemetry OTLP integration** - Direct log forwarding to Alloy
- Multiple API endpoints for log generation
- JSON formatted logs for easy parsing
- Integration with LGTM stack via Alloy
- Graceful shutdown with telemetry cleanup

## Installation

```bash
npm install
```

## Development

```bash
npm run dev
```

## Build & Run

```bash
npm run build
npm start
```

## API Endpoints

### Core Routes

#### `GET /test-otlp`
Tests OTLP log forwarding to Alloy.

**Response:**
```json
{
  "success": true,
  "message": "OTLP test logs sent to Alloy",
  "otlpEndpoint": "http://localhost:4318/v1/logs",
  "timestamp": "2025-01-27T..."
}
```

**Example:**
```bash
curl http://localhost:3001/test-otlp
```

#### `GET /simple-generate-log`
Generates logs with customizable parameters.

**Query Parameters:**
- `level` (optional): Log level (info, warn, error, debug) - default: "info"
- `message` (optional): Custom log message - default: "Simple log generated"
- `count` (optional): Number of logs to generate - default: 1

**Examples:**
```bash
# Generate a single info log
curl http://localhost:3001/simple-generate-log

# Generate 5 error logs
curl "http://localhost:3001/simple-generate-log?level=error&count=5"

# Generate custom message logs
curl "http://localhost:3001/simple-generate-log?message=Custom%20log%20message&count=3"
```

#### `GET /health`
Health check endpoint.

**Response:**
```json
{
  "status": "OK",
  "timestamp": "2025-01-27T...",
  "uptime": 123.456
}
```

### Simulation Routes

#### `GET /user/:id`
Simulates user data retrieval.

**Example:**
```bash
curl http://localhost:3001/user/123
```

#### `GET /products`
Simulates product listing with pagination.

**Query Parameters:**
- `page` (optional): Page number - default: 1
- `limit` (optional): Items per page - default: 10

**Example:**
```bash
curl "http://localhost:3001/products?page=2&limit=5"
```

#### `GET /orders`
Simulates order listing with filtering.

**Query Parameters:**
- `status` (optional): Filter by status (pending, completed, cancelled, all) - default: "all"

**Example:**
```bash
curl "http://localhost:3001/orders?status=pending"
```

#### `GET /stats`
Generates random statistics data.

**Example:**
```bash
curl http://localhost:3001/stats
```

#### `GET /simulate-error`
Simulates different types of errors for testing.

**Query Parameters:**
- `type` (optional): Error type (auth, notfound, server, generic) - default: "generic"

**Examples:**
```bash
# Simulate 401 error
curl "http://localhost:3001/simulate-error?type=auth"

# Simulate 404 error
curl "http://localhost:3001/simulate-error?type=notfound"

# Simulate 500 error
curl "http://localhost:3001/simulate-error?type=server"
```

## Logging

The application uses Winston for structured logging with the following features:

- **Console output**: Colorized, human-readable format for development
- **File output**: JSON format saved to `/var/log/app.log`
- **OTLP output**: Direct log forwarding to Alloy via OpenTelemetry Protocol
- **Request logging**: All HTTP requests are automatically logged
- **Error tracking**: Comprehensive error logging with stack traces

### Log Structure

All logs include:
- Timestamp
- Log level
- Service name (`lgtm-client`)
- Request metadata (method, URL, IP, User-Agent)
- Custom context data

## Environment Variables

- `PORT`: Server port (default: 3001)
- `NODE_ENV`: Environment (development/production)
- `LOG_LEVEL`: Winston log level (default: info)
- `OTLP_LOG_LEVEL`: OTLP transport log level (default: info)
- `OTLP_ENDPOINT`: OpenTelemetry endpoint (default: http://localhost:4318/v1/logs)
- `SERVICE_NAME`: Service identifier (default: lgtm-client)
- `SERVICE_VERSION`: Service version (default: 1.0.0)
- `SERVICE_NAMESPACE`: Service namespace (default: lgtm-stack)

Copy `env.example` to `.env` to configure these values.

## Integration with LGTM Stack

This application integrates with the LGTM observability stack in two ways:

### **Method 1: OTLP Direct Integration (Recommended)**
1. **Application** sends logs directly to **Alloy** via OTLP HTTP (port 4318)
2. **Alloy** receives OTLP logs and forwards them to **Loki**
3. **Grafana** queries and visualizes the logs from **Loki**

### **Method 2: File-based Integration**
1. **Logs** are written to `/var/log/app.log` in JSON format
2. **Alloy** scrapes these logs and forwards them to **Loki**
3. **Grafana** can query and visualize the logs from **Loki**

## Testing the Setup

1. Start the LGTM stack:
   ```bash
   docker-compose up -d
   ```

2. Start the client application:
   ```bash
   npm run dev
   ```

3. Generate some logs:
   ```bash
   # Test OTLP integration
   curl http://localhost:3001/test-otlp
   
   # Generate traditional logs
   curl "http://localhost:3001/simple-generate-log?count=10&level=info"
   curl "http://localhost:3001/simple-generate-log?count=5&level=error"
   ```

4. Check logs in Grafana at `http://localhost:3000`

## Docker Integration

The application logs to `/var/log/app.log` which should be mounted as a volume in the Docker container for Alloy to access.
