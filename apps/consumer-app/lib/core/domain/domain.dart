/// Core domain layer barrel file
///
/// This module contains the domain layer building blocks:
/// - Entities: Objects with identity
/// - Value Objects: Immutable, validated objects
/// - Failures: Domain-specific error types
/// - Use Cases: Business operations
library;

// Entities
export 'entities/entity.dart';

// Failures
export 'failures/failure.dart';
export 'failures/value_failures.dart';

// Use cases
export 'usecases/usecase.dart';

// Value objects
export 'value_objects/value_object.dart';
export 'value_objects/money.dart';
export 'value_objects/phone_number.dart';
export 'value_objects/email.dart';
export 'value_objects/pin.dart';
export 'value_objects/bvn.dart';
