import 'package:json_annotation/json_annotation.dart';

part 'wallet_models.g.dart';

/// Wallet model
@JsonSerializable()
class Wallet {
  final String id;
  final String userId;
  final double balance;
  final double escrowBalance;
  final double savingsBalance;
  final String currency;
  final bool isActive;
  final DateTime createdAt;
  final DateTime updatedAt;

  Wallet({
    required this.id,
    required this.userId,
    required this.balance,
    required this.escrowBalance,
    required this.savingsBalance,
    this.currency = 'NGN',
    this.isActive = true,
    required this.createdAt,
    required this.updatedAt,
  });

  double get totalBalance => balance + escrowBalance + savingsBalance;

  factory Wallet.fromJson(Map<String, dynamic> json) => _$WalletFromJson(json);
  Map<String, dynamic> toJson() => _$WalletToJson(this);
}

/// Transaction types
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
  @JsonValue('gig_escrow')
  gigEscrow,
  @JsonValue('escrow_release')
  escrowRelease,
  @JsonValue('savings_contribution')
  savingsContribution,
  @JsonValue('savings_payout')
  savingsPayout,
  @JsonValue('loan_disbursement')
  loanDisbursement,
  @JsonValue('loan_repayment')
  loanRepayment,
  @JsonValue('refund')
  refund,
  @JsonValue('fee')
  fee,
}

/// Transaction status
enum TransactionStatus {
  @JsonValue('pending')
  pending,
  @JsonValue('completed')
  completed,
  @JsonValue('failed')
  failed,
  @JsonValue('reversed')
  reversed,
}

/// Transaction model
@JsonSerializable()
class Transaction {
  final String id;
  final String walletId;
  final TransactionType type;
  final double amount;
  final double balanceBefore;
  final double balanceAfter;
  final String? reference;
  final String? description;
  final TransactionStatus status;
  final Map<String, dynamic>? metadata;
  final DateTime createdAt;

  Transaction({
    required this.id,
    required this.walletId,
    required this.type,
    required this.amount,
    required this.balanceBefore,
    required this.balanceAfter,
    this.reference,
    this.description,
    required this.status,
    this.metadata,
    required this.createdAt,
  });

  bool get isCredit => amount > 0;
  bool get isDebit => amount < 0;

  factory Transaction.fromJson(Map<String, dynamic> json) =>
      _$TransactionFromJson(json);
  Map<String, dynamic> toJson() => _$TransactionToJson(this);
}

/// Bank account model
@JsonSerializable()
class BankAccount {
  final String id;
  final String userId;
  final String bankCode;
  final String bankName;
  final String accountNumber;
  final String accountName;
  final bool isVerified;
  final bool isPrimary;
  final DateTime createdAt;

  BankAccount({
    required this.id,
    required this.userId,
    required this.bankCode,
    required this.bankName,
    required this.accountNumber,
    required this.accountName,
    this.isVerified = false,
    this.isPrimary = false,
    required this.createdAt,
  });

  String get maskedAccountNumber {
    if (accountNumber.length < 4) return accountNumber;
    return '****${accountNumber.substring(accountNumber.length - 4)}';
  }

  factory BankAccount.fromJson(Map<String, dynamic> json) =>
      _$BankAccountFromJson(json);
  Map<String, dynamic> toJson() => _$BankAccountToJson(this);
}

/// Transfer request
@JsonSerializable()
class TransferRequest {
  final String recipientPhone;
  final double amount;
  final String? note;
  final String pin;

  TransferRequest({
    required this.recipientPhone,
    required this.amount,
    this.note,
    required this.pin,
  });

  factory TransferRequest.fromJson(Map<String, dynamic> json) =>
      _$TransferRequestFromJson(json);
  Map<String, dynamic> toJson() => _$TransferRequestToJson(this);
}

/// Withdrawal request
@JsonSerializable()
class WithdrawalRequest {
  final String bankAccountId;
  final double amount;
  final String pin;

  WithdrawalRequest({
    required this.bankAccountId,
    required this.amount,
    required this.pin,
  });

  factory WithdrawalRequest.fromJson(Map<String, dynamic> json) =>
      _$WithdrawalRequestFromJson(json);
  Map<String, dynamic> toJson() => _$WithdrawalRequestToJson(this);
}

/// Deposit response with Paystack
@JsonSerializable()
class DepositInitiation {
  final String reference;
  final String authorizationUrl;
  final String accessCode;

  DepositInitiation({
    required this.reference,
    required this.authorizationUrl,
    required this.accessCode,
  });

  factory DepositInitiation.fromJson(Map<String, dynamic> json) =>
      _$DepositInitiationFromJson(json);
  Map<String, dynamic> toJson() => _$DepositInitiationToJson(this);
}
