import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/savings_circle.dart';
import '../repositories/savings_repository.dart';

/// Use case to get paginated list of public circles
@injectable
class GetCircles implements UseCase<PaginatedCircles, CircleFilter> {
  final SavingsRepository _repository;

  GetCircles(this._repository);

  @override
  Future<Either<Failure, PaginatedCircles>> call(CircleFilter params) {
    return _repository.getCircles(params);
  }
}

/// Parameters for getting user's circles
class GetMyCirclesParams {
  final CircleStatus? status;

  const GetMyCirclesParams({this.status});
}

/// Use case to get circles that current user is a member of
@injectable
class GetMyCircles implements UseCase<List<SavingsCircle>, GetMyCirclesParams> {
  final SavingsRepository _repository;

  GetMyCircles(this._repository);

  @override
  Future<Either<Failure, List<SavingsCircle>>> call(GetMyCirclesParams params) {
    return _repository.getMyCircles(status: params.status);
  }
}

/// Use case to get a single circle by ID
@injectable
class GetCircle implements UseCase<SavingsCircle, String> {
  final SavingsRepository _repository;

  GetCircle(this._repository);

  @override
  Future<Either<Failure, SavingsCircle>> call(String circleId) {
    return _repository.getCircle(circleId);
  }
}

/// Use case to watch a circle for real-time updates
@injectable
class WatchCircle implements StreamUseCase<SavingsCircle, String> {
  final SavingsRepository _repository;

  WatchCircle(this._repository);

  @override
  Stream<Either<Failure, SavingsCircle>> call(String circleId) {
    return _repository.watchCircle(circleId);
  }
}
