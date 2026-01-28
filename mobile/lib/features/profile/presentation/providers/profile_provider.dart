import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/di/providers.dart';
import '../../../auth/data/models/user_model.dart';
import '../../data/repositories/profile_repository.dart';

/// Profile state
class ProfileState {
  final User? user;
  final KycStatus? kycStatus;
  final UserSettings? settings;
  final ReferralInfo? referralInfo;
  final bool isLoading;
  final String? error;

  const ProfileState({
    this.user,
    this.kycStatus,
    this.settings,
    this.referralInfo,
    this.isLoading = false,
    this.error,
  });

  ProfileState copyWith({
    User? user,
    KycStatus? kycStatus,
    UserSettings? settings,
    ReferralInfo? referralInfo,
    bool? isLoading,
    String? error,
  }) {
    return ProfileState(
      user: user ?? this.user,
      kycStatus: kycStatus ?? this.kycStatus,
      settings: settings ?? this.settings,
      referralInfo: referralInfo ?? this.referralInfo,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }

  bool get isKycVerified => kycStatus?.isVerified ?? false;
  String get kycLevel => kycStatus?.level ?? 'none';
}

/// Profile notifier
class ProfileNotifier extends StateNotifier<ProfileState> {
  final ProfileRepository _repository;

  ProfileNotifier(this._repository) : super(const ProfileState());

  /// Load user profile
  Future<void> loadProfile() async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.getProfile();

