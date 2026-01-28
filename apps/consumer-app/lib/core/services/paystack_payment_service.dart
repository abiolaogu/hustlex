import 'dart:async';
import 'package:flutter/material.dart';
import 'paystack_stub.dart'; // Stub for flutter_paystack
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../config/env_config.dart';

/// Payment result
class PaymentResult {
  final bool success;
  final String? reference;
  final String? message;
  final Map<String, dynamic>? metadata;

  const PaymentResult({
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

  factory PaymentResult.failed({String? message}) {
    return PaymentResult(
      success: false,
      message: message ?? 'Payment failed',
    );
  }

  factory PaymentResult.cancelled() {
    return const PaymentResult(
      success: false,
      message: 'Payment cancelled by user',
    );
  }
}

/// Card details for saved cards
class SavedCard {
  final String id;
  final String last4;
  final String brand;
  final int expMonth;
  final int expYear;
  final String? bank;
  final bool isDefault;

  const SavedCard({
    required this.id,
    required this.last4,
    required this.brand,
    required this.expMonth,
    required this.expYear,
    this.bank,
    this.isDefault = false,
  });

  String get display => '$brand •••• $last4';
  String get expiry => '${expMonth.toString().padLeft(2, '0')}/$expYear';

  factory SavedCard.fromJson(Map<String, dynamic> json) {
    return SavedCard(
      id: json['id'] as String,
      last4: json['last4'] as String,
      brand: json['brand'] as String,
      expMonth: json['exp_month'] as int,
      expYear: json['exp_year'] as int,
      bank: json['bank'] as String?,
      isDefault: json['is_default'] as bool? ?? false,
    );
  }
}

/// Paystack payment service
class PaystackPaymentService {
  final PaystackPlugin _plugin;
  bool _isInitialized = false;

  PaystackPaymentService() : _plugin = PaystackPlugin();

  /// Initialize Paystack plugin
  Future<void> initialize() async {
    if (_isInitialized) return;

    final publicKey = EnvConfig.paystackPublicKey;
    if (publicKey.isEmpty) {
      throw Exception('Paystack public key not configured');
    }

    await _plugin.initialize(publicKey: publicKey);
    _isInitialized = true;
  }

  /// Ensure plugin is initialized
  Future<void> _ensureInitialized() async {
    if (!_isInitialized) {
      await initialize();
    }
  }

  /// Charge card using Paystack checkout
  /// 
  /// [context] - Build context for showing UI
  /// [email] - Customer's email
  /// [amountInKobo] - Amount in kobo (NGN * 100)
  /// [reference] - Unique transaction reference from backend
  /// [accessCode] - Access code from backend initialization
  Future<PaymentResult> chargeCard({
    required BuildContext context,
    required String email,
    required int amountInKobo,
    required String reference,
    String? accessCode,
    Map<String, dynamic>? metadata,
  }) async {
    await _ensureInitialized();

    final charge = Charge()
      ..amount = amountInKobo
      ..email = email
      ..reference = reference
      ..currency = 'NGN';

    if (accessCode != null) {
      charge.accessCode = accessCode;
    }

    if (metadata != null) {
      charge.putMetaData('custom_fields', metadata);
    }

    try {
      final response = await _plugin.checkout(
        context,
        method: CheckoutMethod.card,
        charge: charge,
        fullscreen: false,
        logo: const FlutterLogo(), // Replace with app logo widget
      );

      if (response.status) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
          metadata: {'message': response.message},
        );
      } else if (response.verify) {
        // Transaction needs verification
        return PaymentResult.success(
          reference: response.reference ?? reference,
          metadata: {'needs_verification': true, 'message': response.message},
        );
      } else {
        return PaymentResult.failed(message: response.message);
      }
    } on PaystackException catch (e) {
      return PaymentResult.failed(message: e.message);
    } catch (e) {
      return PaymentResult.failed(message: e.toString());
    }
  }

