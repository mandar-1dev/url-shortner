# URL Shortener API (Go + Gin)

A simple, beginner-friendly URL shortener REST API built with Go and the Gin
web framework. Uses in-memory storage (a Go map) — no database required.

## Features

- Shorten a long URL into a random 6-character code
- Redirect from the short URL to the original URL (HTTP 302)
- List all shortened URLs with click counts
- Delete a shortened URL
- Automatic click tracking

## Project Structure

```
url-shortener/
├── main.go              # Entry point, starts the server
├── go.mod                # Go module definition
├── handlers/
│   └── url.go            # HTTP handlers (business logic)
├── models/
│   └── url.go            # Data model + in-memory store
├── routes/
│   └── routes.go         # Route registration
├── utils/
│   └── generator.go       # Short code generator
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or newer installed ([https://go.dev/dl/](https://go.dev/dl/))

### Setup

```bash
# 1. Clone or download this project
cd url-shortener

# 2. Download dependencies (creates go.sum)
go mod tidy

# 3. Run the server
go run main.go
```

The server starts on `http://localhost:8080`.

## API Endpoints

### 1. Shorten a URL

```
POST /shorten
Content-Type: application/json

{
    "url": "https://google.com"
}
```

**Response (201 Created):**
```json
{
    "short_url": "http://localhost:8080/Ab12Cd"
}
```

### 2. Redirect to Original URL

```
GET /:code
```

Visiting `http://localhost:8080/Ab12Cd` in a browser (or via curl -L)
redirects (302) to `https://google.com` and increments the click count.

### 3. List All URLs

```
GET /urls
```

**Response (200 OK):**
```json
{
    "count": 1,
    "urls": [
        {
            "original_url": "https://google.com",
            "short_code": "Ab12Cd",
            "clicks": 3
        }
    ]
}
```

### 4. Delete a URL

```
DELETE /urls/:code
```

**Response (200 OK):**
```json
{
    "message": "Short URL deleted successfully"
}
```

## Example curl Commands

```bash
# Shorten a URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'

# List all URLs
curl http://localhost:8080/urls

# Follow a redirect
curl -L http://localhost:8080/Ab12Cd

# Delete a URL
curl -X DELETE http://localhost:8080/urls/Ab12Cd
```

## Notes

- Data is stored **in memory only** — restarting the server clears all URLs.
- Short codes are 6-character alphanumeric strings (e.g. `Ab12Cd`).
- URLs must start with `http://` or `https://` to be accepted.

## License

Free to use for learning purposes.
