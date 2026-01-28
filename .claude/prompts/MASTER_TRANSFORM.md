Execute this COMPLETE transformation of HustleX Pro. Do ALL phases in sequence, committing after each.

## PHASE 1: MONOREPO SETUP
1. Create directory structure: apps/{api,consumer-app,provider-app,admin-web,android,ios}, packages/, backend/hasura/, infrastructure/
2. Move existing code appropriately
3. Add Global-FinTech as submodule: `git submodule add https://github.com/abiolaogu/Global-FinTech.git global-fintech`
4. Create infrastructure/fintech/client.go wrapper
COMMIT: "🏗️ Restructure as monorepo with Global-FinTech"

## PHASE 2: DIASPORA FEATURES
1. Create apps/api/internal/domain/diaspora/entity/ with:
   - beneficiary.go (Beneficiary struct with relationships, linked wallet)
   - diaspora_booking.go (dual currency, FX rate locking, verification)
   - remittance.go (multi-currency transfer, recurring, status tracking)
2. Update wallet.go for multi-currency (NGN, GBP, USD, EUR, CAD, GHS, KES)
3. Create fx_service.go with corridor rates (1.5-2.5% spreads)
COMMIT: "✨ Add diaspora features (beneficiaries, remote booking, remittances)"

## PHASE 3: INFRASTRUCTURE
1. Update docker-compose.yml:
   - Replace redis with DragonflyDB (drop-in, same port 6379)
   - Add YugabyteDB (port 5433)
   - Add Hasura (port 8080)
   - Add n8n (port 5678)
2. Create backend/hasura/config.yaml
3. Create backend/hasura/migrations/default/1_init/up.sql with FULL schema:
   - users, profiles, wallets, services, bookings, beneficiaries
   - remittances, savings_circles, circle_members, transactions
   - notifications, audit_logs
   - All indexes and triggers
COMMIT: "🚀 Infrastructure: DragonflyDB, Hasura, YugabyteDB, n8n"

## PHASE 4: VAS INTEGRATION
1. Create apps/api/internal/domain/vas/entity/:
   - ussd_session.go (state machine, menu stack)
   - sms_message.go (templates, delivery status)
   - ivr_session.go (flow definitions)
2. Create apps/api/internal/domain/vas/service/:
   - ussd_service.go (full menu: balance, send, book, savings)
   - notification_service.go (multi-channel routing)
3. Create apps/api/internal/infrastructure/vas/catalyst/:
   - client.go (Catalyst API integration)
   - handler.go (webhooks for USSD, SMS, IVR)
4. Define USSD menu tree (*384*123#)
5. Define SMS templates for all notification types
COMMIT: "📱 VAS integration via Project Catalyst"

## PHASE 5: MULTI-PLATFORM CLIENTS
1. apps/admin-web/: Refine v4 + Ant Design 5.x
   - package.json, codegen.ts, data-provider.ts, App.tsx
   - Resources: users, bookings, remittances, services
2. apps/consumer-app/ (Flutter): Ferry + Riverpod + GoRouter
   - Update pubspec.yaml
   - Create lib/graphql/client.dart
   - Clean Architecture: domain/, data/, presentation/
3. apps/android/: Jetpack Compose + Apollo Kotlin + Hilt + Orbit MVI
   - build.gradle.kts with Apollo plugin
   - NetworkModule.kt, WalletViewModel.kt, WalletScreen.kt
4. apps/ios/: SwiftUI + Apollo iOS + TCA
   - WalletFeature.swift (TCA reducer)
   - WalletView.swift (SwiftUI)
COMMIT: "📱 Multi-platform clients (Web, Flutter, Android, iOS)"

## PHASE 6: AUTONOMOUS PIPELINE
1. Create .github/workflows/codegen.yml:
   - Trigger on backend/hasura/metadata/** changes
   - Introspect schema
   - Generate types for all platforms
   - Quality gates
   - Auto PR
COMMIT: "🤖 Autonomous code generation pipeline"

## FINAL
1. Run: make test
2. Tag: git tag -a v1.0.0-alpha -m "HustleX Pro unified platform"
3. Create comprehensive README.md
COMMIT: "🎉 HustleX Pro v1.0.0-alpha ready"
