import 'dart:async';
import 'dart:io';

import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:local_auth/local_auth.dart';
import 'package:local_auth/error_codes.dart' as auth_error;
import 'package:logger/logger.dart';

import '../storage/secure_storage.dart';

/// Biometric type enum
enum BiometricType {
  fingerprint,
  faceId,
  iris,
  none,
}

/// Biometric authentication result
class BiometricResult {
  final bool success;
  final String? message;
  final BiometricErrorType? errorType;

  BiometricResult({
    required this.success,
    this.message,
    this.errorType,
  });

  factory BiometricResult.success() {
    return BiometricResult(success: true);
  }

  factory BiometricResult.failure(String message, [BiometricErrorType? type]) {
    return BiometricResult(
      success: false,
      message: message,
      errorType: type,
    );
  }

  factory BiometricResult.notAvailable() {
    return BiometricResult(
      success: false,
      message: 'Biometric authentication is not available on this device',
      errorType: BiometricErrorType.notAvailable,
    );
  }

  factory BiometricResult.notEnrolled() {
    return BiometricResult(
      success: false,
      message: 'No biometrics enrolled on this device',
      errorType: BiometricErrorType.notEnrolled,
    );
  }

  factory BiometricResult.cancelled() {
    return BiometricResult(
      success: false,
      message: 'Authentication cancelled',
      errorType: BiometricErrorType.cancelled,
    );
  }

  factory BiometricResult.lockedOut() {
    return BiometricResult(
      success: false,
      message: 'Too many attempts. Please try again later',
      errorType: BiometricErrorType.lockedOut,
    );
  }

  factory BiometricResult.permanentlyLockedOut() {
    return BiometricResult(
      success: false,
      message: 'Biometrics permanently locked. Use PIN instead',
      errorType: BiometricErrorType.permanentlyLockedOut,
    );
  }
}

/// Biometric error types
enum BiometricErrorType {
  notAvailable,
  notEnrolled,
  cancelled,
  failed,
  lockedOut,
  permanentlyLockedOut,
  passcodeNotSet,
  unknown,
}

/// Biometric authentication service
class BiometricService {
  final Logger _logger = Logger();
  final LocalAuthentication _localAuth;
  final SecureStorage _secureStorage;

  static const _biometricEnabledKey = 'biometric_enabled';
  static const _biometricAuthTimeKey = 'biometric_last_auth';

  BiometricService({
    required SecureStorage secureStorage,
    LocalAuthentication? localAuth,
  })  : _secureStorage = secureStorage,
        _localAuth = localAuth ?? LocalAuthentication();

  /// Check if device supports biometrics
  Future<bool> isDeviceSupported() async {
    try {
      return await _localAuth.isDeviceSupported();
    } catch (e) {
      _logger.e('Error checking device support', error: e);
      return false;
    }
  }

  /// Check if biometrics can be used (device supported + enrolled)
  Future<bool> canCheckBiometrics() async {
    try {
      return await _localAuth.canCheckBiometrics;
    } catch (e) {
      _logger.e('Error checking biometrics availability', error: e);
      return false;
    }
  }

  /// Get available biometric types
  Future<List<BiometricType>> getAvailableBiometrics() async {
    try {
      final available = await _localAuth.getAvailableBiometrics();
      return available.map((type) {
        switch (type) {
          case BiometricType.fingerprint:
            return BiometricType.fingerprint;
          case BiometricType.face:
            return BiometricType.faceId;
          case BiometricType.iris:
            return BiometricType.iris;
          default:
            return BiometricType.none;
        }
      }).where((type) => type != BiometricType.none).toList();
    } catch (e) {
      _logger.e('Error getting available biometrics', error: e);
      return [];
    }
  }

  /// Get primary biometric type
  Future<BiometricType> getPrimaryBiometricType() async {
    final available = await getAvailableBiometrics();
    
    if (available.isEmpty) return BiometricType.none;
    
    // Prefer Face ID on iOS, fingerprint elsewhere
    if (Platform.isIOS && available.contains(BiometricType.faceId)) {
      return BiometricType.faceId;
    }
    
    if (available.contains(BiometricType.fingerprint)) {
      return BiometricType.fingerprint;
    }
    
    return available.first;
  }

  /// Get biometric type display name
  String getBiometricTypeName(BiometricType type) {
    switch (type) {
      case BiometricType.fingerprint:
        return 'Fingerprint';
      case BiometricType.faceId:
        return Platform.isIOS ? 'Face ID' : 'Face Recognition';
      case BiometricType.iris:
        return 'Iris';
      case BiometricType.none:
        return 'None';
    }
  }

