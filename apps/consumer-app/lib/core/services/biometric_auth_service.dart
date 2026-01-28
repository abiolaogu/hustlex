import 'dart:io';

import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:local_auth/local_auth.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

/// Biometric authentication type
enum BiometricType {
  fingerprint,
  face,
  iris,
  none,
}

/// Biometric availability status
class BiometricStatus {
  final bool isAvailable;
  final bool isEnrolled;
  final List<BiometricType> availableTypes;
  final String? errorMessage;

  const BiometricStatus({
    this.isAvailable = false,
    this.isEnrolled = false,
    this.availableTypes = const [],
    this.errorMessage,
  });

  bool get canUseBiometrics => isAvailable && isEnrolled;
  
  BiometricType get primaryType {
    if (availableTypes.isEmpty) return BiometricType.none;
    // Prefer fingerprint over face
    if (availableTypes.contains(BiometricType.fingerprint)) {
      return BiometricType.fingerprint;
    }
    return availableTypes.first;
  }

  String get biometricName {
    switch (primaryType) {
      case BiometricType.fingerprint:
        return 'Fingerprint';
      case BiometricType.face:
        return Platform.isIOS ? 'Face ID' : 'Face Recognition';
      case BiometricType.iris:
        return 'Iris';
      case BiometricType.none:
        return 'Biometrics';
    }
  }
}

/// Authentication result
class BiometricAuthResult {
  final bool success;
  final String? errorMessage;
  final BiometricAuthErrorType? errorType;

  const BiometricAuthResult({
    required this.success,
    this.errorMessage,
    this.errorType,
  });

  factory BiometricAuthResult.success() {
    return const BiometricAuthResult(success: true);
  }

  factory BiometricAuthResult.failed({
    String? message,
    BiometricAuthErrorType? errorType,
  }) {
    return BiometricAuthResult(
      success: false,
      errorMessage: message ?? 'Authentication failed',
      errorType: errorType,
    );
  }

  factory BiometricAuthResult.cancelled() {
    return const BiometricAuthResult(
      success: false,
      errorMessage: 'Authentication cancelled',
      errorType: BiometricAuthErrorType.cancelled,
    );
  }
}

/// Authentication error types
enum BiometricAuthErrorType {
  cancelled,
  notAvailable,
  notEnrolled,
  lockedOut,
  permanentlyLockedOut,
  passcodeNotSet,
  unknown,
}

/// Biometric authentication service
class BiometricAuthService {
  final LocalAuthentication _auth;

  BiometricAuthService({LocalAuthentication? auth})
      : _auth = auth ?? LocalAuthentication();

  /// Check if device supports biometrics
  Future<BiometricStatus> checkBiometricStatus() async {
    try {
      // Check if device supports biometrics
      final isAvailable = await _auth.canCheckBiometrics;
      final isDeviceSupported = await _auth.isDeviceSupported();

      if (!isAvailable || !isDeviceSupported) {
        return const BiometricStatus(
          isAvailable: false,
          errorMessage: 'Biometric authentication not available on this device',
        );
      }

      // Get available biometrics
      final availableBiometrics = await _auth.getAvailableBiometrics();
      
      if (availableBiometrics.isEmpty) {
        return const BiometricStatus(
          isAvailable: true,
          isEnrolled: false,
          errorMessage: 'No biometrics enrolled on this device',
        );
      }

      // Map to our types
      final types = availableBiometrics.map(_mapBiometricType).toList();

      return BiometricStatus(
        isAvailable: true,
        isEnrolled: true,
        availableTypes: types,
      );
    } on PlatformException catch (e) {
      return BiometricStatus(
        isAvailable: false,
        errorMessage: e.message ?? 'Failed to check biometric status',
      );
    } catch (e) {
      return BiometricStatus(
        isAvailable: false,
        errorMessage: e.toString(),
      );
    }
  }

  /// Map platform biometric type to our type
  BiometricType _mapBiometricType(BiometricType type) {
    // The local_auth package uses same enum names
    switch (type.name) {
      case 'fingerprint':
        return BiometricType.fingerprint;
      case 'face':
        return BiometricType.face;
      case 'iris':
        return BiometricType.iris;
      default:
        return BiometricType.none;
    }
  }

  /// Authenticate user with biometrics
  Future<BiometricAuthResult> authenticate({
    String reason = 'Authenticate to access HustleX',
    bool biometricOnly = false,
    bool stickyAuth = true,
  }) async {
    try {
      final status = await checkBiometricStatus();
      
      if (!status.canUseBiometrics) {
        return BiometricAuthResult.failed(
          message: status.errorMessage ?? 'Biometrics not available',
          errorType: BiometricAuthErrorType.notAvailable,
        );
      }

      final authenticated = await _auth.authenticate(
        localizedReason: reason,
        options: AuthenticationOptions(
          biometricOnly: biometricOnly,
          stickyAuth: stickyAuth,
          useErrorDialogs: true,
        ),
      );

      if (authenticated) {
        return BiometricAuthResult.success();
      } else {
        return BiometricAuthResult.failed();
      }
    } on PlatformException catch (e) {
      return _handlePlatformException(e);
    } catch (e) {
      return BiometricAuthResult.failed(
        message: e.toString(),
        errorType: BiometricAuthErrorType.unknown,
      );
    }
  }

