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

## Dependencies
- gRPC
- Protocol Buffers
- Google UUID
