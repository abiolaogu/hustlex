import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/base_entity.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'transaction.freezed.dart';

/// Transaction types
enum TransactionType {
  deposit,
  withdrawal,
  transferIn,
  transferOut,
  gigPayment,
  gigEscrow,
  escrowRelease,
  savingsContribution,
  savingsPayout,
  loanDisbursement,
  loanRepayment,
  refund,
  fee,
}

/// Transaction status
enum TransactionStatus {
  pending,
  processing,
  completed,
  failed,
  reversed,
}

/// Transaction domain entity
@freezed
class Transaction with _$Transaction implements Entity {
  const Transaction._();

  const factory Transaction({
    required String id,
    required String walletId,
    required TransactionType type,
    required Money amount,
    required Money balanceBefore,
    required Money balanceAfter,
    required TransactionStatus status,
    String? reference,
    String? description,
    String? recipientName,
    String? recipientPhone,
    Map<String, dynamic>? metadata,
    required DateTime createdAt,
  }) = _Transaction;

  /// Check if this is a credit transaction
  bool get isCredit =>
      type == TransactionType.deposit ||
      type == TransactionType.transferIn ||
      type == TransactionType.escrowRelease ||
      type == TransactionType.savingsPayout ||
      type == TransactionType.loanDisbursement ||
      type == TransactionType.refund;

  /// Check if this is a debit transaction
  bool get isDebit => !isCredit;

  /// Check if transaction is completed
  bool get isCompleted => status == TransactionStatus.completed;

  /// Check if transaction is pending
  bool get isPending =>
      status == TransactionStatus.pending ||
      status == TransactionStatus.processing;

  /// Check if transaction failed
  bool get isFailed =>
      status == TransactionStatus.failed ||
      status == TransactionStatus.reversed;

  /// Get display title for transaction type
  String get displayTitle {
    switch (type) {
      case TransactionType.deposit:
        return 'Deposit';
      case TransactionType.withdrawal:
        return 'Withdrawal';
      case TransactionType.transferIn:
        return 'Received';
      case TransactionType.transferOut:
        return 'Sent';
      case TransactionType.gigPayment:
        return 'Gig Payment';
      case TransactionType.gigEscrow:
        return 'Escrow';
      case TransactionType.escrowRelease:
        return 'Escrow Release';
      case TransactionType.savingsContribution:
        return 'Savings';
      case TransactionType.savingsPayout:
        return 'Savings Payout';
      case TransactionType.loanDisbursement:
        return 'Loan';
      case TransactionType.loanRepayment:
        return 'Loan Repayment';
      case TransactionType.refund:
        return 'Refund';
      case TransactionType.fee:
        return 'Fee';
    }
  }
}
