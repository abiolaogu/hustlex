import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';

import '../core/api/api_client.dart';
import '../core/storage/secure_storage.dart';
import '../core/storage/local_cache_service.dart';

/// Creates a ProviderContainer for testing
ProviderContainer createTestContainer({
  List<Override>? overrides,
}) {
  return ProviderContainer(
    overrides: overrides ?? [],
  );
}

/// Creates a test widget with ProviderScope
Widget createTestWidget({
  required Widget child,
  List<Override>? overrides,
}) {
  return ProviderScope(
    overrides: overrides ?? [],
    child: MaterialApp(
      home: child,
    ),
  );
}

/// Mock API client for testing
class MockApiClient extends ApiClient {
  final Map<String, dynamic> mockResponses;
  final List<ApiCallRecord> callHistory = [];

  MockApiClient({
    this.mockResponses = const {},
    SecureStorage? secureStorage,
  }) : super(secureStorage: secureStorage ?? MockSecureStorage());

  @override
  Future<Map<String, dynamic>> get(
    String endpoint, {
    Map<String, dynamic>? queryParams,
  }) async {
    _recordCall('GET', endpoint, queryParams: queryParams);
    return _getMockResponse('GET:$endpoint');
  }

  @override
  Future<Map<String, dynamic>> post(
    String endpoint, {
    Map<String, dynamic>? data,
  }) async {
    _recordCall('POST', endpoint, data: data);
    return _getMockResponse('POST:$endpoint');
  }

  @override
  Future<Map<String, dynamic>> put(
    String endpoint, {
    Map<String, dynamic>? data,
  }) async {
    _recordCall('PUT', endpoint, data: data);
    return _getMockResponse('PUT:$endpoint');
  }

  @override
  Future<Map<String, dynamic>> patch(
    String endpoint, {
    Map<String, dynamic>? data,
  }) async {
    _recordCall('PATCH', endpoint, data: data);
    return _getMockResponse('PATCH:$endpoint');
  }

  @override
  Future<Map<String, dynamic>> delete(String endpoint) async {
    _recordCall('DELETE', endpoint);
    return _getMockResponse('DELETE:$endpoint');
  }

  void _recordCall(
    String method,
    String endpoint, {
    Map<String, dynamic>? queryParams,
    Map<String, dynamic>? data,
  }) {
    callHistory.add(ApiCallRecord(
      method: method,
      endpoint: endpoint,
      queryParams: queryParams,
      data: data,
      timestamp: DateTime.now(),
    ));
  }

  Map<String, dynamic> _getMockResponse(String key) {
    if (mockResponses.containsKey(key)) {
      return mockResponses[key] as Map<String, dynamic>;
    }
    // Return default success response
    return {'success': true, 'data': {}};
  }

  /// Set a mock response for a specific endpoint
  void setMockResponse(String method, String endpoint, Map<String, dynamic> response) {
    mockResponses['$method:$endpoint'] = response;
  }

  /// Get the last call for a specific endpoint
  ApiCallRecord? getLastCall(String method, String endpoint) {
    return callHistory
        .where((c) => c.method == method && c.endpoint == endpoint)
        .lastOrNull;
  }

  /// Clear call history
  void clearHistory() {
    callHistory.clear();
  }
}

/// Record of an API call
class ApiCallRecord {
  final String method;
  final String endpoint;
  final Map<String, dynamic>? queryParams;
  final Map<String, dynamic>? data;
  final DateTime timestamp;

  ApiCallRecord({
    required this.method,
    required this.endpoint,
    this.queryParams,
    this.data,
    required this.timestamp,
  });

  @override
  String toString() => '$method $endpoint';
}

/// Mock secure storage for testing
class MockSecureStorage extends SecureStorage {
  final Map<String, String> _storage = {};

  @override
  Future<void> write({required String key, required String value}) async {
    _storage[key] = value;
  }

  @override
  Future<String?> read({required String key}) async {
    return _storage[key];
  }

  @override
  Future<void> delete({required String key}) async {
    _storage.remove(key);
  }

  @override
  Future<void> deleteAll() async {
    _storage.clear();
  }

  @override
  Future<bool> containsKey({required String key}) async {
    return _storage.containsKey(key);
  }

  @override
  Future<Map<String, String>> readAll() async {
    return Map.from(_storage);
  }

  /// Get current storage contents for verification
  Map<String, String> get contents => Map.unmodifiable(_storage);

  /// Set a value directly (for test setup)
  void setValue(String key, String value) {
    _storage[key] = value;
  }
}

