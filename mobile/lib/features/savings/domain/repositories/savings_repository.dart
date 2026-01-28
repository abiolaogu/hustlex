import 'package:dartz/dartz.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/savings_circle.dart';
import '../entities/contribution.dart';
import '../entities/payout.dart';
import '../entities/savings_stats.dart';
import '../entities/circle_invite.dart';

/// Filter parameters for circle listings
class CircleFilter {
  final CircleType? type;
  final CircleStatus? status;
  final ContributionFrequency? frequency;
  final Money? minContribution;
  final Money? maxContribution;
  final bool onlyJoinable;
  final String? search;
  final int page;
  final int limit;

  const CircleFilter({
    this.type,
    this.status,
    this.frequency,
    this.minContribution,
    this.maxContribution,
    this.onlyJoinable = false,
    this.search,
    this.page = 1,
    this.limit = 20,
  });

  CircleFilter copyWith({
    CircleType? type,
    CircleStatus? status,
    ContributionFrequency? frequency,
    Money? minContribution,
    Money? maxContribution,
    bool? onlyJoinable,
    String? search,
    int? page,
    int? limit,
  }) {
    return CircleFilter(
      type: type ?? this.type,
      status: status ?? this.status,
      frequency: frequency ?? this.frequency,
      minContribution: minContribution ?? this.minContribution,
      maxContribution: maxContribution ?? this.maxContribution,
      onlyJoinable: onlyJoinable ?? this.onlyJoinable,
      search: search ?? this.search,
      page: page ?? this.page,
      limit: limit ?? this.limit,
    );
  }
}

/// Paginated result for circles
class PaginatedCircles {
  final List<SavingsCircle> circles;
  final int total;
  final int page;
  final int limit;
  final bool hasMore;

  const PaginatedCircles({
    required this.circles,
    required this.total,
    required this.page,
    required this.limit,
    required this.hasMore,
  });
}

/// Abstract repository interface for savings feature
abstract class SavingsRepository {
  // ============ CIRCLE OPERATIONS ============

  /// Get paginated list of public circles with optional filters
  Future<Either<Failure, PaginatedCircles>> getCircles(CircleFilter filter);

  /// Get circles that current user is a member of
  Future<Either<Failure, List<SavingsCircle>>> getMyCircles({
    CircleStatus? status,
  });

  /// Get a single circle by ID
  Future<Either<Failure, SavingsCircle>> getCircle(String circleId);

  /// Watch a circle for real-time updates
  Stream<Either<Failure, SavingsCircle>> watchCircle(String circleId);

  /// Create a new savings circle
  Future<Either<Failure, SavingsCircle>> createCircle({
    required String name,
    String? description,
    required CircleType type,
    required Money contributionAmount,
    required ContributionFrequency frequency,
    required int maxMembers,
    required bool isPrivate,
    DateTime? startDate,
    Money? targetAmount,
    int? durationMonths,
  });

  /// Update circle details (admin only)
  Future<Either<Failure, SavingsCircle>> updateCircle({
    required String circleId,
    String? name,
    String? description,
    bool? isPrivate,
    DateTime? startDate,
  });

  /// Start a pending circle (admin only, requires minimum members)
  Future<Either<Failure, SavingsCircle>> startCircle(String circleId);

  /// Cancel a circle (admin only)
  Future<Either<Failure, Unit>> cancelCircle({
    required String circleId,
    required String reason,
  });

  // ============ MEMBERSHIP OPERATIONS ============

  /// Join a public circle
  Future<Either<Failure, CircleMember>> joinCircle(String circleId);

  /// Join a private circle with invite code
  Future<Either<Failure, CircleMember>> joinCircleWithCode({
    required String circleId,
    required String inviteCode,
  });

  /// Leave a circle
  Future<Either<Failure, Unit>> leaveCircle(String circleId);

  /// Remove a member from circle (admin only)
  Future<Either<Failure, Unit>> removeMember({
    required String circleId,
    required String memberId,
    required String reason,
  });

  /// Get members of a circle
  Future<Either<Failure, List<CircleMember>>> getCircleMembers(String circleId);

  /// Update payout order (admin only, for ajo/rotating circles)
  Future<Either<Failure, Unit>> updatePayoutOrder({
    required String circleId,
    required List<String> memberIds,
  });

  // ============ INVITE OPERATIONS ============

  /// Send invite to join circle
  Future<Either<Failure, CircleInvite>> sendInvite({
    required String circleId,
    required String phoneNumber,
  });

  /// Get pending invites for current user
  Future<Either<Failure, List<CircleInvite>>> getMyInvites();

  /// Accept an invite
  Future<Either<Failure, CircleMember>> acceptInvite(String inviteId);

  /// Decline an invite
  Future<Either<Failure, Unit>> declineInvite(String inviteId);

  /// Generate invite code for circle (admin only)
  Future<Either<Failure, String>> generateInviteCode(String circleId);

  // ============ CONTRIBUTION OPERATIONS ============

  /// Get contributions for a circle
  Future<Either<Failure, List<Contribution>>> getContributions({
    required String circleId,
    int? cycleNumber,
    ContributionStatus? status,
  });

  /// Get current user's pending contributions
  Future<Either<Failure, List<Contribution>>> getMyPendingContributions();

  /// Get current user's contribution history
  Future<Either<Failure, List<Contribution>>> getMyContributionHistory({
    String? circleId,
    int page = 1,
    int limit = 20,
  });

  /// Make a contribution
  Future<Either<Failure, Contribution>> makeContribution({
    required String circleId,
    required String contributionId,
    required Pin pin,
  });

  /// Make contribution with auto-debit from wallet
  Future<Either<Failure, Contribution>> makeContributionFromWallet({
    required String circleId,
    required String contributionId,
    required Pin pin,
  });

  // ============ PAYOUT OPERATIONS ============

  /// Get payouts for a circle
  Future<Either<Failure, List<Payout>>> getPayouts({
    required String circleId,
    PayoutStatus? status,
  });

  /// Get current user's payout history
  Future<Either<Failure, List<Payout>>> getMyPayoutHistory({
    String? circleId,
    int page = 1,
    int limit = 20,
  });

  /// Get next scheduled payout for a circle
  Future<Either<Failure, Payout?>> getNextPayout(String circleId);

  // ============ STATS & ANALYTICS ============

  /// Get current user's savings statistics
  Future<Either<Failure, SavingsStats>> getSavingsStats();

  /// Get circle activity/history
  Future<Either<Failure, List<CircleActivity>>> getCircleActivity({
    required String circleId,
    int page = 1,
    int limit = 20,
  });
}

/// Circle activity item
class CircleActivity {
  final String id;
  final String circleId;
  final String type;
  final String description;
  final String? actorId;
  final String? actorName;
  final Money? amount;
  final DateTime createdAt;

  const CircleActivity({
    required this.id,
    required this.circleId,
    required this.type,
    required this.description,
    this.actorId,
    this.actorName,
    this.amount,
    required this.createdAt,
  });
}
