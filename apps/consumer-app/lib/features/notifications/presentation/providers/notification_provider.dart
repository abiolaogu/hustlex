import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/di/providers.dart';
import '../../data/models/notification_model.dart';
import '../../data/repositories/notification_repository.dart';

/// Notifications state
class NotificationsState {
  final List<AppNotification> notifications;
  final int unreadCount;
  final NotificationPreferences? preferences;
  final bool isLoading;
  final bool hasMore;
  final int currentPage;
  final NotificationFilter filter;
  final String? error;

  const NotificationsState({
    this.notifications = const [],
    this.unreadCount = 0,
    this.preferences,
    this.isLoading = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.filter = const NotificationFilter(),
    this.error,
  });

  NotificationsState copyWith({
    List<AppNotification>? notifications,
    int? unreadCount,
    NotificationPreferences? preferences,
    bool? isLoading,
    bool? hasMore,
    int? currentPage,
    NotificationFilter? filter,
    String? error,
  }) {
    return NotificationsState(
      notifications: notifications ?? this.notifications,
      unreadCount: unreadCount ?? this.unreadCount,
      preferences: preferences ?? this.preferences,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      filter: filter ?? this.filter,
      error: error,
    );
  }

  bool get hasUnread => unreadCount > 0;
}

/// Notifications notifier
class NotificationsNotifier extends StateNotifier<NotificationsState> {
  final NotificationRepository _repository;

  NotificationsNotifier(this._repository) : super(const NotificationsState());

  /// Load notifications with pagination
  Future<void> loadNotifications({bool refresh = false}) async {
    if (state.isLoading) return;
    if (!refresh && !state.hasMore) return;

    state = state.copyWith(
      isLoading: true,
      error: null,
      currentPage: refresh ? 1 : state.currentPage,
    );

    final filter = state.filter.copyWith(page: refresh ? 1 : state.currentPage);
    final result = await _repository.getNotifications(filter);

    result.when(
      success: (data) {
        final newNotifications = refresh
            ? data.notifications
            : [...state.notifications, ...data.notifications];
        state = state.copyWith(
          notifications: newNotifications,
          unreadCount: data.unreadCount,
          isLoading: false,
          hasMore: data.hasMore,
          currentPage: data.page + 1,
        );
      },
      failure: (message, _) {
        state = state.copyWith(isLoading: false, error: message);
      },
    );
  }

  /// Refresh notifications
  Future<void> refresh() => loadNotifications(refresh: true);

  /// Load unread count only (for badge)
  Future<void> loadUnreadCount() async {
    final result = await _repository.getUnreadCount();

    result.when(
      success: (count) {
        state = state.copyWith(unreadCount: count);
      },
      failure: (_, __) {},
    );
  }

