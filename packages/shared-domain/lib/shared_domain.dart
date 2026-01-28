/// Shared domain models for HustleX Pro platform
///
/// This package contains domain models, value objects, and failures
/// shared across consumer and provider Flutter applications.
library shared_domain;

// Entities
export 'src/entities/entity.dart';

// Failures
export 'src/failures/failure.dart';
export 'src/failures/value_failures.dart';

// Use cases
export 'src/usecases/usecase.dart';

// Value objects
export 'src/value_objects/value_object.dart';
export 'src/value_objects/money.dart';
export 'src/value_objects/phone_number.dart';
export 'src/value_objects/email.dart';
export 'src/value_objects/pin.dart';
export 'src/value_objects/bvn.dart';
