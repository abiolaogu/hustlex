import '../../../core/api/api_client.dart';
import '../../../core/repositories/base_repository.dart';
import '../models/notification_model.dart';

class NotificationRepository extends BaseRepository {
  final ApiClient _apiClient;

  NotificationRepository({required ApiClient apiClient}) : _apiClient = apiClient;

  /// Get paginated notifications with filters
  Future<Result<PaginatedNotifications>> getNotifications(NotificationFilter filter) {
    return safeCall(() async {
      final queryParams = <String, dynamic>{
        'page': filter.page,
        'limit': filter.limit,
      };

      if (filter.type != null) {
        queryParams['type'] = filter.type!.name;
      }
      if (filter.isRead != null) {
        queryParams['is_read'] = filter.isRead;
      }
      if (filter.startDate != null) {
        queryParams['start_date'] = filter.startDate!.toIso8601String();
      }
      if (filter.endDate != null) {
        queryParams['end_date'] = filter.endDate!.toIso8601String();
      }

      final response = await _apiClient.get('/notifications', queryParameters: queryParams);
      return PaginatedNotifications.fromJson(response.data['data']);
    });
  }

  /// Get recent notifications (for bell icon badge)
  Future<Result<List<AppNotification>>> getRecentNotifications({int limit = 10}) {
    return safeCall(() async {
      final response = await _apiClient.get(
        '/notifications',
        queryParameters: {'limit': limit},
      );
      final data = response.data['data']['notifications'] as List;
      return data.map((e) => AppNotification.fromJson(e)).toList();
    });
  }

  /// Get unread notification count
  Future<Result<int>> getUnreadCount() {
    return safeCall(() async {
      final response = await _apiClient.get('/notifications/unread-count');
      return response.data['data']['count'] as int;
    });
  }

  /// Mark a notification as read
  Future<Result<void>> markAsRead(String notificationId) {
    return safeVoidCall(() async {
      await _apiClient.post('/notifications/$notificationId/read');
    });
  }

  /// Mark all notifications as read
  Future<Result<void>> markAllAsRead() {
    return safeVoidCall(() async {
      await _apiClient.post('/notifications/read-all');
    });
  }

  /// Delete a notification
  Future<Result<void>> deleteNotification(String notificationId) {
    return safeVoidCall(() async {
      await _apiClient.delete('/notifications/$notificationId');
    });
  }

  /// Delete all notifications
  Future<Result<void>> deleteAllNotifications() {
    return safeVoidCall(() async {
      await _apiClient.delete('/notifications');
    });
  }

  // Notification Preferences

  /// Get notification preferences
  Future<Result<NotificationPreferences>> getPreferences() {
    return safeCall(() async {
      final response = await _apiClient.get('/notifications/preferences');
      return NotificationPreferences.fromJson(response.data['data']);
    });
  }

  /// Update notification preferences
  Future<Result<NotificationPreferences>> updatePreferences(
    NotificationPreferences preferences,
  ) {
    return safeCall(() async {
      final response = await _apiClient.patch(
        '/notifications/preferences',
        data: preferences.toJson(),
      );
      return NotificationPreferences.fromJson(response.data['data']);
    });
  }

  // Push Notifications

  /// Register FCM token
  Future<Result<void>> registerFcmToken(String token) {
    return safeVoidCall(() async {
      await _apiClient.post('/notifications/fcm/register', data: {
        'token': token,
      });
    });
  }

  /// Unregister FCM token (on logout)
  Future<Result<void>> unregisterFcmToken(String token) {
    return safeVoidCall(() async {
      await _apiClient.post('/notifications/fcm/unregister', data: {
        'token': token,
      });
    });
  }
}
