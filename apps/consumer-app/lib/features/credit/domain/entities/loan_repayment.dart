import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/entity.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'loan_repayment.freezed.dart';

/// Repayment status
enum RepaymentStatus {
  pending,
  paid,
  overdue,
  missed,
  partial,
}

/// Extension for RepaymentStatus
extension RepaymentStatusX on RepaymentStatus {
  String get displayName {
    switch (this) {
      case RepaymentStatus.pending:
        return 'Pending';
      case RepaymentStatus.paid:
        return 'Paid';
      case RepaymentStatus.overdue:
        return 'Overdue';
      case RepaymentStatus.missed:
        return 'Missed';
      case RepaymentStatus.partial:
        return 'Partial';
    }
  }

  bool get requiresPayment =>
      this == RepaymentStatus.pending ||
      this == RepaymentStatus.overdue ||
      this == RepaymentStatus.partial;
}

/// Loan repayment entity
@freezed
class LoanRepayment with _$LoanRepayment implements Entity {
  const factory LoanRepayment({
    required String id,
    required String loanId,
    required Money amount,
    required Money principalPortion,
    required Money interestPortion,
    required int installmentNumber,
    required RepaymentStatus status,
    required DateTime dueDate,
    DateTime? paidAt,
    String? transactionId,
    Money? lateFee,
    Money? amountPaid,
    required DateTime createdAt,
  }) = _LoanRepayment;

  const LoanRepayment._();

  /// Check if repayment is paid
  bool get isPaid => status == RepaymentStatus.paid;

  /// Check if repayment is overdue
  bool get isOverdue => status == RepaymentStatus.overdue;

  /// Check if repayment is pending
  bool get isPending => status == RepaymentStatus.pending;

  /// Check if repayment requires payment
  bool get requiresPayment => status.requiresPayment;

  /// Total amount due including late fee
  Money get totalDue => lateFee != null ? amount + lateFee! : amount;

  /// Days until due (negative if overdue)
  int get daysUntilDue {
    if (isPaid) return 0;
    return dueDate.difference(DateTime.now()).inDays;
  }

  /// Check if due today
  bool get isDueToday {
    final now = DateTime.now();
    return dueDate.year == now.year &&
        dueDate.month == now.month &&
        dueDate.day == now.day;
  }
}
