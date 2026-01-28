# HustleX Mobile App

A Flutter-based mobile application for the HustleX gig economy and social savings platform targeting the Nigerian market.

## Features

- **Gig Marketplace**: Find work or hire talent
- **Ajo/Esusu Savings**: Social savings circles (rotational and fixed-target)
- **Digital Wallet**: Send, receive, deposit, and withdraw funds
- **Credit Building**: Alternative credit scoring based on platform activity
- **Microloans**: Access to credit based on your HustleX score

## Tech Stack

- **Framework**: Flutter 3.x
- **State Management**: Riverpod
- **Navigation**: GoRouter
- **HTTP Client**: Dio
- **Local Storage**: Hive, Flutter Secure Storage
- **Payments**: Paystack
- **Push Notifications**: Firebase Cloud Messaging

## Project Structure

```
lib/
├── core/
│   ├── api/              # API client and interceptors
│   ├── constants/        # App constants, colors, typography
│   ├── exceptions/       # Custom exception classes
│   ├── providers/        # Global providers (auth, etc.)
│   ├── utils/            # Helper utilities
│   └── widgets/          # Shared widgets (buttons, cards, inputs)
├── features/
│   ├── auth/             # Authentication (login, OTP, registration)
│   ├── home/             # Home screen and main shell
│   ├── gigs/             # Gig marketplace
│   ├── savings/          # Savings circles
│   ├── wallet/           # Wallet and transactions
│   ├── credit/           # Credit score and loans
│   ├── profile/          # User profile
│   └── notifications/    # Push notifications
├── router/               # App routing configuration
└── main.dart             # App entry point
```

Each feature follows clean architecture:
```
feature/
├── data/
│   ├── models/           # Data models
│   ├── repositories/     # Repository implementations
│   └── datasources/      # API/local data sources
├── domain/
│   ├── entities/         # Business entities
│   ├── repositories/     # Repository interfaces
│   └── usecases/         # Business logic
└── presentation/
    ├── screens/          # UI screens
    ├── widgets/          # Feature-specific widgets
    └── controllers/      # State management
```

## Setup Instructions

### Prerequisites

- Flutter SDK 3.16.0 or higher
- Dart SDK 3.2.0 or higher
- Android Studio / Xcode
- Node.js (for build tools)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/billyronks/hustlex.git
cd hustlex/mobile
```

2. Install dependencies:
```bash
flutter pub get
```

3. Create environment file:
```bash
cp .env.example .env
```

4. Configure environment variables in `.env`:
```env
API_BASE_URL=https://api.hustlex.app
PAYSTACK_PUBLIC_KEY=pk_live_xxxxx
```

5. Generate code (models, routes):
```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

6. Run the app:
```bash
# Development
flutter run --flavor development

# Staging
flutter run --flavor staging

# Production
flutter run --flavor production
```

## Build Commands

### Android

```bash
# Debug APK
flutter build apk --flavor development --debug

# Release APK
flutter build apk --flavor production --release

# App Bundle (for Play Store)
flutter build appbundle --flavor production --release
```

### iOS

```bash
# Debug build
flutter build ios --flavor development --debug

# Release build
flutter build ios --flavor production --release

# Archive for App Store
flutter build ipa --flavor production --release
```

## Environment Configuration

The app supports three environments:

| Environment | API URL | App ID Suffix |
|-------------|---------|---------------|
| Development | localhost:8080 | .dev |
| Staging | staging-api.hustlex.app | .staging |
| Production | api.hustlex.app | (none) |

## Key Dependencies

```yaml
# State Management
flutter_riverpod: ^2.4.9
riverpod_annotation: ^2.3.3

# Navigation
go_router: ^13.1.0

# Networking
dio: ^5.4.0
connectivity_plus: ^5.0.2

# Storage
flutter_secure_storage: ^9.0.0
hive_flutter: ^1.1.0

# UI
google_fonts: ^6.1.0
shimmer: ^3.0.0
fl_chart: ^0.66.0
lottie: ^3.0.0

# Payments
flutter_paystack: ^1.0.7

# Firebase
firebase_core: ^2.25.4
firebase_messaging: ^14.7.15
```

## Architecture Notes

### State Management

The app uses Riverpod for state management:
- `StateNotifier` for complex state
- `AsyncNotifier` for async operations
- `Provider` for simple dependencies

### API Layer

- Dio client with interceptors for auth, logging, and errors
- Automatic token refresh on 401
- Request/response logging in debug mode

### Navigation

- GoRouter for declarative routing
- Shell routes for bottom navigation
- Deep linking support

### Security

- Secure storage for sensitive data (tokens, PIN)
- Certificate pinning for production
- Biometric authentication support

## Testing

```bash
# Unit tests
flutter test

# Integration tests
flutter test integration_test/

# With coverage
flutter test --coverage
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## License

Copyright © 2026 BillyRonks Global Limited. All rights reserved.
