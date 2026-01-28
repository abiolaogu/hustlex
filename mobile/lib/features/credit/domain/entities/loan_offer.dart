import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/value_objects/money.dart';
import 'loan.dart';

part 'loan_offer.freezed.dart';

/// Loan offer/product entity
@freezed
class LoanOffer with _$LoanOffer {
  const factory LoanOffer({
    required String id,
    required String name,
    String? description,
    required Money minAmount,
    required Money maxAmount,
    required double interestRate,
    required int minTenorMonths,
    required int maxTenorMonths,
    required int minCreditScore,
    @Default([]) List<RepaymentFrequency> availableFrequencies,
    Money? processingFee,
    String? terms,
    @Default(true) bool isActive,
  }) = _LoanOffer;

  const LoanOffer._();

  /// Calculate total amount for given principal and tenor
  Money calculateTotalAmount(Money principal, int tenorMonths) {
    final monthlyRate = interestRate / 100 / 12;
    final totalInterest = principal.amount * monthlyRate * tenorMonths;
    final total = principal.amount + totalInterest + (processingFee?.amount ?? 0);
    return Money.fromMajorUnits(amount: total, currency: principal.currency);
  }

  /// Calculate monthly payment
  Money calculateMonthlyPayment(Money principal, int tenorMonths) {
    final total = calculateTotalAmount(principal, tenorMonths);
    return Money.fromMajorUnits(
      amount: total.amount / tenorMonths,
      currency: principal.currency,
    );
  }

  /// Calculate interest amount
  Money calculateInterest(Money principal, int tenorMonths) {
    final monthlyRate = interestRate / 100 / 12;
    final interest = principal.amount * monthlyRate * tenorMonths;
    return Money.fromMajorUnits(amount: interest, currency: principal.currency);
  }
}

/// Loan eligibility result
@freezed
class LoanEligibility with _$LoanEligibility {
  const factory LoanEligibility({
    required bool isEligible,
    required Money maxAmount,
    required Money suggestedAmount,
    required double interestRate,
    required int maxTenorMonths,
    String? reason,
    @Default([]) List<String> suggestions,
    @Default([]) List<LoanOffer> availableOffers,
  }) = _LoanEligibility;

  const LoanEligibility._();

  /// Check if user can borrow
  bool get canBorrow => isEligible && maxAmount.amount > 0;
}
