# HustleX Pro

> Unified platform for gig economy financial services and AI-powered service marketplace

[![CI/CD](https://github.com/abiolaogu/hustlex-pro/actions/workflows/ci.yml/badge.svg)](https://github.com/abiolaogu/hustlex-pro/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

HustleX Pro is a comprehensive monorepo combining:

- **HustleX** - Gig economy financial services (wallet, savings circles, micro-loans)
- **VendorPlatform** - AI-powered service marketplace with intelligent matching

## Architecture

```
hustlex-pro/
├── apps/
│   ├── api/                 # Go backend (unified API)
│   ├── consumer-app/        # Flutter app for consumers
│   ├── provider-app/        # Flutter app for service providers
│   ├── admin-web/           # React admin dashboard
│   └── recommendation/      # Python ML recommendation service
├── packages/
│   ├── shared-domain/       # Shared Dart domain models
│   ├── shared-ui/           # Shared Flutter widgets
│   ├── go-common/           # Shared Go packages
│   └── proto/               # gRPC/protobuf definitions
├── infrastructure/
│   ├── docker/              # Docker configurations
│   ├── k8s/                 # Kubernetes manifests
│   └── terraform/           # Infrastructure as Code
└── docs/
    ├── architecture/        # Architecture decisions
    ├── api/                 # API documentation
    └── business/            # Business documentation
```

## Features

### Consumer App
- Digital wallet with deposits, withdrawals, transfers
- Gig marketplace for finding freelance work
- Ajo/Esusu savings circles (traditional rotating savings)
- Micro-loans with credit scoring
- Real-time notifications

### Provider App
- Service listing management
- Booking and scheduling
- Earnings tracking
- Customer communication
- Performance analytics

### AI Features
- Intelligent service provider matching
- Dynamic pricing recommendations
- Fraud detection
- Credit risk assessment
- Demand forecasting

## Quick Start

### Prerequisites

- Go 1.21+
- Flutter 3.16+
- Docker & Docker Compose
- Node.js 18+ (for admin dashboard)
- Python 3.11+ (for ML services)

### Setup

```bash
# Clone the repository
git clone https://github.com/abiolaogu/hustlex-pro.git
cd hustlex-pro

# Install all dependencies
make setup

# Start infrastructure services
make docker-up

# Run database migrations
make db-migrate

# Start the API server
make run-api

# In separate terminals:
make run-consumer    # Start consumer Flutter app
make run-provider    # Start provider Flutter app
```

### Using Melos (Flutter)

```bash
# Bootstrap all Flutter packages
melos bootstrap

# Run tests across all packages
melos run test

# Generate code (freezed, json_serializable)
melos run generate

# Analyze all packages
melos run analyze
```

## Development

### API Development

```bash
cd apps/api

# Run with hot reload
air

# Run tests
go test -v ./...

# Build binary
go build -o bin/server cmd/server/main.go
```

### Flutter Development

```bash
# Consumer app
cd apps/consumer-app
flutter run

# Provider app
cd apps/provider-app
flutter run
```

### ML Service Development

```bash
cd apps/recommendation

# Install dependencies
pip install -r requirements.txt

# Run service
uvicorn main:app --reload --port 8081
```

## Testing

```bash
# Run all tests
make test

# Run specific tests
make test-api        # API tests
make test-flutter    # Flutter tests
make test-ml         # ML service tests
```

## Docker

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down

# Rebuild images
make docker-build
```

## API Documentation

- **OpenAPI Spec**: `docs/api/openapi.yaml`
- **Postman Collection**: `docs/api/postman/`
- **Swagger UI**: http://localhost:8080/swagger (when running)

## Project Structure

| Directory | Description |
|-----------|-------------|
| `apps/api` | Go backend with Clean Architecture |
| `apps/consumer-app` | Consumer-facing Flutter application |
| `apps/provider-app` | Provider-facing Flutter application |
| `apps/admin-web` | React-based admin dashboard |
| `apps/recommendation` | Python ML recommendation service |
| `packages/shared-domain` | Shared Dart domain models |
| `packages/shared-ui` | Shared Flutter UI components |
| `packages/go-common` | Shared Go utilities |
| `packages/proto` | Protocol buffer definitions |
| `infrastructure/` | DevOps configurations |
| `docs/` | Documentation |

## Contributing

1. Create a feature branch from `main`
2. Make your changes
3. Run tests: `make test`
4. Run linting: `make lint`
5. Submit a pull request

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go, Gin, GORM, PostgreSQL |
| Mobile | Flutter, Riverpod, Freezed |
| Web | React, TypeScript, TailwindCSS |
| ML | Python, FastAPI, scikit-learn, PyTorch |
| Infrastructure | Docker, Kubernetes, Terraform |
| CI/CD | GitHub Actions |
| Monitoring | Prometheus, Grafana |

## License

MIT License - see [LICENSE](LICENSE) for details.

---

Built with ❤️ for the Nigerian gig economy
