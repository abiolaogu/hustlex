import 'package:json_annotation/json_annotation.dart';

part 'savings_models.g.dart';

/// Circle type
enum CircleType {
  @JsonValue('rotational')
  rotational,  // Ajo/Esusu - members take turns receiving
  @JsonValue('fixed_target')
  fixedTarget, // Everyone saves toward a goal
}

/// Circle status
enum CircleStatus {
  @JsonValue('pending')
  pending,     // Waiting for minimum members
  @JsonValue('active')
  active,      // Circle is running
  @JsonValue('completed')
  completed,   // All rounds complete
  @JsonValue('cancelled')
  cancelled,
}

/// Contribution frequency
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

/// Member status
enum MemberStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('active')
  active,
  @JsonValue('defaulted')
  defaulted,
  @JsonValue('left')
  left,
}

/// Savings circle model
@JsonSerializable()
class SavingsCircle {
  final String id;
  final String creatorId;
  final String name;
  final String? description;
  final CircleType type;
  final double contributionAmount;
  final ContributionFrequency frequency;
  final int maxMembers;
  final int currentMembers;
  final double targetAmount;
  final double currentAmount;
  final int currentRound;
  final int totalRounds;
  final CircleStatus status;
  final bool isPrivate;
  final String? inviteCode;
  final DateTime? startDate;
  final DateTime? nextContributionDate;
  final DateTime createdAt;
  final DateTime updatedAt;
  final List<CircleMember>? members;

  SavingsCircle({
    required this.id,
    required this.creatorId,
    required this.name,
    this.description,
    required this.type,
    required this.contributionAmount,
    required this.frequency,
    required this.maxMembers,
    this.currentMembers = 0,
    required this.targetAmount,
    this.currentAmount = 0,
    this.currentRound = 0,
    required this.totalRounds,
    required this.status,
    this.isPrivate = false,
    this.inviteCode,
    this.startDate,
    this.nextContributionDate,
    required this.createdAt,
    required this.updatedAt,
    this.members,
  });

  double get progress => targetAmount > 0 ? currentAmount / targetAmount : 0;
  int get availableSlots => maxMembers - currentMembers;
  bool get isFull => currentMembers >= maxMembers;

  String get frequencyLabel {
    switch (frequency) {
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

  String get typeLabel {
    switch (type) {
      case CircleType.rotational:
        return 'Rotational (Ajo)';
      case CircleType.fixedTarget:
        return 'Fixed Target';
    }
  }

  factory SavingsCircle.fromJson(Map<String, dynamic> json) =>
      _$SavingsCircleFromJson(json);
  Map<String, dynamic> toJson() => _$SavingsCircleToJson(this);
}

/// Circle member model
@JsonSerializable()
class CircleMember {
  final String id;
  final String circleId;
  final String userId;
  final int position;  // For rotational: payout order
  final MemberStatus status;
  final double totalContributed;
  final double totalReceived;
  final bool hasReceivedPayout;  // For rotational circles
  final DateTime joinedAt;
  final CircleMemberUser? user;

  CircleMember({
    required this.id,
    required this.circleId,
    required this.userId,
    required this.position,
    required this.status,
    this.totalContributed = 0,
    this.totalReceived = 0,
    this.hasReceivedPayout = false,
    required this.joinedAt,
    this.user,
  });

  factory CircleMember.fromJson(Map<String, dynamic> json) =>
      _$CircleMemberFromJson(json);
  Map<String, dynamic> toJson() => _$CircleMemberToJson(this);
}

/// Simple user for circle member display
@JsonSerializable()
class CircleMemberUser {
  final String id;
  final String firstName;
  final String lastName;
  final String? avatar;
  final String phone;

  CircleMemberUser({
    required this.id,
    required this.firstName,
    required this.lastName,
    this.avatar,
    required this.phone,
  });

  String get fullName => '$firstName $lastName';
  String get initials => '${firstName[0]}${lastName[0]}'.toUpperCase();

  factory CircleMemberUser.fromJson(Map<String, dynamic> json) =>
      _$CircleMemberUserFromJson(json);
  Map<String, dynamic> toJson() => _$CircleMemberUserToJson(this);
}

/// Contribution model
@JsonSerializable()
class Contribution {
  final String id;
  final String circleId;
  final String memberId;
  final int round;
  final double amount;
  final String status;
  final DateTime dueDate;
  final DateTime? paidAt;
  final String? transactionId;
  final DateTime createdAt;

  Contribution({
    required this.id,
    required this.circleId,
    required this.memberId,
    required this.round,
    required this.amount,
    required this.status,
    required this.dueDate,
    this.paidAt,
    this.transactionId,
    required this.createdAt,
  });

  bool get isPaid => paidAt != null;
  bool get isOverdue => !isPaid && DateTime.now().isAfter(dueDate);

  factory Contribution.fromJson(Map<String, dynamic> json) =>
      _$ContributionFromJson(json);
  Map<String, dynamic> toJson() => _$ContributionToJson(this);
}

/// Payout model (for rotational circles)
@JsonSerializable()
class Payout {
  final String id;
  final String circleId;
  final String recipientId;
  final int round;
  final double amount;
  final String status;
  final DateTime scheduledDate;
  final DateTime? paidAt;
  final String? transactionId;
  final DateTime createdAt;
  final CircleMemberUser? recipient;

  Payout({
    required this.id,
    required this.circleId,
    required this.recipientId,
    required this.round,
    required this.amount,
    required this.status,
    required this.scheduledDate,
    this.paidAt,
    this.transactionId,
    required this.createdAt,
    this.recipient,
  });

  factory Payout.fromJson(Map<String, dynamic> json) => _$PayoutFromJson(json);
  Map<String, dynamic> toJson() => _$PayoutToJson(this);
}

/// Create circle request
@JsonSerializable()
class CreateCircleRequest {
  final String name;
  final String? description;
  final CircleType type;
  final double contributionAmount;
  final ContributionFrequency frequency;
  final int maxMembers;
  final double? targetAmount;
  final bool isPrivate;
  final DateTime? startDate;

  CreateCircleRequest({
    required this.name,
    this.description,
    required this.type,
    required this.contributionAmount,
    required this.frequency,
    required this.maxMembers,
    this.targetAmount,
    this.isPrivate = false,
    this.startDate,
  });

  factory CreateCircleRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateCircleRequestFromJson(json);
  Map<String, dynamic> toJson() => _$CreateCircleRequestToJson(this);
}

/// Join circle request
@JsonSerializable()
class JoinCircleRequest {
  final String circleId;
  final String? inviteCode;

  JoinCircleRequest({
    required this.circleId,
    this.inviteCode,
  });

  factory JoinCircleRequest.fromJson(Map<String, dynamic> json) =>
      _$JoinCircleRequestFromJson(json);
  Map<String, dynamic> toJson() => _$JoinCircleRequestToJson(this);
}

/// Make contribution request
@JsonSerializable()
class MakeContributionRequest {
  final String circleId;
  final String pin;

  MakeContributionRequest({
    required this.circleId,
    required this.pin,
  });

  factory MakeContributionRequest.fromJson(Map<String, dynamic> json) =>
      _$MakeContributionRequestFromJson(json);
  Map<String, dynamic> toJson() => _$MakeContributionRequestToJson(this);
}
