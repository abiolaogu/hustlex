import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';
import 'loan_repayment.dart';

part 'loan.freezed.dart';

/// Loan status
enum LoanStatus {
  pending,
  approved,
  rejected,
  disbursed,
  active,
  overdue,
  defaulted,
  paid,
  cancelled,
}

/// Extension for LoanStatus
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

  bool get isTerminal =>
      this == LoanStatus.paid ||
      this == LoanStatus.rejected ||
      this == LoanStatus.cancelled ||
      this == LoanStatus.defaulted;

  bool get requiresPayment =>
      this == LoanStatus.active ||
      this == LoanStatus.overdue ||
      this == LoanStatus.disbursed;
}

/// Loan purpose
enum LoanPurpose {
  personal,
  business,
  education,
  medical,
  emergency,
  rent,
  other,
}

/// Extension for LoanPurpose
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
}

/// Repayment frequency
enum RepaymentFrequency {
  weekly,
  biweekly,
  monthly,
}

/// Extension for RepaymentFrequency
extension RepaymentFrequencyX on RepaymentFrequency {
  String get displayName {
    switch (this) {
      case RepaymentFrequency.weekly:
        return 'Weekly';
      case RepaymentFrequency.biweekly:
        return 'Bi-weekly';
      case RepaymentFrequency.monthly:
        return 'Monthly';
    }
  }

  int get daysInterval {
    switch (this) {
      case RepaymentFrequency.weekly:
        return 7;
      case RepaymentFrequency.biweekly:
        return 14;
      case RepaymentFrequency.monthly:
        return 30;
    }
  }
}

/// Loan entity
@freezed
class Loan with _$Loan implements Entity {
  const factory Loan({
    required String id,
    required String userId,
    required Money principalAmount,
    required double interestRate,
    required Money totalAmount,
    required Money amountPaid,
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
    Money? nextPaymentAmount,
    @Default(0) int paymentsMade,
    @Default(0) int paymentsTotal,
    @Default(0) int daysOverdue,
    String? rejectionReason,
    @Default([]) List<LoanRepayment> repayments,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _Loan;

  const Loan._();

  /// Outstanding balance
  Money get outstandingBalance => totalAmount - amountPaid;

  /// Repayment progress (0.0 - 1.0)
  double get repaymentProgress =>
      totalAmount.amount > 0 ? amountPaid.amount / totalAmount.amount : 0;

  /// Repayment percentage (0 - 100)
  int get repaymentPercent => (repaymentProgress * 100).round();

  /// Check if loan is active
  bool get isActive =>
      status == LoanStatus.active || status == LoanStatus.disbursed;

  /// Check if loan is overdue
  bool get isOverdue => status == LoanStatus.overdue;

  /// Check if loan is paid off
  bool get isPaidOff => status == LoanStatus.paid;

  /// Days remaining until due
  int get daysRemaining {
    if (dueDate == null || !isActive) return 0;
    final days = dueDate!.difference(DateTime.now()).inDays;
    return days < 0 ? 0 : days;
  }

  /// Interest amount
  Money get interestAmount => totalAmount - principalAmount;
}
