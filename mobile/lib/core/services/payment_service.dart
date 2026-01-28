import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_paystack/flutter_paystack.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:logger/logger.dart';

import '../constants/app_constants.dart';
import '../storage/secure_storage.dart';

/// Payment result
class PaymentResult {
  final bool success;
  final String? reference;
  final String? message;
  final Map<String, dynamic>? metadata;

  PaymentResult({
    required this.success,
    this.reference,
    this.message,
    this.metadata,
  });

  factory PaymentResult.success({
    required String reference,
    Map<String, dynamic>? metadata,
  }) {
    return PaymentResult(
      success: true,
      reference: reference,
      metadata: metadata,
    );
  }

  factory PaymentResult.failure(String message) {
    return PaymentResult(
      success: false,
      message: message,
    );
  }

  factory PaymentResult.cancelled() {
    return PaymentResult(
      success: false,
      message: 'Payment was cancelled',
    );
  }
}

/// Payment service using Paystack
class PaymentService {
  final Logger _logger = Logger();
  final PaystackPlugin _paystackPlugin;
  final SecureStorage _secureStorage;

  bool _isInitialized = false;

  PaymentService({
    required SecureStorage secureStorage,
    PaystackPlugin? paystackPlugin,
  })  : _secureStorage = secureStorage,
        _paystackPlugin = paystackPlugin ?? PaystackPlugin();

  /// Initialize Paystack plugin
  Future<void> initialize() async {
    if (_isInitialized) return;

    try {
      final publicKey = AppConstants.paystackPublicKey;
      await _paystackPlugin.initialize(publicKey: publicKey);
      _isInitialized = true;
      _logger.i('Paystack initialized');
    } catch (e) {
      _logger.e('Failed to initialize Paystack', error: e);
      rethrow;
    }
  }

  /// Ensure plugin is initialized
  Future<void> _ensureInitialized() async {
    if (!_isInitialized) {
      await initialize();
    }
  }

  /// Process payment with Paystack checkout
  /// 
  /// [context] - BuildContext for showing checkout UI
  /// [email] - Customer's email address
  /// [amount] - Amount in Naira (will be converted to Kobo)
  /// [reference] - Unique transaction reference from backend
  /// [accessCode] - Access code from backend initialization
  Future<PaymentResult> processPayment({
    required BuildContext context,
    required String email,
    required double amount,
    required String reference,
    required String accessCode,
    Map<String, dynamic>? metadata,
  }) async {
    await _ensureInitialized();

    try {
      // Convert amount to kobo (Paystack uses minor currency units)
      final amountInKobo = (amount * 100).round();

      final charge = Charge()
        ..amount = amountInKobo
        ..email = email
        ..reference = reference
        ..accessCode = accessCode;

      // Add metadata if provided
      if (metadata != null) {
        charge.putMetaData('custom_fields', metadata);
      }

      final response = await _paystackPlugin.checkout(
        context,
        method: CheckoutMethod.selectable,
        charge: charge,
        fullscreen: true,
        logo: const AssetImage('assets/images/logo.png'),
      );

      if (response.status) {
        _logger.i('Payment successful: ${response.reference}');
        return PaymentResult.success(
          reference: response.reference ?? reference,
          metadata: metadata,
        );
      } else {
        if (response.message.contains('cancelled') ||
            response.message.contains('Cancelled')) {
          _logger.w('Payment cancelled');
          return PaymentResult.cancelled();
        }
        _logger.e('Payment failed: ${response.message}');
        return PaymentResult.failure(response.message);
      }
    } catch (e) {
      _logger.e('Payment error', error: e);
      return PaymentResult.failure('Payment processing failed: $e');
    }
  }

