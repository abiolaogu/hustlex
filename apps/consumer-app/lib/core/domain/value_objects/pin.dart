import 'package:dartz/dartz.dart';
import '../failures/value_failures.dart';
import 'value_object.dart';

/// PIN (Personal Identification Number) value object.
/// Validates PIN format and checks for weak patterns.
class Pin extends ValueObject<String> {
  @override
  final Either<ValueFailure<String>, String> value;

  factory Pin(String input) {
    return Pin._(_validatePin(input));
  }

  const Pin._(this.value);

  static Either<ValueFailure<String>, String> _validatePin(String input) {
    if (input.isEmpty) {
      return left(ValueFailure.empty(failedValue: input));
    }

    // Must be exactly 4 digits
    if (input.length != 4) {
      return left(ValueFailure.invalidPin(failedValue: input));
    }

    // Must contain only digits
    if (!RegExp(r'^\d{4}$').hasMatch(input)) {
      return left(ValueFailure.invalidPin(failedValue: input));
    }

    // Check for weak patterns
    if (_isWeakPin(input)) {
      return left(ValueFailure.weakPin(failedValue: input));
    }

    return right(input);
  }

  static bool _isWeakPin(String pin) {
    // Check for sequential patterns
    const sequentialPatterns = [
      '0123',
      '1234',
      '2345',
      '3456',
      '4567',
      '5678',
      '6789',
      '9876',
      '8765',
      '7654',
      '6543',
      '5432',
      '4321',
      '3210',
    ];

    if (sequentialPatterns.contains(pin)) {
      return true;
    }

    // Check for repeated digits (e.g., 1111, 2222)
    if (pin.split('').toSet().length == 1) {
      return true;
    }

    // Check for common weak pins
    const commonWeakPins = [
      '0000',
      '1111',
      '2222',
      '3333',
      '4444',
      '5555',
      '6666',
      '7777',
      '8888',
      '9999',
      '1212',
      '2121',
      '1122',
      '2211',
      '1100',
      '0011',
      '2000',
      '2020',
      '2023',
      '2024',
      '2025',
      '1990',
      '1991',
      '1992',
      '1993',
      '1994',
      '1995',
      '1996',
      '1997',
      '1998',
      '1999',
    ];

    return commonWeakPins.contains(pin);
  }

  /// Returns masked PIN for display (e.g., ****)
  String get masked => '****';
}
