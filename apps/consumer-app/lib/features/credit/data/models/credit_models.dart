import 'package:json_annotation/json_annotation.dart';

part 'credit_models.g.dart';

/// Credit score tier
enum CreditTier {
  @JsonValue('poor')
  poor,       // 300-549
  @JsonValue('fair')
  fair,       // 550-649
  @JsonValue('good')
  good,       // 650-749
  @JsonValue('excellent')
  excellent,  // 750-850
}

/// Loan status
enum LoanStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('approved')
  approved,
  @JsonValue('disbursed')
  disbursed,
  @JsonValue('active')
  active,
  @JsonValue('completed')
  completed,
  @JsonValue('defaulted')
  defaulted,
  @JsonValue('rejected')
  rejected,
}

/// Credit score model
@JsonSerializable()
class CreditScore {
  final String id;
  final String userId;
  final int score;
  final CreditTier tier;
  final int previousScore;
  final int change;
  final CreditFactors factors;
  final double maxLoanAmount;
  final DateTime calculatedAt;

  CreditScore({
    required this.id,
    required this.userId,
    required this.score,
    required this.tier,
    required this.previousScore,
    required this.change,
    required this.factors,
    required this.maxLoanAmount,
    required this.calculatedAt,
  });

  bool get isImproving => change > 0;
  bool get isDecreasing => change < 0;

  String get tierLabel {
    switch (tier) {
      case CreditTier.poor:
        return 'Poor';
      case CreditTier.fair:
        return 'Fair';
      case CreditTier.good:
        return 'Good';
      case CreditTier.excellent:
        return 'Excellent';
    }
  }

  factory CreditScore.fromJson(Map<String, dynamic> json) =>
      _$CreditScoreFromJson(json);
  Map<String, dynamic> toJson() => _$CreditScoreToJson(this);
}

/// Credit score factors breakdown
@JsonSerializable()
class CreditFactors {
  final int paymentHistory;       // 0-100
  final int savingsConsistency;   // 0-100
  final int gigPerformance;       // 0-100
  final int accountAge;           // 0-100
  final int walletActivity;       // 0-100
  final int communityTrust;       // 0-100

  CreditFactors({
    required this.paymentHistory,
    required this.savingsConsistency,
    required this.gigPerformance,
    required this.accountAge,
    required this.walletActivity,
    required this.communityTrust,
  });

  factory CreditFactors.fromJson(Map<String, dynamic> json) =>
      _$CreditFactorsFromJson(json);
  Map<String, dynamic> toJson() => _$CreditFactorsToJson(this);
}

/// Credit score history entry
@JsonSerializable()
class CreditScoreHistory {
  final String id;
  final String userId;
  final int score;
  final CreditTier tier;
  final DateTime recordedAt;

  CreditScoreHistory({
    required this.id,
    required this.userId,
    required this.score,
    required this.tier,
    required this.recordedAt,
  });

  factory CreditScoreHistory.fromJson(Map<String, dynamic> json) =>
      _$CreditScoreHistoryFromJson(json);
  Map<String, dynamic> toJson() => _$CreditScoreHistoryToJson(this);
}

/// Loan model
@JsonSerializable()
class Loan {
  final String id;
  final String userId;
  final double principalAmount;
  final double interestRate;
  final double totalAmount;
  final double amountPaid;
  final double amountRemaining;
  final int tenureDays;
  final LoanStatus status;
  final String purpose;
  final DateTime? disbursedAt;
  final DateTime dueDate;
  final DateTime? completedAt;
  final DateTime createdAt;
  final List<LoanRepayment>? repayments;

  Loan({
    required this.id,
    required this.userId,
    required this.principalAmount,
    required this.interestRate,
    required this.totalAmount,
    required this.amountPaid,
    required this.amountRemaining,
    required this.tenureDays,
    required this.status,
    required this.purpose,
    this.disbursedAt,
    required this.dueDate,
    this.completedAt,
    required this.createdAt,
    this.repayments,
  });

