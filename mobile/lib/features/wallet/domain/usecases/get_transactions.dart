import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/transaction.dart';
import '../repositories/wallet_repository.dart';

part 'get_transactions.freezed.dart';

/// Parameters for getting transactions
@freezed
class GetTransactionsParams with _$GetTransactionsParams {
  const factory GetTransactionsParams({
    @Default(1) int page,
    @Default(20) int limit,
    TransactionType? type,
    DateTime? startDate,
    DateTime? endDate,
  }) = _GetTransactionsParams;
}

/// Use case to get wallet transactions
@injectable
class GetTransactions implements UseCase<List<Transaction>, GetTransactionsParams> {
  final WalletRepository _repository;

  GetTransactions(this._repository);

  @override
  Future<Either<Failure, List<Transaction>>> call(GetTransactionsParams params) {
    return _repository.getTransactions(
      page: params.page,
      limit: params.limit,
      type: params.type,
      startDate: params.startDate,
      endDate: params.endDate,
    );
  }
}

/// Use case to get a single transaction
@injectable
class GetTransaction implements UseCase<Transaction, String> {
  final WalletRepository _repository;

  GetTransaction(this._repository);

  @override
  Future<Either<Failure, Transaction>> call(String transactionId) {
    return _repository.getTransaction(transactionId);
  }
}