  /// Charge using bank transfer
  Future<PaymentResult> chargeBank({
    required BuildContext context,
    required String email,
    required int amountInKobo,
    required String reference,
    String? accessCode,
  }) async {
    await _ensureInitialized();

    final charge = Charge()
      ..amount = amountInKobo
      ..email = email
      ..reference = reference
      ..currency = 'NGN';

    if (accessCode != null) {
      charge.accessCode = accessCode;
    }

    try {
      final response = await _plugin.checkout(
        context,
        method: CheckoutMethod.bank,
        charge: charge,
        fullscreen: false,
      );

      if (response.status) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
        );
      } else if (response.verify) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
          metadata: {'needs_verification': true},
        );
      } else {
        return PaymentResult.failed(message: response.message);
      }
    } on PaystackException catch (e) {
      return PaymentResult.failed(message: e.message);
    } catch (e) {
      return PaymentResult.failed(message: e.toString());
    }
  }

  /// Charge using USSD
  Future<PaymentResult> chargeUssd({
    required BuildContext context,
    required String email,
    required int amountInKobo,
    required String reference,
    String? accessCode,
  }) async {
    await _ensureInitialized();

    final charge = Charge()
      ..amount = amountInKobo
      ..email = email
      ..reference = reference
      ..currency = 'NGN';

    if (accessCode != null) {
      charge.accessCode = accessCode;
    }

    try {
      final response = await _plugin.checkout(
        context,
        method: CheckoutMethod.selectable,
        charge: charge,
        fullscreen: false,
      );

      if (response.status) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
        );
      } else if (response.verify) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
          metadata: {'needs_verification': true},
        );
      } else {
        return PaymentResult.failed(message: response.message);
      }
    } on PaystackException catch (e) {
      return PaymentResult.failed(message: e.message);
    } catch (e) {
      return PaymentResult.failed(message: e.toString());
    }
  }

  /// Allow user to select payment method
  Future<PaymentResult> chargeWithMethodSelection({
    required BuildContext context,
    required String email,
    required int amountInKobo,
    required String reference,
    String? accessCode,
    Map<String, dynamic>? metadata,
  }) async {
    await _ensureInitialized();

    final charge = Charge()
      ..amount = amountInKobo
      ..email = email
      ..reference = reference
      ..currency = 'NGN';

    if (accessCode != null) {
      charge.accessCode = accessCode;
    }

    if (metadata != null) {
      charge.putMetaData('custom_fields', metadata);
    }

    try {
      final response = await _plugin.checkout(
        context,
        method: CheckoutMethod.selectable,
        charge: charge,
        fullscreen: true,
      );

      if (response.status) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
          metadata: {'method': response.method.name},
        );
      } else if (response.verify) {
        return PaymentResult.success(
          reference: response.reference ?? reference,
          metadata: {'needs_verification': true, 'method': response.method.name},
        );
      } else {
        return PaymentResult.failed(message: response.message);
      }
    } on PaystackException catch (e) {
      return PaymentResult.failed(message: e.message);
    } catch (e) {
      return PaymentResult.failed(message: e.toString());
    }
  }

  /// Get supported banks
  Future<List<Map<String, dynamic>>> getSupportedBanks() async {
    // Banks are typically fetched from backend API
    // This is a fallback list of major Nigerian banks
    return [
      {'code': '044', 'name': 'Access Bank'},
      {'code': '023', 'name': 'Citibank Nigeria'},
      {'code': '063', 'name': 'Diamond Bank'},
      {'code': '050', 'name': 'Ecobank Nigeria'},
      {'code': '084', 'name': 'Enterprise Bank'},
      {'code': '070', 'name': 'Fidelity Bank'},
      {'code': '011', 'name': 'First Bank of Nigeria'},
      {'code': '214', 'name': 'First City Monument Bank'},
      {'code': '058', 'name': 'Guaranty Trust Bank'},
      {'code': '030', 'name': 'Heritage Bank'},
      {'code': '301', 'name': 'Jaiz Bank'},
      {'code': '082', 'name': 'Keystone Bank'},
      {'code': '076', 'name': 'Polaris Bank'},
      {'code': '101', 'name': 'Providus Bank'},
      {'code': '221', 'name': 'Stanbic IBTC Bank'},
      {'code': '068', 'name': 'Standard Chartered Bank'},
      {'code': '232', 'name': 'Sterling Bank'},
      {'code': '100', 'name': 'Suntrust Bank'},
      {'code': '032', 'name': 'Union Bank of Nigeria'},
      {'code': '033', 'name': 'United Bank For Africa'},
      {'code': '215', 'name': 'Unity Bank'},
      {'code': '035', 'name': 'Wema Bank'},
      {'code': '057', 'name': 'Zenith Bank'},
    ];
  }

  /// Convert naira to kobo
  static int nairaToKobo(double naira) {
    return (naira * 100).round();
  }

  /// Convert kobo to naira
  static double koboToNaira(int kobo) {
    return kobo / 100;
  }

  /// Format amount for display
  static String formatAmount(int kobo) {
    final naira = koboToNaira(kobo);
    return '₦${naira.toStringAsFixed(2)}';
  }
}

