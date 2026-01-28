/// =============================================================================
/// API EXCEPTION
/// =============================================================================

class ApiException implements Exception {
  final String message;
  final String? code;
  final int? statusCode;
  final dynamic errors;

  ApiException({
    required this.message,
    this.code,
    this.statusCode,
    this.errors,
  });

  @override
  String toString() => message;

  bool get isUnauthorized => statusCode == 401 || code == 'UNAUTHORIZED';
  bool get isForbidden => statusCode == 403 || code == 'FORBIDDEN';
  bool get isNotFound => statusCode == 404 || code == 'NOT_FOUND';
  bool get isValidationError => statusCode == 422 || code == 'VALIDATION_ERROR';
  bool get isRateLimited => statusCode == 429 || code == 'RATE_LIMITED';
  bool get isServerError => (statusCode ?? 0) >= 500;
  bool get isNetworkError => code == 'NO_CONNECTION' || code == 'TIMEOUT';

  Map<String, List<String>> get validationErrors {
    if (errors == null || errors is! Map) return {};
    final result = <String, List<String>>{};
    (errors as Map).forEach((key, value) {
      if (value is List) {
        result[key.toString()] = value.map((e) => e.toString()).toList();
      } else {
        result[key.toString()] = [value.toString()];
      }
    });
    return result;
  }

  String? getFieldError(String field) {
    final fieldErrors = validationErrors[field];
    return fieldErrors?.isNotEmpty == true ? fieldErrors!.first : null;
  }
}

/// =============================================================================
/// AUTH EXCEPTION
/// =============================================================================

class AuthException implements Exception {
  final String message;
  final AuthErrorType type;

  AuthException({
    required this.message,
    required this.type,
  });

  @override
  String toString() => message;
}

enum AuthErrorType {
  invalidCredentials,
  sessionExpired,
  accountLocked,
  accountNotVerified,
  invalidOtp,
  otpExpired,
  maxOtpAttempts,
  invalidPin,
  pinLocked,
  tokenRefreshFailed,
  unknown,
}

/// =============================================================================
/// WALLET EXCEPTION
/// =============================================================================

class WalletException implements Exception {
  final String message;
  final WalletErrorType type;

  WalletException({
    required this.message,
    required this.type,
  });

  @override
  String toString() => message;
}

enum WalletErrorType {
  insufficientBalance,
  invalidPin,
  pinLocked,
  dailyLimitExceeded,
  weeklyLimitExceeded,
  invalidAmount,
  invalidAccount,
  transferFailed,
  withdrawalFailed,
  depositFailed,
  accountNotVerified,
  unknown,
}

/// =============================================================================
/// SAVINGS EXCEPTION
/// =============================================================================

class SavingsException implements Exception {
  final String message;
  final SavingsErrorType type;

  SavingsException({
    required this.message,
    required this.type,
  });

  @override
  String toString() => message;
}

enum SavingsErrorType {
  circleFull,
  circleNotFound,
  alreadyMember,
  notMember,
  invalidInviteCode,
  circleStarted,
  contributionNotAllowed,
  payoutPending,
  unknown,
}

/// =============================================================================
/// LOAN EXCEPTION
/// =============================================================================

class LoanException implements Exception {
  final String message;
  final LoanErrorType type;

  LoanException({
    required this.message,
    required this.type,
  });

  @override
  String toString() => message;
}

enum LoanErrorType {
  notEligible,
  insufficientCreditScore,
  activeLoanExists,
  invalidAmount,
  invalidTenure,
  loanNotFound,
  repaymentFailed,
  alreadyRepaid,
  unknown,
}

/// =============================================================================
/// GIG EXCEPTION
/// =============================================================================

class GigException implements Exception {
  final String message;
  final GigErrorType type;

  GigException({
    required this.message,
    required this.type,
  });

  @override
  String toString() => message;
}

enum GigErrorType {
  gigNotFound,
  proposalNotFound,
  contractNotFound,
  alreadyApplied,
  notAuthorized,
  gigClosed,
  invalidBudget,
  escrowFailed,
  unknown,
}

/// =============================================================================
/// CACHE EXCEPTION
/// =============================================================================

class CacheException implements Exception {
  final String message;

  CacheException({required this.message});

  @override
  String toString() => message;
}

/// =============================================================================
/// VALIDATION EXCEPTION
/// =============================================================================

class ValidationException implements Exception {
  final String field;
  final String message;

  ValidationException({
    required this.field,
    required this.message,
  });

  @override
  String toString() => message;
}
