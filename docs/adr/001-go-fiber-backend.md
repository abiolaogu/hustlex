# ADR-001: Go with Fiber Framework for Backend

## Status

Accepted

## Date

2024-01-15

## Context

HustleX requires a high-performance backend API capable of handling:
- Financial transactions (wallet operations, payments, escrow)
- Real-time notifications
- Background job processing
- High concurrent user load (targeting 10,000+ concurrent users)
- Nigerian market requirements (low latency, high reliability)

We needed to choose a programming language and web framework that could:
1. Handle high concurrency efficiently
2. Provide excellent performance for API responses
3. Compile to small, efficient binaries for containerized deployment
4. Have good ecosystem support for databases, caching, and queue systems
5. Enable fast development with strong typing

## Decision

We chose **Go (Golang) 1.21+** with the **Fiber** web framework as our backend technology stack.

### Key Reasons:

1. **Performance**: Go's goroutines provide lightweight concurrency, allowing thousands of simultaneous connections with minimal memory overhead.

2. **Fiber Framework**: Built on fasthttp (fastest HTTP engine for Go), Fiber provides Express.js-like syntax while maintaining exceptional performance (~0.1ms overhead per request).

3. **Static Binary Compilation**: Single binary deployment simplifies containerization and reduces runtime dependencies.

4. **Strong Typing**: Catches errors at compile time, reducing production bugs in critical financial operations.

5. **Excellent Standard Library**: net/http, encoding/json, crypto, and other packages reduce external dependencies.

6. **Goroutine-based Concurrency**: Natural fit for handling multiple payment webhook callbacks and background job processing.

## Consequences

### Positive

- **High throughput**: Benchmarks show 100,000+ requests/second capacity
- **Low memory footprint**: ~10MB per instance vs 100MB+ for Node.js/Java
- **Fast compilation**: Rapid development cycles
- **Simple deployment**: Single binary with no runtime dependencies
- **Strong community**: Active ecosystem with mature libraries (GORM, Redis, Asynq)
- **Built-in testing**: go test provides comprehensive testing framework

### Negative

- **Steeper learning curve**: Team members familiar with Node.js/Python need Go training
- **Verbose error handling**: Explicit error checking increases code verbosity
- **Limited generics** (pre-1.18): Some patterns require interface{} type assertions
- **Smaller talent pool**: Fewer Go developers in Nigerian market compared to Node.js

### Neutral

- Go modules for dependency management (standard but different from npm/pip)
- Different ORM patterns compared to ActiveRecord/Sequelize

## Alternatives Considered

### Alternative 1: Node.js with Express/Fastify

**Pros**: Large talent pool, JavaScript ecosystem, rapid prototyping
**Cons**: Single-threaded event loop limits CPU-bound operations, higher memory usage, less suitable for financial applications requiring type safety

**Rejected because**: Performance concerns for high-concurrency financial transactions and lack of compile-time type safety.

### Alternative 2: Python with FastAPI

**Pros**: Excellent for data science integrations (credit scoring), async support, clean syntax
**Cons**: Slower execution speed, GIL limits true parallelism, less suitable for high-throughput APIs

**Rejected because**: Performance requirements for real-time payment processing.

### Alternative 3: Java with Spring Boot

**Pros**: Enterprise-grade, excellent for financial systems, strong typing
**Cons**: Higher memory footprint (300MB+ per instance), slower startup times, verbose boilerplate

**Rejected because**: Resource efficiency requirements for cost-effective cloud deployment.

### Alternative 4: Rust with Actix-web

**Pros**: Maximum performance, memory safety guarantees
**Cons**: Very steep learning curve, longer development time, smaller ecosystem

**Rejected because**: Development velocity requirements and limited talent availability.

## References

- [Go Official Website](https://golang.org/)
- [Fiber Documentation](https://docs.gofiber.io/)
- [TechEmpower Benchmarks](https://www.techempower.com/benchmarks/)
- [Uber's Go Guidelines](https://github.com/uber-go/guide)
