# Distributed Calculator Service - Technical Specifications

## 0. Project Structure and API Design

### Monorepo Architecture
```
proto-buf-experiment/
│
├── buf.work.yaml           # Buf workspace configuration
├── buf.yaml                # Global Buf configuration
├── go.mod                  # Root module definition
│
├── proto/                  # Centralized proto definitions
│   ├── calculator/         # Domain-specific protos
│   │   └── v1/
│   │       ├── calculator.proto
│   │       └── service.proto
│   │
│   └── common/             # Shared proto definitions
│       └── v1/
│           ├── metadata.proto
│           └── errors.proto
│
├── services/               # Microservices
│   ├── web-handler/        # HTTP interface service
│   │   ├── cmd/
│   │   │   └── main.go     # Service entrypoint
│   │   ├── internal/       # Internal implementation
│   │   │   ├── server/
│   │   │   └── handlers/
│   │   └── configs/        # Configuration management
│   │
│   └── calculation/        # gRPC calculation service
│       ├── cmd/
│       │   └── main.go     # Service entrypoint
│       ├── internal/       # Internal implementation
│       │   ├── service/
│       │   └── logic/
│       └── configs/        # Configuration management
│
└── pkg/                    # Shared packages
    ├── logging/
    └── errors/
```

### API Design Principles

#### Versioning Strategy
1. **Proto Versioning**
   - Use semantic versioning in proto package names
   - Explicit `v1/`, `v2/` directories
   - Backward compatibility maintained

2. **API Version Tracking**
   - Include version in gRPC service names
   - Add version header in HTTP endpoints
   - Semantic versioning in module paths

#### Endpoint Design

##### gRPC Service (Calculation Service)
```protobuf
syntax = "proto3";

package calculator.v1;

service CalculatorService {
  // Addition operation with versioned RPC
  rpc AddNumbers(AddRequest) returns (AddResponse) {
    option (google.api.http) = {
      post: "/v1/calculator/add"
      body: "*"
    };
  }
}

message AddRequest {
  // Versioned request with metadata
  string request_id = 1;
  repeated double numbers = 2;
  Metadata metadata = 3;
}

message AddResponse {
  double result = 1;
  string version = 2;
  Metadata metadata = 3;
}

message Metadata {
  string service_version = 1;
  google.protobuf.Timestamp timestamp = 2;
}
```

##### HTTP Endpoint Design
- Base URL: `/v1/calculator/`
- Endpoints:
  - `POST /v1/calculator/add`
  - `GET /v1/calculator/health`
  - `GET /v1/calculator/version`

#### API Contract Guarantees
- Explicit version in every request/response
- Consistent error handling
- Metadata inclusion for traceability
- Semantic versioning of services

### API Evolution Guidelines
1. Never remove fields, only deprecate
2. New fields must be optional
3. Maintain backward compatibility
4. Use semantic versioning (MAJOR.MINOR.PATCH)
5. Provide migration guides between versions

### Error Handling Standardization
```protobuf
message ErrorDetails {
  enum ErrorType {
    UNSPECIFIED = 0;
    VALIDATION_ERROR = 1;
    CALCULATION_ERROR = 2;
    SYSTEM_ERROR = 3;
  }
  
  ErrorType type = 1;
  string code = 2;
  string message = 3;
  string details = 4;
  string request_id = 5;
}
```

### Documentation Requirements
- Comprehensive proto documentation
- Inline comments for all messages and services
- Example request/response scenarios
- Error code explanations

### Compatibility Considerations
- Support for multiple client languages
- Consistent serialization
- Platform-agnostic design

## 1. Requirements Specification

### 1.1 Functional Requirements

#### Core Mathematical Operation
- [x] Support addition operation
  - Add multiple integer and floating-point numbers
  - Validate input types
- [x] Handle single and multiple operands
- [x] Provide clear error handling for invalid inputs

