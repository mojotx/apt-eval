# Apartment Evaluation App

A simple Go application for evaluating apartments with CRUD functionality. This application uses:

- [gin-gonic/gin](https://github.com/gin-gonic/gin) for the web server
- [SQLite](https://www.sqlite.org/) for the database
- [zerolog](https://github.com/rs/zerolog) for structured logging

## Features

- Create, read, update, and delete apartment evaluations
- SQLite database for persistent storage
- RESTful API endpoints

## Requirements

- Go 1.18+ (developed with Go 1.21)
- GCC (for compiling SQLite driver)

## Installation

```bash
# Clone the repository
git clone https://github.com/mojotx/apt-eval.git
cd apt-eval

# Install dependencies
go mod download
```

## Usage

### Starting the server

```bash
go run main.go
```

The server will start on port 8443 by default with HTTPS enabled. You can change the port with the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

### TLS Configuration

The application supports HTTPS using TLS. By default, it looks for certificate files in the `./certs` directory:

1. Place your wildcard certificate at `./certs/wildcard.crt`
2. Place your private key at `./certs/wildcard.key`

Or specify custom paths using environment variables:

```bash
CERT_FILE=/path/to/your/certificate.crt KEY_FILE=/path/to/your/private.key go run main.go
```

### Landing Page

Access the web interface at: [https://localhost:8443/](https://localhost:8443/)

### API Endpoints

#### Create an apartment evaluation

```text
POST /api/apartments
```

Request body:

```json
{
  "address": "123 Main St, Apt 4B",
  "visit_date": "2025-09-05T14:30:00Z",
  "notes": "Nice layout, good natural light",
  "rating": 4,
  "price": 1500
}
```

#### Get all apartment evaluations

```text
GET /api/apartments
```

#### Get a specific apartment evaluation

```text
GET /api/apartments/:id
```

#### Update an apartment evaluation

```text
PUT /api/apartments/:id
```

Request body (same format as create):

```json
{
  "address": "123 Main St, Apt 4B",
  "visit_date": "2025-09-05T14:30:00Z",
  "notes": "Nice layout, good natural light, but noisy neighbors",
  "rating": 3,
  "price": 1500
}
```

#### Delete an apartment evaluation

```text
DELETE /api/apartments/:id
```

### Health Check

```text
GET /health
```

## Environment Variables

- `PORT`: HTTPS server port (default: 8443)
- `HTTP_PORT`: HTTP server port for redirects (default: 8080)
- `DATA_DIR`: Directory for SQLite database (default: ./data)
- `CERT_FILE`: Path to TLS certificate file (default: ./certs/wildcard.crt)
- `KEY_FILE`: Path to TLS private key file (default: ./certs/wildcard.key)

## Building for Production

```bash
go build -o apt-eval .
```

## License

MIT
