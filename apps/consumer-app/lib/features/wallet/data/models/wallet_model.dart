import 'package:freezed_annotation/freezed_annotation.dart';

part 'wallet_model.freezed.dart';
part 'wallet_model.g.dart';

enum TransactionType {
  @JsonValue('deposit')
  deposit,
  @JsonValue('withdrawal')
  withdrawal,
  @JsonValue('transfer_in')
  transferIn,
  @JsonValue('transfer_out')
  transferOut,
  @JsonValue('gig_payment')
  gigPayment,
  @JsonValue('gig_earning')
  gigEarning,
  @JsonValue('savings_contribution')
  savingsContribution,
  @JsonValue('savings_payout')
  savingsPayout,
  @JsonValue('loan_disbursement')
  loanDisbursement,
  @JsonValue('loan_repayment')
  loanRepayment,
  @JsonValue('bill_payment')
  billPayment,
  @JsonValue('airtime')
  airtime,
  @JsonValue('refund')
  refund,
  @JsonValue('fee')
  fee,
  @JsonValue('cashback')
  cashback,
}

enum TransactionStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('processing')
  processing,
  @JsonValue('completed')
  completed,
  @JsonValue('failed')
  failed,
  @JsonValue('cancelled')
  cancelled,
  @JsonValue('reversed')
  reversed,
}

enum PaymentMethod {
  @JsonValue('wallet')
  wallet,
  @JsonValue('card')
  card,
  @JsonValue('bank_transfer')
  bankTransfer,
  @JsonValue('ussd')
  ussd,
  @JsonValue('paystack')
  paystack,
}

enum BillCategory {
  @JsonValue('electricity')
  electricity,
  @JsonValue('cable')
  cable,
  @JsonValue('internet')
  internet,
  @JsonValue('water')
  water,
  @JsonValue('airtime')
  airtime,
  @JsonValue('data')
  data,
}

extension TransactionTypeX on TransactionType {
  String get displayName {
    switch (this) {
      case TransactionType.deposit:
        return 'Deposit';
      case TransactionType.withdrawal:
        return 'Withdrawal';
      case TransactionType.transferIn:
        return 'Transfer Received';
      case TransactionType.transferOut:
        return 'Transfer Sent';
      case TransactionType.gigPayment:
        return 'Gig Payment';
      case TransactionType.gigEarning:
        return 'Gig Earning';
      case TransactionType.savingsContribution:
        return 'Savings Contribution';
      case TransactionType.savingsPayout:
        return 'Savings Payout';
      case TransactionType.loanDisbursement:
        return 'Loan Received';
      case TransactionType.loanRepayment:
        return 'Loan Repayment';
      case TransactionType.billPayment:
        return 'Bill Payment';
      case TransactionType.airtime:
        return 'Airtime Purchase';
      case TransactionType.refund:
        return 'Refund';
      case TransactionType.fee:
        return 'Fee';
      case TransactionType.cashback:
        return 'Cashback';
    }
  }

  bool get isCredit {
    return [
      TransactionType.deposit,
      TransactionType.transferIn,
      TransactionType.gigEarning,
      TransactionType.savingsPayout,
      TransactionType.loanDisbursement,
      TransactionType.refund,
      TransactionType.cashback,
    ].contains(this);
  }

  String get icon {
    switch (this) {
      case TransactionType.deposit:
        return '⬇️';
      case TransactionType.withdrawal:
        return '⬆️';
      case TransactionType.transferIn:
        return '↙️';
      case TransactionType.transferOut:
        return '↗️';
      case TransactionType.gigPayment:
        return '💼';
      case TransactionType.gigEarning:
        return '💰';
      case TransactionType.savingsContribution:
        return '🏦';
      case TransactionType.savingsPayout:
        return '🎉';
      case TransactionType.loanDisbursement:
        return '💵';
      case TransactionType.loanRepayment:
        return '📋';
      case TransactionType.billPayment:
        return '🧾';
      case TransactionType.airtime:
        return '📱';
      case TransactionType.refund:
        return '↩️';
      case TransactionType.fee:
        return '💳';
      case TransactionType.cashback:
        return '🎁';
    }
  }
}

