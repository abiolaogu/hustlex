# ADR-009: Clean Architecture for Mobile App

## Status

Accepted

## Date

2024-01-15

## Context

The HustleX mobile app requires:
- Multiple features (auth, wallet, gigs, savings, credit, profile, notifications)
- Offline-first capability
- Testable codebase
- Maintainable as team grows
- Clear separation between UI and business logic
- Reusable components across features

We needed an architecture that:
1. Scales with feature count
2. Enables parallel development
3. Isolates changes to specific layers
4. Facilitates unit testing
5. Supports dependency injection

## Decision

We adopted a **Feature-First Clean Architecture** approach with the following structure:

```
lib/
├── core/           # Shared infrastructure
│   ├── api/        # HTTP client
│   ├── config/     # App configuration
│   ├── constants/  # Colors, strings, typography
│   ├── di/         # Dependency injection (Riverpod)
│   ├── exceptions/ # Custom exceptions
│   ├── router/     # Navigation (GoRouter)
│   ├── services/   # Platform services
│   ├── storage/    # Local persistence
│   ├── theme/      # App theming
│   ├── utils/      # Helpers, extensions
│   └── widgets/    # Shared UI components
│
└── features/       # Feature modules
    └── {feature}/
        ├── data/
        │   ├── models/       # Data transfer objects
        │   ├── repositories/ # Repository implementations
        │   └── services/     # API service methods
        └── presentation/
            ├── providers/    # Riverpod state
            ├── screens/      # Full-screen UI
            └── widgets/      # Feature-specific widgets
```

### Key Reasons:

1. **Feature Isolation**: Each feature is self-contained, enabling parallel development.

2. **Clear Dependencies**: Core → Features, never Features → Features.

3. **Testability**: Each layer can be mocked independently.

4. **Scalability**: New features follow established patterns.

5. **Maintainability**: Changes localized to specific features.

## Consequences

### Positive

- **Parallel development**: Teams can work on different features simultaneously
- **Code reuse**: Core components shared across features
- **Testing**: Unit tests for repositories, widget tests for screens
- **Onboarding**: Clear patterns for new developers
- **Refactoring**: Changes in one feature don't affect others

### Negative

- **Boilerplate**: More files/folders than simple architecture
- **Learning curve**: Team must understand architecture patterns
- **Over-engineering risk**: Simple features may feel heavy
- **Navigation complexity**: Cross-feature navigation requires care

### Neutral

- More directories than flat structure
- Requires discipline to maintain boundaries
- Initial setup time higher than MVP approach

## Architecture Layers

### 1. Data Layer

**Purpose**: Data access and transformation

```dart
// Model (DTO) - lib/features/wallet/data/models/wallet_model.dart
@freezed
class WalletModel with _$WalletModel {
  const factory WalletModel({
    required String id,
    required double balance,
    required double escrowBalance,
    required double savingsBalance,
    required String currency,
  }) = _WalletModel;

  factory WalletModel.fromJson(Map<String, dynamic> json) =>
      _$WalletModelFromJson(json);
}

// Repository - lib/features/wallet/data/repositories/wallet_repository.dart
class WalletRepository {
  final ApiClient _apiClient;

  WalletRepository(this._apiClient);

  Future<WalletModel> getWallet() async {
    final response = await _apiClient.get('/wallet');
    return WalletModel.fromJson(response.data);
  }

  Future<List<TransactionModel>> getTransactions({int page = 1}) async {
    final response = await _apiClient.get('/wallet/transactions', params: {'page': page});
    return (response.data as List)
        .map((e) => TransactionModel.fromJson(e))
        .toList();
  }

  Future<void> initiateDeposit(double amount) async {
    await _apiClient.post('/wallet/deposit', data: {'amount': amount});
  }
}
```

### 2. Presentation Layer

**Purpose**: UI and state management

```dart
// Provider - lib/features/wallet/presentation/providers/wallet_provider.dart
@riverpod
class WalletNotifier extends _$WalletNotifier {
  @override
  FutureOr<WalletState> build() async {
    final wallet = await ref.watch(walletRepositoryProvider).getWallet();
    return WalletState(wallet: wallet, transactions: []);
  }

  Future<void> refresh() async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      final wallet = await ref.read(walletRepositoryProvider).getWallet();
      final transactions = await ref.read(walletRepositoryProvider).getTransactions();
      return WalletState(wallet: wallet, transactions: transactions);
    });
  }

  Future<void> deposit(double amount) async {
    await ref.read(walletRepositoryProvider).initiateDeposit(amount);
    await refresh(); // Reload after deposit
  }
}

// Screen - lib/features/wallet/presentation/screens/wallet_screen.dart
class WalletScreen extends ConsumerWidget {
  const WalletScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final walletAsync = ref.watch(walletNotifierProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Wallet')),
      body: walletAsync.when(
        data: (state) => WalletContent(state: state),
        loading: () => const WalletSkeleton(),
        error: (error, stack) => ErrorView(
          error: error,
          onRetry: () => ref.invalidate(walletNotifierProvider),
        ),
      ),
    );
  }
}
```

