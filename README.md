# Telegram API Bridge

A clean, maintainable Go service for sending messages to Telegram users via Bot API.

## Features

- Send messages to single or multiple Telegram users/chats
- Concurrent message delivery for better performance
- Structured logging with file rotation
- Environment-based configuration
- Clean architecture with separated concerns
- Health check endpoint

## Project Structure

```
telegram-bridge/
├── main.go                    # Application entry point
├── .env.example              # Environment variables template
├── config/
│   └── config.go             # Configuration management
├── models/
│   └── models.go             # Request/response models
├── telegram/
│   └── client.go             # Telegram API client
├── middleware/
│   └── auth.go               # Authentication middleware
├── handlers/
│   └── handlers.go           # HTTP request handlers
├── router/
│   └── router.go             # Route configuration
└── logger/
    └── logger.go             # Structured logging
```

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Environment

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` with your settings:

```env
# Server Configuration
PORT=8080

# Authentication
AUTH_TOKEN=your-secret-auth-token-here

# Logging Configuration
LOG_LEVEL=info
LOG_DIR=./logs
LOG_MAX_SIZE=100
LOG_MAX_BACKUPS=5
LOG_MAX_AGE=30
LOG_COMPRESS=true
LOG_CONSOLE=true
```

### 3. Run the Service

```bash
go run main.go
```

Or build and run:

```bash
go build -o telegram-bridge
./telegram-bridge
```

## API Endpoints

### POST /send

Send messages to Telegram users/chats.

**Headers:**
```
Authorization: Bearer <your-auth-token>
Content-Type: application/json
```

**Request Body:**

Option 1 - Single chat:
```json
{
  "bot_token": "your-telegram-bot-token",
  "message": "Hello from Telegram Bridge!",
  "chat_id": "123456789"
}
```

Option 2 - Multiple users:
```json
{
  "bot_token": "your-telegram-bot-token",
  "message": "Hello from Telegram Bridge!",
  "user_ids": ["123456789", "987654321"]
}
```

**Response:**
```json
{
  "success": true,
  "total_sent": 2,
  "success_count": 2,
  "fail_count": 0
}
```

### GET /health

Health check endpoint (no authentication required).

**Response:**
```json
{
  "status": "healthy",
  "service": "telegram-api-bridge"
}
```

## Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `AUTH_TOKEN` | (default token) | API authentication token |
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `LOG_DIR` | `./logs` | Directory for log files |
| `LOG_MAX_SIZE` | `100` | Max log file size in MB |
| `LOG_MAX_BACKUPS` | `5` | Max number of old log files to keep |
| `LOG_MAX_AGE` | `30` | Max days to retain old log files |
| `LOG_COMPRESS` | `true` | Compress rotated log files |
| `LOG_CONSOLE` | `true` | Enable console logging |

## Example Usage

### cURL

```bash
curl -X POST http://localhost:8080/send \
  -H "Authorization: Bearer your-auth-token" \
  -H "Content-Type: application/json" \
  -d '{
    "bot_token": "your-bot-token",
    "message": "Test message",
    "chat_id": "123456789"
  }'
```

### Python

```python
import requests

url = "http://localhost:8080/send"
headers = {
    "Authorization": "Bearer your-auth-token",
    "Content-Type": "application/json"
}
data = {
    "bot_token": "your-bot-token",
    "message": "Test message",
    "user_ids": ["123456789", "987654321"]
}

response = requests.post(url, json=data, headers=headers)
print(response.json())
```

## Development

### Run Tests

```bash
go test ./...
```

### Build

```bash
go build -o telegram-bridge
```

### Docker (Optional)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o telegram-bridge

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/telegram-bridge .
CMD ["./telegram-bridge"]
```

## License

MIT
