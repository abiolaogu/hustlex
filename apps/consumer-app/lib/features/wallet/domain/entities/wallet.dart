import 'package:freezed_annotation/freezed_annotation.dart';
import '../../../../core/domain/entities/base_entity.dart';
import '../../../../core/domain/value_objects/money.dart';

part 'wallet.freezed.dart';

/// Wallet status
enum WalletStatus {
  active,
  locked,
  suspended,
}

/// Wallet domain entity
@freezed
class Wallet with _$Wallet implements Entity {
  const Wallet._();

  const factory Wallet({
    required String id,
    required String userId,
    required Money availableBalance,
    required Money escrowBalance,
    required Money savingsBalance,
    required WalletStatus status,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _Wallet;

  /// Total balance across all wallet components
  Money get totalBalance =>
      availableBalance + escrowBalance + savingsBalance;

  /// Check if wallet is active
  bool get isActive => status == WalletStatus.active;

  /// Check if user can perform transactions
  bool get canTransact => isActive;

  /// Check if user has sufficient funds for a transaction
  bool hasSufficientFunds(Money amount) =>
      availableBalance >= amount;
}
