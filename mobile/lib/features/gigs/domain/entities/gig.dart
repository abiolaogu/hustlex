import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'gig.freezed.dart';

/// Gig status enum
enum GigStatus {
  draft,
  active,
  inProgress,
  completed,
  cancelled,
  disputed,
}

/// Gig category enum
enum GigCategory {
  tech,
  design,
  writing,
  marketing,
  video,
  audio,
  translation,
  data,
  admin,
  other,
}

/// Extension for GigStatus display names
extension GigStatusX on GigStatus {
  String get displayName {
    switch (this) {
      case GigStatus.draft:
        return 'Draft';
      case GigStatus.active:
        return 'Active';
      case GigStatus.inProgress:
        return 'In Progress';
      case GigStatus.completed:
        return 'Completed';
      case GigStatus.cancelled:
        return 'Cancelled';
      case GigStatus.disputed:
        return 'Disputed';
    }
  }

  bool get canEdit => this == GigStatus.draft;
  bool get canAcceptProposals => this == GigStatus.active;
  bool get isTerminal =>
      this == GigStatus.completed ||
      this == GigStatus.cancelled;
}

/// Extension for GigCategory display and icons
extension GigCategoryX on GigCategory {
  String get displayName {
    switch (this) {
      case GigCategory.tech:
        return 'Technology';
      case GigCategory.design:
        return 'Design';
      case GigCategory.writing:
        return 'Writing';
      case GigCategory.marketing:
        return 'Marketing';
      case GigCategory.video:
        return 'Video';
      case GigCategory.audio:
        return 'Audio';
      case GigCategory.translation:
        return 'Translation';
      case GigCategory.data:
        return 'Data Entry';
      case GigCategory.admin:
        return 'Admin';
      case GigCategory.other:
        return 'Other';
    }
  }
}

/// Gig participant info
@freezed
class GigParticipant with _$GigParticipant {
  const factory GigParticipant({
    required String id,
    required String firstName,
    required String lastName,
    String? avatar,
    @Default(0.0) double rating,
    @Default(0) int completedGigs,
  }) = _GigParticipant;

  const GigParticipant._();

  String get fullName => '$firstName $lastName';
}

/// Gig entity representing a job posting
@freezed
class Gig with _$Gig implements Entity {
  const factory Gig({
    required String id,
    required String clientId,
    required String title,
    required String description,
    required GigCategory category,
    required GigStatus status,
    required Money budgetMin,
    required Money budgetMax,
    required int durationDays,
    @Default(true) bool isRemote,
    String? location,
    DateTime? deadline,
    @Default([]) List<String> skills,
    @Default([]) List<String> attachments,
    @Default(0) int proposalCount,
    @Default(0) int viewsCount,
    String? freelancerId,
    GigParticipant? client,
    GigParticipant? freelancer,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _Gig;

  const Gig._();

  /// Get formatted budget range
  String get budgetRange {
    if (budgetMin.amount == budgetMax.amount) {
      return budgetMin.formatted;
    }
    return '${budgetMin.formatted} - ${budgetMax.formatted}';
  }

  /// Check if gig is active
  bool get isActive => status == GigStatus.active;

  /// Check if gig is completed
  bool get isCompleted => status == GigStatus.completed;

  /// Check if gig can accept proposals
  bool get canAcceptProposals => status == GigStatus.active;

  /// Check if gig has been assigned
  bool get isAssigned => freelancerId != null;
}
