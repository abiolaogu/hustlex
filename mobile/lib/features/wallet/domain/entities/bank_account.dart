import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/base_entity.dart';
import '../../../../core/domain/value_objects/bvn.dart';

part 'bank_account.freezed.dart';

/// Bank account domain entity
@freezed
class BankAccount with _$BankAccount implements Entity {
  const BankAccount._();

  const factory BankAccount({
    required String id,
    required String userId,
    required String bankCode,
    required String bankName,
    required String accountNumber,
    required String accountName,
    @Default(false) bool isVerified,
    @Default(false) bool isPrimary,
    required DateTime createdAt,
  }) = _BankAccount;

  /// Get masked account number for display
  String get maskedAccountNumber {
    if (accountNumber.length < 4) return accountNumber;
    return '******${accountNumber.substring(accountNumber.length - 4)}';
  }

  /// Create AccountNumber value object
  AccountNumber get accountNumberVO => AccountNumber(accountNumber);
}

/// Bank verification result
@freezed
class BankAccountVerification with _$BankAccountVerification {
  const factory BankAccountVerification({
    required String accountNumber,
    required String accountName,
    required String bankCode,
    required String bankName,
  }) = _BankAccountVerification;
}

/// Supported bank
@freezed
class Bank with _$Bank {
  const factory Bank({
    required String code,
    required String name,
    String? slug,
    String? logo,
  }) = _Bank;
}