/// Mock local cache service for testing
class MockLocalCacheService extends LocalCacheService {
  final Map<String, dynamic> _cache = {};
  final Map<String, DateTime> _timestamps = {};

  @override
  Future<void> initialize() async {
    // No-op for tests
  }

  @override
  Future<void> cache<T>(String key, T data, {Duration? expiry}) async {
    _cache[key] = data;
    _timestamps[key] = DateTime.now();
  }

  @override
  Future<T?> get<T>(String key) async {
    return _cache[key] as T?;
  }

  @override
  Future<bool> isCacheValid(String key, Duration maxAge) async {
    final timestamp = _timestamps[key];
    if (timestamp == null) return false;
    return DateTime.now().difference(timestamp) < maxAge;
  }

  @override
  Future<void> remove(String key) async {
    _cache.remove(key);
    _timestamps.remove(key);
  }

  @override
  Future<void> clearAllCache() async {
    _cache.clear();
    _timestamps.clear();
  }

  /// Set cache data directly (for test setup)
  void setCache(String key, dynamic data) {
    _cache[key] = data;
    _timestamps[key] = DateTime.now();
  }
}

/// Test fixtures for common data
class TestFixtures {
  /// Sample user data
  static Map<String, dynamic> get user => {
        'id': 'user_123',
        'email': 'test@example.com',
        'phone': '+2348012345678',
        'first_name': 'Test',
        'last_name': 'User',
        'is_verified': true,
        'kyc_level': 2,
        'created_at': '2024-01-01T00:00:00Z',
      };

  /// Sample wallet data
  static Map<String, dynamic> get wallet => {
        'id': 'wallet_123',
        'balance': 50000.00,
        'currency': 'NGN',
        'available_balance': 45000.00,
        'escrow_balance': 5000.00,
        'savings_balance': 10000.00,
      };

  /// Sample transaction data
  static Map<String, dynamic> get transaction => {
        'id': 'txn_123',
        'type': 'deposit',
        'amount': 10000.00,
        'status': 'completed',
        'reference': 'REF_123456',
        'description': 'Wallet deposit',
        'created_at': '2024-01-15T10:30:00Z',
      };

  /// Sample gig data
  static Map<String, dynamic> get gig => {
        'id': 'gig_123',
        'title': 'Build a mobile app',
        'description': 'Need a Flutter developer',
        'budget_min': 100000,
        'budget_max': 200000,
        'category': 'Development',
        'status': 'open',
        'created_at': '2024-01-10T00:00:00Z',
      };

  /// Sample savings circle data
  static Map<String, dynamic> get circle => {
        'id': 'circle_123',
        'name': 'Friends Ajo',
        'amount': 10000.00,
        'frequency': 'weekly',
        'total_members': 10,
        'current_members': 5,
        'status': 'active',
        'created_at': '2024-01-01T00:00:00Z',
      };

  /// Sample loan data
  static Map<String, dynamic> get loan => {
        'id': 'loan_123',
        'amount': 50000.00,
        'interest_rate': 5.0,
        'tenure_months': 3,
        'monthly_payment': 17500.00,
        'status': 'active',
        'disbursed_at': '2024-01-01T00:00:00Z',
      };

  /// Sample credit score data
  static Map<String, dynamic> get creditScore => {
        'score': 720,
        'max_score': 850,
        'rating': 'Good',
        'last_updated': '2024-01-15T00:00:00Z',
      };
}

/// Pump widget and settle with timeout
extension WidgetTesterExtensions on WidgetTester {
  /// Pump and settle with a shorter timeout for faster tests
  Future<void> pumpAndSettleQuick() async {
    await pumpAndSettle(const Duration(milliseconds: 100));
  }

  /// Find widget by key and tap
  Future<void> tapByKey(Key key) async {
    await tap(find.byKey(key));
    await pumpAndSettle();
  }

  /// Find widget by text and tap
  Future<void> tapByText(String text) async {
    await tap(find.text(text));
    await pumpAndSettle();
  }

  /// Enter text in a text field by key
  Future<void> enterTextByKey(Key key, String text) async {
    await enterText(find.byKey(key), text);
    await pumpAndSettle();
  }
}

/// Golden test helpers
class GoldenTestHelper {
  /// Compare widget with golden file
  static Future<void> matchGolden(
    WidgetTester tester,
    Widget widget,
    String goldenName, {
    Size? size,
  }) async {
    await tester.pumpWidget(
      MaterialApp(
        debugShowCheckedModeBanner: false,
        home: widget,
      ),
    );
    await tester.pumpAndSettle();

    await expectLater(
      find.byType(MaterialApp),
      matchesGoldenFile('goldens/$goldenName.png'),
    );
  }
}
