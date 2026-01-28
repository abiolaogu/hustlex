import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:connectivity_plus/connectivity_plus.dart';

import '../api/api_client.dart';
import '../storage/secure_storage.dart';
import '../storage/local_cache_service.dart';
import '../../features/auth/data/repositories/auth_repository.dart';
import '../../features/gigs/data/repositories/gigs_repository.dart';
import '../../features/savings/data/repositories/savings_repository.dart';
import '../../features/wallet/data/repositories/wallet_repository.dart';
import '../../features/credit/data/repositories/credit_repository.dart';
import '../../features/profile/data/repositories/profile_repository.dart';
import '../../features/notifications/data/repositories/notification_repository.dart';

// ==================== CORE PROVIDERS ====================

/// Secure storage provider
final secureStorageProvider = Provider<SecureStorage>((ref) {
  return SecureStorage();
});

/// Local cache service provider
final localCacheServiceProvider = Provider<LocalCacheService>((ref) {
  return LocalCacheService();
});

/// Connectivity provider
final connectivityProvider = Provider<Connectivity>((ref) {
  return Connectivity();
});

/// Connectivity stream provider
final connectivityStreamProvider = StreamProvider<ConnectivityResult>((ref) {
  return ref.watch(connectivityProvider).onConnectivityChanged;
});

/// Is online provider
final isOnlineProvider = Provider<bool>((ref) {
  final connectivity = ref.watch(connectivityStreamProvider);
  return connectivity.when(
    data: (result) => result != ConnectivityResult.none,
    loading: () => true, // Assume online while loading
    error: (_, __) => true, // Assume online on error
  );
});

/// API client provider
final apiClientProvider = Provider<ApiClient>((ref) {
  final secureStorage = ref.watch(secureStorageProvider);
  return ApiClient(secureStorage: secureStorage);
});

// ==================== REPOSITORY PROVIDERS ====================

/// Auth repository provider
final authRepositoryProvider = Provider<AuthRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  final secureStorage = ref.watch(secureStorageProvider);
  return AuthRepository(apiClient: apiClient, secureStorage: secureStorage);
});

/// Gigs repository provider
final gigsRepositoryProvider = Provider<GigsRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return GigsRepository(apiClient: apiClient);
});

/// Savings repository provider
final savingsRepositoryProvider = Provider<SavingsRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return SavingsRepository(apiClient: apiClient);
});

/// Wallet repository provider
final walletRepositoryProvider = Provider<WalletRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return WalletRepository(apiClient: apiClient);
});

/// Credit repository provider
final creditRepositoryProvider = Provider<CreditRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return CreditRepository(apiClient: apiClient);
});

/// Profile repository provider
final profileRepositoryProvider = Provider<ProfileRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return ProfileRepository(apiClient: apiClient);
});

/// Notification repository provider
final notificationRepositoryProvider = Provider<NotificationRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return NotificationRepository(apiClient: apiClient);
});
