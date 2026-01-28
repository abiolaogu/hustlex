import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:equatable/equatable.dart';

import '../api/api_client.dart';
import '../constants/app_constants.dart';
import '../exceptions/api_exception.dart';

/// =============================================================================
/// AUTH STATE
/// =============================================================================

class AuthState extends Equatable {
  final bool isAuthenticated;
  final bool isLoading;
  final bool hasCompletedOnboarding;
  final bool hasPinSet;
  final User? user;
  final String? error;

  const AuthState({
    this.isAuthenticated = false,
    this.isLoading = true,
    this.hasCompletedOnboarding = false,
    this.hasPinSet = false,
    this.user,
    this.error,
  });

  AuthState copyWith({
    bool? isAuthenticated,
    bool? isLoading,
    bool? hasCompletedOnboarding,
    bool? hasPinSet,
    User? user,
    String? error,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      isLoading: isLoading ?? this.isLoading,
      hasCompletedOnboarding: hasCompletedOnboarding ?? this.hasCompletedOnboarding,
      hasPinSet: hasPinSet ?? this.hasPinSet,
      user: user ?? this.user,
      error: error,
    );
  }

  @override
  List<Object?> get props => [
        isAuthenticated,
        isLoading,
        hasCompletedOnboarding,
        hasPinSet,
        user,
        error,
      ];
}

/// =============================================================================
/// USER MODEL
/// =============================================================================

class User extends Equatable {
  final String id;
  final String phone;
  final String? email;
  final String firstName;
  final String lastName;
  final String? avatar;
  final bool isVerified;
  final DateTime createdAt;

  const User({
    required this.id,
    required this.phone,
    this.email,
    required this.firstName,
    required this.lastName,
    this.avatar,
    this.isVerified = false,
    required this.createdAt,
  });

  String get fullName => '$firstName $lastName';

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? '',
      phone: json['phone'] ?? '',
      email: json['email'],
      firstName: json['first_name'] ?? '',
      lastName: json['last_name'] ?? '',
      avatar: json['avatar'],
      isVerified: json['is_verified'] ?? false,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'phone': phone,
      'email': email,
      'first_name': firstName,
      'last_name': lastName,
      'avatar': avatar,
      'is_verified': isVerified,
      'created_at': createdAt.toIso8601String(),
    };
  }

  @override
  List<Object?> get props => [
        id,
        phone,
        email,
        firstName,
        lastName,
        avatar,
        isVerified,
        createdAt,
      ];
}

/// =============================================================================
/// AUTH STATE PROVIDER
/// =============================================================================

final authStateProvider = StateNotifierProvider<AuthNotifier, AsyncValue<AuthState>>((ref) {
  return AuthNotifier(ref);
});

class AuthNotifier extends StateNotifier<AsyncValue<AuthState>> {
  final Ref _ref;
  final _storage = const FlutterSecureStorage();

  AuthNotifier(this._ref) : super(const AsyncValue.loading()) {
    _checkAuthStatus();
  }

  ApiClient get _apiClient => _ref.read(apiClientProvider);

  Future<void> _checkAuthStatus() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final hasCompletedOnboarding = prefs.getBool(AppConstants.onboardingKey) ?? false;
      
      final accessToken = await _storage.read(key: AppConstants.accessTokenKey);
      final hasPinSet = prefs.getBool(AppConstants.pinSetKey) ?? false;

      if (accessToken != null && accessToken.isNotEmpty) {
        // Try to fetch user profile
        try {
          final response = await _apiClient.get('/auth/me');
          if (response.success && response.data != null) {
            final user = User.fromJson(response.data);
            state = AsyncValue.data(AuthState(
              isAuthenticated: true,
              isLoading: false,
              hasCompletedOnboarding: hasCompletedOnboarding,
              hasPinSet: hasPinSet,
              user: user,
            ));
            return;
          }
        } catch (e) {
          // Token might be invalid, clear it
          await _storage.delete(key: AppConstants.accessTokenKey);
        }
      }

