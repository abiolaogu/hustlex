import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/di/providers.dart';
import '../data/models/savings_model.dart';
import '../data/repositories/savings_repository.dart';

/// State for user's savings circles
class MySavingsState {
  final List<SavingsCircle> myCircles;
  final List<SavingsCircle> discoverCircles;
  final SavingsStats? stats;
  final List<Contribution> pendingContributions;
  final bool isLoading;
  final String? error;

  const MySavingsState({
    this.myCircles = const [],
    this.discoverCircles = const [],
    this.stats,
    this.pendingContributions = const [],
    this.isLoading = false,
    this.error,
  });

  MySavingsState copyWith({
    List<SavingsCircle>? myCircles,
    List<SavingsCircle>? discoverCircles,
    SavingsStats? stats,
    List<Contribution>? pendingContributions,
    bool? isLoading,
    String? error,
  }) {
    return MySavingsState(
      myCircles: myCircles ?? this.myCircles,
      discoverCircles: discoverCircles ?? this.discoverCircles,
      stats: stats ?? this.stats,
      pendingContributions: pendingContributions ?? this.pendingContributions,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }

  /// Active circles only
  List<SavingsCircle> get activeCircles =>
      myCircles.where((c) => c.status == CircleStatus.active).toList();

  /// Total amount saved
  double get totalSaved => stats?.totalSaved ?? 0;
}

/// Notifier for savings state
class MySavingsNotifier extends StateNotifier<MySavingsState> {
  final SavingsRepository _repository;

  MySavingsNotifier(this._repository) : super(const MySavingsState());

  Future<void> loadAll() async {
    state = state.copyWith(isLoading: true, error: null);

    // Load all in parallel
    final results = await Future.wait([
      _repository.getMyCircles(),
      _repository.getDiscoverCircles(limit: 10),
      _repository.getSavingsStats(),
      _repository.getPendingContributions(),
    ]);

    state = state.copyWith(
      myCircles: (results[0].data as List<SavingsCircle>?) ?? [],
      discoverCircles: (results[1].data as List<SavingsCircle>?) ?? [],
      stats: results[2].data as SavingsStats?,
      pendingContributions: (results[3].data as List<Contribution>?) ?? [],
      isLoading: false,
      error: results[0].error ?? results[1].error ?? results[2].error,
    );
  }

  Future<void> refresh() => loadAll();

  /// Join a circle
  Future<bool> joinCircle(String circleId, {String? inviteCode}) async {
    final result = await _repository.joinCircle(
      JoinCircleRequest(circleId: circleId, inviteCode: inviteCode),
    );
    
    if (result.isSuccess) {
      await loadAll();
      return true;
    }
    
    state = state.copyWith(error: result.error);
    return false;
  }

  /// Leave a circle
  Future<bool> leaveCircle(String circleId) async {
    final result = await _repository.leaveCircle(circleId);
    
    if (result.isSuccess) {
      await loadAll();
      return true;
    }
    
    state = state.copyWith(error: result.error);
    return false;
  }

  /// Make a contribution
  Future<bool> makeContribution(MakeContributionRequest request) async {
    final result = await _repository.makeContribution(request);
    
    if (result.isSuccess) {
      await loadAll();
      return true;
    }
    
    state = state.copyWith(error: result.error);
    return false;
  }
}

/// Provider for user's savings state
final mySavingsProvider = StateNotifierProvider<MySavingsNotifier, MySavingsState>((ref) {
  final repository = ref.watch(savingsRepositoryProvider);
  return MySavingsNotifier(repository);
});

/// State for browsing/discovering circles
class CirclesListState {
  final List<SavingsCircle> circles;
  final bool isLoading;
  final bool hasMore;
  final int currentPage;
  final CircleFilter filter;
  final String? error;

  const CirclesListState({
    this.circles = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.filter = const CircleFilter(onlyJoinable: true),
    this.error,
  });

  CirclesListState copyWith({
    List<SavingsCircle>? circles,
    bool? isLoading,
    bool? hasMore,
    int? currentPage,
    CircleFilter? filter,
    String? error,
  }) {
    return CirclesListState(
      circles: circles ?? this.circles,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      filter: filter ?? this.filter,
      error: error,
    );
  }
}

/// Notifier for browsing circles
class CirclesListNotifier extends StateNotifier<CirclesListState> {
  final SavingsRepository _repository;

  CirclesListNotifier(this._repository) : super(const CirclesListState());

  Future<void> loadCircles({bool refresh = false}) async {
    if (state.isLoading) return;
    if (!refresh && !state.hasMore) return;

    state = state.copyWith(
      isLoading: true,
      error: null,
      currentPage: refresh ? 1 : state.currentPage,
    );

    final filter = state.filter.copyWith(page: refresh ? 1 : state.currentPage);
    final result = await _repository.getCircles(filter);

    result.when(
      success: (data) {
        final newCircles = refresh ? data.circles : [...state.circles, ...data.circles];
        state = state.copyWith(
          circles: newCircles,
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

  Future<void> refresh() => loadCircles(refresh: true);

  void updateFilter(CircleFilter filter) {
    state = state.copyWith(filter: filter);
    loadCircles(refresh: true);
  }

  void setType(CircleType? type) {
    updateFilter(state.filter.copyWith(type: type));
  }

  void setFrequency(ContributionFrequency? frequency) {
    updateFilter(state.filter.copyWith(frequency: frequency));
  }
}

/// Provider for browsing circles
final circlesListProvider = StateNotifierProvider<CirclesListNotifier, CirclesListState>((ref) {
  final repository = ref.watch(savingsRepositoryProvider);
  return CirclesListNotifier(repository);
});

/// Provider for a single circle details
final circleDetailProvider = FutureProvider.family<SavingsCircle?, String>((ref, circleId) async {
  final repository = ref.watch(savingsRepositoryProvider);
  final result = await repository.getCircle(circleId);
  return result.data;
});

/// Provider for circle members
final circleMembersProvider = FutureProvider.family<List<CircleMember>, String>((ref, circleId) async {
  final repository = ref.watch(savingsRepositoryProvider);
  final result = await repository.getCircleMembers(circleId);
  return result.data ?? [];
});

/// Provider for circle contributions
final circleContributionsProvider = FutureProvider.family<List<Contribution>, String>((ref, circleId) async {
  final repository = ref.watch(savingsRepositoryProvider);
  final result = await repository.getCircleContributions(circleId);
  return result.data ?? [];
});

/// Provider for circle payouts
final circlePayoutsProvider = FutureProvider.family<List<Payout>, String>((ref, circleId) async {
  final repository = ref.watch(savingsRepositoryProvider);
  final result = await repository.getCirclePayouts(circleId);
  return result.data ?? [];
});

/// Provider for pending invites
final savingsInvitesProvider = FutureProvider<List<CircleInvite>>((ref) async {
  final repository = ref.watch(savingsRepositoryProvider);
  final result = await repository.getMyInvites();
  return result.data ?? [];
});

/// Provider for active circles (home screen)
final activeCirclesProvider = FutureProvider<List<SavingsCircle>>((ref) async {
  final repository = ref.watch(savingsRepositoryProvider);
  final result = await repository.getActiveCircles(limit: 3);
  return result.data ?? [];
});

/// Provider for savings stats
final savingsStatsProvider = FutureProvider<SavingsStats>((ref) async {
  final repository = ref.watch(savingsRepositoryProvider);
  final result = await repository.getSavingsStats();
  return result.data ?? const SavingsStats();
});
