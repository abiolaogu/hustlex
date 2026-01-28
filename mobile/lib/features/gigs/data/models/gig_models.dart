import 'package:json_annotation/json_annotation.dart';

part 'gig_models.g.dart';

/// Gig status
enum GigStatus {
  @JsonValue('draft')
  draft,
  @JsonValue('open')
  open,
  @JsonValue('in_progress')
  inProgress,
  @JsonValue('completed')
  completed,
  @JsonValue('cancelled')
  cancelled,
  @JsonValue('disputed')
  disputed,
}

/// Proposal status
enum ProposalStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('accepted')
  accepted,
  @JsonValue('rejected')
  rejected,
  @JsonValue('withdrawn')
  withdrawn,
}

/// Gig category
@JsonSerializable()
class GigCategory {
  final String id;
  final String name;
  final String? icon;
  final String? description;

  GigCategory({
    required this.id,
    required this.name,
    this.icon,
    this.description,
  });

  factory GigCategory.fromJson(Map<String, dynamic> json) =>
      _$GigCategoryFromJson(json);
  Map<String, dynamic> toJson() => _$GigCategoryToJson(this);
}

/// Gig model
@JsonSerializable()
class Gig {
  final String id;
  final String clientId;
  final String? freelancerId;
  final String title;
  final String description;
  final String categoryId;
  final GigCategory? category;
  final double budgetMin;
  final double budgetMax;
  final int durationDays;
  final bool isRemote;
  final String? location;
  final List<String> skills;
  final List<String>? attachments;
  final GigStatus status;
  final int proposalCount;
  final DateTime? deadline;
  final DateTime createdAt;
  final DateTime updatedAt;
  final GigUser? client;

  Gig({
    required this.id,
    required this.clientId,
    this.freelancerId,
    required this.title,
    required this.description,
    required this.categoryId,
    this.category,
    required this.budgetMin,
    required this.budgetMax,
    required this.durationDays,
    this.isRemote = true,
    this.location,
    required this.skills,
    this.attachments,
    required this.status,
    this.proposalCount = 0,
    this.deadline,
    required this.createdAt,
    required this.updatedAt,
    this.client,
  });

  String get budgetRange {
    if (budgetMin == budgetMax) {
      return '₦${budgetMin.toStringAsFixed(0)}';
    }
    return '₦${budgetMin.toStringAsFixed(0)} - ₦${budgetMax.toStringAsFixed(0)}';
  }

  factory Gig.fromJson(Map<String, dynamic> json) => _$GigFromJson(json);
  Map<String, dynamic> toJson() => _$GigToJson(this);
}

/// Simple user for gig display
@JsonSerializable()
class GigUser {
  final String id;
  final String firstName;
  final String lastName;
  final String? avatar;
  final double? rating;
  final int completedGigs;

  GigUser({
    required this.id,
    required this.firstName,
    required this.lastName,
    this.avatar,
    this.rating,
    this.completedGigs = 0,
  });

  String get fullName => '$firstName $lastName';

  factory GigUser.fromJson(Map<String, dynamic> json) =>
      _$GigUserFromJson(json);
  Map<String, dynamic> toJson() => _$GigUserToJson(this);
}

/// Proposal model
@JsonSerializable()
class Proposal {
  final String id;
  final String gigId;
  final String freelancerId;
  final double proposedAmount;
  final int deliveryDays;
  final String coverLetter;
  final List<String>? attachments;
  final ProposalStatus status;
  final DateTime createdAt;
  final DateTime updatedAt;
  final Gig? gig;
  final GigUser? freelancer;

  Proposal({
    required this.id,
    required this.gigId,
    required this.freelancerId,
    required this.proposedAmount,
    required this.deliveryDays,
    required this.coverLetter,
    this.attachments,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
    this.gig,
    this.freelancer,
  });

  factory Proposal.fromJson(Map<String, dynamic> json) =>
      _$ProposalFromJson(json);
  Map<String, dynamic> toJson() => _$ProposalToJson(this);
}

/// Contract model
@JsonSerializable()
class Contract {
  final String id;
  final String gigId;
  final String clientId;
  final String freelancerId;
  final double amount;
  final double escrowAmount;
  final double platformFee;
  final int deliveryDays;
  final DateTime startDate;
  final DateTime dueDate;
  final String status;
  final DateTime? completedAt;
  final DateTime createdAt;
  final Gig? gig;
  final GigUser? client;
  final GigUser? freelancer;

  Contract({
    required this.id,
    required this.gigId,
    required this.clientId,
    required this.freelancerId,
    required this.amount,
    required this.escrowAmount,
    required this.platformFee,
    required this.deliveryDays,
    required this.startDate,
    required this.dueDate,
    required this.status,
    this.completedAt,
    required this.createdAt,
    this.gig,
    this.client,
    this.freelancer,
  });

  factory Contract.fromJson(Map<String, dynamic> json) =>
      _$ContractFromJson(json);
  Map<String, dynamic> toJson() => _$ContractToJson(this);
}

/// Create gig request
@JsonSerializable()
class CreateGigRequest {
  final String title;
  final String description;
  final String categoryId;
  final double budgetMin;
  final double budgetMax;
  final int durationDays;
  final bool isRemote;
  final String? location;
  final List<String> skills;
  final DateTime? deadline;

  CreateGigRequest({
    required this.title,
    required this.description,
    required this.categoryId,
    required this.budgetMin,
    required this.budgetMax,
    required this.durationDays,
    this.isRemote = true,
    this.location,
    required this.skills,
    this.deadline,
  });

  factory CreateGigRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateGigRequestFromJson(json);
  Map<String, dynamic> toJson() => _$CreateGigRequestToJson(this);
}

/// Submit proposal request
@JsonSerializable()
class SubmitProposalRequest {
  final String gigId;
  final double proposedAmount;
  final int deliveryDays;
  final String coverLetter;

  SubmitProposalRequest({
    required this.gigId,
    required this.proposedAmount,
    required this.deliveryDays,
    required this.coverLetter,
  });

  factory SubmitProposalRequest.fromJson(Map<String, dynamic> json) =>
      _$SubmitProposalRequestFromJson(json);
  Map<String, dynamic> toJson() => _$SubmitProposalRequestToJson(this);
}

/// Review model
@JsonSerializable()
class Review {
  final String id;
  final String contractId;
  final String reviewerId;
  final String revieweeId;
  final double rating;
  final String? comment;
  final DateTime createdAt;
  final GigUser? reviewer;

  Review({
    required this.id,
    required this.contractId,
    required this.reviewerId,
    required this.revieweeId,
    required this.rating,
    this.comment,
    required this.createdAt,
    this.reviewer,
  });

  factory Review.fromJson(Map<String, dynamic> json) => _$ReviewFromJson(json);
  Map<String, dynamic> toJson() => _$ReviewToJson(this);
}
