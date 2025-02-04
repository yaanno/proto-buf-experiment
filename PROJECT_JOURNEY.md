# Project Journey: Golang gRPC Calculator Service

## 1. System Design Clarification
- **Objective**: Create a distributed calculator service focusing on addition
- **Key Decisions**:
  - Use gRPC for inter-service communication
  - Implement a web handler and calculation service
  - Focus on learning and demonstration
  - Ensure clear API design and error handling

## 2. Technical Specifications Creation
- **Documented in [SPECS.md](SPECS.md)**
- Defined functional and non-functional requirements
- Outlined service interaction flow
- Specified error handling strategies
- Established API design principles

## 3. API Contract Design
- **Defined in Proto Definitions**
- Created `calculator.proto` with:
  - `AddRequest` message
    - Unique request ID
    - Repeated number inputs
  - `AddResponse` message
    - Calculation result
    - Optional error message
    - Request ID for tracing
- Versioned as `v1`
- Supported multiple number addition

## 4. Project Structure and Skeleton Creation
- Established monorepo architecture
- Created directories:
  ```
  proto-buf-experiment/
  ├── proto/
  │   └── calculator/v1/
  ├── services/
  │   ├── calculation/
  │   │   ├── cmd/
  │   │   └── internal/
  │   └── web-handler/
  │       ├── cmd/
  │       └── internal/
  └── gen/
  ```
- Initialized Go module
- Set up Buf configuration files

## 5. Code Generation with Buf
- Configured `buf.yaml` and `buf.gen.yaml`
- Generated Go code from proto definitions
- Created:
  - gRPC service interfaces
  - Message type structs
  - Client and server code

## 6. Service Implementation
### Calculation Service
- Implemented gRPC server
- Added addition logic
- Included:
  - Input validation
  - Overflow detection
  - Request ID generation
  - Error handling

### Web Handler Service
- Created HTTP to gRPC translation layer
- Implemented `/add` endpoint
- Features:
  - Request transformation
  - Timeout handling
  - Error propagation

## 7. Code Integrity and Validation Checks
- Ran Buf linting
- Performed `go vet` checks
- Validated package imports
- Ensured compilation success
- Verified dependency management

## 8. Key Technologies
- Go 1.21+
- gRPC
- Protocol Buffers
- Buf v1
- Gorilla Mux

## 9. Learning Outcomes
- Microservice architecture
- gRPC communication
- Proto-based API design
- Go best practices
- Service error handling

## 10. Potential Future Improvements
- Add comprehensive logging
- Implement more robust error handling
- Create integration tests
- Add authentication
- Implement more mathematical operations

## 11. Useful Commands and Workflows

### Buf Commands
#### Linting and Validation
```bash
# Lint proto files
buf lint

# Check for breaking changes
buf breaking --against .git#branch=main

# Update dependencies
buf dep update
```

#### Code Generation
```bash
# Generate code from proto definitions
buf generate

# Generate code with specific output
buf generate --template buf.gen.yaml
```

### Go Commands
#### Dependency Management
```bash
# Initialize go module
go mod init github.com/yourusername/project

# Add a new dependency
go get github.com/package/name

# Update dependencies
go get -u ./...

# Tidy and clean dependencies
go mod tidy

# Verify module dependencies
go mod verify
```

#### Code Validation and Checking
```bash
# Run static analysis
go vet ./...

# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Compile all packages
go build ./...

# Format all go files
go fmt ./...
```

#### Protobuf and gRPC Specific
```bash
# Install protobuf generator
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# Install gRPC generator
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Troubleshooting Commands
```bash
# Check go version
go version

# List all dependencies
go list -m all

# Show dependency tree
go mod graph

# Diagnose module issues
go mod verify
```

### Best Practices
1. Always run `go mod tidy` after changing dependencies
2. Use `go vet` before committing code
3. Run tests with race detector in CI/CD
4. Regularly update dependencies
5. Use `buf generate` to keep generated code in sync with proto definitions

## Conclusion
A learning-focused, well-structured microservice demonstrating modern Go and gRPC practices.
