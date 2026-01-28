import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import '../theme/app_colors.dart';

/// Widget to display monetary amounts with proper formatting
class MoneyDisplay extends StatelessWidget {
  const MoneyDisplay({
    super.key,
    required this.amount,
    this.currency = 'NGN',
    this.style,
    this.showCurrency = true,
    this.compact = false,
    this.colorize = false,
  });

  final double amount;
  final String currency;
  final TextStyle? style;
  final bool showCurrency;
  final bool compact;
  final bool colorize;

  String get _currencySymbol {
    switch (currency) {
      case 'NGN':
        return '₦';
      case 'USD':
        return '\$';
      case 'EUR':
        return '€';
      case 'GBP':
        return '£';
      default:
        return currency;
    }
  }

  String get _formattedAmount {
    if (compact && amount.abs() >= 1000000) {
      return '${(amount / 1000000).toStringAsFixed(1)}M';
    } else if (compact && amount.abs() >= 1000) {
      return '${(amount / 1000).toStringAsFixed(1)}K';
    }
    return NumberFormat('#,##0.00').format(amount);
  }

  @override
  Widget build(BuildContext context) {
    Color? textColor;
    if (colorize) {
      textColor = amount > 0
          ? AppColors.success
          : amount < 0
              ? AppColors.error
              : null;
    }

    final effectiveStyle = (style ?? Theme.of(context).textTheme.titleLarge)
        ?.copyWith(color: textColor);

    return Text(
      showCurrency ? '$_currencySymbol$_formattedAmount' : _formattedAmount,
      style: effectiveStyle,
    );
  }
}

/// Large balance display widget
class BalanceDisplay extends StatelessWidget {
  const BalanceDisplay({
    super.key,
    required this.balance,
    this.label = 'Available Balance',
    this.currency = 'NGN',
    this.isHidden = false,
  });

  final double balance;
  final String label;
  final String currency;
  final bool isHidden;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall,
        ),
        const SizedBox(height: 4),
        isHidden
            ? Text(
                '••••••',
                style: Theme.of(context).textTheme.headlineMedium,
              )
            : MoneyDisplay(
                amount: balance,
                currency: currency,
                style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
              ),
      ],
    );
  }
}
