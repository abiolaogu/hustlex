import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/credit_score.dart';
import '../repositories/credit_repository.dart';

/// Use case to get current user's credit score
@injectable
class GetCreditScore implements UseCaseNoParams<CreditScore> {
  final CreditRepository _repository;

  GetCreditScore(this._repository);

  @override
  Future<Either<Failure, CreditScore>> call() {
    return _repository.getCreditScore();
  }
}

/// Use case to watch credit score for real-time updates
@injectable
class WatchCreditScore implements StreamUseCaseNoParams<CreditScore> {
  final CreditRepository _repository;

  WatchCreditScore(this._repository);

  @override
  Stream<Either<Failure, CreditScore>> call() {
    return _repository.watchCreditScore();
  }
}

/// Parameters for getting credit score history
class GetCreditHistoryParams {
  final int page;
  final int limit;

  const GetCreditHistoryParams({
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get credit score history
@injectable
class GetCreditScoreHistory
    implements UseCase<List<CreditScoreHistory>, GetCreditHistoryParams> {
  final CreditRepository _repository;

  GetCreditScoreHistory(this._repository);

  @override
  Future<Either<Failure, List<CreditScoreHistory>>> call(
      GetCreditHistoryParams params) {
    return _repository.getCreditScoreHistory(
      page: params.page,
      limit: params.limit,
    );
  }
}

/// Use case to get tips to improve credit score
@injectable
class GetCreditTips implements UseCaseNoParams<List<CreditTip>> {
  final CreditRepository _repository;

  GetCreditTips(this._repository);

  @override
  Future<Either<Failure, List<CreditTip>>> call() {
    return _repository.getCreditTips();
  }
}

/// Use case to refresh credit score
@injectable
class RefreshCreditScore implements UseCaseNoParams<CreditScore> {
  final CreditRepository _repository;

  RefreshCreditScore(this._repository);

  @override
  Future<Either<Failure, CreditScore>> call() {
    return _repository.refreshCreditScore();
  }
}