#### Service Interaction
- [x] Implement gRPC service for addition
- [x] Create a simple web handler for HTTP-based interactions
- [x] Support synchronous calculation requests

#### Input Validation
- [x] Validate input types and ranges
- [x] Prevent invalid number inputs
- [x] Handle potential overflow scenarios

### 1.2 Non-Functional Requirements

#### Performance
- [ ] Latency Goal
  - Target average response time: < 50ms
- [ ] Resource Consumption
  - Memory: Minimal (< 30MB)
  - CPU: Lightweight processing

#### Reliability
- [ ] Basic error handling
- [ ] Structured logging
- [ ] Simple error recovery

#### Observability
- [ ] Basic logging
  - Log calculation inputs
  - Log errors
  - Use structured logging format

#### Compatibility
- [x] Go (primary implementation)
- [ ] Potential gRPC client generation

### 1.3 Constraints
- Use Go 1.21
- Use gRPC and Protocol Buffers
- Buf for proto management
- Minimal external dependencies
- Learning-focused implementation

### 1.4 Assumptions
- Calculations are simple additions
- No complex mathematical requirements
- Single-instance deployment
- Learning and demonstration purpose

### 1.5 Out of Scope
- Multiple mathematical operations
- High availability
- Advanced monitoring
- Distributed tracing
- Authentication
- Persistent storage
- Complex error handling

## 2. Development Tools and Workflow

### Buf v2 Ecosystem Toolchain

#### Buf v2 Key Differences
- Enhanced configuration model
- Improved plugin management
- More granular lint and breaking change rules
- Better monorepo and multi-module support

#### Buf v2 Configuration Files

##### `buf.work.yaml` (Workspace Configuration)
```yaml
version: v2
directories:
  - proto
  - services
```

##### `buf.yaml` (Module Configuration)
```yaml
version: v2
name: buf.build/yourusername/proto-buf-experiment

# V2 Lint Configuration
lint:
  # Default ruleset with custom exceptions
  use:
    - BASIC
    - FILE_LOWER_SNAKE_CASE
    - PACKAGE_LOWER_SNAKE_CASE
    - ENUM_PASCAL_CASE
  except:
    - FIELD_LOWER_SNAKE_CASE  # Allow more flexible naming

# V2 Breaking Change Detection
breaking:
  use:
    - FILE
    - PACKAGE
    - WIRE_JSON
  ignore_unstable_packages: true
```

##### `buf.gen.yaml` (Code Generation Configuration)
```yaml
version: v2
managed:
  enabled: true
  go_package_prefix:
    default: github.com/yourusername/proto-buf-experiment/gen

plugins:
  # Go code generation
  - name: go
    out: gen/go
    opt:
      - paths=source_relative
      - go_opt=paths=source_relative
  
  # gRPC Go plugin
  - name: go-grpc
    out: gen/go
    opt:
      - paths=source_relative
  
  # Optional: Connect RPC support
  - name: connect-go
    out: gen/go
    opt:
      - paths=source_relative
```

#### Buf v2 Development Workflow

##### Proto Management Commands
```bash
# Initialize a new Buf module
buf mod init

# Lint with v2 configuration
buf lint

# Format proto files
buf format -w

# Check for breaking changes
buf breaking --against .git#branch=main

# Generate code with v2 configuration
buf generate

# Validate module configuration
buf mod update
```

### Buf v2 Best Practices

#### Proto Design Principles
1. Use `lower_snake_case` for package and file names
2. Use `PascalCase` for message and enum names
3. Use `lower_snake_case` for field names
4. Add clear, concise comments
5. Prefer optional fields for extensibility

#### Code Generation Workflow
1. Define proto contracts in `proto/` directory
2. Run `buf lint` for validation
3. Run `buf generate` for code generation
4. Implement service logic in generated interfaces
5. Write tests against generated code

#### Version Control Considerations
- Commit generated code
- Use `.gitignore` to manage generated files
- Tag releases with semantic versioning
- Include `buf.lock` for dependency locking