  /// Authenticate with biometrics
  Future<BiometricResult> authenticate({
    String reason = 'Please authenticate to continue',
    bool useErrorDialogs = true,
    bool stickyAuth = true,
    bool sensitiveTransaction = true,
    bool biometricOnly = false,
  }) async {
    try {
      // Check if biometrics are available
      final canCheck = await canCheckBiometrics();
      if (!canCheck) {
        return BiometricResult.notAvailable();
      }

      // Check if any biometrics are enrolled
      final available = await getAvailableBiometrics();
      if (available.isEmpty) {
        return BiometricResult.notEnrolled();
      }

      // Attempt authentication
      final authenticated = await _localAuth.authenticate(
        localizedReason: reason,
        options: AuthenticationOptions(
          useErrorDialogs: useErrorDialogs,
          stickyAuth: stickyAuth,
          sensitiveTransaction: sensitiveTransaction,
          biometricOnly: biometricOnly,
        ),
      );

      if (authenticated) {
        _logger.i('Biometric authentication successful');
        await _recordAuthTime();
        return BiometricResult.success();
      } else {
        _logger.w('Biometric authentication failed');
        return BiometricResult.failure(
          'Authentication failed',
          BiometricErrorType.failed,
        );
      }
    } on PlatformException catch (e) {
      _logger.e('Biometric authentication error', error: e);
      return _handlePlatformException(e);
    } catch (e) {
      _logger.e('Unexpected biometric error', error: e);
      return BiometricResult.failure(
        'An unexpected error occurred',
        BiometricErrorType.unknown,
      );
    }
  }

  /// Handle platform-specific authentication errors
  BiometricResult _handlePlatformException(PlatformException e) {
    switch (e.code) {
      case auth_error.notAvailable:
        return BiometricResult.notAvailable();
      case auth_error.notEnrolled:
        return BiometricResult.notEnrolled();
      case auth_error.passcodeNotSet:
        return BiometricResult.failure(
          'Please set up a device passcode first',
          BiometricErrorType.passcodeNotSet,
        );
      case auth_error.lockedOut:
        return BiometricResult.lockedOut();
      case auth_error.permanentlyLockedOut:
        return BiometricResult.permanentlyLockedOut();
      default:
        if (e.message?.contains('cancel') == true) {
          return BiometricResult.cancelled();
        }
        return BiometricResult.failure(
          e.message ?? 'Authentication failed',
          BiometricErrorType.unknown,
        );
    }
  }

  /// Check if biometric is enabled for app
  Future<bool> isBiometricEnabled() async {
    final enabled = await _secureStorage.read(key: _biometricEnabledKey);
    return enabled == 'true';
  }

  /// Enable/disable biometric authentication
  Future<void> setBiometricEnabled(bool enabled) async {
    await _secureStorage.write(
      key: _biometricEnabledKey,
      value: enabled.toString(),
    );
    _logger.i('Biometric ${enabled ? "enabled" : "disabled"}');
  }

  /// Record last successful authentication time
  Future<void> _recordAuthTime() async {
    await _secureStorage.write(
      key: _biometricAuthTimeKey,
      value: DateTime.now().toIso8601String(),
    );
  }

  /// Get last authentication time
  Future<DateTime?> getLastAuthTime() async {
    final timeStr = await _secureStorage.read(key: _biometricAuthTimeKey);
    if (timeStr == null) return null;
    return DateTime.tryParse(timeStr);
  }

  /// Check if recent authentication is valid (within duration)
  Future<bool> isRecentAuthValid({Duration validity = const Duration(minutes: 5)}) async {
    final lastAuth = await getLastAuthTime();
    if (lastAuth == null) return false;
    return DateTime.now().difference(lastAuth) < validity;
  }

  /// Authenticate for sensitive transaction
  Future<BiometricResult> authenticateForTransaction({
    required String transactionDescription,
    required double amount,
  }) async {
    final formattedAmount = '₦${amount.toStringAsFixed(2)}';
    return authenticate(
      reason: 'Authenticate to $transactionDescription of $formattedAmount',
      sensitiveTransaction: true,
      stickyAuth: true,
    );
  }

  /// Authenticate for app unlock
  Future<BiometricResult> authenticateForAppUnlock() async {
    return authenticate(
      reason: 'Unlock HustleX',
      sensitiveTransaction: false,
      stickyAuth: true,
    );
  }

  /// Authenticate for PIN change
  Future<BiometricResult> authenticateForPinChange() async {
    return authenticate(
      reason: 'Authenticate to change your PIN',
      sensitiveTransaction: true,
      biometricOnly: true,
    );
  }

  /// Stop authentication
  Future<void> stopAuthentication() async {
    await _localAuth.stopAuthentication();
  }
}

/// Biometric state for UI
class BiometricState {
  final bool isAvailable;
  final bool isEnabled;
  final BiometricType primaryType;
  final List<BiometricType> availableTypes;
  final bool isAuthenticating;
  final BiometricResult? lastResult;

