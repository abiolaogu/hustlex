import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../entities/savings_circle.dart';
import '../repositories/savings_repository.dart';

part 'manage_membership.freezed.dart';

/// Use case to join a public circle
@injectable
class JoinCircle implements UseCase<CircleMember, String> {
  final SavingsRepository _repository;

  JoinCircle(this._repository);

  @override
  Future<Either<Failure, CircleMember>> call(String circleId) {
    return _repository.joinCircle(circleId);
  }
}

/// Parameters for joining with invite code
@freezed
class JoinCircleWithCodeParams with _$JoinCircleWithCodeParams {
  const factory JoinCircleWithCodeParams({
    required String circleId,
    required String inviteCode,
  }) = _JoinCircleWithCodeParams;
}

/// Use case to join a private circle with invite code
@injectable
class JoinCircleWithCode
    implements UseCase<CircleMember, JoinCircleWithCodeParams> {
  final SavingsRepository _repository;

  JoinCircleWithCode(this._repository);

  @override
  Future<Either<Failure, CircleMember>> call(
      JoinCircleWithCodeParams params) async {
    // Validate invite code
    if (params.inviteCode.trim().isEmpty) {
      return left(
        const Failure.invalidInput('inviteCode', 'Invite code is required'),
      );
    }

    return _repository.joinCircleWithCode(
      circleId: params.circleId,
      inviteCode: params.inviteCode.trim(),
    );
  }
}

/// Use case to leave a circle
@injectable
class LeaveCircle implements UseCase<Unit, String> {
  final SavingsRepository _repository;

  LeaveCircle(this._repository);

  @override
  Future<Either<Failure, Unit>> call(String circleId) {
    return _repository.leaveCircle(circleId);
  }
}

/// Parameters for removing a member
@freezed
class RemoveMemberParams with _$RemoveMemberParams {
  const factory RemoveMemberParams({
    required String circleId,
    required String memberId,
    required String reason,
  }) = _RemoveMemberParams;
}

/// Use case to remove a member from circle (admin only)
@injectable
class RemoveMember implements UseCase<Unit, RemoveMemberParams> {
  final SavingsRepository _repository;

  RemoveMember(this._repository);

  @override
  Future<Either<Failure, Unit>> call(RemoveMemberParams params) async {
    // Validate reason
    if (params.reason.trim().isEmpty) {
      return left(
        const Failure.invalidInput('reason', 'Removal reason is required'),
      );
    }

    return _repository.removeMember(
      circleId: params.circleId,
      memberId: params.memberId,
      reason: params.reason.trim(),
    );
  }
}

/// Use case to get members of a circle
@injectable
class GetCircleMembers implements UseCase<List<CircleMember>, String> {
  final SavingsRepository _repository;

  GetCircleMembers(this._repository);

  @override
  Future<Either<Failure, List<CircleMember>>> call(String circleId) {
    return _repository.getCircleMembers(circleId);
  }
}

/// Parameters for updating payout order
@freezed
class UpdatePayoutOrderParams with _$UpdatePayoutOrderParams {
  const factory UpdatePayoutOrderParams({
    required String circleId,
    required List<String> memberIds,
  }) = _UpdatePayoutOrderParams;
}

/// Use case to update payout order (admin only, for ajo/rotating circles)
@injectable
class UpdatePayoutOrder implements UseCase<Unit, UpdatePayoutOrderParams> {
  final SavingsRepository _repository;

  UpdatePayoutOrder(this._repository);

  @override
  Future<Either<Failure, Unit>> call(UpdatePayoutOrderParams params) async {
    // Validate member IDs
    if (params.memberIds.isEmpty) {
      return left(
        const Failure.invalidInput('memberIds', 'Member order list cannot be empty'),
      );
    }

    return _repository.updatePayoutOrder(
      circleId: params.circleId,
      memberIds: params.memberIds,
    );
  }
}