  double get progress => totalAmount > 0 ? amountPaid / totalAmount : 0;
  bool get isOverdue =>
      status == LoanStatus.active && DateTime.now().isAfter(dueDate);

  int get daysRemaining {
    if (status != LoanStatus.active) return 0;
    return dueDate.difference(DateTime.now()).inDays;
  }

  factory Loan.fromJson(Map<String, dynamic> json) => _$LoanFromJson(json);
  Map<String, dynamic> toJson() => _$LoanToJson(this);
}

/// Loan repayment record
@JsonSerializable()
class LoanRepayment {
  final String id;
  final String loanId;
  final double amount;
  final double principalPortion;
  final double interestPortion;
  final double balanceAfter;
  final String status;
  final DateTime dueDate;
  final DateTime? paidAt;
  final String? transactionId;
  final DateTime createdAt;

  LoanRepayment({
    required this.id,
    required this.loanId,
    required this.amount,
    required this.principalPortion,
    required this.interestPortion,
    required this.balanceAfter,
    required this.status,
    required this.dueDate,
    this.paidAt,
    this.transactionId,
    required this.createdAt,
  });

  bool get isPaid => paidAt != null;
  bool get isOverdue => !isPaid && DateTime.now().isAfter(dueDate);

  factory LoanRepayment.fromJson(Map<String, dynamic> json) =>
      _$LoanRepaymentFromJson(json);
  Map<String, dynamic> toJson() => _$LoanRepaymentToJson(this);
}

/// Loan product (available loan types)
@JsonSerializable()
class LoanProduct {
  final String id;
  final String name;
  final String description;
  final double minAmount;
  final double maxAmount;
  final double interestRate;
  final int minTenureDays;
  final int maxTenureDays;
  final int minCreditScore;
  final bool isActive;

  LoanProduct({
    required this.id,
    required this.name,
    required this.description,
    required this.minAmount,
    required this.maxAmount,
    required this.interestRate,
    required this.minTenureDays,
    required this.maxTenureDays,
    required this.minCreditScore,
    this.isActive = true,
  });

  factory LoanProduct.fromJson(Map<String, dynamic> json) =>
      _$LoanProductFromJson(json);
  Map<String, dynamic> toJson() => _$LoanProductToJson(this);
}

/// Loan application request
@JsonSerializable()
class LoanApplicationRequest {
  final String productId;
  final double amount;
  final int tenureDays;
  final String purpose;

  LoanApplicationRequest({
    required this.productId,
    required this.amount,
    required this.tenureDays,
    required this.purpose,
  });

  factory LoanApplicationRequest.fromJson(Map<String, dynamic> json) =>
      _$LoanApplicationRequestFromJson(json);
  Map<String, dynamic> toJson() => _$LoanApplicationRequestToJson(this);
}

/// Loan eligibility check response
@JsonSerializable()
class LoanEligibility {
  final bool isEligible;
  final double maxAmount;
  final double suggestedAmount;
  final double interestRate;
  final int maxTenureDays;
  final String? reason;
  final List<String>? suggestions;

  LoanEligibility({
    required this.isEligible,
    required this.maxAmount,
    required this.suggestedAmount,
    required this.interestRate,
    required this.maxTenureDays,
    this.reason,
    this.suggestions,
  });

  factory LoanEligibility.fromJson(Map<String, dynamic> json) =>
      _$LoanEligibilityFromJson(json);
  Map<String, dynamic> toJson() => _$LoanEligibilityToJson(this);
}

/// Repayment request
@JsonSerializable()
class RepaymentRequest {
  final String loanId;
  final double amount;
  final String pin;

  RepaymentRequest({
    required this.loanId,
    required this.amount,
    required this.pin,
  });

  factory RepaymentRequest.fromJson(Map<String, dynamic> json) =>
      _$RepaymentRequestFromJson(json);
  Map<String, dynamic> toJson() => _$RepaymentRequestToJson(this);
}
