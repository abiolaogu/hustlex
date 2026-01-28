import 'package:freezed_annotation/freezed_annotation.dart';

part 'value_failures.freezed.dart';

/// Value object validation failures
@freezed
class ValueFailure<T> with _$ValueFailure<T> {
  // String failures
  const factory ValueFailure.empty({required T failedValue}) = _Empty<T>;
  const factory ValueFailure.exceedingLength({
    required T failedValue,
    required int max,
  }) = _ExceedingLength<T>;
  const factory ValueFailure.tooShort({
    required T failedValue,
    required int min,
  }) = _TooShort<T>;
  const factory ValueFailure.invalidFormat({required T failedValue}) =
      _InvalidFormat<T>;

  // Email failures
  const factory ValueFailure.invalidEmail({required T failedValue}) =
      _InvalidEmail<T>;

  // Phone failures
  const factory ValueFailure.invalidPhoneNumber({required T failedValue}) =
      _InvalidPhoneNumber<T>;

  // Money failures
  const factory ValueFailure.invalidMoney({required T failedValue}) =
      _InvalidMoney<T>;
  const factory ValueFailure.negativeAmount({required T failedValue}) =
      _NegativeAmount<T>;
  const factory ValueFailure.exceedsLimit({
    required T failedValue,
    required int limit,
  }) = _ExceedsLimit<T>;

  // PIN failures
  const factory ValueFailure.invalidPin({required T failedValue}) =
      _InvalidPin<T>;
  const factory ValueFailure.weakPin({required T failedValue}) = _WeakPin<T>;

  // BVN/NIN failures
  const factory ValueFailure.invalidBvn({required T failedValue}) =
      _InvalidBvn<T>;
  const factory ValueFailure.invalidNin({required T failedValue}) =
      _InvalidNin<T>;

  // Account number failures
  const factory ValueFailure.invalidAccountNumber({required T failedValue}) =
      _InvalidAccountNumber<T>;

  // Generic failures
  const factory ValueFailure.multipleFailures({
    required List<ValueFailure<T>> failures,
  }) = _MultipleFailures<T>;
}

extension ValueFailureX<T> on ValueFailure<T> {
  /// Returns a human-readable error message
  String get message => when(
        empty: (_) => 'This field is required',
        exceedingLength: (_, max) => 'Must be at most $max characters',
        tooShort: (_, min) => 'Must be at least $min characters',
        invalidFormat: (_) => 'Invalid format',
        invalidEmail: (_) => 'Invalid email address',
        invalidPhoneNumber: (_) => 'Invalid phone number',
        invalidMoney: (_) => 'Invalid amount',
        negativeAmount: (_) => 'Amount cannot be negative',
        exceedsLimit: (_, limit) => 'Amount exceeds limit of $limit',
        invalidPin: (_) => 'Invalid PIN',
        weakPin: (_) => 'PIN is too weak (avoid sequences like 1234)',
        invalidBvn: (_) => 'Invalid BVN (must be 11 digits)',
        invalidNin: (_) => 'Invalid NIN (must be 11 digits)',
        invalidAccountNumber: (_) => 'Invalid account number (must be 10 digits)',
        multipleFailures: (failures) => failures.map((f) => f.message).join(', '),
      );
}
