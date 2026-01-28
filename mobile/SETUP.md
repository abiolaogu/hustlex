# HustleX Mobile App Setup Guide

## Prerequisites

- Flutter SDK 3.2.0 or higher
- Dart 3.2.0 or higher
- Android Studio (for Android development)
- Xcode 15+ (for iOS development, macOS only)
- VS Code or Android Studio with Flutter extensions

## Quick Start

```bash
# 1. Navigate to mobile directory
cd mobile

# 2. Install dependencies
flutter pub get

# 3. Generate code (JSON serialization, Freezed models)
flutter pub run build_runner build --delete-conflicting-outputs

# 4. Copy environment file
cp .env.example .env
# Edit .env with your actual values

# 5. Run the app
flutter run
```

## Detailed Setup

### 1. Environment Configuration

Copy `.env.example` to `.env` and configure:

```env
API_BASE_URL=https://api.hustlex.ng/v1
PAYSTACK_PUBLIC_KEY=pk_test_xxxxx
```

### 2. Firebase Setup

#### Android
1. Create a Firebase project at [Firebase Console](https://console.firebase.google.com)
2. Add an Android app with package name: `ng.hustlex.app`
3. Download `google-services.json`
4. Place it in `android/app/google-services.json`

#### iOS
1. Add an iOS app in Firebase Console with bundle ID: `ng.hustlex.app`
2. Download `GoogleService-Info.plist`
3. Place it in `ios/Runner/GoogleService-Info.plist`
4. Add it to Xcode project (Runner target)

### 3. Paystack Setup

1. Create an account at [Paystack Dashboard](https://dashboard.paystack.com)
2. Get your test/live public key
3. Add to `.env` file:
   ```
   PAYSTACK_PUBLIC_KEY=pk_test_xxxxx
   ```

### 4. Code Generation

The app uses code generation for:
- JSON serialization (`json_serializable`)
- Immutable models (`freezed`)
- State management (`riverpod_generator`)

```bash
# One-time build
flutter pub run build_runner build --delete-conflicting-outputs

# Watch mode (rebuilds on file changes)
flutter pub run build_runner watch --delete-conflicting-outputs
```

### 5. Platform-Specific Setup

#### Android

1. **Signing Configuration** (for release builds):
   ```bash
   # Generate keystore
   keytool -genkey -v -keystore android/app/hustlex-release.keystore \
     -alias hustlex -keyalg RSA -keysize 2048 -validity 10000
   ```

2. Create `android/key.properties`:
   ```properties
   storePassword=your_store_password
   keyPassword=your_key_password
   keyAlias=hustlex
   storeFile=hustlex-release.keystore
   ```

3. **Minimum SDK**: 21 (Android 5.0)
4. **Target SDK**: 34 (Android 14)

#### iOS

1. Open `ios/Runner.xcworkspace` in Xcode
2. Set your Development Team in Signing & Capabilities
3. Enable required capabilities:
   - Push Notifications
   - Associated Domains (for deep linking)
   - Face ID usage (for biometrics)

4. **Minimum iOS Version**: 12.0

### 6. Running the App

```bash
# Development (debug mode)
flutter run

# Release mode
flutter run --release

# Specific device
flutter run -d <device_id>

# List available devices
flutter devices
```

### 7. Building for Release

#### Android APK
```bash
flutter build apk --release
# Output: build/app/outputs/flutter-apk/app-release.apk
```

#### Android App Bundle (for Play Store)
```bash
flutter build appbundle --release
# Output: build/app/outputs/bundle/release/app-release.aab
```

#### iOS
```bash
flutter build ios --release
# Then archive in Xcode for App Store submission
```

## Project Structure

```
mobile/
├── lib/
│   ├── core/                  # Core app infrastructure
│   │   ├── api/               # API client, interceptors
│   │   ├── config/            # App configuration
│   │   ├── constants/         # Colors, typography, strings
│   │   ├── di/                # Dependency injection
│   │   ├── exceptions/        # Custom exceptions
│   │   ├── providers/         # Core Riverpod providers
│   │   ├── repositories/      # Base repository
│   │   ├── router/            # Go Router setup
│   │   ├── services/          # API services
│   │   ├── storage/           # Local storage, cache
│   │   ├── theme/             # App theming
│   │   ├── utils/             # Utilities, helpers
│   │   └── widgets/           # Shared UI components
│   │
│   ├── features/              # Feature modules
│   │   ├── auth/              # Authentication
│   │   ├── credit/            # Credit scoring, loans
│   │   ├── gigs/              # Gig marketplace
│   │   ├── home/              # Dashboard
│   │   ├── notifications/     # In-app notifications
│   │   ├── profile/           # User profile, settings
│   │   ├── savings/           # Ajo/Esusu circles
│   │   ├── splash/            # Splash, onboarding
│   │   └── wallet/            # Wallet, transactions
│   │
│   ├── router/                # App router configuration
│   └── main.dart              # App entry point
│
├── assets/
│   ├── images/                # Static images
│   ├── icons/                 # App icons, SVGs
│   ├── animations/            # Lottie animations
│   └── fonts/                 # Custom fonts
│
├── android/                   # Android native code
├── ios/                       # iOS native code
├── test/                      # Test files
└── pubspec.yaml               # Dependencies
```

## Architecture

The app follows **Clean Architecture** with feature-based organization:

```
Feature/
├── data/
│   ├── models/          # Data models (Freezed)
│   ├── repositories/    # Repository implementations
│   └── services/        # API service classes
│
└── presentation/
    ├── providers/       # Riverpod state management
    ├── screens/         # UI screens
    └── widgets/         # Feature-specific widgets
```

### State Management

- **Riverpod** for dependency injection and state management
- **StateNotifier** for complex state with side effects
- **FutureProvider/StreamProvider** for async data
- **Freezed** for immutable state classes

### Data Flow

```
UI (Screens) 
    ↓ watches
Providers (StateNotifiers)
    ↓ calls
Repositories (Cache + API)
    ↓ uses
Services (API Client)
    ↓ to
Backend API
```

## Common Tasks

### Adding a New Feature

1. Create feature directory under `lib/features/`
2. Add data models in `data/models/`
3. Create repository in `data/repositories/`
4. Add service in `data/services/` (or `core/services/`)
5. Create providers in `presentation/providers/`
6. Build screens in `presentation/screens/`
7. Add routes in `lib/core/router/app_router.dart`

### Adding a New API Endpoint

1. Add method to relevant service class
2. Update repository to use the service
3. Update provider to expose the data
4. Connect to UI

### Updating Models

1. Modify model class
2. Run code generation:
   ```bash
   flutter pub run build_runner build --delete-conflicting-outputs
   ```

## Testing

```bash
# Run all tests
flutter test

# Run specific test file
flutter test test/features/wallet/wallet_test.dart

# Run with coverage
flutter test --coverage

# Generate coverage report
genhtml coverage/lcov.info -o coverage/html
```

## Debugging

### Common Issues

1. **Build errors after model changes**:
   ```bash
   flutter clean
   flutter pub get
   flutter pub run build_runner build --delete-conflicting-outputs
   ```

2. **iOS build fails**:
   ```bash
   cd ios && pod install && cd ..
   ```

3. **Android build fails**:
   - Check Android SDK version in `android/app/build.gradle`
   - Run `flutter doctor` to verify setup

### Debug Tools

- Flutter DevTools: `flutter pub global activate devtools`
- Riverpod DevTools extension for state inspection
- Network logging in API interceptors

## Performance Tips

1. Use `const` constructors where possible
2. Implement proper `==` and `hashCode` (Freezed handles this)
3. Use `ListView.builder` for long lists
4. Cache API responses appropriately
5. Lazy load images with `cached_network_image`
6. Profile with Flutter DevTools

## Security Checklist

- [ ] API keys not in version control (use .env)
- [ ] Sensitive data in secure storage
- [ ] Certificate pinning enabled
- [ ] ProGuard/R8 enabled for Android release
- [ ] Biometric authentication implemented
- [ ] PIN stored securely (hashed)
- [ ] Network security config for Android

## Deployment

### Google Play Store

1. Build app bundle: `flutter build appbundle --release`
2. Upload to Play Console
3. Complete store listing
4. Submit for review

### Apple App Store

1. Build: `flutter build ios --release`
2. Archive in Xcode
3. Upload to App Store Connect
4. Complete metadata
5. Submit for review

## Support

- Documentation: [docs.hustlex.ng](https://docs.hustlex.ng)
- API Reference: [api.hustlex.ng/docs](https://api.hustlex.ng/docs)
- Email: dev@hustlex.ng
