import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/contribution.dart';
import '../repositories/savings_repository.dart';

part 'manage_contribution.freezed.dart';

/// Parameters for getting contributions
class GetContributionsParams {
  final String circleId;
  final int? cycleNumber;
  final ContributionStatus? status;

  const GetContributionsParams({
    required this.circleId,
    this.cycleNumber,
    this.status,
  });
}

/// Use case to get contributions for a circle
@injectable
class GetContributions
    implements UseCase<List<Contribution>, GetContributionsParams> {
  final SavingsRepository _repository;

  GetContributions(this._repository);

  @override
  Future<Either<Failure, List<Contribution>>> call(
      GetContributionsParams params) {
    return _repository.getContributions(
      circleId: params.circleId,
      cycleNumber: params.cycleNumber,
      status: params.status,
    );
  }
}

/// Use case to get current user's pending contributions
@injectable
class GetMyPendingContributions implements UseCaseNoParams<List<Contribution>> {
  final SavingsRepository _repository;

  GetMyPendingContributions(this._repository);

  @override
  Future<Either<Failure, List<Contribution>>> call() {
    return _repository.getMyPendingContributions();
  }
}

/// Parameters for getting contribution history
class GetContributionHistoryParams {
  final String? circleId;
  final int page;
  final int limit;

  const GetContributionHistoryParams({
    this.circleId,
    this.page = 1,
    this.limit = 20,
  });
}

/// Use case to get current user's contribution history
@injectable
class GetMyContributionHistory
    implements UseCase<List<Contribution>, GetContributionHistoryParams> {
  final SavingsRepository _repository;

  GetMyContributionHistory(this._repository);

  @override
  Future<Either<Failure, List<Contribution>>> call(
      GetContributionHistoryParams params) {
    return _repository.getMyContributionHistory(
      circleId: params.circleId,
      page: params.page,
      limit: params.limit,
    );
  }
}

/// Parameters for making a contribution
@freezed
class MakeContributionParams with _$MakeContributionParams {
  const factory MakeContributionParams({
    required String circleId,
    required String contributionId,
    required Pin pin,
  }) = _MakeContributionParams;
}

/// Use case to make a contribution
@injectable
class MakeContribution implements UseCase<Contribution, MakeContributionParams> {
  final SavingsRepository _repository;

  MakeContribution(this._repository);

  @override
  Future<Either<Failure, Contribution>> call(
      MakeContributionParams params) async {
    // Validate PIN
    if (!params.pin.isValid()) {
      return left(const Failure.invalidPin());
    }

    return _repository.makeContribution(
      circleId: params.circleId,
      contributionId: params.contributionId,
      pin: params.pin,
    );
  }
}

/// Use case to make contribution from wallet balance
@injectable
class MakeContributionFromWallet
    implements UseCase<Contribution, MakeContributionParams> {
  final SavingsRepository _repository;

  MakeContributionFromWallet(this._repository);

  @override
  Future<Either<Failure, Contribution>> call(
      MakeContributionParams params) async {
    // Validate PIN
    if (!params.pin.isValid()) {
      return left(const Failure.invalidPin());
    }

    return _repository.makeContributionFromWallet(
      circleId: params.circleId,
      contributionId: params.contributionId,
      pin: params.pin,
    );
  }
}
