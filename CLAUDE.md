# HustleX Pro - Complete Platform

## Mission
Transform HustleX into a unified platform combining gig marketplace, fintech, and diaspora services.

## Tech Stack
- Backend: Go, Hasura, YugabyteDB, n8n, DragonflyDB
- Flutter: Ferry, Riverpod, GoRouter (Clean Architecture)
- Web: Refine v4, Ant Design, React Query
- Android: Jetpack Compose, Apollo Kotlin, Hilt, Orbit MVI
- iOS: SwiftUI, Apollo iOS, TCA
- VAS: Project Catalyst (USSD, SMS), IVR

## Principles
- XP: Small commits, continuous integration
- DDD: Bounded contexts, domain events
- TDD: Tests before implementation

## Structure
apps/api/ - Go backend
apps/consumer-app/ - Flutter
apps/admin-web/ - React admin
apps/android/ - Native Android
apps/ios/ - Native iOS
backend/hasura/ - GraphQL config
global-fintech/ - Submodule
infrastructure/ - Docker, K8s
