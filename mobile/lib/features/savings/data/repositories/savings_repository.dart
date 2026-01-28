import '../../../core/api/api_client.dart';
import '../../../core/repositories/base_repository.dart';
import '../models/savings_model.dart';

class SavingsRepository extends BaseRepository {
  final ApiClient _apiClient;

  SavingsRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  // Circles

  /// Get paginated list of circles with filters
  Future<Result<PaginatedCircles>> getCircles(CircleFilter filter) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{
        'page': filter.page,
        'limit': filter.limit,
      };

      if (filter.type != null) {
        queryParams['type'] = filter.type!.name;
      }
      if (filter.status != null) {
        queryParams['status'] = filter.status!.name;
      }
      if (filter.frequency != null) {
        queryParams['frequency'] = filter.frequency!.name;
      }
      if (filter.minContribution != null) {
        queryParams['min_contribution'] = filter.minContribution;
      }
      if (filter.maxContribution != null) {
        queryParams['max_contribution'] = filter.maxContribution;
      }
      if (filter.onlyJoinable) {
        queryParams['only_joinable'] = true;
      }
      if (filter.search != null && filter.search!.isNotEmpty) {
        queryParams['search'] = filter.search;
      }

      final response = await _apiClient.get('/circles', queryParameters: queryParams);
      return PaginatedCircles.fromJson(response.data['data']);
    });
  }

  /// Get a single circle by ID
  Future<Result<SavingsCircle>> getCircle(String circleId) {
    return safeCall(() async {
      final response = await _apiClient.get('/circles/$circleId');
      return SavingsCircle.fromJson(response.data['data']);
    });
  }

  /// Create a new savings circle
  Future<Result<SavingsCircle>> createCircle(CreateCircleRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/circles', data: request.toJson());
      return SavingsCircle.fromJson(response.data['data']);
    });
  }

  /// Update circle settings (admin only)
  Future<Result<SavingsCircle>> updateCircle(String circleId, Map<String, dynamic> updates) {
    return safeCall(() async {
      final response = await _apiClient.patch('/circles/$circleId', data: updates);
      return SavingsCircle.fromJson(response.data['data']);
    });
  }

  /// Start a circle (admin only, after minimum members joined)
  Future<Result<SavingsCircle>> startCircle(String circleId) {
    return safeCall(() async {
      final response = await _apiClient.post('/circles/$circleId/start');
      return SavingsCircle.fromJson(response.data['data']);
    });
  }

  /// Get circles the current user is a member of
  Future<Result<List<SavingsCircle>>> getMyCircles({CircleStatus? status}) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (status != null) {
        queryParams['status'] = status.name;
      }

      final response = await _apiClient.get('/circles/my', queryParameters: queryParams);
      final data = response.data['data'] as List;
      return data.map((e) => SavingsCircle.fromJson(e)).toList();
    });
  }

  /// Get active circles (for home screen)
  Future<Result<List<SavingsCircle>>> getActiveCircles({int limit = 5}) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/circles/my',
        queryParameters: {'status': 'active', 'limit': limit},
      );
      final data = response.data['data'] as List;
      return data.map((e) => SavingsCircle.fromJson(e)).toList();
    });
  }

  /// Get discover circles (joinable circles)
  Future<Result<List<SavingsCircle>>> getDiscoverCircles({int limit = 10}) {
    return getCircles(CircleFilter(onlyJoinable: true, limit: limit)).then(
      (result) => result.map((data) => data.circles),
    );
  }

  /// Get savings stats for current user
  Future<Result<SavingsStats>> getSavingsStats() {
    return safeCall(() async {
      final response = await _apiClient.get('/savings/stats');
      return SavingsStats.fromJson(response.data['data']);
    });
  }

  // Membership

  /// Join a circle
  Future<Result<CircleMember>> joinCircle(JoinCircleRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/circles/${request.circleId}/join',
        data: request.inviteCode != null ? {'invite_code': request.inviteCode} : null,
      );
      return CircleMember.fromJson(response.data['data']);
    });
  }

  /// Leave a circle
  Future<Result<void>> leaveCircle(String circleId) {
    return safeVoidCall(() async {
      await _apiClient.post('/circles/$circleId/leave');
    });
  }

  /// Get members of a circle
  Future<Result<List<CircleMember>>> getCircleMembers(String circleId) {
    return safeCall(() async {
      final response = await _apiClient.get('/circles/$circleId/members');
      final data = response.data['data'] as List;
      return data.map((e) => CircleMember.fromJson(e)).toList();
    });
  }

  /// Remove a member (admin only)
  Future<Result<void>> removeMember(String circleId, String memberId) {
    return safeVoidCall(() async {
      await _apiClient.delete('/circles/$circleId/members/$memberId');
    });
  }

  /// Set payout order (admin only, for Ajo circles)
  Future<Result<void>> setPayoutOrder(String circleId, List<String> memberIds) {
    return safeVoidCall(() async {
      await _apiClient.post(
        '/circles/$circleId/payout-order',
        data: {'member_ids': memberIds},
      );
    });
  }

  // Invitations

  /// Generate invite code for a private circle
  Future<Result<String>> generateInviteCode(String circleId) {
    return safeCall(() async {
      final response = await _apiClient.post('/circles/$circleId/invite');
      return response.data['data']['invite_code'] as String;
    });
  }

  /// Invite a user by phone
  Future<Result<CircleInvite>> inviteUser(String circleId, String phone) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/circles/$circleId/invite',
        data: {'phone': phone},
      );
      return CircleInvite.fromJson(response.data['data']);
    });
  }

  /// Get pending invites for current user
  Future<Result<List<CircleInvite>>> getMyInvites() {
    return safeCall(() async {
      final response = await _apiClient.get('/circles/invites');
      final data = response.data['data'] as List;
      return data.map((e) => CircleInvite.fromJson(e)).toList();
    });
  }

  /// Accept an invite
  Future<Result<CircleMember>> acceptInvite(String inviteId) {
    return safeCall(() async {
      final response = await _apiClient.post('/circles/invites/$inviteId/accept');
      return CircleMember.fromJson(response.data['data']);
    });
  }

  /// Decline an invite
  Future<Result<void>> declineInvite(String inviteId) {
    return safeVoidCall(() async {
      await _apiClient.post('/circles/invites/$inviteId/decline');
    });
  }

  // Contributions

  /// Get contributions for a circle
  Future<Result<List<Contribution>>> getCircleContributions(
    String circleId, {
    int? cycleNumber,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (cycleNumber != null) {
        queryParams['cycle'] = cycleNumber;
      }

      final response = await _apiClient.get(
        '/circles/$circleId/contributions',
        queryParameters: queryParams,
      );
      final data = response.data['data'] as List;
      return data.map((e) => Contribution.fromJson(e)).toList();
    });
  }

  /// Get current user's contributions across all circles
  Future<Result<List<Contribution>>> getMyContributions({ContributionStatus? status}) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (status != null) {
        queryParams['status'] = status.name;
      }

      final response = await _apiClient.get('/contributions/my', queryParameters: queryParams);
      final data = response.data['data'] as List;
      return data.map((e) => Contribution.fromJson(e)).toList();
    });
  }

  /// Get pending contributions (for reminders)
  Future<Result<List<Contribution>>> getPendingContributions() {
    return getMyContributions(status: ContributionStatus.pending);
  }

  /// Make a contribution
  Future<Result<Contribution>> makeContribution(MakeContributionRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/circles/${request.circleId}/contributions/${request.contributionId}/pay',
        data: {
          'payment_method': request.paymentMethod,
          if (request.paymentReference != null)
            'payment_reference': request.paymentReference,
        },
      );
      return Contribution.fromJson(response.data['data']);
    });
  }

  // Payouts

  /// Get payouts for a circle
  Future<Result<List<Payout>>> getCirclePayouts(String circleId) {
    return safeCall(() async {
      final response = await _apiClient.get('/circles/$circleId/payouts');
      final data = response.data['data'] as List;
      return data.map((e) => Payout.fromJson(e)).toList();
    });
  }

  /// Get current user's payouts
  Future<Result<List<Payout>>> getMyPayouts() {
    return safeCall(() async {
      final response = await _apiClient.get('/payouts/my');
      final data = response.data['data'] as List;
      return data.map((e) => Payout.fromJson(e)).toList();
    });
  }

  // Circle activity / history

  /// Get circle activity feed
  Future<Result<List<Map<String, dynamic>>> getCircleActivity(
    String circleId, {
    int page = 1,
    int limit = 20,
  }) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/circles/$circleId/activity',
        queryParameters: {'page': page, 'limit': limit},
      );
      return (response.data['data'] as List).cast<Map<String, dynamic>>();
    });
  }
}