    result.when(
      success: (user) {
        state = state.copyWith(user: user, isLoading: false);
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Update profile
  Future<bool> updateProfile(UpdateProfileRequest request) async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.updateProfile(request);

    return result.when(
      success: (user) {
        state = state.copyWith(user: user, isLoading: false);
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
        return false;
      },
    );
  }

  /// Upload profile photo
  Future<String?> uploadPhoto(String filePath) async {
    final result = await _repository.uploadProfilePhoto(filePath);

    return result.when(
      success: (url) {
        if (state.user != null) {
          state = state.copyWith(
            user: state.user!.copyWith(avatarUrl: url),
          );
        }
        return url;
      },
      failure: (_, __) => null,
    );
  }

  /// Delete profile photo
  Future<bool> deletePhoto() async {
    final result = await _repository.deleteProfilePhoto();

    return result.when(
      success: (_) {
        if (state.user != null) {
          state = state.copyWith(
            user: state.user!.copyWith(avatarUrl: null),
          );
        }
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Load KYC status
  Future<void> loadKycStatus() async {
    final result = await _repository.getKycStatus();

    result.when(
      success: (status) {
        state = state.copyWith(kycStatus: status);
      },
      failure: (_, __) {},
    );
  }

  /// Submit KYC
  Future<bool> submitKyc({
    required String bvn,
    String? nin,
    String? idType,
    String? idNumber,
    String? idDocumentPath,
    String? proofOfAddressPath,
    String? selfiePhotoPath,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.submitKyc(
      bvn: bvn,
      nin: nin,
      idType: idType,
      idNumber: idNumber,
      idDocumentPath: idDocumentPath,
      proofOfAddressPath: proofOfAddressPath,
      selfiePhotoPath: selfiePhotoPath,
    );

    return result.when(
      success: (status) {
        state = state.copyWith(kycStatus: status, isLoading: false);
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
        return false;
      },
    );
  }

  /// Load settings
  Future<void> loadSettings() async {
    final result = await _repository.getSettings();

    result.when(
      success: (settings) {
        state = state.copyWith(settings: settings);
      },
      failure: (_, __) {},
    );
  }

  /// Update settings
  Future<bool> updateSettings(UserSettings settings) async {
    final result = await _repository.updateSettings(settings);

    return result.when(
      success: (newSettings) {
        state = state.copyWith(settings: newSettings);
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Load referral info
  Future<void> loadReferralInfo() async {
    final result = await _repository.getReferralInfo();

    result.when(
      success: (info) {
        state = state.copyWith(referralInfo: info);
      },
      failure: (_, __) {},
    );
  }

  /// Load all profile data
  Future<void> loadAll() async {
    state = state.copyWith(isLoading: true, error: null);

    await Future.wait([
      loadProfile(),
      loadKycStatus(),
      loadSettings(),
      loadReferralInfo(),
    ]);

    state = state.copyWith(isLoading: false);
  }

  /// Refresh profile
  Future<void> refresh() => loadProfile();
}

/// Main profile provider
final profileProvider = StateNotifierProvider<ProfileNotifier, ProfileState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return ProfileNotifier(repository);
});

/// User only provider
final currentUserProvider = Provider<User?>((ref) {
  return ref.watch(profileProvider).user;
});

/// KYC status provider
final kycStatusProvider = Provider<KycStatus?>((ref) {
  return ref.watch(profileProvider).kycStatus;
});

/// User settings provider
final userSettingsProvider = Provider<UserSettings?>((ref) {
  return ref.watch(profileProvider).settings;
});

/// Referral info provider
final referralInfoProvider = Provider<ReferralInfo?>((ref) {
  return ref.watch(profileProvider).referralInfo;
});

// ==================== PIN MANAGEMENT ====================

/// PIN change state
class PinChangeState {
  final bool isProcessing;
  final bool isSuccess;
  final String? error;

  const PinChangeState({
    this.isProcessing = false,
    this.isSuccess = false,
    this.error,
  });

  PinChangeState copyWith({
    bool? isProcessing,
    bool? isSuccess,
    String? error,
  }) {
    return PinChangeState(
      isProcessing: isProcessing ?? this.isProcessing,
      isSuccess: isSuccess ?? this.isSuccess,
      error: error,
    );
  }
}

/// PIN change notifier
class PinChangeNotifier extends StateNotifier<PinChangeState> {
  final ProfileRepository _repository;

  PinChangeNotifier(this._repository) : super(const PinChangeState());

  /// Change PIN
  Future<bool> changePin({
    required String currentPin,
    required String newPin,
  }) async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.changePin(
      currentPin: currentPin,
      newPin: newPin,
    );

    return result.when(
      success: (_) {
        state = state.copyWith(isProcessing: false, isSuccess: true);
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(isProcessing: false, error: message);
        return false;
      },
    );
  }

  /// Request PIN reset OTP
  Future<bool> requestResetOtp() async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.requestPinResetOtp();

    return result.when(
      success: (_) {
        state = state.copyWith(isProcessing: false);
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(isProcessing: false, error: message);
        return false;
      },
    );
  }

  /// Reset PIN with OTP
  Future<bool> resetPin({
    required String otp,
    required String newPin,
  }) async {
    state = state.copyWith(isProcessing: true, error: null);

    final result = await _repository.resetPin(otp: otp, newPin: newPin);

    return result.when(
      success: (_) {
        state = state.copyWith(isProcessing: false, isSuccess: true);
        return true;
      },
      failure: (message, _) {
        state = state.copyWith(isProcessing: false, error: message);
        return false;
      },
    );
  }

  /// Reset state
  void reset() {
    state = const PinChangeState();
  }
}

/// PIN change provider
final pinChangeProvider =
    StateNotifierProvider<PinChangeNotifier, PinChangeState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return PinChangeNotifier(repository);
});

// ==================== LOGIN ACTIVITY ====================

/// Login activity provider
final loginActivityProvider = FutureProvider.family<List<LoginActivity>, int>((
  ref,
  page,
) async {
  final repository = ref.watch(profileRepositoryProvider);
  final result = await repository.getLoginActivity(page: page);
  return result.data ?? [];
});

// ==================== DEVICE SESSIONS ====================

/// Active sessions state
class SessionsState {
  final List<DeviceSession> sessions;
  final bool isLoading;
  final String? error;

  const SessionsState({
    this.sessions = const [],
    this.isLoading = false,
    this.error,
  });

  SessionsState copyWith({
    List<DeviceSession>? sessions,
    bool? isLoading,
    String? error,
  }) {
    return SessionsState(
      sessions: sessions ?? this.sessions,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// Sessions notifier
class SessionsNotifier extends StateNotifier<SessionsState> {
  final ProfileRepository _repository;

  SessionsNotifier(this._repository) : super(const SessionsState());

  /// Load active sessions
  Future<void> loadSessions() async {
    state = state.copyWith(isLoading: true, error: null);

    final result = await _repository.getActiveSessions();

    result.when(
      success: (sessions) {
        state = state.copyWith(sessions: sessions, isLoading: false);
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Revoke a session
  Future<bool> revokeSession(String sessionId) async {
    final result = await _repository.revokeSession(sessionId);

    return result.when(
      success: (_) {
        state = state.copyWith(
          sessions: state.sessions.where((s) => s.id != sessionId).toList(),
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Revoke all other sessions
  Future<bool> revokeAllOthers() async {
    final result = await _repository.revokeAllOtherSessions();

    return result.when(
      success: (_) {
        state = state.copyWith(
          sessions: state.sessions.where((s) => s.isCurrent).toList(),
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }
}

/// Sessions provider
final sessionsProvider =
    StateNotifierProvider<SessionsNotifier, SessionsState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return SessionsNotifier(repository);
});

// ==================== REFERRALS ====================

/// Referred users provider with pagination
class ReferredUsersState {
  final List<ReferredUser> users;
  final bool isLoading;
  final bool hasMore;
  final int currentPage;
  final String? error;

  const ReferredUsersState({
    this.users = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.error,
  });

  ReferredUsersState copyWith({
    List<ReferredUser>? users,
    bool? isLoading,
    bool? hasMore,
    int? currentPage,
    String? error,
  }) {
    return ReferredUsersState(
      users: users ?? this.users,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      error: error,
    );
  }
}

/// Referred users notifier
class ReferredUsersNotifier extends StateNotifier<ReferredUsersState> {
  final ProfileRepository _repository;

  ReferredUsersNotifier(this._repository) : super(const ReferredUsersState());

  /// Load referred users
  Future<void> loadUsers({bool refresh = false}) async {
    if (state.isLoading) return;
    if (!refresh && !state.hasMore) return;

    state = state.copyWith(
      isLoading: true,
      error: null,
      currentPage: refresh ? 1 : state.currentPage,
    );

    final result = await _repository.getReferredUsers(
      page: refresh ? 1 : state.currentPage,
    );

    result.when(
      success: (users) {
        final newUsers = refresh ? users : [...state.users, ...users];
        state = state.copyWith(
          users: newUsers,
          isLoading: false,
          hasMore: users.length >= 20,
          currentPage: state.currentPage + 1,
        );
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Refresh
  Future<void> refresh() => loadUsers(refresh: true);
}

/// Referred users provider
final referredUsersProvider =
    StateNotifierProvider<ReferredUsersNotifier, ReferredUsersState>((ref) {
  final repository = ref.watch(profileRepositoryProvider);
  return ReferredUsersNotifier(repository);
});
