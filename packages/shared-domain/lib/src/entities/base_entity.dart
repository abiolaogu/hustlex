import 'package:flutter/foundation.dart';

/// Base class for all domain entities.
///
/// Entities have identity (id) and are compared by identity,
/// not by their attributes.
@immutable
abstract class Entity {
  /// Unique identifier for this entity
  String get id;

  const Entity();

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Entity && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}

/// Base class for aggregate roots.
///
/// Aggregate roots are entities that serve as the entry point
/// to a cluster of related entities and value objects.
@immutable
abstract class AggregateRoot extends Entity {
  const AggregateRoot();
}
