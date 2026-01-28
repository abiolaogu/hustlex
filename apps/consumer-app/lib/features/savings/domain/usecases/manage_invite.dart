import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/phone_number.dart';
import '../entities/savings_circle.dart';
import '../entities/circle_invite.dart';
import '../repositories/savings_repository.dart';

part 'manage_invite.freezed.dart';

/// Parameters for sending an invite
@freezed
class SendInviteParams with _$SendInviteParams {
  const factory SendInviteParams({
    required String circleId,
    required PhoneNumber phoneNumber,
  }) = _SendInviteParams;
}

/// Use case to send invite to join circle
@injectable
class SendInvite implements UseCase<CircleInvite, SendInviteParams> {
  final SavingsRepository _repository;

  SendInvite(this._repository);

  @override
  Future<Either<Failure, CircleInvite>> call(SendInviteParams params) async {
    // Validate phone number
    if (!params.phoneNumber.isValid()) {
      return left(
        const Failure.invalidInput('phoneNumber', 'Invalid phone number'),
      );
    }

    return _repository.sendInvite(
      circleId: params.circleId,
      phoneNumber: params.phoneNumber.getOrCrash(),
    );
  }
}

/// Use case to get pending invites for current user
@injectable
class GetMyInvites implements UseCaseNoParams<List<CircleInvite>> {
  final SavingsRepository _repository;

  GetMyInvites(this._repository);

  @override
  Future<Either<Failure, List<CircleInvite>>> call() {
    return _repository.getMyInvites();
  }
}

/// Use case to accept an invite
@injectable
class AcceptInvite implements UseCase<CircleMember, String> {
  final SavingsRepository _repository;

  AcceptInvite(this._repository);

  @override
  Future<Either<Failure, CircleMember>> call(String inviteId) {
    return _repository.acceptInvite(inviteId);
  }
}

/// Use case to decline an invite
@injectable
class DeclineInvite implements UseCase<Unit, String> {
  final SavingsRepository _repository;

  DeclineInvite(this._repository);

  @override
  Future<Either<Failure, Unit>> call(String inviteId) {
    return _repository.declineInvite(inviteId);
  }
}

/// Use case to generate invite code for circle (admin only)
@injectable
class GenerateInviteCode implements UseCase<String, String> {
  final SavingsRepository _repository;

  GenerateInviteCode(this._repository);

  @override
  Future<Either<Failure, String>> call(String circleId) {
    return _repository.generateInviteCode(circleId);
  }
}
