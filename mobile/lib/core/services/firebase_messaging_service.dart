import 'dart:convert';
import 'dart:io';

import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../features/notifications/data/models/notification_model.dart';

/// Background message handler (must be top-level function)
@pragma('vm:entry-point')
Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  await Firebase.initializeApp();
  // Handle background message
  debugPrint('Handling background message: ${message.messageId}');
}

/// Local notification channel for Android
const AndroidNotificationChannel _channel = AndroidNotificationChannel(
  'hustlex_high_importance_channel',
  'HustleX Notifications',
  description: 'This channel is used for important HustleX notifications.',
  importance: Importance.high,
);

/// Firebase messaging service
class FirebaseMessagingService {
  final FirebaseMessaging _messaging;
  final FlutterLocalNotificationsPlugin _localNotifications;
  
  String? _fcmToken;
  
  // Callbacks
  void Function(String token)? onTokenRefresh;
  void Function(RemoteMessage message)? onMessage;
  void Function(RemoteMessage message)? onMessageOpenedApp;
  void Function(PushNotificationPayload payload)? onNotificationTapped;

  FirebaseMessagingService({
    FirebaseMessaging? messaging,
    FlutterLocalNotificationsPlugin? localNotifications,
  })  : _messaging = messaging ?? FirebaseMessaging.instance,
        _localNotifications = localNotifications ?? FlutterLocalNotificationsPlugin();

  /// Get current FCM token
  String? get fcmToken => _fcmToken;

  /// Initialize the service
  Future<void> initialize() async {
    // Request permission
    await _requestPermission();

    // Initialize local notifications
    await _initializeLocalNotifications();

    // Get initial token
    _fcmToken = await _messaging.getToken();
    debugPrint('FCM Token: $_fcmToken');

    // Listen for token refreshes
    _messaging.onTokenRefresh.listen((token) {
      _fcmToken = token;
      onTokenRefresh?.call(token);
    });

    // Handle foreground messages
    FirebaseMessaging.onMessage.listen(_handleForegroundMessage);

    // Handle background messages
    FirebaseMessaging.onBackgroundMessage(firebaseMessagingBackgroundHandler);

    // Handle notification tap when app was terminated
    final initialMessage = await _messaging.getInitialMessage();
    if (initialMessage != null) {
      _handleNotificationTap(initialMessage);
    }

    // Handle notification tap when app was in background
    FirebaseMessaging.onMessageOpenedApp.listen(_handleNotificationTap);
  }

  /// Request notification permission
  Future<bool> _requestPermission() async {
    final settings = await _messaging.requestPermission(
      alert: true,
      announcement: false,
      badge: true,
      carPlay: false,
      criticalAlert: false,
      provisional: false,
      sound: true,
    );

    debugPrint('Notification permission: ${settings.authorizationStatus}');
    return settings.authorizationStatus == AuthorizationStatus.authorized ||
        settings.authorizationStatus == AuthorizationStatus.provisional;
  }

  /// Initialize local notifications for foreground display
  Future<void> _initializeLocalNotifications() async {
    // Create high importance channel on Android
    if (Platform.isAndroid) {
      await _localNotifications
          .resolvePlatformSpecificImplementation<
              AndroidFlutterLocalNotificationsPlugin>()
          ?.createNotificationChannel(_channel);
    }

    // Initialize plugin
    const androidSettings = AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings(
      requestAlertPermission: false,
      requestBadgePermission: false,
      requestSoundPermission: false,
    );
    
    await _localNotifications.initialize(
      const InitializationSettings(
        android: androidSettings,
        iOS: iosSettings,
      ),
      onDidReceiveNotificationResponse: _onNotificationResponse,
    );
  }

