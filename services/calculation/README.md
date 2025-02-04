# Calculation Service

## Overview
This service provides a gRPC-based addition service that can add multiple numbers.

## Features
- Add multiple numbers via gRPC
- Request ID tracking
- Basic error handling
- Overflow detection

## Running the Service
```bash
go run cmd/main.go
```

## Error Handling
- Returns error if no numbers are provided
- Detects and handles calculation overflow
- Generates a unique request ID if not provided

## Logging
- Structured logging with Zerolog
- Logs calculation inputs and results
- Captures detailed error information
- Configurable log levels
- Logs written to console and file
- JSON-formatted log output
- Contextual logging with request IDs

### Log Configuration
- `DEBUG` environment variable controls log verbosity
- Log files stored in configurable location
- Supports rotation and retention policies

## Dependencies
- gRPC
- Protocol Buffers
- Google UUID
