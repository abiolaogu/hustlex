import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'savings_circle.freezed.dart';

/// Type of savings circle (Nigerian savings traditions)
enum CircleType {
  ajo, // Rotating savings - members take turns receiving pooled contributions
  esusu, // Target savings - fixed duration, withdraw at end
  goal, // Personal goal savings with friends for accountability
}

/// Circle status
enum CircleStatus {
  pending, // Waiting for minimum members
  active, // Circle is running
  completed, // All cycles complete
  cancelled, // Circle was cancelled
}

/// Contribution frequency
enum ContributionFrequency {
  daily,
  weekly,
  biweekly,
  monthly,
}

/// Extension for CircleType
extension CircleTypeX on CircleType {
  String get displayName {
    switch (this) {
      case CircleType.ajo:
        return 'Ajo (Rotating)';
      case CircleType.esusu:
        return 'Esusu (Target)';
      case CircleType.goal:
        return 'Goal Savings';
    }
  }

  String get description {
    switch (this) {
      case CircleType.ajo:
        return 'Members take turns receiving the pooled contributions';
      case CircleType.esusu:
        return 'Save towards a target with everyone withdrawing at the end';
      case CircleType.goal:
        return 'Personal savings with friends for accountability';
    }
  }
}

/// Extension for CircleStatus
extension CircleStatusX on CircleStatus {
  String get displayName {
    switch (this) {
      case CircleStatus.pending:
        return 'Pending';
      case CircleStatus.active:
        return 'Active';
      case CircleStatus.completed:
        return 'Completed';
      case CircleStatus.cancelled:
        return 'Cancelled';
    }
  }

  bool get isTerminal =>
      this == CircleStatus.completed || this == CircleStatus.cancelled;
}

/// Extension for ContributionFrequency
extension ContributionFrequencyX on ContributionFrequency {
  String get displayName {
    switch (this) {
      case ContributionFrequency.daily:
        return 'Daily';
      case ContributionFrequency.weekly:
        return 'Weekly';
      case ContributionFrequency.biweekly:
        return 'Bi-weekly';
      case ContributionFrequency.monthly:
        return 'Monthly';
    }
  }

  int get daysInterval {
    switch (this) {
      case ContributionFrequency.daily:
        return 1;
      case ContributionFrequency.weekly:
        return 7;
      case ContributionFrequency.biweekly:
        return 14;
      case ContributionFrequency.monthly:
        return 30;
    }
  }
}

/// Circle member info
@freezed
class CircleMemberInfo with _$CircleMemberInfo {
  const factory CircleMemberInfo({
    required String id,
    required String firstName,
    required String lastName,
    String? avatar,
    String? phone,
  }) = _CircleMemberInfo;

  const CircleMemberInfo._();

  String get fullName => '$firstName $lastName';
  String get initials =>
      '${firstName.isNotEmpty ? firstName[0] : ''}${lastName.isNotEmpty ? lastName[0] : ''}'
          .toUpperCase();
}

/// Savings circle entity
@freezed
class SavingsCircle with _$SavingsCircle implements Entity {
  const factory SavingsCircle({
    required String id,
    required String creatorId,
    required String name,
    String? description,
    String? imageUrl,
    required CircleType type,
    required CircleStatus status,
    required Money contributionAmount,
    required ContributionFrequency frequency,
    required int maxMembers,
    @Default(0) int currentMembers,
    required Money targetAmount,
    required Money totalSaved,
    @Default(0) int currentCycle,
    @Default(0) int totalCycles,
    @Default(false) bool isPrivate,
    String? inviteCode,
    String? creatorName,
    String? nextPayoutMemberId,
    String? nextPayoutMemberName,
    DateTime? startDate,
    DateTime? endDate,
    DateTime? nextContributionDate,
    DateTime? nextPayoutDate,
    @Default([]) List<CircleMember> members,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _SavingsCircle;

  const SavingsCircle._();

  /// Progress towards target (0.0 - 1.0)
  double get progress {
    if (targetAmount.amount <= 0) return 0;
    return (totalSaved.amount / targetAmount.amount).clamp(0.0, 1.0);
  }

  /// Progress percentage (0 - 100)
  int get progressPercent => (progress * 100).round();

  /// Available spots in the circle
  int get spotsLeft => maxMembers - currentMembers;

  /// Check if circle is full
  bool get isFull => currentMembers >= maxMembers;

  /// Check if circle is active
  bool get isActive => status == CircleStatus.active;

  /// Check if circle can be joined
  bool get canJoin => !isFull && status == CircleStatus.pending;

  /// Get pot amount for current cycle
  Money get potAmount => contributionAmount * currentMembers.toDouble();

  /// Formatted contribution string
  String get formattedContribution =>
      '${contributionAmount.formatted}/${frequency.displayName.toLowerCase()}';
}

/// Member role in circle
enum MemberRole {
  admin,
  member,
}

/// Member status in circle
enum MemberStatus {
  pending,
  active,
  defaulted,
  removed,
  left,
}

/// Extension for MemberStatus
extension MemberStatusX on MemberStatus {
  String get displayName {
    switch (this) {
      case MemberStatus.pending:
        return 'Pending';
      case MemberStatus.active:
        return 'Active';
      case MemberStatus.defaulted:
        return 'Defaulted';
      case MemberStatus.removed:
        return 'Removed';
      case MemberStatus.left:
        return 'Left';
    }
  }
}

/// Circle member entity
@freezed
class CircleMember with _$CircleMember implements Entity {
  const factory CircleMember({
    required String id,
    required String circleId,
    required String userId,
    required String userName,
    String? userAvatar,
    required MemberRole role,
    required MemberStatus status,
    @Default(0) int payoutOrder,
    @Default(false) bool hasReceivedPayout,
    required Money totalContributed,
    required Money totalReceived,
    @Default(0) int contributionsMade,
    @Default(0) int contributionsMissed,
    DateTime? joinedAt,
    DateTime? payoutDate,
    CircleMemberInfo? user,
  }) = _CircleMember;

  const CircleMember._();

  /// Check if member is admin
  bool get isAdmin => role == MemberRole.admin;

  /// Check if member is active
  bool get isActive => status == MemberStatus.active;

  /// Contribution rate percentage
  double get contributionRate {
    final total = contributionsMade + contributionsMissed;
    if (total == 0) return 100;
    return (contributionsMade / total) * 100;
  }

  /// Check if member has good standing (>80% contribution rate)
  bool get inGoodStanding => contributionRate >= 80;
}
