import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:logger/logger.dart';

import '../storage/secure_storage.dart';

/// Background message handler - must be top-level function
@pragma('vm:entry-point')
Future<void> _firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  await Firebase.initializeApp();
  await FirebaseService._handleBackgroundMessage(message);
}

/// Firebase service for push notifications
class FirebaseService {
  final Logger _logger = Logger();
  final FlutterLocalNotificationsPlugin _localNotifications;
  final SecureStorage _secureStorage;

  static final StreamController<RemoteMessage> _messageController =
      StreamController<RemoteMessage>.broadcast();

  /// Stream of incoming push notifications
  static Stream<RemoteMessage> get onMessage => _messageController.stream;

  /// Stream controller for notification taps
  static final StreamController<Map<String, dynamic>> _notificationTapController =
      StreamController<Map<String, dynamic>>.broadcast();

  /// Stream of notification tap events
  static Stream<Map<String, dynamic>> get onNotificationTap =>
      _notificationTapController.stream;

  FirebaseService({
    required SecureStorage secureStorage,
    FlutterLocalNotificationsPlugin? localNotifications,
  })  : _secureStorage = secureStorage,
        _localNotifications =
            localNotifications ?? FlutterLocalNotificationsPlugin();

  /// Initialize Firebase and messaging
  Future<void> initialize() async {
    try {
      // Initialize Firebase
      await Firebase.initializeApp();
      _logger.i('Firebase initialized');

      // Initialize local notifications
      await _initializeLocalNotifications();

      // Set up Firebase Messaging
      await _setupFirebaseMessaging();

      _logger.i('Firebase service fully initialized');
    } catch (e) {
      _logger.e('Failed to initialize Firebase', error: e);
    }
  }

  /// Initialize local notifications plugin
  Future<void> _initializeLocalNotifications() async {
    const androidSettings = AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings(
      requestAlertPermission: true,
      requestBadgePermission: true,
      requestSoundPermission: true,
    );

    const settings = InitializationSettings(
      android: androidSettings,
      iOS: iosSettings,
    );

    await _localNotifications.initialize(
      settings,
      onDidReceiveNotificationResponse: _onNotificationTap,
      onDidReceiveBackgroundNotificationResponse: _onBackgroundNotificationTap,
    );

    // Create notification channel for Android
    if (Platform.isAndroid) {
      await _createNotificationChannels();
    }
  }

  /// Create Android notification channels
  Future<void> _createNotificationChannels() async {
    final androidPlugin =
        _localNotifications.resolvePlatformSpecificImplementation<
            AndroidFlutterLocalNotificationsPlugin>();

    if (androidPlugin == null) return;

    // Transaction channel
    await androidPlugin.createNotificationChannel(
      const AndroidNotificationChannel(
        'transactions',
        'Transactions',
        description: 'Notifications for wallet transactions',
        importance: Importance.high,
        playSound: true,
      ),
    );

    // Savings channel
    await androidPlugin.createNotificationChannel(
      const AndroidNotificationChannel(
        'savings',
        'Savings',
        description: 'Notifications for Ajo/Esusu circles',
        importance: Importance.high,
        playSound: true,
      ),
    );

    // Gigs channel
    await androidPlugin.createNotificationChannel(
      const AndroidNotificationChannel(
        'gigs',
        'Gigs',
        description: 'Notifications for gig opportunities',
        importance: Importance.high,
        playSound: true,
      ),
    );

    // Loans channel
    await androidPlugin.createNotificationChannel(
      const AndroidNotificationChannel(
        'loans',
        'Loans',
        description: 'Notifications for loan updates',
        importance: Importance.high,
        playSound: true,
      ),
    );

    // Promotions channel
    await androidPlugin.createNotificationChannel(
      const AndroidNotificationChannel(
        'promotions',
        'Promotions',
        description: 'Promotional notifications',
        importance: Importance.defaultImportance,
        playSound: false,
      ),
    );

    // General channel
    await androidPlugin.createNotificationChannel(
      const AndroidNotificationChannel(
        'general',
        'General',
        description: 'General notifications',
        importance: Importance.defaultImportance,
        playSound: true,
      ),
    );
  }

