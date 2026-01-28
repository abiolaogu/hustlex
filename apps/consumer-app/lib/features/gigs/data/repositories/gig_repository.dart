import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/api/api_client.dart';
import '../../../../core/repositories/base_repository.dart';
import '../../domain/models/gig_models.dart';

/// =============================================================================
/// GIG REPOSITORY PROVIDER
/// =============================================================================

final gigRepositoryProvider = Provider<GigRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return GigRepository(apiClient);
});

/// =============================================================================
/// GIG REPOSITORY
/// =============================================================================

class GigRepository extends BaseRepository {
  final ApiClient _apiClient;

  GigRepository(this._apiClient);

  // ===================== GIGS =====================

  /// Fetch paginated list of gigs
  Future<Result<PaginatedResponse<Gig>>> getGigs({
    int page = 1,
    int limit = 20,
    String? category,
    String? status,
    double? minBudget,
    double? maxBudget,
    bool? isRemote,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      
      if (category != null) queryParams['category'] = category;
      if (status != null) queryParams['status'] = status;
      if (minBudget != null) queryParams['min_budget'] = minBudget;
      if (maxBudget != null) queryParams['max_budget'] = maxBudget;
      if (isRemote != null) queryParams['is_remote'] = isRemote;
      if (search != null && search.isNotEmpty) queryParams['search'] = search;
      if (sortBy != null) queryParams['sort_by'] = sortBy;
      if (sortOrder != null) queryParams['sort_order'] = sortOrder;

      return await _apiClient.getPaginated<Gig>(
        '/gigs',
        queryParameters: queryParams,
        fromJson: Gig.fromJson,
        page: page,
        limit: limit,
      );
    });
  }

  /// Get a single gig by ID
  Future<Result<Gig>> getGigById(String gigId) {
    return safeCall(() async {
      final response = await _apiClient.get<Map<String, dynamic>>(
        '/gigs/$gigId',
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to get gig');
      }
      
      return Gig.fromJson(response.data!);
    });
  }

  /// Create a new gig
  Future<Result<Gig>> createGig(CreateGigRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/gigs',
        data: request.toJson(),
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to create gig');
      }
      
      return Gig.fromJson(response.data!);
    });
  }

  /// Update an existing gig
  Future<Result<Gig>> updateGig(String gigId, UpdateGigRequest request) {
    return safeCall(() async {
      final response = await _apiClient.put<Map<String, dynamic>>(
        '/gigs/$gigId',
        data: request.toJson(),
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to update gig');
      }
      
      return Gig.fromJson(response.data!);
    });
  }

  /// Delete a gig
  Future<Result<void>> deleteGig(String gigId) {
    return safeVoidCall(() async {
      final response = await _apiClient.delete('/gigs/$gigId');
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to delete gig');
      }
    });
  }

  /// Close a gig (stop accepting proposals)
  Future<Result<Gig>> closeGig(String gigId) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/gigs/$gigId/close',
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to close gig');
      }
      
      return Gig.fromJson(response.data!);
    });
  }

  /// Get my posted gigs
  Future<Result<PaginatedResponse<Gig>>> getMyGigs({
    int page = 1,
    int limit = 20,
    String? status,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (status != null) queryParams['status'] = status;

      return await _apiClient.getPaginated<Gig>(
        '/gigs/my-gigs',
        queryParameters: queryParams,
        fromJson: Gig.fromJson,
        page: page,
        limit: limit,
      );
    });
  }

  // ===================== PROPOSALS =====================

  /// Get proposals for a gig
  Future<Result<PaginatedResponse<Proposal>>> getGigProposals(
    String gigId, {
    int page = 1,
    int limit = 20,
    String? status,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (status != null) queryParams['status'] = status;

      return await _apiClient.getPaginated<Proposal>(
        '/gigs/$gigId/proposals',
        queryParameters: queryParams,
        fromJson: Proposal.fromJson,
        page: page,
        limit: limit,
      );
    });
  }

  /// Submit a proposal for a gig
  Future<Result<Proposal>> submitProposal(
    String gigId,
    SubmitProposalRequest request,
  ) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/gigs/$gigId/proposals',
        data: request.toJson(),
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to submit proposal');
      }
      
      return Proposal.fromJson(response.data!);
    });
  }

  /// Accept a proposal
  Future<Result<Proposal>> acceptProposal(String proposalId) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/proposals/$proposalId/accept',
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to accept proposal');
      }
      
      return Proposal.fromJson(response.data!);
    });
  }

  /// Reject a proposal
  Future<Result<Proposal>> rejectProposal(String proposalId) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/proposals/$proposalId/reject',
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to reject proposal');
      }
      
      return Proposal.fromJson(response.data!);
    });
  }

  /// Withdraw my proposal
  Future<Result<void>> withdrawProposal(String proposalId) {
    return safeVoidCall(() async {
      final response = await _apiClient.delete('/proposals/$proposalId');
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to withdraw proposal');
      }
    });
  }

  /// Get my submitted proposals
  Future<Result<PaginatedResponse<Proposal>>> getMyProposals({
    int page = 1,
    int limit = 20,
    String? status,
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (status != null) queryParams['status'] = status;

      return await _apiClient.getPaginated<Proposal>(
        '/proposals/my-proposals',
        queryParameters: queryParams,
        fromJson: Proposal.fromJson,
        page: page,
        limit: limit,
      );
    });
  }

  // ===================== CONTRACTS =====================

  /// Get contracts (both as client and freelancer)
  Future<Result<PaginatedResponse<GigContract>>> getContracts({
    int page = 1,
    int limit = 20,
    String? status,
    String? role, // 'client' or 'freelancer'
  }) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{};
      if (status != null) queryParams['status'] = status;
      if (role != null) queryParams['role'] = role;

      return await _apiClient.getPaginated<GigContract>(
        '/contracts',
        queryParameters: queryParams,
        fromJson: GigContract.fromJson,
        page: page,
        limit: limit,
      );
    });
  }

  /// Get a single contract
  Future<Result<GigContract>> getContractById(String contractId) {
    return safeCall(() async {
      final response = await _apiClient.get<Map<String, dynamic>>(
        '/contracts/$contractId',
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to get contract');
      }
      
      return GigContract.fromJson(response.data!);
    });
  }

  /// Complete a contract (mark work as done)
  Future<Result<GigContract>> completeContract(String contractId) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/contracts/$contractId/complete',
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to complete contract');
      }
      
      return GigContract.fromJson(response.data!);
    });
  }

  /// Approve contract completion and release payment
  Future<Result<GigContract>> approveContract(String contractId) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/contracts/$contractId/approve',
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to approve contract');
      }
      
      return GigContract.fromJson(response.data!);
    });
  }

  /// Request revision on contract
  Future<Result<GigContract>> requestRevision(
    String contractId, {
    required String reason,
  }) {
    return safeCall(() async {
      final response = await _apiClient.post<Map<String, dynamic>>(
        '/contracts/$contractId/revision',
        data: {'reason': reason},
        fromJson: (data) => data as Map<String, dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to request revision');
      }
      
      return GigContract.fromJson(response.data!);
    });
  }

  /// Rate a completed contract
  Future<Result<void>> rateContract(
    String contractId, {
    required int rating,
    String? review,
  }) {
    return safeVoidCall(() async {
      final response = await _apiClient.post(
        '/contracts/$contractId/rate',
        data: {
          'rating': rating,
          if (review != null) 'review': review,
        },
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to rate contract');
      }
    });
  }

  // ===================== CATEGORIES =====================

  /// Get all gig categories
  Future<Result<List<GigCategory>>> getCategories() {
    return safeCall(() async {
      final response = await _apiClient.get<List<dynamic>>(
        '/gigs/categories',
        fromJson: (data) => data as List<dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to get categories');
      }
      
      return response.data!
          .map((item) => GigCategory.fromJson(item as Map<String, dynamic>))
          .toList();
    });
  }

  /// Get popular skills
  Future<Result<List<String>>> getPopularSkills() {
    return safeCall(() async {
      final response = await _apiClient.get<List<dynamic>>(
        '/gigs/skills',
        fromJson: (data) => data as List<dynamic>,
      );
      
      if (response.hasError) {
        throw Exception(response.error ?? 'Failed to get skills');
      }
      
      return response.data!.map((item) => item.toString()).toList();
    });
  }
}

