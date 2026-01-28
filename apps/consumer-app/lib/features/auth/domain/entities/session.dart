import 'package:freezed_annotation/freezed_annotation.dart';
import 'user.dart';

part 'session.freezed.dart';

/// Authentication tokens
@freezed
class AuthTokens with _$AuthTokens {
  const AuthTokens._();

  const factory AuthTokens({
    required String accessToken,
    required String refreshToken,
    required DateTime expiresAt,
  }) = _AuthTokens;

  /// Check if token has expired
  bool get isExpired => DateTime.now().isAfter(expiresAt);

  /// Check if token should be refreshed (within 5 minutes of expiry)
  bool get shouldRefresh {
    const buffer = Duration(minutes: 5);
    return DateTime.now().isAfter(expiresAt.subtract(buffer));
  }

  /// Time until token expires
  Duration get timeUntilExpiry => expiresAt.difference(DateTime.now());
}

/// User session containing user data and tokens
@freezed
class Session with _$Session {
  const Session._();

  const factory Session({
    required User user,
    required AuthTokens tokens,
    @Default(false) bool isNewUser,
  }) = _Session;

  /// Check if session is valid
  bool get isValid => !tokens.isExpired;

  /// Check if this is a fresh registration
  bool get requiresOnboarding => isNewUser || !user.hasSetPin;
}

/// OTP verification state
enum OtpPurpose {
  login,
  register,
  resetPin,
  verifyPhone,
}

/// OTP session for tracking verification flow
@freezed
class OtpSession with _$OtpSession {
  const OtpSession._();

  const factory OtpSession({
    required String phone,
    required OtpPurpose purpose,
    required DateTime expiresAt,
    @Default(0) int attempts,
    @Default(3) int maxAttempts,
  }) = _OtpSession;

  /// Check if OTP has expired
  bool get isExpired => DateTime.now().isAfter(expiresAt);

  /// Check if max attempts reached
  bool get isLocked => attempts >= maxAttempts;

  /// Time until OTP expires
  Duration get timeUntilExpiry => expiresAt.difference(DateTime.now());

  /// Remaining attempts
  int get remainingAttempts => maxAttempts - attempts;
}
