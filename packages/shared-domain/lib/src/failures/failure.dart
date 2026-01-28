import 'package:freezed_annotation/freezed_annotation.dart';

part 'failure.freezed.dart';

/// Application-level failures for use cases and repositories
@freezed
class Failure with _$Failure {
  // Server failures
  const factory Failure.serverError([String? message]) = ServerError;
  const factory Failure.networkError() = NetworkError;
  const factory Failure.timeout() = TimeoutError;

  // Auth failures
  const factory Failure.unauthenticated() = Unauthenticated;
  const factory Failure.unauthorized() = Unauthorized;
  const factory Failure.invalidCredentials() = InvalidCredentials;
  const factory Failure.sessionExpired() = SessionExpired;
  const factory Failure.invalidOtp() = InvalidOtp;
  const factory Failure.otpExpired() = OtpExpired;
  const factory Failure.invalidPin() = InvalidPin;
  const factory Failure.accountLocked() = AccountLocked;

  // Validation failures
  const factory Failure.validationError(Map<String, String> errors) =
      ValidationError;
  const factory Failure.invalidInput(String field, String message) =
      InvalidInput;

  // Business logic failures
  const factory Failure.insufficientFunds() = InsufficientFunds;
  const factory Failure.transactionFailed(String reason) = TransactionFailed;
  const factory Failure.walletLocked() = WalletLocked;
  const factory Failure.dailyLimitExceeded() = DailyLimitExceeded;
  const factory Failure.withdrawalLimitExceeded() = WithdrawalLimitExceeded;
  const factory Failure.transferLimitExceeded() = TransferLimitExceeded;
  const factory Failure.minimumAmountNotMet(int minimum) = MinimumAmountNotMet;

  // Gig failures
  const factory Failure.gigNotFound() = GigNotFound;
  const factory Failure.gigAlreadyClosed() = GigAlreadyClosed;
  const factory Failure.proposalAlreadySubmitted() = ProposalAlreadySubmitted;
  const factory Failure.cannotBidOnOwnGig() = CannotBidOnOwnGig;

  // Savings circle failures
  const factory Failure.circleNotFound() = CircleNotFound;
  const factory Failure.circleFull() = CircleFull;
  const factory Failure.alreadyInCircle() = AlreadyInCircle;
  const factory Failure.invalidInviteCode() = InvalidInviteCode;
  const factory Failure.contributionNotDue() = ContributionNotDue;

  // Credit failures
  const factory Failure.loanNotEligible(String reason) = LoanNotEligible;
  const factory Failure.activeLoanExists() = ActiveLoanExists;
  const factory Failure.loanNotFound() = LoanNotFound;

  // Resource failures
  const factory Failure.notFound(String resource) = NotFound;
  const factory Failure.alreadyExists(String resource) = AlreadyExists;
  const factory Failure.conflict(String message) = Conflict;

  // Cache failures
  const factory Failure.cacheError() = CacheError;
  const factory Failure.cacheExpired() = CacheExpired;

  // Permission failures
  const factory Failure.permissionDenied(String permission) = PermissionDenied;
  const factory Failure.biometricNotAvailable() = BiometricNotAvailable;
  const factory Failure.biometricNotEnrolled() = BiometricNotEnrolled;

  // Unknown
  const factory Failure.unknown([String? message]) = UnknownError;
}

extension FailureX on Failure {
  /// Returns a human-readable error message
  String get message => when(
        serverError: (msg) => msg ?? 'Server error occurred. Please try again.',
        networkError: () => 'No internet connection. Please check your network.',
        timeout: () => 'Request timed out. Please try again.',
        unauthenticated: () => 'Please log in to continue.',
        unauthorized: () => 'You don\'t have permission to perform this action.',
        invalidCredentials: () => 'Invalid credentials. Please try again.',
        sessionExpired: () => 'Your session has expired. Please log in again.',
        invalidOtp: () => 'Invalid OTP code. Please check and try again.',
        otpExpired: () => 'OTP has expired. Please request a new one.',
        invalidPin: () => 'Invalid PIN. Please try again.',
        accountLocked: () =>
            'Your account has been locked. Please contact support.',
        validationError: (errors) => errors.values.first,
        invalidInput: (field, msg) => msg,
        insufficientFunds: () =>
            'Insufficient funds. Please top up your wallet.',
        transactionFailed: (reason) => reason,
        walletLocked: () => 'Your wallet is locked. Please contact support.',
        dailyLimitExceeded: () =>
            'Daily transaction limit exceeded. Try again tomorrow.',
        withdrawalLimitExceeded: () => 'Withdrawal limit exceeded.',
        transferLimitExceeded: () => 'Transfer limit exceeded.',
        minimumAmountNotMet: (min) => 'Minimum amount is ₦${min / 100}',
        gigNotFound: () => 'Gig not found.',
        gigAlreadyClosed: () => 'This gig is no longer accepting proposals.',
        proposalAlreadySubmitted: () =>
            'You have already submitted a proposal for this gig.',
        cannotBidOnOwnGig: () => 'You cannot bid on your own gig.',
        circleNotFound: () => 'Savings circle not found.',
        circleFull: () => 'This savings circle is full.',
        alreadyInCircle: () => 'You are already a member of this circle.',
        invalidInviteCode: () => 'Invalid invite code.',
        contributionNotDue: () => 'Your contribution is not due yet.',
        loanNotEligible: (reason) => reason,
        activeLoanExists: () =>
            'You already have an active loan. Please repay it first.',
        loanNotFound: () => 'Loan not found.',
        notFound: (resource) => '$resource not found.',
        alreadyExists: (resource) => '$resource already exists.',
        conflict: (msg) => msg,
        cacheError: () => 'Unable to load cached data.',
        cacheExpired: () => 'Cached data has expired.',
        permissionDenied: (permission) => '$permission permission denied.',
        biometricNotAvailable: () =>
            'Biometric authentication is not available on this device.',
        biometricNotEnrolled: () =>
            'No biometrics enrolled. Please set up fingerprint or face ID.',
        unknown: (msg) => msg ?? 'An unexpected error occurred.',
      );

  /// Returns true if this is a network-related error
  bool get isNetworkError => this is NetworkError || this is TimeoutError;

  /// Returns true if the user needs to re-authenticate
  bool get requiresAuth =>
      this is Unauthenticated ||
      this is SessionExpired ||
      this is Unauthorized;

  /// Returns true if this is a validation error
  bool get isValidationError =>
      this is ValidationError || this is InvalidInput;

  /// Returns true if this is a transient error that might succeed on retry
  bool get isRetryable =>
      this is NetworkError || this is TimeoutError || this is ServerError;
}
