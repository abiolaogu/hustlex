import 'dart:async';
import 'dart:io';

import 'package:dio/dio.dart';

import '../../exceptions/api_exception.dart';

/// Base result class for repository operations
sealed class Result<T> {
  const Result();

  bool get isSuccess => this is Success<T>;
  bool get isFailure => this is Failure<T>;

  T? get data => this is Success<T> ? (this as Success<T>).data : null;
  String? get error => this is Failure<T> ? (this as Failure<T>).message : null;
  int? get errorCode => this is Failure<T> ? (this as Failure<T>).code : null;

  R when<R>({
    required R Function(T data) success,
    required R Function(String message, int? code) failure,
  }) {
    if (this is Success<T>) {
      return success((this as Success<T>).data);
    } else {
      final failure_ = this as Failure<T>;
      return failure(failure_.message, failure_.code);
    }
  }

  R? maybeWhen<R>({
    R Function(T data)? success,
    R Function(String message, int? code)? failure,
    R Function()? orElse,
  }) {
    if (this is Success<T> && success != null) {
      return success((this as Success<T>).data);
    } else if (this is Failure<T> && failure != null) {
      final failure_ = this as Failure<T>;
      return failure(failure_.message, failure_.code);
    }
    return orElse?.call();
  }
}

class Success<T> extends Result<T> {
  final T data;
  const Success(this.data);
}

class Failure<T> extends Result<T> {
  final String message;
  final int? code;
  final dynamic originalError;

  const Failure(this.message, {this.code, this.originalError});
}

/// Base repository class with common error handling
abstract class BaseRepository {
  const BaseRepository();

  /// Execute an API call and handle errors uniformly
  Future<Result<T>> safeCall<T>(Future<T> Function() call) async {
    try {
      final result = await call();
      return Success(result);
    } on ApiException catch (e) {
      return Failure(e.message, code: e.statusCode, originalError: e);
    } on DioException catch (e) {
      return Failure(_handleDioError(e), code: e.response?.statusCode, originalError: e);
    } on SocketException {
      return const Failure('No internet connection. Please check your network.');
    } on TimeoutException {
      return const Failure('Request timed out. Please try again.');
    } on FormatException catch (e) {
      return Failure('Invalid response format: ${e.message}');
    } catch (e) {
      return Failure('An unexpected error occurred: $e', originalError: e);
    }
  }

  /// Execute an API call that returns void
  Future<Result<void>> safeVoidCall(Future<void> Function() call) async {
    try {
      await call();
      return const Success(null);
    } on ApiException catch (e) {
      return Failure(e.message, code: e.statusCode, originalError: e);
    } on DioException catch (e) {
      return Failure(_handleDioError(e), code: e.response?.statusCode, originalError: e);
    } on SocketException {
      return const Failure('No internet connection. Please check your network.');
    } on TimeoutException {
      return const Failure('Request timed out. Please try again.');
    } catch (e) {
      return Failure('An unexpected error occurred: $e', originalError: e);
    }
  }

  String _handleDioError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return 'Connection timed out. Please try again.';
      case DioExceptionType.badCertificate:
        return 'Security certificate error. Please contact support.';
      case DioExceptionType.badResponse:
        return _handleBadResponse(error.response);
      case DioExceptionType.cancel:
        return 'Request was cancelled.';
      case DioExceptionType.connectionError:
        return 'Unable to connect to server. Please check your internet.';
      case DioExceptionType.unknown:
        if (error.error is SocketException) {
          return 'No internet connection.';
        }
        return 'An unexpected error occurred.';
    }
  }

  String _handleBadResponse(Response? response) {
    if (response == null) {
      return 'No response from server.';
    }

    final statusCode = response.statusCode ?? 0;
    final data = response.data;

    // Try to extract error message from response body
    String? serverMessage;
    if (data is Map<String, dynamic>) {
      serverMessage = data['message'] as String? ??
          data['error'] as String? ??
          data['detail'] as String?;
    }

    if (serverMessage != null && serverMessage.isNotEmpty) {
      return serverMessage;
    }

    // Default messages based on status code
    switch (statusCode) {
      case 400:
        return 'Invalid request. Please check your input.';
      case 401:
        return 'Session expired. Please log in again.';
      case 403:
        return 'You don\'t have permission to perform this action.';
      case 404:
        return 'The requested resource was not found.';
      case 409:
        return 'This action conflicts with existing data.';
      case 422:
        return 'The submitted data is invalid.';
      case 429:
        return 'Too many requests. Please wait a moment.';
      case 500:
        return 'Server error. Please try again later.';
      case 502:
        return 'Service temporarily unavailable.';
      case 503:
        return 'Service is under maintenance.';
      default:
        return 'An error occurred (Error $statusCode).';
    }
  }
}

/// Extension for chaining repository results
extension ResultExtensions<T> on Result<T> {
  /// Transform successful data
  Result<R> map<R>(R Function(T data) transform) {
    return when(
      success: (data) => Success(transform(data)),
      failure: (message, code) => Failure(message, code: code),
    );
  }

  /// Chain another async operation on success
  Future<Result<R>> flatMap<R>(Future<Result<R>> Function(T data) transform) async {
    return when(
      success: (data) => transform(data),
      failure: (message, code) async => Failure(message, code: code),
    );
  }

  /// Get data or throw
  T getOrThrow() {
    return when(
      success: (data) => data,
      failure: (message, code) => throw ApiException(message, statusCode: code),
    );
  }

  /// Get data or default
  T getOrElse(T defaultValue) {
    return when(
      success: (data) => data,
      failure: (_, __) => defaultValue,
    );
  }
}
