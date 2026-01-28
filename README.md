# HustleX Pro

> Unified platform combining gig marketplace, fintech, and diaspora services for the African market.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://golang.org/)
[![Flutter Version](https://img.shields.io/badge/flutter-3.16+-02569B.svg)](https://flutter.dev/)

## Overview

HustleX Pro is a comprehensive monorepo containing all components of the HustleX platform:

- **Gig Marketplace**: Connect service providers with consumers
- **Fintech**: Multi-currency wallets, payments, and transactions
- **Diaspora Services**: International remittances with competitive FX rates
- **Savings Circles**: Traditional Ajo/Esusu digitized

## Architecture

```
hustlex-pro/
├── apps/
│   ├── api/                 # Go backend (Clean Architecture)
│   ├── consumer-app/        # Flutter consumer app
│   ├── admin-web/           # React admin dashboard (Refine v4)
│   ├── android/             # Native Android (Jetpack Compose)
│   └── ios/                 # Native iOS (SwiftUI)
├── packages/
│   ├── shared-domain/       # Shared domain models (Dart)
│   └── shared-ui/           # Shared UI components (Dart)
├── backend/
│   └── hasura/              # Hasura GraphQL metadata & migrations
├── infrastructure/
│   ├── docker/              # Docker configurations
│   ├── vas/                 # VAS (USSD/SMS/IVR) - Project Catalyst
│   └── n8n/                 # Workflow automation
└── global-fintech/          # Fintech submodule
```

## Tech Stack

### Backend
- **Language**: Go 1.21+
- **Architecture**: Clean Architecture, DDD
- **Database**: PostgreSQL 16, YugabyteDB
- **Cache**: DragonflyDB (Redis-compatible)
- **GraphQL**: Hasura
- **Messaging**: RabbitMQ
- **Workflows**: n8n

### Mobile
- **Flutter**: Ferry (GraphQL), Riverpod, GoRouter
- **Android**: Jetpack Compose, Hilt, Apollo, Orbit MVI
- **iOS**: SwiftUI, TCA, Apollo iOS

### Web
- **Admin**: React, Refine v4, Ant Design

### Infrastructure
- **Container**: Docker, Kubernetes
- **Monitoring**: Prometheus, Grafana
- **VAS**: Project Catalyst (USSD, SMS, IVR)

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Flutter 3.16+
- Node.js 20+

### Development Setup

1. **Clone the repository with submodules**
   ```bash
   git clone --recurse-submodules https://github.com/abiolaogu/hustlex-pro.git
   cd hustlex-pro
   ```

2. **Start infrastructure services**
   ```bash
   docker-compose up -d postgres dragonfly hasura
   ```

3. **Run the API**
   ```bash
   cd apps/api
   go run cmd/server/main.go
   ```

4. **Run the Flutter app**
   ```bash
   cd apps/consumer-app/flutter
   flutter run
   ```

5. **Run the admin dashboard**
   ```bash
   cd apps/admin-web
   npm install
   npm run dev
   ```

### Full Stack

```bash
docker-compose up
```

Access services:
- API: http://localhost:8081
- Hasura Console: http://localhost:8080
- Admin Dashboard: http://localhost:3000
- n8n Workflows: http://localhost:5678
- Grafana: http://localhost:3001

## Features

### Gig Marketplace
- Service provider profiles and portfolios
- Service discovery and booking
- Real-time availability
- Escrow payments
- Reviews and ratings

### Fintech
- Multi-currency wallets (NGN, GBP, USD, EUR, CAD, GHS, KES)
- P2P transfers
- Bill payments
- Airtime/Data purchase
- Transaction history

### Diaspora Services
- International remittances
- Competitive FX rates with transparent pricing
- Multiple delivery methods (Bank, Mobile Wallet, Cash Pickup)
- Beneficiary management
- Recurring transfers
- Remote service booking

### Savings Circles (Ajo/Esusu)
- Digital traditional savings
- Automated contribution tracking
- Payout scheduling
- Default management

### VAS (Value Added Services)
- USSD banking (*347*123#)
- SMS notifications
- IVR support

## Supported Corridors

| From | To | Spread | Delivery |
|------|----|---------|----|
| 🇬🇧 GBP | 🇳🇬 NGN | 1.5% | 24h |
| 🇺🇸 USD | 🇳🇬 NGN | 1.75% | 24h |
| 🇪🇺 EUR | 🇳🇬 NGN | 1.75% | 24h |
| 🇨🇦 CAD | 🇳🇬 NGN | 2.0% | 24h |
| 🇬🇧 GBP | 🇬🇭 GHS | 2.0% | 24h |
| 🇬🇧 GBP | 🇰🇪 KES | 2.25% | 1-2 days |

## Development

### Generate a New Feature

```bash
# Backend feature
./scripts/generate-feature.sh backend notifications

# Flutter feature
./scripts/generate-feature.sh flutter savings

# Admin page
./scripts/generate-feature.sh admin reports

# Android feature
./scripts/generate-feature.sh android loyalty

# iOS feature
./scripts/generate-feature.sh ios loyalty
```

### Run Tests

```bash
# Backend
cd apps/api && go test ./...

# Flutter
cd apps/consumer-app/flutter && flutter test

# Admin
cd apps/admin-web && npm test
```

### Database Migrations

```bash
# Apply migrations
hasura migrate apply --database-name default

# Create new migration
hasura migrate create <name> --database-name default
```

## API Documentation

- GraphQL Playground: http://localhost:8080/console
- REST API (Swagger): http://localhost:8081/swagger

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feat/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feat/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Global-FinTech](https://github.com/abiolaogu/Global-FinTech) - Core fintech APIs
- [Project Catalyst](https://github.com/abiolaogu/Project-Catalyst) - VAS platform
- [Hasura](https://hasura.io/) - GraphQL engine
- [Refine](https://refine.dev/) - Admin framework

---

Built with ❤️ for the African diaspora

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
