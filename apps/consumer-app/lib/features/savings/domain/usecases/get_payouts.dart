import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/payout.dart';
import '../repositories/savings_repository.dart';

/// Parameters for getting payouts
class GetPayoutsParams {
  final String circleId;
  final PayoutStatus? status;

  const GetPayoutsParams({
    required this.circleId,
    this.status,
  });
}

/// Use case to get payouts for a circle
@injectable
class GetPayouts implements UseCase<List<Payout>, GetPayoutsParams> {
  final SavingsRepository _repository;

  GetPayouts(this._repository);

  @override
  Future<Either<Failure, List<Payout>>> call(GetPayoutsParams params) {
    return _repository.getPayouts(
      circleId: params.circleId,
      status: params.status,
    );
  }
}

/// Parameters for getting payout history
class GetPayoutHistoryParams {
  final String? circleId;
  final int page;
  final int limit;

  const GetPayoutHistoryParams({
    this.circleId,
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get current user's payout history
@injectable
class GetMyPayoutHistory
    implements UseCase<List<Payout>, GetPayoutHistoryParams> {
  final SavingsRepository _repository;

  GetMyPayoutHistory(this._repository);

  @override
  Future<Either<Failure, List<Payout>>> call(GetPayoutHistoryParams params) {
    return _repository.getMyPayoutHistory(
      circleId: params.circleId,
      page: params.page,
      limit: params.limit,
    );
  }
}

/// Use case to get next scheduled payout for a circle
@injectable
class GetNextPayout implements UseCase<Payout?, String> {
  final SavingsRepository _repository;

  GetNextPayout(this._repository);

  @override
  Future<Either<Failure, Payout?>> call(String circleId) {
    return _repository.getNextPayout(circleId);
  }
}
