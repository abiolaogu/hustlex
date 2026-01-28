import 'package:dartz/dartz.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../../../../core/domain/value_objects/phone_number.dart';
import '../entities/wallet.dart';
import '../entities/transaction.dart';
import '../entities/bank_account.dart';
import '../entities/deposit.dart';

/// Wallet repository interface - defines the contract for wallet operations.
/// Implementation will be in the data layer.
abstract class WalletRepository {
  /// Get the current user's wallet
  Future<Either<Failure, Wallet>> getWallet();

  /// Watch wallet balance changes (real-time updates)
  Stream<Either<Failure, Wallet>> watchWallet();

  /// Get wallet transactions with pagination
  Future<Either<Failure, List<Transaction>>> getTransactions({
    int page = 1,
    int limit = 20,
    TransactionType? type,
    DateTime? startDate,
    DateTime? endDate,
  });

  /// Get a single transaction by ID
  Future<Either<Failure, Transaction>> getTransaction(String id);

  /// Initialize a deposit via payment gateway
  Future<Either<Failure, DepositInitiation>> initiateDeposit(Money amount);

  /// Verify a deposit after payment
  Future<Either<Failure, DepositVerification>> verifyDeposit(String reference);

  /// Initiate a withdrawal to a bank account
  Future<Either<Failure, Transaction>> initiateWithdrawal({
    required Money amount,
    required String bankAccountId,
    required String pin,
  });

  /// Transfer funds to another user
  Future<Either<Failure, Transaction>> transfer({
    required Money amount,
    required PhoneNumber recipientPhone,
    required String pin,
    String? narration,
  });

  /// Get list of supported banks
  Future<Either<Failure, List<Bank>>> getBanks();

  /// Verify a bank account
  Future<Either<Failure, BankAccountVerification>> verifyBankAccount({
    required String accountNumber,
    required String bankCode,
  });

  /// Add a bank account
  Future<Either<Failure, BankAccount>> addBankAccount({
    required String accountNumber,
    required String bankCode,
    required String accountName,
  });

  /// Get user's saved bank accounts
  Future<Either<Failure, List<BankAccount>>> getBankAccounts();

  /// Delete a bank account
  Future<Either<Failure, Unit>> deleteBankAccount(String id);

  /// Set a bank account as primary
  Future<Either<Failure, BankAccount>> setPrimaryBankAccount(String id);
}