extension TransactionStatusX on TransactionStatus {
  String get displayName {
    switch (this) {
      case TransactionStatus.pending:
        return 'Pending';
      case TransactionStatus.processing:
        return 'Processing';
      case TransactionStatus.completed:
        return 'Completed';
      case TransactionStatus.failed:
        return 'Failed';
      case TransactionStatus.cancelled:
        return 'Cancelled';
      case TransactionStatus.reversed:
        return 'Reversed';
    }
  }
}

@freezed
class Wallet with _$Wallet {
  const factory Wallet({
    required String id,
    required String userId,
    required double availableBalance,
    required double ledgerBalance,
    @Default(0) double escrowBalance,
    @Default(0) double savingsBalance,
    required String currency,
    @Default(false) bool isLocked,
    String? lockReason,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Wallet;

  const Wallet._();

  factory Wallet.fromJson(Map<String, dynamic> json) => _$WalletFromJson(json);

  double get totalBalance => availableBalance + escrowBalance + savingsBalance;

  String formatAmount(double amount) {
    if (amount >= 1000000) {
      return '₦${(amount / 1000000).toStringAsFixed(2)}M';
    } else if (amount >= 1000) {
      return '₦${(amount / 1000).toStringAsFixed(2)}K';
    }
    return '₦${amount.toStringAsFixed(2)}';
  }

  String get formattedAvailable => formatAmount(availableBalance);
  String get formattedTotal => formatAmount(totalBalance);
}

@freezed
class Transaction with _$Transaction {
  const factory Transaction({
    required String id,
    required String walletId,
    required TransactionType type,
    required TransactionStatus status,
    required double amount,
    required String currency,
    @Default(0) double fee,
    double? balanceBefore,
    double? balanceAfter,
    String? description,
    String? reference,
    String? externalReference,
    String? counterpartyId,
    String? counterpartyName,
    String? counterpartyAccount,
    String? counterpartyBank,
    String? relatedEntityId,
    String? relatedEntityType,
    Map<String, dynamic>? metadata,
    DateTime? createdAt,
    DateTime? completedAt,
  }) = _Transaction;

  const Transaction._();

  factory Transaction.fromJson(Map<String, dynamic> json) => _$TransactionFromJson(json);

  bool get isCredit => type.isCredit;
  bool get isDebit => !isCredit;
  bool get isPending => status == TransactionStatus.pending;
  bool get isCompleted => status == TransactionStatus.completed;
  bool get isFailed => status == TransactionStatus.failed;

  String get formattedAmount {
    final sign = isCredit ? '+' : '-';
    return '$sign₦${amount.toStringAsFixed(2)}';
  }

  String get displayTitle {
    if (description != null && description!.isNotEmpty) {
      return description!;
    }
    if (counterpartyName != null) {
      return isCredit ? 'From $counterpartyName' : 'To $counterpartyName';
    }
    return type.displayName;
  }
}

@freezed
class BankAccount with _$BankAccount {
  const factory BankAccount({
    required String id,
    required String userId,
    required String bankCode,
    required String bankName,
    required String accountNumber,
    required String accountName,
    @Default(false) bool isDefault,
    @Default(false) bool isVerified,
    DateTime? createdAt,
  }) = _BankAccount;

  const BankAccount._();

  factory BankAccount.fromJson(Map<String, dynamic> json) => _$BankAccountFromJson(json);

  String get maskedAccountNumber {
    if (accountNumber.length < 4) return accountNumber;
    return '****${accountNumber.substring(accountNumber.length - 4)}';
  }
}

@freezed
class Bank with _$Bank {
  const factory Bank({
    required String code,
    required String name,
    String? logo,
    @Default(true) bool isActive,
  }) = _Bank;

  factory Bank.fromJson(Map<String, dynamic> json) => _$BankFromJson(json);
}

@freezed
class DepositRequest with _$DepositRequest {
  const factory DepositRequest({
    required double amount,
    required PaymentMethod method,
    String? paymentReference,
    Map<String, dynamic>? metadata,
  }) = _DepositRequest;

  factory DepositRequest.fromJson(Map<String, dynamic> json) => _$DepositRequestFromJson(json);
}

@freezed
class WithdrawalRequest with _$WithdrawalRequest {
  const factory WithdrawalRequest({
    required double amount,
    required String bankAccountId,
    required String pin,
    String? narration,
  }) = _WithdrawalRequest;

  factory WithdrawalRequest.fromJson(Map<String, dynamic> json) => _$WithdrawalRequestFromJson(json);
}

@freezed
class TransferRequest with _$TransferRequest {
  const factory TransferRequest({
    required double amount,
    required String recipientPhone,
    required String pin,
    String? narration,
  }) = _TransferRequest;

  factory TransferRequest.fromJson(Map<String, dynamic> json) => _$TransferRequestFromJson(json);
}

@freezed
class BankTransferRequest with _$BankTransferRequest {
  const factory BankTransferRequest({
    required double amount,
    required String bankCode,
    required String accountNumber,
    required String accountName,
    required String pin,
    String? narration,
  }) = _BankTransferRequest;

  factory BankTransferRequest.fromJson(Map<String, dynamic> json) => _$BankTransferRequestFromJson(json);
}

@freezed
class BillPaymentRequest with _$BillPaymentRequest {
  const factory BillPaymentRequest({
    required BillCategory category,
    required String providerId,
    required String customerId,
    required double amount,
    required String pin,
    String? phone,
    Map<String, dynamic>? metadata,
  }) = _BillPaymentRequest;

  factory BillPaymentRequest.fromJson(Map<String, dynamic> json) => _$BillPaymentRequestFromJson(json);
}

@freezed
class AirtimeRequest with _$AirtimeRequest {
  const factory AirtimeRequest({
    required String phone,
    required double amount,
    required String network,
    required String pin,
  }) = _AirtimeRequest;

  factory AirtimeRequest.fromJson(Map<String, dynamic> json) => _$AirtimeRequestFromJson(json);
}

@freezed
class AccountVerification with _$AccountVerification {
  const factory AccountVerification({
    required String accountNumber,
    required String accountName,
    required String bankCode,
    required String bankName,
  }) = _AccountVerification;

  factory AccountVerification.fromJson(Map<String, dynamic> json) => _$AccountVerificationFromJson(json);
}

@freezed
class PaystackInitResponse with _$PaystackInitResponse {
  const factory PaystackInitResponse({
    required String authorizationUrl,
    required String accessCode,
    required String reference,
  }) = _PaystackInitResponse;

  factory PaystackInitResponse.fromJson(Map<String, dynamic> json) => _$PaystackInitResponseFromJson(json);
}

@freezed
class TransactionFilter with _$TransactionFilter {
  const factory TransactionFilter({
    TransactionType? type,
    TransactionStatus? status,
    DateTime? startDate,
    DateTime? endDate,
    double? minAmount,
    double? maxAmount,
    String? search,
    @Default(1) int page,
    @Default(20) int limit,
    @Default('created_at') String sortBy,
    @Default('desc') String sortOrder,
  }) = _TransactionFilter;

  factory TransactionFilter.fromJson(Map<String, dynamic> json) => _$TransactionFilterFromJson(json);
}

@freezed
class PaginatedTransactions with _$PaginatedTransactions {
  const factory PaginatedTransactions({
    required List<Transaction> transactions,
    required int total,
    required int page,
    required int limit,
    required bool hasMore,
  }) = _PaginatedTransactions;

  factory PaginatedTransactions.fromJson(Map<String, dynamic> json) => _$PaginatedTransactionsFromJson(json);
}

@freezed
class WalletStats with _$WalletStats {
  const factory WalletStats({
    @Default(0) double totalIncome,
    @Default(0) double totalExpenses,
    @Default(0) int transactionsCount,
    @Default(0) double avgTransactionAmount,
    Map<String, double>? incomeByType,
    Map<String, double>? expensesByType,
  }) = _WalletStats;

  factory WalletStats.fromJson(Map<String, dynamic> json) => _$WalletStatsFromJson(json);
}
