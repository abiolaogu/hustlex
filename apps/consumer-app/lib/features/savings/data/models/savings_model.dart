import 'package:freezed_annotation/freezed_annotation.dart';

part 'savings_model.freezed.dart';
part 'savings_model.g.dart';

enum CircleType {
  @JsonValue('ajo')
  ajo, // Rotating savings - each member takes turns receiving the pot
  @JsonValue('esusu')
  esusu, // Emergency/target savings - fixed duration, withdraw at end
  @JsonValue('goal')
  goal, // Personal goal savings with friends for accountability
}

enum CircleStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('active')
  active,
  @JsonValue('completed')
  completed,
  @JsonValue('cancelled')
  cancelled,
}

enum MemberRole {
  @JsonValue('admin')
  admin,
  @JsonValue('member')
  member,
}

enum MemberStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('active')
  active,
  @JsonValue('removed')
  removed,
  @JsonValue('left')
  left,
}

enum ContributionFrequency {
  @JsonValue('daily')
  daily,
  @JsonValue('weekly')
  weekly,
  @JsonValue('biweekly')
  biweekly,
  @JsonValue('monthly')
  monthly,
}

enum ContributionStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('paid')
  paid,
  @JsonValue('overdue')
  overdue,
  @JsonValue('missed')
  missed,
}

enum PayoutStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('processing')
  processing,
  @JsonValue('completed')
  completed,
  @JsonValue('failed')
  failed,
}

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

  String get icon {
    switch (this) {
      case CircleType.ajo:
        return '🔄';
      case CircleType.esusu:
        return '🎯';
      case CircleType.goal:
        return '🏆';
    }
  }
}

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

@freezed
class SavingsCircle with _$SavingsCircle {
  const factory SavingsCircle({
    required String id,
    required String name,
    String? description,
    String? imageUrl,
    required CircleType type,
    required CircleStatus status,
    required double contributionAmount,
    required ContributionFrequency frequency,
    required String currency,
    required int maxMembers,
    required int currentMembers,
    required String creatorId,
    String? creatorName,
    @Default(false) bool isPrivate,
    String? inviteCode,
    DateTime? startDate,
    DateTime? endDate,
    @Default(0) double totalSaved,
    @Default(0) double targetAmount,
    @Default(0) int currentCycle,
    @Default(0) int totalCycles,
    String? nextPayoutMemberId,
    String? nextPayoutMemberName,
    DateTime? nextContributionDate,
    DateTime? nextPayoutDate,
    @Default([]) List<CircleMember> members,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _SavingsCircle;

  const SavingsCircle._();

  factory SavingsCircle.fromJson(Map<String, dynamic> json) => _$SavingsCircleFromJson(json);

  double get progress {
    if (targetAmount <= 0) return 0;
    return (totalSaved / targetAmount).clamp(0.0, 1.0);
  }

  int get spotsLeft => maxMembers - currentMembers;
  bool get isFull => currentMembers >= maxMembers;
  bool get isActive => status == CircleStatus.active;
  bool get canJoin => !isFull && status == CircleStatus.pending;

  String get formattedContribution {
    return '₦${contributionAmount.toStringAsFixed(0)}/${frequency.displayName.toLowerCase()}';
  }

  double get potAmount => contributionAmount * currentMembers;
}

@freezed
class CircleMember with _$CircleMember {
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
    @Default(0) double totalContributed,
    @Default(0) int contributionsMade,
    @Default(0) int contributionsMissed,
    DateTime? joinedAt,
    DateTime? payoutDate,
  }) = _CircleMember;

  const CircleMember._();

  factory CircleMember.fromJson(Map<String, dynamic> json) => _$CircleMemberFromJson(json);

  bool get isAdmin => role == MemberRole.admin;
  bool get isActive => status == MemberStatus.active;

  double get contributionRate {
    final total = contributionsMade + contributionsMissed;
    if (total == 0) return 100;
    return (contributionsMade / total) * 100;
  }
}

@freezed
class Contribution with _$Contribution {
  const factory Contribution({
    required String id,
    required String circleId,
    required String memberId,
    String? memberName,
    required double amount,
    required int cycleNumber,
    required ContributionStatus status,
    DateTime? dueDate,
    DateTime? paidAt,
    String? transactionId,
    String? paymentMethod,
  }) = _Contribution;

