import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Secure storage keys
abstract class StorageKeys {
  static const accessToken = 'access_token';
  static const refreshToken = 'refresh_token';
  static const tokenExpiry = 'token_expiry';
  static const userId = 'user_id';
  static const userPhone = 'user_phone';
  static const biometricEnabled = 'biometric_enabled';
  static const pinHash = 'pin_hash';
  static const deviceId = 'device_id';
  static const fcmToken = 'fcm_token';
  static const onboardingComplete = 'onboarding_complete';
}

/// Secure storage wrapper for sensitive data
class SecureStorage {
  final FlutterSecureStorage _storage;

  SecureStorage({FlutterSecureStorage? storage})
      : _storage = storage ??
            const FlutterSecureStorage(
              aOptions: AndroidOptions(
                encryptedSharedPreferences: true,
                sharedPreferencesName: 'hustlex_secure_prefs',
                preferencesKeyPrefix: 'hustlex_',
              ),
              iOptions: IOSOptions(
                accessibility: KeychainAccessibility.first_unlock_this_device,
                accountName: 'HustleX',
              ),
            );

  // Token management
  Future<void> setAccessToken(String token) =>
      _storage.write(key: StorageKeys.accessToken, value: token);

  Future<String?> getAccessToken() =>
      _storage.read(key: StorageKeys.accessToken);

  Future<void> deleteAccessToken() =>
      _storage.delete(key: StorageKeys.accessToken);

  Future<void> setRefreshToken(String token) =>
      _storage.write(key: StorageKeys.refreshToken, value: token);

  Future<String?> getRefreshToken() =>
      _storage.read(key: StorageKeys.refreshToken);

  Future<void> deleteRefreshToken() =>
      _storage.delete(key: StorageKeys.refreshToken);

  Future<void> setTokenExpiry(DateTime expiry) =>
      _storage.write(key: StorageKeys.tokenExpiry, value: expiry.toIso8601String());

  Future<DateTime?> getTokenExpiry() async {
    final value = await _storage.read(key: StorageKeys.tokenExpiry);
    if (value == null) return null;
    return DateTime.tryParse(value);
  }

  // User info
  Future<void> setUserId(String userId) =>
      _storage.write(key: StorageKeys.userId, value: userId);

  Future<String?> getUserId() => _storage.read(key: StorageKeys.userId);

  Future<void> setUserPhone(String phone) =>
      _storage.write(key: StorageKeys.userPhone, value: phone);

  Future<String?> getUserPhone() => _storage.read(key: StorageKeys.userPhone);

  // Biometric settings
  Future<void> setBiometricEnabled(bool enabled) =>
      _storage.write(key: StorageKeys.biometricEnabled, value: enabled.toString());

  Future<bool> getBiometricEnabled() async {
    final value = await _storage.read(key: StorageKeys.biometricEnabled);
    return value == 'true';
  }

  // PIN storage (hashed)
  Future<void> setPinHash(String hash) =>
      _storage.write(key: StorageKeys.pinHash, value: hash);

  Future<String?> getPinHash() => _storage.read(key: StorageKeys.pinHash);

  // Device info
  Future<void> setDeviceId(String deviceId) =>
      _storage.write(key: StorageKeys.deviceId, value: deviceId);

  Future<String?> getDeviceId() => _storage.read(key: StorageKeys.deviceId);

  // FCM token
  Future<void> setFcmToken(String token) =>
      _storage.write(key: StorageKeys.fcmToken, value: token);

  Future<String?> getFcmToken() => _storage.read(key: StorageKeys.fcmToken);

  // Onboarding
  Future<void> setOnboardingComplete(bool complete) =>
      _storage.write(key: StorageKeys.onboardingComplete, value: complete.toString());

  Future<bool> getOnboardingComplete() async {
    final value = await _storage.read(key: StorageKeys.onboardingComplete);
    return value == 'true';
  }

  // Clear all auth data (logout)
  Future<void> clearAuth() async {
    await _storage.delete(key: StorageKeys.accessToken);
    await _storage.delete(key: StorageKeys.refreshToken);
    await _storage.delete(key: StorageKeys.tokenExpiry);
  }

  // Clear all data
  Future<void> clearAll() async {
    await _storage.deleteAll();
  }

  // Check if has valid tokens
  Future<bool> hasValidTokens() async {
    final accessToken = await getAccessToken();
    final expiry = await getTokenExpiry();
    
    if (accessToken == null || accessToken.isEmpty) return false;
    if (expiry == null) return false;
    
    return DateTime.now().isBefore(expiry);
  }

  // Generic read/write for custom data
  Future<void> write(String key, String value) =>
      _storage.write(key: key, value: value);

  Future<String?> read(String key) => _storage.read(key: key);

  Future<void> delete(String key) => _storage.delete(key: key);

  Future<Map<String, String>> readAll() => _storage.readAll();

  Future<bool> containsKey(String key) => _storage.containsKey(key: key);
}
