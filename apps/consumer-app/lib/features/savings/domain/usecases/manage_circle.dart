import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/money.dart';
import '../entities/savings_circle.dart';
import '../repositories/savings_repository.dart';

part 'manage_circle.freezed.dart';

/// Parameters for creating a circle
@freezed
class CreateCircleParams with _$CreateCircleParams {
  const factory CreateCircleParams({
    required String name,
    String? description,
    required CircleType type,
    required Money contributionAmount,
    required ContributionFrequency frequency,
    required int maxMembers,
    @Default(false) bool isPrivate,
    DateTime? startDate,
    Money? targetAmount,
    int? durationMonths,
  }) = _CreateCircleParams;
}

/// Use case to create a new savings circle
@injectable
class CreateCircle implements UseCase<SavingsCircle, CreateCircleParams> {
  final SavingsRepository _repository;

  CreateCircle(this._repository);

  @override
  Future<Either<Failure, SavingsCircle>> call(CreateCircleParams params) async {
    // Validate name
    if (params.name.trim().isEmpty) {
      return left(
        const Failure.invalidInput('name', 'Circle name cannot be empty'),
      );
    }
    if (params.name.trim().length < 3) {
      return left(
        const Failure.invalidInput(
          'name',
          'Circle name must be at least 3 characters',
        ),
      );
    }

    // Validate contribution amount
    if (params.contributionAmount.amount <= 0) {
      return left(
        const Failure.invalidInput(
          'contributionAmount',
          'Contribution amount must be greater than 0',
        ),
      );
    }

    // Minimum contribution of 500 NGN
    if (params.contributionAmount.amount < 500) {
      return left(
        const Failure.invalidInput(
          'contributionAmount',
          'Minimum contribution is ₦500',
        ),
      );
    }

    // Validate max members
    if (params.maxMembers < 2) {
      return left(
        const Failure.invalidInput(
          'maxMembers',
          'Circle must have at least 2 members',
        ),
      );
    }
    if (params.maxMembers > 50) {
      return left(
        const Failure.invalidInput(
          'maxMembers',
          'Circle cannot have more than 50 members',
        ),
      );
    }

    // Validate start date if provided
    if (params.startDate != null) {
      if (params.startDate!.isBefore(DateTime.now())) {
        return left(
          const Failure.invalidInput(
            'startDate',
            'Start date must be in the future',
          ),
        );
      }
    }

    // Validate target amount for esusu/goal circles
    if (params.type == CircleType.esusu || params.type == CircleType.goal) {
      if (params.targetAmount == null || params.targetAmount!.amount <= 0) {
        return left(
          const Failure.invalidInput(
            'targetAmount',
            'Target amount is required for this circle type',
          ),
        );
      }
    }

    return _repository.createCircle(
      name: params.name.trim(),
      description: params.description?.trim(),
      type: params.type,
      contributionAmount: params.contributionAmount,
      frequency: params.frequency,
      maxMembers: params.maxMembers,
      isPrivate: params.isPrivate,
      startDate: params.startDate,
      targetAmount: params.targetAmount,
      durationMonths: params.durationMonths,
    );
  }
}

/// Parameters for updating a circle
@freezed
class UpdateCircleParams with _$UpdateCircleParams {
  const factory UpdateCircleParams({
    required String circleId,
    String? name,
    String? description,
    bool? isPrivate,
    DateTime? startDate,
  }) = _UpdateCircleParams;
}

/// Use case to update circle details (admin only)
@injectable
class UpdateCircle implements UseCase<SavingsCircle, UpdateCircleParams> {
  final SavingsRepository _repository;

  UpdateCircle(this._repository);

  @override
  Future<Either<Failure, SavingsCircle>> call(UpdateCircleParams params) async {
    // Validate name if provided
    if (params.name != null && params.name!.trim().isEmpty) {
      return left(
        const Failure.invalidInput('name', 'Circle name cannot be empty'),
      );
    }
    if (params.name != null && params.name!.trim().length < 3) {
      return left(
        const Failure.invalidInput(
          'name',
          'Circle name must be at least 3 characters',
        ),
      );
    }

    return _repository.updateCircle(
      circleId: params.circleId,
      name: params.name?.trim(),
      description: params.description?.trim(),
      isPrivate: params.isPrivate,
      startDate: params.startDate,
    );
  }
}

/// Use case to start a pending circle (admin only)
@injectable
class StartCircle implements UseCase<SavingsCircle, String> {
  final SavingsRepository _repository;

  StartCircle(this._repository);

  @override
  Future<Either<Failure, SavingsCircle>> call(String circleId) {
    return _repository.startCircle(circleId);
  }
}

/// Parameters for cancelling a circle
@freezed
class CancelCircleParams with _$CancelCircleParams {
  const factory CancelCircleParams({
    required String circleId,
    required String reason,
  }) = _CancelCircleParams;
}

/// Use case to cancel a circle (admin only)
@injectable
class CancelCircle implements UseCase<Unit, CancelCircleParams> {
  final SavingsRepository _repository;

  CancelCircle(this._repository);

  @override
  Future<Either<Failure, Unit>> call(CancelCircleParams params) async {
    // Validate reason
    if (params.reason.trim().isEmpty) {
      return left(
        const Failure.invalidInput('reason', 'Cancellation reason is required'),
      );
    }
    if (params.reason.trim().length < 10) {
      return left(
        const Failure.invalidInput(
          'reason',
          'Please provide a more detailed reason',
        ),
      );
    }

    return _repository.cancelCircle(
      circleId: params.circleId,
      reason: params.reason.trim(),
    );
  }
}