  /// Set up Firebase Messaging
  Future<void> _setupFirebaseMessaging() async {
    final messaging = FirebaseMessaging.instance;

    // Request permission
    final settings = await messaging.requestPermission(
      alert: true,
      announcement: false,
      badge: true,
      carPlay: false,
      criticalAlert: false,
      provisional: false,
      sound: true,
    );

    _logger.i('Notification permission: ${settings.authorizationStatus}');

    if (settings.authorizationStatus == AuthorizationStatus.authorized ||
        settings.authorizationStatus == AuthorizationStatus.provisional) {
      // Get FCM token
      await _getFCMToken();

      // Listen for token refresh
      messaging.onTokenRefresh.listen(_onTokenRefresh);

      // Set up background handler
      FirebaseMessaging.onBackgroundMessage(
        _firebaseMessagingBackgroundHandler,
      );

      // Handle foreground messages
      FirebaseMessaging.onMessage.listen(_handleForegroundMessage);

      // Handle notification tap when app is in background/terminated
      FirebaseMessaging.onMessageOpenedApp.listen(_handleMessageOpenedApp);

      // Check if app was opened from notification
      final initialMessage = await messaging.getInitialMessage();
      if (initialMessage != null) {
        _handleMessageOpenedApp(initialMessage);
      }
    }
  }

  /// Get FCM token
  Future<String?> _getFCMToken() async {
    try {
      final token = await FirebaseMessaging.instance.getToken();
      if (token != null) {
        await _secureStorage.write(key: 'fcm_token', value: token);
        _logger.i('FCM token obtained: ${token.substring(0, 20)}...');
      }
      return token;
    } catch (e) {
      _logger.e('Failed to get FCM token', error: e);
      return null;
    }
  }

  /// Handle token refresh
  Future<void> _onTokenRefresh(String token) async {
    _logger.i('FCM token refreshed');
    await _secureStorage.write(key: 'fcm_token', value: token);
    // TODO: Send updated token to backend
  }

  /// Get current FCM token
  Future<String?> getFCMToken() async {
    return await _secureStorage.read(key: 'fcm_token');
  }

  /// Handle foreground message
  Future<void> _handleForegroundMessage(RemoteMessage message) async {
    _logger.i('Foreground message received: ${message.messageId}');
    _messageController.add(message);

    // Show local notification
    await _showLocalNotification(message);
  }

  /// Handle background message (static method for top-level handler)
  static Future<void> _handleBackgroundMessage(RemoteMessage message) async {
    Logger().i('Background message received: ${message.messageId}');
    _messageController.add(message);
  }

  /// Handle message opened app (notification tap)
  void _handleMessageOpenedApp(RemoteMessage message) {
    _logger.i('Notification tapped: ${message.messageId}');
    final data = message.data;
    _notificationTapController.add(data);
  }

  /// Show local notification
  Future<void> _showLocalNotification(RemoteMessage message) async {
    final notification = message.notification;
    if (notification == null) return;

    final channelId = _getChannelId(message.data['type'] ?? 'general');

    final androidDetails = AndroidNotificationDetails(
      channelId,
      channelId,
      importance: Importance.high,
      priority: Priority.high,
      showWhen: true,
      icon: '@mipmap/ic_launcher',
      largeIcon: const DrawableResourceAndroidBitmap('@mipmap/ic_launcher'),
    );

    const iosDetails = DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );

    final details = NotificationDetails(
      android: androidDetails,
      iOS: iosDetails,
    );

