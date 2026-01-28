import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/api/api_client.dart';
import '../../../../core/exceptions/api_exception.dart';

/// Gigs API Service
/// Handles all gig marketplace-related API calls
class GigsService {
  final ApiClient _apiClient;

  GigsService(this._apiClient);

  /// Get all available gigs with filters
  Future<GigsListResponse> getGigs({
    int page = 1,
    int perPage = 20,
    String? category,
    String? status,
    double? minBudget,
    double? maxBudget,
    String? location,
    String? search,
    String? sortBy,
    String? sortOrder,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (category != null) queryParams['category'] = category;
      if (status != null) queryParams['status'] = status;
      if (minBudget != null) queryParams['min_budget'] = minBudget;
      if (maxBudget != null) queryParams['max_budget'] = maxBudget;
      if (location != null) queryParams['location'] = location;
      if (search != null) queryParams['search'] = search;
      if (sortBy != null) queryParams['sort_by'] = sortBy;
      if (sortOrder != null) queryParams['sort_order'] = sortOrder;

      final response = await _apiClient.get(
        '/api/v1/gigs',
        queryParameters: queryParams,
      );
      return GigsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get single gig details
  Future<GigDetail> getGigDetails(String gigId) async {
    try {
      final response = await _apiClient.get('/api/v1/gigs/$gigId');
      return GigDetail.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Create a new gig
  Future<GigDetail> createGig({
    required String title,
    required String description,
    required String category,
    required double budget,
    String? budgetType, // fixed, hourly, negotiable
    required int durationDays,
    List<String>? skills,
    String? location,
    bool isRemote = true,
    List<String>? attachments,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/gigs', data: {
        'title': title,
        'description': description,
        'category': category,
        'budget': budget,
        'budget_type': budgetType ?? 'fixed',
        'duration_days': durationDays,
        'skills': skills ?? [],
        'location': location,
        'is_remote': isRemote,
        'attachments': attachments ?? [],
      });
      return GigDetail.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Update an existing gig
  Future<GigDetail> updateGig(String gigId, Map<String, dynamic> updates) async {
    try {
      final response = await _apiClient.patch('/api/v1/gigs/$gigId', data: updates);
      return GigDetail.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Delete a gig
  Future<void> deleteGig(String gigId) async {
    try {
      await _apiClient.delete('/api/v1/gigs/$gigId');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get gigs posted by the current user
  Future<GigsListResponse> getMyPostedGigs({
    int page = 1,
    int perPage = 20,
    String? status,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (status != null) queryParams['status'] = status;

      final response = await _apiClient.get(
        '/api/v1/gigs/my-posts',
        queryParameters: queryParams,
      );
      return GigsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Submit a proposal for a gig
  Future<Proposal> submitProposal({
    required String gigId,
    required String coverLetter,
    required double bidAmount,
    required int estimatedDays,
    List<String>? attachments,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/gigs/$gigId/proposals', data: {
        'cover_letter': coverLetter,
        'bid_amount': bidAmount,
        'estimated_days': estimatedDays,
        'attachments': attachments ?? [],
      });
      return Proposal.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get proposals for a gig (for gig owner)
  Future<ProposalsListResponse> getGigProposals(String gigId, {
    int page = 1,
    int perPage = 20,
  }) async {
    try {
      final response = await _apiClient.get(
        '/api/v1/gigs/$gigId/proposals',
        queryParameters: {'page': page, 'per_page': perPage},
      );
      return ProposalsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get my submitted proposals
  Future<ProposalsListResponse> getMyProposals({
    int page = 1,
    int perPage = 20,
    String? status,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (status != null) queryParams['status'] = status;

      final response = await _apiClient.get(
        '/api/v1/proposals/my',
        queryParameters: queryParams,
      );
      return ProposalsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Withdraw a proposal
  Future<void> withdrawProposal(String proposalId) async {
    try {
      await _apiClient.delete('/api/v1/proposals/$proposalId');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Accept a proposal (for gig owner)
  Future<Contract> acceptProposal(String proposalId) async {
    try {
      final response = await _apiClient.post('/api/v1/proposals/$proposalId/accept');
      return Contract.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Reject a proposal (for gig owner)
  Future<void> rejectProposal(String proposalId, {String? reason}) async {
    try {
      await _apiClient.post('/api/v1/proposals/$proposalId/reject', data: {
        if (reason != null) 'reason': reason,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get active contracts
  Future<ContractsListResponse> getContracts({
    int page = 1,
    int perPage = 20,
    String? role, // client, freelancer
    String? status,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (role != null) queryParams['role'] = role;
      if (status != null) queryParams['status'] = status;

      final response = await _apiClient.get(
        '/api/v1/contracts',
        queryParameters: queryParams,
      );
      return ContractsListResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get contract details
  Future<Contract> getContractDetails(String contractId) async {
    try {
      final response = await _apiClient.get('/api/v1/contracts/$contractId');
      return Contract.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Submit work for contract milestone
  Future<void> submitMilestoneWork(String contractId, String milestoneId, {
    required String description,
    List<String>? attachments,
  }) async {
    try {
      await _apiClient.post('/api/v1/contracts/$contractId/milestones/$milestoneId/submit', data: {
        'description': description,
        'attachments': attachments ?? [],
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Approve milestone (for client)
  Future<void> approveMilestone(String contractId, String milestoneId) async {
    try {
      await _apiClient.post('/api/v1/contracts/$contractId/milestones/$milestoneId/approve');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Request revision for milestone (for client)
  Future<void> requestMilestoneRevision(String contractId, String milestoneId, {
    required String feedback,
  }) async {
    try {
      await _apiClient.post('/api/v1/contracts/$contractId/milestones/$milestoneId/revision', data: {
        'feedback': feedback,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Complete contract
  Future<void> completeContract(String contractId) async {
    try {
      await _apiClient.post('/api/v1/contracts/$contractId/complete');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Leave a review
  Future<void> leaveReview(String contractId, {
    required int rating,
    String? comment,
  }) async {
    try {
      await _apiClient.post('/api/v1/contracts/$contractId/review', data: {
        'rating': rating,
        'comment': comment,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get gig categories
  Future<CategoriesResponse> getCategories() async {
    try {
      final response = await _apiClient.get('/api/v1/gigs/categories');
      return CategoriesResponse.fromJson(response.data);
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
class GigsListResponse {
  final List<GigSummary> gigs;
  final PaginationMeta meta;

  GigsListResponse({required this.gigs, required this.meta});

  factory GigsListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return GigsListResponse(
      gigs: data.map((g) => GigSummary.fromJson(g)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class GigSummary {
  final String id;
  final String title;
  final String description;
  final String category;
  final double budget;
  final String budgetType;
  final int durationDays;
  final String status;
  final int proposalCount;
  final int viewCount;
  final String? location;
  final bool isRemote;
  final List<String> skills;
  final GigOwner owner;
  final DateTime createdAt;

  GigSummary({
    required this.id,
    required this.title,
    required this.description,
    required this.category,
    required this.budget,
    required this.budgetType,
    required this.durationDays,
    required this.status,
    required this.proposalCount,
    required this.viewCount,
    this.location,
    required this.isRemote,
    required this.skills,
    required this.owner,
    required this.createdAt,
  });

  factory GigSummary.fromJson(Map<String, dynamic> json) {
    return GigSummary(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      description: json['description'] ?? '',
      category: json['category'] ?? '',
      budget: (json['budget'] ?? 0).toDouble(),
      budgetType: json['budget_type'] ?? 'fixed',
      durationDays: json['duration_days'] ?? 0,
      status: json['status'] ?? 'open',
      proposalCount: json['proposal_count'] ?? 0,
      viewCount: json['view_count'] ?? 0,
      location: json['location'],
      isRemote: json['is_remote'] ?? true,
      skills: List<String>.from(json['skills'] ?? []),
      owner: GigOwner.fromJson(json['owner'] ?? {}),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class GigDetail extends GigSummary {
  final List<String> attachments;
  final List<Milestone>? milestones;
  final bool hasApplied;
  final String? myProposalId;

  GigDetail({
    required super.id,
    required super.title,
    required super.description,
    required super.category,
    required super.budget,
    required super.budgetType,
    required super.durationDays,
    required super.status,
    required super.proposalCount,
    required super.viewCount,
    super.location,
    required super.isRemote,
    required super.skills,
    required super.owner,
    required super.createdAt,
    required this.attachments,
    this.milestones,
    this.hasApplied = false,
    this.myProposalId,
  });

  factory GigDetail.fromJson(Map<String, dynamic> json) {
    return GigDetail(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      description: json['description'] ?? '',
      category: json['category'] ?? '',
      budget: (json['budget'] ?? 0).toDouble(),
      budgetType: json['budget_type'] ?? 'fixed',
      durationDays: json['duration_days'] ?? 0,
      status: json['status'] ?? 'open',
      proposalCount: json['proposal_count'] ?? 0,
      viewCount: json['view_count'] ?? 0,
      location: json['location'],
      isRemote: json['is_remote'] ?? true,
      skills: List<String>.from(json['skills'] ?? []),
      owner: GigOwner.fromJson(json['owner'] ?? {}),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      attachments: List<String>.from(json['attachments'] ?? []),
      milestones: json['milestones'] != null
          ? (json['milestones'] as List).map((m) => Milestone.fromJson(m)).toList()
          : null,
      hasApplied: json['has_applied'] ?? false,
      myProposalId: json['my_proposal_id'],
    );
  }
}

class GigOwner {
  final String id;
  final String name;
  final String? avatar;
  final double? rating;
  final int? completedGigs;
  final DateTime? joinedAt;

  GigOwner({
    required this.id,
    required this.name,
    this.avatar,
    this.rating,
    this.completedGigs,
    this.joinedAt,
  });

  factory GigOwner.fromJson(Map<String, dynamic> json) {
    return GigOwner(
      id: json['id'] ?? '',
      name: json['name'] ?? 'Unknown',
      avatar: json['avatar'],
      rating: json['rating']?.toDouble(),
      completedGigs: json['completed_gigs'],
      joinedAt: json['joined_at'] != null ? DateTime.tryParse(json['joined_at']) : null,
    );
  }
}

class Milestone {
  final String id;
  final String title;
  final String? description;
  final double amount;
  final int order;
  final String status; // pending, in_progress, submitted, approved, paid

  Milestone({
    required this.id,
    required this.title,
    this.description,
    required this.amount,
    required this.order,
    required this.status,
  });

  factory Milestone.fromJson(Map<String, dynamic> json) {
    return Milestone(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      description: json['description'],
      amount: (json['amount'] ?? 0).toDouble(),
      order: json['order'] ?? 0,
      status: json['status'] ?? 'pending',
    );
  }
}

class ProposalsListResponse {
  final List<Proposal> proposals;
  final PaginationMeta meta;

  ProposalsListResponse({required this.proposals, required this.meta});

  factory ProposalsListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return ProposalsListResponse(
      proposals: data.map((p) => Proposal.fromJson(p)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class Proposal {
  final String id;
  final String gigId;
  final String coverLetter;
  final double bidAmount;
  final int estimatedDays;
  final String status; // pending, shortlisted, accepted, rejected, withdrawn
  final List<String> attachments;
  final ProposalFreelancer? freelancer;
  final GigSummary? gig;
  final DateTime createdAt;

  Proposal({
    required this.id,
    required this.gigId,
    required this.coverLetter,
    required this.bidAmount,
    required this.estimatedDays,
    required this.status,
    required this.attachments,
    this.freelancer,
    this.gig,
    required this.createdAt,
  });

  factory Proposal.fromJson(Map<String, dynamic> json) {
    return Proposal(
      id: json['id'] ?? '',
      gigId: json['gig_id'] ?? '',
      coverLetter: json['cover_letter'] ?? '',
      bidAmount: (json['bid_amount'] ?? 0).toDouble(),
      estimatedDays: json['estimated_days'] ?? 0,
      status: json['status'] ?? 'pending',
      attachments: List<String>.from(json['attachments'] ?? []),
      freelancer: json['freelancer'] != null ? ProposalFreelancer.fromJson(json['freelancer']) : null,
      gig: json['gig'] != null ? GigSummary.fromJson(json['gig']) : null,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class ProposalFreelancer {
  final String id;
  final String name;
  final String? avatar;
  final double? rating;
  final int completedGigs;
  final List<String> skills;

  ProposalFreelancer({
    required this.id,
    required this.name,
    this.avatar,
    this.rating,
    required this.completedGigs,
    required this.skills,
  });

  factory ProposalFreelancer.fromJson(Map<String, dynamic> json) {
    return ProposalFreelancer(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      avatar: json['avatar'],
      rating: json['rating']?.toDouble(),
      completedGigs: json['completed_gigs'] ?? 0,
      skills: List<String>.from(json['skills'] ?? []),
    );
  }
}

class ContractsListResponse {
  final List<Contract> contracts;
  final PaginationMeta meta;

  ContractsListResponse({required this.contracts, required this.meta});

  factory ContractsListResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return ContractsListResponse(
      contracts: data.map((c) => Contract.fromJson(c)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class Contract {
  final String id;
  final String gigId;
  final String gigTitle;
  final double amount;
  final int durationDays;
  final String status; // active, completed, cancelled, disputed
  final GigOwner client;
  final ProposalFreelancer freelancer;
  final List<ContractMilestone> milestones;
  final DateTime startDate;
  final DateTime? endDate;
  final DateTime createdAt;

  Contract({
    required this.id,
    required this.gigId,
    required this.gigTitle,
    required this.amount,
    required this.durationDays,
    required this.status,
    required this.client,
    required this.freelancer,
    required this.milestones,
    required this.startDate,
    this.endDate,
    required this.createdAt,
  });

  factory Contract.fromJson(Map<String, dynamic> json) {
    return Contract(
      id: json['id'] ?? '',
      gigId: json['gig_id'] ?? '',
      gigTitle: json['gig_title'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      durationDays: json['duration_days'] ?? 0,
      status: json['status'] ?? 'active',
      client: GigOwner.fromJson(json['client'] ?? {}),
      freelancer: ProposalFreelancer.fromJson(json['freelancer'] ?? {}),
      milestones: (json['milestones'] as List? ?? [])
          .map((m) => ContractMilestone.fromJson(m))
          .toList(),
      startDate: DateTime.tryParse(json['start_date'] ?? '') ?? DateTime.now(),
      endDate: json['end_date'] != null ? DateTime.tryParse(json['end_date']) : null,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class ContractMilestone extends Milestone {
  final String? submissionDescription;
  final List<String>? submissionAttachments;
  final DateTime? submittedAt;
  final DateTime? approvedAt;
  final DateTime? paidAt;

  ContractMilestone({
    required super.id,
    required super.title,
    super.description,
    required super.amount,
    required super.order,
    required super.status,
    this.submissionDescription,
    this.submissionAttachments,
    this.submittedAt,
    this.approvedAt,
    this.paidAt,
  });

  factory ContractMilestone.fromJson(Map<String, dynamic> json) {
    return ContractMilestone(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      description: json['description'],
      amount: (json['amount'] ?? 0).toDouble(),
      order: json['order'] ?? 0,
      status: json['status'] ?? 'pending',
      submissionDescription: json['submission_description'],
      submissionAttachments: json['submission_attachments'] != null
          ? List<String>.from(json['submission_attachments'])
          : null,
      submittedAt: json['submitted_at'] != null ? DateTime.tryParse(json['submitted_at']) : null,
      approvedAt: json['approved_at'] != null ? DateTime.tryParse(json['approved_at']) : null,
      paidAt: json['paid_at'] != null ? DateTime.tryParse(json['paid_at']) : null,
    );
  }
}

class CategoriesResponse {
  final List<GigCategory> categories;

  CategoriesResponse({required this.categories});

  factory CategoriesResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return CategoriesResponse(
      categories: data.map((c) => GigCategory.fromJson(c)).toList(),
    );
  }
}

class GigCategory {
  final String id;
  final String name;
  final String? icon;
  final int gigCount;

  GigCategory({
    required this.id,
    required this.name,
    this.icon,
    required this.gigCount,
  });

  factory GigCategory.fromJson(Map<String, dynamic> json) {
    return GigCategory(
      id: json['id'] ?? '',
      name: json['name'] ?? '',
      icon: json['icon'],
      gigCount: json['gig_count'] ?? 0,
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
final gigsServiceProvider = Provider<GigsService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return GigsService(apiClient);
});
