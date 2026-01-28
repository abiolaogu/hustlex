import 'package:dartz/dartz.dart';
import 'package:flutter/foundation.dart';
import '../failures/value_failures.dart';

/// Base class for all value objects in the domain layer.
/// Value objects are immutable and validated on construction.
@immutable
abstract class ValueObject<T> {
  const ValueObject();

  /// The validated value wrapped in Either
  Either<ValueFailure<T>, T> get value;

  /// Throws [UnexpectedValueError] containing the [ValueFailure]
  T getOrCrash() {
    return value.fold(
      (f) => throw UnexpectedValueError(f),
      (r) => r,
    );
  }

  /// Returns the value or a default if invalid
  T getOrElse(T dflt) {
    return value.getOrElse(() => dflt);
  }

  /// Returns true if the value is valid
  bool isValid() => value.isRight();

  /// Returns the failure if invalid, null otherwise
  ValueFailure<T>? get failureOrNull {
    return value.fold((f) => f, (_) => null);
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is ValueObject<T> && other.value == value;
  }

  @override
  int get hashCode => value.hashCode;

  @override
  String toString() => 'ValueObject($value)';
}

/// Error thrown when trying to access an invalid value
class UnexpectedValueError extends Error {
  final ValueFailure valueFailure;

  UnexpectedValueError(this.valueFailure);

  @override
  String toString() {
    return 'UnexpectedValueError: $valueFailure';
  }
}