### 3. Core Layer

**Purpose**: Shared infrastructure

```dart
// API Client - lib/core/api/api_client.dart
class ApiClient {
  final Dio _dio;

  ApiClient(this._dio);

  Future<Response> get(String path, {Map<String, dynamic>? params}) async {
    try {
      return await _dio.get(path, queryParameters: params);
    } on DioException catch (e) {
      throw _handleError(e);
    }
  }

  AppException _handleError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
        return NetworkException('Connection timeout');
      case DioExceptionType.badResponse:
        return ApiException.fromResponse(e.response);
      default:
        return UnknownException(e.message);
    }
  }
}

// Shared Widget - lib/core/widgets/primary_button.dart
class PrimaryButton extends StatelessWidget {
  final String text;
  final VoidCallback? onPressed;
  final bool isLoading;

  const PrimaryButton({
    required this.text,
    this.onPressed,
    this.isLoading = false,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: isLoading ? null : onPressed,
      child: isLoading
          ? const SizedBox(
              width: 20,
              height: 20,
              child: CircularProgressIndicator(strokeWidth: 2),
            )
          : Text(text),
    );
  }
}
```

## Dependency Flow

```
┌─────────────────────────────────────────────────────┐
│                      Widgets                         │
│  (ConsumerWidget, HookConsumerWidget, StatelessWidget)│
└───────────────────────┬─────────────────────────────┘
                        │ ref.watch / ref.read
                        ▼
┌─────────────────────────────────────────────────────┐
│                     Providers                        │
│    (StateNotifier, FutureProvider, StreamProvider)   │
└───────────────────────┬─────────────────────────────┘
                        │ depends on
                        ▼
┌─────────────────────────────────────────────────────┐
│                    Repositories                      │
│          (Data access, caching, transformation)      │
└───────────────────────┬─────────────────────────────┘
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────┐
│                      Services                        │
│        (ApiClient, LocalStorage, PlatformServices)   │
└───────────────────────┬─────────────────────────────┘
                        │
          ┌─────────────┴─────────────┐
          ▼                           ▼
┌─────────────────┐         ┌─────────────────┐
│   Backend API   │         │  Local Storage  │
│   (REST/HTTP)   │         │  (Hive/SQLite)  │
└─────────────────┘         └─────────────────┘
```

## Feature List

| Feature | Screens | Purpose |
|---------|---------|---------|
| `auth` | Login, OTP, Register, PIN Setup | User authentication |
| `home` | Dashboard | Main navigation hub |
| `wallet` | Wallet, Deposit, Transfer, Withdraw, History | Financial management |
| `gigs` | Browse, Details, Create, Proposals, Contract | Marketplace |
| `savings` | Circles List, Circle Details, Create, Contribute | Ajo/Esusu savings |
| `credit` | Score Dashboard, Loans, Apply, Repay | Credit and loans |
| `profile` | Profile, Edit, Settings, Bank Accounts | User management |
| `notifications` | Notification List, Settings | Push notifications |

## Alternatives Considered

### Alternative 1: Layer-First Architecture

```
lib/
├── data/
│   ├── wallet/
│   ├── gigs/
│   └── ...
├── domain/
│   ├── wallet/
│   └── ...
└── presentation/
    ├── wallet/
    └── ...
```

**Rejected because**: Feature changes require touching multiple top-level directories; harder to maintain feature boundaries.

### Alternative 2: Simple MVC/MVVM

**Rejected because**: Doesn't scale well with feature count; business logic often ends up in UI layer.

### Alternative 3: Redux/BLoC Heavy

**Rejected because**: Too much boilerplate for a mobile app; Riverpod provides simpler patterns with same benefits.

## References

- [Flutter Clean Architecture Guide](https://resocoder.com/flutter-clean-architecture-tdd/)
- [Riverpod Architecture](https://codewithandrea.com/articles/flutter-app-architecture-riverpod-introduction/)
- [Feature-First vs Layer-First](https://codewithandrea.com/articles/flutter-project-structure/)
