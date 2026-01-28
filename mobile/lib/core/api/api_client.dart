import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:logger/logger.dart';

import '../constants/app_constants.dart';
import '../exceptions/api_exception.dart';

/// =============================================================================
/// API CLIENT PROVIDER
/// =============================================================================

final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient();
});

/// =============================================================================
/// API RESPONSE MODEL
/// =============================================================================

class ApiResponse<T> {
  final bool success;
  final T? data;
  final String? message;
  final String? error;
  final Map<String, dynamic>? meta;

  ApiResponse({
    required this.success,
    this.data,
    this.message,
    this.error,
    this.meta,
  });

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic)? fromJsonT,
  ) {
    return ApiResponse(
      success: json['success'] ?? false,
      data: fromJsonT != null && json['data'] != null
          ? fromJsonT(json['data'])
          : json['data'],
      message: json['message'],
      error: json['error'],
      meta: json['meta'],
    );
  }

  bool get hasError => !success || error != null;
}

/// =============================================================================
/// PAGINATED RESPONSE
/// =============================================================================

class PaginatedResponse<T> {
  final List<T> items;
  final int page;
  final int limit;
  final int total;
  final int totalPages;
  final bool hasNext;
  final bool hasPrev;

  PaginatedResponse({
    required this.items,
    required this.page,
    required this.limit,
    required this.total,
    required this.totalPages,
    required this.hasNext,
    required this.hasPrev,
  });

  factory PaginatedResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Map<String, dynamic>) fromJsonT,
  ) {
    final data = json['data'] ?? json;
    final items = (data['items'] as List? ?? [])
        .map((item) => fromJsonT(item as Map<String, dynamic>))
        .toList();
    
    return PaginatedResponse(
      items: items,
      page: data['page'] ?? 1,
      limit: data['limit'] ?? AppConstants.defaultPageSize,
      total: data['total'] ?? items.length,
      totalPages: data['total_pages'] ?? 1,
      hasNext: data['has_next'] ?? false,
      hasPrev: data['has_prev'] ?? false,
    );
  }
}

/// =============================================================================
/// API CLIENT
/// =============================================================================

class ApiClient {
  late final Dio _dio;
  final _storage = const FlutterSecureStorage();
  final _logger = Logger(
    printer: PrettyPrinter(
      methodCount: 0,
      errorMethodCount: 5,
      lineLength: 80,
      colors: true,
      printEmojis: true,
    ),
  );

  ApiClient() {
    _dio = Dio(
      BaseOptions(
        baseUrl: kDebugMode 
            ? AppConstants.devBaseUrl 
            : AppConstants.baseUrl,
        connectTimeout: AppConstants.connectTimeout,
        receiveTimeout: AppConstants.receiveTimeout,
        sendTimeout: AppConstants.sendTimeout,
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
      ),
    );

    _dio.interceptors.addAll([
      _AuthInterceptor(_storage, _dio),
      _LoggingInterceptor(_logger),
      _ErrorInterceptor(),
    ]);
  }

  Dio get dio => _dio;

