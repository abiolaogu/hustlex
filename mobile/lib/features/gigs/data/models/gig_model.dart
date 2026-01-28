import 'package:freezed_annotation/freezed_annotation.dart';

part 'gig_model.freezed.dart';
part 'gig_model.g.dart';

enum GigStatus {
  @JsonValue('draft')
  draft,
  @JsonValue('active')
  active,
  @JsonValue('in_progress')
  inProgress,
  @JsonValue('completed')
  completed,
  @JsonValue('cancelled')
  cancelled,
  @JsonValue('disputed')
  disputed,
}

enum GigCategory {
  @JsonValue('tech')
  tech,
  @JsonValue('design')
  design,
  @JsonValue('writing')
  writing,
  @JsonValue('marketing')
  marketing,
  @JsonValue('video')
  video,
  @JsonValue('audio')
  audio,
  @JsonValue('translation')
  translation,
  @JsonValue('data')
  data,
  @JsonValue('admin')
  admin,
  @JsonValue('other')
  other,
}

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

enum MilestoneStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('in_progress')
  inProgress,
  @JsonValue('submitted')
  submitted,
  @JsonValue('approved')
  approved,
  @JsonValue('revision_requested')
  revisionRequested,
  @JsonValue('paid')
  paid,
}

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
}

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

  String get icon {
    switch (this) {
      case GigCategory.tech:
        return '💻';
      case GigCategory.design:
        return '🎨';
      case GigCategory.writing:
        return '✍️';
      case GigCategory.marketing:
        return '📣';
      case GigCategory.video:
        return '🎬';
      case GigCategory.audio:
        return '🎵';
      case GigCategory.translation:
        return '🌐';
      case GigCategory.data:
        return '📊';
      case GigCategory.admin:
        return '📋';
      case GigCategory.other:
        return '🔧';
    }
  }
}

@freezed
class Gig with _$Gig {
  const factory Gig({
    required String id,
    required String title,
    required String description,
    required GigCategory category,
    required GigStatus status,
    required double budgetMin,
    required double budgetMax,
    required String currency,
    required String clientId,
    String? clientName,
    String? clientAvatar,
    String? freelancerId,
    String? freelancerName,
    @Default(false) bool isRemote,
    String? location,
    DateTime? deadline,
    @Default([]) List<String> skills,
    @Default([]) List<String> attachments,
    @Default([]) List<Milestone> milestones,
    @Default(0) int proposalsCount,
    @Default(0) int viewsCount,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Gig;

  const Gig._();

  factory Gig.fromJson(Map<String, dynamic> json) => _$GigFromJson(json);

  String get budgetRange {
    if (budgetMin == budgetMax) {
      return '₦${_formatAmount(budgetMin)}';
    }
    return '₦${_formatAmount(budgetMin)} - ₦${_formatAmount(budgetMax)}';
  }

  String _formatAmount(double amount) {
    if (amount >= 1000000) {
      return '${(amount / 1000000).toStringAsFixed(1)}M';
    } else if (amount >= 1000) {
      return '${(amount / 1000).toStringAsFixed(0)}K';
    }
    return amount.toStringAsFixed(0);
  }

  bool get isActive => status == GigStatus.active;
  bool get isCompleted => status == GigStatus.completed;
  bool get canAcceptProposals => status == GigStatus.active;
}

@freezed
class Milestone with _$Milestone {
  const factory Milestone({
    required String id,
    required String gigId,
    required String title,
    String? description,
    required double amount,
    required int order,
    required MilestoneStatus status,
    DateTime? dueDate,
    DateTime? submittedAt,
    DateTime? approvedAt,
    DateTime? paidAt,
    @Default([]) List<String> deliverables,
  }) = _Milestone;

  const Milestone._();

  factory Milestone.fromJson(Map<String, dynamic> json) => _$MilestoneFromJson(json);

  bool get isPending => status == MilestoneStatus.pending;
  bool get isInProgress => status == MilestoneStatus.inProgress;
  bool get isCompleted => status == MilestoneStatus.approved || status == MilestoneStatus.paid;
}

@freezed
class Proposal with _$Proposal {
  const factory Proposal({
    required String id,
    required String gigId,
    required String freelancerId,
    String? freelancerName,
    String? freelancerAvatar,
    @Default(0) int freelancerRating,
    @Default(0) int freelancerCompletedGigs,
    required String coverLetter,
    required double bidAmount,
    required int deliveryDays,
    required ProposalStatus status,
    @Default([]) List<String> attachments,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Proposal;

  const Proposal._();

  factory Proposal.fromJson(Map<String, dynamic> json) => _$ProposalFromJson(json);

  bool get isPending => status == ProposalStatus.pending;
  bool get isAccepted => status == ProposalStatus.accepted;
}

@freezed
class GigReview with _$GigReview {
  const factory GigReview({
    required String id,
    required String gigId,
    required String reviewerId,
    required String revieweeId,
    String? reviewerName,
    String? reviewerAvatar,
    required int rating,
    String? comment,
    DateTime? createdAt,
  }) = _GigReview;

  factory GigReview.fromJson(Map<String, dynamic> json) => _$GigReviewFromJson(json);
}

@freezed
class CreateGigRequest with _$CreateGigRequest {
  const factory CreateGigRequest({
    required String title,
    required String description,
    required GigCategory category,
    required double budgetMin,
    required double budgetMax,
    @Default('NGN') String currency,
    @Default(false) bool isRemote,
    String? location,
    DateTime? deadline,
    @Default([]) List<String> skills,
    @Default([]) List<CreateMilestoneRequest> milestones,
  }) = _CreateGigRequest;

  factory CreateGigRequest.fromJson(Map<String, dynamic> json) => _$CreateGigRequestFromJson(json);
}

@freezed
class CreateMilestoneRequest with _$CreateMilestoneRequest {
  const factory CreateMilestoneRequest({
    required String title,
    String? description,
    required double amount,
    required int order,
    DateTime? dueDate,
  }) = _CreateMilestoneRequest;

  factory CreateMilestoneRequest.fromJson(Map<String, dynamic> json) => _$CreateMilestoneRequestFromJson(json);
}

@freezed
class CreateProposalRequest with _$CreateProposalRequest {
  const factory CreateProposalRequest({
    required String gigId,
    required String coverLetter,
    required double bidAmount,
    required int deliveryDays,
    @Default([]) List<String> attachments,
  }) = _CreateProposalRequest;

  factory CreateProposalRequest.fromJson(Map<String, dynamic> json) => _$CreateProposalRequestFromJson(json);
}

@freezed
class GigFilter with _$GigFilter {
  const factory GigFilter({
    GigCategory? category,
    GigStatus? status,
    double? minBudget,
    double? maxBudget,
    bool? isRemote,
    String? search,
    @Default(1) int page,
    @Default(20) int limit,
    @Default('created_at') String sortBy,
    @Default('desc') String sortOrder,
  }) = _GigFilter;

  factory GigFilter.fromJson(Map<String, dynamic> json) => _$GigFilterFromJson(json);
}

@freezed
class PaginatedGigs with _$PaginatedGigs {
  const factory PaginatedGigs({
    required List<Gig> gigs,
    required int total,
    required int page,
    required int limit,
    required bool hasMore,
  }) = _PaginatedGigs;

  factory PaginatedGigs.fromJson(Map<String, dynamic> json) => _$PaginatedGigsFromJson(json);
}
