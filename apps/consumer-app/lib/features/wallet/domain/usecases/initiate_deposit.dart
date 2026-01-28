import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../entities/deposit.dart';
import '../repositories/wallet_repository.dart';

/// Use case to initiate a deposit
@injectable
class InitiateDeposit implements UseCase<DepositInitiation, Money> {
  final WalletRepository _repository;

  InitiateDeposit(this._repository);

  @override
  Future<Either<Failure, DepositInitiation>> call(Money amount) async {
    // Validate amount
    if (!amount.isValid()) {
      return left(const Failure.invalidInput('amount', 'Invalid amount'));
    }

    // Check minimum deposit amount (100 NGN = 10000 kobo)
    if (amount.amountInMinorUnits < 10000) {
      return left(const Failure.minimumAmountNotMet(10000));
    }

    // Check maximum deposit amount (10,000,000 NGN = 1,000,000,000 kobo)
    if (amount.amountInMinorUnits > 1000000000) {
      return left(const Failure.dailyLimitExceeded());
    }

    return _repository.initiateDeposit(amount);
  }
}

/// Use case to verify a deposit
@injectable
class VerifyDeposit implements UseCase<DepositVerification, String> {
  final WalletRepository _repository;

  VerifyDeposit(this._repository);

  @override
  Future<Either<Failure, DepositVerification>> call(String reference) {
    if (reference.isEmpty) {
      return Future.value(
        left(const Failure.invalidInput('reference', 'Reference is required')),
      );
    }
    return _repository.verifyDeposit(reference);
  }
}