### Advanced Buf v2 Features

#### Dependency Management
- Use `buf.lock` for exact dependency pinning
- Support for local and remote proto module dependencies
- Integrated dependency resolution

#### Plugin Management
- More flexible plugin configuration
- Support for custom plugins
- Easy integration with existing toolchains

#### Schema Registry Integration
- Improved Buf Schema Registry support
- Push and pull schemas with enhanced metadata
- Version tracking and discovery

### Troubleshooting and Support
- Use `buf debug` for configuration insights
- Consult Buf v2 migration guide
- Join Buf community channels
- Reference official Buf documentation

### Future Improvements
- Explore Connect RPC integration
- Investigate multi-language generation
- Continuous learning of Buf v2 features

## 3. Detailed Service Specifications

### Web Handler Service
#### Configuration
```go
type WebHandlerConfig struct {
    Port                int      `env:"WEB_PORT" default:"8080"`
    CalculationEndpoint string   `env:"CALC_ENDPOINT" default:"localhost:50051"`
    AllowedOrigins      []string `env:"CORS_ORIGINS" default:"*"`
    RequestTimeout      duration `env:"REQUEST_TIMEOUT" default:"5s"`
}
```

#### Middleware Requirements
- Request logging
- CORS support
- Input sanitization
- Prometheus metrics endpoint

### Calculation Service
#### Configuration
```go
type CalculationServiceConfig struct {
    Port                int      `env:"GRPC_PORT" default:"50051"`
    MaxConcurrentCalls  int      `env:"MAX_CONCURRENT" default:"100"`
    MetricsEnabled      bool     `env:"METRICS_ENABLED" default:"true"`
    LogLevel            string   `env:"LOG_LEVEL" default:"info"`
}
```

## 3. Protocol Buffer Definitions

### 3.1 Addition Service Proto
```protobuf
syntax = "proto3";

package calculator.v1;

service AdditionService {
  // RPC method for addition with explicit result return
  rpc Add(AddRequest) returns (AddResponse) {
    option (google.api.http) = {
      post: "/v1/calculator/add"
      body: "*"
    };
  }
}

message AddRequest {
  // Unique request identifier for tracing
  string request_id = 1;
  
  // Numbers to be added
  repeated double numbers = 2;
}

message AddResponse {
  // Calculation result
  double result = 1;
  
  // Optional error message
  string error = 2;
  
  // Request identifier for correlation
  string request_id = 3;
}
```

### 3.2 Service Interaction Flow
1. Web Handler creates `AddRequest` with:
   - Unique `request_id`
   - List of numbers to add
2. Calculation Service processes request
3. Calculation Service returns `AddResponse` containing:
   - Calculation `result`
   - Optional `error`
   - Original `request_id`
4. Web Handler uses returned result for HTTP response

## 4. Error Handling

### 4.1 Error Categories
- Input validation errors
- Calculation overflow
- Type conversion errors
- Request-response correlation errors

### 4.2 Error Codes
```go
type ErrorCode int

const (
    ErrOK ErrorCode = iota
    ErrInvalidInput
    ErrOverflow
    ErrRequestCorrelation
)
```

### 4.3 Error Handling Strategy
- Each request has a unique identifier
- Errors include descriptive messages
- Errors propagated back to web handler
- Web handler responsible for converting gRPC errors to HTTP errors

## 5. Performance Benchmarks
### Benchmark Scenarios
- Single operation latency
- Memory consumption
- CPU utilization

### Benchmark Targets
- **Latency**: 
  - Target average response time: < 50ms
- **Memory**: 
  - < 30MB per service instance

## 6. Security Specifications
### Input Validation
- Strict type checking
- Limit maximum number of input values
- Sanitize and validate all inputs

## 7. Logging Specifications
### Log Levels
- `DEBUG`: Detailed execution flow
- `INFO`: Operational milestones
- `WARN`: Potential issues
- `ERROR`: Failure scenarios
- `FATAL`: Unrecoverable errors

