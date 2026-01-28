import '../../../core/api/api_client.dart';
import '../../../core/repositories/base_repository.dart';
import '../../auth/data/models/user_model.dart';

/// Profile/User repository for account management operations
class ProfileRepository extends BaseRepository {
  final ApiClient _apiClient;

  ProfileRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  // ==================== PROFILE ====================

  /// Get current user's profile
  Future<Result<User>> getProfile() {
    return safeCall(() async {
      final response = await _apiClient.get('/profile');
      return User.fromJson(response.data['data']);
    });
  }

  /// Update user profile
  Future<Result<User>> updateProfile(UpdateProfileRequest request) {
    return safeCall(() async {
      final response = await _apiClient.patch('/profile', data: {
        if (request.fullName != null) 'full_name': request.fullName,
        if (request.email != null) 'email': request.email,
        if (request.phone != null) 'phone': request.phone,
        if (request.dateOfBirth != null)
          'date_of_birth': request.dateOfBirth!.toIso8601String(),
        if (request.address != null) 'address': request.address,
        if (request.state != null) 'state': request.state,
        if (request.city != null) 'city': request.city,
        if (request.bio != null) 'bio': request.bio,
        if (request.occupation != null) 'occupation': request.occupation,
      });
      return User.fromJson(response.data['data']);
    });
  }

  /// Upload profile photo
  Future<Result<String>> uploadProfilePhoto(String filePath) {
    return safeCall(() async {
      final formData = await _createFormData(filePath, 'photo');
      final response = await _apiClient.post(
        '/profile/photo',
        data: formData,
      );
      return response.data['data']['photo_url'] as String;
    });
  }

  /// Delete profile photo
  Future<Result<void>> deleteProfilePhoto() {
    return safeVoidCall(() async {
      await _apiClient.delete('/profile/photo');
    });
  }

  // ==================== PIN MANAGEMENT ====================

  /// Change transaction PIN
  Future<Result<void>> changePin({
    required String currentPin,
    required String newPin,
  }) {
    return safeVoidCall(() async {
      await _apiClient.post('/profile/change-pin', data: {
        'current_pin': currentPin,
        'new_pin': newPin,
      });
    });
  }

  /// Request PIN reset OTP
  Future<Result<void>> requestPinResetOtp() {
    return safeVoidCall(() async {
      await _apiClient.post('/profile/reset-pin/request');
    });
  }

  /// Reset PIN with OTP
  Future<Result<void>> resetPin({
    required String otp,
    required String newPin,
  }) {
    return safeVoidCall(() async {
      await _apiClient.post('/profile/reset-pin', data: {
        'otp': otp,
        'new_pin': newPin,
      });
    });
  }

  /// Verify PIN (for sensitive operations)
  Future<Result<bool>> verifyPin(String pin) {
    return safeCall(() async {
      final response = await _apiClient.post('/profile/verify-pin', data: {
        'pin': pin,
      });
      return response.data['data']['valid'] as bool;
    });
  }

  // ==================== KYC ====================

  /// Get KYC status
  Future<Result<KycStatus>> getKycStatus() {
    return safeCall(() async {
      final response = await _apiClient.get('/profile/kyc');
      return KycStatus.fromJson(response.data['data']);
    });
  }

  /// Submit KYC documents
  Future<Result<KycStatus>> submitKyc({
    required String bvn,
    String? nin,
    String? idType,
    String? idNumber,
    String? idDocumentPath,
    String? proofOfAddressPath,
    String? selfiePhotoPath,
  }) {
    return safeCall(() async {
      final data = <String, dynamic>{
        'bvn': bvn,
        if (nin != null) 'nin': nin,
        if (idType != null) 'id_type': idType,
        if (idNumber != null) 'id_number': idNumber,
      };

      // Handle file uploads separately if provided
      if (idDocumentPath != null ||
          proofOfAddressPath != null ||
          selfiePhotoPath != null) {
        final formData = await _createMultipartData(
          data,
          {
            if (idDocumentPath != null) 'id_document': idDocumentPath,
            if (proofOfAddressPath != null) 'proof_of_address': proofOfAddressPath,
            if (selfiePhotoPath != null) 'selfie_photo': selfiePhotoPath,
          },
        );
        final response = await _apiClient.post('/profile/kyc', data: formData);
        return KycStatus.fromJson(response.data['data']);
      }

      final response = await _apiClient.post('/profile/kyc', data: data);
      return KycStatus.fromJson(response.data['data']);
    });
  }

