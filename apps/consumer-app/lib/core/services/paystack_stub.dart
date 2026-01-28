/// Stub classes to replace flutter_paystack while it has dependency issues
/// TODO: Replace with actual flutter_paystack or pay_with_paystack when deps are resolved

import 'package:flutter/material.dart';

enum CardType { unknown, visa, masterCard, verve, americanExpress, discover }

class PaymentCard {
  final String? number;
  final int? expiryMonth;
  final int? expiryYear;
  final String? cvc;
  final String? name;
  final CardType? type;

  PaymentCard({
    this.number,
    this.expiryMonth,
    this.expiryYear,
    this.cvc,
    this.name,
  }) : type = _detectCardType(number);

  static CardType _detectCardType(String? number) {
    if (number == null || number.isEmpty) return CardType.unknown;
    final cleaned = number.replaceAll(RegExp(r'\D'), '');
    if (cleaned.startsWith('4')) return CardType.visa;
    if (cleaned.startsWith('5')) return CardType.masterCard;
    if (cleaned.startsWith('506')) return CardType.verve;
    return CardType.unknown;
  }

  bool isValid() {
    return number != null &&
        number!.length >= 13 &&
        expiryMonth != null &&
        expiryYear != null &&
        cvc != null;
  }
}

class Charge {
  int amount = 0;
  String? email;
  String? reference;
  PaymentCard? card;
  String? accessCode;
  String? currency;
  Map<String, dynamic>? metadata;

  Charge();

  void putMetaData(String key, dynamic value) {
    metadata ??= {};
    metadata![key] = value;
  }
}

enum CheckoutMethod { card, bank, selectable }

class CheckoutResponse {
  final bool status;
  final bool verify;
  final String? reference;
  final String? message;
  final CheckoutMethod method;

  CheckoutResponse({
    this.status = false,
    this.verify = false,
    this.reference,
    this.message,
    this.method = CheckoutMethod.selectable,
  });
}

class PaystackException implements Exception {
  final String message;
  PaystackException(this.message);

  @override
  String toString() => message;
}

class PaystackPlugin {
  bool _initialized = false;

  Future<void> initialize({required String publicKey}) async {
    _initialized = true;
  }

  bool get isInitialized => _initialized;

  Future<CheckoutResponse> checkout(
    BuildContext context, {
    required Charge charge,
    bool fullscreen = false,
    CheckoutMethod method = CheckoutMethod.selectable,
    Widget? logo,
    bool hideEmail = false,
    bool hideAmount = false,
  }) async {
    // Stub implementation - in production, use webview or actual SDK
    return CheckoutResponse(
      status: false,
      message: 'Paystack stub - implement with webview checkout',
    );
  }

  bool isCardExpired(int month, int year) {
    final now = DateTime.now();
    final cardExpiry = DateTime(year < 100 ? 2000 + year : year, month + 1, 0);
    return cardExpiry.isBefore(now);
  }

  bool validateCardNumber(String cardNumber) {
    final cleaned = cardNumber.replaceAll(RegExp(r'\D'), '');
    if (cleaned.length < 13 || cleaned.length > 19) return false;
    // Luhn algorithm
    int sum = 0;
    bool alternate = false;
    for (int i = cleaned.length - 1; i >= 0; i--) {
      int digit = int.parse(cleaned[i]);
      if (alternate) {
        digit *= 2;
        if (digit > 9) digit -= 9;
      }
      sum += digit;
      alternate = !alternate;
    }
    return sum % 10 == 0;
  }

  bool validateCvc(String cvc) {
    final cleaned = cvc.replaceAll(RegExp(r'\D'), '');
    return cleaned.length >= 3 && cleaned.length <= 4;
  }
}