  const BiometricState({
    this.isAvailable = false,
    this.isEnabled = false,
    this.primaryType = BiometricType.none,
    this.availableTypes = const [],
    this.isAuthenticating = false,
    this.lastResult,
  });

  BiometricState copyWith({
    bool? isAvailable,
    bool? isEnabled,
    BiometricType? primaryType,
    List<BiometricType>? availableTypes,
    bool? isAuthenticating,
    BiometricResult? lastResult,
  }) {
    return BiometricState(
      isAvailable: isAvailable ?? this.isAvailable,
      isEnabled: isEnabled ?? this.isEnabled,
      primaryType: primaryType ?? this.primaryType,
      availableTypes: availableTypes ?? this.availableTypes,
      isAuthenticating: isAuthenticating ?? this.isAuthenticating,
      lastResult: lastResult,
    );
  }

  String get primaryTypeName {
    switch (primaryType) {
      case BiometricType.fingerprint:
        return 'Fingerprint';
      case BiometricType.faceId:
        return 'Face ID';
      case BiometricType.iris:
        return 'Iris';
      case BiometricType.none:
        return 'None';
    }
  }
}

/// Biometric state notifier
class BiometricNotifier extends StateNotifier<BiometricState> {
  final BiometricService _service;

  BiometricNotifier(this._service) : super(const BiometricState()) {
    _initialize();
  }

  /// Initialize biometric state
  Future<void> _initialize() async {
    final isAvailable = await _service.canCheckBiometrics();
    final isEnabled = await _service.isBiometricEnabled();
    final availableTypes = await _service.getAvailableBiometrics();
    final primaryType = await _service.getPrimaryBiometricType();

    state = state.copyWith(
      isAvailable: isAvailable,
      isEnabled: isEnabled,
      availableTypes: availableTypes,
      primaryType: primaryType,
    );
  }

  /// Refresh biometric status
  Future<void> refresh() async {
    await _initialize();
  }

  /// Enable biometric authentication
  Future<BiometricResult> enable() async {
    state = state.copyWith(isAuthenticating: true);

    // Require biometric auth to enable
    final result = await _service.authenticate(
      reason: 'Authenticate to enable biometric login',
    );

    if (result.success) {
      await _service.setBiometricEnabled(true);
      state = state.copyWith(
        isEnabled: true,
        isAuthenticating: false,
        lastResult: result,
      );
    } else {
      state = state.copyWith(
        isAuthenticating: false,
        lastResult: result,
      );
    }

    return result;
  }

  /// Disable biometric authentication
  Future<void> disable() async {
    await _service.setBiometricEnabled(false);
    state = state.copyWith(isEnabled: false);
  }

  /// Authenticate user
  Future<BiometricResult> authenticate({String? reason}) async {
    if (!state.isEnabled) {
      return BiometricResult.failure('Biometric authentication is not enabled');
    }

    state = state.copyWith(isAuthenticating: true);

    final result = await _service.authenticate(
      reason: reason ?? 'Please authenticate',
    );

    state = state.copyWith(
      isAuthenticating: false,
      lastResult: result,
    );

    return result;
  }

  /// Authenticate for transaction
  Future<BiometricResult> authenticateForTransaction({
    required String description,
    required double amount,
  }) async {
    state = state.copyWith(isAuthenticating: true);

    final result = await _service.authenticateForTransaction(
      transactionDescription: description,
      amount: amount,
    );

    state = state.copyWith(
      isAuthenticating: false,
      lastResult: result,
    );

    return result;
  }

  /// Clear last result
  void clearResult() {
    state = state.copyWith(lastResult: null);
  }
}

/// Biometric service provider
final biometricServiceProvider = Provider<BiometricService>((ref) {
  final secureStorage = ref.watch(secureStorageProvider);
  return BiometricService(secureStorage: secureStorage);
});

/// Secure storage provider (duplicate for standalone usage)
final secureStorageProvider = Provider<SecureStorage>((ref) {
  return SecureStorage();
});

/// Biometric state provider
final biometricProvider =
    StateNotifierProvider<BiometricNotifier, BiometricState>((ref) {
  final service = ref.watch(biometricServiceProvider);
  return BiometricNotifier(service);
});

/// Biometric availability provider
final biometricAvailableProvider = FutureProvider<bool>((ref) async {
  final service = ref.watch(biometricServiceProvider);
  return await service.canCheckBiometrics();
});

/// Biometric enabled provider
final biometricEnabledProvider = FutureProvider<bool>((ref) async {
  final service = ref.watch(biometricServiceProvider);
  return await service.isBiometricEnabled();
});

/// Biometric type provider
final biometricTypeProvider = FutureProvider<BiometricType>((ref) async {
  final service = ref.watch(biometricServiceProvider);
  return await service.getPrimaryBiometricType();
});