  // ==================== SETTINGS ====================

  /// Get user settings
  Future<Result<UserSettings>> getSettings() {
    return safeCall(() async {
      final response = await _apiClient.get('/profile/settings');
      return UserSettings.fromJson(response.data['data']);
    });
  }

  /// Update user settings
  Future<Result<UserSettings>> updateSettings(UserSettings settings) {
    return safeCall(() async {
      final response = await _apiClient.patch(
        '/profile/settings',
        data: settings.toJson(),
      );
      return UserSettings.fromJson(response.data['data']);
    });
  }

  // ==================== SECURITY ====================

  /// Enable/disable biometric authentication
  Future<Result<void>> setBiometricEnabled(bool enabled) {
    return safeVoidCall(() async {
      await _apiClient.post('/profile/biometric', data: {
        'enabled': enabled,
      });
    });
  }

  /// Get login activity/history
  Future<Result<List<LoginActivity>>> getLoginActivity({
    int page = 1,
    int limit = 20,
  }) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/profile/activity',
        queryParameters: {'page': page, 'limit': limit},
      );
      final data = response.data['data']['activities'] as List;
      return data.map((e) => LoginActivity.fromJson(e)).toList();
    });
  }

  /// Get active devices/sessions
  Future<Result<List<DeviceSession>>> getActiveSessions() {
    return safeCall(() async {
      final response = await _apiClient.get('/profile/sessions');
      final data = response.data['data'] as List;
      return data.map((e) => DeviceSession.fromJson(e)).toList();
    });
  }

  /// Revoke a session
  Future<Result<void>> revokeSession(String sessionId) {
    return safeVoidCall(() async {
      await _apiClient.delete('/profile/sessions/$sessionId');
    });
  }

  /// Revoke all other sessions
  Future<Result<void>> revokeAllOtherSessions() {
    return safeVoidCall(() async {
      await _apiClient.delete('/profile/sessions');
    });
  }

  // ==================== ACCOUNT ====================

  /// Request account deletion
  Future<Result<void>> requestAccountDeletion({String? reason}) {
    return safeVoidCall(() async {
      await _apiClient.post('/profile/delete-request', data: {
        if (reason != null) 'reason': reason,
      });
    });
  }

  /// Cancel account deletion request
  Future<Result<void>> cancelAccountDeletion() {
    return safeVoidCall(() async {
      await _apiClient.delete('/profile/delete-request');
    });
  }

  /// Submit feedback
  Future<Result<void>> submitFeedback({
    required String type,
    required String message,
    String? email,
    List<String>? attachmentPaths,
  }) {
    return safeVoidCall(() async {
      if (attachmentPaths != null && attachmentPaths.isNotEmpty) {
        final formData = await _createMultipartData(
          {'type': type, 'message': message, if (email != null) 'email': email},
          {for (int i = 0; i < attachmentPaths.length; i++) 'attachment_$i': attachmentPaths[i]},
        );
        await _apiClient.post('/feedback', data: formData);
      } else {
        await _apiClient.post('/feedback', data: {
          'type': type,
          'message': message,
          if (email != null) 'email': email,
        });
      }
    });
  }

  // ==================== REFERRALS ====================

  /// Get referral info and stats
  Future<Result<ReferralInfo>> getReferralInfo() {
    return safeCall(() async {
      final response = await _apiClient.get('/profile/referral');
      return ReferralInfo.fromJson(response.data['data']);
    });
  }

  /// Get referred users
  Future<Result<List<ReferredUser>>> getReferredUsers({
    int page = 1,
    int limit = 20,
  }) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/profile/referral/users',
        queryParameters: {'page': page, 'limit': limit},
      );
      final data = response.data['data']['users'] as List;
      return data.map((e) => ReferredUser.fromJson(e)).toList();
    });
  }

  // ==================== HELPERS ====================

  /// Create FormData for single file upload
  Future<dynamic> _createFormData(String filePath, String fieldName) async {
    // This would use dio's FormData in actual implementation
    // For now, return a map structure
    return {
      fieldName: filePath,
    };
  }

  /// Create FormData for multiple file uploads
  Future<dynamic> _createMultipartData(
    Map<String, dynamic> fields,
    Map<String, String> files,
  ) async {
    // This would use dio's FormData in actual implementation
    return {
      ...fields,
      ...files,
    };
  }
}

// ==================== REQUEST MODELS ====================