  /// Mark notification as read
  Future<bool> markAsRead(String notificationId) async {
    final result = await _repository.markAsRead(notificationId);

    return result.when(
      success: (_) {
        state = state.copyWith(
          notifications: state.notifications.map((n) {
            if (n.id == notificationId) {
              return n.copyWith(isRead: true);
            }
            return n;
          }).toList(),
          unreadCount: state.unreadCount > 0 ? state.unreadCount - 1 : 0,
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Mark all as read
  Future<bool> markAllAsRead() async {
    final result = await _repository.markAllAsRead();

    return result.when(
      success: (_) {
        state = state.copyWith(
          notifications: state.notifications.map((n) => n.copyWith(isRead: true)).toList(),
          unreadCount: 0,
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Delete notification
  Future<bool> deleteNotification(String notificationId) async {
    // Find the notification before deleting to check read status
    final notificationToDelete = state.notifications.where((n) => n.id == notificationId).toList();
    final wasUnread = notificationToDelete.isNotEmpty && !notificationToDelete.first.isRead;
    
    final result = await _repository.deleteNotification(notificationId);

    return result.when(
      success: (_) {
        state = state.copyWith(
          notifications: state.notifications.where((n) => n.id != notificationId).toList(),
          unreadCount: wasUnread && state.unreadCount > 0
              ? state.unreadCount - 1
              : state.unreadCount,
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Delete all notifications
  Future<bool> deleteAll() async {
    final result = await _repository.deleteAllNotifications();

    return result.when(
      success: (_) {
        state = state.copyWith(
          notifications: [],
          unreadCount: 0,
          hasMore: false,
        );
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Update filter
  void updateFilter(NotificationFilter filter) {
    state = state.copyWith(filter: filter);
    loadNotifications(refresh: true);
  }

  /// Filter by type
  void filterByType(NotificationType? type) {
    updateFilter(state.filter.copyWith(type: type));
  }

  /// Filter by read status
  void filterByReadStatus(bool? isRead) {
    updateFilter(state.filter.copyWith(isRead: isRead));
  }

  /// Show only unread
  void showUnreadOnly() {
    filterByReadStatus(false);
  }

  /// Show all
  void showAll() {
    filterByReadStatus(null);
  }

  /// Clear filters
  void clearFilters() {
    updateFilter(const NotificationFilter());
  }

  // ==================== PREFERENCES ====================

  /// Load notification preferences
  Future<void> loadPreferences() async {
    final result = await _repository.getPreferences();

    result.when(
      success: (preferences) {
        state = state.copyWith(preferences: preferences);
      },
      failure: (_, __) {},
    );
  }

  /// Update preferences
  Future<bool> updatePreferences(NotificationPreferences preferences) async {
    final result = await _repository.updatePreferences(preferences);

    return result.when(
      success: (newPreferences) {
        state = state.copyWith(preferences: newPreferences);
        return true;
      },
      failure: (_, __) => false,
    );
  }

  /// Toggle push notifications
  Future<bool> togglePush(bool enabled) async {
    if (state.preferences == null) return false;
    return updatePreferences(state.preferences!.copyWith(pushEnabled: enabled));
  }

  /// Toggle email notifications
  Future<bool> toggleEmail(bool enabled) async {
    if (state.preferences == null) return false;
    return updatePreferences(state.preferences!.copyWith(emailEnabled: enabled));
  }

  /// Toggle SMS notifications
  Future<bool> toggleSms(bool enabled) async {
    if (state.preferences == null) return false;
    return updatePreferences(state.preferences!.copyWith(smsEnabled: enabled));
  }

  // ==================== FCM ====================

  /// Register FCM token
  Future<bool> registerFcmToken(String token) async {
    final result = await _repository.registerFcmToken(token);
    return result.isSuccess;
  }

  /// Unregister FCM token
  Future<bool> unregisterFcmToken(String token) async {
    final result = await _repository.unregisterFcmToken(token);
    return result.isSuccess;
  }

  /// Handle incoming notification (add to local state)
  void addNotification(AppNotification notification) {
    state = state.copyWith(
      notifications: [notification, ...state.notifications],
      unreadCount: state.unreadCount + 1,
    );
  }
}

/// Main notifications provider
final notificationsProvider =
    StateNotifierProvider<NotificationsNotifier, NotificationsState>((ref) {
  final repository = ref.watch(notificationRepositoryProvider);
  return NotificationsNotifier(repository);
});

/// Unread count provider (for badge)
final unreadNotificationCountProvider = Provider<int>((ref) {
  return ref.watch(notificationsProvider).unreadCount;
});

/// Has unread provider
final hasUnreadNotificationsProvider = Provider<bool>((ref) {
  return ref.watch(unreadNotificationCountProvider) > 0;
});

/// Recent notifications provider (for dropdown/quick view)
final recentNotificationsProvider = FutureProvider<List<AppNotification>>((ref) async {
  final repository = ref.watch(notificationRepositoryProvider);
  final result = await repository.getRecentNotifications(limit: 5);
  return result.data ?? [];
});

/// Notification preferences provider
final notificationPreferencesProvider = Provider<NotificationPreferences?>((ref) {
  return ref.watch(notificationsProvider).preferences;
});
