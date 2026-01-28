import 'package:dartz/dartz.dart';
import '../failures/value_failures.dart';
import 'value_object.dart';

/// Bank Verification Number (BVN) value object.
/// Nigerian BVN is exactly 11 digits.
class Bvn extends ValueObject<String> {
  @override
  final Either<ValueFailure<String>, String> value;

  factory Bvn(String input) {
    return Bvn._(_validateBvn(input));
  }

  const Bvn._(this.value);

  static Either<ValueFailure<String>, String> _validateBvn(String input) {
    if (input.isEmpty) {
      return left(ValueFailure.empty(failedValue: input));
    }

    // Remove any spaces or dashes
    final cleaned = input.replaceAll(RegExp(r'[\s-]'), '');

    // Must be exactly 11 digits
    if (!RegExp(r'^\d{11}$').hasMatch(cleaned)) {
      return left(ValueFailure.invalidBvn(failedValue: input));
    }

    return right(cleaned);
  }

  /// Returns masked BVN for display (e.g., 123****8901)
  String get masked {
    final bvn = getOrCrash();
    return '${bvn.substring(0, 3)}*****${bvn.substring(8)}';
  }
}

/// National Identification Number (NIN) value object.
/// Nigerian NIN is exactly 11 digits.
class Nin extends ValueObject<String> {
  @override
  final Either<ValueFailure<String>, String> value;

  factory Nin(String input) {
    return Nin._(_validateNin(input));
  }

  const Nin._(this.value);

  static Either<ValueFailure<String>, String> _validateNin(String input) {
    if (input.isEmpty) {
      return left(ValueFailure.empty(failedValue: input));
    }

    // Remove any spaces or dashes
    final cleaned = input.replaceAll(RegExp(r'[\s-]'), '');

    // Must be exactly 11 digits
    if (!RegExp(r'^\d{11}$').hasMatch(cleaned)) {
      return left(ValueFailure.invalidNin(failedValue: input));
    }

    return right(cleaned);
  }

  /// Returns masked NIN for display (e.g., 123****8901)
  String get masked {
    final nin = getOrCrash();
    return '${nin.substring(0, 3)}*****${nin.substring(8)}';
  }
}

/// Nigerian Bank Account Number value object.
/// NUBAN account numbers are exactly 10 digits.
class AccountNumber extends ValueObject<String> {
  @override
  final Either<ValueFailure<String>, String> value;

  factory AccountNumber(String input) {
    return AccountNumber._(_validateAccountNumber(input));
  }

  const AccountNumber._(this.value);

  static Either<ValueFailure<String>, String> _validateAccountNumber(
      String input) {
    if (input.isEmpty) {
      return left(ValueFailure.empty(failedValue: input));
    }

    // Remove any spaces or dashes
    final cleaned = input.replaceAll(RegExp(r'[\s-]'), '');

    // Must be exactly 10 digits
    if (!RegExp(r'^\d{10}$').hasMatch(cleaned)) {
      return left(ValueFailure.invalidAccountNumber(failedValue: input));
    }

    return right(cleaned);
  }

  /// Returns masked account number for display (e.g., ******5678)
  String get masked {
    final account = getOrCrash();
    return '******${account.substring(6)}';
  }
}
