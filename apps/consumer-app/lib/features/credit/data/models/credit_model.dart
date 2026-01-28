import 'package:freezed_annotation/freezed_annotation.dart';

part 'credit_model.freezed.dart';
part 'credit_model.g.dart';

enum CreditTier {
  @JsonValue('poor')
  poor,
  @JsonValue('fair')
  fair,
  @JsonValue('good')
  good,
  @JsonValue('excellent')
  excellent,
}

enum LoanStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('approved')
  approved,
  @JsonValue('rejected')
  rejected,
  @JsonValue('disbursed')
  disbursed,
  @JsonValue('active')
  active,
  @JsonValue('overdue')
  overdue,
  @JsonValue('defaulted')
  defaulted,
  @JsonValue('paid')
  paid,
  @JsonValue('cancelled')
  cancelled,
}

enum LoanPurpose {
  @JsonValue('personal')
  personal,
  @JsonValue('business')
  business,
  @JsonValue('education')
  education,
  @JsonValue('medical')
  medical,
  @JsonValue('emergency')
  emergency,
  @JsonValue('rent')
  rent,
  @JsonValue('other')
  other,
}

enum RepaymentFrequency {
  @JsonValue('weekly')
  weekly,
  @JsonValue('biweekly')
  biweekly,
  @JsonValue('monthly')
  monthly,
}

enum RepaymentStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('paid')
  paid,
  @JsonValue('overdue')
  overdue,
  @JsonValue('missed')
  missed,
  @JsonValue('partial')
  partial,
}

