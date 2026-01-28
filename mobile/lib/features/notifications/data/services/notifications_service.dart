import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/api/api_client.dart';
import '../../../../core/exceptions/api_exception.dart';

/// Notifications API Service
/// Handles all notification-related API calls
class NotificationsService {
  final ApiClient _apiClient;

  NotificationsService(this._apiClient);

  /// Get all notifications with pagination
  Future<NotificationsResponse> getNotifications({
    int page = 1,
    int perPage = 20,
    bool? unreadOnly,
    String? type,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };
      if (unreadOnly != null) queryParams['unread_only'] = unreadOnly;
      if (type != null) queryParams['type'] = type;

      final response = await _apiClient.get(
        '/api/v1/notifications',
        queryParameters: queryParams,
      );
      return NotificationsResponse.fromJson(response.data);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Get notification count (unread)
  Future<NotificationCount> getNotificationCount() async {
    try {
      final response = await _apiClient.get('/api/v1/notifications/count');
      return NotificationCount.fromJson(response.data['data']);
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Mark notification as read
  Future<void> markAsRead(String notificationId) async {
    try {
      await _apiClient.patch('/api/v1/notifications/$notificationId/read');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Mark all notifications as read
  Future<void> markAllAsRead() async {
    try {
      await _apiClient.post('/api/v1/notifications/read-all');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Delete notification
  Future<void> deleteNotification(String notificationId) async {
    try {
      await _apiClient.delete('/api/v1/notifications/$notificationId');
    } catch (e) {
      throw _handleError(e);
    }
  }

  /// Delete all notifications
  Future<void> deleteAllNotifications() async {
    try {
      await _apiClient.delete('/api/v1/notifications');
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
class NotificationsResponse {
  final List<AppNotification> notifications;
  final PaginationMeta meta;

  NotificationsResponse({required this.notifications, required this.meta});

  factory NotificationsResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as List? ?? [];
    return NotificationsResponse(
      notifications: data.map((n) => AppNotification.fromJson(n)).toList(),
      meta: PaginationMeta.fromJson(json['meta'] ?? {}),
    );
  }
}

class AppNotification {
  final String id;
  final String type;
  final String title;
  final String body;
  final Map<String, dynamic>? data;
  final String? actionUrl;
  final bool isRead;
  final DateTime createdAt;

  AppNotification({
    required this.id,
    required this.type,
    required this.title,
    required this.body,
    this.data,
    this.actionUrl,
    required this.isRead,
    required this.createdAt,
  });

  factory AppNotification.fromJson(Map<String, dynamic> json) {
    return AppNotification(
      id: json['id'] ?? '',
      type: json['type'] ?? 'general',
      title: json['title'] ?? '',
      body: json['body'] ?? '',
      data: json['data'],
      actionUrl: json['action_url'],
      isRead: json['is_read'] ?? false,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
    );
  }

  // Notification type constants
  static const String typeTransaction = 'transaction';
  static const String typeSavings = 'savings';
  static const String typeLoan = 'loan';
  static const String typeGig = 'gig';
  static const String typePromotion = 'promotion';
  static const String typeSecurity = 'security';
  static const String typeGeneral = 'general';
}

class NotificationCount {
  final int total;
  final int unread;
  final Map<String, int>? byType;

  NotificationCount({
    required this.total,
    required this.unread,
    this.byType,
  });

  factory NotificationCount.fromJson(Map<String, dynamic> json) {
    return NotificationCount(
      total: json['total'] ?? 0,
      unread: json['unread'] ?? 0,
      byType: json['by_type'] != null
          ? Map<String, int>.from(json['by_type'])
          : null,
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
final notificationsServiceProvider = Provider<NotificationsService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return NotificationsService(apiClient);
});
