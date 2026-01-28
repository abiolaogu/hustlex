import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// Local cache service using Hive for offline data storage
/// 
/// This provides a type-safe caching layer for:
/// - User profile data
/// - Wallet balances
/// - Recent transactions
/// - Gigs listings
/// - Savings circles
/// - Credit information
class LocalCacheService {
  static const String _userBox = 'user_cache';
  static const String _walletBox = 'wallet_cache';
  static const String _gigsBox = 'gigs_cache';
  static const String _savingsBox = 'savings_cache';
  static const String _creditBox = 'credit_cache';
  static const String _settingsBox = 'settings_cache';
  static const String _generalBox = 'general_cache';

  late Box<String> _userCacheBox;
  late Box<String> _walletCacheBox;
  late Box<String> _gigsCacheBox;
  late Box<String> _savingsCacheBox;
  late Box<String> _creditCacheBox;
  late Box<String> _settingsCacheBox;
  late Box<String> _generalCacheBox;

  bool _isInitialized = false;

  /// Initialize Hive and open all cache boxes
  Future<void> initialize() async {
    if (_isInitialized) return;

    try {
      await Hive.initFlutter();

      _userCacheBox = await Hive.openBox<String>(_userBox);
      _walletCacheBox = await Hive.openBox<String>(_walletBox);
      _gigsCacheBox = await Hive.openBox<String>(_gigsBox);
      _savingsCacheBox = await Hive.openBox<String>(_savingsBox);
      _creditCacheBox = await Hive.openBox<String>(_creditBox);
      _settingsCacheBox = await Hive.openBox<String>(_settingsBox);
      _generalCacheBox = await Hive.openBox<String>(_generalBox);

      _isInitialized = true;
    } catch (e) {
      debugPrint('LocalCacheService initialization error: $e');
      rethrow;
    }
  }

  // ==================== User Cache ====================

  /// Cache user profile data
  Future<void> cacheUserProfile(Map<String, dynamic> profile) async {
    await _userCacheBox.put('profile', jsonEncode(profile));
    await _userCacheBox.put('profile_timestamp', DateTime.now().toIso8601String());
  }

