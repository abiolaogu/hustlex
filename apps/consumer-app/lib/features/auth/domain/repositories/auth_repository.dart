import 'package:dartz/dartz.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/value_objects/phone_number.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/user.dart';
import '../entities/session.dart';

/// Auth repository interface - defines the contract for authentication operations.
abstract class AuthRepository {
  /// Request OTP for phone number verification
  Future<Either<Failure, OtpSession>> requestOtp({
    required PhoneNumber phone,
    required OtpPurpose purpose,
  });

  /// Verify OTP code
  Future<Either<Failure, Session>> verifyOtp({
    required PhoneNumber phone,
    required String code,
  });

  /// Register a new user
  Future<Either<Failure, Session>> register({
    required PhoneNumber phone,
    required String firstName,
    required String lastName,
    String? email,
    String? referralCode,
  });

  /// Set user's transaction PIN
  Future<Either<Failure, User>> setPin(Pin pin);

  /// Change user's transaction PIN
  Future<Either<Failure, User>> changePin({
    required Pin currentPin,
    required Pin newPin,
  });

  /// Verify PIN (for sensitive operations)
  Future<Either<Failure, bool>> verifyPin(Pin pin);

  /// Refresh authentication tokens
  Future<Either<Failure, AuthTokens>> refreshTokens();

  /// Get current session
  Future<Either<Failure, Session>> getCurrentSession();

  /// Get current user
  Future<Either<Failure, User>> getCurrentUser();

  /// Update user profile
  Future<Either<Failure, User>> updateProfile({
    String? firstName,
    String? lastName,
    String? email,
    String? avatar,
  });

  /// Logout and clear session
  Future<Either<Failure, Unit>> logout();

  /// Check if user is authenticated
  Future<bool> isAuthenticated();

  /// Check if biometric auth is available
  Future<bool> isBiometricAvailable();

  /// Enable biometric authentication
  Future<Either<Failure, Unit>> enableBiometric();

  /// Disable biometric authentication
  Future<Either<Failure, Unit>> disableBiometric();

  /// Authenticate with biometric
  Future<Either<Failure, bool>> authenticateWithBiometric();
}
