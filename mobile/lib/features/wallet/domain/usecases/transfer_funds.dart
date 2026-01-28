import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../../../../core/domain/value_objects/phone_number.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/transaction.dart';
import '../repositories/wallet_repository.dart';

part 'transfer_funds.freezed.dart';

/// Parameters for transferring funds
@freezed
class TransferParams with _$TransferParams {
  const factory TransferParams({
    required Money amount,
    required PhoneNumber recipientPhone,
    required Pin pin,
    String? narration,
  }) = _TransferParams;
}

/// Use case to transfer funds to another user
@injectable
class TransferFunds implements UseCase<Transaction, TransferParams> {
  final WalletRepository _repository;

  TransferFunds(this._repository);

  @override
  Future<Either<Failure, Transaction>> call(TransferParams params) async {
    // Validate amount
    if (!params.amount.isValid()) {
      return left(const Failure.invalidInput('amount', 'Invalid amount'));
    }

    // Check minimum transfer amount (100 NGN = 10000 kobo)
    if (params.amount.amountInMinorUnits < 10000) {
      return left(const Failure.minimumAmountNotMet(10000));
    }

    // Validate recipient phone
    if (!params.recipientPhone.isValid()) {
      return left(
          const Failure.invalidInput('recipient', 'Invalid phone number'));
    }

    // Validate PIN
    if (!params.pin.isValid()) {
      final failure = params.pin.failureOrNull;
      return left(Failure.invalidInput('pin', failure?.message ?? 'Invalid PIN'));
    }

    return _repository.transfer(
      amount: params.amount,
      recipientPhone: params.recipientPhone,
      pin: params.pin.getOrCrash(),
      narration: params.narration,
    );
  }
}
