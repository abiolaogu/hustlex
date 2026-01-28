import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/api/api_client.dart';
import '../../../../core/exceptions/api_exception.dart';

/// User/Profile API Service
/// Handles all user profile and account-related API calls
class UserService {
  final ApiClient _apiClient;

  UserService(this._apiClient);

  /// Get current user profile
  Future<UserProfile> getProfile() async {
    try {
      final response = await _apiClient.get('/api/v1/user/profile');
      return UserProfile.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Update user profile
  Future<UserProfile> updateProfile({
    String? firstName,
    String? lastName,
    String? email,
    String? dateOfBirth,
    String? gender,
    String? address,
    String? city,
    String? state,
    String? occupation,
  }) async {
    try {
      final data = <String, dynamic>{};
      if (firstName != null) data['first_name'] = firstName;
      if (lastName != null) data['last_name'] = lastName;
      if (email != null) data['email'] = email;
      if (dateOfBirth != null) data['date_of_birth'] = dateOfBirth;
      if (gender != null) data['gender'] = gender;
      if (address != null) data['address'] = address;
      if (city != null) data['city'] = city;
      if (state != null) data['state'] = state;
      if (occupation != null) data['occupation'] = occupation;

      final response = await _apiClient.patch('/api/v1/user/profile', data: data);
      return UserProfile.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Upload profile photo
  Future<String> uploadProfilePhoto(String filePath) async {
    try {
      final response = await _apiClient.uploadFile(
        '/api/v1/user/profile/photo',
        filePath: filePath,
        field: 'photo',
      );
      return response.data['data']['photo_url'] ?? '';
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Change transaction PIN
  Future<void> changePin({
    required String currentPin,
    required String newPin,
  }) async {
    try {
      await _apiClient.post('/api/v1/user/pin/change', data: {
        'current_pin': currentPin,
        'new_pin': newPin,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Reset PIN (via OTP)
  Future<void> resetPin({
    required String otp,
    required String newPin,
  }) async {
    try {
      await _apiClient.post('/api/v1/user/pin/reset', data: {
        'otp': otp,
        'new_pin': newPin,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Request PIN reset OTP
  Future<void> requestPinResetOtp() async {
    try {
      await _apiClient.post('/api/v1/user/pin/reset-request');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Verify PIN
  Future<bool> verifyPin(String pin) async {
    try {
      final response = await _apiClient.post('/api/v1/user/pin/verify', data: {
        'pin': pin,
      });
      return response.data['valid'] ?? false;
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get user settings
  Future<UserSettings> getSettings() async {
    try {
      final response = await _apiClient.get('/api/v1/user/settings');
      return UserSettings.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Update user settings
  Future<UserSettings> updateSettings({
    bool? pushNotifications,
    bool? emailNotifications,
    bool? smsNotifications,
    bool? biometricEnabled,
    String? language,
    String? currency,
  }) async {
    try {
      final data = <String, dynamic>{};
      if (pushNotifications != null) data['push_notifications'] = pushNotifications;
      if (emailNotifications != null) data['email_notifications'] = emailNotifications;
      if (smsNotifications != null) data['sms_notifications'] = smsNotifications;
      if (biometricEnabled != null) data['biometric_enabled'] = biometricEnabled;
      if (language != null) data['language'] = language;
      if (currency != null) data['currency'] = currency;

      final response = await _apiClient.patch('/api/v1/user/settings', data: data);
      return UserSettings.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get notification preferences
  Future<NotificationPreferences> getNotificationPreferences() async {
    try {
      final response = await _apiClient.get('/api/v1/user/notifications/preferences');
      return NotificationPreferences.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Update notification preferences
  Future<NotificationPreferences> updateNotificationPreferences({
    bool? transactions,
    bool? savingsReminders,
    bool? loanReminders,
    bool? gigUpdates,
    bool? promotions,
    bool? securityAlerts,
  }) async {
    try {
      final data = <String, dynamic>{};
      if (transactions != null) data['transactions'] = transactions;
      if (savingsReminders != null) data['savings_reminders'] = savingsReminders;
      if (loanReminders != null) data['loan_reminders'] = loanReminders;
      if (gigUpdates != null) data['gig_updates'] = gigUpdates;
      if (promotions != null) data['promotions'] = promotions;
      if (securityAlerts != null) data['security_alerts'] = securityAlerts;

      final response = await _apiClient.patch('/api/v1/user/notifications/preferences', data: data);
      return NotificationPreferences.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get KYC status
  Future<KycStatus> getKycStatus() async {
    try {
      final response = await _apiClient.get('/api/v1/user/kyc/status');
      return KycStatus.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Submit KYC documents
  Future<KycSubmission> submitKycDocuments({
    required String documentType, // nin, bvn, drivers_license, passport
    required String documentNumber,
    String? frontImagePath,
    String? backImagePath,
    String? selfieImagePath,
  }) async {
    try {
      final formData = <String, dynamic>{
        'document_type': documentType,
        'document_number': documentNumber,
      };

      // Handle file uploads separately if needed
      final response = await _apiClient.post('/api/v1/user/kyc/submit', data: formData);
      return KycSubmission.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get user activity log
  Future<ActivityLogResponse> getActivityLog({
    int page = 1,
    int perPage = 20,
    String? type,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (type != null) queryParams['type'] = type;

      final response = await _apiClient.get(
        '/api/v1/user/activity',
        queryParameters: queryParams,
      );
      return ActivityLogResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get user devices (active sessions)
  Future<DevicesResponse> getDevices() async {
    try {
      final response = await _apiClient.get('/api/v1/user/devices');
      return DevicesResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Revoke device session
  Future<void> revokeDevice(String deviceId) async {
    try {
      await _apiClient.delete('/api/v1/user/devices/$deviceId');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Delete account
  Future<void> deleteAccount({
    required String reason,
    required String pin,
  }) async {
    try {
      await _apiClient.post('/api/v1/user/delete', data: {
        'reason': reason,
        'pin': pin,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Register FCM token for push notifications
  Future<void> registerFcmToken(String token) async {
    try {
      await _apiClient.post('/api/v1/user/fcm-token', data: {
        'token': token,
      });
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Submit feedback/support ticket
  Future<SupportTicket> submitFeedback({
    required String category,
    required String subject,
    required String message,
    List<String>? attachments,
  }) async {
    try {
      final response = await _apiClient.post('/api/v1/support/tickets', data: {
        'category': category,
        'subject': subject,
        'message': message,
        'attachments': attachments ?? [],
      });
      return SupportTicket.fromJson(response.data['data']);
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
class UserProfile {
  final String id;
  final String phone;
  final String? email;
  final String firstName;
  final String lastName;
  final String? avatar;
  final String? dateOfBirth;
  final String? gender;
  final String? address;
  final String? city;
  final String? state;
  final String? occupation;
  final String accountTier; // basic, standard, premium
  final bool isVerified;
  final bool hasPin;
  final bool hasBiometric;
  final double walletBalance;
  final int creditScore;
  final DateTime createdAt;
  final DateTime updatedAt;

  UserProfile({
    required this.id,
    required this.phone,
    this.email,
    required this.firstName,
    required this.lastName,
    this.avatar,
    this.dateOfBirth,
    this.gender,
    this.address,
    this.city,
    this.state,
    this.occupation,
    required this.accountTier,
    required this.isVerified,
    required this.hasPin,
    required this.hasBiometric,
    required this.walletBalance,
    required this.creditScore,
    required this.createdAt,
    required this.updatedAt,
  });

  String get fullName => '$firstName $lastName';

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      id: json['id'] ?? '',
      phone: json['phone'] ?? '',
      email: json['email'],
      firstName: json['first_name'] ?? '',
      lastName: json['last_name'] ?? '',
      avatar: json['avatar'],
      dateOfBirth: json['date_of_birth'],
      gender: json['gender'],
      address: json['address'],
      city: json['city'],
      state: json['state'],
      occupation: json['occupation'],
      accountTier: json['account_tier'] ?? 'basic',
      isVerified: json['is_verified'] ?? false,
      hasPin: json['has_pin'] ?? false,
      hasBiometric: json['has_biometric'] ?? false,
      walletBalance: (json['wallet_balance'] ?? 0).toDouble(),
      creditScore: json['credit_score'] ?? 0,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
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
      'date_of_birth': dateOfBirth,
      'gender': gender,
      'address': address,
      'city': city,
      'state': state,
      'occupation': occupation,
      'account_tier': accountTier,
      'is_verified': isVerified,
      'has_pin': hasPin,
      'has_biometric': hasBiometric,
      'wallet_balance': walletBalance,
      'credit_score': creditScore,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}

class UserSettings {
  final bool pushNotifications;
  final bool emailNotifications;
  final bool smsNotifications;
  final bool biometricEnabled;
  final String language;
  final String currency;
  final String theme; // light, dark, system

  UserSettings({
    required this.pushNotifications,
    required this.emailNotifications,
    required this.smsNotifications,
    required this.biometricEnabled,
    required this.language,
    required this.currency,
    required this.theme,
  });

  factory UserSettings.fromJson(Map<String, dynamic> json) {
    return UserSettings(
      pushNotifications: json['push_notifications'] ?? true,
      emailNotifications: json['email_notifications'] ?? true,
      smsNotifications: json['sms_notifications'] ?? true,
      biometricEnabled: json['biometric_enabled'] ?? false,
      language: json['language'] ?? 'en',
      currency: json['currency'] ?? 'NGN',
      theme: json['theme'] ?? 'system',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'push_notifications': pushNotifications,
      'email_notifications': emailNotifications,
      'sms_notifications': smsNotifications,
      'biometric_enabled': biometricEnabled,
      'language': language,
      'currency': currency,
      'theme': theme,
    };
  }
}

class NotificationPreferences {
  final bool transactions;
  final bool savingsReminders;
  final bool loanReminders;
  final bool gigUpdates;
  final bool promotions;
  final bool securityAlerts;

  NotificationPreferences({
    required this.transactions,
    required this.savingsReminders,
    required this.loanReminders,
    required this.gigUpdates,
    required this.promotions,
    required this.securityAlerts,
  });

  factory NotificationPreferences.fromJson(Map<String, dynamic> json) {
    return NotificationPreferences(
      transactions: json['transactions'] ?? true,
      savingsReminders: json['savings_reminders'] ?? true,
      loanReminders: json['loan_reminders'] ?? true,
      gigUpdates: json['gig_updates'] ?? true,
      promotions: json['promotions'] ?? false,
      securityAlerts: json['security_alerts'] ?? true,
    );
  }
}

class KycStatus {
  final String level; // none, basic, standard, full
  final String status; // pending, verified, rejected
  final List<KycDocument> documents;
  final List<String> missingDocuments;
  final DateTime? verifiedAt;
  final String? rejectionReason;

  KycStatus({
    required this.level,
    required this.status,
    required this.documents,
    required this.missingDocuments,
    this.verifiedAt,
    this.rejectionReason,
  });

  factory KycStatus.fromJson(Map<String, dynamic> json) {
    return KycStatus(
      level: json['level'] ?? 'none',
      status: json['status'] ?? 'pending',
      documents: (json['documents'] as List? ?? [])
          .map((d) => KycDocument.fromJson(d))
          .toList(),
      missingDocuments: List<String>.from(json['missing_documents'] ?? []),
      verifiedAt: json['verified_at'] != null ? DateTime.tryParse(json['verified_at']) : null,
      rejectionReason: json['rejection_reason'],
    );
  }
}

class KycDocument {
  final String type;
  final String status; // pending, approved, rejected
  final DateTime submittedAt;
  final DateTime? verifiedAt;
  final String? rejectionReason;

  KycDocument({
    required this.type,
    required this.status,
    required this.submittedAt,
    this.verifiedAt,
    this.rejectionReason,
  });

  factory KycDocument.fromJson(Map<String, dynamic> json) {
    return KycDocument(
      type: json['type'] ?? '',
      status: json['status'] ?? 'pending',
      submittedAt: DateTime.tryParse(json['submitted_at'] ?? '') ?? DateTime.now(),
      verifiedAt: json['verified_at'] != null ? DateTime.tryParse(json['verified_at']) : null,
      rejectionReason: json['rejection_reason'],
    );
  }
}

class KycSubmission {
  final String id;
  final String documentType;
  final String status;
  final DateTime submittedAt;
  final String message;

  KycSubmission({
    required this.id,
    required this.documentType,
    required this.status,
    required this.submittedAt,
    required this.message,
  });

  factory KycSubmission.fromJson(Map<String, dynamic> json) {
    return KycSubmission(
      id: json['id'] ?? '',
      documentType: json['document_type'] ?? '',
      status: json['status'] ?? 'pending',
      submittedAt: DateTime.tryParse(json['submitted_at'] ?? '') ?? DateTime.now(),
      message: json['message'] ?? 'Document submitted for verification',
    );
  }
}

class ActivityLogResponse {
  final List<ActivityLogEntry> activities;
  final PaginationMeta meta;

  ActivityLogResponse({required this.activities, required this.meta});

  factory ActivityLogResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return ActivityLogResponse(
      activities: data.map((a) => ActivityLogEntry.fromJson(a)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class ActivityLogEntry {
  final String id;
  final String type; // login, transaction, settings_change, etc.
  final String description;
  final String? ipAddress;
  final String? device;
  final String? location;
  final DateTime createdAt;

  ActivityLogEntry({
    required this.id,
    required this.type,
    required this.description,
    this.ipAddress,
    this.device,
    this.location,
    required this.createdAt,
  });

  factory ActivityLogEntry.fromJson(Map<String, dynamic> json) {
    return ActivityLogEntry(
      id: json['id'] ?? '',
      type: json['type'] ?? '',
      description: json['description'] ?? '',
      ipAddress: json['ip_address'],
      device: json['device'],
      location: json['location'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class DevicesResponse {
  final List<DeviceInfo> devices;

  DevicesResponse({required this.devices});

  factory DevicesResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return DevicesResponse(
      devices: data.map((d) => DeviceInfo.fromJson(d)).toList(),
    );
  }
}

class DeviceInfo {
  final String id;
  final String name;
  final String type; // mobile, web, desktop
  final String? platform; // iOS, Android, Web
  final String? browser;
  final String? ipAddress;
  final String? location;
  final bool isCurrent;
  final DateTime lastActiveAt;
  final DateTime createdAt;

  DeviceInfo({
    required this.id,
    required this.name,
    required this.type,
    this.platform,
    this.browser,
    this.ipAddress,
    this.location,
    required this.isCurrent,
    required this.lastActiveAt,
    required this.createdAt,
  });

  factory DeviceInfo.fromJson(Map<String, dynamic> json) {
    return DeviceInfo(
      id: json['id'] ?? '',
      name: json['name'] ?? 'Unknown Device',
      type: json['type'] ?? 'mobile',
      platform: json['platform'],
      browser: json['browser'],
      ipAddress: json['ip_address'],
      location: json['location'],
      isCurrent: json['is_current'] ?? false,
      lastActiveAt: DateTime.tryParse(json['last_active_at'] ?? '') ?? DateTime.now(),
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }
}

class SupportTicket {
  final String id;
  final String ticketNumber;
  final String category;
  final String subject;
  final String status; // open, in_progress, resolved, closed
  final DateTime createdAt;

  SupportTicket({
    required this.id,
    required this.ticketNumber,
    required this.category,
    required this.subject,
    required this.status,
    required this.createdAt,
  });

  factory SupportTicket.fromJson(Map<String, dynamic> json) {
    return SupportTicket(
      id: json['id'] ?? '',
      ticketNumber: json['ticket_number'] ?? '',
      category: json['category'] ?? '',
      subject: json['subject'] ?? '',
      status: json['status'] ?? 'open',
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
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
final userServiceProvider = Provider<UserService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return UserService(apiClient);
});