/// Paystack payment service provider
final paystackPaymentServiceProvider = Provider<PaystackPaymentService>((ref) {
  return PaystackPaymentService();
});

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
      result: result ?? this.result,
      error: error,
    );
  }

  bool get isSuccess => result?.success ?? false;
  bool get isFailed => result != null && !result!.success;
}

/// Payment state notifier
class PaymentNotifier extends StateNotifier<PaymentState> {
  final PaystackPaymentService _paymentService;

  PaymentNotifier(this._paymentService) : super(const PaymentState());

  /// Process card payment
  Future<PaymentResult> processCardPayment({
    required BuildContext context,
    required String email,
    required double amountInNaira,
    required String reference,
    String? accessCode,
    Map<String, dynamic>? metadata,
  }) async {
    state = state.copyWith(isProcessing: true, error: null);

    try {
      final result = await _paymentService.chargeCard(
        context: context,
        email: email,
        amountInKobo: PaystackPaymentService.nairaToKobo(amountInNaira),
        reference: reference,
        accessCode: accessCode,
        metadata: metadata,
      );

      state = state.copyWith(
        isProcessing: false,
        result: result,
        error: result.success ? null : result.message,
      );

      return result;
    } catch (e) {
      final result = PaymentResult.failed(message: e.toString());
      state = state.copyWith(
        isProcessing: false,
        result: result,
        error: e.toString(),
      );
      return result;
    }
  }

  /// Process payment with method selection
  Future<PaymentResult> processPayment({
    required BuildContext context,
    required String email,
    required double amountInNaira,
    required String reference,
    String? accessCode,
    Map<String, dynamic>? metadata,
  }) async {
    state = state.copyWith(isProcessing: true, error: null);

    try {
      final result = await _paymentService.chargeWithMethodSelection(
        context: context,
        email: email,
        amountInKobo: PaystackPaymentService.nairaToKobo(amountInNaira),
        reference: reference,
        accessCode: accessCode,
        metadata: metadata,
      );

      state = state.copyWith(
        isProcessing: false,
        result: result,
        error: result.success ? null : result.message,
      );

      return result;
    } catch (e) {
      final result = PaymentResult.failed(message: e.toString());
      state = state.copyWith(
        isProcessing: false,
        result: result,
        error: e.toString(),
      );
      return result;
    }
  }

  /// Reset state
  void reset() {
    state = const PaymentState();
  }
}

/// Payment notifier provider
final paymentProvider = StateNotifierProvider<PaymentNotifier, PaymentState>((ref) {
  final service = ref.watch(paystackPaymentServiceProvider);
  return PaymentNotifier(service);
});
