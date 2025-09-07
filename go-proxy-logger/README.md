# Go Proxy Logger

A simple HTTP proxy server in Go that logs each request and response to SQLite with full metadata.

## Features
- Logs request/response timestamps, headers, bodies, and URLs
- Stores logs in SQLite
- Easy to run: `go run main.go`

## Requirements
- Go 1.21+
- SQLite

## Usage
1. Install dependencies: `go mod tidy`
2. Run the server: `go run main.go`
3. Proxy requests through `localhost:8080`

## TODO
- Capture response headers and body
- Improve error handling
- Add configuration options