class UpdateProfileRequest {
  final String? fullName;
  final String? email;
  final String? phone;
  final DateTime? dateOfBirth;
  final String? address;
  final String? state;
  final String? city;
  final String? bio;
  final String? occupation;

  const UpdateProfileRequest({
    this.fullName,
    this.email,
    this.phone,
    this.dateOfBirth,
    this.address,
    this.state,
    this.city,
    this.bio,
    this.occupation,
  });
}

// ==================== RESPONSE MODELS ====================

class KycStatus {
  final String level; // 'none', 'basic', 'intermediate', 'full'
  final bool bvnVerified;
  final bool ninVerified;
  final bool idVerified;
  final bool addressVerified;
  final bool selfieVerified;
  final DateTime? verifiedAt;
  final String? rejectionReason;
  final Map<String, double> limits;

  const KycStatus({
    required this.level,
    this.bvnVerified = false,
    this.ninVerified = false,
    this.idVerified = false,
    this.addressVerified = false,
    this.selfieVerified = false,
    this.verifiedAt,
    this.rejectionReason,
    this.limits = const {},
  });

  factory KycStatus.fromJson(Map<String, dynamic> json) {
    return KycStatus(
      level: json['level'] as String,
      bvnVerified: json['bvn_verified'] as bool? ?? false,
      ninVerified: json['nin_verified'] as bool? ?? false,
      idVerified: json['id_verified'] as bool? ?? false,
      addressVerified: json['address_verified'] as bool? ?? false,
      selfieVerified: json['selfie_verified'] as bool? ?? false,
      verifiedAt: json['verified_at'] != null
          ? DateTime.parse(json['verified_at'] as String)
          : null,
      rejectionReason: json['rejection_reason'] as String?,
      limits: (json['limits'] as Map<String, dynamic>?)?.map(
            (k, v) => MapEntry(k, (v as num).toDouble()),
          ) ??
          {},
    );
  }

  bool get isVerified => level != 'none';
  bool get isFullyVerified => level == 'full';
}

class UserSettings {
  final bool pushNotifications;
  final bool emailNotifications;
  final bool smsNotifications;
  final bool transactionAlerts;
  final bool savingsReminders;
  final bool marketingMessages;
  final bool biometricEnabled;
  final String language;
  final String currency;
  final String theme; // 'light', 'dark', 'system'

  const UserSettings({
    this.pushNotifications = true,
    this.emailNotifications = true,
    this.smsNotifications = true,
    this.transactionAlerts = true,
    this.savingsReminders = true,
    this.marketingMessages = false,
    this.biometricEnabled = false,
    this.language = 'en',
    this.currency = 'NGN',
    this.theme = 'system',
  });

  factory UserSettings.fromJson(Map<String, dynamic> json) {
    return UserSettings(
      pushNotifications: json['push_notifications'] as bool? ?? true,
      emailNotifications: json['email_notifications'] as bool? ?? true,
      smsNotifications: json['sms_notifications'] as bool? ?? true,
      transactionAlerts: json['transaction_alerts'] as bool? ?? true,
      savingsReminders: json['savings_reminders'] as bool? ?? true,
      marketingMessages: json['marketing_messages'] as bool? ?? false,
      biometricEnabled: json['biometric_enabled'] as bool? ?? false,
      language: json['language'] as String? ?? 'en',
      currency: json['currency'] as String? ?? 'NGN',
      theme: json['theme'] as String? ?? 'system',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'push_notifications': pushNotifications,
      'email_notifications': emailNotifications,
      'sms_notifications': smsNotifications,
      'transaction_alerts': transactionAlerts,
      'savings_reminders': savingsReminders,
      'marketing_messages': marketingMessages,
      'biometric_enabled': biometricEnabled,
      'language': language,
      'currency': currency,
      'theme': theme,
    };
  }

  UserSettings copyWith({
    bool? pushNotifications,
    bool? emailNotifications,
    bool? smsNotifications,
    bool? transactionAlerts,
    bool? savingsReminders,
    bool? marketingMessages,
    bool? biometricEnabled,
    String? language,
    String? currency,
    String? theme,
  }) {
    return UserSettings(
      pushNotifications: pushNotifications ?? this.pushNotifications,
      emailNotifications: emailNotifications ?? this.emailNotifications,
      smsNotifications: smsNotifications ?? this.smsNotifications,
      transactionAlerts: transactionAlerts ?? this.transactionAlerts,
      savingsReminders: savingsReminders ?? this.savingsReminders,
      marketingMessages: marketingMessages ?? this.marketingMessages,
      biometricEnabled: biometricEnabled ?? this.biometricEnabled,
      language: language ?? this.language,
      currency: currency ?? this.currency,
      theme: theme ?? this.theme,
    );
  }
}

