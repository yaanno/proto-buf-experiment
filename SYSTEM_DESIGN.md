# Distributed Calculator Service - System Design

## 1. Executive Summary

### Purpose
A learning-focused microservice project demonstrating:
- Distributed system design
- gRPC communication
- Modern development practices

### Core Concept
A simple calculator service specializing in addition operations, designed to provide hands-on experience with microservice architecture.

## 2. System Goals

### Primary Objectives
- Create a modular, extensible calculator service
- Demonstrate effective inter-service communication
- Implement clean, maintainable code
- Provide a learning platform for distributed systems

### Learning Outcomes
- gRPC service design
- Protocol Buffer usage
- Microservice architecture
- Modern Go development practices

## 3. System Context

### High-Level Architecture
```
[Web Handler Service] <--gRPC--> [Calculation Service]
         |                               |
         V                               V
     HTTP Endpoint               Core Calculation Logic
         ^                               |
         |                               |
         +--- Calculation Result --------+
```

### Result Flow
1. Web Handler receives HTTP request
2. Web Handler transforms request and calls Calculation Service via gRPC
3. Calculation Service performs addition
4. Calculation Service returns result to Web Handler
5. Web Handler presents result to user via HTTP response

### Key Components
1. **Web Handler Service**
   - Receives HTTP requests
   - Validates and transforms inputs
   - Communicates with Calculation Service via gRPC

2. **Calculation Service**
   - Implements core addition logic
   - Provides gRPC interface
   - Stateless design

## 4. Architectural Principles

### Design Philosophy
- Simplicity over complexity
- Clear separation of concerns
- Extensibility
- Learning-driven implementation

### Communication Patterns
- Synchronous gRPC communication
- RESTful HTTP interface
- Explicit versioning

## 5. Scope and Constraints

### Functional Boundaries
- Addition operation only
- Single-instance deployment
- Minimal error handling
- No persistent storage

### Non-Functional Constraints
- Learning-focused performance
- Simplified monitoring
- Basic error handling
- Minimal external dependencies

## 6. Future Potential

### Possible Extensions
- Support for more mathematical operations
- Advanced error handling
- Distributed deployment
- **Completed: Implemented structured logging**
- Authentication mechanisms
- Advanced logging features

### Logging and Observability
- Implemented structured logging using Zerolog
- Supports console and file logging
- Contextual logging with request IDs
- Configurable log levels and output
- JSON-formatted log output for easy parsing

## 7. Decision Log

### Key Design Choices
- gRPC for inter-service communication
- Buf v2 for proto management
- Go 1.21 as primary implementation language
- Monorepo architectural approach

### Rationale
- gRPC: Efficient, typed communication
- Buf v2: Modern proto management
- Go: Performance and simplicity
- Monorepo: Simplified dependency management

## 8. Risk Assessment

### Potential Challenges
- Over-engineering
- Complexity creep
- Divergence from learning goals

### Mitigation Strategies
- Regular design reviews
- Maintain focus on learning
- Keep implementation minimal and clear

## 9. Stakeholder Considerations

#### Target Audience
- Software engineers learning distributed systems
- Go developers exploring microservices
- Students and self-learners
- DevOps and SRE professionals interested in logging practices

#### Success Criteria
- Clear, understandable implementation
- Demonstrable inter-service communication
- Ease of comprehension
- Potential for incremental learning
- **Robust logging and observability**

## 10. Conclusion

A purposefully simple yet instructive microservice project designed to provide hands-on experience with modern distributed system concepts, focusing on clarity, learning, and practical implementation.
