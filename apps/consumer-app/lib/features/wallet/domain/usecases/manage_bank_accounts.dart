import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/bvn.dart';
import '../entities/bank_account.dart';
import '../repositories/wallet_repository.dart';

part 'manage_bank_accounts.freezed.dart';

/// Use case to get all supported banks
@injectable
class GetBanks implements UseCaseNoParams<List<Bank>> {
  final WalletRepository _repository;

  GetBanks(this._repository);

  @override
  Future<Either<Failure, List<Bank>>> call() {
    return _repository.getBanks();
  }
}

/// Use case to get user's bank accounts
@injectable
class GetBankAccounts implements UseCaseNoParams<List<BankAccount>> {
  final WalletRepository _repository;

  GetBankAccounts(this._repository);

  @override
  Future<Either<Failure, List<BankAccount>>> call() {
    return _repository.getBankAccounts();
  }
}

/// Parameters for verifying a bank account
@freezed
class VerifyBankAccountParams with _$VerifyBankAccountParams {
  const factory VerifyBankAccountParams({
    required AccountNumber accountNumber,
    required String bankCode,
  }) = _VerifyBankAccountParams;
}

/// Use case to verify a bank account
@injectable
class VerifyBankAccount
    implements UseCase<BankAccountVerification, VerifyBankAccountParams> {
  final WalletRepository _repository;

  VerifyBankAccount(this._repository);

  @override
  Future<Either<Failure, BankAccountVerification>> call(
    VerifyBankAccountParams params,
  ) async {
    // Validate account number
    if (!params.accountNumber.isValid()) {
      return left(
        const Failure.invalidInput('accountNumber', 'Invalid account number'),
      );
    }

    // Validate bank code
    if (params.bankCode.isEmpty) {
      return left(
        const Failure.invalidInput('bankCode', 'Please select a bank'),
      );
    }

    return _repository.verifyBankAccount(
      accountNumber: params.accountNumber.getOrCrash(),
      bankCode: params.bankCode,
    );
  }
}

/// Parameters for adding a bank account
@freezed
class AddBankAccountParams with _$AddBankAccountParams {
  const factory AddBankAccountParams({
    required AccountNumber accountNumber,
    required String bankCode,
    required String accountName,
  }) = _AddBankAccountParams;
}

/// Use case to add a bank account
@injectable
class AddBankAccount implements UseCase<BankAccount, AddBankAccountParams> {
  final WalletRepository _repository;

  AddBankAccount(this._repository);

  @override
  Future<Either<Failure, BankAccount>> call(AddBankAccountParams params) async {
    // Validate account number
    if (!params.accountNumber.isValid()) {
      return left(
        const Failure.invalidInput('accountNumber', 'Invalid account number'),
      );
    }

    // Validate bank code
    if (params.bankCode.isEmpty) {
      return left(
        const Failure.invalidInput('bankCode', 'Please select a bank'),
      );
    }

    // Validate account name
    if (params.accountName.trim().isEmpty) {
      return left(
        const Failure.invalidInput('accountName', 'Account name is required'),
      );
    }

    return _repository.addBankAccount(
      accountNumber: params.accountNumber.getOrCrash(),
      bankCode: params.bankCode,
      accountName: params.accountName.trim(),
    );
  }
}

/// Use case to delete a bank account
@injectable
class DeleteBankAccount implements UseCase<Unit, String> {
  final WalletRepository _repository;

  DeleteBankAccount(this._repository);

  @override
  Future<Either<Failure, Unit>> call(String bankAccountId) {
    if (bankAccountId.isEmpty) {
      return Future.value(
        left(const Failure.invalidInput('id', 'Bank account ID is required')),
      );
    }
    return _repository.deleteBankAccount(bankAccountId);
  }
}

/// Use case to set primary bank account
@injectable
class SetPrimaryBankAccount implements UseCase<BankAccount, String> {
  final WalletRepository _repository;

  SetPrimaryBankAccount(this._repository);

  @override
  Future<Either<Failure, BankAccount>> call(String bankAccountId) {
    if (bankAccountId.isEmpty) {
      return Future.value(
        left(const Failure.invalidInput('id', 'Bank account ID is required')),
      );
    }
    return _repository.setPrimaryBankAccount(bankAccountId);
  }
}
