import 'package:dartz/dartz.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../entities/gig.dart';
import '../entities/proposal.dart';
import '../entities/contract.dart';
import '../entities/milestone.dart';
import '../entities/review.dart';

/// Filter parameters for gig listings
class GigFilter {
  final GigCategory? category;
  final GigStatus? status;
  final Money? minBudget;
  final Money? maxBudget;
  final bool? isRemote;
  final String? search;
  final int page;
  final int limit;
  final String sortBy;
  final bool ascending;

  const GigFilter({
    this.category,
    this.status,
    this.minBudget,
    this.maxBudget,
    this.isRemote,
    this.search,
    this.page = 1,
    this.limit = 20,
    this.sortBy = 'createdAt',
    this.ascending = false,
  });

  GigFilter copyWith({
    GigCategory? category,
    GigStatus? status,
    Money? minBudget,
    Money? maxBudget,
    bool? isRemote,
    String? search,
    int? page,
    int? limit,
    String? sortBy,
    bool? ascending,
  }) {
    return GigFilter(
      category: category ?? this.category,
      status: status ?? this.status,
      minBudget: minBudget ?? this.minBudget,
      maxBudget: maxBudget ?? this.maxBudget,
      isRemote: isRemote ?? this.isRemote,
      search: search ?? this.search,
      page: page ?? this.page,
      limit: limit ?? this.limit,
      sortBy: sortBy ?? this.sortBy,
      ascending: ascending ?? this.ascending,
    );
  }
}

/// Paginated result for gigs
class PaginatedGigs {
  final List<Gig> gigs;
  final int total;
  final int page;
  final int limit;
  final bool hasMore;

  const PaginatedGigs({
    required this.gigs,
    required this.total,
    required this.page,
    required this.limit,
    required this.hasMore,
  });
}

/// Paginated result for proposals
class PaginatedProposals {
  final List<Proposal> proposals;
  final int total;
  final int page;
  final int limit;
  final bool hasMore;

  const PaginatedProposals({
    required this.proposals,
    required this.total,
    required this.page,
    required this.limit,
    required this.hasMore,
  });
}

/// Paginated result for contracts
class PaginatedContracts {
  final List<Contract> contracts;
  final int total;
  final int page;
  final int limit;
  final bool hasMore;

  const PaginatedContracts({
    required this.contracts,
    required this.total,
    required this.page,
    required this.limit,
    required this.hasMore,
  });
}

/// Abstract repository interface for gigs feature
abstract class GigsRepository {
  // ============ GIG OPERATIONS ============

  /// Get paginated list of gigs with optional filters
  Future<Either<Failure, PaginatedGigs>> getGigs(GigFilter filter);

  /// Get a single gig by ID
  Future<Either<Failure, Gig>> getGig(String gigId);

  /// Watch a gig for real-time updates
  Stream<Either<Failure, Gig>> watchGig(String gigId);

  /// Create a new gig
  Future<Either<Failure, Gig>> createGig({
    required String title,
    required String description,
    required GigCategory category,
    required Money budgetMin,
    required Money budgetMax,
    required int durationDays,
    required bool isRemote,
    String? location,
    DateTime? deadline,
    required List<String> skills,
    List<String>? attachments,
  });

  /// Update an existing gig
  Future<Either<Failure, Gig>> updateGig({
    required String gigId,
    String? title,
    String? description,
    GigCategory? category,
    Money? budgetMin,
    Money? budgetMax,
    int? durationDays,
    bool? isRemote,
    String? location,
    DateTime? deadline,
    List<String>? skills,
    List<String>? attachments,
  });

  /// Publish a draft gig
  Future<Either<Failure, Gig>> publishGig(String gigId);

  /// Cancel a gig
  Future<Either<Failure, Unit>> cancelGig(String gigId);

  /// Get gigs created by current user (client)
  Future<Either<Failure, PaginatedGigs>> getMyGigs({
    GigStatus? status,
    int page = 1,
    int limit = 20,
  });