  /// Handle foreground message
  void _handleForegroundMessage(RemoteMessage message) {
    debugPrint('Foreground message received: ${message.messageId}');
    
    onMessage?.call(message);

    // Show local notification
    final notification = message.notification;
    final android = message.notification?.android;

    if (notification != null) {
      _localNotifications.show(
        notification.hashCode,
        notification.title,
        notification.body,
        NotificationDetails(
          android: AndroidNotificationDetails(
            _channel.id,
            _channel.name,
            channelDescription: _channel.description,
            icon: android?.smallIcon ?? '@mipmap/ic_launcher',
            importance: Importance.high,
            priority: Priority.high,
          ),
          iOS: const DarwinNotificationDetails(
            presentAlert: true,
            presentBadge: true,
            presentSound: true,
          ),
        ),
        payload: jsonEncode(message.data),
      );
    }
  }

  /// Handle notification tap (background/terminated state)
  void _handleNotificationTap(RemoteMessage message) {
    debugPrint('Notification tapped: ${message.messageId}');
    
    onMessageOpenedApp?.call(message);

    final payload = _parsePayload(message);
    if (payload != null) {
      onNotificationTapped?.call(payload);
    }
  }

  /// Handle local notification response
  void _onNotificationResponse(NotificationResponse response) {
    debugPrint('Local notification tapped: ${response.payload}');
    
    if (response.payload != null) {
      try {
        final data = jsonDecode(response.payload!) as Map<String, dynamic>;
        final payload = PushNotificationPayload.fromJson(data);
        onNotificationTapped?.call(payload);
      } catch (e) {
        debugPrint('Error parsing notification payload: $e');
      }
    }
  }

  /// Parse remote message to payload
  PushNotificationPayload? _parsePayload(RemoteMessage message) {
    try {
      return PushNotificationPayload(
        title: message.notification?.title ?? '',
        body: message.notification?.body ?? '',
        type: _parseNotificationType(message.data['type']),
        actionUrl: message.data['action_url'],
        imageUrl: message.notification?.android?.imageUrl ?? 
                  message.notification?.apple?.imageUrl,
        data: message.data,
      );
    } catch (e) {
      debugPrint('Error parsing notification payload: $e');
      return null;
    }
  }

  /// Parse notification type from string
  NotificationType? _parseNotificationType(String? type) {
    if (type == null) return null;
    try {
      return NotificationType.values.firstWhere(
        (t) => t.name == type,
        orElse: () => NotificationType.system,
      );
    } catch (e) {
      return NotificationType.system;
    }
  }

  /// Subscribe to a topic
  Future<void> subscribeToTopic(String topic) async {
    await _messaging.subscribeToTopic(topic);
    debugPrint('Subscribed to topic: $topic');
  }

  /// Unsubscribe from a topic
  Future<void> unsubscribeFromTopic(String topic) async {
    await _messaging.unsubscribeFromTopic(topic);
    debugPrint('Unsubscribed from topic: $topic');
  }

  /// Get APNs token (iOS only)
  Future<String?> getApnsToken() async {
    if (!Platform.isIOS) return null;
    return await _messaging.getAPNSToken();
  }

  /// Delete FCM token (for logout)
  Future<void> deleteToken() async {
    await _messaging.deleteToken();
    _fcmToken = null;
  }

  /// Update notification badge count (iOS)
  Future<void> setBadgeCount(int count) async {
    if (Platform.isIOS) {
      // iOS badge is handled by the system
    }
  }

  /// Clear all notifications
  Future<void> clearAllNotifications() async {
    await _localNotifications.cancelAll();
  }
}

/// Debug print helper (can be replaced with proper logging)
void debugPrint(String message) {
  // ignore: avoid_print
  print('[FirebaseMessaging] $message');
}

/// Firebase messaging service provider
final firebaseMessagingServiceProvider = Provider<FirebaseMessagingService>((ref) {
  return FirebaseMessagingService();
});

/// FCM token provider
final fcmTokenProvider = FutureProvider<String?>((ref) async {
  final service = ref.watch(firebaseMessagingServiceProvider);
  await service.initialize();
  return service.fcmToken;
});
