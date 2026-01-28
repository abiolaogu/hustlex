import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/api/api_client.dart';
import '../../../../core/exceptions/api_exception.dart';

/// Savings API Service
/// Handles all Ajo/Esusu savings circles API calls
class SavingsService {
  final ApiClient _apiClient;

  SavingsService(this._apiClient);

  /// Get my savings circles
  Future<CirclesListResponse> getMyCircles({
    int page = 1,
    int perPage = 20,
    String? status, // active, completed, pending
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (status != null) queryParams['status'] = status;

      final response = await _apiClient.get(
        '/api/v1/savings/circles/my',
        queryParameters: queryParams,
      );
      return CirclesListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get available circles to join
  Future<CirclesListResponse> getAvailableCircles({
    int page = 1,
    int perPage = 20,
    String? type, // rotating, fixed
    String? frequency, // weekly, monthly
    double? minContribution,
    double? maxContribution,
    String? search,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (type != null) queryParams['type'] = type;
      if (frequency != null) queryParams['frequency'] = frequency;
      if (minContribution != null) queryParams['min_contribution'] = minContribution;
      if (maxContribution != null) queryParams['max_contribution'] = maxContribution;
      if (search != null) queryParams['search'] = search;

      final response = await _apiClient.get(
        '/api/v1/savings/circles/available',
        queryParameters: queryParams,
      );
      return CirclesListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get circle details
  Future<CircleDetail> getCircleDetails(String circleId) async {
    try {
      final response = await _apiClient.get('/api/v1/savings/circles/$circleId');
      return CircleDetail.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Create a new savings circle
  Future<CircleDetail> createCircle({
    required String name,
    required String description,
    required String type, // rotating, fixed
    required String frequency, // weekly, biweekly, monthly
    required double contributionAmount,
    required int maxMembers,
    required DateTime startDate,
    bool isPrivate = false,
    String? rules,
    int? penaltyPercentage,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/savings/circles', data: {
        'name': name,
        'description': description,
        'type': type,
        'frequency': frequency,
        'contribution_amount': contributionAmount,
        'max_members': maxMembers,
        'start_date': startDate.toIso8601String(),
        'is_private': isPrivate,
        'rules': rules,
        'penalty_percentage': penaltyPercentage ?? 5,
      });
      return CircleDetail.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Join a savings circle
  Future<CircleMembership> joinCircle(String circleId, {
    int? preferredSlot, // For rotating circles
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/savings/circles/$circleId/join', data: {
        if (preferredSlot != null) 'preferred_slot': preferredSlot,
      });
      return CircleMembership.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Leave a savings circle (before start)
  Future<void> leaveCircle(String circleId) async {
    try {
      await _apiClient.post('/api/v1/savings/circles/$circleId/leave');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Make contribution to circle
  Future<ContributionResponse> makeContribution(String circleId, {
    required double amount,
    String? pin,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/savings/circles/$circleId/contribute', data: {
        'amount': amount,
        if (pin != null) 'pin': pin,
      });
      return ContributionResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get circle contributions history
  Future<ContributionsListResponse> getCircleContributions(String circleId, {
    int page = 1,
    int perPage = 50,
    int? roundNumber,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (roundNumber != null) queryParams['round'] = roundNumber;

      final response = await _apiClient.get(
        '/api/v1/savings/circles/$circleId/contributions',
        queryParameters: queryParams,
      );
      return ContributionsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get my contributions across all circles
  Future<ContributionsListResponse> getMyContributions({
    int page = 1,
    int perPage = 20,
  }) async {
    try {
      final response = await _apiClient.get(
        '/api/v1/savings/contributions/my',
        queryParameters: {'page': page, 'per_page': perPage},
      );
      return ContributionsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get circle payouts history
  Future<PayoutsListResponse> getCirclePayouts(String circleId, {
    int page = 1,
    int perPage = 20,
  }) async {
    try {
      final response = await _apiClient.get(
        '/api/v1/savings/circles/$circleId/payouts',
        queryParameters: {'page': page, 'per_page': perPage},
      );
      return PayoutsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Request early payout (if allowed by circle rules)
  Future<EarlyPayoutResponse> requestEarlyPayout(String circleId, {
    required String reason,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/savings/circles/$circleId/early-payout', data: {
        'reason': reason,
      });
      return EarlyPayoutResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Invite member to private circle
  Future<void> inviteMember(String circleId, {
    required String phone,
    String? message,
  }) async {
    try {
      await _apiClient.post('/api/v1/savings/circles/$circleId/invite', data: {
        'phone': phone,
        if (message != null) 'message': message,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get circle members
  Future<MembersListResponse> getCircleMembers(String circleId) async {
    try {
      final response = await _apiClient.get('/api/v1/savings/circles/$circleId/members');
      return MembersListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get savings stats/summary
  Future<SavingsStats> getSavingsStats() async {
    try {
      final response = await _apiClient.get('/api/v1/savings/stats');
      return SavingsStats.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  ApiException _handleError(dynamic error) {
    if (error is ApiException) return error;
    return ApiException(message: error.toString());
  }
}

// Response models
class CirclesListResponse {
  final List<CircleSummary> circles;
  final PaginationMeta meta;

  CirclesListResponse({required this.circles, required this.meta});

  factory CirclesListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return CirclesListResponse(
      circles: data.map((c) => CircleSummary.fromJson(c)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class CircleSummary {
  final String id;
  final String name;
  final String description;
  final String type; // rotating, fixed
  final String frequency; // weekly, biweekly, monthly
  final double contributionAmount;
  final int maxMembers;
  final int currentMembers;
  final String status; // pending, active, completed
  final int currentRound;
  final int totalRounds;
  final double totalPooled;
  final DateTime startDate;
  final DateTime? endDate;
  final bool isPrivate;
  final CircleAdmin admin;
  final DateTime createdAt;

  CircleSummary({
    required this.id,
    required this.name,
    required this.description,
    required this.type,
    required this.frequency,
    required this.contributionAmount,
    required this.maxMembers,
    required this.currentMembers,
    required this.status,
    required this.currentRound,
    required this.totalRounds,
    required this.totalPooled,
    required this.startDate,
    this.endDate,
    required this.isPrivate,
    required this.admin,
    required this.createdAt,
  });

  factory CircleSummary.fromJson(Map<String, dynamic> json) {
    return CircleSummary(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      type: json['type'] ?? 'rotating',
      frequency: json['frequency'] ?? 'monthly',
      contributionAmount: (json['contribution_amount'] ?? 0).toDouble(),
      maxMembers: json['max_members'] ?? 0,
      currentMembers: json['current_members'] ?? 0,
      status: json['status'] ?? 'pending',
      currentRound: json['current_round'] ?? 0,
      totalRounds: json['total_rounds'] ?? 0,
      totalPooled: (json['total_pooled'] ?? 0).toDouble(),
      startDate: DateTime.tryParse(json['start_date'] ?? '') ?? DateTime.now(),
      endDate: json['end_date'] != null ? DateTime.tryParse(json['end_date']) : null,
      isPrivate: json['is_private'] ?? false,
      admin: CircleAdmin.fromJson(json['admin'] ?? {}),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }

  int get availableSlots => maxMembers - currentMembers;
  bool get isFull => currentMembers >= maxMembers;
}

class CircleDetail extends CircleSummary {
  final String? rules;
  final int penaltyPercentage;
  final double nextContributionAmount;
  final DateTime? nextContributionDate;
  final DateTime? nextPayoutDate;
  final String? nextPayoutRecipient;
  final List<CircleMember> members;
  final CircleMembership? myMembership;
  final List<RoundInfo> rounds;

  CircleDetail({
    required super.id,
    required super.name,
    required super.description,
    required super.type,
    required super.frequency,
    required super.contributionAmount,
    required super.maxMembers,
    required super.currentMembers,
    required super.status,
    required super.currentRound,
    required super.totalRounds,
    required super.totalPooled,
    required super.startDate,
    super.endDate,
    required super.isPrivate,
    required super.admin,
    required super.createdAt,
    this.rules,
    required this.penaltyPercentage,
    required this.nextContributionAmount,
    this.nextContributionDate,
    this.nextPayoutDate,
    this.nextPayoutRecipient,
    required this.members,
    this.myMembership,
    required this.rounds,
  });

  factory CircleDetail.fromJson(Map<String, dynamic> json) {
    return CircleDetail(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      type: json['type'] ?? 'rotating',
      frequency: json['frequency'] ?? 'monthly',
      contributionAmount: (json['contribution_amount'] ?? 0).toDouble(),
      maxMembers: json['max_members'] ?? 0,
      currentMembers: json['current_members'] ?? 0,
      status: json['status'] ?? 'pending',
      currentRound: json['current_round'] ?? 0,
      totalRounds: json['total_rounds'] ?? 0,
      totalPooled: (json['total_pooled'] ?? 0).toDouble(),
      startDate: DateTime.tryParse(json['start_date'] ?? '') ?? DateTime.now(),
      endDate: json['end_date'] != null ? DateTime.tryParse(json['end_date']) : null,
      isPrivate: json['is_private'] ?? false,
      admin: CircleAdmin.fromJson(json['admin'] ?? {}),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      rules: json['rules'],
      penaltyPercentage: json['penalty_percentage'] ?? 5,
      nextContributionAmount: (json['next_contribution_amount'] ?? json['contribution_amount'] ?? 0).toDouble(),
      nextContributionDate: json['next_contribution_date'] != null
          ? DateTime.tryParse(json['next_contribution_date'])
          : null,
      nextPayoutDate: json['next_payout_date'] != null
          ? DateTime.tryParse(json['next_payout_date'])
          : null,
      nextPayoutRecipient: json['next_payout_recipient'],
      members: (json['members'] as List? ?? [])
          .map((m) => CircleMember.fromJson(m))
          .toList(),
      myMembership: json['my_membership'] != null
          ? CircleMembership.fromJson(json['my_membership'])
          : null,
      rounds: (json['rounds'] as List? ?? [])
          .map((r) => RoundInfo.fromJson(r))
          .toList(),
    );
  }
}

class CircleAdmin {
  final String id;
  final String name;
  final String? avatar;
  final String? phone;

  CircleAdmin({
    required this.id,
    required this.name,
    this.avatar,
    this.phone,
  });

  factory CircleAdmin.fromJson(Map<String, dynamic> json) {
    return CircleAdmin(
      id: json['id'] ?? '',
      name: json['name'] ?? 'Unknown',
      avatar: json['avatar'],
      phone: json['phone'],
    );
  }
}

class CircleMember {
  final String id;
  final String userId;
  final String name;
  final String? avatar;
  final int slotNumber; // Position in rotation
  final String status; // active, defaulted, left
  final double totalContributed;
  final int contributionCount;
  final bool hasReceivedPayout;
  final DateTime joinedAt;

  CircleMember({
    required this.id,
    required this.userId,
    required this.name,
    this.avatar,
    required this.slotNumber,
    required this.status,
    required this.totalContributed,
    required this.contributionCount,
    required this.hasReceivedPayout,
    required this.joinedAt,
  });

  factory CircleMember.fromJson(Map<String, dynamic> json) {
    return CircleMember(
      id: json['id'] ?? '',
      userId: json['user_id'] ?? '',
      name: json['name'] ?? '',
      avatar: json['avatar'],
      slotNumber: json['slot_number'] ?? 0,
      status: json['status'] ?? 'active',
      totalContributed: (json['total_contributed'] ?? 0).toDouble(),
      contributionCount: json['contribution_count'] ?? 0,
      hasReceivedPayout: json['has_received_payout'] ?? false,
      joinedAt: DateTime.tryParse(json['joined_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class CircleMembership {
  final String id;
  final String circleId;
  final int slotNumber;
  final String status;
  final double totalContributed;
  final int contributionCount;
  final bool hasReceivedPayout;
  final int? payoutRound;
  final bool hasPendingContribution;
  final DateTime? nextContributionDue;
  final DateTime joinedAt;

  CircleMembership({
    required this.id,
    required this.circleId,
    required this.slotNumber,
    required this.status,
    required this.totalContributed,
    required this.contributionCount,
    required this.hasReceivedPayout,
    this.payoutRound,
    required this.hasPendingContribution,
    this.nextContributionDue,
    required this.joinedAt,
  });

  factory CircleMembership.fromJson(Map<String, dynamic> json) {
    return CircleMembership(
      id: json['id'] ?? '',
      circleId: json['circle_id'] ?? '',
      slotNumber: json['slot_number'] ?? 0,
      status: json['status'] ?? 'active',
      totalContributed: (json['total_contributed'] ?? 0).toDouble(),
      contributionCount: json['contribution_count'] ?? 0,
      hasReceivedPayout: json['has_received_payout'] ?? false,
      payoutRound: json['payout_round'],
      hasPendingContribution: json['has_pending_contribution'] ?? false,
      nextContributionDue: json['next_contribution_due'] != null
          ? DateTime.tryParse(json['next_contribution_due'])
          : null,
      joinedAt: DateTime.tryParse(json['joined_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class RoundInfo {
  final int roundNumber;
  final String status; // pending, active, completed
  final String? recipientId;
  final String? recipientName;
  final double targetAmount;
  final double collectedAmount;
  final DateTime? startDate;
  final DateTime? endDate;
  final DateTime? payoutDate;

  RoundInfo({
    required this.roundNumber,
    required this.status,
    this.recipientId,
    this.recipientName,
    required this.targetAmount,
    required this.collectedAmount,
    this.startDate,
    this.endDate,
    this.payoutDate,
  });

  factory RoundInfo.fromJson(Map<String, dynamic> json) {
    return RoundInfo(
      roundNumber: json['round_number'] ?? 0,
      status: json['status'] ?? 'pending',
      recipientId: json['recipient_id'],
      recipientName: json['recipient_name'],
      targetAmount: (json['target_amount'] ?? 0).toDouble(),
      collectedAmount: (json['collected_amount'] ?? 0).toDouble(),
      startDate: json['start_date'] != null ? DateTime.tryParse(json['start_date']) : null,
      endDate: json['end_date'] != null ? DateTime.tryParse(json['end_date']) : null,
      payoutDate: json['payout_date'] != null ? DateTime.tryParse(json['payout_date']) : null,
    );
  }

  double get progressPercentage => targetAmount > 0 ? (collectedAmount / targetAmount) * 100 : 0;
}

class ContributionResponse {
  final bool success;
  final String message;
  final Contribution? contribution;

  ContributionResponse({
    required this.success,
    required this.message,
    this.contribution,
  });

  factory ContributionResponse.fromJson(Map<String, dynamic> json) {
    return ContributionResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      contribution: json['data'] != null ? Contribution.fromJson(json['data']) : null,
    );
  }
}

class ContributionsListResponse {
  final List<Contribution> contributions;
  final PaginationMeta meta;

  ContributionsListResponse({required this.contributions, required this.meta});

  factory ContributionsListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return ContributionsListResponse(
      contributions: data.map((c) => Contribution.fromJson(c)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class Contribution {
  final String id;
  final String circleId;
  final String? circleName;
  final String memberId;
  final String? memberName;
  final int roundNumber;
  final double amount;
  final String status; // pending, completed, failed
  final String? transactionRef;
  final DateTime createdAt;

  Contribution({
    required this.id,
    required this.circleId,
    this.circleName,
    required this.memberId,
    this.memberName,
    required this.roundNumber,
    required this.amount,
    required this.status,
    this.transactionRef,
    required this.createdAt,
  });

  factory Contribution.fromJson(Map<String, dynamic> json) {
    return Contribution(
      id: json['id'] ?? '',
      circleId: json['circle_id'] ?? '',
      circleName: json['circle_name'],
      memberId: json['member_id'] ?? '',
      memberName: json['member_name'],
      roundNumber: json['round_number'] ?? 0,
      amount: (json['amount'] ?? 0).toDouble(),
      status: json['status'] ?? 'pending',
      transactionRef: json['transaction_ref'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class PayoutsListResponse {
  final List<Payout> payouts;
  final PaginationMeta meta;

  PayoutsListResponse({required this.payouts, required this.meta});

  factory PayoutsListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return PayoutsListResponse(
      payouts: data.map((p) => Payout.fromJson(p)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class Payout {
  final String id;
  final String circleId;
  final String recipientId;
  final String recipientName;
  final int roundNumber;
  final double amount;
  final String status; // pending, processing, completed, failed
  final String? transactionRef;
  final DateTime? paidAt;
  final DateTime createdAt;

  Payout({
    required this.id,
    required this.circleId,
    required this.recipientId,
    required this.recipientName,
    required this.roundNumber,
    required this.amount,
    required this.status,
    this.transactionRef,
    this.paidAt,
    required this.createdAt,
  });

  factory Payout.fromJson(Map<String, dynamic> json) {
    return Payout(
      id: json['id'] ?? '',
      circleId: json['circle_id'] ?? '',
      recipientId: json['recipient_id'] ?? '',
      recipientName: json['recipient_name'] ?? '',
      roundNumber: json['round_number'] ?? 0,
      amount: (json['amount'] ?? 0).toDouble(),
      status: json['status'] ?? 'pending',
      transactionRef: json['transaction_ref'],
      paidAt: json['paid_at'] != null ? DateTime.tryParse(json['paid_at']) : null,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class EarlyPayoutResponse {
  final bool approved;
  final String message;
  final Payout? payout;
  final double? penaltyAmount;

  EarlyPayoutResponse({
    required this.approved,
    required this.message,
    this.payout,
    this.penaltyAmount,
  });

  factory EarlyPayoutResponse.fromJson(Map<String, dynamic> json) {
    return EarlyPayoutResponse(
      approved: json['approved'] ?? false,
      message: json['message'] ?? '',
      payout: json['data'] != null ? Payout.fromJson(json['data']) : null,
      penaltyAmount: json['penalty_amount']?.toDouble(),
    );
  }
}

class MembersListResponse {
  final List<CircleMember> members;

  MembersListResponse({required this.members});

  factory MembersListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return MembersListResponse(
      members: data.map((m) => CircleMember.fromJson(m)).toList(),
    );
  }
}

class SavingsStats {
  final double totalSaved;
  final double totalReceived;
  final int activeCircles;
  final int completedCircles;
  final int upcomingPayouts;
  final double nextPayoutAmount;
  final DateTime? nextPayoutDate;

  SavingsStats({
    required this.totalSaved,
    required this.totalReceived,
    required this.activeCircles,
    required this.completedCircles,
    required this.upcomingPayouts,
    required this.nextPayoutAmount,
    this.nextPayoutDate,
  });

  factory SavingsStats.fromJson(Map<String, dynamic> json) {
    return SavingsStats(
      totalSaved: (json['total_saved'] ?? 0).toDouble(),
      totalReceived: (json['total_received'] ?? 0).toDouble(),
      activeCircles: json['active_circles'] ?? 0,
      completedCircles: json['completed_circles'] ?? 0,
      upcomingPayouts: json['upcoming_payouts'] ?? 0,
      nextPayoutAmount: (json['next_payout_amount'] ?? 0).toDouble(),
      nextPayoutDate: json['next_payout_date'] != null
          ? DateTime.tryParse(json['next_payout_date'])
          : null,
    );
  }
}

class PaginationMeta {
  final int currentPage;
  final int lastPage;
  final int perPage;
  final int total;

  PaginationMeta({
    required this.currentPage,
    required this.lastPage,
    required this.perPage,
    required this.total,
  });

  factory PaginationMeta.fromJson(Map<String, dynamic> json) {
    return PaginationMeta(
      currentPage: json['current_page'] ?? 1,
      lastPage: json['last_page'] ?? 1,
      perPage: json['per_page'] ?? 20,
      total: json['total'] ?? 0,
    );
  }

  bool get hasNextPage => currentPage < lastPage;
}

// Provider
final savingsServiceProvider = Provider<SavingsService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return SavingsService(apiClient);
});
