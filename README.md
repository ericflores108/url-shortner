# URL Shortener

A simple URL shortener service built with Go and Redis that allows you to convert long URLs into short, manageable links.

## Features

- Shorten long URLs to compact short codes
- Redirect from short codes to original URLs  
- Redis-based storage for fast lookups
- RESTful API endpoints
- Health check endpoint
- CORS support
- Graceful shutdown

## API Endpoints

- `POST /shorten` - Shorten a URL
- `GET /r/{shortCode}` - Redirect to original URL
- `GET /health` - Health check
- `GET /` - Hello world

## Prerequisites

- Go 1.24.2+
- Redis
- Docker & Docker Compose (optional)

## Quick Start

### Using Docker Compose

```bash
make docker-run
```

### Local Development

1. Start Redis locally (port 6379)
2. Set environment variables:
   ```bash
   export URLSHORT_DB_ADDRESS=localhost
   export URLSHORT_DB_PORT=6379
   export PORT=8080
   ```
3. Run the application:
   ```bash
   make run
   ```

## Environment Variables

- `PORT` - Server port
- `REDIS_URL` - Redis URL to parse
- `APP_ENV` - Application environment

