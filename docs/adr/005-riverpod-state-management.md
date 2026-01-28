# ADR-005: Riverpod for State Management

## Status

Accepted

## Date

2024-01-15

## Context

The HustleX mobile app requires robust state management for:
- Authentication state (logged in/out, user profile)
- Wallet balances and transaction history
- Gig listings and proposals
- Savings circle membership and contributions
- Credit scores and loan status
- Offline data synchronization

Requirements:
1. Compile-time safety (avoid runtime errors)
2. Testability (easy to mock dependencies)
3. Performance (efficient rebuilds)
4. Code generation support (reduce boilerplate)
5. Scalability (works for large apps)

## Decision

We chose **Riverpod 2.x** as our state management solution.

### Key Reasons:

1. **Compile-Time Safety**: Provider dependencies are validated at compile time, unlike Provider which can throw runtime errors.

2. **No BuildContext Required**: Providers can be accessed anywhere, including in services and repositories.

3. **Built-in Dependency Injection**: Natural DI pattern without additional packages.

4. **Automatic Disposal**: Resources are automatically cleaned up when no longer needed.

5. **Code Generation**: `riverpod_generator` reduces boilerplate with `@riverpod` annotations.

6. **DevTools Support**: Flutter DevTools integration for debugging state.

## Consequences

### Positive

- **Type safety**: Compiler catches missing dependencies
- **Testability**: Easy to override providers in tests
- **Performance**: Efficient widget rebuilds with `ref.watch`
- **Modularity**: Providers are independent, composable units
- **Caching**: Built-in caching with `keepAlive`
- **Async support**: Native `FutureProvider` and `StreamProvider`

### Negative

- **Learning curve**: Different mental model from traditional state management
- **Code generation dependency**: Requires `build_runner` for annotations
- **Migration complexity**: Harder to migrate from other solutions
- **Verbose for simple cases**: Overkill for very simple apps

### Neutral

- Different syntax from Provider package
- Requires understanding of provider lifecycles
- Community still growing compared to BLoC

## Provider Architecture

### Layer Structure

```
Widgets
   │
   ▼ ref.watch / ref.read
Providers (StateNotifier, FutureProvider, StreamProvider)
   │
   ▼ depends on
Repositories (data access abstraction)
   │
   ▼ uses
Services (API client, local storage)
   │
   ▼ connects to
Backend API / Local Database
```

### Provider Types Used

| Type | Use Case | Example |
|------|----------|---------|
| `Provider` | Static dependencies | `apiClientProvider` |
| `StateNotifierProvider` | Mutable state | `walletProvider` |
| `FutureProvider` | Async one-time data | `userProfileProvider` |
| `StreamProvider` | Real-time data | `connectivityStreamProvider` |
| `StateProvider` | Simple mutable values | `selectedTabProvider` |

### Code Example

```dart
// Repository provider
@riverpod
WalletRepository walletRepository(WalletRepositoryRef ref) {
  return WalletRepository(ref.watch(apiClientProvider));
}

// State provider with notifier
@riverpod
class WalletNotifier extends _$WalletNotifier {
  @override
  AsyncValue<Wallet> build() {
    return const AsyncValue.loading();
  }

  Future<void> loadWallet() async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() =>
      ref.read(walletRepositoryProvider).getWallet()
    );
  }

  Future<void> deposit(double amount) async {
    // Optimistic update + API call
  }
}

// Widget usage
class WalletScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final wallet = ref.watch(walletNotifierProvider);

    return wallet.when(
      data: (data) => WalletCard(balance: data.balance),
      loading: () => const LoadingIndicator(),
      error: (e, st) => ErrorWidget(error: e),
    );
  }
}
```

## Alternatives Considered

### Alternative 1: BLoC/Cubit

**Pros**: Event-driven, clear separation, large community
**Cons**: Verbose boilerplate, requires BuildContext for access, steeper learning curve

**Rejected because**: More boilerplate for similar functionality; Riverpod's DI is more intuitive.

### Alternative 2: Provider

**Pros**: Simple, widely used, Flutter team supported
**Cons**: Runtime errors possible, requires BuildContext, limited compile-time checks

**Rejected because**: Riverpod is the evolution of Provider with compile-time safety.

### Alternative 3: GetX

**Pros**: All-in-one solution, simple syntax, minimal boilerplate
**Cons**: Magic singletons, hard to test, controversial architecture

**Rejected because**: Testing difficulties and lack of explicit dependency management.

### Alternative 4: MobX

**Pros**: Reactive programming, minimal boilerplate, observable pattern
**Cons**: Magic annotations, code generation required, less Flutter-native

**Rejected because**: Riverpod feels more native to Flutter ecosystem.

## References

- [Riverpod Official Documentation](https://riverpod.dev/)
- [Riverpod GitHub Repository](https://github.com/rrousselGit/riverpod)
- [Flutter Riverpod Best Practices](https://codewithandrea.com/articles/flutter-state-management-riverpod/)
- [Riverpod Generator](https://pub.dev/packages/riverpod_generator)