    await _localNotifications.show(
      message.hashCode,
      notification.title,
      notification.body,
      details,
      payload: jsonEncode(message.data),
    );
  }

  /// Get notification channel ID based on type
  String _getChannelId(String type) {
    switch (type) {
      case 'transaction':
        return 'transactions';
      case 'savings':
        return 'savings';
      case 'gig':
        return 'gigs';
      case 'loan':
        return 'loans';
      case 'promotion':
        return 'promotions';
      default:
        return 'general';
    }
  }

  /// Handle notification tap from local notification
  static void _onNotificationTap(NotificationResponse response) {
    if (response.payload != null) {
      try {
        final data = jsonDecode(response.payload!) as Map<String, dynamic>;
        _notificationTapController.add(data);
      } catch (e) {
        Logger().e('Failed to parse notification payload', error: e);
      }
    }
  }

  /// Handle background notification tap
  @pragma('vm:entry-point')
  static void _onBackgroundNotificationTap(NotificationResponse response) {
    _onNotificationTap(response);
  }

  /// Subscribe to topic
  Future<void> subscribeToTopic(String topic) async {
    try {
      await FirebaseMessaging.instance.subscribeToTopic(topic);
      _logger.i('Subscribed to topic: $topic');
    } catch (e) {
      _logger.e('Failed to subscribe to topic: $topic', error: e);
    }
  }

  /// Unsubscribe from topic
  Future<void> unsubscribeFromTopic(String topic) async {
    try {
      await FirebaseMessaging.instance.unsubscribeFromTopic(topic);
      _logger.i('Unsubscribed from topic: $topic');
    } catch (e) {
      _logger.e('Failed to unsubscribe from topic: $topic', error: e);
    }
  }

  /// Clear notification badge
  Future<void> clearBadge() async {
    if (Platform.isIOS) {
      await _localNotifications
          .resolvePlatformSpecificImplementation<
              IOSFlutterLocalNotificationsPlugin>()
          ?.cancelAll();
    }
  }

  /// Get notification permission status
  Future<bool> hasPermission() async {
    final settings = await FirebaseMessaging.instance.getNotificationSettings();
    return settings.authorizationStatus == AuthorizationStatus.authorized ||
        settings.authorizationStatus == AuthorizationStatus.provisional;
  }

  /// Request notification permission
  Future<bool> requestPermission() async {
    final settings = await FirebaseMessaging.instance.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );
    return settings.authorizationStatus == AuthorizationStatus.authorized ||
        settings.authorizationStatus == AuthorizationStatus.provisional;
  }

  /// Dispose resources
  void dispose() {
    _messageController.close();
    _notificationTapController.close();
  }
}

/// Notification payload model
class NotificationPayload {
  final String type;
  final String? action;
  final String? resourceId;
  final Map<String, dynamic> extra;

  NotificationPayload({
    required this.type,
    this.action,
    this.resourceId,
    this.extra = const {},
  });

  factory NotificationPayload.fromMap(Map<String, dynamic> map) {
    return NotificationPayload(
      type: map['type'] ?? 'general',
      action: map['action'],
      resourceId: map['resource_id'] ?? map['resourceId'],
      extra: Map<String, dynamic>.from(map)
        ..remove('type')
        ..remove('action')
        ..remove('resource_id')
        ..remove('resourceId'),
    );
  }

  /// Route to navigate based on notification type
  String? get route {
    switch (type) {
      case 'transaction':
        return resourceId != null ? '/transactions/$resourceId' : '/wallet';
      case 'savings':
        return resourceId != null ? '/savings/circle/$resourceId' : '/savings';
      case 'gig':
        return resourceId != null ? '/gigs/$resourceId' : '/gigs';
      case 'loan':
        return resourceId != null ? '/credit/loan/$resourceId' : '/credit';
      case 'proposal':
        return '/gigs/my-gigs';
      case 'contract':
        return resourceId != null ? '/gigs/contract/$resourceId' : '/gigs/my-gigs';
      default:
        return '/notifications';
    }
  }
}

/// Firebase service provider
final firebaseServiceProvider = Provider<FirebaseService>((ref) {
  final secureStorage = ref.watch(secureStorageProvider);
  return FirebaseService(secureStorage: secureStorage);
});

/// Secure storage provider (duplicate for standalone usage)
final secureStorageProvider = Provider<SecureStorage>((ref) {
  return SecureStorage();
});

/// FCM token provider
final fcmTokenProvider = FutureProvider<String?>((ref) async {
  final service = ref.watch(firebaseServiceProvider);
  return await service.getFCMToken();
});

/// Notification permission provider
final notificationPermissionProvider = FutureProvider<bool>((ref) async {
  final service = ref.watch(firebaseServiceProvider);
  return await service.hasPermission();
});

/// Notification stream provider
final notificationStreamProvider = StreamProvider<RemoteMessage>((ref) {
  return FirebaseService.onMessage;
});

/// Notification tap stream provider
final notificationTapProvider = StreamProvider<Map<String, dynamic>>((ref) {
  return FirebaseService.onNotificationTap;
});
