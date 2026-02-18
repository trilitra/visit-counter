# IP Visit Tracker

A simple HTTP service written in Go that:

- Detects client IP address
- Counts visits per IP
- Tracks number of unique IPs
- Supports in-memory and Redis storage
- Provides JSON API endpoints

---

## Features

- `GET /health` — health check endpoint
- `GET /whoami` — returns client IP
- `GET /visit` — increments and returns visit count for IP
- `GET /stats` — returns number of unique IPs

---

## Requirements

- Go 1.20+
- (Optional) Redis if using Redis storage

---

## Installation

    git clone <your-repo-url>
    cd <project-folder>
    go mod tidy

---

## Environment Variables

The application is configured via environment variables.

You can define them either directly in your system or by creating a `.env` file in the project root.

| Variable        | Default          | Description                                |
|-----------------|------------------|--------------------------------------------|
| PORT            | 8080             | HTTP server port                           |
| STORE_TYPE      | memory           | Storage backend: `memory` or `redis`       |
| REDIS_ADDR      | localhost:6379   | Redis server address                       |
| REDIS_PASSWORD  | (empty)          | Redis password (if required)               |

---

## Configuration and Running

### Run with in-memory storage (default)

    go run main.go

Server will start at:

    http://localhost:8080

---

### Run with Redis

1. Start Redis

Using Docker:

    docker run -d -p 6379:6379 --name redis redis

Or if Redis is installed locally:

    redis-server

2. Create a `.env` file (recommended)

Create a file named `.env` in the project root:

    STORE_TYPE=redis
    REDIS_ADDR=localhost:6379
    REDIS_PASSWORD=
    PORT=8080

3. Run the application

   go run main.go

The application automatically loads `.env` using `godotenv`.

---

## API Examples

### Health Check

    curl http://localhost:8080/health

Response:

    {
      "status": "healthy"
    }

---

### Who Am I

    curl http://localhost:8080/whoami

Response:

    {
      "ip": "192.168.1.10"
    }

---

### Visit Counter

    curl http://localhost:8080/visit

Response:

    {
      "ip": "192.168.1.10",
      "visits": 3
    }

Each request increments the visit counter for that IP.

---

### Statistics

    curl http://localhost:8080/stats

Response:

    {
      "unique_ips": 12
    }

---

## Architecture

The application uses a storage abstraction:

Store interface:
- MemoryStore
- RedisStore

Storage implementation is selected at runtime using `STORE_TYPE`.

---

## Production Notes

- Redis connection is validated on startup
- Structured logging with slog
- Clean dependency injection for storage
- JSON error handling

---

## License

MIT 