class LoginActivity {
  final String id;
  final String action; // 'login', 'logout', 'password_change', etc.
  final String device;
  final String browser;
  final String os;
  final String ipAddress;
  final String? location;
  final DateTime timestamp;
  final bool success;

  const LoginActivity({
    required this.id,
    required this.action,
    required this.device,
    required this.browser,
    required this.os,
    required this.ipAddress,
    this.location,
    required this.timestamp,
    this.success = true,
  });

  factory LoginActivity.fromJson(Map<String, dynamic> json) {
    return LoginActivity(
      id: json['id'] as String,
      action: json['action'] as String,
      device: json['device'] as String? ?? 'Unknown',
      browser: json['browser'] as String? ?? 'Unknown',
      os: json['os'] as String? ?? 'Unknown',
      ipAddress: json['ip_address'] as String,
      location: json['location'] as String?,
      timestamp: DateTime.parse(json['timestamp'] as String),
      success: json['success'] as bool? ?? true,
    );
  }
}

class DeviceSession {
  final String id;
  final String deviceName;
  final String deviceType; // 'mobile', 'tablet', 'desktop'
  final String os;
  final String browser;
  final String ipAddress;
  final String? location;
  final DateTime lastActive;
  final bool isCurrent;

  const DeviceSession({
    required this.id,
    required this.deviceName,
    required this.deviceType,
    required this.os,
    required this.browser,
    required this.ipAddress,
    this.location,
    required this.lastActive,
    this.isCurrent = false,
  });

  factory DeviceSession.fromJson(Map<String, dynamic> json) {
    return DeviceSession(
      id: json['id'] as String,
      deviceName: json['device_name'] as String,
      deviceType: json['device_type'] as String? ?? 'mobile',
      os: json['os'] as String? ?? 'Unknown',
      browser: json['browser'] as String? ?? 'Unknown',
      ipAddress: json['ip_address'] as String,
      location: json['location'] as String?,
      lastActive: DateTime.parse(json['last_active'] as String),
      isCurrent: json['is_current'] as bool? ?? false,
    );
  }
}

class ReferralInfo {
  final String referralCode;
  final String referralLink;
  final int totalReferrals;
  final int successfulReferrals;
  final double totalEarnings;
  final double pendingEarnings;
  final double referralBonus; // per successful referral

  const ReferralInfo({
    required this.referralCode,
    required this.referralLink,
    this.totalReferrals = 0,
    this.successfulReferrals = 0,
    this.totalEarnings = 0,
    this.pendingEarnings = 0,
    this.referralBonus = 0,
  });

  factory ReferralInfo.fromJson(Map<String, dynamic> json) {
    return ReferralInfo(
      referralCode: json['referral_code'] as String,
      referralLink: json['referral_link'] as String,
      totalReferrals: json['total_referrals'] as int? ?? 0,
      successfulReferrals: json['successful_referrals'] as int? ?? 0,
      totalEarnings: (json['total_earnings'] as num?)?.toDouble() ?? 0,
      pendingEarnings: (json['pending_earnings'] as num?)?.toDouble() ?? 0,
      referralBonus: (json['referral_bonus'] as num?)?.toDouble() ?? 0,
    );
  }
}

class ReferredUser {
  final String id;
  final String name;
  final String? avatarUrl;
  final DateTime joinedAt;
  final String status; // 'pending', 'active', 'completed'
  final double? bonus;

  const ReferredUser({
    required this.id,
    required this.name,
    this.avatarUrl,
    required this.joinedAt,
    required this.status,
    this.bonus,
  });

  factory ReferredUser.fromJson(Map<String, dynamic> json) {
    return ReferredUser(
      id: json['id'] as String,
      name: json['name'] as String,
      avatarUrl: json['avatar_url'] as String?,
      joinedAt: DateTime.parse(json['joined_at'] as String),
      status: json['status'] as String,
      bonus: (json['bonus'] as num?)?.toDouble(),
    );
  }

  bool get isPending => status == 'pending';
  bool get isActive => status == 'active';
  bool get isCompleted => status == 'completed';
}
