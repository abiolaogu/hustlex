import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/failures/value_failures.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/transaction.dart';
import '../repositories/wallet_repository.dart';

part 'initiate_withdrawal.freezed.dart';

/// Parameters for initiating a withdrawal
@freezed
class WithdrawalParams with _$WithdrawalParams {
  const factory WithdrawalParams({
    required Money amount,
    required String bankAccountId,
    required Pin pin,
  }) = _WithdrawalParams;
}

/// Use case to initiate a withdrawal
@injectable
class InitiateWithdrawal implements UseCase<Transaction, WithdrawalParams> {
  final WalletRepository _repository;

  InitiateWithdrawal(this._repository);

  @override
  Future<Either<Failure, Transaction>> call(WithdrawalParams params) async {
    // Validate amount
    if (!params.amount.isValid()) {
      return left(const Failure.invalidInput('amount', 'Invalid amount'));
    }

    // Check minimum withdrawal amount (100 NGN = 10000 kobo)
    if (params.amount.amountInMinorUnits < 10000) {
      return left(const Failure.minimumAmountNotMet(10000));
    }

    // Validate bank account ID
    if (params.bankAccountId.isEmpty) {
      return left(
        const Failure.invalidInput('bankAccount', 'Please select a bank account'),
      );
    }

    // Validate PIN
    if (!params.pin.isValid()) {
      final failure = params.pin.failureOrNull;
      return left(
        Failure.invalidInput('pin', failure?.message ?? 'Invalid PIN'),
      );
    }

    return _repository.initiateWithdrawal(
      amount: params.amount,
      bankAccountId: params.bankAccountId,
      pin: params.pin.getOrCrash(),
    );
  }
}
