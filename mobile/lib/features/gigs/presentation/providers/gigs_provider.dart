import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/di/providers.dart';
import '../../../core/repositories/base_repository.dart';
import '../data/models/gig_model.dart';
import '../data/repositories/gigs_repository.dart';

/// State for gigs list with pagination
class GigsListState {
  final List<Gig> gigs;
  final bool isLoading;
  final bool hasMore;
  final int currentPage;
  final GigFilter filter;
  final String? error;

  const GigsListState({
    this.gigs = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.filter = const GigFilter(),
    this.error,
  });

  GigsListState copyWith({
    List<Gig>? gigs,
    bool? isLoading,
    bool? hasMore,
    int? currentPage,
    GigFilter? filter,
    String? error,
  }) {
    return GigsListState(
      gigs: gigs ?? this.gigs,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      filter: filter ?? this.filter,
      error: error,
    );
  }
}

/// Notifier for browsing gigs
class GigsListNotifier extends StateNotifier<GigsListState> {
  final GigsRepository _repository;

  GigsListNotifier(this._repository) : super(const GigsListState());

  Future<void> loadGigs({bool refresh = false}) async {
    if (state.isLoading) return;
    if (!refresh && !state.hasMore) return;

    state = state.copyWith(
      isLoading: true,
      error: null,
      currentPage: refresh ? 1 : state.currentPage,
    );

    final filter = state.filter.copyWith(page: refresh ? 1 : state.currentPage);
    final result = await _repository.getGigs(filter);

    result.when(
      success: (data) {
        final newGigs = refresh ? data.gigs : [...state.gigs, ...data.gigs];
        state = state.copyWith(
          gigs: newGigs,
          isLoading: false,
          hasMore: data.hasMore,
          currentPage: data.page + 1,
        );
      },
      failure: (message, _) {
        state = state.copyWith(
          isLoading: false,
          error: message,
        );
      },
    );
  }

  Future<void> refresh() => loadGigs(refresh: true);

  void updateFilter(GigFilter filter) {
    state = state.copyWith(filter: filter);
    loadGigs(refresh: true);
  }

  void setCategory(GigCategory? category) {
    updateFilter(state.filter.copyWith(category: category));
  }

  void setSearch(String? search) {
    updateFilter(state.filter.copyWith(search: search));
  }
}

/// Provider for browsing all gigs
final gigsListProvider = StateNotifierProvider<GigsListNotifier, GigsListState>((ref) {
  final repository = ref.watch(gigsRepositoryProvider);
  return GigsListNotifier(repository);
});

/// State for user's gigs (as client or freelancer)
class MyGigsState {
  final List<Gig> clientGigs;
  final List<Gig> freelancerGigs;
  final List<Proposal> proposals;
  final bool isLoading;
  final String? error;

  const MyGigsState({
    this.clientGigs = const [],
    this.freelancerGigs = const [],
    this.proposals = const [],
    this.isLoading = false,
    this.error,
  });

  MyGigsState copyWith({
    List<Gig>? clientGigs,
    List<Gig>? freelancerGigs,
    List<Proposal>? proposals,
    bool? isLoading,
    String? error,
  }) {
    return MyGigsState(
      clientGigs: clientGigs ?? this.clientGigs,
      freelancerGigs: freelancerGigs ?? this.freelancerGigs,
      proposals: proposals ?? this.proposals,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// Notifier for user's own gigs
class MyGigsNotifier extends StateNotifier<MyGigsState> {
  final GigsRepository _repository;

  MyGigsNotifier(this._repository) : super(const MyGigsState());

  Future<void> loadAll() async {
    state = state.copyWith(isLoading: true, error: null);

    // Load all in parallel
    final results = await Future.wait([
      _repository.getMyClientGigs(),
      _repository.getMyFreelancerGigs(),
      _repository.getMyProposals(),
    ]);

    final clientResult = results[0] as Result<PaginatedGigs>;
    final freelancerResult = results[1] as Result<PaginatedGigs>;
    final proposalsResult = results[2] as Result<List<Proposal>>;

    state = state.copyWith(
      clientGigs: clientResult.data?.gigs ?? [],
      freelancerGigs: freelancerResult.data?.gigs ?? [],
      proposals: proposalsResult.data ?? [],
      isLoading: false,
      error: clientResult.error ?? freelancerResult.error ?? proposalsResult.error,
    );
  }
}

/// Provider for user's own gigs
final myGigsProvider = StateNotifierProvider<MyGigsNotifier, MyGigsState>((ref) {
  final repository = ref.watch(gigsRepositoryProvider);
  return MyGigsNotifier(repository);
});

/// Provider for a single gig details
final gigDetailProvider = FutureProvider.family<Gig?, String>((ref, gigId) async {
  final repository = ref.watch(gigsRepositoryProvider);
  final result = await repository.getGig(gigId);
  return result.data;
});

/// Provider for gig proposals
final gigProposalsProvider = FutureProvider.family<List<Proposal>, String>((ref, gigId) async {
  final repository = ref.watch(gigsRepositoryProvider);
  final result = await repository.getGigProposals(gigId);
  return result.data ?? [];
});

/// Provider for featured gigs (home screen)
final featuredGigsProvider = FutureProvider<List<Gig>>((ref) async {
  final repository = ref.watch(gigsRepositoryProvider);
  final result = await repository.getFeaturedGigs(limit: 5);
  return result.data ?? [];
});

/// Provider for category counts
final gigCategoryCountsProvider = FutureProvider<Map<GigCategory, int>>((ref) async {
  final repository = ref.watch(gigsRepositoryProvider);
  final result = await repository.getCategoryCounts();
  return result.data ?? {};
});
