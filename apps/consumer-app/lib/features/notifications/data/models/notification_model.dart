import 'package:freezed_annotation/freezed_annotation.dart';

part 'notification_model.freezed.dart';
part 'notification_model.g.dart';

enum NotificationType {
  @JsonValue('gig')
  gig,
  @JsonValue('savings')
  savings,
  @JsonValue('wallet')
  wallet,
  @JsonValue('credit')
  credit,
  @JsonValue('loan')
  loan,
  @JsonValue('system')
  system,
  @JsonValue('promo')
  promo,
}

enum NotificationPriority {
  @JsonValue('low')
  low,
  @JsonValue('normal')
  normal,
  @JsonValue('high')
  high,
  @JsonValue('urgent')
  urgent,
}

extension NotificationTypeX on NotificationType {
  String get displayName {
    switch (this) {
      case NotificationType.gig:
        return 'Gigs';
      case NotificationType.savings:
        return 'Savings';
      case NotificationType.wallet:
        return 'Wallet';
      case NotificationType.credit:
        return 'Credit';
      case NotificationType.loan:
        return 'Loans';
      case NotificationType.system:
        return 'System';
      case NotificationType.promo:
        return 'Promotions';
    }
  }

  String get icon {
    switch (this) {
      case NotificationType.gig:
        return '💼';
      case NotificationType.savings:
        return '🏦';
      case NotificationType.wallet:
        return '💰';
      case NotificationType.credit:
        return '📊';
      case NotificationType.loan:
        return '💵';
      case NotificationType.system:
        return '⚙️';
      case NotificationType.promo:
        return '🎁';
    }
  }
}

@freezed
class AppNotification with _$AppNotification {
  const factory AppNotification({
    required String id,
    required String userId,
    required String title,
    required String body,
    required NotificationType type,
    @Default(NotificationPriority.normal) NotificationPriority priority,
    @Default(false) bool isRead,
    String? actionUrl,
    String? relatedEntityId,
    String? relatedEntityType,
    String? imageUrl,
    Map<String, dynamic>? data,
    DateTime? createdAt,
    DateTime? readAt,
    DateTime? expiresAt,
  }) = _AppNotification;

  const AppNotification._();

  factory AppNotification.fromJson(Map<String, dynamic> json) => _$AppNotificationFromJson(json);

  bool get isExpired {
    if (expiresAt == null) return false;
    return DateTime.now().isAfter(expiresAt!);
  }

  String get timeAgo {
    if (createdAt == null) return '';
    final now = DateTime.now();
    final difference = now.difference(createdAt!);

    if (difference.inDays > 7) {
      return '${createdAt!.day}/${createdAt!.month}/${createdAt!.year}';
    } else if (difference.inDays > 0) {
      return '${difference.inDays}d ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours}h ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes}m ago';
    } else {
      return 'Just now';
    }
  }
}

@freezed
class NotificationPreferences with _$NotificationPreferences {
  const factory NotificationPreferences({
    @Default(true) bool pushEnabled,
    @Default(true) bool emailEnabled,
    @Default(true) bool smsEnabled,
    @Default(true) bool gigNotifications,
    @Default(true) bool savingsNotifications,
    @Default(true) bool walletNotifications,
    @Default(true) bool creditNotifications,
    @Default(true) bool loanNotifications,
    @Default(true) bool systemNotifications,
    @Default(true) bool promoNotifications,
    @Default(true) bool soundEnabled,
    @Default(true) bool vibrationEnabled,
  }) = _NotificationPreferences;

  factory NotificationPreferences.fromJson(Map<String, dynamic> json) => _$NotificationPreferencesFromJson(json);
}

@freezed
class NotificationFilter with _$NotificationFilter {
  const factory NotificationFilter({
    NotificationType? type,
    bool? isRead,
    DateTime? startDate,
    DateTime? endDate,
    @Default(1) int page,
    @Default(20) int limit,
  }) = _NotificationFilter;

  factory NotificationFilter.fromJson(Map<String, dynamic> json) => _$NotificationFilterFromJson(json);
}

@freezed
class PaginatedNotifications with _$PaginatedNotifications {
  const factory PaginatedNotifications({
    required List<AppNotification> notifications,
    required int total,
    required int unreadCount,
    required int page,
    required int limit,
    required bool hasMore,
  }) = _PaginatedNotifications;

  factory PaginatedNotifications.fromJson(Map<String, dynamic> json) => _$PaginatedNotificationsFromJson(json);
}

@freezed
class PushNotificationPayload with _$PushNotificationPayload {
  const factory PushNotificationPayload({
    required String title,
    required String body,
    NotificationType? type,
    String? actionUrl,
    String? imageUrl,
    Map<String, dynamic>? data,
  }) = _PushNotificationPayload;

  factory PushNotificationPayload.fromJson(Map<String, dynamic> json) => _$PushNotificationPayloadFromJson(json);
}