  // ===================== GET =====================
  Future<ApiResponse<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    T Function(dynamic)? fromJson,
    Options? options,
  }) async {
    try {
      final response = await _dio.get(
        path,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(response.data, fromJson);
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  // ===================== POST =====================
  Future<ApiResponse<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    T Function(dynamic)? fromJson,
    Options? options,
  }) async {
    try {
      final response = await _dio.post(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(response.data, fromJson);
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  // ===================== PUT =====================
  Future<ApiResponse<T>> put<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    T Function(dynamic)? fromJson,
    Options? options,
  }) async {
    try {
      final response = await _dio.put(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(response.data, fromJson);
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  // ===================== PATCH =====================
  Future<ApiResponse<T>> patch<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    T Function(dynamic)? fromJson,
    Options? options,
  }) async {
    try {
      final response = await _dio.patch(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(response.data, fromJson);
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  // ===================== DELETE =====================
  Future<ApiResponse<T>> delete<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    T Function(dynamic)? fromJson,
    Options? options,
  }) async {
    try {
      final response = await _dio.delete(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      );
      return ApiResponse.fromJson(response.data, fromJson);
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  // ===================== MULTIPART =====================
  Future<ApiResponse<T>> uploadFile<T>(
    String path, {
    required File file,
    required String fieldName,
    Map<String, dynamic>? fields,
    T Function(dynamic)? fromJson,
    void Function(int, int)? onSendProgress,
  }) async {
    try {
      final formData = FormData.fromMap({
        ...?fields,
        fieldName: await MultipartFile.fromFile(
          file.path,
          filename: file.path.split('/').last,
        ),
      });

      final response = await _dio.post(
        path,
        data: formData,
        onSendProgress: onSendProgress,
        options: Options(
          headers: {'Content-Type': 'multipart/form-data'},
        ),
      );

      return ApiResponse.fromJson(response.data, fromJson);
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  // ===================== PAGINATED GET =====================
  Future<PaginatedResponse<T>> getPaginated<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    required T Function(Map<String, dynamic>) fromJson,
    int page = 1,
    int limit = AppConstants.defaultPageSize,
  }) async {
    try {
      final params = {
        ...?queryParameters,
        'page': page,
        'limit': limit,
      };

      final response = await _dio.get(path, queryParameters: params);
      return PaginatedResponse.fromJson(response.data, fromJson);
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  // ===================== ERROR HANDLING =====================
  ApiException _handleDioError(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return ApiException(
          message: 'Connection timed out. Please try again.',
          code: 'TIMEOUT',
        );
      case DioExceptionType.connectionError:
        return ApiException(
          message: 'No internet connection. Please check your network.',
          code: 'NO_CONNECTION',
        );
      case DioExceptionType.badResponse:
        return _handleBadResponse(e.response);
      case DioExceptionType.cancel:
        return ApiException(
          message: 'Request was cancelled.',
          code: 'CANCELLED',
        );
      default:
        return ApiException(
          message: 'Something went wrong. Please try again.',
          code: 'UNKNOWN',
        );
    }
  }

  ApiException _handleBadResponse(Response? response) {
    if (response == null) {
      return ApiException(
        message: 'No response from server',
        code: 'NO_RESPONSE',
      );
    }

    final statusCode = response.statusCode ?? 500;
    final data = response.data;
    String message;
    String? code;

    if (data is Map<String, dynamic>) {
      message = data['error'] ?? data['message'] ?? 'An error occurred';
      code = data['code']?.toString();
    } else {
      message = 'An error occurred';
    }

    switch (statusCode) {
      case 400:
        return ApiException(
          message: message,
          code: code ?? 'BAD_REQUEST',
          statusCode: statusCode,
        );
      case 401:
        return ApiException(
          message: message,
          code: code ?? 'UNAUTHORIZED',
          statusCode: statusCode,
        );
      case 403:
        return ApiException(
          message: message,
          code: code ?? 'FORBIDDEN',
          statusCode: statusCode,
        );
      case 404:
        return ApiException(
          message: message,
          code: code ?? 'NOT_FOUND',
          statusCode: statusCode,
        );
      case 409:
        return ApiException(
          message: message,
          code: code ?? 'CONFLICT',
          statusCode: statusCode,
        );
      case 422:
        return ApiException(
          message: message,
          code: code ?? 'VALIDATION_ERROR',
          statusCode: statusCode,
          errors: data['errors'],
        );
      case 429:
        return ApiException(
          message: 'Too many requests. Please try again later.',
          code: 'RATE_LIMITED',
          statusCode: statusCode,
        );
      case 500:
      case 502:
      case 503:
        return ApiException(
          message: 'Server error. Please try again later.',
          code: 'SERVER_ERROR',
          statusCode: statusCode,
        );
      default:
        return ApiException(
          message: message,
          code: code ?? 'HTTP_ERROR',
          statusCode: statusCode,
        );
    }
  }

  // ===================== TOKEN MANAGEMENT =====================
  Future<void> setAccessToken(String token) async {
    await _storage.write(key: AppConstants.accessTokenKey, value: token);
  }

  Future<void> setRefreshToken(String token) async {
    await _storage.write(key: AppConstants.refreshTokenKey, value: token);
  }

  Future<String?> getAccessToken() async {
    return await _storage.read(key: AppConstants.accessTokenKey);
  }

  Future<String?> getRefreshToken() async {
    return await _storage.read(key: AppConstants.refreshTokenKey);
  }

  Future<void> clearTokens() async {
    await _storage.delete(key: AppConstants.accessTokenKey);
    await _storage.delete(key: AppConstants.refreshTokenKey);
  }

  Future<bool> hasValidToken() async {
    final token = await getAccessToken();
    return token != null && token.isNotEmpty;
  }
}

/// =============================================================================
/// AUTH INTERCEPTOR
/// =============================================================================

class _AuthInterceptor extends Interceptor {
  final FlutterSecureStorage _storage;
  final Dio _dio;
  bool _isRefreshing = false;

  _AuthInterceptor(this._storage, this._dio);

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // Skip auth for public endpoints
    final publicEndpoints = [
      '/auth/request-otp',
      '/auth/verify-otp',
      '/auth/register',
      '/health',
      '/ready',
    ];

    final isPublic = publicEndpoints.any(
      (endpoint) => options.path.contains(endpoint),
    );

    if (!isPublic) {
      final token = await _storage.read(key: AppConstants.accessTokenKey);
      if (token != null) {
        options.headers['Authorization'] = 'Bearer $token';
      }
    }

    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401 && !_isRefreshing) {
      _isRefreshing = true;

      try {
        final refreshToken = await _storage.read(
          key: AppConstants.refreshTokenKey,
        );

        if (refreshToken != null) {
          final response = await _dio.post(
            '/auth/refresh',
            data: {'refresh_token': refreshToken},
          );

          if (response.statusCode == 200) {
            final newAccessToken = response.data['data']['access_token'];
            final newRefreshToken = response.data['data']['refresh_token'];

            await _storage.write(
              key: AppConstants.accessTokenKey,
              value: newAccessToken,
            );
            await _storage.write(
              key: AppConstants.refreshTokenKey,
              value: newRefreshToken,
            );

            // Retry the original request
            err.requestOptions.headers['Authorization'] =
                'Bearer $newAccessToken';

            final retryResponse = await _dio.fetch(err.requestOptions);
            return handler.resolve(retryResponse);
          }
        }

        // Refresh failed - clear tokens and force re-login
        await _storage.delete(key: AppConstants.accessTokenKey);
        await _storage.delete(key: AppConstants.refreshTokenKey);
      } catch (e) {
        // Refresh failed
        await _storage.delete(key: AppConstants.accessTokenKey);
        await _storage.delete(key: AppConstants.refreshTokenKey);
      } finally {
        _isRefreshing = false;
      }
    }

    handler.next(err);
  }
}

/// =============================================================================
/// LOGGING INTERCEPTOR
/// =============================================================================

class _LoggingInterceptor extends Interceptor {
  final Logger _logger;

  _LoggingInterceptor(this._logger);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    if (kDebugMode) {
      _logger.i(
        '📤 REQUEST[${options.method}] => ${options.uri}\n'
        'Headers: ${options.headers}\n'
        'Data: ${options.data}',
      );
    }
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    if (kDebugMode) {
      _logger.i(
        '📥 RESPONSE[${response.statusCode}] <= ${response.requestOptions.uri}\n'
        'Data: ${response.data}',
      );
    }
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (kDebugMode) {
      _logger.e(
        '❌ ERROR[${err.response?.statusCode}] <= ${err.requestOptions.uri}\n'
        'Message: ${err.message}\n'
        'Data: ${err.response?.data}',
      );
    }
    handler.next(err);
  }
}

/// =============================================================================
/// ERROR INTERCEPTOR
/// =============================================================================

class _ErrorInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    // Transform specific error responses
    if (err.response?.statusCode == 401) {
      // Could emit an event for the app to handle logout
    }
    handler.next(err);
  }
}
