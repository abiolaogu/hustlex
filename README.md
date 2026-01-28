# HustleX

> Nigerian Gig Economy + Social Savings (Ajo/Esusu) + Credit Building Super App

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Flutter Version](https://img.shields.io/badge/Flutter-3.16+-02569B?style=flat&logo=flutter)](https://flutter.dev)
[![License](https://img.shields.io/badge/License-Proprietary-red.svg)](LICENSE)

## Overview

HustleX is a comprehensive platform designed for the Nigerian gig economy, combining three key features:

1. **Gig Marketplace** - Connect freelancers with clients for short-term work
2. **Ajo/Esusu Savings Circles** - Traditional social savings digitized
3. **Alternative Credit Scoring** - Build credit through platform activity

## Architecture

```
hustlex/
├── backend/           # Go API server
│   ├── cmd/          # Application entry points
│   │   ├── api/      # Main API server
│   │   └── worker/   # Background job processor
│   └── internal/     # Private packages
│       ├── config/   # Configuration
│       ├── database/ # Database connection
│       ├── handlers/ # HTTP handlers
│       ├── jobs/     # Background jobs
│       ├── middleware/
│       ├── models/   # Database models
│       └── services/ # Business logic
├── mobile/           # Flutter mobile app
│   └── lib/
│       ├── core/     # Shared utilities
│       └── features/ # Feature modules
├── docs/             # Documentation
└── k8s/              # Kubernetes manifests
```

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Chi Router
- **Database**: PostgreSQL 15
- **Cache/Queue**: Redis
- **Background Jobs**: Asynq
- **Payments**: Paystack

### Mobile
- **Framework**: Flutter 3.16+
- **State Management**: Riverpod
- **Navigation**: GoRouter
- **Local Storage**: Hive, Secure Storage

### Infrastructure
- **Container**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: GitHub Actions

## Getting Started

### Prerequisites

- Go 1.21+
- Flutter 3.16+
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)
- Redis 7+ (or use Docker)

### Quick Start

1. **Clone the repository**
```bash
git clone https://github.com/billyronks/hustlex.git
cd hustlex
```

2. **Start backend services**
```bash
cd backend
cp .env.example .env  # Edit with your values
docker-compose up -d
go run cmd/api/main.go
```

3. **Start mobile app**
```bash
cd mobile
cp .env.example .env  # Edit with your values
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
flutter run
```

## Key Features

### Gig Marketplace
- Post gigs with budget and requirements
- Submit proposals with cover letters
- Escrow-based payments for security
- Rating and review system

### Savings Circles (Ajo/Esusu)
- **Rotational (Ajo)**: Members take turns receiving pooled contributions
- **Fixed Target**: Group saves toward individual goals
- Automatic contribution reminders
- Payout scheduling

### Credit Building
- Alternative credit scoring algorithm
- Factors: payment history, savings consistency, gig performance
- Access to microloans based on score
- Credit score tracking over time

### Digital Wallet
- Fund via bank transfer or card (Paystack)
- Transfer to other users
- Withdraw to bank account
- Transaction history

## API Documentation

### Authentication
```
POST /api/v1/auth/otp/request     # Request OTP
POST /api/v1/auth/otp/verify      # Verify OTP
POST /api/v1/auth/register        # Register new user
POST /api/v1/auth/pin/set         # Set transaction PIN
```

### Gigs
```
GET    /api/v1/gigs               # List gigs
POST   /api/v1/gigs               # Create gig
GET    /api/v1/gigs/{id}          # Get gig details
POST   /api/v1/gigs/{id}/proposals # Submit proposal
```

### Savings
```
GET    /api/v1/savings/circles    # List circles
POST   /api/v1/savings/circles    # Create circle
POST   /api/v1/savings/circles/{id}/join
POST   /api/v1/savings/circles/{id}/contribute
```

### Wallet
```
GET    /api/v1/wallet             # Get wallet
POST   /api/v1/wallet/deposit     # Initiate deposit
POST   /api/v1/wallet/transfer    # Transfer to user
POST   /api/v1/wallet/withdraw    # Withdraw to bank
GET    /api/v1/wallet/transactions
```

### Credit
```
GET    /api/v1/credit/score       # Get credit score
GET    /api/v1/credit/history     # Score history
GET    /api/v1/credit/loans       # List loans
POST   /api/v1/credit/loans/apply # Apply for loan
POST   /api/v1/credit/loans/{id}/repay
```

## Environment Variables

### Backend
```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=hustlex
DB_PASSWORD=secret
DB_NAME=hustlex

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-secret-key

# Paystack
PAYSTACK_SECRET_KEY=sk_test_xxx
```

### Mobile
```env
API_BASE_URL=http://localhost:8080
PAYSTACK_PUBLIC_KEY=pk_test_xxx
```

## Development

### Running Tests

```bash
# Backend
cd backend
go test ./...

# Mobile
cd mobile
flutter test
```

### Database Migrations

```bash
cd backend
go run cmd/api/main.go migrate up
```

### Code Generation (Mobile)

```bash
cd mobile
flutter pub run build_runner build --delete-conflicting-outputs
```

## Deployment

### Docker

```bash
# Build images
docker build -t hustlex-api ./backend
docker build -t hustlex-mobile ./mobile

# Run with compose
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes

```bash
kubectl apply -k k8s/overlays/production
```

## Project Status

### Completed ✅
- Product requirements documentation
- System architecture design
- Backend API handlers
- Background job system
- Mobile app foundation
- All main screens (UI)

### In Progress 🚧
- Payment integration (Paystack)
- Push notifications (FCM)
- SMS gateway integration

### Planned 📋
- Admin dashboard
- Analytics system
- API documentation (Swagger)
- E2E testing
- CI/CD pipelines

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open a Pull Request

## License

Copyright © 2026 BillyRonks Global Limited. All rights reserved.

This is proprietary software. Unauthorized copying, modification, distribution, or use is strictly prohibited.

## Contact

- **Company**: BillyRonks Global Limited
- **Email**: tech@billyronks.com
- **Website**: https://hustlex.app
