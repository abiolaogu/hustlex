import '../../api/api_client.dart';
import '../../repositories/base_repository.dart';
import '../../storage/secure_storage.dart';
import '../models/user_model.dart';

class AuthRepository extends BaseRepository {
  final ApiClient _apiClient;
  final SecureStorage _secureStorage;

  AuthRepository({
    required ApiClient apiClient,
    required SecureStorage secureStorage,
  })  : _apiClient = apiClient,
        _secureStorage = secureStorage;

  /// Request OTP for phone number
  Future<Result<void>> requestOtp(String phone) {
    return safeVoidCall(() async {
      await _apiClient.post('/auth/otp/request', data: {'phone': phone});
    });
  }

  /// Verify OTP and get auth tokens
  Future<Result<AuthResponse>> verifyOtp(String phone, String code) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/auth/otp/verify',
        data: {'phone': phone, 'code': code},
      );
      
      final authResponse = AuthResponse.fromJson(response.data['data']);
      
      // Store tokens securely
      await _secureStorage.setAccessToken(authResponse.tokens.accessToken);
      await _secureStorage.setRefreshToken(authResponse.tokens.refreshToken);
      await _secureStorage.setTokenExpiry(authResponse.tokens.expiresAt);
      
      return authResponse;
    });
  }

  /// Register new user
  Future<Result<AuthResponse>> register(RegisterRequest request) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/auth/register',
        data: request.toJson(),
      );
      
      final authResponse = AuthResponse.fromJson(response.data['data']);
      
      // Store tokens securely
      await _secureStorage.setAccessToken(authResponse.tokens.accessToken);
      await _secureStorage.setRefreshToken(authResponse.tokens.refreshToken);
      await _secureStorage.setTokenExpiry(authResponse.tokens.expiresAt);
      
      return authResponse;
    });
  }

  /// Set transaction PIN
  Future<Result<void>> setPin(String pin) {
    return safeVoidCall(() async {
      await _apiClient.post('/auth/pin/set', data: {'pin': pin});
    });
  }

  /// Change transaction PIN
  Future<Result<void>> changePin(String currentPin, String newPin) {
    return safeVoidCall(() async {
      await _apiClient.post('/auth/pin/change', data: {
        'current_pin': currentPin,
        'new_pin': newPin,
      });
    });
  }

  /// Verify transaction PIN
  Future<Result<bool>> verifyPin(String pin) {
    return safeCall(() async {
      final response = await _apiClient.post(
        '/auth/pin/verify',
        data: {'pin': pin},
      );
      return response.data['data']['valid'] as bool;
    });
  }

  /// Get current user profile
  Future<Result<User>> getCurrentUser() {
    return safeCall(() async {
      final response = await _apiClient.get('/users/me');
      return User.fromJson(response.data['data']);
    });
  }

  /// Update user profile
  Future<Result<User>> updateProfile({
    String? firstName,
    String? lastName,
    String? email,
    String? avatar,
  }) {
    return safeCall(() async {
      final data = <String, dynamic>{};
      if (firstName != null) data['first_name'] = firstName;
      if (lastName != null) data['last_name'] = lastName;
      if (email != null) data['email'] = email;
      if (avatar != null) data['avatar'] = avatar;

      final response = await _apiClient.patch('/users/me', data: data);
      return User.fromJson(response.data['data']);
    });
  }

  /// Refresh access token
  Future<Result<AuthTokens>> refreshToken() {
    return safeCall(() async {
      final refreshToken = await _secureStorage.getRefreshToken();
      if (refreshToken == null) {
        throw Exception('No refresh token available');
      }

      final response = await _apiClient.post(
        '/auth/refresh',
        data: {'refresh_token': refreshToken},
      );

      final tokens = AuthTokens.fromJson(response.data['data']);
      
      // Update stored tokens
      await _secureStorage.setAccessToken(tokens.accessToken);
      await _secureStorage.setRefreshToken(tokens.refreshToken);
      await _secureStorage.setTokenExpiry(tokens.expiresAt);

      return tokens;
    });
  }

  /// Logout user
  Future<Result<void>> logout() {
    return safeVoidCall(() async {
      try {
        await _apiClient.post('/auth/logout');
      } catch (_) {
        // Ignore logout API errors
      }
      
      // Clear all stored auth data
      await _secureStorage.clearAll();
    });
  }

  /// Check if user is logged in
  Future<bool> isLoggedIn() async {
    final token = await _secureStorage.getAccessToken();
    return token != null && token.isNotEmpty;
  }

  /// Check if token is expired
  Future<bool> isTokenExpired() async {
    final expiry = await _secureStorage.getTokenExpiry();
    if (expiry == null) return true;
    return DateTime.now().isAfter(expiry);
  }

  /// Get stored access token
  Future<String?> getAccessToken() => _secureStorage.getAccessToken();
}