/// =============================================================================
/// REQUEST MODELS
/// =============================================================================

class CreateGigRequest {
  final String title;
  final String description;
  final String category;
  final double budgetMin;
  final double budgetMax;
  final int duration;
  final List<String> skills;
  final bool isRemote;

  CreateGigRequest({
    required this.title,
    required this.description,
    required this.category,
    required this.budgetMin,
    required this.budgetMax,
    required this.duration,
    required this.skills,
    this.isRemote = true,
  });

  Map<String, dynamic> toJson() => {
        'title': title,
        'description': description,
        'category': category,
        'budget_min': budgetMin,
        'budget_max': budgetMax,
        'duration': duration,
        'skills': skills,
        'is_remote': isRemote,
      };
}

class UpdateGigRequest {
  final String? title;
  final String? description;
  final String? category;
  final double? budgetMin;
  final double? budgetMax;
  final int? duration;
  final List<String>? skills;
  final bool? isRemote;
  final String? status;

  UpdateGigRequest({
    this.title,
    this.description,
    this.category,
    this.budgetMin,
    this.budgetMax,
    this.duration,
    this.skills,
    this.isRemote,
    this.status,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (title != null) json['title'] = title;
    if (description != null) json['description'] = description;
    if (category != null) json['category'] = category;
    if (budgetMin != null) json['budget_min'] = budgetMin;
    if (budgetMax != null) json['budget_max'] = budgetMax;
    if (duration != null) json['duration'] = duration;
    if (skills != null) json['skills'] = skills;
    if (isRemote != null) json['is_remote'] = isRemote;
    if (status != null) json['status'] = status;
    return json;
  }
}

class SubmitProposalRequest {
  final String coverLetter;
  final double amount;
  final int deliveryDays;

  SubmitProposalRequest({
    required this.coverLetter,
    required this.amount,
    required this.deliveryDays,
  });

  Map<String, dynamic> toJson() => {
        'cover_letter': coverLetter,
        'amount': amount,
        'delivery_days': deliveryDays,
      };
}