  /// Authenticate for sensitive operations (transactions, settings)
  Future<BiometricAuthResult> authenticateForTransaction({
    required double amount,
    String? recipient,
  }) async {
    String reason = 'Authenticate to confirm transaction';
    if (recipient != null) {
      reason = 'Authenticate to send ₦${amount.toStringAsFixed(2)} to $recipient';
    } else {
      reason = 'Authenticate to confirm ₦${amount.toStringAsFixed(2)} transaction';
    }

    return authenticate(
      reason: reason,
      biometricOnly: true, // Force biometric only for transactions
      stickyAuth: true,
    );
  }

  /// Authenticate to view sensitive data
  Future<BiometricAuthResult> authenticateToView(String dataType) async {
    return authenticate(
      reason: 'Authenticate to view $dataType',
      biometricOnly: false,
      stickyAuth: false,
    );
  }

  /// Authenticate for login
  Future<BiometricAuthResult> authenticateForLogin() async {
    return authenticate(
      reason: 'Authenticate to login to HustleX',
      biometricOnly: false,
      stickyAuth: true,
    );
  }

  /// Handle platform exceptions
  BiometricAuthResult _handlePlatformException(PlatformException e) {
    switch (e.code) {
      case auth_error.notAvailable:
        return BiometricAuthResult.failed(
          message: 'Biometrics not available',
          errorType: BiometricAuthErrorType.notAvailable,
        );
      case auth_error.notEnrolled:
        return BiometricAuthResult.failed(
          message: 'No biometrics enrolled',
          errorType: BiometricAuthErrorType.notEnrolled,
        );
      case auth_error.lockedOut:
        return BiometricAuthResult.failed(
          message: 'Too many attempts. Please try again later.',
          errorType: BiometricAuthErrorType.lockedOut,
        );
      case auth_error.permanentlyLockedOut:
        return BiometricAuthResult.failed(
          message: 'Biometrics permanently locked. Please use PIN/password.',
          errorType: BiometricAuthErrorType.permanentlyLockedOut,
        );
      case auth_error.passcodeNotSet:
        return BiometricAuthResult.failed(
          message: 'Please set up a device passcode first',
          errorType: BiometricAuthErrorType.passcodeNotSet,
        );
      default:
        if (e.message?.toLowerCase().contains('cancel') ?? false) {
          return BiometricAuthResult.cancelled();
        }
        return BiometricAuthResult.failed(
          message: e.message ?? 'Authentication failed',
          errorType: BiometricAuthErrorType.unknown,
        );
    }
  }

  /// Stop authentication (cancel ongoing)
  Future<void> stopAuthentication() async {
    await _auth.stopAuthentication();
  }
}

/// Biometric auth service provider
final biometricAuthServiceProvider = Provider<BiometricAuthService>((ref) {
  return BiometricAuthService();
});

/// Biometric status provider
final biometricStatusProvider = FutureProvider<BiometricStatus>((ref) async {
  final service = ref.watch(biometricAuthServiceProvider);
  return await service.checkBiometricStatus();
});

/// Can use biometrics provider
final canUseBiometricsProvider = Provider<bool>((ref) {
  final status = ref.watch(biometricStatusProvider);
  return status.when(
    data: (s) => s.canUseBiometrics,
    loading: () => false,
    error: (_, __) => false,
  );
});

/// Biometric state for UI
class BiometricState {
  final bool isAuthenticating;
  final BiometricAuthResult? result;

  const BiometricState({
    this.isAuthenticating = false,
    this.result,
  });

  BiometricState copyWith({
    bool? isAuthenticating,
    BiometricAuthResult? result,
  }) {
    return BiometricState(
      isAuthenticating: isAuthenticating ?? this.isAuthenticating,
      result: result,
    );
  }

  bool get isSuccess => result?.success ?? false;
  bool get isFailed => result != null && !result!.success;
  bool get isCancelled => result?.errorType == BiometricAuthErrorType.cancelled;
}

/// Biometric state notifier
class BiometricNotifier extends StateNotifier<BiometricState> {
  final BiometricAuthService _authService;

  BiometricNotifier(this._authService) : super(const BiometricState());

  /// Authenticate
  Future<bool> authenticate({String? reason}) async {
    state = state.copyWith(isAuthenticating: true, result: null);

    final result = await _authService.authenticate(
      reason: reason ?? 'Authenticate to continue',
    );

    state = state.copyWith(
      isAuthenticating: false,
      result: result,
    );

    return result.success;
  }

  /// Authenticate for transaction
  Future<bool> authenticateForTransaction({
    required double amount,
    String? recipient,
  }) async {
    state = state.copyWith(isAuthenticating: true, result: null);

    final result = await _authService.authenticateForTransaction(
      amount: amount,
      recipient: recipient,
    );

    state = state.copyWith(
      isAuthenticating: false,
      result: result,
    );

    return result.success;
  }

  /// Authenticate for login
  Future<bool> authenticateForLogin() async {
    state = state.copyWith(isAuthenticating: true, result: null);

    final result = await _authService.authenticateForLogin();

    state = state.copyWith(
      isAuthenticating: false,
      result: result,
    );

    return result.success;
  }

  /// Cancel authentication
  Future<void> cancel() async {
    await _authService.stopAuthentication();
    state = state.copyWith(
      isAuthenticating: false,
      result: BiometricAuthResult.cancelled(),
    );
  }

  /// Reset state
  void reset() {
    state = const BiometricState();
  }
}

/// Biometric notifier provider
final biometricProvider = StateNotifierProvider<BiometricNotifier, BiometricState>((ref) {
  final service = ref.watch(biometricAuthServiceProvider);
  return BiometricNotifier(service);
});
