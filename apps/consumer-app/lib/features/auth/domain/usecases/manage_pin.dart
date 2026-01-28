import 'package:dartz/dartz.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/failures/value_failures.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../../../../core/domain/value_objects/pin.dart';
import '../entities/user.dart';
import '../repositories/auth_repository.dart';

part 'manage_pin.freezed.dart';

/// Use case to set user's transaction PIN
@injectable
class SetPin implements UseCase<User, Pin> {
  final AuthRepository _repository;

  SetPin(this._repository);

  @override
  Future<Either<Failure, User>> call(Pin pin) async {
    // Validate PIN
    if (!pin.isValid()) {
      final failure = pin.failureOrNull;
      return left(
        Failure.invalidInput('pin', failure?.message ?? 'Invalid PIN'),
      );
    }

    return _repository.setPin(pin);
  }
}

/// Parameters for changing PIN
@freezed
class ChangePinParams with _$ChangePinParams {
  const factory ChangePinParams({
    required Pin currentPin,
    required Pin newPin,
    required Pin confirmPin,
  }) = _ChangePinParams;
}

/// Use case to change user's transaction PIN
@injectable
class ChangePin implements UseCase<User, ChangePinParams> {
  final AuthRepository _repository;

  ChangePin(this._repository);

  @override
  Future<Either<Failure, User>> call(ChangePinParams params) async {
    // Validate current PIN
    if (!params.currentPin.isValid()) {
      return left(
        const Failure.invalidInput('currentPin', 'Invalid current PIN'),
      );
    }

    // Validate new PIN
    if (!params.newPin.isValid()) {
      final failure = params.newPin.failureOrNull;
      return left(
        Failure.invalidInput('newPin', failure?.message ?? 'Invalid new PIN'),
      );
    }

    // Validate confirmation matches
    if (!params.confirmPin.isValid()) {
      return left(
        const Failure.invalidInput('confirmPin', 'Invalid confirmation PIN'),
      );
    }

    if (params.newPin.getOrCrash() != params.confirmPin.getOrCrash()) {
      return left(
        const Failure.invalidInput('confirmPin', 'PINs do not match'),
      );
    }

    // Check new PIN is different from current
    if (params.currentPin.getOrCrash() == params.newPin.getOrCrash()) {
      return left(
        const Failure.invalidInput(
          'newPin',
          'New PIN must be different from current PIN',
        ),
      );
    }

    return _repository.changePin(
      currentPin: params.currentPin,
      newPin: params.newPin,
    );
  }
}

/// Use case to verify PIN
@injectable
class VerifyPin implements UseCase<bool, Pin> {
  final AuthRepository _repository;

  VerifyPin(this._repository);

  @override
  Future<Either<Failure, bool>> call(Pin pin) async {
    if (!pin.isValid()) {
      return left(const Failure.invalidPin());
    }

    return _repository.verifyPin(pin);
  }
}
