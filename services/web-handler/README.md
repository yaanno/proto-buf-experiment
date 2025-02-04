# Web Handler Service

## Overview
A web service that translates HTTP requests to gRPC calls for the Calculation Service.

## Endpoints
- `POST /add`: Add multiple numbers
  - Request Body: `{"numbers": [1.0, 2.0, 3.0]}`
  - Response: 
    ```json
    {
      "result": 6.0,
      "request_id": "unique-uuid",
      "error": ""
    }
    ```

## Features
- HTTP to gRPC translation
- Request ID generation
- Timeout handling
- Error propagation

## Running the Service
```bash
go run cmd/main.go
```

## Dependencies
- Gorilla Mux (HTTP routing)
- gRPC
- Google UUID

## Configuration
- Calculation Service gRPC endpoint: `localhost:50051`
- Web Handler Port: `8080`