      state = AsyncValue.data(AuthState(
        isAuthenticated: false,
        isLoading: false,
        hasCompletedOnboarding: hasCompletedOnboarding,
        hasPinSet: hasPinSet,
      ));
    } catch (e) {
      state = AsyncValue.data(AuthState(
        isAuthenticated: false,
        isLoading: false,
        error: e.toString(),
      ));
    }
  }

  Future<void> completeOnboarding() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(AppConstants.onboardingKey, true);
    
    final currentState = state.value ?? const AuthState();
    state = AsyncValue.data(currentState.copyWith(hasCompletedOnboarding: true));
  }

  Future<ApiResponse> requestOtp(String phone) async {
    try {
      final response = await _apiClient.post(
        '/auth/request-otp',
        data: {'phone': phone},
      );
      return response;
    } on ApiException {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> verifyOtp(String phone, String otp) async {
    try {
      final response = await _apiClient.post(
        '/auth/verify-otp',
        data: {'phone': phone, 'otp': otp},
      );

      if (!response.success) {
        throw ApiException(message: response.error ?? 'OTP verification failed');
      }

      final data = response.data as Map<String, dynamic>;
      final isNewUser = data['is_new_user'] ?? true;

      if (!isNewUser) {
        // Existing user - save tokens and set authenticated
        final accessToken = data['access_token'];
        final refreshToken = data['refresh_token'];

        if (accessToken != null) {
          await _storage.write(key: AppConstants.accessTokenKey, value: accessToken);
        }
        if (refreshToken != null) {
          await _storage.write(key: AppConstants.refreshTokenKey, value: refreshToken);
        }

        // Fetch user profile
        final profileResponse = await _apiClient.get('/auth/me');
        if (profileResponse.success && profileResponse.data != null) {
          final user = User.fromJson(profileResponse.data);
          final prefs = await SharedPreferences.getInstance();
          final hasPinSet = prefs.getBool(AppConstants.pinSetKey) ?? false;

          state = AsyncValue.data(AuthState(
            isAuthenticated: true,
            isLoading: false,
            hasCompletedOnboarding: true,
            hasPinSet: hasPinSet,
            user: user,
          ));
        }
      }

      return data;
    } on ApiException {
      rethrow;
    }
  }

  Future<void> register({
    required String phone,
    required String firstName,
    required String lastName,
    String? email,
    String? referralCode,
  }) async {
    try {
      final response = await _apiClient.post(
        '/auth/register',
        data: {
          'phone': phone,
          'first_name': firstName,
          'last_name': lastName,
          if (email != null) 'email': email,
          if (referralCode != null) 'referral_code': referralCode,
        },
      );

      if (!response.success) {
        throw ApiException(message: response.error ?? 'Registration failed');
      }

      final data = response.data as Map<String, dynamic>;
      final accessToken = data['access_token'];
      final refreshToken = data['refresh_token'];

      if (accessToken != null) {
        await _storage.write(key: AppConstants.accessTokenKey, value: accessToken);
      }
      if (refreshToken != null) {
        await _storage.write(key: AppConstants.refreshTokenKey, value: refreshToken);
      }

      final user = User.fromJson(data['user']);
      state = AsyncValue.data(AuthState(
        isAuthenticated: true,
        isLoading: false,
        hasCompletedOnboarding: true,
        hasPinSet: false,
        user: user,
      ));
    } on ApiException {
      rethrow;
    }
  }

  Future<void> setTransactionPin(String pin) async {
    try {
      final response = await _apiClient.post(
        '/auth/pin',
        data: {'pin': pin},
      );

      if (!response.success) {
        throw ApiException(message: response.error ?? 'Failed to set PIN');
      }

      final prefs = await SharedPreferences.getInstance();
      await prefs.setBool(AppConstants.pinSetKey, true);

      final currentState = state.value ?? const AuthState();
      state = AsyncValue.data(currentState.copyWith(hasPinSet: true));
    } on ApiException {
      rethrow;
    }
  }

  Future<bool> verifyTransactionPin(String pin) async {
    try {
      final response = await _apiClient.post(
        '/auth/pin/verify',
        data: {'pin': pin},
      );

      return response.success;
    } on ApiException {
      return false;
    }
  }

  Future<void> logout() async {
    try {
      await _apiClient.post('/auth/logout');
    } catch (_) {
      // Ignore errors during logout
    }

    await _storage.delete(key: AppConstants.accessTokenKey);
    await _storage.delete(key: AppConstants.refreshTokenKey);

    final prefs = await SharedPreferences.getInstance();
    final hasCompletedOnboarding = prefs.getBool(AppConstants.onboardingKey) ?? false;

    state = AsyncValue.data(AuthState(
      isAuthenticated: false,
      isLoading: false,
      hasCompletedOnboarding: hasCompletedOnboarding,
      hasPinSet: false,
    ));
  }

  Future<void> refreshUser() async {
    try {
      final response = await _apiClient.get('/auth/me');
      if (response.success && response.data != null) {
        final user = User.fromJson(response.data);
        final currentState = state.value ?? const AuthState();
        state = AsyncValue.data(currentState.copyWith(user: user));
      }
    } catch (_) {
      // Ignore refresh errors
    }
  }
}

/// =============================================================================
/// CURRENT USER PROVIDER
/// =============================================================================

final currentUserProvider = Provider<User?>((ref) {
  final authState = ref.watch(authStateProvider);
  return authState.value?.user;
});

/// =============================================================================
/// IS AUTHENTICATED PROVIDER
/// =============================================================================

final isAuthenticatedProvider = Provider<bool>((ref) {
  final authState = ref.watch(authStateProvider);
  return authState.value?.isAuthenticated ?? false;
});
