/// Base interface for all domain entities.
///
/// Entities are objects that have a distinct identity that runs through time
/// and different representations. They are defined by their identity rather
/// than their attributes.
abstract class Entity {
  /// Unique identifier for the entity.
  /// All entities must have an ID for tracking and comparison.
  String get id;
}
