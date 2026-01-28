/// Core domain layer barrel file
///
/// This module contains the domain layer building blocks:
/// - Entities: Objects with identity
/// - Value Objects: Immutable, validated objects
/// - Failures: Domain-specific error types
/// - Use Cases: Business operations
library;

export 'entities/base_entity.dart';
export 'failures/failures.dart';
export 'usecases/usecase.dart';
export 'value_objects/value_objects.dart';