  const Contribution._();

  factory Contribution.fromJson(Map<String, dynamic> json) => _$ContributionFromJson(json);

  bool get isPaid => status == ContributionStatus.paid;
  bool get isOverdue => status == ContributionStatus.overdue;
}

@freezed
class Payout with _$Payout {
  const factory Payout({
    required String id,
    required String circleId,
    required String memberId,
    String? memberName,
    required double amount,
    required int cycleNumber,
    required PayoutStatus status,
    DateTime? scheduledDate,
    DateTime? processedAt,
    String? transactionId,
    String? failureReason,
  }) = _Payout;

  const Payout._();

  factory Payout.fromJson(Map<String, dynamic> json) => _$PayoutFromJson(json);

  bool get isCompleted => status == PayoutStatus.completed;
  bool get isPending => status == PayoutStatus.pending;
}

@freezed
class CircleInvite with _$CircleInvite {
  const factory CircleInvite({
    required String id,
    required String circleId,
    required String circleName,
    required String inviterId,
    required String inviterName,
    required String inviteePhone,
    @Default(false) bool isAccepted,
    DateTime? expiresAt,
    DateTime? createdAt,
  }) = _CircleInvite;

  factory CircleInvite.fromJson(Map<String, dynamic> json) => _$CircleInviteFromJson(json);
}

@freezed
class CreateCircleRequest with _$CreateCircleRequest {
  const factory CreateCircleRequest({
    required String name,
    String? description,
    required CircleType type,
    required double contributionAmount,
    required ContributionFrequency frequency,
    @Default('NGN') String currency,
    required int maxMembers,
    @Default(false) bool isPrivate,
    DateTime? startDate,
    double? targetAmount, // For esusu/goal type
    int? durationMonths, // For esusu/goal type
  }) = _CreateCircleRequest;

  factory CreateCircleRequest.fromJson(Map<String, dynamic> json) => _$CreateCircleRequestFromJson(json);
}

@freezed
class JoinCircleRequest with _$JoinCircleRequest {
  const factory JoinCircleRequest({
    required String circleId,
    String? inviteCode,
  }) = _JoinCircleRequest;

  factory JoinCircleRequest.fromJson(Map<String, dynamic> json) => _$JoinCircleRequestFromJson(json);
}

@freezed
class MakeContributionRequest with _$MakeContributionRequest {
  const factory MakeContributionRequest({
    required String circleId,
    required String contributionId,
    required String paymentMethod,
    String? paymentReference,
  }) = _MakeContributionRequest;

  factory MakeContributionRequest.fromJson(Map<String, dynamic> json) => _$MakeContributionRequestFromJson(json);
}

@freezed
class SavingsStats with _$SavingsStats {
  const factory SavingsStats({
    @Default(0) double totalSaved,
    @Default(0) int activeCircles,
    @Default(0) int completedCircles,
    @Default(0) double totalEarned,
    @Default(0) double totalPayoutsReceived,
    @Default(100.0) double contributionRate,
  }) = _SavingsStats;

  factory SavingsStats.fromJson(Map<String, dynamic> json) => _$SavingsStatsFromJson(json);
}

@freezed
class CircleFilter with _$CircleFilter {
  const factory CircleFilter({
    CircleType? type,
    CircleStatus? status,
    ContributionFrequency? frequency,
    double? minContribution,
    double? maxContribution,
    @Default(false) bool onlyJoinable,
    String? search,
    @Default(1) int page,
    @Default(20) int limit,
  }) = _CircleFilter;

  factory CircleFilter.fromJson(Map<String, dynamic> json) => _$CircleFilterFromJson(json);
}

@freezed
class PaginatedCircles with _$PaginatedCircles {
  const factory PaginatedCircles({
    required List<SavingsCircle> circles,
    required int total,
    required int page,
    required int limit,
    required bool hasMore,
  }) = _PaginatedCircles;

  factory PaginatedCircles.fromJson(Map<String, dynamic> json) => _$PaginatedCirclesFromJson(json);
}