  /// Process card payment directly
  Future<PaymentResult> processCardPayment({
    required BuildContext context,
    required String email,
    required double amount,
    required String reference,
    required String accessCode,
    PaymentCard? card,
  }) async {
    await _ensureInitialized();

    try {
      final amountInKobo = (amount * 100).round();

      final charge = Charge()
        ..amount = amountInKobo
        ..email = email
        ..reference = reference
        ..accessCode = accessCode;

      if (card != null) {
        charge.card = card;
      }

      final response = await _paystackPlugin.checkout(
        context,
        method: CheckoutMethod.card,
        charge: charge,
        fullscreen: true,
        hideEmail: true,
      );

      if (response.status) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
        );
      } else {
        if (response.message.contains('cancelled')) {
          return PaymentResult.cancelled();
        }
        return PaymentResult.failure(response.message);
      }
    } catch (e) {
      _logger.e('Card payment error', error: e);
      return PaymentResult.failure('Card payment failed: $e');
    }
  }

  /// Process bank transfer payment
  Future<PaymentResult> processBankTransfer({
    required BuildContext context,
    required String email,
    required double amount,
    required String reference,
    required String accessCode,
  }) async {
    await _ensureInitialized();

    try {
      final amountInKobo = (amount * 100).round();

      final charge = Charge()
        ..amount = amountInKobo
        ..email = email
        ..reference = reference
        ..accessCode = accessCode;

      final response = await _paystackPlugin.checkout(
        context,
        method: CheckoutMethod.bank,
        charge: charge,
        fullscreen: true,
      );

      if (response.status) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
        );
      } else {
        if (response.message.contains('cancelled')) {
          return PaymentResult.cancelled();
        }
        return PaymentResult.failure(response.message);
      }
    } catch (e) {
      _logger.e('Bank transfer error', error: e);
      return PaymentResult.failure('Bank transfer failed: $e');
    }
  }

  /// Validate card number
  bool isCardNumberValid(String cardNumber) {
    return _paystackPlugin.validateCardNumber(cardNumber);
  }

  /// Validate CVC
  bool isCvcValid(String cvc) {
    return _paystackPlugin.validateCvc(cvc);
  }

  /// Validate expiry date
  bool isExpiryValid(int month, int year) {
    return !_paystackPlugin.isCardExpired(month, year);
  }

  /// Get card type from card number
  String getCardType(String cardNumber) {
    final card = PaymentCard(number: cardNumber);
    return card.type?.toString().split('.').last ?? 'Unknown';
  }

  /// Create payment card from details
  PaymentCard createCard({
    required String number,
    required int expiryMonth,
    required int expiryYear,
    required String cvc,
    String? name,
  }) {
    return PaymentCard(
      number: number,
      expiryMonth: expiryMonth,
      expiryYear: expiryYear,
      cvc: cvc,
      name: name,
    );
  }

  /// Check if card is valid
  bool isCardValid(PaymentCard card) {
    return card.isValid();
  }
}

/// Payment state for UI
class PaymentState {
  final bool isProcessing;
  final PaymentResult? result;
  final String? error;

  const PaymentState({
    this.isProcessing = false,
    this.result,
    this.error,
  });

  PaymentState copyWith({
    bool? isProcessing,
    PaymentResult? result,
    String? error,
  }) {
    return PaymentState(
      isProcessing: isProcessing ?? this.isProcessing,
      result: result,
      error: error,
    );
  }
}

/// Payment state notifier
class PaymentNotifier extends StateNotifier<PaymentState> {
  final PaymentService _paymentService;

  PaymentNotifier(this._paymentService) : super(const PaymentState());

  /// Process deposit payment
  Future<PaymentResult> processDeposit({
    required BuildContext context,
    required String email,
    required double amount,
    required String reference,
    required String accessCode,
  }) async {
    state = state.copyWith(isProcessing: true, error: null, result: null);

    final result = await _paymentService.processPayment(
      context: context,
      email: email,
      amount: amount,
      reference: reference,
      accessCode: accessCode,
      metadata: {'type': 'deposit'},
    );

    state = state.copyWith(isProcessing: false, result: result);

    if (!result.success && result.message != null) {
      state = state.copyWith(error: result.message);
    }

    return result;
  }

  /// Reset state
  void reset() {
    state = const PaymentState();
  }
}

/// Payment service provider
final paymentServiceProvider = Provider<PaymentService>((ref) {
  final secureStorage = ref.watch(secureStorageProvider);
  return PaymentService(secureStorage: secureStorage);
});

/// Secure storage provider (duplicate for standalone usage)
final secureStorageProvider = Provider<SecureStorage>((ref) {
  return SecureStorage();
});

/// Payment state provider
final paymentProvider =
    StateNotifierProvider<PaymentNotifier, PaymentState>((ref) {
  final service = ref.watch(paymentServiceProvider);
  return PaymentNotifier(service);
});

/// Card validation providers
final cardNumberValidProvider =
    Provider.family<bool, String>((ref, cardNumber) {
  final service = ref.watch(paymentServiceProvider);
  return service.isCardNumberValid(cardNumber);
});

final cardCvcValidProvider = Provider.family<bool, String>((ref, cvc) {
  final service = ref.watch(paymentServiceProvider);
  return service.isCvcValid(cvc);
});

final cardTypeProvider = Provider.family<String, String>((ref, cardNumber) {
  final service = ref.watch(paymentServiceProvider);
  return service.getCardType(cardNumber);
});
