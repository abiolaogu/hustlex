import 'package:flutter_dotenv/flutter_dotenv.dart';

/// Application configuration loaded from environment variables
class AppConfig {
  AppConfig._();

  static const String appName = 'HustleX';
  static const String appVersion = '1.0.0';

  // API Configuration
  static String get apiBaseUrl =>
      dotenv.env['API_BASE_URL'] ?? 'https://api.hustlex.ng/api/v1';

  static String get paystackPublicKey =>
      dotenv.env['PAYSTACK_PUBLIC_KEY'] ?? '';

  // Feature Flags
  static bool get isProduction =>
      dotenv.env['ENVIRONMENT'] == 'production';

  static bool get enableAnalytics =>
      dotenv.env['ENABLE_ANALYTICS'] == 'true';

  static bool get enableCrashlytics =>
      dotenv.env['ENABLE_CRASHLYTICS'] == 'true';

  // Timeouts
  static const Duration connectionTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);

  // Cache durations
  static const Duration tokenCacheDuration = Duration(minutes: 14);
  static const Duration userCacheDuration = Duration(hours: 1);

  // Validation
  static const int minPasswordLength = 8;
  static const int otpLength = 6;
  static const int pinLength = 4;

  // Pagination
  static const int defaultPageSize = 20;

  // File upload
  static const int maxImageSizeBytes = 5 * 1024 * 1024; // 5MB
  static const List<String> allowedImageTypes = ['jpg', 'jpeg', 'png'];

  // Currency
  static const String currencyCode = 'NGN';
  static const String currencySymbol = '₦';
  static const int koboMultiplier = 100; // 1 Naira = 100 Kobo
}
