import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/savings_stats.dart';
import '../repositories/savings_repository.dart';

/// Use case to get current user's savings statistics
@injectable
class GetSavingsStats implements UseCaseNoParams<SavingsStats> {
  final SavingsRepository _repository;

  GetSavingsStats(this._repository);

  @override
  Future<Either<Failure, SavingsStats>> call() {
    return _repository.getSavingsStats();
  }
}

/// Parameters for getting circle activity
class GetCircleActivityParams {
  final String circleId;
  final int page;
  final int limit;

  const GetCircleActivityParams({
    required this.circleId,
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get circle activity/history
@injectable
class GetCircleActivity
    implements UseCase<List<CircleActivity>, GetCircleActivityParams> {
  final SavingsRepository _repository;

  GetCircleActivity(this._repository);

  @override
  Future<Either<Failure, List<CircleActivity>>> call(
      GetCircleActivityParams params) {
    return _repository.getCircleActivity(
      circleId: params.circleId,
      page: params.page,
      limit: params.limit,
    );
  }
}
