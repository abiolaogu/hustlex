import 'package:dartz/dartz.dart';
import 'package:injectable/injectable.dart';
import '../../../../core/domain/failures/failure.dart';
import '../../../../core/domain/usecases/usecase.dart';
import '../repositories/auth_repository.dart';

/// Use case to check if biometric is available
@injectable
class IsBiometricAvailable {
  final AuthRepository _repository;

  IsBiometricAvailable(this._repository);

  Future<bool> call() {
    return _repository.isBiometricAvailable();
  }
}

/// Use case to enable biometric authentication
@injectable
class EnableBiometric implements UseCaseNoParams<Unit> {
  final AuthRepository _repository;

  EnableBiometric(this._repository);

  @override
  Future<Either<Failure, Unit>> call() {
    return _repository.enableBiometric();
  }
}

/// Use case to disable biometric authentication
@injectable
class DisableBiometric implements UseCaseNoParams<Unit> {
  final AuthRepository _repository;

  DisableBiometric(this._repository);

  @override
  Future<Either<Failure, Unit>> call() {
    return _repository.disableBiometric();
  }
}

/// Use case to authenticate with biometric
@injectable
class AuthenticateWithBiometric implements UseCaseNoParams<bool> {
  final AuthRepository _repository;

  AuthenticateWithBiometric(this._repository);

  @override
  Future<Either<Failure, bool>> call() {
    return _repository.authenticateWithBiometric();
  }
}
