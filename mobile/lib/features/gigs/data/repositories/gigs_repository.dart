import '../../../core/api/api_client.dart';
import '../../../core/repositories/base_repository.dart';
import '../models/gig_model.dart';

class GigsRepository extends BaseRepository {
  final ApiClient _apiClient;

  GigsRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Get paginated list of gigs with filters
  Future<Result<PaginatedGigs>> getGigs(GigFilter filter) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{
        'page': filter.page,
        'limit': filter.limit,
        'sort_by': filter.sortBy,
        'sort_order': filter.sortOrder,
      };

      if (filter.category != null) {
        queryParams['category'] = filter.category!.name;
      }
      if (filter.status != null) {
        queryParams['status'] = filter.status!.name;
      }
      if (filter.minBudget != null) {
        queryParams['min_budget'] = filter.minBudget;
      }
      if (filter.maxBudget != null) {
        queryParams['max_budget'] = filter.maxBudget;
      }
      if (filter.isRemote != null) {
        queryParams['is_remote'] = filter.isRemote;
      }
      if (filter.search != null && filter.search!.isNotEmpty) {
        queryParams['search'] = filter.search;
      }

      final response = await _apiClient.get('/gigs', queryParameters: queryParams);
      return PaginatedGigs.fromJson(response.data['data']);
    });
  }

  /// Get a single gig by ID
  Future<Result<Gig>> getGig(String gigId) {
    return safeCall(() async {
      final response = await _apiClient.get('/gigs/$gigId');
      return Gig.fromJson(response.data['data']);
    });
  }

  /// Create a new gig
  Future<Result<Gig>> createGig(CreateGigRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post('/gigs', data: request.toJson());
      return Gig.fromJson(response.data['data']);
    });
  }

  /// Update an existing gig
  Future<Result<Gig>> updateGig(String gigId, Map<String, dynamic> updates) {
    return safeCall(() async {
      final response = await _apiClient.patch('/gigs/$gigId', data: updates);
      return Gig.fromJson(response.data['data']);
    });
  }

  /// Delete a gig
  Future<Result<void>> deleteGig(String gigId) {
    return safeVoidCall(() async {
      await _apiClient.delete('/gigs/$gigId');
    });
  }

  /// Get gigs created by current user (as client)
  Future<Result<PaginatedGigs>> getMyClientGigs({
    GigStatus? status,
    int page = 1,
    int limit = 20,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{
        'page': page,
        'limit': limit,
        'role': 'client',
      };
      if (status != null) {
        queryParams['status'] = status.name;
      }

      final response = await _apiClient.get('/gigs/my', queryParameters: queryParams);
      return PaginatedGigs.fromJson(response.data['data']);
    });
  }

  /// Get gigs where user is freelancer
  Future<Result<PaginatedGigs>> getMyFreelancerGigs({
    GigStatus? status,
    int page = 1,
    int limit = 20,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{
        'page': page,
        'limit': limit,
        'role': 'freelancer',
      };
      if (status != null) {
        queryParams['status'] = status.name;
      }

      final response = await _apiClient.get('/gigs/my', queryParameters: queryParams);
      return PaginatedGigs.fromJson(response.data['data']);
    });
  }

  // Proposals

  /// Get proposals for a gig (as client)
  Future<Result<List<Proposal>>> getGigProposals(String gigId) {
    return safeCall(() async {
      final response = await _apiClient.get('/gigs/$gigId/proposals');
      final data = response.data['data'] as List;
      return data.map((e) => Proposal.fromJson(e)).toList();
    });
  }

  /// Get current user's proposals (as freelancer)
  Future<Result<List<Proposal>>> getMyProposals({ProposalStatus? status}) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (status != null) {
        queryParams['status'] = status.name;
      }

      final response = await _apiClient.get('/proposals/my', queryParameters: queryParams);
      final data = response.data['data'] as List;
      return data.map((e) => Proposal.fromJson(e)).toList();
    });
  }

  /// Submit a proposal for a gig
  Future<Result<Proposal>> submitProposal(CreateProposalRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/gigs/${request.gigId}/proposals',
        data: request.toJson(),
      );
      return Proposal.fromJson(response.data['data']);
    });
  }

  /// Update a proposal
  Future<Result<Proposal>> updateProposal(String proposalId, Map<String, dynamic> updates) {
    return safeCall(() async {
      final response = await _apiClient.patch('/proposals/$proposalId', data: updates);
      return Proposal.fromJson(response.data['data']);
    });
  }

  /// Withdraw a proposal
  Future<Result<void>> withdrawProposal(String proposalId) {
    return safeVoidCall(() async {
      await _apiClient.post('/proposals/$proposalId/withdraw');
    });
  }

  /// Accept a proposal (as client)
  Future<Result<void>> acceptProposal(String gigId, String proposalId) {
    return safeVoidCall(() async {
      await _apiClient.post('/gigs/$gigId/proposals/$proposalId/accept');
    });
  }

  /// Reject a proposal (as client)
  Future<Result<void>> rejectProposal(String gigId, String proposalId) {
    return safeVoidCall(() async {
      await _apiClient.post('/gigs/$gigId/proposals/$proposalId/reject');
    });
  }

  // Milestones

  /// Get milestones for a gig
  Future<Result<List<Milestone>>> getGigMilestones(String gigId) {
    return safeCall(() async {
      final response = await _apiClient.get('/gigs/$gigId/milestones');
      final data = response.data['data'] as List;
      return data.map((e) => Milestone.fromJson(e)).toList();
    });
  }

  /// Submit milestone for review (as freelancer)
  Future<Result<Milestone>> submitMilestone(String milestoneId, List<String> deliverables) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/milestones/$milestoneId/submit',
        data: {'deliverables': deliverables},
      );
      return Milestone.fromJson(response.data['data']);
    });
  }

  /// Approve milestone (as client)
  Future<Result<Milestone>> approveMilestone(String milestoneId) {
    return safeCall(() async {
      final response = await _apiClient.post('/milestones/$milestoneId/approve');
      return Milestone.fromJson(response.data['data']);
    });
  }

  /// Request revision on milestone (as client)
  Future<Result<Milestone>> requestRevision(String milestoneId, String feedback) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/milestones/$milestoneId/revision',
        data: {'feedback': feedback},
      );
      return Milestone.fromJson(response.data['data']);
    });
  }

  /// Release payment for milestone (as client)
  Future<Result<Milestone>> releaseMilestonePayment(String milestoneId, String pin) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/milestones/$milestoneId/release',
        data: {'pin': pin},
      );
      return Milestone.fromJson(response.data['data']);
    });
  }

  // Reviews

  /// Get reviews for a user
  Future<Result<List<GigReview>>> getUserReviews(String userId) {
    return safeCall(() async {
      final response = await _apiClient.get('/users/$userId/reviews');
      final data = response.data['data'] as List;
      return data.map((e) => GigReview.fromJson(e)).toList();
    });
  }

  /// Submit a review for a completed gig
  Future<Result<GigReview>> submitReview(String gigId, int rating, String? comment) {
    return safeCall(() async {
      final response = await _apiClient.post('/gigs/$gigId/review', data: {
        'rating': rating,
        'comment': comment,
      });
      return GigReview.fromJson(response.data['data']);
    });
  }

  // Categories

  /// Get available gig categories with counts
  Future<Result<Map<GigCategory, int>>> getCategoryCounts() {
    return safeCall(() async {
      final response = await _apiClient.get('/gigs/categories');
      final data = response.data['data'] as Map<String, dynamic>;
      return data.map((key, value) {
        final category = GigCategory.values.firstWhere(
          (e) => e.name == key,
          orElse: () => GigCategory.other,
        );
        return MapEntry(category, value as int);
      });
    });
  }

  // Search

  /// Search gigs with text query
  Future<Result<PaginatedGigs>> searchGigs(String query, {int page = 1, int limit = 20}) {
    return getGigs(GigFilter(search: query, page: page, limit: limit));
  }

  /// Get featured/recommended gigs
  Future<Result<List<Gig>>> getFeaturedGigs({int limit = 10}) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/gigs/featured',
        queryParameters: {'limit': limit},
      );
      final data = response.data['data'] as List;
      return data.map((e) => Gig.fromJson(e)).toList();
    });
  }
}
