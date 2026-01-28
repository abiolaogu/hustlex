import 'package:dartz/dartz.dart';
import '../failures/value_failures.dart';
import 'value_object.dart';

/// Email address value object with validation.
class Email extends ValueObject<String> {
  @override
  final Either<ValueFailure<String>, String> value;

  factory Email(String input) {
    return Email._(_validateEmail(input));
  }

  const Email._(this.value);

  static Either<ValueFailure<String>, String> _validateEmail(String input) {
    if (input.isEmpty) {
      return left(ValueFailure.empty(failedValue: input));
    }

    // RFC 5322 simplified email regex
    final emailRegex = RegExp(
      r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
      caseSensitive: false,
    );

    if (!emailRegex.hasMatch(input)) {
      return left(ValueFailure.invalidEmail(failedValue: input));
    }

    // Normalize to lowercase
    return right(input.toLowerCase().trim());
  }

  /// Returns the domain part of the email
  String get domain {
    final email = getOrCrash();
    return email.split('@').last;
  }

  /// Returns the local part (before @)
  String get localPart {
    final email = getOrCrash();
    return email.split('@').first;
  }

  /// Returns masked email for display (e.g., a***e@example.com)
  String get masked {
    final email = getOrCrash();
    final parts = email.split('@');
    final local = parts[0];
    final domain = parts[1];

    if (local.length <= 2) {
      return '${local[0]}***@$domain';
    }
    return '${local[0]}***${local[local.length - 1]}@$domain';
  }
}