### Log Fields
```go
type LogEntry struct {
    Timestamp     time.Time
    Level         string
    ServiceName   string
    RequestID     string
    Operation     string
    InputValues   []float64
    Result        float64
    ErrorCode     ErrorCode
    ErrorMessage  string
}
```

## 8. Deployment Considerations
### Container Specifications
- Alpine Linux base image
- Minimal runtime dependencies
- Non-root container user
- Health check endpoints

### Kubernetes Deployment
- Resource limits
- Horizontal Pod Autoscaler configuration
- Readiness and liveness probes

## 9. Development Workflow
### Git Workflow
- Feature branch model
- Conventional commit messages
- Mandatory code review
- Automated CI/CD pipeline

### Code Quality Gates
- 80%+ test coverage
- No high-severity linter warnings
- Passing integration tests
- Performance benchmark comparisons

## 10. Versioning Strategy
- Semantic Versioning (SemVer)
- Backward compatibility preservation
- Clear deprecation policies

## 11. Future Extension Points
- Support for complex mathematical functions
- Machine learning model integration
- Streaming calculation support
- Multi-language client libraries

## 12. Buf Integration Strategy

### Buf Configuration Philosophy
- Centralized proto definition management
- Strict schema validation
- Cross-language code generation
- Breaking change detection

### Project Structure
```
proto-buf-experiment/
│
├── buf.yaml           # Buf workspace configuration
├── buf.gen.yaml       # Code generation configuration
├── buf.work.yaml      # Monorepo workspace configuration
│
├── proto/             # Centralized proto definitions
│   ├── calculator/    # Domain-specific proto package
│   │   ├── v1/
│   │   │   ├── calculator.proto
│   │   │   └── calculator_service.proto
│   │
│   └── common/        # Shared proto definitions
│       ├── v1/
│       │   ├── metadata.proto
│       │   └── errors.proto
│
└── gen/               # Generated code output directory
    ├── go/
    ├── python/
    └── typescript/
```

### Buf Configuration Files

#### `buf.yaml` (Workspace Configuration)
```yaml
version: v2
name: buf.build/yourusername/proto-buf-experiment
deps:
  - buf.build/googleapis/googleapis
lint:
  use:
    - DEFAULT
  except:
    - PACKAGE_VERSION_SUFFIX
breaking:
  use:
    - FILE
    - PACKAGE
```

#### `buf.gen.yaml` (Code Generation Configuration)
```yaml
version: v2
managed:
  enabled: true
plugins:
  - plugin: buf.build/protocolbuffers/go
    out: gen/go
    opt: 
      - paths=source_relative
  - plugin: buf.build/grpc/go
    out: gen/go
    opt:
      - paths=source_relative
  - plugin: buf.build/protocolbuffers/python
    out: gen/python
  - plugin: buf.build/grpc/python
    out: gen/python
```

#### `buf.work.yaml` (Monorepo Workspace)
```yaml
version: v2
directories:
  - proto
```

### Protobuf Design Principles
1. **Versioning**
   - Use semantic versioning in package names
   - Maintain backward compatibility
   - Clear deprecation paths

2. **Package Organization**
   - Domain-driven proto package structure
   - Separate concerns into distinct proto files
   - Use `v1`, `v2` etc. for major versions

3. **Code Generation Targets**
   - Go (primary implementation)
   - Python (potential client libraries)
   - TypeScript (potential web clients)

### Buf CLI Workflow
```bash
# Initialize new proto package
buf new proto/calculator

# Lint proto files
buf lint

# Break compatibility check
buf breaking --against .git#branch=main

# Generate code
buf generate

# Push to Buf Schema Registry (optional)
buf push
```

### Future Improvements
- Automated CI integration
- Schema registry usage
- Multi-language support
- Advanced breaking change detection