  // ============ PROPOSAL OPERATIONS ============

  /// Get proposals for a gig
  Future<Either<Failure, PaginatedProposals>> getProposals({
    required String gigId,
    ProposalStatus? status,
    int page = 1,
    int limit = 20,
  });

  /// Get proposals submitted by current user (freelancer)
  Future<Either<Failure, PaginatedProposals>> getMyProposals({
    ProposalStatus? status,
    int page = 1,
    int limit = 20,
  });

  /// Submit a proposal for a gig
  Future<Either<Failure, Proposal>> submitProposal({
    required String gigId,
    required Money proposedAmount,
    required int deliveryDays,
    required String coverLetter,
    List<String>? attachments,
  });

  /// Update an existing proposal
  Future<Either<Failure, Proposal>> updateProposal({
    required String proposalId,
    Money? proposedAmount,
    int? deliveryDays,
    String? coverLetter,
    List<String>? attachments,
  });

  /// Withdraw a proposal
  Future<Either<Failure, Unit>> withdrawProposal(String proposalId);

  /// Accept a proposal (creates a contract)
  Future<Either<Failure, Contract>> acceptProposal(String proposalId);

  /// Reject a proposal
  Future<Either<Failure, Unit>> rejectProposal(String proposalId);

  // ============ CONTRACT OPERATIONS ============

  /// Get contracts for current user
  Future<Either<Failure, PaginatedContracts>> getContracts({
    ContractStatus? status,
    int page = 1,
    int limit = 20,
  });

  /// Get a single contract by ID
  Future<Either<Failure, Contract>> getContract(String contractId);

  /// Watch a contract for real-time updates
  Stream<Either<Failure, Contract>> watchContract(String contractId);

  /// Complete a contract
  Future<Either<Failure, Contract>> completeContract(String contractId);

  /// Request contract cancellation
  Future<Either<Failure, Unit>> requestCancellation({
    required String contractId,
    required String reason,
  });

  /// Raise a dispute on a contract
  Future<Either<Failure, Unit>> raiseDispute({
    required String contractId,
    required String reason,
    List<String>? evidence,
  });

  // ============ MILESTONE OPERATIONS ============

  /// Get milestones for a contract
  Future<Either<Failure, List<Milestone>>> getMilestones(String contractId);

  /// Create a milestone
  Future<Either<Failure, Milestone>> createMilestone({
    required String contractId,
    required String title,
    String? description,
    required Money amount,
    required int order,
    DateTime? dueDate,
  });

  /// Start working on a milestone
  Future<Either<Failure, Milestone>> startMilestone(String milestoneId);

  /// Submit a milestone for review
  Future<Either<Failure, Milestone>> submitMilestone({
    required String milestoneId,
    required List<String> deliverables,
  });

  /// Approve a submitted milestone
  Future<Either<Failure, Milestone>> approveMilestone(String milestoneId);

  /// Request revision on a milestone
  Future<Either<Failure, Milestone>> requestRevision({
    required String milestoneId,
    required String feedback,
  });

  /// Release payment for an approved milestone
  Future<Either<Failure, Milestone>> releaseMilestonePayment(String milestoneId);

  // ============ REVIEW OPERATIONS ============

  /// Submit a review for a completed contract
  Future<Either<Failure, Review>> submitReview({
    required String contractId,
    required double rating,
    String? comment,
  });

  /// Get reviews for a user
  Future<Either<Failure, List<Review>>> getUserReviews(String userId);

  /// Get review for a contract
  Future<Either<Failure, Review?>> getContractReview(String contractId);

  // ============ CATEGORY OPERATIONS ============

  /// Get all available gig categories
  Future<Either<Failure, List<GigCategoryInfo>>> getCategories();
}

/// Extended category info with metadata
class GigCategoryInfo {
  final GigCategory category;
  final String name;
  final String? icon;
  final String? description;
  final int gigCount;

  const GigCategoryInfo({
    required this.category,
    required this.name,
    this.icon,
    this.description,
    this.gigCount = 0,
  });
}