  /// Get cached user profile
  Map<String, dynamic>? getCachedUserProfile() {
    final data = _userCacheBox.get('profile');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  /// Check if user profile cache is valid (not older than 1 hour)
  bool isUserProfileCacheValid() {
    final timestamp = _userCacheBox.get('profile_timestamp');
    if (timestamp == null) return false;
    final cached = DateTime.parse(timestamp);
    return DateTime.now().difference(cached).inHours < 1;
  }

  // ==================== Wallet Cache ====================

  /// Cache wallet data
  Future<void> cacheWallet(Map<String, dynamic> wallet) async {
    await _walletCacheBox.put('wallet', jsonEncode(wallet));
    await _walletCacheBox.put('wallet_timestamp', DateTime.now().toIso8601String());
  }

  /// Get cached wallet
  Map<String, dynamic>? getCachedWallet() {
    final data = _walletCacheBox.get('wallet');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  /// Cache recent transactions
  Future<void> cacheTransactions(List<Map<String, dynamic>> transactions) async {
    await _walletCacheBox.put('transactions', jsonEncode(transactions));
    await _walletCacheBox.put('transactions_timestamp', DateTime.now().toIso8601String());
  }

  /// Get cached transactions
  List<Map<String, dynamic>> getCachedTransactions() {
    final data = _walletCacheBox.get('transactions');
    if (data == null) return [];
    final list = jsonDecode(data) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  /// Cache bank accounts
  Future<void> cacheBankAccounts(List<Map<String, dynamic>> accounts) async {
    await _walletCacheBox.put('bank_accounts', jsonEncode(accounts));
  }

  /// Get cached bank accounts
  List<Map<String, dynamic>> getCachedBankAccounts() {
    final data = _walletCacheBox.get('bank_accounts');
    if (data == null) return [];
    final list = jsonDecode(data) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  // ==================== Gigs Cache ====================

  /// Cache gigs list
  Future<void> cacheGigs(List<Map<String, dynamic>> gigs, {String key = 'all'}) async {
    await _gigsCacheBox.put('gigs_$key', jsonEncode(gigs));
    await _gigsCacheBox.put('gigs_${key}_timestamp', DateTime.now().toIso8601String());
  }

  /// Get cached gigs
  List<Map<String, dynamic>> getCachedGigs({String key = 'all'}) {
    final data = _gigsCacheBox.get('gigs_$key');
    if (data == null) return [];
    final list = jsonDecode(data) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  /// Check if gigs cache is valid (not older than 15 minutes)
  bool isGigsCacheValid({String key = 'all'}) {
    final timestamp = _gigsCacheBox.get('gigs_${key}_timestamp');
    if (timestamp == null) return false;
    final cached = DateTime.parse(timestamp);
    return DateTime.now().difference(cached).inMinutes < 15;
  }

  /// Cache single gig details
  Future<void> cacheGigDetails(String gigId, Map<String, dynamic> gig) async {
    await _gigsCacheBox.put('gig_$gigId', jsonEncode(gig));
  }

  /// Get cached gig details
  Map<String, dynamic>? getCachedGigDetails(String gigId) {
    final data = _gigsCacheBox.get('gig_$gigId');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  /// Cache gig categories
  Future<void> cacheGigCategories(List<Map<String, dynamic>> categories) async {
    await _gigsCacheBox.put('categories', jsonEncode(categories));
  }

  /// Get cached gig categories
  List<Map<String, dynamic>> getCachedGigCategories() {
    final data = _gigsCacheBox.get('categories');
    if (data == null) return [];
    final list = jsonDecode(data) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  // ==================== Savings Cache ====================

  /// Cache savings circles
  Future<void> cacheSavingsCircles(List<Map<String, dynamic>> circles, {String key = 'my'}) async {
    await _savingsCacheBox.put('circles_$key', jsonEncode(circles));
    await _savingsCacheBox.put('circles_${key}_timestamp', DateTime.now().toIso8601String());
  }

  /// Get cached savings circles
  List<Map<String, dynamic>> getCachedSavingsCircles({String key = 'my'}) {
    final data = _savingsCacheBox.get('circles_$key');
    if (data == null) return [];
    final list = jsonDecode(data) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  /// Cache single circle details
  Future<void> cacheCircleDetails(String circleId, Map<String, dynamic> circle) async {
    await _savingsCacheBox.put('circle_$circleId', jsonEncode(circle));
  }

  /// Get cached circle details
  Map<String, dynamic>? getCachedCircleDetails(String circleId) {
    final data = _savingsCacheBox.get('circle_$circleId');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  /// Cache savings stats
  Future<void> cacheSavingsStats(Map<String, dynamic> stats) async {
    await _savingsCacheBox.put('stats', jsonEncode(stats));
  }

  /// Get cached savings stats
  Map<String, dynamic>? getCachedSavingsStats() {
    final data = _savingsCacheBox.get('stats');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  // ==================== Credit Cache ====================

  /// Cache credit profile
  Future<void> cacheCreditProfile(Map<String, dynamic> profile) async {
    await _creditCacheBox.put('profile', jsonEncode(profile));
    await _creditCacheBox.put('profile_timestamp', DateTime.now().toIso8601String());
  }

  /// Get cached credit profile
  Map<String, dynamic>? getCachedCreditProfile() {
    final data = _creditCacheBox.get('profile');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  /// Cache loans list
  Future<void> cacheLoans(List<Map<String, dynamic>> loans) async {
    await _creditCacheBox.put('loans', jsonEncode(loans));
  }

  /// Get cached loans
  List<Map<String, dynamic>> getCachedLoans() {
    final data = _creditCacheBox.get('loans');
    if (data == null) return [];
    final list = jsonDecode(data) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  /// Cache loan details
  Future<void> cacheLoanDetails(String loanId, Map<String, dynamic> loan) async {
    await _creditCacheBox.put('loan_$loanId', jsonEncode(loan));
  }

  /// Get cached loan details
  Map<String, dynamic>? getCachedLoanDetails(String loanId) {
    final data = _creditCacheBox.get('loan_$loanId');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  /// Cache credit score history
  Future<void> cacheCreditScoreHistory(List<Map<String, dynamic>> history) async {
    await _creditCacheBox.put('score_history', jsonEncode(history));
  }

  /// Get cached credit score history
  List<Map<String, dynamic>> getCachedCreditScoreHistory() {
    final data = _creditCacheBox.get('score_history');
    if (data == null) return [];
    final list = jsonDecode(data) as List;
    return list.map((e) => e as Map<String, dynamic>).toList();
  }

  // ==================== Settings Cache ====================

  /// Cache user settings
  Future<void> cacheSettings(Map<String, dynamic> settings) async {
    await _settingsCacheBox.put('settings', jsonEncode(settings));
  }

  /// Get cached settings
  Map<String, dynamic>? getCachedSettings() {
    final data = _settingsCacheBox.get('settings');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  /// Cache notification preferences
  Future<void> cacheNotificationPreferences(Map<String, dynamic> prefs) async {
    await _settingsCacheBox.put('notification_prefs', jsonEncode(prefs));
  }

  /// Get cached notification preferences
  Map<String, dynamic>? getCachedNotificationPreferences() {
    final data = _settingsCacheBox.get('notification_prefs');
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  // ==================== General Cache ====================

  /// Cache any key-value data
  Future<void> cache(String key, dynamic value) async {
    await _generalCacheBox.put(key, jsonEncode(value));
    await _generalCacheBox.put('${key}_timestamp', DateTime.now().toIso8601String());
  }

  /// Get cached data
  T? get<T>(String key) {
    final data = _generalCacheBox.get(key);
    if (data == null) return null;
    return jsonDecode(data) as T;
  }

  /// Check if cache is valid
  bool isCacheValid(String key, {int maxMinutes = 30}) {
    final timestamp = _generalCacheBox.get('${key}_timestamp');
    if (timestamp == null) return false;
    final cached = DateTime.parse(timestamp);
    return DateTime.now().difference(cached).inMinutes < maxMinutes;
  }

  /// Remove specific cache
  Future<void> remove(String key) async {
    await _generalCacheBox.delete(key);
    await _generalCacheBox.delete('${key}_timestamp');
  }

  // ==================== Cache Management ====================

  /// Clear all user-specific cache (on logout)
  Future<void> clearUserCache() async {
    await _userCacheBox.clear();
    await _walletCacheBox.clear();
    await _gigsCacheBox.clear();
    await _savingsCacheBox.clear();
    await _creditCacheBox.clear();
    // Keep settings cache for app preferences
  }

  /// Clear all cache
  Future<void> clearAllCache() async {
    await _userCacheBox.clear();
    await _walletCacheBox.clear();
    await _gigsCacheBox.clear();
    await _savingsCacheBox.clear();
    await _creditCacheBox.clear();
    await _settingsCacheBox.clear();
    await _generalCacheBox.clear();
  }

  /// Get cache size in bytes (approximate)
  int getCacheSize() {
    int total = 0;
    for (final box in [
      _userCacheBox,
      _walletCacheBox,
      _gigsCacheBox,
      _savingsCacheBox,
      _creditCacheBox,
      _settingsCacheBox,
      _generalCacheBox,
    ]) {
      for (final key in box.keys) {
        final value = box.get(key);
        if (value != null) {
          total += value.length;
        }
      }
    }
    return total;
  }

  /// Cleanup old cache entries
  Future<void> cleanupOldCache({int maxDays = 7}) async {
    final cutoff = DateTime.now().subtract(Duration(days: maxDays));

    for (final box in [
      _generalCacheBox,
      _gigsCacheBox,
    ]) {
      final keysToDelete = <String>[];
      for (final key in box.keys) {
        if (key.toString().endsWith('_timestamp')) {
          final timestamp = box.get(key);
          if (timestamp != null) {
            final cached = DateTime.tryParse(timestamp);
            if (cached != null && cached.isBefore(cutoff)) {
              keysToDelete.add(key.toString());
              // Also delete the associated data
              final dataKey = key.toString().replaceAll('_timestamp', '');
              keysToDelete.add(dataKey);
            }
          }
        }
      }
      for (final key in keysToDelete) {
        await box.delete(key);
      }
    }
  }
}

/// Provider for LocalCacheService
final localCacheServiceProvider = Provider<LocalCacheService>((ref) {
  return LocalCacheService();
});

/// FutureProvider for initialization
final localCacheInitProvider = FutureProvider<void>((ref) async {
  final cacheService = ref.read(localCacheServiceProvider);
  await cacheService.initialize();
});