extension CreditTierX on CreditTier {
  String get displayName {
    switch (this) {
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

  int get minScore {
    switch (this) {
      case CreditTier.poor:
        return 300;
      case CreditTier.fair:
        return 500;
      case CreditTier.good:
        return 650;
      case CreditTier.excellent:
        return 750;
    }
  }

  int get maxScore {
    switch (this) {
      case CreditTier.poor:
        return 499;
      case CreditTier.fair:
        return 649;
      case CreditTier.good:
        return 749;
      case CreditTier.excellent:
        return 850;
    }
  }

  double get maxLoanMultiplier {
    switch (this) {
      case CreditTier.poor:
        return 0.5;
      case CreditTier.fair:
        return 1.0;
      case CreditTier.good:
        return 2.0;
      case CreditTier.excellent:
        return 3.0;
    }
  }
}

extension LoanStatusX on LoanStatus {
  String get displayName {
    switch (this) {
      case LoanStatus.pending:
        return 'Pending Review';
      case LoanStatus.approved:
        return 'Approved';
      case LoanStatus.rejected:
        return 'Rejected';
      case LoanStatus.disbursed:
        return 'Disbursed';
      case LoanStatus.active:
        return 'Active';
      case LoanStatus.overdue:
        return 'Overdue';
      case LoanStatus.defaulted:
        return 'Defaulted';
      case LoanStatus.paid:
        return 'Paid Off';
      case LoanStatus.cancelled:
        return 'Cancelled';
    }
  }
}

extension LoanPurposeX on LoanPurpose {
  String get displayName {
    switch (this) {
      case LoanPurpose.personal:
        return 'Personal';
      case LoanPurpose.business:
        return 'Business';
      case LoanPurpose.education:
        return 'Education';
      case LoanPurpose.medical:
        return 'Medical';
      case LoanPurpose.emergency:
        return 'Emergency';
      case LoanPurpose.rent:
        return 'Rent';
      case LoanPurpose.other:
        return 'Other';
    }
  }

  String get icon {
    switch (this) {
      case LoanPurpose.personal:
        return '👤';
      case LoanPurpose.business:
        return '💼';
      case LoanPurpose.education:
        return '📚';
      case LoanPurpose.medical:
        return '🏥';
      case LoanPurpose.emergency:
        return '🚨';
      case LoanPurpose.rent:
        return '🏠';
      case LoanPurpose.other:
        return '📋';
    }
  }
}

@freezed
class CreditScore with _$CreditScore {
  const factory CreditScore({
    required String userId,
    required int score,
    required CreditTier tier,
    @Default(0) int scoreDelta,
    required double maxLoanAmount,
    required CreditScoreFactors factors,
    @Default([]) List<CreditScoreHistory> history,
    DateTime? lastUpdated,
    DateTime? nextUpdate,
  }) = _CreditScore;

  const CreditScore._();

  factory CreditScore.fromJson(Map<String, dynamic> json) => _$CreditScoreFromJson(json);

  double get scoreProgress => (score - 300) / 550; // 300-850 range
  bool get isImproving => scoreDelta > 0;
  bool get isDeclining => scoreDelta < 0;
}

@freezed
class CreditScoreFactors with _$CreditScoreFactors {
  const factory CreditScoreFactors({
    @Default(0) int paymentHistory, // 0-100
    @Default(0) int savingsConsistency, // 0-100
    @Default(0) int gigPerformance, // 0-100
    @Default(0) int accountAge, // 0-100
    @Default(0) int walletActivity, // 0-100
    @Default(0) int communityTrust, // 0-100
  }) = _CreditScoreFactors;

  const CreditScoreFactors._();

  factory CreditScoreFactors.fromJson(Map<String, dynamic> json) => _$CreditScoreFactorsFromJson(json);

  List<CreditFactor> toList() {
    return [
      CreditFactor(
        name: 'Payment History',
        description: 'Your track record of paying bills and loans on time',
        score: paymentHistory,
        weight: 0.25,
        icon: '📅',
      ),
      CreditFactor(
        name: 'Savings Consistency',
        description: 'How regularly you contribute to savings circles',
        score: savingsConsistency,
        weight: 0.20,
        icon: '💰',
      ),
      CreditFactor(
        name: 'Gig Performance',
        description: 'Your ratings and completion rate on gigs',
        score: gigPerformance,
        weight: 0.20,
        icon: '⭐',
      ),
      CreditFactor(
        name: 'Account Age',
        description: 'How long you have been using HustleX',
        score: accountAge,
        weight: 0.10,
        icon: '📆',
      ),
      CreditFactor(
        name: 'Wallet Activity',
        description: 'Your transaction history and wallet usage',
        score: walletActivity,
        weight: 0.15,
        icon: '💳',
      ),
      CreditFactor(
        name: 'Community Trust',
        description: 'Trust signals from savings circle members',
        score: communityTrust,
        weight: 0.10,
        icon: '🤝',
      ),
    ];
  }
}

@freezed
class CreditFactor with _$CreditFactor {
  const factory CreditFactor({
    required String name,
    required String description,
    required int score,
    required double weight,
    required String icon,
  }) = _CreditFactor;

  const CreditFactor._();

  factory CreditFactor.fromJson(Map<String, dynamic> json) => _$CreditFactorFromJson(json);

  String get rating {
    if (score >= 80) return 'Excellent';
    if (score >= 60) return 'Good';
    if (score >= 40) return 'Fair';
    return 'Needs Work';
  }
}

@freezed
class CreditScoreHistory with _$CreditScoreHistory {
  const factory CreditScoreHistory({
    required int score,
    required DateTime date,
    int? delta,
    String? reason,
  }) = _CreditScoreHistory;

  factory CreditScoreHistory.fromJson(Map<String, dynamic> json) => _$CreditScoreHistoryFromJson(json);
}

@freezed
class Loan with _$Loan {
  const factory Loan({
    required String id,
    required String userId,
    required double principalAmount,
    required double interestRate,
    required double totalAmount,
    required double amountPaid,
    required int tenorMonths,
    required RepaymentFrequency repaymentFrequency,
    required LoanPurpose purpose,
    required LoanStatus status,
    String? purposeDescription,
    DateTime? applicationDate,
    DateTime? approvalDate,
    DateTime? disbursementDate,
    DateTime? dueDate,
    DateTime? nextPaymentDate,
    double? nextPaymentAmount,
    @Default(0) int paymentsMade,
    @Default(0) int paymentsTotal,
    @Default(0) int daysOverdue,
    @Default([]) List<LoanRepayment> repayments,
    String? rejectionReason,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Loan;

  const Loan._();

  factory Loan.fromJson(Map<String, dynamic> json) => _$LoanFromJson(json);

  double get outstandingBalance => totalAmount - amountPaid;
  double get repaymentProgress => totalAmount > 0 ? amountPaid / totalAmount : 0;
  bool get isActive => status == LoanStatus.active || status == LoanStatus.disbursed;
  bool get isOverdue => status == LoanStatus.overdue;
  bool get isPaidOff => status == LoanStatus.paid;

  String get formattedOutstanding => '₦${outstandingBalance.toStringAsFixed(2)}';
  String get formattedPrincipal => '₦${principalAmount.toStringAsFixed(2)}';
}

@freezed
class LoanRepayment with _$LoanRepayment {
  const factory LoanRepayment({
    required String id,
    required String loanId,
    required double amount,
    required double principalPortion,
    required double interestPortion,
    required int installmentNumber,
    required RepaymentStatus status,
    required DateTime dueDate,
    DateTime? paidAt,
    String? transactionId,
    double? lateFee,
  }) = _LoanRepayment;

  const LoanRepayment._();

  factory LoanRepayment.fromJson(Map<String, dynamic> json) => _$LoanRepaymentFromJson(json);

  bool get isPaid => status == RepaymentStatus.paid;
  bool get isOverdue => status == RepaymentStatus.overdue;
  bool get isPending => status == RepaymentStatus.pending;
}

@freezed
class LoanOffer with _$LoanOffer {
  const factory LoanOffer({
    required double minAmount,
    required double maxAmount,
    required double interestRate,
    required int minTenorMonths,
    required int maxTenorMonths,
    @Default([]) List<RepaymentFrequency> availableFrequencies,
    double? processingFee,
    String? terms,
  }) = _LoanOffer;

  const LoanOffer._();

  factory LoanOffer.fromJson(Map<String, dynamic> json) => _$LoanOfferFromJson(json);

  double calculateTotalAmount(double principal, int tenorMonths) {
    final monthlyRate = interestRate / 100 / 12;
    final totalInterest = principal * monthlyRate * tenorMonths;
    return principal + totalInterest + (processingFee ?? 0);
  }

  double calculateMonthlyPayment(double principal, int tenorMonths) {
    return calculateTotalAmount(principal, tenorMonths) / tenorMonths;
  }
}

@freezed
class LoanApplication with _$LoanApplication {
  const factory LoanApplication({
    required double amount,
    required int tenorMonths,
    required RepaymentFrequency repaymentFrequency,
    required LoanPurpose purpose,
    String? purposeDescription,
    String? employmentStatus,
    double? monthlyIncome,
    String? guarantorName,
    String? guarantorPhone,
  }) = _LoanApplication;

  factory LoanApplication.fromJson(Map<String, dynamic> json) => _$LoanApplicationFromJson(json);
}

@freezed
class MakeLoanPaymentRequest with _$MakeLoanPaymentRequest {
  const factory MakeLoanPaymentRequest({
    required String loanId,
    required double amount,
    required String pin,
    String? repaymentId,
  }) = _MakeLoanPaymentRequest;

  factory MakeLoanPaymentRequest.fromJson(Map<String, dynamic> json) => _$MakeLoanPaymentRequestFromJson(json);
}

@freezed
class CreditTip with _$CreditTip {
  const factory CreditTip({
    required String title,
    required String description,
    required String icon,
    int? potentialScoreIncrease,
  }) = _CreditTip;

  factory CreditTip.fromJson(Map<String, dynamic> json) => _$CreditTipFromJson(json);
}

@freezed
class LoanFilter with _$LoanFilter {
  const factory LoanFilter({
    LoanStatus? status,
    DateTime? startDate,
    DateTime? endDate,
    @Default(1) int page,
    @Default(20) int limit,
  }) = _LoanFilter;

  factory LoanFilter.fromJson(Map<String, dynamic> json) => _$LoanFilterFromJson(json);
}

@freezed
class PaginatedLoans with _$PaginatedLoans {
  const factory PaginatedLoans({
    required List<Loan> loans,
    required int total,
    required int page,
    required int limit,
    required bool hasMore,
  }) = _PaginatedLoans;

  factory PaginatedLoans.fromJson(Map<String, dynamic> json) => _$PaginatedLoansFromJson(json);
}
