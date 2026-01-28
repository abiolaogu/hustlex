import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import '../failures/value_failures.dart';
import 'value_object.dart';

part 'money.freezed.dart';

/// Supported currencies
@freezed
class Currency with _$Currency {
  const Currency._();

  const factory Currency.ngn() = _NGN;
  const factory Currency.usd() = _USD;
  const factory Currency.ghs() = _GHS;
  const factory Currency.gbp() = _GBP;

  factory Currency.fromCode(String code) {
    switch (code.toUpperCase()) {
      case 'NGN':
        return const Currency.ngn();
      case 'USD':
        return const Currency.usd();
      case 'GHS':
        return const Currency.ghs();
      case 'GBP':
        return const Currency.gbp();
      default:
        return const Currency.ngn();
    }
  }

  String get code => when(
        ngn: () => 'NGN',
        usd: () => 'USD',
        ghs: () => 'GHS',
        gbp: () => 'GBP',
      );

  String get symbol => when(
        ngn: () => '₦',
        usd: () => '\$',
        ghs: () => 'GH₵',
        gbp: () => '£',
      );

  String get name => when(
        ngn: () => 'Nigerian Naira',
        usd: () => 'US Dollar',
        ghs: () => 'Ghanaian Cedi',
        gbp: () => 'British Pound',
      );

  int get minorUnitDigits => 2;
}

/// Money value object - immutable, validated monetary value.
/// Amount is stored in minor units (kobo for NGN, cents for USD).
class Money extends ValueObject<int> {
  @override
  final Either<ValueFailure<int>, int> value;
  final Currency currency;

  /// Creates Money from minor units (e.g., kobo, cents)
  factory Money({
    required int amountInMinorUnits,
    Currency currency = const Currency.ngn(),
  }) {
    return Money._(
      _validateAmount(amountInMinorUnits),
      currency,
    );
  }

  /// Creates Money from major units (e.g., Naira, Dollars)
  factory Money.fromMajorUnits({
    required double amount,
    Currency currency = const Currency.ngn(),
  }) {
    final minorUnits = (amount * 100).round();
    return Money(amountInMinorUnits: minorUnits, currency: currency);
  }

  /// Creates Money representing zero
  factory Money.zero({Currency currency = const Currency.ngn()}) {
    return Money(amountInMinorUnits: 0, currency: currency);
  }

  const Money._(this.value, this.currency);

  static Either<ValueFailure<int>, int> _validateAmount(int amount) {
    if (amount < 0) {
      return left(ValueFailure.negativeAmount(failedValue: amount));
    }
    return right(amount);
  }

  /// Returns amount in minor units (e.g., Kobo)
  int get amountInMinorUnits => getOrCrash();

  /// Returns amount in major units (e.g., Naira)
  double get amountInMajorUnits => getOrCrash() / 100;

  /// Alias for amountInMajorUnits for convenience
  double get amount => amountInMajorUnits;

  /// Formatted string with currency symbol (e.g., "₦1,000.00")
  String get formatted {
    final amount = amountInMajorUnits;
    final parts = amount.toStringAsFixed(2).split('.');
    final intPart = parts[0];
    final decPart = parts[1];

    // Add thousand separators
    final buffer = StringBuffer();
    for (int i = 0; i < intPart.length; i++) {
      if (i > 0 && (intPart.length - i) % 3 == 0) {
        buffer.write(',');
      }
      buffer.write(intPart[i]);
    }

    return '${currency.symbol}$buffer.$decPart';
  }

  /// Formatted string without decimal places for whole amounts
  String get formattedCompact {
    if (amountInMinorUnits % 100 == 0) {
      final amount = (amountInMinorUnits / 100).toInt();
      final formatted = _formatWithCommas(amount);
      return '${currency.symbol}$formatted';
    }
    return formatted;
  }

  String _formatWithCommas(int number) {
    final str = number.toString();
    final buffer = StringBuffer();
    for (int i = 0; i < str.length; i++) {
      if (i > 0 && (str.length - i) % 3 == 0) {
        buffer.write(',');
      }
      buffer.write(str[i]);
    }
    return buffer.toString();
  }

  /// Add two Money values (must have same currency)
  Money operator +(Money other) {
    assert(currency == other.currency, 'Currency mismatch');
    return Money(
      amountInMinorUnits: getOrCrash() + other.getOrCrash(),
      currency: currency,
    );
  }

  /// Subtract Money (must have same currency)
  Money operator -(Money other) {
    assert(currency == other.currency, 'Currency mismatch');
    return Money(
      amountInMinorUnits: getOrCrash() - other.getOrCrash(),
      currency: currency,
    );
  }

  /// Multiply by a factor
  Money operator *(num factor) {
    return Money(
      amountInMinorUnits: (getOrCrash() * factor).round(),
      currency: currency,
    );
  }

  /// Compare amounts
  bool operator >(Money other) {
    assert(currency == other.currency, 'Currency mismatch');
    return getOrCrash() > other.getOrCrash();
  }

  bool operator <(Money other) {
    assert(currency == other.currency, 'Currency mismatch');
    return getOrCrash() < other.getOrCrash();
  }

  bool operator >=(Money other) {
    assert(currency == other.currency, 'Currency mismatch');
    return getOrCrash() >= other.getOrCrash();
  }

  bool operator <=(Money other) {
    assert(currency == other.currency, 'Currency mismatch');
    return getOrCrash() <= other.getOrCrash();
  }

  bool get isZero => getOrCrash() == 0;
  bool get isPositive => getOrCrash() > 0;

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Money &&
        other.value == value &&
        other.currency == currency;
  }

  @override
  int get hashCode => value.hashCode ^ currency.hashCode;

  @override
  String toString() => formatted;
}
