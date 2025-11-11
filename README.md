# Telegram Bridge Server

A dynamic TCP server that acts as a bridge to send Telegram messages. The server accepts requests with JWT authentication and dynamically creates Telegram bot instances based on the bot token provided in each request.

## Features

- JWT token authentication for secure access
- Dynamic Telegram bot token handling (no static bot token in config)
- Support for sending text messages
- Support for sending images (base64 encoded)
- Support for sending images with captions
- Structured logging with file rotation
- Multiple user targeting per request

## Setup

### 1. Environment Configuration

Create a `.env` file based on `.env.example`:

```bash
PORT=8080
JWT_SECRET=your-secure-jwt-secret-key-here
LOG_LEVEL=info
```

### 2. Build and Run

```bash
go build
./telegram-bridge
```

## API Usage

### Request Format

Connect to the TCP server and send JSON requests in the following format:

#### Text Message
```json
{
  "JWT": "your-jwt-token",
  "TELEGRAM_BOT_TOKEN": "bot-token-from-telegram-botfather",
  "USER_IDS": [123456789, 987654321],
  "MESSAGE": "Hello from the bridge!",
  "TYPE": "message"
}
```

#### Image Only
```json
{
  "JWT": "your-jwt-token",
  "TELEGRAM_BOT_TOKEN": "bot-token-from-telegram-botfather",
  "USER_IDS": [123456789],
  "IMAGE_BASE64": "base64-encoded-image-data",
  "TYPE": "image"
}
```

#### Image with Caption
```json
{
  "JWT": "your-jwt-token",
  "TELEGRAM_BOT_TOKEN": "bot-token-from-telegram-botfather",
  "USER_IDS": [123456789],
  "MESSAGE": "Check out this image!",
  "IMAGE_BASE64": "base64-encoded-image-data",
  "TYPE": "image_and_message"
}
```

### Response Format

```json
{
  "success": true,
  "message": "Message sent successfully",
  "error": ""
}
```

## JWT Token Generation

### Using the CLI Tool (Recommended)

The project includes a built-in JWT token generator that automatically reads the `JWT_SECRET` from your `.env` file:

```bash
# Build the JWT generator
go build -o jwt-gen ./cmd/jwt-gen

# Generate a token with no expiration
./jwt-gen

# Generate a token that expires in 24 hours
./jwt-gen -expires 24h

# Generate a token that expires in 7 days
./jwt-gen -expires 168h

# Generate a token with a custom subject
./jwt-gen -subject "my-app" -expires 30d

# Show help
./jwt-gen -help
```

The CLI tool will:
- Automatically load `JWT_SECRET` from your `.env` file
- Generate a properly signed JWT token
- Display token details including expiration time
- Support custom claims (subject, issuer)

### Manual Generation

Alternatively, you can generate JWT tokens manually using Go:

```go
package main

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

func main() {
    secret := "your-jwt-secret"

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
        "iat": time.Now().Unix(),
    })

    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        panic(err)
    }

    fmt.Println(tokenString)
}
```

Or using online tools like [jwt.io](https://jwt.io) with your JWT_SECRET.

## Testing

### Using netcat (Linux/Mac)
```bash
echo '{"JWT":"your-jwt-token","TELEGRAM_BOT_TOKEN":"your-bot-token","USER_IDS":[123456789],"MESSAGE":"Test message","TYPE":"message"}' | nc localhost 8080
```

### Using PowerShell (Windows)
```powershell
$json = '{"JWT":"your-jwt-token","TELEGRAM_BOT_TOKEN":"your-bot-token","USER_IDS":[123456789],"MESSAGE":"Test message","TYPE":"message"}'
$bytes = [System.Text.Encoding]::UTF8.GetBytes($json + "`n")
$client = New-Object System.Net.Sockets.TcpClient("localhost", 8080)
$stream = $client.GetStream()
$stream.Write($bytes, 0, $bytes.Length)
$reader = New-Object System.IO.StreamReader($stream)
$response = $reader.ReadLine()
Write-Host $response
$client.Close()
```

## Project Structure

```
telegram-bridge/
├── cmd/             # Command-line tools
│   └── jwt-gen/     # JWT token generator CLI
│       └── main.go
├── config/          # Configuration management
│   ├── config.go    # Config struct and loader
│   └── env.go       # Environment variable utilities
├── logger/          # Logging system
│   └── logger.go    # Structured logger with rotation
├── server/          # TCP server
│   └── server.go    # Server implementation with JWT validation
├── telegram/        # Telegram API integration
│   └── telegram.go  # Dynamic bot service
├── main.go          # Application entry point
├── .env.example     # Environment variables template
└── README.md        # Documentation
```

## Security Notes

- Keep your JWT_SECRET secure and never commit it to version control
- Use strong, randomly generated JWT secrets
- Telegram bot tokens are never stored on the server
- All requests are validated with JWT before processing
- Failed authentication attempts are logged

## Error Handling

The server returns detailed error messages:
- Invalid JWT token
- Missing required fields
- Invalid message type
- Telegram API errors
- Base64 decoding errors

All errors are logged for audit purposes.
